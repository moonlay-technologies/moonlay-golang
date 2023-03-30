package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/global/utils/helper"
	baseModel "order-service/global/utils/model"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			respondWithError(http.StatusUnauthorized, c)
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		if string(payload) != fmt.Sprintf("%s:%s", os.Getenv("AUTHBASIC_USERNAME"), os.Getenv("AUTHBASIC_PASSWORD")) {
			respondWithError(http.StatusUnauthorized, c)
			return
		}
		c.Request.Header.Set("token", "eyJ1c2VyX2lkIjo1NzgzLCJhZ2VudF9pZCI6NCwidXNlcl9lbWFpbCI6ImNtc19wcmluc2lwYWxAZ21haWwuY29tIiwidXNlcl9yb2xlX3NsdWciOiJkYm8tYWRtaW5pc3RyYXRvciIsInVzZXJfcm9sZV9jYXRlZ29yeSI6ImFkbWluIiwiZmlyc3RfbmFtZSI6IkNNUyBQcmluc2lwYWwiLCJsYXN0X25hbWUiOiIiLCJhbGciOiJIUzI1NiJ9==")
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !authenticateUser(c.Request.Header.Get("token"), c) {
			respondWithError(http.StatusUnauthorized, c)
			return
		}

		c.Next()
	}
}

func authenticateUser(tokenString string, c *gin.Context) bool {
	user := models.UserClaims{}
	token, _ := base64.StdEncoding.DecodeString(tokenString)
	json.Unmarshal(token, &user)
	c.Set("user", &user)
	return true
}

func respondWithError(code int, c *gin.Context) {
	var result baseModel.Response

	err := helper.NewError(http.StatusText(http.StatusUnauthorized))
	errorLogData := helper.WriteLog(err, http.StatusUnauthorized, nil)
	result.StatusCode = code
	result.Error = errorLogData
	c.JSON(result.StatusCode, result)
	c.Abort()
}
