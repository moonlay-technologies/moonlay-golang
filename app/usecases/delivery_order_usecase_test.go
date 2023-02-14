package usecases

import (
	"context"
	"fmt"
	"order-service/app/models"
	repos "order-service/app/repositories"
	mongoRepo "order-service/app/repositories/mongod"
	repositories "order-service/app/repositories/open_search"
	"order-service/app/usecases/mocks"
	"order-service/global/utils/helper"
	kafkadbo "order-service/global/utils/kafka"
	"order-service/global/utils/mongodb"
	"order-service/global/utils/opensearch_dbo"
	"order-service/global/utils/redisdb"
	"order-service/global/utils/sqldb"
	"os"
	"strings"
	"testing"

	envConfig "github.com/joho/godotenv"

	"github.com/bxcodec/dbresolver"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
)

func newDeliveryOrderUsecase(status bool) deliveryOrderUseCase {
	ctx := context.Background()
	// mockDeliveryOrderRepository := &mocks.DeliveryOrderRepositoryInterface{}
	// mockDeliveryOrderDetailRepository := &mocks.DeliveryOrderDetailRepositoryInterface{}
	mockSalesOrderRepository := &mocks.SalesOrderRepositoryInterface{}
	mockSalesOrderDetailRepository := &mocks.SalesOrderDetailRepositoryInterface{}
	// mockOrderStatusRepository := &mocks.OrderStatusRepositoryInterface{}
	// mockOrderSourceRepository := &mocks.OrderSourceRepositoryInterface{}
	// mockWarehouseRepository := &mocks.WarehouseRepositoryInterface{}
	mockBrandRepository := &mocks.BrandRepositoryInterface{}
	mockUomRepository := &mocks.UomRepositoryInterface{}
	mockProductRepository := &mocks.ProductRepositoryInterface{}
	mockAgentRepository := &mocks.AgentRepositoryInterface{}
	mockStoreRepository := &mocks.StoreRepositoryInterface{}
	// mockDeliveryOrderLogRepository := &mocks.DeliveryOrderLogRepositoryInterface{}
	mockSalesOrderUseCase := &mocks.SalesOrderUseCaseInterface{}
	// mockKafkaClient := &mocks.KafkaClientInterface{}
	mockSalesOrderOpenSearchRepository := &mocks.SalesOrderOpenSearchRepositoryInterface{}
	openSearchHosts := []string{os.Getenv("OPENSEARCH_HOST_01")}
	openSearchClient := opensearch_dbo.InitOpenSearchClientInterface(openSearchHosts, os.Getenv("OPENSEARCH_USERNAME"), os.Getenv("OPENSEARCH_PASSWORD"), ctx)
	deliveryOrderOpenSearch := repositories.InitDeliveryOrderOpenSearchRepository(openSearchClient)

	if err := envConfig.Load(".env"); err != nil {
		errStr := fmt.Sprintf(".env not load properly %s", err.Error())
		helper.SetSentryError(err, errStr, sentry.LevelError)
		panic(err)
	}
	//mysql write
	mysqlWrite, err := sqldb.InitSql("mysql", os.Getenv("MYSQL_WRITE_HOST"), os.Getenv("MYSQL_WRITE_PORT"), os.Getenv("MYSQL_WRITE_USERNAME"), os.Getenv("MYSQL_WRITE_PASSWORD"), os.Getenv("MYSQL_WRITE_DATABASE"))
	if err != nil {
		errStr := fmt.Sprintf("Error mysql write connection %s", err.Error())
		helper.SetSentryError(err, errStr, sentry.LevelError)
		panic(err)
	}

	//mysql read
	mysqlRead, err := sqldb.InitSql("mysql", os.Getenv("MYSQL_READ_01_HOST"), os.Getenv("MYSQL_READ_01_PORT"), os.Getenv("MYSQL_READ_01_USERNAME"), os.Getenv("MYSQL_READ_01_PASSWORD"), os.Getenv("MYSQL_READ_01_DATABASE"))
	if err != nil {
		errStr := fmt.Sprintf("Error mysql read onnection %s", err.Error())
		helper.SetSentryError(err, errStr, sentry.LevelError)
		panic(err)
	}

	// MongoDB
	mongoDb := mongodb.InitMongoDB(os.Getenv("MONGO_HOST"), os.Getenv("MONGO_DATABASE"), os.Getenv("MONGO_USER"), os.Getenv("MONGO_PASSWORD"), os.Getenv("MONGO_PORT"), ctx)

	// Kafka
	kafkaHosts := strings.Split(os.Getenv("KAFKA_HOSTS"), ",")

	dbConnection := dbresolver.WrapDBs(mysqlWrite.DB(), mysqlRead.DB())
	redisDb := redisdb.InitRedis(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_DATABASE"))
	deliveryOrderRepository := repos.InitDeliveryRepository(dbConnection, redisDb)
	deliveryOrderDetailRepository := repos.InitDeliveryOrderDetailRepository(dbConnection, redisDb)
	warehouseRepository := repos.InitWarehouseRepository(dbConnection, redisDb)
	orderSourceRepository := repos.InitOrderSourceRepository(dbConnection, redisDb)
	orderStatusRepository := repos.InitOrderStatusRepository(dbConnection, redisDb)
	deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongoDb)
	kafkaClient := kafkadbo.InitKafkaClientInterface(ctx, kafkaHosts)
	// mocking
	return deliveryOrderUseCase{
		deliveryOrderRepository:           deliveryOrderRepository,
		deliveryOrderDetailRepository:     deliveryOrderDetailRepository,
		salesOrderRepository:              mockSalesOrderRepository,
		salesOrderDetailRepository:        mockSalesOrderDetailRepository,
		orderStatusRepository:             orderStatusRepository,
		orderSourceRepository:             orderSourceRepository,
		warehouseRepository:               warehouseRepository,
		brandRepository:                   mockBrandRepository,
		uomRepository:                     mockUomRepository,
		productRepository:                 mockProductRepository,
		agentRepository:                   mockAgentRepository,
		storeRepository:                   mockStoreRepository,
		deliveryOrderLogRepository:        deliveryOrderLogRepository,
		salesOrderUseCase:                 mockSalesOrderUseCase,
		kafkaClient:                       kafkaClient,
		deliveryOrderOpenSearchRepository: deliveryOrderOpenSearch,
		salesOrderOpenSearchRepository:    mockSalesOrderOpenSearchRepository,
		db:                                dbConnection,
		ctx:                               ctx,
	}
}

func Test_DeliveryOrderUseCase_InitDeliveryOrderUseCaseInterface_ShouldSuccess(t *testing.T) {
	// Arrange
	deliveryOrderUseCase := newDeliveryOrderUsecase(false)
	// deliveryOrderUseCaseInterface := InitDeliveryOrderUseCaseInterface(deliveryOrderUseCase.deliveryOrderRepository, deliveryOrderUseCase.deliveryOrderDetailRepository, deliveryOrderUseCase.salesOrderRepository, deliveryOrderUseCase.salesOrderDetailRepository, deliveryOrderUseCase.orderStatusRepository, deliveryOrderUseCase.orderSourceRepository, deliveryOrderUseCase.warehouseRepository, deliveryOrderUseCase.brandRepository, deliveryOrderUseCase.uomRepository, deliveryOrderUseCase.agentRepository, deliveryOrderUseCase.storeRepository, deliveryOrderUseCase.productRepository, deliveryOrderUseCase.userRepository, deliveryOrderUseCase.salesmanRepository, deliveryOrderUseCase.deliveryOrderLogRepository, deliveryOrderUseCase.deliveryOrderOpenSearchRepository, deliveryOrderUseCase.salesOrderOpenSearchRepository, deliveryOrderUseCase.salesOrderUseCase, deliveryOrderUseCase.kafkaClient, deliveryOrderUseCase.db, deliveryOrderUseCase.ctx)
	// Act
	dataDeliveryOrderUseCaseInit := InitDeliveryOrderUseCaseInterface(deliveryOrderUseCase.deliveryOrderRepository, deliveryOrderUseCase.deliveryOrderDetailRepository, deliveryOrderUseCase.salesOrderRepository, deliveryOrderUseCase.salesOrderDetailRepository, deliveryOrderUseCase.orderStatusRepository, deliveryOrderUseCase.orderSourceRepository, deliveryOrderUseCase.warehouseRepository, deliveryOrderUseCase.brandRepository, deliveryOrderUseCase.uomRepository, deliveryOrderUseCase.agentRepository, deliveryOrderUseCase.storeRepository, deliveryOrderUseCase.productRepository, deliveryOrderUseCase.userRepository, deliveryOrderUseCase.salesmanRepository, deliveryOrderUseCase.deliveryOrderLogRepository, deliveryOrderUseCase.deliveryOrderOpenSearchRepository, deliveryOrderUseCase.salesOrderOpenSearchRepository, deliveryOrderUseCase.salesOrderUseCase, deliveryOrderUseCase.kafkaClient, deliveryOrderUseCase.db, deliveryOrderUseCase.ctx)

	// Assert
	assert.NotNil(t, dataDeliveryOrderUseCaseInit)
}

func Test_DeliveryOrderUseCase_UpdateByID_ShouldSuccess(t *testing.T) {
	// Arrange
	deliveryOrderUsecase := newDeliveryOrderUsecase(false)
	ctx := context.Background()
	db := deliveryOrderUsecase.db
	sqlTx, _ := db.Begin()
	request := &models.DeliveryOrderUpdateByIDRequest{
		WarehouseID:   10,
		OrderSourceID: 2,
		OrderStatusID: 17,
		DeliveryOrderDetails: []*models.DeliveryOrderDetailUpdateByIDRequest{
			{
				Qty:  8,
				Note: "Kirim Segera",
			},
		},
	}

	// Act
	_, err := deliveryOrderUsecase.UpdateByID(90, request, sqlTx, ctx)
	// Assert
	assert.Nil(t, err)
}

func Test_DeliveryOrderUseCase_UpdateDODetailByID_ShouldSuccess(t *testing.T) {
	// Arrange
	deliveryOrderUsecase := newDeliveryOrderUsecase(false)
	ctx := context.Background()
	db := deliveryOrderUsecase.db
	sqlTx, _ := db.Begin()
	request := &models.DeliveryOrderDetailUpdateByIDRequest{
		Qty:  8,
		Note: "Kirim Segera",
	}

	// Act
	_, err := deliveryOrderUsecase.UpdateDODetailByID(90, request, sqlTx, ctx)
	// Assert
	assert.Nil(t, err)
}

func Test_DeliveryOrderUseCase_UpdateDoDetailByDeliveryOrderID_ShouldSuccess(t *testing.T) {
	// Arrange
	deliveryOrderUsecase := newDeliveryOrderUsecase(false)
	ctx := context.Background()
	db := deliveryOrderUsecase.db
	sqlTx, _ := db.Begin()
	request := []*models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest{
		{
			ID:   90,
			Qty:  8,
			Note: "Kirim Segera",
		},
	}

	// Act
	_, err := deliveryOrderUsecase.UpdateDoDetailByDeliveryOrderID(90, request, sqlTx, ctx)
	// Assert
	assert.Nil(t, err)
}

func Test_DeliveryOrderUseCase_Get_ShouldError(t *testing.T) {
	// Arrange
	deliveryOrderUsecase := newDeliveryOrderUsecase(false)
	// var ctx context.Context
	request := &models.DeliveryOrderRequest{}
	request.ID = 1

	// Act
	_, err := deliveryOrderUsecase.Get(request)
	// Assert
	assert.NotNil(t, err)
}

func Test_DeliveryOrderUseCase_GetByID_ShouldError(t *testing.T) {
	// Arrange
	deliveryOrderUsecase := newDeliveryOrderUsecase(false)
	var ctx context.Context
	request := &models.DeliveryOrderRequest{}
	request.ID = 1

	// Act
	_, err := deliveryOrderUsecase.GetByID(request, ctx)
	// Assert
	assert.NotNil(t, err)
}

func Test_DeliveryOrderUseCase_GetAgentID_ShouldError(t *testing.T) {
	// Arrange
	deliveryOrderUsecase := newDeliveryOrderUsecase(false)
	request := &models.DeliveryOrderRequest{}
	request.ID = 1

	// Act
	_, err := deliveryOrderUsecase.GetByAgentID(request)
	// Assert
	assert.NotNil(t, err)
}

func Test_DeliveryOrderUseCase_GetByStoreID_ShouldError(t *testing.T) {
	// Arrange
	deliveryOrderUsecase := newDeliveryOrderUsecase(false)
	request := &models.DeliveryOrderRequest{}
	request.ID = 1

	// Act
	_, err := deliveryOrderUsecase.GetByStoreID(request)
	// Assert
	assert.NotNil(t, err)
}

func Test_DeliveryOrderUseCase_GetBySalesmanID_ShouldError(t *testing.T) {
	// Arrange
	deliveryOrderUsecase := newDeliveryOrderUsecase(false)
	request := &models.DeliveryOrderRequest{}
	request.ID = 1

	// Act
	_, err := deliveryOrderUsecase.GetBySalesmanID(request)
	// Assert
	assert.NotNil(t, err)
}

func Test_DeliveryOrderUseCase_GetByOrderStatusID_ShouldError(t *testing.T) {
	// Arrange
	deliveryOrderUsecase := newDeliveryOrderUsecase(false)
	request := &models.DeliveryOrderRequest{}
	request.ID = 1

	// Act
	_, err := deliveryOrderUsecase.GetByOrderStatusID(request)
	// Assert
	assert.NotNil(t, err)
}

func Test_DeliveryOrderUseCase_GetByOrderSourceID_ShouldError(t *testing.T) {
	// Arrange
	deliveryOrderUsecase := newDeliveryOrderUsecase(false)
	request := &models.DeliveryOrderRequest{}
	request.ID = 1

	// Act
	_, err := deliveryOrderUsecase.GetByOrderSourceID(request)
	// Assert
	assert.NotNil(t, err)
}
