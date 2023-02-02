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
	Create(request *models.SalesOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) ([]*models.SalesOrder, *model.ErrorLog)
	Get(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetByID(request *models.SalesOrderRequest, withDetail bool, ctx context.Context) (*models.SalesOrder, *model.ErrorLog)
	GetByAgentID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetByStoreID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetBySalesmanID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetByOrderStatusID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetByOrderSourceID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	SyncToOpenSearchFromCreateEvent(salesOrder *models.SalesOrder, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromUpdateEvent(salesOrder *models.SalesOrder, ctx context.Context) *model.ErrorLog
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
	salesOrderOpenSearchRepository    openSearchRepositories.SalesOrderOpenSearchRepositoryInterface
	deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface
	kafkaClient                       kafkadbo.KafkaClientInterface
	db                                dbresolver.DB
	ctx                               context.Context
}

func InitSalesOrderUseCaseInterface(salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, orderSourceRepository repositories.OrderSourceRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, uomRepository repositories.UomRepositoryInterface, deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, salesOrderLogRepository mongoRepositories.SalesOrderLogRepositoryInterface, userRepository repositories.UserRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface, kafkaClient kafkadbo.KafkaClientInterface, db dbresolver.DB, ctx context.Context) SalesOrderUseCaseInterface {
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
		salesOrderOpenSearchRepository:    salesOrderOpenSearchRepository,
		deliveryOrderOpenSearchRepository: deliveryOrderOpenSearchRepository,
		kafkaClient:                       kafkaClient,
		db:                                db,
		ctx:                               ctx,
	}
}

func (u *salesOrderUseCase) Create(request *models.SalesOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) ([]*models.SalesOrder, *model.ErrorLog) {
	now := time.Now()
	var soCode string

	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType("open", "sales_order", false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return []*models.SalesOrder{}, getOrderStatusResult.ErrorLog
	}

	getOrderDetailStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType("open", "sales_order_detail", false, ctx, getOrderDetailStatusResultChan)
	getOrderDetailStatusResult := <-getOrderDetailStatusResultChan

	if getOrderDetailStatusResult.Error != nil {
		return []*models.SalesOrder{}, getOrderDetailStatusResult.ErrorLog
	}

	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(request.OrderSourceID, false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getOrderSourceResult.Error != nil {
		return []*models.SalesOrder{}, getOrderSourceResult.ErrorLog
	}

	brandIds := []int{}
	var salesOrderBrands map[int]*models.SalesOrder
	salesOrderBrands = map[int]*models.SalesOrder{}

	for _, v := range request.SalesOrderDetails {
		checkIfBrandExist := helper.InSliceInt(brandIds, v.BrandID)

		if checkIfBrandExist {
			salesOrderDetail := &models.SalesOrderDetail{
				ProductID:         v.ProductID,
				UomID:             v.UomID,
				OrderStatusID:     getOrderDetailStatusResult.OrderStatus.ID,
				Qty:               v.Qty,
				Price:             v.Price,
				ResidualQty:       0,
				SentQty:           0,
				IsDoneSyncToEs:    "0",
				Note:              models.NullString{NullString: sql.NullString{String: v.Note, Valid: true}},
				StartDateSyncToEs: &now,
				CreatedAt:         &now,
			}

			salesOrder := salesOrderBrands[v.BrandID]
			salesOrderDetails := salesOrder.SalesOrderDetails
			salesOrderDetails = append(salesOrderDetails, salesOrderDetail)
			salesOrder.SalesOrderDetails = salesOrderDetails
			salesOrderBrands[v.BrandID] = salesOrder

		} else {
			soCode = helper.GenerateSOCode(request.AgentID, getOrderSourceResult.OrderSource.Code)
			brandIds = append(brandIds, v.BrandID)

			salesOrderDetails := []*models.SalesOrderDetail{}
			salesOrderDetail := &models.SalesOrderDetail{
				ProductID:         v.ProductID,
				UomID:             v.UomID,
				OrderStatusID:     getOrderDetailStatusResult.OrderStatus.ID,
				Qty:               v.Qty,
				Price:             v.Price,
				ResidualQty:       0,
				SentQty:           0,
				IsDoneSyncToEs:    "0",
				Note:              models.NullString{NullString: sql.NullString{String: v.Note, Valid: true}},
				StartDateSyncToEs: &now,
				CreatedAt:         &now,
			}

			salesOrderDetails = append(salesOrderDetails, salesOrderDetail)

			salesOrder := &models.SalesOrder{
				CartID:            request.CartID,
				AgentID:           request.AgentID,
				StoreID:           request.StoreID,
				BrandID:           v.BrandID,
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
				IsDoneSyncToEs:    "0",
				CreatedAt:         &now,
				StartDateSyncToEs: &now,
				StartCreatedDate:  &now,
				SalesOrderDetails: salesOrderDetails,
			}

			salesOrderBrands[v.BrandID] = salesOrder
		}
	}

	salesOrders := []*models.SalesOrder{}

	for _, v := range salesOrderBrands {
		createSalesOrderResultChan := make(chan *models.SalesOrderChan)
		go u.salesOrderRepository.Insert(v, sqlTransaction, ctx, createSalesOrderResultChan)
		createSalesOrderResult := <-createSalesOrderResultChan

		if createSalesOrderResult.Error != nil {
			return []*models.SalesOrder{}, createSalesOrderResult.ErrorLog
		}

		for _, x := range v.SalesOrderDetails {
			soDetailCode, _ := helper.GenerateSODetailCode(int(createSalesOrderResult.ID), v.AgentID, x.ProductID, x.UomID)
			x.SalesOrderID = int(createSalesOrderResult.ID)
			x.SoDetailCode = soDetailCode
			createSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
			go u.salesOrderDetailRepository.Insert(x, sqlTransaction, ctx, createSalesOrderDetailResultChan)
			createSalesOrderDetailResult := <-createSalesOrderDetailResultChan

			if createSalesOrderDetailResult.Error != nil {
				return []*models.SalesOrder{}, createSalesOrderDetailResult.ErrorLog
			}
		}

		salesOrderLog := &models.SalesOrderLog{
			RequestID: request.RequestID,
			SoCode:    soCode,
			Data:      v,
			Status:    "0",
			CreatedAt: &now,
		}

		createSalesOrderLogResultChan := make(chan *models.SalesOrderLogChan)
		go u.salesOrderLogRepository.Insert(salesOrderLog, ctx, createSalesOrderLogResultChan)
		createSalesOrderLogResult := <-createSalesOrderLogResultChan

		if createSalesOrderLogResult.Error != nil {
			return []*models.SalesOrder{}, createSalesOrderLogResult.ErrorLog
		}

		getAgentResultChan := make(chan *models.AgentChan)
		go u.agentRepository.GetByID(v.AgentID, false, ctx, getAgentResultChan)
		getAgentResult := <-getAgentResultChan

		if getAgentResult.Error != nil {
			errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
			return []*models.SalesOrder{}, errorLogData
		}

		v.Agent = getAgentResult.Agent
		v.AgentName = models.NullString{NullString: sql.NullString{String: getAgentResult.Agent.Name, Valid: true}}
		v.AgentEmail = getAgentResult.Agent.Email
		v.AgentProvinceName = getAgentResult.Agent.ProvinceName
		v.AgentCityName = getAgentResult.Agent.CityName
		v.AgentDistrictName = getAgentResult.Agent.DistrictName
		v.AgentVillageName = getAgentResult.Agent.VillageName
		v.AgentAddress = getAgentResult.Agent.Address
		v.AgentPhone = getAgentResult.Agent.Phone
		v.AgentMainMobilePhone = getAgentResult.Agent.MainMobilePhone

		getStoreResultChan := make(chan *models.StoreChan)
		go u.storeRepository.GetByID(v.StoreID, false, ctx, getStoreResultChan)
		getStoreResult := <-getStoreResultChan

		if getStoreResult.Error != nil {
			errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
			return []*models.SalesOrder{}, errorLogData
		}

		v.Store = getStoreResult.Store
		v.StoreName = getStoreResult.Store.Name
		v.StoreCode = getStoreResult.Store.StoreCode
		v.StoreEmail = getStoreResult.Store.Email
		v.StoreProvinceName = getStoreResult.Store.ProvinceName
		v.StoreCityName = getStoreResult.Store.CityName
		v.StoreDistrictName = getStoreResult.Store.DistrictName
		v.StoreVillageName = getStoreResult.Store.VillageName
		v.StoreAddress = getStoreResult.Store.Address
		v.StorePhone = getStoreResult.Store.Phone
		v.StoreMainMobilePhone = getStoreResult.Store.MainMobilePhone

		getBrandResultChan := make(chan *models.BrandChan)
		go u.brandRepository.GetByID(v.BrandID, false, ctx, getBrandResultChan)
		getBrandResult := <-getBrandResultChan

		if getBrandResult.Error != nil {
			errorLogData := helper.WriteLog(getBrandResult.Error, http.StatusInternalServerError, nil)
			return []*models.SalesOrder{}, errorLogData
		}

		v.Brand = getBrandResult.Brand
		v.BrandName = getBrandResult.Brand.Name
		v.OrderSource = getOrderSourceResult.OrderSource
		v.OrderSourceName = getOrderSourceResult.OrderSource.SourceName
		v.OrderStatus = getOrderStatusResult.OrderStatus
		v.OrderStatusName = getOrderStatusResult.OrderStatus.Name

		getUserResultChan := make(chan *models.UserChan)
		go u.userRepository.GetByID(v.UserID, false, ctx, getUserResultChan)
		getUserResult := <-getUserResultChan

		if getUserResult.Error != nil {
			errorLogData := helper.WriteLog(getUserResult.Error, http.StatusInternalServerError, nil)
			return []*models.SalesOrder{}, errorLogData
		}

		v.User = getUserResult.User
		v.UserFirstName = getUserResult.User.FirstName
		v.UserLastName = getUserResult.User.LastName
		v.UserEmail = models.NullString{NullString: sql.NullString{String: getUserResult.User.Email, Valid: true}}

		if getUserResult.User.RoleID.String == "3" {
			getSalesmanResultChan := make(chan *models.SalesmanChan)
			go u.salesmanRepository.GetByEmail(getUserResult.User.Email, false, ctx, getSalesmanResultChan)
			getSalesmanResult := <-getSalesmanResultChan

			if getSalesmanResult.Error != nil {
				errorLogData := helper.WriteLog(getSalesmanResult.Error, http.StatusInternalServerError, nil)
				return []*models.SalesOrder{}, errorLogData
			}

			v.Salesman = getSalesmanResult.Salesman
			v.SalesmanName = models.NullString{NullString: sql.NullString{String: getSalesmanResult.Salesman.Name, Valid: true}}
			v.SalesmanEmail = getSalesmanResult.Salesman.Email
		}

		keyKafka := []byte(v.SoCode)
		messageKafka, _ := json.Marshal(v)
		err := u.kafkaClient.WriteToTopic(constants.CREATE_SALES_ORDER_TOPIC, keyKafka, messageKafka)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return []*models.SalesOrder{}, errorLogData
		}

		salesOrders = append(salesOrders, v)
	}

	return salesOrders, nil
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
