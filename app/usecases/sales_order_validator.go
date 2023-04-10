package usecases

import (
	"context"
	"fmt"
	"net/http"
	"order-service/app/middlewares"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/app/repositories"
	"order-service/global/utils/helper"
	"order-service/global/utils/model"
	"strconv"
	"strings"
	"time"

	"github.com/bxcodec/dbresolver"
	"github.com/gin-gonic/gin"
)

type SalesOrderValidatorInterface interface {
	CreateSalesOrderValidator(insertRequest *models.SalesOrderStoreRequest, ctx *gin.Context) error
	GetSalesOrderValidator(*gin.Context) (*models.SalesOrderRequest, error)
	GetSalesOrderDetailValidator(*gin.Context) (*models.GetSalesOrderDetailRequest, error)
	GetSalesOrderSyncToKafkaHistoriesValidator(*gin.Context) (*models.SalesOrderEventLogRequest, error)
	GetSalesOrderJourneysValidator(*gin.Context) (*models.SalesOrderJourneyRequest, error)
	GetSOUploadHistoriesValidator(*gin.Context) (*models.GetSoUploadHistoriesRequest, error)
	GetSosjUploadHistoriesValidator(*gin.Context) (*models.GetSosjUploadHistoriesRequest, error)
	ExportSalesOrderValidator(ctx *gin.Context) (*models.SalesOrderExportRequest, error)
	ExportSalesOrderDetailValidator(ctx *gin.Context) (*models.SalesOrderDetailExportRequest, error)
	UpdateSalesOrderByIdValidator(updateRequest *models.SalesOrderUpdateRequest, ctx *gin.Context) error
	UpdateSalesOrderDetailBySoIdValidator(updateRequest *models.SalesOrderUpdateRequest, ctx *gin.Context) error
	UpdateSalesOrderDetailByIdValidator(updateRequest *models.UpdateSalesOrderDetailByIdRequest, ctx *gin.Context) error
	DeleteSalesOrderByIdValidator(string, *gin.Context) (int, error)
	DeleteSalesOrderDetailByIdValidator(string, *gin.Context) (int, error)
	DeleteSalesOrderDetailBySoIdValidator(string, *gin.Context) (int, error)
}

type SalesOrderValidator struct {
	requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface
	orderSourceRepository       repositories.OrderSourceRepositoryInterface
	salesmanRepository          repositories.SalesmanRepositoryInterface
	orderStatusRepository       repositories.OrderStatusRepositoryInterface
	salesOrderRepository        repositories.SalesOrderRepositoryInterface
	salesOrderDetailRepository  repositories.SalesOrderDetailRepositoryInterface
	deliveryOrderRepository     repositories.DeliveryOrderRepositoryInterface
	db                          dbresolver.DB
	ctx                         context.Context
}

func InitSalesOrderValidator(requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, orderSourceRepository repositories.OrderSourceRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, db dbresolver.DB, ctx context.Context) SalesOrderValidatorInterface {
	return &SalesOrderValidator{
		requestValidationMiddleware: requestValidationMiddleware,
		orderSourceRepository:       orderSourceRepository,
		salesmanRepository:          salesmanRepository,
		orderStatusRepository:       orderStatusRepository,
		salesOrderRepository:        salesOrderRepository,
		salesOrderDetailRepository:  salesOrderDetailRepository,
		deliveryOrderRepository:     deliveryOrderRepository,
		db:                          db,
		ctx:                         ctx,
	}
}

func (c *SalesOrderValidator) CreateSalesOrderValidator(insertRequest *models.SalesOrderStoreRequest, ctx *gin.Context) error {
	dateField := []*models.DateInputRequest{
		{
			Field: "so_date",
			Value: insertRequest.SoDate,
		},
		{
			Field: "so_ref_date",
			Value: insertRequest.SoRefDate,
		},
	}
	err := c.requestValidationMiddleware.DateInputValidation(ctx, dateField, constants.ERROR_ACTION_NAME_CREATE)
	if err != nil {
		return err
	}

	uniqueField := []*models.UniqueRequest{}
	if len(insertRequest.SoRefCode) > 0 {
		err = c.requestValidationMiddleware.OrderSourceValidation(ctx, insertRequest.OrderSourceID, insertRequest.SoRefCode, constants.ERROR_ACTION_NAME_CREATE)
		if err != nil {
			return err
		}

		uniqueField = append(uniqueField, &models.UniqueRequest{
			Table: constants.SALES_ORDERS_TABLE,
			Field: "so_ref_code",
			Value: insertRequest.SoRefCode,
		})
	}

	if len(uniqueField) > 0 {
		err = c.requestValidationMiddleware.UniqueValidation(ctx, uniqueField)
		if err != nil {
			return err
		}
	}

	brandIds := []int{}
	mustActiveField := []*models.MustActiveRequest{
		helper.GenerateMustActive("agents", "agent_id", insertRequest.AgentID, "active"),
		helper.GenerateMustActive("stores", "store_id", insertRequest.StoreID, "active"),
		helper.GenerateMustActive("users", "user_id", insertRequest.UserID, "ACTIVE"),
	}
	for i, v := range insertRequest.SalesOrderDetails {
		mustActiveField = append(mustActiveField, &models.MustActiveRequest{
			Table:    "products",
			ReqField: fmt.Sprintf("sales_order_details[%d].product_id", i),
			Clause:   fmt.Sprintf("id = %d AND isActive = %d", v.ProductID, 1),
		})
		mustActiveField = append(mustActiveField, &models.MustActiveRequest{
			Table:    "uoms",
			ReqField: fmt.Sprintf("sales_order_details[%d].uom_id", i),
			Clause:   fmt.Sprintf(constants.CLAUSE_ID_VALIDATION, v.UomID),
		})
		mustActiveField = append(mustActiveField, &models.MustActiveRequest{
			Table:    "brands",
			ReqField: fmt.Sprintf("sales_order_details[%d].brand_id", i),
			Clause:   fmt.Sprintf("id = %d AND status_active = %d", v.BrandID, 1),
		})

		brandIds = append(brandIds, v.BrandID)
	}
	err = c.requestValidationMiddleware.MustActiveValidation(ctx, mustActiveField)
	if err != nil {
		return err
	}

	// Get Order Source Status By Id
	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go c.orderSourceRepository.GetByID(insertRequest.OrderSourceID, false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getOrderSourceResult.Error != nil {
		ctx.JSON(getOrderSourceResult.ErrorLog.StatusCode, helper.GenerateResultByErrorLog(getOrderSourceResult.ErrorLog))
		return getOrderSourceResult.Error
	}
	sourceName := getOrderSourceResult.OrderSource.SourceName
	if sourceName != "manager" && len(insertRequest.DeviceId) < 1 {
		errorLog := helper.NewWriteLog(model.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage("create", "device_id tidak boleh kosong")},
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		helper.GenerateResultByErrorLog(errorLog)
		ctx.JSON(http.StatusUnprocessableEntity, helper.GenerateResultByErrorLog(errorLog))

		err = helper.NewError("device_id tidak boleh kosong")
		return err
	}

	if sourceName == "store" && len(insertRequest.ReferralCode) > 0 {
		// Get Salesmans By Agent Id
		getSalesmanResultChan := make(chan *models.SalesmansChan)
		go c.salesmanRepository.GetByAgentId(insertRequest.AgentID, false, ctx, getSalesmanResultChan)
		getSalesmanResult := <-getSalesmanResultChan

		if getSalesmanResult.Error != nil {
			ctx.JSON(getSalesmanResult.ErrorLog.StatusCode, helper.GenerateResultByErrorLog(getOrderSourceResult.ErrorLog))
			return getSalesmanResult.Error
		}

		isExist := false
		for _, v := range getSalesmanResult.Salesmans {
			if v.ReferralCode == insertRequest.ReferralCode {
				isExist = true
				break
			}
		}

		if !isExist {
			errorLog := helper.NewWriteLog(model.ErrorLog{
				Message:       []string{helper.GenerateUnprocessableErrorMessage("create", "referral code tidak terdaftar")},
				SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
				StatusCode:    http.StatusUnprocessableEntity,
			})
			ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))

			err = helper.NewError("referral code tidak terdaftar")
			return err
		}
	}

	now := time.Now().UTC().Add(7 * time.Hour)
	parseSoDate, _ := time.Parse(constants.DATE_FORMAT_COMMON, insertRequest.SoDate)
	parseSoRefDate, _ := time.Parse(constants.DATE_FORMAT_COMMON, insertRequest.SoRefDate)
	duration := time.Hour*time.Duration(now.Hour()) + time.Minute*time.Duration(now.Minute()) + time.Second*time.Duration(now.Second()) + time.Nanosecond*time.Duration(now.Nanosecond())

	soDate := parseSoDate.UTC().Add(duration)
	soRefDate := parseSoRefDate.UTC().Add(duration)
	nowUTC := now.UTC()
	if now.UTC().Hour() < soRefDate.Hour() {
		nowUTC = nowUTC.Add(7 * time.Hour)
	}

	if sourceName == "manager" && !(soDate.Add(1*time.Minute).After(soRefDate) && soRefDate.Add(-1*time.Minute).Before(nowUTC) && soDate.Add(-1*time.Minute).Before(nowUTC) && soRefDate.Month() == nowUTC.Month() && soRefDate.UTC().Year() == nowUTC.Year()) {

		err = helper.NewError("so_date dan so_ref_date harus sama dengan kurang dari hari ini dan harus di bulan berjalan")
		errorLog := helper.NewWriteLog(model.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage("create", "so_date dan so_ref_date harus sama dengan kurang dari hari ini dan harus di bulan berjalan")},
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))

		return err

	} else if (sourceName == "salesman" || sourceName == "store") && !(soDate.Equal(now.Local()) && soRefDate.Equal(now.Local())) {

		err = helper.NewError("so_date dan so_ref_date harus sama dengan hari ini")
		errorLog := helper.NewWriteLog(model.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage("create", "so_date dan so_ref_date harus sama dengan hari ini")},
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))

		return err

	}

	err = c.requestValidationMiddleware.AgentIdValidation(ctx, insertRequest.AgentID, insertRequest.UserID, constants.ERROR_ACTION_NAME_CREATE)
	if err != nil {
		return err
	}

	err = c.requestValidationMiddleware.StoreIdValidation(ctx, insertRequest.StoreID, insertRequest.AgentID, constants.ERROR_ACTION_NAME_CREATE)
	if err != nil {
		return err
	}

	if insertRequest.SalesmanID > 0 {
		err = c.requestValidationMiddleware.SalesmanIdValidation(ctx, insertRequest.SalesmanID, insertRequest.AgentID, constants.ERROR_ACTION_NAME_CREATE)
		if err != nil {
			return err
		}
	}

	err = c.requestValidationMiddleware.BrandIdValidation(ctx, brandIds, insertRequest.AgentID, constants.ERROR_ACTION_NAME_CREATE)
	if err != nil {
		return err
	}

	return nil
}

func (c *SalesOrderValidator) GetSalesOrderValidator(ctx *gin.Context) (*models.SalesOrderRequest, error) {
	pageInt, err := c.getIntQueryWithDefault("page", "1", true, ctx)
	if err != nil {
		return nil, err
	}

	perPageInt, err := c.getIntQueryWithDefault("per_page", "10", true, ctx)
	if err != nil {
		return nil, err
	}

	sortField := c.getQueryWithDefault("sort_field", "created_at", ctx)

	if sortField != "order_status_id" && sortField != "so_date" && sortField != "so_ref_code" && sortField != "so_code" && sortField != "store_code" && sortField != "store_name" && sortField != "created_at" && sortField != "updated_at" {
		err = helper.NewError("Parameter 'sort_field' harus bernilai 'order_status_id' or 'so_date' or 'so_ref_code' or 'so_code' or 'store_code' or 'store_name' or 'created_at' or 'updated_at' ")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return nil, err
	}

	intAgentID, err := c.getIntQueryWithDefault("agent_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intStoreID, err := c.getIntQueryWithDefault("store_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intBrandID, err := c.getIntQueryWithDefault("brand_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intOrderSourceID, err := c.getIntQueryWithDefault("order_source_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intOrderStatusID, err := c.getIntQueryWithDefault("order_status_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intID, err := c.getIntQueryWithDefault("id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intProductID, err := c.getIntQueryWithDefault("product_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intCategoryID, err := c.getIntQueryWithDefault("category_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intSalesmanID, err := c.getIntQueryWithDefault("salesman_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intProvinceID, err := c.getIntQueryWithDefault("province_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intCityID, err := c.getIntQueryWithDefault("city_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intDistrictID, err := c.getIntQueryWithDefault("district_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intVillageID, err := c.getIntQueryWithDefault("village_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	dateFields := []*models.DateInputRequest{}

	startSoDate, dateFields := c.getQueryWithDateValidation("start_so_date", "", dateFields, ctx)

	endSoDate, dateFields := c.getQueryWithDateValidation("end_so_date", "", dateFields, ctx)

	startCreatedAt, dateFields := c.getQueryWithDateValidation("start_created_at", "", dateFields, ctx)

	endCreatedAt, dateFields := c.getQueryWithDateValidation("end_created_at", "", dateFields, ctx)

	err = c.requestValidationMiddleware.DateInputValidation(ctx, dateFields, constants.ERROR_ACTION_NAME_GET)
	if err != nil {
		return nil, err
	}

	salesOrderRequest := &models.SalesOrderRequest{
		Page:              pageInt,
		PerPage:           perPageInt,
		SortField:         sortField,
		SortValue:         c.getQueryWithDefault("sort_value", "desc", ctx),
		GlobalSearchValue: c.getQueryWithDefault("global_search_value", "", ctx),
		AgentID:           intAgentID,
		StoreID:           intStoreID,
		BrandID:           intBrandID,
		OrderSourceID:     intOrderSourceID,
		OrderStatusID:     intOrderStatusID,
		StartSoDate:       startSoDate,
		EndSoDate:         endSoDate,
		ID:                intID,
		ProductID:         intProductID,
		CategoryID:        intCategoryID,
		SalesmanID:        intSalesmanID,
		ProvinceID:        intProvinceID,
		CityID:            intCityID,
		DistrictID:        intDistrictID,
		VillageID:         intVillageID,
		StartCreatedAt:    startCreatedAt,
		EndCreatedAt:      endCreatedAt,
	}
	return salesOrderRequest, nil
}

func (c *SalesOrderValidator) GetSalesOrderDetailValidator(ctx *gin.Context) (*models.GetSalesOrderDetailRequest, error) {
	pageInt, err := c.getIntQueryWithDefault("page", "1", true, ctx)
	if err != nil {
		return nil, err
	}

	perPageInt, err := c.getIntQueryWithDefault("per_page", "10", true, ctx)
	if err != nil {
		return nil, err
	}

	sortField := c.getQueryWithDefault("sort_field", "created_at", ctx)

	if sortField != "order_status_id" && sortField != "so_date" && sortField != "so_ref_code" && sortField != "so_code" && sortField != "store_code" && sortField != "store_name" && sortField != "created_at" && sortField != "updated_at" {
		err = helper.NewError("Parameter 'sort_field' harus bernilai 'order_status_id' or 'so_date' or 'so_ref_code' or 'so_code' or 'store_code' or 'store_name' or 'created_at' or 'updated_at' ")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return nil, err
	}

	intAgentID, err := c.getIntQueryWithDefault("agent_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intStoreID, err := c.getIntQueryWithDefault("store_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intBrandID, err := c.getIntQueryWithDefault("brand_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intOrderSourceID, err := c.getIntQueryWithDefault("order_source_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intOrderStatusID, err := c.getIntQueryWithDefault("order_status_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intID, err := c.getIntQueryWithDefault("id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intProductID, err := c.getIntQueryWithDefault("product_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intCategoryID, err := c.getIntQueryWithDefault("category_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intSalesmanID, err := c.getIntQueryWithDefault("salesman_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intProvinceID, err := c.getIntQueryWithDefault("province_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intCityID, err := c.getIntQueryWithDefault("city_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intDistrictID, err := c.getIntQueryWithDefault("district_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intVillageID, err := c.getIntQueryWithDefault("village_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	dateFields := []*models.DateInputRequest{}

	startSoDate, dateFields := c.getQueryWithDateValidation("start_so_date", "", dateFields, ctx)

	endSoDate, dateFields := c.getQueryWithDateValidation("end_so_date", "", dateFields, ctx)

	startCreatedAt, dateFields := c.getQueryWithDateValidation("start_created_at", "", dateFields, ctx)

	endCreatedAt, dateFields := c.getQueryWithDateValidation("end_created_at", "", dateFields, ctx)

	err = c.requestValidationMiddleware.DateInputValidation(ctx, dateFields, constants.ERROR_ACTION_NAME_GET)
	if err != nil {
		return nil, err
	}

	salesOrderRequest := &models.GetSalesOrderDetailRequest{
		Page:              pageInt,
		PerPage:           perPageInt,
		SortField:         sortField,
		SortValue:         c.getQueryWithDefault("sort_value", "desc", ctx),
		GlobalSearchValue: c.getQueryWithDefault("global_search_value", "", ctx),
		AgentID:           intAgentID,
		StoreID:           intStoreID,
		BrandID:           intBrandID,
		OrderSourceID:     intOrderSourceID,
		OrderStatusID:     intOrderStatusID,
		StartSoDate:       startSoDate,
		EndSoDate:         endSoDate,
		ID:                intID,
		ProductID:         intProductID,
		CategoryID:        intCategoryID,
		SalesmanID:        intSalesmanID,
		ProvinceID:        intProvinceID,
		CityID:            intCityID,
		DistrictID:        intDistrictID,
		VillageID:         intVillageID,
		StartCreatedAt:    startCreatedAt,
		EndCreatedAt:      endCreatedAt,
	}
	return salesOrderRequest, nil
}

func (c *SalesOrderValidator) GetSalesOrderSyncToKafkaHistoriesValidator(ctx *gin.Context) (*models.SalesOrderEventLogRequest, error) {
	pageInt, err := c.getIntQueryWithDefault("page", "1", true, ctx)
	if err != nil {
		return nil, err
	}

	perPageInt, err := c.getIntQueryWithDefault("per_page", "10", true, ctx)
	if err != nil {
		return nil, err
	}

	sortField := c.getQueryWithDefault("sort_field", "created_at", ctx)

	if sortField != "so_code" && sortField != "status" && sortField != "agent_name" && sortField != "created_at" {
		err = helper.NewError("Parameter 'sort_field' harus bernilai 'so_code' or 'status' or 'agent_name' or 'created_at' ")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return nil, err
	}

	intAgentID, err := c.getIntQueryWithDefault("agent_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	salesOrderRequest := &models.SalesOrderEventLogRequest{
		Page:              pageInt,
		PerPage:           perPageInt,
		SortField:         sortField,
		SortValue:         strings.ToLower(c.getQueryWithDefault("sort_value", "desc", ctx)),
		GlobalSearchValue: strings.ToLower(c.getQueryWithDefault("global_search_value", "", ctx)),
		RequestID:         c.getQueryWithDefault("request_id", "", ctx),
		SoCode:            c.getQueryWithDefault("so_code", "", ctx),
		Status:            strings.ToLower(c.getQueryWithDefault("status", "", ctx)),
		Action:            c.getQueryWithDefault("action", "", ctx),
		AgentID:           intAgentID,
	}
	return salesOrderRequest, nil
}

func (c *SalesOrderValidator) GetSalesOrderJourneysValidator(ctx *gin.Context) (*models.SalesOrderJourneyRequest, error) {
	pageInt, err := c.getIntQueryWithDefault("page", "1", true, ctx)
	if err != nil {
		return nil, err
	}

	perPageInt, err := c.getIntQueryWithDefault("per_page", "10", true, ctx)
	if err != nil {
		return nil, err
	}

	intSoID, err := c.getIntQueryWithDefault("so_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	sortField := c.getQueryWithDefault("sort_field", "created_at", ctx)

	if sortField != "so_code" && sortField != "status" && sortField != "created_at" && sortField != "action" {
		err = helper.NewError("Parameter 'sort_field' harus bernilai 'so_code' or 'status' or 'created_at' or 'action' ")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return nil, err
	}

	salesOrderJourneysRequest := models.SalesOrderJourneyRequest{
		Page:              pageInt,
		PerPage:           perPageInt,
		SortField:         sortField,
		SortValue:         strings.ToLower(c.getQueryWithDefault("sort_value", "desc", ctx)),
		GlobalSearchValue: c.getQueryWithDefault("global_search_value", "", ctx),
		SoId:              intSoID,
		SoCode:            c.getQueryWithDefault("so_code", "", ctx),
		Status:            c.getQueryWithDefault("status", "", ctx),
		Action:            c.getQueryWithDefault("action", "", ctx),
		StartDate:         c.getQueryWithDefault("start_date", "", ctx),
		EndDate:           c.getQueryWithDefault("end_date", "", ctx),
	}

	return &salesOrderJourneysRequest, nil
}

func (c *SalesOrderValidator) GetSOUploadHistoriesValidator(ctx *gin.Context) (*models.GetSoUploadHistoriesRequest, error) {
	pageInt, err := c.getIntQueryWithDefault("page", "1", true, ctx)
	if err != nil {
		return nil, err
	}

	perPageInt, err := c.getIntQueryWithDefault("per_page", "10", true, ctx)
	if err != nil {
		return nil, err
	}

	sortField := c.getQueryWithDefault("sort_field", "created_at", ctx)

	if sortField != "agent_name" && sortField != "file_name" && sortField != "status" && sortField != "created_at" {
		err = helper.NewError("Parameter 'sort_field' harus bernilai 'agent_name' or 'file_name' or 'status' or 'created_at' ")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return nil, err
	}

	intAgentID, err := c.getIntQueryWithDefault("agent_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intUploadedBy, err := c.getIntQueryWithDefault("uploaded_by", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	dateFields := []*models.DateInputRequest{}

	startUploadAt, dateFields := c.getQueryWithDateValidation("start_upload_at", "", dateFields, ctx)

	endUploadAt, dateFields := c.getQueryWithDateValidation("end_upload_at", "", dateFields, ctx)

	finishProcessDateStart, dateFields := c.getQueryWithDateValidation("finish_process_date_start", "", dateFields, ctx)

	finishProcessDateEnd, dateFields := c.getQueryWithDateValidation("finish_process_date_end", "", dateFields, ctx)

	err = c.requestValidationMiddleware.DateInputValidation(ctx, dateFields, constants.ERROR_ACTION_NAME_GET)
	if err != nil {
		return nil, err
	}

	getSoUploadHistoriesRequest := &models.GetSoUploadHistoriesRequest{
		ID:                     c.getQueryWithDefault("id", "", ctx),
		Page:                   pageInt,
		PerPage:                perPageInt,
		SortField:              sortField,
		SortValue:              strings.ToLower(c.getQueryWithDefault("sort_value", "desc", ctx)),
		GlobalSearchValue:      strings.ToLower(c.getQueryWithDefault("global_search_value", "", ctx)),
		RequestID:              c.getQueryWithDefault("request_id", "", ctx),
		FileName:               c.getQueryWithDefault("file_name", "", ctx),
		BulkCode:               c.getQueryWithDefault("bulk_code", "", ctx),
		AgentID:                intAgentID,
		Status:                 strings.ToLower(c.getQueryWithDefault("status", "", ctx)),
		UploadedBy:             intUploadedBy,
		StartUploadAt:          startUploadAt,
		EndUploadAt:            endUploadAt,
		FinishProcessDateStart: finishProcessDateStart,
		FinishProcessDateEnd:   finishProcessDateEnd,
	}
	return getSoUploadHistoriesRequest, nil
}

func (c *SalesOrderValidator) GetSosjUploadHistoriesValidator(ctx *gin.Context) (*models.GetSosjUploadHistoriesRequest, error) {
	pageInt, err := c.getIntQueryWithDefault("page", "1", true, ctx)
	if err != nil {
		return nil, err
	}

	perPageInt, err := c.getIntQueryWithDefault("per_page", "10", true, ctx)
	if err != nil {
		return nil, err
	}

	sortField := c.getQueryWithDefault("sort_field", "created_at", ctx)

	if sortField != "agent_name" && sortField != "file_name" && sortField != "status" && sortField != "uploaded_by_name" && sortField != "created_at" {
		err = helper.NewError("Parameter 'sort_field' harus bernilai 'agent_name' or 'file_name' or 'status' or `uploaded_by_name` or 'created_at' ")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return nil, err
	}

	intAgentID, err := c.getIntQueryWithDefault("agent_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intUploadedBy, err := c.getIntQueryWithDefault("uploaded_by", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	dateFields := []*models.DateInputRequest{}

	createdAt, dateFields := c.getQueryWithDateValidation("created_at", "", dateFields, ctx)

	startUploadAt, dateFields := c.getQueryWithDateValidation("start_upload_at", "", dateFields, ctx)

	endUploadAt, dateFields := c.getQueryWithDateValidation("end_upload_at", "", dateFields, ctx)

	finishProcessDateStart, dateFields := c.getQueryWithDateValidation("finish_process_date_start", "", dateFields, ctx)

	finishProcessDateEnd, dateFields := c.getQueryWithDateValidation("finish_process_date_end", "", dateFields, ctx)

	err = c.requestValidationMiddleware.DateInputValidation(ctx, dateFields, constants.ERROR_ACTION_NAME_GET)
	if err != nil {
		return nil, err
	}

	getSosjUploadHistoriesRequest := &models.GetSosjUploadHistoriesRequest{
		ID:                     c.getQueryWithDefault("id", "", ctx),
		Page:                   pageInt,
		PerPage:                perPageInt,
		SortField:              sortField,
		SortValue:              strings.ToLower(c.getQueryWithDefault("sort_value", "desc", ctx)),
		GlobalSearchValue:      strings.ToLower(c.getQueryWithDefault("global_search_value", "", ctx)),
		RequestID:              c.getQueryWithDefault("request_id", "", ctx),
		FileName:               c.getQueryWithDefault("file_name", "", ctx),
		BulkCode:               c.getQueryWithDefault("bulk_code", "", ctx),
		AgentID:                intAgentID,
		Status:                 strings.ToLower(c.getQueryWithDefault("status", "", ctx)),
		UploadedBy:             intUploadedBy,
		UploadedByName:         c.getQueryWithDefault("uploaded_by_name", "", ctx),
		UploadedByEmail:        c.getQueryWithDefault("uploaded_by_email", "", ctx),
		CreatedAt:              createdAt,
		StartUploadAt:          startUploadAt,
		EndUploadAt:            endUploadAt,
		FinishProcessDateStart: finishProcessDateStart,
		FinishProcessDateEnd:   finishProcessDateEnd,
	}
	return getSosjUploadHistoriesRequest, nil
}

func (d *SalesOrderValidator) getQueryWithDefault(param string, empty string, ctx *gin.Context) string {
	result, isStartCreatedAt := ctx.GetQuery(param)
	if isStartCreatedAt == false {
		result = empty
	}
	return result
}

func (d *SalesOrderValidator) getQueryWithDateValidation(param string, empty string, dateFields []*models.DateInputRequest, ctx *gin.Context) (string, []*models.DateInputRequest) {
	result, isResultExist := ctx.GetQuery(param)
	if isResultExist == false {
		result = empty
	} else {
		dateFields = append(dateFields, &models.DateInputRequest{
			Field: param,
			Value: result,
		})
	}
	return result, dateFields
}

func (d *SalesOrderValidator) getIntQueryWithDefault(param string, empty string, isNotZero bool, ctx *gin.Context) (int, error) {
	sResult := d.getQueryWithDefault(param, empty, ctx)
	result, err := strconv.Atoi(sResult)
	if err != nil {
		err = helper.NewError(fmt.Sprintf("Parameter '%s' harus bernilai integer", param))
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return 0, err
	}
	if result == 0 && isNotZero {
		err = helper.NewError(fmt.Sprintf("Parameter '%s' harus bernilai integer > 0", param))
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return 0, err
	}
	return result, nil
}

func (d *SalesOrderValidator) getIntQueryWithMustActive(param string, empty string, isNotZero bool, table string, clause string, ctx *gin.Context) (int, *models.MustActiveRequest, error) {
	sResult := d.getQueryWithDefault(param, empty, ctx)
	result, err := strconv.Atoi(sResult)
	if err != nil {
		err = helper.NewError(fmt.Sprintf("Parameter '%s' harus bernilai integer", param))
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return 0, nil, err
	}
	if result == 0 && isNotZero {
		err = helper.NewError(fmt.Sprintf("Parameter '%s' harus bernilai integer > 0", param))
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return 0, nil, err
	}
	mustActiveField := &models.MustActiveRequest{Table: table, ReqField: param, Clause: fmt.Sprintf(clause, result)}

	return result, mustActiveField, nil
}

func (d *SalesOrderValidator) ExportSalesOrderValidator(ctx *gin.Context) (*models.SalesOrderExportRequest, error) {
	sortField := d.getQueryWithDefault("sort_field", "created_at", ctx)

	if sortField != "order_status_id" && sortField != "do_date" && sortField != "do_ref_code" && sortField != "store_id" && sortField != "created_at" && sortField != "updated_at" {
		err := helper.NewError("Parameter 'sort_field' harus bernilai 'order_status_id' or 'do_date' or 'do_ref_code' or 'created_at' or 'updated_at'")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return nil, err
	}

	intID, err := d.getIntQueryWithDefault("id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	mustActiveFields := []*models.MustActiveRequest{}

	intAgentID, m, err := d.getIntQueryWithMustActive("agent_id", "0", false, "agents", constants.CLAUSE_ID_VALIDATION, ctx)
	if err != nil {
		return nil, err
	}
	if intAgentID > 0 {
		mustActiveFields = append(mustActiveFields, m)
	}

	intStoreID, m, err := d.getIntQueryWithMustActive("store_id", "0", false, "stores", constants.CLAUSE_ID_VALIDATION, ctx)
	if err != nil {
		return nil, err
	}
	if intStoreID > 0 {
		mustActiveFields = append(mustActiveFields, m)
	}

	intBrandID, m, err := d.getIntQueryWithMustActive("brand_id", "0", false, "brands", constants.CLAUSE_ID_VALIDATION, ctx)
	if err != nil {
		return nil, err
	}
	if intBrandID > 0 {
		mustActiveFields = append(mustActiveFields, m)
	}

	intOrderSourceID, m, err := d.getIntQueryWithMustActive("order_source_id", "0", false, "order_sources", constants.CLAUSE_ID_VALIDATION, ctx)
	if err != nil {
		return nil, err
	}
	if intOrderSourceID > 0 {
		mustActiveFields = append(mustActiveFields, m)
	}

	intOrderStatusID, m, err := d.getIntQueryWithMustActive("order_status_id", "0", false, "order_statuses", constants.CLAUSE_ID_VALIDATION, ctx)
	if err != nil {
		return nil, err
	}
	if intOrderStatusID > 0 {
		mustActiveFields = append(mustActiveFields, m)
	}

	intProductID, m, err := d.getIntQueryWithMustActive("product_id", "0", false, "products", constants.CLAUSE_ID_VALIDATION, ctx)
	if err != nil {
		return nil, err
	}
	if intProductID > 0 {
		mustActiveFields = append(mustActiveFields, m)
	}

	intCategoryID, m, err := d.getIntQueryWithMustActive("category_id", "0", false, "categories", constants.CLAUSE_ID_VALIDATION, ctx)
	if err != nil {
		return nil, err
	}
	if intCategoryID > 0 {
		mustActiveFields = append(mustActiveFields, m)
	}

	intSalesmanID, m, err := d.getIntQueryWithMustActive("salesman_id", "0", false, "salesmans", constants.CLAUSE_ID_VALIDATION, ctx)
	if err != nil {
		return nil, err
	}
	if intSalesmanID > 0 {
		mustActiveFields = append(mustActiveFields, m)
	}

	err = d.requestValidationMiddleware.MustActiveValidationCustomCode(404, ctx, mustActiveFields)
	if err != nil {
		return nil, err
	}

	intProvinceID, err := d.getIntQueryWithDefault("province_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intCityID, err := d.getIntQueryWithDefault("city_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intDistrictID, err := d.getIntQueryWithDefault("district_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intVillageID, err := d.getIntQueryWithDefault("village_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	dateFields := []*models.DateInputRequest{}

	startSoDate, dateFields := d.getQueryWithDateValidation("start_do_date", "", dateFields, ctx)

	endSoDate, dateFields := d.getQueryWithDateValidation("end_do_date", "", dateFields, ctx)

	startCreatedAt, dateFields := d.getQueryWithDateValidation("start_created_at", time.Now().AddDate(0, -1, 0).Format(constants.DATE_FORMAT_COMMON), dateFields, ctx)

	endCreatedAt, dateFields := d.getQueryWithDateValidation("end_created_at", time.Now().Format(constants.DATE_FORMAT_COMMON), dateFields, ctx)

	err = d.requestValidationMiddleware.DateInputValidation(ctx, dateFields, constants.ERROR_ACTION_NAME_GET)
	if err != nil {
		return nil, err
	}

	dStartDate, err := time.Parse(constants.DATE_FORMAT_COMMON, startCreatedAt)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return nil, err
	}

	dEndDate, err := time.Parse(constants.DATE_FORMAT_COMMON, endCreatedAt)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return nil, err
	}

	if dStartDate.Before(time.Now().AddDate(0, -3, 0)) {
		err = helper.NewError("Proses export tidak dapat dilakukan karena file yang akan di-export lebih dari 3 bulan dari periode tanggal buat. Silahkan cek file export kembali")
		ctx.JSON(http.StatusUnprocessableEntity, helper.GenerateResultByError(err, http.StatusUnprocessableEntity, constants.ERROR_INVALID_PROCESS))
		return nil, err
	}

	if dStartDate.After(dEndDate) {
		err = helper.NewError("Proses export tidak dapat dilakukan karena tanggal selesai melebihi tanggal mulai. Silahkan cek file export kembali")
		ctx.JSON(http.StatusUnprocessableEntity, helper.GenerateResultByError(err, http.StatusUnprocessableEntity, constants.ERROR_INVALID_PROCESS))
		return nil, err
	}

	salesOrderReqeuest := &models.SalesOrderExportRequest{
		SortField:         sortField,
		SortValue:         d.getQueryWithDefault("sort_value", "desc", ctx),
		GlobalSearchValue: d.getQueryWithDefault("global_search_value", "", ctx),
		FileType:          d.getQueryWithDefault("file_type", "xlsx", ctx),
		ID:                intID,
		AgentID:           intAgentID,
		StoreID:           intStoreID,
		BrandID:           intBrandID,
		OrderSourceID:     intOrderSourceID,
		OrderStatusID:     intOrderStatusID,
		StartSoDate:       startSoDate,
		EndSoDate:         endSoDate,
		ProductID:         intProductID,
		CategoryID:        intCategoryID,
		SalesmanID:        intSalesmanID,
		ProvinceID:        intProvinceID,
		CityID:            intCityID,
		DistrictID:        intDistrictID,
		VillageID:         intVillageID,
		StartCreatedAt:    startCreatedAt,
		EndCreatedAt:      endCreatedAt,
		UpdatedAt:         d.getQueryWithDefault("updated_at", "", ctx),
	}
	return salesOrderReqeuest, nil
}

func (d *SalesOrderValidator) ExportSalesOrderDetailValidator(ctx *gin.Context) (*models.SalesOrderDetailExportRequest, error) {
	sortField := d.getQueryWithDefault("sort_field", "created_at", ctx)

	if sortField != "order_status_id" && sortField != "do_date" && sortField != "do_ref_code" && sortField != "store_id" && sortField != "created_at" && sortField != "updated_at" {
		err := helper.NewError("Parameter 'sort_field' harus bernilai 'order_status_id' or 'do_date' or 'do_ref_code' or 'created_at' or 'updated_at'")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return nil, err
	}

	intID, err := d.getIntQueryWithDefault("id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	mustActiveFields := []*models.MustActiveRequest{}

	intAgentID, m, err := d.getIntQueryWithMustActive("agent_id", "0", false, "agents", constants.CLAUSE_ID_VALIDATION, ctx)
	if err != nil {
		return nil, err
	}
	if intAgentID > 0 {
		mustActiveFields = append(mustActiveFields, m)
	}

	intStoreID, m, err := d.getIntQueryWithMustActive("store_id", "0", false, "stores", constants.CLAUSE_ID_VALIDATION, ctx)
	if err != nil {
		return nil, err
	}
	if intStoreID > 0 {
		mustActiveFields = append(mustActiveFields, m)
	}

	intBrandID, m, err := d.getIntQueryWithMustActive("brand_id", "0", false, "brands", constants.CLAUSE_ID_VALIDATION, ctx)
	if err != nil {
		return nil, err
	}
	if intBrandID > 0 {
		mustActiveFields = append(mustActiveFields, m)
	}

	intOrderSourceID, m, err := d.getIntQueryWithMustActive("order_source_id", "0", false, "order_sources", constants.CLAUSE_ID_VALIDATION, ctx)
	if err != nil {
		return nil, err
	}
	if intOrderSourceID > 0 {
		mustActiveFields = append(mustActiveFields, m)
	}

	intOrderStatusID, m, err := d.getIntQueryWithMustActive("order_status_id", "0", false, "order_statuses", constants.CLAUSE_ID_VALIDATION, ctx)
	if err != nil {
		return nil, err
	}
	if intOrderStatusID > 0 {
		mustActiveFields = append(mustActiveFields, m)
	}

	intProductID, m, err := d.getIntQueryWithMustActive("product_id", "0", false, "products", constants.CLAUSE_ID_VALIDATION, ctx)
	if err != nil {
		return nil, err
	}
	if intProductID > 0 {
		mustActiveFields = append(mustActiveFields, m)
	}

	intCategoryID, m, err := d.getIntQueryWithMustActive("category_id", "0", false, "categories", constants.CLAUSE_ID_VALIDATION, ctx)
	if err != nil {
		return nil, err
	}
	if intCategoryID > 0 {
		mustActiveFields = append(mustActiveFields, m)
	}

	intSalesmanID, m, err := d.getIntQueryWithMustActive("salesman_id", "0", false, "salesmans", constants.CLAUSE_ID_VALIDATION, ctx)
	if err != nil {
		return nil, err
	}
	if intSalesmanID > 0 {
		mustActiveFields = append(mustActiveFields, m)
	}

	err = d.requestValidationMiddleware.MustActiveValidationCustomCode(404, ctx, mustActiveFields)
	if err != nil {
		return nil, err
	}

	intProvinceID, err := d.getIntQueryWithDefault("province_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intCityID, err := d.getIntQueryWithDefault("city_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intDistrictID, err := d.getIntQueryWithDefault("district_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	intVillageID, err := d.getIntQueryWithDefault("village_id", "0", false, ctx)
	if err != nil {
		return nil, err
	}

	dateFields := []*models.DateInputRequest{}

	startSoDate, dateFields := d.getQueryWithDateValidation("start_do_date", "", dateFields, ctx)

	endSoDate, dateFields := d.getQueryWithDateValidation("end_do_date", "", dateFields, ctx)

	startCreatedAt, dateFields := d.getQueryWithDateValidation("start_created_at", time.Now().AddDate(0, -1, 0).Format(constants.DATE_FORMAT_COMMON), dateFields, ctx)

	endCreatedAt, dateFields := d.getQueryWithDateValidation("end_created_at", time.Now().Format(constants.DATE_FORMAT_COMMON), dateFields, ctx)

	err = d.requestValidationMiddleware.DateInputValidation(ctx, dateFields, constants.ERROR_ACTION_NAME_GET)
	if err != nil {
		return nil, err
	}

	dStartDate, err := time.Parse(constants.DATE_FORMAT_COMMON, startCreatedAt)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return nil, err
	}

	dEndDate, err := time.Parse(constants.DATE_FORMAT_COMMON, endCreatedAt)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return nil, err
	}

	if dStartDate.Before(time.Now().AddDate(0, -3, 0)) {
		err = helper.NewError("Proses export tidak dapat dilakukan karena file yang akan di-export lebih dari 3 bulan dari periode tanggal buat. Silahkan cek file export kembali")
		ctx.JSON(http.StatusUnprocessableEntity, helper.GenerateResultByError(err, http.StatusUnprocessableEntity, constants.ERROR_INVALID_PROCESS))
		return nil, err
	}

	if dStartDate.After(dEndDate) {
		err = helper.NewError("Proses export tidak dapat dilakukan karena tanggal selesai melebihi tanggal mulai. Silahkan cek file export kembali")
		ctx.JSON(http.StatusUnprocessableEntity, helper.GenerateResultByError(err, http.StatusUnprocessableEntity, constants.ERROR_INVALID_PROCESS))
		return nil, err
	}

	salesOrderDetailReqeuest := &models.SalesOrderDetailExportRequest{
		SortField:         sortField,
		SortValue:         d.getQueryWithDefault("sort_value", "desc", ctx),
		GlobalSearchValue: d.getQueryWithDefault("global_search_value", "", ctx),
		FileType:          d.getQueryWithDefault("file_type", "xlsx", ctx),
		ID:                intID,
		AgentID:           intAgentID,
		StoreID:           intStoreID,
		BrandID:           intBrandID,
		OrderSourceID:     intOrderSourceID,
		OrderStatusID:     intOrderStatusID,
		StartSoDate:       startSoDate,
		EndSoDate:         endSoDate,
		ProductID:         intProductID,
		CategoryID:        intCategoryID,
		SalesmanID:        intSalesmanID,
		ProvinceID:        intProvinceID,
		CityID:            intCityID,
		DistrictID:        intDistrictID,
		VillageID:         intVillageID,
		StartCreatedAt:    startCreatedAt,
		EndCreatedAt:      endCreatedAt,
	}
	return salesOrderDetailReqeuest, nil
}

func (c *SalesOrderValidator) UpdateSalesOrderByIdValidator(updateRequest *models.SalesOrderUpdateRequest, ctx *gin.Context) error {
	ids := ctx.Param("so-id")
	id, _ := strconv.Atoi(ids)

	var errors []string
	if updateRequest.OrderStatusID != 5 && updateRequest.OrderStatusID != 9 && updateRequest.OrderStatusID != 10 {
		errors = append(errors, helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_UPDATE, fmt.Sprintf("status %d tidak terdaftar", updateRequest.OrderStatusID)))
	}
	for _, v := range updateRequest.SalesOrderDetails {
		if v.OrderStatusID != 11 && v.OrderStatusID != 15 && v.OrderStatusID != 16 {
			errors = append(errors, helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_UPDATE, fmt.Sprintf("status %d tidak terdaftar", v.OrderStatusID)))
		}
	}
	if len(errors) > 0 {
		errorLog := helper.NewWriteLog(model.ErrorLog{
			Message:       errors,
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))

		err := helper.NewError("")
		return err
	}

	// Get Sales Order By Id
	getSalesOrderByIDResultChan := make(chan *models.SalesOrderChan)
	go c.salesOrderRepository.GetByID(id, false, ctx, getSalesOrderByIDResultChan)
	getSalesOrderByIDResult := <-getSalesOrderByIDResultChan

	if getSalesOrderByIDResult.Error != nil {
		ctx.JSON(getSalesOrderByIDResult.ErrorLog.StatusCode, helper.GenerateResultByErrorLog(getSalesOrderByIDResult.ErrorLog))
		return getSalesOrderByIDResult.Error
	}

	// Get Sales Order By So Id
	getSalesOrderDetailBySoIDResultChan := make(chan *models.SalesOrderDetailsChan)
	go c.salesOrderDetailRepository.GetBySalesOrderID(id, false, ctx, getSalesOrderDetailBySoIDResultChan)
	getSalesOrderDetailBySoIDResult := <-getSalesOrderDetailBySoIDResultChan

	if getSalesOrderDetailBySoIDResult.Error != nil {
		ctx.JSON(getSalesOrderByIDResult.ErrorLog.StatusCode, helper.GenerateResultByErrorLog(getSalesOrderByIDResult.ErrorLog))
		return getSalesOrderDetailBySoIDResult.Error
	}

	if len(updateRequest.SalesOrderDetails) != int(getSalesOrderDetailBySoIDResult.Total) {
		errorLog := helper.NewWriteLog(model.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage("update", fmt.Sprintf("jumlah request so detail tidak sesuai dengan jumlah so detail sales order %d", id))},
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))

		err := helper.NewError("")
		return err
	}

	errors = []string{}
	for _, v := range updateRequest.SalesOrderDetails {
		isExist := false
		for _, y := range getSalesOrderDetailBySoIDResult.SalesOrderDetails {
			if v.ID == y.ID {
				isExist = true
			}
		}

		if !isExist {
			errors = append(errors, helper.GenerateUnprocessableErrorMessage("update", fmt.Sprintf("Sales Order Detail Id = %d tidak terdaftar pada Sales Order Id = %d", v.ID, id)))
		}
	}
	if len(errors) > 0 {
		errorLog := helper.NewWriteLog(model.ErrorLog{
			Message:       errors,
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))

		err := helper.NewError(strings.Join(errors, ""))
		return err
	}

	// Get Order Status By Id
	getOrderStatusResultChan := make(chan *models.OrderStatusChan)
	go c.orderStatusRepository.GetByID(getSalesOrderByIDResult.SalesOrder.OrderStatusID, false, ctx, getOrderStatusResultChan)
	getOrderStatusResult := <-getOrderStatusResultChan

	if getOrderStatusResult.Error != nil {
		errorLog := helper.NewWriteLog(model.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage("update", getOrderStatusResult.Error.Error())},
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))

		return getOrderStatusResult.Error
	}

	errorValidation := c.updateSOValidation(getSalesOrderByIDResult.SalesOrder.ID, getOrderStatusResult.OrderStatus.Name, ctx)

	if len(errorValidation) > 0 {
		err := errorValidation
		errorLog := helper.NewWriteLog(model.ErrorLog{
			Message:       errorValidation,
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))

		return fmt.Errorf(strings.Join(err, ";"))
	}

	return nil
}

func (c *SalesOrderValidator) UpdateSalesOrderDetailBySoIdValidator(updateRequest *models.SalesOrderUpdateRequest, ctx *gin.Context) error {
	ids := ctx.Param("so-id")
	id, _ := strconv.Atoi(ids)

	var errors []string
	if updateRequest.OrderStatusID != 5 && updateRequest.OrderStatusID != 9 && updateRequest.OrderStatusID != 10 {
		errors = append(errors, helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_UPDATE, fmt.Sprintf("status %d tidak terdaftar", updateRequest.OrderStatusID)))
	}
	for _, v := range updateRequest.SalesOrderDetails {
		if v.OrderStatusID != 11 && v.OrderStatusID != 15 && v.OrderStatusID != 16 {
			errors = append(errors, helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_UPDATE, fmt.Sprintf("status %d tidak terdaftar", v.OrderStatusID)))
		}
	}
	if len(errors) > 0 {
		errorLog := helper.NewWriteLog(model.ErrorLog{
			Message:       errors,
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))

		return errorLog.Err
	}

	// Get Sales Order By Id
	getSalesOrderByIDResultChan := make(chan *models.SalesOrderChan)
	go c.salesOrderRepository.GetByID(id, false, ctx, getSalesOrderByIDResultChan)
	getSalesOrderByIDResult := <-getSalesOrderByIDResultChan

	if getSalesOrderByIDResult.Error != nil {
		ctx.JSON(getSalesOrderByIDResult.ErrorLog.StatusCode, helper.GenerateResultByErrorLog(getSalesOrderByIDResult.ErrorLog))
		return getSalesOrderByIDResult.Error
	}

	// Get Sales Order By So Id
	getSalesOrderDetailBySoIDResultChan := make(chan *models.SalesOrderDetailsChan)
	go c.salesOrderDetailRepository.GetBySalesOrderID(id, false, ctx, getSalesOrderDetailBySoIDResultChan)
	getSalesOrderDetailBySoIDResult := <-getSalesOrderDetailBySoIDResultChan

	if getSalesOrderDetailBySoIDResult.Error != nil {
		ctx.JSON(getSalesOrderDetailBySoIDResult.ErrorLog.StatusCode, helper.GenerateResultByErrorLog(getSalesOrderDetailBySoIDResult.ErrorLog))
		return getSalesOrderDetailBySoIDResult.Error
	}

	errors = []string{}
	for _, v := range updateRequest.SalesOrderDetails {
		isExist := false
		for _, y := range getSalesOrderDetailBySoIDResult.SalesOrderDetails {
			if v.ID == y.ID {
				isExist = true
			}
		}

		if !isExist {
			errors = append(errors, helper.GenerateUnprocessableErrorMessage("update", fmt.Sprintf("Sales Order Detail Id = %d tidak terdaftar pada Sales Order Id = %d", v.ID, id)))
		}
	}
	if len(errors) > 0 {
		errorLog := helper.NewWriteLog(model.ErrorLog{
			Message:       errors,
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))

		err := helper.NewError(strings.Join(errors, ""))
		return err
	}

	errorValidation := c.updateSOValidation(id, getSalesOrderByIDResult.SalesOrder.OrderStatusName, ctx)
	if errorValidation != nil {
		err := errorValidation
		errorLog := helper.NewWriteLog(model.ErrorLog{
			Message:       err,
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))

		return fmt.Errorf(strings.Join(err, ";"))
	}

	return nil
}

func (c *SalesOrderValidator) UpdateSalesOrderDetailByIdValidator(updateRequest *models.UpdateSalesOrderDetailByIdRequest, ctx *gin.Context) error {
	var errors []string
	if updateRequest.OrderStatusID != 8 && updateRequest.OrderStatusID != 10 {
		errors = append(errors, helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_UPDATE, fmt.Sprintf("status %d tidak dikenal", updateRequest.OrderStatusID)))
	}

	if updateRequest.SalesOrderDetails.OrderStatusID != 14 && updateRequest.SalesOrderDetails.OrderStatusID != 16 {
		errors = append(errors, helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_UPDATE, fmt.Sprintf("status %d tidak dikenal", updateRequest.SalesOrderDetails.OrderStatusID)))
	}

	if len(errors) > 0 {
		errorLog := helper.NewWriteLog(model.ErrorLog{
			Message:       errors,
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))

		err := helper.NewError(strings.Join(errors, ""))
		return err
	}

	soIds := ctx.Param("so-id")
	soId, _ := strconv.Atoi(soIds)

	// Get Sales Order By Id
	getSalesOrderByIDResultChan := make(chan *models.SalesOrderChan)
	go c.salesOrderRepository.GetByID(soId, false, ctx, getSalesOrderByIDResultChan)
	getSalesOrderByIDResult := <-getSalesOrderByIDResultChan

	if getSalesOrderByIDResult.Error != nil {
		ctx.JSON(getSalesOrderByIDResult.ErrorLog.StatusCode, helper.GenerateResultByErrorLog(getSalesOrderByIDResult.ErrorLog))
		return getSalesOrderByIDResult.Error
	}

	soDetailIds := ctx.Param("so-detail-id")
	soDetailId, _ := strconv.Atoi(soDetailIds)

	// Get Sales Order By  Id
	getSalesOrderDetailByIDResultChan := make(chan *models.SalesOrderDetailChan)
	go c.salesOrderDetailRepository.GetByID(soDetailId, false, ctx, getSalesOrderDetailByIDResultChan)
	getSalesOrderDetailByIDResult := <-getSalesOrderDetailByIDResultChan

	if getSalesOrderDetailByIDResult.Error != nil {
		ctx.JSON(getSalesOrderDetailByIDResult.ErrorLog.StatusCode, helper.GenerateResultByErrorLog(getSalesOrderDetailByIDResult.ErrorLog))
		return getSalesOrderDetailByIDResult.Error
	}

	if getSalesOrderDetailByIDResult.SalesOrderDetail.SalesOrderID != soId {
		errorLog := helper.NewWriteLog(model.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage("update", fmt.Sprintf("Sales Order Detail Id = %d tidak terdaftar pada Sales Order Id = %d", soDetailId, soId))},
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))

		err := helper.NewError(strings.Join(errors, ""))
		return err
	}

	errorValidation := c.updateSOValidation(soId, getSalesOrderByIDResult.SalesOrder.OrderStatusName, ctx)
	if errorValidation != nil {
		err := errorValidation
		errorLog := helper.NewWriteLog(model.ErrorLog{
			Message:       err,
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))

		return fmt.Errorf(strings.Join(err, ";"))
	}

	salesOrder := getSalesOrderByIDResult.SalesOrder
	totalSoDetail := getSalesOrderDetailByIDResult.Total
	var soStatus string
	var soDetailStatus string

	if salesOrder.OrderStatusName == constants.ORDER_STATUS_OPEN {
		if totalSoDetail == 1 {
			soStatus = constants.ORDER_STATUS_CANCELLED
			soDetailStatus = constants.ORDER_STATUS_CANCELLED
		} else if helper.CheckSalesOrderDetailStatus(soDetailId, true, constants.ORDER_STATUS_CANCELLED, salesOrder.SalesOrderDetails) > 0 || helper.CheckSalesOrderDetailStatus(soDetailId, false, constants.ORDER_STATUS_CANCELLED, salesOrder.SalesOrderDetails) > 0 {
			soStatus = constants.ORDER_STATUS_OPEN
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

	if len(soStatus) < 1 || len(soDetailStatus) < 1 {
		errorLog := helper.NewWriteLog(model.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_UPDATE, fmt.Sprintf("tidak memenuhi syarat"))},
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})

		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))

		err := helper.NewError(strings.Join(errors, ""))
		return err

	}

	return nil
}

func (c *SalesOrderValidator) updateSOValidation(salesOrderId int, orderStatusName string, ctx context.Context) []string {

	if orderStatusName != constants.ORDER_STATUS_OPEN && orderStatusName != constants.ORDER_STATUS_PENDING && orderStatusName != constants.ORDER_STATUS_PARTIAL {
		return []string{helper.GenerateUnprocessableErrorMessage("update", fmt.Sprintf("status sales order %s", orderStatusName))}
	}

	getDeliveryOrderByIDResultChan := make(chan *models.DeliveryOrdersChan)
	go c.deliveryOrderRepository.GetBySalesOrderID(salesOrderId, false, ctx, getDeliveryOrderByIDResultChan)
	getDeliveryOrderByIDResult := <-getDeliveryOrderByIDResultChan

	if getDeliveryOrderByIDResult.Total == 0 {

		return nil

	} else {

		errors := []string{}
		for _, v := range getDeliveryOrderByIDResult.DeliveryOrders {

			if v.OrderStatusName != "cancel" {
				errors = append(errors, fmt.Sprintf("Sales Order tidak dapat diupdate dikarenakan ada Delivery Order dengan status %s", v.OrderStatusName))
			}
		}

		if len(errors) > 0 {
			return errors
		}

	}

	return nil
}

func (c *SalesOrderValidator) DeleteSalesOrderByIdValidator(sId string, ctx *gin.Context) (int, error) {
	id, err := strconv.Atoi(sId)

	if err != nil {
		err = helper.NewError(constants.ERROR_BAD_REQUEST_INT_ID_PARAMS)
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return 0, err
	}
	mustActiveField := []*models.MustActiveRequest{
		{
			Table:    "sales_orders",
			ReqField: "id",
			Clause:   fmt.Sprintf("id = %d", id),
		},
	}
	mustActiveField422 := []*models.MustActiveRequest{
		{
			Table:    "sales_orders",
			ReqField: "id",
			Clause:   fmt.Sprintf(constants.CLAUSE_ID_VALIDATION, id),
		},
	}

	err = c.requestValidationMiddleware.MustActiveValidation(ctx, mustActiveField)
	if err != nil {
		return 0, err
	}
	err = c.requestValidationMiddleware.MustActiveValidationCustomCode(422, ctx, mustActiveField422)
	if err != nil {
		return 0, err
	}
	mustEmpties := []*models.MustEmptyValidationRequest{
		{
			Table:           "sales_orders s JOIN order_statuses o ON s.order_status_id = o.id",
			SelectedCollumn: "o.name",
			Clause:          fmt.Sprintf("o.id = %d AND s.order_status_id NOT IN (5,6,9,10)", id),
			MessageFormat:   "Status Sales Order <result>",
		},
		{
			Table:           "delivery_orders d JOIN sales_orders s ON d.sales_order_id = s.id",
			SelectedCollumn: "d.id",
			Clause:          fmt.Sprintf("s.id = %d AND d.deleted_at IS NULL", id),
			MessageFormat:   constants.ERROR_UPDATE_SO_MESSAGE,
		},
	}
	err = c.requestValidationMiddleware.MustEmptyValidation(ctx, mustEmpties)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (c *SalesOrderValidator) DeleteSalesOrderDetailByIdValidator(sId string, ctx *gin.Context) (int, error) {
	id, err := strconv.Atoi(sId)

	if err != nil {
		err = helper.NewError(constants.ERROR_BAD_REQUEST_INT_ID_PARAMS)
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
		return 0, err
	}
	mustActiveField := []*models.MustActiveRequest{
		{
			Table:    "sales_order_details",
			ReqField: "id",
			Clause:   fmt.Sprintf("id = %d", id),
		},
	}
	mustActiveField422 := []*models.MustActiveRequest{
		{
			Table:    "sales_order_details",
			ReqField: "id",
			Clause:   fmt.Sprintf(constants.CLAUSE_ID_VALIDATION, id),
		},
	}

	err = c.requestValidationMiddleware.MustActiveValidation(ctx, mustActiveField)
	if err != nil {
		return id, err
	}
	err = c.requestValidationMiddleware.MustActiveValidationCustomCode(422, ctx, mustActiveField422)
	if err != nil {
		return id, err
	}
	mustEmpties := []*models.MustEmptyValidationRequest{
		{
			Table:           "sales_order_details s JOIN order_statuses o ON s.order_status_id = o.id",
			SelectedCollumn: "o.name",
			Clause:          fmt.Sprintf("s.id = %d AND s.order_status_id NOT IN (16)", id),
			MessageFormat:   "Hanya status cancelled pada Sales Order Detail yang dapat di delete",
		},
		{
			Table:           "sales_orders s JOIN sales_order_details sd ON s.id = sd.sales_order_id JOIN delivery_orders d ON d.sales_order_id = s.id",
			SelectedCollumn: "d.id",
			Clause:          fmt.Sprintf("sd.id = %d AND d.deleted_at IS NULL", id),
			MessageFormat:   constants.ERROR_UPDATE_SO_MESSAGE,
		},
	}
	err = c.requestValidationMiddleware.MustEmptyValidation(ctx, mustEmpties)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (c *SalesOrderValidator) DeleteSalesOrderDetailBySoIdValidator(sId string, ctx *gin.Context) (int, error) {
	id, err := strconv.Atoi(sId)

	if err != nil {
		if err != nil {
			err = helper.NewError(constants.ERROR_BAD_REQUEST_INT_ID_PARAMS)
			ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest, ""))
			return 0, err
		}
	}
	mustActiveField := []*models.MustActiveRequest{
		{
			Table:    "sales_orders",
			ReqField: "id",
			Clause:   fmt.Sprintf("id = %d", id),
		},
	}
	mustActiveField422 := []*models.MustActiveRequest{
		{
			Table:    "sales_orders",
			ReqField: "id",
			Clause:   fmt.Sprintf(constants.CLAUSE_ID_VALIDATION, id),
		},
	}

	err = c.requestValidationMiddleware.MustActiveValidation(ctx, mustActiveField)
	if err != nil {
		return id, err
	}
	err = c.requestValidationMiddleware.MustActiveValidationCustomCode(422, ctx, mustActiveField422)
	if err != nil {
		return id, err
	}
	mustEmpties := []*models.MustEmptyValidationRequest{
		{
			Table:           "sales_orders s JOIN order_statuses o ON s.order_status_id = o.id",
			SelectedCollumn: "o.name",
			Clause:          fmt.Sprintf("o.id = %d AND s.order_status_id NOT IN (5,6,9,10)", id),
			MessageFormat:   "Status Sales Order <result>",
		},
		{
			Table:           "sales_orders s JOIN sales_order_details sd ON s.id = sd.sales_order_id JOIN order_statuses o ON sd.order_status_id = o.id",
			SelectedCollumn: "o.name",
			Clause:          fmt.Sprintf("s.id = %d AND sd.order_status_id NOT IN (16)", id),
			MessageFormat:   "Hanya status cancelled pada Sales Order Detail yang dapat di delete",
		},
		{
			Table:           "delivery_orders d JOIN sales_orders s ON d.sales_order_id = s.id",
			SelectedCollumn: "d.id",
			Clause:          fmt.Sprintf("s.id = %d AND d.deleted_at IS NULL", id),
			MessageFormat:   constants.ERROR_UPDATE_SO_MESSAGE,
		},
	}
	err = c.requestValidationMiddleware.MustEmptyValidation(ctx, mustEmpties)
	if err != nil {
		return id, err
	}
	return id, nil
}
