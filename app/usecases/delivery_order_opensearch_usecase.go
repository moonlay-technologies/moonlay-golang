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

type DeliveryOrderOpenSearchUseCaseInterface interface {
	SyncToOpenSearchFromCreateEvent(deliveryOrder *models.DeliveryOrder, salesOrderUseCase SalesOrderUseCaseInterface, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromUpdateEvent(deliveryOrder *models.DeliveryOrder, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromDeleteEvent(deliveryOrderId *int, ctx context.Context) *model.ErrorLog
}

type deliveryOrderOpenSearchUseCase struct {
	salesOrderRepository              repositories.SalesOrderRepositoryInterface
	salesOrderDetailRepository        repositories.SalesOrderDetailRepositoryInterface
	orderStatusRepository             repositories.OrderStatusRepositoryInterface
	brandRepository                   repositories.BrandRepositoryInterface
	uomRepository                     repositories.UomRepositoryInterface
	agentRepository                   repositories.AgentRepositoryInterface
	storeRepository                   repositories.StoreRepositoryInterface
	productRepository                 repositories.ProductRepositoryInterface
	deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface
	salesOrderOpenSearchRepository    openSearchRepositories.SalesOrderOpenSearchRepositoryInterface
	SalesOrderOpenSearchUseCase       SalesOrderOpenSearchUseCaseInterface
	db                                dbresolver.DB
	ctx                               context.Context
}

func InitDeliveryOrderOpenSearchUseCaseInterface(salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, uomRepository repositories.UomRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, salesOrderOpenSearchUseCase SalesOrderOpenSearchUseCaseInterface, db dbresolver.DB, ctx context.Context) DeliveryOrderOpenSearchUseCaseInterface {
	return &deliveryOrderOpenSearchUseCase{
		salesOrderRepository:              salesOrderRepository,
		salesOrderDetailRepository:        salesOrderDetailRepository,
		orderStatusRepository:             orderStatusRepository,
		brandRepository:                   brandRepository,
		uomRepository:                     uomRepository,
		productRepository:                 productRepository,
		agentRepository:                   agentRepository,
		storeRepository:                   storeRepository,
		deliveryOrderOpenSearchRepository: deliveryOrderOpenSearchRepository,
		salesOrderOpenSearchRepository:    salesOrderOpenSearchRepository,
		SalesOrderOpenSearchUseCase:       salesOrderOpenSearchUseCase,
		db:                                db,
		ctx:                               ctx,
	}
}

func (u *deliveryOrderOpenSearchUseCase) SyncToOpenSearchFromCreateEvent(deliveryOrder *models.DeliveryOrder, salesOrderUseCase SalesOrderUseCaseInterface, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog {
	now := time.Now()

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
		removeCacheSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.RemoveCacheByID(v.SoDetailID, ctx, removeCacheSalesOrderDetailResultChan)
		removeCacheSalesOrderDetailResult := <-removeCacheSalesOrderDetailResultChan

		if removeCacheSalesOrderDetailResult.Error != nil {
			errorLogData := helper.WriteLog(removeCacheSalesOrderDetailResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		v.EndDateSyncToEs = &now
		v.IsDoneSyncToEs = "1"
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

	return &model.ErrorLog{}
}

func (u *deliveryOrderOpenSearchUseCase) SyncToOpenSearchFromUpdateEvent(deliveryOrder *models.DeliveryOrder, ctx context.Context) *model.ErrorLog {
	now := time.Now()

	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(deliveryOrder.AgentID, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	deliveryOrder.Agent = getAgentResult.Agent

	getStoreResultChan := make(chan *models.StoreChan)
	go u.storeRepository.GetByID(deliveryOrder.StoreID, false, ctx, getStoreResultChan)
	getStoreResult := <-getStoreResultChan

	if getStoreResult.Error != nil {
		errorLogData := helper.WriteLog(getStoreResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	deliveryOrder.Store = getStoreResult.Store

	for _, v := range deliveryOrder.DeliveryOrderDetails {
		getSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.GetByID(v.SoDetailID, false, ctx, getSalesOrderDetailResultChan)
		getSalesOrderDetailResult := <-getSalesOrderDetailResultChan

		if getSalesOrderDetailResult.Error != nil {
			errorLogData := helper.WriteLog(getSalesOrderDetailResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		getProductResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetByID(v.ProductID, false, ctx, getProductResultChan)
		getProductResult := <-getProductResultChan

		if getProductResult.Error != nil {
			errorLogData := helper.WriteLog(getProductResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		v.Product = getProductResult.Product

		getUomResultChan := make(chan *models.UomChan)
		go u.uomRepository.GetByID(v.UomID, false, ctx, getUomResultChan)
		getUomResult := <-getUomResultChan

		if getUomResult.Error != nil {
			errorLogData := helper.WriteLog(getUomResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		v.Uom = getUomResult.Uom

		getBrandResultChan := make(chan *models.BrandChan)
		go u.brandRepository.GetByID(v.UomID, false, ctx, getBrandResultChan)
		getBrandResult := <-getBrandResultChan

		if getBrandResult.Error != nil {
			errorLogData := helper.WriteLog(getBrandResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		v.Brand = getBrandResult.Brand
	}

	deliveryOrder.UpdatedAt = &now

	createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.Create(deliveryOrder, createDeliveryOrderResultChan)
	createDeliveryOrderResult := <-createDeliveryOrderResultChan

	if createDeliveryOrderResult.Error != nil {
		fmt.Println(createDeliveryOrderResult.ErrorLog.Err.Error())
		return createDeliveryOrderResult.ErrorLog
	}

	return &model.ErrorLog{}
}

func (u *deliveryOrderOpenSearchUseCase) SyncToOpenSearchFromDeleteEvent(deliveryOrderId *int, ctx context.Context) *model.ErrorLog {
	now := time.Now()
	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.GetByID(&models.DeliveryOrderRequest{ID: *deliveryOrderId}, getDeliveryOrdersResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrdersResultChan

	if getDeliveryOrdersResult.Error != nil {
		if !strings.Contains(getDeliveryOrdersResult.Error.Error(), "not found") {
			errorLogData := helper.WriteLog(getDeliveryOrdersResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}
	}
	deliveryOrder := getDeliveryOrdersResult.DeliveryOrder

	for _, v := range deliveryOrder.DeliveryOrderDetails {
		v.DeletedAt = &now
		v.UpdatedAt = &now
	}

	deliveryOrder.DeletedAt = &now
	deliveryOrder.UpdatedAt = &now

	createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.Create(deliveryOrder, createDeliveryOrderResultChan)
	deleteDeliveryOrderResult := <-createDeliveryOrderResultChan

	if deleteDeliveryOrderResult.Error != nil {
		fmt.Println(deleteDeliveryOrderResult.ErrorLog.Err.Error())
		return deleteDeliveryOrderResult.ErrorLog
	}

	salesOrderRequest := &models.SalesOrderRequest{
		ID:            deliveryOrder.SalesOrderID,
		OrderSourceID: deliveryOrder.OrderSourceID,
	}

	getSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderOpenSearchRepository.GetByID(salesOrderRequest, getSalesOrderResultChan)
	getSalesOrderResult := <-getSalesOrderResultChan

	if getSalesOrderResult.ErrorLog != nil {
		deleteDeliveryOrderResult.ErrorLog = getSalesOrderResult.ErrorLog
		fmt.Println(deleteDeliveryOrderResult.ErrorLog.Err.Error())
		return deleteDeliveryOrderResult.ErrorLog
	}

	salesOrderWithDetail := getSalesOrderResult.SalesOrder
	for k, v := range salesOrderWithDetail.SalesOrderDetails {
		v.SentQty -= deliveryOrder.DeliveryOrderDetails[k].Qty
		v.ResidualQty += deliveryOrder.DeliveryOrderDetails[k].Qty
	}
	deleteDeliveryOrderResult.DeliveryOrder.SalesOrder = salesOrderWithDetail
	deleteDeliveryOrderResult.ErrorLog = u.SalesOrderOpenSearchUseCase.SyncToOpenSearchFromUpdateEvent(salesOrderWithDetail, ctx)

	if deleteDeliveryOrderResult.ErrorLog != nil {
		fmt.Println(deleteDeliveryOrderResult.ErrorLog.Err.Error())
		return deleteDeliveryOrderResult.ErrorLog
	}

	deleteDeliveryOrderResult.ErrorLog = u.SyncToOpenSearchFromUpdateEvent(deleteDeliveryOrderResult.DeliveryOrder, ctx)

	if deleteDeliveryOrderResult.ErrorLog != nil {
		fmt.Println(deleteDeliveryOrderResult.ErrorLog.Err.Error())
		return deleteDeliveryOrderResult.ErrorLog
	}

	return &model.ErrorLog{}
}
