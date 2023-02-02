package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"poc-order-service/global/utils/helper"
	baseModel "poc-order-service/global/utils/model"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func DataTypeValidation(ctx *gin.Context, err error, unmarshalTypeError *json.UnmarshalTypeError) {
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
}

func MandatoryValidation(ctx *gin.Context, err error) {
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
}
