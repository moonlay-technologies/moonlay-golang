package middlewares

//
//import (
//	"github.com/gin-gonic/gin"
//	"github.com/sirupsen/logrus"
//	"net/http"
//	"order-service/app/models"
//	"order-service/app/repositories"
//	"order-service/app/usecases"
//	"order-service/global/utils/helper"
//	baseModel "order-service/global/utils/model"
//	"order-service/global/utils/mongod"
//	"strconv"
//)
//
//func OauthMiddleware(mongod mongod.MongoDBInterface) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		var result baseModel.Response
//		logRepository := repositories.InitLogRepository(logrus.Logger{}, mongod)
//		logUseCase := usecases.InitLogUseCaseInterface(logRepository)
//		oauthUseCase := usecases.InitApiOauthUseCaseInterface(logUseCase)
//
//		authorizationToken := c.Request.Header.Get("Authorization")
//
//		if len(authorizationToken) == 0 {
//			result.StatusCode = http.StatusUnauthorized
//			result.Error = helper.NewError("Unauthorized").Error()
//			c.Status(http.StatusUnauthorized)
//			c.AbortWithStatusJSON(http.StatusUnauthorized, result)
//			return
//		}
//
//		token := helper.GetAuthorizationValue(authorizationToken)
//		requestID, _ := c.Get("RequestId")
//
//		resultOauthValidateResponseChan := make(chan *models.ValidAccessTokenWithClient)
//		go oauthUseCase.ValidateToken(token, requestID.(string), resultOauthValidateResponseChan)
//		resultOauthValidateResponse := <-resultOauthValidateResponseChan
//
//		if resultOauthValidateResponse.Error != nil {
//			result.StatusCode = resultOauthValidateResponse.StatusCode
//			result.Error = resultOauthValidateResponse.Error.Error()
//			c.Status(resultOauthValidateResponse.StatusCode)
//			c.AbortWithStatusJSON(resultOauthValidateResponse.StatusCode, result)
//			return
//		}
//
//		managerAppAgentIDStr := strconv.Itoa(resultOauthValidateResponse.Client.ManagerAppAgentID)
//		c.Set("ClientID", resultOauthValidateResponse.Client.ClientID)
//		c.Set("AccessToken", token)
//		c.Set("ManagerAppAgentID", managerAppAgentIDStr)
//		c.Writer.Header().Set("Content-Type", "application/json")
//		c.Next()
//	}
//}
