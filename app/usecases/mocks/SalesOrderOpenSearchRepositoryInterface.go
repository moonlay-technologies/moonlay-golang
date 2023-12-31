// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	model "order-service/global/utils/model"

	mock "github.com/stretchr/testify/mock"

	models "order-service/app/models"
	repositories "order-service/app/repositories/open_search"
)

// SalesOrderOpenSearchRepositoryInterface is an autogenerated mock type for the SalesOrderOpenSearchRepositoryInterface type
type SalesOrderOpenSearchRepositoryInterface struct {
	mock.Mock
	repositories.SalesOrderOpenSearchRepositoryInterface
}

// Create provides a mock function with given fields: request, result
func (_m *SalesOrderOpenSearchRepositoryInterface) Create(request *models.SalesOrder, result chan *models.SalesOrderChan) {
	_m.Called(request, result)
}

// Get provides a mock function with given fields: request, result
func (_m *SalesOrderOpenSearchRepositoryInterface) Get(request *models.SalesOrderRequest, result chan *models.SalesOrdersChan) {
	_m.Called(request, result)
}

// GetByAgentID provides a mock function with given fields: request, result
func (_m *SalesOrderOpenSearchRepositoryInterface) GetByAgentID(request *models.SalesOrderRequest, result chan *models.SalesOrdersChan) {
	_m.Called(request, result)
}

// GetByID provides a mock function with given fields: request, result
func (_m *SalesOrderOpenSearchRepositoryInterface) GetByID(request *models.SalesOrderRequest, result chan *models.SalesOrderChan) {
	_m.Called(request, result)
}

// GetByOrderSourceID provides a mock function with given fields: request, result
func (_m *SalesOrderOpenSearchRepositoryInterface) GetByOrderSourceID(request *models.SalesOrderRequest, result chan *models.SalesOrdersChan) {
	_m.Called(request, result)
}

// GetByOrderStatusID provides a mock function with given fields: request, result
func (_m *SalesOrderOpenSearchRepositoryInterface) GetByOrderStatusID(request *models.SalesOrderRequest, result chan *models.SalesOrdersChan) {
	_m.Called(request, result)
}

// GetBySalesmanID provides a mock function with given fields: request, result
func (_m *SalesOrderOpenSearchRepositoryInterface) GetBySalesmanID(request *models.SalesOrderRequest, result chan *models.SalesOrdersChan) {
	_m.Called(request, result)
}

// GetByStoreID provides a mock function with given fields: request, result
func (_m *SalesOrderOpenSearchRepositoryInterface) GetByStoreID(request *models.SalesOrderRequest, result chan *models.SalesOrdersChan) {
	_m.Called(request, result)
}

// generateSalesOrderQueryOpenSearchResult provides a mock function with given fields: openSearchQueryJson, withSalesOrderDetails
func (_m *SalesOrderOpenSearchRepositoryInterface) generateSalesOrderQueryOpenSearchResult(openSearchQueryJson []byte, withSalesOrderDetails bool) (*models.SalesOrders, *model.ErrorLog) {
	ret := _m.Called(openSearchQueryJson, withSalesOrderDetails)

	var r0 *models.SalesOrders
	if rf, ok := ret.Get(0).(func([]byte, bool) *models.SalesOrders); ok {
		r0 = rf(openSearchQueryJson, withSalesOrderDetails)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SalesOrders)
		}
	}

	var r1 *model.ErrorLog
	if rf, ok := ret.Get(1).(func([]byte, bool) *model.ErrorLog); ok {
		r1 = rf(openSearchQueryJson, withSalesOrderDetails)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// generateSalesOrderQueryOpenSearchTermRequest provides a mock function with given fields: term_field, term_value, request
func (_m *SalesOrderOpenSearchRepositoryInterface) generateSalesOrderQueryOpenSearchTermRequest(term_field string, term_value interface{}, request *models.SalesOrderRequest) []byte {
	ret := _m.Called(term_field, term_value, request)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, interface{}, *models.SalesOrderRequest) []byte); ok {
		r0 = rf(term_field, term_value, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

type mockConstructorTestingTNewSalesOrderOpenSearchRepositoryInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewSalesOrderOpenSearchRepositoryInterface creates a new instance of SalesOrderOpenSearchRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSalesOrderOpenSearchRepositoryInterface(t mockConstructorTestingTNewSalesOrderOpenSearchRepositoryInterface) *SalesOrderOpenSearchRepositoryInterface {
	mock := &SalesOrderOpenSearchRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
