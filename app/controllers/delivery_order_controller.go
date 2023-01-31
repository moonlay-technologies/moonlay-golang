package controllers

import (
	"context"
	"fmt"
	"github.com/bxcodec/dbresolver"
	"github.com/gin-gonic/gin"
	"net/http"
	"poc-order-service/app/models"
	"poc-order-service/app/usecases"
	"poc-order-service/global/utils/helper"
	baseModel "poc-order-service/global/utils/model"
	"strconv"
)

type DeliveryOrderControllerInterface interface {
	Create(ctx *gin.Context)
	Get(ctx *gin.Context)
	GetByID(ctx *gin.Context)
}

type deliveryOrderController struct {
	deliveryOrderUseCase usecases.DeliveryOrderUseCaseInterface
	db                   dbresolver.DB
	ctx                  context.Context
}

func InitDeliveryOrderController(deliveryOrderUseCase usecases.DeliveryOrderUseCaseInterface, db dbresolver.DB, ctx context.Context) DeliveryOrderControllerInterface {
	return &deliveryOrderController{
		deliveryOrderUseCase: deliveryOrderUseCase,
		db:                   db,
		ctx:                  ctx,
	}
}

func (c *deliveryOrderController) Create(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	insertRequest := &models.DeliveryOrderStoreRequest{}

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	err := ctx.BindJSON(insertRequest)

	if err != nil {
		fmt.Println("error")
		errorLog := helper.WriteLog(err, http.StatusBadRequest, "Ada kesalahan, silahkan coba lagi nanti")
		result.StatusCode = http.StatusBadRequest
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	dbTransaction, err := c.db.BeginTx(ctx, nil)

	if err != nil {
		errorLog := helper.WriteLog(err, http.StatusInternalServerError, "Ada kesalahan, silahkan coba lagi nanti")
		resultErrorLog = errorLog
		result.StatusCode = http.StatusInternalServerError
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	deliveryOrder, errorLog := c.deliveryOrderUseCase.Create(insertRequest, dbTransaction, ctx)
	fmt.Println(deliveryOrder)
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

	result.Data = deliveryOrder
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}

func (c *deliveryOrderController) Get(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	var pageInt, perPageInt int

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

	startCreatedAt, isStartCreatedAt := ctx.GetQuery("start_created_at")
	if isStartCreatedAt == false {
		startCreatedAt = ""
	}

	endCreatedAt, isEndCreatedAt := ctx.GetQuery("end_created_at")
	if isEndCreatedAt == false {
		endCreatedAt = ""
	}

	startSoDate, isStartSoDate := ctx.GetQuery("start_so_date")
	if isStartSoDate == false {
		startSoDate = ""
	}

	endSoDate, isEndSoDate := ctx.GetQuery("end_so_date")
	if isEndSoDate == false {
		endSoDate = ""
	}

	deliveryOrderReqeuest := &models.DeliveryOrderRequest{
		Page:           pageInt,
		PerPage:        perPageInt,
		StartCreatedAt: startCreatedAt,
		EndCreatedAt:   endCreatedAt,
		StartSoDate:    startSoDate,
		EndSoDate:      endSoDate,
		SortField:      sortField,
		SortValue:      sortValue,
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
		result.StatusCode = 400
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
