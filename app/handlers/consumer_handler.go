package handlers

import (
	"context"
	"fmt"
	"github.com/bxcodec/dbresolver"
	"os"
	"poc-order-service/app/consumer"
	kafkadbo "poc-order-service/global/utils/kafka"
	"poc-order-service/global/utils/mongodb"
	"poc-order-service/global/utils/opensearch_dbo"
	"poc-order-service/global/utils/redisdb"
	"sync"
)

func MainConsumerHandler(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) {
	wg := sync.WaitGroup{}
	switch args[1] {
	case "create-sales-order":
		wg.Add(1)
		salesOrderConsumer := consumer.InitCreateSalesOrderConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go salesOrderConsumer.ProcessMessage()
		break
	case "update-sales-order":
		wg.Add(1)
		salesOrderConsumer := consumer.InitUpdateSalesOrderConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go salesOrderConsumer.ProcessMessage()
		break
	case "create-delivery-order":
		wg.Add(1)
		deliveryOrderConsumer := consumer.InitCreateDeliveryOrderConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go deliveryOrderConsumer.ProcessMessage()
		break
	default:
		fmt.Println("Choose Command Type You Want")
		os.Exit(0)
	}

	wg.Wait()

}
