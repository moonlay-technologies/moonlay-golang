package models

import (
	"database/sql"
	"order-service/app/models/constants"
	"strings"
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

func (d *DeliveryOrder) DeliveryOrderUpdateMap(r *DeliveryOrder) {
	defaultDo := &DeliveryOrder{}
	if defaultDo.ID != d.ID {
		d.ID = r.ID
	}
	if defaultDo.SalesOrderID != r.SalesOrderID {
		d.SalesOrderID = r.SalesOrderID
	}
	if defaultDo.SalesOrder != r.SalesOrder {
		d.SalesOrder = r.SalesOrder
	}
	if defaultDo.Brand != r.Brand {
		d.Brand = r.Brand
	}
	if defaultDo.SalesOrderCode != r.SalesOrderCode {
		d.SalesOrderCode = r.SalesOrderCode
	}
	if defaultDo.SalesOrderDate != r.SalesOrderDate {
		d.SalesOrderDate = r.SalesOrderDate
	}
	if defaultDo.Salesman != r.Salesman {
		d.Salesman = r.Salesman
	}
	if defaultDo.WarehouseID != r.WarehouseID {
		d.WarehouseID = r.WarehouseID
	}
	if defaultDo.Warehouse != r.Warehouse {
		d.Warehouse = r.Warehouse
	}
	if defaultDo.WarehouseName != r.WarehouseName {
		d.WarehouseName = r.WarehouseName
	}
	if defaultDo.WarehouseAddress != r.WarehouseAddress {
		d.WarehouseAddress = r.WarehouseAddress
	}
	if defaultDo.WarehouseCode != r.WarehouseCode {
		d.WarehouseCode = r.WarehouseCode
	}
	if defaultDo.WarehouseProvinceID != r.WarehouseProvinceID {
		d.WarehouseProvinceID = r.WarehouseProvinceID
	}
	if defaultDo.WarehouseProvinceName != r.WarehouseProvinceName {
		d.WarehouseProvinceName = r.WarehouseProvinceName
	}
	if defaultDo.WarehouseCityID != r.WarehouseCityID {
		d.WarehouseCityID = r.WarehouseCityID
	}
	if defaultDo.WarehouseCityName != r.WarehouseCityName {
		d.WarehouseCityName = r.WarehouseCityName
	}
	if defaultDo.WarehouseDistrictID != r.WarehouseDistrictID {
		d.WarehouseDistrictID = r.WarehouseDistrictID
	}
	if defaultDo.WarehouseDistrictName != r.WarehouseDistrictName {
		d.WarehouseDistrictName = r.WarehouseDistrictName
	}
	if defaultDo.WarehouseVillageID != r.WarehouseVillageID {
		d.WarehouseVillageID = r.WarehouseVillageID
	}
	if defaultDo.WarehouseVillageName != r.WarehouseVillageName {
		d.WarehouseVillageName = r.WarehouseVillageName
	}
	if defaultDo.OrderStatusID != r.OrderStatusID {
		d.OrderStatusID = r.OrderStatusID
	}
	if defaultDo.OrderStatus != r.OrderStatus {
		d.OrderStatus = r.OrderStatus
	}
	if defaultDo.OrderStatusName != r.OrderStatusName {
		d.OrderStatusName = r.OrderStatusName
	}
	if defaultDo.OrderSourceID != r.OrderSourceID {
		d.OrderSourceID = r.OrderSourceID
	}
	if defaultDo.OrderSource != r.OrderSource {
		d.OrderSource = r.OrderSource
	}
	if defaultDo.OrderSourceName != r.OrderSourceName {
		d.OrderSourceName = r.OrderSourceName
	}
	if defaultDo.AgentID != r.AgentID {
		d.AgentID = r.AgentID
	}
	if defaultDo.AgentName != r.AgentName {
		d.AgentName = r.AgentName
	}
	if defaultDo.Agent != r.Agent {
		d.Agent = r.Agent
	}
	if defaultDo.StoreID != r.StoreID {
		d.StoreID = r.StoreID
	}
	if defaultDo.Store != r.Store {
		d.Store = r.Store
	}
	if defaultDo.DoCode != r.DoCode {
		d.DoCode = r.DoCode
	}
	if defaultDo.DoDate != r.DoDate {
		d.DoDate = r.DoDate
	}
	if defaultDo.DoRefCode != r.DoRefCode {
		d.DoRefCode = r.DoRefCode
	}
	if defaultDo.DoRefDate != r.DoRefDate {
		d.DoRefDate = r.DoRefDate
	}
	if defaultDo.DriverName != r.DriverName {
		d.DriverName = r.DriverName
	}
	if defaultDo.PlatNumber != r.PlatNumber {
		d.PlatNumber = r.PlatNumber
	}
	if defaultDo.Note != r.Note {
		d.Note = r.Note
	}
	if defaultDo.IsDoneSyncToEs != r.IsDoneSyncToEs {
		d.IsDoneSyncToEs = r.IsDoneSyncToEs
	}
	if defaultDo.StartDateSyncToEs != r.StartDateSyncToEs {
		d.StartDateSyncToEs = r.StartDateSyncToEs
	}
	if defaultDo.EndDateSyncToEs != r.EndDateSyncToEs {
		d.EndDateSyncToEs = r.EndDateSyncToEs
	}
	if defaultDo.CreatedBy != r.CreatedBy {
		d.CreatedBy = r.CreatedBy
	}
	if defaultDo.LatestUpdatedBy != r.LatestUpdatedBy {
		d.LatestUpdatedBy = r.LatestUpdatedBy
	}
	if defaultDo.StartCreatedDate != r.StartCreatedDate {
		d.StartCreatedDate = r.StartCreatedDate
	}
	if defaultDo.EndCreatedDate != r.EndCreatedDate {
		d.EndCreatedDate = r.EndCreatedDate
	}
	if defaultDo.CreatedAt != r.CreatedAt {
		d.CreatedAt = r.CreatedAt
	}
	if defaultDo.UpdatedAt != r.UpdatedAt {
		d.UpdatedAt = r.UpdatedAt
	}
	if defaultDo.DeletedAt != r.DeletedAt {
		d.DeletedAt = r.DeletedAt
	}
	for _, v := range r.DeliveryOrderDetails {
		for k, y := range d.DeliveryOrderDetails {
			if y.ID == v.ID {
				d.DeliveryOrderDetails[k].DoDetailUpdateMap(v)
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
	deliveryOrderJourney.DoId = &request.DoId
	deliveryOrderJourney.DoCode = &request.DoCode
	deliveryOrderJourney.DoDate = &request.DoDate
	deliveryOrderJourney.Status = &request.Status
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
	deliveryOrderEventLog.DoCode = request.Data.DoCode
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
	dataDOEventLog.AgentID = &request.Data.SalesOrder.AgentID
	dataDOEventLog.AgentName = &request.Data.SalesOrder.AgentName.String
	dataDOEventLog.SoCode = &request.Data.SalesOrder.SoCode
	dataDOEventLog.DoDate = &request.Data.DoDate
	dataDOEventLog.DoRefCode = &request.Data.DoRefCode.String
	dataDOEventLog.Note = NullString{sql.NullString{String: request.Data.Note.String, Valid: true}}
	dataDOEventLog.InternalComment = NullString{sql.NullString{String: request.Data.SalesOrder.InternalComment.String, Valid: true}}
	dataDOEventLog.DriverName = NullString{sql.NullString{String: request.Data.DriverName.String, Valid: true}}
	dataDOEventLog.PlatNumber = NullString{sql.NullString{String: request.Data.PlatNumber.String, Valid: true}}
	dataDOEventLog.BrandID = &request.Data.SalesOrder.BrandID
	dataDOEventLog.BrandName = &request.Data.SalesOrder.BrandName
	dataDOEventLog.WarehouseCode = &request.Data.WarehouseCode
	dataDOEventLog.WarehouseName = &request.Data.WarehouseName
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

func (deliveryOrderLog *GetDeliveryOrderLog) DeliveryOrderLogBinaryMap(request *DeliveryOrderLog, data DeliveryOrder) {
	deliveryOrderLog.Data = &data
	deliveryOrderLog.Action = request.Action
	deliveryOrderLog.Status = request.Status
	deliveryOrderLog.CreatedAt = request.CreatedAt
	return
}

func (data *DeliveryOrder) MapToCsvRow() []interface{} {
	deliveryOrderCsv := DeliveryOrderCsvResponse{}
	deliveryOrderCsv.DeliveryOrderMap(data)
	if len(data.DoDate) > 9 {
		data.DoDate = data.DoDate[0:10]
	}
	return []interface{}{
		deliveryOrderCsv.DoStatus,
		strings.ReplaceAll(deliveryOrderCsv.DoDate, "T00:00:00Z", ""),
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
	d.OrderNo = r.SalesOrder.SoRefCode.String
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
