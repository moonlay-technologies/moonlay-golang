package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"order-service/app/models"
	mocks1 "order-service/app/usecases/mocks"
	"order-service/global/utils/helper"
	"order-service/global/utils/sqldb"
	mocks2 "order-service/mocks/app/repositories"
	"os"
	"testing"

	"github.com/bxcodec/dbresolver"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newCartUsecase(condition string) cartUseCase {
	ctx := context.Background()
	mockOrderStatusRepository := &mocks1.OrderStatusRepositoryInterface{}
	mockCartRepository := &mocks1.CartRepositoryInterface{}
	// mockCartDetailRepository := &mocks.CartDetailRepositoryInterface{}
	mockCartDetailRepo := &mocks2.CartDetailRepositoryInterface{}

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

	if condition == "1" {
		mockOrderStatusRepository.On("GetByNameAndType", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			output := args.Get(4).(chan *models.OrderStatusChan)
			output <- &models.OrderStatusChan{
				OrderStatus: &models.OrderStatus{
					ID: 1,
				},
				Error:    nil,
				ErrorLog: nil,
				ID:       int64(1),
			}
		})

		mockCartRepository.On("Insert", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			output := args.Get(3).(chan *models.CartChan)
			output <- &models.CartChan{
				Cart: &models.Cart{
					ID: 1,
				},
				Error:    nil,
				ErrorLog: nil,
				ID:       int64(1),
			}
		})

		mockCartDetailRepo.On("Insert", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			output := args.Get(3).(chan *models.CartDetailChan)
			output <- &models.CartDetailChan{
				Error:    nil,
				ErrorLog: nil,
				ID:       int64(1),
			}
		})
	}

	if condition == "2" {
		mockOrderStatusRepository.On("GetByNameAndType", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			output := args.Get(4).(chan *models.OrderStatusChan)
			output <- &models.OrderStatusChan{
				OrderStatus: &models.OrderStatus{
					ID: 1,
				},
				Error:    fmt.Errorf("fail to get data"),
				ErrorLog: nil,
				ID:       int64(1),
			}
		})

		mockCartRepository.On("Insert", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			output := args.Get(3).(chan *models.CartChan)
			output <- &models.CartChan{
				Cart: &models.Cart{
					ID: 1,
				},
				Error:    nil,
				ErrorLog: nil,
				ID:       int64(1),
			}
		})

		mockCartDetailRepo.On("Insert", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			output := args.Get(3).(chan *models.CartDetailChan)
			output <- &models.CartDetailChan{
				Error:    nil,
				ErrorLog: nil,
				ID:       int64(1),
			}
		})
	}

	if condition == "3" {
		mockOrderStatusRepository.On("GetByNameAndType", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			output := args.Get(4).(chan *models.OrderStatusChan)
			output <- &models.OrderStatusChan{
				OrderStatus: &models.OrderStatus{
					ID: 1,
				},
				Error:    nil,
				ErrorLog: nil,
				ID:       int64(1),
			}
		})

		mockCartRepository.On("Insert", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			output := args.Get(3).(chan *models.CartChan)
			output <- &models.CartChan{
				Cart: &models.Cart{
					ID: 1,
				},
				Error:    fmt.Errorf("fail to insert data"),
				ErrorLog: nil,
				ID:       int64(1),
			}
		})

		mockCartDetailRepo.On("Insert", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			output := args.Get(3).(chan *models.CartDetailChan)
			output <- &models.CartDetailChan{
				Error:    nil,
				ErrorLog: nil,
				ID:       int64(1),
			}
		})
	}

	if condition == "4" {
		mockOrderStatusRepository.On("GetByNameAndType", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			output := args.Get(4).(chan *models.OrderStatusChan)
			output <- &models.OrderStatusChan{
				OrderStatus: &models.OrderStatus{
					ID: 1,
				},
				Error:    nil,
				ErrorLog: nil,
				ID:       int64(1),
			}
		}).Once()

		mockCartRepository.On("Insert", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			output := args.Get(3).(chan *models.CartChan)
			output <- &models.CartChan{
				Cart: &models.Cart{
					ID: 1,
				},
				Error:    nil,
				ErrorLog: nil,
				ID:       int64(1),
			}
		})

		mockOrderStatusRepository.On("GetByNameAndType", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			output := args.Get(4).(chan *models.OrderStatusChan)
			output <- &models.OrderStatusChan{
				OrderStatus: &models.OrderStatus{
					ID: 1,
				},
				Error:    fmt.Errorf("fail to get data"),
				ErrorLog: nil,
				ID:       int64(1),
			}
		})
	}

	if condition == "5" {
		mockOrderStatusRepository.On("GetByNameAndType", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			output := args.Get(4).(chan *models.OrderStatusChan)
			output <- &models.OrderStatusChan{
				OrderStatus: &models.OrderStatus{
					ID: 1,
				},
				Error:    nil,
				ErrorLog: nil,
				ID:       int64(1),
			}
		}).Once()

		mockCartRepository.On("Insert", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			output := args.Get(3).(chan *models.CartChan)
			output <- &models.CartChan{
				Cart: &models.Cart{
					ID: 1,
				},
				Error:    nil,
				ErrorLog: nil,
				ID:       int64(1),
			}
		})

		mockOrderStatusRepository.On("GetByNameAndType", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			output := args.Get(4).(chan *models.OrderStatusChan)
			output <- &models.OrderStatusChan{
				OrderStatus: &models.OrderStatus{
					ID: 1,
				},
				Error:    nil,
				ErrorLog: nil,
				ID:       int64(1),
			}
		})

		mockCartDetailRepo.On("Insert", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			output := args.Get(3).(chan *models.CartDetailChan)
			output <- &models.CartDetailChan{
				Error:    fmt.Errorf("fail to insert data"),
				ErrorLog: nil,
				ID:       int64(1),
			}
		})
	}

	dbConnection := dbresolver.WrapDBs(mysqlWrite.DB(), mysqlRead.DB())

	return cartUseCase{
		orderStatusRepository: mockOrderStatusRepository,
		cartRepository:        mockCartRepository,
		cartDetailRepository:  mockCartDetailRepo,
		db:                    dbConnection,
		ctx:                   ctx,
	}
}

func Test_CartUseCase_InitCartUseCaseInterface_ShouldSuccess(t *testing.T) {
	// Arrange
	cartUseCase := newCartUsecase("1")
	// Act
	dataDeliveryOrderUseCaseInit := InitCartUseCaseInterface(cartUseCase.cartRepository, cartUseCase.cartDetailRepository, cartUseCase.orderStatusRepository, cartUseCase.db, cartUseCase.ctx)

	// Assert
	assert.NotNil(t, dataDeliveryOrderUseCaseInit)
}

func Test_CartUseCase_Create_ShouldSuccess(t *testing.T) {
	cartUseCase := newCartUsecase("1")
	// Act
	dataDeliveryOrderUseCaseInit, _ := cartUseCase.Create(&models.SalesOrderStoreRequest{
		SalesOrderDetails: []*models.SalesOrderDetailStoreRequest{
			{
				BrandID: 1,
			},
		},
	}, &sql.Tx{}, cartUseCase.ctx)
	assert.NotNil(t, dataDeliveryOrderUseCaseInit)
}

func Test_CartUseCase_Create_ShouldErrGetOrderStatusCart(t *testing.T) {
	cartUseCase := newCartUsecase("2")
	// Act
	dataDeliveryOrderUseCaseInit, _ := cartUseCase.Create(&models.SalesOrderStoreRequest{
		SalesOrderDetails: []*models.SalesOrderDetailStoreRequest{
			{
				BrandID: 9999,
			},
		},
	}, &sql.Tx{}, cartUseCase.ctx)
	assert.NotNil(t, dataDeliveryOrderUseCaseInit)
}

func Test_CartUseCase_Create_ShouldErrInsert(t *testing.T) {
	cartUseCase := newCartUsecase("3")
	// Act
	dataDeliveryOrderUseCaseInit, _ := cartUseCase.Create(&models.SalesOrderStoreRequest{
		SalesOrderDetails: []*models.SalesOrderDetailStoreRequest{
			{
				BrandID: 9999,
			},
		},
	}, &sql.Tx{}, cartUseCase.ctx)
	assert.NotNil(t, dataDeliveryOrderUseCaseInit)
}

func Test_CartUseCase_Create_ShouldErrGetOrderStatusCartDetail(t *testing.T) {
	cartUseCase := newCartUsecase("4")
	// Act
	dataDeliveryOrderUseCaseInit, _ := cartUseCase.Create(&models.SalesOrderStoreRequest{
		SalesOrderDetails: []*models.SalesOrderDetailStoreRequest{
			{
				BrandID: 9999,
			},
		},
	}, &sql.Tx{}, cartUseCase.ctx)
	assert.NotNil(t, dataDeliveryOrderUseCaseInit)
}

func Test_CartUseCase_Create_ShouldErrInsertCartDetail(t *testing.T) {
	cartUseCase := newCartUsecase("5")
	// Act
	dataDeliveryOrderUseCaseInit, _ := cartUseCase.Create(&models.SalesOrderStoreRequest{
		SalesOrderDetails: []*models.SalesOrderDetailStoreRequest{
			{
				BrandID: 9999,
			},
		},
	}, &sql.Tx{}, cartUseCase.ctx)
	assert.NotNil(t, dataDeliveryOrderUseCaseInit)
}
