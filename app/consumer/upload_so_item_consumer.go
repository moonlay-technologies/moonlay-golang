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

type UploadSOItemConsumerHandlerInterface interface {
	ProcessMessage()
}

type uploadSOItemConsumerHandler struct {
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
	soUploadHistoriesRepository        mongoRepositories.SoUploadHistoriesRepositoryInterface
	soUploadErrorLogsRepository        mongoRepositories.SoUploadErrorLogsRepositoryInterface
	kafkaClient                        kafkadbo.KafkaClientInterface
	ctx                                context.Context
	args                               []interface{}
	db                                 dbresolver.DB
}

func InitUploadSOItemConsumerHandlerInterface(orderSourceRepository repositories.OrderSourceRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, productRepository repositories.ProductRepositoryInterface, uomRepository repositories.UomRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, userRepository repositories.UserRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, salesOrderLogRepository mongoRepositories.SalesOrderLogRepositoryInterface, salesOrderJourneysRepository mongoRepositories.SalesOrderJourneysRepositoryInterface, salesOrderDetailJourneysRepository mongoRepositories.SalesOrderDetailJourneysRepositoryInterface, soUploadHistoriesRepository mongoRepositories.SoUploadHistoriesRepositoryInterface, soUploadErrorLogsRepository mongoRepositories.SoUploadErrorLogsRepositoryInterface, kafkaClient kafkadbo.KafkaClientInterface, db dbresolver.DB, ctx context.Context, args []interface{}) UploadSOItemConsumerHandlerInterface {
	return &uploadSOItemConsumerHandler{
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
		soUploadHistoriesRepository:        soUploadHistoriesRepository,
		soUploadErrorLogsRepository:        soUploadErrorLogsRepository,
		kafkaClient:                        kafkaClient,
		ctx:                                ctx,
		args:                               args,
		db:                                 db,
	}
}

func (c *uploadSOItemConsumerHandler) ProcessMessage() {
	fmt.Println("process ", constants.UPLOAD_SO_ITEM_TOPIC)
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

		var uploadSOFields []*models.UploadSOField
		err = json.Unmarshal(m.Value, &uploadSOFields)

		requestId := string(m.Key[:])

		soRefCodes := []string{}
		salesOrderSoRefCodes := map[string]*models.SalesOrder{}

		// Get Order Source Status By Id
		getOrderSourceResultChan := make(chan *models.OrderSourceChan)
		go c.orderSourceRepository.GetBySourceName("upload_so", false, c.ctx, getOrderSourceResultChan)
		getOrderSourceResult := <-getOrderSourceResultChan

		for i, v := range uploadSOFields {

			var errors []string

			if getOrderSourceResult.Error != nil {
				fmt.Println(getOrderSourceResult.Error.Error())
				errors = append(errors, getOrderSourceResult.Error.Error())
			}

			checkIfSoRefCodeExist := helper.InSliceString(soRefCodes, v.NoOrder)

			// Get SO Status By Name
			getSalesOrderStatusResultChan := make(chan *models.OrderStatusChan)
			go c.orderStatusRepository.GetByNameAndType("open", "sales_order", false, c.ctx, getSalesOrderStatusResultChan)
			getSalesOrderStatusResult := <-getSalesOrderStatusResultChan

			if getSalesOrderStatusResult.Error != nil {
				fmt.Println(getSalesOrderStatusResult.Error.Error())
				errors = append(errors, getOrderSourceResult.Error.Error())
			}

			// Get SO Detail Status By Name
			getSalesOrderDetailStatusResultChan := make(chan *models.OrderStatusChan)
			go c.orderStatusRepository.GetByNameAndType("open", "sales_order_detail", false, c.ctx, getSalesOrderDetailStatusResultChan)
			getSalesOrderDetailStatusResult := <-getSalesOrderDetailStatusResultChan

			if getSalesOrderDetailStatusResult.Error != nil {
				fmt.Println(getSalesOrderDetailStatusResult.Error.Error())
				errors = append(errors, getOrderSourceResult.Error.Error())
			}

			// Check Product By Id
			getProductResultChan := make(chan *models.ProductChan)
			go c.productRepository.GetBySKU(v.KodeProduk, false, c.ctx, getProductResultChan)
			getProductResult := <-getProductResultChan

			if getProductResult.Error != nil {
				fmt.Println(getProductResult.Error.Error())
				errors = append(errors, fmt.Sprintf("Kode SKU = %s dengan Merek %s Tidak Ditemukan. Silahkan gunakan Kode SKU yang lain", v.KodeProduk, v.NamaMerk))
			}

			// Check Uom By Id
			getUomResultChan := make(chan *models.UomChan)
			go c.uomRepository.GetByCode(v.UnitProduk, false, c.ctx, getUomResultChan)
			getUomResult := <-getUomResultChan

			if getUomResult.Error != nil {
				fmt.Println(getUomResult.Error.Error())
				errors = append(errors, getOrderSourceResult.Error.Error())
			}

			var price float64

			if v.UnitProduk == getProductResult.Product.UnitMeasurementSmall.String {
				price = getProductResult.Product.PriceSmall
			} else if v.UnitProduk == getProductResult.Product.UnitMeasurementMedium.String {
				price = getProductResult.Product.PriceMedium
			} else {
				price = getProductResult.Product.PriceBig
			}

			if price < 1 {
				errors = append(errors, fmt.Sprintf("Produk dengan Kode SKU %s Belum Ada Harga atau Harga = 0. Silahkan gunakan Kode SKU Produk yang lain.", getProductResult.Product.Sku.String))
			}

			if checkIfSoRefCodeExist {

				salesOrder := salesOrderSoRefCodes[v.NoOrder]

				salesOrder.TotalAmount = salesOrder.TotalAmount + (price * float64(v.QTYOrder))
				salesOrder.TotalTonase = salesOrder.TotalTonase + (float64(v.QTYOrder) * getProductResult.Product.NettWeight)

				// ### Sales Order Detail ###
				salesOrderDetail := &models.SalesOrderDetail{}
				salesOrderDetail.SalesOrderDetailUploadSOMap(v, now)
				salesOrderDetail.SalesOrderDetailStatusChanMap(getSalesOrderDetailStatusResult)
				salesOrderDetail.OrderStatusID = getSalesOrderDetailStatusResult.OrderStatus.ID
				salesOrderDetail.Price = price
				salesOrderDetail.ProductID = getProductResult.Product.ID
				salesOrderDetail.Product = getProductResult.Product
				salesOrderDetail.UomID = getUomResult.Uom.ID
				salesOrderDetail.Uom = getUomResult.Uom

				salesOrder.SalesOrderDetails = append(salesOrder.SalesOrderDetails, salesOrderDetail)

				salesOrderSoRefCodes[v.NoOrder] = salesOrder

			} else {

				soRefCodes = append(soRefCodes, v.NoOrder)

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
				go c.storeRepository.GetIdByStoreCode(v.KodeToko, false, c.ctx, getStoreIdResultChan)
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
					errors = append(errors, getAgentResult.Error.Error())
				}

				// Check User By Id
				getUserResultChan := make(chan *models.UserChan)
				go c.userRepository.GetByID(v.IDUser, false, c.ctx, getUserResultChan)
				getUserResult := <-getUserResultChan

				if getUserResult.Error != nil {
					fmt.Println(getUserResult.Error.Error())
					errors = append(errors, getAgentResult.Error.Error())
				}

				// Check Salesman By Id
				getSalesmanResultChan := make(chan *models.SalesmanChan)
				go c.salesmanRepository.GetByID(v.IDSalesman, false, c.ctx, getSalesmanResultChan)
				getSalesmanResult := <-getSalesmanResultChan

				if getSalesmanResult.Error != nil {
					fmt.Println(getSalesmanResult.Error.Error())
					errors = append(errors, getAgentResult.Error.Error())
				}

				// Check Brand By Id
				getBrandResultChan := make(chan *models.BrandChan)
				go c.brandRepository.GetByID(v.KodeMerk, false, c.ctx, getBrandResultChan)
				getBrandResult := <-getBrandResultChan

				if getBrandResult.Error != nil {
					fmt.Println(getBrandResult.Error.Error())
					errors = append(errors, getAgentResult.Error.Error())
				}

				// ### Sales Order ###
				salesOrder := &models.SalesOrder{}

				salesOrder.SalesOrderUploadSOMap(v, now)
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
				salesOrder.LatestUpdatedBy = v.IDUser
				salesOrder.SoCode = helper.GenerateSOCode(v.IDDistributor, getOrderSourceResult.OrderSource.Code)
				salesOrder.StoreID = getStoreResult.Store.ID
				salesOrder.OrderStatusID = getSalesOrderStatusResult.OrderStatus.ID
				salesOrder.OrderSourceID = getOrderSourceResult.OrderSource.ID
				salesOrder.TotalAmount = price * float64(v.QTYOrder)
				salesOrder.TotalTonase = float64(v.QTYOrder) * getProductResult.Product.NettWeight
				salesOrder.SoRefCode = models.NullString{NullString: sql.NullString{String: v.NoOrder, Valid: true}}

				// ### Sales Order Detail ###
				salesOrderDetail := &models.SalesOrderDetail{}
				salesOrderDetail.SalesOrderDetailUploadSOMap(v, now)
				salesOrderDetail.SalesOrderDetailStatusChanMap(getSalesOrderDetailStatusResult)
				salesOrderDetail.OrderStatusID = getSalesOrderDetailStatusResult.OrderStatus.ID
				salesOrderDetail.Price = price
				salesOrderDetail.ProductID = getProductResult.Product.ID
				salesOrderDetail.Product = getProductResult.Product
				salesOrderDetail.UomID = getUomResult.Uom.ID
				salesOrderDetail.Uom = getUomResult.Uom

				salesOrderDetails := []*models.SalesOrderDetail{}
				salesOrderDetails = append(salesOrderDetails, salesOrderDetail)

				salesOrder.SalesOrderDetails = salesOrderDetails

				if len(errors) < 1 {
					salesOrderSoRefCodes[v.NoOrder] = salesOrder
				} else {

					getUploadSOHistoriesResultChan := make(chan *models.SoUploadHistoryChan)
					go c.soUploadHistoriesRepository.GetByID(v.SoUploadHistoryId, false, c.ctx, getUploadSOHistoriesResultChan)
					getUploadSOHistoriesResult := <-getUploadSOHistoriesResultChan
					message := getUploadSOHistoriesResult.SoUploadHistory

					if v.UploadType == "retry" {

						message.Status = constants.UPLOAD_STATUS_HISTORY_FAILED
						uploadSOHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
						go c.soUploadHistoriesRepository.UpdateByID(message.ID.Hex(), message, c.ctx, uploadSOHistoryJourneysResultChan)

						salesOrderSoRefCodes = nil
						break
					} else {

						var myMap map[string]string
						data, _ := json.Marshal(v)
						json.Unmarshal(data, &myMap)

						rowData := &models.RowDataSoUploadErrorLog{}
						rowData.RowDataSoUploadErrorLogMap(myMap)

						soUploadErrorLog := &models.SoUploadErrorLog{}
						soUploadErrorLog.SoUploadErrorLogsMap(i+2, message.ID.Hex(), message.RequestId, message.BulkCode, errors, &now)
						soUploadErrorLog.RowData = *rowData

						soUploadErrorLogssResultChan := make(chan *models.SoUploadErrorLogChan)
						go c.soUploadErrorLogsRepository.Insert(soUploadErrorLog, c.ctx, soUploadErrorLogssResultChan)

						continue
					}
				}
			}
		}

		if salesOrderSoRefCodes == nil {
			continue
		}

		for _, v := range salesOrderSoRefCodes {

			sqlTransaction, _ := c.db.BeginTx(c.ctx, nil)

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

				salesOrderDetailJourneys := &models.SalesOrderDetailJourneys{
					SoDetailId:   createSalesOrderDetailResult.SalesOrderDetail.ID,
					SoDetailCode: soDetailCode,
					Status:       constants.SO_STATUS_OPEN,
					Remark:       "",
					Reason:       "",
					CreatedAt:    &now,
					UpdatedAt:    &now,
				}

				createSalesOrderDetailJourneysResultChan := make(chan *models.SalesOrderDetailJourneysChan)
				go c.salesOrderDetailJourneysRepository.Insert(salesOrderDetailJourneys, c.ctx, createSalesOrderDetailJourneysResultChan)
				createSalesOrderDetailJourneysResult := <-createSalesOrderDetailJourneysResultChan

				if createSalesOrderDetailJourneysResult.Error != nil {
					fmt.Println(createSalesOrderDetailJourneysResult.Error.Error())
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

			salesOrderJourneys := &models.SalesOrderJourneys{
				SoCode:    v.SoCode,
				SoId:      v.ID,
				SoDate:    v.SoDate,
				Status:    constants.SO_STATUS_APPV,
				Remark:    "",
				Reason:    "",
				CreatedAt: &now,
				UpdatedAt: &now,
			}

			createSalesOrderJourneysResultChan := make(chan *models.SalesOrderJourneysChan)
			go c.salesOrderJourneysRepository.Insert(salesOrderJourneys, c.ctx, createSalesOrderJourneysResultChan)
			createSalesOrderJourneysResult := <-createSalesOrderJourneysResultChan

			if createSalesOrderJourneysResult.Error != nil {
				fmt.Println(createSalesOrderJourneysResult.Error.Error())
			}

			keyKafka := []byte(v.SoCode)
			messageKafka, _ := json.Marshal(v)

			err := c.kafkaClient.WriteToTopic(constants.CREATE_SALES_ORDER_TOPIC, keyKafka, messageKafka)

			if err != nil {
				fmt.Println(err.Error())
			}
		}

	}

	if err := reader.Close(); err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
}
