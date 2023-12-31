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
	ExportDetail(ctx *gin.Context)
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

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	ids := ctx.Param("so-id")
	id, err := strconv.Atoi(ids)

	if err != nil {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_GET, constants.ERROR_BAD_REQUEST_INT_SO_ID_PARAMS)},
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
			Message:       []string{helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_UPDATE, constants.ERROR_BAD_REQUEST_INT_SO_ID_PARAMS)},
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

	err = c.salesOrderValidator.UpdateSalesOrderByIdValidator(updateRequest, ctx)
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
			errorLog = helper.WriteLog(err, http.StatusInternalServerError, constants.ERROR_INTERNAL_SERVER_1)
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
		errorLog = helper.WriteLog(err, http.StatusInternalServerError, constants.ERROR_INTERNAL_SERVER_1)
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
	updateRequest := &models.SalesOrderDetailUpdateByIdRequest{}

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	soIds := ctx.Param("so-id")
	soId, err := strconv.Atoi(soIds)

	if err != nil {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       []string{helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_UPDATE, constants.ERROR_BAD_REQUEST_INT_SO_ID_PARAMS)},
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

	err = c.salesOrderValidator.UpdateSalesOrderDetailByIdValidator(updateRequest, ctx)
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

	salesOrderDetail, errorLog := c.salesOrderUseCase.UpdateSODetailById(soId, soDetailId, updateRequest, dbTransaction, ctx)

	if errorLog != nil {
		err = dbTransaction.Rollback()

		if err != nil {
			errorLog = helper.WriteLog(err, http.StatusInternalServerError, constants.ERROR_INTERNAL_SERVER_1)
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
		errorLog = helper.WriteLog(err, http.StatusInternalServerError, constants.ERROR_INTERNAL_SERVER_1)
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
			Message:       []string{helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_UPDATE, constants.ERROR_BAD_REQUEST_INT_SO_ID_PARAMS)},
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

	err = c.salesOrderValidator.UpdateSalesOrderDetailBySoIdValidator(updateRequest, ctx)
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
			errorLog = helper.WriteLog(err, http.StatusInternalServerError, constants.ERROR_INTERNAL_SERVER_1)
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
		errorLog = helper.WriteLog(err, http.StatusInternalServerError, constants.ERROR_INTERNAL_SERVER_1)
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
	var id int

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	sId := ctx.Param("so-id")
	id, err := c.salesOrderValidator.DeleteSalesOrderByIdValidator(sId, ctx)
	if err != nil {
		return
	}

	dbTransaction, err := c.db.BeginTx(ctx, nil)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError, ""))
		return
	}
	errorLog := c.salesOrderUseCase.DeleteById(id, dbTransaction)

	if errorLog != nil {
		err = dbTransaction.Rollback()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError, ""))
			return
		}

		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	err = dbTransaction.Commit()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError, constants.ERROR_INTERNAL_SERVER_1))
		return
	}

	ctx.JSON(http.StatusOK, baseModel.Response{StatusCode: http.StatusOK})
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
	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

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
	var id int

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	sId := ctx.Param("so-detail-id")
	id, err := c.salesOrderValidator.DeleteSalesOrderDetailByIdValidator(sId, ctx)
	if err != nil {
		return
	}

	dbTransaction, err := c.db.BeginTx(ctx, nil)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError, ""))
		return
	}
	errorLog := c.salesOrderUseCase.DeleteDetailById(id, dbTransaction)

	if errorLog != nil {
		err = dbTransaction.Rollback()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError, ""))
			return
		}

		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	err = dbTransaction.Commit()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError, constants.ERROR_INTERNAL_SERVER_1))
		return
	}

	ctx.JSON(http.StatusOK, baseModel.Response{StatusCode: http.StatusOK})
}

func (c *salesOrderController) DeleteDetailBySOID(ctx *gin.Context) {
	var id int

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	sId := ctx.Param("so-id")
	id, err := c.salesOrderValidator.DeleteSalesOrderDetailBySoIdValidator(sId, ctx)
	if err != nil {
		return
	}

	dbTransaction, err := c.db.BeginTx(ctx, nil)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError, ""))
		return
	}
	errorLog := c.salesOrderUseCase.DeleteDetailBySOId(id, dbTransaction)

	if errorLog != nil {
		err = dbTransaction.Rollback()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError, ""))
			return
		}

		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}

	err = dbTransaction.Commit()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.GenerateResultByError(err, http.StatusInternalServerError, constants.ERROR_INTERNAL_SERVER_1))
		return
	}

	ctx.JSON(http.StatusOK, baseModel.Response{StatusCode: http.StatusOK})

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

	fileName, errorLog := c.salesOrderUseCase.Export(salesOrderRequest, ctx)
	if errorLog != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}
	response := &models.SalesOrderExportResponse{
		StatusCode: http.StatusOK,
		UrlFile:    fmt.Sprintf("%s/%s.%s", constants.SALES_ORDER_EXPORT_PATH, fileName, salesOrderRequest.FileType),
	}
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *salesOrderController) ExportDetail(ctx *gin.Context) {
	salesOrderRequest, err := c.salesOrderValidator.ExportSalesOrderDetailValidator(ctx)
	if err != nil {
		return
	}

	fileName, errorLog := c.salesOrderUseCase.ExportDetail(salesOrderRequest, ctx)
	if errorLog != nil {
		ctx.JSON(errorLog.StatusCode, helper.GenerateResultByErrorLog(errorLog))
		return
	}
	response := &models.SalesOrderExportResponse{
		StatusCode: http.StatusOK,
		UrlFile:    fmt.Sprintf("%s/%s.%s", constants.SALES_ORDER_DETAIL_EXPORT_PATH, fileName, salesOrderRequest.FileType),
	}
	ctx.JSON(http.StatusOK, response)
	return
}
