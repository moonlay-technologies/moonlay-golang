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
	GetBySalesmanID(ctx *gin.Context)
	DeleteByID(ctx *gin.Context)
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

	dateField := []*models.DateInputRequest{
		{
			Field: "so_date",
			Value: insertRequest.DoRefDate,
		},
	}
	err = c.requestValidationMiddleware.DateInputValidation(ctx, dateField, constants.ERROR_ACTION_NAME_CREATE)
	if err != nil {
		return
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
	mustActiveField417 := []*models.MustActiveRequest{
		{
			Table:    "sales_orders",
			ReqField: "sales_order_id",
			Clause:   fmt.Sprintf("id = %d AND deleted_at IS NULL", insertRequest.SalesOrderID),
		},
		{
			Table:    "warehouses a JOIN agents b JOIN sales_orders c ON a.owner_id = b.id AND c.agent_id = b.id",
			ReqField: "warehouse_owner_id",
			Clause:   fmt.Sprintf("c.id = %d AND b.deleted_at IS NULL AND a.`status` = 1", insertRequest.SalesOrderID),
		},
		{
			Table:    "stores a JOIN sales_orders b ON b.store_id = a.id",
			ReqField: "stores_id",
			Clause:   fmt.Sprintf("b.id = %d AND a.deleted_at IS NULL AND b.deleted_at IS NULL", insertRequest.SalesOrderID),
		},
		{
			Table:    "brands a JOIN sales_orders b ON b.brand_id = a.id",
			ReqField: "brands_id",
			Clause:   fmt.Sprintf("b.id = %d AND a.status_active = 1 AND b.deleted_at IS NULL", insertRequest.SalesOrderID),
		},
	}

	sDoDate := "NOW()"
	sDoDateEqualMonth := fmt.Sprintf("MONTH(so_date) = MONTH('%s') AND MONTH(so_date) = MONTH(%s)", insertRequest.DoRefDate, sDoDate)
	sDoDateHigherOrEqualSoDate := fmt.Sprintf("DAY(so_date) <= DAY('%s') AND DAY(so_date) <= DAY(%s)", insertRequest.DoRefDate, sDoDate)
	// sDoDateLowerOrEqualSoRefDate := fmt.Sprintf("DAY(so_ref_date) >= DAY('%s') AND DAY(so_ref_date) >= DAY(%s)", insertRequest.DoRefDate, sDoDate)
	sDoDateLowerOrEqualToday := fmt.Sprintf("DAY('%s') <= DAY(%s)", insertRequest.DoRefDate, sDoDate)
	sSoDateEqualDoDate := fmt.Sprintf("IF(DAY(so_date) = DAY(%[2]s), IF(DAY(%[2]s) = DAY('%[1]s'), TRUE, FALSE), TRUE)", insertRequest.DoRefDate, sDoDate)
	mustActiveField422 := []*models.MustActiveRequest{
		{
			Table:         "sales_orders",
			ReqField:      "so_date",
			Clause:        fmt.Sprintf("id = %d AND %s AND %s AND %s AND %s ", insertRequest.SalesOrderID, sDoDateEqualMonth, sDoDateHigherOrEqualSoDate, sDoDateLowerOrEqualToday, sSoDateEqualDoDate),
			CustomMessage: "do_date and do_ref_date must be equal less than today, must be equal more than so_date and must be in the current month",
		},
	}
	mustEmpties := []*models.MustEmptyValidationRequest{
		{
			Table:           "sales_orders",
			TableJoin:       "order_statuses",
			ForeignKey:      "order_status_id",
			SelectedCollumn: "order_statuses.name",
			Clause:          fmt.Sprintf("sales_orders.id = %d AND sales_orders.order_status_id NOT IN (5,7)", insertRequest.SalesOrderID),
			MessageFormat:   "Status Sales Order <result>",
		},
	}
	totalQty := 0
	for _, x := range insertRequest.DeliveryOrderDetails {
		if x.Qty < 0 {
			errorLog := helper.WriteLog(err, http.StatusBadRequest, fmt.Sprintf("qty sales order detail %d must equal or higher than 0", x.SoDetailID))
			resultErrorLog = errorLog
			result.StatusCode = http.StatusUnprocessableEntity
			result.Error = resultErrorLog
			ctx.JSON(result.StatusCode, result)
			return
		}

		totalQty += x.Qty
		mustActiveField417 = append(mustActiveField417, &models.MustActiveRequest{
			Table:         "sales_orders a JOIN sales_order_details b ON b.sales_order_id = a.id",
			ReqField:      "sales_order_id",
			Clause:        fmt.Sprintf("b.id = %d AND a.id = %d AND a.deleted_at IS NULL AND b.deleted_at IS NULL", x.SoDetailID, insertRequest.SalesOrderID),
			CustomMessage: fmt.Sprintf("sales order detail = %d tidak ditemukan", x.SoDetailID),
		})
		mustActiveField417 = append(mustActiveField417, &models.MustActiveRequest{
			Table:    "products a JOIN sales_order_details b ON b.product_id = a.id",
			ReqField: "product_id",
			Clause:   fmt.Sprintf("b.id = %d AND a.deleted_at IS NULL AND b.deleted_at IS NULL", x.SoDetailID),
		})
		mustActiveField417 = append(mustActiveField417, &models.MustActiveRequest{
			Table:    "uoms a JOIN sales_order_details b ON b.uom_id = a.id",
			ReqField: "uoms_id",
			Clause:   fmt.Sprintf("b.id = %d AND a.deleted_at IS NULL AND b.deleted_at IS NULL", x.SoDetailID),
		})
		mustActiveField422 = append(mustActiveField422, &models.MustActiveRequest{
			Table:         "sales_order_details",
			ReqField:      "residual_qty",
			Clause:        fmt.Sprintf("id = %d AND deleted_at IS NULL AND residual_qty >= %d", x.SoDetailID, x.Qty),
			CustomMessage: fmt.Sprintf("Residual Qty SO Detail %d must be higher than or equal delivery order qty", x.SoDetailID),
		})
		mustEmpties = append(mustEmpties, &models.MustEmptyValidationRequest{
			Table:           "sales_order_details",
			TableJoin:       "order_statuses",
			ForeignKey:      "order_status_id",
			SelectedCollumn: "order_statuses.name",
			Clause:          fmt.Sprintf("sales_order_details.id = %d AND sales_order_details.order_status_id NOT IN (11, 13)", x.SoDetailID),
			MessageFormat:   fmt.Sprintf("Status Sales Order Detail %d <result>", x.SoDetailID),
		})
	}

	err = c.requestValidationMiddleware.MustActiveValidation(ctx, mustActiveField417)
	if err != nil {
		return
	}

	err = c.requestValidationMiddleware.MustEmptyValidation(ctx, mustEmpties)
	if err != nil {
		return
	}

	err = c.requestValidationMiddleware.MustActiveValidation422(ctx, mustActiveField422)
	if err != nil {
		return
	}

	if totalQty <= 0 {
		errorLog := helper.WriteLog(err, http.StatusBadRequest, "total qty must higher than 0")
		resultErrorLog = errorLog
		result.StatusCode = http.StatusUnprocessableEntity
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
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

	result.Data = deliveryOrder
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
	var pageInt, perPageInt, intID, intSalesOrderID, intAgentID, intStoreID, intBrandID, intOrderStatusID, intProductID, intCategoryID, intSalesmanID, intProvinceID, intCityID, intDistrictID, intVillageID int

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

	if sortField != "order_status_id" && sortField != "do_date" && sortField != "do_ref_code" && sortField != "created_at" && sortField != "updated_at" {
		err = helper.NewError("Parameter 'sort_field' harus bernilai 'order_status_id' or 'do_date' or 'do_ref_code' or 'created_at' or 'updated_at'")
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, err.Error())
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLogData
		ctx.JSON(result.StatusCode, result)
		return
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

func (c *deliveryOrderController) GetBySalesmanID(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	var pageInt, perPageInt, intID, intSalesOrderID, intAgentID, intStoreID, intBrandID, intOrderStatusID, intProductID, intCategoryID, intSalesmanID, intProvinceID, intCityID, intDistrictID, intVillageID int

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

	deliveryOrders, errorLog := c.deliveryOrderUseCase.GetBySalesmansID(deliveryOrderReqeuest)

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

	deliveryOrder, errorLog := c.deliveryOrderUseCase.GetByIDWithDetail(deliveryOrderRequest, ctx)

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
func (c *deliveryOrderController) DeleteByID(ctx *gin.Context) {
	var result baseModel.Response
	var id int

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	sId := ctx.Param("id")
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
			Table:    "delivery_orders",
			ReqField: "id",
			Clause:   fmt.Sprintf("id = %d AND deleted_at IS NULL", id),
		},
	}

	err = c.requestValidationMiddleware.MustActiveValidation(ctx, mustActiveField)
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

	errorLog := c.deliveryOrderUseCase.DeleteByID(id, dbTransaction)

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

	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
}
