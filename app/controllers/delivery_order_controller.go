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

type DeliveryOrderControllerInterface interface {
	Create(ctx *gin.Context)
	UpdateByID(ctx *gin.Context)
	UpdateDeliveryOrderDetailByID(ctx *gin.Context)
	UpdateDeliveryOrderDetailByDeliveryOrderID(ctx *gin.Context)
	Get(ctx *gin.Context)
	GetByID(ctx *gin.Context)
}

type deliveryOrderController struct {
	deliveryOrderUseCase        usecases.DeliveryOrderUseCaseInterface
	requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface
	db                          dbresolver.DB
	ctx                         context.Context
}

func InitDeliveryOrderController(deliveryOrderUseCase usecases.DeliveryOrderUseCaseInterface, requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, db dbresolver.DB, ctx context.Context) DeliveryOrderControllerInterface {
	return &deliveryOrderController{
		deliveryOrderUseCase:        deliveryOrderUseCase,
		requestValidationMiddleware: requestValidationMiddleware,
		db:                          db,
		ctx:                         ctx,
	}
}

func (c *deliveryOrderController) Create(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	insertRequest := &models.DeliveryOrderStoreRequest{}

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	err := ctx.ShouldBindJSON(insertRequest)
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
			Clause:   fmt.Sprintf("id = %d AND status = '%s'", insertRequest.AgentID, "active"),
		},
		{
			Table:    "stores",
			ReqField: "store_id",
			Clause:   fmt.Sprintf("id = %d AND status = '%s'", insertRequest.StoreID, "active"),
		},
	}

	for i, v := range insertRequest.DeliveryOrderDetails {
		mustActiveField = append(mustActiveField, &models.MustActiveRequest{
			Table:    "brands",
			ReqField: fmt.Sprintf("delivery_order_details[%d].brand_id", i),
			Clause:   fmt.Sprintf("id = %d AND status_active = %d", v.BrandID, 1),
		})
		mustActiveField = append(mustActiveField, &models.MustActiveRequest{
			Table:    "products",
			ReqField: fmt.Sprintf("delivery_order_details[%d].product_id", i),
			Clause:   fmt.Sprintf("id = %d AND isActive = %d", v.ProductID, 1),
		})
		mustActiveField = append(mustActiveField, &models.MustActiveRequest{
			Table:    "uoms",
			ReqField: fmt.Sprintf("delivery_order_details[%d].uom_id", i),
			Clause:   fmt.Sprintf("id = %d AND deleted_at IS NULL", v.UomID),
		})
	}

	err = c.requestValidationMiddleware.MustActiveValidation(ctx, mustActiveField)
	if err != nil {
		fmt.Println("Error active validation", err)
		return
	}

	uniqueField := []*models.UniqueRequest{
		{
			Table: constants.DELIVERY_ORDERS_TABLE,
			Field: "do_code",
			Value: insertRequest.DoCode,
		},
		{
			Table: constants.DELIVERY_ORDERS_TABLE,
			Field: "do_ref_code",
			Value: insertRequest.DoRefCode,
		},
	}

	err = c.requestValidationMiddleware.UniqueValidation(ctx, uniqueField)
	if err != nil {
		fmt.Println("Error unique validation", err)
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

	deliveryOrder, errorLog := c.deliveryOrderUseCase.Create(insertRequest, dbTransaction, ctx)
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

	deliveryOrderDetailResults := []*models.DeliveryOrderDetailStoreResponse{}
	for _, v := range deliveryOrder.DeliveryOrderDetails {
		deliveryOrderDetailResult := models.DeliveryOrderDetailStoreResponse{
			DeliveryOrderID: v.DeliveryOrderID,
			SoDetailID:      v.SoDetailID,
			ProductSku:      v.ProductSKU,
			ProductName:     v.ProductName,
			UomCode:         v.Uom.Code.String,
			Qty:             v.Qty,
			ResidualQty:     v.SoDetail.ResidualQty,
			Note:            v.Note.String,
		}
		deliveryOrderDetailResults = append(deliveryOrderDetailResults, &deliveryOrderDetailResult)
	}

	deliveryOrderResult := &models.DeliveryOrderStoreResponse{
		SalesOrderID:              deliveryOrder.SalesOrderID,
		SalesOrderSoCode:          deliveryOrder.SalesOrder.SoCode,
		SalesOrderSoDate:          deliveryOrder.SalesOrder.SoDate,
		SalesOrderReferralCode:    deliveryOrder.SalesOrder.SoRefCode.String,
		SalesOrderNote:            deliveryOrder.SalesOrder.Note.String,
		SalesOrderInternalComment: deliveryOrder.SalesOrder.InternalComment.String,
		SalesmanName:              deliveryOrder.Salesman.Name,
		StoreName:                 deliveryOrder.Store.Name.String,
		StoreCityName:             deliveryOrder.Store.Name.String,
		StoreProvinceName:         deliveryOrder.Store.ProvinceName.String,
		TotalAmount:               int(deliveryOrder.SalesOrder.TotalAmount),
		WarehouseID:               deliveryOrder.WarehouseID,
		WarehouseAddress:          deliveryOrder.Warehouse.Address.String,
		OrderSourceID:             deliveryOrder.OrderSourceID,
		OrderStatusID:             deliveryOrder.OrderStatusID,
		AgentID:                   deliveryOrder.AgentID,
		StoreID:                   deliveryOrder.StoreID,
		DoCode:                    deliveryOrder.DoCode,
		DoDate:                    deliveryOrder.DoDate,
		DoRefCode:                 deliveryOrder.DoRefCode.String,
		DoRefDate:                 deliveryOrder.DoRefDate.String,
		DriverName:                deliveryOrder.DriverName.String,
		PlatNumber:                deliveryOrder.PlatNumber.String,
		Note:                      deliveryOrder.Note.String,
		DeliveryOrderDetails:      deliveryOrderDetailResults,
	}

	result.Data = deliveryOrderResult
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}

func (c *deliveryOrderController) UpdateByID(ctx *gin.Context) {
	id := ctx.Param("id")
	intID, _ := strconv.Atoi(id)

	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	insertRequest := &models.DeliveryOrderUpdateByIDRequest{}

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	err := ctx.ShouldBindJSON(insertRequest)
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

	uniqueField := []*models.UniqueRequest{
		{
			Table: constants.DELIVERY_ORDERS_TABLE,
			Field: "do_ref_code",
			Value: insertRequest.DoRefCode,
		},
	}

	err = c.requestValidationMiddleware.UniqueValidation(ctx, uniqueField)
	if err != nil {
		fmt.Println("Error unique validation", err)
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

	deliveryOrder, errorLog := c.deliveryOrderUseCase.UpdateByID(intID, insertRequest, dbTransaction, ctx)
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

	deliveryOrderDetailResults := []*models.DeliveryOrderDetailUpdateByIDRequest{}
	for _, v := range deliveryOrder.DeliveryOrderDetails {
		deliveryOrderDetailResult := models.DeliveryOrderDetailUpdateByIDRequest{
			Qty:  v.Qty,
			Note: v.Note.String,
		}
		deliveryOrderDetailResults = append(deliveryOrderDetailResults, &deliveryOrderDetailResult)
	}

	deliveryOrderResult := &models.DeliveryOrderUpdateByIDRequest{
		WarehouseID:          deliveryOrder.WarehouseID,
		OrderSourceID:        deliveryOrder.OrderSourceID,
		OrderStatusID:        deliveryOrder.OrderStatusID,
		DoRefCode:            deliveryOrder.DoRefCode.String,
		DoRefDate:            deliveryOrder.DoRefDate.String,
		DriverName:           deliveryOrder.DriverName.String,
		PlatNumber:           deliveryOrder.PlatNumber.String,
		Note:                 deliveryOrder.Note.String,
		DeliveryOrderDetails: deliveryOrderDetailResults,
	}

	result.Data = deliveryOrderResult
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}

func (c *deliveryOrderController) UpdateDeliveryOrderDetailByID(ctx *gin.Context) {
	id := ctx.Param("id")
	intID, _ := strconv.Atoi(id)

	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	insertRequest := &models.DeliveryOrderDetailUpdateByIDRequest{}

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	err := ctx.ShouldBindJSON(insertRequest)
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

	dbTransaction, err := c.db.BeginTx(ctx, nil)

	if err != nil {
		errorLog := helper.WriteLog(err, http.StatusInternalServerError, nil)
		resultErrorLog = errorLog
		result.StatusCode = http.StatusInternalServerError
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	deliveryOrderDetail, errorLog := c.deliveryOrderUseCase.UpdateDODetailByID(intID, insertRequest, dbTransaction, ctx)
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

	deliveryOrderDetailResult := models.DeliveryOrderDetailUpdateByIDRequest{
		Qty:  deliveryOrderDetail.Qty,
		Note: deliveryOrderDetail.Note.String,
	}

	result.Data = deliveryOrderDetailResult
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}

func (c *deliveryOrderController) UpdateDeliveryOrderDetailByDeliveryOrderID(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	insertRequest := []*models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest{}

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	id := ctx.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		err = helper.NewError("Parameter 'id' harus bernilai integer")
		resultErrorLog.Message = err.Error()
		result.StatusCode = http.StatusBadRequest
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	err = ctx.BindJSON(&insertRequest)

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

	dbTransaction, err := c.db.BeginTx(ctx, nil)

	if err != nil {
		errorLog := helper.WriteLog(err, http.StatusInternalServerError, nil)
		resultErrorLog = errorLog
		result.StatusCode = http.StatusInternalServerError
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	deliveryOrderDetail, errorLog := c.deliveryOrderUseCase.UpdateDoDetailByDeliveryOrderID(intID, insertRequest, dbTransaction, ctx)
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

	deliveryOrderResults := []models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest{}
	for _, v := range deliveryOrderDetail.DeliveryOrderDetails {
		deliveryOrderResult := models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest{
			ID:   v.ID,
			Qty:  v.Qty,
			Note: v.Note.String,
		}
		deliveryOrderResults = append(deliveryOrderResults, deliveryOrderResult)
	}

	result.Data = deliveryOrderResults
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}

func (c *deliveryOrderController) Get(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	var pageInt, perPageInt, intAgentID, intStoreID, intBrandID, intOrderSourceID, intOrderStatusID, intCategoryID, intSalesmanID, intProvinceID, intCityID, intDistrictID, intVillageID int
	var floatTotalAmount, floatTotalTonase float64

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

	agentName, isAgentNameExist := ctx.GetQuery("agent_name")
	if isAgentNameExist == false {
		agentName = ""
	}

	storeCode, isStoreCodeExist := ctx.GetQuery("store_code")
	if isStoreCodeExist == false {
		storeCode = ""
	}

	storeName, isStoreNameExist := ctx.GetQuery("store_name")
	if isStoreNameExist == false {
		storeName = ""
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

	brandName, isBrandNameExist := ctx.GetQuery("brand_name")
	if isBrandNameExist == false {
		brandName = ""
	}

	orderSourceID, isOrderSourceIDExist := ctx.GetQuery("order_source_id")
	if isOrderSourceIDExist == false {
		orderSourceID = "0"
	}

	intOrderSourceID, err = strconv.Atoi(orderSourceID)
	if err != nil {
		err = helper.NewError("Parameter 'order_source_id' harus bernilai integer")
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

	startSoDate, isStartSoDate := ctx.GetQuery("start_so_date")
	if isStartSoDate == false {
		startSoDate = ""
	}

	endSoDate, isEndSoDate := ctx.GetQuery("end_so_date")
	if isEndSoDate == false {
		endSoDate = ""
	}

	doRefCode, isDoRefCodeExist := ctx.GetQuery("do_ref_code")
	if isDoRefCodeExist == false {
		doRefCode = ""
	}

	doRefDate, isDoRefDateExist := ctx.GetQuery("do_ref_date")
	if isDoRefDateExist == false {
		doRefDate = ""
	}

	doRefferalCode, isDoRefferalCodeExist := ctx.GetQuery("do_refferal_code")
	if isDoRefferalCodeExist == false {
		doRefferalCode = ""
	}

	totalAmount, isTotalAmountExist := ctx.GetQuery("total_amount")
	if isTotalAmountExist == false {
		totalAmount = "0"
	}

	floatTotalAmount, err = strconv.ParseFloat(totalAmount, 64)
	if err != nil {
		err = helper.NewError("Parameter 'total_amount' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	totalTonase, isTotalTonaseExist := ctx.GetQuery("total_tonase")
	if isTotalTonaseExist == false {
		totalTonase = "0"
	}

	floatTotalTonase, err = strconv.ParseFloat(totalTonase, 64)
	if err != nil {
		err = helper.NewError("Parameter 'total_tonase' harus bernilai integer")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
	}

	productSku, isProductSkuExist := ctx.GetQuery("product_sku")
	if isProductSkuExist == false {
		productSku = ""
	}

	productName, isProductNameExist := ctx.GetQuery("product_name")
	if isProductNameExist == false {
		productName = ""
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

	deliveryOrderReqeuest := &models.DeliveryOrderRequest{
		Page:              pageInt,
		PerPage:           perPageInt,
		SortField:         sortField,
		SortValue:         sortValue,
		GlobalSearchValue: globalSearchValue,
		AgentID:           intAgentID,
		StoreID:           intStoreID,
		AgentName:         agentName,
		StoreCode:         storeCode,
		StoreName:         storeName,
		BrandID:           intBrandID,
		BrandName:         brandName,
		OrderSourceID:     intOrderSourceID,
		OrderStatusID:     intOrderStatusID,
		DoCode:            doCode,
		StartSoDate:       startSoDate,
		EndSoDate:         endSoDate,
		DoRefCode:         doRefCode,
		DoRefDate:         doRefDate,
		DoRefferalCode:    doRefferalCode,
		TotalAmount:       floatTotalAmount,
		TotalTonase:       floatTotalTonase,
		ProductSKU:        productSku,
		ProductName:       productName,
		CategoryID:        intCategoryID,
		SalesmanID:        intSalesmanID,
		ProvinceID:        intProvinceID,
		CityID:            intCityID,
		DistrictID:        intDistrictID,
		VillageID:         intVillageID,
		StartCreatedAt:    startCreatedAt,
		EndCreatedAt:      endCreatedAt,
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

func (c *deliveryOrderController) GetByID(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	var id int

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

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

	deliveryOrderRequest := &models.DeliveryOrderRequest{
		ID: id,
	}

	deliveryOrder, errorLog := c.deliveryOrderUseCase.GetByID(deliveryOrderRequest, ctx)

	if errorLog.Err != nil {
		resultErrorLog = errorLog
		result.StatusCode = resultErrorLog.StatusCode
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = deliveryOrder
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}
