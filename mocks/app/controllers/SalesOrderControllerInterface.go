// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"
)

// SalesOrderControllerInterface is an autogenerated mock type for the SalesOrderControllerInterface type
type SalesOrderControllerInterface struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx
func (_m *SalesOrderControllerInterface) Create(ctx *gin.Context) {
	_m.Called(ctx)
}

// Get provides a mock function with given fields: ctx
func (_m *SalesOrderControllerInterface) Get(ctx *gin.Context) {
	_m.Called(ctx)
}

// GetByID provides a mock function with given fields: ctx
func (_m *SalesOrderControllerInterface) GetByID(ctx *gin.Context) {
	_m.Called(ctx)
}

type mockConstructorTestingTNewSalesOrderControllerInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewSalesOrderControllerInterface creates a new instance of SalesOrderControllerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSalesOrderControllerInterface(t mockConstructorTestingTNewSalesOrderControllerInterface) *SalesOrderControllerInterface {
	mock := &SalesOrderControllerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
