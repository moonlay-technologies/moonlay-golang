package routes

import (
	"context"
	"net/http"
	"order-service/app/controllers"
	"order-service/app/middlewares"
	"order-service/app/models/constants"
	kafkadbo "order-service/global/utils/kafka"
	baseModel "order-service/global/utils/model"
	"order-service/global/utils/mongodb"
	"order-service/global/utils/opensearch_dbo"
	"order-service/global/utils/redisdb"

	"github.com/bxcodec/dbresolver"
	"github.com/gin-gonic/gin"
)

func InitHTTPRoute(g *gin.Engine, database dbresolver.DB, redisdb redisdb.RedisInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, kafkaClient kafkadbo.KafkaClientInterface, ctx context.Context) {
	g.GET(constants.HEALTH_CHECK_PATH, func(context *gin.Context) {
		context.JSON(http.StatusOK, map[string]interface{}{"status": http.StatusText(http.StatusOK)})
	})

	basicAuthRootGroup := g.Group("", middlewares.BasicAuthMiddleware())

	// basicAuthRootGroup := g.Group("", gin.BasicAuth(gin.Accounts{
	// 	os.Getenv("AUTHBASIC_USERNAME"): os.Getenv("AUTHBASIC_PASSWORD"),
	// }))

	salesOrderController := controllers.InitHTTPSalesOrderController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		salesOrderControllerGroup := basicAuthRootGroup.Group(constants.SALES_ORDERS_PATH)
		salesOrderControllerGroup.Use()
		{
			salesOrderControllerGroup.POST("", salesOrderController.Create)
			salesOrderControllerGroup.GET("/upload-histories", salesOrderController.GetSOUploadHistories)
			salesOrderControllerGroup.GET("/upload-histories/:id", salesOrderController.GetSoUploadHistoriesById)
			salesOrderControllerGroup.GET("/upload-histories-by-request-id/:id/error-items", salesOrderController.GetSoUploadErrorLogByReqId)
			salesOrderControllerGroup.GET("/upload-histories/:id/error-items", salesOrderController.GetSoUploadErrorLogBySoUploadHistoryId)
			salesOrderControllerGroup.GET("", salesOrderController.Get)
			salesOrderControllerGroup.GET("export", salesOrderController.Export)
			salesOrderControllerGroup.GET("export/sales-order-details", salesOrderController.ExportDetail)
			salesOrderControllerGroup.GET(":so-id", salesOrderController.GetByID)
			salesOrderControllerGroup.GET("details/:so-detail-id", salesOrderController.GetDetailsById)
			salesOrderControllerGroup.GET(":so-id/details", salesOrderController.GetDetailsBySoId)
			salesOrderControllerGroup.PUT(":so-id", salesOrderController.UpdateByID)
			salesOrderControllerGroup.PUT(":so-id/details", salesOrderController.UpdateSODetailBySOID)
			salesOrderControllerGroup.PUT(":so-id/details/:so-detail-id", salesOrderController.UpdateSODetailByID)
			salesOrderControllerGroup.DELETE(":so-id", salesOrderController.DeleteByID)
			salesOrderControllerGroup.DELETE(":so-id/details", salesOrderController.DeleteDetailBySOID)
			salesOrderControllerGroup.DELETE("details/:so-detail-id", salesOrderController.DeleteDetailByID)
			salesOrderControllerGroup.GET("event-logs", salesOrderController.GetSyncToKafkaHistories)
			salesOrderControllerGroup.GET("/journeys", salesOrderController.GetSOJourneys)
			salesOrderControllerGroup.GET(":so-id/journeys", salesOrderController.GetSOJourneyBySoId)
			salesOrderControllerGroup.GET("/retry-to-sync-kafka/:log-id", salesOrderController.RetrySyncToKafka)
		}

		salesOrderDetailControllerGroup := basicAuthRootGroup.Group(constants.SALES_ORDER_DETAIL)
		salesOrderControllerGroup.Use()
		{
			salesOrderDetailControllerGroup.GET("", salesOrderController.GetDetails)
		}
	}

	deliveryOrderController := controllers.InitHTTPDeliveryOrderController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		deliveryOrderControllerGroup := basicAuthRootGroup.Group(constants.DELIVERY_ORDERS_PATH)
		deliveryOrderControllerGroup.Use()
		{
			deliveryOrderControllerGroup.POST("", deliveryOrderController.Create)
			deliveryOrderControllerGroup.GET("upload-histories-all", deliveryOrderController.GetDoUploadHistories)
			deliveryOrderControllerGroup.GET("/upload-histories/:sj-id", deliveryOrderController.GetDoUploadHistoriesById)
			deliveryOrderControllerGroup.GET("/upload-histories-by-request-id/:sj-id/error-items", deliveryOrderController.GetDoUploadErrorLogByReqId)
			deliveryOrderControllerGroup.GET("/upload-histories/:sj-id/error-items", deliveryOrderController.GetDoUploadErrorLogByDoUploadHistoryId)
			deliveryOrderControllerGroup.GET(":id", deliveryOrderController.GetByID)
			deliveryOrderControllerGroup.GET("", deliveryOrderController.Get)
			deliveryOrderControllerGroup.GET("export", deliveryOrderController.Export)
			deliveryOrderControllerGroup.GET("export/delivery-order-details", deliveryOrderController.ExportDetail)
			deliveryOrderControllerGroup.GET("/details", deliveryOrderController.GetDetails)
			deliveryOrderControllerGroup.GET(":id/details", deliveryOrderController.GetDetailsByDoId)
			deliveryOrderControllerGroup.GET(":id/details/:do-detail-id", deliveryOrderController.GetDetailById)
			deliveryOrderControllerGroup.GET("salesmans", deliveryOrderController.GetBySalesmanID)
			deliveryOrderControllerGroup.GET("/sync-to-kafka-histories", deliveryOrderController.GetSyncToKafkaHistories)
			deliveryOrderControllerGroup.GET("/journeys", deliveryOrderController.GetJourneys)
			deliveryOrderControllerGroup.GET(":id/journeys", deliveryOrderController.GetDOJourneysByDoID)
			deliveryOrderControllerGroup.PUT(":id", deliveryOrderController.UpdateByID)
			deliveryOrderControllerGroup.PUT(":id/details", deliveryOrderController.UpdateDeliveryOrderDetailByDeliveryOrderID)
			deliveryOrderControllerGroup.PUT("details/:id", deliveryOrderController.UpdateDeliveryOrderDetailByID)
			deliveryOrderControllerGroup.DELETE(":id", deliveryOrderController.DeleteByID)
			deliveryOrderControllerGroup.DELETE(":id/details", deliveryOrderController.DeleteByID)
			deliveryOrderControllerGroup.DELETE("details/:id", deliveryOrderController.DeleteByID)
			deliveryOrderControllerGroup.GET("/retry-to-sync-kafka/:log-id", deliveryOrderController.RetrySyncToKafka)
		}
	}

	agentController := controllers.InitHTTPAgentController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		agentControllerGroup := basicAuthRootGroup.Group(constants.AGENT_PATH)
		agentControllerGroup.Use()
		{
			agentControllerGroup.GET(":id/"+constants.SALES_ORDERS_PATH, agentController.GetSalesOrders)
			agentControllerGroup.GET(":id/"+constants.DELIVERY_ORDERS_PATH, agentController.GetDeliveryOrders)
		}
	}

	storeController := controllers.InitHTTPStoreController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		storeControllerGroup := basicAuthRootGroup.Group(constants.STORES_PATH)
		storeControllerGroup.Use()
		{
			storeControllerGroup.GET(":id/"+constants.SALES_ORDERS_PATH, storeController.GetSalesOrders)
			storeControllerGroup.GET(":id/"+constants.DELIVERY_ORDERS_PATH, storeController.GetDeliveryOrders)
		}
	}

	salesmanController := controllers.InitHTTPSalesmanController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		salesmanControllerGroup := basicAuthRootGroup.Group(constants.SALESMANS_PATH)
		salesmanControllerGroup.Use()
		{
			salesmanControllerGroup.GET(":id/"+constants.SALES_ORDERS_PATH, salesmanController.GetSalesOrders)
			salesmanControllerGroup.GET(":id/"+constants.DELIVERY_ORDERS_PATH, salesmanController.GetDeliveryOrders)
		}
	}

	hostToHostController := controllers.InitHTTPHostToHostController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		hostToHostControllerGroup := basicAuthRootGroup.Group(constants.HOST_TO_HOST_PATH)
		hostToHostControllerGroup.Use()
		{
			hostToHostControllerGroup.GET("/"+constants.SALES_ORDERS_PATH, hostToHostController.GetSalesOrders)
			hostToHostControllerGroup.GET("/"+constants.DELIVERY_ORDERS_PATH, hostToHostController.GetDeliveryOrders)
		}
	}

	uploadController := controllers.InitHTTPUploadController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		uploadControllerGroup := basicAuthRootGroup.Group("")
		uploadControllerGroup.Use()
		{
			uploadControllerGroup.POST(constants.UPLOAD_SOSJ_PATH, uploadController.UploadSOSJ)
			uploadControllerGroup.POST(constants.UPLOAD_DO_PATH, uploadController.UploadDO)
			uploadControllerGroup.POST(constants.UPLOAD_SO_PATH, uploadController.UploadSO)
			uploadControllerGroup.GET(constants.UPLOAD_SO_PATH+"/retry/:so-upload-history-id", uploadController.RetryUploadSO)
			uploadControllerGroup.GET(constants.UPLOAD_DO_PATH+"/retry/:sj-upload-history-id", uploadController.RetryUploadDO)
			uploadControllerGroup.GET(constants.UPLOAD_SOSJ_PATH+"/retry/:sosj-upload-history-id", uploadController.RetryUploadSOSJ)
			uploadControllerGroup.GET(constants.SOSJ_PATH+"/"+constants.UPLOAD_HISTORIES_PATH, uploadController.GetSosjUploadHistories)
			uploadControllerGroup.GET(constants.SOSJ_PATH+"/"+constants.UPLOAD_HISTORIES_PATH+"/:id/error-items", uploadController.GetSoUploadErrorLogsByReqId)
			uploadControllerGroup.GET(constants.SOSJ_PATH+"/"+constants.UPLOAD_HISTORIES_PATH+"/:id", uploadController.GetSosjUploadHistoryById)
			uploadControllerGroup.GET(constants.SOSJ_PATH+"/"+constants.UPLOAD_HISTORIES_PATH+"/items/:sosj-upload-history-id", uploadController.GetSosjUploadErrorLogsBySosjUploadHistoryId)
		}
	}

	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, baseModel.Response{
			StatusCode: http.StatusNotFound,
			Error: &baseModel.ErrorLog{
				Message:       "Not Found",
				SystemMessage: "Not Found",
			},
		})
	})
	//
	//oauthRootGroup.Use(middlewares.OauthMiddleware(mongod))
	//{
	//
	//}
}
