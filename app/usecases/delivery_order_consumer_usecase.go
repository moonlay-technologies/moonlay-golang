package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/app/repositories"
	openSearchRepositories "order-service/app/repositories/open_search"
	"order-service/global/utils/helper"
	"order-service/global/utils/model"
	"strings"
	"time"

	"github.com/bxcodec/dbresolver"
)

type DeliveryOrderConsumerUseCaseInterface interface {
	SyncToOpenSearchFromCreateEvent(deliveryOrder *models.DeliveryOrder, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromUpdateEvent(deliveryOrder *models.DeliveryOrder, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromDeleteEvent(deliveryOrderId *int, deliveryOrderDetailId []*int, ctx context.Context) *model.ErrorLog
}

type deliveryOrderConsumerUseCase struct {
	salesOrderRepository                    repositories.SalesOrderRepositoryInterface
	salesOrderDetailRepository              repositories.SalesOrderDetailRepositoryInterface
	orderStatusRepository                   repositories.OrderStatusRepositoryInterface
	brandRepository                         repositories.BrandRepositoryInterface
	uomRepository                           repositories.UomRepositoryInterface
	agentRepository                         repositories.AgentRepositoryInterface
	storeRepository                         repositories.StoreRepositoryInterface
	productRepository                       repositories.ProductRepositoryInterface
	deliveryOrderOpenSearchRepository       openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface
	deliveryOrderDetailOpenSearchRepository openSearchRepositories.DeliveryOrderDetailOpenSearchRepositoryInterface
	salesOrderOpenSearchRepository          openSearchRepositories.SalesOrderOpenSearchRepositoryInterface
	SalesOrderOpenSearchUseCase             SalesOrderOpenSearchUseCaseInterface
	db                                      dbresolver.DB
	ctx                                     context.Context
}

func InitDeliveryOrderConsumerUseCaseInterface(salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, uomRepository repositories.UomRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface, deliveryOrderDetailOpenSearchRepository openSearchRepositories.DeliveryOrderDetailOpenSearchRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, salesOrderOpenSearchUseCase SalesOrderOpenSearchUseCaseInterface, db dbresolver.DB, ctx context.Context) DeliveryOrderConsumerUseCaseInterface {
	return &deliveryOrderConsumerUseCase{
		salesOrderRepository:                    salesOrderRepository,
		salesOrderDetailRepository:              salesOrderDetailRepository,
		orderStatusRepository:                   orderStatusRepository,
		brandRepository:                         brandRepository,
		uomRepository:                           uomRepository,
		productRepository:                       productRepository,
		agentRepository:                         agentRepository,
		storeRepository:                         storeRepository,
		deliveryOrderOpenSearchRepository:       deliveryOrderOpenSearchRepository,
		deliveryOrderDetailOpenSearchRepository: deliveryOrderDetailOpenSearchRepository,
		salesOrderOpenSearchRepository:          salesOrderOpenSearchRepository,
		SalesOrderOpenSearchUseCase:             salesOrderOpenSearchUseCase,
		db:                                      db,
		ctx:                                     ctx,
	}
}

func (u *deliveryOrderConsumerUseCase) SyncToOpenSearchFromCreateEvent(deliveryOrder *models.DeliveryOrder, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog {
	now := time.Now()

	getSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.GetByID(deliveryOrder.SalesOrderID, false, ctx, getSalesOrderResultChan)
	getSalesOrderResult := <-getSalesOrderResultChan

	if getSalesOrderResult.Error != nil {
		errorLogData := helper.WriteLog(getSalesOrderResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	deliveryOrder.SalesOrder = getSalesOrderResult.SalesOrder
	deliveryOrder.SalesOrder.SoCode = getSalesOrderResult.SalesOrder.SoCode
	deliveryOrder.SalesOrder.SoDate = getSalesOrderResult.SalesOrder.SoDate
	deliveryOrder.SalesOrder.SoRefDate = getSalesOrderResult.SalesOrder.SoRefDate

	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByID(deliveryOrder.OrderStatusID, false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		errorLogData := helper.WriteLog(getOrderStatusResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	deliveryOrder.OrderStatusID = getOrderStatusResult.OrderStatus.ID
	deliveryOrder.OrderStatusName = getOrderStatusResult.OrderStatus.Name
	deliveryOrder.OrderStatus = getOrderStatusResult.OrderStatus

	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(deliveryOrder.AgentID, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	deliveryOrder.Agent = getAgentResult.Agent
	deliveryOrder.AgentName = getAgentResult.Agent.Name

	getStoreResultChan := make(chan *models.StoreChan)
	go u.storeRepository.GetByID(deliveryOrder.StoreID, false, ctx, getStoreResultChan)
	getStoreResult := <-getStoreResultChan

	if getStoreResult.Error != nil {
		errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	deliveryOrder.Store = getStoreResult.Store

	for _, v := range deliveryOrder.DeliveryOrderDetails {

		getBrandResultChan := make(chan *models.BrandChan)
		go u.brandRepository.GetByID(deliveryOrder.SalesOrder.BrandID, false, ctx, getBrandResultChan)
		getBrandResult := <-getBrandResultChan

		if getBrandResult.Error != nil {
			errorLogData := helper.WriteLog(getBrandResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}
		v.Brand = getBrandResult.Brand

		getProductResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetByID(v.ProductID, false, ctx, getProductResultChan)
		getProductResult := <-getProductResultChan

		if getProductResult.Error != nil {
			errorLogData := helper.WriteLog(getProductResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}
		v.Product = getProductResult.Product

		v.EndDateSyncToEs = &now
		v.UpdatedAt = &now
		v.CreatedAt = &now
		v.IsDoneSyncToEs = "1"

		doDetailOpenSearch := &models.DeliveryOrderDetailOpenSearch{}
		doDetailOpenSearch.DoDetailMap(deliveryOrder, v)
		createDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailOpenSearchChan)
		go u.deliveryOrderDetailOpenSearchRepository.Create(doDetailOpenSearch, createDeliveryOrderDetailResultChan)
		createDeliveryOrderDetailResult := <-createDeliveryOrderDetailResultChan

		if createDeliveryOrderDetailResult.Error != nil {
			return createDeliveryOrderDetailResult.ErrorLog
		}
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

	createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.Create(deliveryOrder, createDeliveryOrderResultChan)
	createDeliveryOrderResult := <-createDeliveryOrderResultChan

	if createDeliveryOrderResult.Error != nil {
		return createDeliveryOrderResult.ErrorLog
	}

	deliveryOrder.SalesOrder.ID = deliveryOrder.SalesOrderID
	createDeliveryOrderResult.ErrorLog = u.SalesOrderOpenSearchUseCase.SyncToOpenSearchFromUpdateEvent(deliveryOrder.SalesOrder, ctx)

	if createDeliveryOrderResult.Error != nil {
		fmt.Println(createDeliveryOrderResult.Error)
		return createDeliveryOrderResult.ErrorLog
	}

	return &model.ErrorLog{}
}

func (u *deliveryOrderConsumerUseCase) SyncToOpenSearchFromUpdateEvent(request *models.DeliveryOrder, ctx context.Context) *model.ErrorLog {
	now := time.Now()

	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.GetByID(&models.DeliveryOrderRequest{ID: request.ID}, getDeliveryOrdersResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrdersResultChan

	if getDeliveryOrdersResult.Error != nil {
		if !strings.Contains(getDeliveryOrdersResult.Error.Error(), "not found") {
			errorLogData := helper.WriteLog(getDeliveryOrdersResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}
	}
	deliveryOrder := getDeliveryOrdersResult.DeliveryOrder
	deliveryOrder.DeliveryOrderUpdateMap(request)
	deliveryOrder.UpdatedAt = &now
	deliveryOrder.IsDoneSyncToEs = "1"
	deliveryOrder.EndDateSyncToEs = &now

	for _, v := range deliveryOrder.DeliveryOrderDetails {
		v.UpdatedAt = &now
		v.IsDoneSyncToEs = "1"
		v.EndDateSyncToEs = &now

		doDetailOpenSearch := &models.DeliveryOrderDetailOpenSearch{}
		doDetailOpenSearch.DoDetailMap(deliveryOrder, v)

		createDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailOpenSearchChan)
		go u.deliveryOrderDetailOpenSearchRepository.Create(doDetailOpenSearch, createDeliveryOrderDetailResultChan)
		createDeliveryOrderDetailResult := <-createDeliveryOrderDetailResultChan

		if createDeliveryOrderDetailResult.Error != nil {
			return createDeliveryOrderDetailResult.ErrorLog
		}
	}

	updateDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.Create(deliveryOrder, updateDeliveryOrderResultChan)
	updateDeliveryOrderResult := <-updateDeliveryOrderResultChan

	if updateDeliveryOrderResult.Error != nil {
		fmt.Println(updateDeliveryOrderResult.Error)
		return updateDeliveryOrderResult.ErrorLog
	}
	request.SalesOrder.ID = request.SalesOrderID
	updateDeliveryOrderResult.ErrorLog = u.SalesOrderOpenSearchUseCase.SyncToOpenSearchFromUpdateEvent(request.SalesOrder, ctx)

	if updateDeliveryOrderResult.Error != nil {
		fmt.Println(updateDeliveryOrderResult.ErrorLog.Err.Error())
		return updateDeliveryOrderResult.ErrorLog
	}

	return &model.ErrorLog{}
}

func (u *deliveryOrderConsumerUseCase) SyncToOpenSearchFromDeleteEvent(deliveryOrderId *int, deliveryOrderDetailIds []*int, ctx context.Context) *model.ErrorLog {
	now := time.Now()
	isDeleteParent := deliveryOrderDetailIds == nil
	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.GetByID(&models.DeliveryOrderRequest{ID: *deliveryOrderId}, getDeliveryOrdersResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrdersResultChan

	if getDeliveryOrdersResult.Error != nil {
		errorLogData := helper.WriteLog(getDeliveryOrdersResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}
	deliveryOrder := getDeliveryOrdersResult.DeliveryOrder

	salesOrderRequest := &models.SalesOrderRequest{
		ID:            deliveryOrder.SalesOrderID,
		OrderSourceID: deliveryOrder.OrderSourceID,
	}

	getSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderOpenSearchRepository.GetByID(salesOrderRequest, getSalesOrderResultChan)
	getSalesOrderResult := <-getSalesOrderResultChan

	if getSalesOrderResult.ErrorLog != nil {
		fmt.Println(getSalesOrderResult.ErrorLog.Err.Error())
		return getSalesOrderResult.ErrorLog
	}

	salesOrderWithDetail := getSalesOrderResult.SalesOrder
	for _, v := range salesOrderWithDetail.SalesOrderDetails {
		for _, w := range deliveryOrder.DeliveryOrderDetails {
			if v.ID == w.SoDetailID {
				if isDeleteParent {
					v.SentQty -= w.Qty
					v.ResidualQty += w.Qty
				} else {
					for _, x := range deliveryOrderDetailIds {
						if x == &w.ID {
							v.SentQty -= w.Qty
							v.ResidualQty += w.Qty
						}
					}
				}
			}
		}
	}

	errorLog := u.SalesOrderOpenSearchUseCase.SyncToOpenSearchFromUpdateEvent(salesOrderWithDetail, ctx)

	if errorLog != nil {
		fmt.Println(errorLog.Err.Error())
		return errorLog
	}

	for _, v := range deliveryOrder.DeliveryOrderDetails {
		isAllowUpdate := false
		if isDeleteParent {
			isAllowUpdate = true
		} else {
			for _, w := range deliveryOrderDetailIds {
				if *w == v.ID {
					isAllowUpdate = true
				}
			}
		}
		if isAllowUpdate {
			v.DeletedAt = &now
			v.Qty = 0
			v.IsDoneSyncToEs = "1"
			v.EndDateSyncToEs = &now

			v.OrderStatusID = 17

			getOrderStatusDetailChan := make(chan *models.OrderStatusChan)
			go u.orderStatusRepository.GetByID(v.OrderStatusID, false, ctx, getOrderStatusDetailChan)
			getOrderStatusDetailResult := <-getOrderStatusDetailChan

			v.OrderStatus = getOrderStatusDetailResult.OrderStatus
			v.OrderStatusName = getOrderStatusDetailResult.OrderStatus.Name

			doDetailOpenSearch := &models.DeliveryOrderDetailOpenSearch{}
			doDetailOpenSearch.DoDetailMap(deliveryOrder, v)

			createDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailOpenSearchChan)
			go u.deliveryOrderDetailOpenSearchRepository.Create(doDetailOpenSearch, createDeliveryOrderDetailResultChan)
			createDeliveryOrderDetailResult := <-createDeliveryOrderDetailResultChan

			if createDeliveryOrderDetailResult.Error != nil {
				return createDeliveryOrderDetailResult.ErrorLog
			}
		}
	}

	if isDeleteParent {
		deliveryOrder.DeletedAt = &now
		deliveryOrder.OrderStatusID = 17

		getOrderStatusChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByID(deliveryOrder.OrderStatusID, false, ctx, getOrderStatusChan)
		getOrderStatusResult := <-getOrderStatusChan

		deliveryOrder.OrderStatusName = getOrderStatusResult.OrderStatus.Name
		deliveryOrder.OrderStatus = getOrderStatusResult.OrderStatus
	}

	createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.Create(deliveryOrder, createDeliveryOrderResultChan)
	deleteDeliveryOrderResult := <-createDeliveryOrderResultChan

	if deleteDeliveryOrderResult.Error != nil {
		fmt.Println(deleteDeliveryOrderResult.ErrorLog.Err.Error())
		return deleteDeliveryOrderResult.ErrorLog
	}

	// deleteDeliveryOrderResult.ErrorLog = u.SyncToOpenSearchFromUpdateEvent(deleteDeliveryOrderResult.DeliveryOrder, ctx)

	// if deleteDeliveryOrderResult.ErrorLog != nil {
	// 	fmt.Println(deleteDeliveryOrderResult.ErrorLog.Err.Error())
	// 	return deleteDeliveryOrderResult.ErrorLog
	// }

	return &model.ErrorLog{}
}

func (u *deliveryOrderConsumerUseCase) Get(request *models.DeliveryOrderExportRequest) (*models.DeliveryOrdersOpenSearchResponse, *model.ErrorLog) {
	doRequest := &models.DeliveryOrderRequest{}
	doRequest.DeliveryOrderExportMap(request)
	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
	go u.deliveryOrderOpenSearchRepository.Get(doRequest, false, getDeliveryOrdersResultChan)
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
