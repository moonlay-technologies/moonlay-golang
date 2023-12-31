// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"
	models "order-service/app/models"

	mock "github.com/stretchr/testify/mock"
)

// StoreRepositoryInterface is an autogenerated mock type for the StoreRepositoryInterface type
type StoreRepositoryInterface struct {
	mock.Mock
}

// GetByID provides a mock function with given fields: ID, countOnly, ctx, result
func (_m *StoreRepositoryInterface) GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.StoreChan) {
	_m.Called(ID, countOnly, ctx, result)
}

func (_m *StoreRepositoryInterface) GetIdByStoreCode(storeCode string, countOnly bool, ctx context.Context, result chan *models.StoreChan) {
	_m.Called(storeCode, countOnly, ctx, result)
}

type mockConstructorTestingTNewStoreRepositoryInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewStoreRepositoryInterface creates a new instance of StoreRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStoreRepositoryInterface(t mockConstructorTestingTNewStoreRepositoryInterface) *StoreRepositoryInterface {
	mock := &StoreRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
