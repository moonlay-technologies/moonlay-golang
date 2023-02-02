package helper

import (
	"fmt"
	"order-service/global/utils/model"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func WriteLog(err error, errorCode int, message interface{}) *model.ErrorLog {
	if pc, file, line, ok := runtime.Caller(1); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		var output *model.ErrorLog
		output.StatusCode = errorCode
		output.Err = err

		if errorCode == 500 {
			output.Line = fmt.Sprintf("%d", line)
			output.Filename = fmt.Sprintf("%s", file)
			output.Function = fmt.Sprintf("%s", funcName)
			if message == nil {
				output.Message = "Something went wrong, please try again later"
			} else {
				output.Message = message
			}
			output.SystemMessage = err.Error()
		} else if errorCode == 401 || errorCode == 403 || errorCode == 404 || errorCode == 409 || errorCode == 400 {
			output.Message = message
			output.SystemMessage = err.Error()

		} else if errorCode == 422 {
			output.Message = "Field Validation"
			output.Fields = message
		}

		log := map[string]interface{}{}
		_ = DecodeMapType(output, &log)
		logrus.WithFields(log).Error(message)
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
