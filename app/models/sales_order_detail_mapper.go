package models

import (
	"database/sql"
	"order-service/app/models/constants"
	"strconv"
	"time"
)

func (v *SalesOrderDetail) SalesOrderDetailStoreRequestMap(soDetail *SalesOrderDetailStoreRequest, now time.Time) {
	v.ProductID = soDetail.ProductID
	v.UomID = soDetail.UomID
	v.Qty = soDetail.Qty
	v.SentQty = 0
	v.ResidualQty = soDetail.Qty
	v.Price = soDetail.Price
	v.Note = NullString{NullString: sql.NullString{String: soDetail.Note, Valid: true}}
	v.IsDoneSyncToEs = "0"
	v.StartDateSyncToEs = &now
	v.CreatedAt = &now
	v.UpdatedAt = &now
	return
}

func (v *SalesOrderDetailStoreResponse) SalesOrderDetailStoreResponseMap(soDetail *SalesOrderDetail) {
	v.ID = soDetail.ID
	v.SalesOrderId = soDetail.SalesOrderID
	v.ProductID = soDetail.ProductID
	v.UomID = soDetail.UomID
	v.OrderStatusId = soDetail.OrderStatusID
	v.OrderStatusName = soDetail.OrderStatusName
	v.SoDetailCode = soDetail.SoDetailCode
	v.Qty = soDetail.Qty
	v.SentQty = NullInt64{NullInt64: sql.NullInt64{Int64: int64(soDetail.SentQty), Valid: true}}
	v.ResidualQty = NullInt64{NullInt64: sql.NullInt64{Int64: int64(soDetail.ResidualQty), Valid: true}}
	v.Price = soDetail.Price
	v.Note = soDetail.Note.String
	return
}

func (v *SalesOrderDetailStoreResponse) UpdateSoDetailResponseMap(soDetail *SalesOrderDetail) {
	v.ID = soDetail.ID
	v.SalesOrderId = soDetail.SalesOrderID
	v.ProductID = soDetail.ProductID
	v.UomID = soDetail.UomID
	v.OrderStatusId = soDetail.OrderStatusID
	v.SoDetailCode = soDetail.SoDetailCode
	v.Qty = soDetail.Qty
	v.Price = soDetail.Price
	v.SentQty = NullInt64{NullInt64: sql.NullInt64{Int64: int64(soDetail.SentQty), Valid: true}}
	v.ResidualQty = NullInt64{NullInt64: sql.NullInt64{Int64: int64(soDetail.ResidualQty), Valid: true}}
	v.Note = soDetail.Note.String
	v.CreatedAt = soDetail.CreatedAt
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

func (salesOrderDetail *SalesOrderDetailOpenSearchResponse) SalesOrderDetailOpenSearchMap(request *SalesOrderDetailOpenSearch) {
	salesOrderDetail.ID = request.ID
	salesOrderDetail.SalesOrderID = request.SoID
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
	salesOrderDetail.SalesmanID = int(request.SalesmanID.Int64)
	salesOrderDetail.Salesman = request.Salesman
	salesOrderDetail.OrderSourceID = request.OrderSourceID
	salesOrderDetail.OrderSource = request.OrderSource
	salesOrderDetail.OrderStatus = request.OrderStatus
	salesOrderDetail.GLat = request.GLat
	salesOrderDetail.GLong = request.GLong
}

func (salesOrderDetail *SalesOrderDetailOpenSearch) SalesOrderDetailOpenSearchMap(requestSalesOrder *SalesOrder, requestSalesOrderDetail *SalesOrderDetail) {
	salesOrderDetail.ID = requestSalesOrderDetail.ID
	salesOrderDetail.SoDetailCode = requestSalesOrderDetail.SoDetailCode
	salesOrderDetail.SoID = requestSalesOrder.ID
	salesOrderDetail.SoCode = requestSalesOrder.SoCode
	salesOrderDetail.SoDate = requestSalesOrder.SoDate
	salesOrderDetail.SoRefCode = requestSalesOrder.SoRefCode.String
	salesOrderDetail.SoRefDate = requestSalesOrder.SoRefDate
	salesOrderDetail.ReferralCode = requestSalesOrder.ReferralCode

	salesOrderDetail.AgentId = requestSalesOrder.AgentID
	salesOrderDetail.AgentName = requestSalesOrder.AgentName.String
	salesOrderDetail.AgentProvinceID, _ = strconv.Atoi(requestSalesOrder.Agent.ProvinceID.String)
	salesOrderDetail.AgentProvinceName = requestSalesOrder.AgentProvinceName.String
	salesOrderDetail.AgentCityID, _ = strconv.Atoi(requestSalesOrder.Agent.CityID.String)
	salesOrderDetail.AgentCityName = requestSalesOrder.AgentCityName.String
	salesOrderDetail.AgentDistrictID, _ = strconv.Atoi(requestSalesOrder.Agent.DistrictID.String)
	salesOrderDetail.AgentDistrictName = requestSalesOrder.AgentDistrictName.String
	salesOrderDetail.AgentVillageID, _ = strconv.Atoi(requestSalesOrder.Agent.VillageID.String)
	salesOrderDetail.AgentVillageName = requestSalesOrder.AgentVillageName.String
	salesOrderDetail.AgentPhone = requestSalesOrder.AgentPhone.String
	salesOrderDetail.AgentAddress = requestSalesOrder.AgentAddress.String

	salesOrderDetail.StoreID = requestSalesOrder.StoreID
	salesOrderDetail.StoreCode = requestSalesOrder.StoreCode.String
	salesOrderDetail.StoreName = requestSalesOrder.StoreName.String
	salesOrderDetail.StoreProvinceID, _ = strconv.Atoi(requestSalesOrder.Store.ProvinceID.String)
	salesOrderDetail.StoreProvinceName = requestSalesOrder.StoreProvinceName.String
	salesOrderDetail.StoreCityID, _ = strconv.Atoi(requestSalesOrder.Store.CityID.String)
	salesOrderDetail.StoreCityName = requestSalesOrder.StoreCityName.String
	salesOrderDetail.StoreDistrictID, _ = strconv.Atoi(requestSalesOrder.Store.DistrictID.String)
	salesOrderDetail.StoreDistrictName = requestSalesOrder.StoreDistrictName.String
	salesOrderDetail.StoreVillageID, _ = strconv.Atoi(requestSalesOrder.Store.VillageID.String)
	salesOrderDetail.StoreVillageName = requestSalesOrder.StoreVillageName.String
	salesOrderDetail.StoreAddress = requestSalesOrder.StoreAddress.String
	salesOrderDetail.StorePhone = requestSalesOrder.StorePhone.String
	salesOrderDetail.StoreMainMobilePhone = requestSalesOrder.StoreMainMobilePhone.String

	salesOrderDetail.BrandID = requestSalesOrder.BrandID
	salesOrderDetail.BrandName = requestSalesOrder.BrandName

	salesOrderDetail.UserID = requestSalesOrder.UserID
	salesOrderDetail.UserFirstName = requestSalesOrder.UserFirstName.String
	salesOrderDetail.UserLastName = requestSalesOrder.UserLastName.String
	salesOrderDetail.UserRoleID, _ = strconv.Atoi(requestSalesOrder.User.RoleID.String)
	salesOrderDetail.UserEmail = requestSalesOrder.UserEmail.String

	salesOrderDetail.SalesmanID = int(requestSalesOrder.SalesmanID.Int64)
	salesOrderDetail.SalesmanName = requestSalesOrder.SalesmanName.String
	salesOrderDetail.SalesmanEmail = requestSalesOrder.SalesmanEmail.String

	salesOrderDetail.OrderSourceID = requestSalesOrder.OrderSourceID
	salesOrderDetail.OrderSourceName = requestSalesOrder.OrderSourceName
	salesOrderDetail.OrderStatusID = requestSalesOrder.OrderStatusID
	salesOrderDetail.OrderStatusName = requestSalesOrder.OrderStatusName
	salesOrderDetail.OrderStatus = requestSalesOrderDetail.OrderStatus

	salesOrderDetail.GLat = requestSalesOrder.GLat
	salesOrderDetail.GLong = requestSalesOrder.GLong

	salesOrderDetail.ProductID = requestSalesOrderDetail.ProductID
	salesOrderDetail.ProductSKU = requestSalesOrderDetail.Product.Sku.String
	salesOrderDetail.ProductName = requestSalesOrderDetail.Product.ProductName.String
	salesOrderDetail.ProductDescription = requestSalesOrderDetail.Product.Description.String
	salesOrderDetail.CategoryID = requestSalesOrderDetail.Product.CategoryID
	salesOrderDetail.Product = requestSalesOrderDetail.Product

	salesOrderDetail.UomID = requestSalesOrderDetail.UomID
	salesOrderDetail.UomCode = requestSalesOrderDetail.Uom.Code.String
	salesOrderDetail.UomName = requestSalesOrderDetail.Uom.Name.String
	salesOrderDetail.Uom = requestSalesOrderDetail.Uom

	salesOrderDetail.Qty = requestSalesOrderDetail.Qty
	salesOrderDetail.SentQty = requestSalesOrderDetail.SentQty
	salesOrderDetail.ResidualQty = requestSalesOrderDetail.ResidualQty
	salesOrderDetail.Price = requestSalesOrderDetail.Price
	salesOrderDetail.Note = requestSalesOrderDetail.Note

	salesOrderDetail.FirstCategoryId = requestSalesOrderDetail.FirstCategoryId
	salesOrderDetail.FirstCategoryName = requestSalesOrderDetail.FirstCategoryName
	salesOrderDetail.LastCategoryId = requestSalesOrderDetail.LastCategoryId
	salesOrderDetail.LastCategoryName = requestSalesOrderDetail.LastCategoryName

	salesOrderDetail.CreatedBy = requestSalesOrderDetail.CreatedBy
	salesOrderDetail.UpdatedBy = requestSalesOrderDetail.LatestUpdatedBy
	salesOrderDetail.DeletedBy = requestSalesOrderDetail.DeletedBy

	salesOrderDetail.CreatedAt = requestSalesOrderDetail.CreatedAt
	salesOrderDetail.UpdatedAt = requestSalesOrderDetail.UpdatedAt
	salesOrderDetail.DeletedAt = requestSalesOrderDetail.DeletedAt
}

func (v *SalesOrderDetailStoreResponse) UpdateSalesOrderDetailByIdResponseMap(request *SalesOrderDetail, brandId int) {
	v.BrandID = brandId
	v.ProductID = request.ProductID
	v.UomID = request.UomID
	v.Note = request.Note.String
	v.Price = request.Price
	v.Qty = request.Qty
	v.SentQty = NullInt64{NullInt64: sql.NullInt64{Int64: int64(request.SentQty), Valid: true}}
	v.ResidualQty = NullInt64{NullInt64: sql.NullInt64{Int64: int64(request.ResidualQty), Valid: true}}
	return
}

func (v *SalesOrderDetail) SalesOrderDetailUploadSOSJMap(soDetail *UploadSOSJField, now time.Time) {
	v.UomID = soDetail.Unit
	v.Qty = soDetail.Qty
	v.SentQty = soDetail.Qty
	v.ResidualQty = 0
	v.Note = NullString{NullString: sql.NullString{String: soDetail.Catatan, Valid: true}}
	v.IsDoneSyncToEs = "0"
	v.CreatedAt = &now
	v.StartDateSyncToEs = &now
	v.UpdatedAt = &now
}

func (v *SalesOrderDetail) SalesOrderDetailUploadSOMap(soDetail *UploadSOField, now time.Time) {
	v.Qty = soDetail.QTYOrder
	v.SentQty = 0
	v.ResidualQty = soDetail.QTYOrder
	v.Note = NullString{NullString: sql.NullString{String: soDetail.CatatanOrder, Valid: true}}
	v.IsDoneSyncToEs = "0"
	v.CreatedAt = &now
	v.StartDateSyncToEs = &now
	v.UpdatedAt = &now
}

func (d *GetSalesOrderDetailRequest) SalesOrderDetailExportMap(r *SalesOrderDetailExportRequest) {
	d.ID = r.ID
	d.SoID = r.SoID
	d.PerPage = 0
	d.Page = 0
	d.SortField = r.SortField
	d.SortValue = r.SortValue
	d.GlobalSearchValue = r.GlobalSearchValue
	d.AgentID = r.AgentID
	d.StoreID = r.StoreID
	d.BrandID = r.BrandID
	d.OrderSourceID = r.OrderSourceID
	d.OrderStatusID = r.OrderStatusID
	d.StartSoDate = r.StartSoDate
	d.EndSoDate = r.EndSoDate
	d.ProductID = r.ProductID
	d.CategoryID = r.CategoryID
	d.SalesmanID = r.SalesmanID
	d.ProvinceID = r.ProvinceID
	d.CityID = r.CityID
	d.DistrictID = r.DistrictID
	d.VillageID = r.VillageID
	d.StartCreatedAt = r.StartCreatedAt
	d.EndCreatedAt = r.EndCreatedAt

}

func (data *SalesOrderDetailOpenSearch) MapToCsvRow(journey *SalesOrderJourneys, s *SalesOrder, d []*DeliveryOrder) []interface{} {
	cancelReason := ""
	if data.OrderStatusID == 10 && journey != nil {
		cancelReason = journey.Reason
	}
	doQty := 0
	for _, v := range d {
		for _, x := range v.DeliveryOrderDetails {
			if x.SoDetailID == data.ID {
				doQty += x.Qty
			}
		}
	}
	if data.Store == nil {
		data.Store = &Store{}
	}
	updatedAt := ""
	if data.UpdatedAt != nil {
		updatedAt = data.UpdatedAt.Format(constants.DATE_FORMAT_EXPORT_CREATED_AT)
	}
	if len(data.SoDate) > 9 {
		data.SoDate = data.SoDate[0:10]
	}
	soRefDate := data.SoRefDate.String
	if len(soRefDate) > 9 {
		soRefDate = soRefDate[0:10]
	}
	return []interface{}{
		s.OrderStatusName,
		s.OrderSourceName,
		data.ReferralCode.String,
		data.SoRefCode,
		data.SoCode,
		data.SoDate,
		data.AgentId,
		data.AgentName,
		data.SalesmanID,
		data.SalesmanName,
		data.StoreID,
		data.StoreCode,
		data.Store.AliasCode.String,
		data.StoreName,
		data.StoreDistrictID,
		data.StoreDistrictName,
		data.StoreCityID,
		data.StoreCityName,
		data.StoreProvinceID,
		data.StoreProvinceName,
		data.BrandID,
		data.BrandName,
		data.FirstCategoryId,
		// *data.FirstCategoryName,
		nil,
		data.LastCategoryId,
		// *data.LastCategoryName,
		nil,
		data.ProductSKU,
		data.ProductName,
		data.UomCode,
		data.Price,
		data.Qty,
		float64(data.Qty) * data.Price,
		doQty,
		float64(doQty) * data.Price,
		data.OrderStatusName,
		cancelReason,
		data.CreatedAt.Format(constants.DATE_FORMAT_EXPORT_CREATED_AT),
		updatedAt,
		data.CreatedBy,
		s.LatestUpdatedBy,
		s.InternalComment.String,
	}
}
