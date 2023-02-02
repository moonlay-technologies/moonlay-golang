package routes

import (
	"context"
	"os"
	"poc-order-service/app/controllers"
	"poc-order-service/app/models/constants"
	kafkadbo "poc-order-service/global/utils/kafka"
	"poc-order-service/global/utils/mongodb"
	"poc-order-service/global/utils/opensearch_dbo"
	"poc-order-service/global/utils/redisdb"

	"github.com/bxcodec/dbresolver"
	"github.com/gin-gonic/gin"
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
		salesOrderControllerGroup := basicAuthRootGroup.Group(constants.SALES_ORDERS_PATH)
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
		deliveryOrderControllerGroup := basicAuthRootGroup.Group(constants.DELIVERY_ORDERS_PATH)
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
			agentControllerGroup.GET(":id/"+constants.SALES_ORDERS_PATH, agentController.GetSalesOrders)
			agentControllerGroup.GET(":id/"+constants.DELIVERY_ORDERS_PATH, agentController.GetDeliveryOrders)
		}
	}

	storeController := controllers.InitHTTPStoreController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		storeControllerGroup := basicAuthRootGroup.Group("stores")
		storeControllerGroup.Use()
		{
			storeControllerGroup.GET(":id/"+constants.SALES_ORDERS_PATH, storeController.GetSalesOrders)
			storeControllerGroup.GET(":id/"+constants.DELIVERY_ORDERS_PATH, storeController.GetDeliveryOrders)
		}
	}

	salesmanController := controllers.InitHTTPSalesmanController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		salesmanControllerGroup := basicAuthRootGroup.Group("salesmans")
		salesmanControllerGroup.Use()
		{
			salesmanControllerGroup.GET(":id/"+constants.SALES_ORDERS_PATH, salesmanController.GetSalesOrders)
			salesmanControllerGroup.GET(":id/"+constants.DELIVERY_ORDERS_PATH, salesmanController.GetDeliveryOrders)
		}
	}

	//
	//oauthRootGroup.Use(middlewares.OauthMiddleware(mongod))
	//{
	//
	//}
}
