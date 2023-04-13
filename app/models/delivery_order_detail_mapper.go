package models

import (
	"database/sql"
	"order-service/app/models/constants"
	"strconv"
)

func (deliveryOrderDetail *DeliveryOrderDetailOpenSearchDetailResponse) DeliveryOrderDetailOpenSearchResponseMap(request *DeliveryOrderDetail) {
	deliveryOrderDetail.SoDetailID = request.SoDetailID
	deliveryOrderDetail.Qty = NullInt64{NullInt64: sql.NullInt64{Int64: int64(request.Qty), Valid: true}}
	return
}

func (deliveryOrderDetail *DeliveryOrderDetailsOpenSearchResponse) DeliveryOrderDetailsByDoIDOpenSearchResponseMap(request *DeliveryOrderDetail) {
	deliveryOrderDetail.ID = request.ID
	deliveryOrderDetail.DeliveryOrderID = request.DeliveryOrderID
	deliveryOrderDetail.SoDetailID = request.SoDetailID
	deliveryOrderDetail.Qty = NullInt64{NullInt64: sql.NullInt64{Int64: int64(request.Qty), Valid: true}}
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
	defaultDoDetail := &DeliveryOrderDetail{}
	defaultDo := &DeliveryOrder{}
	if defaultDoDetail.ID != rd.ID {
		d.ID = rd.ID
	}
	if defaultDoDetail.DeliveryOrderID != r.ID {
		d.DeliveryOrderID = r.ID
	}
	if defaultDo.DoCode != r.DoCode {
		d.DoCode = r.DoCode
	}
	if defaultDo.DoDate != r.DoDate {
		d.DoDate = r.DoDate
	}
	if defaultDo.DoRefCode != r.DoRefCode {
		d.DoRefCode = r.DoRefCode.String
	}
	if defaultDo.DoRefDate != r.DoRefDate {
		d.DoRefDate = r.DoRefDate.String
	}
	if defaultDo.DriverName != r.DriverName {
		d.DriverName = r.DriverName
	}
	if defaultDo.PlatNumber != r.PlatNumber {
		d.PlatNumber = r.PlatNumber
	}
	if defaultDo.SalesOrderID != r.SalesOrderID {
		d.SalesOrderID = r.SalesOrderID
	}
	if r.SalesOrder != nil {
		if "" != r.SalesOrder.SoCode {
			d.SoCode = NullString{sql.NullString{String: r.SalesOrder.SoCode, Valid: true}}
		}
		if "" != r.SalesOrder.SoDate {
			d.SoDate = NullString{sql.NullString{String: r.SalesOrder.SoDate, Valid: true}}
		}
		if "" != r.SalesOrder.SoRefDate.String {
			d.SoRefDate = r.SalesOrder.SoRefDate
		}
		if 0 != r.SalesOrder.StoreID {
			d.StoreID = r.StoreID
		}
		if nil != r.SalesOrder.Store {
			d.Store = r.Store
		}
	}

	if defaultDoDetail.SoDetailID != rd.SoDetailID {
		d.SoDetailID = rd.SoDetailID
	}
	if rd.SoDetail != nil {
		d.SoDetailCode = rd.SoDetail.SoDetailCode
		d.SoDetail = rd.SoDetail
	}
	if defaultDo.AgentID != r.AgentID {
		d.AgentID = r.AgentID
	}
	if defaultDo.Agent != r.Agent {
		d.Agent = r.Agent
	}
	if defaultDo.WarehouseID != r.WarehouseID {
		d.WarehouseID = r.WarehouseID
	}
	if defaultDo.WarehouseCode != r.WarehouseCode {
		d.WarehouseCode = r.WarehouseCode
	}
	if defaultDo.WarehouseName != r.WarehouseName {
		d.WarehouseName = r.WarehouseName
	}
	if r.Salesman != nil {
		d.SalesmanID = r.Salesman.ID
		d.SalesmanName = r.Salesman.Name
		d.Salesman = r.Salesman
	}
	if defaultDoDetail.BrandID != rd.BrandID {
		d.BrandID = rd.BrandID
	}
	if r.Brand != nil {
		d.BrandName = rd.Brand.Name
		d.Brand = rd.Brand
	}
	if defaultDoDetail.ProductID != rd.ProductID {
		d.ProductID = rd.ProductID
	}
	if rd.Product != nil {
		d.Product = rd.Product
	}
	if defaultDoDetail.UomID != rd.UomID {
		d.UomID = rd.UomID
	}
	if defaultDoDetail.Uom != rd.Uom {
		d.Uom = rd.Uom
	}
	if defaultDoDetail.DoDetailCode != rd.DoDetailCode {
		d.DoDetailCode = rd.DoDetailCode
	}
	if defaultDo.OrderSourceID != r.OrderSourceID {
		d.OrderSourceID = r.OrderSourceID
	}
	if defaultDo.OrderSourceName != r.OrderSourceName {
		d.OrderSourceName = r.OrderSourceName
	}
	if defaultDo.OrderSource != r.OrderSource {
		d.OrderSource = r.OrderSource
	}
	if defaultDo.OrderStatusID != r.OrderStatusID {
		d.OrderStatusID = r.OrderStatusID
	}
	if defaultDoDetail.OrderStatusName != r.OrderStatusName {
		d.OrderStatusName = r.OrderStatusName
	}
	if defaultDoDetail.OrderStatus != r.OrderStatus {
		d.OrderStatus = r.OrderStatus
	}
	if defaultDoDetail.Qty != rd.Qty {
		d.Qty = rd.Qty
	}
	if defaultDoDetail.Note != rd.Note {
		d.Note = rd.Note
	}
	if defaultDoDetail.CreatedAt != rd.CreatedAt {
		d.CreatedAt = rd.CreatedAt
	}
	if defaultDoDetail.UpdatedAt != rd.UpdatedAt {
		d.UpdatedAt = rd.UpdatedAt
	}
	if defaultDoDetail.DeletedAt != rd.DeletedAt {
		d.DeletedAt = rd.DeletedAt
	}
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
	if d.SoDetail == nil {
		d.SoDetail = &SalesOrderDetail{}
	}
	if d.Agent == nil {
		d.Agent = &Agent{}
	}
	if d.Product == nil {
		d.Product = &Product{}
	}
	if d.Uom == nil {
		d.Uom = &Uom{}
	}
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
		d.OrderStatusName,
		d.DoDate,
		d.DoRefCode,
		d.DoCode,
		d.SoDate.String,
		d.SoCode.String,
		dd.SalesOrder.OrderSourceName,
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
		d.Uom.Code.String,
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
