package models

import "database/sql"

func (salesOrderDetail *SalesOrderDetailOpenSearchResponse) SalesOrderDetailOpenSearchResponseMap(request *SalesOrderDetail) {
	salesOrderDetail.ID = request.ID
	salesOrderDetail.SalesOrderID = request.SalesOrderID
	salesOrderDetail.ProductID = request.ProductID
	salesOrderDetail.UomID = request.UomID
	salesOrderDetail.OrderStatusID = request.OrderStatusID
	salesOrderDetail.SoDetailCode = request.SoDetailCode
	salesOrderDetail.Qty = NullInt64{NullInt64: sql.NullInt64{Int64: int64(request.Qty), Valid: true}}
	salesOrderDetail.SentQty = NullInt64{NullInt64: sql.NullInt64{Int64: int64(request.SentQty), Valid: true}}
	salesOrderDetail.ResidualQty = NullInt64{NullInt64: sql.NullInt64{Int64: int64(request.ResidualQty), Valid: true}}
	salesOrderDetail.Price = request.Price
	salesOrderDetail.Note = request.Note
	salesOrderDetail.CreatedAt = request.CreatedAt

	salesOrderDetail.Product = &ProductOpenSearchResponse{}
	salesOrderDetail.Product.ProductOpenSearchResponseMap(request.Product)
	salesOrderDetail.Uom = &UomOpenSearchResponse{}
	salesOrderDetail.Uom.UomOpenSearchResponseMap(request.Uom)
	salesOrderDetail.OrderStatus = &OrderStatusOpenSearchResponse{}
	salesOrderDetail.OrderStatus.OrderStatusOpenSearchResponseMap(request.OrderStatus)
	return
}
