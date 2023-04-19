// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	models "order-service/app/models"

	mock "github.com/stretchr/testify/mock"

	sql "database/sql"
)

// SalesOrderDetailRepositoryInterface is an autogenerated mock type for the SalesOrderDetailRepositoryInterface type
type SalesOrderDetailRepositoryInterface struct {
	mock.Mock
}

// DeleteByID provides a mock function with given fields: request, sqlTransaction, ctx, result
func (_m *SalesOrderDetailRepositoryInterface) DeleteByID(request *models.SalesOrderDetail, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.SalesOrderDetailChan) {
	_m.Called(request, sqlTransaction, ctx, result)
}

// GetByID provides a mock function with given fields: salesOrderDetailID, countOnly, ctx, result
func (_m *SalesOrderDetailRepositoryInterface) GetByID(salesOrderDetailID int, countOnly bool, ctx context.Context, result chan *models.SalesOrderDetailChan) {
	_m.Called(salesOrderDetailID, countOnly, ctx, result)
}

// GetBySOIDAndSku provides a mock function with given fields: salesOrderID, sku, countOnly, ctx, resultChan
func (_m *SalesOrderDetailRepositoryInterface) GetBySOIDAndSku(salesOrderID int, sku string, countOnly bool, ctx context.Context, resultChan chan *models.SalesOrderDetailsChan) {
	_m.Called(salesOrderID, sku, countOnly, ctx, resultChan)
}

// GetBySOIDSkuAndUomCode provides a mock function with given fields: salesOrderID, sku, uomCode, countOnly, ctx, result
func (_m *SalesOrderDetailRepositoryInterface) GetBySOIDSkuAndUomCode(salesOrderID int, sku string, uomCode string, countOnly bool, ctx context.Context, result chan *models.SalesOrderDetailChan) {
	_m.Called(salesOrderID, sku, uomCode, countOnly, ctx, result)
}

// GetBySalesOrderID provides a mock function with given fields: salesOrderID, countOnly, ctx, result
func (_m *SalesOrderDetailRepositoryInterface) GetBySalesOrderID(salesOrderID int, countOnly bool, ctx context.Context, result chan *models.SalesOrderDetailsChan) {
	_m.Called(salesOrderID, countOnly, ctx, result)
}

// Insert provides a mock function with given fields: request, sqlTransaction, ctx, result
func (_m *SalesOrderDetailRepositoryInterface) Insert(request *models.SalesOrderDetail, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.SalesOrderDetailChan) {
	_m.Called(request, sqlTransaction, ctx, result)
}

// RemoveCacheByID provides a mock function with given fields: id, ctx, resultChan
func (_m *SalesOrderDetailRepositoryInterface) RemoveCacheByID(id int, ctx context.Context, resultChan chan *models.SalesOrderDetailChan) {
	_m.Called(id, ctx, resultChan)
}

// UpdateByID provides a mock function with given fields: id, request, isInsertToJourney, reason, sqlTransaction, ctx, result
func (_m *SalesOrderDetailRepositoryInterface) UpdateByID(id int, request *models.SalesOrderDetail, isInsertToJourney bool, reason string, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.SalesOrderDetailChan) {
	_m.Called(id, request, isInsertToJourney, reason, sqlTransaction, ctx, result)
}

type mockConstructorTestingTNewSalesOrderDetailRepositoryInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewSalesOrderDetailRepositoryInterface creates a new instance of SalesOrderDetailRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSalesOrderDetailRepositoryInterface(t mockConstructorTestingTNewSalesOrderDetailRepositoryInterface) *SalesOrderDetailRepositoryInterface {
	mock := &SalesOrderDetailRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
