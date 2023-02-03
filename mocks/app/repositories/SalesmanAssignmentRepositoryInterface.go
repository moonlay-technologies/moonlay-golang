// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"
	models "order-service/app/models"

	mock "github.com/stretchr/testify/mock"
)

// SalesmanAssignmentRepositoryInterface is an autogenerated mock type for the SalesmanAssignmentRepositoryInterface type
type SalesmanAssignmentRepositoryInterface struct {
	mock.Mock
}

// GetBySalesmanIDAndAgentID provides a mock function with given fields: agentID, salesmanID, countOnly, ctx, result
func (_m *SalesmanAssignmentRepositoryInterface) GetBySalesmanIDAndAgentID(agentID int, salesmanID int, countOnly bool, ctx context.Context, result chan *models.SalesmanAssignmentChan) {
	_m.Called(agentID, salesmanID, countOnly, ctx, result)
}

type mockConstructorTestingTNewSalesmanAssignmentRepositoryInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewSalesmanAssignmentRepositoryInterface creates a new instance of SalesmanAssignmentRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSalesmanAssignmentRepositoryInterface(t mockConstructorTestingTNewSalesmanAssignmentRepositoryInterface) *SalesmanAssignmentRepositoryInterface {
	mock := &SalesmanAssignmentRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
