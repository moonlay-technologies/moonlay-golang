// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"
	models "order-service/app/models"

	mock "github.com/stretchr/testify/mock"
)

// ProductRepositoryInterface is an autogenerated mock type for the ProductRepositoryInterface type
type ProductRepositoryInterface struct {
	mock.Mock
}

// GetByID provides a mock function with given fields: ID, countOnly, ctx, result
func (_m *ProductRepositoryInterface) GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.ProductChan) {
	_m.Called(ID, countOnly, ctx, result)
}

type mockConstructorTestingTNewProductRepositoryInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewProductRepositoryInterface creates a new instance of ProductRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProductRepositoryInterface(t mockConstructorTestingTNewProductRepositoryInterface) *ProductRepositoryInterface {
	mock := &ProductRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
