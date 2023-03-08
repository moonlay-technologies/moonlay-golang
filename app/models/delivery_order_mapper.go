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
