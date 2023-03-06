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

type DeliveryOrderDetailOpenSearchUseCaseInterface interface {
	SyncToOpenSearchFromCreateDoEvent(deliveryOrder *models.DeliveryOrder, salesOrderUseCase SalesOrderUseCaseInterface, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromUpdateDoEvent(deliveryOrder *models.DeliveryOrder, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromDeleteDoEvent(deliveryOrderId *int, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromCreateEvent(deliveryOrderDetail *models.DeliveryOrderDetail, salesOrderUseCase SalesOrderUseCaseInterface, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromUpdateEvent(deliveryOrderDetail *models.DeliveryOrderDetail, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromDeleteEvent(Id *int, ctx context.Context) *model.ErrorLog
}

type DeliveryOrderDetailOpenSearchUseCase struct {
	orderStatusRepository                   repositories.OrderStatusRepositoryInterface
	brandRepository                         repositories.BrandRepositoryInterface
	uomRepository                           repositories.UomRepositoryInterface
	agentRepository                         repositories.AgentRepositoryInterface
	storeRepository                         repositories.StoreRepositoryInterface
	productRepository                       repositories.ProductRepositoryInterface
	deliveryOrderDetailOpenSearchRepository openSearchRepositories.DeliveryOrderDetailOpenSearchRepositoryInterface
	salesOrderOpenSearchRepository          openSearchRepositories.SalesOrderOpenSearchRepositoryInterface
	SalesOrderOpenSearchUseCase             SalesOrderOpenSearchUseCaseInterface
	db                                      dbresolver.DB
	ctx                                     context.Context
}

func InitDeliveryOrderDetailOpenSearchUseCaseInterface(salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, uomRepository repositories.UomRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, deliveryOrderDetailOpenSearchRepository openSearchRepositories.DeliveryOrderDetailOpenSearchRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, salesOrderOpenSearchUseCase SalesOrderOpenSearchUseCaseInterface, db dbresolver.DB, ctx context.Context) DeliveryOrderDetailOpenSearchUseCaseInterface {
	return &DeliveryOrderDetailOpenSearchUseCase{
		orderStatusRepository:                   orderStatusRepository,
		brandRepository:                         brandRepository,
		uomRepository:                           uomRepository,
		productRepository:                       productRepository,
		agentRepository:                         agentRepository,
		storeRepository:                         storeRepository,
		deliveryOrderDetailOpenSearchRepository: deliveryOrderDetailOpenSearchRepository,
		salesOrderOpenSearchRepository:          salesOrderOpenSearchRepository,
		SalesOrderOpenSearchUseCase:             salesOrderOpenSearchUseCase,
		db:                                      db,
		ctx:                                     ctx,
	}
}

func (u *DeliveryOrderDetailOpenSearchUseCase) SyncToOpenSearchFromCreateDoEvent(deliveryOrder *models.DeliveryOrder, salesOrderUseCase SalesOrderUseCaseInterface, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog {
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
		v.EndDateSyncToEs = &now
		v.IsDoneSyncToEs = "1"

		createDeliveryDetailOrderResultChan := make(chan *models.DeliveryOrderDetailChan)
		go u.deliveryOrderDetailOpenSearchRepository.Create(v, createDeliveryDetailOrderResultChan)
		createDeliveryOrderDetailResult := <-createDeliveryDetailOrderResultChan

		if createDeliveryOrderDetailResult.Error != nil {
			return createDeliveryOrderDetailResult.ErrorLog
		}

		return &model.ErrorLog{}
	}

	return &model.ErrorLog{}
}

func (u *DeliveryOrderDetailOpenSearchUseCase) SyncToOpenSearchFromUpdateDoEvent(deliveryOrder *models.DeliveryOrder, ctx context.Context) *model.ErrorLog {
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

		createDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
		go u.deliveryOrderDetailOpenSearchRepository.Create(v, createDeliveryOrderDetailResultChan)
		createDeliveryOrderDetailResult := <-createDeliveryOrderDetailResultChan

		if createDeliveryOrderDetailResult.Error != nil {
			fmt.Println(createDeliveryOrderDetailResult.ErrorLog.Err.Error())
			return createDeliveryOrderDetailResult.ErrorLog
		}
	}

	return &model.ErrorLog{}
}

func (u *DeliveryOrderDetailOpenSearchUseCase) SyncToOpenSearchFromDeleteDoEvent(deliveryOrderId *int, ctx context.Context) *model.ErrorLog {
	now := time.Now()
	getDeliveryOrderDetailsResultChan := make(chan *models.DeliveryOrderDetailsChan)
	go u.deliveryOrderDetailOpenSearchRepository.GetByDoID(&models.DeliveryOrderRequest{ID: *deliveryOrderId}, getDeliveryOrderDetailsResultChan)
	getDeliveryOrderDetailsResult := <-getDeliveryOrderDetailsResultChan

	if getDeliveryOrderDetailsResult.Error != nil {
		if !strings.Contains(getDeliveryOrderDetailsResult.Error.Error(), "not found") {
			errorLogData := helper.WriteLog(getDeliveryOrderDetailsResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}
	}
	deliveryOrderDetails := getDeliveryOrderDetailsResult.DeliveryOrderDetails

	for _, v := range deliveryOrderDetails {
		v.DeletedAt = &now
		v.UpdatedAt = &now

		createDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
		go u.deliveryOrderDetailOpenSearchRepository.Create(v, createDeliveryOrderDetailResultChan)
		deleteDeliveryOrderDetailResult := <-createDeliveryOrderDetailResultChan

		if deleteDeliveryOrderDetailResult.Error != nil {
			fmt.Println(deleteDeliveryOrderDetailResult.ErrorLog.Err.Error())
			return deleteDeliveryOrderDetailResult.ErrorLog
		}
	}

	return &model.ErrorLog{}
}
func (u *DeliveryOrderDetailOpenSearchUseCase) SyncToOpenSearchFromCreateEvent(deliveryOrderDetail *models.DeliveryOrderDetail, salesOrderUseCase SalesOrderUseCaseInterface, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog {
	return &model.ErrorLog{}
}
func (u *DeliveryOrderDetailOpenSearchUseCase) SyncToOpenSearchFromUpdateEvent(deliveryOrderDetail *models.DeliveryOrderDetail, ctx context.Context) *model.ErrorLog {
	return &model.ErrorLog{}
}
func (u *DeliveryOrderDetailOpenSearchUseCase) SyncToOpenSearchFromDeleteEvent(Id *int, ctx context.Context) *model.ErrorLog {
	return &model.ErrorLog{}
}
