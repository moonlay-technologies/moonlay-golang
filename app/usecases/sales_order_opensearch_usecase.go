package usecases

import (
	"context"
	"database/sql"
	"net/http"
	"order-service/app/models"
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

type SalesOrderOpenSearchUseCaseInterface interface {
	SyncToOpenSearchFromCreateEvent(salesOrder *models.SalesOrder, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromUpdateEvent(salesOrder *models.SalesOrder, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromDeleteEvent(salesOrder *models.SalesOrder, ctx context.Context) *model.ErrorLog
}

type SalesOrderOpenSearchUseCase struct {
	salesOrderRepository              repositories.SalesOrderRepositoryInterface
	salesOrderDetailRepository        repositories.SalesOrderDetailRepositoryInterface
	orderStatusRepository             repositories.OrderStatusRepositoryInterface
	orderSourceRepository             repositories.OrderSourceRepositoryInterface
	agentRepository                   repositories.AgentRepositoryInterface
	brandRepository                   repositories.BrandRepositoryInterface
	storeRepository                   repositories.StoreRepositoryInterface
	productRepository                 repositories.ProductRepositoryInterface
	uomRepository                     repositories.UomRepositoryInterface
	deliveryOrderRepository           repositories.DeliveryOrderRepositoryInterface
	salesOrderLogRepository           mongoRepositories.SalesOrderLogRepositoryInterface
	userRepository                    repositories.UserRepositoryInterface
	salesmanRepository                repositories.SalesmanRepositoryInterface
	categoryRepository                repositories.CategoryRepositoryInterface
	salesOrderOpenSearchRepository    openSearchRepositories.SalesOrderOpenSearchRepositoryInterface
	deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface
	kafkaClient                       kafkadbo.KafkaClientInterface
	db                                dbresolver.DB
	ctx                               context.Context
}

func InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, orderSourceRepository repositories.OrderSourceRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, uomRepository repositories.UomRepositoryInterface, deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, salesOrderLogRepository mongoRepositories.SalesOrderLogRepositoryInterface, userRepository repositories.UserRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, categoryRepository repositories.CategoryRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface, kafkaClient kafkadbo.KafkaClientInterface, db dbresolver.DB, ctx context.Context) SalesOrderOpenSearchUseCaseInterface {
	return &SalesOrderOpenSearchUseCase{
		salesOrderRepository:              salesOrderRepository,
		salesOrderDetailRepository:        salesOrderDetailRepository,
		orderStatusRepository:             orderStatusRepository,
		orderSourceRepository:             orderSourceRepository,
		agentRepository:                   agentRepository,
		brandRepository:                   brandRepository,
		storeRepository:                   storeRepository,
		productRepository:                 productRepository,
		uomRepository:                     uomRepository,
		deliveryOrderRepository:           deliveryOrderRepository,
		salesOrderLogRepository:           salesOrderLogRepository,
		userRepository:                    userRepository,
		salesmanRepository:                salesmanRepository,
		categoryRepository:                categoryRepository,
		salesOrderOpenSearchRepository:    salesOrderOpenSearchRepository,
		deliveryOrderOpenSearchRepository: deliveryOrderOpenSearchRepository,
		kafkaClient:                       kafkaClient,
		db:                                db,
		ctx:                               ctx,
	}
}

func (u *SalesOrderOpenSearchUseCase) SyncToOpenSearchFromCreateEvent(salesOrder *models.SalesOrder, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog {
	now := time.Now()

	// deliveryOrdersRequest := &models.DeliveryOrderRequest{
	// 	SalesOrderID: salesOrder.ID,
	// }

	// getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
	// go u.deliveryOrderOpenSearchRepository.GetBySalesOrderID(deliveryOrdersRequest, getDeliveryOrdersResultChan)
	// getDeliveryOrdersResult := <-getDeliveryOrdersResultChan
	// deliveryOrdersFound := true

	// if getDeliveryOrdersResult.Error != nil {
	// 	deliveryOrdersFound = false
	// 	if !strings.Contains(getDeliveryOrdersResult.Error.Error(), "not found") {
	// 		errorLogData := helper.WriteLog(getDeliveryOrdersResult.Error, http.StatusInternalServerError, nil)
	// 		return errorLogData
	// 	}
	// }

	// if deliveryOrdersFound == true {
	// 	salesOrder.DeliveryOrders = getDeliveryOrdersResult.DeliveryOrders
	// }

	for k, v := range salesOrder.SalesOrderDetails {
		getProductResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetByID(v.ProductID, false, ctx, getProductResultChan)
		getProductResult := <-getProductResultChan

		if getProductResult.Error != nil {
			errorLogData := helper.WriteLog(getProductResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		salesOrder.SalesOrderDetails[k].Product = getProductResult.Product

		getUomResultChan := make(chan *models.UomChan)
		go u.uomRepository.GetByID(v.UomID, false, ctx, getUomResultChan)
		getUomResult := <-getUomResultChan

		if getUomResult.Error != nil {
			errorLogData := helper.WriteLog(getUomResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		salesOrder.SalesOrderDetails[k].Uom = getUomResult.Uom

		getOrderStatusDetailResultChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByID(v.OrderStatusID, false, ctx, getOrderStatusDetailResultChan)
		getOrderStatusDetailResult := <-getOrderStatusDetailResultChan

		if getOrderStatusDetailResult.Error != nil {
			errorLogData := helper.WriteLog(getOrderStatusDetailResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		salesOrder.SalesOrderDetails[k].EndDateSyncToEs = &now
		salesOrder.SalesOrderDetails[k].IsDoneSyncToEs = "1"
		salesOrder.SalesOrderDetails[k].OrderStatus = getOrderStatusDetailResult.OrderStatus

		salesOrderDetailUpdateData := &models.SalesOrderDetail{
			UpdatedAt:       &now,
			IsDoneSyncToEs:  "1",
			EndDateSyncToEs: &now,
		}

		updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.UpdateByID(v.ID, salesOrderDetailUpdateData, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
		updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

		if updateSalesOrderDetailResult.Error != nil {
			errorLogData := helper.WriteLog(updateSalesOrderDetailResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}
	}

	salesOrder.IsDoneSyncToEs = "1"
	salesOrder.EndDateSyncToEs = &now
	salesOrder.UpdatedAt = &now
	salesOrder.EndCreatedDate = &now

	createSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderOpenSearchRepository.Create(salesOrder, createSalesOrderResultChan)
	createSalesOrderResult := <-createSalesOrderResultChan

	if createSalesOrderResult.Error != nil {
		return createSalesOrderResult.ErrorLog
	}

	salesOrderUpdateData := &models.SalesOrder{
		UpdatedAt:       &now,
		IsDoneSyncToEs:  "1",
		EndCreatedDate:  &now,
		EndDateSyncToEs: &now,
	}

	updateSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.UpdateByID(salesOrder.ID, salesOrderUpdateData, sqlTransaction, ctx, updateSalesOrderResultChan)
	updateSalesOrderResult := <-updateSalesOrderResultChan

	if updateSalesOrderResult.Error != nil {
		errorLogData := helper.WriteLog(updateSalesOrderResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return &model.ErrorLog{}
}

func (u *SalesOrderOpenSearchUseCase) SyncToOpenSearchFromUpdateEvent(salesOrder *models.SalesOrder, ctx context.Context) *model.ErrorLog {
	now := time.Now()

	deliveryOrdersRequest := &models.DeliveryOrderRequest{
		SalesOrderID: salesOrder.ID,
	}

	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
	go u.deliveryOrderOpenSearchRepository.GetBySalesOrderID(deliveryOrdersRequest, getDeliveryOrdersResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrdersResultChan
	deliveryOrdersFound := true

	if getDeliveryOrdersResult.Error != nil {
		deliveryOrdersFound = false
		if !strings.Contains(getDeliveryOrdersResult.Error.Error(), "not found") {
			errorLogData := helper.WriteLog(getDeliveryOrdersResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}
	}

	if deliveryOrdersFound == true {

		for x := range getDeliveryOrdersResult.DeliveryOrders {
			getDeliveryOrdersResult.DeliveryOrders[x].SalesOrder = nil
		}

		salesOrder.DeliveryOrders = getDeliveryOrdersResult.DeliveryOrders
	}

	for k, v := range salesOrder.SalesOrderDetails {
		getProductResultChan := make(chan *models.ProductChan)
		go u.productRepository.GetByID(v.ProductID, false, ctx, getProductResultChan)
		getProductResult := <-getProductResultChan

		if getProductResult.Error != nil {
			errorLogData := helper.WriteLog(getProductResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		salesOrder.SalesOrderDetails[k].Product = getProductResult.Product

		getUomResultChan := make(chan *models.UomChan)
		go u.uomRepository.GetByID(v.UomID, false, ctx, getUomResultChan)
		getUomResult := <-getUomResultChan

		if getUomResult.Error != nil {
			errorLogData := helper.WriteLog(getUomResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		salesOrder.SalesOrderDetails[k].Uom = getUomResult.Uom

		getOrderStatusDetailResultChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByID(v.OrderStatusID, false, ctx, getOrderStatusDetailResultChan)
		getOrderStatusDetailResult := <-getOrderStatusDetailResultChan

		if getOrderStatusDetailResult.Error != nil {
			errorLogData := helper.WriteLog(getOrderStatusDetailResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		salesOrder.SalesOrderDetails[k].OrderStatus = getOrderStatusDetailResult.OrderStatus
	}

	salesOrder.UpdatedAt = &now

	updateSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderOpenSearchRepository.Create(salesOrder, updateSalesOrderResultChan)
	updateSalesOrderResult := <-updateSalesOrderResultChan

	if updateSalesOrderResult.Error != nil {
		return updateSalesOrderResult.ErrorLog
	}

	return &model.ErrorLog{}
}

func (u *SalesOrderOpenSearchUseCase) SyncToOpenSearchFromDeleteEvent(salesOrder *models.SalesOrder, ctx context.Context) *model.ErrorLog {
	now := time.Now()
	salesOrderRequest := &models.SalesOrderRequest{ID: salesOrder.ID}

	getSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderOpenSearchRepository.GetByID(salesOrderRequest, getSalesOrderResultChan)
	getSalesOrderResult := <-getSalesOrderResultChan

	salesOrder.SalesOrderOpenSearchChanMap(getSalesOrderResult)

	for k := range salesOrder.SalesOrderDetails {
		salesOrder.SalesOrderDetails[k].IsDoneSyncToEs = "1"
		salesOrder.SalesOrderDetails[k].EndDateSyncToEs = &now
	}

	salesOrder.UpdatedAt = &now
	salesOrder.IsDoneSyncToEs = "1"
	salesOrder.EndDateSyncToEs = &now

	updateSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderOpenSearchRepository.Create(salesOrder, updateSalesOrderResultChan)
	updateSalesOrderResult := <-updateSalesOrderResultChan

	if updateSalesOrderResult.Error != nil {
		return updateSalesOrderResult.ErrorLog
	}

	return &model.ErrorLog{}
}
