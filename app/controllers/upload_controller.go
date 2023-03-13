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

	c.uploadUseCase.UploadSOSJ(ctx)

}

func (c *uploadController) UploadDO(ctx *gin.Context) {
	var result baseModel.Response
	uploadRequest := &models.UploadDORequest{}

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

	errorLog := c.uploadUseCase.UploadDO(uploadRequest, ctx)

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
