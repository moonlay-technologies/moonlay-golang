package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/app/repositories"
	"strconv"

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
	OrderSourceValidation(ctx *gin.Context, orderSourceId int, soRefCode, actionName string) error
	UniqueValidation(ctx *gin.Context, value []*models.UniqueRequest) error
	MustActiveValidation(ctx *gin.Context, value []*models.MustActiveRequest) error
	MustActiveValidation422(ctx *gin.Context, value []*models.MustActiveRequest) error
	DateInputValidation(ctx *gin.Context, value []*models.DateInputRequest, actionName string) error
	MustEmptyValidation(ctx *gin.Context, value []*models.MustEmptyValidationRequest) error
	AgentIdValidation(ctx *gin.Context, agentId, userId int, actionName string) error
	StoreIdValidation(ctx *gin.Context, storeId, agentId int, actionName string) error
	SalesmanIdValidation(ctx *gin.Context, salesmanId, agentId int, actionName string) error
	BrandIdValidation(ctx *gin.Context, brandId []int, agentId int, actionName string) error
	UploadMandatoryValidation(request []*models.TemplateRequest) []string
	UploadIntTypeValidation(request []*models.TemplateRequest) (map[string]int, []string)
	UploadMustActiveValidation(request []*models.MustActiveRequest) []string
}

type requestValidationMiddleware struct {
	requestValidationRepository repositories.RequestValidationRepositoryInterface
	orderSourceRepository       repositories.OrderSourceRepositoryInterface
}

func InitRequestValidationMiddlewareInterface(requestValidationRepository repositories.RequestValidationRepositoryInterface, orderSourceRepository repositories.OrderSourceRepositoryInterface) RequestValidationMiddlewareInterface {
	return &requestValidationMiddleware{
		requestValidationRepository: requestValidationRepository,
		orderSourceRepository:       orderSourceRepository,
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

func (u *requestValidationMiddleware) OrderSourceValidation(ctx *gin.Context, orderSourceId int, soRefCode, actionName string) error {
	var result baseModel.Response
	messages := []string{}
	systemMessages := []string{}
	var error error

	getOrderSourceResultChan := make(chan *models.OrderSourceChan)
	go u.orderSourceRepository.GetByID(orderSourceId, false, ctx, getOrderSourceResultChan)
	getOrderSourceResult := <-getOrderSourceResultChan

	if getOrderSourceResult.Error != nil {

		message := helper.GenerateUnprocessableErrorMessage(actionName, "order_source_id tidak terdaftar!")
		messages = append(messages, message)
		systemMessages = []string{constants.ERROR_INVALID_PROCESS}

	} else if len(soRefCode) > 0 && getOrderSourceResult.OrderSource.SourceName != "manager" {

		message := helper.GenerateUnprocessableErrorMessage(actionName, fmt.Sprintf("order_source_id %d tidak dapat memasukkan so_ref_code", orderSourceId))
		messages = append(messages, message)
		systemMessages = []string{constants.ERROR_INVALID_PROCESS}

	}

	if len(messages) > 0 {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       messages,
			SystemMessage: systemMessages,
			StatusCode:    http.StatusUnprocessableEntity,
		})
		result.StatusCode = http.StatusUnprocessableEntity
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		error = fmt.Errorf(constants.ERROR_INVALID_PROCESS)
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
			StatusCode:    http.StatusConflict,
		})
		result.StatusCode = http.StatusConflict
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		error = fmt.Errorf("Duplicate value!")
	}

	return error
}

func (u *requestValidationMiddleware) MustActiveValidation422(ctx *gin.Context, value []*models.MustActiveRequest) error {
	return u.BaseMustActiveValidation(422, ctx, value)
}
func (u *requestValidationMiddleware) MustActiveValidation(ctx *gin.Context, value []*models.MustActiveRequest) error {
	return u.BaseMustActiveValidation(417, ctx, value)
}

func (u *requestValidationMiddleware) BaseMustActiveValidation(responseCode int, ctx *gin.Context, value []*models.MustActiveRequest) error {
	var result baseModel.Response
	messages := []string{}
	systemMessages := []string{}
	var error error
	mustActive := make(chan *models.MustActiveRequestChan)
	go u.requestValidationRepository.MustActiveValidation(value, mustActive)
	mustActiveResult := <-mustActive
	for k, v := range mustActiveResult.Total {
		if v < 1 {
			if value[k].CustomMessage != "" {
				messages = append(messages, value[k].CustomMessage)
				systemMessages = append(systemMessages, value[k].CustomMessage)
			} else {
				message := fmt.Sprintf("Data %s tidak ditemukan", value[k].ReqField)
				messages = append(messages, message)
				systemMessage := fmt.Sprintf("%s Not Found", value[k].ReqField)
				systemMessages = append(systemMessages, systemMessage)
			}
		}
	}

	if len(messages) > 0 {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       messages,
			SystemMessage: systemMessages,
			StatusCode:    responseCode,
		})

		result.StatusCode = responseCode
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
			systemMessage := constants.ERROR_INVALID_PROCESS
			systemMessages = append(systemMessages, systemMessage)
		}
	}

	if len(messages) > 0 {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       messages,
			SystemMessage: systemMessages,
			StatusCode:    http.StatusUnprocessableEntity,
		})
		result.StatusCode = http.StatusUnprocessableEntity
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		error = fmt.Errorf(constants.ERROR_INVALID_PROCESS)
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
			StatusCode:    http.StatusUnprocessableEntity,
		})
		result.StatusCode = http.StatusUnprocessableEntity
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		error = fmt.Errorf(constants.ERROR_INVALID_PROCESS)
	}

	return error
}

func (u *requestValidationMiddleware) AgentIdValidation(ctx *gin.Context, agentId, userId int, actionName string) error {
	var result baseModel.Response
	messages := []string{}
	systemMessages := []string{}
	var error error

	agentIdValidationResultChan := make(chan *models.RequestIdValidationChan)
	go u.requestValidationRepository.AgentIdValidation(agentId, userId, agentIdValidationResultChan)
	agentIdValidationResult := <-agentIdValidationResultChan

	if agentIdValidationResult.Total < 1 {

		message := helper.GenerateUnprocessableErrorMessage(actionName, fmt.Sprintf("agent_id %d tidak terdaftar untuk user_id %d", agentId, userId))
		messages = append(messages, message)
		systemMessages = []string{constants.ERROR_INVALID_PROCESS}

	}

	if len(messages) > 0 {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       messages,
			SystemMessage: systemMessages,
			StatusCode:    http.StatusUnprocessableEntity,
		})
		result.StatusCode = http.StatusUnprocessableEntity
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		error = fmt.Errorf(constants.ERROR_INVALID_PROCESS)
	}

	return error
}

func (u *requestValidationMiddleware) StoreIdValidation(ctx *gin.Context, storeId, agentId int, actionName string) error {
	var result baseModel.Response
	messages := []string{}
	systemMessages := []string{}
	var error error

	storeIdValidationResultChan := make(chan *models.RequestIdValidationChan)
	go u.requestValidationRepository.StoreIdValidation(storeId, agentId, storeIdValidationResultChan)
	storeIdValidationResult := <-storeIdValidationResultChan

	if storeIdValidationResult.Total < 1 {

		message := helper.GenerateUnprocessableErrorMessage(actionName, fmt.Sprintf("store_id %d tidak terdaftar untuk agent_id %d", storeId, agentId))
		messages = append(messages, message)
		systemMessages = []string{constants.ERROR_INVALID_PROCESS}

	}

	if len(messages) > 0 {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       messages,
			SystemMessage: systemMessages,
			StatusCode:    http.StatusUnprocessableEntity,
		})
		result.StatusCode = http.StatusUnprocessableEntity
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		error = fmt.Errorf(constants.ERROR_INVALID_PROCESS)
	}

	return error
}

func (u *requestValidationMiddleware) SalesmanIdValidation(ctx *gin.Context, salesmanId, agentId int, actionName string) error {
	var result baseModel.Response
	messages := []string{}
	systemMessages := []string{}
	var error error

	salesmanIdValidationResultChan := make(chan *models.RequestIdValidationChan)
	go u.requestValidationRepository.SalesmanIdValidation(salesmanId, agentId, salesmanIdValidationResultChan)
	salesmanIdValidationResult := <-salesmanIdValidationResultChan

	if salesmanIdValidationResult.Total < 1 {

		message := helper.GenerateUnprocessableErrorMessage(actionName, fmt.Sprintf("salesman_id %d tidak terdaftar untuk agent_id %d", salesmanId, agentId))
		messages = append(messages, message)
		systemMessages = []string{constants.ERROR_INVALID_PROCESS}

	}

	if len(messages) > 0 {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       messages,
			SystemMessage: systemMessages,
			StatusCode:    http.StatusUnprocessableEntity,
		})
		result.StatusCode = http.StatusUnprocessableEntity
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		error = fmt.Errorf(constants.ERROR_INVALID_PROCESS)
	}

	return error
}

func (u *requestValidationMiddleware) BrandIdValidation(ctx *gin.Context, brandId []int, agentId int, actionName string) error {
	var result baseModel.Response
	messages := []string{}
	systemMessages := []string{}
	var error error

	for _, v := range brandId {

		brandIdValidationResultChan := make(chan *models.RequestIdValidationChan)
		go u.requestValidationRepository.BrandIdValidation(v, agentId, brandIdValidationResultChan)
		brandIdValidationResult := <-brandIdValidationResultChan

		if brandIdValidationResult.Total < 1 {

			message := helper.GenerateUnprocessableErrorMessage(actionName, fmt.Sprintf("brand_id %d tidak terdaftar untuk agent_id %d", v, agentId))
			messages = append(messages, message)
			systemMessages = []string{constants.ERROR_INVALID_PROCESS}

		}
	}

	if len(messages) > 0 {
		errorLog := helper.NewWriteLog(baseModel.ErrorLog{
			Message:       messages,
			SystemMessage: systemMessages,
			StatusCode:    http.StatusUnprocessableEntity,
		})
		result.StatusCode = http.StatusUnprocessableEntity
		result.Error = errorLog
		ctx.JSON(result.StatusCode, result)
		error = fmt.Errorf(constants.ERROR_INVALID_PROCESS)
	}

	return error
}

func (u *requestValidationMiddleware) UploadMandatoryValidation(request []*models.TemplateRequest) []string {
	errors := []string{}

	for _, value := range request {
		if len(value.Value) < 1 {
			error := fmt.Sprintf("Data %s tidak boleh kosong", value.Field)
			errors = append(errors, error)
		}
	}

	return errors
}

func (u *requestValidationMiddleware) UploadIntTypeValidation(request []*models.TemplateRequest) (map[string]int, []string) {
	result := map[string]int{}
	errors := []string{}

	for _, v := range request {
		parseInt, error := strconv.Atoi(v.Value)

		if error != nil {
			error := fmt.Sprintf("Data %s harus bertipe data integer", v.Value)
			errors = append(errors, error)
		} else {
			result[v.Field] = parseInt
		}
	}

	return result, errors
}

func (u *requestValidationMiddleware) UploadMustActiveValidation(request []*models.MustActiveRequest) []string {

	errors := []string{}

	mustActive := make(chan *models.MustActiveRequestChan)
	go u.requestValidationRepository.MustActiveValidation(request, mustActive)
	mustActiveResult := <-mustActive

	for k, v := range mustActiveResult.Total {
		if v < 1 {
			var error string
			if request[k].CustomMessage != "" {
				error = request[k].CustomMessage
			} else {
				error = fmt.Sprintf("Kode = %s sudah Tidak Aktif. Silahkan gunakan Kode yang lain.", request[k].Id)
			}

			errors = append(errors, error)
		}
	}

	return errors
}
