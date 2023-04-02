package usecases

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/app/repositories"
	openSearchRepositories "order-service/app/repositories/open_search"
	"order-service/global/utils/helper"
	"order-service/global/utils/model"
	"strconv"
	"strings"
	"time"
)

type DeliveryOrderConsumerUseCaseInterface interface {
	SyncToOpenSearchFromCreateEvent(deliveryOrder *models.DeliveryOrder, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromUpdateEvent(deliveryOrder *models.DeliveryOrder, ctx context.Context) *model.ErrorLog
	SyncToOpenSearchFromDeleteEvent(deliveryOrderId *int, deliveryOrderDetailId []*int, ctx context.Context) *model.ErrorLog
	Get(request *models.DeliveryOrderExportRequest) (*models.DeliveryOrdersOpenSearchResponse, *model.ErrorLog)
}

type deliveryOrderConsumerUseCase struct {
	salesOrderRepository                    repositories.SalesOrderRepositoryInterface
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
	ctx                                     context.Context
}

func InitDeliveryOrderConsumerUseCaseInterface(salesOrderRepository repositories.SalesOrderRepositoryInterface, uploadRepository repositories.UploadRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, deliveryOrderOpenSearchRepository openSearchRepositories.DeliveryOrderOpenSearchRepositoryInterface, deliveryOrderDetailOpenSearchRepository openSearchRepositories.DeliveryOrderDetailOpenSearchRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, salesOrderOpenSearchUseCase SalesOrderOpenSearchUseCaseInterface, ctx context.Context) DeliveryOrderConsumerUseCaseInterface {
	return &deliveryOrderConsumerUseCase{
		salesOrderRepository:                    salesOrderRepository,
		uploadRepository:                        uploadRepository,
		orderStatusRepository:                   orderStatusRepository,
		brandRepository:                         brandRepository,
		productRepository:                       productRepository,
		agentRepository:                         agentRepository,
		storeRepository:                         storeRepository,
		deliveryOrderOpenSearchRepository:       deliveryOrderOpenSearchRepository,
		deliveryOrderDetailOpenSearchRepository: deliveryOrderDetailOpenSearchRepository,
		salesOrderOpenSearchRepository:          salesOrderOpenSearchRepository,
		SalesOrderOpenSearchUseCase:             salesOrderOpenSearchUseCase,
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

	return &model.ErrorLog{}
}

func (u *deliveryOrderConsumerUseCase) Get(request *models.DeliveryOrderExportRequest) (*models.DeliveryOrdersOpenSearchResponse, *model.ErrorLog) {
	doRequest := &models.DeliveryOrderRequest{}
	doRequest.DeliveryOrderExportMap(request)
	// getDeliveryOrdersCountResultChan := make(chan *models.DeliveryOrdersChan)
	// go u.deliveryOrderOpenSearchRepository.Get(doRequest, true, getDeliveryOrdersCountResultChan)
	// getDeliveryOrdersCountResult := <-getDeliveryOrdersCountResultChan

	// if getDeliveryOrdersCountResult.Error != nil {
	// 	fmt.Println("error = ", getDeliveryOrdersCountResult.Error)
	// 	return &models.DeliveryOrdersOpenSearchResponse{}, getDeliveryOrdersCountResult.ErrorLog
	// }
	// fmt.Println("cekk 2")
	// x := math.Ceil(float64(getDeliveryOrdersCountResult.Total / 50))
	// fmt.Println(x)

	getDeliveryOrdersResultChan := make(chan *models.DeliveryOrdersChan)
	go u.deliveryOrderOpenSearchRepository.Get(doRequest, false, getDeliveryOrdersResultChan)
	getDeliveryOrdersResult := <-getDeliveryOrdersResultChan

	if getDeliveryOrdersResult.Error != nil {
		return &models.DeliveryOrdersOpenSearchResponse{}, getDeliveryOrdersResult.ErrorLog
	}
	deliveryOrderResults := []*models.DeliveryOrderOpenSearchResponse{}
	deliveryOrdersCsv := []*models.DeliveryOrderCsvResponse{}

	b := new(bytes.Buffer)
	writer := csv.NewWriter(b)
	defer writer.Flush()

	header := []string{
		"do_status",
		"do_date",
		"sj_no",
		"do_no",
		"order_no",
		"so_date",
		"so_no",
		"so_source",
		"agent_id",
		"agent_name",
		"gudang_id",
		"gudang_name",
		"brand_id",
		"brand_name",
		"kode_salesman",
		"salesman",
		"kategory_toko",
		"kode_toko_dbo",
		"kode_toko",
		"nama_toko",
		"kode_kecamatan",
		"kecamatan",
		"kode_city",
		"city",
		"kode_province",
		"province",
		"do_amount",
		"nama_supir",
		"plat_no",
		"catatan",
		"created_date",
		"updated_date",
		"user_id_created",
		"user_id_modified"}

	if err := writer.Write(header); err != nil {
		fmt.Println("error writer", err)
	}
	for _, v := range getDeliveryOrdersResult.DeliveryOrders {
		deliveryOrder := models.DeliveryOrderOpenSearchResponse{}
		deliveryOrder.DeliveryOrderOpenSearchResponseMap(v)

		deliveryOrderResults = append(deliveryOrderResults, &deliveryOrder)
		deliveryOrderCsv := &models.DeliveryOrderCsvResponse{}
		deliveryOrderCsv.DoDetailMap(v)
		deliveryOrdersCsv = append(deliveryOrdersCsv, deliveryOrderCsv)

		var csvRow []string
		csvRow = append(csvRow,
			strconv.Itoa(deliveryOrderCsv.DoStatus),
			deliveryOrderCsv.DoDate,
			deliveryOrderCsv.SjNo.String,
			deliveryOrderCsv.DoNo,
			deliveryOrderCsv.OrderNo,
			deliveryOrderCsv.SoDate,
			deliveryOrderCsv.SoNo,
			strconv.Itoa(deliveryOrderCsv.SoSource),
			strconv.Itoa(deliveryOrderCsv.AgentID),
			deliveryOrderCsv.AgentName,
			strconv.Itoa(deliveryOrderCsv.GudangID),
			deliveryOrderCsv.GudangName,
			strconv.Itoa(deliveryOrderCsv.BrandID),
			deliveryOrderCsv.BrandName,
			strconv.Itoa(int(deliveryOrderCsv.KodeSalesman.Int64)),
			deliveryOrderCsv.Salesman.String,
			deliveryOrderCsv.KategoryToko.String,
			deliveryOrderCsv.KodeTokoDbo.String,
			deliveryOrderCsv.KodeToko.String,
			deliveryOrderCsv.NamaToko.String,
			strconv.Itoa(deliveryOrderCsv.KodeKecamatan),
			deliveryOrderCsv.Kecamatan.String,
			strconv.Itoa(deliveryOrderCsv.KodeCity),
			deliveryOrderCsv.City.String,
			strconv.Itoa(deliveryOrderCsv.KodeProvince),
			deliveryOrderCsv.Province.String,
			strconv.FormatFloat(deliveryOrderCsv.DoAmount, 'f', 6, 64),
			deliveryOrderCsv.NamaSupir.String,
			deliveryOrderCsv.PlatNo.String,
			deliveryOrderCsv.Catatan.String,
			deliveryOrderCsv.CreatedDate.String(),
			deliveryOrderCsv.UpdatedDate.String(),
			strconv.Itoa(deliveryOrderCsv.UserIDCreated),
			strconv.Itoa(deliveryOrderCsv.UserIDModified))

		if err := writer.Write(csvRow); err != nil {
			fmt.Println("error fill", err)
		}

	}

	// Upload Files
	err := u.uploadRepository.UploadFile(b, request.FileDate, request.FileType)
	if err != nil {
		fmt.Println("error upload", err)
	}

	deliveryOrders := &models.DeliveryOrdersOpenSearchResponse{
		DeliveryOrders: deliveryOrderResults,
		Total:          getDeliveryOrdersResult.Total,
	}

	return deliveryOrders, &model.ErrorLog{}
}
