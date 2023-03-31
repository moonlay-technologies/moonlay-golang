package consumer

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/app/repositories"
	mongoRepositories "order-service/app/repositories/mongod"
	"order-service/global/utils/helper"
	kafkadbo "order-service/global/utils/kafka"
	"time"

	"github.com/bxcodec/dbresolver"
)

type UploadDOItemConsumerHandlerInterface interface {
	ProcessMessage()
}

type uploadDOItemConsumerHandler struct {
	deliveryOrderRepository       repositories.DeliveryOrderRepositoryInterface
	deliveryOrderDetailRepository repositories.DeliveryOrderDetailRepositoryInterface
	salesOrderRepository          repositories.SalesOrderRepositoryInterface
	salesOrderDetailRepository    repositories.SalesOrderDetailRepositoryInterface
	orderStatusRepository         repositories.OrderStatusRepositoryInterface
	orderSourceRepository         repositories.OrderSourceRepositoryInterface
	warehouseRepository           repositories.WarehouseRepositoryInterface
	brandRepository               repositories.BrandRepositoryInterface
	uomRepository                 repositories.UomRepositoryInterface
	agentRepository               repositories.AgentRepositoryInterface
	storeRepository               repositories.StoreRepositoryInterface
	productRepository             repositories.ProductRepositoryInterface
	userRepository                repositories.UserRepositoryInterface
	salesmanRepository            repositories.SalesmanRepositoryInterface
	deliveryOrderLogRepository    mongoRepositories.DeliveryOrderLogRepositoryInterface
	salesOrderLogRepository       mongoRepositories.SalesOrderLogRepositoryInterface
	salesOrderJourneyRepository   mongoRepositories.SalesOrderJourneysRepositoryInterface
	kafkaClient                   kafkadbo.KafkaClientInterface
	ctx                           context.Context
	args                          []interface{}
	db                            dbresolver.DB
}

func InitUploadDOItemConsumerHandlerInterface(deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, deliveryOrderDetailRepository repositories.DeliveryOrderDetailRepositoryInterface, salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, orderSourceRepository repositories.OrderSourceRepositoryInterface, warehouseRepository repositories.WarehouseRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, uomRepository repositories.UomRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, userRepository repositories.UserRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, deliveryOrderLogRepository mongoRepositories.DeliveryOrderLogRepositoryInterface, salesOrderLogRepository mongoRepositories.SalesOrderLogRepositoryInterface, salesOrderJourneyRepository mongoRepositories.SalesOrderJourneysRepositoryInterface, kafkaClient kafkadbo.KafkaClientInterface, ctx context.Context, args []interface{}, db dbresolver.DB) UploadDOItemConsumerHandlerInterface {
	return &uploadDOItemConsumerHandler{
		deliveryOrderRepository:       deliveryOrderRepository,
		deliveryOrderDetailRepository: deliveryOrderDetailRepository,
		salesOrderRepository:          salesOrderRepository,
		salesOrderDetailRepository:    salesOrderDetailRepository,
		orderStatusRepository:         orderStatusRepository,
		orderSourceRepository:         orderSourceRepository,
		warehouseRepository:           warehouseRepository,
		brandRepository:               brandRepository,
		uomRepository:                 uomRepository,
		productRepository:             productRepository,
		userRepository:                userRepository,
		salesmanRepository:            salesmanRepository,
		agentRepository:               agentRepository,
		storeRepository:               storeRepository,
		deliveryOrderLogRepository:    deliveryOrderLogRepository,
		salesOrderLogRepository:       salesOrderLogRepository,
		salesOrderJourneyRepository:   salesOrderJourneyRepository,
		kafkaClient:                   kafkaClient,
		ctx:                           ctx,
		args:                          args,
		db:                            db,
	}
}

// func (c *uploadDOItemConsumerHandler) ProcessMessage() {
// 	fmt.Println("process ", constants.UPLOAD_DO_ITEM_TOPIC)
// 	topic := c.args[1].(string)
// 	groupID := c.args[2].(string)
// 	reader := c.kafkaClient.SetConsumerGroupReader(topic, groupID)

// 	for {
// 		m, err := reader.ReadMessage(c.ctx)
// 		if err != nil {
// 			break
// 		}

// 		fmt.Printf("message do at topic/partition/offset %v/%v/%v \n", m.Topic, m.Partition, m.Offset)
// 		now := time.Now()

// 		var UploadDOFields []*models.UploadDOField
// 		err = json.Unmarshal(m.Value, &UploadDOFields)

// 		requestId := string(m.Key[:])

// 		// Get Order Status for DO
// 		getOrderStatusResultChan := make(chan *models.OrderStatusChan)
// 		go c.orderStatusRepository.GetByNameAndType("open", "delivery_order", false, c.ctx, getOrderStatusResultChan)
// 		getOrderStatusResult := <-getOrderStatusResultChan

// 		// Get Order Source for DO
// 		getOrderSourceResultChan := make(chan *models.OrderSourceChan)
// 		go c.orderSourceRepository.GetBySourceName("upload_sj", false, c.ctx, getOrderSourceResultChan)
// 		getOrderSourceResult := <-getOrderSourceResultChan

// 		for _, v := range UploadDOFields {

// 			var errors []string

// 			if getOrderStatusResult.Error != nil {
// 				fmt.Println(getOrderStatusResult.Error.Error())
// 				errors = append(errors, getOrderStatusResult.Error.Error())
// 			}

// 			if getOrderSourceResult.Error != nil {
// 				fmt.Println(getOrderSourceResult.Error.Error())
// 				errors = append(errors, getOrderSourceResult.Error.Error())
// 			}

// 			deliveryOrder := &models.DeliveryOrder{}

// 			// Get Sales Order By SoCode / NoOrder
// 			getSalesOrderResultChan := make(chan *models.SalesOrderChan)
// 			go c.salesOrderRepository.GetByCode(v.NoOrder, false, c.ctx, getSalesOrderResultChan)
// 			getSalesOrderResult := <-getSalesOrderResultChan

// 			if getSalesOrderResult.Error != nil {
// 				fmt.Println(getSalesOrderResult.Error.Error())
// 				errors = append(errors, getOrderSourceResult.Error.Error())
// 			}

// 			// Get Sales Order Details by SOID, Sku and uomCode (upload data)
// 			getSODetailBySoIdSkuAndUomCodeResultChan := make(chan *models.SalesOrderDetailChan)
// 			go c.salesOrderDetailRepository.GetBySOIDSkuAndUomCode(getSalesOrderResult.SalesOrder.ID, v.KodeProduk, v.Unit, false, c.ctx, getSODetailBySoIdSkuAndUomCodeResultChan)
// 			getSODetailBySoIdSkuAndUomCodeResult := <-getSODetailBySoIdSkuAndUomCodeResultChan

// 			if getSODetailBySoIdSkuAndUomCodeResult.Error != nil {
// 				fmt.Println(getSODetailBySoIdSkuAndUomCodeResult.Error.Error())
// 				errors = append(errors, getOrderSourceResult.Error.Error())
// 			}

// 			// Get Brand by ID / KodeMerk
// 			getBrandResultChan := make(chan *models.BrandChan)
// 			go c.brandRepository.GetByID(v.KodeMerk, false, c.ctx, getBrandResultChan)
// 			getBrandResult := <-getBrandResultChan

// 			if getBrandResult.Error != nil {
// 				fmt.Println(getBrandResult.Error.Error())
// 				errors = append(errors, getOrderSourceResult.Error.Error())
// 			}

// 			// Get Agent By ID / IDDistributor
// 			getAgentResultChan := make(chan *models.AgentChan)
// 			go c.agentRepository.GetByID(v.IDDistributor, false, c.ctx, getAgentResultChan)
// 			getAgentResult := <-getAgentResultChan

// 			if getAgentResult.Error != nil {
// 				fmt.Println(getAgentResult.Error.Error())
// 				errors = append(errors, getOrderSourceResult.Error.Error())
// 			}

// 			// Get Warehouse
// 			getWarehouseResultChan := make(chan *models.WarehouseChan)
// 			if v.KodeGudang == "" {
// 				go c.warehouseRepository.GetByAgentID(getAgentResult.Agent.ID, false, c.ctx, getWarehouseResultChan)
// 			} else {
// 				go c.warehouseRepository.GetByCode(v.KodeGudang, false, c.ctx, getWarehouseResultChan)
// 			}
// 			getWarehouseResult := <-getWarehouseResultChan

// 			if getWarehouseResult.Error != nil {
// 				fmt.Println(getWarehouseResult.Error.Error())
// 				errors = append(errors, getOrderSourceResult.Error.Error())
// 			}

// 			// Get Store By ID
// 			getStoreResultChan := make(chan *models.StoreChan)
// 			go c.storeRepository.GetByID(getSalesOrderResult.SalesOrder.StoreID, false, c.ctx, getStoreResultChan)
// 			getStoreResult := <-getStoreResultChan

// 			if getStoreResult.Error != nil {
// 				fmt.Println(getStoreResult.Error.Error())
// 				errors = append(errors, getOrderSourceResult.Error.Error())
// 			}

// 			// Get User By ID
// 			getUserResultChan := make(chan *models.UserChan)
// 			go c.userRepository.GetByID(getSalesOrderResult.SalesOrder.UserID, false, c.ctx, getUserResultChan)
// 			getUserResult := <-getUserResultChan

// 			if getUserResult.Error != nil {
// 				fmt.Println(getUserResult.Error.Error())
// 				errors = append(errors, getOrderSourceResult.Error.Error())
// 			}

// 			// Get Salesman
// 			getSalesmanResultChan := make(chan *models.SalesmanChan)
// 			if getSalesOrderResult.SalesOrder.SalesmanID.Int64 > 0 {
// 				go c.salesmanRepository.GetByID(int(getSalesOrderResult.SalesOrder.SalesmanID.Int64), false, c.ctx, getSalesmanResultChan)
// 			} else {
// 				go c.salesmanRepository.GetByEmail(getUserResult.User.Email, false, c.ctx, getSalesmanResultChan)
// 			}
// 			getSalesmanResult := <-getSalesmanResultChan

// 			if getSalesmanResult.Error != nil {
// 				fmt.Println(getSalesmanResult.Error.Error())
// 				errors = append(errors, getOrderSourceResult.Error.Error())
// 			}

// 			latestUpdatedBy := c.ctx.Value("user").(*models.UserClaims)
// 			deliveryOrder.SalesOrderID = getSalesOrderResult.SalesOrder.ID
// 			deliveryOrder.DoRefCode = models.NullString{NullString: sql.NullString{String: v.NoSJ, Valid: true}}
// 			deliveryOrder.DoRefDate = models.NullString{NullString: sql.NullString{String: v.TanggalSJ, Valid: true}}
// 			deliveryOrder.DriverName = models.NullString{NullString: sql.NullString{String: v.NamaSupir, Valid: true}}
// 			deliveryOrder.PlatNumber = models.NullString{NullString: sql.NullString{String: v.PlatNo, Valid: true}}
// 			deliveryOrder.Note = models.NullString{NullString: sql.NullString{String: v.Catatan, Valid: true}}
// 			deliveryOrder.IsDoneSyncToEs = "0"
// 			deliveryOrder.StartDateSyncToEs = &now
// 			deliveryOrder.EndDateSyncToEs = &now
// 			deliveryOrder.StartCreatedDate = &now
// 			deliveryOrder.EndCreatedDate = &now
// 			deliveryOrder.LatestUpdatedBy = latestUpdatedBy.UserID
// 			deliveryOrder.CreatedAt = &now
// 			deliveryOrder.UpdatedAt = &now
// 			deliveryOrder.DeletedAt = nil

// 			deliveryOrder.WarehouseChanMap(getWarehouseResult)
// 			deliveryOrder.AgentMap(getAgentResult.Agent)
// 			deliveryOrder.DoCode = helper.GenerateDOCode(getAgentResult.Agent.ID, getOrderSourceResult.OrderSource.Code)
// 			deliveryOrder.DoDate = now.Format("2006-01-02")
// 			deliveryOrder.OrderStatus = getOrderStatusResult.OrderStatus
// 			deliveryOrder.OrderStatusID = getOrderStatusResult.OrderStatus.ID
// 			deliveryOrder.OrderSource = getOrderSourceResult.OrderSource
// 			deliveryOrder.OrderSourceID = getOrderSourceResult.OrderSource.ID
// 			deliveryOrder.Store = getStoreResult.Store
// 			deliveryOrder.CreatedBy = v.IDUser
// 			deliveryOrder.SalesOrder = getSalesOrderResult.SalesOrder
// 			deliveryOrder.Brand = getBrandResult.Brand
// 			if getSalesmanResult.Salesman != nil {
// 				deliveryOrder.Salesman = getSalesmanResult.Salesman
// 			}

// 			sqlTransaction, _ := c.db.BeginTx(c.ctx, nil)

// 			// Insert to DB, table delivery_orders
// 			createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
// 			go c.deliveryOrderRepository.Insert(deliveryOrder, sqlTransaction, c.ctx, createDeliveryOrderResultChan)
// 			createDeliveryOrderResult := <-createDeliveryOrderResultChan

// 			if createDeliveryOrderResult.Error != nil {
// 				sqlTransaction.Rollback()
// 				// return
// 			}

// 			// Delivery Order Detail
// 			deliveryOrderDetails := []*models.DeliveryOrderDetail{}
// 			totalResidualQty := 0
// 			for _, x := range getSalesOrderResult.SalesOrder.SalesOrderDetails {
// 				if getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail.ID == x.ID {
// 					// Get Sales Order Detail By ID
// 					getSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
// 					go c.salesOrderDetailRepository.GetByID(x.ID, false, c.ctx, getSalesOrderDetailResultChan)
// 					getSalesOrderDetailResult := <-getSalesOrderDetailResultChan

// 					// Get SO Detail Order Status by ID
// 					getOrderStatusDetailResultChan := make(chan *models.OrderStatusChan)
// 					go c.orderStatusRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.ID, false, c.ctx, getOrderStatusDetailResultChan)
// 					getOrderStatusDetailResult := <-getOrderStatusDetailResultChan

// 					getSalesOrderDetailResult.SalesOrderDetail.UpdatedAt = &now
// 					getSalesOrderDetailResult.SalesOrderDetail.SentQty += v.QTYShip
// 					getSalesOrderDetailResult.SalesOrderDetail.ResidualQty -= v.QTYShip
// 					totalResidualQty += getSalesOrderDetailResult.SalesOrderDetail.ResidualQty
// 					statusName := "partial"

// 					if getSalesOrderDetailResult.SalesOrderDetail.ResidualQty == 0 {
// 						statusName = "closed"
// 					}

// 					// Get SO Detail Order Status By statusName
// 					getOrderStatusSODetailResultChan := make(chan *models.OrderStatusChan)
// 					go c.orderStatusRepository.GetByNameAndType(statusName, "sales_order_detail", false, c.ctx, getOrderStatusSODetailResultChan)
// 					getOrderStatusSODetailResult := <-getOrderStatusSODetailResultChan

// 					getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID = getOrderStatusSODetailResult.OrderStatus.ID
// 					getSalesOrderDetailResult.SalesOrderDetail.OrderStatusName =
// 						getOrderStatusSODetailResult.OrderStatus.Name
// 					getSalesOrderDetailResult.SalesOrderDetail.OrderStatus = getOrderStatusSODetailResult.OrderStatus

// 					// Get Product Detail by ID
// 					getProductDetailResultChan := make(chan *models.ProductChan)
// 					go c.productRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.ProductID, false, c.ctx, getProductDetailResultChan)
// 					getProductDetailResult := <-getProductDetailResultChan

// 					// Get Uom Detail By ID
// 					getUomDetailResultChan := make(chan *models.UomChan)
// 					go c.uomRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.UomID, false, c.ctx, getUomDetailResultChan)
// 					getUomDetailResult := <-getUomDetailResultChan

// 					doDetailCode, _ := helper.GenerateDODetailCode(createDeliveryOrderResult.DeliveryOrder.ID, getAgentResult.Agent.ID, getProductDetailResult.Product.ID, getUomDetailResult.Uom.ID)

// 					deliveryOrderDetail := &models.DeliveryOrderDetail{}
// 					deliveryOrderDetail.SoDetailID = x.ID
// 					deliveryOrderDetail.Qty = v.QTYShip
// 					deliveryOrderDetail.IsDoneSyncToEs = "0"
// 					deliveryOrderDetail.StartDateSyncToEs = &now
// 					deliveryOrderDetail.EndDateSyncToEs = &now
// 					deliveryOrderDetail.CreatedAt = &now
// 					deliveryOrderDetail.UpdatedAt = &now
// 					deliveryOrderDetail.DeletedAt = nil

// 					deliveryOrderDetail.DeliveryOrderID = int(createDeliveryOrderResult.ID)
// 					deliveryOrderDetail.BrandID = getBrandResult.Brand.ID
// 					deliveryOrderDetail.DoDetailCode = doDetailCode
// 					deliveryOrderDetail.Note = models.NullString{NullString: sql.NullString{String: v.Catatan, Valid: true}}
// 					deliveryOrderDetail.Uom = getUomDetailResult.Uom
// 					deliveryOrderDetail.ProductChanMap(getProductDetailResult)
// 					deliveryOrderDetail.SalesOrderDetailChanMap(getSalesOrderDetailResult)
// 					deliveryOrderDetail.OrderStatusID = deliveryOrder.OrderStatusID
// 					deliveryOrderDetail.OrderStatusName = deliveryOrder.OrderStatusName
// 					deliveryOrderDetail.OrderStatus = getOrderStatusDetailResult.OrderStatus

// 					// Insert to DB, Delivery Order Detail
// 					createDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
// 					go c.deliveryOrderDetailRepository.Insert(deliveryOrderDetail, sqlTransaction, c.ctx, createDeliveryOrderDetailResultChan)
// 					createDeliveryOrderDetailResult := <-createDeliveryOrderDetailResultChan

// 					if createDeliveryOrderDetailResult.Error != nil {
// 						sqlTransaction.Rollback()
// 					}

// 					// Update to DB, Sales Order Detail
// 					updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
// 					go c.salesOrderDetailRepository.UpdateByID(getSalesOrderDetailResult.SalesOrderDetail.ID, getSalesOrderDetailResult.SalesOrderDetail, sqlTransaction, c.ctx, updateSalesOrderDetailResultChan)
// 					updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

// 					if updateSalesOrderDetailResult.Error != nil {
// 						sqlTransaction.Rollback()
// 					}

// 					deliveryOrderDetail.ID = int(createDeliveryOrderDetailResult.ID)
// 					deliveryOrderDetails = append(deliveryOrderDetails, deliveryOrderDetail)
// 					getSalesOrderResult.SalesOrder.SalesOrderDetails = append(getSalesOrderResult.SalesOrder.SalesOrderDetails, getSalesOrderDetailResult.SalesOrderDetail)
// 				}

// 				deliveryOrder.DeliveryOrderDetails = deliveryOrderDetails
// 				if totalResidualQty == 0 {
// 					getSalesOrderResult.SalesOrder.OrderStatusID = 8
// 				} else {
// 					getSalesOrderResult.SalesOrder.OrderStatusID = 7
// 				}

// 				// Get updated Order Status Sales Order
// 				getOrderStatusSOResultChan := make(chan *models.OrderStatusChan)
// 				go c.orderStatusRepository.GetByID(getSalesOrderResult.SalesOrder.OrderStatusID, false, c.ctx, getOrderStatusSOResultChan)
// 				getOrderStatusSODetailResult := <-getOrderStatusSOResultChan

// 				getSalesOrderResult.SalesOrder.OrderStatus = getOrderStatusSODetailResult.OrderStatus
// 				getSalesOrderResult.SalesOrder.OrderStatusName = getOrderStatusSODetailResult.OrderStatus.Name
// 				getSalesOrderResult.SalesOrder.SoDate = ""
// 				getSalesOrderResult.SalesOrder.SoRefDate = models.NullString{}
// 				getSalesOrderResult.SalesOrder.UpdatedAt = &now
// 				deliveryOrder.SalesOrder = getSalesOrderResult.SalesOrder

// 				// Update to DB, Sales Order
// 				updateSalesOrderChan := make(chan *models.SalesOrderChan)
// 				go c.salesOrderRepository.UpdateByID(getSalesOrderResult.SalesOrder.ID, getSalesOrderResult.SalesOrder, sqlTransaction, c.ctx, updateSalesOrderChan)
// 				updateSalesOrderResult := <-updateSalesOrderChan

// 				if updateSalesOrderResult.Error != nil {
// 					sqlTransaction.Rollback()
// 				}
// 			}

// 			sqlTransaction.Commit()

// 			deliveryOrderLog := &models.DeliveryOrderLog{
// 				RequestID: requestId,
// 				DoCode:    deliveryOrder.DoCode,
// 				Data:      deliveryOrder,
// 				Status:    constants.LOG_STATUS_MONGO_DEFAULT,
// 				Action:    constants.LOG_ACTION_MONGO_INSERT,
// 				CreatedAt: &now,
// 			}

// 			createDeliveryOrderLogResultChan := make(chan *models.DeliveryOrderLogChan)
// 			go c.deliveryOrderLogRepository.Insert(deliveryOrderLog, c.ctx, createDeliveryOrderLogResultChan)
// 			createDeliveryOrderLogResult := <-createDeliveryOrderLogResultChan

// 			if createDeliveryOrderLogResult.Error != nil {
// 				fmt.Println(createDeliveryOrderLogResult.Error)
// 			}

// 			deliveryOrderJourney := &models.DeliveryOrderJourney{
// 				DoId:      deliveryOrder.ID,
// 				DoCode:    deliveryOrder.DoCode,
// 				Status:    constants.LOG_STATUS_MONGO_DEFAULT,
// 				Remark:    "",
// 				Reason:    "",
// 				CreatedAt: &now,
// 				UpdatedAt: &now,
// 			}

// 			createDeliveryOrderJourneyChan := make(chan *models.DeliveryOrderJourneyChan)
// 			go c.deliveryOrderLogRepository.InsertJourney(deliveryOrderJourney, c.ctx, createDeliveryOrderJourneyChan)
// 			createDeliveryOrderJourneysResult := <-createDeliveryOrderJourneyChan

// 			if createDeliveryOrderJourneysResult.Error != nil {
// 				// return &models.DeliveryOrderStoreResponse{}, createDeliveryOrderJourneysResult.ErrorLog
// 			}

// 			keyKafka := []byte(deliveryOrder.DoCode)
// 			messageKafka, _ := json.Marshal(deliveryOrder)
// 			err := c.kafkaClient.WriteToTopic(constants.CREATE_DELIVERY_ORDER_TOPIC, keyKafka, messageKafka)

// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 		}
// 	}
// }

func (c *uploadDOItemConsumerHandler) ProcessMessage() {
	fmt.Println("process ", constants.UPLOAD_DO_ITEM_TOPIC)
	topic := c.args[1].(string)
	groupID := c.args[2].(string)
	reader := c.kafkaClient.SetConsumerGroupReader(topic, groupID)

	for {
		m, err := reader.ReadMessage(c.ctx)
		if err != nil {
			break
		}

		fmt.Printf("message do at topic/partition/offset %v/%v/%v \n", m.Topic, m.Partition, m.Offset)
		now := time.Now()

		var UploadDOFields []*models.UploadDOField
		err = json.Unmarshal(m.Value, &UploadDOFields)

		requestId := string(m.Key[:])

		doRefCodes := map[string]*models.DeliveryOrder{}

		// Get Order Status for DO
		getOrderStatusResultChan := make(chan *models.OrderStatusChan)
		go c.orderStatusRepository.GetByNameAndType("open", "delivery_order", false, c.ctx, getOrderStatusResultChan)
		getOrderStatusResult := <-getOrderStatusResultChan

		// Get Order Source for DO
		getOrderSourceResultChan := make(chan *models.OrderSourceChan)
		go c.orderSourceRepository.GetBySourceName("upload_sj", false, c.ctx, getOrderSourceResultChan)
		getOrderSourceResult := <-getOrderSourceResultChan

		for _, v := range UploadDOFields {

			var errors []string

			if getOrderStatusResult.Error != nil {
				fmt.Println(getOrderStatusResult.Error.Error())
				errors = append(errors, getOrderStatusResult.Error.Error())
			}

			if getOrderSourceResult.Error != nil {
				fmt.Println(getOrderSourceResult.Error.Error())
				errors = append(errors, getOrderSourceResult.Error.Error())
			}

			// Get Sales Order By SoCode / NoOrder
			getSalesOrderResultChan := make(chan *models.SalesOrderChan)
			go c.salesOrderRepository.GetByCode(v.NoOrder, false, c.ctx, getSalesOrderResultChan)
			getSalesOrderResult := <-getSalesOrderResultChan

			if getSalesOrderResult.Error != nil {
				fmt.Println(getSalesOrderResult.Error.Error())
				errors = append(errors, getSalesOrderResult.Error.Error())
			}

			// Get Sales Order Detail By ID
			getSalesOrderDetailBySoIdResultChan := make(chan *models.SalesOrderDetailsChan)
			go c.salesOrderDetailRepository.GetBySalesOrderID(getSalesOrderResult.SalesOrder.ID, false, c.ctx, getSalesOrderDetailBySoIdResultChan)
			getSalesOrderDetailBySoIdResult := <-getSalesOrderDetailBySoIdResultChan

			if getSalesOrderDetailBySoIdResult.Error != nil {
				fmt.Println(getSalesOrderDetailBySoIdResult.Error.Error())
				errors = append(errors, getSalesOrderDetailBySoIdResult.Error.Error())
			}

			// Get Sales Order Details by SOID, Sku and uomCode (upload data)
			getSODetailBySoIdSkuAndUomCodeResultChan := make(chan *models.SalesOrderDetailChan)
			go c.salesOrderDetailRepository.GetBySOIDSkuAndUomCode(getSalesOrderResult.SalesOrder.ID, v.KodeProduk, v.Unit, false, c.ctx, getSODetailBySoIdSkuAndUomCodeResultChan)
			getSODetailBySoIdSkuAndUomCodeResult := <-getSODetailBySoIdSkuAndUomCodeResultChan

			if getSODetailBySoIdSkuAndUomCodeResult.Error != nil {
				fmt.Println(getSODetailBySoIdSkuAndUomCodeResult.Error.Error())
				errors = append(errors, getSODetailBySoIdSkuAndUomCodeResult.Error.Error())
			}

			// Get Brand by ID / KodeMerk
			getBrandResultChan := make(chan *models.BrandChan)
			go c.brandRepository.GetByID(v.KodeMerk, false, c.ctx, getBrandResultChan)
			getBrandResult := <-getBrandResultChan

			if getBrandResult.Error != nil {
				fmt.Println(getBrandResult.Error.Error())
				errors = append(errors, getBrandResult.Error.Error())
			}

			// Get Agent By ID / IDDistributor
			getAgentResultChan := make(chan *models.AgentChan)
			go c.agentRepository.GetByID(v.IDDistributor, false, c.ctx, getAgentResultChan)
			getAgentResult := <-getAgentResultChan

			if getAgentResult.Error != nil {
				fmt.Println(getAgentResult.Error.Error())
				errors = append(errors, getAgentResult.Error.Error())
			}

			// Get Warehouse
			getWarehouseResultChan := make(chan *models.WarehouseChan)
			if v.KodeGudang == "" {
				go c.warehouseRepository.GetByAgentID(getAgentResult.Agent.ID, false, c.ctx, getWarehouseResultChan)
			} else {
				go c.warehouseRepository.GetByCode(v.KodeGudang, false, c.ctx, getWarehouseResultChan)
			}
			getWarehouseResult := <-getWarehouseResultChan

			if getWarehouseResult.Error != nil {
				fmt.Println(getWarehouseResult.Error.Error())
				errors = append(errors, getWarehouseResult.Error.Error())
			}

			// Get Store By ID
			getStoreResultChan := make(chan *models.StoreChan)
			go c.storeRepository.GetByID(getSalesOrderResult.SalesOrder.StoreID, false, c.ctx, getStoreResultChan)
			getStoreResult := <-getStoreResultChan

			if getStoreResult.Error != nil {
				fmt.Println(getStoreResult.Error.Error())
				errors = append(errors, getOrderSourceResult.Error.Error())
			}

			// Get User By ID
			getUserResultChan := make(chan *models.UserChan)
			go c.userRepository.GetByID(getSalesOrderResult.SalesOrder.UserID, false, c.ctx, getUserResultChan)
			getUserResult := <-getUserResultChan

			if getUserResult.Error != nil {
				fmt.Println(getUserResult.Error.Error())
				errors = append(errors, getUserResult.Error.Error())
			}

			// Get Salesman
			getSalesmanResultChan := make(chan *models.SalesmanChan)
			if getSalesOrderResult.SalesOrder.SalesmanID.Int64 > 0 {
				go c.salesmanRepository.GetByID(int(getSalesOrderResult.SalesOrder.SalesmanID.Int64), false, c.ctx, getSalesmanResultChan)
			} else {
				go c.salesmanRepository.GetByEmail(getUserResult.User.Email, false, c.ctx, getSalesmanResultChan)
			}
			getSalesmanResult := <-getSalesmanResultChan

			if getSalesmanResult.Error != nil {
				fmt.Println(getSalesmanResult.Error.Error())
				errors = append(errors, getSalesmanResult.Error.Error())
			}

			deliveryOrder := &models.DeliveryOrder{}

			deliveryOrder.DeliveryOrderUploadMap(v, getSalesOrderResult.SalesOrder.ID, getWarehouseResult.Warehouse.ID, now)
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
			deliveryOrder.CreatedBy = v.IDUser
			deliveryOrder.SalesOrder = getSalesOrderResult.SalesOrder
			deliveryOrder.Brand = getBrandResult.Brand

			if getSalesmanResult.Salesman != nil {
				deliveryOrder.Salesman = getSalesmanResult.Salesman
			}

			// Delivery Order Detail
			deliveryOrderDetails := []*models.DeliveryOrderDetail{}
			for _, x := range getSalesOrderDetailBySoIdResult.SalesOrderDetails {
				if getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail.ID == x.ID {
					// Get Sales Order Detail By ID
					getSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
					go c.salesOrderDetailRepository.GetByID(x.ID, false, c.ctx, getSalesOrderDetailResultChan)
					getSalesOrderDetailResult := <-getSalesOrderDetailResultChan

					if getSalesOrderDetailResult.Error != nil {
						fmt.Println(getSalesOrderDetailResult.Error.Error())
						errors = append(errors, getSalesOrderDetailResult.Error.Error())
					}

					// Get SO Detail Order Status by ID
					getOrderStatusDetailResultChan := make(chan *models.OrderStatusChan)
					go c.orderStatusRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.ID, false, c.ctx, getOrderStatusDetailResultChan)
					getOrderStatusDetailResult := <-getOrderStatusDetailResultChan

					if getOrderStatusDetailResult.Error != nil {
						fmt.Println(getOrderStatusDetailResult.Error.Error())
						errors = append(errors, getOrderStatusDetailResult.Error.Error())
					}

					getSalesOrderDetailResult.SalesOrderDetail.UpdatedAt = &now
					getSalesOrderDetailResult.SalesOrderDetail.SentQty += v.QTYShip
					getSalesOrderDetailResult.SalesOrderDetail.ResidualQty -= v.QTYShip
					statusName := "partial"

					if getSalesOrderDetailResult.SalesOrderDetail.ResidualQty == 0 {
						statusName = "closed"
					}

					// Get SO Detail Order Status By statusName
					getOrderStatusSODetailResultChan := make(chan *models.OrderStatusChan)
					go c.orderStatusRepository.GetByNameAndType(statusName, "sales_order_detail", false, c.ctx, getOrderStatusSODetailResultChan)
					getOrderStatusSODetailResult := <-getOrderStatusSODetailResultChan

					if getOrderStatusSODetailResult.Error != nil {
						fmt.Println(getOrderStatusSODetailResult.Error.Error())
						errors = append(errors, getOrderStatusSODetailResult.Error.Error())
					}

					getSalesOrderDetailResult.SalesOrderDetail.OrderStatusID = getOrderStatusSODetailResult.OrderStatus.ID
					getSalesOrderDetailResult.SalesOrderDetail.OrderStatusName = getOrderStatusSODetailResult.OrderStatus.Name
					getSalesOrderDetailResult.SalesOrderDetail.OrderStatus = getOrderStatusSODetailResult.OrderStatus

					// Get Product Detail by ID
					getProductDetailResultChan := make(chan *models.ProductChan)
					go c.productRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.ProductID, false, c.ctx, getProductDetailResultChan)
					getProductDetailResult := <-getProductDetailResultChan

					if getProductDetailResult.Error != nil {
						fmt.Println(getProductDetailResult.Error.Error())
						errors = append(errors, getProductDetailResult.Error.Error())
					}

					// Get Uom Detail By ID
					getUomDetailResultChan := make(chan *models.UomChan)
					go c.uomRepository.GetByID(getSalesOrderDetailResult.SalesOrderDetail.UomID, false, c.ctx, getUomDetailResultChan)
					getUomDetailResult := <-getUomDetailResultChan

					if getUomDetailResult.Error != nil {
						fmt.Println(getUomDetailResult.Error.Error())
						errors = append(errors, getUomDetailResult.Error.Error())
					}

					deliveryOrderDetail := &models.DeliveryOrderDetail{}
					deliveryOrderDetail.DeliveryOrderDetailUploadMap(x.ID, v.QTYShip, now)
					deliveryOrderDetail.BrandID = getBrandResult.Brand.ID
					deliveryOrderDetail.Note = models.NullString{NullString: sql.NullString{String: v.Catatan, Valid: true}}
					deliveryOrderDetail.Uom = getUomDetailResult.Uom
					deliveryOrderDetail.ProductChanMap(getProductDetailResult)
					deliveryOrderDetail.SalesOrderDetailChanMap(getSalesOrderDetailResult)
					deliveryOrderDetail.OrderStatusID = deliveryOrder.OrderStatusID
					deliveryOrderDetail.OrderStatusName = deliveryOrder.OrderStatusName
					deliveryOrderDetail.OrderStatus = getOrderStatusDetailResult.OrderStatus

					deliveryOrderDetails = append(deliveryOrderDetails, deliveryOrderDetail)
					getSalesOrderResult.SalesOrder.SalesOrderDetails = append(getSalesOrderResult.SalesOrder.SalesOrderDetails, getSalesOrderDetailResult.SalesOrderDetail)
				}

				deliveryOrder.DeliveryOrderDetails = deliveryOrderDetails

			}
			a, _ := json.Marshal(deliveryOrder)
			fmt.Println("Disini", string(a))
			doRefCodes[deliveryOrder.DoRefCode.String] = deliveryOrder
		}

		for _, v := range doRefCodes {

			sqlTransaction, err := c.db.BeginTx(c.ctx, nil)

			if err != nil {
				fmt.Println(err.Error())
			}

			// Insert to DB, table delivery_orders
			createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
			go c.deliveryOrderRepository.Insert(v, sqlTransaction, c.ctx, createDeliveryOrderResultChan)
			createDeliveryOrderResult := <-createDeliveryOrderResultChan

			if createDeliveryOrderResult.Error != nil {
				sqlTransaction.Rollback()
			}

			// Delivery Order Detail
			totalResidualQty := 0
			for _, x := range v.DeliveryOrderDetails {
				x.DeliveryOrderID = int(createDeliveryOrderResult.ID)
				doDetailCode, _ := helper.GenerateDODetailCode(createDeliveryOrderResult.DeliveryOrder.ID, v.AgentID, x.Product.ID, x.Uom.ID)
				x.DoDetailCode = doDetailCode

				// Insert to DB, Delivery Order Detail
				createDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
				go c.deliveryOrderDetailRepository.Insert(x, sqlTransaction, c.ctx, createDeliveryOrderDetailResultChan)
				createDeliveryOrderDetailResult := <-createDeliveryOrderDetailResultChan

				if createDeliveryOrderDetailResult.Error != nil {
					sqlTransaction.Rollback()
				}

				x.ID = int(createDeliveryOrderDetailResult.ID)

				// Update to DB, Sales Order Detail
				updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
				go c.salesOrderDetailRepository.UpdateByID(x.SoDetailID, x.SoDetail, sqlTransaction, c.ctx, updateSalesOrderDetailResultChan)
				updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

				if updateSalesOrderDetailResult.Error != nil {
					sqlTransaction.Rollback()
				}

				totalResidualQty += x.SoDetail.ResidualQty

			}

			// Get updated Order Status Sales Order
			getOrderStatusSOResultChan := make(chan *models.OrderStatusChan)
			go c.orderStatusRepository.GetByID(v.SalesOrder.OrderStatusID, false, c.ctx, getOrderStatusSOResultChan)
			getOrderStatusSODetailResult := <-getOrderStatusSOResultChan

			// if getOrderStatusSODetailResult.Error != nil {
			// 	fmt.Println(getOrderStatusSODetailResult.Error.Error())
			// 	errors = append(errors, getOrderStatusSODetailResult.Error.Error())
			// }

			v.SalesOrder.OrderStatus = getOrderStatusSODetailResult.OrderStatus
			v.SalesOrder.OrderStatusName = getOrderStatusSODetailResult.OrderStatus.Name
			v.SalesOrder.SoDate = ""
			v.SalesOrder.SoRefDate = models.NullString{}
			v.SalesOrder.UpdatedAt = &now
			v.SalesOrder = v.SalesOrder

			var statusSoJourney string
			if totalResidualQty == 0 {
				statusSoJourney = constants.SO_STATUS_ORDCLS
				v.SalesOrder.OrderStatusID = 8
			} else {
				statusSoJourney = constants.SO_STATUS_ORDPRT
				v.SalesOrder.OrderStatusID = 7
			}

			// Update to DB, Sales Order
			updateSalesOrderChan := make(chan *models.SalesOrderChan)
			go c.salesOrderRepository.UpdateByID(v.SalesOrder.ID, v.SalesOrder, sqlTransaction, c.ctx, updateSalesOrderChan)
			updateSalesOrderResult := <-updateSalesOrderChan

			if updateSalesOrderResult.Error != nil {
				sqlTransaction.Rollback()
			}

			sqlTransaction.Commit()

			deliveryOrderLog := &models.DeliveryOrderLog{
				RequestID: requestId,
				DoCode:    v.DoCode,
				Data:      v,
				Status:    constants.LOG_STATUS_MONGO_DEFAULT,
				Action:    constants.LOG_ACTION_MONGO_INSERT,
				CreatedAt: &now,
			}

			createDeliveryOrderLogResultChan := make(chan *models.DeliveryOrderLogChan)
			go c.deliveryOrderLogRepository.Insert(deliveryOrderLog, c.ctx, createDeliveryOrderLogResultChan)
			createDeliveryOrderLogResult := <-createDeliveryOrderLogResultChan

			if createDeliveryOrderLogResult.Error != nil {
				fmt.Println(createDeliveryOrderLogResult.Error)
			}

			deliveryOrderJourney := &models.DeliveryOrderJourney{
				DoId:      v.ID,
				DoCode:    v.DoCode,
				Status:    constants.DO_STATUS_OPEN,
				Remark:    "",
				Reason:    "",
				CreatedAt: &now,
				UpdatedAt: &now,
			}

			createDeliveryOrderJourneyChan := make(chan *models.DeliveryOrderJourneyChan)
			go c.deliveryOrderLogRepository.InsertJourney(deliveryOrderJourney, c.ctx, createDeliveryOrderJourneyChan)
			createDeliveryOrderJourneysResult := <-createDeliveryOrderJourneyChan

			if createDeliveryOrderJourneysResult.Error != nil {
				// return &models.DeliveryOrderStoreResponse{}, createDeliveryOrderJourneysResult.ErrorLog
			}

			salesOrderLog := &models.SalesOrderLog{
				RequestID: requestId,
				SoCode:    v.SalesOrder.SoCode,
				Data:      v.SalesOrder,
				Status:    constants.LOG_STATUS_MONGO_DEFAULT,
				Action:    constants.LOG_ACTION_MONGO_INSERT,
				CreatedAt: &now,
			}

			createSalesOrderLogResultChan := make(chan *models.SalesOrderLogChan)
			go c.salesOrderLogRepository.Insert(salesOrderLog, c.ctx, createSalesOrderLogResultChan)
			createSalesOrderLogResult := <-createSalesOrderLogResultChan

			if createSalesOrderLogResult.Error != nil {
				fmt.Println(createSalesOrderLogResult.Error)
			}

			salesOrderJourney := &models.SalesOrderJourneys{
				SoId:      v.SalesOrder.ID,
				SoCode:    v.SalesOrder.SoCode,
				Status:    statusSoJourney,
				Remark:    "",
				Reason:    "",
				CreatedAt: &now,
				UpdatedAt: &now,
			}

			createSalesOrderJourneyChan := make(chan *models.SalesOrderJourneysChan)
			go c.salesOrderJourneyRepository.Insert(salesOrderJourney, c.ctx, createSalesOrderJourneyChan)
			createSalesOrderJourneysResult := <-createSalesOrderJourneyChan

			if createSalesOrderJourneysResult.Error != nil {
				// return &models.DeliveryOrderStoreResponse{}, createSalesOrderJourneysResult.ErrorLog
			}

			keyKafka := []byte(v.DoCode)
			messageKafka, _ := json.Marshal(v)
			fmt.Println("kafka", string(messageKafka))
			err = c.kafkaClient.WriteToTopic(constants.CREATE_DELIVERY_ORDER_TOPIC, keyKafka, messageKafka)

			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
