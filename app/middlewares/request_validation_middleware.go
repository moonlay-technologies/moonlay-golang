package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/app/repositories"
	"order-service/global/utils/helper"
	baseModel "order-service/global/utils/model"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RequestValidationMiddlewareInterface interface {
	DataTypeValidation(ctx *gin.Context, err error, unmarshalTypeError *json.UnmarshalTypeError)
	MandatoryValidation(ctx *gin.Context, err error)
	UniqueValidation(ctx *gin.Context, value []*models.UniqueRequest) error
	MustActiveValidation(ctx *gin.Context, value []*models.MustActiveRequest) error
	InvalidValidation(ctx *gin.Context, err error, unmarshalTypeError *json.InvalidUnmarshalError)
}

type requestValidationMiddleware struct {
	requestValidationRepository repositories.RequestValidationRepositoryInterface
}

func InitRequestValidationMiddlewareInterface(requestValidationRepository repositories.RequestValidationRepositoryInterface) RequestValidationMiddlewareInterface {
	return &requestValidationMiddleware{
		requestValidationRepository: requestValidationRepository,
	}
}

func (u *requestValidationMiddleware) DataTypeValidation(ctx *gin.Context, err error, unmarshalTypeError *json.UnmarshalTypeError) {
	var result baseModel.Response
	messages := []string{}

	message := fmt.Sprintf("Data %s harus bertipe data %s", unmarshalTypeError.Field, unmarshalTypeError.Type)
	messages = append(messages, message)

	errorLog := helper.NewWriteLog(baseModel.ErrorLog{
		Message:       messages,
		SystemMessage: []string{err.Error()},
		StatusCode:    http.StatusBadRequest,
	})
	result.StatusCode = http.StatusBadRequest
	result.Error = errorLog
	ctx.JSON(result.StatusCode, result)
	return
}

func (u *requestValidationMiddleware) MandatoryValidation(ctx *gin.Context, err error) {
	var result baseModel.Response
	messages := []string{}

	for _, value := range err.(validator.ValidationErrors) {
		message := fmt.Sprintf("Data %s tidak boleh kosong", value.Field())
		messages = append(messages, message)
	}
	errorLog := helper.NewWriteLog(baseModel.ErrorLog{
		Message:       messages,
		SystemMessage: strings.Split(err.Error(), "\n"),
		StatusCode:    http.StatusBadRequest,
	})
	result.StatusCode = http.StatusBadRequest
	result.Error = errorLog
	ctx.JSON(result.StatusCode, result)
	return
}

func (u *requestValidationMiddleware) UniqueValidation(ctx *gin.Context, value []*models.UniqueRequest) error {
	var result baseModel.Response
	messages := []string{}
	systemMessages := []string{}
	var error error

	for _, v := range value {
		checkUnique := make(chan *models.UniqueRequestChan)
		go u.requestValidationRepository.UniqueValidation(v, checkUnique)
		checkUniqueResult := <-checkUnique

		if checkUniqueResult.Total > 0 {
			message := fmt.Sprintf("Data %s duplikat", v.Field)
			messages = append(messages, message)
			systemMessage := fmt.Sprintf("%s Duplicate id for %s", v.Field, v.Field)
			systemMessages = append(systemMessages, systemMessage)
		}
	}

	if len(messages) > 0 {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       messages,
			SystemMessage: systemMessages,
			StatusCode:    http.StatusBadRequest,
		})
		result.StatusCode = http.StatusConflict
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		error = fmt.Errorf("Duplicate value!")
	}

	return error
}

func (u *requestValidationMiddleware) MustActiveValidation(ctx *gin.Context, value []*models.MustActiveRequest) error {
	var result baseModel.Response
	messages := []string{}
	systemMessages := []string{}
	var error error

	for _, v := range value {
		mustActive := make(chan *models.MustActiveRequestChan)
		go u.requestValidationRepository.MustActiveValidation(v, mustActive)
		mustActiveResult := <-mustActive

		if mustActiveResult.Total < 1 {
			message := fmt.Sprintf("Data %s tidak ditemukan", v.ReqField)
			messages = append(messages, message)
			systemMessage := fmt.Sprintf("%s Not Found", v.ReqField)
			systemMessages = append(systemMessages, systemMessage)
		}
	}

	if len(messages) > 0 {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       messages,
			SystemMessage: systemMessages,
			StatusCode:    http.StatusBadRequest,
		})
		result.StatusCode = http.StatusExpectationFailed
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		error = fmt.Errorf("Inactive value!")
	}

	return error
}

func (u *requestValidationMiddleware) InvalidValidation(ctx *gin.Context, err error, unmarshalTypeError *json.InvalidUnmarshalError) {
	var result baseModel.Response
	messages := []string{}

	message := fmt.Sprintf("Error: %s", unmarshalTypeError.Error())
	messages = append(messages, message)

	errorLog := helper.NewWriteLog(baseModel.ErrorLog{
		Message:       messages,
		SystemMessage: []string{err.Error()},
		StatusCode:    http.StatusBadRequest,
	})
	result.StatusCode = http.StatusBadRequest
	result.Error = errorLog
	ctx.JSON(result.StatusCode, result)
	return
}
