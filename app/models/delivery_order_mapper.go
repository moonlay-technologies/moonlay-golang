package models

import (
	"database/sql"
	"order-service/app/models/constants"
	"time"
)

func (deliveryOrder *DeliveryOrder) DeliveryOrderStoreRequestMap(request *DeliveryOrderStoreRequest, now time.Time, user *UserClaims) {
	deliveryOrder.SalesOrderID = request.SalesOrderID
	deliveryOrder.WarehouseID = request.WarehouseID
	deliveryOrder.DoRefCode = NullString{NullString: sql.NullString{String: request.DoRefCode, Valid: true}}
	deliveryOrder.DoRefDate = NullString{NullString: sql.NullString{String: request.DoRefDate, Valid: true}}
	deliveryOrder.DriverName = NullString{NullString: sql.NullString{String: request.DriverName, Valid: true}}
	deliveryOrder.PlatNumber = NullString{NullString: sql.NullString{String: request.PlatNumber, Valid: true}}
	deliveryOrder.Note = NullString{NullString: sql.NullString{String: request.Note, Valid: true}}
	deliveryOrder.CreatedBy = user.UserID
	deliveryOrder.LatestUpdatedBy = user.UserID
	deliveryOrder.IsDoneSyncToEs = "0"
	deliveryOrder.StartDateSyncToEs = &now
	deliveryOrder.EndDateSyncToEs = &now
	deliveryOrder.StartCreatedDate = &now
	deliveryOrder.EndCreatedDate = &now
	deliveryOrder.CreatedAt = &now
	deliveryOrder.UpdatedAt = &now
	deliveryOrder.DeletedAt = nil
	return
}

func (deliveryOrder *DeliveryOrder) DeliveryOrderUploadMap(request *UploadDOField, soId, warehouseId int, now time.Time) {
	deliveryOrder.SalesOrderID = soId
	deliveryOrder.WarehouseID = warehouseId
	deliveryOrder.DoRefCode = NullString{NullString: sql.NullString{String: request.NoSJ, Valid: true}}
	deliveryOrder.DoRefDate = NullString{NullString: sql.NullString{String: request.TanggalSJ, Valid: true}}
	deliveryOrder.DriverName = NullString{NullString: sql.NullString{String: request.NamaSupir, Valid: true}}
	deliveryOrder.PlatNumber = NullString{NullString: sql.NullString{String: request.PlatNo, Valid: true}}
	deliveryOrder.Note = NullString{NullString: sql.NullString{String: request.Catatan, Valid: true}}
	deliveryOrder.CreatedBy = request.IDUser
	deliveryOrder.LatestUpdatedBy = request.IDUser
	deliveryOrder.IsDoneSyncToEs = "0"
	deliveryOrder.StartDateSyncToEs = &now
	deliveryOrder.EndDateSyncToEs = &now
	deliveryOrder.StartCreatedDate = &now
	deliveryOrder.EndCreatedDate = &now
	deliveryOrder.CreatedAt = &now
	deliveryOrder.UpdatedAt = &now
	deliveryOrder.DeletedAt = nil
	return
}

func (deliveryOrder *DeliveryOrder) DeliveryOrderUpdateByIDRequestMap(request *DeliveryOrderUpdateByIDRequest, now time.Time, user *UserClaims) {
	if request.WarehouseID > 0 {
		deliveryOrder.WarehouseID = request.WarehouseID
	}
	if request.DoRefCode != "" {
		deliveryOrder.DoRefCode = NullString{NullString: sql.NullString{String: request.DoRefCode, Valid: true}}
	}
	if request.DoRefDate != "" {
		deliveryOrder.DoRefDate = NullString{NullString: sql.NullString{String: request.DoRefDate, Valid: true}}
	}
	if request.DriverName != "" {
		deliveryOrder.DriverName = NullString{NullString: sql.NullString{String: request.DriverName, Valid: true}}
	}
	if request.PlatNumber != "" {
		deliveryOrder.PlatNumber = NullString{NullString: sql.NullString{String: request.PlatNumber, Valid: true}}
	}
	if request.Note != "" {
		deliveryOrder.Note = NullString{NullString: sql.NullString{String: request.Note, Valid: true}}
	}
	deliveryOrder.IsDoneSyncToEs = "0"
	deliveryOrder.StartDateSyncToEs = &now
	deliveryOrder.EndDateSyncToEs = &now
	deliveryOrder.LatestUpdatedBy = user.UserID
	deliveryOrder.UpdatedAt = &now
	deliveryOrder.DeletedAt = nil
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

func (deliveryOrderDetail *DeliveryOrderDetail) DeliveryOrderDetailUploadMap(soDetailId, qty int, now time.Time) {
	deliveryOrderDetail.SoDetailID = soDetailId
	deliveryOrderDetail.Qty = qty
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

func (deliveryOrderDetail *DeliveryOrderDetail) SalesOrderDetailMap(request *SalesOrderDetail) {
	deliveryOrderDetail.ProductID = request.ProductID
	deliveryOrderDetail.UomID = request.UomID
	deliveryOrderDetail.SoDetail = request
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
func (deliveryOrder *DeliveryOrder) DeliveryOrderUploadSOSJMap(request *UploadSOSJField, now time.Time) {
	deliveryOrder.AgentID = request.IDDistributor
	deliveryOrder.WarehouseID = request.KodeGudang
	deliveryOrder.DoDate = request.TglSuratJalan
	deliveryOrder.DoRefDate = NullString{NullString: sql.NullString{String: request.TglSuratJalan, Valid: true}}
	deliveryOrder.DriverName = NullString{NullString: sql.NullString{String: request.NamaSupir, Valid: true}}
	deliveryOrder.PlatNumber = NullString{NullString: sql.NullString{String: request.PlatNo, Valid: true}}
	deliveryOrder.Note = NullString{NullString: sql.NullString{String: request.Catatan, Valid: true}}
	deliveryOrder.IsDoneSyncToEs = "0"
	deliveryOrder.StartDateSyncToEs = &now
	deliveryOrder.EndDateSyncToEs = &now
	deliveryOrder.StartCreatedDate = &now
	deliveryOrder.EndCreatedDate = &now
	deliveryOrder.CreatedAt = &now
	deliveryOrder.UpdatedAt = &now
	deliveryOrder.DeletedAt = nil
	return
}

func (deliveryOrderDetail *DeliveryOrderDetail) DeliveryOrderDetailUploadSOSJMap(request *UploadSOSJField, now time.Time) {
	deliveryOrderDetail.Qty = request.Qty
	deliveryOrderDetail.Note = NullString{NullString: sql.NullString{String: request.Catatan, Valid: true}}
	deliveryOrderDetail.IsDoneSyncToEs = "0"
	deliveryOrderDetail.StartDateSyncToEs = &now
	deliveryOrderDetail.EndDateSyncToEs = &now
	deliveryOrderDetail.CreatedAt = &now
	deliveryOrderDetail.UpdatedAt = &now
	deliveryOrderDetail.DeletedAt = nil
	return
}

func (deliveryOrderJourney *DeliveryOrderJourneysResponse) DeliveryOrderJourneyResponseMap(request *DeliveryOrderJourney) {
	deliveryOrderJourney.ID = request.ID
	deliveryOrderJourney.DoId = request.DoId
	deliveryOrderJourney.DoCode = request.DoCode
	deliveryOrderJourney.DoDate = request.DoDate
	deliveryOrderJourney.Status = request.Status
	deliveryOrderJourney.Remark = NullString{sql.NullString{String: request.Remark, Valid: true}}
	deliveryOrderJourney.Reason = NullString{sql.NullString{String: request.Reason, Valid: true}}
	deliveryOrderJourney.CreatedAt = request.CreatedAt
	deliveryOrderJourney.UpdatedAt = request.UpdatedAt
	return
}

func (d *DeliveryOrderRequest) DeliveryOrderExportMap(r *DeliveryOrderExportRequest) {
	d.ID = r.ID
	d.PerPage = 0
	d.Page = 0
	d.SortField = r.SortField
	d.SortValue = r.SortValue
	d.GlobalSearchValue = r.GlobalSearchValue
	d.Keyword = ""
	d.AgentID = r.AgentID
	d.AgentName = ""
	d.StoreID = r.StoreID
	d.StoreName = ""
	d.BrandID = r.BrandID
	d.BrandName = ""
	d.ProductID = r.ProductID
	d.OrderSourceID = r.OrderSourceID
	d.OrderStatusID = r.OrderStatusID
	d.SalesOrderID = r.SalesOrderID
	d.SoCode = ""
	d.WarehouseID = r.WarehouseID
	d.WarehouseCode = r.WarehouseCode
	d.DoCode = r.DoCode
	d.DoDate = r.DoDate
	d.DoRefCode = r.DoRefCode
	d.DoRefDate = r.DoRefDate
	d.DoRefferalCode = ""
	d.TotalAmount = 0
	d.TotalTonase = 0
	d.ProductSKU = ""
	d.ProductCode = ""
	d.ProductName = ""
	d.CategoryID = r.CategoryID
	d.SalesmanID = r.SalesmanID
	d.ProvinceID = r.ProvinceID
	d.CityID = r.CityID
	d.DistrictID = r.DistrictID
	d.VillageID = r.VillageID
	d.StoreProvinceID = r.StoreProvinceID
	d.StoreCityID = r.StoreCityID
	d.StoreDistrictID = r.StoreDistrictID
	d.StoreVillageID = r.StoreVillageID
	d.StoreCode = r.StoreCode
	d.StartCreatedAt = r.StartCreatedAt
	d.EndCreatedAt = r.EndCreatedAt
	d.UpdatedAt = r.UpdatedAt
	d.StartDoDate = r.StartDoDate
	d.EndDoDate = r.EndDoDate
}

func (deliveryOrderEventLog *DeliveryOrderEventLogResponse) DeliveryOrderEventLogResponseMap(request *GetDeliveryOrderLog, status string) {
	deliveryOrderEventLog.ID = request.ID
	deliveryOrderEventLog.RequestID = request.RequestID
	deliveryOrderEventLog.DoID = &request.Data.ID
	deliveryOrderEventLog.DoCode = request.DoCode
	deliveryOrderEventLog.Status = status
	deliveryOrderEventLog.Action = request.Action
	deliveryOrderEventLog.CreatedAt = request.CreatedAt
	if request.UpdatedAt == nil {
		deliveryOrderEventLog.UpdatedAt = request.CreatedAt
	} else {
		deliveryOrderEventLog.UpdatedAt = request.UpdatedAt
	}
	return
}

func (dataDOEventLog *DataDOEventLogResponse) DataDOEventLogResponseMap(request *GetDeliveryOrderLog) {
	dataDOEventLog.AgentID = &request.Data.AgentID
	dataDOEventLog.AgentName = &request.Data.AgentName
	dataDOEventLog.SoCode = &request.Data.SalesOrder.SoCode
	dataDOEventLog.DoDate = request.Data.DoDate
	dataDOEventLog.DoRefCode = request.Data.DoRefCode.String
	dataDOEventLog.Note = request.Data.Note
	dataDOEventLog.InternalComment = NullString{sql.NullString{String: request.Data.SalesOrder.InternalComment.String, Valid: true}}
	dataDOEventLog.DriverName = request.Data.DriverName
	dataDOEventLog.PlatNumber = request.Data.PlatNumber
	dataDOEventLog.BrandID = &request.Data.Brand.ID
	dataDOEventLog.BrandName = &request.Data.Brand.Name
	dataDOEventLog.WarehouseCode = request.Data.WarehouseCode
	dataDOEventLog.WarehouseName = request.Data.WarehouseName
	return
}

func (doDetailEventLogResponse *DODetailEventLogResponse) DoDetailEventLogResponse(request *DeliveryOrderDetail) {
	doDetailEventLogResponse.ID = request.ID
	doDetailEventLogResponse.ProductCode = &request.ProductSKU
	doDetailEventLogResponse.ProductName = &request.ProductName
	doDetailEventLogResponse.DeliveryQty = NullInt64{NullInt64: sql.NullInt64{Int64: int64(request.Qty), Valid: true}}
	doDetailEventLogResponse.ProductUnit = NullString{sql.NullString{String: request.Product.UnitMeasurementSmall.String, Valid: true}}
	return
}

func (data *DeliveryOrder) MapToCsvRow() []interface{} {
	deliveryOrderCsv := DeliveryOrderCsvResponse{}
	deliveryOrderCsv.DeliveryOrderMap(data)
	return []interface{}{
		deliveryOrderCsv.DoStatus,
		deliveryOrderCsv.DoDate,
		deliveryOrderCsv.SjNo.String,
		deliveryOrderCsv.DoNo,
		deliveryOrderCsv.OrderNo,
		deliveryOrderCsv.SoDate,
		deliveryOrderCsv.SoNo,
		deliveryOrderCsv.SoSource,
		deliveryOrderCsv.AgentID,
		deliveryOrderCsv.AgentName,
		deliveryOrderCsv.GudangID,
		deliveryOrderCsv.GudangName,
		deliveryOrderCsv.BrandID,
		deliveryOrderCsv.BrandName,
		int(deliveryOrderCsv.KodeSalesman.Int64),
		deliveryOrderCsv.Salesman.String,
		deliveryOrderCsv.KategoryToko.String,
		deliveryOrderCsv.KodeTokoDbo.String,
		deliveryOrderCsv.KodeToko.String,
		deliveryOrderCsv.NamaToko.String,
		deliveryOrderCsv.KodeKecamatan,
		deliveryOrderCsv.Kecamatan.String,
		deliveryOrderCsv.KodeCity,
		deliveryOrderCsv.City.String,
		deliveryOrderCsv.KodeProvince,
		deliveryOrderCsv.Province.String,
		deliveryOrderCsv.DoAmount,
		deliveryOrderCsv.NamaSupir.String,
		deliveryOrderCsv.PlatNo.String,
		deliveryOrderCsv.Catatan.String,
		deliveryOrderCsv.CreatedDate.Format(constants.DATE_FORMAT_EXPORT_CREATED_AT),
		deliveryOrderCsv.UpdatedDate.Format(constants.DATE_FORMAT_EXPORT_CREATED_AT),
		deliveryOrderCsv.UserIDCreated,
		deliveryOrderCsv.UserIDModified,
	}
}

func (d *DeliveryOrderCsvResponse) DeliveryOrderMap(r *DeliveryOrder) {
	d.DoStatus = r.OrderStatusName
	d.DoDate = r.DoDate
	d.SjNo = r.DoRefCode
	d.DoNo = r.DoCode
	d.OrderNo = r.SalesOrder.SoCode
	d.SoDate = r.SalesOrder.SoDate
	d.SoNo = r.SalesOrder.SoCode
	d.SoSource = r.SalesOrder.OrderSourceName
	d.AgentID = r.AgentID
	d.AgentName = r.AgentName
	d.GudangID = r.WarehouseCode
	d.GudangName = r.WarehouseName
	d.BrandID = r.SalesOrder.BrandID
	d.BrandName = r.SalesOrder.BrandName
	d.KodeSalesman = r.SalesOrder.SalesmanID
	d.Salesman = r.SalesOrder.SalesmanName
	if r.SalesOrder.Store != nil {
		d.KategoryToko = r.SalesOrder.Store.StoreCategory
		d.KodeToko = r.SalesOrder.Store.AliasCode
	}
	d.KodeTokoDbo = r.SalesOrder.StoreCode
	d.NamaToko = r.SalesOrder.StoreName
	d.KodeKecamatan = r.SalesOrder.StoreDistrictID
	d.Kecamatan = r.SalesOrder.StoreDistrictName
	d.KodeCity = r.SalesOrder.StoreCityID
	d.City = r.SalesOrder.StoreCityName
	d.KodeProvince = r.SalesOrder.StoreProvinceID
	d.Province = r.SalesOrder.StoreProvinceName
	amount := 0
	for _, v := range r.DeliveryOrderDetails {
		if v.SoDetail == nil {
			for _, x := range r.SalesOrder.SalesOrderDetails {
				if x.ID == v.SoDetailID {
					v.SoDetail = x
				}
			}
		}
		if v.SoDetail != nil {
			amount += int(v.SoDetail.Price) * v.Qty
		}
	}
	d.DoAmount = float64(amount)
	d.NamaSupir = r.DriverName
	d.PlatNo = r.PlatNumber
	d.Catatan = r.Note
	d.CreatedDate = r.CreatedAt
	d.UpdatedDate = r.UpdatedAt
	d.UserIDCreated = r.CreatedBy
	d.UserIDModified = r.LatestUpdatedBy
	return
}
