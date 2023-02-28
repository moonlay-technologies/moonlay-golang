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
)

type SalesOrderUseCaseInterface interface {
	Create(request *models.SalesOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) ([]*models.SalesOrderResponse, *model.ErrorLog)
	Get(request *models.SalesOrderRequest) (*models.SalesOrdersOpenSearchResponse, *model.ErrorLog)
	GetByID(request *models.SalesOrderRequest, ctx context.Context) ([]*models.SalesOrderOpenSearchResponse, *model.ErrorLog)
	GetByIDWithDetail(request *models.SalesOrderRequest, ctx context.Context) (*models.SalesOrder, *model.ErrorLog)
	GetByAgentID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetByStoreID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetBySalesmanID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetByOrderStatusID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetByOrderSourceID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	UpdateById(id int, request *models.SalesOrderUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderResponse, *model.ErrorLog)
	UpdateSODetailById(soId, id int, request *models.SalesOrderDetailUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderDetail, *model.ErrorLog)
	UpdateSODetailBySOId(SoId int, request []*models.SalesOrderDetailUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) ([]*models.SalesOrder, *model.ErrorLog)
	GetDetails(request *models.SalesOrderRequest) (*models.SalesOrderDetailsOpenSearchResponse, *model.ErrorLog)
	GetDetailById(id int) (*models.SalesOrderDetailOpenSearchResponse, *model.ErrorLog)
	DeleteById(id int, sqlTransaction *sql.Tx) *model.ErrorLog
}

type salesOrderUseCase struct {
	salesOrderRepository               repositories.SalesOrderRepositoryInterface
	salesOrderDetailRepository         repositories.SalesOrderDetailRepositoryInterface
	orderStatusRepository              repositories.OrderStatusRepositoryInterface
	orderSourceRepository              repositories.OrderSourceRepositoryInterface
	agentRepository                    repositories.AgentRepositoryInterface
	brandRepository                    repositories.BrandRepositoryInterface
	storeRepository                    repositories.StoreRepositoryInterface
	productRepository                  repositories.ProductRepositoryInterface
	uomRepository                      repositories.UomRepositoryInterface
	deliveryOrderRepository            repositories.DeliveryOrderRepositoryInterface
	salesOrderLogRepository            mongoRepositories.SalesOrderLogRepositoryInterface
	salesOrderJourneysRepository       mongoRepositories.SalesOrderJourneysRepositoryInterface
	salesOrderDetailJourneysRepository mongoRepositories.SalesOrderDetailJourneysRepositoryInterface
	userRepository                     repositories.UserRepositoryInterface
	salesmanRepository                 repositories.SalesmanRepositoryInterface
	categoryRepository                 repositories.CategoryRepositoryInterface
	salesOrderOpenSearchRepository     openSearchRepositories.SalesOrderOpenSearchRepositoryInterface
	deliveryOrderOpenSearchRepository  openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface
	kafkaClient                        kafkadbo.KafkaClientInterface
	db                                 dbresolver.DB
	ctx                                context.Context
}

func InitSalesOrderUseCaseInterface(salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, orderSourceRepository repositories.OrderSourceRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, uomRepository repositories.UomRepositoryInterface, deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, salesOrderLogRepository mongoRepositories.SalesOrderLogRepositoryInterface, salesOrderJourneysRepository mongoRepositories.SalesOrderJourneysRepositoryInterface, salesOrderDetailJourneysRepository mongoRepositories.SalesOrderDetailJourneysRepositoryInterface, userRepository repositories.UserRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, categoryRepository repositories.CategoryRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface, kafkaClient kafkadbo.KafkaClientInterface, db dbresolver.DB, ctx context.Context) SalesOrderUseCaseInterface {
	return &salesOrderUseCase{
		salesOrderRepository:               salesOrderRepository,
		salesOrderDetailRepository:         salesOrderDetailRepository,
		orderStatusRepository:              orderStatusRepository,
		orderSourceRepository:              orderSourceRepository,
		agentRepository:                    agentRepository,
		brandRepository:                    brandRepository,
		storeRepository:                    storeRepository,
		productRepository:                  productRepository,
		uomRepository:                      uomRepository,
		deliveryOrderRepository:            deliveryOrderRepository,
		salesOrderLogRepository:            salesOrderLogRepository,
		salesOrderJourneysRepository:       salesOrderJourneysRepository,
		salesOrderDetailJourneysRepository: salesOrderDetailJourneysRepository,
		userRepository:                     userRepository,
		salesmanRepository:                 salesmanRepository,
		categoryRepository:                 categoryRepository,
		salesOrderOpenSearchRepository:     salesOrderOpenSearchRepository,
		deliveryOrderOpenSearchRepository:  deliveryOrderOpenSearchRepository,
		kafkaClient:                        kafkaClient,
		db:                                 db,
		ctx:                                ctx,
	}
}

func (u *salesOrderUseCase) Create(request *models.SalesOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) ([]*models.SalesOrderResponse, *model.ErrorLog) {
	now := time.Now()
	var soCode string

	// Get Order Status By Name
	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType("open", "sales_order", false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return []*models.SalesOrderResponse{}, getOrderStatusResult.ErrorLog
	}

	// Get Order Detail Status By Name
	getOrderDetailStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType("open", "sales_order_detail", false, ctx, getOrderDetailStatusResultChan)
	getOrderDetailStatusResult := <-getOrderDetailStatusResultChan

	if getOrderDetailStatusResult.Error != nil {
		return []*models.SalesOrderResponse{}, getOrderDetailStatusResult.ErrorLog
	}

	// Get Order Source Status By Id
	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(request.OrderSourceID, false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getOrderSourceResult.Error != nil {
		return []*models.SalesOrderResponse{}, getOrderSourceResult.ErrorLog
	}

	parseSoDate, _ := time.Parse("2006-01-02", request.SoDate)
	parseSoRefDate, _ := time.Parse("2006-01-02", request.SoRefDate)
	duration := time.Hour*time.Duration(now.Hour()) + time.Minute*time.Duration(now.Minute()) + time.Second*time.Duration(now.Second()) + time.Nanosecond*time.Duration(now.Nanosecond())

	soDate := parseSoDate.Local().Add(duration)
	soRefDate := parseSoRefDate.Local().Add(duration)
	sourceName := getOrderSourceResult.OrderSource.SourceName

	if sourceName == "manager" && !(soDate.Add(1*time.Hour).After(soRefDate) && soRefDate.Add(-1*time.Hour).Before(now.Local()) && soDate.Add(-1*time.Hour).Before(now.Local()) && soRefDate.Month() == now.Local().Month() && soRefDate.Local().Year() == now.Local().Year()) {

		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage("create", "so_date dan so_ref_date harus sama dengan kurang dari hari ini dan harus di bulan berjalan")},
			SystemMessage: []string{"Invalid Process"},
			StatusCode:    http.StatusUnprocessableEntity,
		})

		return []*models.SalesOrderResponse{}, errorLog

	} else if (sourceName == "salesman" || sourceName == "store") && !(soDate.Equal(now.Local()) && soRefDate.Equal(now.Local())) {

		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage("create", "so_date dan so_ref_date harus sama dengan hari ini")},
			SystemMessage: []string{"Invalid Process"},
			StatusCode:    http.StatusUnprocessableEntity,
		})

		return []*models.SalesOrderResponse{}, errorLog

	}

	// Check Agent By Id
	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(request.AgentID, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return []*models.SalesOrderResponse{}, errorLogData
	}

	// Check Store By Id
	getStoreResultChan := make(chan *models.StoreChan)
	go u.storeRepository.GetByID(request.StoreID, false, ctx, getStoreResultChan)
	getStoreResult := <-getStoreResultChan

	if getStoreResult.Error != nil {
		errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
		return []*models.SalesOrderResponse{}, errorLogData
	}

	// Check User By Id
	getUserResultChan := make(chan *models.UserChan)
	go u.userRepository.GetByID(request.UserID, false, ctx, getUserResultChan)
	getUserResult := <-getUserResultChan

	if getUserResult.Error != nil {
		errorLogData := helper.WriteLog(getUserResult.Error, http.StatusInternalServerError, nil)
		return []*models.SalesOrderResponse{}, errorLogData
	}

	// Check Salesman By Id
	getSalesmanResult := &models.SalesmanChan{}
	if request.SalesmanID > 0 {
		getSalesmanResultChan := make(chan *models.SalesmanChan)
		go u.salesmanRepository.GetByID(request.SalesmanID, false, ctx, getSalesmanResultChan)
		getSalesmanResult = <-getSalesmanResultChan

		if getSalesmanResult.Error != nil {
			errorLogData := helper.WriteLog(getSalesmanResult.Error, http.StatusInternalServerError, nil)
			return []*models.SalesOrderResponse{}, errorLogData
		}
	}

	brandIds := []int{}
	var salesOrderBrands map[int]*models.SalesOrder
	salesOrderBrands = map[int]*models.SalesOrder{}

	for _, v := range request.SalesOrderDetails {

		checkIfBrandExist := helper.InSliceInt(brandIds, v.BrandID)

		getProductResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetByID(v.ProductID, false, ctx, getProductResultChan)
		getProductResult := <-getProductResultChan

		if getProductResult.Error != nil {
			errorLogData := helper.WriteLog(getProductResult.Error, http.StatusInternalServerError, nil)
			return []*models.SalesOrderResponse{}, errorLogData
		}

		if checkIfBrandExist {

			salesOrderDetail := &models.SalesOrderDetail{}
			salesOrderDetail.SalesOrderDetailStoreRequestMap(v, now)
			salesOrderDetail.Note = models.NullString{NullString: sql.NullString{String: request.Note, Valid: true}}
			salesOrderDetail.OrderStatusID = getOrderDetailStatusResult.OrderStatus.ID

			salesOrderBrand := salesOrderBrands[v.BrandID]
			salesOrderBrand.TotalAmount = salesOrderBrand.TotalAmount + (v.Price * float64(v.Qty))
			salesOrderBrand.TotalTonase = salesOrderBrand.TotalTonase + (float64(v.Qty) * getProductResult.Product.NettWeight)

			salesOrderDetails := salesOrderBrand.SalesOrderDetails
			salesOrderDetails = append(salesOrderDetails, salesOrderDetail)
			salesOrderBrand.SalesOrderDetails = salesOrderDetails
			salesOrderBrands[v.BrandID] = salesOrderBrand

		} else {

			soCode = helper.GenerateSOCode(request.AgentID, getOrderSourceResult.OrderSource.Code)
			brandIds = append(brandIds, v.BrandID)

			salesOrder := &models.SalesOrder{}
			salesOrder.SalesOrderRequestMap(request, now)

			salesOrder.OrderSourceChanMap(getOrderSourceResult)

			salesOrder.SalesOrderStatusChanMap(getOrderStatusResult)

			salesOrder.AgentChanMap(getAgentResult)

			salesOrder.StoreChanMap(getStoreResult)

			salesOrder.UserChanMap(getUserResult)

			if request.SalesmanID > 0 {
				salesOrder.SalesmanChanMap(getSalesmanResult)
			}

			salesOrder.SoCode = soCode
			salesOrder.BrandID = v.BrandID
			salesOrder.TotalAmount = v.Price * float64(v.Qty)
			salesOrder.TotalTonase = float64(v.Qty) * getProductResult.Product.NettWeight

			salesOrderDetail := &models.SalesOrderDetail{}
			salesOrderDetail.SalesOrderDetailStoreRequestMap(v, now)
			salesOrderDetail.SalesOrderDetailStatusChanMap(getOrderDetailStatusResult)
			salesOrderDetail.OrderStatusID = getOrderDetailStatusResult.OrderStatus.ID

			salesOrderDetails := []*models.SalesOrderDetail{}
			salesOrderDetails = append(salesOrderDetails, salesOrderDetail)

			salesOrder.SalesOrderDetails = salesOrderDetails

			// Check Brand
			getBrandResultChan := make(chan *models.BrandChan)
			go u.brandRepository.GetByID(v.BrandID, false, ctx, getBrandResultChan)
			getBrandResult := <-getBrandResultChan

			if getBrandResult.Error != nil {
				return []*models.SalesOrderResponse{}, getBrandResult.ErrorLog
			}

			salesOrder.BrandChanMap(getBrandResult)

			salesOrderBrands[v.BrandID] = salesOrder
		}
	}

	salesOrdersResponse := []*models.SalesOrderResponse{}

	for _, v := range salesOrderBrands {

		salesOrderResponse := &models.SalesOrderResponse{}

		createSalesOrderResultChan := make(chan *models.SalesOrderChan)
		go u.salesOrderRepository.Insert(v, sqlTransaction, ctx, createSalesOrderResultChan)
		createSalesOrderResult := <-createSalesOrderResultChan

		if createSalesOrderResult.Error != nil {
			return []*models.SalesOrderResponse{}, createSalesOrderResult.ErrorLog
		}

		salesOrderResponse.CreateSoResponseMap(v)

		var salesOrderDetailsResponse []*models.SalesOrderDetailStoreResponse
		salesOrderDetails := []*models.SalesOrderDetail{}
		for _, x := range v.SalesOrderDetails {

			soDetailCode, _ := helper.GenerateSODetailCode(int(createSalesOrderResult.ID), v.AgentID, x.ProductID, x.UomID)
			x.SalesOrderID = int(createSalesOrderResult.ID)
			x.SoDetailCode = soDetailCode

			createSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
			go u.salesOrderDetailRepository.Insert(x, sqlTransaction, ctx, createSalesOrderDetailResultChan)
			createSalesOrderDetailResult := <-createSalesOrderDetailResultChan

			if createSalesOrderDetailResult.Error != nil {
				return []*models.SalesOrderResponse{}, createSalesOrderDetailResult.ErrorLog
			}

			salesOrderDetailResponse := &models.SalesOrderDetailStoreResponse{}
			salesOrderDetailResponse.SalesOrderDetailStoreResponseMap(x)

			getProductResultChan := make(chan *models.ProductChan)
			go u.productRepository.GetByID(x.ProductID, false, ctx, getProductResultChan)
			getProductResult := <-getProductResultChan

			if getProductResult.Error != nil {
				errorLogData := helper.WriteLog(getProductResult.Error, http.StatusInternalServerError, nil)
				return []*models.SalesOrderResponse{}, errorLogData
			}

			salesOrderDetailResponse.ProductSKU = getProductResult.Product.Sku.String
			salesOrderDetailResponse.ProductName = getProductResult.Product.ProductName.String
			salesOrderDetailResponse.UnitMeasurementSmall = getProductResult.Product.UnitMeasurementSmall
			salesOrderDetailResponse.UnitMeasurementMedium = getProductResult.Product.UnitMeasurementMedium
			salesOrderDetailResponse.UnitMeasurementBig = getProductResult.Product.UnitMeasurementBig

			getCategoryResultChan := make(chan *models.CategoryChan)
			go u.categoryRepository.GetByID(getProductResult.Product.CategoryID, false, ctx, getCategoryResultChan)
			getCategoryResult := <-getCategoryResultChan

			if getCategoryResult.Error != nil {
				errorLogData := helper.WriteLog(getCategoryResult.Error, http.StatusInternalServerError, nil)
				return []*models.SalesOrderResponse{}, errorLogData
			}

			salesOrderDetailResponse.CategoryId = getProductResult.Product.CategoryID
			salesOrderDetailResponse.CategoryName = getCategoryResult.Category.Name

			getUomResultChan := make(chan *models.UomChan)
			go u.uomRepository.GetByID(x.UomID, false, ctx, getUomResultChan)
			getUomResult := <-getUomResultChan

			if getUomResult.Error != nil {
				errorLogData := helper.WriteLog(getUomResult.Error, http.StatusInternalServerError, nil)
				return []*models.SalesOrderResponse{}, errorLogData
			}

			salesOrderDetailResponse.UomCode = getUomResult.Uom.Code.String

			salesOrderDetailsResponse = append(salesOrderDetailsResponse, salesOrderDetailResponse)
			salesOrderDetails = append(salesOrderDetails, x)

			salesOrderDetailJourneys := &models.SalesOrderDetailJourneys{
				SoDetailId:   createSalesOrderDetailResult.SalesOrderDetail.ID,
				SoDetailCode: soDetailCode,
				Status:       constants.LOG_STATUS_MONGO_DEFAULT,
				Remark:       "",
				Reason:       "",
				CreatedAt:    &now,
				UpdatedAt:    &now,
			}

			createSalesOrderDetailJourneysResultChan := make(chan *models.SalesOrderDetailJourneysChan)
			go u.salesOrderDetailJourneysRepository.Insert(salesOrderDetailJourneys, ctx, createSalesOrderDetailJourneysResultChan)
			createSalesOrderDetailJourneysResult := <-createSalesOrderDetailJourneysResultChan

			if createSalesOrderDetailJourneysResult.Error != nil {
				return []*models.SalesOrderResponse{}, createSalesOrderDetailJourneysResult.ErrorLog
			}
		}

		v.SalesOrderDetails = salesOrderDetails
		salesOrderResponse.SalesOrderDetails = salesOrderDetailsResponse

		salesOrdersResponse = append(salesOrdersResponse, salesOrderResponse)

		salesOrderLog := &models.SalesOrderLog{
			RequestID: request.RequestID,
			SoCode:    salesOrderResponse.SoCode,
			Data:      v,
			Status:    constants.LOG_STATUS_MONGO_DEFAULT,
			Action:    constants.LOG_ACTION_MONGO_INSERT,
			CreatedAt: &now,
			UpdatedAt: &now,
		}

		createSalesOrderLogResultChan := make(chan *models.SalesOrderLogChan)
		go u.salesOrderLogRepository.Insert(salesOrderLog, ctx, createSalesOrderLogResultChan)
		createSalesOrderLogResult := <-createSalesOrderLogResultChan

		if createSalesOrderLogResult.Error != nil {
			return []*models.SalesOrderResponse{}, createSalesOrderLogResult.ErrorLog
		}

		salesOrderJourneys := &models.SalesOrderJourneys{
			SoCode:    salesOrderResponse.SoCode,
			SoId:      createSalesOrderResult.SalesOrder.ID,
			SoDate:    createSalesOrderResult.SalesOrder.SoDate,
			Status:    constants.LOG_STATUS_MONGO_DEFAULT,
			Remark:    "",
			Reason:    "",
			CreatedAt: &now,
			UpdatedAt: &now,
		}

		createSalesOrderJourneysResultChan := make(chan *models.SalesOrderJourneysChan)
		go u.salesOrderJourneysRepository.Insert(salesOrderJourneys, ctx, createSalesOrderJourneysResultChan)
		createSalesOrderJourneysResult := <-createSalesOrderJourneysResultChan

		if createSalesOrderJourneysResult.Error != nil {
			return []*models.SalesOrderResponse{}, createSalesOrderJourneysResult.ErrorLog
		}

		keyKafka := []byte(v.SoCode)
		messageKafka, _ := json.Marshal(v)

		err := u.kafkaClient.WriteToTopic(constants.CREATE_SALES_ORDER_TOPIC, keyKafka, messageKafka)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return []*models.SalesOrderResponse{}, errorLogData
		}

	}

	return salesOrdersResponse, nil
}

func (u *salesOrderUseCase) Get(request *models.SalesOrderRequest) (*models.SalesOrdersOpenSearchResponse, *model.ErrorLog) {
	getSalesOrdersResultChan := make(chan *models.SalesOrdersChan)
	go u.salesOrderOpenSearchRepository.Get(request, getSalesOrdersResultChan)
	getSalesOrdersResult := <-getSalesOrdersResultChan

	if getSalesOrdersResult.Error != nil {
		return &models.SalesOrdersOpenSearchResponse{}, getSalesOrdersResult.ErrorLog
	}

	var salesOrdersResult []*models.SalesOrderOpenSearchResponse
	for _, v := range getSalesOrdersResult.SalesOrders {
		var salesOrder models.SalesOrderOpenSearchResponse

		salesOrder.SalesOrderOpenSearchResponseMap(v)

		var salesOrderDetails []*models.SalesOrderDetailOpenSearchResponse
		for _, x := range v.SalesOrderDetails {
			var salesOrderDetail models.SalesOrderDetailOpenSearchResponse

			salesOrderDetail.SalesOrderDetailOpenSearchResponseMap(x)

			salesOrderDetails = append(salesOrderDetails, &salesOrderDetail)
		}

		salesOrder.SalesOrderDetails = salesOrderDetails

		salesOrdersResult = append(salesOrdersResult, &salesOrder)
	}
	salesOrders := &models.SalesOrdersOpenSearchResponse{
		SalesOrders: salesOrdersResult,
		Total:       getSalesOrdersResult.Total,
	}

	return salesOrders, &model.ErrorLog{}
}

func (u *salesOrderUseCase) GetByID(request *models.SalesOrderRequest, ctx context.Context) ([]*models.SalesOrderOpenSearchResponse, *model.ErrorLog) {

	getSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderOpenSearchRepository.GetByID(request, getSalesOrderResultChan)
	getSalesOrderResult := <-getSalesOrderResultChan

	if getSalesOrderResult.Error != nil {
		return []*models.SalesOrderOpenSearchResponse{}, getSalesOrderResult.ErrorLog
	}

	var salesOrder models.SalesOrderOpenSearchResponse
	salesOrder.SalesOrderOpenSearchResponseMap(getSalesOrderResult.SalesOrder)

	var salesOrderDetails []*models.SalesOrderDetailOpenSearchResponse
	for _, x := range getSalesOrderResult.SalesOrder.SalesOrderDetails {
		var salesOrderDetail models.SalesOrderDetailOpenSearchResponse

		salesOrderDetail.SalesOrderDetailOpenSearchResponseMap(x)

		salesOrderDetails = append(salesOrderDetails, &salesOrderDetail)
	}

	salesOrder.SalesOrderDetails = salesOrderDetails

	var salesOrders []*models.SalesOrderOpenSearchResponse
	salesOrders = append(salesOrders, &salesOrder)

	return salesOrders, &model.ErrorLog{}

}

func (u *salesOrderUseCase) GetByIDWithDetail(request *models.SalesOrderRequest, ctx context.Context) (*models.SalesOrder, *model.ErrorLog) {

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

func (u *salesOrderUseCase) UpdateById(id int, request *models.SalesOrderUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderResponse, *model.ErrorLog) {
	// now := time.Now()
	// var soCode string

	// getSalesOrderByIDResultChan := make(chan *models.SalesOrderChan)
	// go u.salesOrderRepository.GetByID(id, false, ctx, getSalesOrderByIDResultChan)
	// getSalesOrderByIDResult := <-getSalesOrderByIDResultChan

	// if getSalesOrderByIDResult.Error != nil {
	// 	return &models.SalesOrderResponse{}, getSalesOrderByIDResult.ErrorLog
	// }

	// // Check Order Status
	// getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	// go u.orderStatusRepository.GetByID(getSalesOrderByIDResult.SalesOrder.OrderStatusID, false, ctx, getOrderStatusResultChan)
	// getOrderStatusResult := <-getOrderStatusResultChan

	// if getOrderStatusResult.Error != nil {
	// 	return &models.SalesOrderResponse{}, getOrderStatusResult.ErrorLog
	// }
	// getSalesOrderByIDResult.SalesOrder.OrderStatusChanMap(getOrderStatusResult)

	// errorValidation := u.updateSOValidation(getSalesOrderByIDResult.SalesOrder.ID, getOrderStatusResult.OrderStatus.Name, ctx)

	// if errorValidation != nil {
	// 	errorLogData := helper.WriteLog(errorValidation, http.StatusBadRequest, "Ada kesalahan, silahkan coba lagi nanti")
	// 	return &models.SalesOrderResponse{}, errorLogData
	// }

	// // Check Order Detail Status
	// getOrderDetailStatusResultChan := make(chan *models.OrderStatusChan)
	// go u.orderStatusRepository.GetByNameAndType("open", "sales_order_detail", false, ctx, getOrderDetailStatusResultChan)
	// getOrderDetailStatusResult := <-getOrderDetailStatusResultChan

	// if getOrderDetailStatusResult.Error != nil {
	// 	return &models.SalesOrderResponse{}, getOrderDetailStatusResult.ErrorLog
	// }

	// // Check Order Source
	// getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	// go u.orderSourceRepository.GetByID(request.OrderSourceID, false, ctx, getOrderSourceResultChan)
	// getOrderSourceResult := <-getOrderSourceResultChan

	// if getOrderSourceResult.Error != nil {
	// 	return &models.SalesOrderResponse{}, getOrderSourceResult.ErrorLog
	// }
	// getSalesOrderByIDResult.SalesOrder.OrderSourceChanMap(getOrderSourceResult)

	salesOrdersResponse := &models.SalesOrderResponse{}

	// // Check Brand
	// brandIds := []int{}

	// for _, v := range request.SalesOrderDetails {
	// 	brandIds = append(brandIds, v.BrandID)
	// }
	// checkIfBrandSame := helper.InSliceInt(brandIds, request.SalesOrderDetails[0].BrandID)

	// if !checkIfBrandSame {
	// 	errorLogData := helper.WriteLog(fmt.Errorf("The brand id must be the same"), http.StatusBadRequest, "Ada kesalahan, silahkan coba lagi nanti")
	// 	return &models.SalesOrderResponse{}, errorLogData
	// }

	// getBrandResultChan := make(chan *models.BrandChan)
	// go u.brandRepository.GetByID(request.SalesOrderDetails[0].BrandID, false, ctx, getBrandResultChan)
	// getBrandResult := <-getBrandResultChan

	// if getBrandResult.Error != nil {
	// 	return &models.SalesOrderResponse{}, getBrandResult.ErrorLog
	// }

	// getSalesOrderByIDResult.SalesOrder.BrandChanMap(getBrandResult)

	// // Check Agent
	// getAgentResultChan := make(chan *models.AgentChan)
	// go u.agentRepository.GetByID(request.AgentID, false, ctx, getAgentResultChan)
	// getAgentResult := <-getAgentResultChan

	// if getAgentResult.Error != nil {
	// 	errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
	// 	return &models.SalesOrderResponse{}, errorLogData
	// }

	// getSalesOrderByIDResult.SalesOrder.AgentChanMap(getAgentResult)

	// // Check Store
	// getStoreResultChan := make(chan *models.StoreChan)
	// go u.storeRepository.GetByID(request.StoreID, false, ctx, getStoreResultChan)
	// getStoreResult := <-getStoreResultChan

	// if getStoreResult.Error != nil {
	// 	errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
	// 	return &models.SalesOrderResponse{}, errorLogData
	// }

	// getSalesOrderByIDResult.SalesOrder.StoreChanMap(getStoreResult)

	// // Check User Result
	// getUserResultChan := make(chan *models.UserChan)
	// go u.userRepository.GetByID(request.UserID, false, ctx, getUserResultChan)
	// getUserResult := <-getUserResultChan

	// if getUserResult.Error != nil {
	// 	errorLogData := helper.WriteLog(getUserResult.Error, http.StatusInternalServerError, nil)
	// 	return &models.SalesOrderResponse{}, errorLogData
	// }

	// getSalesOrderByIDResult.SalesOrder.UserChanMap(getUserResult)

	// // Check Salesman
	// getSalesmanResultChan := make(chan *models.SalesmanChan)
	// go u.salesmanRepository.GetByEmail(getUserResult.User.Email, false, ctx, getSalesmanResultChan)
	// getSalesmanResult := <-getSalesmanResultChan

	// if getSalesmanResult.Error != nil {
	// 	errorLogData := helper.WriteLog(getSalesmanResult.Error, http.StatusInternalServerError, nil)
	// 	return &models.SalesOrderResponse{}, errorLogData
	// }

	// getSalesOrderByIDResult.SalesOrder.SalesmanChanMap(getSalesmanResult)

	// var salesOrderDetailResponses []*models.SalesOrderDetailStoreResponse
	// var soDetails []*models.SalesOrderDetail
	// var totalAmount float64
	// var totalTonase float64
	// for _, v := range request.SalesOrderDetails {

	// 	getSalesOrderDetailByIDResultChan := make(chan *models.SalesOrderDetailChan)
	// 	go u.salesOrderDetailRepository.GetByID(v.ID, false, ctx, getSalesOrderDetailByIDResultChan)
	// 	getSalesOrderDetailByIDResult := <-getSalesOrderDetailByIDResultChan

	// 	if getSalesOrderDetailByIDResult.Error != nil {
	// 		return &models.SalesOrderResponse{}, getSalesOrderDetailByIDResult.ErrorLog
	// 	}
	// 	salesOrderDetail := &models.SalesOrderDetail{}
	// 	soDetail := &models.SalesOrderDetailStoreRequest{
	// 		SalesOrderDetailTemplate: v.SalesOrderDetailTemplate,
	// 		SalesOrderId:             id,
	// 		Price:                    v.Price,
	// 	}
	// 	salesOrderDetail.SalesOrderDetailStoreRequestMap(soDetail, now)

	// 	updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
	// 	go u.salesOrderDetailRepository.UpdateByID(v.ID, salesOrderDetail, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
	// 	updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

	// 	if updateSalesOrderDetailResult.Error != nil {
	// 		return &models.SalesOrderResponse{}, updateSalesOrderDetailResult.ErrorLog
	// 	}

	// 	soDetail.OrderStatusId = getSalesOrderDetailByIDResult.SalesOrderDetail.OrderStatusID
	// 	soDetail.SoDetailCode = getSalesOrderDetailByIDResult.SalesOrderDetail.SoDetailCode
	// 	salesOrderDetailResponse := &models.SalesOrderDetailStoreResponse{
	// 		ID:                           v.ID,
	// 		SalesOrderDetailStoreRequest: *soDetail,
	// 		CreatedAt:                    getSalesOrderDetailByIDResult.SalesOrderDetail.CreatedAt,
	// 	}

	// 	getProductResultChan := make(chan *models.ProductChan)
	// 	go u.productRepository.GetByID(v.ProductID, false, ctx, getProductResultChan)
	// 	getProductResult := <-getProductResultChan

	// 	if getProductResult.Error != nil {
	// 		errorLogData := helper.WriteLog(getProductResult.Error, http.StatusInternalServerError, nil)
	// 		return &models.SalesOrderResponse{}, errorLogData
	// 	}

	// 	salesOrderDetailResponses = append(salesOrderDetailResponses, salesOrderDetailResponse)
	// 	salesOrderDetail.ID = v.ID
	// 	salesOrderDetail.OrderStatusID = getSalesOrderDetailByIDResult.SalesOrderDetail.OrderStatusID
	// 	soDetails = append(soDetails, salesOrderDetail)

	// 	totalAmount = totalAmount + (v.Price * float64(v.Qty))
	// 	totalTonase = totalTonase + (float64(v.Qty) * getProductResult.Product.NettWeight)

	// }

	// getSalesOrderByIDResult.SalesOrder.SalesOrderDetails = soDetails
	// getSalesOrderByIDResult.SalesOrder.TotalAmount = totalAmount
	// getSalesOrderByIDResult.SalesOrder.TotalTonase = totalTonase

	// salesOrdersResponse.SalesOrderDetails = salesOrderDetailResponses

	// salesOrderUpdateReq := &models.SalesOrder{
	// 	OrderSourceID:   request.OrderSourceID,
	// 	AgentID:         request.AgentID,
	// 	StoreID:         request.StoreID,
	// 	BrandID:         request.SalesOrderDetails[0].BrandID,
	// 	UserID:          request.UserID,
	// 	GLat:            models.NullFloat64{NullFloat64: sql.NullFloat64{Float64: request.GLat, Valid: true}},
	// 	GLong:           models.NullFloat64{NullFloat64: sql.NullFloat64{Float64: request.GLong, Valid: true}},
	// 	SoRefCode:       models.NullString{NullString: sql.NullString{String: request.SoRefCode, Valid: true}},
	// 	SoDate:          request.SoDate,
	// 	SoRefDate:       models.NullString{NullString: sql.NullString{String: request.SoRefDate, Valid: true}},
	// 	Note:            models.NullString{NullString: sql.NullString{String: request.Note, Valid: true}},
	// 	InternalComment: models.NullString{NullString: sql.NullString{String: request.InternalComment, Valid: true}},
	// 	TotalAmount:     totalAmount,
	// 	TotalTonase:     totalTonase,
	// 	DeviceId:        models.NullString{NullString: sql.NullString{String: request.DeviceId, Valid: true}},
	// 	ReferralCode:    models.NullString{NullString: sql.NullString{String: request.ReferralCode, Valid: true}},
	// 	UpdatedAt:       &now,
	// 	LatestUpdatedBy: request.UserID,
	// }

	// updateSalesOrderResultChan := make(chan *models.SalesOrderChan)
	// go u.salesOrderRepository.UpdateByID(id, salesOrderUpdateReq, sqlTransaction, ctx, updateSalesOrderResultChan)
	// updateSalesOrderResult := <-updateSalesOrderResultChan

	// if updateSalesOrderResult.Error != nil {
	// 	return &models.SalesOrderResponse{}, updateSalesOrderResult.ErrorLog
	// }

	// getSalesOrderByIDResult.SalesOrder.UpdateSalesOrderChanMap(updateSalesOrderResult)
	// salesOrdersResponse.SoUpdateByIdResponseMap(getSalesOrderByIDResult.SalesOrder)

	// soCode = getSalesOrderByIDResult.SalesOrder.SoCode

	// salesOrderLog := &models.SalesOrderLog{
	// 	RequestID: request.RequestID,
	// 	SoCode:    soCode,
	// 	Data:      getSalesOrderByIDResult.SalesOrder,
	// 	Status:    "0",
	// 	CreatedAt: &now,
	// }

	// createSalesOrderLogResultChan := make(chan *models.SalesOrderLogChan)
	// go u.salesOrderLogRepository.Insert(salesOrderLog, ctx, createSalesOrderLogResultChan)
	// createSalesOrderLogResult := <-createSalesOrderLogResultChan

	// if createSalesOrderLogResult.Error != nil {
	// 	return &models.SalesOrderResponse{}, createSalesOrderLogResult.ErrorLog
	// }

	// keyKafka := []byte(getSalesOrderByIDResult.SalesOrder.SoCode)
	// messageKafka, _ := json.Marshal(getSalesOrderByIDResult.SalesOrder)
	// err := u.kafkaClient.WriteToTopic(constants.UPDATE_SALES_ORDER_TOPIC, keyKafka, messageKafka)

	// if err != nil {
	// 	errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
	// 	return &models.SalesOrderResponse{}, errorLogData
	// }

	return salesOrdersResponse, nil
}

func (u *salesOrderUseCase) UpdateSODetailById(soId, id int, request *models.SalesOrderDetailUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderDetail, *model.ErrorLog) {
	// now := time.Now()
	// var soCode string

	// getSalesOrderByIDResultChan := make(chan *models.SalesOrderChan)
	// go u.salesOrderRepository.GetByID(soId, false, ctx, getSalesOrderByIDResultChan)
	// getSalesOrderByIDResult := <-getSalesOrderByIDResultChan

	// if getSalesOrderByIDResult.Error != nil {
	// 	return &models.SalesOrderDetail{}, getSalesOrderByIDResult.ErrorLog
	// }

	// // Check Order Status
	// getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	// go u.orderStatusRepository.GetByID(getSalesOrderByIDResult.SalesOrder.OrderStatusID, false, ctx, getOrderStatusResultChan)
	// getOrderStatusResult := <-getOrderStatusResultChan

	// if getOrderStatusResult.Error != nil {
	// 	return &models.SalesOrderDetail{}, getOrderStatusResult.ErrorLog
	// }

	// errorValidation := u.updateSOValidation(getSalesOrderByIDResult.SalesOrder.ID, getOrderStatusResult.OrderStatus.Name, ctx)

	// if errorValidation != nil {
	// 	errorLogData := helper.WriteLog(errorValidation, http.StatusBadRequest, "Ada kesalahan, silahkan coba lagi nanti")
	// 	return &models.SalesOrderDetail{}, errorLogData
	// }

	// // Check Order Detail Status
	// getOrderDetailStatusResultChan := make(chan *models.OrderStatusChan)
	// go u.orderStatusRepository.GetByNameAndType("open", "sales_order_detail", false, ctx, getOrderDetailStatusResultChan)
	// getOrderDetailStatusResult := <-getOrderDetailStatusResultChan

	// if getOrderDetailStatusResult.Error != nil {
	// 	return &models.SalesOrderDetail{}, getOrderDetailStatusResult.ErrorLog
	// }

	// // Check Brand
	// getBrandResultChan := make(chan *models.BrandChan)
	// go u.brandRepository.GetByID(request.BrandID, false, ctx, getBrandResultChan)
	// getBrandResult := <-getBrandResultChan

	// if getBrandResult.Error != nil {
	// 	return &models.SalesOrderDetail{}, getBrandResult.ErrorLog
	// }

	salesOrderDetail := &models.SalesOrderDetail{}
	// salesOrderDetail.SalesOrderDetailUpdateRequestMap(request, now)

	// updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
	// go u.salesOrderDetailRepository.UpdateByID(request.ID, salesOrderDetail, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
	// updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

	// if updateSalesOrderDetailResult.Error != nil {
	// 	return &models.SalesOrderDetail{}, updateSalesOrderDetailResult.ErrorLog
	// }

	// salesOrderDetail.BrandID = request.BrandID

	// salesOrderLog := &models.SalesOrderLog{
	// 	RequestID: "",
	// 	SoCode:    soCode,
	// 	Data:      salesOrderDetail,
	// 	Status:    "0",
	// 	CreatedAt: &now,
	// }

	// createSalesOrderLogResultChan := make(chan *models.SalesOrderLogChan)
	// go u.salesOrderLogRepository.Insert(salesOrderLog, ctx, createSalesOrderLogResultChan)
	// createSalesOrderLogResult := <-createSalesOrderLogResultChan

	// if createSalesOrderLogResult.Error != nil {
	// 	return &models.SalesOrderDetail{}, createSalesOrderLogResult.ErrorLog
	// }

	// keyKafka := []byte(soCode)
	// messageKafka, _ := json.Marshal(salesOrderDetail)
	// err := u.kafkaClient.WriteToTopic(constants.UPDATE_SALES_ORDER_TOPIC, keyKafka, messageKafka)

	// if err != nil {
	// 	errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
	// 	return &models.SalesOrderDetail{}, errorLogData
	// }

	return salesOrderDetail, nil

}

func (u *salesOrderUseCase) UpdateSODetailBySOId(SoId int, request []*models.SalesOrderDetailUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) ([]*models.SalesOrder, *model.ErrorLog) {
	// now := time.Now()
	// var soCode string
	var response []*models.SalesOrder
	// var salesOrder *models.SalesOrder

	// getSalesOrderByIDResultChan := make(chan *models.SalesOrderChan)
	// go u.salesOrderRepository.GetByID(SoId, false, ctx, getSalesOrderByIDResultChan)
	// getSalesOrderByIDResult := <-getSalesOrderByIDResultChan

	// if getSalesOrderByIDResult.Error != nil {
	// 	return []*models.SalesOrder{}, getSalesOrderByIDResult.ErrorLog
	// }

	// salesOrder = getSalesOrderByIDResult.SalesOrder
	// salesOrder.AgentProvinceID = 0
	// salesOrder.AgentCityID = 0
	// salesOrder.AgentDistrictID = 0
	// salesOrder.AgentVillageID = 0
	// salesOrder.StoreProvinceID = 0
	// salesOrder.StoreCityID = 0
	// salesOrder.StoreDistrictID = 0
	// salesOrder.StoreVillageID = 0
	// salesOrder.ReferralCode = models.NullString{NullString: sql.NullString{String: "", Valid: true}}
	// salesOrder.DeviceId = models.NullString{NullString: sql.NullString{String: "", Valid: true}}
	// salesOrder.SalesmanID = 0

	// var soDetails []*models.SalesOrderDetail
	// totalAmount := salesOrder.TotalAmount
	// totalTonase := salesOrder.TotalTonase

	// for _, v := range request {

	// 	salesOrderDetail := &models.SalesOrderDetail{}
	// 	salesOrderDetail.SalesOrderDetailUpdateRequestMap(v, now)

	// 	updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
	// 	go u.salesOrderDetailRepository.UpdateByID(v.ID, salesOrderDetail, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
	// 	updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

	// 	if updateSalesOrderDetailResult.Error != nil {
	// 		return []*models.SalesOrder{}, updateSalesOrderDetailResult.ErrorLog
	// 	}

	// 	getSalesOrderDetailByIDResultChan := make(chan *models.SalesOrderDetailChan)
	// 	go u.salesOrderDetailRepository.GetByID(v.ID, false, ctx, getSalesOrderDetailByIDResultChan)
	// 	getSalesOrderDetailByIDResult := <-getSalesOrderDetailByIDResultChan

	// 	if getSalesOrderDetailByIDResult.Error != nil {
	// 		return []*models.SalesOrder{}, getSalesOrderDetailByIDResult.ErrorLog
	// 	}

	// 	soDetails = append(soDetails, &models.SalesOrderDetail{
	// 		ID:            v.ID,
	// 		SalesOrderID:  SoId,
	// 		ProductID:     v.ProductID,
	// 		UomID:         v.UomID,
	// 		OrderStatusID: getSalesOrderDetailByIDResult.SalesOrderDetail.OrderStatusID,
	// 		SoDetailCode:  getSalesOrderDetailByIDResult.SalesOrderDetail.SoDetailCode,
	// 		Qty:           v.Qty,
	// 		ResidualQty:   v.ResidualQty,
	// 		Price:         v.Price,
	// 		Note:          models.NullString{NullString: sql.NullString{String: v.Note, Valid: true}},
	// 		CreatedAt:     getSalesOrderDetailByIDResult.SalesOrderDetail.CreatedAt,
	// 	})

	// 	getProductByIDResultChan := make(chan *models.ProductChan)
	// 	go u.productRepository.GetByID(v.ID, false, ctx, getProductByIDResultChan)
	// 	getProductByIDResult := <-getProductByIDResultChan

	// 	if getProductByIDResult.Error != nil {
	// 		return []*models.SalesOrder{}, getProductByIDResult.ErrorLog
	// 	}

	// 	totalAmount = totalAmount + (v.Price * float64(v.Qty))
	// 	totalTonase = totalTonase + (float64(v.Qty) * getProductByIDResult.Product.NettWeight)

	// }

	// salesOrder.SalesOrderDetails = soDetails
	// salesOrder.TotalAmount = totalAmount
	// salesOrder.TotalTonase = totalTonase

	// salesOrderUpdateReq := &models.SalesOrder{
	// 	TotalAmount:     totalAmount,
	// 	TotalTonase:     totalTonase,
	// 	UpdatedAt:       &now,
	// 	LatestUpdatedBy: getSalesOrderByIDResult.SalesOrder.UserID,
	// }

	// updateSalesOrderResultChan := make(chan *models.SalesOrderChan)
	// go u.salesOrderRepository.UpdateByID(SoId, salesOrderUpdateReq, sqlTransaction, ctx, updateSalesOrderResultChan)
	// updateSalesOrderResult := <-updateSalesOrderResultChan

	// if updateSalesOrderResult.Error != nil {
	// 	return []*models.SalesOrder{}, updateSalesOrderResult.ErrorLog
	// }

	// // Check Order Status
	// getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	// go u.orderStatusRepository.GetByID(salesOrder.OrderStatusID, false, ctx, getOrderStatusResultChan)
	// getOrderStatusResult := <-getOrderStatusResultChan

	// if getOrderStatusResult.Error != nil {
	// 	return []*models.SalesOrder{}, getOrderStatusResult.ErrorLog
	// }

	// salesOrder.OrderStatusName = getOrderStatusResult.OrderStatus.Name

	// // Check Order Source
	// getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	// go u.orderSourceRepository.GetByID(salesOrder.OrderSourceID, false, ctx, getOrderSourceResultChan)
	// getOrderSourceResult := <-getOrderSourceResultChan

	// if getOrderSourceResult.Error != nil {
	// 	return []*models.SalesOrder{}, getOrderSourceResult.ErrorLog
	// }

	// salesOrder.OrderSourceName = getOrderSourceResult.OrderSource.SourceName

	// getBrandResultChan := make(chan *models.BrandChan)
	// go u.brandRepository.GetByID(salesOrder.BrandID, false, ctx, getBrandResultChan)
	// getBrandResult := <-getBrandResultChan

	// if getBrandResult.Error != nil {
	// 	return []*models.SalesOrder{}, getBrandResult.ErrorLog
	// }

	// salesOrder.BrandName = getBrandResult.Brand.Name

	// // Check Agent
	// getAgentResultChan := make(chan *models.AgentChan)
	// go u.agentRepository.GetByID(salesOrder.AgentID, false, ctx, getAgentResultChan)
	// getAgentResult := <-getAgentResultChan

	// if getAgentResult.Error != nil {
	// 	errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
	// 	return []*models.SalesOrder{}, errorLogData
	// }

	// salesOrder.AgentName = models.NullString{NullString: sql.NullString{String: getAgentResult.Agent.Name, Valid: true}}
	// salesOrder.AgentEmail = getAgentResult.Agent.Email
	// salesOrder.AgentProvinceName = getAgentResult.Agent.ProvinceName
	// salesOrder.AgentCityName = getAgentResult.Agent.CityName
	// salesOrder.AgentDistrictName = getAgentResult.Agent.DistrictName
	// salesOrder.AgentVillageName = getAgentResult.Agent.VillageName
	// salesOrder.AgentAddress = getAgentResult.Agent.Address
	// salesOrder.AgentPhone = getAgentResult.Agent.Phone
	// salesOrder.AgentMainMobilePhone = getAgentResult.Agent.MainMobilePhone

	// // Check Store
	// getStoreResultChan := make(chan *models.StoreChan)
	// go u.storeRepository.GetByID(salesOrder.StoreID, false, ctx, getStoreResultChan)
	// getStoreResult := <-getStoreResultChan

	// if getStoreResult.Error != nil {
	// 	errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
	// 	return []*models.SalesOrder{}, errorLogData
	// }

	// salesOrder.StoreName = getStoreResult.Store.Name
	// salesOrder.StoreCode = getStoreResult.Store.StoreCode
	// salesOrder.StoreEmail = getStoreResult.Store.Email
	// salesOrder.StoreProvinceName = getStoreResult.Store.ProvinceName
	// salesOrder.StoreCityName = getStoreResult.Store.CityName
	// salesOrder.StoreDistrictName = getStoreResult.Store.DistrictName
	// salesOrder.StoreVillageName = getStoreResult.Store.VillageName
	// salesOrder.StoreAddress = getStoreResult.Store.Address
	// salesOrder.StorePhone = getStoreResult.Store.Phone
	// salesOrder.StoreMainMobilePhone = getStoreResult.Store.MainMobilePhone

	// // Check User Result
	// getUserResultChan := make(chan *models.UserChan)
	// go u.userRepository.GetByID(salesOrder.UserID, false, ctx, getUserResultChan)
	// getUserResult := <-getUserResultChan

	// if getUserResult.Error != nil {
	// 	errorLogData := helper.WriteLog(getUserResult.Error, http.StatusInternalServerError, nil)
	// 	return []*models.SalesOrder{}, errorLogData
	// }

	// salesOrder.UserFirstName = getUserResult.User.FirstName
	// salesOrder.UserLastName = getUserResult.User.LastName
	// salesOrder.UserEmail = models.NullString{NullString: sql.NullString{String: getUserResult.User.Email, Valid: true}}

	// // Check Salesman
	// getSalesmanResultChan := make(chan *models.SalesmanChan)
	// go u.salesmanRepository.GetByEmail(getUserResult.User.Email, false, ctx, getSalesmanResultChan)
	// getSalesmanResult := <-getSalesmanResultChan

	// if getSalesmanResult.Error != nil {
	// 	errorLogData := helper.WriteLog(getSalesmanResult.Error, http.StatusInternalServerError, nil)
	// 	return []*models.SalesOrder{}, errorLogData
	// }

	// salesOrder.SalesmanName = models.NullString{NullString: sql.NullString{String: getSalesmanResult.Salesman.Name, Valid: true}}
	// salesOrder.SalesmanEmail = getSalesmanResult.Salesman.Email

	// response = append(response, salesOrder)

	// salesOrderLog := &models.SalesOrderLog{
	// 	RequestID: "",
	// 	SoCode:    soCode,
	// 	Data:      salesOrder,
	// 	Status:    "0",
	// 	CreatedAt: &now,
	// }

	// createSalesOrderLogResultChan := make(chan *models.SalesOrderLogChan)
	// go u.salesOrderLogRepository.Insert(salesOrderLog, ctx, createSalesOrderLogResultChan)
	// createSalesOrderLogResult := <-createSalesOrderLogResultChan

	// if createSalesOrderLogResult.Error != nil {
	// 	return []*models.SalesOrder{}, createSalesOrderLogResult.ErrorLog
	// }

	// keyKafka := []byte(soCode)
	// messageKafka, _ := json.Marshal(salesOrder)
	// err := u.kafkaClient.WriteToTopic(constants.UPDATE_SALES_ORDER_TOPIC, keyKafka, messageKafka)

	// if err != nil {
	// 	errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
	// 	return []*models.SalesOrder{}, errorLogData
	// }

	return response, nil
}

func (u *salesOrderUseCase) updateSOValidation(salesOrderId int, orderStatusName string, ctx context.Context) error {

	if orderStatusName == "closed" {
		return fmt.Errorf("Cannot update. Sales order status are close")
	}

	getDeliveryOrderByIDResultChan := make(chan *models.DeliveryOrdersChan)
	go u.deliveryOrderRepository.GetBySalesOrderID(salesOrderId, true, ctx, getDeliveryOrderByIDResultChan)
	getDeliveryOrderByIDResult := <-getDeliveryOrderByIDResultChan

	if getDeliveryOrderByIDResult.Total == 0 {

		return nil

	} else {

		for _, v := range getDeliveryOrderByIDResult.DeliveryOrders {

			getOrderStatusResultChan := make(chan *models.OrderStatusChan)
			go u.orderStatusRepository.GetByID(v.OrderStatusID, false, ctx, getOrderStatusResultChan)
			getOrderStatusResult := <-getOrderStatusResultChan

			if getOrderStatusResult.OrderStatus.Name != "cancel" {
				return fmt.Errorf("Cannot update. Order delivery must be cancel first")
			}
		}

	}

	return nil
}

func (u *salesOrderUseCase) GetDetails(request *models.SalesOrderRequest) (*models.SalesOrderDetailsOpenSearchResponse, *model.ErrorLog) {
	getSalesOrdersResultChan := make(chan *models.SalesOrdersChan)
	go u.salesOrderOpenSearchRepository.Get(request, getSalesOrdersResultChan)
	getSalesOrdersResult := <-getSalesOrdersResultChan

	if getSalesOrdersResult.Error != nil {
		return &models.SalesOrderDetailsOpenSearchResponse{}, getSalesOrdersResult.ErrorLog
	}

	var salesOrderDetails []*models.SalesOrderDetailOpenSearchResponse
	for _, v := range getSalesOrdersResult.SalesOrders {
		for _, x := range v.SalesOrderDetails {
			var salesOrderDetail models.SalesOrderDetailOpenSearchResponse

			salesOrderDetail.SalesOrderDetailOpenSearchResponseMap(x)

			salesOrderDetails = append(salesOrderDetails, &salesOrderDetail)
		}
	}

	salesOrders := &models.SalesOrderDetailsOpenSearchResponse{
		SalesOrderDetails: salesOrderDetails,
		Total:             int64(len(salesOrderDetails)),
	}

	return salesOrders, &model.ErrorLog{}
}

func (u *salesOrderUseCase) DeleteById(id int, sqlTransaction *sql.Tx) *model.ErrorLog {
	now := time.Now()

	getSalesOrderByIDResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.GetByID(id, false, u.ctx, getSalesOrderByIDResultChan)
	getSalesOrderByIDResult := <-getSalesOrderByIDResultChan

	if getSalesOrderByIDResult.Error != nil {
		return getSalesOrderByIDResult.ErrorLog
	}

	getSalesOrderDetailsByIDResultChan := make(chan *models.SalesOrderDetailsChan)
	go u.salesOrderDetailRepository.GetBySalesOrderID(getSalesOrderByIDResult.SalesOrder.ID, false, u.ctx, getSalesOrderDetailsByIDResultChan)
	getSalesOrderDetailsByIDResult := <-getSalesOrderDetailsByIDResultChan

	if getSalesOrderDetailsByIDResult.Error != nil {
		return getSalesOrderDetailsByIDResult.ErrorLog
	}

	getSalesOrderByIDResult.SalesOrder.SalesOrderDetails = getSalesOrderDetailsByIDResult.SalesOrderDetails

	var soDetails []*models.SalesOrderDetail
	for _, v := range getSalesOrderByIDResult.SalesOrder.SalesOrderDetails {
		deleteSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.DeleteByID(v, sqlTransaction, u.ctx, deleteSalesOrderDetailResultChan)
		updateSalesOrderDetailResult := <-deleteSalesOrderDetailResultChan

		if updateSalesOrderDetailResult.Error != nil {
			return updateSalesOrderDetailResult.ErrorLog
		}
		soDetails = append(soDetails, updateSalesOrderDetailResult.SalesOrderDetail)
	}

	deleteSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.DeleteByID(getSalesOrderByIDResult.SalesOrder, sqlTransaction, u.ctx, deleteSalesOrderResultChan)
	deleteSalesOrderResult := <-deleteSalesOrderResultChan
	if deleteSalesOrderResult.Error != nil {
		return deleteSalesOrderResult.ErrorLog
	}

	getSalesOrderByIDResult.SalesOrder.SalesOrderDetails = soDetails
	getSalesOrderByIDResult.SalesOrder.UpdateSalesOrderChanMap(deleteSalesOrderResult)

	salesOrderLog := &models.SalesOrderLog{
		RequestID: "",
		SoCode:    getSalesOrderByIDResult.SalesOrder.SoCode,
		Data:      getSalesOrderByIDResult.SalesOrder,
		Action:    constants.LOG_ACTION_MONGO_DELETE,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		CreatedAt: &now,
	}
	createSalesOrderLogResultChan := make(chan *models.SalesOrderLogChan)
	go u.salesOrderLogRepository.Insert(salesOrderLog, u.ctx, createSalesOrderLogResultChan)
	createSalesOrderLogResult := <-createSalesOrderLogResultChan

	if createSalesOrderLogResult.Error != nil {
		return createSalesOrderLogResult.ErrorLog
	}
	keyKafka := []byte(getSalesOrderByIDResult.SalesOrder.SoCode)
	messageKafka, _ := json.Marshal(
		&models.SalesOrder{
			ID:        id,
			SoCode:    salesOrderLog.SoCode,
			UpdatedAt: getSalesOrderByIDResult.SalesOrder.UpdatedAt,
			DeletedAt: getSalesOrderByIDResult.SalesOrder.DeletedAt,
		},
	)
	err := u.kafkaClient.WriteToTopic(constants.DELETE_SALES_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}

func (u *salesOrderUseCase) GetDetailById(id int) (*models.SalesOrderDetailOpenSearchResponse, *model.ErrorLog) {
	result := &models.SalesOrderDetailOpenSearchResponse{}
	getSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderOpenSearchRepository.GetDetailByID(id, getSalesOrderResultChan)
	getSalesOrderResult := <-getSalesOrderResultChan

	if getSalesOrderResult.Error != nil {
		return &models.SalesOrderDetailOpenSearchResponse{}, getSalesOrderResult.ErrorLog
	}

	var salesOrder models.SalesOrderOpenSearchResponse
	salesOrder.SalesOrderOpenSearchResponseMap(getSalesOrderResult.SalesOrder)

	for _, x := range getSalesOrderResult.SalesOrder.SalesOrderDetails {
		fmt.Println("id = ", x.ID, x.ID == id)
		if x.ID == id {
			result.SalesOrderDetailOpenSearchResponseMap(x)
		}
		fmt.Println("result = ", result.ID)
	}
	fmt.Println("result = ", result.ID)
	return result, &model.ErrorLog{}

}
