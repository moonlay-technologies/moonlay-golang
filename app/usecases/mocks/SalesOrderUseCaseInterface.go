// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"
	model "order-service/global/utils/model"

	mock "github.com/stretchr/testify/mock"

	models "order-service/app/models"

	sql "database/sql"
)

// SalesOrderUseCaseInterface is an autogenerated mock type for the SalesOrderUseCaseInterface type
type SalesOrderUseCaseInterface struct {
	mock.Mock
}

// Create provides a mock function with given fields: request, sqlTransaction, ctx
func (_m *SalesOrderUseCaseInterface) Create(request *models.SalesOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderResponse, *model.ErrorLog) {
	ret := _m.Called(request, sqlTransaction, ctx)

	var r0 *models.SalesOrderResponse
	// if rf, ok := ret.Get(0).(func(*models.SalesOrderStoreRequest, *sql.Tx, context.Context) []*models.SalesOrder); ok {
	// 	r0 = rf(request, sqlTransaction, ctx)
	// } else {
	// 	if ret.Get(0) != nil {
	// 		r0 = ret.Get(0).([]*models.SalesOrder)
	// 	}
	// }

	var r1 *model.ErrorLog
	if rf, ok := ret.Get(1).(func(*models.SalesOrderStoreRequest, *sql.Tx, context.Context) *model.ErrorLog); ok {
		r1 = rf(request, sqlTransaction, ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// Get provides a mock function with given fields: request
func (_m *SalesOrderUseCaseInterface) Get(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog) {
	ret := _m.Called(request)

	var r0 *models.SalesOrders
	if rf, ok := ret.Get(0).(func(*models.SalesOrderRequest) *models.SalesOrders); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SalesOrders)
		}
	}

	var r1 *model.ErrorLog
	if rf, ok := ret.Get(1).(func(*models.SalesOrderRequest) *model.ErrorLog); ok {
		r1 = rf(request)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetByAgentID provides a mock function with given fields: request
func (_m *SalesOrderUseCaseInterface) GetByAgentID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog) {
	ret := _m.Called(request)

	var r0 *models.SalesOrders
	if rf, ok := ret.Get(0).(func(*models.SalesOrderRequest) *models.SalesOrders); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SalesOrders)
		}
	}

	var r1 *model.ErrorLog
	if rf, ok := ret.Get(1).(func(*models.SalesOrderRequest) *model.ErrorLog); ok {
		r1 = rf(request)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: request, withDetail, ctx
func (_m *SalesOrderUseCaseInterface) GetByID(request *models.SalesOrderRequest, withDetail bool, ctx context.Context) (*models.SalesOrder, *model.ErrorLog) {
	ret := _m.Called(request, withDetail, ctx)

	var r0 *models.SalesOrder
	if rf, ok := ret.Get(0).(func(*models.SalesOrderRequest, bool, context.Context) *models.SalesOrder); ok {
		r0 = rf(request, withDetail, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SalesOrder)
		}
	}

	var r1 *model.ErrorLog
	if rf, ok := ret.Get(1).(func(*models.SalesOrderRequest, bool, context.Context) *model.ErrorLog); ok {
		r1 = rf(request, withDetail, ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetByOrderSourceID provides a mock function with given fields: request
func (_m *SalesOrderUseCaseInterface) GetByOrderSourceID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog) {
	ret := _m.Called(request)

	var r0 *models.SalesOrders
	if rf, ok := ret.Get(0).(func(*models.SalesOrderRequest) *models.SalesOrders); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SalesOrders)
		}
	}

	var r1 *model.ErrorLog
	if rf, ok := ret.Get(1).(func(*models.SalesOrderRequest) *model.ErrorLog); ok {
		r1 = rf(request)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetByOrderStatusID provides a mock function with given fields: request
func (_m *SalesOrderUseCaseInterface) GetByOrderStatusID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog) {
	ret := _m.Called(request)

	var r0 *models.SalesOrders
	if rf, ok := ret.Get(0).(func(*models.SalesOrderRequest) *models.SalesOrders); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SalesOrders)
		}
	}

	var r1 *model.ErrorLog
	if rf, ok := ret.Get(1).(func(*models.SalesOrderRequest) *model.ErrorLog); ok {
		r1 = rf(request)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetBySalesmanID provides a mock function with given fields: request
func (_m *SalesOrderUseCaseInterface) GetBySalesmanID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog) {
	ret := _m.Called(request)

	var r0 *models.SalesOrders
	if rf, ok := ret.Get(0).(func(*models.SalesOrderRequest) *models.SalesOrders); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SalesOrders)
		}
	}

	var r1 *model.ErrorLog
	if rf, ok := ret.Get(1).(func(*models.SalesOrderRequest) *model.ErrorLog); ok {
		r1 = rf(request)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetByStoreID provides a mock function with given fields: request
func (_m *SalesOrderUseCaseInterface) GetByStoreID(request *models.SalesOrderRequest) (*models.SalesOrders, *model.ErrorLog) {
	ret := _m.Called(request)

	var r0 *models.SalesOrders
	if rf, ok := ret.Get(0).(func(*models.SalesOrderRequest) *models.SalesOrders); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SalesOrders)
		}
	}

	var r1 *model.ErrorLog
	if rf, ok := ret.Get(1).(func(*models.SalesOrderRequest) *model.ErrorLog); ok {
		r1 = rf(request)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// SyncToOpenSearchFromCreateEvent provides a mock function with given fields: salesOrder, sqlTransaction, ctx
func (_m *SalesOrderUseCaseInterface) SyncToOpenSearchFromCreateEvent(salesOrder *models.SalesOrder, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog {
	ret := _m.Called(salesOrder, sqlTransaction, ctx)

	var r0 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.SalesOrder, *sql.Tx, context.Context) *model.ErrorLog); ok {
		r0 = rf(salesOrder, sqlTransaction, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ErrorLog)
		}
	}

	return r0
}

// SyncToOpenSearchFromUpdateEvent provides a mock function with given fields: salesOrder, ctx
func (_m *SalesOrderUseCaseInterface) SyncToOpenSearchFromUpdateEvent(salesOrder *models.SalesOrder, ctx context.Context) *model.ErrorLog {
	ret := _m.Called(salesOrder, ctx)

	var r0 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.SalesOrder, context.Context) *model.ErrorLog); ok {
		r0 = rf(salesOrder, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ErrorLog)
		}
	}

	return r0
}

func (_m *SalesOrderUseCaseInterface) UpdateById(id int, request *models.SalesOrderUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrder, *model.ErrorLog) {
	return nil, nil
}

func (_m *SalesOrderUseCaseInterface) UpdateSODetailById(id int, request *models.SalesOrderDetailUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderDetail, *model.ErrorLog) {
	return nil, nil
}

func (_m *SalesOrderUseCaseInterface) UpdateSODetailBySOId(SoId int, request []*models.SalesOrderDetailUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) ([]*models.SalesOrder, *model.ErrorLog) {
	return nil, nil
}

type mockConstructorTestingTNewSalesOrderUseCaseInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewSalesOrderUseCaseInterface creates a new instance of SalesOrderUseCaseInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSalesOrderUseCaseInterface(t mockConstructorTestingTNewSalesOrderUseCaseInterface) *SalesOrderUseCaseInterface {
	mock := &SalesOrderUseCaseInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}


