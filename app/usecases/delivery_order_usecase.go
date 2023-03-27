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
	baseModel "order-service/global/utils/model"
	"time"

	"github.com/bxcodec/dbresolver"
	"github.com/google/uuid"
)

type DeliveryOrderUseCaseInterface interface {
	Create(request *models.DeliveryOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderStoreResponse, *model.ErrorLog)
	UpdateByID(ID int, request *models.DeliveryOrderUpdateByIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderUpdateByIDRequest, *model.ErrorLog)
	UpdateDODetailByID(ID int, request *models.DeliveryOrderDetailUpdateByIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderDetailUpdateByIDRequest, *model.ErrorLog)
	UpdateDoDetailByDeliveryOrderID(deliveryOrderID int, request []*models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderDetails, *model.ErrorLog)
	Get(request *models.DeliveryOrderRequest) (*models.DeliveryOrdersOpenSearchResponse, *model.ErrorLog)
	GetByID(request *models.DeliveryOrderRequest, ctx context.Context) (*models.DeliveryOrder, *model.ErrorLog)
	GetDetails(request *models.DeliveryOrderDetailOpenSearchRequest) (*models.DeliveryOrderDetailsOpenSearchResponses, *model.ErrorLog)
	GetDetailsByDoID(request *models.DeliveryOrderDetailRequest) (*models.DeliveryOrderOpenSearchResponse, *model.ErrorLog)
	GetDetailByID(doDetailID int, doID int) (*models.DeliveryOrderDetailsOpenSearchResponse, *model.ErrorLog)
	GetByIDWithDetail(request *models.DeliveryOrderRequest, ctx context.Context) (*models.DeliveryOrderOpenSearchResponse, *model.ErrorLog)
	GetByAgentID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)
	GetByStoreID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)
	GetBySalesmanID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)
	GetBySalesmansID(request *models.DeliveryOrderRequest) (*models.DeliveryOrdersOpenSearchResponses, *model.ErrorLog)
	GetByOrderStatusID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)
	GetByOrderSourceID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)
	GetSyncToKafkaHistories(request *models.DeliveryOrderEventLogRequest, ctx context.Context) ([]*models.DeliveryOrderEventLogResponse, *model.ErrorLog)
	GetDOJourneys(request *models.DeliveryOrderJourneysRequest, ctx context.Context) (*models.DeliveryOrderJourneysResponses, *model.ErrorLog)
	GetDOJourneysByDoID(doId int, ctx context.Context) (*models.DeliveryOrderJourneysResponses, *model.ErrorLog)
	DeleteByID(deliveryOrderId int, sqlTransaction *sql.Tx) *model.ErrorLog
	DeleteDetailByID(deliveryOrderDetailId int, sqlTransaction *sql.Tx) *model.ErrorLog
	DeleteDetailByDoID(deliveryOrderId int, sqlTransaction *sql.Tx) *model.ErrorLog
	Export(request *models.DeliveryOrderExportRequest) *model.ErrorLog
}

type deliveryOrderUseCase struct {
	deliveryOrderRepository                 repositories.DeliveryOrderRepositoryInterface
	deliveryOrderDetailRepository           repositories.DeliveryOrderDetailRepositoryInterface
	salesOrderRepository                    repositories.SalesOrderRepositoryInterface
	salesOrderDetailRepository              repositories.SalesOrderDetailRepositoryInterface
	orderStatusRepository                   repositories.OrderStatusRepositoryInterface
	orderSourceRepository                   repositories.OrderSourceRepositoryInterface
	warehouseRepository                     repositories.WarehouseRepositoryInterface
	brandRepository                         repositories.BrandRepositoryInterface
	uomRepository                           repositories.UomRepositoryInterface
	agentRepository                         repositories.AgentRepositoryInterface
	storeRepository                         repositories.StoreRepositoryInterface
	productRepository                       repositories.ProductRepositoryInterface
	userRepository                          repositories.UserRepositoryInterface
	salesmanRepository                      repositories.SalesmanRepositoryInterface
	deliveryOrderLogRepository              mongoRepositories.DeliveryOrderLogRepositoryInterface
	deliveryOrderJourneysRepository         mongoRepositories.DeliveryOrderJourneysRepositoryInterface
	deliveryOrderOpenSearchRepository       openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface
	deliveryOrderDetailOpenSearchRepository openSearchRepositories.DeliveryOrderDetailOpenSearchRepositoryInterface
	SalesOrderOpenSearchUseCase             SalesOrderOpenSearchUseCaseInterface
	kafkaClient                             kafkadbo.KafkaClientInterface
	db                                      dbresolver.DB
	ctx                                     context.Context
}

func InitDeliveryOrderUseCaseInterface(deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, deliveryOrderDetailRepository repositories.DeliveryOrderDetailRepositoryInterface, salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, orderSourceRepository repositories.OrderSourceRepositoryInterface, warehouseRepository repositories.WarehouseRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, uomRepository repositories.UomRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, userRepository repositories.UserRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, deliveryOrderLogRepository mongoRepositories.DeliveryOrderLogRepositoryInterface, deliveryOrderJourneysRepository mongoRepositories.DeliveryOrderJourneysRepositoryInterface, deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface, deliveryOrderDetailOpenSearchRepository openSearchRepositories.DeliveryOrderDetailOpenSearchRepositoryInterface, salesOrderOpenSearchUseCase SalesOrderOpenSearchUseCaseInterface, kafkaClient kafkadbo.KafkaClientInterface, db dbresolver.DB, ctx context.Context) DeliveryOrderUseCaseInterface {
	return &deliveryOrderUseCase{
		deliveryOrderRepository:                 deliveryOrderRepository,
		deliveryOrderDetailRepository:           deliveryOrderDetailRepository,
		salesOrderRepository:                    salesOrderRepository,
		salesOrderDetailRepository:              salesOrderDetailRepository,
		orderStatusRepository:                   orderStatusRepository,
		orderSourceRepository:                   orderSourceRepository,
		warehouseRepository:                     warehouseRepository,
		brandRepository:                         brandRepository,
		uomRepository:                           uomRepository,
		productRepository:                       productRepository,
		userRepository:                          userRepository,
		salesmanRepository:                      salesmanRepository,
		agentRepository:                         agentRepository,
		storeRepository:                         storeRepository,
		deliveryOrderLogRepository:              deliveryOrderLogRepository,
		deliveryOrderJourneysRepository:         deliveryOrderJourneysRepository,
		deliveryOrderOpenSearchRepository:       deliveryOrderOpenSearchRepository,
		deliveryOrderDetailOpenSearchRepository: deliveryOrderDetailOpenSearchRepository,
		SalesOrderOpenSearchUseCase:             salesOrderOpenSearchUseCase,
		kafkaClient:                             kafkaClient,
		db:                                      db,
		ctx:                                     ctx,
	}
}

func (u *deliveryOrderUseCase) Create(request *models.DeliveryOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderStoreResponse, *model.ErrorLog) {
	now := time.Now()
	getSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.GetByID(request.SalesOrderID, false, ctx, getSalesOrderResultChan)
	getSalesOrderResult := <-getSalesOrderResultChan

	if getSalesOrderResult.Error != nil {
		return &models.DeliveryOrderStoreResponse{}, getSalesOrderResult.ErrorLog
	}

	getBrandResultChan := make(chan *models.BrandChan)
	go u.brandRepository.GetByID(getSalesOrderResult.SalesOrder.BrandID, false, ctx, getBrandResultChan)
	getBrandResult := <-getBrandResultChan

	if getBrandResult.Error != nil {
		return &models.DeliveryOrderStoreResponse{}, getBrandResult.ErrorLog
	}

	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType("open", "delivery_order", false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return &models.DeliveryOrderStoreResponse{}, getOrderStatusResult.ErrorLog
	}

	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetBySourceName("manager", false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getOrderSourceResult.Error != nil {
		return &models.DeliveryOrderStoreResponse{}, getOrderSourceResult.ErrorLog
	}

	getWarehouseResultChan := make(chan *models.WarehouseChan)
	go u.warehouseRepository.GetByID(request.WarehouseID, false, ctx, getWarehouseResultChan)
	getWarehouseResult := <-getWarehouseResultChan

	if getWarehouseResult.Error != nil {
		return &models.DeliveryOrderStoreResponse{}, getWarehouseResult.ErrorLog
	}

	getSalesOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(getSalesOrderResult.SalesOrder.OrderSourceID, false, ctx, getSalesOrderSourceResultChan)
	getSalesOrderSourceResult := <-getSalesOrderSourceResultChan

	if getSalesOrderSourceResult.Error != nil {
		return &models.DeliveryOrderStoreResponse{}, getSalesOrderSourceResult.ErrorLog
	}

	getSalesOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByID(getSalesOrderResult.SalesOrder.OrderStatusID, false, ctx, getSalesOrderStatusResultChan)
	getSalesOrderStatusResult := <-getSalesOrderStatusResultChan

	if getSalesOrderStatusResult.Error != nil {
		return &models.DeliveryOrderStoreResponse{}, getSalesOrderStatusResult.ErrorLog
	}

	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(getSalesOrderResult.SalesOrder.AgentID, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		return &models.DeliveryOrderStoreResponse{}, getAgentResult.ErrorLog
	}

	getStoreResultChan := make(chan *models.StoreChan)
	go u.storeRepository.GetByID(getSalesOrderResult.SalesOrder.StoreID, false, ctx, getStoreResultChan)
	getStoreResult := <-getStoreResultChan

	if getStoreResult.Error != nil {
		return &models.DeliveryOrderStoreResponse{}, getStoreResult.ErrorLog
	}

	getUserResultChan := make(chan *models.UserChan)
	go u.userRepository.GetByID(getSalesOrderResult.SalesOrder.UserID, false, ctx, getUserResultChan)
	getUserResult := <-getUserResultChan

	if getUserResult.Error != nil {
		return &models.DeliveryOrderStoreResponse{}, getUserResult.ErrorLog
	}

	getSalesmanResultChan := make(chan *models.SalesmanChan)
	if getSalesOrderResult.SalesOrder.SalesmanID.Int64 > 0 {
		go u.salesmanRepository.GetByID(int(getSalesOrderResult.SalesOrder.SalesmanID.Int64), false, ctx, getSalesmanResultChan)
	} else {
		go u.salesmanRepository.GetByEmail(getUserResult.User.Email, false, ctx, getSalesmanResultChan)
	}
	getSalesmanResult := <-getSalesmanResultChan

	if getSalesmanResult.Error != nil {
		// ignore null salesman
		if getSalesmanResult.Error.Error() != "salesman data not found" {
			return &models.DeliveryOrderStoreResponse{}, getSalesmanResult.ErrorLog
		}
	}

	deliveryOrder := &models.DeliveryOrder{}

	deliveryOrder.DeliveryOrderStoreRequestMap(request, now, ctx.Value("user").(*models.UserClaims))
	deliveryOrder.WarehouseChanMap(getWarehouseResult)
	deliveryOrder.AgentMap(getAgentResult.Agent)
	deliveryOrder.DoCode = helper.GenerateDOCode(getAgentResult.Agent.ID, getOrderSourceResult.OrderSource.Code)
	deliveryOrder.DoDate = now.Format("2006-01-02")
	deliveryOrder.OrderStatus = getOrderStatusResult.OrderStatus
	deliveryOrder.OrderStatusID = getOrderStatusResult.OrderStatus.ID
	deliveryOrder.OrderSource = getOrderSourceResult.OrderSource
	deliveryOrder.OrderSourceID = getOrderSourceResult.OrderSource.ID
	deliveryOrder.Store = getStoreResult.Store
	deliveryOrder.StoreID = getStoreResult.Store.ID
	deliveryOrder.SalesOrder = getSalesOrderResult.SalesOrder
	deliveryOrder.Brand = getBrandResult.Brand

	if getSalesmanResult.Salesman != nil {
		deliveryOrder.Salesman = getSalesmanResult.Salesman
	}

	createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderRepository.Insert(deliveryOrder, sqlTransaction, ctx, createDeliveryOrderResultChan)
	createDeliveryOrderResult := <-createDeliveryOrderResultChan

	if createDeliveryOrderResult.Error != nil {
		return &models.DeliveryOrderStoreResponse{}, createDeliveryOrderResult.ErrorLog
	}

	deliveryOrderDetails := []*models.DeliveryOrderDetail{}
	totalResidualQty := 0
	for _, doDetail := range request.DeliveryOrderDetails {
		getSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.GetByID(doDetail.SoDetailID, false, ctx, getSalesOrderDetailResultChan)
		getSalesOrderDetailResult := <-getSalesOrderDetailResultChan

		if getSalesOrderDetailResult.Error != nil {
			return &models.DeliveryOrderStoreResponse{}, getSalesOrderDetailResult.ErrorLog
		}

		getOrderStatusDetailResultChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID, false, ctx, getOrderStatusDetailResultChan)
		getOrderStatusDetailResult := <-getOrderStatusDetailResultChan

		if getOrderStatusDetailResult.Error != nil {
			return &models.DeliveryOrderStoreResponse{}, getOrderStatusDetailResult.ErrorLog
		}

		getSalesOrderDetailResult.SalesOrderDetail.UpdatedAt = &now
		getSalesOrderDetailResult.SalesOrderDetail.SentQty += doDetail.Qty
		getSalesOrderDetailResult.SalesOrderDetail.ResidualQty -= doDetail.Qty
		totalResidualQty += getSalesOrderDetailResult.SalesOrderDetail.ResidualQty
		statusName := "partial"

		if getSalesOrderDetailResult.SalesOrderDetail.ResidualQty == 0 {
			statusName = "closed"
		}

		getOrderStatusSODetailResultChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByNameAndType(statusName, "sales_order_detail", false, ctx, getOrderStatusSODetailResultChan)
		getOrderStatusSODetailResult := <-getOrderStatusSODetailResultChan

		if getOrderStatusSODetailResult.Error != nil {
			return &models.DeliveryOrderStoreResponse{}, getOrderStatusSODetailResult.ErrorLog
		}
		getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID = getOrderStatusSODetailResult.OrderStatus.ID
		getSalesOrderDetailResult.SalesOrderDetail.OrderStatusName = getOrderStatusSODetailResult.OrderStatus.Name
		getSalesOrderDetailResult.SalesOrderDetail.OrderStatus = getOrderStatusSODetailResult.OrderStatus

		getProductDetailResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.ProductID, false, ctx, getProductDetailResultChan)
		getProductDetailResult := <-getProductDetailResultChan

		if getProductDetailResult.Error != nil {
			return &models.DeliveryOrderStoreResponse{}, getProductDetailResult.ErrorLog
		}

		getUomDetailResultChan := make(chan *models.UomChan)
		go u.uomRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.UomID, false, ctx, getUomDetailResultChan)
		getUomDetailResult := <-getUomDetailResultChan

		if getUomDetailResult.Error != nil {
			return &models.DeliveryOrderStoreResponse{}, getUomDetailResult.ErrorLog
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
		deliveryOrderDetail.OrderStatusID = deliveryOrder.OrderStatusID
		deliveryOrderDetail.OrderStatusName = deliveryOrder.OrderStatusName
		deliveryOrderDetail.OrderStatus = getOrderStatusDetailResult.OrderStatus

		createDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
		go u.deliveryOrderDetailRepository.Insert(deliveryOrderDetail, sqlTransaction, ctx, createDeliveryOrderDetailResultChan)
		createDeliveryOrderDetailResult := <-createDeliveryOrderDetailResultChan

		if createDeliveryOrderDetailResult.Error != nil {
			return &models.DeliveryOrderStoreResponse{}, createDeliveryOrderDetailResult.ErrorLog
		}

		updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.UpdateByID(getSalesOrderDetailResult.SalesOrderDetail.ID, getSalesOrderDetailResult.SalesOrderDetail, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
		updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

		if updateSalesOrderDetailResult.Error != nil {
			return &models.DeliveryOrderStoreResponse{}, updateSalesOrderDetailResult.ErrorLog
		}

		deliveryOrderDetail.ID = int(createDeliveryOrderDetailResult.ID)
		deliveryOrderDetails = append(deliveryOrderDetails, deliveryOrderDetail)
		getSalesOrderResult.SalesOrder.SalesOrderDetails = append(getSalesOrderResult.SalesOrder.SalesOrderDetails, getSalesOrderDetailResult.SalesOrderDetail)
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
		return &models.DeliveryOrderStoreResponse{}, getOrderStatusSODetailResult.ErrorLog
	}
	getSalesOrderResult.SalesOrder.OrderStatus = getOrderStatusSODetailResult.OrderStatus
	getSalesOrderResult.SalesOrder.OrderStatusName = getOrderStatusSODetailResult.OrderStatus.Name
	getSalesOrderResult.SalesOrder.SoDate = ""
	getSalesOrderResult.SalesOrder.SoRefDate = models.NullString{}
	getSalesOrderResult.SalesOrder.UpdatedAt = &now
	deliveryOrder.SalesOrder = getSalesOrderResult.SalesOrder

	updateSalesOrderChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.UpdateByID(getSalesOrderResult.SalesOrder.ID, getSalesOrderResult.SalesOrder, sqlTransaction, u.ctx, updateSalesOrderChan)
	updateSalesOrderResult := <-updateSalesOrderChan

	if updateSalesOrderResult.ErrorLog != nil {
		err := sqlTransaction.Rollback()
		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return &models.DeliveryOrderStoreResponse{}, errorLogData
		}
		return &models.DeliveryOrderStoreResponse{}, updateSalesOrderResult.ErrorLog
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
		return &models.DeliveryOrderStoreResponse{}, createDeliveryOrderLogResult.ErrorLog
	}

	deliveryOrderJourney := &models.DeliveryOrderJourney{
		DoId:      deliveryOrder.ID,
		DoCode:    deliveryOrder.DoCode,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		Remark:    "",
		Reason:    "",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	createDeliveryOrderJourneyChan := make(chan *models.DeliveryOrderJourneyChan)
	go u.deliveryOrderLogRepository.InsertJourney(deliveryOrderJourney, ctx, createDeliveryOrderJourneyChan)
	createDeliveryOrderJourneysResult := <-createDeliveryOrderJourneyChan

	if createDeliveryOrderJourneysResult.Error != nil {
		return &models.DeliveryOrderStoreResponse{}, createDeliveryOrderJourneysResult.ErrorLog
	}

	keyKafka := []byte(deliveryOrder.DoCode)
	messageKafka, _ := json.Marshal(deliveryOrder)
	err := u.kafkaClient.WriteToTopic(constants.CREATE_DELIVERY_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.DeliveryOrderStoreResponse{}, errorLogData
	}

	deliveryOrderDetailResults := []*models.DeliveryOrderDetailStoreResponse{}
	for _, v := range deliveryOrder.DeliveryOrderDetails {
		deliveryOrderDetailResult := &models.DeliveryOrderDetailStoreResponse{}
		deliveryOrderDetailResult.DeliveryOrderDetailMap(v)
		deliveryOrderDetailResults = append(deliveryOrderDetailResults, deliveryOrderDetailResult)
	}

	deliveryOrderResult := &models.DeliveryOrderStoreResponse{}
	deliveryOrderResult.DeliveryOrderMap(deliveryOrder)
	return deliveryOrderResult, nil
}

func (u *deliveryOrderUseCase) UpdateByID(ID int, request *models.DeliveryOrderUpdateByIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderUpdateByIDRequest, *model.ErrorLog) {
	now := time.Now()

	getDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderRepository.GetByID(ID, false, ctx, getDeliveryOrderResultChan)
	getDeliveryOrderResult := <-getDeliveryOrderResultChan

	if getDeliveryOrderResult.Error != nil {
		return &models.DeliveryOrderUpdateByIDRequest{}, getDeliveryOrderResult.ErrorLog
	}

	getSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.GetByID(getDeliveryOrderResult.DeliveryOrder.SalesOrderID, false, ctx, getSalesOrderResultChan)
	getSalesOrderResult := <-getSalesOrderResultChan

	if getSalesOrderResult.Error != nil {
		return &models.DeliveryOrderUpdateByIDRequest{}, getDeliveryOrderResult.ErrorLog
	}

	warehouseId := request.WarehouseID
	if warehouseId == 0 {
		warehouseId = getDeliveryOrderResult.DeliveryOrder.WarehouseID
	}
	getWarehouseResultChan := make(chan *models.WarehouseChan)
	go u.warehouseRepository.GetByID(warehouseId, false, ctx, getWarehouseResultChan)
	getWarehouseResult := <-getWarehouseResultChan

	if getWarehouseResult.Error != nil {
		return &models.DeliveryOrderUpdateByIDRequest{}, getWarehouseResult.ErrorLog
	}

	orderSourceID := request.OrderSourceID
	if orderSourceID == 0 {
		orderSourceID = getDeliveryOrderResult.DeliveryOrder.OrderSourceID
	}
	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(orderSourceID, false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getDeliveryOrderResult.Error != nil {
		return &models.DeliveryOrderUpdateByIDRequest{}, getOrderSourceResult.ErrorLog
	}

	orderStatusID := request.OrderStatusID
	if orderStatusID == 0 {
		orderStatusID = getDeliveryOrderResult.DeliveryOrder.OrderStatusID
	}
	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByID(orderStatusID, false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return &models.DeliveryOrderUpdateByIDRequest{}, getOrderStatusResult.ErrorLog
	}

	deliveryOrder := getDeliveryOrderResult.DeliveryOrder
	deliveryOrder.DeliveryOrderUpdateByIDRequestMap(request, now, ctx.Value("user").(*models.UserClaims))
	deliveryOrder.WarehouseChanMap(getWarehouseResult)
	deliveryOrder.OrderStatus = getOrderStatusResult.OrderStatus
	deliveryOrder.OrderStatusID = getOrderStatusResult.OrderStatus.ID
	deliveryOrder.OrderSource = getOrderSourceResult.OrderSource
	deliveryOrder.OrderSourceID = getOrderSourceResult.OrderSource.ID
	deliveryOrder.CreatedBy = ctx.Value("user").(*models.UserClaims).UserID
	deliveryOrder.SalesOrder = getSalesOrderResult.SalesOrder

	getDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailsChan)
	go u.deliveryOrderDetailRepository.GetByDeliveryOrderID(ID, false, ctx, getDeliveryOrderDetailResultChan)
	getDeliveryOrderDetailResult := <-getDeliveryOrderDetailResultChan

	if getDeliveryOrderDetailResult.Error != nil {
		return &models.DeliveryOrderUpdateByIDRequest{}, getDeliveryOrderDetailResult.ErrorLog
	}

	deliveryOrderDetails := []*models.DeliveryOrderDetail{}
	deliveryOrderDetailResults := []*models.DeliveryOrderDetailUpdateByIDRequest{}
	salesOrderDetails := []*models.SalesOrderDetail{}
	totalSentQty := 0
	totalQty := 0
	for _, x := range request.DeliveryOrderDetails {
		for _, v := range getDeliveryOrderDetailResult.DeliveryOrderDetails {
			if x.ID == v.ID {
				balanceQty := x.Qty - v.Qty
				v.Qty = x.Qty
				orderStatusID := 19
				if v.Qty > 0 {
					orderStatusID = 18
				}
				getOrderStatusResultChan := make(chan *models.OrderStatusChan)
				go u.orderStatusRepository.GetByID(orderStatusID, false, ctx, getOrderStatusResultChan)
				getOrderStatusResult := <-getOrderStatusResultChan

				if getOrderStatusResult.Error != nil {
					return &models.DeliveryOrderUpdateByIDRequest{}, getOrderStatusResult.ErrorLog
				}

				deliveryOrder.OrderStatus = getOrderStatusResult.OrderStatus
				deliveryOrder.OrderStatusID = getOrderStatusResult.OrderStatus.ID
				deliveryOrder.OrderStatusName = getOrderStatusResult.OrderStatus.Name

				v.OrderStatus = getOrderStatusResult.OrderStatus
				v.OrderStatusID = getOrderStatusResult.OrderStatus.ID
				v.OrderStatusName = getOrderStatusResult.OrderStatus.Name

				updateDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
				go u.deliveryOrderDetailRepository.UpdateByID(v.ID, v, sqlTransaction, ctx, updateDeliveryOrderDetailResultChan)
				updateDeliveryOrderDetailResult := <-updateDeliveryOrderDetailResultChan

				if updateDeliveryOrderDetailResult.Error != nil {
					return &models.DeliveryOrderUpdateByIDRequest{}, updateDeliveryOrderDetailResult.ErrorLog
				}

				getSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
				go u.salesOrderDetailRepository.GetByID(v.SoDetailID, false, ctx, getSalesOrderDetailResultChan)
				getSalesOrderDetailResult := <-getSalesOrderDetailResultChan

				if getSalesOrderDetailResult.Error != nil {
					return &models.DeliveryOrderUpdateByIDRequest{}, getSalesOrderDetailResult.ErrorLog
				}

				totalQty += getSalesOrderDetailResult.SalesOrderDetail.Qty

				if balanceQty != 0 {

					getOrderStatusDetailResultChan := make(chan *models.OrderStatusChan)
					go u.orderStatusRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID, false, ctx, getOrderStatusDetailResultChan)
					getOrderStatusDetailResult := <-getOrderStatusDetailResultChan

					if getOrderStatusDetailResult.Error != nil {
						return &models.DeliveryOrderUpdateByIDRequest{}, getOrderStatusDetailResult.ErrorLog
					}

					getSalesOrderDetailResult.SalesOrderDetail.UpdatedAt = &now
					getSalesOrderDetailResult.SalesOrderDetail.SentQty += balanceQty
					getSalesOrderDetailResult.SalesOrderDetail.ResidualQty -= balanceQty
					totalSentQty += getSalesOrderDetailResult.SalesOrderDetail.SentQty

					if getSalesOrderDetailResult.SalesOrderDetail.SentQty == 0 {
						getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID = 11
					} else if getSalesOrderDetailResult.SalesOrderDetail.SentQty == getSalesOrderDetailResult.SalesOrderDetail.Qty {
						getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID = 14
					} else {
						getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID = 13
					}

					getSalesOrderDetailResult.SalesOrderDetail.IsDoneSyncToEs = "0"
					getSalesOrderDetailResult.SalesOrderDetail.StartDateSyncToEs = &now

					getOrderStatusSODetailResultChan := make(chan *models.OrderStatusChan)
					go u.orderStatusRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID, false, ctx, getOrderStatusSODetailResultChan)
					getOrderStatusSODetailResult := <-getOrderStatusSODetailResultChan

					if getOrderStatusSODetailResult.Error != nil {
						return &models.DeliveryOrderUpdateByIDRequest{}, getOrderStatusSODetailResult.ErrorLog
					}

					getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID = getOrderStatusSODetailResult.OrderStatus.ID
					getSalesOrderDetailResult.SalesOrderDetail.OrderStatusName = getOrderStatusSODetailResult.OrderStatus.Name
					getSalesOrderDetailResult.SalesOrderDetail.OrderStatus = getOrderStatusSODetailResult.OrderStatus

					updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
					go u.salesOrderDetailRepository.UpdateByID(getSalesOrderDetailResult.SalesOrderDetail.ID, getSalesOrderDetailResult.SalesOrderDetail, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
					updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

					if updateSalesOrderDetailResult.Error != nil {
						return &models.DeliveryOrderUpdateByIDRequest{}, updateSalesOrderDetailResult.ErrorLog
					}
				} else {
					totalSentQty += v.Qty
				}
				deliveryOrderDetails = append(deliveryOrderDetails, v)
				deliveryOrderDetailResults = append(deliveryOrderDetailResults, &models.DeliveryOrderDetailUpdateByIDRequest{
					Qty:  v.Qty,
					Note: v.Note.String,
				})
				salesOrderDetails = append(salesOrderDetails, getSalesOrderDetailResult.SalesOrderDetail)
			}
		}
	}

	deliveryOrder.DeliveryOrderDetails = deliveryOrderDetails

	updateDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderRepository.UpdateByID(getDeliveryOrderResult.DeliveryOrder.ID, deliveryOrder, sqlTransaction, ctx, updateDeliveryOrderResultChan)
	updateDeliveryOrderResult := <-updateDeliveryOrderResultChan

	if updateDeliveryOrderResult.Error != nil {
		return &models.DeliveryOrderUpdateByIDRequest{}, updateDeliveryOrderResult.ErrorLog
	}

	if totalSentQty == 0 {
		getSalesOrderResult.SalesOrder.OrderStatusID = 5
	} else if totalSentQty == totalQty {
		getSalesOrderResult.SalesOrder.OrderStatusID = 8
	} else {
		getSalesOrderResult.SalesOrder.OrderStatusID = 7
	}

	getOrderStatusSOResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByID(getSalesOrderResult.SalesOrder.OrderStatusID, false, ctx, getOrderStatusSOResultChan)
	getOrderStatusSODetailResult := <-getOrderStatusSOResultChan

	if getOrderStatusSODetailResult.Error != nil {
		return &models.DeliveryOrderUpdateByIDRequest{}, getOrderStatusSODetailResult.ErrorLog
	}
	getSalesOrderResult.SalesOrder.OrderStatus = getOrderStatusSODetailResult.OrderStatus
	getSalesOrderResult.SalesOrder.OrderStatusName = getOrderStatusSODetailResult.OrderStatus.Name
	getSalesOrderResult.SalesOrder.SoDate = ""
	getSalesOrderResult.SalesOrder.SoRefDate = models.NullString{}
	getSalesOrderResult.SalesOrder.UpdatedAt = &now
	getSalesOrderResult.SalesOrder.SalesOrderDetails = salesOrderDetails
	deliveryOrder.SalesOrder = getSalesOrderResult.SalesOrder

	updateSalesOrderChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.UpdateByID(getSalesOrderResult.SalesOrder.ID, getSalesOrderResult.SalesOrder, sqlTransaction, u.ctx, updateSalesOrderChan)
	updateSalesOrderResult := <-updateSalesOrderChan

	if updateSalesOrderResult.ErrorLog != nil {
		err := sqlTransaction.Rollback()
		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return &models.DeliveryOrderUpdateByIDRequest{}, errorLogData
		}
		return &models.DeliveryOrderUpdateByIDRequest{}, updateSalesOrderResult.ErrorLog
	}

	deliveryOrderLog := &models.DeliveryOrderLog{
		RequestID: request.RequestID,
		DoCode:    deliveryOrder.DoCode,
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
		return &models.DeliveryOrderUpdateByIDRequest{}, errorLogData
	}

	deliveryOrderJourney := &models.DeliveryOrderJourney{
		DoId:      deliveryOrder.ID,
		DoCode:    deliveryOrder.DoCode,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		Remark:    "",
		Reason:    "",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	createDeliveryOrderJourneyChan := make(chan *models.DeliveryOrderJourneyChan)
	go u.deliveryOrderLogRepository.InsertJourney(deliveryOrderJourney, ctx, createDeliveryOrderJourneyChan)
	createDeliveryOrderJourneysResult := <-createDeliveryOrderJourneyChan

	if createDeliveryOrderJourneysResult.Error != nil {
		return &models.DeliveryOrderUpdateByIDRequest{}, createDeliveryOrderJourneysResult.ErrorLog
	}

	keyKafka := []byte(deliveryOrder.DoCode)
	messageKafka, _ := json.Marshal(deliveryOrder)
	err := u.kafkaClient.WriteToTopic(constants.UPDATE_DELIVERY_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.DeliveryOrderUpdateByIDRequest{}, errorLogData
	}

	deliveryOrderResult := &models.DeliveryOrderUpdateByIDRequest{
		WarehouseID:          deliveryOrder.WarehouseID,
		OrderSourceID:        deliveryOrder.OrderSourceID,
		OrderStatusID:        deliveryOrder.OrderStatusID,
		DoRefCode:            deliveryOrder.DoRefCode.String,
		DoRefDate:            deliveryOrder.DoRefDate.String,
		DriverName:           deliveryOrder.DriverName.String,
		PlatNumber:           deliveryOrder.PlatNumber.String,
		Note:                 deliveryOrder.Note.String,
		DeliveryOrderDetails: deliveryOrderDetailResults,
	}

	return deliveryOrderResult, nil
}

func (u *deliveryOrderUseCase) UpdateDODetailByID(id int, request *models.DeliveryOrderDetailUpdateByIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderDetailUpdateByIDRequest, *model.ErrorLog) {
	now := time.Now()

	getDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
	go u.deliveryOrderDetailRepository.GetByID(id, false, ctx, getDeliveryOrderDetailResultChan)
	getDeliveryOrderDetailResult := <-getDeliveryOrderDetailResultChan

	if getDeliveryOrderDetailResult.Error != nil {
		return &models.DeliveryOrderDetailUpdateByIDRequest{}, getDeliveryOrderDetailResult.ErrorLog
	}

	getDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderRepository.GetByID(getDeliveryOrderDetailResult.DeliveryOrderDetail.DeliveryOrderID, false, ctx, getDeliveryOrderResultChan)
	getDeliveryOrderResult := <-getDeliveryOrderResultChan

	if getDeliveryOrderResult.Error != nil {
		return &models.DeliveryOrderDetailUpdateByIDRequest{}, getDeliveryOrderResult.ErrorLog
	}

	getSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.GetByID(getDeliveryOrderResult.DeliveryOrder.SalesOrderID, false, ctx, getSalesOrderResultChan)
	getSalesOrderResult := <-getSalesOrderResultChan

	if getSalesOrderResult.Error != nil {
		return &models.DeliveryOrderDetailUpdateByIDRequest{}, getDeliveryOrderResult.ErrorLog
	}

	getWarehouseResultChan := make(chan *models.WarehouseChan)
	go u.warehouseRepository.GetByID(getDeliveryOrderResult.DeliveryOrder.WarehouseID, false, ctx, getWarehouseResultChan)
	getWarehouseResult := <-getWarehouseResultChan

	if getWarehouseResult.Error != nil {
		return &models.DeliveryOrderDetailUpdateByIDRequest{}, getWarehouseResult.ErrorLog
	}
	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(getDeliveryOrderResult.DeliveryOrder.OrderSourceID, false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getDeliveryOrderResult.Error != nil {
		return &models.DeliveryOrderDetailUpdateByIDRequest{}, getOrderSourceResult.ErrorLog
	}
	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByID(getDeliveryOrderResult.DeliveryOrder.OrderStatusID, false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return &models.DeliveryOrderDetailUpdateByIDRequest{}, getOrderStatusResult.ErrorLog
	}

	deliveryOrder := getDeliveryOrderResult.DeliveryOrder
	deliveryOrder.WarehouseChanMap(getWarehouseResult)
	deliveryOrder.IsDoneSyncToEs = "0"
	deliveryOrder.StartDateSyncToEs = &now
	deliveryOrder.EndDateSyncToEs = &now
	deliveryOrder.LatestUpdatedBy = ctx.Value("user").(*models.UserClaims).UserID
	deliveryOrder.UpdatedAt = &now
	deliveryOrder.DeletedAt = nil
	deliveryOrder.OrderStatus = getOrderStatusResult.OrderStatus
	deliveryOrder.OrderStatusID = getOrderStatusResult.OrderStatus.ID
	deliveryOrder.OrderSource = getOrderSourceResult.OrderSource
	deliveryOrder.OrderSourceID = getOrderSourceResult.OrderSource.ID
	deliveryOrder.CreatedBy = ctx.Value("user").(*models.UserClaims).UserID
	deliveryOrder.SalesOrder = getSalesOrderResult.SalesOrder

	getDeliveryOrderDetailsResultChan := make(chan *models.DeliveryOrderDetailsChan)
	go u.deliveryOrderDetailRepository.GetByDeliveryOrderID(deliveryOrder.ID, false, ctx, getDeliveryOrderDetailsResultChan)
	getDeliveryOrderDetailsResult := <-getDeliveryOrderDetailsResultChan

	if getDeliveryOrderDetailResult.Error != nil {
		return &models.DeliveryOrderDetailUpdateByIDRequest{}, getDeliveryOrderDetailResult.ErrorLog
	}

	deliveryOrderDetails := []*models.DeliveryOrderDetail{}
	deliveryOrderDetailResults := []*models.DeliveryOrderDetailUpdateByIDRequest{}
	salesOrderDetails := []*models.SalesOrderDetail{}
	totalSentQty := 0
	totalQty := 0
	for _, v := range getDeliveryOrderDetailsResult.DeliveryOrderDetails {
		if request.ID == v.ID {
			balanceQty := request.Qty - v.Qty
			v.Qty = request.Qty
			orderStatusID := 19
			if v.Qty > 0 {
				orderStatusID = 18
			}
			getOrderStatusResultChan := make(chan *models.OrderStatusChan)
			go u.orderStatusRepository.GetByID(orderStatusID, false, ctx, getOrderStatusResultChan)
			getOrderStatusResult := <-getOrderStatusResultChan

			if getOrderStatusResult.Error != nil {
				return &models.DeliveryOrderDetailUpdateByIDRequest{}, getOrderStatusResult.ErrorLog
			}

			deliveryOrder.OrderStatus = getOrderStatusResult.OrderStatus
			deliveryOrder.OrderStatusID = getOrderStatusResult.OrderStatus.ID
			deliveryOrder.OrderStatusName = getOrderStatusResult.OrderStatus.Name

			v.OrderStatus = getOrderStatusResult.OrderStatus
			v.OrderStatusID = getOrderStatusResult.OrderStatus.ID
			v.OrderStatusName = getOrderStatusResult.OrderStatus.Name

			updateDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
			go u.deliveryOrderDetailRepository.UpdateByID(v.ID, v, sqlTransaction, ctx, updateDeliveryOrderDetailResultChan)
			updateDeliveryOrderDetailResult := <-updateDeliveryOrderDetailResultChan

			if updateDeliveryOrderDetailResult.Error != nil {
				return &models.DeliveryOrderDetailUpdateByIDRequest{}, updateDeliveryOrderDetailResult.ErrorLog
			}

			getSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
			go u.salesOrderDetailRepository.GetByID(v.SoDetailID, false, ctx, getSalesOrderDetailResultChan)
			getSalesOrderDetailResult := <-getSalesOrderDetailResultChan

			if getSalesOrderDetailResult.Error != nil {
				return &models.DeliveryOrderDetailUpdateByIDRequest{}, getSalesOrderDetailResult.ErrorLog
			}

			totalQty += getSalesOrderDetailResult.SalesOrderDetail.Qty

			if balanceQty != 0 {

				getOrderStatusDetailResultChan := make(chan *models.OrderStatusChan)
				go u.orderStatusRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID, false, ctx, getOrderStatusDetailResultChan)
				getOrderStatusDetailResult := <-getOrderStatusDetailResultChan

				if getOrderStatusDetailResult.Error != nil {
					return &models.DeliveryOrderDetailUpdateByIDRequest{}, getOrderStatusDetailResult.ErrorLog
				}

				getSalesOrderDetailResult.SalesOrderDetail.UpdatedAt = &now
				getSalesOrderDetailResult.SalesOrderDetail.SentQty += balanceQty
				getSalesOrderDetailResult.SalesOrderDetail.ResidualQty -= balanceQty
				totalSentQty += getSalesOrderDetailResult.SalesOrderDetail.SentQty

				if getSalesOrderDetailResult.SalesOrderDetail.SentQty == 0 {
					getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID = 11
				} else if getSalesOrderDetailResult.SalesOrderDetail.SentQty == getSalesOrderDetailResult.SalesOrderDetail.Qty {
					getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID = 14
				} else {
					getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID = 13
				}

				getSalesOrderDetailResult.SalesOrderDetail.IsDoneSyncToEs = "0"
				getSalesOrderDetailResult.SalesOrderDetail.StartDateSyncToEs = &now

				getOrderStatusSODetailResultChan := make(chan *models.OrderStatusChan)
				go u.orderStatusRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID, false, ctx, getOrderStatusSODetailResultChan)
				getOrderStatusSODetailResult := <-getOrderStatusSODetailResultChan

				if getOrderStatusSODetailResult.Error != nil {
					return &models.DeliveryOrderDetailUpdateByIDRequest{}, getOrderStatusSODetailResult.ErrorLog
				}

				getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID = getOrderStatusSODetailResult.OrderStatus.ID
				getSalesOrderDetailResult.SalesOrderDetail.OrderStatusName = getOrderStatusSODetailResult.OrderStatus.Name
				getSalesOrderDetailResult.SalesOrderDetail.OrderStatus = getOrderStatusSODetailResult.OrderStatus

				updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
				go u.salesOrderDetailRepository.UpdateByID(getSalesOrderDetailResult.SalesOrderDetail.ID, getSalesOrderDetailResult.SalesOrderDetail, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
				updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

				if updateSalesOrderDetailResult.Error != nil {
					return &models.DeliveryOrderDetailUpdateByIDRequest{}, updateSalesOrderDetailResult.ErrorLog
				}
			} else {
				totalSentQty += v.Qty
			}
			deliveryOrderDetails = append(deliveryOrderDetails, v)
			deliveryOrderDetailResults = append(deliveryOrderDetailResults, &models.DeliveryOrderDetailUpdateByIDRequest{
				Qty:  v.Qty,
				Note: v.Note.String,
			})
			salesOrderDetails = append(salesOrderDetails, getSalesOrderDetailResult.SalesOrderDetail)
		}

	}

	deliveryOrder.DeliveryOrderDetails = deliveryOrderDetails

	updateDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderRepository.UpdateByID(getDeliveryOrderResult.DeliveryOrder.ID, deliveryOrder, sqlTransaction, ctx, updateDeliveryOrderResultChan)
	updateDeliveryOrderResult := <-updateDeliveryOrderResultChan

	if updateDeliveryOrderResult.Error != nil {
		return &models.DeliveryOrderDetailUpdateByIDRequest{}, updateDeliveryOrderResult.ErrorLog
	}

	if totalSentQty == 0 {
		getSalesOrderResult.SalesOrder.OrderStatusID = 5
	} else if totalSentQty == totalQty {
		getSalesOrderResult.SalesOrder.OrderStatusID = 8
	} else {
		getSalesOrderResult.SalesOrder.OrderStatusID = 7
	}

	getOrderStatusSOResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByID(getSalesOrderResult.SalesOrder.OrderStatusID, false, ctx, getOrderStatusSOResultChan)
	getOrderStatusSODetailResult := <-getOrderStatusSOResultChan

	if getOrderStatusSODetailResult.Error != nil {
		return &models.DeliveryOrderDetailUpdateByIDRequest{}, getOrderStatusSODetailResult.ErrorLog
	}
	getSalesOrderResult.SalesOrder.OrderStatus = getOrderStatusSODetailResult.OrderStatus
	getSalesOrderResult.SalesOrder.OrderStatusName = getOrderStatusSODetailResult.OrderStatus.Name
	getSalesOrderResult.SalesOrder.SoDate = ""
	getSalesOrderResult.SalesOrder.SoRefDate = models.NullString{}
	getSalesOrderResult.SalesOrder.UpdatedAt = &now
	getSalesOrderResult.SalesOrder.SalesOrderDetails = salesOrderDetails
	deliveryOrder.SalesOrder = getSalesOrderResult.SalesOrder

	updateSalesOrderChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.UpdateByID(getSalesOrderResult.SalesOrder.ID, getSalesOrderResult.SalesOrder, sqlTransaction, u.ctx, updateSalesOrderChan)
	updateSalesOrderResult := <-updateSalesOrderChan

	if updateSalesOrderResult.ErrorLog != nil {
		err := sqlTransaction.Rollback()
		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return &models.DeliveryOrderDetailUpdateByIDRequest{}, errorLogData
		}
		return &models.DeliveryOrderDetailUpdateByIDRequest{}, updateSalesOrderResult.ErrorLog
	}

	deliveryOrderLog := &models.DeliveryOrderLog{
		RequestID: request.RequestID,
		DoCode:    deliveryOrder.DoCode,
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
		return &models.DeliveryOrderDetailUpdateByIDRequest{}, errorLogData
	}

	deliveryOrderJourney := &models.DeliveryOrderJourney{
		DoId:      deliveryOrder.ID,
		DoCode:    deliveryOrder.DoCode,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		Remark:    "",
		Reason:    "",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	createDeliveryOrderJourneyChan := make(chan *models.DeliveryOrderJourneyChan)
	go u.deliveryOrderLogRepository.InsertJourney(deliveryOrderJourney, ctx, createDeliveryOrderJourneyChan)
	createDeliveryOrderJourneysResult := <-createDeliveryOrderJourneyChan

	if createDeliveryOrderJourneysResult.Error != nil {
		return &models.DeliveryOrderDetailUpdateByIDRequest{}, createDeliveryOrderJourneysResult.ErrorLog
	}

	keyKafka := []byte(deliveryOrder.DoCode)
	messageKafka, _ := json.Marshal(deliveryOrder)
	err := u.kafkaClient.WriteToTopic(constants.UPDATE_DELIVERY_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.DeliveryOrderDetailUpdateByIDRequest{}, errorLogData
	}

	deliveryOrderResult := &models.DeliveryOrderDetailUpdateByIDRequest{
		Qty:  request.Qty,
		Note: request.Note,
	}

	return deliveryOrderResult, nil
}

func (u *deliveryOrderUseCase) UpdateDoDetailByDeliveryOrderID(deliveryOrderID int, request []*models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderDetails, *model.ErrorLog) {
	now := time.Now()

	getDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderRepository.GetByID(deliveryOrderID, false, ctx, getDeliveryOrderResultChan)
	getDeliveryOrderResult := <-getDeliveryOrderResultChan

	if getDeliveryOrderResult.Error != nil {
		return &models.DeliveryOrderDetails{}, getDeliveryOrderResult.ErrorLog
	}

	getSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.GetByID(getDeliveryOrderResult.DeliveryOrder.SalesOrderID, false, ctx, getSalesOrderResultChan)
	getSalesOrderResult := <-getSalesOrderResultChan

	if getSalesOrderResult.Error != nil {
		return &models.DeliveryOrderDetails{}, getDeliveryOrderResult.ErrorLog
	}

	deliveryOrder := getDeliveryOrderResult.DeliveryOrder

	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByID(deliveryOrder.OrderStatusID, false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	deliveryOrder.OrderStatusID = getOrderStatusResult.OrderStatus.ID
	deliveryOrder.OrderStatusName = getOrderStatusResult.OrderStatus.Name
	deliveryOrder.OrderStatus = getOrderStatusResult.OrderStatus

	getDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailsChan)
	go u.deliveryOrderDetailRepository.GetByDeliveryOrderID(deliveryOrderID, false, ctx, getDeliveryOrderDetailResultChan)
	getDeliveryOrderDetailResult := <-getDeliveryOrderDetailResultChan

	if getDeliveryOrderDetailResult.Error != nil {
		return &models.DeliveryOrderDetails{}, getDeliveryOrderDetailResult.ErrorLog
	}

	deliveryOrderDetails := []*models.DeliveryOrderDetail{}
	deliveryOrderDetailResults := []*models.DeliveryOrderDetailUpdateByIDRequest{}
	totalSentQty := 0
	totalQty := 0
	for _, x := range request {
		for _, v := range getDeliveryOrderDetailResult.DeliveryOrderDetails {
			if x.ID == v.ID {
				balanceQty := x.Qty - v.Qty
				v.Qty = x.Qty
				orderStatusID := 19
				if v.Qty > 0 {
					orderStatusID = 18
				}
				getOrderStatusResultChan := make(chan *models.OrderStatusChan)
				go u.orderStatusRepository.GetByID(orderStatusID, false, ctx, getOrderStatusResultChan)
				getOrderStatusResult := <-getOrderStatusResultChan

				if getOrderStatusResult.Error != nil {
					return &models.DeliveryOrderDetails{}, getOrderStatusResult.ErrorLog
				}

				v.OrderStatus = getOrderStatusResult.OrderStatus
				v.OrderStatusID = getOrderStatusResult.OrderStatus.ID
				v.OrderStatusName = getOrderStatusResult.OrderStatus.Name

				updateDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
				go u.deliveryOrderDetailRepository.UpdateByID(v.ID, v, sqlTransaction, ctx, updateDeliveryOrderDetailResultChan)
				updateDeliveryOrderDetailResult := <-updateDeliveryOrderDetailResultChan

				if updateDeliveryOrderDetailResult.Error != nil {
					return &models.DeliveryOrderDetails{}, updateDeliveryOrderDetailResult.ErrorLog
				}

				getSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
				go u.salesOrderDetailRepository.GetByID(v.SoDetailID, false, ctx, getSalesOrderDetailResultChan)
				getSalesOrderDetailResult := <-getSalesOrderDetailResultChan

				if getSalesOrderDetailResult.Error != nil {
					return &models.DeliveryOrderDetails{}, getSalesOrderDetailResult.ErrorLog
				}

				totalQty += getSalesOrderDetailResult.SalesOrderDetail.Qty

				if balanceQty != 0 {

					getOrderStatusDetailResultChan := make(chan *models.OrderStatusChan)
					go u.orderStatusRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID, false, ctx, getOrderStatusDetailResultChan)
					getOrderStatusDetailResult := <-getOrderStatusDetailResultChan

					if getOrderStatusDetailResult.Error != nil {
						return &models.DeliveryOrderDetails{}, getOrderStatusDetailResult.ErrorLog
					}

					getSalesOrderDetailResult.SalesOrderDetail.UpdatedAt = &now
					getSalesOrderDetailResult.SalesOrderDetail.SentQty += balanceQty
					getSalesOrderDetailResult.SalesOrderDetail.ResidualQty -= balanceQty

					if getSalesOrderDetailResult.SalesOrderDetail.SentQty == 0 {
						getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID = 11
					} else if getSalesOrderDetailResult.SalesOrderDetail.SentQty == getSalesOrderDetailResult.SalesOrderDetail.Qty {
						getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID = 14
					} else {
						getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID = 13
					}

					getOrderStatusSODetailResultChan := make(chan *models.OrderStatusChan)
					go u.orderStatusRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID, false, ctx, getOrderStatusSODetailResultChan)
					getOrderStatusSODetailResult := <-getOrderStatusSODetailResultChan

					if getOrderStatusSODetailResult.Error != nil {
						return &models.DeliveryOrderDetails{}, getOrderStatusSODetailResult.ErrorLog
					}

					getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID = getOrderStatusSODetailResult.OrderStatus.ID
					getSalesOrderDetailResult.SalesOrderDetail.OrderStatusName = getOrderStatusSODetailResult.OrderStatus.Name
					getSalesOrderDetailResult.SalesOrderDetail.OrderStatus = getOrderStatusSODetailResult.OrderStatus

					updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
					go u.salesOrderDetailRepository.UpdateByID(getSalesOrderDetailResult.SalesOrderDetail.ID, getSalesOrderDetailResult.SalesOrderDetail, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
					updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

					if updateSalesOrderDetailResult.Error != nil {
						return &models.DeliveryOrderDetails{}, updateSalesOrderDetailResult.ErrorLog
					}
				} else {
					totalSentQty += v.Qty
				}

				deliveryOrderDetails = append(deliveryOrderDetails, v)
				deliveryOrderDetailResults = append(deliveryOrderDetailResults, &models.DeliveryOrderDetailUpdateByIDRequest{
					Qty:  v.Qty,
					Note: v.Note.String,
				})
			}
		}
	}

	deliveryOrder.DeliveryOrderDetails = deliveryOrderDetails

	if totalSentQty == 0 {
		getSalesOrderResult.SalesOrder.OrderStatusID = 5
	} else if totalSentQty == totalQty {
		getSalesOrderResult.SalesOrder.OrderStatusID = 8
	} else {
		getSalesOrderResult.SalesOrder.OrderStatusID = 7
	}

	getOrderStatusSOResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByID(getSalesOrderResult.SalesOrder.OrderStatusID, false, ctx, getOrderStatusSOResultChan)
	getOrderStatusSODetailResult := <-getOrderStatusSOResultChan

	if getOrderStatusSODetailResult.Error != nil {
		return &models.DeliveryOrderDetails{}, getOrderStatusSODetailResult.ErrorLog
	}
	getSalesOrderResult.SalesOrder.OrderStatus = getOrderStatusSODetailResult.OrderStatus
	getSalesOrderResult.SalesOrder.OrderStatusName = getOrderStatusSODetailResult.OrderStatus.Name
	getSalesOrderResult.SalesOrder.SoDate = ""
	getSalesOrderResult.SalesOrder.SoRefDate = models.NullString{}
	getSalesOrderResult.SalesOrder.UpdatedAt = &now
	deliveryOrder.SalesOrder = getSalesOrderResult.SalesOrder

	updateSalesOrderChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.UpdateByID(getSalesOrderResult.SalesOrder.ID, getSalesOrderResult.SalesOrder, sqlTransaction, u.ctx, updateSalesOrderChan)
	updateSalesOrderResult := <-updateSalesOrderChan

	if updateSalesOrderResult.ErrorLog != nil {
		err := sqlTransaction.Rollback()
		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return &models.DeliveryOrderDetails{}, errorLogData
		}
		return &models.DeliveryOrderDetails{}, updateSalesOrderResult.ErrorLog
	}

	deliveryOrderLog := &models.DeliveryOrderLog{
		RequestID: "",
		DoCode:    deliveryOrder.DoCode,
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
		return &models.DeliveryOrderDetails{}, errorLogData
	}

	deliveryOrderJourney := &models.DeliveryOrderJourney{
		DoId:      deliveryOrder.ID,
		DoCode:    deliveryOrder.DoCode,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		Remark:    "",
		Reason:    "",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	createDeliveryOrderJourneyChan := make(chan *models.DeliveryOrderJourneyChan)
	go u.deliveryOrderLogRepository.InsertJourney(deliveryOrderJourney, ctx, createDeliveryOrderJourneyChan)
	createDeliveryOrderJourneysResult := <-createDeliveryOrderJourneyChan

	if createDeliveryOrderJourneysResult.Error != nil {
		return &models.DeliveryOrderDetails{}, createDeliveryOrderJourneysResult.ErrorLog
	}

	keyKafka := []byte(deliveryOrder.DoCode)
	messageKafka, _ := json.Marshal(deliveryOrder)
	err := u.kafkaClient.WriteToTopic(constants.UPDATE_DELIVERY_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.DeliveryOrderDetails{}, errorLogData
	}

	result := &models.DeliveryOrderDetails{
		Total:                0,
		DeliveryOrderDetails: deliveryOrderDetails,
	}

	return result, nil
}

func (u *deliveryOrderUseCase) Get(request *models.DeliveryOrderRequest) (*models.DeliveryOrdersOpenSearchResponse, *model.ErrorLog) {
	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
	go u.deliveryOrderOpenSearchRepository.Get(request, false, getDeliveryOrdersResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrdersResultChan

	if getDeliveryOrdersResult.Error != nil {
		return &models.DeliveryOrdersOpenSearchResponse{}, getDeliveryOrdersResult.ErrorLog
	}

	deliveryOrderResults := []*models.DeliveryOrderOpenSearchResponse{}
	for _, v := range getDeliveryOrdersResult.DeliveryOrders {
		deliveryOrder := models.DeliveryOrderOpenSearchResponse{}
		deliveryOrder.DeliveryOrderOpenSearchResponseMap(v)

		deliveryOrderResults = append(deliveryOrderResults, &deliveryOrder)

		deliveryOrderDetails := []*models.DeliveryOrderDetailOpenSearchDetailResponse{}
		for _, x := range v.DeliveryOrderDetails {
			deliveryOrderDetail := models.DeliveryOrderDetailOpenSearchDetailResponse{}
			deliveryOrderDetail.DeliveryOrderDetailOpenSearchResponseMap(x)

			deliveryOrderDetails = append(deliveryOrderDetails, &deliveryOrderDetail)
		}
		deliveryOrder.DeliveryOrderDetail = deliveryOrderDetails
	}

	deliveryOrders := &models.DeliveryOrdersOpenSearchResponse{
		DeliveryOrders: deliveryOrderResults,
		Total:          getDeliveryOrdersResult.Total,
	}

	return deliveryOrders, &model.ErrorLog{}
}

func (u *deliveryOrderUseCase) Export(request *models.DeliveryOrderExportRequest) *model.ErrorLog {

	keyKafka := []byte(uuid.New().String())
	messageKafka, _ := json.Marshal(request)
	err := u.kafkaClient.WriteToTopic(constants.EXPORT_DELIVERY_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}

func (u *deliveryOrderUseCase) GetDetails(request *models.DeliveryOrderDetailOpenSearchRequest) (*models.DeliveryOrderDetailsOpenSearchResponses, *model.ErrorLog) {
	getDeliveryOrderDetailsResultChan := make(chan *models.DeliveryOrderDetailsOpenSearchChan)
	go u.deliveryOrderDetailOpenSearchRepository.Get(request, getDeliveryOrderDetailsResultChan)
	getDeliveryOrderDetailsResult := <-getDeliveryOrderDetailsResultChan

	if getDeliveryOrderDetailsResult.Error != nil {
		return &models.DeliveryOrderDetailsOpenSearchResponses{}, getDeliveryOrderDetailsResult.ErrorLog
	}

	deliveryOrderDetails := []*models.DeliveryOrderDetailsOpenSearchResponse{}
	for _, v := range getDeliveryOrderDetailsResult.DeliveryOrderDetailOpenSearch {
		deliveryOrderDetail := models.DeliveryOrderDetailsOpenSearchResponse{}
		deliveryOrderDetail.DeliveryOrderDetailsOpenSearchResponseMap(v)

		deliveryOrderDetails = append(deliveryOrderDetails, &deliveryOrderDetail)
	}

	deliveryOrderDetailsResult := &models.DeliveryOrderDetailsOpenSearchResponses{
		DeliveryOrderDetails: deliveryOrderDetails,
		Total:                getDeliveryOrderDetailsResult.Total,
	}

	return deliveryOrderDetailsResult, &model.ErrorLog{}
}

func (u *deliveryOrderUseCase) GetDetailsByDoID(request *models.DeliveryOrderDetailRequest) (*models.DeliveryOrderOpenSearchResponse, *model.ErrorLog) {
	getDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.GetDetailsByDoID(request, getDeliveryOrderResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrderResultChan

	if getDeliveryOrdersResult.Error != nil {
		return &models.DeliveryOrderOpenSearchResponse{}, getDeliveryOrdersResult.ErrorLog
	}

	deliveryOrder := models.DeliveryOrderOpenSearchResponse{}
	deliveryOrder.DeliveryOrderOpenSearchResponseMap(getDeliveryOrdersResult.DeliveryOrder)

	deliveryOrderDetails := []*models.DeliveryOrderDetailOpenSearchDetailResponse{}
	for _, v := range getDeliveryOrdersResult.DeliveryOrder.DeliveryOrderDetails {
		deliveryOrderDetail := models.DeliveryOrderDetailOpenSearchDetailResponse{}
		deliveryOrderDetail.DeliveryOrderDetailOpenSearchResponseMap(v)

		deliveryOrderDetails = append(deliveryOrderDetails, &deliveryOrderDetail)
	}
	deliveryOrder.DeliveryOrderDetail = deliveryOrderDetails

	return &deliveryOrder, &model.ErrorLog{}
}

func (u *deliveryOrderUseCase) GetDetailByID(doDetailID int, doID int) (*models.DeliveryOrderDetailsOpenSearchResponse, *model.ErrorLog) {
	getDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.GetDetailByID(doDetailID, getDeliveryOrderResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrderResultChan

	if getDeliveryOrdersResult.Error != nil {
		return &models.DeliveryOrderDetailsOpenSearchResponse{}, getDeliveryOrdersResult.ErrorLog
	}

	if doID != getDeliveryOrdersResult.DeliveryOrder.ID {
		errorLogData := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       "Ada kesalahan pada request data, silahkan dicek kembali",
			SystemMessage: helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_GET, fmt.Sprintf("DO Detail Id %d tidak terdaftar di DO Id %d", doDetailID, doID)),
			StatusCode:    http.StatusBadRequest,
			Err:           fmt.Errorf("invalid Process"),
		})
		return &models.DeliveryOrderDetailsOpenSearchResponse{}, errorLogData
	}

	deliveryOrderDetail := &models.DeliveryOrderDetailsOpenSearchResponse{}
	for _, v := range getDeliveryOrdersResult.DeliveryOrder.DeliveryOrderDetails {
		if v.ID == doDetailID {
			deliveryOrderDetail.DeliveryOrderDetailsByDoIDOpenSearchResponseMap(v)
		}
	}

	return deliveryOrderDetail, &model.ErrorLog{}
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

func (u *deliveryOrderUseCase) GetSyncToKafkaHistories(request *models.DeliveryOrderEventLogRequest, ctx context.Context) ([]*models.DeliveryOrderEventLogResponse, *model.ErrorLog) {
	getDeliveryOrderLogResultChan := make(chan *models.GetDeliveryOrderLogsChan)
	go u.deliveryOrderLogRepository.Get(request, false, ctx, getDeliveryOrderLogResultChan)
	getDeliveryOrderLogResult := <-getDeliveryOrderLogResultChan

	if getDeliveryOrderLogResult.Error != nil {
		return []*models.DeliveryOrderEventLogResponse{}, getDeliveryOrderLogResult.ErrorLog
	}

	deliveryOrderEventLogs := []*models.DeliveryOrderEventLogResponse{}
	for _, v := range getDeliveryOrderLogResult.DeliveryOrderLog {
		deliveryOrderEventLog := models.DeliveryOrderEventLogResponse{}
		deliveryOrderEventLog.DeliveryOrderEventLogResponseMap(v)

		dataDOEventLog := models.DataDOEventLogResponse{}
		dataDOEventLog.DataDOEventLogResponseMap(v)

		deliveryOrderEventLog.Data = &dataDOEventLog

		for _, x := range v.Data.DeliveryOrderDetails {
			doDetailEventLog := models.DODetailEventLogResponse{}
			doDetailEventLog.DoDetailEventLogResponse(x)

			dataDOEventLog.DeliveryOrderDetails = append(dataDOEventLog.DeliveryOrderDetails, &doDetailEventLog)
		}

		deliveryOrderEventLogs = append(deliveryOrderEventLogs, &deliveryOrderEventLog)
	}

	return deliveryOrderEventLogs, nil
}

func (u *deliveryOrderUseCase) GetDOJourneys(request *models.DeliveryOrderJourneysRequest, ctx context.Context) (*models.DeliveryOrderJourneysResponses, *model.ErrorLog) {
	getJourneysResultChan := make(chan *models.DeliveryOrderJourneysChan)
	go u.deliveryOrderJourneysRepository.Get(request, false, ctx, getJourneysResultChan)
	getJourneysResult := <-getJourneysResultChan

	if getJourneysResult.Error != nil {
		return &models.DeliveryOrderJourneysResponses{}, getJourneysResult.ErrorLog
	}

	deliveryOrderJourneys := []*models.DeliveryOrderJourneysResponse{}
	for _, v := range getJourneysResult.DeliveryOrderJourney {
		deliveryOrderJourney := models.DeliveryOrderJourneysResponse{}
		deliveryOrderJourney.DeliveryOrderJourneyResponseMap(v)

		deliveryOrderJourneys = append(deliveryOrderJourneys, &deliveryOrderJourney)
	}

	deliveryOrderJourneyResult := models.DeliveryOrderJourneysResponses{
		DeliveryOrderJourneys: deliveryOrderJourneys,
		Total:                 getJourneysResult.Total,
	}

	return &deliveryOrderJourneyResult, nil
}

func (u *deliveryOrderUseCase) GetDOJourneysByDoID(doId int, ctx context.Context) (*models.DeliveryOrderJourneysResponses, *model.ErrorLog) {
	getDeliveryOrderJourneysResultChan := make(chan *models.DeliveryOrderJourneysChan)
	go u.deliveryOrderJourneysRepository.GetByDoID(doId, false, ctx, getDeliveryOrderJourneysResultChan)
	getDeliveryOrderJourneysResult := <-getDeliveryOrderJourneysResultChan

	if getDeliveryOrderJourneysResult.Error != nil {
		return &models.DeliveryOrderJourneysResponses{}, getDeliveryOrderJourneysResult.ErrorLog
	}

	deliveryOrderJourneys := []*models.DeliveryOrderJourneysResponse{}
	for _, v := range getDeliveryOrderJourneysResult.DeliveryOrderJourney {
		deliveryOrderJourney := models.DeliveryOrderJourneysResponse{}
		deliveryOrderJourney.DeliveryOrderJourneyResponseMap(v)

		deliveryOrderJourneys = append(deliveryOrderJourneys, &deliveryOrderJourney)
	}

	deliveryOrderJourneysResult := &models.DeliveryOrderJourneysResponses{
		DeliveryOrderJourneys: deliveryOrderJourneys,
		Total:                 getDeliveryOrderJourneysResult.Total,
	}

	return deliveryOrderJourneysResult, nil
}

func (u deliveryOrderUseCase) DeleteByID(id int, sqlTransaction *sql.Tx) *model.ErrorLog {
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
	for _, v := range getDeliveryOrderDetailsByIDResult.DeliveryOrderDetails {

		getSalesOrderDetailByIDResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.GetByID(v.SoDetailID, false, u.ctx, getSalesOrderDetailByIDResultChan)
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
			return deleteDeliveryOrderDetailResult.ErrorLog
		}

		updateSalesOrderDetailChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.UpdateByID(v.SoDetailID, getSalesOrderDetailsByIDResult.SalesOrderDetail, sqlTransaction, u.ctx, updateSalesOrderDetailChan)
		updateSalesOrderDetailResult := <-updateSalesOrderDetailChan

		if updateSalesOrderDetailResult.ErrorLog != nil {
			return updateSalesOrderDetailResult.ErrorLog
		}
	}

	deleteDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderRepository.DeleteByID(getDeliveryOrderByIDResult.DeliveryOrder, u.ctx, deleteDeliveryOrderResultChan)
	deleteDeliveryOrderResult := <-deleteDeliveryOrderResultChan

	if deleteDeliveryOrderResult.ErrorLog != nil {
		return deleteDeliveryOrderResult.ErrorLog
	}
	if totalSentQty > 0 {
		getSalesOrderByIDResult.SalesOrder.OrderStatusID = 7
	} else {
		getSalesOrderByIDResult.SalesOrder.OrderStatusID = 5
	}
	updateSalesOrderChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.UpdateByID(getSalesOrderByIDResult.SalesOrder.ID, getSalesOrderByIDResult.SalesOrder, sqlTransaction, u.ctx, updateSalesOrderChan)
	updateSalesOrderResult := <-updateSalesOrderChan

	deiveryOrderLog := &models.DeliveryOrderLog{
		RequestID: "",
		DoCode:    getDeliveryOrderByIDResult.DeliveryOrder.DoCode,
		Data:      getDeliveryOrderByIDResult.DeliveryOrder,
		Action:    constants.LOG_ACTION_MONGO_DELETE,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		CreatedAt: &now,
	}
	createDeliveryOrderLogResultChan := make(chan *models.DeliveryOrderLogChan)
	go u.deliveryOrderLogRepository.Insert(deiveryOrderLog, u.ctx, createDeliveryOrderLogResultChan)
	createDeliveryOrderLogResult := <-createDeliveryOrderLogResultChan

	if createDeliveryOrderLogResult.Error != nil {
		return createDeliveryOrderLogResult.ErrorLog
	}

	deliveryOrderJourney := &models.DeliveryOrderJourney{
		DoId:      getDeliveryOrderByIDResult.DeliveryOrder.ID,
		DoCode:    getDeliveryOrderByIDResult.DeliveryOrder.DoCode,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		Remark:    "",
		Reason:    "",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	createDeliveryOrderJourneyChan := make(chan *models.DeliveryOrderJourneyChan)
	go u.deliveryOrderLogRepository.InsertJourney(deliveryOrderJourney, u.ctx, createDeliveryOrderJourneyChan)
	createDeliveryOrderJourneysResult := <-createDeliveryOrderJourneyChan

	if createDeliveryOrderJourneysResult.Error != nil {
		return createDeliveryOrderJourneysResult.ErrorLog
	}

	keyKafka := []byte(getDeliveryOrderByIDResult.DeliveryOrder.DoCode)
	messageKafka, _ := json.Marshal(
		&models.DeliveryOrder{
			ID:        id,
			DoCode:    deiveryOrderLog.DoCode,
			UpdatedAt: getDeliveryOrderByIDResult.DeliveryOrder.UpdatedAt,
			DeletedAt: getDeliveryOrderByIDResult.DeliveryOrder.DeletedAt,
		},
	)
	err := u.kafkaClient.WriteToTopic(constants.DELETE_DELIVERY_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	if updateSalesOrderResult.ErrorLog != nil {
		return updateSalesOrderResult.ErrorLog
	}

	return nil
}

func (u deliveryOrderUseCase) DeleteDetailByID(id int, sqlTransaction *sql.Tx) *model.ErrorLog {
	now := time.Now()
	getDeliveryOrderDetailByIDResultChan := make(chan *models.DeliveryOrderDetailChan)
	go u.deliveryOrderDetailRepository.GetByID(id, false, u.ctx, getDeliveryOrderDetailByIDResultChan)
	getDeliveryOrderDetailByIDResult := <-getDeliveryOrderDetailByIDResultChan

	if getDeliveryOrderDetailByIDResult.Error != nil {
		return getDeliveryOrderDetailByIDResult.ErrorLog
	}

	getDeliveryOrderByIDResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderRepository.GetByID(getDeliveryOrderDetailByIDResult.DeliveryOrderDetail.DeliveryOrderID, false, u.ctx, getDeliveryOrderByIDResultChan)
	getDeliveryOrderByIDResult := <-getDeliveryOrderByIDResultChan

	if getDeliveryOrderByIDResult.Error != nil {
		return getDeliveryOrderByIDResult.ErrorLog
	}

	getSalesOrderDetailByIDResultChan := make(chan *models.SalesOrderDetailChan)
	go u.salesOrderDetailRepository.GetByID(getDeliveryOrderDetailByIDResult.DeliveryOrderDetail.SoDetailID, false, u.ctx, getSalesOrderDetailByIDResultChan)
	getSalesOrderDetailsByIDResult := <-getSalesOrderDetailByIDResultChan

	if getSalesOrderDetailsByIDResult.Error != nil {
		return getSalesOrderDetailsByIDResult.ErrorLog
	}

	getSalesOrderDetailsByIDResult.SalesOrderDetail.SentQty -= getDeliveryOrderDetailByIDResult.DeliveryOrderDetail.Qty
	getSalesOrderDetailsByIDResult.SalesOrderDetail.ResidualQty += getDeliveryOrderDetailByIDResult.DeliveryOrderDetail.Qty
	getSalesOrderDetailsByIDResult.SalesOrderDetail.UpdatedAt = &now

	updateSalesOrderDetailChan := make(chan *models.SalesOrderDetailChan)
	go u.salesOrderDetailRepository.UpdateByID(getDeliveryOrderDetailByIDResult.DeliveryOrderDetail.SoDetailID, getSalesOrderDetailsByIDResult.SalesOrderDetail, sqlTransaction, u.ctx, updateSalesOrderDetailChan)
	updateSalesOrderDetailResult := <-updateSalesOrderDetailChan

	if updateSalesOrderDetailResult.ErrorLog != nil {
		return updateSalesOrderDetailResult.ErrorLog
	}

	deleteDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
	go u.deliveryOrderDetailRepository.DeleteByID(getDeliveryOrderDetailByIDResult.DeliveryOrderDetail, sqlTransaction, u.ctx, deleteDeliveryOrderDetailResultChan)
	deleteDeliveryOrderDetailResult := <-deleteDeliveryOrderDetailResultChan

	if deleteDeliveryOrderDetailResult.ErrorLog != nil {
		return deleteDeliveryOrderDetailResult.ErrorLog
	}

	doDetailLogData := &models.DeliveryOrderDetailLogData{}
	doDetailLogData.DoDetailMap(getDeliveryOrderByIDResult.DeliveryOrder, getDeliveryOrderDetailByIDResult.DeliveryOrderDetail)

	deiveryOrderLog := &models.DeliveryOrderLog{
		RequestID: "",
		DoCode:    getDeliveryOrderByIDResult.DeliveryOrder.DoCode,
		Data:      doDetailLogData,
		Action:    constants.LOG_ACTION_MONGO_DELETE,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		CreatedAt: &now,
	}
	createDeliveryOrderLogResultChan := make(chan *models.DeliveryOrderLogChan)
	go u.deliveryOrderLogRepository.Insert(deiveryOrderLog, u.ctx, createDeliveryOrderLogResultChan)
	createDeliveryOrderLogResult := <-createDeliveryOrderLogResultChan

	if createDeliveryOrderLogResult.Error != nil {
		return createDeliveryOrderLogResult.ErrorLog
	}

	deliveryOrderJourney := &models.DeliveryOrderJourney{
		DoId:      getDeliveryOrderByIDResult.DeliveryOrder.ID,
		DoCode:    getDeliveryOrderByIDResult.DeliveryOrder.DoCode,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		Remark:    "",
		Reason:    "",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	createDeliveryOrderJourneyChan := make(chan *models.DeliveryOrderJourneyChan)
	go u.deliveryOrderLogRepository.InsertJourney(deliveryOrderJourney, u.ctx, createDeliveryOrderJourneyChan)
	createDeliveryOrderJourneysResult := <-createDeliveryOrderJourneyChan

	if createDeliveryOrderJourneysResult.Error != nil {
		return createDeliveryOrderJourneysResult.ErrorLog
	}

	keyKafka := []byte(getDeliveryOrderDetailByIDResult.DeliveryOrderDetail.DoDetailCode)
	messageKafka, _ := json.Marshal(
		&models.DeliveryOrderDetail{
			ID:              id,
			DeliveryOrderID: getDeliveryOrderByIDResult.DeliveryOrder.ID,
			DoDetailCode:    doDetailLogData.DoDetailCode,
			UpdatedAt:       getDeliveryOrderDetailByIDResult.DeliveryOrderDetail.UpdatedAt,
			DeletedAt:       getDeliveryOrderDetailByIDResult.DeliveryOrderDetail.DeletedAt,
		},
	)
	err := u.kafkaClient.WriteToTopic(constants.DELETE_DELIVERY_ORDER_DETAIL_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}

func (u deliveryOrderUseCase) DeleteDetailByDoID(id int, sqlTransaction *sql.Tx) *model.ErrorLog {
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

	getSalesOrderByIDResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.GetByID(getDeliveryOrderByIDResult.DeliveryOrder.SalesOrderID, false, u.ctx, getSalesOrderByIDResultChan)
	getSalesOrderByIDResult := <-getSalesOrderByIDResultChan

	if getSalesOrderByIDResult.Error != nil {
		return getSalesOrderByIDResult.ErrorLog
	}
	totalSentQty := 0
	deliveryOrderDetails := []*models.DeliveryOrderDetail{}
	for _, v := range getDeliveryOrderDetailsByIDResult.DeliveryOrderDetails {

		getSalesOrderDetailByIDResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.GetByID(v.SoDetailID, false, u.ctx, getSalesOrderDetailByIDResultChan)
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
			return deleteDeliveryOrderDetailResult.ErrorLog
		}

		updateSalesOrderDetailChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.UpdateByID(v.SoDetailID, getSalesOrderDetailsByIDResult.SalesOrderDetail, sqlTransaction, u.ctx, updateSalesOrderDetailChan)
		updateSalesOrderDetailResult := <-updateSalesOrderDetailChan

		if updateSalesOrderDetailResult.ErrorLog != nil {
			return updateSalesOrderDetailResult.ErrorLog
		}
		deliveryOrderDetails = append(deliveryOrderDetails, &models.DeliveryOrderDetail{ID: v.ID})
	}
	if totalSentQty > 0 {
		getSalesOrderByIDResult.SalesOrder.OrderStatusID = 7
	} else {
		getSalesOrderByIDResult.SalesOrder.OrderStatusID = 5
	}
	updateSalesOrderChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.UpdateByID(getSalesOrderByIDResult.SalesOrder.ID, getSalesOrderByIDResult.SalesOrder, sqlTransaction, u.ctx, updateSalesOrderChan)
	updateSalesOrderResult := <-updateSalesOrderChan

	deiveryOrderLog := &models.DeliveryOrderLog{
		RequestID: "",
		DoCode:    getDeliveryOrderByIDResult.DeliveryOrder.DoCode,
		Data:      getDeliveryOrderByIDResult.DeliveryOrder,
		Action:    constants.LOG_ACTION_MONGO_DELETE,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		CreatedAt: &now,
	}
	createDeliveryOrderLogResultChan := make(chan *models.DeliveryOrderLogChan)
	go u.deliveryOrderLogRepository.Insert(deiveryOrderLog, u.ctx, createDeliveryOrderLogResultChan)
	createDeliveryOrderLogResult := <-createDeliveryOrderLogResultChan

	if createDeliveryOrderLogResult.Error != nil {
		return createDeliveryOrderLogResult.ErrorLog
	}

	deliveryOrderJourney := &models.DeliveryOrderJourney{
		DoId:      getDeliveryOrderByIDResult.DeliveryOrder.ID,
		DoCode:    getDeliveryOrderByIDResult.DeliveryOrder.DoCode,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		Remark:    "",
		Reason:    "",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	createDeliveryOrderJourneyChan := make(chan *models.DeliveryOrderJourneyChan)
	go u.deliveryOrderLogRepository.InsertJourney(deliveryOrderJourney, u.ctx, createDeliveryOrderJourneyChan)
	createDeliveryOrderJourneysResult := <-createDeliveryOrderJourneyChan

	if createDeliveryOrderJourneysResult.Error != nil {
		return createDeliveryOrderJourneysResult.ErrorLog
	}

	keyKafka := []byte(getDeliveryOrderByIDResult.DeliveryOrder.DoCode)
	messageKafka, _ := json.Marshal(
		&models.DeliveryOrder{
			ID:                   id,
			DeliveryOrderDetails: deliveryOrderDetails,
			DoCode:               deiveryOrderLog.DoCode,
			UpdatedAt:            getDeliveryOrderByIDResult.DeliveryOrder.UpdatedAt,
			DeletedAt:            getDeliveryOrderByIDResult.DeliveryOrder.DeletedAt,
		},
	)
	err := u.kafkaClient.WriteToTopic(constants.DELETE_DELIVERY_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	if updateSalesOrderResult.ErrorLog != nil {
		return updateSalesOrderResult.ErrorLog
	}

	return nil
}
