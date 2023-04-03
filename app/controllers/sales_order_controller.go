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
	"time"

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
	GetSyncToKafkaHistories(ctx *gin.Context)
	GetSOJourneys(ctx *gin.Context)
	GetSOJourneyBySoId(ctx *gin.Context)
	GetSOUploadHistories(ctx *gin.Context)
	GetSoUploadHistoriesById(ctx *gin.Context)
	GetSoUploadErrorLogByReqId(ctx *gin.Context)
	GetSoUploadErrorLogBySoUploadHistoryId(ctx *gin.Context)
	DeleteByID(ctx *gin.Context)
	DeleteDetailByID(ctx *gin.Context)
	DeleteDetailBySOID(ctx *gin.Context)
	RetrySyncToKafka(ctx *gin.Context)
	Export(ctx *gin.Context)
}

type salesOrderController struct {
	cartUseCase                 usecases.CartUseCaseInterface
	salesOrderUseCase           usecases.SalesOrderUseCaseInterface
	salesOrderValidator         usecases.SalesOrderValidatorInterface
	requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface
	db                          dbresolver.DB
	ctx                         context.Context
}

func InitSalesOrderController(cartUseCase usecases.CartUseCaseInterface, salesOrderUseCase usecases.SalesOrderUseCaseInterface, salesOrderValidator usecases.SalesOrderValidatorInterface, requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, db dbresolver.DB, ctx context.Context) SalesOrderControllerInterface {
	return &salesOrderController{
		cartUseCase:                 cartUseCase,
		salesOrderUseCase:           salesOrderUseCase,
		salesOrderValidator:         salesOrderValidator,
		requestValidationMiddleware: requestValidationMiddleware,
		db:                          db,
		ctx:                         ctx,
	}
}

func (c *salesOrderController) Get(ctx *gin.Context) {
	salesOrderRequest, err := c.salesOrderValidator.GetSalesOrderValidator(ctx)
	if err != nil {
		return
	}

	salesOrders, errorLog := c.salesOrderUseCase.Get(salesOrderRequest)

	if errorLog.Err != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, baseModel.Response{Data: salesOrders.SalesOrders, Total: salesOrders.Total, StatusCode: http.StatusOK})
}

func (c *salesOrderController) GetByID(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	var id int

	// ctx.Set("full_path", ctx.FullPath())
	// ctx.Set("method", ctx.Request.Method)

	ids := ctx.Param("so-id")
	id, err := strconv.Atoi(ids)

	if err != nil {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_GET, "sales order id harus bernilai integer")},
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		result.StatusCode = http.StatusUnprocessableEntity
		result.Error = errorLog
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

	err = c.salesOrderValidator.CreateSalesOrderValidator(insertRequest, ctx)
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
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_UPDATE, "sales order id harus bernilai integer")},
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		result.StatusCode = http.StatusUnprocessableEntity
		result.Error = errorLog
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

}

func (c *salesOrderController) UpdateSODetailByID(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	updateRequest := &models.UpdateSalesOrderDetailByIdRequest{}

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	soIds := ctx.Param("so-id")
	soId, err := strconv.Atoi(soIds)

	if err != nil {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_UPDATE, "sales order id harus bernilai integer")},
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		result.StatusCode = http.StatusUnprocessableEntity
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	soDetailIds := ctx.Param("so-detail-id")
	soDetailId, err := strconv.Atoi(soDetailIds)

	if err != nil {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_UPDATE, "sales order detail id harus bernilai integer")},
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		result.StatusCode = http.StatusUnprocessableEntity
		result.Error = errorLog
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

	dbTransaction, err := c.db.BeginTx(ctx, nil)

	if err != nil {
		errorLog := helper.WriteLog(err, http.StatusInternalServerError, nil)
		resultErrorLog = errorLog
		result.StatusCode = http.StatusInternalServerError
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	salesOrderDetail, errorLog := c.salesOrderUseCase.UpdateSODetailById(soId, soDetailId, updateRequest, dbTransaction, ctx)

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

}

func (c *salesOrderController) UpdateSODetailBySOID(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog
	var soId int
	updateRequest := &models.SalesOrderUpdateRequest{}

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	soIds := ctx.Param("so-id")
	soId, err := strconv.Atoi(soIds)

	if err != nil {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_UPDATE, "sales order id harus bernilai integer")},
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		result.StatusCode = http.StatusUnprocessableEntity
		result.Error = errorLog
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

}

func (c *salesOrderController) GetDetails(ctx *gin.Context) {
	salesOrderDetailRequest, err := c.salesOrderValidator.GetSalesOrderDetailValidator(ctx)
	if err != nil {
		return
	}

	salesOrders, errorLog := c.salesOrderUseCase.GetDetails(salesOrderDetailRequest)

	if errorLog.Err != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, baseModel.Response{Data: salesOrders.SalesOrderDetails, Total: salesOrders.Total, StatusCode: http.StatusOK})

}

func (c *salesOrderController) DeleteByID(ctx *gin.Context) {
	var result baseModel.Response
	var id int

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	sId := ctx.Param("so-id")
	id, err := strconv.Atoi(sId)

	if err != nil {
		err = helper.NewError(constants.ERROR_BAD_REQUEST_INT_ID_PARAMS)
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
			Table:           "sales_orders s JOIN order_statuses o ON s.order_status_id = o.id",
			SelectedCollumn: "o.name",
			Clause:          fmt.Sprintf("o.id = %d AND s.order_status_id NOT IN (5,6,9,10)", id),
			MessageFormat:   "Status Sales Order <result>",
		},
		{
			Table:           "delivery_orders d JOIN sales_orders s ON d.sales_order_id = s.id",
			SelectedCollumn: "d.id",
			Clause:          fmt.Sprintf("s.id = %d AND d.deleted_at IS NULL", id),
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

}

func (c *salesOrderController) GetDetailsBySoId(ctx *gin.Context) {
	soIds := ctx.Param("so-id")
	soId, err := strconv.Atoi(soIds)

	if err != nil {
		err = helper.NewError("Parameter 'so-id' harus bernilai integer")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByErrorWithMessage(err, http.StatusBadRequest, nil))
		return
	}

	salesOrderDetailRequest, err := c.salesOrderValidator.GetSalesOrderDetailValidator(ctx)
	if err != nil {
		return
	}

	salesOrderDetailRequest.SoID = soId

	salesOrders, errorLog := c.salesOrderUseCase.GetDetails(salesOrderDetailRequest)

	if errorLog.Err != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, baseModel.Response{Data: salesOrders.SalesOrderDetails, Total: salesOrders.Total, StatusCode: http.StatusOK})

}

func (c *salesOrderController) GetDetailsById(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog

	soDetailIds := ctx.Param("so-detail-id")
	soDetailId, err := strconv.Atoi(soDetailIds)

	if err != nil {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_UPDATE, "sales order detail id harus bernilai integer")},
			SystemMessage: []string{constants.ERROR_INVALID_PROCESS},
			StatusCode:    http.StatusUnprocessableEntity,
		})
		result.StatusCode = http.StatusUnprocessableEntity
		result.Error = errorLog
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

}

func (c *salesOrderController) GetSyncToKafkaHistories(ctx *gin.Context) {

	salesOrderEventLogRequest, err := c.salesOrderValidator.GetSalesOrderSyncToKafkaHistoriesValidator(ctx)
	if err != nil {
		return
	}

	salesOrders, errorLog := c.salesOrderUseCase.GetSyncToKafkaHistories(salesOrderEventLogRequest, ctx)

	if errorLog != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, baseModel.Response{Data: salesOrders, StatusCode: http.StatusOK})

}

func (c *salesOrderController) GetSOJourneys(ctx *gin.Context) {
	salesOrderJourneysRequest, err := c.salesOrderValidator.GetSalesOrderJourneysValidator(ctx)
	if err != nil {
		return
	}

	salesOrderJourneys, errorLog := c.salesOrderUseCase.GetSOJourneys(salesOrderJourneysRequest, ctx)

	if errorLog != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, baseModel.Response{Data: salesOrderJourneys.SalesOrderJourneys, Total: salesOrderJourneys.Total, StatusCode: http.StatusOK})

}

func (c *salesOrderController) GetSOJourneyBySoId(ctx *gin.Context) {
	soIds := ctx.Param("so-id")
	soId, err := strconv.Atoi(soIds)

	if err != nil {
		err = helper.NewError("Parameter 'sales order id' harus bernilai integer")
		ctx.JSON(http.StatusBadRequest, helper.GenerateResultByErrorWithMessage(err, http.StatusBadRequest, nil))
		return
	}

	salesOrderJourney, errorLog := c.salesOrderUseCase.GetSOJourneyBySOId(soId, ctx)

	if errorLog != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, baseModel.Response{Data: salesOrderJourney.SalesOrderJourneys, Total: salesOrderJourney.Total, StatusCode: http.StatusOK})

}

func (c *salesOrderController) GetSOUploadHistories(ctx *gin.Context) {
	soUploadHistoriesRequest, err := c.salesOrderValidator.GetSOUploadHistoriesValidator(ctx)
	if err != nil {
		return
	}

	soUploadHistories, errorLog := c.salesOrderUseCase.GetSOUploadHistories(soUploadHistoriesRequest, ctx)

	if errorLog != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, baseModel.Response{Data: soUploadHistories.SoUploadHistories, Total: soUploadHistories.Total, StatusCode: http.StatusOK})
}

func (c *salesOrderController) GetSoUploadHistoriesById(ctx *gin.Context) {
	soUploadHistoriesId := ctx.Param("id")

	soUploadHistories, errorLog := c.salesOrderUseCase.GetSOUploadHistoriesByid(soUploadHistoriesId, ctx)

	if errorLog != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, baseModel.Response{Data: soUploadHistories, StatusCode: http.StatusOK})

}

func (c *salesOrderController) GetSoUploadErrorLogByReqId(ctx *gin.Context) {
	soUploadRequestId := ctx.Param("id")

	request := &models.GetSoUploadErrorLogsRequest{
		RequestID: soUploadRequestId,
	}

	soUploadErrorLogs, errorLog := c.salesOrderUseCase.GetSOUploadErrorLogsByReqId(request, ctx)

	if errorLog != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, baseModel.Response{Data: soUploadErrorLogs.SoUploadErrosLogs, Total: soUploadErrorLogs.Total, StatusCode: http.StatusOK})

}

func (c *salesOrderController) GetSoUploadErrorLogBySoUploadHistoryId(ctx *gin.Context) {
	soUploadHistoryId := ctx.Param("id")

	request := &models.GetSoUploadErrorLogsRequest{
		SoUploadHistoryID: soUploadHistoryId,
	}

	soUploadErrorLogs, errorLog := c.salesOrderUseCase.GetSOUploadErrorLogsBySoUploadHistoryId(request, ctx)

	if errorLog != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	ctx.JSON(http.StatusOK, baseModel.Response{Data: soUploadErrorLogs.SoUploadErrosLogs, Total: soUploadErrorLogs.Total, StatusCode: http.StatusOK})

}

func (c *salesOrderController) DeleteDetailByID(ctx *gin.Context) {
	var result baseModel.Response
	var id int

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	sId := ctx.Param("so-detail-id")
	id, err := strconv.Atoi(sId)

	if err != nil {
		err = helper.NewError(constants.ERROR_BAD_REQUEST_INT_ID_PARAMS)
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, result.StatusCode, nil)
		ctx.JSON(result.StatusCode, result)
		return
	}
	mustActiveField := []*models.MustActiveRequest{
		{
			Table:    "sales_order_details",
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
			Table:           "sales_order_details s JOIN order_statuses o ON s.order_status_id = o.id",
			SelectedCollumn: "o.name",
			Clause:          fmt.Sprintf("s.id = %d AND s.order_status_id NOT IN (16)", id),
			MessageFormat:   "Hanya status cancelled pada Sales Order Detail yang dapat di delete",
		},
		{
			Table:           "sales_orders s JOIN sales_order_details sd ON s.id = sd.sales_order_id JOIN delivery_orders d ON d.sales_order_id = s.id",
			SelectedCollumn: "d.id",
			Clause:          fmt.Sprintf("sd.id = %d AND d.deleted_at IS NULL", id),
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
	errorLog := c.salesOrderUseCase.DeleteDetailById(id, dbTransaction)

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

}

func (c *salesOrderController) DeleteDetailBySOID(ctx *gin.Context) {
	var result baseModel.Response
	var id int

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	sId := ctx.Param("so-id")
	id, err := strconv.Atoi(sId)

	if err != nil {
		err = helper.NewError(constants.ERROR_BAD_REQUEST_INT_ID_PARAMS)
		result.StatusCode = http.StatusBadRequest
		result.Error = helper.WriteLog(err, result.StatusCode, nil)
		ctx.JSON(result.StatusCode, result)
		return
	}
	mustActiveField := []*models.MustActiveRequest{
		{
			Table:    "sales_order_details",
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
	errorLog := c.salesOrderUseCase.DeleteDetailBySOId(id, dbTransaction)

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

}

func (c *salesOrderController) RetrySyncToKafka(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	logId := ctx.Param("log-id")

	retryToKafka, errorLog := c.salesOrderUseCase.RetrySyncToKafka(logId)

	if errorLog != nil {
		resultErrorLog = errorLog
		result.StatusCode = resultErrorLog.StatusCode
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = retryToKafka
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
}
func (c *salesOrderController) Export(ctx *gin.Context) {
	salesOrderRequest, err := c.salesOrderValidator.ExportSalesOrderValidator(ctx)
	if err != nil {
		return
	}

	fileDate := time.Now().Format("2_January_2006")
	fmt.Println(fileDate)

	fileName, errorLog := c.salesOrderUseCase.Export(salesOrderRequest, ctx)
	fmt.Println(fileName)
	if errorLog != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}
	ctx.JSON(http.StatusOK, fmt.Sprintf("%s_%s.%s", constants.SALES_ORDER_EXPORT_PATH, fileDate, salesOrderRequest.FileType)) // Makesure dor response pattern
	// ctx.JSON(http.StatusOK, fmt.Sprintf("%s/%s.%s", constants.SALES_ORDER_EXPORT_PATH, fileName, deliveryOrderRequest.FileType))
	return
}
