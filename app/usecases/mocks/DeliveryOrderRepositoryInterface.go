// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"
	models "order-service/app/models"

	mock "github.com/stretchr/testify/mock"

	sql "database/sql"
)

// DeliveryOrderRepositoryInterface is an autogenerated mock type for the DeliveryOrderRepositoryInterface type
type DeliveryOrderRepositoryInterface struct {
	mock.Mock
}

// GetByID provides a mock function with given fields: id, countOnly, ctx, result
func (_m *DeliveryOrderRepositoryInterface) GetByID(id int, countOnly bool, ctx context.Context, result chan *models.DeliveryOrderChan) {
	_m.Called(id, countOnly, ctx, result)
}

// GetBySalesOrderID provides a mock function with given fields: deliveryOrderID, countOnly, ctx, result
func (_m *DeliveryOrderRepositoryInterface) GetBySalesOrderID(deliveryOrderID int, countOnly bool, ctx context.Context, result chan *models.DeliveryOrdersChan) {
	_m.Called(deliveryOrderID, countOnly, ctx, result)
}

// Insert provides a mock function with given fields: request, sqlTransaction, ctx, result
func (_m *DeliveryOrderRepositoryInterface) Insert(request *models.DeliveryOrder, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.DeliveryOrderChan) {
	_m.Called(request, sqlTransaction, ctx, result)
}

// UpdateByID provides a mock function with given fields: id, deliveryOrder, sqlTransaction, ctx, result
func (_m *DeliveryOrderRepositoryInterface) UpdateByID(id int, deliveryOrder *models.DeliveryOrder, jouneyRemarks string, isInsertToJourney bool, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.DeliveryOrderChan) {
	_m.Called(id, deliveryOrder, sqlTransaction, ctx, result)
}

// UpdateByID provides a mock function with given fields: id, deliveryOrder, sqlTransaction, ctx, result
func (_m *DeliveryOrderRepositoryInterface) GetByDoRefCode(doRefCode string, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderChan) {
	_m.Called(doRefCode, countOnly, ctx, resultChan)
}
type mockConstructorTestingTNewDeliveryOrderRepositoryInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewDeliveryOrderRepositoryInterface creates a new instance of DeliveryOrderRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDeliveryOrderRepositoryInterface(t mockConstructorTestingTNewDeliveryOrderRepositoryInterface) *DeliveryOrderRepositoryInterface {
	mock := &DeliveryOrderRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
func (_m *DeliveryOrderRepositoryInterface) DeleteByID(request *models.DeliveryOrder, ctx context.Context, resultChan chan *models.DeliveryOrderChan) {
	_m.Called(request, ctx, resultChan)
}
