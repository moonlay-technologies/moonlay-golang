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
	case constants.DELETE_SALES_ORDER_DETAIL_TOPIC:
		wg.Add(1)
		salesOrderConsumer := consumer.InitDeleteSalesOrderDetailConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go salesOrderConsumer.ProcessMessage()
		break
	case constants.CREATE_DELIVERY_ORDER_TOPIC_TMP:
		wg.Add(1)
		deliveryOrderConsumer := consumer.InitCreateDeliveryOrderConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go deliveryOrderConsumer.ProcessMessage()
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
		deliveryOrderDetailConsumer := consumer.InitUpdateDeliveryOrderDetailConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go deliveryOrderDetailConsumer.ProcessMessage()
		break
	case constants.EXPORT_SALES_ORDER_TOPIC:
		wg.Add(1)
		salesOrderConsumer := consumer.InitExportSalesOrderConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go salesOrderConsumer.ProcessMessage()
	case constants.EXPORT_SALES_ORDER_DETAIL_TOPIC:
		wg.Add(1)
		salesOrderConsumer := consumer.InitExportSalesOrderDetailConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go salesOrderConsumer.ProcessMessage()
	case constants.UPLOAD_SO_FILE_TOPIC:
		wg.Add(1)
		uploadSOFileConsumer := consumer.InitUploadSOFileConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go uploadSOFileConsumer.ProcessMessage()
		break
	case constants.UPLOAD_SO_ITEM_TOPIC:
		wg.Add(1)
		uploadSOItemConsumer := consumer.InitUploadSOItemConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go uploadSOItemConsumer.ProcessMessage()
		break
	case constants.UPLOAD_SOSJ_FILE_TOPIC:
		wg.Add(1)
		uploadSOSJFileConsumer := consumer.InitUploadSOSJFileConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go uploadSOSJFileConsumer.ProcessMessage()
		break
	case constants.UPLOAD_SOSJ_ITEM_TOPIC:
		wg.Add(1)
		uploadSOSJItemConsumer := consumer.InitUploadSOSJItemConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go uploadSOSJItemConsumer.ProcessMessage()
	case constants.DELETE_DELIVERY_ORDER_DETAIL_TOPIC:
		wg.Add(1)
		deliveryOrderDetailConsumer := consumer.InitDeleteDeliveryOrderDetailConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go deliveryOrderDetailConsumer.ProcessMessage()
		break
	case constants.EXPORT_DELIVERY_ORDER_TOPIC:
		wg.Add(1)
		deliveryOrderConsumer := consumer.InitExportDeliveryOrderConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go deliveryOrderConsumer.ProcessMessage()
	case constants.EXPORT_DELIVERY_ORDER_DETAIL_TOPIC:
		wg.Add(1)
		deliveryOrderConsumer := consumer.InitExportDeliveryOrderDetailConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go deliveryOrderConsumer.ProcessMessage()
	case constants.UPLOAD_DO_FILE_TOPIC:
		wg.Add(1)
		uploadDOFileConsumer := consumer.InitUploadDOFileConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go uploadDOFileConsumer.ProcessMessage()
		break
	case constants.UPLOAD_DO_ITEM_TOPIC:
		wg.Add(1)
		uploadDOItemConsumer := consumer.InitUploadDOItemConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
		go uploadDOItemConsumer.ProcessMessage()
		break
	default:
		fmt.Println("Choose Command Type You Want")
		os.Exit(0)
	}

	wg.Wait()

}
