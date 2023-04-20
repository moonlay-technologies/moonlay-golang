// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	models "order-service/app/models"

	mock "github.com/stretchr/testify/mock"
)

// SOSJUploadHistoriesRepositoryInterface is an autogenerated mock type for the SOSJUploadHistoriesRepositoryInterface type
type SOSJUploadHistoriesRepositoryInterface struct {
	mock.Mock
}

// Get provides a mock function with given fields: request, countOnly, ctx, resultChan
func (_m *SOSJUploadHistoriesRepositoryInterface) Get(request *models.GetSosjUploadHistoriesRequest, countOnly bool, ctx context.Context, resultChan chan *models.UploadHistoriesChan) {
	_m.Called(request, countOnly, ctx, resultChan)
}

// GetByID provides a mock function with given fields: ID, countOnly, ctx, resultChan
func (_m *SOSJUploadHistoriesRepositoryInterface) GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.GetSosjUploadHistoryResponseChan) {
	_m.Called(ID, countOnly, ctx, resultChan)
}

// Insert provides a mock function with given fields: request, ctx, resultChan
func (_m *SOSJUploadHistoriesRepositoryInterface) Insert(request *models.UploadHistory, ctx context.Context, resultChan chan *models.UploadHistoryChan) {
	_m.Called(request, ctx, resultChan)
}

// UpdateByID provides a mock function with given fields: ID, request, ctx, resultChan
func (_m *SOSJUploadHistoriesRepositoryInterface) UpdateByID(ID string, request *models.UploadHistory, ctx context.Context, resultChan chan *models.UploadHistoryChan) {
	_m.Called(ID, request, ctx, resultChan)
}

type mockConstructorTestingTNewSOSJUploadHistoriesRepositoryInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewSOSJUploadHistoriesRepositoryInterface creates a new instance of SOSJUploadHistoriesRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSOSJUploadHistoriesRepositoryInterface(t mockConstructorTestingTNewSOSJUploadHistoriesRepositoryInterface) *SOSJUploadHistoriesRepositoryInterface {
	mock := &SOSJUploadHistoriesRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}