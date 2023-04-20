// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"

	models "order-service/app/models"
)

// SalesOrderValidatorInterface is an autogenerated mock type for the SalesOrderValidatorInterface type
type SalesOrderValidatorInterface struct {
	mock.Mock
}

// CreateSalesOrderValidator provides a mock function with given fields: insertRequest, ctx
func (_m *SalesOrderValidatorInterface) CreateSalesOrderValidator(insertRequest *models.SalesOrderStoreRequest, ctx *gin.Context) error {
	ret := _m.Called(insertRequest, ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.SalesOrderStoreRequest, *gin.Context) error); ok {
		r0 = rf(insertRequest, ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSalesOrderByIdValidator provides a mock function with given fields: _a0, _a1
func (_m *SalesOrderValidatorInterface) DeleteSalesOrderByIdValidator(_a0 string, _a1 *gin.Context) (int, error) {
	ret := _m.Called(_a0, _a1)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *gin.Context) (int, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(string, *gin.Context) int); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(string, *gin.Context) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteSalesOrderDetailByIdValidator provides a mock function with given fields: _a0, _a1
func (_m *SalesOrderValidatorInterface) DeleteSalesOrderDetailByIdValidator(_a0 string, _a1 *gin.Context) (int, error) {
	ret := _m.Called(_a0, _a1)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *gin.Context) (int, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(string, *gin.Context) int); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(string, *gin.Context) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteSalesOrderDetailBySoIdValidator provides a mock function with given fields: _a0, _a1
func (_m *SalesOrderValidatorInterface) DeleteSalesOrderDetailBySoIdValidator(_a0 string, _a1 *gin.Context) (int, error) {
	ret := _m.Called(_a0, _a1)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *gin.Context) (int, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(string, *gin.Context) int); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(string, *gin.Context) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExportSalesOrderDetailValidator provides a mock function with given fields: ctx
func (_m *SalesOrderValidatorInterface) ExportSalesOrderDetailValidator(ctx *gin.Context) (*models.SalesOrderDetailExportRequest, error) {
	ret := _m.Called(ctx)

	var r0 *models.SalesOrderDetailExportRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(*gin.Context) (*models.SalesOrderDetailExportRequest, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(*gin.Context) *models.SalesOrderDetailExportRequest); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SalesOrderDetailExportRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(*gin.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExportSalesOrderValidator provides a mock function with given fields: ctx
func (_m *SalesOrderValidatorInterface) ExportSalesOrderValidator(ctx *gin.Context) (*models.SalesOrderExportRequest, error) {
	ret := _m.Called(ctx)

	var r0 *models.SalesOrderExportRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(*gin.Context) (*models.SalesOrderExportRequest, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(*gin.Context) *models.SalesOrderExportRequest); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SalesOrderExportRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(*gin.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSOUploadHistoriesValidator provides a mock function with given fields: _a0
func (_m *SalesOrderValidatorInterface) GetSOUploadHistoriesValidator(_a0 *gin.Context) (*models.GetSoUploadHistoriesRequest, error) {
	ret := _m.Called(_a0)

	var r0 *models.GetSoUploadHistoriesRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(*gin.Context) (*models.GetSoUploadHistoriesRequest, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*gin.Context) *models.GetSoUploadHistoriesRequest); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.GetSoUploadHistoriesRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(*gin.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSalesOrderDetailValidator provides a mock function with given fields: _a0
func (_m *SalesOrderValidatorInterface) GetSalesOrderDetailValidator(_a0 *gin.Context) (*models.GetSalesOrderDetailRequest, error) {
	ret := _m.Called(_a0)

	var r0 *models.GetSalesOrderDetailRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(*gin.Context) (*models.GetSalesOrderDetailRequest, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*gin.Context) *models.GetSalesOrderDetailRequest); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.GetSalesOrderDetailRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(*gin.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSalesOrderJourneysValidator provides a mock function with given fields: _a0
func (_m *SalesOrderValidatorInterface) GetSalesOrderJourneysValidator(_a0 *gin.Context) (*models.SalesOrderJourneyRequest, error) {
	ret := _m.Called(_a0)

	var r0 *models.SalesOrderJourneyRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(*gin.Context) (*models.SalesOrderJourneyRequest, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*gin.Context) *models.SalesOrderJourneyRequest); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SalesOrderJourneyRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(*gin.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSalesOrderSyncToKafkaHistoriesValidator provides a mock function with given fields: _a0
func (_m *SalesOrderValidatorInterface) GetSalesOrderSyncToKafkaHistoriesValidator(_a0 *gin.Context) (*models.SalesOrderEventLogRequest, error) {
	ret := _m.Called(_a0)

	var r0 *models.SalesOrderEventLogRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(*gin.Context) (*models.SalesOrderEventLogRequest, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*gin.Context) *models.SalesOrderEventLogRequest); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SalesOrderEventLogRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(*gin.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSalesOrderValidator provides a mock function with given fields: _a0
func (_m *SalesOrderValidatorInterface) GetSalesOrderValidator(_a0 *gin.Context) (*models.SalesOrderRequest, error) {
	ret := _m.Called(_a0)

	var r0 *models.SalesOrderRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(*gin.Context) (*models.SalesOrderRequest, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*gin.Context) *models.SalesOrderRequest); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SalesOrderRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(*gin.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSosjUploadHistoriesValidator provides a mock function with given fields: _a0
func (_m *SalesOrderValidatorInterface) GetSosjUploadHistoriesValidator(_a0 *gin.Context) (*models.GetSosjUploadHistoriesRequest, error) {
	ret := _m.Called(_a0)

	var r0 *models.GetSosjUploadHistoriesRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(*gin.Context) (*models.GetSosjUploadHistoriesRequest, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*gin.Context) *models.GetSosjUploadHistoriesRequest); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.GetSosjUploadHistoriesRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(*gin.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSalesOrderByIdValidator provides a mock function with given fields: updateRequest, ctx
func (_m *SalesOrderValidatorInterface) UpdateSalesOrderByIdValidator(updateRequest *models.SalesOrderUpdateRequest, ctx *gin.Context) error {
	ret := _m.Called(updateRequest, ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.SalesOrderUpdateRequest, *gin.Context) error); ok {
		r0 = rf(updateRequest, ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateSalesOrderDetailByIdValidator provides a mock function with given fields: updateRequest, ctx
func (_m *SalesOrderValidatorInterface) UpdateSalesOrderDetailByIdValidator(updateRequest *models.SalesOrderDetailUpdateByIdRequest, ctx *gin.Context) error {
	ret := _m.Called(updateRequest, ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.SalesOrderDetailUpdateByIdRequest, *gin.Context) error); ok {
		r0 = rf(updateRequest, ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateSalesOrderDetailBySoIdValidator provides a mock function with given fields: updateRequest, ctx
func (_m *SalesOrderValidatorInterface) UpdateSalesOrderDetailBySoIdValidator(updateRequest *models.SalesOrderUpdateRequest, ctx *gin.Context) error {
	ret := _m.Called(updateRequest, ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.SalesOrderUpdateRequest, *gin.Context) error); ok {
		r0 = rf(updateRequest, ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewSalesOrderValidatorInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewSalesOrderValidatorInterface creates a new instance of SalesOrderValidatorInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSalesOrderValidatorInterface(t mockConstructorTestingTNewSalesOrderValidatorInterface) *SalesOrderValidatorInterface {
	mock := &SalesOrderValidatorInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}