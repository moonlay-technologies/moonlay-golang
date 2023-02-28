package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"order-service/app/middlewares"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/app/usecases"
	"order-service/global/utils/helper"
	baseModel "order-service/global/utils/model"
	"strconv"

	"github.com/bxcodec/dbresolver"
	"github.com/gin-gonic/gin"
)

type SalesOrderControllerInterface interface {
	GetByID(ctx *gin.Context)
	Create(ctx *gin.Context)
	Get(ctx *gin.Context)
	UpdateByID(ctx *gin.Context)
	UpdateSODetailByID(ctx *gin.Context)
	UpdateSODetailBySOID(ctx *gin.Context)
	GetDetails(ctx *gin.Context)
	GetDetailsBySoId(ctx *gin.Context)
	GetDetailsById(ctx *gin.Context)
	DeleteByID(ctx *gin.Context)
}

type salesOrderController struct {
	cartUseCase                 usecases.CartUseCaseInterface
	salesOrderUseCase           usecases.SalesOrderUseCaseInterface
	requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface
	db                          dbresolver.DB
	ctx                         context.Context
}

func InitSalesOrderController(cartUseCase usecases.CartUseCaseInterface, salesOrderUseCase usecases.SalesOrderUseCaseInterface, requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, db dbresolver.DB, ctx context.Context) SalesOrderControllerInterface {
	return &salesOrderController{
		cartUseCase:                 cartUseCase,
		salesOrderUseCase:           salesOrderUseCase,
		requestValidationMiddleware: requestValidationMiddleware,
		db:                          db,
		ctx:                         ctx,
	}
}

func (c *salesOrderController) Get(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	var pageInt, perPageInt, agentIdInt, storeIdInt, brandIdInt, orderSourceIdInt, orderStatusIdInt, productIdInt, categoryIdInt, salesmanIdInt, provinceIdInt, cityIdInt, districtIdInt, villageIdInt, idInt int

	page, isPageExist := ctx.GetQuery("page")
	if isPageExist == false {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)

	if err != nil {
		err = helper.NewError("Parameter 'page' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	if pageInt == 0 {
		err = helper.NewError("Parameter 'page' harus bernilai integer > 0")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	perPage, isPerPageExist := ctx.GetQuery("per_page")
	if isPerPageExist == false {
		perPage = "10"
	}

	perPageInt, err = strconv.Atoi(perPage)
	if err != nil {
		err = helper.NewError("Parameter 'per_page' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	if perPageInt == 0 {
		err = helper.NewError("Parameter 'per_page' harus bernilai integer > 0")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	sortField, isSortFieldExist := ctx.GetQuery("sort_field")
	if isSortFieldExist == false {
		sortField = "created_at"
	}

	if sortField != "order_status_id" && sortField != "so_date" && sortField != "so_ref_code" && sortField != "so_code" && sortField != "store_code" && sortField != "store_name" && sortField != "created_at" && sortField != "updated_at" {
		err = helper.NewError("Parameter 'sort_field' harus bernilai 'order_status_id' or 'so_date' or 'so_ref_code' or 'so_code' or 'store_code' or 'store_name' or 'created_at' or 'updated_at' ")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	sortValue, isSortValueExist := ctx.GetQuery("sort_value")
	if isSortValueExist == false {
		sortValue = "desc"
	}

	globalSearchValue, isGlobalSearchValueExist := ctx.GetQuery("global_search_value")
	if isGlobalSearchValueExist == false {
		globalSearchValue = ""
	}

	agentId, isAgentIdExist := ctx.GetQuery("agent_id")
	if isAgentIdExist == false {
		agentId = "0"
	}

	agentIdInt, err = strconv.Atoi(agentId)
	if err != nil {
		err = helper.NewError("Parameter 'agent_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	storeId, isStoreIdExist := ctx.GetQuery("store_id")
	if isStoreIdExist == false {
		storeId = "0"
	}

	storeIdInt, err = strconv.Atoi(storeId)
	if err != nil {
		err = helper.NewError("Parameter 'store_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	brandId, isBrandIdExist := ctx.GetQuery("brand_id")
	if isBrandIdExist == false {
		brandId = "0"
	}

	brandIdInt, err = strconv.Atoi(brandId)
	if err != nil {
		err = helper.NewError("Parameter 'brand_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	orderSourceId, isOrderSourceIdExist := ctx.GetQuery("order_source_id")
	if isOrderSourceIdExist == false {
		orderSourceId = "0"
	}

	orderSourceIdInt, err = strconv.Atoi(orderSourceId)
	if err != nil {
		err = helper.NewError("Parameter 'order_source_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	orderStatusId, isOrderStatusIdExist := ctx.GetQuery("order_status_id")
	if isOrderStatusIdExist == false {
		orderStatusId = "0"
	}

	orderStatusIdInt, err = strconv.Atoi(orderStatusId)
	if err != nil {
		err = helper.NewError("Parameter 'order_status_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	startSoDate, isStartSoDate := ctx.GetQuery("start_so_date")
	if isStartSoDate == false {
		startSoDate = ""
	}

	endSoDate, isEndSoDate := ctx.GetQuery("end_so_date")
	if isEndSoDate == false {
		endSoDate = ""
	}

	id, isIdExist := ctx.GetQuery("id")
	if isIdExist == false {
		id = "0"
	}

	idInt, err = strconv.Atoi(id)
	if err != nil {
		err = helper.NewError("Parameter 'id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	productId, isProductIdExist := ctx.GetQuery("product_id")
	if isProductIdExist == false {
		productId = "0"
	}

	productIdInt, err = strconv.Atoi(productId)
	if err != nil {
		err = helper.NewError("Parameter 'product_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	categoryId, isCategoryIdExist := ctx.GetQuery("category_id")
	if isCategoryIdExist == false {
		categoryId = "0"
	}

	categoryIdInt, err = strconv.Atoi(categoryId)
	if err != nil {
		err = helper.NewError("Parameter 'category_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	salesmanId, isSalesmanIdExist := ctx.GetQuery("salesman_id")
	if isSalesmanIdExist == false {
		salesmanId = "0"
	}

	salesmanIdInt, err = strconv.Atoi(salesmanId)
	if err != nil {
		err = helper.NewError("Parameter 'salesman_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	provinceId, isProvinceIdExist := ctx.GetQuery("province_id")
	if isProvinceIdExist == false {
		provinceId = "0"
	}

	provinceIdInt, err = strconv.Atoi(provinceId)
	if err != nil {
		err = helper.NewError("Parameter 'province_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	cityId, isCityIdExist := ctx.GetQuery("city_id")
	if isCityIdExist == false {
		cityId = "0"
	}

	cityIdInt, err = strconv.Atoi(cityId)
	if err != nil {
		err = helper.NewError("Parameter 'city_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	districtId, isDistrictIdExist := ctx.GetQuery("district_id")
	if isDistrictIdExist == false {
		districtId = "0"
	}

	districtIdInt, err = strconv.Atoi(districtId)
	if err != nil {
		err = helper.NewError("Parameter 'district_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	villageId, isVillageIdExist := ctx.GetQuery("village_id")
	if isVillageIdExist == false {
		villageId = "0"
	}

	villageIdInt, err = strconv.Atoi(villageId)
	if err != nil {
		err = helper.NewError("Parameter 'village_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	startCreatedAt, isStartCreatedAt := ctx.GetQuery("start_created_at")
	if isStartCreatedAt == false {
		startCreatedAt = ""
	}

	endCreatedAt, isEndCreatedAt := ctx.GetQuery("end_created_at")
	if isEndCreatedAt == false {
		endCreatedAt = ""
	}

	salesOrderRequest := &models.SalesOrderRequest{
		Page:              pageInt,
		PerPage:           perPageInt,
		SortField:         sortField,
		SortValue:         sortValue,
		GlobalSearchValue: globalSearchValue,
		AgentID:           agentIdInt,
		StoreID:           storeIdInt,
		BrandID:           brandIdInt,
		OrderSourceID:     orderSourceIdInt,
		OrderStatusID:     orderStatusIdInt,
		StartSoDate:       startSoDate,
		EndSoDate:         endSoDate,
		ID:                idInt,
		ProductID:         productIdInt,
		CategoryID:        categoryIdInt,
		SalesmanID:        salesmanIdInt,
		ProvinceID:        provinceIdInt,
		CityID:            cityIdInt,
		DistrictID:        districtIdInt,
		VillageID:         villageIdInt,
		StartCreatedAt:    startCreatedAt,
		EndCreatedAt:      endCreatedAt,
	}

	salesOrders, errorLog := c.salesOrderUseCase.Get(salesOrderRequest)

	if errorLog.Err != nil {
		resultErrorLog = errorLog
		result.StatusCode = resultErrorLog.StatusCode
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = salesOrders.SalesOrders
	result.Total = salesOrders.Total
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}

func (c *salesOrderController) GetByID(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	var id int

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	ids := ctx.Param("so-id")
	id, err := strconv.Atoi(ids)

	if err != nil {
		err = helper.NewError("Parameter 'id' harus bernilai integer")
		resultErrorLog.Message = err.Error()
		result.StatusCode = http.StatusBadRequest
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	salesOrderRequest := &models.SalesOrderRequest{
		ID: id,
	}

	salesOrder, errorLog := c.salesOrderUseCase.GetByID(salesOrderRequest, ctx)

	if errorLog.Err != nil {
		resultErrorLog = errorLog
		result.StatusCode = resultErrorLog.StatusCode
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = salesOrder
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}

func (c *salesOrderController) Create(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	insertRequest := &models.SalesOrderStoreRequest{}

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	err := ctx.BindJSON(insertRequest)

	if err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError

		if errors.As(err, &unmarshalTypeError) {
			c.requestValidationMiddleware.DataTypeValidation(ctx, err, unmarshalTypeError)
			return
		} else {
			c.requestValidationMiddleware.MandatoryValidation(ctx, err)
			return
		}
	}

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
	err = c.requestValidationMiddleware.DateInputValidation(ctx, dateField, constants.ERROR_ACTION_NAME_CREATE)
	if err != nil {
		return
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
		return
	}

	uniqueField := []*models.UniqueRequest{}
	if len(insertRequest.SoRefCode) > 0 {
		err = c.requestValidationMiddleware.OrderSourceValidation(ctx, insertRequest.OrderSourceID, insertRequest.SoRefCode, constants.ERROR_ACTION_NAME_CREATE)
		if err != nil {
			return
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
			return
		}
	}

	err = c.requestValidationMiddleware.AgentIdValidation(ctx, insertRequest.AgentID, insertRequest.UserID, constants.ERROR_ACTION_NAME_CREATE)
	if err != nil {
		return
	}

	err = c.requestValidationMiddleware.StoreIdValidation(ctx, insertRequest.StoreID, insertRequest.AgentID, constants.ERROR_ACTION_NAME_CREATE)
	if err != nil {
		return
	}

	if insertRequest.SalesmanID > 0 {
		err = c.requestValidationMiddleware.SalesmanIdValidation(ctx, insertRequest.SalesmanID, insertRequest.AgentID, constants.ERROR_ACTION_NAME_CREATE)
		if err != nil {
			return
		}
	}

	err = c.requestValidationMiddleware.BrandIdValidation(ctx, brandIds, insertRequest.AgentID, constants.ERROR_ACTION_NAME_CREATE)
	if err != nil {
		return
	}

	dbTransaction, err := c.db.BeginTx(ctx, nil)

	if err != nil {
		errorLog := helper.WriteLog(err, http.StatusInternalServerError, nil)
		resultErrorLog = errorLog
		result.StatusCode = http.StatusInternalServerError
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	salesOrders, errorLog := c.salesOrderUseCase.Create(insertRequest, dbTransaction, ctx)

	if errorLog != nil {
		err = dbTransaction.Rollback()

		if err != nil {
			errorLog = helper.WriteLog(err, http.StatusInternalServerError, nil)
			resultErrorLog = errorLog
			result.StatusCode = http.StatusInternalServerError
			result.Error = resultErrorLog
			ctx.JSON(result.StatusCode, result)
			return
		}

		result.StatusCode = errorLog.StatusCode
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	err = dbTransaction.Commit()

	if err != nil {
		errorLog = helper.WriteLog(err, http.StatusInternalServerError, nil)
		resultErrorLog = errorLog
		result.StatusCode = http.StatusInternalServerError
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = salesOrders
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}

func (c *salesOrderController) UpdateByID(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	var id int
	updateRequest := &models.SalesOrderUpdateRequest{}

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	ids := ctx.Param("so-id")
	id, err := strconv.Atoi(ids)

	if err != nil {
		err = helper.NewError("Parameter 'id' harus bernilai integer")
		resultErrorLog.Message = err.Error()
		result.StatusCode = http.StatusBadRequest
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	err = ctx.BindJSON(updateRequest)

	if err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError

		if errors.As(err, &unmarshalTypeError) {
			c.requestValidationMiddleware.DataTypeValidation(ctx, err, unmarshalTypeError)
			return
		} else {
			c.requestValidationMiddleware.MandatoryValidation(ctx, err)
			return
		}
	}

	mustActiveField := []*models.MustActiveRequest{
		{
			Table:    "agents",
			ReqField: "agent_id",
			Clause:   fmt.Sprintf("id = %d AND status = '%s'", updateRequest.AgentID, "active"),
		},
		{
			Table:    "stores",
			ReqField: "store_id",
			Clause:   fmt.Sprintf("id = %d AND status = '%s'", updateRequest.StoreID, "active"),
		},
		{
			Table:    "users",
			ReqField: "user_id",
			Clause:   fmt.Sprintf("id = %d AND status = '%s'", updateRequest.UserID, "ACTIVE"),
		},
	}
	for i, v := range updateRequest.SalesOrderDetails {
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
	}
	err = c.requestValidationMiddleware.MustActiveValidation(ctx, mustActiveField)
	if err != nil {
		return
	}

	dbTransaction, err := c.db.BeginTx(ctx, nil)

	if err != nil {
		errorLog := helper.WriteLog(err, http.StatusInternalServerError, nil)
		resultErrorLog = errorLog
		result.StatusCode = http.StatusInternalServerError
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	salesOrder, errorLog := c.salesOrderUseCase.UpdateById(id, updateRequest, dbTransaction, ctx)

	if errorLog != nil {
		err = dbTransaction.Rollback()

		if err != nil {
			errorLog = helper.WriteLog(err, http.StatusInternalServerError, "Ada kesalahan, silahkan coba lagi nanti")
			resultErrorLog = errorLog
			result.StatusCode = http.StatusInternalServerError
			result.Error = resultErrorLog
			ctx.JSON(result.StatusCode, result)
			return
		}

		result.StatusCode = errorLog.StatusCode
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	err = dbTransaction.Commit()

	if err != nil {
		errorLog = helper.WriteLog(err, http.StatusInternalServerError, "Ada kesalahan, silahkan coba lagi nanti")
		resultErrorLog = errorLog
		result.StatusCode = http.StatusInternalServerError
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = salesOrder
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}

func (c *salesOrderController) UpdateSODetailByID(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	updateRequest := &models.SalesOrderDetailUpdateRequest{}

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	soIds := ctx.Param("so-id")
	soId, err := strconv.Atoi(soIds)

	if err != nil {
		err = helper.NewError("Parameter 'so-id' harus bernilai integer")
		resultErrorLog.Message = err.Error()
		result.StatusCode = http.StatusBadRequest
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	ids := ctx.Param("id")
	id, err := strconv.Atoi(ids)

	if err != nil {
		err = helper.NewError("Parameter 'id' harus bernilai integer")
		resultErrorLog.Message = err.Error()
		result.StatusCode = http.StatusBadRequest
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	err = ctx.BindJSON(updateRequest)

	if err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError

		if errors.As(err, &unmarshalTypeError) {
			c.requestValidationMiddleware.DataTypeValidation(ctx, err, unmarshalTypeError)
			return
		} else {
			c.requestValidationMiddleware.MandatoryValidation(ctx, err)
			return
		}
	}

	mustActiveField := []*models.MustActiveRequest{
		{
			Table:    "products",
			ReqField: "product_id",
			Clause:   fmt.Sprintf("id = %d AND isActive = %d", updateRequest.ProductID, 1),
		},
		{
			Table:    "uoms",
			ReqField: "uom_id",
			Clause:   fmt.Sprintf("id = %d AND deleted_at IS NULL", updateRequest.UomID),
		},
		{
			Table:    "brands",
			ReqField: "brand_id",
			Clause:   fmt.Sprintf("id = %d AND status_active = %d", updateRequest.BrandID, 1),
		},
	}

	err = c.requestValidationMiddleware.MustActiveValidation(ctx, mustActiveField)
	if err != nil {
		return
	}

	dbTransaction, err := c.db.BeginTx(ctx, nil)

	if err != nil {
		errorLog := helper.WriteLog(err, http.StatusInternalServerError, nil)
		resultErrorLog = errorLog
		result.StatusCode = http.StatusInternalServerError
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	salesOrderDetail, errorLog := c.salesOrderUseCase.UpdateSODetailById(soId, id, updateRequest, dbTransaction, ctx)

	if errorLog != nil {
		err = dbTransaction.Rollback()

		if err != nil {
			errorLog = helper.WriteLog(err, http.StatusInternalServerError, "Ada kesalahan, silahkan coba lagi nanti")
			resultErrorLog = errorLog
			result.StatusCode = http.StatusInternalServerError
			result.Error = resultErrorLog
			ctx.JSON(result.StatusCode, result)
			return
		}

		result.StatusCode = errorLog.StatusCode
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	err = dbTransaction.Commit()

	if err != nil {
		errorLog = helper.WriteLog(err, http.StatusInternalServerError, "Ada kesalahan, silahkan coba lagi nanti")
		resultErrorLog = errorLog
		result.StatusCode = http.StatusInternalServerError
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = salesOrderDetail
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}

func (c *salesOrderController) UpdateSODetailBySOID(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	var soId int
	var updateRequest []*models.SalesOrderDetailUpdateRequest

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	ids := ctx.Param("so-id")
	soId, err := strconv.Atoi(ids)

	if err != nil {
		err = helper.NewError("Parameter 'id' harus bernilai integer")
		resultErrorLog.Message = err.Error()
		result.StatusCode = http.StatusBadRequest
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	err = ctx.BindJSON(&updateRequest)

	if err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError

		if errors.As(err, &unmarshalTypeError) {
			c.requestValidationMiddleware.DataTypeValidation(ctx, err, unmarshalTypeError)
			return
		} else {
			c.requestValidationMiddleware.MandatoryValidation(ctx, err)
			return
		}
	}

	mustActiveField := []*models.MustActiveRequest{}
	for i, v := range updateRequest {
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
	}

	err = c.requestValidationMiddleware.MustActiveValidation(ctx, mustActiveField)
	if err != nil {
		return
	}

	dbTransaction, err := c.db.BeginTx(ctx, nil)

	if err != nil {
		errorLog := helper.WriteLog(err, http.StatusInternalServerError, nil)
		resultErrorLog = errorLog
		result.StatusCode = http.StatusInternalServerError
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	salesOrderDetail, errorLog := c.salesOrderUseCase.UpdateSODetailBySOId(soId, updateRequest, dbTransaction, ctx)

	if errorLog != nil {
		err = dbTransaction.Rollback()

		if err != nil {
			errorLog = helper.WriteLog(err, http.StatusInternalServerError, "Ada kesalahan, silahkan coba lagi nanti")
			resultErrorLog = errorLog
			result.StatusCode = http.StatusInternalServerError
			result.Error = resultErrorLog
			ctx.JSON(result.StatusCode, result)
			return
		}

		result.StatusCode = errorLog.StatusCode
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	err = dbTransaction.Commit()

	if err != nil {
		errorLog = helper.WriteLog(err, http.StatusInternalServerError, "Ada kesalahan, silahkan coba lagi nanti")
		resultErrorLog = errorLog
		result.StatusCode = http.StatusInternalServerError
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = salesOrderDetail
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}

func (c *salesOrderController) GetDetails(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	var pageInt, perPageInt, agentIdInt, storeIdInt, brandIdInt, orderSourceIdInt, orderStatusIdInt, productIdInt, categoryIdInt, salesmanIdInt, provinceIdInt, cityIdInt, districtIdInt, villageIdInt, idInt int

	page, isPageExist := ctx.GetQuery("page")
	if isPageExist == false {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)

	if err != nil {
		err = helper.NewError("Parameter 'page' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	if pageInt == 0 {
		err = helper.NewError("Parameter 'page' harus bernilai integer > 0")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	perPage, isPerPageExist := ctx.GetQuery("per_page")
	if isPerPageExist == false {
		perPage = "10"
	}

	perPageInt, err = strconv.Atoi(perPage)
	if err != nil {
		err = helper.NewError("Parameter 'per_page' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	if perPageInt == 0 {
		err = helper.NewError("Parameter 'per_page' harus bernilai integer > 0")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	sortField, isSortFieldExist := ctx.GetQuery("sort_field")
	if isSortFieldExist == false {
		sortField = "created_at"
	}

	if sortField != "order_status_id" && sortField != "so_date" && sortField != "so_ref_code" && sortField != "so_code" && sortField != "store_code" && sortField != "store_name" && sortField != "created_at" && sortField != "updated_at" {
		err = helper.NewError("Parameter 'sort_field' harus bernilai 'order_status_id' or 'so_date' or 'so_ref_code' or 'so_code' or 'store_code' or 'store_name' or 'created_at' or 'updated_at' ")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	sortValue, isSortValueExist := ctx.GetQuery("sort_value")
	if isSortValueExist == false {
		sortValue = "desc"
	}

	globalSearchValue, isGlobalSearchValueExist := ctx.GetQuery("global_search_value")
	if isGlobalSearchValueExist == false {
		globalSearchValue = ""
	}

	agentId, isAgentIdExist := ctx.GetQuery("agent_id")
	if isAgentIdExist == false {
		agentId = "0"
	}

	agentIdInt, err = strconv.Atoi(agentId)
	if err != nil {
		err = helper.NewError("Parameter 'agent_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	storeId, isStoreIdExist := ctx.GetQuery("store_id")
	if isStoreIdExist == false {
		storeId = "0"
	}

	storeIdInt, err = strconv.Atoi(storeId)
	if err != nil {
		err = helper.NewError("Parameter 'store_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	brandId, isBrandIdExist := ctx.GetQuery("brand_id")
	if isBrandIdExist == false {
		brandId = "0"
	}

	brandIdInt, err = strconv.Atoi(brandId)
	if err != nil {
		err = helper.NewError("Parameter 'brand_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	orderSourceId, isOrderSourceIdExist := ctx.GetQuery("order_source_id")
	if isOrderSourceIdExist == false {
		orderSourceId = "0"
	}

	orderSourceIdInt, err = strconv.Atoi(orderSourceId)
	if err != nil {
		err = helper.NewError("Parameter 'order_source_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	orderStatusId, isOrderStatusIdExist := ctx.GetQuery("order_status_id")
	if isOrderStatusIdExist == false {
		orderStatusId = "0"
	}

	orderStatusIdInt, err = strconv.Atoi(orderStatusId)
	if err != nil {
		err = helper.NewError("Parameter 'order_status_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	startSoDate, isStartSoDate := ctx.GetQuery("start_so_date")
	if isStartSoDate == false {
		startSoDate = ""
	}

	endSoDate, isEndSoDate := ctx.GetQuery("end_so_date")
	if isEndSoDate == false {
		endSoDate = ""
	}

	id, isIdExist := ctx.GetQuery("id")
	if isIdExist == false {
		id = "0"
	}

	idInt, err = strconv.Atoi(id)
	if err != nil {
		err = helper.NewError("Parameter 'id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	productId, isProductIdExist := ctx.GetQuery("product_id")
	if isProductIdExist == false {
		productId = "0"
	}

	productIdInt, err = strconv.Atoi(productId)
	if err != nil {
		err = helper.NewError("Parameter 'product_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	categoryId, isCategoryIdExist := ctx.GetQuery("category_id")
	if isCategoryIdExist == false {
		categoryId = "0"
	}

	categoryIdInt, err = strconv.Atoi(categoryId)
	if err != nil {
		err = helper.NewError("Parameter 'category_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	salesmanId, isSalesmanIdExist := ctx.GetQuery("salesman_id")
	if isSalesmanIdExist == false {
		salesmanId = "0"
	}

	salesmanIdInt, err = strconv.Atoi(salesmanId)
	if err != nil {
		err = helper.NewError("Parameter 'salesman_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	provinceId, isProvinceIdExist := ctx.GetQuery("province_id")
	if isProvinceIdExist == false {
		provinceId = "0"
	}

	provinceIdInt, err = strconv.Atoi(provinceId)
	if err != nil {
		err = helper.NewError("Parameter 'province_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	cityId, isCityIdExist := ctx.GetQuery("city_id")
	if isCityIdExist == false {
		cityId = "0"
	}

	cityIdInt, err = strconv.Atoi(cityId)
	if err != nil {
		err = helper.NewError("Parameter 'city_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	districtId, isDistrictIdExist := ctx.GetQuery("district_id")
	if isDistrictIdExist == false {
		districtId = "0"
	}

	districtIdInt, err = strconv.Atoi(districtId)
	if err != nil {
		err = helper.NewError("Parameter 'district_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	villageId, isVillageIdExist := ctx.GetQuery("village_id")
	if isVillageIdExist == false {
		villageId = "0"
	}

	villageIdInt, err = strconv.Atoi(villageId)
	if err != nil {
		err = helper.NewError("Parameter 'village_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	startCreatedAt, isStartCreatedAt := ctx.GetQuery("start_created_at")
	if isStartCreatedAt == false {
		startCreatedAt = ""
	}

	endCreatedAt, isEndCreatedAt := ctx.GetQuery("end_created_at")
	if isEndCreatedAt == false {
		endCreatedAt = ""
	}

	salesOrderRequest := &models.SalesOrderRequest{
		Page:              pageInt,
		PerPage:           perPageInt,
		SortField:         sortField,
		SortValue:         sortValue,
		GlobalSearchValue: globalSearchValue,
		AgentID:           agentIdInt,
		StoreID:           storeIdInt,
		BrandID:           brandIdInt,
		OrderSourceID:     orderSourceIdInt,
		OrderStatusID:     orderStatusIdInt,
		StartSoDate:       startSoDate,
		EndSoDate:         endSoDate,
		ID:                idInt,
		ProductID:         productIdInt,
		CategoryID:        categoryIdInt,
		SalesmanID:        salesmanIdInt,
		ProvinceID:        provinceIdInt,
		CityID:            cityIdInt,
		DistrictID:        districtIdInt,
		VillageID:         villageIdInt,
		StartCreatedAt:    startCreatedAt,
		EndCreatedAt:      endCreatedAt,
	}

	salesOrders, errorLog := c.salesOrderUseCase.GetDetails(salesOrderRequest)

	if errorLog.Err != nil {
		resultErrorLog = errorLog
		result.StatusCode = resultErrorLog.StatusCode
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = salesOrders.SalesOrderDetails
	result.Total = salesOrders.Total
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}

func (c *salesOrderController) DeleteByID(ctx *gin.Context) {
	var result baseModel.Response
	var id int

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	sId := ctx.Param("so-id")
	id, err := strconv.Atoi(sId)

	if err != nil {
		err = helper.NewError("Parameter 'id' harus bernilai integer")
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, result.StatusCode, nil)
		ctx.JSON(result.StatusCode, result)
		return
	}
	mustActiveField := []*models.MustActiveRequest{
		{
			Table:    "sales_orders",
			ReqField: "id",
			Clause:   fmt.Sprintf("id = %d AND deleted_at IS NULL", id),
		},
	}

	err = c.requestValidationMiddleware.MustActiveValidation(ctx, mustActiveField)
	if err != nil {
		return
	}
	mustEmpties := []*models.MustEmptyValidationRequest{
		{
			Table:           "sales_orders",
			TableJoin:       "order_statuses",
			ForeignKey:      "order_status_id",
			SelectedCollumn: "order_statuses.name",
			Clause:          fmt.Sprintf("sales_orders.id = %d AND sales_orders.order_status_id NOT IN (5,6,9,10)", id),
			MessageFormat:   "Status Sales Order <result>",
		},
		{
			Table:           "delivery_orders",
			TableJoin:       "sales_orders",
			ForeignKey:      "sales_order_id",
			SelectedCollumn: "delivery_orders.id",
			Clause:          fmt.Sprintf("sales_orders.id = %d AND delivery_orders.deleted_at IS NULL", id),
			MessageFormat:   "Sales Order Has Delivery Order <result>, Please Delete it First",
		},
	}
	err = c.requestValidationMiddleware.MustEmptyValidation(ctx, mustEmpties)
	if err != nil {
		return
	}

	dbTransaction, err := c.db.BeginTx(ctx, nil)

	if err != nil {
		result.StatusCode = http.StatusInternalServerError
		result.Error = helper.WriteLog(err, result.StatusCode, nil)
		ctx.JSON(result.StatusCode, result)
		return
	}
	errorLog := c.salesOrderUseCase.DeleteById(id, dbTransaction)

	if errorLog != nil {
		err = dbTransaction.Rollback()

		if err != nil {
			result.StatusCode = http.StatusInternalServerError
			result.Error = helper.WriteLog(err, result.StatusCode, nil)
			ctx.JSON(result.StatusCode, result)
			return
		}

		result.StatusCode = errorLog.StatusCode
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	err = dbTransaction.Commit()

	if err != nil {
		result.StatusCode = http.StatusInternalServerError
		result.Error = helper.WriteLog(err, result.StatusCode, "Ada kesalahan, silahkan coba lagi nanti")
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}

func (c *salesOrderController) GetDetailsBySoId(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	var pageInt, perPageInt, agentIdInt, storeIdInt, brandIdInt, orderSourceIdInt, orderStatusIdInt, productIdInt, categoryIdInt, salesmanIdInt, provinceIdInt, cityIdInt, districtIdInt, villageIdInt int

	soIds := ctx.Param("so-id")
	soId, err := strconv.Atoi(soIds)

	if err != nil {
		err = helper.NewError("Parameter 'so-id' harus bernilai integer")
		resultErrorLog.Message = err.Error()
		result.StatusCode = http.StatusBadRequest
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	page, isPageExist := ctx.GetQuery("page")
	if isPageExist == false {
		page = "1"
	}

	pageInt, err = strconv.Atoi(page)

	if err != nil {
		err = helper.NewError("Parameter 'page' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	if pageInt == 0 {
		err = helper.NewError("Parameter 'page' harus bernilai integer > 0")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	perPage, isPerPageExist := ctx.GetQuery("per_page")
	if isPerPageExist == false {
		perPage = "10"
	}

	perPageInt, err = strconv.Atoi(perPage)
	if err != nil {
		err = helper.NewError("Parameter 'per_page' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	if perPageInt == 0 {
		err = helper.NewError("Parameter 'per_page' harus bernilai integer > 0")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	sortField, isSortFieldExist := ctx.GetQuery("sort_field")
	if isSortFieldExist == false {
		sortField = "created_at"
	}

	if sortField != "order_status_id" && sortField != "so_date" && sortField != "so_ref_code" && sortField != "so_code" && sortField != "store_code" && sortField != "store_name" && sortField != "created_at" && sortField != "updated_at" {
		err = helper.NewError("Parameter 'sort_field' harus bernilai 'order_status_id' or 'so_date' or 'so_ref_code' or 'so_code' or 'store_code' or 'store_name' or 'created_at' or 'updated_at' ")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	sortValue, isSortValueExist := ctx.GetQuery("sort_value")
	if isSortValueExist == false {
		sortValue = "desc"
	}

	globalSearchValue, isGlobalSearchValueExist := ctx.GetQuery("global_search_value")
	if isGlobalSearchValueExist == false {
		globalSearchValue = ""
	}

	agentId, isAgentIdExist := ctx.GetQuery("agent_id")
	if isAgentIdExist == false {
		agentId = "0"
	}

	agentIdInt, err = strconv.Atoi(agentId)
	if err != nil {
		err = helper.NewError("Parameter 'agent_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	storeId, isStoreIdExist := ctx.GetQuery("store_id")
	if isStoreIdExist == false {
		storeId = "0"
	}

	storeIdInt, err = strconv.Atoi(storeId)
	if err != nil {
		err = helper.NewError("Parameter 'store_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	brandId, isBrandIdExist := ctx.GetQuery("brand_id")
	if isBrandIdExist == false {
		brandId = "0"
	}

	brandIdInt, err = strconv.Atoi(brandId)
	if err != nil {
		err = helper.NewError("Parameter 'brand_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	orderSourceId, isOrderSourceIdExist := ctx.GetQuery("order_source_id")
	if isOrderSourceIdExist == false {
		orderSourceId = "0"
	}

	orderSourceIdInt, err = strconv.Atoi(orderSourceId)
	if err != nil {
		err = helper.NewError("Parameter 'order_source_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	orderStatusId, isOrderStatusIdExist := ctx.GetQuery("order_status_id")
	if isOrderStatusIdExist == false {
		orderStatusId = "0"
	}

	orderStatusIdInt, err = strconv.Atoi(orderStatusId)
	if err != nil {
		err = helper.NewError("Parameter 'order_status_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	startSoDate, isStartSoDate := ctx.GetQuery("start_so_date")
	if isStartSoDate == false {
		startSoDate = ""
	}

	endSoDate, isEndSoDate := ctx.GetQuery("end_so_date")
	if isEndSoDate == false {
		endSoDate = ""
	}

	productId, isProductIdExist := ctx.GetQuery("product_id")
	if isProductIdExist == false {
		productId = "0"
	}

	productIdInt, err = strconv.Atoi(productId)
	if err != nil {
		err = helper.NewError("Parameter 'product_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	categoryId, isCategoryIdExist := ctx.GetQuery("category_id")
	if isCategoryIdExist == false {
		categoryId = "0"
	}

	categoryIdInt, err = strconv.Atoi(categoryId)
	if err != nil {
		err = helper.NewError("Parameter 'category_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	salesmanId, isSalesmanIdExist := ctx.GetQuery("salesman_id")
	if isSalesmanIdExist == false {
		salesmanId = "0"
	}

	salesmanIdInt, err = strconv.Atoi(salesmanId)
	if err != nil {
		err = helper.NewError("Parameter 'salesman_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	provinceId, isProvinceIdExist := ctx.GetQuery("province_id")
	if isProvinceIdExist == false {
		provinceId = "0"
	}

	provinceIdInt, err = strconv.Atoi(provinceId)
	if err != nil {
		err = helper.NewError("Parameter 'province_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	cityId, isCityIdExist := ctx.GetQuery("city_id")
	if isCityIdExist == false {
		cityId = "0"
	}

	cityIdInt, err = strconv.Atoi(cityId)
	if err != nil {
		err = helper.NewError("Parameter 'city_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	districtId, isDistrictIdExist := ctx.GetQuery("district_id")
	if isDistrictIdExist == false {
		districtId = "0"
	}

	districtIdInt, err = strconv.Atoi(districtId)
	if err != nil {
		err = helper.NewError("Parameter 'district_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	villageId, isVillageIdExist := ctx.GetQuery("village_id")
	if isVillageIdExist == false {
		villageId = "0"
	}

	villageIdInt, err = strconv.Atoi(villageId)
	if err != nil {
		err = helper.NewError("Parameter 'village_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	startCreatedAt, isStartCreatedAt := ctx.GetQuery("start_created_at")
	if isStartCreatedAt == false {
		startCreatedAt = ""
	}

	endCreatedAt, isEndCreatedAt := ctx.GetQuery("end_created_at")
	if isEndCreatedAt == false {
		endCreatedAt = ""
	}

	salesOrderRequest := &models.SalesOrderRequest{
		Page:              pageInt,
		PerPage:           perPageInt,
		SortField:         sortField,
		SortValue:         sortValue,
		GlobalSearchValue: globalSearchValue,
		AgentID:           agentIdInt,
		StoreID:           storeIdInt,
		BrandID:           brandIdInt,
		OrderSourceID:     orderSourceIdInt,
		OrderStatusID:     orderStatusIdInt,
		StartSoDate:       startSoDate,
		EndSoDate:         endSoDate,
		ID:                soId,
		ProductID:         productIdInt,
		CategoryID:        categoryIdInt,
		SalesmanID:        salesmanIdInt,
		ProvinceID:        provinceIdInt,
		CityID:            cityIdInt,
		DistrictID:        districtIdInt,
		VillageID:         villageIdInt,
		StartCreatedAt:    startCreatedAt,
		EndCreatedAt:      endCreatedAt,
	}

	salesOrders, errorLog := c.salesOrderUseCase.GetDetails(salesOrderRequest)

	if errorLog.Err != nil {
		resultErrorLog = errorLog
		result.StatusCode = resultErrorLog.StatusCode
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = salesOrders.SalesOrderDetails
	result.Total = salesOrders.Total
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}

func (c *salesOrderController) GetDetailsById(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog

	soDetailIds := ctx.Param("id")
	soDetailId, err := strconv.Atoi(soDetailIds)

	if err != nil {
		err = helper.NewError("Parameter 'so-id' harus bernilai integer")
		resultErrorLog.Message = err.Error()
		result.StatusCode = http.StatusBadRequest
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	data, errorLog := c.salesOrderUseCase.GetDetailById(soDetailId)

	if errorLog.Err != nil {
		resultErrorLog = errorLog
		result.StatusCode = resultErrorLog.StatusCode
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}
	result.Data = data
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}
