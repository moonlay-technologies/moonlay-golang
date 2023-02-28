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
	"strconv"
	"strings"
	"time"

	"github.com/bxcodec/dbresolver"
)

type DeliveryOrderUseCaseInterface interface {
	Create(request *models.DeliveryOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrder, *model.ErrorLog)
	UpdateByID(ID int, request *models.DeliveryOrderUpdateByIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrder, *model.ErrorLog)
	UpdateDODetailByID(ID int, request *models.DeliveryOrderDetailUpdateByIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderDetail, *model.ErrorLog)
	UpdateDoDetailByDeliveryOrderID(deliveryOrderID int, request []*models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderDetails, *model.ErrorLog)
	Get(request *models.DeliveryOrderRequest) (*models.DeliveryOrdersOpenSearchResponse, *model.ErrorLog)
	GetByID(request *models.DeliveryOrderRequest, ctx context.Context) (*models.DeliveryOrder, *model.ErrorLog)
	GetByIDWithDetail(request *models.DeliveryOrderRequest, ctx context.Context) (*models.DeliveryOrderOpenSearchResponse, *model.ErrorLog)
	GetByAgentID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)
	GetByStoreID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)
	GetBySalesmanID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)
	GetBySalesmansID(request *models.DeliveryOrderRequest) (*models.DeliveryOrdersOpenSearchResponses, *model.ErrorLog)
	GetByOrderStatusID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)
	GetByOrderSourceID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)
	DeleteByID(deliveryOrderId int) *model.ErrorLog
	SyncToOpenSearchFromCreateEvent(deliveryOrder *models.DeliveryOrder, salesOrderUseCase SalesOrderUseCaseInterface, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromUpdateEvent(deliveryOrder *models.DeliveryOrder, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromDeleteEvent(deliveryOrderId *int, ctx context.Context) *model.ErrorLog
}

type deliveryOrderUseCase struct {
	deliveryOrderRepository           repositories.DeliveryOrderRepositoryInterface
	deliveryOrderDetailRepository     repositories.DeliveryOrderDetailRepositoryInterface
	salesOrderRepository              repositories.SalesOrderRepositoryInterface
	salesOrderDetailRepository        repositories.SalesOrderDetailRepositoryInterface
	orderStatusRepository             repositories.OrderStatusRepositoryInterface
	orderSourceRepository             repositories.OrderSourceRepositoryInterface
	warehouseRepository               repositories.WarehouseRepositoryInterface
	brandRepository                   repositories.BrandRepositoryInterface
	uomRepository                     repositories.UomRepositoryInterface
	agentRepository                   repositories.AgentRepositoryInterface
	storeRepository                   repositories.StoreRepositoryInterface
	productRepository                 repositories.ProductRepositoryInterface
	userRepository                    repositories.UserRepositoryInterface
	salesmanRepository                repositories.SalesmanRepositoryInterface
	deliveryOrderLogRepository        mongoRepositories.DeliveryOrderLogRepositoryInterface
	deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface
	salesOrderOpenSearchRepository    openSearchRepositories.SalesOrderOpenSearchRepositoryInterface
	salesOrderUseCase                 SalesOrderUseCaseInterface
	SalesOrderOpenSearchUseCase       SalesOrderOpenSearchUseCaseInterface
	kafkaClient                       kafkadbo.KafkaClientInterface
	ValidationRepository              repositories.RequestValidationRepositoryInterface
	db                                dbresolver.DB
	ctx                               context.Context
}

func InitDeliveryOrderUseCaseInterface(deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, deliveryOrderDetailRepository repositories.DeliveryOrderDetailRepositoryInterface, salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, orderSourceRepository repositories.OrderSourceRepositoryInterface, warehouseRepository repositories.WarehouseRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, uomRepository repositories.UomRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, userRepository repositories.UserRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, deliveryOrderLogRepository mongoRepositories.DeliveryOrderLogRepositoryInterface, deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, salesOrderUseCase SalesOrderUseCaseInterface, salesOrderOpenSearchUseCase SalesOrderOpenSearchUseCaseInterface, kafkaClient kafkadbo.KafkaClientInterface, ValidationRepository repositories.RequestValidationRepositoryInterface, db dbresolver.DB, ctx context.Context) DeliveryOrderUseCaseInterface {
	return &deliveryOrderUseCase{
		deliveryOrderRepository:           deliveryOrderRepository,
		deliveryOrderDetailRepository:     deliveryOrderDetailRepository,
		salesOrderRepository:              salesOrderRepository,
		salesOrderDetailRepository:        salesOrderDetailRepository,
		orderStatusRepository:             orderStatusRepository,
		orderSourceRepository:             orderSourceRepository,
		warehouseRepository:               warehouseRepository,
		brandRepository:                   brandRepository,
		uomRepository:                     uomRepository,
		productRepository:                 productRepository,
		userRepository:                    userRepository,
		salesmanRepository:                salesmanRepository,
		agentRepository:                   agentRepository,
		storeRepository:                   storeRepository,
		deliveryOrderLogRepository:        deliveryOrderLogRepository,
		deliveryOrderOpenSearchRepository: deliveryOrderOpenSearchRepository,
		salesOrderOpenSearchRepository:    salesOrderOpenSearchRepository,
		salesOrderUseCase:                 salesOrderUseCase,
		SalesOrderOpenSearchUseCase:       salesOrderOpenSearchUseCase,
		kafkaClient:                       kafkaClient,
		ValidationRepository:              ValidationRepository,
		db:                                db,
		ctx:                               ctx,
	}
}

func (u *deliveryOrderUseCase) Create(request *models.DeliveryOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrder, *model.ErrorLog) {
	now := time.Now()

	getSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.GetByID(request.SalesOrderID, false, ctx, getSalesOrderResultChan)
	getSalesOrderResult := <-getSalesOrderResultChan

	if getSalesOrderResult.Error != nil {
		return &models.DeliveryOrder{}, getSalesOrderResult.ErrorLog
	}

	getBrandResultChan := make(chan *models.BrandChan)
	go u.brandRepository.GetByID(getSalesOrderResult.SalesOrder.BrandID, false, ctx, getBrandResultChan)
	getBrandResult := <-getBrandResultChan

	if getBrandResult.Error != nil {
		return &models.DeliveryOrder{}, getBrandResult.ErrorLog
	}

	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType("open", "delivery_order", false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return &models.DeliveryOrder{}, getOrderStatusResult.ErrorLog
	}

	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetBySourceName("manager", false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getOrderSourceResult.Error != nil {
		return &models.DeliveryOrder{}, getOrderSourceResult.ErrorLog
	}

	getWarehouseResultChan := make(chan *models.WarehouseChan)
	go u.warehouseRepository.GetByID(request.WarehouseID, false, ctx, getWarehouseResultChan)
	getWarehouseResult := <-getWarehouseResultChan

	if getWarehouseResult.Error != nil {
		return &models.DeliveryOrder{}, getWarehouseResult.ErrorLog
	}

	getSalesOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(getSalesOrderResult.SalesOrder.OrderSourceID, false, ctx, getSalesOrderSourceResultChan)
	getSalesOrderSourceResult := <-getSalesOrderSourceResultChan

	if getSalesOrderSourceResult.Error != nil {
		return &models.DeliveryOrder{}, getSalesOrderSourceResult.ErrorLog
	}

	getSalesOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByID(getSalesOrderResult.SalesOrder.OrderStatusID, false, ctx, getSalesOrderStatusResultChan)
	getSalesOrderStatusResult := <-getSalesOrderStatusResultChan

	if getSalesOrderStatusResult.Error != nil {
		return &models.DeliveryOrder{}, getSalesOrderStatusResult.ErrorLog
	}

	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(getSalesOrderResult.SalesOrder.AgentID, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		return &models.DeliveryOrder{}, getAgentResult.ErrorLog
	}

	getStoreResultChan := make(chan *models.StoreChan)
	go u.storeRepository.GetByID(getSalesOrderResult.SalesOrder.StoreID, false, ctx, getStoreResultChan)
	getStoreResult := <-getStoreResultChan

	if getStoreResult.Error != nil {
		return &models.DeliveryOrder{}, getStoreResult.ErrorLog
	}

	getUserResultChan := make(chan *models.UserChan)
	go u.userRepository.GetByID(getSalesOrderResult.SalesOrder.UserID, false, ctx, getUserResultChan)
	getUserResult := <-getUserResultChan

	if getUserResult.Error != nil {
		return &models.DeliveryOrder{}, getUserResult.ErrorLog
	}

	getSalesmanResultChan := make(chan *models.SalesmanChan)
	go u.salesmanRepository.GetByEmail(getUserResult.User.Email, false, ctx, getSalesmanResultChan)
	getSalesmanResult := <-getSalesmanResultChan

	if getSalesmanResult.Error != nil {
		// ignore null salesman
		if getSalesmanResult.Error.Error() != "salesman data not found" {
			return &models.DeliveryOrder{}, getSalesmanResult.ErrorLog
		}
	}

	deliveryOrder := &models.DeliveryOrder{}

	deliveryOrder.DeliveryOrderStoreRequestMap(request, now)
	deliveryOrder.WarehouseChanMap(getWarehouseResult)
	deliveryOrder.OrderStatus = getOrderStatusResult.OrderStatus
	deliveryOrder.OrderStatusID = getOrderStatusResult.OrderStatus.ID
	deliveryOrder.OrderSource = getOrderSourceResult.OrderSource
	deliveryOrder.OrderSourceID = getOrderSourceResult.OrderSource.ID
	deliveryOrder.Agent = getAgentResult.Agent
	deliveryOrder.AgentID = getAgentResult.Agent.ID
	deliveryOrder.AgentName = getAgentResult.Agent.Name
	deliveryOrder.Store = getStoreResult.Store
	deliveryOrder.StoreID = getStoreResult.Store.ID
	deliveryOrder.DoCode = helper.GenerateDOCode(getAgentResult.Agent.ID, getOrderSourceResult.OrderSource.Code)
	deliveryOrder.DoDate = now.Format("2006-01-02")
	deliveryOrder.CreatedBy = getUserResult.User.ID
	deliveryOrder.SalesOrder = getSalesOrderResult.SalesOrder
	deliveryOrder.Brand = getBrandResult.Brand
	if getSalesmanResult.Salesman != nil {
		deliveryOrder.Salesman = getSalesmanResult.Salesman
	}

	createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderRepository.Insert(deliveryOrder, sqlTransaction, ctx, createDeliveryOrderResultChan)
	createDeliveryOrderResult := <-createDeliveryOrderResultChan

	if createDeliveryOrderResult.Error != nil {
		return &models.DeliveryOrder{}, createDeliveryOrderResult.ErrorLog
	}

	deliveryOrderDetails := []*models.DeliveryOrderDetail{}
	totalResidualQty := 0
	for _, doDetail := range request.DeliveryOrderDetails {
		getSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.GetByID(doDetail.SoDetailID, false, ctx, getSalesOrderDetailResultChan)
		getSalesOrderDetailResult := <-getSalesOrderDetailResultChan

		if getSalesOrderDetailResult.Error != nil {
			return &models.DeliveryOrder{}, getSalesOrderDetailResult.ErrorLog
		}

		getOrderStatusDetailResultChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID, false, ctx, getOrderStatusDetailResultChan)
		getOrderStatusDetailResult := <-getOrderStatusDetailResultChan

		if getOrderStatusDetailResult.Error != nil {
			return &models.DeliveryOrder{}, getOrderStatusDetailResult.ErrorLog
		}

		getSalesOrderDetailResult.SalesOrderDetail.UpdatedAt = &now
		getSalesOrderDetailResult.SalesOrderDetail.SentQty += doDetail.Qty
		getSalesOrderDetailResult.SalesOrderDetail.ResidualQty -= doDetail.Qty
		totalResidualQty += getSalesOrderDetailResult.SalesOrderDetail.ResidualQty

		if getSalesOrderDetailResult.SalesOrderDetail.ResidualQty == 0 {
			getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID = 14
		} else {
			getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID = 13
		}

		getOrderStatusSODetailResultChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID, false, ctx, getOrderStatusSODetailResultChan)
		getOrderStatusSODetailResult := <-getOrderStatusSODetailResultChan

		if getOrderStatusSODetailResult.Error != nil {
			return &models.DeliveryOrder{}, getOrderStatusSODetailResult.ErrorLog
		}
		getSalesOrderDetailResult.SalesOrderDetail.OrderStatusName = getOrderStatusSODetailResult.OrderStatus.Name
		getSalesOrderDetailResult.SalesOrderDetail.OrderStatus = getOrderStatusSODetailResult.OrderStatus

		getProductDetailResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.ProductID, false, ctx, getProductDetailResultChan)
		getProductDetailResult := <-getProductDetailResultChan

		if getProductDetailResult.Error != nil {
			return &models.DeliveryOrder{}, getProductDetailResult.ErrorLog
		}

		getUomDetailResultChan := make(chan *models.UomChan)
		go u.uomRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.UomID, false, ctx, getUomDetailResultChan)
		getUomDetailResult := <-getUomDetailResultChan

		if getUomDetailResult.Error != nil {
			return &models.DeliveryOrder{}, getUomDetailResult.ErrorLog
		}

		doDetailCode, _ := helper.GenerateDODetailCode(createDeliveryOrderResult.DeliveryOrder.ID, getAgentResult.Agent.ID, getProductDetailResult.Product.ID, getUomDetailResult.Uom.ID)

		deliveryOrderDetail := &models.DeliveryOrderDetail{}
		deliveryOrderDetail.DeliveryOrderDetailStoreRequestMap(doDetail, now)
		deliveryOrderDetail.DeliveryOrderID = int(createDeliveryOrderResult.ID)
		deliveryOrderDetail.BrandID = getBrandResult.Brand.ID
		deliveryOrderDetail.DoDetailCode = doDetailCode
		deliveryOrderDetail.Note = models.NullString{NullString: sql.NullString{String: doDetail.Note, Valid: true}}
		deliveryOrderDetail.Uom = getUomDetailResult.Uom
		deliveryOrderDetail.ProductChanMap(getProductDetailResult)
		deliveryOrderDetail.SalesOrderDetailChanMap(getSalesOrderDetailResult)
		deliveryOrderDetail.OrderStatusID = getOrderStatusDetailResult.OrderStatus.ID
		deliveryOrderDetail.OrderStatusName = getOrderStatusDetailResult.OrderStatus.Name
		deliveryOrderDetail.OrderStatus = getOrderStatusDetailResult.OrderStatus

		createDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
		go u.deliveryOrderDetailRepository.Insert(deliveryOrderDetail, sqlTransaction, ctx, createDeliveryOrderDetailResultChan)
		createDeliveryOrderDetailResult := <-createDeliveryOrderDetailResultChan

		if createDeliveryOrderDetailResult.Error != nil {
			return &models.DeliveryOrder{}, createDeliveryOrderDetailResult.ErrorLog
		}

		updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.UpdateByID(getSalesOrderDetailResult.SalesOrderDetail.ID, getSalesOrderDetailResult.SalesOrderDetail, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
		updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

		if updateSalesOrderDetailResult.Error != nil {
			return &models.DeliveryOrder{}, updateSalesOrderDetailResult.ErrorLog
		}

		deliveryOrderDetail.ID = int(createDeliveryOrderDetailResult.ID)
		deliveryOrderDetails = append(deliveryOrderDetails, deliveryOrderDetail)
	}

	deliveryOrder.DeliveryOrderDetails = deliveryOrderDetails
	if totalResidualQty == 0 {
		getSalesOrderResult.SalesOrder.OrderStatusID = 8
	} else {
		getSalesOrderResult.SalesOrder.OrderStatusID = 7
	}

	getOrderStatusSOResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByID(getSalesOrderResult.SalesOrder.OrderStatusID, false, ctx, getOrderStatusSOResultChan)
	getOrderStatusSODetailResult := <-getOrderStatusSOResultChan

	if getOrderStatusSODetailResult.Error != nil {
		return &models.DeliveryOrder{}, getOrderStatusSODetailResult.ErrorLog
	}
	getSalesOrderResult.SalesOrder.OrderStatus = getOrderStatusSODetailResult.OrderStatus
	getSalesOrderResult.SalesOrder.SoDate = ""
	getSalesOrderResult.SalesOrder.SoRefDate = models.NullString{}
	getSalesOrderResult.SalesOrder.UpdatedAt = &now

	updateSalesOrderChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.UpdateByID(getSalesOrderResult.SalesOrder.ID, getSalesOrderResult.SalesOrder, sqlTransaction, u.ctx, updateSalesOrderChan)
	updateSalesOrderResult := <-updateSalesOrderChan

	if updateSalesOrderResult.ErrorLog != nil {
		err := sqlTransaction.Rollback()
		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return &models.DeliveryOrder{}, errorLogData
		}
		return &models.DeliveryOrder{}, updateSalesOrderResult.ErrorLog
	}

	deliveryOrderLog := &models.DeliveryOrderLog{
		RequestID: request.RequestID,
		DoCode:    deliveryOrder.DoCode,
		Data:      deliveryOrder,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		Action:    constants.LOG_ACTION_MONGO_INSERT,
		CreatedAt: &now,
	}

	createDeliveryOrderLogResultChan := make(chan *models.DeliveryOrderLogChan)
	go u.deliveryOrderLogRepository.Insert(deliveryOrderLog, ctx, createDeliveryOrderLogResultChan)
	createDeliveryOrderLogResult := <-createDeliveryOrderLogResultChan

	if createDeliveryOrderLogResult.Error != nil {
		return &models.DeliveryOrder{}, createDeliveryOrderLogResult.ErrorLog
	}

	keyKafka := []byte(deliveryOrder.DoCode)
	messageKafka, _ := json.Marshal(deliveryOrder)
	err := u.kafkaClient.WriteToTopic(constants.CREATE_DELIVERY_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.DeliveryOrder{}, errorLogData
	}

	return deliveryOrder, nil
}

func (u *deliveryOrderUseCase) UpdateByID(ID int, request *models.DeliveryOrderUpdateByIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrder, *model.ErrorLog) {
	now := time.Now()

	getDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderRepository.GetByID(ID, false, ctx, getDeliveryOrderResultChan)
	getDeliveryOrderResult := <-getDeliveryOrderResultChan

	if getDeliveryOrderResult.Error != nil {
		return &models.DeliveryOrder{}, getDeliveryOrderResult.ErrorLog
	}

	getWarehouseResultChan := make(chan *models.WarehouseChan)
	go u.warehouseRepository.GetByID(request.WarehouseID, false, ctx, getWarehouseResultChan)
	getWarehouseResult := <-getWarehouseResultChan

	if getWarehouseResult.Error != nil {
		return &models.DeliveryOrder{}, getWarehouseResult.ErrorLog
	}

	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(request.OrderSourceID, false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getDeliveryOrderResult.Error != nil {
		return &models.DeliveryOrder{}, getOrderSourceResult.ErrorLog
	}

	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByID(request.OrderStatusID, false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return &models.DeliveryOrder{}, getOrderStatusResult.ErrorLog
	}

	orderStatusID := int(getOrderStatusResult.OrderStatus.ID)
	deliveryOrder := &models.DeliveryOrder{
		WarehouseID:       int(getWarehouseResult.Warehouse.ID),
		OrderSourceID:     int(getOrderSourceResult.OrderSource.ID),
		OrderStatusID:     orderStatusID,
		DoRefCode:         models.NullString{NullString: sql.NullString{String: request.DoRefCode, Valid: true}},
		DoRefDate:         models.NullString{NullString: sql.NullString{String: request.DoRefDate, Valid: true}},
		DriverName:        models.NullString{NullString: sql.NullString{String: request.DriverName, Valid: true}},
		PlatNumber:        models.NullString{NullString: sql.NullString{String: request.PlatNumber, Valid: true}},
		IsDoneSyncToEs:    "0",
		StartDateSyncToEs: &now,
		EndCreatedDate:    &now,
		UpdatedAt:         &now,
	}

	createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderRepository.UpdateByID(getDeliveryOrderResult.DeliveryOrder.ID, deliveryOrder, sqlTransaction, ctx, createDeliveryOrderResultChan)
	createDeliveryOrderResult := <-createDeliveryOrderResultChan

	if createDeliveryOrderResult.Error != nil {
		return &models.DeliveryOrder{}, createDeliveryOrderResult.ErrorLog
	}

	getDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailsChan)
	go u.deliveryOrderDetailRepository.GetByDeliveryOrderID(ID, false, ctx, getDeliveryOrderDetailResultChan)
	getDeliveryOrderDetailResult := <-getDeliveryOrderDetailResultChan

	if getDeliveryOrderDetailResult.Error != nil {
		return &models.DeliveryOrder{}, getDeliveryOrderDetailResult.ErrorLog
	}

	deliveryOrderDetails := []*models.DeliveryOrderDetail{}
	for _, v := range request.DeliveryOrderDetails {
		deliveryOrderDetail := &models.DeliveryOrderDetail{
			Qty:       v.Qty,
			Note:      models.NullString{NullString: sql.NullString{String: v.Note, Valid: true}},
			UpdatedAt: &now,
		}
		for _, x := range getDeliveryOrderDetailResult.DeliveryOrderDetails {
			createDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
			go u.deliveryOrderDetailRepository.UpdateByID(int(x.ID), deliveryOrderDetail, sqlTransaction, ctx, createDeliveryOrderDetailResultChan)
			createDeliveryOrderDetailResult := <-createDeliveryOrderDetailResultChan

			if createDeliveryOrderDetailResult.Error != nil {
				return &models.DeliveryOrder{}, createDeliveryOrderDetailResult.ErrorLog
			}
			deliveryOrderDetail.ID = int(createDeliveryOrderDetailResult.ID)
			deliveryOrderDetails = append(deliveryOrderDetails, deliveryOrderDetail)
		}
	}

	deliveryOrder.DeliveryOrderDetails = deliveryOrderDetails
	deliveryOrderLog := &models.DeliveryOrderLog{
		RequestID: request.RequestID,
		DoCode:    createDeliveryOrderResult.DeliveryOrder.DoCode,
		Data:      deliveryOrder,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		Action:    constants.LOG_ACTION_MONGO_UPDATE,
		CreatedAt: &now,
	}

	createDeliveryOrderLogResultChan := make(chan *models.DeliveryOrderLogChan)
	go u.deliveryOrderLogRepository.Insert(deliveryOrderLog, ctx, createDeliveryOrderLogResultChan)
	createDeliveryOrderLogResult := <-createDeliveryOrderLogResultChan

	if createDeliveryOrderLogResult.Error != nil {
		errorLogData := helper.WriteLog(createDeliveryOrderLogResult.Error, http.StatusInternalServerError, nil)
		return &models.DeliveryOrder{}, errorLogData
	}

	keyKafka := []byte(deliveryOrder.DoCode)
	messageKafka, _ := json.Marshal(deliveryOrder)
	err := u.kafkaClient.WriteToTopic(constants.UPDATE_DELIVERY_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.DeliveryOrder{}, errorLogData
	}

	return deliveryOrder, nil
}

func (u *deliveryOrderUseCase) UpdateDODetailByID(ID int, request *models.DeliveryOrderDetailUpdateByIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderDetail, *model.ErrorLog) {
	now := time.Now()

	getDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailsChan)
	go u.deliveryOrderDetailRepository.GetByID(ID, false, ctx, getDeliveryOrderDetailResultChan)
	getDeliveryOrderDetailResult := <-getDeliveryOrderDetailResultChan

	if getDeliveryOrderDetailResult.Error != nil {
		return &models.DeliveryOrderDetail{}, getDeliveryOrderDetailResult.ErrorLog
	}

	deliveryOrderDetail := &models.DeliveryOrderDetail{
		ID:                ID,
		Qty:               request.Qty,
		Note:              models.NullString{NullString: sql.NullString{String: request.Note, Valid: true}},
		IsDoneSyncToEs:    "0",
		StartDateSyncToEs: &now,
		EndDateSyncToEs:   &now,
		UpdatedAt:         &now,
	}

	updateDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
	go u.deliveryOrderDetailRepository.UpdateByID(ID, deliveryOrderDetail, sqlTransaction, ctx, updateDeliveryOrderDetailResultChan)
	updateDeliveryOrderDetailResult := <-updateDeliveryOrderDetailResultChan

	if updateDeliveryOrderDetailResult.Error != nil {
		return &models.DeliveryOrderDetail{}, updateDeliveryOrderDetailResult.ErrorLog
	}

	deliveryOrder := &models.DeliveryOrder{}
	for _, v := range getDeliveryOrderDetailResult.DeliveryOrderDetails {
		getDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
		go u.deliveryOrderRepository.GetByID(v.DeliveryOrderID, false, ctx, getDeliveryOrderResultChan)
		getDeliveryOrderResult := <-getDeliveryOrderResultChan

		if getDeliveryOrderResult.Error != nil {
			return &models.DeliveryOrderDetail{}, getDeliveryOrderDetailResult.ErrorLog
		}
		deliveryOrder = getDeliveryOrderResult.DeliveryOrder
	}

	deliveryOrderLog := &models.DeliveryOrderLog{
		RequestID: request.RequestID,
		DoCode:    deliveryOrder.DoCode,
		Data:      deliveryOrderDetail,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		Action:    constants.LOG_ACTION_MONGO_UPDATE,
		CreatedAt: &now,
	}

	createDeliveryOrderLogResultChan := make(chan *models.DeliveryOrderLogChan)
	go u.deliveryOrderLogRepository.Insert(deliveryOrderLog, ctx, createDeliveryOrderLogResultChan)
	createDeliveryOrderLogResult := <-createDeliveryOrderLogResultChan

	if createDeliveryOrderLogResult.Error != nil {
		errorLogData := helper.WriteLog(createDeliveryOrderLogResult.Error, http.StatusInternalServerError, nil)
		return &models.DeliveryOrderDetail{}, errorLogData
	}

	keyKafka := []byte(deliveryOrder.DoCode)
	messageKafka, _ := json.Marshal(deliveryOrderDetail)
	err := u.kafkaClient.WriteToTopic(constants.UPDATE_DELIVERY_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.DeliveryOrderDetail{}, errorLogData
	}

	return deliveryOrderDetail, nil
}

func (u *deliveryOrderUseCase) UpdateDoDetailByDeliveryOrderID(deliveryOrderID int, request []*models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderDetails, *model.ErrorLog) {
	now := time.Now()

	getDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailsChan)
	go u.deliveryOrderDetailRepository.GetByDeliveryOrderID(deliveryOrderID, false, ctx, getDeliveryOrderDetailResultChan)
	getDeliveryOrderDetailResult := <-getDeliveryOrderDetailResultChan

	if getDeliveryOrderDetailResult.Error != nil {
		return &models.DeliveryOrderDetails{}, getDeliveryOrderDetailResult.ErrorLog
	}

	deliveryOrderDetails := &models.DeliveryOrderDetails{}
	for _, v := range request {
		checkDeliveryOrderDetailIDResultChan := make(chan *models.DeliveryOrderDetailsChan)
		go u.deliveryOrderDetailRepository.GetByID(v.ID, false, ctx, checkDeliveryOrderDetailIDResultChan)
		checkDeliveryOrderDetailIDResult := <-checkDeliveryOrderDetailIDResultChan

		if checkDeliveryOrderDetailIDResult.Error != nil {
			return &models.DeliveryOrderDetails{}, checkDeliveryOrderDetailIDResult.ErrorLog
		}

		deliveryOrderDetail := &models.DeliveryOrderDetail{
			ID:                v.ID,
			Qty:               v.Qty,
			Note:              models.NullString{NullString: sql.NullString{String: v.Note, Valid: true}},
			IsDoneSyncToEs:    "0",
			UpdatedAt:         &now,
			StartDateSyncToEs: &now,
			EndDateSyncToEs:   &now,
		}
		deliveryOrderDetails.DeliveryOrderDetails = append(deliveryOrderDetails.DeliveryOrderDetails, deliveryOrderDetail)
	}

	deliveryOrderDetailss := &models.DeliveryOrderDetails{}
	for i, x := range getDeliveryOrderDetailResult.DeliveryOrderDetails {
		updateDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
		go u.deliveryOrderDetailRepository.UpdateByID(x.ID, deliveryOrderDetails.DeliveryOrderDetails[i], sqlTransaction, ctx, updateDeliveryOrderDetailResultChan)
		updateDeliveryOrderDetailResult := <-updateDeliveryOrderDetailResultChan

		if updateDeliveryOrderDetailResult.Error != nil {
			return &models.DeliveryOrderDetails{}, updateDeliveryOrderDetailResult.ErrorLog
		}
		updateDeliveryOrderDetailResult.DeliveryOrderDetail.ID = x.ID
		deliveryOrderDetailss.DeliveryOrderDetails = append(deliveryOrderDetailss.DeliveryOrderDetails, updateDeliveryOrderDetailResult.DeliveryOrderDetail)
	}

	getDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderRepository.GetByID(deliveryOrderID, false, ctx, getDeliveryOrderResultChan)
	getDeliveryOrderResult := <-getDeliveryOrderResultChan

	if getDeliveryOrderResult.Error != nil {
		return &models.DeliveryOrderDetails{}, getDeliveryOrderDetailResult.ErrorLog
	}

	requestID := strconv.Itoa(deliveryOrderID)
	deliveryOrderLog := &models.DeliveryOrderLog{
		RequestID: requestID,
		DoCode:    getDeliveryOrderResult.DeliveryOrder.DoCode,
		Data:      deliveryOrderDetailss,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		Action:    constants.LOG_ACTION_MONGO_UPDATE,
		CreatedAt: &now,
	}

	createDeliveryOrderLogResultChan := make(chan *models.DeliveryOrderLogChan)
	go u.deliveryOrderLogRepository.Insert(deliveryOrderLog, ctx, createDeliveryOrderLogResultChan)
	createDeliveryOrderLogResult := <-createDeliveryOrderLogResultChan

	if createDeliveryOrderLogResult.Error != nil {
		errorLogData := helper.WriteLog(createDeliveryOrderLogResult.Error, http.StatusInternalServerError, nil)
		return &models.DeliveryOrderDetails{}, errorLogData
	}

	keyKafka := []byte(getDeliveryOrderResult.DeliveryOrder.DoCode)
	messageKafka, _ := json.Marshal(deliveryOrderDetails)
	err := u.kafkaClient.WriteToTopic(constants.UPDATE_DELIVERY_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.DeliveryOrderDetails{}, errorLogData
	}

	return deliveryOrderDetailss, nil
}

func (u *deliveryOrderUseCase) Get(request *models.DeliveryOrderRequest) (*models.DeliveryOrdersOpenSearchResponse, *model.ErrorLog) {
	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
	go u.deliveryOrderOpenSearchRepository.Get(request, getDeliveryOrdersResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrdersResultChan

	if getDeliveryOrdersResult.Error != nil {
		return &models.DeliveryOrdersOpenSearchResponse{}, getDeliveryOrdersResult.ErrorLog
	}

	deliveryOrderResult := []*models.DeliveryOrderOpenSearchResponse{}
	for _, v := range getDeliveryOrdersResult.DeliveryOrders {
		deliveryOrder := models.DeliveryOrderOpenSearchResponse{
			ID:            v.ID,
			SalesOrderID:  v.SalesOrderID,
			WarehouseID:   v.WarehouseID,
			OrderSourceID: v.OrderSourceID,
			AgentName:     v.AgentName,
			AgentID:       v.AgentID,
			StoreID:       v.StoreID,
			DoCode:        v.DoCode,
			DoDate:        v.DoDate,
			DoRefCode:     v.DoRefCode,
			DoRefDate:     v.DoRefDate,
			DriverName:    v.DriverName,
			PlatNumber:    v.PlatNumber,
			Note:          v.Note,
		}
		deliveryOrderResult = append(deliveryOrderResult, &deliveryOrder)
		deliveryOrderDetails := []*models.DeliveryOrderDetailOpenSearchDetailResponse{}
		for _, x := range v.DeliveryOrderDetails {
			deliveryOrderDetail := models.DeliveryOrderDetailOpenSearchDetailResponse{
				SoDetailID: x.SoDetailID,
				Qty:        x.Qty,
			}
			deliveryOrderDetails = append(deliveryOrderDetails, &deliveryOrderDetail)
		}
		deliveryOrder.DeliveryOrderDetail = deliveryOrderDetails
	}

	deliveryOrders := &models.DeliveryOrdersOpenSearchResponse{
		DeliveryOrders: deliveryOrderResult,
		Total:          getDeliveryOrdersResult.Total,
	}

	return deliveryOrders, &model.ErrorLog{}
}

func (u *deliveryOrderUseCase) GetByID(request *models.DeliveryOrderRequest, ctx context.Context) (*models.DeliveryOrder, *model.ErrorLog) {
	getDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.GetByID(request, getDeliveryOrderResultChan)
	getDeliveryOrderResult := <-getDeliveryOrderResultChan

	if getDeliveryOrderResult.Error != nil {
		return &models.DeliveryOrder{}, getDeliveryOrderResult.ErrorLog
	}

	return getDeliveryOrderResult.DeliveryOrder, &model.ErrorLog{}
}

func (u *deliveryOrderUseCase) GetByIDWithDetail(request *models.DeliveryOrderRequest, ctx context.Context) (*models.DeliveryOrderOpenSearchResponse, *model.ErrorLog) {
	getDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.GetByID(request, getDeliveryOrderResultChan)
	getDeliveryOrderResult := <-getDeliveryOrderResultChan

	if getDeliveryOrderResult.Error != nil {
		return &models.DeliveryOrderOpenSearchResponse{}, getDeliveryOrderResult.ErrorLog
	}

	deliveryOrderResult := models.DeliveryOrderOpenSearchResponse{}
	deliveryOrderResult.DeliveryOrderOpenSearchResponseMap(getDeliveryOrderResult.DeliveryOrder)

	deliveryOrderDetails := []*models.DeliveryOrderDetailOpenSearchDetailResponse{}
	for _, x := range getDeliveryOrderResult.DeliveryOrder.DeliveryOrderDetails {
		deliveryOrderDetail := models.DeliveryOrderDetailOpenSearchDetailResponse{}
		deliveryOrderDetail.DeliveryOrderDetailOpenSearchResponseMap(x)

		deliveryOrderDetails = append(deliveryOrderDetails, &deliveryOrderDetail)
	}
	deliveryOrderResult.DeliveryOrderDetail = deliveryOrderDetails

	return &deliveryOrderResult, &model.ErrorLog{}
}

func (u *deliveryOrderUseCase) GetByAgentID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog) {
	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
	go u.deliveryOrderOpenSearchRepository.GetByAgentID(request, getDeliveryOrdersResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrdersResultChan

	if getDeliveryOrdersResult.Error != nil {
		return &models.DeliveryOrders{}, getDeliveryOrdersResult.ErrorLog
	}

	deliveryOrders := &models.DeliveryOrders{
		DeliveryOrders: getDeliveryOrdersResult.DeliveryOrders,
		Total:          getDeliveryOrdersResult.Total,
	}

	return deliveryOrders, &model.ErrorLog{}
}

func (u *deliveryOrderUseCase) GetByStoreID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog) {
	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
	go u.deliveryOrderOpenSearchRepository.GetByStoreID(request, getDeliveryOrdersResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrdersResultChan

	if getDeliveryOrdersResult.Error != nil {
		return &models.DeliveryOrders{}, getDeliveryOrdersResult.ErrorLog
	}

	deliveryOrders := &models.DeliveryOrders{
		DeliveryOrders: getDeliveryOrdersResult.DeliveryOrders,
		Total:          getDeliveryOrdersResult.Total,
	}

	return deliveryOrders, &model.ErrorLog{}
}

func (u *deliveryOrderUseCase) GetBySalesmanID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog) {
	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
	go u.deliveryOrderOpenSearchRepository.GetBySalesmanID(request, getDeliveryOrdersResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrdersResultChan

	if getDeliveryOrdersResult.Error != nil {
		return &models.DeliveryOrders{}, getDeliveryOrdersResult.ErrorLog
	}

	deliveryOrders := &models.DeliveryOrders{
		DeliveryOrders: getDeliveryOrdersResult.DeliveryOrders,
		Total:          getDeliveryOrdersResult.Total,
	}

	return deliveryOrders, &model.ErrorLog{}
}

func (u *deliveryOrderUseCase) GetBySalesmansID(request *models.DeliveryOrderRequest) (*models.DeliveryOrdersOpenSearchResponses, *model.ErrorLog) {
	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
	go u.deliveryOrderOpenSearchRepository.GetBySalesmansID(request, getDeliveryOrdersResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrdersResultChan

	if getDeliveryOrdersResult.Error != nil {
		return &models.DeliveryOrdersOpenSearchResponses{}, getDeliveryOrdersResult.ErrorLog
	}

	deliveryOrderResult := []*models.DeliveryOrderOpenSearchResponses{}
	for _, v := range getDeliveryOrdersResult.DeliveryOrders {
		deliveryOrder := models.DeliveryOrderOpenSearchResponses{
			ID:                    v.ID,
			SoCode:                v.SalesOrder.SoCode,
			SoDate:                v.SalesOrder.SoDate,
			WarehouseName:         models.NullString{NullString: sql.NullString{String: v.Warehouse.Name, Valid: true}},
			WarehouseCode:         models.NullString{NullString: sql.NullString{String: v.Warehouse.Code, Valid: true}},
			WarehouseProvinceName: models.NullString{NullString: sql.NullString{String: v.Warehouse.ProvinceName.String, Valid: true}},
			WarehouseCityName:     models.NullString{NullString: sql.NullString{String: v.Warehouse.CityName.String, Valid: true}},
			WarehouseDistrictName: models.NullString{NullString: sql.NullString{String: v.Warehouse.DistrictName.String, Valid: true}},
			WarehouseVillageName:  v.Warehouse.VillageName,
			DriverName:            models.NullString{NullString: sql.NullString{String: v.DriverName.String, Valid: true}},
			PlatNumber:            v.PlatNumber,
			AgentName:             models.NullString{NullString: sql.NullString{String: v.Agent.Name, Valid: true}},
			AgentEmail:            models.NullString{NullString: sql.NullString{String: v.Agent.Email.String, Valid: true}},
			AgentProvinceName:     models.NullString{NullString: sql.NullString{String: v.Agent.ProvinceName.String, Valid: true}},
			AgentCityName:         models.NullString{NullString: sql.NullString{String: v.Agent.CityName.String, Valid: true}},
			AgentDistrictName:     models.NullString{NullString: sql.NullString{String: v.Agent.DistrictName.String, Valid: true}},
			AgentVillageName:      models.NullString{NullString: sql.NullString{String: v.Agent.VillageName.String, Valid: true}},
			AgentAddress:          models.NullString{NullString: sql.NullString{String: v.Agent.Address.String, Valid: true}},
			AgentPhone:            models.NullString{NullString: sql.NullString{String: v.Agent.Phone.String, Valid: true}},
			AgentMainMobilePhone:  models.NullString{NullString: sql.NullString{String: v.Agent.MainMobilePhone.String, Valid: true}},
			StoreName:             models.NullString{NullString: sql.NullString{String: v.Store.Name.String, Valid: true}},
			StoreCode:             models.NullString{NullString: sql.NullString{String: v.Store.StoreCode.String, Valid: true}},
			StoreEmail:            models.NullString{NullString: sql.NullString{String: v.Store.Email.String, Valid: true}},
			StoreProvinceName:     models.NullString{NullString: sql.NullString{String: v.Store.ProvinceName.String, Valid: true}},
			StoreCityName:         models.NullString{NullString: sql.NullString{String: v.Store.CityName.String, Valid: true}},
			StoreDistrictName:     models.NullString{NullString: sql.NullString{String: v.Store.DistrictName.String, Valid: true}},
			StoreVillageName:      models.NullString{NullString: sql.NullString{String: v.Store.VillageName.String, Valid: true}},
			StoreAddress:          models.NullString{NullString: sql.NullString{String: v.Store.Address.String, Valid: true}},
			StorePhone:            models.NullString{NullString: sql.NullString{String: v.Store.Phone.String, Valid: true}},
			StoreMainMobilePhone:  models.NullString{NullString: sql.NullString{String: v.Store.MainMobilePhone.String, Valid: true}},
			BrandName:             v.SalesOrder.BrandName,
			UserFirstName:         models.NullString{NullString: sql.NullString{String: v.SalesOrder.UserFirstName.String, Valid: true}},
			UserLastName:          models.NullString{NullString: sql.NullString{String: v.SalesOrder.UserLastName.String, Valid: true}},
			UserEmail:             models.NullString{NullString: sql.NullString{String: v.SalesOrder.UserEmail.String, Valid: true}},
			OrderSourceName:       v.OrderSource.SourceName,
			OrderStatusName:       v.OrderStatus.Name,
			DoCode:                v.DoCode,
			DoDate:                v.DoDate,
			DoRefCode:             v.DoRefCode,
			DoRefDate:             v.DoRefDate,
			Note:                  v.Note,
		}
		deliveryOrderResult = append(deliveryOrderResult, &deliveryOrder)
		var deliveryOrderDetails []*models.DeliveryOrderDetailOpenSearchResponse
		for _, x := range v.DeliveryOrderDetails {
			deliveryOrderDetail := models.DeliveryOrderDetailOpenSearchResponse{
				ID:              x.ID,
				DeliveryOrderID: x.DeliveryOrderID,
				ProductID:       x.ProductID,
				Product: &models.ProductOpenSearchDeliveryOrderResponse{
					ID:                    x.Product.ID,
					ProductSku:            x.Product.Sku,
					AliasSku:              x.Product.AliasSku,
					ProductName:           x.Product.ProductName,
					Description:           x.Product.Description,
					UnitMeasurementSmall:  x.Product.UnitMeasurementSmall,
					UnitMeasurementMedium: x.Product.UnitMeasurementMedium,
					UnitMeasurementBig:    x.Product.UnitMeasurementBig,
					Ukuran:                x.Product.Ukuran,
					NettWeightUm:          x.Product.NettWeightUm,
					Currency:              x.Product.Currency,
					DataType:              x.Product.DataType,
					Image:                 x.Product.Image,
				},
				UomID:           x.UomID,
				UomName:         x.Uom.Name.String,
				UomCode:         x.Uom.Code.String,
				OrderStatusID:   x.OrderStatusID,
				OrderStatusName: x.OrderStatusName,
				DoDetailCode:    x.DoDetailCode,
				Qty:             x.Qty,
				SentQty:         x.SentQty,
				ResidualQty:     x.ResidualQty,
				Price:           x.Price,
				Note:            x.Note,
				CreatedAt:       x.CreatedAt,
			}
			deliveryOrderDetails = append(deliveryOrderDetails, &deliveryOrderDetail)
		}
		deliveryOrder.DeliveryOrderDetails = deliveryOrderDetails
		deliveryOrder.CreatedAt = v.CreatedAt
		deliveryOrder.UpdatedAt = v.UpdatedAt
	}

	deliveryOrders := &models.DeliveryOrdersOpenSearchResponses{
		DeliveryOrders: deliveryOrderResult,
		Total:          getDeliveryOrdersResult.Total,
	}

	return deliveryOrders, &model.ErrorLog{}
}

func (u *deliveryOrderUseCase) GetByOrderStatusID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog) {
	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
	go u.deliveryOrderOpenSearchRepository.GetBySalesmanID(request, getDeliveryOrdersResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrdersResultChan

	if getDeliveryOrdersResult.Error != nil {
		return &models.DeliveryOrders{}, getDeliveryOrdersResult.ErrorLog
	}

	deliveryOrders := &models.DeliveryOrders{
		DeliveryOrders: getDeliveryOrdersResult.DeliveryOrders,
		Total:          getDeliveryOrdersResult.Total,
	}

	return deliveryOrders, &model.ErrorLog{}
}

func (u *deliveryOrderUseCase) GetByOrderSourceID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog) {
	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
	go u.deliveryOrderOpenSearchRepository.GetBySalesmanID(request, getDeliveryOrdersResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrdersResultChan

	if getDeliveryOrdersResult.Error != nil {
		return &models.DeliveryOrders{}, getDeliveryOrdersResult.ErrorLog
	}

	deliveryOrders := &models.DeliveryOrders{
		DeliveryOrders: getDeliveryOrdersResult.DeliveryOrders,
		Total:          getDeliveryOrdersResult.Total,
	}

	return deliveryOrders, &model.ErrorLog{}
}

func (u *deliveryOrderUseCase) SyncToOpenSearchFromCreateEvent(deliveryOrder *models.DeliveryOrder, salesOrderUseCase SalesOrderUseCaseInterface, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog {
	now := time.Now()

	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(deliveryOrder.AgentID, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	deliveryOrder.Agent = getAgentResult.Agent
	deliveryOrder.AgentName = getAgentResult.Agent.Name

	getStoreResultChan := make(chan *models.StoreChan)
	go u.storeRepository.GetByID(deliveryOrder.StoreID, false, ctx, getStoreResultChan)
	getStoreResult := <-getStoreResultChan

	if getStoreResult.Error != nil {
		errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	deliveryOrder.Store = getStoreResult.Store

	for k, v := range deliveryOrder.DeliveryOrderDetails {
		getSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.GetByID(v.SoDetailID, false, ctx, getSalesOrderDetailResultChan)
		getSalesOrderDetailResult := <-getSalesOrderDetailResultChan

		if getSalesOrderDetailResult.Error != nil {
			errorLogData := helper.WriteLog(getSalesOrderDetailResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		getProductResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetByID(v.ProductID, false, ctx, getProductResultChan)
		getProductResult := <-getProductResultChan

		if getProductResult.Error != nil {
			errorLogData := helper.WriteLog(getProductResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		deliveryOrder.DeliveryOrderDetails[k].Product = getProductResult.Product

		getUomResultChan := make(chan *models.UomChan)
		go u.uomRepository.GetByID(v.UomID, false, ctx, getUomResultChan)
		getUomResult := <-getUomResultChan

		if getUomResult.Error != nil {
			errorLogData := helper.WriteLog(getUomResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		deliveryOrder.DeliveryOrderDetails[k].Uom = getUomResult.Uom

		getBrandResultChan := make(chan *models.BrandChan)
		go u.brandRepository.GetByID(v.UomID, false, ctx, getBrandResultChan)
		getBrandResult := <-getBrandResultChan

		if getBrandResult.Error != nil {
			errorLogData := helper.WriteLog(getBrandResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		deliveryOrder.DeliveryOrderDetails[k].Brand = getBrandResult.Brand

		residualQty := (getSalesOrderDetailResult.SalesOrderDetail.Qty - getSalesOrderDetailResult.SalesOrderDetail.SentQty) - v.Qty
		getSalesOrderDetailResult.SalesOrderDetail.ResidualQty = residualQty
		getSalesOrderDetailResult.SalesOrderDetail.SentQty = v.Qty

		statusName := "partial"

		if residualQty == 0 {
			statusName = "closed"
		}

		getStatusSalesOrderDetailResultChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByNameAndType(statusName, "sales_order_detail", false, ctx, getStatusSalesOrderDetailResultChan)
		getStatusSalesOrderDetailResult := <-getStatusSalesOrderDetailResultChan

		if getStatusSalesOrderDetailResult.Error != nil {
			errorLogData := helper.WriteLog(getStatusSalesOrderDetailResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		salesOrderDetailDataUpdate := &models.SalesOrderDetail{
			UpdatedAt:       &now,
			OrderStatusID:   getStatusSalesOrderDetailResult.OrderStatus.ID,
			ResidualQty:     residualQty,
			SentQty:         v.Qty,
			EndDateSyncToEs: &now,
		}

		updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.UpdateByID(getSalesOrderDetailResult.SalesOrderDetail.ID, salesOrderDetailDataUpdate, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
		updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

		if updateSalesOrderDetailResult.Error != nil {
			errorLogData := helper.WriteLog(updateSalesOrderDetailResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		removeCacheSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.RemoveCacheByID(getSalesOrderDetailResult.SalesOrderDetail.ID, ctx, removeCacheSalesOrderDetailResultChan)
		removeCacheSalesOrderDetailResult := <-removeCacheSalesOrderDetailResultChan

		if removeCacheSalesOrderDetailResult.Error != nil {
			errorLogData := helper.WriteLog(removeCacheSalesOrderDetailResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		deliveryOrder.DeliveryOrderDetails[k].EndDateSyncToEs = &now
		deliveryOrder.DeliveryOrderDetails[k].IsDoneSyncToEs = "1"
	}

	getSalesOrderDetailsResultChan := make(chan *models.SalesOrderDetailsChan)
	go u.salesOrderDetailRepository.GetBySalesOrderID(deliveryOrder.SalesOrderID, false, ctx, getSalesOrderDetailsResultChan)
	getSalesOrderDetailsResult := <-getSalesOrderDetailsResultChan

	if getSalesOrderDetailsResult.Error != nil {
		errorLogData := helper.WriteLog(getSalesOrderDetailsResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	deliveryOrder.SalesOrder.SalesOrderDetails = getSalesOrderDetailsResult.SalesOrderDetails

	countPartialStatusOnSalesOrderDetail := 0
	countClosedStatusOnSalesOrderDetail := 0
	for _, v := range getSalesOrderDetailsResult.SalesOrderDetails {
		getStatusSalesOrderDetailResultChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByID(v.OrderStatusID, false, ctx, getStatusSalesOrderDetailResultChan)
		getStatusSalesOrderDetailResult := <-getStatusSalesOrderDetailResultChan

		if getStatusSalesOrderDetailResult.Error != nil {
			errorLogData := helper.WriteLog(getStatusSalesOrderDetailResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		if getStatusSalesOrderDetailResult.OrderStatus.Name == "partial" {
			countPartialStatusOnSalesOrderDetail++
		} else {
			countClosedStatusOnSalesOrderDetail++
		}
	}

	salesOrderStatus := "partial"
	if countPartialStatusOnSalesOrderDetail == 0 {
		salesOrderStatus = "closed"
	}

	getStatusSalesOrderResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType(salesOrderStatus, "sales_order", false, ctx, getStatusSalesOrderResultChan)
	getStatusSalesOrderResult := <-getStatusSalesOrderResultChan

	if getStatusSalesOrderResult.Error != nil {
		errorLogData := helper.WriteLog(getStatusSalesOrderResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	salesOrderUpdateData := &models.SalesOrder{
		UpdatedAt:     &now,
		OrderStatusID: getStatusSalesOrderResult.OrderStatus.ID,
	}

	updateSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.UpdateByID(deliveryOrder.SalesOrderID, salesOrderUpdateData, sqlTransaction, ctx, updateSalesOrderResultChan)
	updateSalesOrderResult := <-updateSalesOrderResultChan

	if updateSalesOrderResult.Error != nil {
		errorLogData := helper.WriteLog(updateSalesOrderResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	removeCacheSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.RemoveCacheByID(deliveryOrder.SalesOrderID, ctx, removeCacheSalesOrderResultChan)
	removeCacheSalesOrderResult := <-removeCacheSalesOrderResultChan

	if removeCacheSalesOrderResult.Error != nil {
		errorLogData := helper.WriteLog(removeCacheSalesOrderResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	deliveryOrder.IsDoneSyncToEs = "1"
	deliveryOrder.EndDateSyncToEs = &now
	deliveryOrder.UpdatedAt = &now
	deliveryOrder.EndCreatedDate = &now

	deliveryOrderUpdateData := &models.DeliveryOrder{
		UpdatedAt:       &now,
		IsDoneSyncToEs:  "1",
		EndDateSyncToEs: &now,
	}

	updateDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderRepository.UpdateByID(deliveryOrder.ID, deliveryOrderUpdateData, sqlTransaction, ctx, updateDeliveryOrderResultChan)
	updateDeliveryOrderResult := <-updateDeliveryOrderResultChan

	if updateDeliveryOrderResult.Error != nil {
		errorLogData := helper.WriteLog(updateDeliveryOrderResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.Create(deliveryOrder, createDeliveryOrderResultChan)
	createDeliveryOrderResult := <-createDeliveryOrderResultChan

	if createDeliveryOrderResult.Error != nil {
		return createDeliveryOrderResult.ErrorLog
	}

	return &model.ErrorLog{}
}

func (u *deliveryOrderUseCase) SyncToOpenSearchFromUpdateEvent(deliveryOrder *models.DeliveryOrder, ctx context.Context) *model.ErrorLog {
	now := time.Now()

	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(deliveryOrder.AgentID, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	deliveryOrder.Agent = getAgentResult.Agent

	getStoreResultChan := make(chan *models.StoreChan)
	go u.storeRepository.GetByID(deliveryOrder.StoreID, false, ctx, getStoreResultChan)
	getStoreResult := <-getStoreResultChan

	if getStoreResult.Error != nil {
		errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	deliveryOrder.Store = getStoreResult.Store

	for k, v := range deliveryOrder.DeliveryOrderDetails {
		getSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.GetByID(v.SoDetailID, false, ctx, getSalesOrderDetailResultChan)
		getSalesOrderDetailResult := <-getSalesOrderDetailResultChan

		if getSalesOrderDetailResult.Error != nil {
			errorLogData := helper.WriteLog(getSalesOrderDetailResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		getProductResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetByID(v.ProductID, false, ctx, getProductResultChan)
		getProductResult := <-getProductResultChan

		if getProductResult.Error != nil {
			errorLogData := helper.WriteLog(getProductResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		deliveryOrder.DeliveryOrderDetails[k].Product = getProductResult.Product

		getUomResultChan := make(chan *models.UomChan)
		go u.uomRepository.GetByID(v.UomID, false, ctx, getUomResultChan)
		getUomResult := <-getUomResultChan

		if getUomResult.Error != nil {
			errorLogData := helper.WriteLog(getUomResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		deliveryOrder.DeliveryOrderDetails[k].Uom = getUomResult.Uom

		getBrandResultChan := make(chan *models.BrandChan)
		go u.brandRepository.GetByID(v.UomID, false, ctx, getBrandResultChan)
		getBrandResult := <-getBrandResultChan

		if getBrandResult.Error != nil {
			errorLogData := helper.WriteLog(getBrandResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		deliveryOrder.DeliveryOrderDetails[k].Brand = getBrandResult.Brand
	}

	deliveryOrder.UpdatedAt = &now

	createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.Create(deliveryOrder, createDeliveryOrderResultChan)
	createDeliveryOrderResult := <-createDeliveryOrderResultChan

	if createDeliveryOrderResult.Error != nil {
		fmt.Println(createDeliveryOrderResult.ErrorLog.Err.Error())
		return createDeliveryOrderResult.ErrorLog
	}

	return &model.ErrorLog{}
}

func (u *deliveryOrderUseCase) SyncToOpenSearchFromDeleteEvent(deliveryOrderId *int, ctx context.Context) *model.ErrorLog {
	now := time.Now()
	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.GetByID(&models.DeliveryOrderRequest{ID: *deliveryOrderId}, getDeliveryOrdersResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrdersResultChan

	if getDeliveryOrdersResult.Error != nil {
		if !strings.Contains(getDeliveryOrdersResult.Error.Error(), "not found") {
			errorLogData := helper.WriteLog(getDeliveryOrdersResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}
	}
	deliveryOrder := getDeliveryOrdersResult.DeliveryOrder

	for k := range deliveryOrder.DeliveryOrderDetails {
		deliveryOrder.DeliveryOrderDetails[k].DeletedAt = &now
		deliveryOrder.DeliveryOrderDetails[k].UpdatedAt = &now
	}

	deliveryOrder.DeletedAt = &now
	deliveryOrder.UpdatedAt = &now

	createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.Create(deliveryOrder, createDeliveryOrderResultChan)
	deleteDeliveryOrderResult := <-createDeliveryOrderResultChan

	if deleteDeliveryOrderResult.Error != nil {
		fmt.Println(deleteDeliveryOrderResult.ErrorLog.Err.Error())
		return deleteDeliveryOrderResult.ErrorLog
	}

	salesOrderRequest := &models.SalesOrderRequest{
		ID:            deliveryOrder.SalesOrderID,
		OrderSourceID: deliveryOrder.OrderSourceID,
	}

	getSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderOpenSearchRepository.GetByID(salesOrderRequest, getSalesOrderResultChan)
	getSalesOrderResult := <-getSalesOrderResultChan

	if getSalesOrderResult.ErrorLog != nil {
		deleteDeliveryOrderResult.ErrorLog = getSalesOrderResult.ErrorLog
		fmt.Println(deleteDeliveryOrderResult.ErrorLog.Err.Error())
		return deleteDeliveryOrderResult.ErrorLog
	}

	salesOrderWithDetail := getSalesOrderResult.SalesOrder
	for k, v := range salesOrderWithDetail.SalesOrderDetails {
		v.SentQty -= deliveryOrder.DeliveryOrderDetails[k].Qty
		v.ResidualQty += deliveryOrder.DeliveryOrderDetails[k].Qty
	}
	deleteDeliveryOrderResult.DeliveryOrder.SalesOrder = salesOrderWithDetail
	deleteDeliveryOrderResult.ErrorLog = u.SalesOrderOpenSearchUseCase.SyncToOpenSearchFromUpdateEvent(salesOrderWithDetail, ctx)

	if deleteDeliveryOrderResult.ErrorLog != nil {
		fmt.Println(deleteDeliveryOrderResult.ErrorLog.Err.Error())
		return deleteDeliveryOrderResult.ErrorLog
	}

	deleteDeliveryOrderResult.ErrorLog = u.SyncToOpenSearchFromUpdateEvent(deleteDeliveryOrderResult.DeliveryOrder, ctx)

	if deleteDeliveryOrderResult.ErrorLog != nil {
		fmt.Println(deleteDeliveryOrderResult.ErrorLog.Err.Error())
		return deleteDeliveryOrderResult.ErrorLog
	}

	return &model.ErrorLog{}
}

func (u deliveryOrderUseCase) DeleteByID(id int) *model.ErrorLog {
	now := time.Now()
	getDeliveryOrderByIDResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderRepository.GetByID(id, false, u.ctx, getDeliveryOrderByIDResultChan)
	getDeliveryOrderByIDResult := <-getDeliveryOrderByIDResultChan

	if getDeliveryOrderByIDResult.Error != nil {
		return getDeliveryOrderByIDResult.ErrorLog
	}

	getDeliveryOrderDetailByIDResultChan := make(chan *models.DeliveryOrderDetailsChan)
	go u.deliveryOrderDetailRepository.GetByDeliveryOrderID(id, false, u.ctx, getDeliveryOrderDetailByIDResultChan)
	getDeliveryOrderDetailsByIDResult := <-getDeliveryOrderDetailByIDResultChan

	if getDeliveryOrderDetailsByIDResult.Error != nil {
		return getDeliveryOrderDetailsByIDResult.ErrorLog
	}

	getDeliveryOrderByIDResult.DeliveryOrder.DeliveryOrderDetails = getDeliveryOrderDetailsByIDResult.DeliveryOrderDetails

	getSalesOrderByIDResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.GetByID(getDeliveryOrderByIDResult.DeliveryOrder.SalesOrderID, false, u.ctx, getSalesOrderByIDResultChan)
	getSalesOrderByIDResult := <-getSalesOrderByIDResultChan

	if getSalesOrderByIDResult.Error != nil {
		return getSalesOrderByIDResult.ErrorLog
	}
	totalSentQty := 0
	sqlTransaction, err := u.db.Begin()
	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}
	for _, v := range getDeliveryOrderDetailsByIDResult.DeliveryOrderDetails {

		getSalesOrderDetailByIDResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.GetByID(getSalesOrderByIDResult.SalesOrder.ID, false, u.ctx, getSalesOrderDetailByIDResultChan)
		getSalesOrderDetailsByIDResult := <-getSalesOrderDetailByIDResultChan

		if getSalesOrderDetailsByIDResult.Error != nil {
			return getSalesOrderDetailsByIDResult.ErrorLog
		}

		getSalesOrderDetailsByIDResult.SalesOrderDetail.SentQty -= v.Qty
		getSalesOrderDetailsByIDResult.SalesOrderDetail.ResidualQty += v.Qty
		getSalesOrderDetailsByIDResult.SalesOrderDetail.UpdatedAt = &now

		totalSentQty += getSalesOrderDetailsByIDResult.SalesOrderDetail.SentQty

		deleteDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
		go u.deliveryOrderDetailRepository.DeleteByID(v, sqlTransaction, u.ctx, deleteDeliveryOrderDetailResultChan)
		deleteDeliveryOrderDetailResult := <-deleteDeliveryOrderDetailResultChan

		if deleteDeliveryOrderDetailResult.ErrorLog != nil {
			err = sqlTransaction.Rollback()
			if err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
				return errorLogData
			}
			return deleteDeliveryOrderDetailResult.ErrorLog
		}

		updateSalesOrderDetailChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.UpdateByID(v.SoDetailID, getSalesOrderDetailsByIDResult.SalesOrderDetail, sqlTransaction, u.ctx, updateSalesOrderDetailChan)
		updateSalesOrderDetailResult := <-updateSalesOrderDetailChan

		if updateSalesOrderDetailResult.ErrorLog != nil {
			err = sqlTransaction.Rollback()
			if err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
				return errorLogData
			}
			return updateSalesOrderDetailResult.ErrorLog
		}
	}

	deleteDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderRepository.DeleteByID(getDeliveryOrderByIDResult.DeliveryOrder, u.ctx, deleteDeliveryOrderResultChan)
	deleteDeliveryOrderResult := <-deleteDeliveryOrderResultChan

	if deleteDeliveryOrderResult.ErrorLog != nil {
		err = sqlTransaction.Rollback()
		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return errorLogData
		}
		return deleteDeliveryOrderResult.ErrorLog
	}
	fmt.Println(getSalesOrderByIDResult.SalesOrder.OrderStatusID)
	if totalSentQty > 0 {
		getSalesOrderByIDResult.SalesOrder.OrderStatusID = 7
	} else {
		getSalesOrderByIDResult.SalesOrder.OrderStatusID = 5
	}
	updateSalesOrderChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.UpdateByID(getSalesOrderByIDResult.SalesOrder.ID, getSalesOrderByIDResult.SalesOrder, sqlTransaction, u.ctx, updateSalesOrderChan)
	updateSalesOrderResult := <-updateSalesOrderChan

	if updateSalesOrderResult.ErrorLog != nil {
		err = sqlTransaction.Rollback()
		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return errorLogData
		}
		return updateSalesOrderResult.ErrorLog
	}

	err = sqlTransaction.Commit()

	if err != nil {
		sqlTransaction.Rollback()
		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return errorLogData
		}
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}
