package usecases

//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/getsentry/sentry-go"
//	"net/http"
//	"os"
//	"order-service/app/models"
//	"order-service/global/utils/helper"
//	baseModel "order-service/global/utils/model"
//	"time"
//)
//
//type ApiOauthUseCaseInterface interface {
//	ValidateToken(token string, requestID string, result chan *models.ValidAccessTokenWithClient)
//}
//
//type apiOauthUseCase struct {
//	logUseCase LogUseCaseInterface
//}
//
//func InitApiOauthUseCaseInterface(logUseCase LogUseCaseInterface) ApiOauthUseCaseInterface {
//	return &apiOauthUseCase{
//		logUseCase: logUseCase,
//	}
//}
//
//func (u *apiOauthUseCase) ValidateToken(token string, requestID string, resultChan chan *models.ValidAccessTokenWithClient) {
//	now := time.Now()
//	response := &models.ValidAccessTokenWithClient{}
//	//get client
//	tokenAuthBasic := helper.BasicAuth(os.Getenv("AUTHBASIC_OAUTH_USERNAME"), os.Getenv("AUTHBASIC_OAUTH_PASSWORD"))
//	httpBody := map[string]interface{}{
//		"access_token": token,
//	}
//
//	httpJsonBody, err := json.Marshal(httpBody)
//
//	if err != nil {
//		errStr := fmt.Sprintf("Error Encode Json Body Oauth to Oauth API %s", err.Error())
//		helper.SetSentryError(err, errStr, sentry.LevelError)
//		response.Error = err
//		response.StatusCode = http.StatusBadRequest
//		resultChan <- response
//		return
//	}
//
//	httpRequestOption := helper.Options{
//		Method: "POST",
//		Headers: map[string]string{
//			"Authorization": fmt.Sprintf("Basic %s", tokenAuthBasic),
//		},
//		Body:        httpJsonBody,
//		ContentType: "application/json",
//		URL:         fmt.Sprintf("%s/api/v1/oauth/verify-access-token", os.Getenv("OAUTH_URL")),
//	}
//
//	validateTokenRest := helper.POST(&httpRequestOption)
//	validateTokenMapResponse := &baseModel.Response{}
//
//	if validateTokenRest.StatusCode == 200 {
//		validateAccessTokenWithClient := models.ValidAccessTokenWithClient{}
//		err = json.Unmarshal(validateTokenRest.Body, &validateTokenMapResponse)
//
//		if err != nil {
//			errStr := fmt.Sprintf("Error Decode Json Response Validate Token HTTP to API Oauth %s", err.Error())
//			helper.SetSentryError(err, errStr, sentry.LevelError)
//			response.Error = err
//			response.StatusCode = http.StatusBadRequest
//			resultChan <- response
//			return
//		}
//
//		err = helper.DecodeMapType(validateTokenMapResponse.Data, &validateAccessTokenWithClient)
//
//		if err != nil {
//			errStr := fmt.Sprintf("Error Decode Map to Struct Response Validate Token HTTP to API Oauth %s", err.Error())
//			helper.SetSentryError(err, errStr, sentry.LevelError)
//			response.Error = err
//			resultChan <- response
//			return
//		}
//
//		log := &models.Log{
//			RequestID:    requestID,
//			LogType:      "http_request_validate_oauth_token_log",
//			ServiceName:  "thirdparty-service",
//			FromService:  "thirdpary-service",
//			ToService:    "oauth-service",
//			RequestBody:  nil,
//			ErrorMessage: "",
//			StatusCode:   validateTokenRest.StatusCode,
//			Response:     nil,
//			CreatedAt:    &now,
//		}
//
//		_, err = u.logUseCase.Insert(log)
//
//		if err != nil {
//			response = nil
//			response.Error = err
//			response.StatusCode = http.StatusInternalServerError
//			resultChan <- response
//			return
//		}
//
//		response = &validateAccessTokenWithClient
//		response.StatusCode = http.StatusOK
//		response.Error = nil
//		resultChan <- response
//		return
//	} else {
//		err = json.Unmarshal(validateTokenRest.Body, validateTokenMapResponse)
//
//		if err != nil {
//			errStr := fmt.Sprintf("Error Decode Json Response Validate Token HTTP to API Oauth %s", err.Error())
//			helper.SetSentryError(err, errStr, sentry.LevelError)
//			fmt.Println(err)
//			response.Error = err
//			response.StatusCode = http.StatusBadRequest
//			resultChan <- response
//			return
//		}
//
//		errorMessage := helper.NewError(validateTokenMapResponse.Error.(string))
//		validateTokenRest.Error = errorMessage
//
//		log := &models.Log{
//			RequestID:    requestID,
//			LogType:      "http_request_validate_oauth_token_log",
//			ServiceName:  "thirdparty-service",
//			FromService:  "thirdpary-service",
//			ToService:    "oauth-service",
//			RequestBody:  nil,
//			ErrorMessage: errorMessage.Error(),
//			StatusCode:   validateTokenRest.StatusCode,
//			Response:     nil,
//			CreatedAt:    &now,
//		}
//
//		_, err = u.logUseCase.Insert(log)
//
//		if err != nil {
//			response.StatusCode = http.StatusInternalServerError
//			response.Error = err
//			resultChan <- response
//			return
//		}
//
//		errStr := fmt.Sprintf("Error validate token oauth http request usecase %s", validateTokenRest.Error.Error())
//		helper.SetSentryError(validateTokenRest.Error, errStr, sentry.LevelError)
//		response.Error = errorMessage
//		response.StatusCode = validateTokenRest.StatusCode
//		resultChan <- response
//		return
//	}
//}
