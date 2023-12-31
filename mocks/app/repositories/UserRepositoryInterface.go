// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	models "order-service/app/models"

	mock "github.com/stretchr/testify/mock"
)

// UserRepositoryInterface is an autogenerated mock type for the UserRepositoryInterface type
type UserRepositoryInterface struct {
	mock.Mock
}

// GetByID provides a mock function with given fields: ID, countOnly, ctx, result
func (_m *UserRepositoryInterface) GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.UserChan) {
	_m.Called(ID, countOnly, ctx, result)
}

type mockConstructorTestingTNewUserRepositoryInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserRepositoryInterface creates a new instance of UserRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserRepositoryInterface(t mockConstructorTestingTNewUserRepositoryInterface) *UserRepositoryInterface {
	mock := &UserRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
