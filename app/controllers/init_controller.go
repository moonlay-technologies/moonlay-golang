package controllers

import (
	"context"
	"github.com/bxcodec/dbresolver"
	"poc-order-service/app/repositories"
	mongoRepo "poc-order-service/app/repositories/mongod"
	openSearchRepo "poc-order-service/app/repositories/open_search"
	"poc-order-service/app/usecases"
	kafkadbo "poc-order-service/global/utils/kafka"
	"poc-order-service/global/utils/mongodb"
	"poc-order-service/global/utils/opensearch_dbo"
	"poc-order-service/global/utils/redisdb"
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
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	cartUseCase := usecases.InitCartUseCaseInterface(cartRepository, cartDetailRepository, orderStatusRepository, database, ctx)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, userRepository, salesmanRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	handler := InitSalesOrderController(cartUseCase, salesOrderUseCase, database, ctx)
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
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, userRepository, salesmanRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	deliveryOrderUseCase := usecases.InitDeliveryOrderUseCaseInterface(deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, warehouseRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, deliveryOrderLogRepository, deliveryOrderOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderUseCase, kafkaClient, database, ctx)
	handler := InitDeliveryOrderController(deliveryOrderUseCase, database, ctx)
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
	deliveryOrderRepository := repositories.InitDeliveryRepository(database, redisdb)
	userRepository := repositories.InitUserRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	warehouseRepository := repositories.InitWarehouseRepository(database, redisdb)
	deliveryOrderDetailRepository := repositories.InitDeliveryOrderDetailRepository(database, redisdb)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, userRepository, salesmanRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	deliveryOrderUseCase := usecases.InitDeliveryOrderUseCaseInterface(deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, warehouseRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, deliveryOrderLogRepository, deliveryOrderOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderUseCase, kafkaClient, database, ctx)

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
	userRepository := repositories.InitUserRepository(database, redisdb)
	warehouseRepository := repositories.InitWarehouseRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	deliveryOrderDetailRepository := repositories.InitDeliveryOrderDetailRepository(database, redisdb)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, userRepository, salesmanRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	deliveryOrderUseCase := usecases.InitDeliveryOrderUseCaseInterface(deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, warehouseRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, deliveryOrderLogRepository, deliveryOrderOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderUseCase, kafkaClient, database, ctx)

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
	userRepository := repositories.InitUserRepository(database, redisdb)
	warehouseRepository := repositories.InitWarehouseRepository(database, redisdb)
	salesmanRepository := repositories.InitSalesmanRepository(database, redisdb)
	deliveryOrderDetailRepository := repositories.InitDeliveryOrderDetailRepository(database, redisdb)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongodbClient)
	deliveryOrderOpenSearchRepository := openSearchRepo.InitDeliveryOrderOpenSearchRepository(opensearchClient)
	salesOrderOpenSearchRepository := openSearchRepo.InitSalesOrderOpenSearchRepository(opensearchClient)
	salesOrderUseCase := usecases.InitSalesOrderUseCaseInterface(salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, agentRepository, brandRepository, storeRepository, productRepository, uomRepository, deliveryOrderRepository, salesOrderLogRepository, userRepository, salesmanRepository, salesOrderOpenSearchRepository, deliveryOrderOpenSearchRepository, kafkaClient, database, ctx)
	deliveryOrderUseCase := usecases.InitDeliveryOrderUseCaseInterface(deliveryOrderRepository, deliveryOrderDetailRepository, salesOrderRepository, salesOrderDetailRepository, orderStatusRepository, orderSourceRepository, warehouseRepository, brandRepository, uomRepository, agentRepository, storeRepository, productRepository, deliveryOrderLogRepository, deliveryOrderOpenSearchRepository, salesOrderOpenSearchRepository, salesOrderUseCase, kafkaClient, database, ctx)
	handler := InitSalesmanController(salesOrderUseCase, deliveryOrderUseCase, database, ctx)
	return handler
}
