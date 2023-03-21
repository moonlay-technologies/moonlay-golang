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
	"order-service/global/utils/model"
	"strconv"

	"github.com/bxcodec/dbresolver"
	"github.com/gin-gonic/gin"
)

type DeliveryOrderControllerInterface interface {
	Create(ctx *gin.Context)
	UpdateByID(ctx *gin.Context)
	UpdateDeliveryOrderDetailByID(ctx *gin.Context)
	UpdateDeliveryOrderDetailByDeliveryOrderID(ctx *gin.Context)
	DeleteByID(ctx *gin.Context)
	DeleteDetailByID(ctx *gin.Context)
	DeleteDetailByDoID(ctx *gin.Context)

	Get(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	GetDetails(ctx *gin.Context)
	GetDetailsByDoId(ctx *gin.Context)
	GetDetailById(ctx *gin.Context)
	GetBySalesmanID(ctx *gin.Context)
	GetDOJourneysByDoID(ctx *gin.Context)
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
		ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
		return
	}

	deliveryOrder, errorLog := c.deliveryOrderUseCase.Create(insertRequest, dbTransaction, ctx)
	if errorLog != nil {
		err = dbTransaction.Rollback()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
			return
		}

		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	err = dbTransaction.Commit()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, model.Response{Data: deliveryOrder, StatusCode: http.StatusOK})
	return
}

func (c *deliveryOrderController) UpdateByID(ctx *gin.Context) {
	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	id := ctx.Param("id")
	intID, err := strconv.Atoi(id)

	if err != nil {
		err = helper.NewError("Parameter 'id' harus bernilai integer")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest))
		return
	}

	updateRequest := &models.DeliveryOrderUpdateByIDRequest{}

	err = ctx.ShouldBindJSON(updateRequest)
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
		ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
		return
	}

	deliveryOrder, errorLog := c.deliveryOrderUseCase.UpdateByID(intID, updateRequest, dbTransaction, ctx)
	if errorLog != nil {
		err = dbTransaction.Rollback()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
			return
		}

		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	err = dbTransaction.Commit()

	if err != nil {
		errorLog = helper.WriteLog(err, http.StatusInternalServerError, nil)
		ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, model.Response{Data: deliveryOrder, StatusCode: http.StatusOK})
	return
}

func (c *deliveryOrderController) UpdateDeliveryOrderDetailByID(ctx *gin.Context) {
	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	id := ctx.Param("id")
	intID, err := strconv.Atoi(id)

	if err != nil {
		err = helper.NewError("Parameter 'id' harus bernilai integer")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest))
		return
	}
	updateRequest := &models.DeliveryOrderDetailUpdateByIDRequest{}

	err = ctx.ShouldBindJSON(updateRequest)
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
		ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
		return
	}

	deliveryOrderDetail, errorLog := c.deliveryOrderUseCase.UpdateDODetailByID(intID, updateRequest, dbTransaction, ctx)
	if errorLog != nil {
		err = dbTransaction.Rollback()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
			return
		}

		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	err = dbTransaction.Commit()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, model.Response{Data: deliveryOrderDetail, StatusCode: http.StatusOK})
	return
}

func (c *deliveryOrderController) UpdateDeliveryOrderDetailByDeliveryOrderID(ctx *gin.Context) {
	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	id := ctx.Param("id")
	intID, err := strconv.Atoi(id)

	if err != nil {
		err = helper.NewError("Parameter 'id' harus bernilai integer")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest))
		return
	}

	updateRequest := []*models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest{}
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
		ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
		return
	}

	deliveryOrderDetail, errorLog := c.deliveryOrderUseCase.UpdateDoDetailByDeliveryOrderID(intID, updateRequest, dbTransaction, ctx)
	if errorLog != nil {
		err = dbTransaction.Rollback()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
			return
		}

		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	err = dbTransaction.Commit()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
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

	ctx.JSON(http.StatusOK, model.Response{Data: deliveryOrderResults, StatusCode: http.StatusOK})
	return
}

func (c *deliveryOrderController) Get(ctx *gin.Context) {
	deliveryOrderRequest, err := c.deliveryOrderValidator.GetDeliveryOrderValidator(ctx)
	if err != nil {
		return
	}

	deliveryOrders, errorLog := c.deliveryOrderUseCase.Get(deliveryOrderRequest)

	if errorLog.Err != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, model.Response{Data: deliveryOrders.DeliveryOrders, Total: deliveryOrders.Total, StatusCode: http.StatusOK})
	return
}

func (c *deliveryOrderController) GetDetails(ctx *gin.Context) {
	deliveryOrderDetailRequest, err := c.deliveryOrderValidator.GetDeliveryOrderDetailValidator(ctx)
	if err != nil {
		return
	}

	deliveryOrderDetails, errorLog := c.deliveryOrderUseCase.GetDetails(deliveryOrderDetailRequest)

	if errorLog.Err != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, model.Response{Data: deliveryOrderDetails.DeliveryOrderDetails, Total: deliveryOrderDetails.Total, StatusCode: http.StatusOK})
	return
}

func (c *deliveryOrderController) GetDetailsByDoId(ctx *gin.Context) {
	doIds := ctx.Param("id")
	doId, err := strconv.Atoi(doIds)

	if err != nil {
		err = helper.NewError("Parameter 'delivery order id' harus bernilai integer")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByErrorWithMessage(err, http.StatusBadRequest, nil))
		return
	}

	deliveryOrderRequest, err := c.deliveryOrderValidator.GetDeliveryOrderDetailByDoIDValidator(ctx)
	if err != nil {
		return
	}

	deliveryOrderRequest.ID = doId

	deliveryOrders, errorLog := c.deliveryOrderUseCase.GetDetailsByDoID(deliveryOrderRequest)

	if errorLog.Err != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, model.Response{Data: deliveryOrders, StatusCode: http.StatusOK})
	return
}

func (c *deliveryOrderController) GetDetailById(ctx *gin.Context) {
	doDetailIds := ctx.Param("do-detail-id")
	doDetailId, err := strconv.Atoi(doDetailIds)

	if err != nil {
		err = helper.NewError("Parameter 'delivery order detail id' harus bernilai integer")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByErrorWithMessage(err, http.StatusBadRequest, nil))
		return
	}

	doIds := ctx.Param("id")
	doId, err := strconv.Atoi(doIds)

	if err != nil {
		err = helper.NewError("Parameter 'delivery order id' harus bernilai integer")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByErrorWithMessage(err, http.StatusBadRequest, nil))
		return
	}

	deliveryOrders, errorLog := c.deliveryOrderUseCase.GetDetailByID(doDetailId, doId)

	if errorLog.Err != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, model.Response{Data: deliveryOrders, StatusCode: http.StatusOK})
	return
}

func (c *deliveryOrderController) GetBySalesmanID(ctx *gin.Context) {
	deliveryOrderReqeuest, err := c.deliveryOrderValidator.GetDeliveryOrderBySalesmanIDValidator(ctx)
	if err != nil {
		return
	}

	deliveryOrders, errorLog := c.deliveryOrderUseCase.GetBySalesmansID(deliveryOrderReqeuest)

	if errorLog.Err != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, model.Response{Data: deliveryOrders.DeliveryOrders, Total: deliveryOrders.Total, StatusCode: http.StatusOK})
	return
}

func (c *deliveryOrderController) GetByID(ctx *gin.Context) {
	var id int

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	ids := ctx.Param("id")
	id, err := strconv.Atoi(ids)

	if err != nil {
		err = helper.NewError("Parameter 'id' harus bernilai integer")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByError(err, http.StatusBadRequest))
		return
	}

	deliveryOrderRequest := &models.DeliveryOrderRequest{
		ID: id,
	}

	deliveryOrder, errorLog := c.deliveryOrderUseCase.GetByIDWithDetail(deliveryOrderRequest, ctx)

	if errorLog.Err != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, model.Response{Data: deliveryOrder, StatusCode: http.StatusOK})
	return
}

func (c *deliveryOrderController) GetDOJourneysByDoID(ctx *gin.Context) {
	doIds := ctx.Param("id")
	doId, err := strconv.Atoi(doIds)
	if err != nil {
		err = helper.NewError("Parameter 'delivery order id' harus bernilai integer")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByErrorWithMessage(err, http.StatusBadRequest, nil))
		return
	}
	deliveryOrderJourneys, errorLog := c.deliveryOrderUseCase.GetDOJourneysByDoID(doId, ctx)
	if errorLog != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}
	ctx.JSON(http.StatusOK, model.Response{Data: deliveryOrderJourneys.DeliveryOrderJourneys, Total: deliveryOrderJourneys.Total, StatusCode: http.StatusOK})
	return
}

func (c *deliveryOrderController) DeleteByID(ctx *gin.Context) {
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
		ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
		return
	}

	errorLog := c.deliveryOrderUseCase.DeleteByID(id, dbTransaction)

	if errorLog != nil {
		err = dbTransaction.Rollback()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
			return
		}

		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK})
}

func (c *deliveryOrderController) DeleteDetailByID(ctx *gin.Context) {
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
		ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
		return
	}

	errorLog := c.deliveryOrderUseCase.DeleteDetailByID(id, dbTransaction)

	if errorLog != nil {
		err = dbTransaction.Rollback()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
			return
		}

		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK})
}

func (c *deliveryOrderController) DeleteDetailByDoID(ctx *gin.Context) {
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
		ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
		return
	}

	errorLog := c.deliveryOrderUseCase.DeleteDetailByDoID(id, dbTransaction)

	if errorLog != nil {
		err = dbTransaction.Rollback()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError))
			return
		}

		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK})
}
