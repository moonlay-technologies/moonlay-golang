// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"
)

// StoreControllerInterface is an autogenerated mock type for the StoreControllerInterface type
type StoreControllerInterface struct {
	mock.Mock
}

// GetDeliveryOrders provides a mock function with given fields: ctx
func (_m *StoreControllerInterface) GetDeliveryOrders(ctx *gin.Context) {
	_m.Called(ctx)
}

// GetSalesOrders provides a mock function with given fields: ctx
func (_m *StoreControllerInterface) GetSalesOrders(ctx *gin.Context) {
	_m.Called(ctx)
}

type mockConstructorTestingTNewStoreControllerInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewStoreControllerInterface creates a new instance of StoreControllerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStoreControllerInterface(t mockConstructorTestingTNewStoreControllerInterface) *StoreControllerInterface {
	mock := &StoreControllerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
