// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UploadDOItemConsumerHandlerInterface is an autogenerated mock type for the UploadDOItemConsumerHandlerInterface type
type UploadDOItemConsumerHandlerInterface struct {
	mock.Mock
}

// ProcessMessage provides a mock function with given fields:
func (_m *UploadDOItemConsumerHandlerInterface) ProcessMessage() {
	_m.Called()
}

type mockConstructorTestingTNewUploadDOItemConsumerHandlerInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewUploadDOItemConsumerHandlerInterface creates a new instance of UploadDOItemConsumerHandlerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUploadDOItemConsumerHandlerInterface(t mockConstructorTestingTNewUploadDOItemConsumerHandlerInterface) *UploadDOItemConsumerHandlerInterface {
	mock := &UploadDOItemConsumerHandlerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
