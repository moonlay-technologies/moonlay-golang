package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"poc-order-service/app/middlewares"
	"poc-order-service/app/routes"
	kafkadbo "poc-order-service/global/utils/kafka"
	"poc-order-service/global/utils/mongodb"
	"poc-order-service/global/utils/opensearch_dbo"
	"poc-order-service/global/utils/redisdb"
	"strconv"

	"github.com/bxcodec/dbresolver"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func MainHttpHandler(database dbresolver.DB, redisdb redisdb.RedisInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, kafkaClient kafkadbo.KafkaClientInterface, ctx context.Context) {
	g := gin.Default()
	g.Use(middlewares.CORSMiddleware(), middlewares.JSONMiddleware(), RequestId())

	routes.InitHTTPRoute(g, database, redisdb, mongodbClient, opensearchClient, kafkaClient, ctx)
	useSSL, err := strconv.ParseBool(os.Getenv("USE_SSL"))
	addr := fmt.Sprintf(":%s", os.Getenv("MAIN_PORT"))

	if err != nil || useSSL {
		g.RunTLS(addr, os.Getenv("PUBLIC_SSL_PATH"), os.Getenv("PRIVATE_SSL_PATH"))
	} else {
		err = http.ListenAndServe(addr, g)
	}
}

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestID := c.Request.Header.Get("X-Request-Id")

		// Create request id with UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Expose it for use in the application
		c.Set("RequestId", requestID)
		// Set X-Request-Id header
		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Next()
	}
}
