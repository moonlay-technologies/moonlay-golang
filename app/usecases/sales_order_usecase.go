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
	"strings"
	"time"

	"github.com/bxcodec/dbresolver"
	"github.com/google/uuid"
)

type SalesOrderUseCaseInterface interface {
	Create(request *models.SalesOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) ([]*models.SalesOrderResponse, *model.ErrorLog)
	Get(request *models.SalesOrderRequest) (*models.SalesOrdersOpenSearchResponse, *model.ErrorLog)
	GetByID(request *models.SalesOrderRequest, ctx context.Context) ([]*models.SalesOrderOpenSearchResponse, *model.ErrorLog)
	GetByIDWithDetail(request *models.SalesOrderRequest, ctx context.Context) (*models.SalesOrder, *model.ErrorLog)
	GetByAgentID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetByStoreID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetBySalesmanID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog)
	GetSyncToKafkaHistories(request *models.SalesOrderEventLogRequest, ctx context.Context) ([]*models.SalesOrderEventLogResponse, *model.ErrorLog)
	GetSOJourneys(request *models.SalesOrderJourneyRequest, ctx context.Context) (*models.SalesOrderJourneyResponses, *model.ErrorLog)
	GetSOJourneyBySOId(soId int, ctx context.Context) (*models.SalesOrderJourneyResponses, *model.ErrorLog)
	GetSOUploadHistories(request *models.GetSoUploadHistoriesRequest, ctx context.Context) (*models.GetSoUploadHistoryResponses, *model.ErrorLog)
	GetSOUploadHistoriesByid(id string, ctx context.Context) (*models.GetSoUploadHistoryResponse, *model.ErrorLog)
	GetSOUploadErrorLogsByReqId(request *models.GetSoUploadErrorLogsRequest, ctx context.Context) (*models.GetSoUploadErrorLogsResponse, *model.ErrorLog)
	GetSOUploadErrorLogsBySoUploadHistoryId(request *models.GetSoUploadErrorLogsRequest, ctx context.Context) (*models.GetSoUploadErrorLogsResponse, *model.ErrorLog)
	UpdateById(id int, request *models.SalesOrderUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderResponse, *model.ErrorLog)
	UpdateSODetailById(soId, soDetailId int, request *models.SalesOrderDetailUpdateByIdRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderDetailStoreResponse, *model.ErrorLog)
	UpdateSODetailBySOId(soId int, request *models.SalesOrderUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderResponse, *model.ErrorLog)
	GetDetails(request *models.GetSalesOrderDetailRequest) (*models.SalesOrderDetailsOpenSearchResponse, *model.ErrorLog)
	GetDetailById(id int) (*models.SalesOrderDetailOpenSearchResponse, *model.ErrorLog)
	DeleteById(id int, sqlTransaction *sql.Tx) *model.ErrorLog
	DeleteDetailBySOId(id int, sqlTransaction *sql.Tx) *model.ErrorLog
	DeleteDetailById(id int, sqlTransaction *sql.Tx) *model.ErrorLog
	RetrySyncToKafka(logId string) (*models.SORetryProcessSyncToKafkaResponse, *model.ErrorLog)
	Export(request *models.SalesOrderExportRequest, ctx context.Context) (string, *model.ErrorLog)
	ExportDetail(request *models.SalesOrderDetailExportRequest, ctx context.Context) (string, *model.ErrorLog)
}

type salesOrderUseCase struct {
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
	salesOrderLogRepository              mongoRepositories.SalesOrderLogRepositoryInterface
	salesOrderJourneysRepository         mongoRepositories.SalesOrderJourneysRepositoryInterface
	salesOrderDetailJourneysRepository   mongoRepositories.SalesOrderDetailJourneysRepositoryInterface
	soUploadHistoriesRepository          mongoRepositories.SoUploadHistoriesRepositoryInterface
	soUploadErrorLogsRepository          mongoRepositories.SoUploadErrorLogsRepositoryInterface
	userRepository                       repositories.UserRepositoryInterface
	salesmanRepository                   repositories.SalesmanRepositoryInterface
	categoryRepository                   repositories.CategoryRepositoryInterface
	salesOrderOpenSearchRepository       openSearchRepositories.SalesOrderOpenSearchRepositoryInterface
	salesOrderDetailOpenSearchRepository openSearchRepositories.SalesOrderDetailOpenSearchRepositoryInterface
	kafkaClient                          kafkadbo.KafkaClientInterface
	db                                   dbresolver.DB
	ctx                                  context.Context
}

func InitSalesOrderUseCaseInterface(salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, orderSourceRepository repositories.OrderSourceRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, uomRepository repositories.UomRepositoryInterface, deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, salesOrderLogRepository mongoRepositories.SalesOrderLogRepositoryInterface, salesOrderJourneysRepository mongoRepositories.SalesOrderJourneysRepositoryInterface, salesOrderDetailJourneysRepository mongoRepositories.SalesOrderDetailJourneysRepositoryInterface, soUploadHistoriesRepository mongoRepositories.SoUploadHistoriesRepositoryInterface, soUploadErrorLogsRepository mongoRepositories.SoUploadErrorLogsRepositoryInterface, userRepository repositories.UserRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, categoryRepository repositories.CategoryRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, salesOrderDetailOpenSearchRepository openSearchRepositories.SalesOrderDetailOpenSearchRepositoryInterface, kafkaClient kafkadbo.KafkaClientInterface, db dbresolver.DB, ctx context.Context) SalesOrderUseCaseInterface {
	return &salesOrderUseCase{
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
		salesOrderLogRepository:              salesOrderLogRepository,
		salesOrderJourneysRepository:         salesOrderJourneysRepository,
		salesOrderDetailJourneysRepository:   salesOrderDetailJourneysRepository,
		soUploadHistoriesRepository:          soUploadHistoriesRepository,
		soUploadErrorLogsRepository:          soUploadErrorLogsRepository,
		userRepository:                       userRepository,
		salesmanRepository:                   salesmanRepository,
		categoryRepository:                   categoryRepository,
		salesOrderOpenSearchRepository:       salesOrderOpenSearchRepository,
		salesOrderDetailOpenSearchRepository: salesOrderDetailOpenSearchRepository,
		kafkaClient:                          kafkaClient,
		db:                                   db,
		ctx:                                  ctx,
	}
}

func (u *salesOrderUseCase) Create(request *models.SalesOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) ([]*models.SalesOrderResponse, *baseModel.ErrorLog) {
	now := time.Now()
	var soCode string

	// Get Order Source Status By Id
	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(request.OrderSourceID, false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getOrderSourceResult.Error != nil {
		return []*models.SalesOrderResponse{}, getOrderSourceResult.ErrorLog
	}

	var status string
	if getOrderSourceResult.OrderSource.SourceName == "manager" {
		status = "open"
	} else {
		status = "pending"
	}

	// Get Order Status By Name
	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType(status, "sales_order", false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return []*models.SalesOrderResponse{}, getOrderStatusResult.ErrorLog
	}

	// Get Order Detail Status By Name
	getOrderDetailStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType(status, "sales_order_detail", false, ctx, getOrderDetailStatusResultChan)
	getOrderDetailStatusResult := <-getOrderDetailStatusResultChan

	if getOrderDetailStatusResult.Error != nil {
		return []*models.SalesOrderResponse{}, getOrderDetailStatusResult.ErrorLog
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

	if len(request.SoRefCode) < 1 {
		request.SoRefCode = helper.GenerateSORefCode(request.AgentID, request.SoDate)
	}

	brandIds := []int{}
	salesOrderBrands := map[int]*models.SalesOrder{}

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
			salesOrderDetail.SalesOrderDetailStatusChanMap(getOrderDetailStatusResult)
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
			salesOrderDetail.Note = models.NullString{NullString: sql.NullString{String: request.Note, Valid: true}}
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
			salesOrderDetailResponse.UnitMeasurementSmall = getProductResult.Product.UnitMeasurementSmall.String
			salesOrderDetailResponse.UnitMeasurementMedium = getProductResult.Product.UnitMeasurementMedium.String
			salesOrderDetailResponse.UnitMeasurementBig = getProductResult.Product.UnitMeasurementBig.String

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

		}

		v.SalesOrderDetails = salesOrderDetails
		salesOrderResponse.SalesOrderDetails = salesOrderDetailsResponse

		salesOrdersResponse = append(salesOrdersResponse, salesOrderResponse)

		salesOrderLog := &models.SalesOrderLog{
			RequestID: ctx.Value("RequestId").(string),
			SoCode:    v.SoCode,
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

func (u *salesOrderUseCase) Get(request *models.SalesOrderRequest) (*models.SalesOrdersOpenSearchResponse, *baseModel.ErrorLog) {
	getSalesOrdersResultChan := make(chan *models.SalesOrdersChan)
	go u.salesOrderOpenSearchRepository.Get(request, false, getSalesOrdersResultChan)
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

	return salesOrders, &baseModel.ErrorLog{}
}

func (u *salesOrderUseCase) GetByID(request *models.SalesOrderRequest, ctx context.Context) ([]*models.SalesOrderOpenSearchResponse, *baseModel.ErrorLog) {

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

	return salesOrders, &baseModel.ErrorLog{}

}

func (u *salesOrderUseCase) GetByIDWithDetail(request *models.SalesOrderRequest, ctx context.Context) (*models.SalesOrder, *baseModel.ErrorLog) {

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

	return getSalesOrderResult.SalesOrder, &baseModel.ErrorLog{}

}

func (u *salesOrderUseCase) GetByAgentID(request *models.SalesOrderRequest) (*models.SalesOrders, *baseModel.ErrorLog) {
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

	return salesOrders, &baseModel.ErrorLog{}
}

func (u *salesOrderUseCase) GetByStoreID(request *models.SalesOrderRequest) (*models.SalesOrders, *baseModel.ErrorLog) {
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

	return salesOrders, &baseModel.ErrorLog{}
}

func (u *salesOrderUseCase) GetBySalesmanID(request *models.SalesOrderRequest) (*models.SalesOrders, *baseModel.ErrorLog) {
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

	return salesOrders, &baseModel.ErrorLog{}
}

func (u *salesOrderUseCase) GetSyncToKafkaHistories(request *models.SalesOrderEventLogRequest, ctx context.Context) ([]*models.SalesOrderEventLogResponse, *baseModel.ErrorLog) {
	getSalesOrderLogResultChan := make(chan *models.GetSalesOrderLogsChan)
	go u.salesOrderLogRepository.Get(request, false, ctx, getSalesOrderLogResultChan)
	getSalesOrderLogResult := <-getSalesOrderLogResultChan

	if getSalesOrderLogResult.Error != nil {
		return []*models.SalesOrderEventLogResponse{}, getSalesOrderLogResult.ErrorLog
	}

	salesOrderEventLogs := []*models.SalesOrderEventLogResponse{}
	for _, v := range getSalesOrderLogResult.SalesOrderLogs {
		var status string
		switch v.Status {
		case constants.LOG_STATUS_MONGO_DEFAULT:
			status = "In Progress"
		case constants.LOG_STATUS_MONGO_SUCCESS:
			status = "Success"
		case constants.LOG_STATUS_MONGO_ERROR:
			status = "Failed"
		default:
			status = ""
		}

		salesOrderEventLog := models.SalesOrderEventLogResponse{}
		salesOrderEventLog.SalesOrderEventLogResponseMap(v, status)
		dataSOEventLog := models.DataSOEventLogResponse{}
		dataSOEventLog.DataSOEventLogResponseMap(v)
		salesOrderEventLog.Data = &dataSOEventLog

		for _, x := range v.Data.SalesOrderDetails {
			soDetailEventLog := models.SODetailEventLogResponse{}
			soDetailEventLog.SoDetailEventLogResponse(x)

			dataSOEventLog.SalesOrderDetails = append(dataSOEventLog.SalesOrderDetails, &soDetailEventLog)
		}

		salesOrderEventLogs = append(salesOrderEventLogs, &salesOrderEventLog)
	}

	return salesOrderEventLogs, nil
}

func (u *salesOrderUseCase) GetSOJourneys(request *models.SalesOrderJourneyRequest, ctx context.Context) (*models.SalesOrderJourneyResponses, *model.ErrorLog) {
	getSalesOrderJourneyResultChan := make(chan *models.SalesOrdersJourneysChan)
	go u.salesOrderJourneysRepository.Get(request, false, ctx, getSalesOrderJourneyResultChan)
	getSalesOrderJourneyResult := <-getSalesOrderJourneyResultChan

	if getSalesOrderJourneyResult.Error != nil {
		return &models.SalesOrderJourneyResponses{}, getSalesOrderJourneyResult.ErrorLog
	}

	salesOrderJourneys := []*models.SalesOrderJourneyResponse{}
	for _, v := range getSalesOrderJourneyResult.SalesOrderJourneys {
		salesOrderJourney := models.SalesOrderJourneyResponse{}
		salesOrderJourney.SalesOrderJourneyResponseMap(v)

		salesOrderJourneys = append(salesOrderJourneys, &salesOrderJourney)
	}

	salesOrderJourneysResult := models.SalesOrderJourneyResponses{
		SalesOrderJourneys: salesOrderJourneys,
		Total:              getSalesOrderJourneyResult.Total,
	}

	return &salesOrderJourneysResult, nil
}

func (u *salesOrderUseCase) GetSOJourneyBySOId(soId int, ctx context.Context) (*models.SalesOrderJourneyResponses, *baseModel.ErrorLog) {
	getSalesOrderJourneyResultChan := make(chan *models.SalesOrdersJourneysChan)
	go u.salesOrderJourneysRepository.GetBySoId(soId, false, ctx, getSalesOrderJourneyResultChan)
	getSalesOrderJourneyResult := <-getSalesOrderJourneyResultChan

	if getSalesOrderJourneyResult.Error != nil {
		return &models.SalesOrderJourneyResponses{}, getSalesOrderJourneyResult.ErrorLog
	}

	salesOrderJourneys := []*models.SalesOrderJourneyResponse{}
	for _, v := range getSalesOrderJourneyResult.SalesOrderJourneys {
		orderStatusID := 0
		switch v.Status {
		case constants.SO_STATUS_APPV:
			orderStatusID = 5
		case constants.SO_STATUS_REAPPV:
			orderStatusID = 5
		case constants.SO_STATUS_RJC:
			orderStatusID = 9
		case constants.SO_STATUS_CNCL:
			orderStatusID = 10
		case constants.SO_STATUS_ORDPRT:
			orderStatusID = 7
		case constants.SO_STATUS_ORDCLS:
			orderStatusID = 8
		case constants.SO_STATUS_PEND:
			orderStatusID = 6
		default:
			orderStatusID = 0
		}

		salesOrderJourney := models.SalesOrderJourneyResponse{}
		salesOrderJourney.SalesOrderJourneyResponseMap(v)
		salesOrderJourney.OrderStatusID = &orderStatusID

		salesOrderJourneys = append(salesOrderJourneys, &salesOrderJourney)
	}

	salesOrderJourneysResult := &models.SalesOrderJourneyResponses{
		SalesOrderJourneys: salesOrderJourneys,
		Total:              getSalesOrderJourneyResult.Total,
	}

	return salesOrderJourneysResult, nil
}

func (u *salesOrderUseCase) GetSOUploadHistories(request *models.GetSoUploadHistoriesRequest, ctx context.Context) (*models.GetSoUploadHistoryResponses, *model.ErrorLog) {

	getSoUploadHistoriesResultChan := make(chan *models.SoUploadHistoriesChan)
	go u.soUploadHistoriesRepository.Get(request, false, ctx, getSoUploadHistoriesResultChan)
	getSoUploadHistoriesResult := <-getSoUploadHistoriesResultChan

	if getSoUploadHistoriesResult.Error != nil {
		return &models.GetSoUploadHistoryResponses{}, getSoUploadHistoriesResult.ErrorLog
	}

	result := models.GetSoUploadHistoryResponses{
		SoUploadHistories: getSoUploadHistoriesResult.SoUploadHistories,
		Total:             getSoUploadHistoriesResult.Total,
	}

	return &result, nil
}

func (u *salesOrderUseCase) GetSOUploadHistoriesByid(id string, ctx context.Context) (*models.GetSoUploadHistoryResponse, *model.ErrorLog) {

	getSoUploadHistoryByIdResultChan := make(chan *models.GetSoUploadHistoryResponseChan)
	go u.soUploadHistoriesRepository.GetByHistoryID(id, false, ctx, getSoUploadHistoryByIdResultChan)
	getSoUploadHistoryByIdResult := <-getSoUploadHistoryByIdResultChan

	if getSoUploadHistoryByIdResult.Error != nil {
		return &models.GetSoUploadHistoryResponse{}, getSoUploadHistoryByIdResult.ErrorLog
	}

	return getSoUploadHistoryByIdResult.SoUploadHistories, nil
}

func (u *salesOrderUseCase) GetSOUploadErrorLogsByReqId(request *models.GetSoUploadErrorLogsRequest, ctx context.Context) (*models.GetSoUploadErrorLogsResponse, *model.ErrorLog) {

	getSoUploadErrorLogsResultChan := make(chan *models.SoUploadErrorLogsChan)
	go u.soUploadErrorLogsRepository.Get(request, false, ctx, getSoUploadErrorLogsResultChan)
	getSoUploadErrorLogsResult := <-getSoUploadErrorLogsResultChan

	if getSoUploadErrorLogsResult.Error != nil {
		return &models.GetSoUploadErrorLogsResponse{}, getSoUploadErrorLogsResult.ErrorLog
	}

	result := models.GetSoUploadErrorLogsResponse{
		SoUploadErrosLogs: getSoUploadErrorLogsResult.SoUploadErrorLogs,
		Total:             getSoUploadErrorLogsResult.Total,
	}

	return &result, nil

}

func (u *salesOrderUseCase) GetSOUploadErrorLogsBySoUploadHistoryId(request *models.GetSoUploadErrorLogsRequest, ctx context.Context) (*models.GetSoUploadErrorLogsResponse, *model.ErrorLog) {

	getSoUploadErrorLogsResultChan := make(chan *models.SoUploadErrorLogsChan)
	go u.soUploadErrorLogsRepository.Get(request, false, ctx, getSoUploadErrorLogsResultChan)
	getSoUploadErrorLogsResult := <-getSoUploadErrorLogsResultChan

	if getSoUploadErrorLogsResult.Error != nil {
		return &models.GetSoUploadErrorLogsResponse{}, getSoUploadErrorLogsResult.ErrorLog
	}

	result := models.GetSoUploadErrorLogsResponse{
		SoUploadErrosLogs: getSoUploadErrorLogsResult.SoUploadErrorLogs,
		Total:             getSoUploadErrorLogsResult.Total,
	}

	return &result, nil

}

func (u *salesOrderUseCase) UpdateById(id int, request *models.SalesOrderUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderResponse, *baseModel.ErrorLog) {
	now := time.Now()
	salesOrder := &models.SalesOrder{}

	// Get Sales Order By Id
	getSalesOrderByIDResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.GetByID(id, false, ctx, getSalesOrderByIDResultChan)
	getSalesOrderByIDResult := <-getSalesOrderByIDResultChan

	if getSalesOrderByIDResult.Error != nil {
		return &models.SalesOrderResponse{}, getSalesOrderByIDResult.ErrorLog
	}
	salesOrder = getSalesOrderByIDResult.SalesOrder

	// Get Order Status By Id
	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByID(request.OrderStatusID, false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return &models.SalesOrderResponse{}, getOrderStatusResult.ErrorLog
	}

	salesOrder.SalesOrderStatusChanMap(getOrderStatusResult)

	// Get Order Source Status By Id
	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(salesOrder.OrderSourceID, false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getOrderSourceResult.Error != nil {
		return &models.SalesOrderResponse{}, getOrderSourceResult.ErrorLog
	}

	salesOrder.OrderSourceChanMap(getOrderSourceResult)

	// Check Agent By Id
	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(salesOrder.AgentID, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	salesOrder.AgentChanMap(getAgentResult)

	// Check Store By Id
	getStoreResultChan := make(chan *models.StoreChan)
	go u.storeRepository.GetByID(salesOrder.StoreID, false, ctx, getStoreResultChan)
	getStoreResult := <-getStoreResultChan

	if getStoreResult.Error != nil {
		errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	salesOrder.StoreChanMap(getStoreResult)

	// Check User By Id
	getUserResultChan := make(chan *models.UserChan)
	go u.userRepository.GetByID(salesOrder.UserID, false, ctx, getUserResultChan)
	getUserResult := <-getUserResultChan

	if getUserResult.Error != nil {
		errorLogData := helper.WriteLog(getUserResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	salesOrder.UserChanMap(getUserResult)

	// Check Salesman By Id
	getSalesmanResult := &models.SalesmanChan{}

	if salesOrder.SalesmanID.Int64 > 0 {
		getSalesmanResultChan := make(chan *models.SalesmanChan)
		go u.salesmanRepository.GetByID(int(salesOrder.SalesmanID.Int64), false, ctx, getSalesmanResultChan)
		getSalesmanResult = <-getSalesmanResultChan

		if getSalesmanResult.Error != nil {
			errorLogData := helper.WriteLog(getSalesmanResult.Error, http.StatusInternalServerError, nil)
			return &models.SalesOrderResponse{}, errorLogData
		}

		salesOrder.SalesmanChanMap(getSalesmanResult)
	}

	salesOrderUpdateReq := &models.SalesOrder{
		OrderStatusID:   getOrderStatusResult.OrderStatus.ID,
		SoCode:          getSalesOrderByIDResult.SalesOrder.SoCode,
		SoDate:          strings.ReplaceAll(getSalesOrderByIDResult.SalesOrder.SoDate, "T00:00:00Z", ""),
		LatestUpdatedBy: ctx.Value("user").(*models.UserClaims).UserID,
		UpdatedAt:       &now,
	}

	// Update Sales Order
	updateSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.UpdateByID(id, salesOrderUpdateReq, true, request.Reason, sqlTransaction, ctx, updateSalesOrderResultChan)
	updateSalesOrderResult := <-updateSalesOrderResultChan

	if updateSalesOrderResult.Error != nil {
		return &models.SalesOrderResponse{}, updateSalesOrderResult.ErrorLog
	}

	// Remove Cache Sales Order
	removeCacheSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.RemoveCacheByID(id, ctx, removeCacheSalesOrderResultChan)
	removeCacheSalesOrderResult := <-removeCacheSalesOrderResultChan

	if removeCacheSalesOrderResult.Error != nil {
		return &models.SalesOrderResponse{}, removeCacheSalesOrderResult.ErrorLog
	}

	salesOrdersResponse := &models.SalesOrderResponse{}
	salesOrdersResponse.UpdateSoResponseMap(salesOrder)

	var salesOrderDetailResponses []*models.SalesOrderDetailStoreResponse
	var salesOrderDetails []*models.SalesOrderDetail
	for _, v := range request.SalesOrderDetails {

		// Get Sales Order Detail By Id
		getSalesOrderDetailByIDResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.GetByID(v.ID, false, ctx, getSalesOrderDetailByIDResultChan)
		getSalesOrderDetailByIDResult := <-getSalesOrderDetailByIDResultChan

		if getSalesOrderDetailByIDResult.Error != nil {
			return &models.SalesOrderResponse{}, getSalesOrderDetailByIDResult.ErrorLog
		}

		// Get Order Detail Status By Id
		getOrderDetailStatusResultChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByID(v.OrderStatusID, false, ctx, getOrderDetailStatusResultChan)
		getOrderDetailStatusResult := <-getOrderDetailStatusResultChan

		if getOrderDetailStatusResult.Error != nil {
			return &models.SalesOrderResponse{}, getOrderDetailStatusResult.ErrorLog
		}

		salesOrderDetail := &models.SalesOrderDetail{
			OrderStatusID:   getOrderDetailStatusResult.OrderStatus.ID,
			SoDetailCode:    getSalesOrderDetailByIDResult.SalesOrderDetail.SoDetailCode,
			LatestUpdatedBy: ctx.Value("user").(*models.UserClaims).UserID,
			UpdatedAt:       &now,
		}

		// Update Sales Order Detail
		updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.UpdateByID(v.ID, salesOrderDetail, true, v.Reason, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
		updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

		if updateSalesOrderDetailResult.Error != nil {
			return &models.SalesOrderResponse{}, updateSalesOrderDetailResult.ErrorLog
		}

		// Remove Cache Sales Order Detail
		clearCacheSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.RemoveCacheByID(v.ID, ctx, clearCacheSalesOrderDetailResultChan)
		clearCacheSalesOrderDetailResult := <-clearCacheSalesOrderDetailResultChan

		if clearCacheSalesOrderDetailResult.Error != nil {
			return &models.SalesOrderResponse{}, clearCacheSalesOrderDetailResult.ErrorLog
		}

		// Get Sales Order Detail by Id
		getSalesOrderDetailByIDResultChan = make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.GetByID(v.ID, false, ctx, getSalesOrderDetailByIDResultChan)
		getSalesOrderDetailByIDResult = <-getSalesOrderDetailByIDResultChan

		if getSalesOrderDetailByIDResult.Error != nil {
			return &models.SalesOrderResponse{}, getSalesOrderDetailByIDResult.ErrorLog
		}

		getSalesOrderDetailByIDResult.SalesOrderDetail.OrderStatusID = getOrderDetailStatusResult.OrderStatus.ID
		salesOrderDetails = append(salesOrderDetails, getSalesOrderDetailByIDResult.SalesOrderDetail)

		salesOrderDetailResponse := &models.SalesOrderDetailStoreResponse{}
		salesOrderDetailResponse.UpdateSoDetailResponseMap(getSalesOrderDetailByIDResult.SalesOrderDetail)
		salesOrderDetailResponse.OrderStatusId = getOrderDetailStatusResult.OrderStatus.ID
		salesOrderDetailResponses = append(salesOrderDetailResponses, salesOrderDetailResponse)
	}

	salesOrder.SalesOrderDetails = salesOrderDetails
	salesOrdersResponse.SalesOrderDetails = salesOrderDetailResponses

	salesOrderLog := &models.SalesOrderLog{
		RequestID: ctx.Value("RequestId").(string),
		SoCode:    salesOrder.SoCode,
		Data:      salesOrder,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		Action:    constants.LOG_ACTION_MONGO_UPDATE,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	createSalesOrderLogResultChan := make(chan *models.SalesOrderLogChan)
	go u.salesOrderLogRepository.Insert(salesOrderLog, ctx, createSalesOrderLogResultChan)
	createSalesOrderLogResult := <-createSalesOrderLogResultChan

	if createSalesOrderLogResult.Error != nil {
		return &models.SalesOrderResponse{}, createSalesOrderLogResult.ErrorLog
	}

	keyKafka := []byte(salesOrder.SoCode)
	messageKafka, _ := json.Marshal(salesOrder)

	err := u.kafkaClient.WriteToTopic(constants.UPDATE_SALES_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	return salesOrdersResponse, nil
}

func (u *salesOrderUseCase) UpdateSODetailById(soId, soDetailId int, request *models.SalesOrderDetailUpdateByIdRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderDetailStoreResponse, *baseModel.ErrorLog) {
	now := time.Now()
	salesOrder := &models.SalesOrder{}

	// Get Sales Order Detail by Id
	getSalesOrderDetailByIDResultChan := make(chan *models.SalesOrderDetailChan)
	go u.salesOrderDetailRepository.GetByID(soDetailId, false, ctx, getSalesOrderDetailByIDResultChan)
	getSalesOrderDetailByIDResult := <-getSalesOrderDetailByIDResultChan

	if getSalesOrderDetailByIDResult.Error != nil {
		return &models.SalesOrderDetailStoreResponse{}, getSalesOrderDetailByIDResult.ErrorLog
	}

	// Get Sales Order By Id
	getSalesOrderByIDResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.GetByID(soId, false, ctx, getSalesOrderByIDResultChan)
	getSalesOrderByIDResult := <-getSalesOrderByIDResultChan

	if getSalesOrderByIDResult.Error != nil {
		return &models.SalesOrderDetailStoreResponse{}, getSalesOrderByIDResult.ErrorLog
	}

	salesOrder = getSalesOrderByIDResult.SalesOrder

	// Get Sales Order Detail by So Id
	getSalesOrderDetailsBySoIdResultChan := make(chan *models.SalesOrderDetailsChan)
	go u.salesOrderDetailRepository.GetBySalesOrderID(soId, false, ctx, getSalesOrderDetailsBySoIdResultChan)
	getSalesOrderDetailsBySoIdResult := <-getSalesOrderDetailsBySoIdResultChan

	if getSalesOrderDetailsBySoIdResult.Error != nil {
		return &models.SalesOrderDetailStoreResponse{}, getSalesOrderDetailsBySoIdResult.ErrorLog
	}

	salesOrder.SalesOrderDetails = getSalesOrderDetailsBySoIdResult.SalesOrderDetails
	totalSoDetail := getSalesOrderDetailsBySoIdResult.Total
	var soStatus string
	var soDetailStatus string

	if salesOrder.OrderStatusName == constants.ORDER_STATUS_OPEN {
		if totalSoDetail == 1 {
			soStatus = constants.ORDER_STATUS_CANCELLED
			soDetailStatus = constants.ORDER_STATUS_CANCELLED
		} else if helper.CheckSalesOrderDetailStatus(soDetailId, true, constants.ORDER_STATUS_CANCELLED, salesOrder.SalesOrderDetails) > 0 {
			soStatus = constants.ORDER_STATUS_OPEN
			soDetailStatus = constants.ORDER_STATUS_CANCELLED
		} else if helper.CheckSalesOrderDetailStatus(soDetailId, false, constants.ORDER_STATUS_CANCELLED, salesOrder.SalesOrderDetails) > 0 {
			soStatus = constants.ORDER_STATUS_CANCELLED
			soDetailStatus = constants.ORDER_STATUS_CANCELLED
		}
	} else if salesOrder.OrderStatusName == constants.ORDER_STATUS_PARTIAL {
		if totalSoDetail == 1 {
			soStatus = constants.ORDER_STATUS_CLOSED
			soDetailStatus = constants.ORDER_STATUS_CLOSED
		} else if totalSoDetail > 1 && getSalesOrderDetailByIDResult.SalesOrderDetail.SentQty > 0 {
			soStatus = constants.ORDER_STATUS_PARTIAL
			soDetailStatus = constants.ORDER_STATUS_CLOSED
		} else if totalSoDetail > 1 && getSalesOrderDetailByIDResult.SalesOrderDetail.SentQty == 0 {
			soStatus = constants.ORDER_STATUS_PARTIAL
			soDetailStatus = constants.ORDER_STATUS_CANCELLED
		}
	}

	// Check Brand
	getBrandResultChan := make(chan *models.BrandChan)
	go u.brandRepository.GetByID(salesOrder.BrandID, false, ctx, getBrandResultChan)
	getBrandResult := <-getBrandResultChan

	if getBrandResult.Error != nil {
		return &models.SalesOrderDetailStoreResponse{}, getBrandResult.ErrorLog
	}

	salesOrder.BrandChanMap(getBrandResult)

	// Get Order Status By Name
	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType(soStatus, "sales_order", false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return &models.SalesOrderDetailStoreResponse{}, getOrderStatusResult.ErrorLog
	}

	// Get Order Detail Status By Name
	getOrderDetailStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByNameAndType(soDetailStatus, "sales_order_detail", false, ctx, getOrderDetailStatusResultChan)
	getOrderDetailStatusResult := <-getOrderDetailStatusResultChan

	if getOrderDetailStatusResult.Error != nil {
		return &models.SalesOrderDetailStoreResponse{}, getOrderDetailStatusResult.ErrorLog
	}

	// Get Order Source Status By Id
	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(salesOrder.OrderSourceID, false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getOrderSourceResult.Error != nil {
		return &models.SalesOrderDetailStoreResponse{}, getOrderSourceResult.ErrorLog
	}

	// Check Agent By Id
	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(salesOrder.AgentID, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderDetailStoreResponse{}, errorLogData
	}

	// Check Store By Id
	getStoreResultChan := make(chan *models.StoreChan)
	go u.storeRepository.GetByID(salesOrder.StoreID, false, ctx, getStoreResultChan)
	getStoreResult := <-getStoreResultChan

	if getStoreResult.Error != nil {
		errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderDetailStoreResponse{}, errorLogData
	}

	// Check User By Id
	getUserResultChan := make(chan *models.UserChan)
	go u.userRepository.GetByID(salesOrder.UserID, false, ctx, getUserResultChan)
	getUserResult := <-getUserResultChan

	if getUserResult.Error != nil {
		errorLogData := helper.WriteLog(getUserResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderDetailStoreResponse{}, errorLogData
	}

	// Check Salesman By Id
	getSalesmanResult := &models.SalesmanChan{}

	if salesOrder.SalesmanID.Int64 > 0 {
		getSalesmanResultChan := make(chan *models.SalesmanChan)
		go u.salesmanRepository.GetByID(int(salesOrder.SalesmanID.Int64), false, ctx, getSalesmanResultChan)
		getSalesmanResult = <-getSalesmanResultChan

		if getSalesmanResult.Error != nil {
			errorLogData := helper.WriteLog(getSalesmanResult.Error, http.StatusInternalServerError, nil)
			return &models.SalesOrderDetailStoreResponse{}, errorLogData
		}

		salesOrder.SalesmanChanMap(getSalesmanResult)
	}

	if soStatus != salesOrder.OrderStatusName {

		salesOrderUpdateReq := &models.SalesOrder{
			OrderStatusID:   getOrderStatusResult.OrderStatus.ID,
			SoCode:          getSalesOrderByIDResult.SalesOrder.SoCode,
			SoDate:          strings.ReplaceAll(getSalesOrderByIDResult.SalesOrder.SoDate, "T00:00:00Z", ""),
			LatestUpdatedBy: ctx.Value("user").(*models.UserClaims).UserID,
			UpdatedAt:       &now,
		}

		// Update Sales Order
		updateSalesOrderResultChan := make(chan *models.SalesOrderChan)
		go u.salesOrderRepository.UpdateByID(soId, salesOrderUpdateReq, true, request.Reason, sqlTransaction, ctx, updateSalesOrderResultChan)
		updateSalesOrderResult := <-updateSalesOrderResultChan

		if updateSalesOrderResult.Error != nil {
			return &models.SalesOrderDetailStoreResponse{}, updateSalesOrderResult.ErrorLog
		}

		// Remove Cache Sales Order
		removeCacheSalesOrderResultChan := make(chan *models.SalesOrderChan)
		go u.salesOrderRepository.RemoveCacheByID(soId, ctx, removeCacheSalesOrderResultChan)
		removeCacheSalesOrderResult := <-removeCacheSalesOrderResultChan

		if removeCacheSalesOrderResult.Error != nil {
			return &models.SalesOrderDetailStoreResponse{}, removeCacheSalesOrderResult.ErrorLog
		}

	}

	salesOrder.OrderSourceChanMap(getOrderSourceResult)
	salesOrder.SalesOrderStatusChanMap(getOrderStatusResult)
	salesOrder.AgentChanMap(getAgentResult)
	salesOrder.StoreChanMap(getStoreResult)
	salesOrder.UserChanMap(getUserResult)

	// Get Sales Order Detail by Id
	getSalesOrderDetailByIdResultChan := make(chan *models.SalesOrderDetailChan)
	go u.salesOrderDetailRepository.GetByID(soDetailId, false, ctx, getSalesOrderDetailByIdResultChan)
	getSalesOrderDetailByIdResult := <-getSalesOrderDetailByIdResultChan

	if getSalesOrderDetailByIdResult.Error != nil {
		return &models.SalesOrderDetailStoreResponse{}, getSalesOrderDetailsBySoIdResult.ErrorLog
	}

	salesOrderDetailReq := &models.SalesOrderDetail{
		OrderStatusID:   getOrderDetailStatusResult.OrderStatus.ID,
		SoDetailCode:    getSalesOrderDetailByIdResult.SalesOrderDetail.SoDetailCode,
		LatestUpdatedBy: ctx.Value("user").(*models.UserClaims).UserID,
		UpdatedAt:       &now,
	}

	// Update Sales Order Detail
	updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
	go u.salesOrderDetailRepository.UpdateByID(soDetailId, salesOrderDetailReq, true, request.Reason, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
	updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

	if updateSalesOrderDetailResult.Error != nil {
		return &models.SalesOrderDetailStoreResponse{}, updateSalesOrderDetailResult.ErrorLog
	}

	for _, v := range salesOrder.SalesOrderDetails {
		if v.ID == soDetailId {
			v.OrderStatusID = getOrderDetailStatusResult.OrderStatus.ID
			v.LatestUpdatedBy = ctx.Value("user").(*models.UserClaims).UserID
			v.UpdatedAt = &now
			break
		}
	}

	salesOrderDetail := &models.SalesOrderDetailStoreResponse{}
	salesOrderDetail.UpdateSalesOrderDetailByIdResponseMap(getSalesOrderDetailByIDResult.SalesOrderDetail, salesOrder.BrandID)

	// Remove Cache Sales Order Detail
	clearCacheSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
	go u.salesOrderDetailRepository.RemoveCacheByID(soDetailId, ctx, clearCacheSalesOrderDetailResultChan)
	clearCacheSalesOrderDetailResult := <-clearCacheSalesOrderDetailResultChan

	if clearCacheSalesOrderDetailResult.Error != nil {
		return &models.SalesOrderDetailStoreResponse{}, clearCacheSalesOrderDetailResult.ErrorLog

	}

	salesOrderLog := &models.SalesOrderLog{
		RequestID: ctx.Value("RequestId").(string),
		SoCode:    salesOrder.SoCode,
		Data:      salesOrder,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		Action:    constants.LOG_ACTION_MONGO_UPDATE,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	createSalesOrderLogResultChan := make(chan *models.SalesOrderLogChan)
	go u.salesOrderLogRepository.Insert(salesOrderLog, ctx, createSalesOrderLogResultChan)
	createSalesOrderLogResult := <-createSalesOrderLogResultChan

	if createSalesOrderLogResult.Error != nil {
		return &models.SalesOrderDetailStoreResponse{}, createSalesOrderLogResult.ErrorLog
	}

	keyKafka := []byte(getSalesOrderDetailByIDResult.SalesOrderDetail.SoDetailCode)
	messageKafka, _ := json.Marshal(salesOrder)

	err := u.kafkaClient.WriteToTopic(constants.UPDATE_SALES_ORDER_DETAIL_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.SalesOrderDetailStoreResponse{}, errorLogData
	}

	return salesOrderDetail, nil
}

func (u *salesOrderUseCase) UpdateSODetailBySOId(soId int, request *models.SalesOrderUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderResponse, *baseModel.ErrorLog) {
	now := time.Now()
	salesOrder := &models.SalesOrder{}

	// Get Sales Order By Id
	getSalesOrderByIDResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.GetByID(soId, false, ctx, getSalesOrderByIDResultChan)
	getSalesOrderByIDResult := <-getSalesOrderByIDResultChan

	if getSalesOrderByIDResult.Error != nil {
		return &models.SalesOrderResponse{}, getSalesOrderByIDResult.ErrorLog
	}
	salesOrder = getSalesOrderByIDResult.SalesOrder

	// Get Order Status By Id
	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByID(request.OrderStatusID, false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		return &models.SalesOrderResponse{}, getOrderStatusResult.ErrorLog
	}

	salesOrder.SalesOrderStatusChanMap(getOrderStatusResult)

	// Get Order Source Status By Id
	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(salesOrder.OrderSourceID, false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getOrderSourceResult.Error != nil {
		return &models.SalesOrderResponse{}, getOrderSourceResult.ErrorLog
	}

	salesOrder.OrderSourceChanMap(getOrderSourceResult)

	// Check Agent By Id
	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(salesOrder.AgentID, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	salesOrder.AgentChanMap(getAgentResult)

	// Check Store By Id
	getStoreResultChan := make(chan *models.StoreChan)
	go u.storeRepository.GetByID(salesOrder.StoreID, false, ctx, getStoreResultChan)
	getStoreResult := <-getStoreResultChan

	if getStoreResult.Error != nil {
		errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	salesOrder.StoreChanMap(getStoreResult)

	// Check User By Id
	getUserResultChan := make(chan *models.UserChan)
	go u.userRepository.GetByID(salesOrder.UserID, false, ctx, getUserResultChan)
	getUserResult := <-getUserResultChan

	if getUserResult.Error != nil {
		errorLogData := helper.WriteLog(getUserResult.Error, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	salesOrder.UserChanMap(getUserResult)

	// Check Salesman By Id
	getSalesmanResult := &models.SalesmanChan{}

	if salesOrder.SalesmanID.Int64 > 0 {
		getSalesmanResultChan := make(chan *models.SalesmanChan)
		go u.salesmanRepository.GetByID(int(salesOrder.SalesmanID.Int64), false, ctx, getSalesmanResultChan)
		getSalesmanResult = <-getSalesmanResultChan

		if getSalesmanResult.Error != nil {
			errorLogData := helper.WriteLog(getSalesmanResult.Error, http.StatusInternalServerError, nil)
			return &models.SalesOrderResponse{}, errorLogData
		}

		salesOrder.SalesmanChanMap(getSalesmanResult)
	}

	salesOrderUpdateReq := &models.SalesOrder{
		OrderStatusID:   getOrderStatusResult.OrderStatus.ID,
		SoCode:          getSalesOrderByIDResult.SalesOrder.SoCode,
		SoDate:          strings.ReplaceAll(getSalesOrderByIDResult.SalesOrder.SoDate, "T00:00:00Z", ""),
		LatestUpdatedBy: ctx.Value("user").(*models.UserClaims).UserID,
		UpdatedAt:       &now,
	}

	// Update Sales Order
	updateSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.UpdateByID(soId, salesOrderUpdateReq, true, request.Reason, sqlTransaction, ctx, updateSalesOrderResultChan)
	updateSalesOrderResult := <-updateSalesOrderResultChan

	if updateSalesOrderResult.Error != nil {
		return &models.SalesOrderResponse{}, updateSalesOrderResult.ErrorLog
	}

	// Remove Cache Sales Order
	removeCacheSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.RemoveCacheByID(soId, ctx, removeCacheSalesOrderResultChan)
	removeCacheSalesOrderResult := <-removeCacheSalesOrderResultChan

	if removeCacheSalesOrderResult.Error != nil {
		return &models.SalesOrderResponse{}, removeCacheSalesOrderResult.ErrorLog
	}

	salesOrdersResponse := &models.SalesOrderResponse{}
	salesOrdersResponse.UpdateSoResponseMap(salesOrder)

	var salesOrderDetailResponses []*models.SalesOrderDetailStoreResponse
	var salesOrderDetails []*models.SalesOrderDetail
	for _, v := range request.SalesOrderDetails {

		// Get Sales Order Detail By Id
		getSalesOrderDetailByIDResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.GetByID(v.ID, false, ctx, getSalesOrderDetailByIDResultChan)
		getSalesOrderDetailByIDResult := <-getSalesOrderDetailByIDResultChan

		if getSalesOrderDetailByIDResult.Error != nil {
			return &models.SalesOrderResponse{}, getSalesOrderDetailByIDResult.ErrorLog
		}

		// Get Order Detail Status By Name
		getOrderDetailStatusResultChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByID(v.OrderStatusID, false, ctx, getOrderDetailStatusResultChan)
		getOrderDetailStatusResult := <-getOrderDetailStatusResultChan

		if getOrderDetailStatusResult.Error != nil {
			return &models.SalesOrderResponse{}, getOrderDetailStatusResult.ErrorLog
		}

		salesOrderDetail := &models.SalesOrderDetail{
			OrderStatusID:   getOrderDetailStatusResult.OrderStatus.ID,
			SoDetailCode:    getSalesOrderDetailByIDResult.SalesOrderDetail.SoDetailCode,
			LatestUpdatedBy: ctx.Value("user").(*models.UserClaims).UserID,
			UpdatedAt:       &now,
		}

		// Update Sales Order Detail
		updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.UpdateByID(v.ID, salesOrderDetail, true, v.Reason, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
		updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

		if updateSalesOrderDetailResult.Error != nil {
			return &models.SalesOrderResponse{}, updateSalesOrderDetailResult.ErrorLog
		}

		// Remove Cache Sales Order Detail
		clearCacheSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.RemoveCacheByID(v.ID, ctx, clearCacheSalesOrderDetailResultChan)
		clearCacheSalesOrderDetailResult := <-clearCacheSalesOrderDetailResultChan

		if clearCacheSalesOrderDetailResult.Error != nil {
			return &models.SalesOrderResponse{}, clearCacheSalesOrderDetailResult.ErrorLog
		}

		// Get Sales Order Detail by Id
		getSalesOrderDetailByIDResultChan = make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.GetByID(v.ID, false, ctx, getSalesOrderDetailByIDResultChan)
		getSalesOrderDetailByIDResult = <-getSalesOrderDetailByIDResultChan

		if getSalesOrderDetailByIDResult.Error != nil {
			return &models.SalesOrderResponse{}, getSalesOrderDetailByIDResult.ErrorLog
		}

		getSalesOrderDetailByIDResult.SalesOrderDetail.OrderStatusID = getOrderDetailStatusResult.OrderStatus.ID
		salesOrderDetails = append(salesOrderDetails, getSalesOrderDetailByIDResult.SalesOrderDetail)

		salesOrderDetailResponse := &models.SalesOrderDetailStoreResponse{}
		salesOrderDetailResponse.UpdateSoDetailResponseMap(getSalesOrderDetailByIDResult.SalesOrderDetail)
		salesOrderDetailResponse.OrderStatusId = getOrderDetailStatusResult.OrderStatus.ID
		salesOrderDetailResponses = append(salesOrderDetailResponses, salesOrderDetailResponse)
	}

	salesOrder.SalesOrderDetails = salesOrderDetails
	salesOrdersResponse.SalesOrderDetails = salesOrderDetailResponses

	salesOrderLog := &models.SalesOrderLog{
		RequestID: ctx.Value("RequestId").(string),
		SoCode:    salesOrder.SoCode,
		Data:      salesOrder,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		Action:    constants.LOG_ACTION_MONGO_UPDATE,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	createSalesOrderLogResultChan := make(chan *models.SalesOrderLogChan)
	go u.salesOrderLogRepository.Insert(salesOrderLog, ctx, createSalesOrderLogResultChan)
	createSalesOrderLogResult := <-createSalesOrderLogResultChan

	if createSalesOrderLogResult.Error != nil {
		return &models.SalesOrderResponse{}, createSalesOrderLogResult.ErrorLog
	}

	keyKafka := []byte(salesOrder.SoCode)
	messageKafka, _ := json.Marshal(salesOrder)

	err := u.kafkaClient.WriteToTopic(constants.UPDATE_SALES_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.SalesOrderResponse{}, errorLogData
	}

	return salesOrdersResponse, nil
}

func (u *salesOrderUseCase) GetDetails(request *models.GetSalesOrderDetailRequest) (*models.SalesOrderDetailsOpenSearchResponse, *baseModel.ErrorLog) {
	getSalesOrderDetailsResultChan := make(chan *models.SalesOrderDetailsOpenSearchChan)
	go u.salesOrderDetailOpenSearchRepository.Get(request, false, getSalesOrderDetailsResultChan)
	getSalesOrderDetailsResult := <-getSalesOrderDetailsResultChan

	if getSalesOrderDetailsResult.Error != nil {
		return &models.SalesOrderDetailsOpenSearchResponse{}, getSalesOrderDetailsResult.ErrorLog
	}

	var salesOrderDetails []*models.SalesOrderDetailOpenSearchResponse
	for _, v := range getSalesOrderDetailsResult.SalesOrderDetails {

		var salesOrderDetail models.SalesOrderDetailOpenSearchResponse

		salesOrderDetail.SalesOrderDetailOpenSearchMap(v)

		salesOrderDetails = append(salesOrderDetails, &salesOrderDetail)
	}

	salesOrders := &models.SalesOrderDetailsOpenSearchResponse{
		SalesOrderDetails: salesOrderDetails,
		Total:             getSalesOrderDetailsResult.Total,
	}

	return salesOrders, &baseModel.ErrorLog{}
}

func (u *salesOrderUseCase) DeleteById(id int, sqlTransaction *sql.Tx) *baseModel.ErrorLog {
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

func (u *salesOrderUseCase) DeleteDetailById(id int, sqlTransaction *sql.Tx) *model.ErrorLog {
	now := time.Now()

	getSalesOrderDetailByIDResultChan := make(chan *models.SalesOrderDetailChan)
	go u.salesOrderDetailRepository.GetByID(id, false, u.ctx, getSalesOrderDetailByIDResultChan)
	getSalesOrderDetailByIDResult := <-getSalesOrderDetailByIDResultChan

	if getSalesOrderDetailByIDResult.Error != nil {
		return getSalesOrderDetailByIDResult.ErrorLog
	}

	getSalesOrderByIDResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.GetByID(getSalesOrderDetailByIDResult.SalesOrderDetail.SalesOrderID, false, u.ctx, getSalesOrderByIDResultChan)
	getSalesOrderByIDResult := <-getSalesOrderByIDResultChan

	if getSalesOrderByIDResult.Error != nil {
		return getSalesOrderByIDResult.ErrorLog
	}

	deleteSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
	go u.salesOrderDetailRepository.DeleteByID(getSalesOrderDetailByIDResult.SalesOrderDetail, sqlTransaction, u.ctx, deleteSalesOrderDetailResultChan)
	deleteSalesOrderDetailResult := <-deleteSalesOrderDetailResultChan

	if deleteSalesOrderDetailResult.Error != nil {
		return deleteSalesOrderDetailResult.ErrorLog
	}

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
		&models.SalesOrderDetail{
			ID:           id,
			SalesOrderID: getSalesOrderByIDResult.SalesOrder.ID,
			UpdatedAt:    getSalesOrderByIDResult.SalesOrder.UpdatedAt,
			DeletedAt:    getSalesOrderByIDResult.SalesOrder.DeletedAt,
		},
	)
	err := u.kafkaClient.WriteToTopic(constants.DELETE_SALES_ORDER_DETAIL_TOPIC, keyKafka, messageKafka)

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
		if x.ID == id {
			result.SalesOrderDetailOpenSearchResponseMap(x)
		}
	}
	return result, &baseModel.ErrorLog{}

}

func (u *salesOrderUseCase) RetrySyncToKafka(logId string) (*models.SORetryProcessSyncToKafkaResponse, *model.ErrorLog) {

	now := time.Now()

	getSalesOrderLogByIdResultChan := make(chan *models.GetSalesOrderLogChan)
	go u.salesOrderLogRepository.GetByID(logId, false, u.ctx, getSalesOrderLogByIdResultChan)
	getSalesOrderLogByIdResult := <-getSalesOrderLogByIdResultChan

	if getSalesOrderLogByIdResult.Error != nil {
		return &models.SORetryProcessSyncToKafkaResponse{}, getSalesOrderLogByIdResult.ErrorLog
	}

	if getSalesOrderLogByIdResult.SalesOrderLog.Status != "2" {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage("retry", fmt.Sprintf("status log dengan id %s bukan gagal", logId))},
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})

		return &models.SORetryProcessSyncToKafkaResponse{}, errorLog
	}

	keyKafka := []byte(getSalesOrderLogByIdResult.SalesOrderLog.SoCode)
	messageKafka, _ := json.Marshal(getSalesOrderLogByIdResult.SalesOrderLog.Data)

	var topic string
	switch getSalesOrderLogByIdResult.SalesOrderLog.Action {
	case constants.LOG_ACTION_MONGO_INSERT:
		topic = constants.CREATE_SALES_ORDER_TOPIC
	case constants.LOG_ACTION_MONGO_UPDATE:
		topic = constants.UPDATE_SALES_ORDER_TOPIC
	case constants.LOG_ACTION_MONGO_DELETE:
		topic = constants.DELETE_SALES_ORDER_TOPIC
	}

	err := u.kafkaClient.WriteToTopic(topic, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.SORetryProcessSyncToKafkaResponse{}, errorLogData
	}

	salesOrderLog := &models.SalesOrderLog{
		RequestID: getSalesOrderLogByIdResult.SalesOrderLog.RequestID,
		SoCode:    getSalesOrderLogByIdResult.SalesOrderLog.SoCode,
		Data:      getSalesOrderLogByIdResult.SalesOrderLog.Data,
		Action:    getSalesOrderLogByIdResult.SalesOrderLog.Action,
		Status:    constants.LOG_STATUS_MONGO_DEFAULT,
		CreatedAt: getSalesOrderLogByIdResult.SalesOrderLog.CreatedAt,
		UpdatedAt: &now,
	}

	salesOrderLogResultChan := make(chan *models.SalesOrderLogChan)
	go u.salesOrderLogRepository.UpdateByID(logId, salesOrderLog, u.ctx, salesOrderLogResultChan)

	result := models.SORetryProcessSyncToKafkaResponse{
		SalesOrderLogEventId: logId,
		Status:               constants.LOG_STATUS_MONGO_DEFAULT,
		Message:              "on progres",
	}
	return &result, nil

}

func (u *salesOrderUseCase) DeleteDetailBySOId(id int, sqlTransaction *sql.Tx) *model.ErrorLog {
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

	// deleteSalesOrderResultChan := make(chan *models.SalesOrderChan)
	// go u.salesOrderRepository.DeleteByID(getSalesOrderByIDResult.SalesOrder, sqlTransaction, u.ctx, deleteSalesOrderResultChan)
	// deleteSalesOrderResult := <-deleteSalesOrderResultChan
	// if deleteSalesOrderResult.Error != nil {
	// 	return deleteSalesOrderResult.ErrorLog
	// }

	getSalesOrderByIDResult.SalesOrder.SalesOrderDetails = soDetails

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
			ID:                id,
			SoCode:            salesOrderLog.SoCode,
			SalesOrderDetails: soDetails,
			UpdatedAt:         getSalesOrderByIDResult.SalesOrder.UpdatedAt,
			DeletedAt:         getSalesOrderByIDResult.SalesOrder.DeletedAt,
		},
	)
	err := u.kafkaClient.WriteToTopic(constants.DELETE_SALES_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}

func (u *salesOrderUseCase) Export(request *models.SalesOrderExportRequest, ctx context.Context) (string, *model.ErrorLog) {
	soRequest := &models.SalesOrderRequest{}
	soRequest.SalesOrderExportMap(request)
	if soRequest.SortField == "order_status" {
		request.SortField = "order_status.name"
		soRequest.SortField = "order_status.name"
	}
	getSalesOrdersCountResultChan := make(chan *models.SalesOrdersChan)
	go u.salesOrderOpenSearchRepository.Get(soRequest, true, getSalesOrdersCountResultChan)
	getSalesOrdersCountResult := <-getSalesOrdersCountResultChan

	if getSalesOrdersCountResult.Error != nil {
		fmt.Println("error = ", getSalesOrdersCountResult.Error)
		errorLogData := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{"Data tidak ditemukan"},
			SystemMessage: []string{"sales Orders data not found"},
			StatusCode:    http.StatusNotFound,
		})
		return "", errorLogData
	}

	if getSalesOrdersCountResult.Total == 0 {
		err := helper.NewError("Data tidak ditemukan")
		errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
		errorLogData = helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{"Data tidak ditemukan"},
			SystemMessage: []string{"sales Orders data not found"},
			StatusCode:    http.StatusNotFound,
		})
		return "", errorLogData
	}

	rand, err := helper.Generate(`[A-Za-z]{12}`)
	loc, err := time.LoadLocation("Asia/Jakarta")
	fileHour := time.Now().In(loc).Format(constants.DATE_FORMAT_EXPORT)
	if ctx == nil {
		err = fmt.Errorf("nil context")
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return "", errorLogData
	}
	fileName := fmt.Sprintf("SO-LIST-SUMMARY-%s-%d-%s", fileHour, ctx.Value("user").(*models.UserClaims).UserID, rand)
	request.FileName = fileName
	request.UserID = ctx.Value("user").(*models.UserClaims).UserID
	keyKafka := []byte(uuid.New().String())
	messageKafka, _ := json.Marshal(request)

	err = u.kafkaClient.WriteToTopic(constants.EXPORT_SALES_ORDER_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return "", errorLogData
	}

	return request.FileName, nil
}

func (u *salesOrderUseCase) ExportDetail(request *models.SalesOrderDetailExportRequest, ctx context.Context) (string, *model.ErrorLog) {
	soDetailRequest := &models.GetSalesOrderDetailRequest{}
	soDetailRequest.SalesOrderDetailExportMap(request)
	if soDetailRequest.SortField == "order_status" {
		request.SortField = "order_status.name"
		soDetailRequest.SortField = "order_status.name"
	}
	getSalesOrderDetailsCountResultChan := make(chan *models.SalesOrderDetailsOpenSearchChan)
	go u.salesOrderDetailOpenSearchRepository.Get(soDetailRequest, true, getSalesOrderDetailsCountResultChan)
	getSalesOrderDetailsCountResult := <-getSalesOrderDetailsCountResultChan

	if getSalesOrderDetailsCountResult.Error != nil {
		fmt.Println("error = ", getSalesOrderDetailsCountResult.Error)
		errorLogData := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{"Data tidak ditemukan"},
			SystemMessage: []string{"Sales Order Details data not found"},
			StatusCode:    http.StatusNotFound,
		})
		return "", errorLogData
	}

	if getSalesOrderDetailsCountResult.Total == 0 {
		err := helper.NewError("Data tidak ditemukan")
		errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
		errorLogData = helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{"Data tidak ditemukan"},
			SystemMessage: []string{"Sales Order Details data not found"},
			StatusCode:    http.StatusNotFound,
		})
		return "", errorLogData
	}
	rand, err := helper.Generate(`[A-Za-z]{12}`)
	loc, err := time.LoadLocation("Asia/Jakarta")
	fileHour := time.Now().In(loc).Format(constants.DATE_FORMAT_EXPORT)
	if ctx == nil {
		err = fmt.Errorf("nil context")
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return "", errorLogData
	}
	fileName := fmt.Sprintf("SO-LIST-DETAIL-%s-%d-%s", fileHour, ctx.Value("user").(*models.UserClaims).UserID, rand)
	request.FileName = fileName
	request.UserID = ctx.Value("user").(*models.UserClaims).UserID
	keyKafka := []byte(uuid.New().String())
	messageKafka, _ := json.Marshal(request)

	err = u.kafkaClient.WriteToTopic(constants.EXPORT_SALES_ORDER_DETAIL_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return "", errorLogData
	}

	return request.FileName, nil
}
