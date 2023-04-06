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
	"strconv"
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
	doUploadHistoriesRepository   mongoRepositories.DoUploadHistoriesRepositoryInterface
	doUploadErrorLogsRepository   mongoRepositories.DoUploadErrorLogsRepositoryInterface
	kafkaClient                   kafkadbo.KafkaClientInterface
	ctx                           context.Context
	args                          []interface{}
	db                            dbresolver.DB
}

func InitUploadDOItemConsumerHandlerInterface(deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, deliveryOrderDetailRepository repositories.DeliveryOrderDetailRepositoryInterface, salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, orderSourceRepository repositories.OrderSourceRepositoryInterface, warehouseRepository repositories.WarehouseRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, uomRepository repositories.UomRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, userRepository repositories.UserRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, deliveryOrderLogRepository mongoRepositories.DeliveryOrderLogRepositoryInterface, salesOrderLogRepository mongoRepositories.SalesOrderLogRepositoryInterface, salesOrderJourneyRepository mongoRepositories.SalesOrderJourneysRepositoryInterface, doUploadHistoriesRepository mongoRepositories.DoUploadHistoriesRepositoryInterface,
	doUploadErrorLogsRepository mongoRepositories.DoUploadErrorLogsRepositoryInterface, kafkaClient kafkadbo.KafkaClientInterface, ctx context.Context, args []interface{}, db dbresolver.DB) UploadDOItemConsumerHandlerInterface {
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
		doUploadHistoriesRepository:   doUploadHistoriesRepository,
		doUploadErrorLogsRepository:   doUploadErrorLogsRepository,
		kafkaClient:                   kafkaClient,
		ctx:                           ctx,
		args:                          args,
		db:                            db,
	}
}

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

		doRefCodes := []string{}
		deliveryOrderRefCodes := map[string]*models.DeliveryOrder{}

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

			checkIfDoRefCodeExist := helper.InSliceString(doRefCodes, v.NoSJ)

			// Get Agent By ID / IDDistributor
			getAgentResultChan := make(chan *models.AgentChan)
			go c.agentRepository.GetByID(v.IDDistributor, false, c.ctx, getAgentResultChan)
			getAgentResult := <-getAgentResultChan

			if getAgentResult.Error != nil {
				fmt.Println(getAgentResult.Error)
				errors = append(errors, getAgentResult.Error.Error())
			}

			// Get Warehouse
			getWarehouseResultChan := make(chan *models.WarehouseChan)
			if v.KodeGudang == "" {
				go c.warehouseRepository.GetByAgentID(v.IDDistributor, false, c.ctx, getWarehouseResultChan)
			} else {
				go c.warehouseRepository.GetByCode(v.KodeGudang, false, c.ctx, getWarehouseResultChan)
			}
			getWarehouseResult := <-getWarehouseResultChan

			if getWarehouseResult.Error != nil {
				fmt.Println(getWarehouseResult.Error.Error())
				errors = append(errors, getWarehouseResult.Error.Error())
			}

			// Get Sales Order By SoCode / NoOrder
			getSalesOrderResultChan := make(chan *models.SalesOrderChan)
			go c.salesOrderRepository.GetByCode(v.NoOrder, false, c.ctx, getSalesOrderResultChan)
			getSalesOrderResult := <-getSalesOrderResultChan

			if getSalesOrderResult.Error != nil {
				fmt.Println(getSalesOrderResult.Error.Error())
				errors = append(errors, getSalesOrderResult.Error.Error())
			}

			// Get Order Status for SO
			getOrderStatusSOResultChan := make(chan *models.OrderStatusChan)
			go c.orderStatusRepository.GetByID(getSalesOrderResult.SalesOrder.OrderStatusID, false, c.ctx, getOrderStatusSOResultChan)
			getOrderStatusResult := <-getOrderStatusSOResultChan

			if getOrderStatusResult.Error != nil {
				fmt.Println(getOrderStatusResult.Error.Error())
				errors = append(errors, getOrderStatusResult.Error.Error())
			} else if getOrderStatusResult.OrderStatus.Name != "open" && getOrderStatusResult.OrderStatus.Name != "partial" {
				errorMessage := fmt.Sprintf("Status Sales Order %s. Mohon sesuaikan kembali.", getOrderStatusResult.OrderStatus.Name)
				fmt.Println(errorMessage)
				errors = append(errors, errorMessage)
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

			getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail.UpdatedAt = &now
			getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail.SentQty += v.QTYShip
			getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail.ResidualQty -= v.QTYShip
			statusName := "partial"

			if getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail.ResidualQty == 0 {
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

			// Get Product Detail by ID
			getProductDetailResultChan := make(chan *models.ProductChan)
			go c.productRepository.GetByID(getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail.ProductID, false, c.ctx, getProductDetailResultChan)
			getProductDetailResult := <-getProductDetailResultChan

			if getProductDetailResult.Error != nil {
				fmt.Println(getProductDetailResult.Error.Error())
				errors = append(errors, getProductDetailResult.Error.Error())
			} else if v.Unit != getProductDetailResult.Product.UnitMeasurementSmall.String && v.Unit != getProductDetailResult.Product.UnitMeasurementMedium.String && v.Unit != getProductDetailResult.Product.UnitMeasurementBig.String {
				errorMessage := fmt.Sprintf("Satuan produk untuk SKU %s di data order dan file upload surat jalan tidak sesuai. Mohon sesuaikan kembali", v.KodeProduk)
				fmt.Println(errorMessage)
				errors = append(errors, errorMessage)
			}

			// Get Uom Detail By ID
			getUomDetailResultChan := make(chan *models.UomChan)
			go c.uomRepository.GetByID(getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail.UomID, false, c.ctx, getUomDetailResultChan)
			getUomDetailResult := <-getUomDetailResultChan

			if getUomDetailResult.Error != nil {
				fmt.Println(getUomDetailResult.Error.Error())
				errors = append(errors, getUomDetailResult.Error.Error())
			}

			if checkIfDoRefCodeExist {

				deliveryOrder := deliveryOrderRefCodes[v.NoSJ]

				if v.NoOrder != deliveryOrder.SalesOrder.SoCode {
					errorMessage := fmt.Sprintf("No SJ %s hanya bisa memiliki 1 No Order", v.NoSJ)
					fmt.Println(errorMessage)
					errors = append(errors, errorMessage)
				}

				if v.KodeMerk != deliveryOrder.Brand.ID {
					errorMessage := fmt.Sprintf("No SJ %d hanya bisa memiliki 1 Kode Merk", v.KodeMerk)
					fmt.Println(errorMessage)
					errors = append(errors, errorMessage)
				}

				if len(errors) < 1 {
					getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail.OrderStatusID = getOrderStatusSODetailResult.OrderStatus.ID
					getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail.OrderStatusName = getOrderStatusSODetailResult.OrderStatus.Name
					getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail.OrderStatus = getOrderStatusSODetailResult.OrderStatus

					deliveryOrderDetail := &models.DeliveryOrderDetail{}
					deliveryOrderDetail.DeliveryOrderDetailUploadMap(getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail.ID, v.QTYShip, now)
					deliveryOrderDetail.BrandID = getBrandResult.Brand.ID
					deliveryOrderDetail.Note = models.NullString{NullString: sql.NullString{String: v.Catatan, Valid: true}}
					deliveryOrderDetail.Uom = getUomDetailResult.Uom
					deliveryOrderDetail.ProductChanMap(getProductDetailResult)
					deliveryOrderDetail.SalesOrderDetailMap(getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail)
					deliveryOrderDetail.OrderStatusID = deliveryOrder.OrderStatusID
					deliveryOrderDetail.OrderStatusName = deliveryOrder.OrderStatusName
					deliveryOrderDetail.OrderStatus = deliveryOrder.OrderStatus

					// Delivery Order Detail
					deliveryOrder.DeliveryOrderDetails = append(deliveryOrder.DeliveryOrderDetails, deliveryOrderDetail)
					getSalesOrderResult.SalesOrder.SalesOrderDetails = append(getSalesOrderResult.SalesOrder.SalesOrderDetails, getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail)

					deliveryOrderRefCodes[v.NoSJ] = deliveryOrder
				} else {
					getUploadSOHistoriesResultChan := make(chan *models.DoUploadHistoryChan)
					go c.doUploadHistoriesRepository.GetByID(v.SjUploadHistoryId, false, c.ctx, getUploadSOHistoriesResultChan)
					getUploadSOHistoriesResult := <-getUploadSOHistoriesResultChan
					message := getUploadSOHistoriesResult.DoUploadHistory

					if v.UploadType == "retry" {

						c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
						deliveryOrderRefCodes = nil
						break
					} else {

						var myMap map[string]string
						data, _ := json.Marshal(v)
						json.Unmarshal(data, &myMap)

						c.createSjUploadErrorLog(v.ErrorLine, strconv.Itoa(v.IDDistributor), v.SjUploadHistoryId, message.RequestId, getAgentResult.Agent.Name, message.BulkCode, getWarehouseResult.Warehouse.Name, errors, &now, myMap)
						continue
					}
				}
			} else {

				doRefCodes = append(doRefCodes, v.NoSJ)

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

				if len(errors) < 1 {
					deliveryOrder := &models.DeliveryOrder{}

					deliveryOrder.DeliveryOrderUploadMap(v, getSalesOrderResult.SalesOrder.ID, getWarehouseResult.Warehouse.ID, now)
					deliveryOrder.WarehouseChanMap(getWarehouseResult)
					deliveryOrder.AgentMap(getAgentResult.Agent)
					deliveryOrder.DoCode = helper.GenerateDOCode(getAgentResult.Agent.ID, getOrderSourceResult.OrderSource.Code)
					deliveryOrder.DoDate = v.TanggalSJ
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

					getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail.OrderStatusID = getOrderStatusSODetailResult.OrderStatus.ID
					getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail.OrderStatusName = getOrderStatusSODetailResult.OrderStatus.Name
					getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail.OrderStatus = getOrderStatusSODetailResult.OrderStatus

					deliveryOrderDetail := &models.DeliveryOrderDetail{}
					deliveryOrderDetail.DeliveryOrderDetailUploadMap(getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail.ID, v.QTYShip, now)
					deliveryOrderDetail.BrandID = getBrandResult.Brand.ID
					deliveryOrderDetail.Note = models.NullString{NullString: sql.NullString{String: v.Catatan, Valid: true}}
					deliveryOrderDetail.Uom = getUomDetailResult.Uom
					deliveryOrderDetail.ProductChanMap(getProductDetailResult)
					deliveryOrderDetail.SalesOrderDetailMap(getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail)
					deliveryOrderDetail.OrderStatusID = deliveryOrder.OrderStatusID
					deliveryOrderDetail.OrderStatusName = deliveryOrder.OrderStatusName
					deliveryOrderDetail.OrderStatus = deliveryOrder.OrderStatus

					// Delivery Order Detail
					deliveryOrderDetails := []*models.DeliveryOrderDetail{}
					deliveryOrderDetails = append(deliveryOrderDetails, deliveryOrderDetail)
					getSalesOrderResult.SalesOrder.SalesOrderDetails = append(getSalesOrderResult.SalesOrder.SalesOrderDetails, getSODetailBySoIdSkuAndUomCodeResult.SalesOrderDetail)

					deliveryOrder.DeliveryOrderDetails = deliveryOrderDetails

					deliveryOrderRefCodes[v.NoSJ] = deliveryOrder
				} else {
					getUploadSOHistoriesResultChan := make(chan *models.DoUploadHistoryChan)
					go c.doUploadHistoriesRepository.GetByID(v.SjUploadHistoryId, false, c.ctx, getUploadSOHistoriesResultChan)
					getUploadSOHistoriesResult := <-getUploadSOHistoriesResultChan
					message := getUploadSOHistoriesResult.DoUploadHistory

					if v.UploadType == "retry" {

						c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
						deliveryOrderRefCodes = nil
						break
					} else {

						var myMap map[string]string
						data, _ := json.Marshal(v)
						json.Unmarshal(data, &myMap)

						c.createSjUploadErrorLog(v.ErrorLine, strconv.Itoa(v.IDDistributor), v.SjUploadHistoryId, message.RequestId, getAgentResult.Agent.Name, message.BulkCode, getWarehouseResult.Warehouse.Name, errors, &now, myMap)
						continue
					}
				}
			}
		}

		for _, v := range deliveryOrderRefCodes {

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
					fmt.Println(createDeliveryOrderDetailResult.Error.Error())
				}

				x.ID = int(createDeliveryOrderDetailResult.ID)

				// Update to DB, Sales Order Detail
				updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
				go c.salesOrderDetailRepository.UpdateByID(x.SoDetailID, x.SoDetail, sqlTransaction, c.ctx, updateSalesOrderDetailResultChan)
				updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

				if updateSalesOrderDetailResult.Error != nil {
					sqlTransaction.Rollback()
					fmt.Println(updateSalesOrderDetailResult.Error.Error())
				}

				totalResidualQty += x.SoDetail.ResidualQty

			}

			// Get updated Order Status Sales Order
			getOrderStatusSOResultChan := make(chan *models.OrderStatusChan)
			go c.orderStatusRepository.GetByID(v.SalesOrder.OrderStatusID, false, c.ctx, getOrderStatusSOResultChan)
			getOrderStatusSODetailResult := <-getOrderStatusSOResultChan

			if getOrderStatusSODetailResult.Error != nil {
				sqlTransaction.Rollback()
				fmt.Println(getOrderStatusSODetailResult.Error.Error())
			}

			v.SalesOrder.OrderStatus = getOrderStatusSODetailResult.OrderStatus
			v.SalesOrder.OrderStatusName = getOrderStatusSODetailResult.OrderStatus.Name
			v.SalesOrder.SoDate = ""
			v.SalesOrder.SoRefDate = models.NullString{}
			v.SalesOrder.UpdatedAt = &now

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
				fmt.Println(updateSalesOrderResult.Error.Error())
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
				fmt.Println(createDeliveryOrderLogResult.Error.Error())
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
				fmt.Println(createDeliveryOrderLogResult.Error.Error())
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
				fmt.Println(createSalesOrderLogResult.Error.Error())
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
				fmt.Println(createSalesOrderJourneysResult.Error.Error())
			}

			keyKafka := []byte(v.DoCode)
			messageKafka, _ := json.Marshal(v)

			err = c.kafkaClient.WriteToTopic(constants.CREATE_DELIVERY_ORDER_TOPIC, keyKafka, messageKafka)

			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (c *uploadDOItemConsumerHandler) createSjUploadErrorLog(errorLine int, agentId, sjUploadHistoryId, requestId, agentName, bulkCode, warehouseName string, errors []string, now *time.Time, item map[string]string) {

	rowData := &models.RowDataDoUploadErrorLog{}
	rowData.RowDataDoUploadErrorLogMap(item, agentName, warehouseName)

	doUploadErrorLog := &models.DoUploadErrorLog{}
	doUploadErrorLog.DoUploadErrorLogsMap(errorLine, sjUploadHistoryId, requestId, bulkCode, errors, now)
	doUploadErrorLog.RowData = *rowData

	doUploadErrorLogssResultChan := make(chan *models.DoUploadErrorLogChan)
	go c.doUploadErrorLogsRepository.Insert(doUploadErrorLog, c.ctx, doUploadErrorLogssResultChan)
}

func (c *uploadDOItemConsumerHandler) updateSjUploadHistories(message *models.DoUploadHistory, status string) {
	message.UpdatedAt = time.Now()
	message.Status = status
	uploadDOHistoryJourneysResultChan := make(chan *models.DoUploadHistoryChan)
	go c.doUploadHistoriesRepository.UpdateByID(message.ID.Hex(), message, c.ctx, uploadDOHistoryJourneysResultChan)
}
