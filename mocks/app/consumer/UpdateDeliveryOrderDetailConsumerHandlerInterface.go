// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UpdateDeliveryOrderDetailConsumerHandlerInterface is an autogenerated mock type for the UpdateDeliveryOrderDetailConsumerHandlerInterface type
type UpdateDeliveryOrderDetailConsumerHandlerInterface struct {
	mock.Mock
}

// ProcessMessage provides a mock function with given fields:
func (_m *UpdateDeliveryOrderDetailConsumerHandlerInterface) ProcessMessage() {
	_m.Called()
}

type mockConstructorTestingTNewUpdateDeliveryOrderDetailConsumerHandlerInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewUpdateDeliveryOrderDetailConsumerHandlerInterface creates a new instance of UpdateDeliveryOrderDetailConsumerHandlerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUpdateDeliveryOrderDetailConsumerHandlerInterface(t mockConstructorTestingTNewUpdateDeliveryOrderDetailConsumerHandlerInterface) *UpdateDeliveryOrderDetailConsumerHandlerInterface {
	mock := &UpdateDeliveryOrderDetailConsumerHandlerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
