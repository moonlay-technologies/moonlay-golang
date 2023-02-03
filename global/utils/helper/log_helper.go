package helper

import (
	"fmt"
	"poc-order-service/global/utils/model"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func WriteLog(err error, errorCode int, message interface{}) *model.ErrorLog {
	if pc, file, line, ok := runtime.Caller(1); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		var output *model.ErrorLog

		if errorCode == 500 {
			output = &model.ErrorLog{
				Line:          fmt.Sprintf("%d", line),
				Filename:      fmt.Sprintf("%s", file),
				Function:      fmt.Sprintf("%s", funcName),
				Message:       message,
				SystemMessage: err.Error(),
				Err:           err,
				StatusCode:    errorCode,
			}
		} else if errorCode == 401 || errorCode == 403 || errorCode == 404 || errorCode == 409 || errorCode == 400 {
			output = &model.ErrorLog{
				Message:       message,
				SystemMessage: err.Error(),
				Err:           err,
				StatusCode:    errorCode,
			}

		} else if errorCode == 422 {
			output = &model.ErrorLog{
				Message:    "Field Validation",
				Err:        err,
				StatusCode: errorCode,
				Fields:     message,
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
