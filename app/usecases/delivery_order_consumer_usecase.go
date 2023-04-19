package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/app/repositories"
	openSearchRepositories "order-service/app/repositories/open_search"
	"order-service/global/utils/helper"
	"order-service/global/utils/model"
	"os"
	"strconv"
	"strings"
	"time"
)

type DeliveryOrderConsumerUseCaseInterface interface {
	SyncToOpenSearchFromCreateEvent(deliveryOrder *models.DeliveryOrder, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromUpdateEvent(deliveryOrder *models.DeliveryOrder, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromDeleteEvent(deliveryOrderId *int, deliveryOrderDetailId []*int, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	Get(request *models.DeliveryOrderExportRequest) *model.ErrorLog
	GetDetail(request *models.DeliveryOrderDetailExportRequest) *model.ErrorLog
}

type deliveryOrderConsumerUseCase struct {
	uploadRepository                        repositories.UploadRepositoryInterface
	orderStatusRepository                   repositories.OrderStatusRepositoryInterface
	brandRepository                         repositories.BrandRepositoryInterface
	agentRepository                         repositories.AgentRepositoryInterface
	storeRepository                         repositories.StoreRepositoryInterface
	productRepository                       repositories.ProductRepositoryInterface
	deliveryOrderOpenSearchRepository       openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface
	deliveryOrderDetailOpenSearchRepository openSearchRepositories.DeliveryOrderDetailOpenSearchRepositoryInterface
	salesOrderOpenSearchRepository          openSearchRepositories.SalesOrderOpenSearchRepositoryInterface
	SalesOrderOpenSearchUseCase             SalesOrderOpenSearchUseCaseInterface
	pusherRepository                        repositories.PusherRepositoryInterface
	ctx                                     context.Context
}

func InitDeliveryOrderConsumerUseCaseInterface(uploadRepository repositories.UploadRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface, deliveryOrderDetailOpenSearchRepository openSearchRepositories.DeliveryOrderDetailOpenSearchRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, salesOrderOpenSearchUseCase SalesOrderOpenSearchUseCaseInterface, ctx context.Context) DeliveryOrderConsumerUseCaseInterface {
	return &deliveryOrderConsumerUseCase{
		uploadRepository:                        uploadRepository,
		orderStatusRepository:                   orderStatusRepository,
		brandRepository:                         brandRepository,
		productRepository:                       productRepository,
		agentRepository:                         agentRepository,
		storeRepository:                         storeRepository,
		deliveryOrderOpenSearchRepository:       deliveryOrderOpenSearchRepository,
		deliveryOrderDetailOpenSearchRepository: deliveryOrderDetailOpenSearchRepository,
		salesOrderOpenSearchRepository:          salesOrderOpenSearchRepository,
		pusherRepository:                        repositories.InitPusherRepository(),
		SalesOrderOpenSearchUseCase:             salesOrderOpenSearchUseCase,
		ctx:                                     ctx,
	}
}

func (u *deliveryOrderConsumerUseCase) SyncToOpenSearchFromCreateEvent(deliveryOrder *models.DeliveryOrder, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog {
	now := time.Now()

	salesOrderRequest := &models.SalesOrderRequest{
		ID:            deliveryOrder.SalesOrderID,
		OrderSourceID: deliveryOrder.OrderSourceID,
	}

	getSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderOpenSearchRepository.GetByID(salesOrderRequest, getSalesOrderResultChan)
	getSalesOrderResult := <-getSalesOrderResultChan

	if getSalesOrderResult.Error != nil {
		errorLogData := helper.WriteLog(getSalesOrderResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}
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
		go u.deliveryOrderDetailOpenSearchRepository.Create(doDetailOpenSearch, v, sqlTransaction, ctx, createDeliveryOrderDetailResultChan)
		createDeliveryOrderDetailResult := <-createDeliveryOrderDetailResultChan

		if createDeliveryOrderDetailResult.Error != nil {
			return createDeliveryOrderDetailResult.ErrorLog
		}
	}

	deliveryOrder.IsDoneSyncToEs = "1"
	deliveryOrder.EndDateSyncToEs = &now
	deliveryOrder.UpdatedAt = &now
	deliveryOrder.EndCreatedDate = &now

	createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.Create(deliveryOrder, sqlTransaction, ctx, createDeliveryOrderResultChan)
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

func (u *deliveryOrderConsumerUseCase) SyncToOpenSearchFromUpdateEvent(request *models.DeliveryOrder, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog {
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
		go u.deliveryOrderDetailOpenSearchRepository.Create(doDetailOpenSearch, v, sqlTransaction, ctx, createDeliveryOrderDetailResultChan)
		createDeliveryOrderDetailResult := <-createDeliveryOrderDetailResultChan

		if createDeliveryOrderDetailResult.Error != nil {
			return createDeliveryOrderDetailResult.ErrorLog
		}
	}

	updateDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.Create(deliveryOrder, sqlTransaction, ctx, updateDeliveryOrderResultChan)
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

func (u *deliveryOrderConsumerUseCase) SyncToOpenSearchFromDeleteEvent(deliveryOrderId *int, deliveryOrderDetailIds []*int, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog {
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

	if errorLog.Err != nil {
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

			v.OrderStatusID = 19

			getOrderStatusDetailChan := make(chan *models.OrderStatusChan)
			go u.orderStatusRepository.GetByID(v.OrderStatusID, false, ctx, getOrderStatusDetailChan)
			getOrderStatusDetailResult := <-getOrderStatusDetailChan

			v.OrderStatus = getOrderStatusDetailResult.OrderStatus
			v.OrderStatusName = getOrderStatusDetailResult.OrderStatus.Name

			doDetailOpenSearch := &models.DeliveryOrderDetailOpenSearch{}
			doDetailOpenSearch.DoDetailMap(deliveryOrder, v)

			createDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailOpenSearchChan)
			go u.deliveryOrderDetailOpenSearchRepository.Create(doDetailOpenSearch, v, sqlTransaction, ctx, createDeliveryOrderDetailResultChan)
			createDeliveryOrderDetailResult := <-createDeliveryOrderDetailResultChan

			if createDeliveryOrderDetailResult.Error != nil {
				return createDeliveryOrderDetailResult.ErrorLog
			}
		}
	}

	if isDeleteParent {
		deliveryOrder.DeletedAt = &now
		deliveryOrder.OrderStatusID = 19

		getOrderStatusChan := make(chan *models.OrderStatusChan)
		go u.orderStatusRepository.GetByID(deliveryOrder.OrderStatusID, false, ctx, getOrderStatusChan)
		getOrderStatusResult := <-getOrderStatusChan

		deliveryOrder.OrderStatusName = getOrderStatusResult.OrderStatus.Name
		deliveryOrder.OrderStatus = getOrderStatusResult.OrderStatus
	}

	createDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go u.deliveryOrderOpenSearchRepository.Create(deliveryOrder, sqlTransaction, ctx, createDeliveryOrderResultChan)
	deleteDeliveryOrderResult := <-createDeliveryOrderResultChan

	if deleteDeliveryOrderResult.Error != nil {
		fmt.Println(deleteDeliveryOrderResult.ErrorLog.Err.Error())
		return deleteDeliveryOrderResult.ErrorLog
	}

	return &model.ErrorLog{}
}

func (u *deliveryOrderConsumerUseCase) Get(request *models.DeliveryOrderExportRequest) *model.ErrorLog {
	countDoRequest := &models.DeliveryOrderRequest{}
	countDoRequest.DeliveryOrderExportMap(request)

	getDeliveryOrdersCountResultChan := make(chan *models.DeliveryOrdersChan)
	go u.deliveryOrderOpenSearchRepository.Get(countDoRequest, true, getDeliveryOrdersCountResultChan)
	getDeliveryOrdersCountResult := <-getDeliveryOrdersCountResultChan

	if getDeliveryOrdersCountResult.Error != nil {
		fmt.Println("error = ", getDeliveryOrdersCountResult.Error)
		return getDeliveryOrdersCountResult.ErrorLog
	}

	doRequest := &models.DeliveryOrderRequest{}
	doRequest.DeliveryOrderExportMap(request)
	perPage, err := strconv.Atoi(os.Getenv("EXPORT_PARTIAL"))
	if err != nil {
		perPage = constants.EXPORT_PARTIAL_DEFAULT
	}
	doRequest.PerPage = perPage
	instalmentData := math.Ceil(float64(getDeliveryOrdersCountResult.Total) / float64(doRequest.PerPage))
	data := [][]interface{}{constants.DELIVERY_ORDER_EXPORT_HEADER()}

	for i := 0; i < int(instalmentData); i++ {
		doRequest.Page = i + 1
		getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
		go u.deliveryOrderOpenSearchRepository.Get(doRequest, false, getDeliveryOrdersResultChan)
		getDeliveryOrdersResult := <-getDeliveryOrdersResultChan

		if getDeliveryOrdersResult.Error != nil {
			return getDeliveryOrdersResult.ErrorLog
		}
		for _, v := range getDeliveryOrdersResult.DeliveryOrders {
			data = append(data, v.MapToCsvRow())
		}
		progres := math.Round((float64((i*doRequest.PerPage)+len(getDeliveryOrdersResult.DeliveryOrders)) / float64(getDeliveryOrdersCountResult.Total)) * 100)
		fmt.Println(request.FileName, progres, "%")
	}

	b, err := helper.GenerateExportBufferFile(data, request.FileType)

	// Upload Files
	err = u.uploadRepository.UploadFile(b, constants.S3_EXPORT_DO_PATH, request.FileName, request.FileType)
	if err != nil {
		fmt.Println("error upload", err)
		return helper.WriteLog(err, http.StatusInternalServerError, nil)
	}

	pusherData := &models.Pusher{
		Subject: "Export DO Summary",
		Link:    fmt.Sprintf("%s/%s.%s", constants.DELIVERY_ORDER_EXPORT_PATH, request.FileName, request.FileType),
		Type:    "export",
		UserId:  strconv.Itoa(request.UserID),
	}

	err = u.pusherRepository.Publish(pusherData)
	if err != nil {
		fmt.Println(err.Error())
		return helper.WriteLog(err, http.StatusInternalServerError, nil)
	}

	return &model.ErrorLog{}
}

func (u *deliveryOrderConsumerUseCase) GetDetail(request *models.DeliveryOrderDetailExportRequest) *model.ErrorLog {
	countDoDetailRequest := &models.DeliveryOrderDetailOpenSearchRequest{}
	countDoDetailRequest.DeliveryOrderDetailExportMap(request)

	getDeliveryOrderDetailsCountResultChan := make(chan *models.DeliveryOrderDetailsOpenSearchChan)
	go u.deliveryOrderDetailOpenSearchRepository.Get(countDoDetailRequest, true, getDeliveryOrderDetailsCountResultChan)
	getDeliveryOrderDetailsCountResult := <-getDeliveryOrderDetailsCountResultChan

	if getDeliveryOrderDetailsCountResult.Error != nil {
		fmt.Println("error = ", getDeliveryOrderDetailsCountResult.Error)
		return getDeliveryOrderDetailsCountResult.ErrorLog
	}

	doDetailRequest := &models.DeliveryOrderDetailOpenSearchRequest{}
	doDetailRequest.DeliveryOrderDetailExportMap(request)
	perPage, err := strconv.Atoi(os.Getenv("EXPORT_PARTIAL"))
	if err != nil {
		perPage = constants.EXPORT_PARTIAL_DEFAULT
	}
	doDetailRequest.PerPage = perPage
	instalmentData := math.Ceil(float64(getDeliveryOrderDetailsCountResult.Total) / float64(doDetailRequest.PerPage))

	data := [][]interface{}{constants.DELIVERY_ORDER_DETAIL_EXPORT_HEADER()}

	for i := 0; i < int(instalmentData); i++ {
		doDetailRequest.Page = i + 1
		getDeliveryOrderDetailsResultChan := make(chan *models.DeliveryOrderDetailsOpenSearchChan)
		go u.deliveryOrderDetailOpenSearchRepository.Get(doDetailRequest, false, getDeliveryOrderDetailsResultChan)
		getDeliveryOrderDetailsResult := <-getDeliveryOrderDetailsResultChan

		if getDeliveryOrderDetailsResult.Error != nil {
			return getDeliveryOrderDetailsResult.ErrorLog
		}
		for _, v := range getDeliveryOrderDetailsResult.DeliveryOrderDetailOpenSearch {
			getDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
			go u.deliveryOrderOpenSearchRepository.GetByID(&models.DeliveryOrderRequest{ID: v.DeliveryOrderID}, getDeliveryOrderResultChan)
			getDeliveryOrdersResult := <-getDeliveryOrderResultChan

			if getDeliveryOrdersResult.Error != nil {
				return getDeliveryOrdersResult.ErrorLog
			}
			data = append(data, v.MapToCsvRow(getDeliveryOrdersResult.DeliveryOrder))
		}
		progres := math.Round((float64((i*doDetailRequest.PerPage)+len(getDeliveryOrderDetailsResult.DeliveryOrderDetailOpenSearch)) / float64(getDeliveryOrderDetailsCountResult.Total)) * 100)
		fmt.Println(request.FileName, progres, "%")
	}

	b, err := helper.GenerateExportBufferFile(data, request.FileType)

	// Upload Files
	err = u.uploadRepository.UploadFile(b, constants.S3_EXPORT_DO_DETAIL_PATH, request.FileName, request.FileType)
	if err != nil {
		fmt.Println("error upload", err)
		return helper.WriteLog(err, http.StatusInternalServerError, nil)
	}

	pusherData := &models.Pusher{
		Subject: "Export DO Detail",
		Link:    fmt.Sprintf("%s/%s.%s", constants.DELIVERY_ORDER_DETAIL_EXPORT_PATH, request.FileName, request.FileType),
		Type:    "export",
		UserId:  strconv.Itoa(request.UserID),
	}

	err = u.pusherRepository.Publish(pusherData)
	if err != nil {
		fmt.Println("pusher error = ", err.Error())
		return helper.WriteLog(err, http.StatusInternalServerError, nil)
	}

	return &model.ErrorLog{}
}
