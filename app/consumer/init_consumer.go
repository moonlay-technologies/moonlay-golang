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
	orderSourceRepository := repositories.InitOrderSourceRepository(database, redisdb)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	handler := InitCreateSalesOrderConsumerHandlerInterface(kafkaClient, salesOrderLogRepository, salesOrderUseCase, database, ctx, args)
	return handler
}

func InitUpdateSalesOrderConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UpdateSalesOrderConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	orderSourceRepository := repositories.InitOrderSourceRepository(database, redisdb)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	handler := InitUpdateSalesOrderConsumerHandlerInterface(kafkaClient, salesOrderLogRepository, salesOrderUseCase, database, ctx, args)
	return handler
}

func InitCreateDeliveryOrderConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) CreateDeliveryOrderConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	orderSourceRepository := repositories.InitOrderSourceRepository(database, redisdb)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	deliveryOrderDetailRepository := repositories.InitDeliveryOrderDetailRepository(database, redisdb)
	warehouseRepository := repositories.InitWarehouseRepository(database, redisdb)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	deliveryOrderUseCase := usecases.InitDeliveryOrderUseCaseInterface(deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, warehouseRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, userRepository, salesmanRepository, deliveryOrderLogRepository, deliveryOrderOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderUseCase, kafkaClient, database, ctx)
	handler := InitCreateDeliveryOrderConsumerHandlerInterface(kafkaClient, deliveryOrderLogRepository, salesOrderUseCase, deliveryOrderUseCase, database, ctx, args)
	return handler
}

func InitUpdateDeliveryOrderConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UpdateSalesOrderConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	orderSourceRepository := repositories.InitOrderSourceRepository(database, redisdb)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	deliveryOrderDetailRepository := repositories.InitDeliveryOrderDetailRepository(database, redisdb)
	warehouseRepository := repositories.InitWarehouseRepository(database, redisdb)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	deliveryOrderUseCase := usecases.InitDeliveryOrderUseCaseInterface(deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, warehouseRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, userRepository, salesmanRepository, deliveryOrderLogRepository, deliveryOrderOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderUseCase, kafkaClient, database, ctx)
	handler := InitUpdateDeliveryOrderConsumerHandlerInterface(kafkaClient, deliveryOrderLogRepository, salesOrderUseCase, deliveryOrderUseCase, database, ctx, args)
	return handler
}
