package models

import (
	"database/sql"
	"time"
)

func (v *SalesOrderDetail) SalesOrderDetailStoreRequestMap(soDetail *SalesOrderDetailStoreRequest, now time.Time) {
	v.ProductID = soDetail.ProductID
	v.UomID = soDetail.UomID
	v.OrderStatusID = soDetail.OrderStatusId
	v.Qty = soDetail.Qty
	v.SentQty = soDetail.SentQty
	v.ResidualQty = soDetail.ResidualQty
	v.Price = soDetail.Price
	v.Note = NullString{NullString: sql.NullString{String: soDetail.Note, Valid: true}}
	v.IsDoneSyncToEs = "0"
	v.StartDateSyncToEs = &now
	v.CreatedAt = &now
	v.UpdatedAt = &now
	return
}

func (v *SalesOrderDetail) SalesOrderDetailUpdateRequestMap(soDetail *SalesOrderDetailUpdateRequest, now time.Time) {
	v.ID = soDetail.ID
	v.ProductID = soDetail.ProductID
	v.UomID = soDetail.UomID
	v.Qty = soDetail.Qty
	v.SentQty = soDetail.SentQty
	v.ResidualQty = soDetail.ResidualQty
	v.Price = soDetail.Price
	v.Note = NullString{NullString: sql.NullString{String: soDetail.Note, Valid: true}}
	v.UpdatedAt = &now
	return
}

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

func (salesOrderDetail *SalesOrderDetailOpenSearch) SalesOrderDetailMap(request *SalesOrderDetail) {
	salesOrderDetail.ID = request.ID
	salesOrderDetail.SoDetailCode = request.SoDetailCode
	salesOrderDetail.Qty = request.Qty
	salesOrderDetail.SentQty = request.SentQty
	salesOrderDetail.ResidualQty = request.ResidualQty
	salesOrderDetail.Price = request.Price
	salesOrderDetail.Note = request.Note
	salesOrderDetail.OrderStatusID = request.OrderStatusID
	salesOrderDetail.UomID = request.UomID
	salesOrderDetail.IsDoneSyncToEs = request.IsDoneSyncToEs
	salesOrderDetail.StartDateSyncToEs = request.StartDateSyncToEs
	salesOrderDetail.EndDateSyncToEs = request.EndDateSyncToEs
	salesOrderDetail.Subtotal = request.Subtotal
	salesOrderDetail.CreatedAt = request.CreatedAt
	salesOrderDetail.UpdatedAt = request.UpdatedAt
	salesOrderDetail.DeletedAt = request.DeletedAt
	salesOrderDetail.ProductID = request.ProductID
	salesOrderDetail.Product = request.Product
	salesOrderDetail.Uom = request.Uom
	return
}

func (salesOrderDetail *SalesOrderDetailOpenSearch) SalesOrderMap(request *SalesOrder) {
	salesOrderDetail.SalesOrderID = request.ID
	salesOrderDetail.SoCode = request.SoCode
	salesOrderDetail.SoDate = request.SoDate
	salesOrderDetail.SoRefCode = request.ReferralCode.String
	salesOrderDetail.SoRefDate = request.SoRefDate
	salesOrderDetail.SoReferralCode = request.ReferralCode
	salesOrderDetail.AgentId = request.AgentID
	salesOrderDetail.Agent = request.Agent
	salesOrderDetail.StoreID = request.StoreID
	salesOrderDetail.Store = request.Store
	salesOrderDetail.BrandID = request.BrandID
	salesOrderDetail.BrandName = request.BrandName
	salesOrderDetail.Brand = request.Brand
	salesOrderDetail.UserID = request.UserID
	salesOrderDetail.User = request.User
	salesOrderDetail.SalesmanID = request.SalesmanID
	salesOrderDetail.Salesman = request.Salesman
	salesOrderDetail.OrderSourceID = request.OrderSourceID
	salesOrderDetail.OrderSource = request.OrderSource
	salesOrderDetail.OrderStatus = request.OrderStatus
	salesOrderDetail.GLat = request.GLat
	salesOrderDetail.GLong = request.GLong
}
