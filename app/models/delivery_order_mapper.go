package models

import (
	"database/sql"
	"time"
)

func (delicveryOrder *DeliveryOrder) DeliveryOrderStoreRequestMap(request *DeliveryOrderStoreRequest, now time.Time) {
	delicveryOrder.SalesOrderID = request.SalesOrderID
	delicveryOrder.WarehouseID = request.WarehouseID
	delicveryOrder.DoRefCode = NullString{NullString: sql.NullString{String: request.DoRefCode, Valid: true}}
	delicveryOrder.DoRefDate = NullString{NullString: sql.NullString{String: request.DoRefDate, Valid: true}}
	delicveryOrder.DriverName = NullString{NullString: sql.NullString{String: request.DriverName, Valid: true}}
	delicveryOrder.PlatNumber = NullString{NullString: sql.NullString{String: request.PlatNumber, Valid: true}}
	delicveryOrder.Note = NullString{NullString: sql.NullString{String: request.Note, Valid: true}}
	delicveryOrder.IsDoneSyncToEs = "0"
	delicveryOrder.StartDateSyncToEs = &now
	delicveryOrder.EndDateSyncToEs = &now
	delicveryOrder.StartCreatedDate = &now
	delicveryOrder.EndCreatedDate = &now
	delicveryOrder.LatestUpdatedBy = &now
	delicveryOrder.CreatedAt = &now
	delicveryOrder.UpdatedAt = &now
	delicveryOrder.DeletedAt = nil
	return
}

func (delicveryOrder *DeliveryOrder) DeliveryOrderUpdateByIDRequestMap(request *DeliveryOrderUpdateByIDRequest, now time.Time) {
	if request.WarehouseID > 0 {
		delicveryOrder.WarehouseID = request.WarehouseID
	}
	if request.DoRefCode != "" {
		delicveryOrder.DoRefCode = NullString{NullString: sql.NullString{String: request.DoRefCode, Valid: true}}
	}
	if request.DoRefDate != "" {
		delicveryOrder.DoRefDate = NullString{NullString: sql.NullString{String: request.DoRefDate, Valid: true}}
	}
	if request.DriverName != "" {
		delicveryOrder.DriverName = NullString{NullString: sql.NullString{String: request.DriverName, Valid: true}}
	}
	if request.PlatNumber != "" {
		delicveryOrder.PlatNumber = NullString{NullString: sql.NullString{String: request.PlatNumber, Valid: true}}
	}
	if request.Note != "" {
		delicveryOrder.Note = NullString{NullString: sql.NullString{String: request.Note, Valid: true}}
	}
	delicveryOrder.IsDoneSyncToEs = "0"
	delicveryOrder.StartDateSyncToEs = &now
	delicveryOrder.EndDateSyncToEs = &now
	delicveryOrder.LatestUpdatedBy = &now
	delicveryOrder.UpdatedAt = &now
	delicveryOrder.DeletedAt = nil
	return
}

func (deliveryOrder *DeliveryOrder) WarehouseChanMap(request *WarehouseChan) {
	deliveryOrder.Warehouse = request.Warehouse
	deliveryOrder.WarehouseID = request.Warehouse.ID
	deliveryOrder.WarehouseName = request.Warehouse.Name
	deliveryOrder.WarehouseAddress = request.Warehouse.Address
	deliveryOrder.WarehouseCode = request.Warehouse.Code
	deliveryOrder.WarehouseProvinceName = request.Warehouse.ProvinceName
	deliveryOrder.WarehouseCityName = request.Warehouse.CityName
	deliveryOrder.WarehouseDistrictName = request.Warehouse.DistrictName
	deliveryOrder.WarehouseVillageName = request.Warehouse.VillageName
	return
}

func (deliveryOrderDetail *DeliveryOrderDetail) DeliveryOrderDetailStoreRequestMap(request *DeliveryOrderDetailStoreRequest, now time.Time) {
	deliveryOrderDetail.SoDetailID = request.SoDetailID
	deliveryOrderDetail.Qty = request.Qty
	deliveryOrderDetail.IsDoneSyncToEs = "0"
	deliveryOrderDetail.StartDateSyncToEs = &now
	deliveryOrderDetail.EndDateSyncToEs = &now
	deliveryOrderDetail.CreatedAt = &now
	deliveryOrderDetail.UpdatedAt = &now
	deliveryOrderDetail.DeletedAt = nil
	return
}

func (deliveryOrderDetail *DeliveryOrderDetail) ProductChanMap(request *ProductChan) {
	deliveryOrderDetail.Product = request.Product
	deliveryOrderDetail.ProductSKU = request.Product.Sku.String
	deliveryOrderDetail.ProductName = request.Product.ProductName.String
	return
}

func (deliveryOrderDetail *DeliveryOrderDetail) SalesOrderDetailChanMap(request *SalesOrderDetailChan) {
	deliveryOrderDetail.ProductID = request.SalesOrderDetail.ProductID
	deliveryOrderDetail.UomID = request.SalesOrderDetail.UomID
	deliveryOrderDetail.SoDetail = request.SalesOrderDetail
	return
}

func (deliveryOrder *DeliveryOrderOpenSearchResponse) DeliveryOrderOpenSearchResponseMap(request *DeliveryOrder) {
	deliveryOrder.ID = request.ID
	deliveryOrder.SalesOrderID = request.SalesOrderID
	deliveryOrder.WarehouseID = request.WarehouseID
	deliveryOrder.OrderSourceID = request.OrderSourceID
	deliveryOrder.AgentName = request.AgentName
	deliveryOrder.AgentID = request.AgentID
	deliveryOrder.StoreID = request.StoreID
	deliveryOrder.DoCode = request.DoCode
	deliveryOrder.DoDate = request.DoDate
	deliveryOrder.DoRefCode = request.DoRefCode
	deliveryOrder.DoRefDate = request.DoRefDate
	deliveryOrder.DriverName = request.DriverName
	deliveryOrder.PlatNumber = request.PlatNumber
	deliveryOrder.Note = request.Note
	return
}

func (deliveryOrder *DeliveryOrder) DeliveryOrderUpdateMap(request *DeliveryOrder) {
	defaultDo := &DeliveryOrder{}
	if defaultDo.ID != deliveryOrder.ID {
		deliveryOrder.ID = request.ID
	}
	if defaultDo.SalesOrderID != deliveryOrder.SalesOrderID {
		deliveryOrder.SalesOrderID = request.SalesOrderID
	}
	if defaultDo.SalesOrder != deliveryOrder.SalesOrder {
		deliveryOrder.SalesOrder = request.SalesOrder
	}
	if defaultDo.Brand != deliveryOrder.Brand {
		deliveryOrder.Brand = request.Brand
	}
	if defaultDo.SalesOrderCode != deliveryOrder.SalesOrderCode {
		deliveryOrder.SalesOrderCode = request.SalesOrderCode
	}
	if defaultDo.SalesOrderDate != deliveryOrder.SalesOrderDate {
		deliveryOrder.SalesOrderDate = request.SalesOrderDate
	}
	if defaultDo.Salesman != deliveryOrder.Salesman {
		deliveryOrder.Salesman = request.Salesman
	}
	if defaultDo.WarehouseID != deliveryOrder.WarehouseID {
		deliveryOrder.WarehouseID = request.WarehouseID
	}
	if defaultDo.Warehouse != deliveryOrder.Warehouse {
		deliveryOrder.Warehouse = request.Warehouse
	}
	if defaultDo.WarehouseName != deliveryOrder.WarehouseName {
		deliveryOrder.WarehouseName = request.WarehouseName
	}
	if defaultDo.WarehouseAddress != deliveryOrder.WarehouseAddress {
		deliveryOrder.WarehouseAddress = request.WarehouseAddress
	}
	if defaultDo.WarehouseCode != deliveryOrder.WarehouseCode {
		deliveryOrder.WarehouseCode = request.WarehouseCode
	}
	if defaultDo.WarehouseProvinceID != deliveryOrder.WarehouseProvinceID {
		deliveryOrder.WarehouseProvinceID = request.WarehouseProvinceID
	}
	if defaultDo.WarehouseProvinceName != deliveryOrder.WarehouseProvinceName {
		deliveryOrder.WarehouseProvinceName = request.WarehouseProvinceName
	}
	if defaultDo.WarehouseCityID != deliveryOrder.WarehouseCityID {
		deliveryOrder.WarehouseCityID = request.WarehouseCityID
	}
	if defaultDo.WarehouseCityName != deliveryOrder.WarehouseCityName {
		deliveryOrder.WarehouseCityName = request.WarehouseCityName
	}
	if defaultDo.WarehouseDistrictID != deliveryOrder.WarehouseDistrictID {
		deliveryOrder.WarehouseDistrictID = request.WarehouseDistrictID
	}
	if defaultDo.WarehouseDistrictName != deliveryOrder.WarehouseDistrictName {
		deliveryOrder.WarehouseDistrictName = request.WarehouseDistrictName
	}
	if defaultDo.WarehouseVillageID != deliveryOrder.WarehouseVillageID {
		deliveryOrder.WarehouseVillageID = request.WarehouseVillageID
	}
	if defaultDo.WarehouseVillageName != deliveryOrder.WarehouseVillageName {
		deliveryOrder.WarehouseVillageName = request.WarehouseVillageName
	}
	if defaultDo.OrderStatusID != deliveryOrder.OrderStatusID {
		deliveryOrder.OrderStatusID = request.OrderStatusID
	}
	if defaultDo.OrderStatus != deliveryOrder.OrderStatus {
		deliveryOrder.OrderStatus = request.OrderStatus
	}
	if defaultDo.OrderStatusName != deliveryOrder.OrderStatusName {
		deliveryOrder.OrderStatusName = request.OrderStatusName
	}
	if defaultDo.OrderSourceID != deliveryOrder.OrderSourceID {
		deliveryOrder.OrderSourceID = request.OrderSourceID
	}
	if defaultDo.OrderSource != deliveryOrder.OrderSource {
		deliveryOrder.OrderSource = request.OrderSource
	}
	if defaultDo.OrderSourceName != deliveryOrder.OrderSourceName {
		deliveryOrder.OrderSourceName = request.OrderSourceName
	}
	if defaultDo.AgentID != deliveryOrder.AgentID {
		deliveryOrder.AgentID = request.AgentID
	}
	if defaultDo.AgentName != deliveryOrder.AgentName {
		deliveryOrder.AgentName = request.AgentName
	}
	if defaultDo.Agent != deliveryOrder.Agent {
		deliveryOrder.Agent = request.Agent
	}
	if defaultDo.StoreID != deliveryOrder.StoreID {
		deliveryOrder.StoreID = request.StoreID
	}
	if defaultDo.Store != deliveryOrder.Store {
		deliveryOrder.Store = request.Store
	}
	if defaultDo.DoCode != deliveryOrder.DoCode {
		deliveryOrder.DoCode = request.DoCode
	}
	if defaultDo.DoDate != deliveryOrder.DoDate {
		deliveryOrder.DoDate = request.DoDate
	}
	if defaultDo.DoRefCode != deliveryOrder.DoRefCode {
		deliveryOrder.DoRefCode = request.DoRefCode
	}
	if defaultDo.DoRefDate != deliveryOrder.DoRefDate {
		deliveryOrder.DoRefDate = request.DoRefDate
	}
	if defaultDo.DriverName != deliveryOrder.DriverName {
		deliveryOrder.DriverName = request.DriverName
	}
	if defaultDo.PlatNumber != deliveryOrder.PlatNumber {
		deliveryOrder.PlatNumber = request.PlatNumber
	}
	if defaultDo.Note != deliveryOrder.Note {
		deliveryOrder.Note = request.Note
	}
	if defaultDo.IsDoneSyncToEs != deliveryOrder.IsDoneSyncToEs {
		deliveryOrder.IsDoneSyncToEs = request.IsDoneSyncToEs
	}
	if defaultDo.StartDateSyncToEs != deliveryOrder.StartDateSyncToEs {
		deliveryOrder.StartDateSyncToEs = request.StartDateSyncToEs
	}
	if defaultDo.EndDateSyncToEs != deliveryOrder.EndDateSyncToEs {
		deliveryOrder.EndDateSyncToEs = request.EndDateSyncToEs
	}
	if defaultDo.CreatedBy != deliveryOrder.CreatedBy {
		deliveryOrder.CreatedBy = request.CreatedBy
	}
	if defaultDo.LatestUpdatedBy != deliveryOrder.LatestUpdatedBy {
		deliveryOrder.LatestUpdatedBy = request.LatestUpdatedBy
	}
	if defaultDo.StartCreatedDate != deliveryOrder.StartCreatedDate {
		deliveryOrder.StartCreatedDate = request.StartCreatedDate
	}
	if defaultDo.EndCreatedDate != deliveryOrder.EndCreatedDate {
		deliveryOrder.EndCreatedDate = request.EndCreatedDate
	}
	if defaultDo.CreatedAt != deliveryOrder.CreatedAt {
		deliveryOrder.CreatedAt = request.CreatedAt
	}
	if defaultDo.UpdatedAt != deliveryOrder.UpdatedAt {
		deliveryOrder.UpdatedAt = request.UpdatedAt
	}
	if defaultDo.DeletedAt != deliveryOrder.DeletedAt {
		deliveryOrder.DeletedAt = request.DeletedAt
	}
	for _, v := range request.DeliveryOrderDetails {
		for k, y := range deliveryOrder.DeliveryOrderDetails {
			if y.ID == v.ID {
				deliveryOrder.DeliveryOrderDetails[k] = v
			}
		}
	}
	return
}
