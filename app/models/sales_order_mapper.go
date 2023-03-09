package models

import (
	"database/sql"
	"strconv"
	"time"
)

func (salesOrder *SalesOrder) SalesOrderRequestMap(request *SalesOrderStoreRequest, now time.Time) {
	salesOrder.CartID = request.CartID
	salesOrder.AgentID = request.AgentID
	salesOrder.StoreID = request.StoreID
	salesOrder.UserID = request.UserID
	salesOrder.VisitationID = request.VisitationID
	salesOrder.OrderSourceID = request.OrderSourceID
	salesOrder.GLat = NullFloat64{NullFloat64: sql.NullFloat64{Float64: request.GLat, Valid: true}}
	salesOrder.GLong = NullFloat64{NullFloat64: sql.NullFloat64{Float64: request.GLong, Valid: true}}
	salesOrder.SoRefCode = NullString{NullString: sql.NullString{String: request.SoRefCode, Valid: true}}
	salesOrder.SoDate = request.SoDate
	salesOrder.SoRefDate = NullString{NullString: sql.NullString{String: request.SoRefDate, Valid: true}}
	salesOrder.Note = NullString{NullString: sql.NullString{String: request.Note, Valid: true}}
	salesOrder.InternalComment = NullString{NullString: sql.NullString{String: request.InternalComment, Valid: true}}
	salesOrder.TotalAmount = request.TotalAmount
	salesOrder.TotalTonase = request.TotalTonase
	salesOrder.DeviceId = NullString{NullString: sql.NullString{String: request.DeviceId, Valid: true}}
	salesOrder.ReferralCode = NullString{NullString: sql.NullString{String: request.ReferralCode, Valid: true}}
	salesOrder.IsDoneSyncToEs = "0"
	salesOrder.CreatedAt = &now
	salesOrder.StartDateSyncToEs = &now
	salesOrder.StartCreatedDate = &now
	salesOrder.CreatedBy = request.UserID
	return
}

func (salesOrder *SalesOrder) SalesOrderOpenSearchChanMap(request *SalesOrderChan) {
	x := &SalesOrder{}
	if salesOrder.ID == x.ID {
		salesOrder.ID = request.SalesOrder.ID
	}
	if salesOrder.CartID == x.CartID {
		salesOrder.CartID = request.SalesOrder.CartID
	}
	if salesOrder.AgentID == x.AgentID {
		salesOrder.AgentID = request.SalesOrder.AgentID
	}
	if salesOrder.AgentName == x.AgentName {
		salesOrder.AgentName = request.SalesOrder.AgentName
	}
	if salesOrder.AgentEmail == x.AgentEmail {
		salesOrder.AgentEmail = request.SalesOrder.AgentEmail
	}
	if salesOrder.AgentProvinceID == x.AgentProvinceID {
		salesOrder.AgentProvinceID = request.SalesOrder.AgentProvinceID
	}
	if salesOrder.AgentProvinceName == x.AgentProvinceName {
		salesOrder.AgentProvinceName = request.SalesOrder.AgentProvinceName
	}
	if salesOrder.AgentCityID == x.AgentCityID {
		salesOrder.AgentCityID = request.SalesOrder.AgentCityID
	}
	if salesOrder.AgentCityName == x.AgentCityName {
		salesOrder.AgentCityName = request.SalesOrder.AgentCityName
	}
	if salesOrder.AgentDistrictID == x.AgentDistrictID {
		salesOrder.AgentDistrictID = request.SalesOrder.AgentDistrictID
	}
	if salesOrder.AgentDistrictName == x.AgentDistrictName {
		salesOrder.AgentDistrictName = request.SalesOrder.AgentDistrictName
	}
	if salesOrder.AgentVillageID == x.AgentVillageID {
		salesOrder.AgentVillageID = request.SalesOrder.AgentVillageID
	}
	if salesOrder.AgentVillageName == x.AgentVillageName {
		salesOrder.AgentVillageName = request.SalesOrder.AgentVillageName
	}
	if salesOrder.AgentAddress == x.AgentAddress {
		salesOrder.AgentAddress = request.SalesOrder.AgentAddress
	}
	if salesOrder.AgentPhone == x.AgentPhone {
		salesOrder.AgentPhone = request.SalesOrder.AgentPhone
	}
	if salesOrder.AgentMainMobilePhone == x.AgentMainMobilePhone {
		salesOrder.AgentMainMobilePhone = request.SalesOrder.AgentMainMobilePhone
	}
	if salesOrder.Agent == x.Agent {
		salesOrder.Agent = request.SalesOrder.Agent
	}
	if salesOrder.StoreID == x.StoreID {
		salesOrder.StoreID = request.SalesOrder.StoreID
	}
	if salesOrder.StoreName == x.StoreName {
		salesOrder.StoreName = request.SalesOrder.StoreName
	}
	if salesOrder.StoreCode == x.StoreCode {
		salesOrder.StoreCode = request.SalesOrder.StoreCode
	}
	if salesOrder.StoreEmail == x.StoreEmail {
		salesOrder.StoreEmail = request.SalesOrder.StoreEmail
	}
	if salesOrder.StoreProvinceID == x.StoreProvinceID {
		salesOrder.StoreProvinceID = request.SalesOrder.StoreProvinceID
	}
	if salesOrder.StoreProvinceName == x.StoreProvinceName {
		salesOrder.StoreProvinceName = request.SalesOrder.StoreProvinceName
	}
	if salesOrder.StoreCityID == x.StoreCityID {
		salesOrder.StoreCityID = request.SalesOrder.StoreCityID
	}
	if salesOrder.StoreCityName == x.StoreCityName {
		salesOrder.StoreCityName = request.SalesOrder.StoreCityName
	}
	if salesOrder.StoreDistrictID == x.StoreDistrictID {
		salesOrder.StoreDistrictID = request.SalesOrder.StoreDistrictID
	}
	if salesOrder.StoreDistrictName == x.StoreDistrictName {
		salesOrder.StoreDistrictName = request.SalesOrder.StoreDistrictName
	}
	if salesOrder.StoreVillageID == x.StoreVillageID {
		salesOrder.StoreVillageID = request.SalesOrder.StoreVillageID
	}
	if salesOrder.StoreVillageName == x.StoreVillageName {
		salesOrder.StoreVillageName = request.SalesOrder.StoreVillageName
	}
	if salesOrder.StoreAddress == x.StoreAddress {
		salesOrder.StoreAddress = request.SalesOrder.StoreAddress
	}
	if salesOrder.StorePhone == x.StorePhone {
		salesOrder.StorePhone = request.SalesOrder.StorePhone
	}
	if salesOrder.StoreMainMobilePhone == x.StoreMainMobilePhone {
		salesOrder.StoreMainMobilePhone = request.SalesOrder.StoreMainMobilePhone
	}
	if salesOrder.Store == x.Store {
		salesOrder.Store = request.SalesOrder.Store
	}
	if salesOrder.BrandID == x.BrandID {
		salesOrder.BrandID = request.SalesOrder.BrandID
	}
	if salesOrder.BrandName == x.BrandName {
		salesOrder.BrandName = request.SalesOrder.BrandName
	}
	if salesOrder.Brand == x.Brand {
		salesOrder.Brand = request.SalesOrder.Brand
	}
	if salesOrder.UserID == x.UserID {
		salesOrder.UserID = request.SalesOrder.UserID
	}
	if salesOrder.UserFirstName == x.UserFirstName {
		salesOrder.UserFirstName = request.SalesOrder.UserFirstName
	}
	if salesOrder.UserLastName == x.UserLastName {
		salesOrder.UserLastName = request.SalesOrder.UserLastName
	}
	if salesOrder.UserRoleID == x.UserRoleID {
		salesOrder.UserRoleID = request.SalesOrder.UserRoleID
	}
	if salesOrder.UserEmail == x.UserEmail {
		salesOrder.UserEmail = request.SalesOrder.UserEmail
	}
	if salesOrder.User == x.User {
		salesOrder.User = request.SalesOrder.User
	}
	if salesOrder.Salesman == x.Salesman {
		salesOrder.Salesman = request.SalesOrder.Salesman
	}
	if salesOrder.VisitationID == x.VisitationID {
		salesOrder.VisitationID = request.SalesOrder.VisitationID
	}
	if salesOrder.OrderSourceID == x.OrderSourceID {
		salesOrder.OrderSourceID = request.SalesOrder.OrderSourceID
	}
	if salesOrder.OrderSourceName == x.OrderSourceName {
		salesOrder.OrderSourceName = request.SalesOrder.OrderSourceName
	}
	if salesOrder.OrderSource == x.OrderSource {
		salesOrder.OrderSource = request.SalesOrder.OrderSource
	}
	if salesOrder.OrderStatusID == x.OrderStatusID {
		salesOrder.OrderStatusID = request.SalesOrder.OrderStatusID
	}
	if salesOrder.OrderStatus == x.OrderStatus {
		salesOrder.OrderStatus = request.SalesOrder.OrderStatus
	}
	if salesOrder.OrderStatusName == x.OrderStatusName {
		salesOrder.OrderStatusName = request.SalesOrder.OrderStatusName
	}
	if salesOrder.SoCode == x.SoCode {
		salesOrder.SoCode = request.SalesOrder.SoCode
	}
	if salesOrder.SoDate == x.SoDate {
		salesOrder.SoDate = request.SalesOrder.SoDate
	}
	if salesOrder.SoRefCode == x.SoRefCode {
		salesOrder.SoRefCode = request.SalesOrder.SoRefCode
	}
	if salesOrder.SoRefDate == x.SoRefDate {
		salesOrder.SoRefDate = request.SalesOrder.SoRefDate
	}
	if salesOrder.ReferralCode == x.ReferralCode {
		salesOrder.ReferralCode = request.SalesOrder.ReferralCode
	}
	if salesOrder.GLat == x.GLat {
		salesOrder.GLat = request.SalesOrder.GLat
	}
	if salesOrder.GLong == x.GLong {
		salesOrder.GLong = request.SalesOrder.GLong
	}
	if salesOrder.DeviceId == x.DeviceId {
		salesOrder.DeviceId = request.SalesOrder.DeviceId
	}
	if salesOrder.Note == x.Note {
		salesOrder.Note = request.SalesOrder.Note
	}
	if salesOrder.InternalComment == x.InternalComment {
		salesOrder.InternalComment = request.SalesOrder.InternalComment
	}
	if salesOrder.TotalAmount == x.TotalAmount {
		salesOrder.TotalAmount = request.SalesOrder.TotalAmount
	}
	if salesOrder.TotalTonase == x.TotalTonase {
		salesOrder.TotalTonase = request.SalesOrder.TotalTonase
	}
	if salesOrder.IsDoneSyncToEs == x.IsDoneSyncToEs {
		salesOrder.IsDoneSyncToEs = request.SalesOrder.IsDoneSyncToEs
	}
	if salesOrder.StartDateSyncToEs == x.StartDateSyncToEs {
		salesOrder.StartDateSyncToEs = request.SalesOrder.StartDateSyncToEs
	}
	if salesOrder.EndDateSyncToEs == x.EndDateSyncToEs {
		salesOrder.EndDateSyncToEs = request.SalesOrder.EndDateSyncToEs
	}
	if salesOrder.StartCreatedDate == x.StartCreatedDate {
		salesOrder.StartCreatedDate = request.SalesOrder.StartCreatedDate
	}
	if salesOrder.EndCreatedDate == x.EndCreatedDate {
		salesOrder.EndCreatedDate = request.SalesOrder.EndCreatedDate
	}
	if salesOrder.SalesOrderDetails == nil {
		salesOrder.SalesOrderDetails = request.SalesOrder.SalesOrderDetails
	}
	if salesOrder.DeliveryOrders == nil {
		salesOrder.DeliveryOrders = request.SalesOrder.DeliveryOrders
	}
	if salesOrder.SalesmanID == x.SalesmanID {
		salesOrder.SalesmanID = request.SalesOrder.SalesmanID
	}
	if salesOrder.SalesmanName == x.SalesmanName {
		salesOrder.SalesmanName = request.SalesOrder.SalesmanName
	}
	if salesOrder.SalesmanEmail == x.SalesmanEmail {
		salesOrder.SalesmanEmail = request.SalesOrder.SalesmanEmail
	}
	if salesOrder.SalesOrderLogID == x.SalesOrderLogID {
		salesOrder.SalesOrderLogID = request.SalesOrder.SalesOrderLogID
	}
	if salesOrder.CreatedBy == x.CreatedBy {
		salesOrder.CreatedBy = request.SalesOrder.CreatedBy
	}
	if salesOrder.LatestUpdatedBy == x.LatestUpdatedBy {
		salesOrder.LatestUpdatedBy = request.SalesOrder.LatestUpdatedBy
	}
	if salesOrder.CreatedAt == x.CreatedAt {
		salesOrder.CreatedAt = request.SalesOrder.CreatedAt
	}
	if salesOrder.UpdatedAt == x.UpdatedAt {
		salesOrder.UpdatedAt = request.SalesOrder.UpdatedAt
	}
	if salesOrder.DeletedAt == x.DeletedAt {
		salesOrder.DeletedAt = request.SalesOrder.DeletedAt
	}
	for _, v := range request.SalesOrder.SalesOrderDetails {
		for k, y := range salesOrder.SalesOrderDetails {
			if y.ID == v.ID {
				salesOrder.SalesOrderDetails[k] = v
			}
		}
	}
}

func (salesOrder *SalesOrder) UpdateSalesOrderChanMap(request *SalesOrderChan) {
	salesOrder.AgentID = request.SalesOrder.AgentID
	salesOrder.StoreID = request.SalesOrder.StoreID
	salesOrder.BrandID = request.SalesOrder.BrandID
	salesOrder.UserID = request.SalesOrder.UserID
	salesOrder.OrderSourceID = request.SalesOrder.OrderSourceID
	salesOrder.GLat = request.SalesOrder.GLat
	salesOrder.GLong = request.SalesOrder.GLong
	salesOrder.SoRefCode = request.SalesOrder.SoRefCode
	salesOrder.SoDate = request.SalesOrder.SoDate
	salesOrder.SoRefDate = request.SalesOrder.SoRefDate
	salesOrder.Note = request.SalesOrder.Note
	salesOrder.InternalComment = request.SalesOrder.InternalComment
	salesOrder.TotalAmount = request.SalesOrder.TotalAmount
	salesOrder.TotalTonase = request.SalesOrder.TotalTonase
	salesOrder.DeviceId = request.SalesOrder.DeviceId
	salesOrder.ReferralCode = request.SalesOrder.ReferralCode
	salesOrder.UpdatedAt = request.SalesOrder.UpdatedAt
	salesOrder.LatestUpdatedBy = request.SalesOrder.UserID
	return
}

func (salesOrder *SalesOrder) OrderSourceChanMap(request *OrderSourceChan) {
	salesOrder.OrderSource = request.OrderSource
	salesOrder.OrderSourceName = request.OrderSource.SourceName
	return
}

func (salesOrder *SalesOrder) BrandChanMap(request *BrandChan) {
	salesOrder.Brand = request.Brand
	salesOrder.BrandName = request.Brand.Name
	return
}

func (salesOrder *SalesOrder) AgentChanMap(request *AgentChan) {
	salesOrder.Agent = request.Agent
	salesOrder.AgentName = NullString{sql.NullString{String: request.Agent.Name, Valid: true}}
	salesOrder.AgentEmail = request.Agent.Email
	salesOrder.AgentProvinceName = request.Agent.ProvinceName
	salesOrder.AgentCityName = request.Agent.CityName
	salesOrder.AgentDistrictName = request.Agent.DistrictName
	salesOrder.AgentVillageName = request.Agent.VillageName
	salesOrder.AgentAddress = request.Agent.Address
	salesOrder.AgentPhone = request.Agent.Phone
	salesOrder.AgentMainMobilePhone = request.Agent.MainMobilePhone
	return
}

func (salesOrder *SalesOrder) StoreChanMap(request *StoreChan) {
	salesOrder.Store = request.Store
	salesOrder.StoreName = request.Store.Name
	salesOrder.StoreCode = request.Store.StoreCode
	salesOrder.StoreEmail = request.Store.Email
	salesOrder.StoreProvinceName = request.Store.ProvinceName
	salesOrder.StoreCityName = request.Store.CityName
	salesOrder.StoreDistrictName = request.Store.DistrictName
	salesOrder.StoreVillageName = request.Store.VillageName
	salesOrder.StoreAddress = request.Store.Address
	salesOrder.StorePhone = request.Store.Phone
	salesOrder.StoreMainMobilePhone = request.Store.MainMobilePhone
	return
}

func (salesOrder *SalesOrder) UserChanMap(request *UserChan) {
	salesOrder.User = request.User
	salesOrder.UserFirstName = request.User.FirstName
	salesOrder.UserLastName = request.User.LastName
	salesOrder.UserEmail = NullString{sql.NullString{String: request.User.Email, Valid: true}}
	return
}

func (salesOrder *SalesOrder) SalesmanChanMap(request *SalesmanChan) {
	salesOrder.Salesman = request.Salesman
	salesOrder.SalesmanID = NullInt64{sql.NullInt64{Int64: int64(request.Salesman.ID), Valid: true}}
	salesOrder.SalesmanName = NullString{sql.NullString{String: request.Salesman.Name, Valid: true}}
	salesOrder.SalesmanEmail = request.Salesman.Email
	return
}

func (result *SalesOrderResponse) CreateSoResponseMap(request *SalesOrder) {
	result.CartID = request.CartID
	result.AgentID = request.AgentID
	result.StoreID = request.StoreID
	result.StoreCode = request.StoreCode.String
	result.StoreName = request.StoreName.String
	result.StoreStatus = request.Store.Status.String
	result.StorePhone = request.StorePhone
	result.StoreOwner = request.AgentName.String
	result.StoreProvinceId, _ = strconv.Atoi(request.Store.ProvinceID.String)
	result.StoreProvinceName = request.StoreProvinceName.String
	result.StoreCityId, _ = strconv.Atoi(request.Store.CityID.String)
	result.StoreCityName = request.StoreCityName.String
	result.StoreDistrictId, _ = strconv.Atoi(request.Store.DistrictID.String)
	result.StoreDistrictName = request.StoreDistrictName.String
	result.StoreAddress = request.StoreAddress.String
	result.BrandID = request.BrandID
	result.BrandName = request.BrandName
	result.UserID = request.UserID
	result.SalesmanID = int(request.SalesmanID.Int64)
	result.SalesmanName = request.SalesmanName.String
	result.VisitationID = request.VisitationID
	result.OrderSourceID = request.OrderSourceID
	result.OrderSourceName = request.OrderSource.SourceName
	result.OrderStatusID = request.OrderStatusID
	result.OrderStatusName = request.OrderStatus.Name
	result.SoCode = request.SoCode
	result.SoDate = request.SoDate
	result.SoRefCode = request.SoRefCode.String
	result.SoRefDate = request.SoRefDate.String
	result.GLong = request.GLong.Float64
	result.GLat = request.GLat.Float64
	result.Note = request.Note.String
	result.InternalComment = request.InternalComment.String
	result.TotalAmount = request.TotalAmount
	result.TotalTonase = request.TotalTonase
	result.DeviceId = request.DeviceId.String
	result.ReferralCode = request.ReferralCode.String
	return
}

func (result *SalesOrderResponse) UpdateSoResponseMap(request *SalesOrder) {
	result.ID = request.ID
	result.AgentID = request.AgentID
	result.AgentName = request.AgentName.String
	result.AgentEmail = request.AgentEmail.String
	result.AgentProvinceName = request.AgentProvinceName.String
	result.AgentCityName = request.AgentCityName.String
	result.AgentDistrictName = request.AgentDistrictName.String
	result.AgentVillageName = request.AgentVillageName.String
	result.AgentAddress = request.AgentAddress.String
	result.AgentPhone = request.AgentPhone.String
	result.AgentMainMobilePhone = request.AgentMainMobilePhone.String
	result.StoreID = request.StoreID
	result.StoreName = request.StoreName.String
	result.StoreCode = request.StoreCode.String
	result.StoreEmail = request.StoreEmail.String
	result.StoreProvinceName = request.StoreProvinceName.String
	result.StoreCityName = request.StoreCityName.String
	result.StoreDistrictName = request.StoreDistrictName.String
	result.StoreVillageName = request.StoreVillageName.String
	result.StoreAddress = request.StoreAddress.String
	result.StorePhone = request.StorePhone
	result.StoreMainMobilePhone = request.StoreMainMobilePhone.String
	result.BrandID = request.BrandID
	result.BrandName = request.BrandName
	result.UserID = request.UserID
	result.UserFirstName = request.UserFirstName.String
	result.UserLastName = request.UserLastName.String
	result.UserEmail = request.UserEmail.String
	result.OrderSourceID = request.OrderSourceID
	result.OrderSourceName = request.OrderSource.SourceName
	result.OrderStatusID = request.OrderStatusID
	result.OrderStatusName = request.OrderStatus.Name
	result.SoCode = request.SoCode
	result.SoDate = request.SoDate
	result.SoRefCode = request.SoRefCode.String
	result.SoRefDate = request.SoRefDate.String
	result.GLong = request.GLong.Float64
	result.GLat = request.GLat.Float64
	result.Note = request.Note.String
	result.InternalComment = request.InternalComment.String
	result.TotalAmount = request.TotalAmount
	result.TotalTonase = request.TotalTonase
	result.StartCreatedDate = request.StartCreatedDate
	result.SalesmanName = request.SalesmanName.String
	result.SalesmanEmail = request.SalesmanEmail.String
	result.CreatedAt = request.CreatedAt
	return
}

func (salesOrder *SalesOrderOpenSearchResponse) SalesOrderOpenSearchResponseMap(request *SalesOrder) {
	salesOrder.ID = request.ID
	salesOrder.AgentName = request.AgentName
	salesOrder.AgentEmail = request.AgentEmail
	salesOrder.AgentName = request.AgentName
	salesOrder.AgentCityName = request.AgentCityName
	salesOrder.AgentVillageName = request.AgentVillageName
	salesOrder.AgentAddress = request.AgentAddress
	salesOrder.AgentPhone = request.AgentPhone
	salesOrder.AgentMainMobilePhone = request.AgentMainMobilePhone

	salesOrder.StoreName = request.StoreName
	salesOrder.StoreCode = request.StoreCode
	salesOrder.StoreEmail = request.StoreEmail
	salesOrder.StoreName = request.StoreName
	salesOrder.StoreCityName = request.StoreCityName
	salesOrder.StoreDistrictName = request.StoreDistrictName
	salesOrder.StoreVillageName = request.StoreVillageName
	salesOrder.StoreAddress = request.StoreAddress
	salesOrder.StorePhone = request.StorePhone
	salesOrder.StoreMainMobilePhone = request.StoreMainMobilePhone

	salesOrder.BrandName = request.BrandName

	salesOrder.UserFirstName = request.UserFirstName
	salesOrder.UserLastName = request.UserLastName
	salesOrder.UserEmail = request.UserEmail

	salesOrder.OrderSourceName = request.OrderSourceName
	salesOrder.OrderStatusName = request.OrderStatusName

	salesOrder.SoCode = request.SoCode
	salesOrder.SoDate = request.SoDate
	salesOrder.SoRefCode = request.SoRefCode
	salesOrder.SoRefDate = request.SoRefDate
	salesOrder.GLat = request.GLat
	salesOrder.GLong = request.GLong
	salesOrder.Note = request.Note
	salesOrder.ReferralCode = request.ReferralCode
	salesOrder.InternalComment = request.InternalComment
	salesOrder.TotalAmount = request.TotalAmount
	salesOrder.TotalTonase = request.TotalTonase

	salesOrder.SalesmanName = request.SalesmanName
	salesOrder.SalesmanEmail = request.SalesmanEmail
	salesOrder.CreatedAt = request.CreatedAt
	salesOrder.UpdatedAt = request.UpdatedAt

	return
}

func (salesOrderEventLog *SalesOrderEventLogResponse) SalesOrderEventLogResponseMap(request *GetSalesOrderLog) {
	salesOrderEventLog.ID = request.ID
	salesOrderEventLog.RequestID = request.RequestID
	salesOrderEventLog.SoCode = request.SoCode
	salesOrderEventLog.Status = request.Status
	salesOrderEventLog.Action = request.Action
	salesOrderEventLog.CreatedAt = request.CreatedAt
	salesOrderEventLog.UpdatedAt = request.UpdatedAt
	return
}

func (dataSOEventLog *DataSOEventLogResponse) DataSOEventLogResponseMap(request *GetSalesOrderLog) {
	dataSOEventLog.AgentID = request.Data.AgentID
	dataSOEventLog.AgentName = request.Data.AgentName.String
	dataSOEventLog.StoreCode = request.Data.StoreCode.String
	dataSOEventLog.StoreName = request.Data.StoreName.String
	dataSOEventLog.SalesID = int(request.Data.SalesmanID.Int64)
	dataSOEventLog.SalesName = request.Data.SalesmanName.String
	dataSOEventLog.OrderDate = request.Data.CreatedAt
	dataSOEventLog.StartOrderAt = request.Data.StartCreatedDate
	dataSOEventLog.OrderNote = NullString{sql.NullString{String: request.Data.Note.String, Valid: true}}
	dataSOEventLog.InternalNote = NullString{sql.NullString{String: request.Data.InternalComment.String, Valid: true}}
	dataSOEventLog.BrandCode = request.Data.BrandID
	dataSOEventLog.BrandName = request.Data.BrandName
	return
}

func (soDetailEventLogResponse *SODetailEventLogResponse) SoDetailEventLogResponse(request *SalesOrderDetail) {
	soDetailEventLogResponse.ID = request.ID
	soDetailEventLogResponse.SalesOrderID = request.SalesOrderID
	soDetailEventLogResponse.ProductID = request.ProductID
	soDetailEventLogResponse.OrderQty = request.Qty
	soDetailEventLogResponse.UomID = request.UomID
	return
}

func (salesOrderJourney *SalesOrderJourneyResponse) SalesOrderJourneyResponseMap(request *SalesOrderJourneys) {
	salesOrderJourney.ID = request.ID
	salesOrderJourney.SoId = request.SoId
	salesOrderJourney.SoCode = request.SoCode
	salesOrderJourney.SoDate = request.SoDate
	salesOrderJourney.OrderStatusName = request.Status
	salesOrderJourney.Remark = NullString{sql.NullString{String: request.Remark, Valid: true}}
	salesOrderJourney.Reason = NullString{sql.NullString{String: request.Reason, Valid: true}}
	salesOrderJourney.CreatedAt = request.CreatedAt
	salesOrderJourney.UpdatedAt = request.UpdatedAt
	return
}

func (salesOrder *SalesOrder) SalesOrderUploadSOSJMap(request *UploadSOSJField, now time.Time) {

	salesOrder.AgentID = request.IDDistributor
	salesOrder.StoreID = request.KodeTokoDBO
	salesOrder.BrandID = request.IDMerk
	salesOrder.SalesmanID = NullInt64{sql.NullInt64{Int64: int64(request.IDSalesman), Valid: true}}
	salesOrder.SoDate = request.TglSuratJalan
	salesOrder.SoRefDate = NullString{sql.NullString{String: request.TglSuratJalan, Valid: true}}
	salesOrder.Note = NullString{sql.NullString{String: request.Catatan, Valid: true}}
	salesOrder.InternalComment = NullString{sql.NullString{String: request.CatatanInternal, Valid: true}}
	salesOrder.IsDoneSyncToEs = "0"
	salesOrder.CreatedAt = &now
	salesOrder.StartDateSyncToEs = &now
	salesOrder.StartCreatedDate = &now
}
