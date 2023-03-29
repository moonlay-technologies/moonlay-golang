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
	GetDeliveryOrders(ctx *gin.Context)
}

type hostToHostController struct {
	cartUseCase                 usecases.CartUseCaseInterface
	salesOrderUseCase           usecases.SalesOrderUseCaseInterface
	deliveryOrderUseCase        usecases.DeliveryOrderUseCaseInterface
	requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface
	db                          dbresolver.DB
	ctx                         context.Context
}

func InitHostToHostController(cartUseCase usecases.CartUseCaseInterface, salesOrderUseCase usecases.SalesOrderUseCaseInterface, deliveryOrderUseCase usecases.DeliveryOrderUseCaseInterface, requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, db dbresolver.DB, ctx context.Context) HostToHostControllerInterface {
	return &hostToHostController{
		cartUseCase:                 cartUseCase,
		salesOrderUseCase:           salesOrderUseCase,
		deliveryOrderUseCase:        deliveryOrderUseCase,
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
	if !isPageExist {
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
	if !isPerPageExist {
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
	if !isSortFieldExist {
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

func (c *hostToHostController) GetDeliveryOrders(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	var pageInt, perPageInt, intID, intSalesOrderID, intAgentID, intStoreID, intBrandID, intOrderStatusID, intProductID, intCategoryID, intSalesmanID, intProvinceID, intCityID, intDistrictID, intVillageID int

	page, isPageExist := ctx.GetQuery("page")
	if !isPageExist {
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

	perPageInt, err = strconv.Atoi(perPage)
	if err != nil {
		err = helper.NewError("Parameter 'per_page' harus bernilai integer")
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

	sortField, isSortFieldExist := ctx.GetQuery("sort_field")
	if isSortFieldExist == false {
		sortField = "created_at"
	}

	globalSearchValue, isGlobalSearchValueExist := ctx.GetQuery("global_search_value")
	if isGlobalSearchValueExist == false {
		globalSearchValue = ""
	}

	id, isIdExist := ctx.GetQuery("id")
	if isIdExist == false {
		id = "0"
	}

	intID, err = strconv.Atoi(id)
	if err != nil {
		err = helper.NewError("Parameter 'id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	salesOrderID, isSalesOrderIDExist := ctx.GetQuery("sales_order_id")
	if isSalesOrderIDExist == false {
		salesOrderID = "0"
	}

	intSalesOrderID, err = strconv.Atoi(salesOrderID)
	if err != nil {
		err = helper.NewError("Parameter 'sales_order_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	agentID, isAgentIDExist := ctx.GetQuery("agent_id")
	if isAgentIDExist == false {
		agentID = "0"
	}

	intAgentID, err = strconv.Atoi(agentID)
	if err != nil {
		err = helper.NewError("Parameter 'agent_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	storeID, isStoreIDExist := ctx.GetQuery("store_id")
	if isStoreIDExist == false {
		storeID = "0"
	}

	intStoreID, err = strconv.Atoi(storeID)
	if err != nil {
		err = helper.NewError("Parameter 'store_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	brandID, isBrandIDExist := ctx.GetQuery("brand_id")
	if isBrandIDExist == false {
		brandID = "0"
	}

	intBrandID, err = strconv.Atoi(brandID)
	if err != nil {
		err = helper.NewError("Parameter 'brand_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	orderStatusID, isOrderStatusIDExist := ctx.GetQuery("order_status_id")
	if isOrderStatusIDExist == false {
		orderStatusID = "0"
	}

	intOrderStatusID, err = strconv.Atoi(orderStatusID)
	if err != nil {
		err = helper.NewError("Parameter 'order_status_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	doCode, isDoCodeExist := ctx.GetQuery("do_code")
	if isDoCodeExist == false {
		doCode = ""
	}

	soCode, isSoCodeExist := ctx.GetQuery("so_code")
	if isSoCodeExist == false {
		soCode = ""
	}

	startDoDate, isStartSoDate := ctx.GetQuery("start_do_date")
	if isStartSoDate == false {
		startDoDate = ""
	}

	endDoDate, isEndSoDate := ctx.GetQuery("end_do_date")
	if isEndSoDate == false {
		endDoDate = ""
	}

	doRefCode, isDoRefCodeExist := ctx.GetQuery("do_ref_code")
	if isDoRefCodeExist == false {
		doRefCode = ""
	}

	doRefDate, isDoRefDateExist := ctx.GetQuery("do_ref_date")
	if isDoRefDateExist == false {
		doRefDate = ""
	}

	productID, isProductIDExist := ctx.GetQuery("product_id")
	if isProductIDExist == false {
		productID = "0"
	}

	intProductID, err = strconv.Atoi(productID)
	if err != nil {
		err = helper.NewError("Parameter 'product_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	categoryID, isCategoryIDExist := ctx.GetQuery("category_id")
	if isCategoryIDExist == false {
		categoryID = "0"
	}

	intCategoryID, err = strconv.Atoi(categoryID)
	if err != nil {
		err = helper.NewError("Parameter 'category_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	salesmanID, isSalesmanIDExist := ctx.GetQuery("salesman_id")
	if isSalesmanIDExist == false {
		salesmanID = "0"
	}

	intSalesmanID, err = strconv.Atoi(salesmanID)
	if err != nil {
		err = helper.NewError("Parameter 'salesman_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	provinceID, isProvinceIDExist := ctx.GetQuery("province_id")
	if isProvinceIDExist == false {
		provinceID = "0"
	}

	intProvinceID, err = strconv.Atoi(provinceID)
	if err != nil {
		err = helper.NewError("Parameter 'province_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	cityID, isCityIDExist := ctx.GetQuery("city_id")
	if isCityIDExist == false {
		cityID = "0"
	}

	intCityID, err = strconv.Atoi(cityID)
	if err != nil {
		err = helper.NewError("Parameter 'city_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	districtID, isDistrictIDExist := ctx.GetQuery("district_id")
	if isDistrictIDExist == false {
		districtID = "0"
	}

	intDistrictID, err = strconv.Atoi(districtID)
	if err != nil {
		err = helper.NewError("Parameter 'district_id' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	villageID, isVillageIDExist := ctx.GetQuery("village_id")
	if isVillageIDExist == false {
		villageID = "0"
	}

	intVillageID, err = strconv.Atoi(villageID)
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

	updatedAt, isUpdatedAtExist := ctx.GetQuery("updated_at")
	if isUpdatedAtExist == false {
		updatedAt = ""
	}

	deliveryOrderReqeuest := &models.DeliveryOrderRequest{
		Page:              pageInt,
		PerPage:           perPageInt,
		SortField:         sortField,
		SortValue:         sortValue,
		GlobalSearchValue: globalSearchValue,
		ID:                intID,
		SalesOrderID:      intSalesOrderID,
		AgentID:           intAgentID,
		StoreID:           intStoreID,
		BrandID:           intBrandID,
		OrderStatusID:     intOrderStatusID,
		DoCode:            doCode,
		SoCode:            soCode,
		StartDoDate:       startDoDate,
		EndDoDate:         endDoDate,
		DoRefCode:         doRefCode,
		DoRefDate:         doRefDate,
		ProductID:         intProductID,
		CategoryID:        intCategoryID,
		SalesmanID:        intSalesmanID,
		ProvinceID:        intProvinceID,
		CityID:            intCityID,
		DistrictID:        intDistrictID,
		VillageID:         intVillageID,
		StartCreatedAt:    startCreatedAt,
		EndCreatedAt:      endCreatedAt,
		UpdatedAt:         updatedAt,
	}

	deliveryOrders, errorLog := c.deliveryOrderUseCase.Get(deliveryOrderReqeuest)

	if errorLog.Err != nil {
		resultErrorLog = errorLog
		result.StatusCode = resultErrorLog.StatusCode
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = deliveryOrders.DeliveryOrders
	result.Total = deliveryOrders.Total
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}
