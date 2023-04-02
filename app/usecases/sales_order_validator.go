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
	baseModel "order-service/global/utils/model"
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
}

type SalesOrderValidator struct {
	requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface
	orderSourceRepository       repositories.OrderSourceRepositoryInterface
	salesmanRepository          repositories.SalesmanRepositoryInterface
	db                          dbresolver.DB
	ctx                         context.Context
}

func InitSalesOrderValidator(requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, orderSourceRepository repositories.OrderSourceRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, db dbresolver.DB, ctx context.Context) SalesOrderValidatorInterface {
	return &SalesOrderValidator{
		requestValidationMiddleware: requestValidationMiddleware,
		orderSourceRepository:       orderSourceRepository,
		salesmanRepository:          salesmanRepository,
		db:                          db,
		ctx:                         ctx,
	}
}

func (c *SalesOrderValidator) CreateSalesOrderValidator(insertRequest *models.SalesOrderStoreRequest, ctx *gin.Context) error {
	var result baseModel.Response
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
			Clause:   fmt.Sprintf("id = %d AND deleted_at IS NULL", v.UomID),
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
		result.StatusCode = getOrderSourceResult.ErrorLog.StatusCode
		result.Error = getOrderSourceResult.ErrorLog
		ctx.JSON(result.StatusCode, result)
		return getOrderSourceResult.Error
	}
	sourceName := getOrderSourceResult.OrderSource.SourceName
	if sourceName != "manager" && len(insertRequest.DeviceId) < 1 {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage("create", "device_id tidak boleh kosong")},
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		result.StatusCode = errorLog.StatusCode
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)

		err = helper.NewError("device_id tidak boleh kosong")
		return err
	}

	if sourceName == "store" && len(insertRequest.ReferralCode) > 0 {
		// Get Salesmans By Agent Id
		getSalesmanResultChan := make(chan *models.SalesmansChan)
		go c.salesmanRepository.GetByAgentId(insertRequest.AgentID, false, ctx, getSalesmanResultChan)
		getSalesmanResult := <-getSalesmanResultChan

		if getSalesmanResult.Error != nil {
			result.StatusCode = getSalesmanResult.ErrorLog.StatusCode
			result.Error = getSalesmanResult.ErrorLog
			ctx.JSON(result.StatusCode, result)
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
			errorLog := helper.NewWriteLog(baseModel.ErrorLog{
				Message:       []string{helper.GenerateUnprocessableErrorMessage("create", "referral code tidak terdaftar")},
				SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
				StatusCode:    http.StatusUnprocessableEntity,
			})
			result.StatusCode = errorLog.StatusCode
			result.Error = errorLog
			ctx.JSON(result.StatusCode, result)

			err = helper.NewError("referral code tidak terdaftar")
			return err
		}
	}

	now := time.Now().UTC().Add(7 * time.Hour)
	parseSoDate, _ := time.Parse("2006-01-02", insertRequest.SoDate)
	parseSoRefDate, _ := time.Parse("2006-01-02", insertRequest.SoRefDate)
	duration := time.Hour*time.Duration(now.Hour()) + time.Minute*time.Duration(now.Minute()) + time.Second*time.Duration(now.Second()) + time.Nanosecond*time.Duration(now.Nanosecond())

	soDate := parseSoDate.UTC().Add(duration)
	soRefDate := parseSoRefDate.UTC().Add(duration)
	nowUTC := now.UTC()
	if now.UTC().Hour() < soRefDate.Hour() {
		nowUTC = nowUTC.Add(7 * time.Hour)
	}

	if sourceName == "manager" && !(soDate.Add(1*time.Minute).After(soRefDate) && soRefDate.Add(-1*time.Minute).Before(nowUTC) && soDate.Add(-1*time.Minute).Before(nowUTC) && soRefDate.Month() == nowUTC.Month() && soRefDate.UTC().Year() == nowUTC.Year()) {

		err = helper.NewError("so_date dan so_ref_date harus sama dengan kurang dari hari ini dan harus di bulan berjalan")
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage("create", "so_date dan so_ref_date harus sama dengan kurang dari hari ini dan harus di bulan berjalan")},
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		result.StatusCode = errorLog.StatusCode
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)

		return err

	} else if (sourceName == "salesman" || sourceName == "store") && !(soDate.Equal(now.Local()) && soRefDate.Equal(now.Local())) {

		err = helper.NewError("so_date dan so_ref_date harus sama dengan hari ini")
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage("create", "so_date dan so_ref_date harus sama dengan hari ini")},
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		result.StatusCode = errorLog.StatusCode
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)

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
	var result baseModel.Response

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
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(result.StatusCode, result)
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
	var result baseModel.Response

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
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(result.StatusCode, result)
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
	var result baseModel.Response

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
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(result.StatusCode, result)
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
	var result baseModel.Response
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
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(result.StatusCode, result)
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
	var result baseModel.Response

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
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(result.StatusCode, result)
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
	var result baseModel.Response

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
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(result.StatusCode, result)
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
	var response baseModel.Response
	sResult := d.getQueryWithDefault(param, empty, ctx)
	result, err := strconv.Atoi(sResult)
	if err != nil {
		err = helper.NewError(fmt.Sprintf("Parameter '%s' harus bernilai integer", param))
		response.StatusCode = http.StatusBadRequest
		response.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(response.StatusCode, response)
		return 0, err
	}
	if result == 0 && isNotZero {
		err = helper.NewError(fmt.Sprintf("Parameter '%s' harus bernilai integer > 0", param))
		response.StatusCode = http.StatusBadRequest
		response.Error = helper.WriteLog(err, http.StatusBadRequest, err.Error())
		ctx.JSON(response.StatusCode, response)
		return 0, err
	}
	return result, nil
}
