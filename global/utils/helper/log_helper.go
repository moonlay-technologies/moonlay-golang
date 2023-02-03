package helper

import (
	"fmt"
	"net/http"
	"order-service/global/utils/model"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// var DefaultStatusText = map[int]string{
// 	http.StatusInternalServerError: "Something went wrong, please try again later",
// 	http.StatusNotFound:            "data not found",
// }

var DefaultStatusText = map[int]string{
	http.StatusInternalServerError: "Terjadi Kesalahan, Silahkan Coba lagi Nanti",
	http.StatusNotFound:            "Data tidak Ditemukan",
}

func WriteLog(err error, errorCode int, message interface{}) *model.ErrorLog {
	if pc, file, line, ok := runtime.Caller(1); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		var output *model.ErrorLog
		output.StatusCode = errorCode
		output.Err = err

		if errorCode == 422 {
			output.Message = "Field Validation"
			output.Fields = message
		} else {
			output.SystemMessage = err.Error()
			if message == nil {
				output.Message = DefaultStatusText[errorCode]
				if message == "" {
					output.Message = http.StatusText(errorCode)
				}
			} else {
				output.Message = message
			}
			if errorCode == http.StatusInternalServerError {
				output.Line = fmt.Sprintf("%d", line)
				output.Filename = fmt.Sprintf("%s", file)
				output.Function = fmt.Sprintf("%s", funcName)
			}
		}

		log := map[string]interface{}{}
		_ = DecodeMapType(output, &log)
		logrus.WithFields(log).Error(message)
		return output
	}

	return nil
}

func NewWriteLog(errorLog model.ErrorLog) *model.ErrorLog {
	if pc, file, line, ok := runtime.Caller(1); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		var output *model.ErrorLog

		if errorLog.StatusCode == 500 {
			output = &model.ErrorLog{
				Line:          fmt.Sprintf("%d", line),
				Filename:      fmt.Sprintf("%s", file),
				Function:      fmt.Sprintf("%s", funcName),
				Message:       errorLog.Message,
				SystemMessage: errorLog.Err.Error(),
				Err:           errorLog.Err,
				StatusCode:    errorLog.StatusCode,
			}
		} else if errorLog.StatusCode == 401 || errorLog.StatusCode == 403 || errorLog.StatusCode == 404 || errorLog.StatusCode == 409 || errorLog.StatusCode == 400 {
			output = &model.ErrorLog{
				Message:       errorLog.Message,
				SystemMessage: errorLog.SystemMessage,
				Err:           errorLog.Err,
				StatusCode:    errorLog.StatusCode,
			}

		} else if errorLog.StatusCode == 422 {
			output = &model.ErrorLog{
				Message:    "Field Validation",
				Err:        errorLog.Err,
				StatusCode: errorLog.StatusCode,
				Fields:     errorLog.Message,
			}
		}

		log := map[string]interface{}{}
		_ = DecodeMapType(output, &log)
		logrus.WithFields(log).Error(errorLog.Message)
		return output
	}

	return nil
}

func NewWriteLog(errorLog model.ErrorLog) *model.ErrorLog {
	if pc, file, line, ok := runtime.Caller(1); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		var output *model.ErrorLog

		if errorLog.StatusCode == 500 {
			output = &model.ErrorLog{
				Line:          fmt.Sprintf("%d", line),
				Filename:      fmt.Sprintf("%s", file),
				Function:      fmt.Sprintf("%s", funcName),
				Message:       errorLog.Message,
				SystemMessage: errorLog.Err.Error(),
				Err:           errorLog.Err,
				StatusCode:    errorLog.StatusCode,
			}
		} else if errorLog.StatusCode == 401 || errorLog.StatusCode == 403 || errorLog.StatusCode == 404 || errorLog.StatusCode == 409 || errorLog.StatusCode == 400 {
			output = &model.ErrorLog{
				Message:       errorLog.Message,
				SystemMessage: errorLog.SystemMessage,
				Err:           errorLog.Err,
				StatusCode:    errorLog.StatusCode,
			}

		} else if errorLog.StatusCode == 422 {
			output = &model.ErrorLog{
				Message:    "Field Validation",
				Err:        errorLog.Err,
				StatusCode: errorLog.StatusCode,
				Fields:     errorLog.Message,
			}
		}

		log := map[string]interface{}{}
		_ = DecodeMapType(output, &log)
		logrus.WithFields(log).Error(errorLog.Message)
		return output
	}

	return nil
}

func WriteLogConsumer(consumerName string, consumerTopic string, consumerPartition int, consumerOffset int64, consumerKey string, err error, errorCode int, message interface{}) *model.ErrorLog {
	if pc, file, line, ok := runtime.Caller(1); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		var output *model.ErrorLog

		if errorCode == 500 {
			output = &model.ErrorLog{
				Line:              fmt.Sprintf("%d", line),
				Filename:          fmt.Sprintf("%s", file),
				Function:          fmt.Sprintf("%s", funcName),
				Message:           message,
				SystemMessage:     err.Error(),
				ConsumerName:      consumerName,
				ConsumerTopic:     consumerTopic,
				ConsumerPartition: consumerPartition,
				ConsumerOffset:    consumerOffset,
				ConsumerKey:       consumerKey,
				Err:               err,
				StatusCode:        errorCode,
			}
		}

		log := map[string]interface{}{}
		_ = DecodeMapType(output, &log)
		logrus.WithFields(log).Error(message)
		return output
	}

	return nil
}
