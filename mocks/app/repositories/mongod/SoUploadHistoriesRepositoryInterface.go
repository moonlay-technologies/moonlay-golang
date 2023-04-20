// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	models "order-service/app/models"

	mock "github.com/stretchr/testify/mock"
)

// SoUploadHistoriesRepositoryInterface is an autogenerated mock type for the SoUploadHistoriesRepositoryInterface type
type SoUploadHistoriesRepositoryInterface struct {
	mock.Mock
}

// Get provides a mock function with given fields: request, countOnly, ctx, resultChan
func (_m *SoUploadHistoriesRepositoryInterface) Get(request *models.GetSoUploadHistoriesRequest, countOnly bool, ctx context.Context, resultChan chan *models.SoUploadHistoriesChan) {
	_m.Called(request, countOnly, ctx, resultChan)
}

// GetByHistoryID provides a mock function with given fields: ID, countOnly, ctx, resultChan
func (_m *SoUploadHistoriesRepositoryInterface) GetByHistoryID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.GetSoUploadHistoryResponseChan) {
	_m.Called(ID, countOnly, ctx, resultChan)
}

// GetByID provides a mock function with given fields: ID, countOnly, ctx, resultChan
func (_m *SoUploadHistoriesRepositoryInterface) GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.SoUploadHistoryChan) {
	_m.Called(ID, countOnly, ctx, resultChan)
}

// Insert provides a mock function with given fields: request, ctx, resultChan
func (_m *SoUploadHistoriesRepositoryInterface) Insert(request *models.SoUploadHistory, ctx context.Context, resultChan chan *models.SoUploadHistoryChan) {
	_m.Called(request, ctx, resultChan)
}

// UpdateByID provides a mock function with given fields: ID, request, ctx, resultChan
func (_m *SoUploadHistoriesRepositoryInterface) UpdateByID(ID string, request *models.SoUploadHistory, ctx context.Context, resultChan chan *models.SoUploadHistoryChan) {
	_m.Called(ID, request, ctx, resultChan)
}

type mockConstructorTestingTNewSoUploadHistoriesRepositoryInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewSoUploadHistoriesRepositoryInterface creates a new instance of SoUploadHistoriesRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSoUploadHistoriesRepositoryInterface(t mockConstructorTestingTNewSoUploadHistoriesRepositoryInterface) *SoUploadHistoriesRepositoryInterface {
	mock := &SoUploadHistoriesRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}