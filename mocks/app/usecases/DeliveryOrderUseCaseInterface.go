// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	model "order-service/global/utils/model"

	mock "github.com/stretchr/testify/mock"

	models "order-service/app/models"

	sql "database/sql"
)

// DeliveryOrderUseCaseInterface is an autogenerated mock type for the DeliveryOrderUseCaseInterface type
type DeliveryOrderUseCaseInterface struct {
	mock.Mock
}

// Create provides a mock function with given fields: request, sqlTransaction, ctx
func (_m *DeliveryOrderUseCaseInterface) Create(request *models.DeliveryOrderStoreRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderStoreResponse, *model.ErrorLog) {
	ret := _m.Called(request, sqlTransaction, ctx)

	var r0 *models.DeliveryOrderStoreResponse
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderStoreRequest, *sql.Tx, context.Context) (*models.DeliveryOrderStoreResponse, *model.ErrorLog)); ok {
		return rf(request, sqlTransaction, ctx)
	}
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderStoreRequest, *sql.Tx, context.Context) *models.DeliveryOrderStoreResponse); ok {
		r0 = rf(request, sqlTransaction, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrderStoreResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.DeliveryOrderStoreRequest, *sql.Tx, context.Context) *model.ErrorLog); ok {
		r1 = rf(request, sqlTransaction, ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// DeleteByID provides a mock function with given fields: deliveryOrderId, sqlTransaction, ctx
func (_m *DeliveryOrderUseCaseInterface) DeleteByID(deliveryOrderId int, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog {
	ret := _m.Called(deliveryOrderId, sqlTransaction, ctx)

	var r0 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(int, *sql.Tx, context.Context) *model.ErrorLog); ok {
		r0 = rf(deliveryOrderId, sqlTransaction, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ErrorLog)
		}
	}

	return r0
}

// DeleteDetailByDoID provides a mock function with given fields: deliveryOrderId, sqlTransaction, ctx
func (_m *DeliveryOrderUseCaseInterface) DeleteDetailByDoID(deliveryOrderId int, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog {
	ret := _m.Called(deliveryOrderId, sqlTransaction, ctx)

	var r0 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(int, *sql.Tx, context.Context) *model.ErrorLog); ok {
		r0 = rf(deliveryOrderId, sqlTransaction, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ErrorLog)
		}
	}

	return r0
}

// DeleteDetailByID provides a mock function with given fields: deliveryOrderDetailId, sqlTransaction, ctx
func (_m *DeliveryOrderUseCaseInterface) DeleteDetailByID(deliveryOrderDetailId int, sqlTransaction *sql.Tx, ctx context.Context) *model.ErrorLog {
	ret := _m.Called(deliveryOrderDetailId, sqlTransaction, ctx)

	var r0 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(int, *sql.Tx, context.Context) *model.ErrorLog); ok {
		r0 = rf(deliveryOrderDetailId, sqlTransaction, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ErrorLog)
		}
	}

	return r0
}

// Export provides a mock function with given fields: request, ctx
func (_m *DeliveryOrderUseCaseInterface) Export(request *models.DeliveryOrderExportRequest, ctx context.Context) (string, *model.ErrorLog) {
	ret := _m.Called(request, ctx)

	var r0 string
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderExportRequest, context.Context) (string, *model.ErrorLog)); ok {
		return rf(request, ctx)
	}
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderExportRequest, context.Context) string); ok {
		r0 = rf(request, ctx)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*models.DeliveryOrderExportRequest, context.Context) *model.ErrorLog); ok {
		r1 = rf(request, ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// ExportDetail provides a mock function with given fields: request, ctx
func (_m *DeliveryOrderUseCaseInterface) ExportDetail(request *models.DeliveryOrderDetailExportRequest, ctx context.Context) (string, *model.ErrorLog) {
	ret := _m.Called(request, ctx)

	var r0 string
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderDetailExportRequest, context.Context) (string, *model.ErrorLog)); ok {
		return rf(request, ctx)
	}
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderDetailExportRequest, context.Context) string); ok {
		r0 = rf(request, ctx)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*models.DeliveryOrderDetailExportRequest, context.Context) *model.ErrorLog); ok {
		r1 = rf(request, ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// Get provides a mock function with given fields: request
func (_m *DeliveryOrderUseCaseInterface) Get(request *models.DeliveryOrderRequest) (*models.DeliveryOrdersOpenSearchResponse, *model.ErrorLog) {
	ret := _m.Called(request)

	var r0 *models.DeliveryOrdersOpenSearchResponse
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest) (*models.DeliveryOrdersOpenSearchResponse, *model.ErrorLog)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest) *models.DeliveryOrdersOpenSearchResponse); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrdersOpenSearchResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.DeliveryOrderRequest) *model.ErrorLog); ok {
		r1 = rf(request)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetByAgentID provides a mock function with given fields: request
func (_m *DeliveryOrderUseCaseInterface) GetByAgentID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog) {
	ret := _m.Called(request)

	var r0 *models.DeliveryOrders
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest) *models.DeliveryOrders); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrders)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.DeliveryOrderRequest) *model.ErrorLog); ok {
		r1 = rf(request)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: request, ctx
func (_m *DeliveryOrderUseCaseInterface) GetByID(request *models.DeliveryOrderRequest, ctx context.Context) (*models.DeliveryOrder, *model.ErrorLog) {
	ret := _m.Called(request, ctx)

	var r0 *models.DeliveryOrder
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest, context.Context) (*models.DeliveryOrder, *model.ErrorLog)); ok {
		return rf(request, ctx)
	}
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest, context.Context) *models.DeliveryOrder); ok {
		r0 = rf(request, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrder)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.DeliveryOrderRequest, context.Context) *model.ErrorLog); ok {
		r1 = rf(request, ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetByIDWithDetail provides a mock function with given fields: request, ctx
func (_m *DeliveryOrderUseCaseInterface) GetByIDWithDetail(request *models.DeliveryOrderRequest, ctx context.Context) (*models.DeliveryOrderOpenSearchResponse, *model.ErrorLog) {
	ret := _m.Called(request, ctx)

	var r0 *models.DeliveryOrderOpenSearchResponse
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest, context.Context) (*models.DeliveryOrderOpenSearchResponse, *model.ErrorLog)); ok {
		return rf(request, ctx)
	}
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest, context.Context) *models.DeliveryOrderOpenSearchResponse); ok {
		r0 = rf(request, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrderOpenSearchResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.DeliveryOrderRequest, context.Context) *model.ErrorLog); ok {
		r1 = rf(request, ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetByOrderSourceID provides a mock function with given fields: request
func (_m *DeliveryOrderUseCaseInterface) GetByOrderSourceID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog) {
	ret := _m.Called(request)

	var r0 *models.DeliveryOrders
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest) *models.DeliveryOrders); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrders)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.DeliveryOrderRequest) *model.ErrorLog); ok {
		r1 = rf(request)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetByOrderStatusID provides a mock function with given fields: request
func (_m *DeliveryOrderUseCaseInterface) GetByOrderStatusID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog) {
	ret := _m.Called(request)

	var r0 *models.DeliveryOrders
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest) *models.DeliveryOrders); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrders)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.DeliveryOrderRequest) *model.ErrorLog); ok {
		r1 = rf(request)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetBySalesmanID provides a mock function with given fields: request
func (_m *DeliveryOrderUseCaseInterface) GetBySalesmanID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog) {
	ret := _m.Called(request)

	var r0 *models.DeliveryOrders
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest) *models.DeliveryOrders); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrders)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.DeliveryOrderRequest) *model.ErrorLog); ok {
		r1 = rf(request)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetBySalesmansID provides a mock function with given fields: request
func (_m *DeliveryOrderUseCaseInterface) GetBySalesmansID(request *models.DeliveryOrderRequest) (*models.DeliveryOrdersOpenSearchResponses, *model.ErrorLog) {
	ret := _m.Called(request)

	var r0 *models.DeliveryOrdersOpenSearchResponses
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest) (*models.DeliveryOrdersOpenSearchResponses, *model.ErrorLog)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest) *models.DeliveryOrdersOpenSearchResponses); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrdersOpenSearchResponses)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.DeliveryOrderRequest) *model.ErrorLog); ok {
		r1 = rf(request)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetByStoreID provides a mock function with given fields: request
func (_m *DeliveryOrderUseCaseInterface) GetByStoreID(request *models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog) {
	ret := _m.Called(request)

	var r0 *models.DeliveryOrders
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest) (*models.DeliveryOrders, *model.ErrorLog)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderRequest) *models.DeliveryOrders); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrders)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.DeliveryOrderRequest) *model.ErrorLog); ok {
		r1 = rf(request)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetDOJourneys provides a mock function with given fields: request, ctx
func (_m *DeliveryOrderUseCaseInterface) GetDOJourneys(request *models.DeliveryOrderJourneysRequest, ctx context.Context) (*models.DeliveryOrderJourneysResponses, *model.ErrorLog) {
	ret := _m.Called(request, ctx)

	var r0 *models.DeliveryOrderJourneysResponses
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderJourneysRequest, context.Context) (*models.DeliveryOrderJourneysResponses, *model.ErrorLog)); ok {
		return rf(request, ctx)
	}
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderJourneysRequest, context.Context) *models.DeliveryOrderJourneysResponses); ok {
		r0 = rf(request, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrderJourneysResponses)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.DeliveryOrderJourneysRequest, context.Context) *model.ErrorLog); ok {
		r1 = rf(request, ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetDOJourneysByDoID provides a mock function with given fields: doId, ctx
func (_m *DeliveryOrderUseCaseInterface) GetDOJourneysByDoID(doId int, ctx context.Context) (*models.DeliveryOrderJourneysResponses, *model.ErrorLog) {
	ret := _m.Called(doId, ctx)

	var r0 *models.DeliveryOrderJourneysResponses
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(int, context.Context) (*models.DeliveryOrderJourneysResponses, *model.ErrorLog)); ok {
		return rf(doId, ctx)
	}
	if rf, ok := ret.Get(0).(func(int, context.Context) *models.DeliveryOrderJourneysResponses); ok {
		r0 = rf(doId, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrderJourneysResponses)
		}
	}

	if rf, ok := ret.Get(1).(func(int, context.Context) *model.ErrorLog); ok {
		r1 = rf(doId, ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetDOUploadErrorLogsByDoUploadHistoryId provides a mock function with given fields: request, ctx
func (_m *DeliveryOrderUseCaseInterface) GetDOUploadErrorLogsByDoUploadHistoryId(request *models.GetDoUploadErrorLogsRequest, ctx context.Context) (*models.GetDoUploadErrorLogsResponse, *model.ErrorLog) {
	ret := _m.Called(request, ctx)

	var r0 *models.GetDoUploadErrorLogsResponse
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.GetDoUploadErrorLogsRequest, context.Context) (*models.GetDoUploadErrorLogsResponse, *model.ErrorLog)); ok {
		return rf(request, ctx)
	}
	if rf, ok := ret.Get(0).(func(*models.GetDoUploadErrorLogsRequest, context.Context) *models.GetDoUploadErrorLogsResponse); ok {
		r0 = rf(request, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.GetDoUploadErrorLogsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.GetDoUploadErrorLogsRequest, context.Context) *model.ErrorLog); ok {
		r1 = rf(request, ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetDOUploadErrorLogsByReqId provides a mock function with given fields: request, ctx
func (_m *DeliveryOrderUseCaseInterface) GetDOUploadErrorLogsByReqId(request *models.GetDoUploadErrorLogsRequest, ctx context.Context) (*models.GetDoUploadErrorLogsResponse, *model.ErrorLog) {
	ret := _m.Called(request, ctx)

	var r0 *models.GetDoUploadErrorLogsResponse
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.GetDoUploadErrorLogsRequest, context.Context) (*models.GetDoUploadErrorLogsResponse, *model.ErrorLog)); ok {
		return rf(request, ctx)
	}
	if rf, ok := ret.Get(0).(func(*models.GetDoUploadErrorLogsRequest, context.Context) *models.GetDoUploadErrorLogsResponse); ok {
		r0 = rf(request, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.GetDoUploadErrorLogsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.GetDoUploadErrorLogsRequest, context.Context) *model.ErrorLog); ok {
		r1 = rf(request, ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetDOUploadHistories provides a mock function with given fields: request, ctx
func (_m *DeliveryOrderUseCaseInterface) GetDOUploadHistories(request *models.GetDoUploadHistoriesRequest, ctx context.Context) (*models.GetDoUploadHistoryResponses, *model.ErrorLog) {
	ret := _m.Called(request, ctx)

	var r0 *models.GetDoUploadHistoryResponses
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.GetDoUploadHistoriesRequest, context.Context) (*models.GetDoUploadHistoryResponses, *model.ErrorLog)); ok {
		return rf(request, ctx)
	}
	if rf, ok := ret.Get(0).(func(*models.GetDoUploadHistoriesRequest, context.Context) *models.GetDoUploadHistoryResponses); ok {
		r0 = rf(request, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.GetDoUploadHistoryResponses)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.GetDoUploadHistoriesRequest, context.Context) *model.ErrorLog); ok {
		r1 = rf(request, ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetDOUploadHistoriesById provides a mock function with given fields: id, ctx
func (_m *DeliveryOrderUseCaseInterface) GetDOUploadHistoriesById(id string, ctx context.Context) (*models.GetDoUploadHistoryResponse, *model.ErrorLog) {
	ret := _m.Called(id, ctx)

	var r0 *models.GetDoUploadHistoryResponse
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(string, context.Context) (*models.GetDoUploadHistoryResponse, *model.ErrorLog)); ok {
		return rf(id, ctx)
	}
	if rf, ok := ret.Get(0).(func(string, context.Context) *models.GetDoUploadHistoryResponse); ok {
		r0 = rf(id, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.GetDoUploadHistoryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, context.Context) *model.ErrorLog); ok {
		r1 = rf(id, ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetDetailByID provides a mock function with given fields: doDetailID, doID
func (_m *DeliveryOrderUseCaseInterface) GetDetailByID(doDetailID int, doID int) (*models.DeliveryOrderDetailsOpenSearchResponse, *model.ErrorLog) {
	ret := _m.Called(doDetailID, doID)

	var r0 *models.DeliveryOrderDetailsOpenSearchResponse
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(int, int) (*models.DeliveryOrderDetailsOpenSearchResponse, *model.ErrorLog)); ok {
		return rf(doDetailID, doID)
	}
	if rf, ok := ret.Get(0).(func(int, int) *models.DeliveryOrderDetailsOpenSearchResponse); ok {
		r0 = rf(doDetailID, doID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrderDetailsOpenSearchResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int) *model.ErrorLog); ok {
		r1 = rf(doDetailID, doID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetDetails provides a mock function with given fields: request
func (_m *DeliveryOrderUseCaseInterface) GetDetails(request *models.DeliveryOrderDetailOpenSearchRequest) (*models.DeliveryOrderDetailsOpenSearchResponses, *model.ErrorLog) {
	ret := _m.Called(request)

	var r0 *models.DeliveryOrderDetailsOpenSearchResponses
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderDetailOpenSearchRequest) (*models.DeliveryOrderDetailsOpenSearchResponses, *model.ErrorLog)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderDetailOpenSearchRequest) *models.DeliveryOrderDetailsOpenSearchResponses); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrderDetailsOpenSearchResponses)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.DeliveryOrderDetailOpenSearchRequest) *model.ErrorLog); ok {
		r1 = rf(request)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetDetailsByDoID provides a mock function with given fields: request
func (_m *DeliveryOrderUseCaseInterface) GetDetailsByDoID(request *models.DeliveryOrderDetailRequest) (*models.DeliveryOrderOpenSearchResponse, *model.ErrorLog) {
	ret := _m.Called(request)

	var r0 *models.DeliveryOrderOpenSearchResponse
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderDetailRequest) (*models.DeliveryOrderOpenSearchResponse, *model.ErrorLog)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderDetailRequest) *models.DeliveryOrderOpenSearchResponse); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrderOpenSearchResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.DeliveryOrderDetailRequest) *model.ErrorLog); ok {
		r1 = rf(request)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// GetSyncToKafkaHistories provides a mock function with given fields: request, ctx
func (_m *DeliveryOrderUseCaseInterface) GetSyncToKafkaHistories(request *models.DeliveryOrderEventLogRequest, ctx context.Context) ([]*models.DeliveryOrderEventLogResponse, *model.ErrorLog) {
	ret := _m.Called(request, ctx)

	var r0 []*models.DeliveryOrderEventLogResponse
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderEventLogRequest, context.Context) ([]*models.DeliveryOrderEventLogResponse, *model.ErrorLog)); ok {
		return rf(request, ctx)
	}
	if rf, ok := ret.Get(0).(func(*models.DeliveryOrderEventLogRequest, context.Context) []*models.DeliveryOrderEventLogResponse); ok {
		r0 = rf(request, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.DeliveryOrderEventLogResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.DeliveryOrderEventLogRequest, context.Context) *model.ErrorLog); ok {
		r1 = rf(request, ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// RetrySyncToKafka provides a mock function with given fields: logId
func (_m *DeliveryOrderUseCaseInterface) RetrySyncToKafka(logId string) (*models.DORetryProcessSyncToKafkaResponse, *model.ErrorLog) {
	ret := _m.Called(logId)

	var r0 *models.DORetryProcessSyncToKafkaResponse
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(string) (*models.DORetryProcessSyncToKafkaResponse, *model.ErrorLog)); ok {
		return rf(logId)
	}
	if rf, ok := ret.Get(0).(func(string) *models.DORetryProcessSyncToKafkaResponse); ok {
		r0 = rf(logId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DORetryProcessSyncToKafkaResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string) *model.ErrorLog); ok {
		r1 = rf(logId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// UpdateByID provides a mock function with given fields: ID, request, sqlTransaction, ctx
func (_m *DeliveryOrderUseCaseInterface) UpdateByID(ID int, request *models.DeliveryOrderUpdateByIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderUpdateByIDRequest, *model.ErrorLog) {
	ret := _m.Called(ID, request, sqlTransaction, ctx)

	var r0 *models.DeliveryOrderUpdateByIDRequest
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(int, *models.DeliveryOrderUpdateByIDRequest, *sql.Tx, context.Context) (*models.DeliveryOrderUpdateByIDRequest, *model.ErrorLog)); ok {
		return rf(ID, request, sqlTransaction, ctx)
	}
	if rf, ok := ret.Get(0).(func(int, *models.DeliveryOrderUpdateByIDRequest, *sql.Tx, context.Context) *models.DeliveryOrderUpdateByIDRequest); ok {
		r0 = rf(ID, request, sqlTransaction, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrderUpdateByIDRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(int, *models.DeliveryOrderUpdateByIDRequest, *sql.Tx, context.Context) *model.ErrorLog); ok {
		r1 = rf(ID, request, sqlTransaction, ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// UpdateDODetailByID provides a mock function with given fields: ID, request, sqlTransaction, ctx
func (_m *DeliveryOrderUseCaseInterface) UpdateDODetailByID(ID int, request *models.DeliveryOrderDetailUpdateByIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderDetailUpdateByIDRequest, *model.ErrorLog) {
	ret := _m.Called(ID, request, sqlTransaction, ctx)

	var r0 *models.DeliveryOrderDetailUpdateByIDRequest
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(int, *models.DeliveryOrderDetailUpdateByIDRequest, *sql.Tx, context.Context) (*models.DeliveryOrderDetailUpdateByIDRequest, *model.ErrorLog)); ok {
		return rf(ID, request, sqlTransaction, ctx)
	}
	if rf, ok := ret.Get(0).(func(int, *models.DeliveryOrderDetailUpdateByIDRequest, *sql.Tx, context.Context) *models.DeliveryOrderDetailUpdateByIDRequest); ok {
		r0 = rf(ID, request, sqlTransaction, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrderDetailUpdateByIDRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(int, *models.DeliveryOrderDetailUpdateByIDRequest, *sql.Tx, context.Context) *model.ErrorLog); ok {
		r1 = rf(ID, request, sqlTransaction, ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

// UpdateDoDetailByDeliveryOrderID provides a mock function with given fields: deliveryOrderID, request, sqlTransaction, ctx
func (_m *DeliveryOrderUseCaseInterface) UpdateDoDetailByDeliveryOrderID(deliveryOrderID int, request []*models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest, sqlTransaction *sql.Tx, ctx context.Context) (*models.DeliveryOrderDetails, *model.ErrorLog) {
	ret := _m.Called(deliveryOrderID, request, sqlTransaction, ctx)

	var r0 *models.DeliveryOrderDetails
	var r1 *model.ErrorLog
	if rf, ok := ret.Get(0).(func(int, []*models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest, *sql.Tx, context.Context) (*models.DeliveryOrderDetails, *model.ErrorLog)); ok {
		return rf(deliveryOrderID, request, sqlTransaction, ctx)
	}
	if rf, ok := ret.Get(0).(func(int, []*models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest, *sql.Tx, context.Context) *models.DeliveryOrderDetails); ok {
		r0 = rf(deliveryOrderID, request, sqlTransaction, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeliveryOrderDetails)
		}
	}

	if rf, ok := ret.Get(1).(func(int, []*models.DeliveryOrderDetailUpdateByDeliveryOrderIDRequest, *sql.Tx, context.Context) *model.ErrorLog); ok {
		r1 = rf(deliveryOrderID, request, sqlTransaction, ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ErrorLog)
		}
	}

	return r0, r1
}

type mockConstructorTestingTNewDeliveryOrderUseCaseInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewDeliveryOrderUseCaseInterface creates a new instance of DeliveryOrderUseCaseInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDeliveryOrderUseCaseInterface(t mockConstructorTestingTNewDeliveryOrderUseCaseInterface) *DeliveryOrderUseCaseInterface {
	mock := &DeliveryOrderUseCaseInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
