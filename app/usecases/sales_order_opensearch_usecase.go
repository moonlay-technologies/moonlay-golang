package usecases

import (
	"context"
	"database/sql"
	"net/http"
	"order-service/app/models"
	"order-service/app/repositories"
	openSearchRepositories "order-service/app/repositories/open_search"
	"order-service/global/utils/helper"
	"order-service/global/utils/model"
	"strings"
	"time"
)

type SalesOrderOpenSearchUseCaseInterface interface {
	SyncToOpenSearchFromCreateEvent(salesOrder *models.SalesOrder, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromUpdateEvent(salesOrder *models.SalesOrder, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromDeleteEvent(salesOrder *models.SalesOrder, ctx context.Context) *model.ErrorLog
	SyncDetailToOpenSearchFromCreateEvent(salesOrderDetail *models.SalesOrderDetail, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncDetailToOpenSearchFromUpdateEvent(salesOrder *models.SalesOrder, salesOrderDetail *models.SalesOrderDetail, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncDetailToOpenSearchFromDeleteEvent(salesOrderDetail *models.SalesOrderDetail, ctx context.Context) *model.ErrorLog
}

type SalesOrderOpenSearchUseCase struct {
	salesOrderRepository                 repositories.SalesOrderRepositoryInterface
	salesOrderDetailRepository           repositories.SalesOrderDetailRepositoryInterface
	orderStatusRepository                repositories.OrderStatusRepositoryInterface
	productRepository                    repositories.ProductRepositoryInterface
	uomRepository                        repositories.UomRepositoryInterface
	salesOrderOpenSearchRepository       openSearchRepositories.SalesOrderOpenSearchRepositoryInterface
	salesOrderDetailOpenSearchRepository openSearchRepositories.SalesOrderDetailOpenSearchRepositoryInterface
	deliveryOrderOpenSearchRepository    openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface
	ctx                                  context.Context
}

func InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, productRepository repositories.ProductRepositoryInterface, uomRepository repositories.UomRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, salesOrderDetailOpenSearchRepository openSearchRepositories.SalesOrderDetailOpenSearchRepositoryInterface, deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface) SalesOrderOpenSearchUseCaseInterface {
	return &SalesOrderOpenSearchUseCase{
		salesOrderRepository:                 salesOrderRepository,
		salesOrderDetailRepository:           salesOrderDetailRepository,
		orderStatusRepository:                orderStatusRepository,
		productRepository:                    productRepository,
		uomRepository:                        uomRepository,
		salesOrderOpenSearchRepository:       salesOrderOpenSearchRepository,
		salesOrderDetailOpenSearchRepository: salesOrderDetailOpenSearchRepository,
		deliveryOrderOpenSearchRepository:    deliveryOrderOpenSearchRepository,
	}
}

func (u *SalesOrderOpenSearchUseCase) SyncToOpenSearchFromCreateEvent(salesOrder *models.SalesOrder, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog {
	now := time.Now()

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

		salesOrderDetail := &models.SalesOrderDetailOpenSearch{}
		salesOrderDetail.SalesOrderDetailOpenSearchMap(salesOrder, v)

		createSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailOpenSearchChan)
		go u.salesOrderDetailOpenSearchRepository.Create(salesOrderDetail, createSalesOrderDetailResultChan)
		createSalesOrderDetailResult := <-createSalesOrderDetailResultChan

		if createSalesOrderDetailResult.Error != nil {
			return createSalesOrderDetailResult.ErrorLog
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
func (u *SalesOrderOpenSearchUseCase) SyncDetailToOpenSearchFromCreateEvent(salesOrderDetail *models.SalesOrderDetail, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog {
	return &model.ErrorLog{}
}

func (u *SalesOrderOpenSearchUseCase) SyncDetailToOpenSearchFromUpdateEvent(salesOrder *models.SalesOrder, salesOrderDetail *models.SalesOrderDetail, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog {
	now := time.Now()

	getProductResultChan := make(chan *models.ProductChan)
	go u.productRepository.GetByID(salesOrderDetail.ProductID, false, ctx, getProductResultChan)
	getProductResult := <-getProductResultChan

	if getProductResult.Error != nil {
		errorLogData := helper.WriteLog(getProductResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	salesOrderDetail.Product = getProductResult.Product

	getUomResultChan := make(chan *models.UomChan)
	go u.uomRepository.GetByID(salesOrderDetail.UomID, false, ctx, getUomResultChan)
	getUomResult := <-getUomResultChan

	if getUomResult.Error != nil {
		errorLogData := helper.WriteLog(getUomResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	salesOrderDetail.Uom = getUomResult.Uom

	getOrderStatusDetailResultChan := make(chan *models.OrderStatusChan)
	go u.orderStatusRepository.GetByID(salesOrderDetail.OrderStatusID, false, ctx, getOrderStatusDetailResultChan)
	getOrderStatusDetailResult := <-getOrderStatusDetailResultChan

	if getOrderStatusDetailResult.Error != nil {
		errorLogData := helper.WriteLog(getOrderStatusDetailResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	salesOrderDetailOpenSearch := &models.SalesOrderDetailOpenSearch{}
	salesOrderDetailOpenSearch.SalesOrderDetailOpenSearchMap(salesOrder, salesOrderDetail)

	createSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailOpenSearchChan)
	go u.salesOrderDetailOpenSearchRepository.Create(salesOrderDetailOpenSearch, createSalesOrderDetailResultChan)
	createSalesOrderDetailResult := <-createSalesOrderDetailResultChan

	if createSalesOrderDetailResult.Error != nil {
		return createSalesOrderDetailResult.ErrorLog
	}

	salesOrderDetailUpdateData := &models.SalesOrderDetail{
		UpdatedAt:       &now,
		IsDoneSyncToEs:  "1",
		EndDateSyncToEs: &now,
	}

	updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
	go u.salesOrderDetailRepository.UpdateByID(salesOrderDetail.ID, salesOrderDetailUpdateData, sqlTransaction, ctx, updateSalesOrderDetailResultChan)
	updateSalesOrderDetailResult := <-updateSalesOrderDetailResultChan

	if updateSalesOrderDetailResult.Error != nil {
		errorLogData := helper.WriteLog(updateSalesOrderDetailResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return &model.ErrorLog{}
}

func (u *SalesOrderOpenSearchUseCase) SyncDetailToOpenSearchFromDeleteEvent(salesOrderDetail *models.SalesOrderDetail, ctx context.Context) *model.ErrorLog {
	now := time.Now()
	salesOrderRequest := &models.SalesOrderRequest{ID: salesOrderDetail.SalesOrderID}
	x := models.SalesOrderDetailOpenSearch{}

	getSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderOpenSearchRepository.GetByID(salesOrderRequest, getSalesOrderResultChan)
	getSalesOrderResult := <-getSalesOrderResultChan
	x.SalesOrderMap(getSalesOrderResult.SalesOrder)

	for v, k := range getSalesOrderResult.SalesOrder.SalesOrderDetails {
		if salesOrderDetail.ID == k.ID {
			k.DeletedAt = &now
			k.IsDoneSyncToEs = "1"
			k.EndDateSyncToEs = &now
			x.SalesOrderDetailMap(k)
		}
		getSalesOrderResult.SalesOrder.SalesOrderDetails[v] = k
	}

	updateSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderOpenSearchRepository.Create(getSalesOrderResult.SalesOrder, updateSalesOrderResultChan)
	updateSalesOrderResult := <-updateSalesOrderResultChan

	if updateSalesOrderResult.Error != nil {
		return updateSalesOrderResult.ErrorLog
	}

	return &model.ErrorLog{}
}
