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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RequestValidationMiddlewareInterface interface {
	DataTypeValidation(ctx *gin.Context, err error, unmarshalTypeError *json.UnmarshalTypeError)
	MandatoryValidation(ctx *gin.Context, err error)
	OrderSourceValidation(ctx *gin.Context, orderSourceId int, actionName string) error
	UniqueValidation(ctx *gin.Context, value []*models.UniqueRequest) error
	MustActiveValidation(ctx *gin.Context, value []*models.MustActiveRequest) error
	DateInputValidation(ctx *gin.Context, value []*models.DateInputRequest, actionName string) error
	MustEmptyValidation(ctx *gin.Context, value []*models.MustEmptyValidationRequest) error
}

type requestValidationMiddleware struct {
	requestValidationRepository repositories.RequestValidationRepositoryInterface
	orderSourceRepository       repositories.OrderSourceRepositoryInterface
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

func (u *requestValidationMiddleware) OrderSourceValidation(ctx *gin.Context, orderSourceId int, actionName string) error {
	var result baseModel.Response
	messages := []string{}
	systemMessages := []string{}
	var error error

	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(orderSourceId, false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getOrderSourceResult.Error != nil {

		message := helper.GenerateUnprocessableErrorMessage(actionName, "order_sourde_id tidak terdaftar!")
		messages = append(messages, message)
		systemMessages = []string{"Invalid Process"}

	} else {

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
		error = fmt.Errorf("Invalid Process")
	}

	return error
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

func (u *requestValidationMiddleware) DateInputValidation(ctx *gin.Context, value []*models.DateInputRequest, actionName string) error {
	var result baseModel.Response
	messages := []string{}
	systemMessages := []string{}
	var error error

	for _, v := range value {
		_, err := time.Parse("2006-01-02", v.Value)
		if err != nil {
			message := helper.GenerateUnprocessableErrorMessage(actionName, fmt.Sprintf("field %s harus memiliki format yyyy-mm-dd", v.Field))
			messages = append(messages, message)
			systemMessage := "Invalid Process"
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
		error = fmt.Errorf("Invalid Process")
	}

	return error
}

func (u *requestValidationMiddleware) MustEmptyValidation(ctx *gin.Context, value []*models.MustEmptyValidationRequest) error {
	var result baseModel.Response
	messages := []string{}
	systemMessages := []string{}
	var error error

	for _, v := range value {
		mustContains := make(chan *models.MustEmptyValidationRequestChan)
		go u.requestValidationRepository.MustEmptyValidation(v, mustContains)
		mustContainsResult := <-mustContains
		if !mustContainsResult.Result {
			message := mustContainsResult.Message
			messages = append(messages, message)
			systemMessage := message
			systemMessages = append(systemMessages, systemMessage)
		}
	}

	if len(messages) > 0 {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       messages,
			SystemMessage: systemMessages,
			StatusCode:    http.StatusBadRequest,
		})
		result.StatusCode = http.StatusNotAcceptable
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		error = fmt.Errorf("Inactive value!")
	}

	return error
}
