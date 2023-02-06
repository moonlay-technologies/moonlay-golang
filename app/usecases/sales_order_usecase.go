package usecases

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/app/repositories"
	mongoRepositories "order-service/app/repositories/mongod"
	openSearchRepositories "order-service/app/repositories/open_search"
	"order-service/global/utils/helper"
	kafkadbo "order-service/global/utils/kafka"
	"order-service/global/utils/model"
	"strings"
	"time"

	"github.com/bxcodec/dbresolver"
)

type SalesOrderUseCaseInterface interface {
	Create(request *models.SalesOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderResponse, *model.ErrorLog)
	Get(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetByID(request *models.SalesOrderRequest, withDetail bool, ctx context.Context) (*models.SalesOrder, *model.ErrorLog)
	GetByAgentID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetByStoreID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetBySalesmanID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetByOrderStatusID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetByOrderSourceID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	SyncToOpenSearchFromCreateEvent(salesOrder *models.SalesOrder, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromUpdateEvent(salesOrder *models.SalesOrder, ctx context.Context) *model.ErrorLog
	UpdateBydId(id int, request *models.SalesOrderUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrder, *model.ErrorLog)
}

type salesOrderUseCase struct {
	salesOrderRepository              repositories.SalesOrderRepositoryInterface
	salesOrderDetailRepository        repositories.SalesOrderDetailRepositoryInterface
	orderStatusRepository             repositories.OrderStatusRepositoryInterface
	orderSourceRepository             repositories.OrderSourceRepositoryInterface
	agentRepository                   repositories.AgentRepositoryInterface
	brandRepository                   repositories.BrandRepositoryInterface
	storeRepository                   repositories.StoreRepositoryInterface
	productRepository                 repositories.ProductRepositoryInterface
	uomRepository                     repositories.UomRepositoryInterface
	deliveryOrderRepository           repositories.DeliveryOrderRepositoryInterface
	salesOrderLogRepository           mongoRepositories.SalesOrderLogRepositoryInterface
	userRepository                    repositories.UserRepositoryInterface
	salesmanRepository                repositories.SalesmanRepositoryInterface
	categoryRepository                repositories.CategoryRepositoryInterface
	salesOrderOpenSearchRepository    openSearchRepositories.SalesOrderOpenSearchRepositoryInterface
	deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface
	kafkaClient                       kafkadbo.KafkaClientInterface
	db                                dbresolver.DB
	ctx                               context.Context
}

func InitSalesOrderUseCaseInterface(salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, orderSourceRepository repositories.OrderSourceRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, uomRepository repositories.UomRepositoryInterface, deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, salesOrderLogRepository mongoRepositories.SalesOrderLogRepositoryInterface, userRepository repositories.UserRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, categoryRepository repositories.CategoryRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface, kafkaClient kafkadbo.KafkaClientInterface, db dbresolver.DB, ctx context.Context) SalesOrderUseCaseInterface {
	return &salesOrderUseCase{
		salesOrderRepository:              salesOrderRepository,
		salesOrderDetailRepository:        salesOrderDetailRepository,
		orderStatusRepository:             orderStatusRepository,
		orderSourceRepository:             orderSourceRepository,
		agentRepository:                   agentRepository,
		brandRepository:                   brandRepository,
		storeRepository:                   storeRepository,
		productRepository:                 productRepository,
		uomRepository:                     uomRepository,
		deliveryOrderRepository:           deliveryOrderRepository,
		salesOrderLogRepository:           salesOrderLogRepository,
		userRepository:                    userRepository,
		salesmanRepository:                salesmanRepository,
		categoryRepository:                categoryRepository,
		salesOrderOpenSearchRepository:    salesOrderOpenSearchRepository,
		deliveryOrderOpenSearchRepository: deliveryOrderOpenSearchRepository,
		kafkaClient:                       kafkaClient,
		db:                                db,
		ctx:                               ctx,
	}
}

func (u *salesOrderUseCase) Create(request *models.SalesOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderResponse, *model.ErrorLog) {
	now := time.Now()
	var soCode string

	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType("open", "sales_order", false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return &models.SalesOrderResponse{}, getOrderStatusResult.ErrorLog
	}

	getOrderDetailStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType("open", "sales_order_detail", false, ctx, getOrderDetailStatusResultChan)
	getOrderDetailStatusResult := <-getOrderDetailStatusResultChan

	if getOrderDetailStatusResult.Error != nil {
		return &models.SalesOrderResponse{}, getOrderDetailStatusResult.ErrorLog
	}

	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(request.OrderSourceID, false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getOrderSourceResult.Error != nil {
		return &models.SalesOrderResponse{}, getOrderSourceResult.ErrorLog
	}

	salesOrdersResponse := &models.SalesOrderResponse{
		SalesOrderStoreRequest: *request,
	}

	getBrandResultChan := make(chan *models.BrandChan)
	go u.brandRepository.GetByID(request.BrandID, false, ctx, getBrandResultChan)
	getBrandResult := <-getBrandResultChan

	if getBrandResult.Error != nil {
		return &models.SalesOrderResponse{}, getBrandResult.ErrorLog
	}

	salesOrdersResponse.BrandName = getBrandResult.Brand.Name

	soCode = helper.GenerateSOCode(request.AgentID, getOrderSourceResult.OrderSource.Code)

	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(request.AgentID, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	getStoreResultChan := make(chan *models.StoreChan)
	go u.storeRepository.GetByID(request.StoreID, false, ctx, getStoreResultChan)
	getStoreResult := <-getStoreResultChan

	if getStoreResult.Error != nil {
		errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	salesOrdersResponse.StoreCode = getStoreResult.Store.StoreCode
	salesOrdersResponse.StoreName = getStoreResult.Store.Name
	salesOrdersResponse.StoreAddress = getStoreResult.Store.Address
	salesOrdersResponse.StoreCityName = getStoreResult.Store.CityName
	salesOrdersResponse.StoreProvinceName = getStoreResult.Store.ProvinceName

	getUserResultChan := make(chan *models.UserChan)
	go u.userRepository.GetByID(request.UserID, false, ctx, getUserResultChan)
	getUserResult := <-getUserResultChan

	if getUserResult.Error != nil {
		errorLogData := helper.WriteLog(getUserResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	getSalesmanResultChan := make(chan *models.SalesmanChan)
	go u.salesmanRepository.GetByEmail(getUserResult.User.Email, false, ctx, getSalesmanResultChan)
	getSalesmanResult := <-getSalesmanResultChan

	if getSalesmanResult.Error != nil {
		errorLogData := helper.WriteLog(getSalesmanResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	salesOrdersResponse.SalesmanName = models.NullString{NullString: sql.NullString{String: getSalesmanResult.Salesman.Name, Valid: true}}

	salesOrder := &models.SalesOrder{
		CartID:            request.CartID,
		AgentID:           request.AgentID,
		StoreID:           request.StoreID,
		BrandID:           request.BrandID,
		UserID:            request.UserID,
		VisitationID:      request.VisitationID,
		OrderStatusID:     getOrderStatusResult.OrderStatus.ID,
		OrderSourceID:     request.OrderSourceID,
		GLat:              models.NullFloat64{NullFloat64: sql.NullFloat64{Float64: request.GLat, Valid: true}},
		GLong:             models.NullFloat64{NullFloat64: sql.NullFloat64{Float64: request.GLong, Valid: true}},
		SoCode:            soCode,
		SoRefCode:         models.NullString{NullString: sql.NullString{String: request.SoRefCode, Valid: true}},
		SoDate:            now.Format("2006-01-02"),
		SoRefDate:         models.NullString{NullString: sql.NullString{String: request.SoRefDate, Valid: true}},
		Note:              models.NullString{NullString: sql.NullString{String: request.Note, Valid: true}},
		InternalComment:   models.NullString{NullString: sql.NullString{String: request.InternalComment, Valid: true}},
		TotalAmount:       request.TotalAmount,
		TotalTonase:       request.TotalTonase,
		DeviceId:          models.NullString{NullString: sql.NullString{String: request.DeviceId, Valid: true}},
		ReferralCode:      models.NullString{NullString: sql.NullString{String: request.ReferralCode, Valid: true}},
		IsDoneSyncToEs:    "0",
		CreatedAt:         &now,
		StartDateSyncToEs: &now,
		StartCreatedDate:  &now,
		CreatedBy:         request.UserID,
	}

	createSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.Insert(salesOrder, sqlTransaction, ctx, createSalesOrderResultChan)
	createSalesOrderResult := <-createSalesOrderResultChan

	if createSalesOrderResult.Error != nil {
		return &models.SalesOrderResponse{}, createSalesOrderResult.ErrorLog
	}

	var salesOrderDetailResponses []*models.SalesOrderDetailStoreResponse
	for _, v := range request.SalesOrderDetails {
		soDetailCode, _ := helper.GenerateSODetailCode(int(createSalesOrderResult.ID), request.AgentID, v.ProductID, v.UomID)
		salesOrderDetail := &models.SalesOrderDetail{
			SalesOrderID:      int(createSalesOrderResult.ID),
			ProductID:         v.ProductID,
			UomID:             v.UomID,
			OrderStatusID:     v.OrderStatusId,
			SoDetailCode:      soDetailCode,
			Qty:               v.Qty,
			SentQty:           v.SentQty,
			ResidualQty:       v.ResidualQty,
			Price:             v.Price,
			Note:              models.NullString{NullString: sql.NullString{String: request.Note, Valid: true}},
			IsDoneSyncToEs:    "0",
			StartDateSyncToEs: &now,
			CreatedAt:         &now,
		}

		createSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.Insert(salesOrderDetail, sqlTransaction, ctx, createSalesOrderDetailResultChan)
		createSalesOrderDetailResult := <-createSalesOrderDetailResultChan

		if createSalesOrderDetailResult.Error != nil {
			return &models.SalesOrderResponse{}, createSalesOrderDetailResult.ErrorLog
		}

		salesOrderDetailResponse := &models.SalesOrderDetailStoreResponse{
			SalesOrderDetailStoreRequest: *v,
		}

		getProductResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetByID(v.ProductID, false, ctx, getProductResultChan)
		getProductResult := <-getProductResultChan

		if getProductResult.Error != nil {
			errorLogData := helper.WriteLog(getProductResult.Error, http.StatusInternalServerError, nil)
			return &models.SalesOrderResponse{}, errorLogData
		}

		salesOrderDetailResponse.ProductSKU = getProductResult.Product.Sku.String
		salesOrderDetailResponse.ProductName = getProductResult.Product.ProductName.String

		getCategoryResultChan := make(chan *models.CategoryChan)
		go u.categoryRepository.GetByID(getProductResult.Product.CategoryID, false, ctx, getCategoryResultChan)
		getCategoryResult := <-getCategoryResultChan

		if getCategoryResult.Error != nil {
			errorLogData := helper.WriteLog(getCategoryResult.Error, http.StatusInternalServerError, nil)
			return &models.SalesOrderResponse{}, errorLogData
		}

		salesOrderDetailResponse.CategoryName = getCategoryResult.Category.Name

		getUomResultChan := make(chan *models.UomChan)
		go u.uomRepository.GetByID(v.UomID, false, ctx, getUomResultChan)
		getUomResult := <-getUomResultChan

		if getUomResult.Error != nil {
			errorLogData := helper.WriteLog(getUomResult.Error, http.StatusInternalServerError, nil)
			return &models.SalesOrderResponse{}, errorLogData
		}

		salesOrderDetailResponse.UomCode = getUomResult.Uom.Code.String

		salesOrderDetailResponses = append(salesOrderDetailResponses, salesOrderDetailResponse)
	}

	salesOrdersResponse.SalesOrderDetails = salesOrderDetailResponses

	salesOrderLog := &models.SalesOrderLog{
		RequestID: request.RequestID,
		SoCode:    soCode,
		Data:      salesOrdersResponse,
		Status:    "0",
		CreatedAt: &now,
	}

	createSalesOrderLogResultChan := make(chan *models.SalesOrderLogChan)
	go u.salesOrderLogRepository.Insert(salesOrderLog, ctx, createSalesOrderLogResultChan)
	createSalesOrderLogResult := <-createSalesOrderLogResultChan

	if createSalesOrderLogResult.Error != nil {
		return &models.SalesOrderResponse{}, createSalesOrderLogResult.ErrorLog
	}

	keyKafka := []byte(salesOrder.SoCode)
	messageKafka, _ := json.Marshal(salesOrder)
	err := u.kafkaClient.WriteToTopic("create-sales-order", keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	return salesOrdersResponse, nil
}

func (u *salesOrderUseCase) Get(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog) {
	getSalesOrdersResultChan := make(chan *models.SalesOrdersChan)
	go u.salesOrderOpenSearchRepository.Get(request, getSalesOrdersResultChan)
	getSalesOrdersResult := <-getSalesOrdersResultChan

	if getSalesOrdersResult.Error != nil {
		return &models.SalesOrders{}, getSalesOrdersResult.ErrorLog
	}

	salesOrders := &models.SalesOrders{
		SalesOrders: getSalesOrdersResult.SalesOrders,
		Total:       getSalesOrdersResult.Total,
	}

	return salesOrders, &model.ErrorLog{}
}

func (u *salesOrderUseCase) GetByID(request *models.SalesOrderRequest, withDetail bool, ctx context.Context) (*models.SalesOrder, *model.ErrorLog) {

	if withDetail == true {
		now := time.Now()
		getSalesOrderResultChan := make(chan *models.SalesOrderChan)
		go u.salesOrderRepository.GetByID(request.ID, false, ctx, getSalesOrderResultChan)
		getSalesOrderResult := <-getSalesOrderResultChan

		if getSalesOrderResult.Error != nil {
			return &models.SalesOrder{}, getSalesOrderResult.ErrorLog
		}

		getSalesOrderDetailsResultChan := make(chan *models.SalesOrderDetailsChan)
		go u.salesOrderDetailRepository.GetBySalesOrderID(request.ID, false, ctx, getSalesOrderDetailsResultChan)
		getSalesOrderDetailsResult := <-getSalesOrderDetailsResultChan

		if getSalesOrderDetailsResult.Error != nil {
			return &models.SalesOrder{}, getSalesOrderDetailsResult.ErrorLog
		}

		getSalesOrderResult.SalesOrder.SalesOrderDetails = getSalesOrderDetailsResult.SalesOrderDetails

		for k, v := range getSalesOrderDetailsResult.SalesOrderDetails {
			getProductResultChan := make(chan *models.ProductChan)
			go u.productRepository.GetByID(v.ProductID, false, ctx, getProductResultChan)
			getProductResult := <-getProductResultChan

			if getProductResult.Error != nil {
				return &models.SalesOrder{}, getProductResult.ErrorLog
			}

			getSalesOrderDetailsResult.SalesOrderDetails[k].Product = getProductResult.Product

			getUomResultChan := make(chan *models.UomChan)
			go u.uomRepository.GetByID(v.UomID, false, ctx, getUomResultChan)
			getUomResult := <-getUomResultChan

			if getUomResult.Error != nil {
				return &models.SalesOrder{}, getUomResult.ErrorLog
			}

			getSalesOrderDetailsResult.SalesOrderDetails[k].Uom = getUomResult.Uom

			getOrderStatusDetailResultChan := make(chan *models.OrderStatusChan)
			go u.orderStatusRepository.GetByID(v.OrderStatusID, false, ctx, getOrderStatusDetailResultChan)
			getOrderStatusDetailResult := <-getOrderStatusDetailResultChan

			if getOrderStatusDetailResult.Error != nil {
				return &models.SalesOrder{}, getOrderStatusDetailResult.ErrorLog
			}

			getSalesOrderResult.SalesOrder.SalesOrderDetails[k].EndDateSyncToEs = &now
			getSalesOrderResult.SalesOrder.SalesOrderDetails[k].IsDoneSyncToEs = "1"
			getSalesOrderDetailsResult.SalesOrderDetails[k].OrderStatus = getOrderStatusDetailResult.OrderStatus
		}

		getSalesOrderResult.SalesOrder.SalesOrderDetails = getSalesOrderDetailsResult.SalesOrderDetails

		getOrderStatusResultChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByID(getSalesOrderResult.SalesOrder.OrderStatusID, false, ctx, getOrderStatusResultChan)
		getOrderStatusResult := <-getOrderStatusResultChan

		if getOrderStatusResult.Error != nil {
			return &models.SalesOrder{}, getOrderStatusResult.ErrorLog
		}

		getSalesOrderResult.SalesOrder.OrderStatus = getOrderStatusResult.OrderStatus

		getOrderSourceResultChan := make(chan *models.OrderSourceChan)
		go u.orderSourceRepository.GetByID(request.OrderSourceID, false, ctx, getOrderSourceResultChan)
		getOrderSourceResult := <-getOrderSourceResultChan

		if getOrderSourceResult.Error != nil {
			return &models.SalesOrder{}, getOrderSourceResult.ErrorLog
		}

		getSalesOrderResult.SalesOrder.OrderSource = getOrderSourceResult.OrderSource

		getAgentResultChan := make(chan *models.AgentChan)
		go u.agentRepository.GetByID(getSalesOrderResult.SalesOrder.AgentID, false, ctx, getAgentResultChan)
		getAgentResult := <-getAgentResultChan

		if getAgentResult.Error != nil {
			return &models.SalesOrder{}, getAgentResult.ErrorLog
		}

		getSalesOrderResult.SalesOrder.Agent = getAgentResult.Agent
		getSalesOrderResult.SalesOrder.AgentName = models.NullString{NullString: sql.NullString{String: getAgentResult.Agent.Name, Valid: true}}
		getSalesOrderResult.SalesOrder.AgentEmail = getAgentResult.Agent.Email
		getSalesOrderResult.SalesOrder.AgentProvinceName = getAgentResult.Agent.ProvinceName
		getSalesOrderResult.SalesOrder.AgentCityName = getAgentResult.Agent.CityName
		getSalesOrderResult.SalesOrder.AgentDistrictName = getAgentResult.Agent.DistrictName
		getSalesOrderResult.SalesOrder.AgentVillageName = getAgentResult.Agent.VillageName
		getSalesOrderResult.SalesOrder.AgentAddress = getAgentResult.Agent.Address
		getSalesOrderResult.SalesOrder.AgentPhone = getAgentResult.Agent.Phone
		getSalesOrderResult.SalesOrder.AgentMainMobilePhone = getAgentResult.Agent.MainMobilePhone

		getStoreResultChan := make(chan *models.StoreChan)
		go u.storeRepository.GetByID(getSalesOrderResult.SalesOrder.StoreID, false, ctx, getStoreResultChan)
		getStoreResult := <-getStoreResultChan

		if getStoreResult.Error != nil {
			return &models.SalesOrder{}, getStoreResult.ErrorLog
		}

		getSalesOrderResult.SalesOrder.Store = getStoreResult.Store
		getSalesOrderResult.SalesOrder.StoreName = getStoreResult.Store.Name
		getSalesOrderResult.SalesOrder.StoreCode = getStoreResult.Store.StoreCode
		getSalesOrderResult.SalesOrder.StoreEmail = getStoreResult.Store.Email
		getSalesOrderResult.SalesOrder.StoreProvinceName = getStoreResult.Store.ProvinceName
		getSalesOrderResult.SalesOrder.StoreCityName = getStoreResult.Store.CityName
		getSalesOrderResult.SalesOrder.StoreDistrictName = getStoreResult.Store.DistrictName
		getSalesOrderResult.SalesOrder.StoreVillageName = getStoreResult.Store.VillageName
		getSalesOrderResult.SalesOrder.StoreAddress = getStoreResult.Store.Address
		getSalesOrderResult.SalesOrder.StorePhone = getStoreResult.Store.Phone
		getSalesOrderResult.SalesOrder.StoreMainMobilePhone = getStoreResult.Store.MainMobilePhone

		getBrandResultChan := make(chan *models.BrandChan)
		go u.brandRepository.GetByID(getSalesOrderResult.SalesOrder.BrandID, false, ctx, getBrandResultChan)
		getBrandResult := <-getBrandResultChan

		if getBrandResult.Error != nil {
			return &models.SalesOrder{}, getBrandResult.ErrorLog
		}

		getSalesOrderResult.SalesOrder.Brand = getBrandResult.Brand
		getSalesOrderResult.SalesOrder.BrandName = getBrandResult.Brand.Name
		getSalesOrderResult.SalesOrder.OrderSource = getOrderSourceResult.OrderSource
		getSalesOrderResult.SalesOrder.OrderSourceName = getOrderSourceResult.OrderSource.SourceName
		getSalesOrderResult.SalesOrder.OrderStatus = getOrderStatusResult.OrderStatus
		getSalesOrderResult.SalesOrder.OrderStatusName = getOrderStatusResult.OrderStatus.Name

		getUserResultChan := make(chan *models.UserChan)
		go u.userRepository.GetByID(getSalesOrderResult.SalesOrder.UserID, false, ctx, getUserResultChan)
		getUserResult := <-getUserResultChan

		if getUserResult.Error != nil {
			return &models.SalesOrder{}, getUserResult.ErrorLog
		}

		getSalesOrderResult.SalesOrder.User = getUserResult.User
		getSalesOrderResult.SalesOrder.UserFirstName = getUserResult.User.FirstName
		getSalesOrderResult.SalesOrder.UserLastName = getUserResult.User.LastName
		getSalesOrderResult.SalesOrder.UserEmail = models.NullString{NullString: sql.NullString{String: getUserResult.User.Email, Valid: true}}

		if getUserResult.User.RoleID.String == "3" {
			getSalesmanResultChan := make(chan *models.SalesmanChan)
			go u.salesmanRepository.GetByEmail(getUserResult.User.Email, false, ctx, getSalesmanResultChan)
			getSalesmanResult := <-getSalesmanResultChan

			if getSalesmanResult.Error != nil {
				return &models.SalesOrder{}, getSalesmanResult.ErrorLog
			}

			getSalesOrderResult.SalesOrder.Salesman = getSalesmanResult.Salesman
			getSalesOrderResult.SalesOrder.SalesmanName = models.NullString{NullString: sql.NullString{String: getSalesmanResult.Salesman.Name, Valid: true}}
			getSalesOrderResult.SalesOrder.SalesmanEmail = getSalesmanResult.Salesman.Email
		}

		return getSalesOrderResult.SalesOrder, &model.ErrorLog{}
	} else {
		getSalesOrderResultChan := make(chan *models.SalesOrderChan)
		go u.salesOrderOpenSearchRepository.GetByID(request, getSalesOrderResultChan)
		getSalesOrderResult := <-getSalesOrderResultChan

		if getSalesOrderResult.Error != nil {
			return &models.SalesOrder{}, getSalesOrderResult.ErrorLog
		}

		return getSalesOrderResult.SalesOrder, &model.ErrorLog{}
	}
}

func (u *salesOrderUseCase) GetByAgentID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog) {
	getSalesOrdersResultChan := make(chan *models.SalesOrdersChan)
	go u.salesOrderOpenSearchRepository.GetByAgentID(request, getSalesOrdersResultChan)
	getSalesOrdersResult := <-getSalesOrdersResultChan

	if getSalesOrdersResult.Error != nil {
		return &models.SalesOrders{}, getSalesOrdersResult.ErrorLog
	}

	salesOrders := &models.SalesOrders{
		SalesOrders: getSalesOrdersResult.SalesOrders,
		Total:       getSalesOrdersResult.Total,
	}

	return salesOrders, &model.ErrorLog{}
}

func (u *salesOrderUseCase) GetByStoreID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog) {
	getSalesOrdersResultChan := make(chan *models.SalesOrdersChan)
	go u.salesOrderOpenSearchRepository.GetByStoreID(request, getSalesOrdersResultChan)
	getSalesOrdersResult := <-getSalesOrdersResultChan

	if getSalesOrdersResult.Error != nil {
		return &models.SalesOrders{}, getSalesOrdersResult.ErrorLog
	}

	salesOrders := &models.SalesOrders{
		SalesOrders: getSalesOrdersResult.SalesOrders,
		Total:       getSalesOrdersResult.Total,
	}

	return salesOrders, &model.ErrorLog{}
}

func (u *salesOrderUseCase) GetBySalesmanID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog) {
	getSalesOrdersResultChan := make(chan *models.SalesOrdersChan)
	go u.salesOrderOpenSearchRepository.GetBySalesmanID(request, getSalesOrdersResultChan)
	getSalesOrdersResult := <-getSalesOrdersResultChan

	if getSalesOrdersResult.Error != nil {
		return &models.SalesOrders{}, getSalesOrdersResult.ErrorLog
	}

	salesOrders := &models.SalesOrders{
		SalesOrders: getSalesOrdersResult.SalesOrders,
		Total:       getSalesOrdersResult.Total,
	}

	return salesOrders, &model.ErrorLog{}
}

func (u *salesOrderUseCase) GetByOrderStatusID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog) {
	getSalesOrdersResultChan := make(chan *models.SalesOrdersChan)
	go u.salesOrderOpenSearchRepository.GetByOrderStatusID(request, getSalesOrdersResultChan)
	getSalesOrdersResult := <-getSalesOrdersResultChan

	if getSalesOrdersResult.Error != nil {
		return &models.SalesOrders{}, getSalesOrdersResult.ErrorLog
	}

	salesOrders := &models.SalesOrders{
		SalesOrders: getSalesOrdersResult.SalesOrders,
		Total:       getSalesOrdersResult.Total,
	}

	return salesOrders, &model.ErrorLog{}
}

func (u *salesOrderUseCase) GetByOrderSourceID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog) {
	getSalesOrdersResultChan := make(chan *models.SalesOrdersChan)
	go u.salesOrderOpenSearchRepository.GetByOrderSourceID(request, getSalesOrdersResultChan)
	getSalesOrdersResult := <-getSalesOrdersResultChan

	if getSalesOrdersResult.Error != nil {
		return &models.SalesOrders{}, getSalesOrdersResult.ErrorLog
	}

	salesOrders := &models.SalesOrders{
		SalesOrders: getSalesOrdersResult.SalesOrders,
		Total:       getSalesOrdersResult.Total,
	}

	return salesOrders, &model.ErrorLog{}
}

func (u *salesOrderUseCase) SyncToOpenSearchFromCreateEvent(salesOrder *models.SalesOrder, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog {
	now := time.Now()

	deliveryOrdersRequest := &models.DeliveryOrderRequest{
		SalesOrderID: salesOrder.ID,
	}

	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
	go u.deliveryOrderOpenSearchRepository.GetBySalesOrderID(deliveryOrdersRequest, getDeliveryOrdersResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrdersResultChan
	deliveryOrdersFound := true

	if getDeliveryOrdersResult.Error != nil {
		deliveryOrdersFound = false
		if !strings.Contains(getDeliveryOrdersResult.Error.Error(), "not found") {
			errorLogData := helper.WriteLog(getDeliveryOrdersResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}
	}

	if deliveryOrdersFound == true {
		salesOrder.DeliveryOrders = getDeliveryOrdersResult.DeliveryOrders
	}

	for k, v := range salesOrder.SalesOrderDetails {
		getProductResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetByID(v.ProductID, false, ctx, getProductResultChan)
		getProductResult := <-getProductResultChan

		if getProductResult.Error != nil {
			errorLogData := helper.WriteLog(getProductResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		salesOrder.SalesOrderDetails[k].Product = getProductResult.Product

		getUomResultChan := make(chan *models.UomChan)
		go u.uomRepository.GetByID(v.UomID, false, ctx, getUomResultChan)
		getUomResult := <-getUomResultChan

		if getUomResult.Error != nil {
			errorLogData := helper.WriteLog(getUomResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		salesOrder.SalesOrderDetails[k].Uom = getUomResult.Uom

		getOrderStatusDetailResultChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByID(v.OrderStatusID, false, ctx, getOrderStatusDetailResultChan)
		getOrderStatusDetailResult := <-getOrderStatusDetailResultChan

		if getOrderStatusDetailResult.Error != nil {
			errorLogData := helper.WriteLog(getOrderStatusDetailResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		salesOrder.SalesOrderDetails[k].EndDateSyncToEs = &now
		salesOrder.SalesOrderDetails[k].IsDoneSyncToEs = "1"
		salesOrder.SalesOrderDetails[k].OrderStatus = getOrderStatusDetailResult.OrderStatus

		salesOrderDetailUpdateData := &models.SalesOrderDetail{
			UpdatedAt:       &now,
			IsDoneSyncToEs:  "1",
			EndDateSyncToEs: &now,
		}

		updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.UpdateByID(v.ID, salesOrderDetailUpdateData, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
		updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

		if updateSalesOrderDetailResult.Error != nil {
			errorLogData := helper.WriteLog(updateSalesOrderDetailResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}
	}

	salesOrder.IsDoneSyncToEs = "1"
	salesOrder.EndDateSyncToEs = &now
	salesOrder.UpdatedAt = &now
	salesOrder.EndCreatedDate = &now

	createSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderOpenSearchRepository.Create(salesOrder, createSalesOrderResultChan)
	createSalesOrderResult := <-createSalesOrderResultChan

	if createSalesOrderResult.Error != nil {
		return createSalesOrderResult.ErrorLog
	}

	salesOrderUpdateData := &models.SalesOrder{
		UpdatedAt:       &now,
		IsDoneSyncToEs:  "1",
		EndCreatedDate:  &now,
		EndDateSyncToEs: &now,
	}

	updateSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.UpdateByID(salesOrder.ID, salesOrderUpdateData, sqlTransaction, ctx, updateSalesOrderResultChan)
	updateSalesOrderResult := <-updateSalesOrderResultChan

	if updateSalesOrderResult.Error != nil {
		errorLogData := helper.WriteLog(updateSalesOrderResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return &model.ErrorLog{}
}

func (u *salesOrderUseCase) SyncToOpenSearchFromUpdateEvent(salesOrder *models.SalesOrder, ctx context.Context) *model.ErrorLog {
	now := time.Now()

	deliveryOrdersRequest := &models.DeliveryOrderRequest{
		SalesOrderID: salesOrder.ID,
	}

	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
	go u.deliveryOrderOpenSearchRepository.GetBySalesOrderID(deliveryOrdersRequest, getDeliveryOrdersResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrdersResultChan
	deliveryOrdersFound := true

	if getDeliveryOrdersResult.Error != nil {
		deliveryOrdersFound = false
		if !strings.Contains(getDeliveryOrdersResult.Error.Error(), "not found") {
			errorLogData := helper.WriteLog(getDeliveryOrdersResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}
	}

	if deliveryOrdersFound == true {
		fmt.Println("sj ktm")
		for x := range getDeliveryOrdersResult.DeliveryOrders {
			getDeliveryOrdersResult.DeliveryOrders[x].SalesOrder = nil
		}

		salesOrder.DeliveryOrders = getDeliveryOrdersResult.DeliveryOrders
	}

	for k, v := range salesOrder.SalesOrderDetails {
		getProductResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetByID(v.ProductID, false, ctx, getProductResultChan)
		getProductResult := <-getProductResultChan

		if getProductResult.Error != nil {
			errorLogData := helper.WriteLog(getProductResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		salesOrder.SalesOrderDetails[k].Product = getProductResult.Product

		getUomResultChan := make(chan *models.UomChan)
		go u.uomRepository.GetByID(v.UomID, false, ctx, getUomResultChan)
		getUomResult := <-getUomResultChan

		if getUomResult.Error != nil {
			errorLogData := helper.WriteLog(getUomResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		salesOrder.SalesOrderDetails[k].Uom = getUomResult.Uom

		getOrderStatusDetailResultChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByID(v.OrderStatusID, false, ctx, getOrderStatusDetailResultChan)
		getOrderStatusDetailResult := <-getOrderStatusDetailResultChan

		if getOrderStatusDetailResult.Error != nil {
			errorLogData := helper.WriteLog(getOrderStatusDetailResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		salesOrder.SalesOrderDetails[k].OrderStatus = getOrderStatusDetailResult.OrderStatus
	}

	salesOrder.UpdatedAt = &now

	updateSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderOpenSearchRepository.Create(salesOrder, updateSalesOrderResultChan)
	updateSalesOrderResult := <-updateSalesOrderResultChan

	if updateSalesOrderResult.Error != nil {
		return updateSalesOrderResult.ErrorLog
	}

	return &model.ErrorLog{}
}

func (u *salesOrderUseCase) UpdateBydId(id int, request *models.SalesOrderUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrder, *model.ErrorLog) {
	now := time.Now()
	var soCode string

	// Check Order Status
	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType("open", "sales_order", false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return &models.SalesOrder{}, getOrderStatusResult.ErrorLog
	}

	// Check Order Detail Status
	getOrderDetailStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType("open", "sales_order_detail", false, ctx, getOrderDetailStatusResultChan)
	getOrderDetailStatusResult := <-getOrderDetailStatusResultChan

	if getOrderDetailStatusResult.Error != nil {
		return &models.SalesOrder{}, getOrderDetailStatusResult.ErrorLog
	}

	// Check Order Source
	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(request.OrderSourceID, false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getOrderSourceResult.Error != nil {
		return &models.SalesOrder{}, getOrderSourceResult.ErrorLog
	}

	salesOrder := &models.SalesOrder{
		OrderStatusID:   getOrderStatusResult.OrderStatus.ID,
		OrderStatusName: getOrderStatusResult.OrderStatus.Name,
		OrderSourceName: getOrderSourceResult.OrderSource.SourceName,
		OrderSourceID:   request.OrderSourceID,
		AgentID:         request.AgentID,
		StoreID:         request.StoreID,
		BrandID:         request.SalesOrderDetails[0].BrandID,
		UserID:          request.UserID,
		GLat:            models.NullFloat64{sql.NullFloat64{Float64: request.GLat, Valid: true}},
		GLong:           models.NullFloat64{sql.NullFloat64{Float64: request.GLong, Valid: true}},
		SoRefCode:       models.NullString{sql.NullString{String: request.SoRefCode, Valid: true}},
		SoDate:          request.SoDate,
		SoRefDate:       models.NullString{sql.NullString{String: request.SoRefDate, Valid: true}},
		Note:            models.NullString{sql.NullString{String: request.Note, Valid: true}},
		InternalComment: models.NullString{sql.NullString{String: request.InternalComment, Valid: true}},
		TotalAmount:     request.TotalAmount,
		TotalTonase:     request.TotalTonase,
		DeviceId:        models.NullString{sql.NullString{String: request.DeviceId, Valid: true}},
		ReferralCode:    models.NullString{sql.NullString{String: request.ReferralCode, Valid: true}},
		UpdatedAt:       &now,
		LatestUpdatedBy: request.UserID,
	}

	// Check Brand
	brandIds := []int{}

	for _, v := range request.SalesOrderDetails {
		brandIds = append(brandIds, v.BrandID)
	}
	checkIfBrandSame := helper.InSliceInt(brandIds, request.SalesOrderDetails[0].BrandID)

	if !checkIfBrandSame {
		errorLogData := helper.WriteLog(fmt.Errorf("The brand id must be the same"), http.StatusInternalServerError, "Ada kesalahan, silahkan coba lagi nanti")
		return &models.SalesOrder{}, errorLogData
	}

	getBrandResultChan := make(chan *models.BrandChan)
	go u.brandRepository.GetByID(request.SalesOrderDetails[0].BrandID, false, ctx, getBrandResultChan)
	getBrandResult := <-getBrandResultChan

	if getBrandResult.Error != nil {
		return &models.SalesOrder{}, getBrandResult.ErrorLog
	}

	salesOrder.BrandName = getBrandResult.Brand.Name

	// Check Agent
	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(request.AgentID, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrder{}, errorLogData
	}

	salesOrder.AgentID = request.AgentID
	salesOrder.AgentName = models.NullString{sql.NullString{String: getAgentResult.Agent.Name, Valid: true}}
	salesOrder.AgentEmail = getAgentResult.Agent.Email
	salesOrder.AgentProvinceName = getAgentResult.Agent.ProvinceName
	salesOrder.AgentCityName = getAgentResult.Agent.CityName
	salesOrder.AgentDistrictName = getAgentResult.Agent.DistrictName
	salesOrder.AgentVillageName = getAgentResult.Agent.VillageName
	salesOrder.AgentAddress = getAgentResult.Agent.Address
	salesOrder.AgentPhone = getAgentResult.Agent.Phone
	salesOrder.AgentMainMobilePhone = getAgentResult.Agent.MainMobilePhone

	// Check Store
	getStoreResultChan := make(chan *models.StoreChan)
	go u.storeRepository.GetByID(request.StoreID, false, ctx, getStoreResultChan)
	getStoreResult := <-getStoreResultChan

	if getStoreResult.Error != nil {
		errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrder{}, errorLogData
	}

	salesOrder.StoreID = request.StoreID
	salesOrder.StoreName = getStoreResult.Store.Name
	salesOrder.StoreCode = getStoreResult.Store.StoreCode
	salesOrder.StoreEmail = getStoreResult.Store.Email
	salesOrder.StoreProvinceName = getStoreResult.Store.ProvinceName
	salesOrder.StoreCityName = getStoreResult.Store.CityName
	salesOrder.StoreDistrictName = getStoreResult.Store.DistrictName
	salesOrder.StoreVillageName = getStoreResult.Store.VillageName
	salesOrder.StoreAddress = getStoreResult.Store.Address
	salesOrder.StorePhone = getStoreResult.Store.Phone
	salesOrder.StoreMainMobilePhone = getStoreResult.Store.MainMobilePhone

	// Check User Result
	getUserResultChan := make(chan *models.UserChan)
	go u.userRepository.GetByID(request.UserID, false, ctx, getUserResultChan)
	getUserResult := <-getUserResultChan

	if getUserResult.Error != nil {
		errorLogData := helper.WriteLog(getUserResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrder{}, errorLogData
	}

	salesOrder.UserID = request.UserID
	salesOrder.UserFirstName = getUserResult.User.FirstName
	salesOrder.UserLastName = getUserResult.User.LastName
	salesOrder.UserEmail = models.NullString{sql.NullString{String: getUserResult.User.Email, Valid: true}}

	// Check Salesman
	getSalesmanResultChan := make(chan *models.SalesmanChan)
	go u.salesmanRepository.GetByEmail(getUserResult.User.Email, false, ctx, getSalesmanResultChan)
	getSalesmanResult := <-getSalesmanResultChan

	if getSalesmanResult.Error != nil {
		errorLogData := helper.WriteLog(getSalesmanResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrder{}, errorLogData
	}

	salesOrder.SalesmanName = models.NullString{sql.NullString{String: getSalesmanResult.Salesman.Name, Valid: true}}
	salesOrder.SalesmanEmail = getSalesmanResult.Salesman.Email

	salesOrderUpdateReq := &models.SalesOrder{
		OrderSourceID:   request.OrderSourceID,
		AgentID:         request.AgentID,
		StoreID:         request.StoreID,
		BrandID:         request.SalesOrderDetails[0].BrandID,
		UserID:          request.UserID,
		GLat:            models.NullFloat64{sql.NullFloat64{Float64: request.GLat, Valid: true}},
		GLong:           models.NullFloat64{sql.NullFloat64{Float64: request.GLong, Valid: true}},
		SoRefCode:       models.NullString{sql.NullString{String: request.SoRefCode, Valid: true}},
		SoDate:          request.SoDate,
		SoRefDate:       models.NullString{sql.NullString{String: request.SoRefDate, Valid: true}},
		Note:            models.NullString{sql.NullString{String: request.Note, Valid: true}},
		InternalComment: models.NullString{sql.NullString{String: request.InternalComment, Valid: true}},
		TotalAmount:     request.TotalAmount,
		TotalTonase:     request.TotalTonase,
		DeviceId:        models.NullString{sql.NullString{String: request.DeviceId, Valid: true}},
		ReferralCode:    models.NullString{sql.NullString{String: request.ReferralCode, Valid: true}},
		UpdatedAt:       &now,
		LatestUpdatedBy: request.UserID,
	}
	updateSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.UpdateByID(id, salesOrderUpdateReq, sqlTransaction, ctx, updateSalesOrderResultChan)
	updateSalesOrderResult := <-updateSalesOrderResultChan

	if updateSalesOrderResult.Error != nil {
		return &models.SalesOrder{}, updateSalesOrderResult.ErrorLog
	}

	for _, v := range request.SalesOrderDetails {
		salesOrderDetail := &models.SalesOrderDetail{
			ID:          v.ID,
			ProductID:   v.ProductID,
			UomID:       v.UomID,
			Qty:         v.Qty,
			SentQty:     v.SentQty,
			ResidualQty: v.ResidualQty,
			Price:       v.Price,
			Note:        models.NullString{sql.NullString{String: request.Note, Valid: true}},
			UpdatedAt:   &now,
		}

		updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.UpdateByID(v.ID, salesOrderDetail, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
		updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

		if updateSalesOrderDetailResult.Error != nil {
			return &models.SalesOrder{}, updateSalesOrderDetailResult.ErrorLog
		}

	}

	getSalesOrderDetailBySOIdResultChan := make(chan *models.SalesOrderDetailsChan)
	go u.salesOrderDetailRepository.GetBySalesOrderID(id, false, ctx, getSalesOrderDetailBySOIdResultChan)
	getSalesOrderDetailBySOIdResult := <-getSalesOrderDetailBySOIdResultChan

	if getSalesOrderDetailBySOIdResult.Error != nil {
		return &models.SalesOrder{}, getSalesOrderDetailBySOIdResult.ErrorLog
	}

	salesOrder.SalesOrderDetails = getSalesOrderDetailBySOIdResult.SalesOrderDetails

	salesOrderLog := &models.SalesOrderLog{
		RequestID: request.RequestID,
		SoCode:    soCode,
		Data:      salesOrder,
		Status:    "0",
		CreatedAt: &now,
	}

	createSalesOrderLogResultChan := make(chan *models.SalesOrderLogChan)
	go u.salesOrderLogRepository.Insert(salesOrderLog, ctx, createSalesOrderLogResultChan)
	createSalesOrderLogResult := <-createSalesOrderLogResultChan

	if createSalesOrderLogResult.Error != nil {
		return &models.SalesOrder{}, createSalesOrderLogResult.ErrorLog
	}

	keyKafka := []byte(salesOrder.SoCode)
	messageKafka, _ := json.Marshal(salesOrder)
	err := u.kafkaClient.WriteToTopic("create-sales-order", keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.SalesOrder{}, errorLogData
	}

	return salesOrder, nil
}
