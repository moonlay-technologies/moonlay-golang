package consumer

import (
	"context"
	"order-service/app/repositories"
	mongoRepo "order-service/app/repositories/mongod"
	openSearchRepo "order-service/app/repositories/open_search"
	"order-service/app/usecases"
	kafkadbo "order-service/global/utils/kafka"
	"order-service/global/utils/mongodb"
	"order-service/global/utils/opensearch_dbo"
	"order-service/global/utils/redisdb"

	"github.com/bxcodec/dbresolver"
)

func InitCreateSalesOrderConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) CreateSalesOrderConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository)
	handler := InitCreateSalesOrderConsumerHandlerInterface(kafkaClient, salesOrderLogRepository, salesOrderOpenSearchUseCase, database, ctx, args)
	return handler
}

func InitUpdateSalesOrderConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UpdateSalesOrderConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository)
	handler := InitUpdateSalesOrderConsumerHandlerInterface(kafkaClient, salesOrderLogRepository, salesOrderOpenSearchUseCase, database, ctx, args)
	return handler
}

func InitDeleteSalesOrderConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UpdateSalesOrderConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository)
	handler := InitDeleteSalesOrderConsumerHandlerInterface(kafkaClient, salesOrderLogRepository, salesOrderOpenSearchUseCase, database, ctx, args)
	return handler
}

func InitDeleteSalesOrderDetailConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UpdateSalesOrderConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository)
	handler := InitDeleteSalesOrderDetailConsumerHandlerInterface(kafkaClient, salesOrderLogRepository, salesOrderRepository, salesOrderOpenSearchUseCase, database, ctx, args)
	return handler
}

func InitCreateDeliveryOrderConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) CreateDeliveryOrderConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	deliveryOrderDetailOpenSearchRepository := openSearchRepo.InitDeliveryOrderDetailOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository)
	deliveryOrderOpenSearchUseCase := usecases.InitDeliveryOrderConsumerUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, deliveryOrderOpenSearchRepository, deliveryOrderDetailOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderOpenSearchUseCase, database, ctx)
	handler := InitCreateDeliveryOrderConsumerHandlerInterface(kafkaClient, deliveryOrderLogRepository, deliveryOrderOpenSearchUseCase, database, ctx, args)
	return handler
}

func InitUpdateDeliveryOrderConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UpdateSalesOrderConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	deliveryOrderDetailOpenSearchRepository := openSearchRepo.InitDeliveryOrderDetailOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository)
	deliveryOrderOpenSearchUseCase := usecases.InitDeliveryOrderConsumerUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, deliveryOrderOpenSearchRepository, deliveryOrderDetailOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderOpenSearchUseCase, database, ctx)
	handler := InitUpdateDeliveryOrderConsumerHandlerInterface(kafkaClient, deliveryOrderLogRepository, deliveryOrderOpenSearchUseCase, database, ctx, args)
	return handler
}

func InitDeleteDeliveryOrderConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UpdateSalesOrderConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	deliveryOrderDetailOpenSearchRepository := openSearchRepo.InitDeliveryOrderDetailOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository)
	deliveryOrderOpenSearchUseCase := usecases.InitDeliveryOrderConsumerUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, deliveryOrderOpenSearchRepository, deliveryOrderDetailOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderOpenSearchUseCase, database, ctx)
	handler := InitDeleteDeliveryOrderConsumerHandlerInterface(kafkaClient, deliveryOrderLogRepository, deliveryOrderOpenSearchUseCase, database, ctx, args)
	return handler
}

func InitUpdateSalesOrderDetailConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UpdateSalesOrderDetailConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository)
	handler := InitUpdateSalesOrderDetailConsumerHandlerInterface(kafkaClient, salesOrderLogRepository, salesOrderOpenSearchUseCase, database, ctx, args)
	return handler
}

func InitUpdateDeliveryOrderDetailConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UpdateSalesOrderConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	deliveryOrderDetailOpenSearchRepository := openSearchRepo.InitDeliveryOrderDetailOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository)
	deliveryOrderOpenSearchUseCase := usecases.InitDeliveryOrderConsumerUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, deliveryOrderOpenSearchRepository, deliveryOrderDetailOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderOpenSearchUseCase, database, ctx)
	handler := InitUpdateDeliveryOrderDetailConsumerHandlerInterface(kafkaClient, deliveryOrderLogRepository, deliveryOrderOpenSearchUseCase, database, ctx, args)
	return handler
}

func InitDeleteDeliveryOrderDetailConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UpdateSalesOrderConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	deliveryOrderDetailOpenSearchRepository := openSearchRepo.InitDeliveryOrderDetailOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository)
	deliveryOrderOpenSearchUseCase := usecases.InitDeliveryOrderConsumerUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, deliveryOrderOpenSearchRepository, deliveryOrderDetailOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderOpenSearchUseCase, database, ctx)
	handler := InitDeleteDeliveryOrderDetailConsumerHandlerInterface(kafkaClient, deliveryOrderLogRepository, deliveryOrderOpenSearchUseCase, database, ctx, args)
	return handler
}
