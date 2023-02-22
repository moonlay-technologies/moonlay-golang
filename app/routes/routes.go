package routes

import (
	"context"
	"net/http"
	"order-service/app/controllers"
	"order-service/app/models/constants"
	kafkadbo "order-service/global/utils/kafka"
	"order-service/global/utils/mongodb"
	"order-service/global/utils/opensearch_dbo"
	"order-service/global/utils/redisdb"
	"os"

	"github.com/bxcodec/dbresolver"
	"github.com/gin-gonic/gin"
)

func InitHTTPRoute(g *gin.Engine, database dbresolver.DB, redisdb redisdb.RedisInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, kafkaClient kafkadbo.KafkaClientInterface, ctx context.Context) {
	g.GET(constants.HEALTH_CHECK, func(context *gin.Context) {
		context.JSON(http.StatusOK, map[string]interface{}{"status": http.StatusText(http.StatusOK)})
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
			salesOrderControllerGroup.GET(":so-id", salesOrderController.GetByID)
			salesOrderControllerGroup.POST("", salesOrderController.Create)
			salesOrderControllerGroup.PUT(":so-id", salesOrderController.UpdateByID)
			salesOrderControllerGroup.PUT(":so-id/details/:id", salesOrderController.UpdateSODetailByID)
			salesOrderControllerGroup.PUT(":so-id/details", salesOrderController.UpdateSODetailBySOID)
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
			deliveryOrderControllerGroup.PUT(":id", deliveryOrderController.UpdateByID)
			deliveryOrderControllerGroup.PUT("/details/:id", deliveryOrderController.UpdateDeliveryOrderDetailByID)
			deliveryOrderControllerGroup.PUT("/:id/details", deliveryOrderController.UpdateDeliveryOrderDetailByDeliveryOrderID)
			deliveryOrderControllerGroup.GET(":id", deliveryOrderController.GetByID)
			deliveryOrderControllerGroup.GET("", deliveryOrderController.Get)
			deliveryOrderControllerGroup.GET("/salesmans", deliveryOrderController.GetBySalesmanID)
		}
	}

	agentController := controllers.InitHTTPAgentController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		agentControllerGroup := basicAuthRootGroup.Group(constants.AGENT)
		agentControllerGroup.Use()
		{
			agentControllerGroup.GET(":id/"+constants.SALES_ORDERS_PATH, agentController.GetSalesOrders)
			agentControllerGroup.GET(":id/"+constants.DELIVERY_ORDERS_PATH, agentController.GetDeliveryOrders)
		}
	}

	storeController := controllers.InitHTTPStoreController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		storeControllerGroup := basicAuthRootGroup.Group(constants.STORES)
		storeControllerGroup.Use()
		{
			storeControllerGroup.GET(":id/"+constants.SALES_ORDERS_PATH, storeController.GetSalesOrders)
			storeControllerGroup.GET(":id/"+constants.DELIVERY_ORDERS_PATH, storeController.GetDeliveryOrders)
		}
	}

	salesmanController := controllers.InitHTTPSalesmanController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		salesmanControllerGroup := basicAuthRootGroup.Group(constants.SALESMANS)
		salesmanControllerGroup.Use()
		{
			salesmanControllerGroup.GET(":id/"+constants.SALES_ORDERS_PATH, salesmanController.GetSalesOrders)
			salesmanControllerGroup.GET(":id/"+constants.DELIVERY_ORDERS_PATH, salesmanController.GetDeliveryOrders)
		}
	}

	hostToHostController := controllers.InitHTTPHostToHostController(database, redisdb, mongodbClient, kafkaClient, opensearchClient, ctx)
	basicAuthRootGroup.Use()
	{
		hostToHostControllerGroup := basicAuthRootGroup.Group(constants.HOST_TO_HOST)
		hostToHostControllerGroup.Use()
		{
			hostToHostControllerGroup.GET("/"+constants.SALES_ORDERS_PATH, hostToHostController.GetSalesOrders)
		}
	}

	//
	//oauthRootGroup.Use(middlewares.OauthMiddleware(mongod))
	//{
	//
	//}
}
