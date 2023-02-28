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
func (_m *SalesOrderUseCaseInterface) Create(request *models.SalesOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) ([]*models.SalesOrderResponse, *model.ErrorLog) {
	ret := _m.Called(request, sqlTransaction, ctx)

	var r0 []*models.SalesOrderResponse
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
func (_m *SalesOrderUseCaseInterface) Get(request *models.SalesOrderRequest) (*models.SalesOrdersOpenSearchResponse, *model.ErrorLog) {
	ret := _m.Called(request)

	var r0 *models.SalesOrdersOpenSearchResponse
	if rf, ok := ret.Get(0).(func(*models.SalesOrderRequest) *models.SalesOrdersOpenSearchResponse); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SalesOrdersOpenSearchResponse)
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
func (_m *SalesOrderUseCaseInterface) GetByID(request *models.SalesOrderRequest,  ctx context.Context) ([]*models.SalesOrderOpenSearchResponse, *model.ErrorLog) {
	ret := _m.Called(request,  ctx)

	var r0 []*models.SalesOrderOpenSearchResponse
	if rf, ok := ret.Get(0).(func(*models.SalesOrderRequest, context.Context) []*models.SalesOrderOpenSearchResponse); ok {
		r0 = rf(request,  ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.SalesOrderOpenSearchResponse)
		}
	}

	var r1 *model.ErrorLog
	if rf, ok := ret.Get(1).(func(*models.SalesOrderRequest, context.Context) *model.ErrorLog); ok {
		r1 = rf(request,  ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

func (_m *SalesOrderUseCaseInterface) GetByIDWithDetail(request *models.SalesOrderRequest, ctx context.Context) (*models.SalesOrder, *model.ErrorLog) {
	ret := _m.Called(request, ctx)

	var r0 *models.SalesOrder
	if rf, ok := ret.Get(0).(func(*models.SalesOrderRequest, context.Context) *models.SalesOrder); ok {
		r0 = rf(request, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SalesOrder)
		}
	}

	var r1 *model.ErrorLog
	if rf, ok := ret.Get(1).(func(*models.SalesOrderRequest, context.Context) *model.ErrorLog); ok {
		r1 = rf(request, ctx)
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

func (_m *SalesOrderUseCaseInterface) UpdateById(id int, request *models.SalesOrderUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderResponse, *model.ErrorLog) {
	return nil, nil
}

func (_m *SalesOrderUseCaseInterface) UpdateSODetailById(soId, id int, request *models.SalesOrderDetailUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderDetail, *model.ErrorLog) {
	return nil, nil
}

func (_m *SalesOrderUseCaseInterface) UpdateSODetailBySOId(soId int, request *models.SalesOrderUpdateRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.SalesOrderResponse, *model.ErrorLog) {
	return nil, nil
}

func (_m *SalesOrderUseCaseInterface) GetDetails(request *models.SalesOrderRequest) (*models.SalesOrderDetailsOpenSearchResponse, *model.ErrorLog)  {
	ret := _m.Called(request)

	var r0 *models.SalesOrderDetailsOpenSearchResponse
	if rf, ok := ret.Get(0).(func(*models.SalesOrderRequest) *models.SalesOrderDetailsOpenSearchResponse); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SalesOrderDetailsOpenSearchResponse)
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

func (_m *SalesOrderUseCaseInterface) DeleteById(id int, sqlTransaction *sql.Tx) (*model.ErrorLog) {
	return nil
}

func (_m *SalesOrderUseCaseInterface) GetDetailById(id int) (*models.SalesOrderDetailOpenSearchResponse, *model.ErrorLog) {
	return nil, nil
}