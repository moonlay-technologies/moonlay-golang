package handlers

import (
	"context"
	"fmt"
	"order-service/app/consumer"
	"order-service/app/models/constants"
	kafkadbo "order-service/global/utils/kafka"
	"order-service/global/utils/mongodb"
	"order-service/global/utils/opensearch_dbo"
	"order-service/global/utils/redisdb"
	"os"
	"sync"

	"github.com/bxcodec/dbresolver"
)

func MainConsumerHandler(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) {
	wg := sync.WaitGroup{}
	switch args[1] {
	case constants.CREATE_SALES_ORDER_TOPIC:
		wg.Add(1)
		salesOrderConsumer := consumer.InitCreateSalesOrderConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go salesOrderConsumer.ProcessMessage()
		break
	case constants.UPDATE_SALES_ORDER_TOPIC:
		wg.Add(1)
		salesOrderConsumer := consumer.InitUpdateSalesOrderConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go salesOrderConsumer.ProcessMessage()
		break
	case constants.DELETE_SALES_ORDER_TOPIC:
		wg.Add(1)
		salesOrderConsumer := consumer.InitDeleteSalesOrderConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go salesOrderConsumer.ProcessMessage()
		break
	case constants.CREATE_DELIVERY_ORDER_TOPIC:
		wg.Add(1)
		deliveryOrderConsumer := consumer.InitCreateDeliveryOrderConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go deliveryOrderConsumer.ProcessMessage()
		break
	case constants.UPDATE_DELIVERY_ORDER_TOPIC:
		wg.Add(1)
		deliveryOrderConsumer := consumer.InitUpdateDeliveryOrderConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go deliveryOrderConsumer.ProcessMessage()
		break
	case constants.DELETE_DELIVERY_ORDER_TOPIC:
		wg.Add(1)
		deliveryOrderConsumer := consumer.InitDeleteDeliveryOrderConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go deliveryOrderConsumer.ProcessMessage()
		break
	case constants.UPDATE_SALES_ORDER_DETAIL_TOPIC:
		wg.Add(1)
		salesOrderDetailConsumer := consumer.InitUpdateSalesOrderDetailConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go salesOrderDetailConsumer.ProcessMessage()
		break
	case constants.UPDATE_DELIVERY_ORDER_DETAIL_TOPIC:
		wg.Add(1)
		salesOrderDetailConsumer := consumer.InitUpdateDeliveryOrderDetailConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go salesOrderDetailConsumer.ProcessMessage()
		break
	default:
		fmt.Println("Choose Command Type You Want")
		os.Exit(0)
	}

	wg.Wait()

}
