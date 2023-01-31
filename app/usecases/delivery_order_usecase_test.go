package usecases

import (
	"context"
	"fmt"
	"os"
	"poc-order-service/app/models"
	repositories "poc-order-service/app/repositories/open_search"
	"poc-order-service/app/usecases/mocks"
	"poc-order-service/global/utils/helper"
	"poc-order-service/global/utils/opensearch_dbo"
	"poc-order-service/global/utils/sqldb"
	"testing"

	"github.com/bxcodec/dbresolver"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
)

func newDeliveryOrderUsecase(status bool) deliveryOrderUseCase {
	var ctx context.Context
	mockDeliveryOrderRepository := &mocks.DeliveryOrderRepositoryInterface{}
	mockDeliveryOrderDetailRepository := &mocks.DeliveryOrderDetailRepositoryInterface{}
	mockSalesOrderRepository := &mocks.SalesOrderRepositoryInterface{}
	mockSalesOrderDetailRepository := &mocks.SalesOrderDetailRepositoryInterface{}
	mockOrderStatusRepository := &mocks.OrderStatusRepositoryInterface{}
	mockOrderSourceRepository := &mocks.OrderSourceRepositoryInterface{}
	mockWarehouseRepository := &mocks.WarehouseRepositoryInterface{}
	mockBrandRepository := &mocks.BrandRepositoryInterface{}
	mockUomRepository := &mocks.UomRepositoryInterface{}
	mockProductRepository := &mocks.ProductRepositoryInterface{}
	mockAgentRepository := &mocks.AgentRepositoryInterface{}
	mockStoreRepository := &mocks.StoreRepositoryInterface{}
	mockDeliveryOrderLogRepository := &mocks.DeliveryOrderLogRepositoryInterface{}
	mockSalesOrderUseCase := &mocks.SalesOrderUseCaseInterface{}
	mockKafkaClient := &mocks.KafkaClientInterface{}
	openSearchHosts := []string{os.Getenv("OPENSEARCH_HOST_01")}
	openSearchClient := opensearch_dbo.InitOpenSearchClientInterface(openSearchHosts, os.Getenv("OPENSEARCH_USERNAME"), os.Getenv("OPENSEARCH_PASSWORD"), ctx)
	deliveryOrderOpenSearch := repositories.InitDeliveryOrderOpenSearchRepository(openSearchClient)

	// mockDeliveryOrderOpenSearchRepository := &mocks.DeliveryOrderOpenSearchRepositoryInterface{
	// 	DeliveryOrderOpenSearchRepositoryInterface: repositories.InitDeliveryOrderOpenSearchRepository(openSearchClient),
	// }
	mockSalesOrderOpenSearchRepository := &mocks.SalesOrderOpenSearchRepositoryInterface{}

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

	dbConnection := dbresolver.WrapDBs(mysqlWrite.DB(), mysqlRead.DB())
	defer dbConnection.Close()

	// mocking
	// mockDeliveryOrderOpenSearchRepository.On("GetByID", mock.Anything, mock.Anything).Return(mockDeliveryOrderOpenSearchRepository)

	return deliveryOrderUseCase{
		deliveryOrderRepository:           mockDeliveryOrderRepository,
		deliveryOrderDetailRepository:     mockDeliveryOrderDetailRepository,
		salesOrderRepository:              mockSalesOrderRepository,
		salesOrderDetailRepository:        mockSalesOrderDetailRepository,
		orderStatusRepository:             mockOrderStatusRepository,
		orderSourceRepository:             mockOrderSourceRepository,
		warehouseRepository:               mockWarehouseRepository,
		brandRepository:                   mockBrandRepository,
		uomRepository:                     mockUomRepository,
		productRepository:                 mockProductRepository,
		agentRepository:                   mockAgentRepository,
		storeRepository:                   mockStoreRepository,
		deliveryOrderLogRepository:        mockDeliveryOrderLogRepository,
		salesOrderUseCase:                 mockSalesOrderUseCase,
		kafkaClient:                       mockKafkaClient,
		deliveryOrderOpenSearchRepository: deliveryOrderOpenSearch,
		salesOrderOpenSearchRepository:    mockSalesOrderOpenSearchRepository,
		db:                                dbConnection,
		ctx:                               ctx,
	}
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
