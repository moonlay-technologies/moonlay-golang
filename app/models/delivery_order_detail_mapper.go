package models

import (
	"database/sql"
	"order-service/app/models/constants"
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
	d.SoDetailID = rd.SoDetailID
	if rd.SoDetail != nil {
		d.SoDetailCode = rd.SoDetail.SoDetailCode
		d.SoDetail = rd.SoDetail
	}
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
	d.BrandID = rd.BrandID
	if r.Brand != nil {
		d.BrandName = rd.Brand.Name
		d.Brand = rd.Brand
	}
	d.ProductID = rd.ProductID
	if rd.Product != nil {
		d.Product = rd.Product
	}
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

func (d *DeliveryOrderDetailOpenSearchRequest) DeliveryOrderDetailExportMap(r *DeliveryOrderDetailExportRequest) {
	d.ID = r.ID
	d.DoDetailID = r.DoDetailID
	d.PerPage = 0
	d.Page = 0
	d.SortField = r.SortField
	d.SortValue = r.SortValue
	d.GlobalSearchValue = r.GlobalSearchValue
	d.AgentID = r.AgentID
	d.AgentName = ""
	d.StoreID = r.StoreID
	d.StoreName = ""
	d.BrandID = r.BrandID
	d.BrandName = ""
	d.ProductID = r.ProductID
	d.OrderSourceID = 0
	d.OrderStatusID = r.OrderStatusID
	d.SalesOrderID = r.SalesOrderID
	d.SoCode = ""
	d.WarehouseID = 0
	d.DoCode = r.DoCode
	d.DoDate = r.DoDate
	d.DoRefCode = r.DoRefCode
	d.DoRefDate = r.DoRefDate
	d.DoRefferalCode = ""
	d.TotalAmount = 0
	d.TotalTonase = 0
	d.CategoryID = r.CategoryID
	d.SalesmanID = 0
	d.ProvinceID = 0
	d.CityID = 0
	d.DistrictID = 0
	d.VillageID = 0
	d.StoreProvinceID = r.ProvinceID
	d.StoreCityID = r.CityID
	d.StoreDistrictID = r.DistrictID
	d.StoreVillageID = 0
	d.StartCreatedAt = r.StartCreatedAt
	d.EndCreatedAt = r.EndCreatedAt
	d.UpdatedAt = r.UpdatedAt
	d.StartDoDate = r.StartDoDate
	d.EndDoDate = r.EndDoDate
}
func (d *DeliveryOrderDetailOpenSearch) MapToCsvRow(dd *DeliveryOrder) []interface{} {
	store := Store{}
	if d.Store != nil {
		store.StoreCategory = d.Store.StoreCategory
		store.StoreCode = d.Store.StoreCode
		store.AliasCode = d.Store.AliasCode
		store.Name = d.Store.Name
		store.DistrictID = d.Store.DistrictID
		store.DistrictName = d.Store.DistrictName
		store.CityID = d.Store.CityID
		store.CityName = d.Store.CityName
		store.ProvinceID = d.Store.ProvinceID
		store.ProvinceName = d.Store.ProvinceName
	} else {
		store.StoreCategory = NullString{NullString: sql.NullString{String: "", Valid: true}}
		store.StoreCode = NullString{NullString: sql.NullString{String: "", Valid: true}}
		store.AliasCode = NullString{NullString: sql.NullString{String: "", Valid: true}}
		store.Name = NullString{NullString: sql.NullString{String: "", Valid: true}}
		store.DistrictID = NullString{NullString: sql.NullString{String: "", Valid: true}}
		store.DistrictName = NullString{NullString: sql.NullString{String: "", Valid: true}}
		store.CityID = NullString{NullString: sql.NullString{String: "", Valid: true}}
		store.CityName = NullString{NullString: sql.NullString{String: "", Valid: true}}
		store.ProvinceID = NullString{NullString: sql.NullString{String: "", Valid: true}}
		store.ProvinceName = NullString{NullString: sql.NullString{String: "", Valid: true}}
	}
	return []interface{}{
		d.OrderStatusID,
		d.DoDate,
		d.DoRefCode,
		d.DoCode,
		d.SoDate.String,
		d.SoCode.String,
		d.OrderSourceID,
		d.AgentID,
		d.Agent.Name,
		d.WarehouseCode,
		d.WarehouseName,
		d.BrandID,
		d.BrandName,
		d.SalesmanID,
		d.SalesmanName,
		store.StoreCategory.String,
		store.StoreCode.String,
		store.AliasCode.String,
		store.Name.String,
		store.DistrictID.String,
		store.DistrictName.String,
		store.CityID.String,
		store.CityName.String,
		store.ProvinceID.String,
		store.ProvinceName.String,
		d.BrandID,
		d.BrandName,
		d.SoDetail.FirstCategoryId,
		// "*d.SoDetail.FirstCategoryName",
		nil,
		d.SoDetail.LastCategoryId,
		// "*d.SoDetail.LastCategoryName",
		nil,
		d.SoDetail.ProductSKU,
		d.Product.ProductName.String,
		d.UomName,
		d.SoDetail.Price,
		d.SoDetail.Qty,
		d.SoDetail.ResidualQty,
		d.Qty,
		int(d.SoDetail.Price) * d.Qty,
		d.CreatedAt.Format(constants.DATE_FORMAT_EXPORT_CREATED_AT),
		d.UpdatedAt.Format(constants.DATE_FORMAT_EXPORT_CREATED_AT),
		dd.CreatedBy,
		dd.LatestUpdatedBy,
	}
}
