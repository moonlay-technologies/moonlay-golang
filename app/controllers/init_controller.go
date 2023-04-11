package controllers

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

func InitHTTPSalesOrderController(database dbresolver.DB, redisdb redisdb.RedisInterface, mongodbClient mongodb.MongoDBInterface, kafkaClient kafkadbo.KafkaClientInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, ctx context.Context) SalesOrderControllerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	cartRepository := repositories.InitCartRepository(database, redisdb)
	cartDetailRepository := repositories.InitCartDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	orderSourceRepository := repositories.InitOrderSourceRepository(database, redisdb)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderJourneysRepository := mongoRepo.InitSalesOrderJourneysRepository(mongodbClient)
	salesOrderDetailJourneysRepository := mongoRepo.InitSalesOrderDetailJourneysRepository(mongodbClient)
	soUploadHistoriesRepository := mongoRepo.InitSoUploadHistoriesRepositoryInterface(mongodbClient)
	soUploadErrorLogsRepository := mongoRepo.InitSoUploadErrorLogsRepositoryInterface(mongodbClient)
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	cartUseCase := usecases.InitCartUseCaseInterface(cartRepository, cartDetailRepository, orderStatusRepository, database, ctx)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, salesOrderJourneysRepository, salesOrderDetailJourneysRepository, soUploadHistoriesRepository, soUploadErrorLogsRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, kafkaClient, database, ctx)
	requestValidationRepository := repositories.InitRequestValidationRepository(database)
	requestValidationMiddleware := middlewares.InitRequestValidationMiddlewareInterface(requestValidationRepository, orderSourceRepository)
	salesOrderValidator := usecases.InitSalesOrderValidator(requestValidationMiddleware, orderSourceRepository, salesmanRepository, orderStatusRepository, salesOrderRepository, salesOrderDetailRepository, deliveryOrderRepository, database, ctx)
	handler := InitSalesOrderController(cartUseCase, salesOrderUseCase, salesOrderValidator, requestValidationMiddleware, database, ctx)
	return handler
}

func InitHTTPDeliveryOrderController(database dbresolver.DB, redisdb redisdb.RedisInterface, mongodbClient mongodb.MongoDBInterface, kafkaClient kafkadbo.KafkaClientInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, ctx context.Context) DeliveryOrderControllerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	orderSourceRepository := repositories.InitOrderSourceRepository(database, redisdb)
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	deliveryOrderDetailRepository := repositories.InitDeliveryOrderDetailRepository(database, redisdb)
	warehouseRepository := repositories.InitWarehouseRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	deliveryOrderDetailOpenSearchRepository := openSearchRepo.InitDeliveryOrderDetailOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderJourneysRepository := mongoRepo.InitDeliveryOrderJourneysRepository(mongodbClient)
	salesOrderJourneysRepository := mongoRepo.InitSalesOrderJourneysRepository(mongodbClient)
	doUploadHistoriesRepository := mongoRepo.InitDoUploadHistoriesRepositoryInterface(mongodbClient)
	doUploadErrorLogsRepository := mongoRepo.InitDoUploadErrorLogsRepositoryInterface(mongodbClient)
	uploadRepositories := repositories.InitUploadRepository(database)
	salesOrderJourneyRepository := mongoRepo.InitSalesOrderJourneysRepository(mongodbClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, categoryRepository, salesOrderJourneyRepository, uploadRepositories, ctx)
	deliveryOrderUseCase := usecases.InitDeliveryOrderUseCaseInterface(deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, warehouseRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, userRepository, salesmanRepository, salesOrderJourneysRepository, deliveryOrderLogRepository, deliveryOrderJourneysRepository, doUploadHistoriesRepository, doUploadErrorLogsRepository, deliveryOrderOpenSearchRepository, deliveryOrderDetailOpenSearchRepository, salesOrderOpenSearchUseCase, kafkaClient, database, ctx)
	requestValidationRepository := repositories.InitRequestValidationRepository(database)
	requestValidationMiddleware := middlewares.InitRequestValidationMiddlewareInterface(requestValidationRepository, orderSourceRepository)
	deliveryOrderValidator := usecases.InitDeliveryOrderValidator(requestValidationMiddleware, database, ctx)
	handler := InitDeliveryOrderController(deliveryOrderUseCase, deliveryOrderValidator, requestValidationMiddleware, database, ctx)
	return handler
}

func InitHTTPAgentController(database dbresolver.DB, redisdb redisdb.RedisInterface, mongodbClient mongodb.MongoDBInterface, kafkaClient kafkadbo.KafkaClientInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, ctx context.Context) AgentControllerInterface {
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
	salesOrderJourneysRepository := mongoRepo.InitSalesOrderJourneysRepository(mongodbClient)
	salesOrderDetailJourneysRepository := mongoRepo.InitSalesOrderDetailJourneysRepository(mongodbClient)
	soUploadHistoriesRepository := mongoRepo.InitSoUploadHistoriesRepositoryInterface(mongodbClient)
	soUploadErrorLogsRepository := mongoRepo.InitSoUploadErrorLogsRepositoryInterface(mongodbClient)
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderJourneysRepository := mongoRepo.InitDeliveryOrderJourneysRepository(mongodbClient)
	doUploadHistoriesRepository := mongoRepo.InitDoUploadHistoriesRepositoryInterface(mongodbClient)
	doUploadErrorLogsRepository := mongoRepo.InitDoUploadErrorLogsRepositoryInterface(mongodbClient)
	warehouseRepository := repositories.InitWarehouseRepository(database, redisdb)
	deliveryOrderDetailRepository := repositories.InitDeliveryOrderDetailRepository(database, redisdb)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	deliveryOrderDetailOpenSearchRepository := openSearchRepo.InitDeliveryOrderDetailOpenSearchRepository(opensearchClient)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, salesOrderJourneysRepository, salesOrderDetailJourneysRepository, soUploadHistoriesRepository, soUploadErrorLogsRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, kafkaClient, database, ctx)
	uploadRepositories := repositories.InitUploadRepository(database)
	salesOrderJourneyRepository := mongoRepo.InitSalesOrderJourneysRepository(mongodbClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, categoryRepository, salesOrderJourneyRepository, uploadRepositories, ctx)
	deliveryOrderUseCase := usecases.InitDeliveryOrderUseCaseInterface(deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, warehouseRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, userRepository, salesmanRepository, salesOrderJourneysRepository, deliveryOrderLogRepository, deliveryOrderJourneysRepository, doUploadHistoriesRepository, doUploadErrorLogsRepository, deliveryOrderOpenSearchRepository, deliveryOrderDetailOpenSearchRepository, salesOrderOpenSearchUseCase, kafkaClient, database, ctx)
	requestValidationRepository := repositories.InitRequestValidationRepository(database)
	requestValidationMiddleware := middlewares.InitRequestValidationMiddlewareInterface(requestValidationRepository, orderSourceRepository)
	salesOrderValidator := usecases.InitSalesOrderValidator(requestValidationMiddleware, orderSourceRepository, salesmanRepository, orderStatusRepository, salesOrderRepository, salesOrderDetailRepository, deliveryOrderRepository, database, ctx)
	deliveryOrderValidator := usecases.InitDeliveryOrderValidator(requestValidationMiddleware, database, ctx)
	handler := InitAgentController(salesOrderUseCase, deliveryOrderUseCase, salesOrderValidator, deliveryOrderValidator, requestValidationMiddleware, database, ctx)
	return handler
}

func InitHTTPStoreController(database dbresolver.DB, redisdb redisdb.RedisInterface, mongodbClient mongodb.MongoDBInterface, kafkaClient kafkadbo.KafkaClientInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, ctx context.Context) StoreControllerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	orderSourceRepository := repositories.InitOrderSourceRepository(database, redisdb)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderJourneysRepository := mongoRepo.InitSalesOrderJourneysRepository(mongodbClient)
	salesOrderDetailJourneysRepository := mongoRepo.InitSalesOrderDetailJourneysRepository(mongodbClient)
	soUploadHistoriesRepository := mongoRepo.InitSoUploadHistoriesRepositoryInterface(mongodbClient)
	soUploadErrorLogsRepository := mongoRepo.InitSoUploadErrorLogsRepositoryInterface(mongodbClient)
	userRepository := repositories.InitUserRepository(database, redisdb)
	warehouseRepository := repositories.InitWarehouseRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	deliveryOrderDetailRepository := repositories.InitDeliveryOrderDetailRepository(database, redisdb)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderJourneysRepository := mongoRepo.InitDeliveryOrderJourneysRepository(mongodbClient)
	doUploadHistoriesRepository := mongoRepo.InitDoUploadHistoriesRepositoryInterface(mongodbClient)
	doUploadErrorLogsRepository := mongoRepo.InitDoUploadErrorLogsRepositoryInterface(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	deliveryOrderDetailOpenSearchRepository := openSearchRepo.InitDeliveryOrderDetailOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, salesOrderJourneysRepository, salesOrderDetailJourneysRepository, soUploadHistoriesRepository, soUploadErrorLogsRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, kafkaClient, database, ctx)
	uploadRepositories := repositories.InitUploadRepository(database)
	salesOrderJourneyRepository := mongoRepo.InitSalesOrderJourneysRepository(mongodbClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, categoryRepository, salesOrderJourneyRepository, uploadRepositories, ctx)
	deliveryOrderUseCase := usecases.InitDeliveryOrderUseCaseInterface(deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, warehouseRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, userRepository, salesmanRepository, salesOrderJourneysRepository, deliveryOrderLogRepository, deliveryOrderJourneysRepository, doUploadHistoriesRepository, doUploadErrorLogsRepository, deliveryOrderOpenSearchRepository, deliveryOrderDetailOpenSearchRepository, salesOrderOpenSearchUseCase, kafkaClient, database, ctx)
	requestValidationRepository := repositories.InitRequestValidationRepository(database)
	requestValidationMiddleware := middlewares.InitRequestValidationMiddlewareInterface(requestValidationRepository, orderSourceRepository)
	salesOrderValidator := usecases.InitSalesOrderValidator(requestValidationMiddleware, orderSourceRepository, salesmanRepository, orderStatusRepository, salesOrderRepository, salesOrderDetailRepository, deliveryOrderRepository, database, ctx)
	deliveryOrderValidator := usecases.InitDeliveryOrderValidator(requestValidationMiddleware, database, ctx)
	handler := InitStoreController(salesOrderUseCase, deliveryOrderUseCase, salesOrderValidator, deliveryOrderValidator, requestValidationMiddleware, database, ctx)
	return handler
}

func InitHTTPSalesmanController(database dbresolver.DB, redisdb redisdb.RedisInterface, mongodbClient mongodb.MongoDBInterface, kafkaClient kafkadbo.KafkaClientInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, ctx context.Context) SalesmanControllerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	orderSourceRepository := repositories.InitOrderSourceRepository(database, redisdb)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderJourneysRepository := mongoRepo.InitSalesOrderJourneysRepository(mongodbClient)
	salesOrderDetailJourneysRepository := mongoRepo.InitSalesOrderDetailJourneysRepository(mongodbClient)
	soUploadHistoriesRepository := mongoRepo.InitSoUploadHistoriesRepositoryInterface(mongodbClient)
	soUploadErrorLogsRepository := mongoRepo.InitSoUploadErrorLogsRepositoryInterface(mongodbClient)
	userRepository := repositories.InitUserRepository(database, redisdb)
	warehouseRepository := repositories.InitWarehouseRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	deliveryOrderDetailRepository := repositories.InitDeliveryOrderDetailRepository(database, redisdb)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderJourneysRepository := mongoRepo.InitDeliveryOrderJourneysRepository(mongodbClient)
	doUploadHistoriesRepository := mongoRepo.InitDoUploadHistoriesRepositoryInterface(mongodbClient)
	doUploadErrorLogsRepository := mongoRepo.InitDoUploadErrorLogsRepositoryInterface(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	deliveryOrderDetailOpenSearchRepository := openSearchRepo.InitDeliveryOrderDetailOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, salesOrderJourneysRepository, salesOrderDetailJourneysRepository, soUploadHistoriesRepository, soUploadErrorLogsRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, kafkaClient, database, ctx)
	uploadRepositories := repositories.InitUploadRepository(database)
	salesOrderJourneyRepository := mongoRepo.InitSalesOrderJourneysRepository(mongodbClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, categoryRepository, salesOrderJourneyRepository, uploadRepositories, ctx)
	deliveryOrderUseCase := usecases.InitDeliveryOrderUseCaseInterface(deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, warehouseRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, userRepository, salesmanRepository, salesOrderJourneysRepository, deliveryOrderLogRepository, deliveryOrderJourneysRepository, doUploadHistoriesRepository, doUploadErrorLogsRepository, deliveryOrderOpenSearchRepository, deliveryOrderDetailOpenSearchRepository, salesOrderOpenSearchUseCase, kafkaClient, database, ctx)
	requestValidationRepository := repositories.InitRequestValidationRepository(database)
	requestValidationMiddleware := middlewares.InitRequestValidationMiddlewareInterface(requestValidationRepository, orderSourceRepository)
	salesOrderValidator := usecases.InitSalesOrderValidator(requestValidationMiddleware, orderSourceRepository, salesmanRepository, orderStatusRepository, salesOrderRepository, salesOrderDetailRepository, deliveryOrderRepository, database, ctx)
	deliveryOrderValidator := usecases.InitDeliveryOrderValidator(requestValidationMiddleware, database, ctx)
	handler := InitSalesmanController(salesOrderUseCase, deliveryOrderUseCase, salesOrderValidator, deliveryOrderValidator, requestValidationMiddleware, database, ctx)
	return handler
}

func InitHTTPHostToHostController(database dbresolver.DB, redisdb redisdb.RedisInterface, mongodbClient mongodb.MongoDBInterface, kafkaClient kafkadbo.KafkaClientInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, ctx context.Context) HostToHostControllerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	cartRepository := repositories.InitCartRepository(database, redisdb)
	cartDetailRepository := repositories.InitCartDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	orderSourceRepository := repositories.InitOrderSourceRepository(database, redisdb)
	warehouseRepository := repositories.InitWarehouseRepository(database, redisdb)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderJourneysRepository := mongoRepo.InitSalesOrderJourneysRepository(mongodbClient)
	salesOrderDetailJourneysRepository := mongoRepo.InitSalesOrderDetailJourneysRepository(mongodbClient)
	soUploadHistoriesRepository := mongoRepo.InitSoUploadHistoriesRepositoryInterface(mongodbClient)
	soUploadErrorLogsRepository := mongoRepo.InitSoUploadErrorLogsRepositoryInterface(mongodbClient)
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	deliveryOrderDetailRepository := repositories.InitDeliveryOrderDetailRepository(database, redisdb)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderJourneysRepository := mongoRepo.InitDeliveryOrderJourneysRepository(mongodbClient)
	doUploadHistoriesRepository := mongoRepo.InitDoUploadHistoriesRepositoryInterface(mongodbClient)
	doUploadErrorLogsRepository := mongoRepo.InitDoUploadErrorLogsRepositoryInterface(mongodbClient)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	deliveryOrderDetailOpenSearchRepository := openSearchRepo.InitDeliveryOrderDetailOpenSearchRepository(opensearchClient)
	cartUseCase := usecases.InitCartUseCaseInterface(cartRepository, cartDetailRepository, orderStatusRepository, database, ctx)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, salesOrderJourneysRepository, salesOrderDetailJourneysRepository, soUploadHistoriesRepository, soUploadErrorLogsRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, kafkaClient, database, ctx)
	uploadRepositories := repositories.InitUploadRepository(database)
	salesOrderJourneyRepository := mongoRepo.InitSalesOrderJourneysRepository(mongodbClient)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, productRepository, uomRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, categoryRepository, salesOrderJourneyRepository, uploadRepositories, ctx)
	requestValidationRepository := repositories.InitRequestValidationRepository(database)
	requestValidationMiddleware := middlewares.InitRequestValidationMiddlewareInterface(requestValidationRepository, orderSourceRepository)
	deliveryOrderUseCase := usecases.InitDeliveryOrderUseCaseInterface(deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, warehouseRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, userRepository, salesmanRepository, salesOrderJourneysRepository, deliveryOrderLogRepository, deliveryOrderJourneysRepository, doUploadHistoriesRepository, doUploadErrorLogsRepository, deliveryOrderOpenSearchRepository, deliveryOrderDetailOpenSearchRepository, salesOrderOpenSearchUseCase, kafkaClient, database, ctx)
	handler := InitHostToHostController(cartUseCase, salesOrderUseCase, deliveryOrderUseCase, requestValidationMiddleware, database, ctx)
	return handler
}

func InitHTTPUploadController(database dbresolver.DB, redisdb redisdb.RedisInterface, mongodbClient mongodb.MongoDBInterface, kafkaClient kafkadbo.KafkaClientInterface, opensearchClient opensearch_dbo.OpenSearchClientInterface, ctx context.Context) UploadControllerInterface {
	salesOrderRepository := repositories.InitSalesOrderRepository(database, redisdb)
	salesOrderDetailRepository := repositories.InitSalesOrderDetailRepository(database, redisdb)
	orderStatusRepository := repositories.InitOrderStatusRepository(database, redisdb)
	orderSourceRepository := repositories.InitOrderSourceRepository(database, redisdb)
	requestValidationRepository := repositories.InitRequestValidationRepository(database)
	uploadRepositories := repositories.InitUploadRepository(database)
	agentRepository := repositories.InitAgentRepository(database, redisdb)
	brandRepository := repositories.InitBrandRepository(database, redisdb)
	storeRepository := repositories.InitStoreRepository(database, redisdb)
	productRepository := repositories.InitProductRepository(database, redisdb)
	uomRepository := repositories.InitUomRepository(database, redisdb)
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	deliveryOrderDetailRepository := repositories.InitDeliveryOrderDetailRepository(database, redisdb)
	warehouseRepository := repositories.InitWarehouseRepository(database, redisdb)
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderJourneysRepository := mongoRepo.InitSalesOrderJourneysRepository(mongodbClient)
	salesOrderDetailJourneysRepository := mongoRepo.InitSalesOrderDetailJourneysRepository(mongodbClient)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	requestValidationMiddleware := middlewares.InitRequestValidationMiddlewareInterface(requestValidationRepository, orderSourceRepository)
	salesOrderValidator := usecases.InitSalesOrderValidator(requestValidationMiddleware, orderSourceRepository, salesmanRepository, orderStatusRepository, salesOrderRepository, salesOrderDetailRepository, deliveryOrderRepository, database, ctx)
	soUploadHistoriesRepository := mongoRepo.InitSoUploadHistoriesRepositoryInterface(mongodbClient)
	soUploadErrorLogsRepository := mongoRepo.InitSoUploadErrorLogsRepositoryInterface(mongodbClient)
	sosjUploadHistoriesRepository := mongoRepo.InitSOSJUploadHistoriesRepositoryInterface(mongodbClient)
	sosjUploadErrorLogsRepository := mongoRepo.InitSOSJUploadErrorLogsRepositoryInterface(mongodbClient)
	doUploadHistoriesRepository := mongoRepo.InitDoUploadHistoriesRepositoryInterface(mongodbClient)
	uploadUseCase := usecases.InitUploadUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderLogRepository, salesOrderJourneysRepository, salesOrderDetailJourneysRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, uploadRepositories, warehouseRepository, soUploadHistoriesRepository, soUploadErrorLogsRepository, sosjUploadHistoriesRepository, sosjUploadErrorLogsRepository, doUploadHistoriesRepository, kafkaClient, database, ctx)
	handler := InitUploadController(uploadUseCase, salesOrderValidator, requestValidationMiddleware, ctx)
	return handler
}
