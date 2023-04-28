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
	mongoRepositories "order-service/app/repositories/mongod"
	openSearchRepositories "order-service/app/repositories/open_search"
	"order-service/global/utils/helper"
	"order-service/global/utils/model"
	"os"
	"strconv"
	"time"
)

type SalesOrderOpenSearchUseCaseInterface interface {
	SyncToOpenSearchFromCreateEvent(salesOrder *models.SalesOrder, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromUpdateEvent(salesOrder *models.SalesOrder, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromDeleteEvent(salesOrder *models.SalesOrder, ctx context.Context) *model.ErrorLog
	SyncDetailToOpenSearchFromCreateEvent(salesOrderDetail *models.SalesOrderDetail, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncDetailToOpenSearchFromUpdateEvent(salesOrder *models.SalesOrder, salesOrderDetail *models.SalesOrderDetail, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncDetailToOpenSearchFromDeleteEvent(salesOrderDetail *models.SalesOrderDetail, ctx context.Context) *model.ErrorLog
	Get(request *models.SalesOrderExportRequest) *model.ErrorLog
	GetDetails(request *models.SalesOrderDetailExportRequest) *model.ErrorLog
}

type SalesOrderOpenSearchUseCase struct {
	salesOrderRepository                 repositories.SalesOrderRepositoryInterface
	salesOrderDetailRepository           repositories.SalesOrderDetailRepositoryInterface
	orderStatusRepository                repositories.OrderStatusRepositoryInterface
	productRepository                    repositories.ProductRepositoryInterface
	uomRepository                        repositories.UomRepositoryInterface
	categoryRepository                   repositories.CategoryRepositoryInterface
	salesOrderOpenSearchRepository       openSearchRepositories.SalesOrderOpenSearchRepositoryInterface
	deliveryOrderOpenSearchRepository    openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface
	salesOrderDetailOpenSearchRepository openSearchRepositories.SalesOrderDetailOpenSearchRepositoryInterface
	salesOrderJourneyRepository          mongoRepositories.SalesOrderJourneysRepositoryInterface
	uploadRepository                     repositories.UploadRepositoryInterface
	pusherRepository                     repositories.PusherRepositoryInterface
	ctx                                  context.Context
}

func InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, productRepository repositories.ProductRepositoryInterface, uomRepository repositories.UomRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface, salesOrderDetailOpenSearchRepository openSearchRepositories.SalesOrderDetailOpenSearchRepositoryInterface, categoryRepository repositories.CategoryRepositoryInterface, salesOrderJourneyRepository mongoRepositories.SalesOrderJourneysRepositoryInterface, uploadRepository repositories.UploadRepositoryInterface, ctx context.Context) SalesOrderOpenSearchUseCaseInterface {
	return &SalesOrderOpenSearchUseCase{
		salesOrderRepository:                 salesOrderRepository,
		salesOrderDetailRepository:           salesOrderDetailRepository,
		orderStatusRepository:                orderStatusRepository,
		productRepository:                    productRepository,
		uomRepository:                        uomRepository,
		salesOrderOpenSearchRepository:       salesOrderOpenSearchRepository,
		deliveryOrderOpenSearchRepository:    deliveryOrderOpenSearchRepository,
		salesOrderDetailOpenSearchRepository: salesOrderDetailOpenSearchRepository,
		categoryRepository:                   categoryRepository,
		salesOrderJourneyRepository:          salesOrderJourneyRepository,
		uploadRepository:                     uploadRepository,
		pusherRepository:                     repositories.InitPusherRepository(),
		ctx:                                  ctx,
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

		getLastCategoryResultChan := make(chan *models.CategoryChan)
		go u.categoryRepository.GetByID(getProductResult.Product.CategoryID, false, ctx, getLastCategoryResultChan)
		getLastCategoryResult := <-getLastCategoryResultChan

		if getLastCategoryResult.Error != nil && getLastCategoryResult.ErrorLog.StatusCode != http.StatusNotFound {
			errorLogData := helper.WriteLog(getLastCategoryResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		salesOrder.SalesOrderDetails[k].LastCategoryId = getProductResult.Product.CategoryID
		if getLastCategoryResult.Category != nil {
			salesOrder.SalesOrderDetails[k].LastCategoryName = &getLastCategoryResult.Category.Name
		}

		getFirstCategoryResultChan := make(chan *models.CategoryChan)
		go u.categoryRepository.GetByID(int(getLastCategoryResult.Category.ParentID.Int64), false, ctx, getFirstCategoryResultChan)
		getFirstCategoryResult := <-getFirstCategoryResultChan

		if getFirstCategoryResult.Error != nil && getFirstCategoryResult.ErrorLog.StatusCode != http.StatusNotFound {
			errorLogData := helper.WriteLog(getFirstCategoryResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		salesOrder.SalesOrderDetails[k].FirstCategoryId = int(getLastCategoryResult.Category.ParentID.Int64)
		if getFirstCategoryResult.Category != nil {
			salesOrder.SalesOrderDetails[k].FirstCategoryName = &getFirstCategoryResult.Category.Name
		}

		salesOrder.SalesOrderDetails[k].CreatedBy = salesOrder.CreatedBy
		salesOrder.SalesOrderDetails[k].LatestUpdatedBy = salesOrder.CreatedBy

		salesOrderDetailUpdateData := &models.SalesOrderDetail{
			UpdatedAt:       &now,
			IsDoneSyncToEs:  "1",
			EndDateSyncToEs: &now,
		}

		updateSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailChan)
		go u.salesOrderDetailRepository.UpdateByID(v.ID, salesOrderDetailUpdateData, false, "", sqlTransaction, ctx, updateSalesOrderDetailResultChan)
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
	salesOrder.LatestUpdatedBy = salesOrder.CreatedBy
	salesOrder.UpdatedAt = &now
	salesOrder.EndCreatedDate = &now

	createSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderOpenSearchRepository.Create(salesOrder, createSalesOrderResultChan)
	createSalesOrderResult := <-createSalesOrderResultChan

	if createSalesOrderResult.Error != nil {
		return createSalesOrderResult.ErrorLog
	}

	salesOrderUpdateData := &models.SalesOrder{
		LatestUpdatedBy: salesOrder.CreatedBy,
		UpdatedAt:       &now,
		IsDoneSyncToEs:  "1",
		EndCreatedDate:  &now,
		EndDateSyncToEs: &now,
	}

	updateSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderRepository.UpdateByID(salesOrder.ID, salesOrderUpdateData, false, "", sqlTransaction, ctx, updateSalesOrderResultChan)
	updateSalesOrderResult := <-updateSalesOrderResultChan

	if updateSalesOrderResult.Error != nil {
		errorLogData := helper.WriteLog(updateSalesOrderResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return &model.ErrorLog{}
}

func (u *SalesOrderOpenSearchUseCase) SyncToOpenSearchFromUpdateEvent(salesOrder *models.SalesOrder, ctx context.Context) *model.ErrorLog {
	now := time.Now()

	salesOrderRequest := &models.SalesOrderRequest{ID: salesOrder.ID}
	getSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderOpenSearchRepository.GetByID(salesOrderRequest, getSalesOrderResultChan)
	getSalesOrderResult := <-getSalesOrderResultChan

	if getSalesOrderResult.ErrorLog != nil {
		return getSalesOrderResult.ErrorLog
	}
	salesOrder.SalesOrderOpenSearchChanMap(getSalesOrderResult)

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

		v.IsDoneSyncToEs = "1"
		v.EndDateSyncToEs = &now
		salesOrder.SalesOrderDetails[k].OrderStatus = getOrderStatusDetailResult.OrderStatus

		getLastCategoryResultChan := make(chan *models.CategoryChan)
		go u.categoryRepository.GetByID(getProductResult.Product.CategoryID, false, ctx, getLastCategoryResultChan)
		getLastCategoryResult := <-getLastCategoryResultChan

		if getLastCategoryResult.Error != nil && getLastCategoryResult.ErrorLog.StatusCode != http.StatusNotFound {
			errorLogData := helper.WriteLog(getLastCategoryResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		salesOrder.SalesOrderDetails[k].LastCategoryId = getProductResult.Product.CategoryID
		if getLastCategoryResult.Category != nil {
			salesOrder.SalesOrderDetails[k].LastCategoryName = &getLastCategoryResult.Category.Name
		}

		getFirstCategoryResultChan := make(chan *models.CategoryChan)
		go u.categoryRepository.GetByID(int(getLastCategoryResult.Category.ParentID.Int64), false, ctx, getFirstCategoryResultChan)
		getFirstCategoryResult := <-getFirstCategoryResultChan

		if getFirstCategoryResult.Error != nil && getFirstCategoryResult.ErrorLog.StatusCode != http.StatusNotFound {
			errorLogData := helper.WriteLog(getFirstCategoryResult.Error, http.StatusInternalServerError, nil)
			return errorLogData
		}

		salesOrder.SalesOrderDetails[k].FirstCategoryId = int(getLastCategoryResult.Category.ParentID.Int64)
		if getFirstCategoryResult.Category != nil {
			salesOrder.SalesOrderDetails[k].FirstCategoryName = &getFirstCategoryResult.Category.Name
		}

		salesOrder.SalesOrderDetails[k].CreatedBy = salesOrder.CreatedBy
		salesOrder.SalesOrderDetails[k].LatestUpdatedBy = salesOrder.LatestUpdatedBy

		salesOrderDetail := &models.SalesOrderDetailOpenSearch{}
		salesOrderDetail.SalesOrderDetailOpenSearchMap(salesOrder, v)

		createSalesOrderDetailResultChan := make(chan *models.SalesOrderDetailOpenSearchChan)
		go u.salesOrderDetailOpenSearchRepository.Create(salesOrderDetail, createSalesOrderDetailResultChan)
		createSalesOrderDetailResult := <-createSalesOrderDetailResultChan

		if createSalesOrderDetailResult.Error != nil {
			return createSalesOrderDetailResult.ErrorLog
		}
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
	isParentDelete := salesOrder.SalesOrderDetails == nil
	salesOrderRequest := &models.SalesOrderRequest{ID: salesOrder.ID}

	if isParentDelete {
		salesOrder.DeletedAt = &now
		salesOrder.UpdatedAt = &now
	}

	getSalesOrderResultChan := make(chan *models.SalesOrderChan)
	go u.salesOrderOpenSearchRepository.GetByID(salesOrderRequest, getSalesOrderResultChan)
	getSalesOrderResult := <-getSalesOrderResultChan

	salesOrder.SalesOrderOpenSearchChanMap(getSalesOrderResult)

	for k := range salesOrder.SalesOrderDetails {
		salesOrder.SalesOrderDetails[k].DeletedAt = &now
		salesOrder.SalesOrderDetails[k].UpdatedAt = &now
		salesOrder.SalesOrderDetails[k].IsDoneSyncToEs = "1"
		salesOrder.SalesOrderDetails[k].EndDateSyncToEs = &now
	}

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
	go u.salesOrderDetailRepository.UpdateByID(salesOrderDetail.ID, salesOrderDetailUpdateData, false, "", sqlTransaction, ctx, updateSalesOrderDetailResultChan)
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

func (u *SalesOrderOpenSearchUseCase) Get(request *models.SalesOrderExportRequest) *model.ErrorLog {
	countSoRequest := &models.SalesOrderRequest{}
	countSoRequest.SalesOrderExportMap(request)

	getSalesOrdersCountResultChan := make(chan *models.SalesOrdersChan)
	go u.salesOrderOpenSearchRepository.Get(countSoRequest, true, getSalesOrdersCountResultChan)
	getSalesOrdersCountResult := <-getSalesOrdersCountResultChan

	if getSalesOrdersCountResult.Error != nil {
		fmt.Println("error = ", getSalesOrdersCountResult.Error)
		return getSalesOrdersCountResult.ErrorLog
	}

	soRequest := &models.SalesOrderRequest{}
	soRequest.SalesOrderExportMap(request)
	perPage, err := strconv.Atoi(os.Getenv("EXPORT_PARTIAL"))
	if err != nil {
		perPage = constants.EXPORT_PARTIAL_DEFAULT
	}
	soRequest.PerPage = perPage
	instalmentData := math.Ceil(float64(getSalesOrdersCountResult.Total) / float64(soRequest.PerPage))
	data := [][]interface{}{constants.SALES_ORDER_EXPORT_HEADER()}

	for i := 0; i < int(instalmentData); i++ {
		soRequest.Page = i + 1
		getSalesOrdersResultChan := make(chan *models.SalesOrdersChan)
		go u.salesOrderOpenSearchRepository.Get(soRequest, false, getSalesOrdersResultChan)
		getSalesOrdersResult := <-getSalesOrdersResultChan

		if getSalesOrdersResult.Error != nil {
			return getSalesOrdersResult.ErrorLog
		}
		for _, v := range getSalesOrdersResult.SalesOrders {
			getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
			go u.deliveryOrderOpenSearchRepository.GetBySalesOrderID(&models.DeliveryOrderRequest{SalesOrderID: v.ID}, getDeliveryOrdersResultChan)
			getDeliveryOrdersResult := <-getDeliveryOrdersResultChan

			getSalesOrdersJourneyChan := make(chan *models.SalesOrdersJourneysChan)
			go u.salesOrderJourneyRepository.GetBySoId(v.ID, false, u.ctx, getSalesOrdersJourneyChan)
			getSalesOrderJourneyResult := <-getSalesOrdersJourneyChan

			var journeys *models.SalesOrderJourneys = nil
			if getSalesOrderJourneyResult.SalesOrderJourneys != nil {
				journeys = getSalesOrderJourneyResult.SalesOrderJourneys[0]
			}

			data = append(data, v.MapToCsvRow(journeys, getDeliveryOrdersResult.DeliveryOrders))
		}
		progres := math.Round((float64((i*soRequest.PerPage)+len(getSalesOrdersResult.SalesOrders)) / float64(getSalesOrdersCountResult.Total)) * 100)
		fmt.Println(request.FileName, progres, "%")
	}

	b, err := helper.GenerateExportBufferFile(data, request.FileType)

	// Upload Files
	err = u.uploadRepository.UploadFile(b, constants.S3_EXPORT_SO_PATH, request.FileName, request.FileType)
	if err != nil {
		fmt.Println("error upload", err)
		return helper.WriteLog(err, http.StatusInternalServerError, nil)
	}

	pusherData := &models.Pusher{
		Subject: "Export SO Summary",
		Link:    fmt.Sprintf("%s/%s.%s", constants.SALES_ORDER_EXPORT_PATH, request.FileName, request.FileType),
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

func (u *SalesOrderOpenSearchUseCase) GetDetails(request *models.SalesOrderDetailExportRequest) *model.ErrorLog {
	countSoDetailRequest := &models.GetSalesOrderDetailRequest{}
	countSoDetailRequest.SalesOrderDetailExportMap(request)

	getSalesOrderDetailsCountResultChan := make(chan *models.SalesOrderDetailsOpenSearchChan)
	go u.salesOrderDetailOpenSearchRepository.Get(countSoDetailRequest, false, getSalesOrderDetailsCountResultChan)
	getSalesOrderDetailsCountResult := <-getSalesOrderDetailsCountResultChan

	if getSalesOrderDetailsCountResult.Error != nil {
		fmt.Println("error = ", getSalesOrderDetailsCountResult.Error)
		return getSalesOrderDetailsCountResult.ErrorLog
	}

	soDetailRequest := &models.GetSalesOrderDetailRequest{}
	soDetailRequest.SalesOrderDetailExportMap(request)
	perPage, err := strconv.Atoi(os.Getenv("EXPORT_PARTIAL"))
	if err != nil {
		perPage = constants.EXPORT_PARTIAL_DEFAULT
	}
	soDetailRequest.PerPage = perPage
	instalmentData := math.Ceil(float64(getSalesOrderDetailsCountResult.Total) / float64(soDetailRequest.PerPage))

	data := [][]interface{}{constants.SALES_ORDER_DETAIL_EXPORT_HEADER()}

	for i := 0; i < int(instalmentData); i++ {
		soDetailRequest.Page = i + 1
		getSalesOrderDetailsResultChan := make(chan *models.SalesOrderDetailsOpenSearchChan)
		go u.salesOrderDetailOpenSearchRepository.Get(soDetailRequest, false, getSalesOrderDetailsResultChan)
		getSalesDetailOrdersResult := <-getSalesOrderDetailsResultChan

		if getSalesDetailOrdersResult.Error != nil {
			return getSalesDetailOrdersResult.ErrorLog
		}

		for _, v := range getSalesDetailOrdersResult.SalesOrderDetails {
			getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
			go u.deliveryOrderOpenSearchRepository.GetBySalesOrderID(&models.DeliveryOrderRequest{SalesOrderID: v.SoID}, getDeliveryOrdersResultChan)
			getDeliveryOrdersResult := <-getDeliveryOrdersResultChan
			getSalesOrdersJourneyChan := make(chan *models.SalesOrdersJourneysChan)
			go u.salesOrderJourneyRepository.GetBySoId(v.ID, false, u.ctx, getSalesOrdersJourneyChan)
			getSalesOrderJourneyResult := <-getSalesOrdersJourneyChan
			var journeys *models.SalesOrderJourneys = nil
			if getSalesOrderJourneyResult.SalesOrderJourneys != nil {
				journeys = getSalesOrderJourneyResult.SalesOrderJourneys[0]
			}
			soRequest := &models.SalesOrderRequest{ID: v.SoID}
			getSalesOrderResultChan := make(chan *models.SalesOrderChan)
			go u.salesOrderOpenSearchRepository.GetByID(soRequest, getSalesOrderResultChan)
			getSalesOrdersResult := <-getSalesOrderResultChan
			if getSalesOrdersResult.Error != nil {
				return getSalesOrdersResult.ErrorLog
			}
			data = append(data, v.MapToCsvRow(journeys, getSalesOrdersResult.SalesOrder, getDeliveryOrdersResult.DeliveryOrders))
		}
		progres := math.Round((float64((i*soDetailRequest.PerPage)+len(getSalesDetailOrdersResult.SalesOrderDetails)) / float64(getSalesOrderDetailsCountResult.Total)) * 100)
		fmt.Println(request.FileName, progres, "%")
	}

	b, err := helper.GenerateExportBufferFile(data, request.FileType)

	// Upload Files
	err = u.uploadRepository.UploadFile(b, constants.S3_EXPORT_SO_DETAIL_PATH, request.FileName, request.FileType)
	if err != nil {
		fmt.Println("error upload", err)
		return helper.WriteLog(err, http.StatusInternalServerError, nil)
	}

	pusherData := &models.Pusher{
		Subject: "Export SO Detail",
		Link:    fmt.Sprintf("%s/%s.%s", constants.SALES_ORDER_DETAIL_EXPORT_PATH, request.FileName, request.FileType),
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
