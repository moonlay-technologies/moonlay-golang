package routes

import (
	"context"
	"github.com/bxcodec/dbresolver"
	"github.com/gin-gonic/gin"
	"os"
	"poc-order-service/app/controllers"
	kafkadbo "poc-order-service/global/utils/kafka"
	"poc-order-service/global/utils/mongodb"
	"poc-order-service/global/utils/opensearch_dbo"
	"poc-order-service/global/utils/redisdb"
)

func InitHTTPRoute(g *gin.Engine, database dbresolver.DB, redisdb redisdb.RedisInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, kafkaClient kafkadbo.KafkaClientInterface, ctx context.Context) {
	g.GET("health-check", func(context *gin.Context) {
		context.JSON(200, map[string]interface{}{"status": "OK"})
	})

	basicAuthRootGroup := g.Group("", gin.BasicAuth(gin.Accounts{
		os.Getenv("AUTHBASIC_USERNAME"): os.Getenv("AUTHBASIC_PASSWORD"),
	}))

	salesOrderController := controllers.InitHTTPSalesOrderController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		salesOrderControllerGroup := basicAuthRootGroup.Group("sales-orders")
		salesOrderControllerGroup.Use()
		{
			salesOrderControllerGroup.GET("", salesOrderController.Get)
			salesOrderControllerGroup.GET(":id", salesOrderController.GetByID)
			salesOrderControllerGroup.POST("", salesOrderController.Create)
		}
	}

	deliveryOrderController := controllers.InitHTTPDeliveryOrderController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		deliveryOrderControllerGroup := basicAuthRootGroup.Group("delivery-orders")
		deliveryOrderControllerGroup.Use()
		{
			deliveryOrderControllerGroup.POST("", deliveryOrderController.Create)
			deliveryOrderControllerGroup.GET(":id", deliveryOrderController.GetByID)
			deliveryOrderControllerGroup.GET("", deliveryOrderController.Get)
		}
	}

	agentController := controllers.InitHTTPAgentController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		agentControllerGroup := basicAuthRootGroup.Group("agents")
		agentControllerGroup.Use()
		{
			agentControllerGroup.GET(":id/sales-orders", agentController.GetSalesOrders)
			agentControllerGroup.GET(":id/delivery-orders", agentController.GetDeliveryOrders)
		}
	}

	storeController := controllers.InitHTTPStoreController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		storeControllerGroup := basicAuthRootGroup.Group("stores")
		storeControllerGroup.Use()
		{
			storeControllerGroup.GET(":id/sales-orders", storeController.GetSalesOrders)
			storeControllerGroup.GET(":id/delivery-orders", storeController.GetDeliveryOrders)
		}
	}

	salesmanController := controllers.InitHTTPSalesmanController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		salesmanControllerGroup := basicAuthRootGroup.Group("salesmans")
		salesmanControllerGroup.Use()
		{
			salesmanControllerGroup.GET(":id/sales-orders", salesmanController.GetSalesOrders)
			salesmanControllerGroup.GET(":id/delivery-orders", salesmanController.GetDeliveryOrders)
		}
	}

	//
	//oauthRootGroup.Use(middlewares.OauthMiddleware(mongod))
	//{
	//
	//}
}
