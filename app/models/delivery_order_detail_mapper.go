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
