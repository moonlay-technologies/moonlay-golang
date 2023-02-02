// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	model "order-service/global/utils/model"

	mock "github.com/stretchr/testify/mock"

	models "order-service/app/models"
)

// DeliveryOrderOpenSearchRepositoryInterface is an autogenerated mock type for the DeliveryOrderOpenSearchRepositoryInterface type
type DeliveryOrderOpenSearchRepositoryInterface struct {
	mock.Mock
}

// Create provides a mock function with given fields: request, result
func (_m *DeliveryOrderOpenSearchRepositoryInterface) Create(request *models.DeliveryOrder, result chan *models.DeliveryOrderChan) {
	_m.Called(request, result)
}

// Get provides a mock function with given fields: request, result
func (_m *DeliveryOrderOpenSearchRepositoryInterface) Get(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan) {
	_m.Called(request, result)
}

// GetByAgentID provides a mock function with given fields: request, result
func (_m *DeliveryOrderOpenSearchRepositoryInterface) GetByAgentID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan) {
	_m.Called(request, result)
}

// GetByID provides a mock function with given fields: request, result
func (_m *DeliveryOrderOpenSearchRepositoryInterface) GetByID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrderChan) {
	_m.Called(request, result)
}

// GetBySalesOrderID provides a mock function with given fields: request, result
func (_m *DeliveryOrderOpenSearchRepositoryInterface) GetBySalesOrderID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan) {
	_m.Called(request, result)
}

// GetBySalesmanID provides a mock function with given fields: request, result
func (_m *DeliveryOrderOpenSearchRepositoryInterface) GetBySalesmanID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan) {
	_m.Called(request, result)
}

// GetByStoreID provides a mock function with given fields: request, result
func (_m *DeliveryOrderOpenSearchRepositoryInterface) GetByStoreID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan) {
	_m.Called(request, result)
}

// generateDeliveryOrderQueryOpenSearchResult provides a mock function with given fields: openSearchQueryJson, withDeliveryOrderDetails
func (_m *DeliveryOrderOpenSearchRepositoryInterface) generateDeliveryOrderQueryOpenSearchResult(openSearchQueryJson []byte, withDeliveryOrderDetails bool) (*models.DeliveryOrders, *model.ErrorLog) {
	ret := _m.Called(openSearchQueryJson, withDeliveryOrderDetails)

	var r0 *models.DeliveryOrders
	if rf, ok := ret.Get(0).(func([]byte, bool) *models.DeliveryOrders); ok {
		r0 = rf(openSearchQueryJson, withDeliveryOrderDetails)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrders)
		}
	}

	var r1 *model.ErrorLog
	if rf, ok := ret.Get(1).(func([]byte, bool) *model.ErrorLog); ok {
		r1 = rf(openSearchQueryJson, withDeliveryOrderDetails)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// generateDeliveryOrderQueryOpenSearchTermRequest provides a mock function with given fields: term_field, term_value, request
func (_m *DeliveryOrderOpenSearchRepositoryInterface) generateDeliveryOrderQueryOpenSearchTermRequest(term_field string, term_value interface{}, request *models.DeliveryOrderRequest) []byte {
	ret := _m.Called(term_field, term_value, request)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, interface{}, *models.DeliveryOrderRequest) []byte); ok {
		r0 = rf(term_field, term_value, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

type mockConstructorTestingTNewDeliveryOrderOpenSearchRepositoryInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewDeliveryOrderOpenSearchRepositoryInterface creates a new instance of DeliveryOrderOpenSearchRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDeliveryOrderOpenSearchRepositoryInterface(t mockConstructorTestingTNewDeliveryOrderOpenSearchRepositoryInterface) *DeliveryOrderOpenSearchRepositoryInterface {
	mock := &DeliveryOrderOpenSearchRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
