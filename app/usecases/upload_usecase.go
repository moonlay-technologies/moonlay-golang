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

type UploadUseCaseInterface interface {
	UploadSOSJ(ctx context.Context) *model.ErrorLog
	UploadDO(request *models.UploadDORequest, ctx context.Context) *model.ErrorLog
	UploadSO(request *models.UploadSORequest, ctx context.Context) *model.ErrorLog
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

// func (u *uploadUseCase) UploadDO(ctx context.Context) *model.ErrorLog {
// 	now := time.Now()

// 	uploadDOResultChan := make(chan *models.UploadDOFieldsChan)
// 	go u.uploadRepository.UploadDO("be-so-service", "upload-service/do/format-file-upload-data-DO-V2 (1).csv", "ap-southeast-1", uploadDOResultChan)
// 	uploadDOResult := <-uploadDOResultChan

// 	// Get Order Source Status By Name
// 	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
// 	go u.orderSourceRepository.GetBySourceName("upload_sj", false, ctx, getOrderSourceResultChan)
// 	getOrderSourceResult := <-getOrderSourceResultChan

// 	if getOrderSourceResult.Error != nil {
// 		return getOrderSourceResult.ErrorLog
// 	}

// 	for _, v := range uploadDOResult.UploadDOFields {
// 		a, _ := json.Marshal(v)
// 		fmt.Println("test", string(a))

// 		// Get Sales Order By So Code
// 		getSalesOrderResultChan := make(chan *models.SalesOrderChan)
// 		go u.salesOrderRepository.GetByCode(v.NoOrder, false, ctx, getSalesOrderResultChan)
// 		getSalesOrderResult := <-getSalesOrderResultChan

// 		if getSalesOrderResult.Error != nil {
// 			return getSalesOrderResult.ErrorLog
// 		}

// 		errorDOValidation := u.uploadDOValidation(getSalesOrderResult.SalesOrder.ID, getSalesOrderResult.SalesOrder.OrderStatusName, v.KodeProduk, v.Unit, v.QTYShip, ctx)

// 		if errorDOValidation != nil {
// 			errorLogData := helper.NewWriteLog(baseModel.ErrorLog{
// 				Message:       []string{helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_UPLOAD, errorDOValidation.Error())},
// 				SystemMessage: []string{"Invalid Process"},
// 				StatusCode:    http.StatusUnprocessableEntity,
// 			})
// 			return errorLogData
// 		}

// 		// Get Brand By KodeMerk/brand_id
// 		getBrandResultChan := make(chan *models.BrandChan)
// 		go u.brandRepository.GetByID(v.KodeMerk, false, ctx, getBrandResultChan)
// 		getBrandResult := <-getBrandResultChan

// 		// Get Order Status By Name and Type
// 		getOrderStatusResultChan := make(chan *models.OrderStatusChan)
// 		go u.orderStatusRepository.GetByNameAndType("open", "delivery_order", false, ctx, getOrderStatusResultChan)
// 		getOrderStatusResult := <-getOrderStatusResultChan

// 		// Get Warehouse By KodeGudang/id
// 		getWarehouseResultChan := make(chan *models.WarehouseChan)
// 		go u.warehouseRepository.GetByID(v.KodeGudang, false, ctx, getWarehouseResultChan)
// 		getWarehouseResult := <-getWarehouseResultChan
// 		// if v.KodeGudang > 0 {
// 		// 	getWarehouseResultChan := make(chan *models.WarehouseChan)
// 		// 	go u.warehouseRepository.GetByID(v.KodeGudang, false, ctx, getWarehouseResultChan)
// 		// 	getWarehouseResult := <-getWarehouseResultChan

// 		// } else {
// 		// 	getWarehouseResultChan := make(chan *models.WarehouseChan)
// 		// 	go u.warehouseRepository.GetByID(10, false, ctx, getWarehouseResultChan)
// 		// 	getWarehouseResult := <-getWarehouseResultChan

// 		// }

// 		// Get Sales Order Source By ID
// 		// getSalesOrderSourceResultChan := make(chan *models.OrderSourceChan)
// 		// go u.orderSourceRepository.GetByID(getSalesOrderResult.SalesOrder.OrderSourceID, false, ctx, getSalesOrderSourceResultChan)
// 		// getSalesOrderSourceResult := <-getSalesOrderSourceResultChan

// 		// Get Sales Order Status by ID
// 		// getSalesOrderStatusResultChan := make(chan *models.OrderStatusChan)
// 		// go u.orderStatusRepository.GetByID(getSalesOrderResult.SalesOrder.OrderStatusID, false, ctx, getSalesOrderStatusResultChan)
// 		// getSalesOrderStatusResult := <-getSalesOrderStatusResultChan

// 		// Get Agent By IDDistributor/id
// 		getAgentResultChan := make(chan *models.AgentChan)
// 		go u.agentRepository.GetByID(v.IDDistributor, false, ctx, getAgentResultChan)
// 		getAgentResult := <-getAgentResultChan

// 		// Get Store By ID
// 		getStoreResultChan := make(chan *models.StoreChan)
// 		go u.storeRepository.GetByID(getSalesOrderResult.SalesOrder.StoreID, false, ctx, getStoreResultChan)
// 		getStoreResult := <-getStoreResultChan

// 		// Get User By Email
// 		getUserResultChan := make(chan *models.UserChan)
// 		go u.userRepository.GetByID(getSalesOrderResult.SalesOrder.UserID, false, ctx, getUserResultChan)
// 		getUserResult := <-getUserResultChan

// 		getSalesmanResultChan := make(chan *models.SalesmanChan)
// 		go u.salesmanRepository.GetByEmail(getUserResult.User.Email, false, ctx, getSalesmanResultChan)
// 		getSalesmanResult := <-getSalesmanResultChan

// 		deliveryOrder := &models.DeliveryOrder{}
// 		latestUpdatedBy := ctx.Value("user").(*models.UserClaims)
// 		deliveryOrder.SalesOrderID = getSalesOrderResult.SalesOrder.ID
// 		deliveryOrder.DoRefCode = models.NullString{NullString: sql.NullString{String: v.NoSJ, Valid: true}}
// 		deliveryOrder.DoRefDate = models.NullString{NullString: sql.NullString{String: v.TanggalSJ, Valid: true}}
// 		deliveryOrder.DriverName = models.NullString{NullString: sql.NullString{String: v.NamaSupir, Valid: true}}
// 		deliveryOrder.PlatNumber = models.NullString{NullString: sql.NullString{String: v.PlatNo, Valid: true}}
// 		deliveryOrder.Note = models.NullString{NullString: sql.NullString{String: v.Catatan, Valid: true}}
// 		deliveryOrder.IsDoneSyncToEs = "0"
// 		deliveryOrder.StartDateSyncToEs = &now
// 		deliveryOrder.EndDateSyncToEs = &now
// 		deliveryOrder.StartCreatedDate = &now
// 		deliveryOrder.EndCreatedDate = &now
// 		deliveryOrder.LatestUpdatedBy = latestUpdatedBy.UserID
// 		deliveryOrder.CreatedAt = &now
// 		deliveryOrder.UpdatedAt = &now
// 		deliveryOrder.DeletedAt = nil

// 		deliveryOrder.WarehouseChanMap(getWarehouseResult)
// 		deliveryOrder.AgentMap(getAgentResult.Agent)
// 		deliveryOrder.DoCode = helper.GenerateDOCode(getAgentResult.Agent.ID, getOrderSourceResult.OrderSource.Code)
// 		deliveryOrder.DoDate = now.Format("2006-01-02")
// 		deliveryOrder.OrderStatus = getOrderStatusResult.OrderStatus
// 		deliveryOrder.OrderStatusID = getOrderStatusResult.OrderStatus.ID
// 		deliveryOrder.OrderSource = getOrderSourceResult.OrderSource
// 		deliveryOrder.OrderSourceID = getOrderSourceResult.OrderSource.ID
// 		deliveryOrder.Store = getStoreResult.Store
// 		deliveryOrder.CreatedBy = ctx.Value("user").(*models.UserClaims).UserID
// 		deliveryOrder.SalesOrder = getSalesOrderResult.SalesOrder
// 		deliveryOrder.Brand = getBrandResult.Brand
// 		if getSalesmanResult.Salesman != nil {
// 			deliveryOrder.Salesman = getSalesmanResult.Salesman
// 		}

// 		sqlTransaction, _ := u.db.BeginTx(ctx, nil)

// 		// Insert to DB, table delivery_orders
// 		createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
// 		go u.deliveryOrderRepository.Insert(deliveryOrder, sqlTransaction, ctx, createDeliveryOrderResultChan)
// 		createDeliveryOrderResult := <-createDeliveryOrderResultChan

// 		if createDeliveryOrderResult.Error != nil {
// 			sqlTransaction.Rollback()
// 			return createDeliveryOrderResult.ErrorLog
// 		}

// 		// getSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailsChan)
// 		// go u.salesOrderDetailRepository.GetBySalesOrderID(getSalesOrderResult.SalesOrder.ID, false, ctx, getSalesOrderDetailResultChan)
// 		// getSalesOrderDetailResult := <- getSalesOrderDetailResultChan

// 		// for _, x := range getSalesOrderDetailResult.SalesOrderDetails {
// 		// 	x.ID = v
// 		// }
// 	}

// 	// Get Sales Order Detail
// 	// Get SO Detail Product

// 	return nil
// }

func (u *uploadUseCase) UploadDO(request *models.UploadDORequest, ctx context.Context) *model.ErrorLog {

	user := ctx.Value("user").(*models.UserClaims)

	// Check Agent By Id
	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(4, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	message := &models.UploadSJHistory{
		RequestId:       ctx.Value("RequestId").(string),
		FileName:        request.File,
		FilePath:        "upload-service/do/" + request.File,
		AgentId:         int64(user.AgentID),
		AgentName:       getAgentResult.Agent.Name,
		UploadedBy:      int64(user.UserID),
		UploadedByName:  user.FirstName + " " + user.LastName,
		UploadedByEmail: user.UserEmail,
	}

	keyKafka := []byte(ctx.Value("RequestId").(string))
	messageKafka, _ := json.Marshal(message)

	err := u.kafkaClient.WriteToTopic(constants.UPLOAD_DO_FILE_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}

func (u *uploadUseCase) uploadDOValidation(salesOrderId int, orderStatusName string, kodeProduk int, unitProduk string, qtyShip int, ctx context.Context) error {

	if orderStatusName != constants.ORDER_STATUS_OPEN && orderStatusName != constants.ORDER_STATUS_PARTIAL {
		return fmt.Errorf("status sales orders %s, not allowed to upload", orderStatusName)
	}

	getSalesOrderDetailBySoIdResultChan := make(chan *models.SalesOrderDetailsChan)
	go u.salesOrderDetailRepository.GetBySalesOrderID(salesOrderId, false, ctx, getSalesOrderDetailBySoIdResultChan)
	getSalesOrderDetailBySoIdResult := <-getSalesOrderDetailBySoIdResultChan

	if getSalesOrderDetailBySoIdResult.Total == 0 {

		return fmt.Errorf("data sales order detail %d", getSalesOrderDetailBySoIdResult.Total)

	} else {
		for _, v := range getSalesOrderDetailBySoIdResult.SalesOrderDetails {

			getOrderStatusResultChan := make(chan *models.OrderStatusChan)
			go u.orderStatusRepository.GetByID(v.OrderStatusID, false, ctx, getOrderStatusResultChan)
			getOrderStatusResult := <-getOrderStatusResultChan

			if getOrderStatusResult.OrderStatus.Name != constants.ORDER_STATUS_OPEN && getOrderStatusResult.OrderStatus.Name != constants.ORDER_STATUS_PARTIAL {
				return fmt.Errorf("status sales order detail %s, not allowed to upload", getOrderStatusResult.OrderStatus.Name)
			}

			allowedQtyUpload := v.ResidualQty - qtyShip
			if allowedQtyUpload <= 0 {
				return fmt.Errorf("residual qty equal or less than zero")
			}

			getProductResultChan := make(chan *models.ProductChan)
			go u.productRepository.GetByID(v.ProductID, false, ctx, getProductResultChan)
			getProductResult := <-getProductResultChan

			productSKU, err := strconv.Atoi(getProductResult.Product.Sku.String)
			if err != nil {
				return fmt.Errorf("fail to convert productSKU from string to int")
			}

			if productSKU != kodeProduk {
				return fmt.Errorf("productSKU %d from database doesn't match with kodeProduk %d to be uploaded", productSKU, kodeProduk)
			}

			getUomResultChan := make(chan *models.UomChan)
			go u.uomRepository.GetByID(v.UomID, false, ctx, getUomResultChan)
			getUomResult := <-getUomResultChan

			if getUomResult.Uom.Code.String != unitProduk {
				return fmt.Errorf("uom code %s from database doesn't match with unitProduk %s to be uploaded", getUomResult.Uom.Code.String, unitProduk)
			}
		}
	}

	return nil
}

func (u *uploadUseCase) UploadSO(request *models.UploadSORequest, ctx context.Context) *model.ErrorLog {

	user := ctx.Value("user").(*models.UserClaims)

	// Check Agent By Id
	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(4, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	message := &models.UploadSOHistory{
		RequestId:       ctx.Value("RequestId").(string),
		FileName:        request.File,
		FilePath:        "upload-service/so/" + request.File,
		AgentId:         int64(user.AgentID),
		AgentName:       getAgentResult.Agent.Name,
		UploadedBy:      int64(user.UserID),
		UploadedByName:  user.FirstName + " " + user.LastName,
		UploadedByEmail: user.UserEmail,
	}

	keyKafka := []byte(ctx.Value("RequestId").(string))
	messageKafka, _ := json.Marshal(message)

	err := u.kafkaClient.WriteToTopic(constants.UPLOAD_SO_FILE_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}
