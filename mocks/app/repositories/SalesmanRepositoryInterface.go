// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	models "order-service/app/models"

	mock "github.com/stretchr/testify/mock"
)

// SalesmanRepositoryInterface is an autogenerated mock type for the SalesmanRepositoryInterface type
type SalesmanRepositoryInterface struct {
	mock.Mock
}

// GetByAgentId provides a mock function with given fields: agentId, countOnly, ctx, resultChan
func (_m *SalesmanRepositoryInterface) GetByAgentId(agentId int, countOnly bool, ctx context.Context, resultChan chan *models.SalesmansChan) {
	_m.Called(agentId, countOnly, ctx, resultChan)
}

// GetByEmail provides a mock function with given fields: email, countOnly, ctx, result
func (_m *SalesmanRepositoryInterface) GetByEmail(email string, countOnly bool, ctx context.Context, result chan *models.SalesmanChan) {
	_m.Called(email, countOnly, ctx, result)
}

// GetByID provides a mock function with given fields: ID, countOnly, ctx, result
func (_m *SalesmanRepositoryInterface) GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.SalesmanChan) {
	_m.Called(ID, countOnly, ctx, result)
}

type mockConstructorTestingTNewSalesmanRepositoryInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewSalesmanRepositoryInterface creates a new instance of SalesmanRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSalesmanRepositoryInterface(t mockConstructorTestingTNewSalesmanRepositoryInterface) *SalesmanRepositoryInterface {
	mock := &SalesmanRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
