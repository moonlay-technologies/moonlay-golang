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
	"strings"
	"time"

	"github.com/bxcodec/dbresolver"
)

type UploadSOSJItemConsumerHandlerInterface interface {
	ProcessMessage()
}

type uploadSOSJItemConsumerHandler struct {
	orderSourceRepository              repositories.OrderSourceRepositoryInterface
	orderStatusRepository              repositories.OrderStatusRepositoryInterface
	productRepository                  repositories.ProductRepositoryInterface
	uomRepository                      repositories.UomRepositoryInterface
	agentRepository                    repositories.AgentRepositoryInterface
	storeRepository                    repositories.StoreRepositoryInterface
	userRepository                     repositories.UserRepositoryInterface
	salesmanRepository                 repositories.SalesmanRepositoryInterface
	brandRepository                    repositories.BrandRepositoryInterface
	salesOrderRepository               repositories.SalesOrderRepositoryInterface
	salesOrderDetailRepository         repositories.SalesOrderDetailRepositoryInterface
	salesOrderLogRepository            mongoRepositories.SalesOrderLogRepositoryInterface
	salesOrderJourneysRepository       mongoRepositories.SalesOrderJourneysRepositoryInterface
	salesOrderDetailJourneysRepository mongoRepositories.SalesOrderDetailJourneysRepositoryInterface
	warehouseRepository                repositories.WarehouseRepositoryInterface
	deliveryOrderRepository            repositories.DeliveryOrderRepositoryInterface
	deliveryOrderDetailRepository      repositories.DeliveryOrderDetailRepositoryInterface
	deliveryOrderLogRepository         mongoRepositories.DeliveryOrderLogRepositoryInterface
	sosjUploadHistoriesRepository      mongoRepositories.SOSJUploadHistoriesRepositoryInterface
	sosjUploadErrorLogsRepository      mongoRepositories.SosjUploadErrorLogsRepositoryInterface
	uploadRepository                   repositories.UploadRepositoryInterface
	kafkaClient                        kafkadbo.KafkaClientInterface
	createSalesOrderConsumer           CreateSalesOrderConsumerHandlerInterface
	ctx                                context.Context
	args                               []interface{}
	db                                 dbresolver.DB
}

func InitUploadSOSJItemConsumerHandlerInterface(orderSourceRepository repositories.OrderSourceRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, productRepository repositories.ProductRepositoryInterface, uomRepository repositories.UomRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, userRepository repositories.UserRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, salesOrderLogRepository mongoRepositories.SalesOrderLogRepositoryInterface, salesOrderJourneysRepository mongoRepositories.SalesOrderJourneysRepositoryInterface, salesOrderDetailJourneysRepository mongoRepositories.SalesOrderDetailJourneysRepositoryInterface, warehouseRepository repositories.WarehouseRepositoryInterface, deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, deliveryOrderDetailRepository repositories.DeliveryOrderDetailRepositoryInterface, deliveryOrderLogRepository mongoRepositories.DeliveryOrderLogRepositoryInterface, sosjUploadHistoriesRepository mongoRepositories.SOSJUploadHistoriesRepositoryInterface, sosjUploadErrorLogsRepository mongoRepositories.SosjUploadErrorLogsRepositoryInterface, uploadRepository repositories.UploadRepositoryInterface, createSalesOrderConsumer CreateSalesOrderConsumerHandlerInterface, kafkaClient kafkadbo.KafkaClientInterface, db dbresolver.DB, ctx context.Context, args []interface{}) UploadSOSJItemConsumerHandlerInterface {
	return &uploadSOSJItemConsumerHandler{
		orderSourceRepository:              orderSourceRepository,
		orderStatusRepository:              orderStatusRepository,
		productRepository:                  productRepository,
		uomRepository:                      uomRepository,
		agentRepository:                    agentRepository,
		storeRepository:                    storeRepository,
		userRepository:                     userRepository,
		salesmanRepository:                 salesmanRepository,
		brandRepository:                    brandRepository,
		salesOrderRepository:               salesOrderRepository,
		salesOrderDetailRepository:         salesOrderDetailRepository,
		salesOrderLogRepository:            salesOrderLogRepository,
		salesOrderJourneysRepository:       salesOrderJourneysRepository,
		salesOrderDetailJourneysRepository: salesOrderDetailJourneysRepository,
		warehouseRepository:                warehouseRepository,
		deliveryOrderRepository:            deliveryOrderRepository,
		deliveryOrderDetailRepository:      deliveryOrderDetailRepository,
		deliveryOrderLogRepository:         deliveryOrderLogRepository,
		sosjUploadHistoriesRepository:      sosjUploadHistoriesRepository,
		sosjUploadErrorLogsRepository:      sosjUploadErrorLogsRepository,
		uploadRepository:                   uploadRepository,
		createSalesOrderConsumer:           createSalesOrderConsumer,
		kafkaClient:                        kafkaClient,
		ctx:                                ctx,
		args:                               args,
		db:                                 db,
	}
}

func (c *uploadSOSJItemConsumerHandler) ProcessMessage() {
	fmt.Println("process ", constants.UPLOAD_SOSJ_ITEM_TOPIC)
	topic := c.args[1].(string)
	groupID := c.args[2].(string)
	reader := c.kafkaClient.SetConsumerGroupReader(topic, groupID)

	for {
		m, err := reader.ReadMessage(c.ctx)
		if err != nil {
			break
		}

		fmt.Printf("message so at topic/partition/offset %v/%v/%v \n", m.Topic, m.Partition, m.Offset)
		now := time.Now()

		var uploadSOSJFields []*models.UploadSOSJField
		err = json.Unmarshal(m.Value, &uploadSOSJFields)

		requestId := string(m.Key[:])

		// Get Order Source Status By Id
		getOrderSourceResultChan := make(chan *models.OrderSourceChan)
		go c.orderSourceRepository.GetBySourceName("upload_sosj", false, c.ctx, getOrderSourceResultChan)
		getOrderSourceResult := <-getOrderSourceResultChan

		soRefCodes := []string{}
		salesOrderSoRefCodes := map[string]*models.SalesOrder{}

		var soStatus string
		var doStatus string
		for _, v := range uploadSOSJFields {

			var noSuratJalan = v.NoSuratJalan
			var errors []string

			if getOrderSourceResult.Error != nil {
				fmt.Println(getOrderSourceResult.Error.Error())
				errors = append(errors, getOrderSourceResult.Error.Error())
			}

			checkIfSoRefCodeExist := helper.InSliceString(soRefCodes, noSuratJalan)

			switch strings.ToLower(v.Status) {
			case "open":
				soStatus = constants.ORDER_STATUS_CLOSED
				doStatus = constants.ORDER_STATUS_OPEN
			case "closed":
				soStatus = constants.ORDER_STATUS_CLOSED
				doStatus = constants.ORDER_STATUS_CLOSED
			}

			// Get SO Status By Name
			getSalesOrderStatusResultChan := make(chan *models.OrderStatusChan)
			go c.orderStatusRepository.GetByNameAndType(soStatus, "sales_order", false, c.ctx, getSalesOrderStatusResultChan)
			getSalesOrderStatusResult := <-getSalesOrderStatusResultChan

			if getSalesOrderStatusResult.Error != nil {
				fmt.Println(getSalesOrderStatusResult.Error.Error())
				errors = append(errors, getSalesOrderStatusResult.Error.Error())
			}

			// Get SO Status By Name
			getSalesOrderDetailStatusResultChan := make(chan *models.OrderStatusChan)
			go c.orderStatusRepository.GetByNameAndType(soStatus, "sales_order_detail", false, c.ctx, getSalesOrderDetailStatusResultChan)
			getSalesOrderDetailStatusResult := <-getSalesOrderDetailStatusResultChan

			if getSalesOrderDetailStatusResult.Error != nil {
				fmt.Println(getSalesOrderDetailStatusResult.Error.Error())
				errors = append(errors, getSalesOrderDetailStatusResult.Error.Error())
			}

			// Get DO Status By Name
			getDeliveryOrderStatusResultChan := make(chan *models.OrderStatusChan)
			go c.orderStatusRepository.GetByNameAndType(doStatus, "delivery_order", false, c.ctx, getDeliveryOrderStatusResultChan)
			getDeliveryOrderStatusResult := <-getDeliveryOrderStatusResultChan

			if getDeliveryOrderStatusResult.Error != nil {
				fmt.Println(getDeliveryOrderStatusResult.Error.Error())
				errors = append(errors, getDeliveryOrderStatusResult.Error.Error())
			}

			// Check Product By Id
			getProductResultChan := make(chan *models.ProductChan)
			go c.productRepository.GetBySKU(v.KodeProdukDBO, false, c.ctx, getProductResultChan)
			getProductResult := <-getProductResultChan

			if getProductResult.Error != nil {
				fmt.Println(getProductResult.Error.Error())
				errors = append(errors, getProductResult.Error.Error())
			}

			// Check Uom By Id
			getUomResultChan := make(chan *models.UomChan)
			go c.uomRepository.GetByID(v.Unit, false, c.ctx, getUomResultChan)
			getUomResult := <-getUomResultChan

			if getUomResult.Error != nil {
				fmt.Println(getUomResult.Error.Error())
				errors = append(errors, getUomResult.Error.Error())
			}

			var price float64

			if getUomResult.Uom.Code.String == getProductResult.Product.UnitMeasurementSmall.String {
				price = getProductResult.Product.PriceSmall
			} else if getUomResult.Uom.Code.String == getProductResult.Product.UnitMeasurementMedium.String {
				price = getProductResult.Product.PriceMedium
			} else {
				price = getProductResult.Product.PriceBig
			}

			if checkIfSoRefCodeExist {
				salesOrder := salesOrderSoRefCodes[noSuratJalan]

				salesOrder.TotalAmount = salesOrder.TotalAmount + (price * float64(v.Qty))
				salesOrder.TotalTonase = salesOrder.TotalTonase + (float64(v.Qty) * getProductResult.Product.NettWeight)

				// ### Sales Order Detail ###
				salesOrderDetail := &models.SalesOrderDetail{}
				salesOrderDetail.SalesOrderDetailUploadSOSJMap(v, now)
				salesOrderDetail.SalesOrderDetailStatusChanMap(getSalesOrderDetailStatusResult)
				salesOrderDetail.ProductID = getProductResult.Product.ID
				salesOrderDetail.OrderStatusID = getSalesOrderDetailStatusResult.OrderStatus.ID
				salesOrderDetail.Price = price
				salesOrderDetail.Product = getProductResult.Product
				salesOrderDetail.Uom = getUomResult.Uom

				salesOrder.SalesOrderDetails = append(salesOrder.SalesOrderDetails, salesOrderDetail)

				// ### Delivery Order Detail ###
				deliveryOrderDetail := &models.DeliveryOrderDetail{}
				deliveryOrderDetail.DeliveryOrderDetailUploadSOSJMap(v, now)

				deliveryOrderDetail.BrandID = v.IDMerk
				deliveryOrderDetail.Note = models.NullString{NullString: sql.NullString{String: v.Catatan, Valid: true}}
				deliveryOrderDetail.UomID = v.Unit
				deliveryOrderDetail.Uom = getUomResult.Uom
				deliveryOrderDetail.ProductID = getProductResult.Product.ID
				deliveryOrderDetail.ProductChanMap(getProductResult)
				deliveryOrderDetail.OrderStatusID = getDeliveryOrderStatusResult.OrderStatus.ID
				deliveryOrderDetail.OrderStatusName = getDeliveryOrderStatusResult.OrderStatus.Name
				deliveryOrderDetail.OrderStatus = getDeliveryOrderStatusResult.OrderStatus

				deliveryOrder := salesOrder.DeliveryOrders[0]
				deliveryOrder.DeliveryOrderDetails = append(deliveryOrder.DeliveryOrderDetails, deliveryOrderDetail)

				salesOrderSoRefCodes[noSuratJalan] = salesOrder

			} else {

				soRefCodes = append(soRefCodes, noSuratJalan)

				// Check Agent By Id
				getAgentResultChan := make(chan *models.AgentChan)
				go c.agentRepository.GetByID(v.IDDistributor, false, c.ctx, getAgentResultChan)
				getAgentResult := <-getAgentResultChan

				if getAgentResult.Error != nil {
					fmt.Println(getAgentResult.Error.Error())
					errors = append(errors, getAgentResult.Error.Error())
				}

				// Get store id by store code
				getStoreIdResultChan := make(chan *models.StoreChan)
				go c.storeRepository.GetIdByStoreCode(v.KodeTokoDBO, false, c.ctx, getStoreIdResultChan)
				getStoreIdResult := <-getStoreIdResultChan

				if getStoreIdResult.Error != nil {
					fmt.Println(getStoreIdResult.Error.Error())
					errors = append(errors, getStoreIdResult.Error.Error())
				}

				// Check Store By Id
				getStoreResultChan := make(chan *models.StoreChan)
				go c.storeRepository.GetByID(getStoreIdResult.Store.ID, false, c.ctx, getStoreResultChan)
				getStoreResult := <-getStoreResultChan

				if getStoreResult.Error != nil {
					fmt.Println(getStoreResult.Error.Error())
					errors = append(errors, getStoreResult.Error.Error())
				}

				// Check User By Id
				getUserResultChan := make(chan *models.UserChan)
				go c.userRepository.GetByID(v.IDUser, false, c.ctx, getUserResultChan)
				getUserResult := <-getUserResultChan

				if getUserResult.Error != nil {
					fmt.Println(getUserResult.Error.Error())
					errors = append(errors, getUserResult.Error.Error())
				}

				// Check Salesman By Id
				getSalesmanResult := &models.SalesmanChan{}
				if v.IDSalesman > 0 {
					getSalesmanResultChan := make(chan *models.SalesmanChan)
					go c.salesmanRepository.GetByID(v.IDSalesman, false, c.ctx, getSalesmanResultChan)
					getSalesmanResult = <-getSalesmanResultChan

					if getSalesmanResult.Error != nil {
						fmt.Println(getSalesmanResult.Error.Error())
						errors = append(errors, getSalesmanResult.Error.Error())
					}
				}

				// Check Brand By Id
				getBrandResultChan := make(chan *models.BrandChan)
				go c.brandRepository.GetByID(v.IDMerk, false, c.ctx, getBrandResultChan)
				getBrandResult := <-getBrandResultChan

				if getBrandResult.Error != nil {
					fmt.Println(getBrandResult.Error.Error())
					errors = append(errors, getBrandResult.Error.Error())
				}

				getWarehouseResultChan := make(chan *models.WarehouseChan)
				go c.warehouseRepository.GetByID(v.KodeGudang, false, c.ctx, getWarehouseResultChan)
				getWarehouseResult := <-getWarehouseResultChan

				if getWarehouseResult.Error != nil {
					fmt.Println(getWarehouseResult.Error.Error())
					errors = append(errors, getWarehouseResult.Error.Error())
				}

				// ### Sales Order ###
				salesOrder := &models.SalesOrder{}
				// soRefCodes = append(soRefCodes, noSuratJalan)

				salesOrder.SalesOrderUploadSOSJMap(v, now)
				salesOrder.OrderSourceChanMap(getOrderSourceResult)
				salesOrder.SalesOrderStatusChanMap(getSalesOrderStatusResult)
				salesOrder.AgentChanMap(getAgentResult)
				salesOrder.StoreChanMap(getStoreResult)
				salesOrder.UserChanMap(getUserResult)
				if v.IDSalesman > 0 {
					salesOrder.SalesmanChanMap(getSalesmanResult)
				}
				salesOrder.BrandChanMap(getBrandResult)

				salesOrder.UserID = v.IDUser
				salesOrder.CreatedBy = v.IDUser
				salesOrder.SoCode = helper.GenerateSOCode(v.IDDistributor, getOrderSourceResult.OrderSource.Code)
				salesOrder.OrderStatusID = getSalesOrderStatusResult.OrderStatus.ID
				salesOrder.OrderSourceID = getOrderSourceResult.OrderSource.ID
				salesOrder.TotalAmount = price * float64(v.Qty)
				salesOrder.TotalTonase = float64(v.Qty) * getProductResult.Product.NettWeight
				salesOrder.SoRefCode = models.NullString{NullString: sql.NullString{String: noSuratJalan, Valid: true}}

				// ### Sales Order Detail ###
				salesOrderDetail := &models.SalesOrderDetail{}
				salesOrderDetail.SalesOrderDetailUploadSOSJMap(v, now)
				salesOrderDetail.SalesOrderDetailStatusChanMap(getSalesOrderDetailStatusResult)
				salesOrderDetail.ProductID = getProductResult.Product.ID
				salesOrderDetail.OrderStatusID = getSalesOrderDetailStatusResult.OrderStatus.ID
				salesOrderDetail.Price = price
				salesOrderDetail.Product = getProductResult.Product
				salesOrderDetail.Uom = getUomResult.Uom

				salesOrderDetails := []*models.SalesOrderDetail{}
				salesOrderDetails = append(salesOrderDetails, salesOrderDetail)

				salesOrder.SalesOrderDetails = salesOrderDetails

				// ### Delivery Order ###
				deliveryOrder := &models.DeliveryOrder{}
				deliveryOrder.DeliveryOrderUploadSOSJMap(v, now)
				deliveryOrder.WarehouseChanMap(getWarehouseResult)
				deliveryOrder.AgentMap(getAgentResult.Agent)

				deliveryOrder.DoCode = helper.GenerateDOCode(v.IDDistributor, getOrderSourceResult.OrderSource.Code)
				deliveryOrder.DoRefCode = models.NullString{NullString: sql.NullString{String: noSuratJalan, Valid: true}}
				deliveryOrder.OrderStatus = getDeliveryOrderStatusResult.OrderStatus
				deliveryOrder.OrderStatusID = getDeliveryOrderStatusResult.OrderStatus.ID
				deliveryOrder.OrderSource = getOrderSourceResult.OrderSource
				deliveryOrder.OrderSourceID = getOrderSourceResult.OrderSource.ID
				deliveryOrder.Store = getStoreResult.Store
				deliveryOrder.StoreID = getStoreResult.Store.ID
				deliveryOrder.CreatedBy = v.IDUser
				deliveryOrder.LatestUpdatedBy = v.IDUser
				deliveryOrder.Brand = getBrandResult.Brand
				if getSalesmanResult.Salesman != nil {
					deliveryOrder.Salesman = getSalesmanResult.Salesman
				}

				// ### Delivery Order Detail ###
				deliveryOrderDetail := &models.DeliveryOrderDetail{}
				deliveryOrderDetail.DeliveryOrderDetailUploadSOSJMap(v, now)

				deliveryOrderDetails := []*models.DeliveryOrderDetail{}
				deliveryOrderDetails = append(deliveryOrderDetails, deliveryOrderDetail)
				deliveryOrderDetail.BrandID = v.IDMerk
				deliveryOrderDetail.Note = models.NullString{NullString: sql.NullString{String: v.Catatan, Valid: true}}
				deliveryOrderDetail.UomID = v.Unit
				deliveryOrderDetail.Uom = getUomResult.Uom
				deliveryOrderDetail.ProductID = getProductResult.Product.ID
				deliveryOrderDetail.ProductChanMap(getProductResult)
				deliveryOrderDetail.OrderStatusID = getDeliveryOrderStatusResult.OrderStatus.ID
				deliveryOrderDetail.OrderStatusName = getDeliveryOrderStatusResult.OrderStatus.Name
				deliveryOrderDetail.OrderStatus = getDeliveryOrderStatusResult.OrderStatus

				deliveryOrder.DeliveryOrderDetails = deliveryOrderDetails

				salesOrder.DeliveryOrders = append(salesOrder.DeliveryOrders, deliveryOrder)

				if len(errors) < 1 {
					salesOrderSoRefCodes[noSuratJalan] = salesOrder
				} else {
					sosjUploadHistoryJourneysResultChan := make(chan *models.GetSosjUploadHistoryResponseChan)
					go c.sosjUploadHistoriesRepository.GetByID(v.SosjUploadHistoryId, false, c.ctx, sosjUploadHistoryJourneysResultChan)
					sosjUploadHistoryJourneysResult := <-sosjUploadHistoryJourneysResultChan
					message := &sosjUploadHistoryJourneysResult.SosjUploadHistories.UploadHistory

					if v.UploadType == "retry" {
						message.Status = constants.UPLOAD_STATUS_HISTORY_FAILED
						sosjUploadHistoryJourneysResultChan := make(chan *models.UploadHistoryChan)
						go c.sosjUploadHistoriesRepository.UpdateByID(message.ID.Hex(), message, c.ctx, sosjUploadHistoryJourneysResultChan)

						salesOrderSoRefCodes = nil
						break
					} else {

						var myMap map[string]string
						data, _ := json.Marshal(v)
						json.Unmarshal(data, &myMap)

						rowDataSosjUploadErrorLogResultChan := make(chan *models.RowDataSosjUploadErrorLogChan)
						go c.uploadRepository.GetSosjRowData(strconv.Itoa(v.IDDistributor), v.KodeTokoDBO, strconv.Itoa(v.IDMerk), v.KodeProdukDBO, strconv.Itoa(v.KodeGudang), strconv.Itoa(v.IDSalesman), v.IDAlamat, rowDataSosjUploadErrorLogResultChan)
						rowDataSosjUploadErrorLogResult := <-rowDataSosjUploadErrorLogResultChan
						rowData := &models.RowDataSosjUploadErrorLog{}
						rowData.RowDataSosjUploadErrorLogMap2(*rowDataSosjUploadErrorLogResult.RowDataSosjUploadErrorLog, v)

						sosjUploadErrorLog := &models.SosjUploadErrorLog{}
						sosjUploadErrorLog.SosjUploadErrorLogsMap(v.ErrorLine, strconv.Itoa(v.IDDistributor), v.SosjUploadHistoryId, requestId, getAgentResult.Agent.Name, v.BulkCode, errors, &now)
						sosjUploadErrorLog.RowData = *rowData

						sosjUploadErrorLogResultChan := make(chan *models.SosjUploadErrorLogChan)
						go c.sosjUploadErrorLogsRepository.Insert(sosjUploadErrorLog, c.ctx, sosjUploadErrorLogResultChan)

						continue
					}
				}
			}
		}

		if salesOrderSoRefCodes == nil {
			continue
		}

		for _, v := range salesOrderSoRefCodes {

			sqlTransaction, err := c.db.BeginTx(c.ctx, nil)

			if err != nil {
				fmt.Println(err.Error())
			}

			createSalesOrderResultChan := make(chan *models.SalesOrderChan)
			go c.salesOrderRepository.Insert(v, sqlTransaction, c.ctx, createSalesOrderResultChan)
			createSalesOrderResult := <-createSalesOrderResultChan

			if createSalesOrderResult.Error != nil {
				sqlTransaction.Rollback()
				fmt.Println(createSalesOrderResult.Error.Error())
			}

			v.ID = createSalesOrderResult.SalesOrder.ID

			for _, x := range v.SalesOrderDetails {

				soDetailCode, _ := helper.GenerateSODetailCode(int(createSalesOrderResult.ID), v.AgentID, x.ProductID, x.UomID)
				x.SalesOrderID = int(createSalesOrderResult.ID)
				x.SoDetailCode = soDetailCode

				createSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
				go c.salesOrderDetailRepository.Insert(x, sqlTransaction, c.ctx, createSalesOrderDetailResultChan)
				createSalesOrderDetailResult := <-createSalesOrderDetailResultChan

				if createSalesOrderDetailResult.Error != nil {
					sqlTransaction.Rollback()
					fmt.Println(createSalesOrderDetailResult.Error.Error())
				}

				x.ID = createSalesOrderDetailResult.SalesOrderDetail.ID

			}

			for _, x := range v.DeliveryOrders {

				salesOrder := &models.SalesOrder{}
				salesOrder.SalesOrderForDOMap(v)

				x.SalesOrderID = v.ID
				x.SalesOrder = salesOrder

				createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
				go c.deliveryOrderRepository.Insert(x, sqlTransaction, c.ctx, createDeliveryOrderResultChan)
				createDeliveryOrderResult := <-createDeliveryOrderResultChan

				if createDeliveryOrderResult.Error != nil {
					sqlTransaction.Rollback()
					fmt.Println(createDeliveryOrderResult.Error.Error())
				}

				x.ID = createDeliveryOrderResult.DeliveryOrder.ID

				for i, doDetail := range x.DeliveryOrderDetails {

					doDetailCode, _ := helper.GenerateDODetailCode(createDeliveryOrderResult.DeliveryOrder.ID, v.AgentID, v.SalesOrderDetails[i].Product.ID, v.SalesOrderDetails[i].Uom.ID)

					doDetail.DeliveryOrderID = createDeliveryOrderResult.DeliveryOrder.ID
					doDetail.SoDetailID = v.SalesOrderDetails[i].ID
					doDetail.SoDetail = v.SalesOrderDetails[i]
					doDetail.DoDetailCode = doDetailCode

					createDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
					go c.deliveryOrderDetailRepository.Insert(doDetail, sqlTransaction, c.ctx, createDeliveryOrderDetailResultChan)
					createDeliveryOrderDetailResult := <-createDeliveryOrderDetailResultChan

					if createDeliveryOrderDetailResult.Error != nil {
						sqlTransaction.Rollback()
						fmt.Println(createDeliveryOrderDetailResult.Error.Error())
					}

					doDetail.ID = createDeliveryOrderDetailResult.DeliveryOrderDetail.ID

					x.DeliveryOrderDetails[i] = doDetail
				}

			}

			sqlTransaction.Commit()

			salesOrderLog := &models.SalesOrderLog{
				RequestID: requestId,
				SoCode:    v.SoCode,
				Data:      v,
				Status:    constants.LOG_STATUS_MONGO_DEFAULT,
				Action:    constants.LOG_ACTION_MONGO_INSERT,
				CreatedAt: &now,
				UpdatedAt: &now,
			}

			createSalesOrderLogResultChan := make(chan *models.SalesOrderLogChan)
			go c.salesOrderLogRepository.Insert(salesOrderLog, c.ctx, createSalesOrderLogResultChan)
			createSalesOrderLogResult := <-createSalesOrderLogResultChan

			if createSalesOrderLogResult.Error != nil {
				fmt.Println(createSalesOrderLogResult.Error.Error())
			}

			keyKafka := []byte(v.SoCode)
			messageKafka, _ := json.Marshal(v)
			fmt.Println("data so", messageKafka)
			err = c.createSalesOrderConsumer.CreateSoConsumer(v, messageKafka, constants.CREATE_SALES_ORDER_TOPIC, 0, 0, string(keyKafka), nil)
			if err != nil {
				fmt.Println(err.Error())
			}

			for _, x := range v.DeliveryOrders {

				deliveryOrderLog := &models.DeliveryOrderLog{
					RequestID: requestId,
					DoCode:    x.DoCode,
					Data:      x,
					Status:    constants.LOG_STATUS_MONGO_DEFAULT,
					Action:    constants.LOG_ACTION_MONGO_INSERT,
					CreatedAt: &now,
				}

				createDeliveryOrderLogResultChan := make(chan *models.DeliveryOrderLogChan)
				go c.deliveryOrderLogRepository.Insert(deliveryOrderLog, c.ctx, createDeliveryOrderLogResultChan)
				createDeliveryOrderLogResult := <-createDeliveryOrderLogResultChan

				if createDeliveryOrderLogResult.Error != nil {
					fmt.Println("error log do", createDeliveryOrderLogResult.ErrorLog)
				}

				keyKafka := []byte(x.DoCode)
				messageKafka, _ := json.Marshal(x)
				fmt.Println("data do", messageKafka)
				err := c.kafkaClient.WriteToTopic(constants.CREATE_DELIVERY_ORDER_TOPIC, keyKafka, messageKafka)

				if err != nil {
					fmt.Println(err.Error())
				}
			}

		}

	}
}
