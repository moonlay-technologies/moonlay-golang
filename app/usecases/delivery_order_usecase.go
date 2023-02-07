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
	"time"

	"github.com/bxcodec/dbresolver"
)

type DeliveryOrderUseCaseInterface interface {
	Create(request *models.DeliveryOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrder, *model.ErrorLog)
	UpdateByID(ID int, request *models.DeliveryOrderUpdateByIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrder, *model.ErrorLog)
	UpdateDODetailByID(ID int, request *models.DeliveryOrderDetailUpdateByIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderDetail, *model.ErrorLog)
	UpdateDoDetailByDeliveryOrderID(deliveryOrderID int, request []*models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderDetails, *model.ErrorLog)
	Get(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)
	GetByID(request *models.DeliveryOrderRequest, ctx context.Context) (*models.DeliveryOrder, *model.ErrorLog)
	GetByAgentID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)
	GetByStoreID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)
	GetBySalesmanID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)
	GetByOrderStatusID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)
	GetByOrderSourceID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)
	SyncToOpenSearchFromCreateEvent(deliveryOrder *models.DeliveryOrder, salesOrderUseCase SalesOrderUseCaseInterface, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromUpdateEvent(deliveryOrder *models.DeliveryOrder, ctx context.Context) *model.ErrorLog
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
	kafkaClient                       kafkadbo.KafkaClientInterface
	db                                dbresolver.DB
	ctx                               context.Context
}

func InitDeliveryOrderUseCaseInterface(deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, deliveryOrderDetailRepository repositories.DeliveryOrderDetailRepositoryInterface, salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, orderSourceRepository repositories.OrderSourceRepositoryInterface, warehouseRepository repositories.WarehouseRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, uomRepository repositories.UomRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, userRepository repositories.UserRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, deliveryOrderLogRepository mongoRepositories.DeliveryOrderLogRepositoryInterface, deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, salesOrderUseCase SalesOrderUseCaseInterface, kafkaClient kafkadbo.KafkaClientInterface, db dbresolver.DB, ctx context.Context) DeliveryOrderUseCaseInterface {
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
		kafkaClient:                       kafkaClient,
		db:                                db,
		ctx:                               ctx,
	}
}

func (u *deliveryOrderUseCase) Create(request *models.DeliveryOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrder, *model.ErrorLog) {
	now := time.Now()
	unixTimestamp := now.Unix()
	unixTimestampInt := int(unixTimestamp)
	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType("open", "delivery_order", false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return &models.DeliveryOrder{}, getOrderStatusResult.ErrorLog
	}

	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(request.OrderSourceID, false, ctx, getOrderSourceResultChan)
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

	getSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.GetByID(request.SalesOrderID, false, ctx, getSalesOrderResultChan)
	getSalesOrderResult := <-getSalesOrderResultChan

	if getSalesOrderResult.Error != nil {
		return &models.DeliveryOrder{}, getSalesOrderResult.ErrorLog
	}

	getStoreResultChan := make(chan *models.StoreChan)
	go u.storeRepository.GetByID(request.StoreID, false, ctx, getStoreResultChan)
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
		return &models.DeliveryOrder{}, getSalesmanResult.ErrorLog
	}

	deliveryOrder := &models.DeliveryOrder{
		SalesOrder:            getSalesOrderResult.SalesOrder,
		SalesOrderID:          request.SalesOrderID,
		Salesman:              getSalesmanResult.Salesman,
		Warehouse:             getWarehouseResult.Warehouse,
		Store:                 getStoreResult.Store,
		StoreID:               request.StoreID,
		AgentID:               request.AgentID,
		WarehouseID:           request.WarehouseID,
		OrderStatus:           getOrderStatusResult.OrderStatus,
		OrderStatusID:         request.OrderStatusID,
		OrderSource:           getOrderSourceResult.OrderSource,
		OrderSourceID:         getOrderSourceResult.OrderSource.ID,
		DoCode:                request.DoCode,
		DoDate:                request.DoDate,
		DoRefCode:             models.NullString{NullString: sql.NullString{String: request.DoRefCode, Valid: true}},
		DoRefDate:             models.NullString{NullString: sql.NullString{String: request.DoRefDate, Valid: true}},
		DriverName:            models.NullString{NullString: sql.NullString{String: request.DriverName, Valid: true}},
		PlatNumber:            models.NullString{NullString: sql.NullString{String: request.PlatNumber, Valid: true}},
		WarehouseName:         getWarehouseResult.Warehouse.Name,
		WarehouseCode:         getWarehouseResult.Warehouse.Code,
		WarehouseProvinceName: getWarehouseResult.Warehouse.ProvinceName,
		WarehouseCityName:     getWarehouseResult.Warehouse.CityName,
		WarehouseDistrictName: getWarehouseResult.Warehouse.DistrictName,
		WarehouseVillageName:  getWarehouseResult.Warehouse.VillageName,
		IsDoneSyncToEs:        "0",
		Note:                  models.NullString{NullString: sql.NullString{String: request.Note, Valid: true}},
		StartDateSyncToEs:     &now,
		EndDateSyncToEs:       &now,
		StartCreatedDate:      &now,
		EndCreatedDate:        &now,
		CreatedBy:             request.SalesOrderID,
		LatestUpdatedBy:       unixTimestampInt,
		CreatedAt:             &now,
		UpdatedAt:             &now,
		DeletedAt:             nil,
	}

	createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderRepository.Insert(deliveryOrder, sqlTransaction, ctx, createDeliveryOrderResultChan)
	createDeliveryOrderResult := <-createDeliveryOrderResultChan

	if createDeliveryOrderResult.Error != nil {
		return &models.DeliveryOrder{}, createDeliveryOrderResult.ErrorLog
	}

	deliveryOrderDetails := []*models.DeliveryOrderDetail{}

	for _, v := range request.DeliveryOrderDetails {
		getSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.GetByID(v.SoDetailID, false, ctx, getSalesOrderDetailResultChan)
		getSalesOrderDetailResult := <-getSalesOrderDetailResultChan

		if getSalesOrderDetailResult.Error != nil {
			return &models.DeliveryOrder{}, getSalesOrderDetailResult.ErrorLog
		}

		getProductOrderResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetByID(v.ProductID, false, ctx, getProductOrderResultChan)
		getProductResult := <-getProductOrderResultChan

		if getProductResult.Error != nil {
			return &models.DeliveryOrder{}, getProductResult.ErrorLog
		}

		getUomResultChan := make(chan *models.UomChan)
		go u.uomRepository.GetByID(v.UomID, false, ctx, getUomResultChan)
		getUomResult := <-getUomResultChan

		if getUomResult.Error != nil {
			return &models.DeliveryOrder{}, getUomResult.ErrorLog
		}

		DoDetailCode, _ := helper.GenerateDODetailCode(int(createDeliveryOrderResult.ID), getSalesOrderResult.SalesOrder.AgentID, getSalesOrderDetailResult.SalesOrderDetail.ProductID, getSalesOrderDetailResult.SalesOrderDetail.UomID)

		deliveryOrderDetail := &models.DeliveryOrderDetail{
			DeliveryOrderID:   int(createDeliveryOrderResult.ID),
			SoDetailID:        v.SoDetailID,
			BrandID:           getSalesOrderResult.SalesOrder.BrandID,
			ProductID:         getSalesOrderDetailResult.SalesOrderDetail.ProductID,
			UomID:             getSalesOrderDetailResult.SalesOrderDetail.UomID,
			OrderStatusID:     getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID,
			DoDetailCode:      DoDetailCode,
			Qty:               v.Qty,
			ProductSKU:        getProductResult.Product.Sku.String,
			ProductName:       getProductResult.Product.ProductName.String,
			Note:              models.NullString{NullString: sql.NullString{String: v.Note, Valid: true}},
			Product:           getProductResult.Product,
			SoDetail:          getSalesOrderDetailResult.SalesOrderDetail,
			Uom:               getUomResult.Uom,
			IsDoneSyncToEs:    "0",
			StartDateSyncToEs: &now,
			EndDateSyncToEs:   &now,
			CreatedAt:         &now,
			UpdatedAt:         &now,
			DeletedAt:         nil,
		}

		createDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
		go u.deliveryOrderDetailRepository.Insert(deliveryOrderDetail, sqlTransaction, ctx, createDeliveryOrderDetailResultChan)
		createDeliveryOrderDetailResult := <-createDeliveryOrderDetailResultChan

		if createDeliveryOrderDetailResult.Error != nil {
			return &models.DeliveryOrder{}, createDeliveryOrderDetailResult.ErrorLog
		}

		deliveryOrderDetail.ID = int(createDeliveryOrderDetailResult.ID)
		deliveryOrderDetails = append(deliveryOrderDetails, deliveryOrderDetail)
	}

	deliveryOrder.DeliveryOrderDetails = deliveryOrderDetails

	deliveryOrderLog := &models.DeliveryOrderLog{
		RequestID: request.RequestID,
		DoCode:    request.DoCode,
		Data:      deliveryOrder,
		Status:    "0",
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
		Status:    "0",
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
		Status:    "0",
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
		Status:    "0",
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

func (u *deliveryOrderUseCase) Get(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog) {
	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
	go u.deliveryOrderOpenSearchRepository.Get(request, getDeliveryOrdersResultChan)
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

func (u *deliveryOrderUseCase) GetByID(request *models.DeliveryOrderRequest, ctx context.Context) (*models.DeliveryOrder, *model.ErrorLog) {
	getDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.GetByID(request, getDeliveryOrderResultChan)
	getDeliveryOrderResult := <-getDeliveryOrderResultChan

	if getDeliveryOrderResult.Error != nil {
		return &models.DeliveryOrder{}, getDeliveryOrderResult.ErrorLog
	}

	return getDeliveryOrderResult.DeliveryOrder, &model.ErrorLog{}
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
