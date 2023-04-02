package usecases

import (
	"context"
	"fmt"
	"order-service/app/models"
	repositories "order-service/app/repositories/open_search"
	"order-service/app/usecases/mocks"
	"order-service/global/utils/helper"
	"order-service/global/utils/opensearch_dbo"
	"order-service/global/utils/sqldb"
	"os"
	"testing"

	"github.com/bxcodec/dbresolver"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
)

func newSalesOrderUseCase(status bool) salesOrderUseCase {
	ctx := context.Background()
	mockDeliveryOrderRepository := &mocks.DeliveryOrderRepositoryInterface{}
	mockKafkaClient := &mocks.KafkaClientInterface{}
	mockOrderStatusRepository := &mocks.OrderStatusRepositoryInterface{}
	mockOrderSourceRepository := &mocks.OrderSourceRepositoryInterface{}
	mockUserRepository := &mocks.UserRepositoryInterface{}
	mockSalesmanRepository := &mocks.SalesmanRepositoryInterface{}
	mockSalesOrderLogRepository := &mocks.SalesOrderLogRepositoryInterface{}
	mockSalesOrderRepository := &mocks.SalesOrderRepositoryInterface{}
	mockSalesOrderDetailRepository := &mocks.SalesOrderDetailRepositoryInterface{}
	mockBrandRepository := &mocks.BrandRepositoryInterface{}
	mockUomRepository := &mocks.UomRepositoryInterface{}
	mockProductRepository := &mocks.ProductRepositoryInterface{}
	mockAgentRepository := &mocks.AgentRepositoryInterface{}
	mockStoreRepository := &mocks.StoreRepositoryInterface{}
	openSearchHosts := []string{os.Getenv("OPENSEARCH_HOST_01")}
	openSearchClient := opensearch_dbo.InitOpenSearchClientInterface(openSearchHosts, os.Getenv("OPENSEARCH_USERNAME"), os.Getenv("OPENSEARCH_PASSWORD"), ctx)
	salesOrderOpenSearch := repositories.InitSalesOrderOpenSearchRepository(openSearchClient)
	salesOrderDetailOpenSearch := repositories.InitSalesOrderDetailOpenSearchRepository(openSearchClient)

	// access to config
	// if err := envConfig.Load(".env"); err != nil {
	// 	errStr := fmt.Sprintf(".env not load properly %s", err.Error())
	// 	helper.SetSentryError(err, errStr, sentry.LevelError)
	// 	panic(err)
	// }

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
	// mongoDb := mongodb.InitMongoDB(os.Getenv("MONGO_HOST"), os.Getenv("MONGO_DATABASE"), os.Getenv("MONGO_USER"), os.Getenv("MONGO_PASSWORD"), os.Getenv("MONGO_PORT"), ctx)

	// Kafka
	// kafkaHosts := strings.Split(os.Getenv("KAFKA_HOSTS"), ",")

	dbConnection := dbresolver.WrapDBs(mysqlWrite.DB(), mysqlRead.DB())
	// redisDb := redisdb.InitRedis(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_DATABASE"))
	// deliveryOrderRepository := repos.InitDeliveryRepository(dbConnection, redisDb)
	// deliveryOrderDetailRepository := repos.InitDeliveryOrderDetailRepository(dbConnection, redisDb)
	// warehouseRepository := repos.InitWarehouseRepository(dbConnection, redisDb)
	// orderSourceRepository := repos.InitOrderSourceRepository(dbConnection, redisDb)
	// orderStatusRepository := repos.InitOrderStatusRepository(dbConnection, redisDb)
	// deliveryOrderLogRepository := mongoRepo.InitDeliveryOrderLogRepository(mongoDb)
	// kafkaClient := kafkadbo.InitKafkaClientInterface(ctx, kafkaHosts)
	return salesOrderUseCase{
		salesOrderRepository:                 mockSalesOrderRepository,
		salesOrderDetailRepository:           mockSalesOrderDetailRepository,
		orderStatusRepository:                mockOrderStatusRepository,
		orderSourceRepository:                mockOrderSourceRepository,
		agentRepository:                      mockAgentRepository,
		brandRepository:                      mockBrandRepository,
		storeRepository:                      mockStoreRepository,
		productRepository:                    mockProductRepository,
		uomRepository:                        mockUomRepository,
		deliveryOrderRepository:              mockDeliveryOrderRepository,
		salesOrderLogRepository:              mockSalesOrderLogRepository,
		salesOrderJourneysRepository:         nil,
		salesOrderDetailJourneysRepository:   nil,
		soUploadHistoriesRepository:          nil,
		soUploadErrorLogsRepository:          nil,
		userRepository:                       mockUserRepository,
		salesmanRepository:                   mockSalesmanRepository,
		categoryRepository:                   nil,
		salesOrderOpenSearchRepository:       salesOrderOpenSearch,
		salesOrderDetailOpenSearchRepository: salesOrderDetailOpenSearch,
		kafkaClient:                          mockKafkaClient,
		db:                                   dbConnection,
		ctx:                                  ctx,
	}
}

func Test_SalesOrderUseCase_InitSalesOrderUseCaseInterface_ShouldSuccess(t *testing.T) {
	// Arrange
	salesOrderUseCase := newSalesOrderUseCase(false)
	// salesOrderUseCaseInterface := InitSalesOrderUseCaseInterface(salesOrderUseCase.deliveryOrderRepository, salesOrderUseCase.deliveryOrderDetailRepository, salesOrderUseCase.salesOrderRepository, salesOrderUseCase.salesOrderDetailRepository, salesOrderUseCase.orderStatusRepository, salesOrderUseCase.orderSourceRepository, salesOrderUseCase.warehouseRepository, salesOrderUseCase.brandRepository, salesOrderUseCase.uomRepository, salesOrderUseCase.agentRepository, salesOrderUseCase.storeRepository, salesOrderUseCase.productRepository, salesOrderUseCase.userRepository, salesOrderUseCase.salesmanRepository, salesOrderUseCase.deliveryOrderLogRepository, salesOrderUseCase.deliveryOrderOpenSearchRepository, salesOrderUseCase.salesOrderOpenSearchRepository, salesOrderUseCase.salesOrderUseCase, salesOrderUseCase.kafkaClient, salesOrderUseCase.db, salesOrderUseCase.ctx)
	// Act
	dataSalesOrderUseCaseInit := InitSalesOrderUseCaseInterface(salesOrderUseCase.salesOrderRepository, salesOrderUseCase.salesOrderDetailRepository, salesOrderUseCase.orderStatusRepository, salesOrderUseCase.orderSourceRepository, salesOrderUseCase.agentRepository, salesOrderUseCase.brandRepository, salesOrderUseCase.storeRepository, salesOrderUseCase.productRepository, salesOrderUseCase.uomRepository, salesOrderUseCase.deliveryOrderRepository, salesOrderUseCase.salesOrderLogRepository, salesOrderUseCase.salesOrderJourneysRepository, salesOrderUseCase.salesOrderDetailJourneysRepository, salesOrderUseCase.soUploadHistoriesRepository, salesOrderUseCase.soUploadErrorLogsRepository, salesOrderUseCase.userRepository, salesOrderUseCase.salesmanRepository, salesOrderUseCase.categoryRepository, salesOrderUseCase.salesOrderOpenSearchRepository, salesOrderUseCase.salesOrderDetailOpenSearchRepository, salesOrderUseCase.kafkaClient, salesOrderUseCase.db, salesOrderUseCase.ctx)

	// Assert
	assert.NotNil(t, dataSalesOrderUseCaseInit)
}

// func Test_SalesOrderUseCase_UpdateByID_ShouldSuccess(t *testing.T) {
// 	// Arrange
// 	salesOrderUseCase := newSalesOrderUseCase(false)
// 	ctx := context.Background()
// 	db := salesOrderUseCase.db
// 	sqlTx, _ := db.Begin()
// 	request := &models.DeliveryOrderUpdateByIDRequest{
// 		WarehouseID:   10,
// 		OrderSourceID: 2,
// 		OrderStatusID: 17,
// 		DeliveryOrderDetails: []*models.DeliveryOrderDetailUpdateByIDRequest{
// 			{
// 				Qty:  8,
// 				Note: "Kirim Segera",
// 			},
// 		},
// 	}

// 	// Act
// 	_, err := salesOrderUseCase.UpdateByID(90, request, sqlTx, ctx)
// 	// Assert
// 	assert.Nil(t, err)
// }

// func Test_SalesOrderUseCase_UpdateDODetailByID_ShouldSuccess(t *testing.T) {
// 	// Arrange
// 	salesOrderUseCase := newSalesOrderUseCase(false)
// 	ctx := context.Background()
// 	db := salesOrderUseCase.db
// 	sqlTx, _ := db.Begin()
// 	request := &models.DeliveryOrderDetailUpdateByIDRequest{
// 		Qty:  8,
// 		Note: "Kirim Segera",
// 	}

// 	// Act
// 	_, err := salesOrderUseCase.UpdateDODetailByID(90, request, sqlTx, ctx)
// 	// Assert
// 	assert.Nil(t, err)
// }

// func Test_SalesOrderUseCase_UpdateDoDetailByDeliveryOrderID_ShouldSuccess(t *testing.T) {
// 	// Arrange
// 	salesOrderUseCase := newSalesOrderUseCase(false)
// 	ctx := context.Background()
// 	db := salesOrderUseCase.db
// 	sqlTx, _ := db.Begin()
// 	request := []*models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest{
// 		{
// 			ID:   90,
// 			Qty:  8,
// 			Note: "Kirim Segera",
// 		},
// 	}

// 	// Act
// 	_, err := salesOrderUseCase.UpdateDoDetailByDeliveryOrderID(90, request, sqlTx, ctx)
// 	// Assert
// 	assert.Nil(t, err)
// }

func Test_SalesOrderUseCase_Get_ShouldError(t *testing.T) {
	// Arrange
	salesOrderUseCase := newSalesOrderUseCase(false)
	request := &models.SalesOrderRequest{}
	request.ID = 1

	// Act
	_, err := salesOrderUseCase.Get(request)
	// Assert
	assert.NotNil(t, err)
}

// func Test_SalesOrderUseCase_GetByID_ShouldError(t *testing.T) {
// 	// Arrange
// 	salesOrderUseCase := newSalesOrderUseCase(false)
// 	var ctx context.Context
// 	request := &models.SalesOrderRequest{}
// 	request.ID = 1

// 	// Act
// 	_, err := salesOrderUseCase.GetByID(request, ctx)
// 	// Assert
// 	assert.NotNil(t, err)
// }

// func Test_SalesOrderUseCase_GetByIDWithDetail_ShouldError(t *testing.T) {
// 	// Arrange
// 	salesOrderUseCase := newSalesOrderUseCase(false)
// 	var ctx context.Context
// 	request := &models.SalesOrderRequest{}
// 	request.ID = 1

// 	// Act
// 	_, err := salesOrderUseCase.GetByIDWithDetail(request, ctx)
// 	// Assert
// 	assert.NotNil(t, err)
// }

func Test_SalesOrderUseCase_GetAgentID_ShouldError(t *testing.T) {
	// Arrange
	salesOrderUseCase := newSalesOrderUseCase(false)
	request := &models.SalesOrderRequest{}
	request.ID = 1

	// Act
	_, err := salesOrderUseCase.GetByAgentID(request)
	// Assert
	assert.NotNil(t, err)
}

func Test_SalesOrderUseCase_GetByStoreID_ShouldError(t *testing.T) {
	// Arrange
	salesOrderUseCase := newSalesOrderUseCase(false)
	request := &models.SalesOrderRequest{}
	request.ID = 1

	// Act
	_, err := salesOrderUseCase.GetByStoreID(request)
	// Assert
	assert.NotNil(t, err)
}

func Test_SalesOrderUseCase_GetBySalesmanID_ShouldError(t *testing.T) {
	// Arrange
	salesOrderUseCase := newSalesOrderUseCase(false)
	request := &models.SalesOrderRequest{}
	request.ID = 1

	// Act
	_, err := salesOrderUseCase.GetBySalesmanID(request)
	// Assert
	assert.NotNil(t, err)
}

func Test_SalesOrderUseCase_GetByOrderStatusID_ShouldError(t *testing.T) {
	// Arrange
	salesOrderUseCase := newSalesOrderUseCase(false)
	request := &models.SalesOrderRequest{}
	request.ID = 1

	// Act
	_, err := salesOrderUseCase.GetByOrderStatusID(request)
	// Assert
	assert.NotNil(t, err)
}

func Test_SalesOrderUseCase_GetByOrderSourceID_ShouldError(t *testing.T) {
	// Arrange
	salesOrderUseCase := newSalesOrderUseCase(false)
	request := &models.SalesOrderRequest{}
	request.ID = 1

	// Act
	_, err := salesOrderUseCase.GetByOrderSourceID(request)
	// Assert
	assert.NotNil(t, err)
}

func Test_SalesOrderUseCase_Export_ShouldError(t *testing.T) {
	// Arrange
	salesOrderUseCase := newSalesOrderUseCase(false)
	request := &models.SalesOrderExportRequest{}
	request.ID = 1

	// Act
	_, err := salesOrderUseCase.Export(request, salesOrderUseCase.ctx)
	// Assert
	assert.NotNil(t, err)
}
