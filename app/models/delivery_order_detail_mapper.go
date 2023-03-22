package models

import (
	"database/sql"
	"strconv"
)

func (deliveryOrderDetail *DeliveryOrderDetailOpenSearchDetailResponse) DeliveryOrderDetailOpenSearchResponseMap(request *DeliveryOrderDetail) {
	deliveryOrderDetail.SoDetailID = request.SoDetailID
	deliveryOrderDetail.Qty = NullInt64{sql.NullInt64{Int64: int64(request.Qty), Valid: true}}
	return
}

func (deliveryOrderDetail *DeliveryOrderDetailsOpenSearchResponse) DeliveryOrderDetailsByDoIDOpenSearchResponseMap(request *DeliveryOrderDetail) {
	deliveryOrderDetail.ID = request.ID
	deliveryOrderDetail.DeliveryOrderID = request.DeliveryOrderID
	deliveryOrderDetail.SoDetailID = request.SoDetailID
	deliveryOrderDetail.Qty = NullInt64{sql.NullInt64{Int64: int64(request.Qty), Valid: true}}
	return
}

func (deliveryOrderDetail *DeliveryOrderDetailsOpenSearchResponse) DeliveryOrderDetailsOpenSearchResponseMap(request *DeliveryOrderDetailOpenSearch) {
	deliveryOrderDetail.ID = request.ID
	deliveryOrderDetail.DeliveryOrderID = request.DeliveryOrderID
	deliveryOrderDetail.SoDetailID = request.SoDetailID
	deliveryOrderDetail.Qty = NullInt64{sql.NullInt64{Int64: int64(request.Qty), Valid: true}}
	return
}

func (deliveryOrder *DeliveryOrder) AgentMap(request *Agent) {
	deliveryOrder.Agent = request
	deliveryOrder.AgentID = request.ID
	deliveryOrder.AgentName = request.Name
}

func (deliveryOrderDetailResponse *DeliveryOrderDetailStoreResponse) DeliveryOrderDetailMap(request *DeliveryOrderDetail) {
	deliveryOrderDetailResponse.DeliveryOrderID = request.DeliveryOrderID
	deliveryOrderDetailResponse.OrderStatusID = request.OrderStatus.ID
	deliveryOrderDetailResponse.SoDetailID = request.SoDetailID
	deliveryOrderDetailResponse.ProductSku = request.ProductSKU
	deliveryOrderDetailResponse.ProductName = request.ProductName
	deliveryOrderDetailResponse.SalesOrderQty = request.SoDetail.Qty
	deliveryOrderDetailResponse.SentQty = request.SoDetail.SentQty
	deliveryOrderDetailResponse.ResidualQty = request.SoDetail.ResidualQty
	deliveryOrderDetailResponse.UomCode = request.Uom.Code.String
	deliveryOrderDetailResponse.Price = int(request.SoDetail.Price)
	deliveryOrderDetailResponse.Qty = request.Qty
	deliveryOrderDetailResponse.Note = request.Note.String
	return
}

func (deliveryOrderResponse *DeliveryOrderStoreResponse) DeliveryOrderMap(deliveryOrder *DeliveryOrder) {
	storeProvinceID, _ := strconv.Atoi(deliveryOrder.Store.ProvinceID.String)
	storeCityID, _ := strconv.Atoi(deliveryOrder.Store.CityID.String)
	deliveryOrderResponse.ID = deliveryOrder.ID
	deliveryOrderResponse.SalesOrderID = deliveryOrder.SalesOrderID
	deliveryOrderResponse.SalesOrderOrderStatusID = deliveryOrder.SalesOrder.OrderStatusID
	deliveryOrderResponse.SalesOrderOrderStatusName = deliveryOrder.SalesOrder.OrderStatusName
	deliveryOrderResponse.SalesOrderSoCode = deliveryOrder.SalesOrder.SoCode
	deliveryOrderResponse.SalesOrderSoDate = deliveryOrder.SalesOrder.SoDate
	deliveryOrderResponse.SalesOrderReferralCode = deliveryOrder.SalesOrder.SoRefCode.String
	deliveryOrderResponse.SalesOrderNote = deliveryOrder.SalesOrder.Note.String
	deliveryOrderResponse.SalesOrderInternalComment = deliveryOrder.SalesOrder.InternalComment.String
	deliveryOrderResponse.StoreName = deliveryOrder.Store.Name.String
	deliveryOrderResponse.StoreProvinceID = storeProvinceID
	deliveryOrderResponse.StoreProvince = deliveryOrder.Store.ProvinceName.String
	deliveryOrderResponse.StoreCityID = storeCityID
	deliveryOrderResponse.StoreCity = deliveryOrder.Store.CityName.String
	deliveryOrderResponse.TotalAmount = int(deliveryOrder.SalesOrder.TotalAmount)
	deliveryOrderResponse.WarehouseID = deliveryOrder.WarehouseID
	deliveryOrderResponse.WarehouseName = deliveryOrder.Warehouse.Name
	deliveryOrderResponse.WarehouseAddress = deliveryOrder.Warehouse.Address.String
	deliveryOrderResponse.OrderSourceID = deliveryOrder.OrderSourceID
	deliveryOrderResponse.OrderSourceName = deliveryOrder.OrderSource.SourceName
	deliveryOrderResponse.OrderStatusID = deliveryOrder.OrderStatusID
	deliveryOrderResponse.OrderStatusName = deliveryOrder.OrderStatus.Name
	deliveryOrderResponse.DoCode = deliveryOrder.DoCode
	deliveryOrderResponse.DoDate = deliveryOrder.DoDate
	deliveryOrderResponse.DoRefCode = deliveryOrder.DoRefCode.String
	deliveryOrderResponse.DoRefDate = deliveryOrder.DoRefDate.String
	deliveryOrderResponse.DriverName = deliveryOrder.DriverName.String
	deliveryOrderResponse.PlatNumber = deliveryOrder.PlatNumber.String
	deliveryOrderResponse.Note = deliveryOrder.Note.String
	deliveryOrderResponse.InternalComment = deliveryOrder.SalesOrder.InternalComment.String
	if deliveryOrder.Salesman != nil {
		deliveryOrderResponse.SalesmanID = deliveryOrder.Salesman.ID
		deliveryOrderResponse.SalesmanName = deliveryOrder.Salesman.Name
	}
}
func (d *DeliveryOrderDetailLogData) DoDetailMap(r *DeliveryOrder, rd *DeliveryOrderDetail) {
	d.ID = rd.ID
	d.AgentID = r.AgentID
	d.AgentName = r.AgentName
	d.DoRefCode = r.DoRefCode
	d.DoDate = r.DoDate
	d.DoNumber = r.DoRefCode.String
	d.DoDetailCode = rd.DoDetailCode
	d.SoDetailID = rd.SoDetailID
	d.Note = r.Note
	d.InternalNote = rd.Note
	d.DriverName = r.DriverName
	d.PlatNumber = r.PlatNumber
	d.BrandID = rd.BrandID
	d.BrandName = rd.BrandName
	d.ProductID = rd.ProductID
	d.ProductName = rd.ProductName
	d.DeliveryQty = rd.Qty
	d.UomCode = rd.UomCode
	d.WarehouseID = r.WarehouseID
	d.WarehouseName = r.WarehouseName
}
func (d *DeliveryOrderDetailOpenSearch) DoDetailMap(r *DeliveryOrder, rd *DeliveryOrderDetail) {
	d.ID = rd.ID
	d.DeliveryOrderID = r.ID
	d.DoCode = r.DoCode
	d.DoDate = r.DoDate
	d.DoRefCode = r.DoRefCode.String
	d.DoRefDate = r.DoRefDate.String
	d.DriverName = r.DriverName
	d.PlatNumber = r.PlatNumber
	d.SalesOrderID = r.SalesOrderID
	d.SoCode = NullString{sql.NullString{String: r.SalesOrder.SoCode, Valid: true}}
	d.SoDate = NullString{sql.NullString{String: r.SalesOrder.SoDate, Valid: true}}
	d.SoRefDate = r.SalesOrder.SoRefDate
	d.SoDetailID = rd.SoDetail.ID
	d.SoDetailCode = rd.SoDetail.SoDetailCode
	d.SoDetail = rd.SoDetail
	d.AgentID = r.AgentID
	d.Agent = r.Agent
	d.StoreID = r.StoreID
	d.Store = r.Store
	d.WarehouseID = r.WarehouseID
	d.WarehouseCode = r.WarehouseCode
	d.WarehouseName = r.WarehouseName
	if r.Salesman != nil {
		d.SalesmanID = r.Salesman.ID
		d.SalesmanName = r.Salesman.Name
		d.Salesman = r.Salesman
	}
	d.BrandID = rd.Brand.ID
	d.BrandName = rd.Brand.Name
	d.Brand = rd.Brand
	d.ProductID = rd.ProductID
	d.Product = rd.Product
	d.UomID = rd.UomID
	d.Uom = rd.Uom
	d.DoDetailCode = rd.DoDetailCode
	d.OrderSourceID = r.OrderSourceID
	d.OrderSourceName = r.OrderSourceName
	d.OrderSource = r.OrderSource
	d.OrderStatusID = r.OrderStatusID
	d.OrderStatusName = r.OrderStatusName
	d.OrderStatus = r.OrderStatus
	d.Qty = rd.Qty
	d.Note = rd.Note
	d.CreatedAt = rd.CreatedAt
	d.UpdatedAt = rd.UpdatedAt
	d.DeletedAt = rd.DeletedAt
}
