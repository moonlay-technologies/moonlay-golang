package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"order-service/app/middlewares"
	"order-service/app/models"
	"order-service/app/usecases"
	baseModel "order-service/global/utils/model"

	"github.com/gin-gonic/gin"
)

type UploadControllerInterface interface {
	UploadSOSJ(ctx *gin.Context)
	UploadDO(ctx *gin.Context)
	UploadSO(ctx *gin.Context)
	RetryUploadSO(ctx *gin.Context)
	RetryUploadSOSJ(ctx *gin.Context)
	GetSoUploadErrorLogsByReqId(ctx *gin.Context)
}

type uploadController struct {
	uploadUseCase               usecases.UploadUseCaseInterface
	requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface
	ctx                         context.Context
}

func InitUploadController(uploadUseCase usecases.UploadUseCaseInterface, requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, ctx context.Context) UploadControllerInterface {
	return &uploadController{
		uploadUseCase:               uploadUseCase,
		requestValidationMiddleware: requestValidationMiddleware,
		ctx:                         ctx,
	}
}

func (c *uploadController) UploadSOSJ(ctx *gin.Context) {

	var result baseModel.Response
	uploadRequest := &models.UploadSOSJRequest{}

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	err := ctx.BindJSON(uploadRequest)

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

	errorLog := c.uploadUseCase.UploadSOSJ(uploadRequest, ctx)

	if errorLog != nil {
		result.StatusCode = errorLog.StatusCode
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = map[string]string{
		"request_id": ctx.Value("RequestId").(string),
	}

	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return

}

func (c *uploadController) UploadDO(ctx *gin.Context) {

	c.uploadUseCase.UploadDO(ctx)

}

func (c *uploadController) UploadSO(ctx *gin.Context) {

	var result baseModel.Response
	uploadRequest := &models.UploadSORequest{}

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	err := ctx.BindJSON(uploadRequest)

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

	errorLog := c.uploadUseCase.UploadSO(uploadRequest, ctx)

	if errorLog != nil {
		result.StatusCode = errorLog.StatusCode
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = map[string]string{
		"request_id": ctx.Value("RequestId").(string),
	}

	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return

}

func (c *uploadController) RetryUploadSO(ctx *gin.Context) {

	var result baseModel.Response

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	id := ctx.Param("so-upload-history-id")

	errorLog := c.uploadUseCase.RetryUploadSO(id, ctx)

	if errorLog != nil {
		result.StatusCode = errorLog.StatusCode
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = map[string]string{
		"so_upload_history_id": id,
		"message":              "upload on progress",
	}

	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return

}

func (c *uploadController) RetryUploadSOSJ(ctx *gin.Context) {

	var result baseModel.Response

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	id := ctx.Param("sosj-upload-history-id")

	errorLog := c.uploadUseCase.RetryUploadSOSJ(id, ctx)

	if errorLog != nil {
		result.StatusCode = errorLog.StatusCode
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = map[string]string{
		"sosj_upload_history_id": id,
		"message":                "upload on progress",
	}

	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return

}

func (c *uploadController) GetSoUploadErrorLogsByReqId(ctx *gin.Context) {
	var result baseModel.Response
	var resultErrorLog *baseModel.ErrorLog

	ctx.Set("full_path", ctx.FullPath())
	ctx.Set("method", ctx.Request.Method)

	requestId := ctx.Param("request-id")

	request := &models.GetSosjUploadErrorLogsRequest{
		RequestID: requestId,
	}

	sosjUploadErrorLogs, errorLog := c.uploadUseCase.GetSosjUploadErrorLogs(request, ctx)

	if errorLog != nil {
		resultErrorLog = errorLog
		result.StatusCode = resultErrorLog.StatusCode
		result.Error = resultErrorLog
		ctx.JSON(result.StatusCode, result)
		return
	}

	result.Data = sosjUploadErrorLogs.SosjUploadErrorLogs
	result.StatusCode = http.StatusOK
	ctx.JSON(http.StatusOK, result)
	return

}
