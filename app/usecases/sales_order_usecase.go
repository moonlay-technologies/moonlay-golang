package usecases

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
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
	Get(request *models.SalesOrderRequest) (*models.SalesOrdersOpenSearchResponse, *model.ErrorLog)
	GetByID(request *models.SalesOrderRequest, withDetail bool, ctx context.Context) (*models.SalesOrder, *model.ErrorLog)
	GetByAgentID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetByStoreID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetBySalesmanID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetByOrderStatusID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetByOrderSourceID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	SyncToOpenSearchFromCreateEvent(salesOrder *models.SalesOrder, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromUpdateEvent(salesOrder *models.SalesOrder, ctx context.Context) *model.ErrorLog
	UpdateById(id int, request *models.SalesOrderUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderResponse, *model.ErrorLog)
	UpdateSODetailById(id int, request *models.SalesOrderDetailUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderDetail, *model.ErrorLog)
	UpdateSODetailBySOId(SoId int, request []*models.SalesOrderDetailUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) ([]*models.SalesOrder, *model.ErrorLog)
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
	salesOrder := &models.SalesOrder{}
	salesOrder.SalesOrderRequestMap(request, now)

	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType("open", "sales_order", false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return &models.SalesOrderResponse{}, getOrderStatusResult.ErrorLog
	}
	salesOrder.OrderStatusChanMap(getOrderStatusResult)

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

	salesOrder.OrderSourceChanMap(getOrderSourceResult)

	salesOrdersResponse := &models.SalesOrderResponse{
		SalesOrderStoreRequest: *request,
	}

	getBrandResultChan := make(chan *models.BrandChan)
	go u.brandRepository.GetByID(request.BrandID, false, ctx, getBrandResultChan)
	getBrandResult := <-getBrandResultChan

	if getBrandResult.Error != nil {
		return &models.SalesOrderResponse{}, getBrandResult.ErrorLog
	}

	salesOrder.BrandChanMap(getBrandResult)

	soCode = helper.GenerateSOCode(request.AgentID, getOrderSourceResult.OrderSource.Code)
	salesOrder.SoCode = soCode

	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(request.AgentID, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}
	salesOrder.AgentChanMap(getAgentResult)

	getStoreResultChan := make(chan *models.StoreChan)
	go u.storeRepository.GetByID(request.StoreID, false, ctx, getStoreResultChan)
	getStoreResult := <-getStoreResultChan

	if getStoreResult.Error != nil {
		errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}
	salesOrder.StoreChanMap(getStoreResult)

	getUserResultChan := make(chan *models.UserChan)
	go u.userRepository.GetByID(request.UserID, false, ctx, getUserResultChan)
	getUserResult := <-getUserResultChan

	if getUserResult.Error != nil {
		errorLogData := helper.WriteLog(getUserResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}
	salesOrder.UserChanMap(getUserResult)

	getSalesmanResultChan := make(chan *models.SalesmanChan)
	go u.salesmanRepository.GetByEmail(getUserResult.User.Email, false, ctx, getSalesmanResultChan)
	getSalesmanResult := <-getSalesmanResultChan

	if getSalesmanResult.Error != nil {
		errorLogData := helper.WriteLog(getSalesmanResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}
	salesOrder.SalesmanChanMap(getSalesmanResult)

	createSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.Insert(salesOrder, sqlTransaction, ctx, createSalesOrderResultChan)
	createSalesOrderResult := <-createSalesOrderResultChan

	if createSalesOrderResult.Error != nil {
		return &models.SalesOrderResponse{}, createSalesOrderResult.ErrorLog
	}
	salesOrdersResponse.SoResponseMap(salesOrder)

	var salesOrderDetailResponses []*models.SalesOrderDetailStoreResponse
	salesOrderDetails := []*models.SalesOrderDetail{}
	for _, soDetail := range request.SalesOrderDetails {
		soDetailCode, _ := helper.GenerateSODetailCode(int(createSalesOrderResult.ID), request.AgentID, soDetail.ProductID, soDetail.UomID)
		salesOrderDetail := &models.SalesOrderDetail{}
		salesOrderDetail.SalesOrderDetailStoreRequestMap(soDetail, now)
		salesOrderDetail.SalesOrderID = int(createSalesOrderResult.ID)
		salesOrderDetail.SoDetailCode = soDetailCode
		salesOrderDetail.Note = models.NullString{NullString: sql.NullString{String: request.Note, Valid: true}}

		createSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.Insert(salesOrderDetail, sqlTransaction, ctx, createSalesOrderDetailResultChan)
		createSalesOrderDetailResult := <-createSalesOrderDetailResultChan

		if createSalesOrderDetailResult.Error != nil {
			return &models.SalesOrderResponse{}, createSalesOrderDetailResult.ErrorLog
		}

		salesOrderDetailResponse := &models.SalesOrderDetailStoreResponse{
			SalesOrderDetailStoreRequest: *soDetail,
		}

		getProductResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetByID(soDetail.ProductID, false, ctx, getProductResultChan)
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
		go u.uomRepository.GetByID(soDetail.UomID, false, ctx, getUomResultChan)
		getUomResult := <-getUomResultChan

		if getUomResult.Error != nil {
			errorLogData := helper.WriteLog(getUomResult.Error, http.StatusInternalServerError, nil)
			return &models.SalesOrderResponse{}, errorLogData
		}

		salesOrderDetailResponse.UomCode = getUomResult.Uom.Code.String

		salesOrderDetailResponses = append(salesOrderDetailResponses, salesOrderDetailResponse)
		salesOrderDetails = append(salesOrderDetails, salesOrderDetail)
	}
	salesOrder.SalesOrderDetails = salesOrderDetails

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
	fmt.Println("message Create SO = ", string(messageKafka))
	err := u.kafkaClient.WriteToTopic(constants.CREATE_SALES_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	return salesOrdersResponse, nil
}

func (u *salesOrderUseCase) Get(request *models.SalesOrderRequest) (*models.SalesOrdersOpenSearchResponse, *model.ErrorLog) {
	getSalesOrdersResultChan := make(chan *models.SalesOrdersChan)
	go u.salesOrderOpenSearchRepository.Get(request, getSalesOrdersResultChan)
	getSalesOrdersResult := <-getSalesOrdersResultChan

	if getSalesOrdersResult.Error != nil {
		return &models.SalesOrdersOpenSearchResponse{}, getSalesOrdersResult.ErrorLog
	}

	var salesOrdersResult []*models.SalesOrderOpenSearchResponse
	for _, v := range getSalesOrdersResult.SalesOrders {
		var salesOrder models.SalesOrderOpenSearchResponse

		salesOrder.ID = v.ID
		salesOrder.AgentName = v.AgentName
		salesOrder.AgentEmail = v.AgentEmail
		salesOrder.AgentProvinceName = v.AgentProvinceName
		salesOrder.AgentCityName = v.AgentCityName
		salesOrder.AgentVillageName = v.AgentVillageName
		salesOrder.AgentAddress = v.AgentAddress
		salesOrder.AgentPhone = v.AgentPhone
		salesOrder.AgentMainMobilePhone = v.AgentMainMobilePhone

		salesOrder.StoreName = v.StoreName
		salesOrder.StoreCode = v.StoreCode
		salesOrder.StoreEmail = v.StoreEmail
		salesOrder.StoreProvinceName = v.StoreProvinceName
		salesOrder.StoreCityName = v.StoreCityName
		salesOrder.StoreDistrictName = v.StoreDistrictName
		salesOrder.StoreVillageName = v.StoreVillageName
		salesOrder.StoreAddress = v.StoreAddress
		salesOrder.StorePhone = v.StorePhone
		salesOrder.StoreMainMobilePhone = v.StoreMainMobilePhone

		salesOrder.BrandID = v.BrandID
		salesOrder.BrandName = v.BrandName

		salesOrder.UserFirstName = v.UserFirstName
		salesOrder.UserLastName = v.UserLastName
		salesOrder.UserEmail = v.UserEmail

		salesOrder.OrderSourceName = v.OrderSourceName
		salesOrder.OrderStatusName = v.OrderStatusName

		salesOrder.SoCode = v.SoCode
		salesOrder.SoDate = v.SoDate
		salesOrder.SoRefCode = v.SoRefCode
		salesOrder.SoRefDate = v.SoRefDate
		salesOrder.GLat = v.GLat
		salesOrder.GLong = v.GLong
		salesOrder.Note = v.Note
		salesOrder.ReferralCode = v.ReferralCode
		salesOrder.InternalComment = v.InternalComment
		salesOrder.TotalAmount = v.TotalAmount
		salesOrder.TotalTonase = v.TotalTonase

		var salesOrderDetails []*models.SalesOrderDetailOpenSearchResponse
		for _, x := range v.SalesOrderDetails {
			var salesOrderDetail models.SalesOrderDetailOpenSearchResponse

			salesOrderDetail.ID = x.ID
			salesOrderDetail.SalesOrderID = x.SalesOrderID
			salesOrderDetail.ProductID = x.ProductID

			var product models.ProductOpenSearchResponse
			product.ID = x.Product.ID
			product.Sku = x.Product.Sku
			product.AliasSku = x.Product.AliasSku
			product.ProductName = x.Product.ProductName
			product.Description = x.Product.Description
			product.CategoryID = x.Product.CategoryID

			salesOrderDetail.Product = &product

			var uom models.UomOpenSearchResponse
			salesOrderDetail.UomID = x.UomID
			uom.Name = x.Uom.Name
			uom.Code = x.Uom.Code

			salesOrderDetail.Uom = &uom

			salesOrderDetail.OrderStatusID = x.OrderStatusID
			var orderStatus models.OrderStatusOpenSearchResponse
			orderStatus.ID = x.OrderStatus.ID
			orderStatus.Name = x.OrderStatus.Name

			salesOrderDetail.OrderStatus = &orderStatus

			salesOrderDetail.SoDetailCode = x.SoDetailCode
			salesOrderDetail.Qty = x.Qty
			salesOrderDetail.SentQty = x.SentQty
			salesOrderDetail.ResidualQty = x.ResidualQty
			salesOrderDetail.Price = x.Price
			salesOrderDetail.Note = x.Note
			salesOrderDetail.CreatedAt = x.CreatedAt

			salesOrderDetails = append(salesOrderDetails, &salesOrderDetail)
		}

		salesOrder.SalesOrderDetails = salesOrderDetails

		salesOrder.SalesmanName = v.SalesmanName
		salesOrder.SalesmanEmail = v.SalesmanEmail
		salesOrder.CreatedAt = v.CreatedAt
		salesOrder.UpdatedAt = v.UpdatedAt

		salesOrdersResult = append(salesOrdersResult, &salesOrder)
	}
	salesOrders := &models.SalesOrdersOpenSearchResponse{
		SalesOrders: salesOrdersResult,
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

func (u *salesOrderUseCase) UpdateById(id int, request *models.SalesOrderUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderResponse, *model.ErrorLog) {
	now := time.Now()
	var soCode string

	getSalesOrderByIDResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.GetByID(id, false, ctx, getSalesOrderByIDResultChan)
	getSalesOrderByIDResult := <-getSalesOrderByIDResultChan

	if getSalesOrderByIDResult.Error != nil {
		return &models.SalesOrderResponse{}, getSalesOrderByIDResult.ErrorLog
	}

	// Check Order Status
	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByID(getSalesOrderByIDResult.SalesOrder.OrderStatusID, false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return &models.SalesOrderResponse{}, getOrderStatusResult.ErrorLog
	}
	getSalesOrderByIDResult.SalesOrder.OrderStatusChanMap(getOrderStatusResult)

	errorValidation := u.updateSOValidation(getSalesOrderByIDResult.SalesOrder.ID, getOrderStatusResult.OrderStatus.Name, ctx)

	if errorValidation != nil {
		errorLogData := helper.WriteLog(errorValidation, http.StatusBadRequest, "Ada kesalahan, silahkan coba lagi nanti")
		return &models.SalesOrderResponse{}, errorLogData
	}

	// Check Order Detail Status
	getOrderDetailStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType("open", "sales_order_detail", false, ctx, getOrderDetailStatusResultChan)
	getOrderDetailStatusResult := <-getOrderDetailStatusResultChan

	if getOrderDetailStatusResult.Error != nil {
		return &models.SalesOrderResponse{}, getOrderDetailStatusResult.ErrorLog
	}

	// Check Order Source
	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(request.OrderSourceID, false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getOrderSourceResult.Error != nil {
		return &models.SalesOrderResponse{}, getOrderSourceResult.ErrorLog
	}
	getSalesOrderByIDResult.SalesOrder.OrderSourceChanMap(getOrderSourceResult)

	salesOrdersResponse := &models.SalesOrderResponse{
		ID: id,
		// SalesOrderStoreRequest: models.SalesOrderStoreRequest{
		// 	SalesOrderTemplate: models.SalesOrderTemplate{
		// 		OrderSourceID:   request.OrderSourceID,
		// 		AgentID:         request.AgentID,
		// 		StoreID:         request.StoreID,
		// 		UserID:          request.UserID,
		// 		GLat:            request.GLat,
		// 		GLong:           request.GLong,
		// 		SoRefCode:       request.SoRefCode,
		// 		Note:            request.Note,
		// 		InternalComment: request.InternalComment,
		// 	},
		// 	OrderStatusID: getOrderStatusResult.OrderStatus.ID,
		// 	BrandID:       request.SalesOrderDetails[0].BrandID,
		// 	SoDate:        request.SoDate,
		// 	SoRefDate:     request.SoRefDate,
		// },
	}

	// Check Brand
	brandIds := []int{}

	for _, v := range request.SalesOrderDetails {
		brandIds = append(brandIds, v.BrandID)
	}
	checkIfBrandSame := helper.InSliceInt(brandIds, request.SalesOrderDetails[0].BrandID)

	if !checkIfBrandSame {
		errorLogData := helper.WriteLog(fmt.Errorf("The brand id must be the same"), http.StatusBadRequest, "Ada kesalahan, silahkan coba lagi nanti")
		return &models.SalesOrderResponse{}, errorLogData
	}

	getBrandResultChan := make(chan *models.BrandChan)
	go u.brandRepository.GetByID(request.SalesOrderDetails[0].BrandID, false, ctx, getBrandResultChan)
	getBrandResult := <-getBrandResultChan

	if getBrandResult.Error != nil {
		return &models.SalesOrderResponse{}, getBrandResult.ErrorLog
	}

	getSalesOrderByIDResult.SalesOrder.BrandChanMap(getBrandResult)

	// Check Agent
	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(request.AgentID, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	getSalesOrderByIDResult.SalesOrder.AgentChanMap(getAgentResult)

	// Check Store
	getStoreResultChan := make(chan *models.StoreChan)
	go u.storeRepository.GetByID(request.StoreID, false, ctx, getStoreResultChan)
	getStoreResult := <-getStoreResultChan

	if getStoreResult.Error != nil {
		errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	getSalesOrderByIDResult.SalesOrder.StoreChanMap(getStoreResult)

	// Check User Result
	getUserResultChan := make(chan *models.UserChan)
	go u.userRepository.GetByID(request.UserID, false, ctx, getUserResultChan)
	getUserResult := <-getUserResultChan

	if getUserResult.Error != nil {
		errorLogData := helper.WriteLog(getUserResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	getSalesOrderByIDResult.SalesOrder.UserChanMap(getUserResult)

	// Check Salesman
	getSalesmanResultChan := make(chan *models.SalesmanChan)
	go u.salesmanRepository.GetByEmail(getUserResult.User.Email, false, ctx, getSalesmanResultChan)
	getSalesmanResult := <-getSalesmanResultChan

	if getSalesmanResult.Error != nil {
		errorLogData := helper.WriteLog(getSalesmanResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	getSalesOrderByIDResult.SalesOrder.SalesmanChanMap(getSalesmanResult)

	var salesOrderDetailResponses []*models.SalesOrderDetailStoreResponse
	var soDetails []*models.SalesOrderDetail
	var totalAmount float64
	var totalTonase float64
	for _, v := range request.SalesOrderDetails {

		getSalesOrderDetailByIDResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.GetByID(v.ID, false, ctx, getSalesOrderDetailByIDResultChan)
		getSalesOrderDetailByIDResult := <-getSalesOrderDetailByIDResultChan

		if getSalesOrderDetailByIDResult.Error != nil {
			return &models.SalesOrderResponse{}, getSalesOrderDetailByIDResult.ErrorLog
		}

		salesOrderDetail := &models.SalesOrderDetail{
			UpdatedAt: &now,
		}
		soDetail := &models.SalesOrderDetailStoreRequest{
			SalesOrderDetailTemplate: v.SalesOrderDetailTemplate,
			SalesOrderId:             id,
			Price:                    v.Price,
		}
		salesOrderDetail.SalesOrderDetailStoreRequestMap(soDetail, now)

		updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.UpdateByID(v.ID, salesOrderDetail, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
		updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

		if updateSalesOrderDetailResult.Error != nil {
			return &models.SalesOrderResponse{}, updateSalesOrderDetailResult.ErrorLog
		}

		soDetail.OrderStatusId = getSalesOrderDetailByIDResult.SalesOrderDetail.OrderStatusID
		soDetail.SoDetailCode = getSalesOrderDetailByIDResult.SalesOrderDetail.SoDetailCode
		salesOrderDetailResponse := &models.SalesOrderDetailStoreResponse{
			ID:                           v.ID,
			SalesOrderDetailStoreRequest: *soDetail,
			CreatedAt:                    getSalesOrderDetailByIDResult.SalesOrderDetail.CreatedAt,
		}

		getProductResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetByID(v.ProductID, false, ctx, getProductResultChan)
		getProductResult := <-getProductResultChan

		if getProductResult.Error != nil {
			errorLogData := helper.WriteLog(getProductResult.Error, http.StatusInternalServerError, nil)
			return &models.SalesOrderResponse{}, errorLogData
		}

		salesOrderDetailResponses = append(salesOrderDetailResponses, salesOrderDetailResponse)
		soDetails = append(soDetails, salesOrderDetail)

		totalAmount = totalAmount + (v.Price * float64(v.Qty))
		totalTonase = totalTonase + (float64(v.Qty) * getProductResult.Product.NettWeight)

	}

	getSalesOrderByIDResult.SalesOrder.SalesOrderDetails = soDetails
	getSalesOrderByIDResult.SalesOrder.TotalAmount = totalAmount
	getSalesOrderByIDResult.SalesOrder.TotalTonase = totalTonase

	salesOrdersResponse.SalesOrderDetails = salesOrderDetailResponses

	salesOrderUpdateReq := &models.SalesOrder{
		OrderSourceID:   request.OrderSourceID,
		AgentID:         request.AgentID,
		StoreID:         request.StoreID,
		BrandID:         request.SalesOrderDetails[0].BrandID,
		UserID:          request.UserID,
		GLat:            models.NullFloat64{NullFloat64: sql.NullFloat64{Float64: request.GLat, Valid: true}},
		GLong:           models.NullFloat64{NullFloat64: sql.NullFloat64{Float64: request.GLong, Valid: true}},
		SoRefCode:       models.NullString{NullString: sql.NullString{String: request.SoRefCode, Valid: true}},
		SoDate:          request.SoDate,
		SoRefDate:       models.NullString{NullString: sql.NullString{String: request.SoRefDate, Valid: true}},
		Note:            models.NullString{NullString: sql.NullString{String: request.Note, Valid: true}},
		InternalComment: models.NullString{NullString: sql.NullString{String: request.InternalComment, Valid: true}},
		TotalAmount:     totalAmount,
		TotalTonase:     totalTonase,
		DeviceId:        models.NullString{NullString: sql.NullString{String: request.DeviceId, Valid: true}},
		ReferralCode:    models.NullString{NullString: sql.NullString{String: request.ReferralCode, Valid: true}},
		UpdatedAt:       &now,
		LatestUpdatedBy: request.UserID,
	}
	updateSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.UpdateByID(id, salesOrderUpdateReq, sqlTransaction, ctx, updateSalesOrderResultChan)
	updateSalesOrderResult := <-updateSalesOrderResultChan

	if updateSalesOrderResult.Error != nil {
		return &models.SalesOrderResponse{}, updateSalesOrderResult.ErrorLog
	}

	salesOrdersResponse.SoResponseMap(getSalesOrderByIDResult.SalesOrder)

	soCode = getSalesOrderByIDResult.SalesOrder.SoCode

	salesOrderLog := &models.SalesOrderLog{
		RequestID: request.RequestID,
		SoCode:    soCode,
		Data:      getSalesOrderByIDResult.SalesOrder,
		Status:    "0",
		CreatedAt: &now,
	}

	createSalesOrderLogResultChan := make(chan *models.SalesOrderLogChan)
	go u.salesOrderLogRepository.Insert(salesOrderLog, ctx, createSalesOrderLogResultChan)
	createSalesOrderLogResult := <-createSalesOrderLogResultChan

	if createSalesOrderLogResult.Error != nil {
		return &models.SalesOrderResponse{}, createSalesOrderLogResult.ErrorLog
	}

	keyKafka := []byte(getSalesOrderByIDResult.SalesOrder.SoCode)
	messageKafka, _ := json.Marshal(getSalesOrderByIDResult)
	err := u.kafkaClient.WriteToTopic(constants.UPDATE_SALES_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	return salesOrdersResponse, nil
}

func (u *salesOrderUseCase) UpdateSODetailById(id int, request *models.SalesOrderDetailUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderDetail, *model.ErrorLog) {
	now := time.Now()
	var soCode string

	// Check Order Status
	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType("open", "sales_order", false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return &models.SalesOrderDetail{}, getOrderStatusResult.ErrorLog
	}

	// Check Order Detail Status
	getOrderDetailStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType("open", "sales_order_detail", false, ctx, getOrderDetailStatusResultChan)
	getOrderDetailStatusResult := <-getOrderDetailStatusResultChan

	if getOrderDetailStatusResult.Error != nil {
		return &models.SalesOrderDetail{}, getOrderDetailStatusResult.ErrorLog
	}

	// Check Brand
	getBrandResultChan := make(chan *models.BrandChan)
	go u.brandRepository.GetByID(request.BrandID, false, ctx, getBrandResultChan)
	getBrandResult := <-getBrandResultChan

	if getBrandResult.Error != nil {
		return &models.SalesOrderDetail{}, getBrandResult.ErrorLog
	}

	salesOrderDetail := &models.SalesOrderDetail{
		ID:          request.ID,
		ProductID:   request.ProductID,
		UomID:       request.UomID,
		Qty:         request.Qty,
		SentQty:     request.SentQty,
		ResidualQty: request.ResidualQty,
		Price:       request.Price,
		Note:        models.NullString{NullString: sql.NullString{String: request.Note, Valid: true}},
		UpdatedAt:   &now,
	}

	updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
	go u.salesOrderDetailRepository.UpdateByID(request.ID, salesOrderDetail, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
	updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

	if updateSalesOrderDetailResult.Error != nil {
		return &models.SalesOrderDetail{}, updateSalesOrderDetailResult.ErrorLog
	}

	salesOrderDetail.BrandID = request.BrandID

	salesOrderLog := &models.SalesOrderLog{
		RequestID: "",
		SoCode:    soCode,
		Data:      salesOrderDetail,
		Status:    "0",
		CreatedAt: &now,
	}

	createSalesOrderLogResultChan := make(chan *models.SalesOrderLogChan)
	go u.salesOrderLogRepository.Insert(salesOrderLog, ctx, createSalesOrderLogResultChan)
	createSalesOrderLogResult := <-createSalesOrderLogResultChan

	if createSalesOrderLogResult.Error != nil {
		return &models.SalesOrderDetail{}, createSalesOrderLogResult.ErrorLog
	}

	keyKafka := []byte(soCode)
	messageKafka, _ := json.Marshal(salesOrderDetail)
	err := u.kafkaClient.WriteToTopic(constants.UPDATE_SALES_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.SalesOrderDetail{}, errorLogData
	}

	return salesOrderDetail, nil

}

func (u *salesOrderUseCase) UpdateSODetailBySOId(SoId int, request []*models.SalesOrderDetailUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) ([]*models.SalesOrder, *model.ErrorLog) {
	now := time.Now()
	var soCode string
	var response []*models.SalesOrder
	var salesOrder *models.SalesOrder

	getSalesOrderByIDResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.GetByID(SoId, false, ctx, getSalesOrderByIDResultChan)
	getSalesOrderByIDResult := <-getSalesOrderByIDResultChan

	if getSalesOrderByIDResult.Error != nil {
		return []*models.SalesOrder{}, getSalesOrderByIDResult.ErrorLog
	}

	salesOrder = getSalesOrderByIDResult.SalesOrder
	salesOrder.AgentProvinceID = 0
	salesOrder.AgentCityID = 0
	salesOrder.AgentDistrictID = 0
	salesOrder.AgentVillageID = 0
	salesOrder.StoreProvinceID = 0
	salesOrder.StoreCityID = 0
	salesOrder.StoreDistrictID = 0
	salesOrder.StoreVillageID = 0
	salesOrder.ReferralCode = models.NullString{NullString: sql.NullString{String: "", Valid: true}}
	salesOrder.DeviceId = models.NullString{NullString: sql.NullString{String: "", Valid: true}}
	salesOrder.SalesmanID = 0

	var soDetails []*models.SalesOrderDetail
	totalAmount := salesOrder.TotalAmount
	totalTonase := salesOrder.TotalTonase

	for _, v := range request {
		salesOrderDetail := &models.SalesOrderDetail{
			ID:          v.ID,
			ProductID:   v.ProductID,
			UomID:       v.UomID,
			Qty:         v.Qty,
			SentQty:     v.SentQty,
			ResidualQty: v.ResidualQty,
			Price:       v.Price,
			Note:        models.NullString{NullString: sql.NullString{String: v.Note, Valid: true}},
			UpdatedAt:   &now,
		}

		updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.UpdateByID(v.ID, salesOrderDetail, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
		updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

		if updateSalesOrderDetailResult.Error != nil {
			return []*models.SalesOrder{}, updateSalesOrderDetailResult.ErrorLog
		}

		getSalesOrderDetailByIDResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.GetByID(v.ID, false, ctx, getSalesOrderDetailByIDResultChan)
		getSalesOrderDetailByIDResult := <-getSalesOrderDetailByIDResultChan

		if getSalesOrderDetailByIDResult.Error != nil {
			return []*models.SalesOrder{}, getSalesOrderDetailByIDResult.ErrorLog
		}

		soDetails = append(soDetails, &models.SalesOrderDetail{
			ID:            v.ID,
			SalesOrderID:  SoId,
			ProductID:     v.ProductID,
			UomID:         v.UomID,
			OrderStatusID: getSalesOrderDetailByIDResult.SalesOrderDetail.OrderStatusID,
			SoDetailCode:  getSalesOrderDetailByIDResult.SalesOrderDetail.SoDetailCode,
			Qty:           v.Qty,
			ResidualQty:   v.ResidualQty,
			Price:         v.Price,
			Note:          models.NullString{NullString: sql.NullString{String: v.Note, Valid: true}},
			CreatedAt:     getSalesOrderDetailByIDResult.SalesOrderDetail.CreatedAt,
		})

		getProductByIDResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetByID(v.ID, false, ctx, getProductByIDResultChan)
		getProductByIDResult := <-getProductByIDResultChan

		if getProductByIDResult.Error != nil {
			return []*models.SalesOrder{}, getProductByIDResult.ErrorLog
		}

		totalAmount = totalAmount + (v.Price * float64(v.Qty))
		totalTonase = totalTonase + (float64(v.Qty) * getProductByIDResult.Product.NettWeight)

	}

	salesOrder.SalesOrderDetails = soDetails
	salesOrder.TotalAmount = totalAmount
	salesOrder.TotalTonase = totalTonase

	salesOrderUpdateReq := &models.SalesOrder{
		TotalAmount:     totalAmount,
		TotalTonase:     totalTonase,
		UpdatedAt:       &now,
		LatestUpdatedBy: getSalesOrderByIDResult.SalesOrder.UserID,
	}

	updateSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.UpdateByID(SoId, salesOrderUpdateReq, sqlTransaction, ctx, updateSalesOrderResultChan)
	updateSalesOrderResult := <-updateSalesOrderResultChan

	if updateSalesOrderResult.Error != nil {
		return []*models.SalesOrder{}, updateSalesOrderResult.ErrorLog
	}

	// Check Order Status
	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByID(salesOrder.OrderStatusID, false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return []*models.SalesOrder{}, getOrderStatusResult.ErrorLog
	}

	salesOrder.OrderStatusName = getOrderStatusResult.OrderStatus.Name

	// Check Order Source
	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(salesOrder.OrderSourceID, false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getOrderSourceResult.Error != nil {
		return []*models.SalesOrder{}, getOrderSourceResult.ErrorLog
	}

	salesOrder.OrderSourceName = getOrderSourceResult.OrderSource.SourceName

	getBrandResultChan := make(chan *models.BrandChan)
	go u.brandRepository.GetByID(salesOrder.BrandID, false, ctx, getBrandResultChan)
	getBrandResult := <-getBrandResultChan

	if getBrandResult.Error != nil {
		return []*models.SalesOrder{}, getBrandResult.ErrorLog
	}

	salesOrder.BrandName = getBrandResult.Brand.Name

	// Check Agent
	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(salesOrder.AgentID, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return []*models.SalesOrder{}, errorLogData
	}

	salesOrder.AgentName = models.NullString{NullString: sql.NullString{String: getAgentResult.Agent.Name, Valid: true}}
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
	go u.storeRepository.GetByID(salesOrder.StoreID, false, ctx, getStoreResultChan)
	getStoreResult := <-getStoreResultChan

	if getStoreResult.Error != nil {
		errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
		return []*models.SalesOrder{}, errorLogData
	}

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
	go u.userRepository.GetByID(salesOrder.UserID, false, ctx, getUserResultChan)
	getUserResult := <-getUserResultChan

	if getUserResult.Error != nil {
		errorLogData := helper.WriteLog(getUserResult.Error, http.StatusInternalServerError, nil)
		return []*models.SalesOrder{}, errorLogData
	}

	salesOrder.UserFirstName = getUserResult.User.FirstName
	salesOrder.UserLastName = getUserResult.User.LastName
	salesOrder.UserEmail = models.NullString{NullString: sql.NullString{String: getUserResult.User.Email, Valid: true}}

	// Check Salesman
	getSalesmanResultChan := make(chan *models.SalesmanChan)
	go u.salesmanRepository.GetByEmail(getUserResult.User.Email, false, ctx, getSalesmanResultChan)
	getSalesmanResult := <-getSalesmanResultChan

	if getSalesmanResult.Error != nil {
		errorLogData := helper.WriteLog(getSalesmanResult.Error, http.StatusInternalServerError, nil)
		return []*models.SalesOrder{}, errorLogData
	}

	salesOrder.SalesmanName = models.NullString{NullString: sql.NullString{String: getSalesmanResult.Salesman.Name, Valid: true}}
	salesOrder.SalesmanEmail = getSalesmanResult.Salesman.Email

	response = append(response, salesOrder)

	salesOrderLog := &models.SalesOrderLog{
		RequestID: "",
		SoCode:    soCode,
		Data:      salesOrder,
		Status:    "0",
		CreatedAt: &now,
	}

	createSalesOrderLogResultChan := make(chan *models.SalesOrderLogChan)
	go u.salesOrderLogRepository.Insert(salesOrderLog, ctx, createSalesOrderLogResultChan)
	createSalesOrderLogResult := <-createSalesOrderLogResultChan

	if createSalesOrderLogResult.Error != nil {
		return []*models.SalesOrder{}, createSalesOrderLogResult.ErrorLog
	}

	keyKafka := []byte(soCode)
	messageKafka, _ := json.Marshal(salesOrder)
	err := u.kafkaClient.WriteToTopic(constants.UPDATE_SALES_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return []*models.SalesOrder{}, errorLogData
	}

	return response, nil
}

func (u *salesOrderUseCase) updateSOValidation(salesOrderId int, orderStatusName string, ctx context.Context) error {

	if orderStatusName == "close" {
		return fmt.Errorf("Cannot update. Sales order status are close")
	}

	getDeliveryOrderByIDResultChan := make(chan *models.DeliveryOrdersChan)
	go u.deliveryOrderRepository.GetBySalesOrderID(salesOrderId, true, ctx, getDeliveryOrderByIDResultChan)
	getDeliveryOrderByIDResult := <-getDeliveryOrderByIDResultChan

	if getDeliveryOrderByIDResult.Total == 0 {

		return nil

	} else {

		for _, v := range getDeliveryOrderByIDResult.DeliveryOrders {

			getOrderStatusResultChan := make(chan *models.OrderStatusChan)
			go u.orderStatusRepository.GetByID(v.OrderStatusID, false, ctx, getOrderStatusResultChan)
			getOrderStatusResult := <-getOrderStatusResultChan

			if getOrderStatusResult.OrderStatus.Name != "cancel" {
				return fmt.Errorf("Cannot update. Order delivery must be cancel first")
			}
		}

	}

	return nil
}
