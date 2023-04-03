package consumer

import (
	"context"
	"order-service/app/middlewares"
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
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository, categoryRepository)
	handler := InitCreateSalesOrderConsumerHandlerInterface(kafkaClient, salesOrderLogRepository, salesOrderOpenSearchUseCase, database, ctx, args)
	return handler
}

func InitUpdateSalesOrderConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UpdateSalesOrderConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository, categoryRepository)
	handler := InitUpdateSalesOrderConsumerHandlerInterface(kafkaClient, salesOrderLogRepository, salesOrderOpenSearchUseCase, database, ctx, args)
	return handler
}

func InitDeleteSalesOrderConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UpdateSalesOrderConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository, categoryRepository)
	handler := InitDeleteSalesOrderConsumerHandlerInterface(kafkaClient, salesOrderLogRepository, salesOrderOpenSearchUseCase, database, ctx, args)
	return handler
}

func InitDeleteSalesOrderDetailConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UpdateSalesOrderConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository, categoryRepository)
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
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	deliveryOrderDetailOpenSearchRepository := openSearchRepo.InitDeliveryOrderDetailOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository, categoryRepository)
	deliveryOrderOpenSearchUseCase := usecases.InitDeliveryOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, deliveryOrderOpenSearchRepository, deliveryOrderDetailOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderOpenSearchUseCase, database, ctx)
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
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	deliveryOrderDetailOpenSearchRepository := openSearchRepo.InitDeliveryOrderDetailOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository, categoryRepository)
	deliveryOrderOpenSearchUseCase := usecases.InitDeliveryOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, deliveryOrderOpenSearchRepository, deliveryOrderDetailOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderOpenSearchUseCase, database, ctx)
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
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	deliveryOrderDetailOpenSearchRepository := openSearchRepo.InitDeliveryOrderDetailOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository, categoryRepository)
	deliveryOrderOpenSearchUseCase := usecases.InitDeliveryOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, deliveryOrderOpenSearchRepository, deliveryOrderDetailOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderOpenSearchUseCase, database, ctx)
	handler := InitDeleteDeliveryOrderConsumerHandlerInterface(kafkaClient, deliveryOrderLogRepository, deliveryOrderOpenSearchUseCase, database, ctx, args)
	return handler
}

func InitUpdateSalesOrderDetailConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UpdateSalesOrderDetailConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository, categoryRepository)
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
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	deliveryOrderDetailOpenSearchRepository := openSearchRepo.InitDeliveryOrderDetailOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository, categoryRepository)
	deliveryOrderOpenSearchUseCase := usecases.InitDeliveryOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, deliveryOrderOpenSearchRepository, deliveryOrderDetailOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderOpenSearchUseCase, database, ctx)
	handler := InitUpdateDeliveryOrderDetailConsumerHandlerInterface(kafkaClient, deliveryOrderLogRepository, deliveryOrderOpenSearchUseCase, database, ctx, args)
	return handler
}

func InitUploadSOFileConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UploadSOFileConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	orderSourceRepository := repositories.InitOrderSourceRepository(database, redisdb)
	requestValidationRepository := repositories.InitRequestValidationRepository(database)
	uploadRepository := repositories.InitUploadRepository(requestValidationRepository, database)
	requestValidationMiddleware := middlewares.InitRequestValidationMiddlewareInterface(requestValidationRepository, orderSourceRepository)
	uploadSoHistoriesRepository := mongoRepo.InitSoUploadHistoriesRepositoryInterface(mongodbClient)
	uploadSoErrorLogsRepository := mongoRepo.InitSoUploadErrorLogsRepositoryInterface(mongodbClient)
	handler := InitUploadSOFileConsumerHandlerInterface(kafkaClient, uploadRepository, requestValidationMiddleware, requestValidationRepository, uploadSoHistoriesRepository, uploadSoErrorLogsRepository, salesOrderRepository, database, ctx, args)
	return handler
}

func InitUploadSOItemConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UploadSOItemConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	orderSourceRepository := repositories.InitOrderSourceRepository(database, redisdb)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderJourneysRepository := mongoRepo.InitSalesOrderJourneysRepository(mongodbClient)
	salesOrderDetailJourneysRepository := mongoRepo.InitSalesOrderDetailJourneysRepository(mongodbClient)
	uploadSOHistoriesRepository := mongoRepo.InitSoUploadHistoriesRepositoryInterface(mongodbClient)
	soUploadErrorLogsRepository := mongoRepo.InitSoUploadErrorLogsRepositoryInterface(mongodbClient)
	handler := InitUploadSOItemConsumerHandlerInterface(orderSourceRepository, orderStatusRepository, productRepository, uomRepository, agentRepository, storeRepository, userRepository, salesmanRepository, brandRepository, salesOrderRepository, salesOrderDetailRepository, salesOrderLogRepository, salesOrderJourneysRepository, salesOrderDetailJourneysRepository, uploadSOHistoriesRepository, soUploadErrorLogsRepository, kafkaClient, database, ctx, args)
	return handler
}

func InitUploadSOSJFileConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UploadSOSJFileConsumerHandlerInterface {
	orderSourceRepository := repositories.InitOrderSourceRepository(database, redisdb)
	requestValidationRepository := repositories.InitRequestValidationRepository(database)
	uploadRepository := repositories.InitUploadRepository(requestValidationRepository, database)
	requestValidationMiddleware := middlewares.InitRequestValidationMiddlewareInterface(requestValidationRepository, orderSourceRepository)
	sosjUploadHistoriesRepository := mongoRepo.InitSOSJUploadHistoriesRepositoryInterface(mongodbClient)
	sosjUploadErrorLogsRepository := mongoRepo.InitSOSJUploadErrorLogsRepositoryInterface(mongodbClient)
	handler := InitUploadSOSJFileConsumerHandlerInterface(kafkaClient, uploadRepository, requestValidationMiddleware, requestValidationRepository, sosjUploadHistoriesRepository, sosjUploadErrorLogsRepository, database, ctx, args)
	return handler
}

func InitUploadSOSJItemConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UploadSOItemConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	orderSourceRepository := repositories.InitOrderSourceRepository(database, redisdb)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderJourneysRepository := mongoRepo.InitSalesOrderJourneysRepository(mongodbClient)
	salesOrderDetailJourneysRepository := mongoRepo.InitSalesOrderDetailJourneysRepository(mongodbClient)
	warehouseRepository := repositories.InitWarehouseRepository(database, redisdb)
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	deliveryOrderDetailRepository := repositories.InitDeliveryOrderDetailRepository(database, redisdb)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	requestValidationRepository := repositories.InitRequestValidationRepository(database)
	sosjUploadHistoriesRepository := mongoRepo.InitSOSJUploadHistoriesRepositoryInterface(mongodbClient)
	sosjUploadErrorLogsRepository := mongoRepo.InitSOSJUploadErrorLogsRepositoryInterface(mongodbClient)
	uploadRepository := repositories.InitUploadRepository(requestValidationRepository, database)
	createSalesOrderConsumer := InitCreateSalesOrderConsumer(kafkaClient, mongodbClient, opensearchClient, database, redisdb, ctx, args)
	handler := InitUploadSOSJItemConsumerHandlerInterface(orderSourceRepository, orderStatusRepository, productRepository, uomRepository, agentRepository, storeRepository, userRepository, salesmanRepository, brandRepository, salesOrderRepository, salesOrderDetailRepository, salesOrderLogRepository, salesOrderJourneysRepository, salesOrderDetailJourneysRepository, warehouseRepository, deliveryOrderRepository, deliveryOrderDetailRepository, deliveryOrderLogRepository, sosjUploadHistoriesRepository, sosjUploadErrorLogsRepository, uploadRepository, createSalesOrderConsumer, kafkaClient, database, ctx, args)
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
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	deliveryOrderDetailOpenSearchRepository := openSearchRepo.InitDeliveryOrderDetailOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository, categoryRepository)
	deliveryOrderOpenSearchUseCase := usecases.InitDeliveryOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, deliveryOrderOpenSearchRepository, deliveryOrderDetailOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderOpenSearchUseCase, database, ctx)
	handler := InitDeleteDeliveryOrderDetailConsumerHandlerInterface(kafkaClient, deliveryOrderLogRepository, deliveryOrderOpenSearchUseCase, database, ctx, args)
	return handler
}

func InitUploadDOFileConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UploadSOSJFileConsumerHandlerInterface {
	orderSourceRepository := repositories.InitOrderSourceRepository(database, redisdb)
	requestValidationRepository := repositories.InitRequestValidationRepository(database)
	uploadRepository := repositories.InitUploadRepository(requestValidationRepository, database)
	requestValidationMiddleware := middlewares.InitRequestValidationMiddlewareInterface(requestValidationRepository, orderSourceRepository)
	uploadSJHistoriesRepository := mongoRepo.InitDoUploadHistoriesRepositoryInterface(mongodbClient)
	uploadSJErrorLogsRepository := mongoRepo.InitDoUploadErrorLogsRepositoryInterface(mongodbClient)
	warehouseRepository := repositories.InitWarehouseRepository(database, redisdb)
	handler := InitUploadDOFileConsumerHandlerInterface(kafkaClient, uploadRepository, requestValidationMiddleware, requestValidationRepository, uploadSJHistoriesRepository, uploadSJErrorLogsRepository, warehouseRepository, ctx, args, database)
	return handler
}

func InitUploadDOItemConsumer(kafkaClient kafkadbo.KafkaClientInterface, mongodbClient mongodb.MongoDBInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, database dbresolver.DB, redisdb redisdb.RedisInterface, ctx context.Context, args []interface{}) UploadSOItemConsumerHandlerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	deliveryOrderDetailRepository := repositories.InitDeliveryOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	orderSourceRepository := repositories.InitOrderSourceRepository(database, redisdb)
	warehouseRepository := repositories.InitWarehouseRepository(database, redisdb)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderJourneysRepository := mongoRepo.InitSalesOrderJourneysRepository(mongodbClient)
	uploadSJHistoriesRepository := mongoRepo.InitDoUploadHistoriesRepositoryInterface(mongodbClient)
	uploadSJErrorLogsRepository := mongoRepo.InitDoUploadErrorLogsRepositoryInterface(mongodbClient)
	handler := InitUploadDOItemConsumerHandlerInterface(deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, warehouseRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, userRepository, salesmanRepository, deliveryOrderLogRepository, salesOrderLogRepository, salesOrderJourneysRepository, uploadSJHistoriesRepository, uploadSJErrorLogsRepository, kafkaClient, ctx, args, database)
	return handler
}
