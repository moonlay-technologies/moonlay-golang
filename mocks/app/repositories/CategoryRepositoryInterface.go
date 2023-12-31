// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	models "order-service/app/models"

	mock "github.com/stretchr/testify/mock"
)

// CategoryRepositoryInterface is an autogenerated mock type for the CategoryRepositoryInterface type
type CategoryRepositoryInterface struct {
	mock.Mock
}

// GetByID provides a mock function with given fields: ID, countOnly, ctx, result
func (_m *CategoryRepositoryInterface) GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.CategoryChan) {
	_m.Called(ID, countOnly, ctx, result)
}

// GetByParentID provides a mock function with given fields: parentId, countOnly, ctx, result
func (_m *CategoryRepositoryInterface) GetByParentID(parentId int, countOnly bool, ctx context.Context, result chan *models.CategoryChan) {
	_m.Called(parentId, countOnly, ctx, result)
}

type mockConstructorTestingTNewCategoryRepositoryInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewCategoryRepositoryInterface creates a new instance of CategoryRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCategoryRepositoryInterface(t mockConstructorTestingTNewCategoryRepositoryInterface) *CategoryRepositoryInterface {
	mock := &CategoryRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
