package controllers

import (
	"context"
	"net/http"
	"order-service/app/middlewares"
	"order-service/app/models"
	"order-service/app/usecases"
	"order-service/global/utils/helper"
	baseModel "order-service/global/utils/model"
	"strconv"

	"github.com/bxcodec/dbresolver"
	"github.com/gin-gonic/gin"
)

type HostToHostControllerInterface interface {
	GetSalesOrders(ctx *gin.Context)
}

type hostToHostController struct {
	cartUseCase                 usecases.CartUseCaseInterface
	salesOrderUseCase           usecases.SalesOrderUseCaseInterface
	requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface
	db                          dbresolver.DB
	ctx                         context.Context
}

func InitHostToHostController(cartUseCase usecases.CartUseCaseInterface, salesOrderUseCase usecases.SalesOrderUseCaseInterface, requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, db dbresolver.DB, ctx context.Context) HostToHostControllerInterface {
	return &hostToHostController{
		cartUseCase:                 cartUseCase,
		salesOrderUseCase:           salesOrderUseCase,
		requestValidationMiddleware: requestValidationMiddleware,
		db:                          db,
		ctx:                         ctx,
	}
}

func (c *hostToHostController) GetSalesOrders(ctx *gin.Context) {
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
