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
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	cartUseCase := usecases.InitCartUseCaseInterface(cartRepository, cartDetailRepository, orderStatusRepository, database, ctx)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, salesOrderJourneysRepository, salesOrderDetailJourneysRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	requestValidationRepository := repositories.InitUniqueRequestValidationRepository(database)
	requestValidationMiddleware := middlewares.InitRequestValidationMiddlewareInterface(requestValidationRepository, orderSourceRepository)
	handler := InitSalesOrderController(cartUseCase, salesOrderUseCase, requestValidationMiddleware, database, ctx)
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
	salesOrderLogRepository := mongoRepo.InitSalesOrderLogRepository(mongodbClient)
	salesOrderJourneysRepository := mongoRepo.InitSalesOrderJourneysRepository(mongodbClient)
	salesOrderDetailJourneysRepository := mongoRepo.InitSalesOrderDetailJourneysRepository(mongodbClient)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	ValidationRepository := repositories.InitUniqueRequestValidationRepository(database)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, salesOrderJourneysRepository, salesOrderDetailJourneysRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	deliveryOrderUseCase := usecases.InitDeliveryOrderUseCaseInterface(deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, warehouseRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, userRepository, salesmanRepository, deliveryOrderLogRepository, deliveryOrderOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderUseCase, salesOrderOpenSearchUseCase, kafkaClient, ValidationRepository, database, ctx)
	requestValidationRepository := repositories.InitUniqueRequestValidationRepository(database)
	requestValidationMiddleware := middlewares.InitRequestValidationMiddlewareInterface(requestValidationRepository, orderSourceRepository)
	handler := InitDeliveryOrderController(deliveryOrderUseCase, requestValidationMiddleware, database, ctx)
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
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	warehouseRepository := repositories.InitWarehouseRepository(database, redisdb)
	deliveryOrderDetailRepository := repositories.InitDeliveryOrderDetailRepository(database, redisdb)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	ValidationRepository := repositories.InitUniqueRequestValidationRepository(database)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, salesOrderJourneysRepository, salesOrderDetailJourneysRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	deliveryOrderUseCase := usecases.InitDeliveryOrderUseCaseInterface(deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, warehouseRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, userRepository, salesmanRepository, deliveryOrderLogRepository, deliveryOrderOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderUseCase, salesOrderOpenSearchUseCase, kafkaClient, ValidationRepository, database, ctx)

	handler := InitAgentController(salesOrderUseCase, deliveryOrderUseCase, database, ctx)
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
	userRepository := repositories.InitUserRepository(database, redisdb)
	warehouseRepository := repositories.InitWarehouseRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	deliveryOrderDetailRepository := repositories.InitDeliveryOrderDetailRepository(database, redisdb)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	ValidationRepository := repositories.InitUniqueRequestValidationRepository(database)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, salesOrderJourneysRepository, salesOrderDetailJourneysRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	deliveryOrderUseCase := usecases.InitDeliveryOrderUseCaseInterface(deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, warehouseRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, userRepository, salesmanRepository, deliveryOrderLogRepository, deliveryOrderOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderUseCase, salesOrderOpenSearchUseCase, kafkaClient, ValidationRepository, database, ctx)

	handler := InitStoreController(salesOrderUseCase, deliveryOrderUseCase, database, ctx)
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
	userRepository := repositories.InitUserRepository(database, redisdb)
	warehouseRepository := repositories.InitWarehouseRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	deliveryOrderDetailRepository := repositories.InitDeliveryOrderDetailRepository(database, redisdb)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	ValidationRepository := repositories.InitUniqueRequestValidationRepository(database)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, salesOrderJourneysRepository, salesOrderDetailJourneysRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	deliveryOrderUseCase := usecases.InitDeliveryOrderUseCaseInterface(deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, warehouseRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, userRepository, salesmanRepository, deliveryOrderLogRepository, deliveryOrderOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderUseCase, salesOrderOpenSearchUseCase, kafkaClient, ValidationRepository, database, ctx)
	handler := InitSalesmanController(salesOrderUseCase, deliveryOrderUseCase, database, ctx)
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
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	deliveryOrderDetailRepository := repositories.InitDeliveryOrderDetailRepository(database, redisdb)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	categoryRepository := repositories.InitCategoryRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderDetailOpenSearchRepository := openSearchRepo.InitSalesOrderDetailOpenSearchRepository(opensearchClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	cartUseCase := usecases.InitCartUseCaseInterface(cartRepository, cartDetailRepository, orderStatusRepository, database, ctx)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, salesOrderJourneysRepository, salesOrderDetailJourneysRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	salesOrderOpenSearchUseCase := usecases.InitSalesOrderOpenSearchUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, userRepository, salesmanRepository, categoryRepository, salesOrderOpenSearchRepository, salesOrderDetailOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	requestValidationRepository := repositories.InitUniqueRequestValidationRepository(database)
	requestValidationMiddleware := middlewares.InitRequestValidationMiddlewareInterface(requestValidationRepository, orderSourceRepository)
	deliveryOrderUseCase := usecases.InitDeliveryOrderUseCaseInterface(deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, warehouseRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, userRepository, salesmanRepository, deliveryOrderLogRepository, deliveryOrderOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderUseCase, salesOrderOpenSearchUseCase, kafkaClient, requestValidationRepository, database, ctx)
	handler := InitHostToHostController(cartUseCase, salesOrderUseCase, deliveryOrderUseCase, requestValidationMiddleware, database, ctx)
	return handler
}
