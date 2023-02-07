package models

import (
	"database/sql"
	"time"
)

func (delicveryOrder *DeliveryOrder) DeliveryOrderStoreRequestMap(request *DeliveryOrderStoreRequest, now time.Time) {
	delicveryOrder.SalesOrderID = request.SalesOrderID
	delicveryOrder.StoreID = request.StoreID
	delicveryOrder.AgentID = request.AgentID
	delicveryOrder.WarehouseID = request.WarehouseID
	delicveryOrder.OrderStatusID = request.OrderStatusID
	delicveryOrder.DoCode = request.DoCode
	delicveryOrder.DoDate = request.DoDate
	delicveryOrder.DoRefCode = NullString{NullString: sql.NullString{String: request.DoRefCode, Valid: true}}
	delicveryOrder.DoRefDate = NullString{NullString: sql.NullString{String: request.DoRefDate, Valid: true}}
	delicveryOrder.DriverName = NullString{NullString: sql.NullString{String: request.DriverName, Valid: true}}
	delicveryOrder.PlatNumber = NullString{NullString: sql.NullString{String: request.PlatNumber, Valid: true}}
	delicveryOrder.IsDoneSyncToEs = "0"
	delicveryOrder.Note = NullString{NullString: sql.NullString{String: request.Note, Valid: true}}
	delicveryOrder.StartDateSyncToEs = &now
	delicveryOrder.EndDateSyncToEs = &now
	delicveryOrder.StartCreatedDate = &now
	delicveryOrder.EndCreatedDate = &now
	delicveryOrder.CreatedBy = request.SalesOrderID
	delicveryOrder.LatestUpdatedBy = int(now.Unix())
	delicveryOrder.CreatedAt = &now
	delicveryOrder.UpdatedAt = &now
	delicveryOrder.DeletedAt = nil
	return
}

func (deliveryOrder *DeliveryOrder) WarehouseChanMap(request *WarehouseChan) {
	deliveryOrder.Warehouse = request.Warehouse
	deliveryOrder.WarehouseName = request.Warehouse.Name
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
	deliveryOrderDetail.OrderStatusID = request.SalesOrderDetail.OrderStatusID
	deliveryOrderDetail.SoDetail = request.SalesOrderDetail
	return
}
