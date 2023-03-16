package controllers

import (
	"context"
	"encoding/json"
	"errors"
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

type DeliveryOrderControllerInterface interface {
	Create(ctx *gin.Context)
	UpdateByID(ctx *gin.Context)
	UpdateDeliveryOrderDetailByID(ctx *gin.Context)
	UpdateDeliveryOrderDetailByDeliveryOrderID(ctx *gin.Context)
	Get(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	GetDetailsByDoId(ctx *gin.Context)
	GetBySalesmanID(ctx *gin.Context)
	DeleteByID(ctx *gin.Context)
}

type deliveryOrderController struct {
	deliveryOrderUseCase        usecases.DeliveryOrderUseCaseInterface
	deliveryOrderValidator      usecases.DeliveryOrderValidatorInterface
	requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface
	db                          dbresolver.DB
	ctx                         context.Context
}

func InitDeliveryOrderController(deliveryOrderUseCase usecases.DeliveryOrderUseCaseInterface, deliveryOrderValidator usecases.DeliveryOrderValidatorInterface, requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, db dbresolver.DB, ctx context.Context) DeliveryOrderControllerInterface {
	return &deliveryOrderController{
		deliveryOrderUseCase:        deliveryOrderUseCase,
		deliveryOrderValidator:      deliveryOrderValidator,
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
	err = c.deliveryOrderValidator.CreateDeliveryOrderValidator(insertRequest, ctx)
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
	updateRequest := &models.DeliveryOrderUpdateByIDRequest{}

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	err := ctx.ShouldBindJSON(updateRequest)
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

	updateRequest.WarehouseID = 0
	updateRequest.OrderStatusID = 0
	updateRequest.OrderSourceID = 0
	updateRequest.DoRefCode = ""
	updateRequest.DoRefDate = ""
	updateRequest.DriverName = ""
	updateRequest.PlatNumber = ""

	err = c.deliveryOrderValidator.UpdateDeliveryOrderByIDValidator(intID, updateRequest, ctx)
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

	deliveryOrder, errorLog := c.deliveryOrderUseCase.UpdateByID(intID, updateRequest, dbTransaction, ctx)
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

func (c *deliveryOrderController) UpdateDeliveryOrderDetailByID(ctx *gin.Context) {
	id := ctx.Param("id")
	intID, _ := strconv.Atoi(id)

	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	updateRequest := &models.DeliveryOrderDetailUpdateByIDRequest{}

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	err := ctx.ShouldBindJSON(updateRequest)
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

	err = c.deliveryOrderValidator.UpdateDeliveryOrderDetailByIDValidator(intID, updateRequest, ctx)
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

	deliveryOrderDetail, errorLog := c.deliveryOrderUseCase.UpdateDODetailByID(intID, updateRequest, dbTransaction, ctx)
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

	result.Data = deliveryOrderDetail
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return
}

func (c *deliveryOrderController) UpdateDeliveryOrderDetailByDeliveryOrderID(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	updateRequest := []*models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest{}

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
	err = c.deliveryOrderValidator.UpdateDeliveryOrderDetailByDoIDValidator(intID, updateRequest, ctx)
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

	deliveryOrderDetail, errorLog := c.deliveryOrderUseCase.UpdateDoDetailByDeliveryOrderID(intID, updateRequest, dbTransaction, ctx)
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

	deliveryOrderRequest, err := c.deliveryOrderValidator.GetDeliveryOrderValidator(ctx)
	if err != nil {
		return
	}

	deliveryOrders, errorLog := c.deliveryOrderUseCase.Get(deliveryOrderRequest)

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

func (c *deliveryOrderController) GetDetailsByDoId(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog

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

	deliveryOrderRequest, err := c.deliveryOrderValidator.GetDeliveryOrderDetailValidator(ctx)
	if err != nil {
		return
	}

	deliveryOrderRequest.ID = id

	deliveryOrders, errorLog := c.deliveryOrderUseCase.GetByDoID(deliveryOrderRequest)

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

	deliveryOrderReqeuest, err := c.deliveryOrderValidator.GetDeliveryOrderBySalesmanIDValidator(ctx)
	if err != nil {
		return
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
	err := c.deliveryOrderValidator.DeleteDeliveryOrderByIDValidator(sId, ctx)
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

func (c *deliveryOrderController) DeleteDetailByID(ctx *gin.Context) {
	var result baseModel.Response
	var id int

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	sId := ctx.Param("id")
	err := c.deliveryOrderValidator.DeleteDeliveryOrderDetailByIDValidator(sId, ctx)
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

	errorLog := c.deliveryOrderUseCase.DeleteDetailByID(id, dbTransaction)

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
