package usecases

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

type UploadUseCaseInterface interface {
	UploadSOSJ(ctx context.Context) *model.ErrorLog
	UploadDO(ctx context.Context) *model.ErrorLog
	UploadSO(ctx context.Context) *model.ErrorLog
}

type uploadUseCase struct {
	uploadRepository                     repositories.UploadRepositoryInterface
	salesOrderRepository                 repositories.SalesOrderRepositoryInterface
	salesOrderDetailRepository           repositories.SalesOrderDetailRepositoryInterface
	orderStatusRepository                repositories.OrderStatusRepositoryInterface
	orderSourceRepository                repositories.OrderSourceRepositoryInterface
	agentRepository                      repositories.AgentRepositoryInterface
	brandRepository                      repositories.BrandRepositoryInterface
	storeRepository                      repositories.StoreRepositoryInterface
	productRepository                    repositories.ProductRepositoryInterface
	uomRepository                        repositories.UomRepositoryInterface
	deliveryOrderRepository              repositories.DeliveryOrderRepositoryInterface
	deliveryOrderDetailRepository        repositories.DeliveryOrderDetailRepositoryInterface
	salesOrderLogRepository              mongoRepositories.SalesOrderLogRepositoryInterface
	salesOrderJourneysRepository         mongoRepositories.SalesOrderJourneysRepositoryInterface
	salesOrderDetailJourneysRepository   mongoRepositories.SalesOrderDetailJourneysRepositoryInterface
	userRepository                       repositories.UserRepositoryInterface
	salesmanRepository                   repositories.SalesmanRepositoryInterface
	categoryRepository                   repositories.CategoryRepositoryInterface
	salesOrderOpenSearchRepository       openSearchRepositories.SalesOrderOpenSearchRepositoryInterface
	salesOrderDetailOpenSearchRepository openSearchRepositories.SalesOrderDetailOpenSearchRepositoryInterface
	warehouseRepository                  repositories.WarehouseRepositoryInterface
	kafkaClient                          kafkadbo.KafkaClientInterface
	db                                   dbresolver.DB
	ctx                                  context.Context
}

func InitUploadUseCaseInterface(salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, orderSourceRepository repositories.OrderSourceRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, uomRepository repositories.UomRepositoryInterface, deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, deliveryOrderDetailRepository repositories.DeliveryOrderDetailRepositoryInterface, salesOrderLogRepository mongoRepositories.SalesOrderLogRepositoryInterface, salesOrderJourneysRepository mongoRepositories.SalesOrderJourneysRepositoryInterface, salesOrderDetailJourneysRepository mongoRepositories.SalesOrderDetailJourneysRepositoryInterface, userRepository repositories.UserRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, categoryRepository repositories.CategoryRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, salesOrderDetailOpenSearchRepository openSearchRepositories.SalesOrderDetailOpenSearchRepositoryInterface, uploadRepository repositories.UploadRepositoryInterface, warehouseRepository repositories.WarehouseRepositoryInterface, kafkaClient kafkadbo.KafkaClientInterface, db dbresolver.DB, ctx context.Context) UploadUseCaseInterface {
	return &uploadUseCase{
		salesOrderRepository:                 salesOrderRepository,
		salesOrderDetailRepository:           salesOrderDetailRepository,
		orderStatusRepository:                orderStatusRepository,
		orderSourceRepository:                orderSourceRepository,
		agentRepository:                      agentRepository,
		brandRepository:                      brandRepository,
		storeRepository:                      storeRepository,
		productRepository:                    productRepository,
		uomRepository:                        uomRepository,
		deliveryOrderRepository:              deliveryOrderRepository,
		deliveryOrderDetailRepository:        deliveryOrderDetailRepository,
		salesOrderLogRepository:              salesOrderLogRepository,
		salesOrderJourneysRepository:         salesOrderJourneysRepository,
		salesOrderDetailJourneysRepository:   salesOrderDetailJourneysRepository,
		userRepository:                       userRepository,
		salesmanRepository:                   salesmanRepository,
		categoryRepository:                   categoryRepository,
		salesOrderOpenSearchRepository:       salesOrderOpenSearchRepository,
		salesOrderDetailOpenSearchRepository: salesOrderDetailOpenSearchRepository,
		uploadRepository:                     uploadRepository,
		warehouseRepository:                  warehouseRepository,
		kafkaClient:                          kafkaClient,
		db:                                   db,
		ctx:                                  ctx,
	}
}

func (u *uploadUseCase) UploadSOSJ(ctx context.Context) *model.ErrorLog {
	now := time.Now()
	var user_id = ctx.Value("user").(*models.UserClaims).UserID

	uploadSOSJResultChan := make(chan *models.UploadSOSJFieldChan)
	go u.uploadRepository.UploadSOSJ("be-so-service", "upload-service/sosj/format-file-upload-data-SOSJ-V2.csv", "ap-southeast-1", user_id, uploadSOSJResultChan)
	uploadSOSJResult := <-uploadSOSJResultChan

	// Get Order Source Status By Id
	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetBySourceName("upload_sosj", false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getOrderSourceResult.Error != nil {
		return getOrderSourceResult.ErrorLog
	}

	soRefCodes := []string{}
	salesOrderSoRefCodes := map[string]*models.SalesOrder{}

	salesOrders := []*models.SalesOrder{}
	for _, v := range uploadSOSJResult.UploadSOSJFields {

		var noSuratJalan = v.NoSuratJalan

		checkIfSoRefCodeExist := helper.InSliceString(soRefCodes, noSuratJalan)

		var soStatus string
		var doStatus string
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
		go u.orderStatusRepository.GetByNameAndType(soStatus, "sales_order", false, ctx, getSalesOrderStatusResultChan)
		getSalesOrderStatusResult := <-getSalesOrderStatusResultChan

		// Get SO Status By Name
		getSalesOrderDetailStatusResultChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByNameAndType(soStatus, "sales_order_detail", false, ctx, getSalesOrderDetailStatusResultChan)
		getSalesOrderDetailStatusResult := <-getSalesOrderDetailStatusResultChan

		// Get DO Status By Name
		getDeliveryOrderStatusResultChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByNameAndType(doStatus, "delivery_order", false, ctx, getDeliveryOrderStatusResultChan)
		getDeliveryOrderStatusResult := <-getDeliveryOrderStatusResultChan

		// if getOrderStatusResult.Error != nil {
		// 	return getOrderStatusResult.ErrorLog
		// }

		// Check Product By Id
		getProductResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetByID(v.KodeProdukDBO, false, ctx, getProductResultChan)
		getProductResult := <-getProductResultChan

		// if getProductResult.Error != nil {
		// 	errorLogData := helper.WriteLog(getProductResult.Error, http.StatusInternalServerError, nil)
		// 	return errorLogData
		// }

		// Check Uom By Id
		getUomResultChan := make(chan *models.UomChan)
		go u.uomRepository.GetByID(v.Unit, false, ctx, getUomResultChan)
		getUomResult := <-getUomResultChan

		// if getUomResult.Error != nil {
		// 	errorLogData := helper.WriteLog(getUomResult.Error, http.StatusInternalServerError, nil)
		// 	return errorLogData
		// }

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
			salesOrderDetail.OrderStatusID = getSalesOrderDetailStatusResult.OrderStatus.ID
			salesOrderDetail.Price = price
			salesOrderDetail.Product = getProductResult.Product
			salesOrderDetail.Uom = getUomResult.Uom

			salesOrder.SalesOrderDetails = append(salesOrder.SalesOrderDetails, salesOrderDetail)

			// ### Delivery Order Detail ###
			deliveryOrderDetail := &models.DeliveryOrderDetail{}
			deliveryOrderDetail.DeliveryOrderDetailUploadSOSJMap(v, now)

			deliveryOrderDetail.BrandID = v.KodeProdukDBO
			deliveryOrderDetail.Note = models.NullString{NullString: sql.NullString{String: v.Catatan, Valid: true}}
			deliveryOrderDetail.UomID = v.Unit
			deliveryOrderDetail.Uom = getUomResult.Uom
			deliveryOrderDetail.ProductID = v.KodeProdukDBO
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
			go u.agentRepository.GetByID(v.IDDistributor, false, ctx, getAgentResultChan)
			getAgentResult := <-getAgentResultChan

			// if getAgentResult.Error != nil {
			// 	errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
			// 	return errorLogData
			// }

			// Check Store By Id
			getStoreResultChan := make(chan *models.StoreChan)
			go u.storeRepository.GetByID(v.KodeTokoDBO, false, ctx, getStoreResultChan)
			getStoreResult := <-getStoreResultChan

			// if getStoreResult.Error != nil {
			// 	errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
			// 	return errorLogData
			// }

			// Check User By Id
			getUserResultChan := make(chan *models.UserChan)
			go u.userRepository.GetByID(user_id, false, ctx, getUserResultChan)
			getUserResult := <-getUserResultChan

			// if getUserResult.Error != nil {
			// 	errorLogData := helper.WriteLog(getUserResult.Error, http.StatusInternalServerError, nil)
			// 	return errorLogData
			// }

			// Check Salesman By Id
			getSalesmanResult := &models.SalesmanChan{}
			if v.IDSalesman > 0 {
				getSalesmanResultChan := make(chan *models.SalesmanChan)
				go u.salesmanRepository.GetByID(v.IDSalesman, false, ctx, getSalesmanResultChan)
				getSalesmanResult = <-getSalesmanResultChan

				// if getSalesmanResult.Error != nil {
				// 	errorLogData := helper.WriteLog(getSalesmanResult.Error, http.StatusInternalServerError, nil)
				// 	return errorLogData
				// }
			}

			// Check Brand By Id
			getBrandResultChan := make(chan *models.BrandChan)
			go u.brandRepository.GetByID(v.IDMerk, false, ctx, getBrandResultChan)
			getBrandResult := <-getBrandResultChan

			// if getBrandResult.Error != nil {
			// 	return getBrandResult.ErrorLog
			// }

			getWarehouseResultChan := make(chan *models.WarehouseChan)
			go u.warehouseRepository.GetByID(v.KodeGudang, false, ctx, getWarehouseResultChan)
			getWarehouseResult := <-getWarehouseResultChan

			// if getWarehouseResult.Error != nil {
			// 	return getWarehouseResult.ErrorLog
			// }

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

			salesOrder.UserID = user_id
			salesOrder.CreatedBy = user_id
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
			deliveryOrder.CreatedBy = ctx.Value("user").(*models.UserClaims).UserID
			deliveryOrder.LatestUpdatedBy = ctx.Value("user").(*models.UserClaims).UserID
			deliveryOrder.Brand = getBrandResult.Brand
			if getSalesmanResult.Salesman != nil {
				deliveryOrder.Salesman = getSalesmanResult.Salesman
			}

			// ### Delivery Order Detail ###
			deliveryOrderDetail := &models.DeliveryOrderDetail{}
			deliveryOrderDetail.DeliveryOrderDetailUploadSOSJMap(v, now)

			deliveryOrderDetails := []*models.DeliveryOrderDetail{}
			deliveryOrderDetails = append(deliveryOrderDetails, deliveryOrderDetail)
			deliveryOrderDetail.BrandID = v.KodeProdukDBO
			deliveryOrderDetail.Note = models.NullString{NullString: sql.NullString{String: v.Catatan, Valid: true}}
			deliveryOrderDetail.UomID = v.Unit
			deliveryOrderDetail.Uom = getUomResult.Uom
			deliveryOrderDetail.ProductID = v.KodeProdukDBO
			deliveryOrderDetail.ProductChanMap(getProductResult)
			deliveryOrderDetail.OrderStatusID = getDeliveryOrderStatusResult.OrderStatus.ID
			deliveryOrderDetail.OrderStatusName = getDeliveryOrderStatusResult.OrderStatus.Name
			deliveryOrderDetail.OrderStatus = getDeliveryOrderStatusResult.OrderStatus

			deliveryOrder.DeliveryOrderDetails = deliveryOrderDetails

			deliveryOrders := []*models.DeliveryOrder{}
			deliveryOrders = append(deliveryOrders, deliveryOrder)

			salesOrder.DeliveryOrders = append(salesOrder.DeliveryOrders, deliveryOrder)

			salesOrderSoRefCodes[noSuratJalan] = salesOrder

			salesOrders = append(salesOrders, salesOrder)
		}
	}

	for _, v := range salesOrderSoRefCodes {

		sqlTransaction, _ := u.db.BeginTx(ctx, nil)

		// if err != nil {
		// 	errorLog := helper.WriteLog(err, http.StatusInternalServerError, nil)

		// 	return errorLog
		// }

		createSalesOrderResultChan := make(chan *models.SalesOrderChan)
		go u.salesOrderRepository.Insert(v, sqlTransaction, ctx, createSalesOrderResultChan)
		createSalesOrderResult := <-createSalesOrderResultChan

		if createSalesOrderResult.Error != nil {
			sqlTransaction.Rollback()
			return createSalesOrderResult.ErrorLog
		}

		v.ID = createSalesOrderResult.SalesOrder.ID

		for _, x := range v.SalesOrderDetails {

			soDetailCode, _ := helper.GenerateSODetailCode(int(createSalesOrderResult.ID), v.AgentID, x.ProductID, x.UomID)
			x.SalesOrderID = int(createSalesOrderResult.ID)
			x.SoDetailCode = soDetailCode

			createSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
			go u.salesOrderDetailRepository.Insert(x, sqlTransaction, ctx, createSalesOrderDetailResultChan)
			createSalesOrderDetailResult := <-createSalesOrderDetailResultChan

			if createSalesOrderDetailResult.Error != nil {
				sqlTransaction.Rollback()
				return createSalesOrderDetailResult.ErrorLog
			}

			x.ID = createSalesOrderDetailResult.SalesOrderDetail.ID

			// salesOrderDetailJourneys := &models.SalesOrderDetailJourneys{
			// 	SoDetailId:   createSalesOrderDetailResult.SalesOrderDetail.ID,
			// 	SoDetailCode: soDetailCode,
			// 	Status:       constants.ORDER_STATUS_OPEN,
			// 	Remark:       "",
			// 	Reason:       "",
			// 	CreatedAt:    &now,
			// 	UpdatedAt:    &now,
			// }

			// createSalesOrderDetailJourneysResultChan := make(chan *models.SalesOrderDetailJourneysChan)
			// go u.salesOrderDetailJourneysRepository.Insert(salesOrderDetailJourneys, ctx, createSalesOrderDetailJourneysResultChan)
			// createSalesOrderDetailJourneysResult := <-createSalesOrderDetailJourneysResultChan

			// if createSalesOrderDetailJourneysResult.Error != nil {
			// 	return createSalesOrderDetailJourneysResult.ErrorLog
			// }
		}

		for _, x := range v.DeliveryOrders {
			x.SalesOrderID = v.ID
			x.SalesOrder = createSalesOrderResult.SalesOrder

			createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
			go u.deliveryOrderRepository.Insert(x, sqlTransaction, ctx, createDeliveryOrderResultChan)
			createDeliveryOrderResult := <-createDeliveryOrderResultChan

			if createDeliveryOrderResult.Error != nil {
				sqlTransaction.Rollback()
				// return createDeliveryOrderResult.ErrorLog
			}

			x.ID = createDeliveryOrderResult.DeliveryOrder.ID

			for i, doDetail := range x.DeliveryOrderDetails {

				doDetailCode, _ := helper.GenerateDODetailCode(createDeliveryOrderResult.DeliveryOrder.ID, v.AgentID, v.SalesOrderDetails[i].Product.ID, v.SalesOrderDetails[i].Uom.ID)

				doDetail.DeliveryOrderID = createDeliveryOrderResult.DeliveryOrder.ID
				doDetail.SoDetailID = v.SalesOrderDetails[i].ID
				doDetail.DoDetailCode = doDetailCode

				createDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
				go u.deliveryOrderDetailRepository.Insert(doDetail, sqlTransaction, ctx, createDeliveryOrderDetailResultChan)
				createDeliveryOrderDetailResult := <-createDeliveryOrderDetailResultChan

				if createDeliveryOrderDetailResult.Error != nil {
					sqlTransaction.Rollback()
					// return createDeliveryOrderDetailResult.ErrorLog
				}

				doDetail.ID = createDeliveryOrderDetailResult.DeliveryOrderDetail.ID
			}
		}

		sqlTransaction.Commit()

	}

	return nil
}

func (u *uploadUseCase) UploadDO(ctx context.Context) *model.ErrorLog {

	uploadDOResultChan := make(chan *models.UploadDOFieldsChan)
	go u.uploadRepository.UploadDO("be-so-service", "upload-service/do/format-file-upload-data-DO-V2 (1).csv", "ap-southeast-1", uploadDOResultChan)
	uploadDOResult := <-uploadDOResultChan

	for _, v := range uploadDOResult.UploadDOFields {
		a, _ := json.Marshal(v)
		fmt.Println(string(a))
	}
	return nil
}

func (u *uploadUseCase) UploadSO(ctx context.Context) *model.ErrorLog {
	now := time.Now()
	var user_id = ctx.Value("user").(*models.UserClaims).UserID

	uploadSOResultChan := make(chan *models.UploadSOFieldsChan)
	go u.uploadRepository.UploadSO("be-so-service", "upload-service/so/format-file-upload-data-SO-V2 (2).csv", "ap-southeast-1", uploadSOResultChan)
	uploadSOResult := <-uploadSOResultChan

	// Get Order Source Status By Id
	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetBySourceName("upload_sosj", false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getOrderSourceResult.Error != nil {
		return getOrderSourceResult.ErrorLog
	}

	soRefCodes := []string{}
	salesOrderSoRefCodes := map[string]*models.SalesOrder{}

	for _, v := range uploadSOResult.UploadSOFields {

		checkIfSoRefCodeExist := helper.InSliceString(soRefCodes, v.NoOrder)

		// Get SO Status By Name
		getSalesOrderStatusResultChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByNameAndType("open", "sales_order", false, ctx, getSalesOrderStatusResultChan)
		getSalesOrderStatusResult := <-getSalesOrderStatusResultChan

		// Get SO Status By Name
		getSalesOrderDetailStatusResultChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByNameAndType("open", "sales_order_detail", false, ctx, getSalesOrderDetailStatusResultChan)
		getSalesOrderDetailStatusResult := <-getSalesOrderDetailStatusResultChan

		// Check Product By Id
		getProductResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetBySKU(v.KodeProduk, false, ctx, getProductResultChan)
		getProductResult := <-getProductResultChan

		// if getProductResult.Error != nil {
		// 	errorLogData := helper.WriteLog(getProductResult.Error, http.StatusInternalServerError, nil)
		// 	return errorLogData
		// }

		// Check Uom By Id
		getUomResultChan := make(chan *models.UomChan)
		go u.uomRepository.GetByCode(v.UnitProduk, false, ctx, getUomResultChan)
		getUomResult := <-getUomResultChan

		// if getUomResult.Error != nil {
		// 	errorLogData := helper.WriteLog(getUomResult.Error, http.StatusInternalServerError, nil)
		// 	return errorLogData
		// }

		var price float64

		if v.UnitProduk == getProductResult.Product.UnitMeasurementSmall.String {
			price = getProductResult.Product.PriceSmall
		} else if v.UnitProduk == getProductResult.Product.UnitMeasurementMedium.String {
			price = getProductResult.Product.PriceMedium
		} else {
			price = getProductResult.Product.PriceBig
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
			go u.agentRepository.GetByID(v.IDDistributor, false, ctx, getAgentResultChan)
			getAgentResult := <-getAgentResultChan

			// if getAgentResult.Error != nil {
			// 	errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
			// 	return errorLogData
			// }

			// Check Store By Id
			getStoreResultChan := make(chan *models.StoreChan)
			go u.storeRepository.GetByID(v.KodeToko, false, ctx, getStoreResultChan)
			getStoreResult := <-getStoreResultChan

			// if getStoreResult.Error != nil {
			// 	errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
			// 	return errorLogData
			// }

			// Check User By Id
			getUserResultChan := make(chan *models.UserChan)
			go u.userRepository.GetByID(user_id, false, ctx, getUserResultChan)
			getUserResult := <-getUserResultChan

			// if getUserResult.Error != nil {
			// 	errorLogData := helper.WriteLog(getUserResult.Error, http.StatusInternalServerError, nil)
			// 	return errorLogData
			// }

			// Check Salesman By Id

			getSalesmanResultChan := make(chan *models.SalesmanChan)
			go u.salesmanRepository.GetByID(v.IDSalesman, false, ctx, getSalesmanResultChan)
			getSalesmanResult := <-getSalesmanResultChan

			// if getSalesmanResult.Error != nil {
			// 	errorLogData := helper.WriteLog(getSalesmanResult.Error, http.StatusInternalServerError, nil)
			// 	return errorLogData
			// }

			// Check Brand By Id
			getBrandResultChan := make(chan *models.BrandChan)
			go u.brandRepository.GetByID(v.KodeMerk, false, ctx, getBrandResultChan)
			getBrandResult := <-getBrandResultChan

			// if getBrandResult.Error != nil {
			// 	return getBrandResult.ErrorLog
			// }

			// ### Sales Order ###
			salesOrder := &models.SalesOrder{}
			// soRefCodes = append(soRefCodes, noSuratJalan)

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

			salesOrder.UserID = user_id
			salesOrder.CreatedBy = user_id
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

			salesOrderSoRefCodes[v.NoOrder] = salesOrder
		}
	}

	for _, v := range salesOrderSoRefCodes {

		sqlTransaction, _ := u.db.BeginTx(ctx, nil)

		// if err != nil {
		// 	errorLog := helper.WriteLog(err, http.StatusInternalServerError, nil)

		// 	return errorLog
		// }

		createSalesOrderResultChan := make(chan *models.SalesOrderChan)
		go u.salesOrderRepository.Insert(v, sqlTransaction, ctx, createSalesOrderResultChan)
		createSalesOrderResult := <-createSalesOrderResultChan

		if createSalesOrderResult.Error != nil {
			sqlTransaction.Rollback()
			return createSalesOrderResult.ErrorLog
		}

		v.ID = createSalesOrderResult.SalesOrder.ID

		for _, x := range v.SalesOrderDetails {

			soDetailCode, _ := helper.GenerateSODetailCode(int(createSalesOrderResult.ID), v.AgentID, x.ProductID, x.UomID)
			x.SalesOrderID = int(createSalesOrderResult.ID)
			x.SoDetailCode = soDetailCode

			createSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
			go u.salesOrderDetailRepository.Insert(x, sqlTransaction, ctx, createSalesOrderDetailResultChan)
			createSalesOrderDetailResult := <-createSalesOrderDetailResultChan

			if createSalesOrderDetailResult.Error != nil {
				sqlTransaction.Rollback()
				return createSalesOrderDetailResult.ErrorLog
			}

			x.ID = createSalesOrderDetailResult.SalesOrderDetail.ID

			// salesOrderDetailJourneys := &models.SalesOrderDetailJourneys{
			// 	SoDetailId:   createSalesOrderDetailResult.SalesOrderDetail.ID,
			// 	SoDetailCode: soDetailCode,
			// 	Status:       constants.ORDER_STATUS_OPEN,
			// 	Remark:       "",
			// 	Reason:       "",
			// 	CreatedAt:    &now,
			// 	UpdatedAt:    &now,
			// }

			// createSalesOrderDetailJourneysResultChan := make(chan *models.SalesOrderDetailJourneysChan)
			// go u.salesOrderDetailJourneysRepository.Insert(salesOrderDetailJourneys, ctx, createSalesOrderDetailJourneysResultChan)
			// createSalesOrderDetailJourneysResult := <-createSalesOrderDetailJourneysResultChan

			// if createSalesOrderDetailJourneysResult.Error != nil {
			// 	return createSalesOrderDetailJourneysResult.ErrorLog
			// }
		}

		sqlTransaction.Commit()
	}

	return nil
}
