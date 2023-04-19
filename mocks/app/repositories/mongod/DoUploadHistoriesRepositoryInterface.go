// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	models "order-service/app/models"

	mock "github.com/stretchr/testify/mock"
)

// DoUploadHistoriesRepositoryInterface is an autogenerated mock type for the DoUploadHistoriesRepositoryInterface type
type DoUploadHistoriesRepositoryInterface struct {
	mock.Mock
}

// Get provides a mock function with given fields: request, countOnly, ctx, resultChan
func (_m *DoUploadHistoriesRepositoryInterface) Get(request *models.GetDoUploadHistoriesRequest, countOnly bool, ctx context.Context, resultChan chan *models.DoUploadHistoriesChan) {
	_m.Called(request, countOnly, ctx, resultChan)
}

// GetByHistoryID provides a mock function with given fields: ID, countOnly, ctx, resultChan
func (_m *DoUploadHistoriesRepositoryInterface) GetByHistoryID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.GetDoUploadHistoryResponseChan) {
	_m.Called(ID, countOnly, ctx, resultChan)
}

// GetByID provides a mock function with given fields: ID, countOnly, ctx, resultChan
func (_m *DoUploadHistoriesRepositoryInterface) GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.DoUploadHistoryChan) {
	_m.Called(ID, countOnly, ctx, resultChan)
}

// Insert provides a mock function with given fields: request, ctx, resultChan
func (_m *DoUploadHistoriesRepositoryInterface) Insert(request *models.DoUploadHistory, ctx context.Context, resultChan chan *models.DoUploadHistoryChan) {
	_m.Called(request, ctx, resultChan)
}

// UpdateByID provides a mock function with given fields: ID, request, ctx, resultChan
func (_m *DoUploadHistoriesRepositoryInterface) UpdateByID(ID string, request *models.DoUploadHistory, ctx context.Context, resultChan chan *models.DoUploadHistoryChan) {
	_m.Called(ID, request, ctx, resultChan)
}

type mockConstructorTestingTNewDoUploadHistoriesRepositoryInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewDoUploadHistoriesRepositoryInterface creates a new instance of DoUploadHistoriesRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDoUploadHistoriesRepositoryInterface(t mockConstructorTestingTNewDoUploadHistoriesRepositoryInterface) *DoUploadHistoriesRepositoryInterface {
	mock := &DoUploadHistoriesRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
