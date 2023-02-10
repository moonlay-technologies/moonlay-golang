// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"
)

// AgentControllerInterface is an autogenerated mock type for the AgentControllerInterface type
type AgentControllerInterface struct {
	mock.Mock
}

// GetDeliveryOrders provides a mock function with given fields: ctx
func (_m *AgentControllerInterface) GetDeliveryOrders(ctx *gin.Context) {
	_m.Called(ctx)
}

// GetSalesOrders provides a mock function with given fields: ctx
func (_m *AgentControllerInterface) GetSalesOrders(ctx *gin.Context) {
	_m.Called(ctx)
}

type mockConstructorTestingTNewAgentControllerInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewAgentControllerInterface creates a new instance of AgentControllerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAgentControllerInterface(t mockConstructorTestingTNewAgentControllerInterface) *AgentControllerInterface {
	mock := &AgentControllerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}