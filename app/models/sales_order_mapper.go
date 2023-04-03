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
	salesOrder.LatestUpdatedBy = request.UserID
	return
}

func (salesOrder *SalesOrder) SalesOrderOpenSearchChanMap(request *SalesOrderChan) {
	defaultSo := &SalesOrder{}
	if salesOrder.ID == defaultSo.ID {
		salesOrder.ID = request.SalesOrder.ID
	}
	if salesOrder.CartID == defaultSo.CartID {
		salesOrder.CartID = request.SalesOrder.CartID
	}
	if salesOrder.AgentID == defaultSo.AgentID {
		salesOrder.AgentID = request.SalesOrder.AgentID
	}
	if salesOrder.AgentName == defaultSo.AgentName {
		salesOrder.AgentName = request.SalesOrder.AgentName
	}
	if salesOrder.AgentEmail == defaultSo.AgentEmail {
		salesOrder.AgentEmail = request.SalesOrder.AgentEmail
	}
	if salesOrder.AgentProvinceID == defaultSo.AgentProvinceID {
		salesOrder.AgentProvinceID = request.SalesOrder.AgentProvinceID
	}
	if salesOrder.AgentProvinceName == defaultSo.AgentProvinceName {
		salesOrder.AgentProvinceName = request.SalesOrder.AgentProvinceName
	}
	if salesOrder.AgentCityID == defaultSo.AgentCityID {
		salesOrder.AgentCityID = request.SalesOrder.AgentCityID
	}
	if salesOrder.AgentCityName == defaultSo.AgentCityName {
		salesOrder.AgentCityName = request.SalesOrder.AgentCityName
	}
	if salesOrder.AgentDistrictID == defaultSo.AgentDistrictID {
		salesOrder.AgentDistrictID = request.SalesOrder.AgentDistrictID
	}
	if salesOrder.AgentDistrictName == defaultSo.AgentDistrictName {
		salesOrder.AgentDistrictName = request.SalesOrder.AgentDistrictName
	}
	if salesOrder.AgentVillageID == defaultSo.AgentVillageID {
		salesOrder.AgentVillageID = request.SalesOrder.AgentVillageID
	}
	if salesOrder.AgentVillageName == defaultSo.AgentVillageName {
		salesOrder.AgentVillageName = request.SalesOrder.AgentVillageName
	}
	if salesOrder.AgentAddress == defaultSo.AgentAddress {
		salesOrder.AgentAddress = request.SalesOrder.AgentAddress
	}
	if salesOrder.AgentPhone == defaultSo.AgentPhone {
		salesOrder.AgentPhone = request.SalesOrder.AgentPhone
	}
	if salesOrder.AgentMainMobilePhone == defaultSo.AgentMainMobilePhone {
		salesOrder.AgentMainMobilePhone = request.SalesOrder.AgentMainMobilePhone
	}
	if salesOrder.Agent == defaultSo.Agent {
		salesOrder.Agent = request.SalesOrder.Agent
	}
	if salesOrder.StoreID == defaultSo.StoreID {
		salesOrder.StoreID = request.SalesOrder.StoreID
	}
	if salesOrder.StoreName == defaultSo.StoreName {
		salesOrder.StoreName = request.SalesOrder.StoreName
	}
	if salesOrder.StoreCode == defaultSo.StoreCode {
		salesOrder.StoreCode = request.SalesOrder.StoreCode
	}
	if salesOrder.StoreEmail == defaultSo.StoreEmail {
		salesOrder.StoreEmail = request.SalesOrder.StoreEmail
	}
	if salesOrder.StoreProvinceID == defaultSo.StoreProvinceID {
		salesOrder.StoreProvinceID = request.SalesOrder.StoreProvinceID
	}
	if salesOrder.StoreProvinceName == defaultSo.StoreProvinceName {
		salesOrder.StoreProvinceName = request.SalesOrder.StoreProvinceName
	}
	if salesOrder.StoreCityID == defaultSo.StoreCityID {
		salesOrder.StoreCityID = request.SalesOrder.StoreCityID
	}
	if salesOrder.StoreCityName == defaultSo.StoreCityName {
		salesOrder.StoreCityName = request.SalesOrder.StoreCityName
	}
	if salesOrder.StoreDistrictID == defaultSo.StoreDistrictID {
		salesOrder.StoreDistrictID = request.SalesOrder.StoreDistrictID
	}
	if salesOrder.StoreDistrictName == defaultSo.StoreDistrictName {
		salesOrder.StoreDistrictName = request.SalesOrder.StoreDistrictName
	}
	if salesOrder.StoreVillageID == defaultSo.StoreVillageID {
		salesOrder.StoreVillageID = request.SalesOrder.StoreVillageID
	}
	if salesOrder.StoreVillageName == defaultSo.StoreVillageName {
		salesOrder.StoreVillageName = request.SalesOrder.StoreVillageName
	}
	if salesOrder.StoreAddress == defaultSo.StoreAddress {
		salesOrder.StoreAddress = request.SalesOrder.StoreAddress
	}
	if salesOrder.StorePhone == defaultSo.StorePhone {
		salesOrder.StorePhone = request.SalesOrder.StorePhone
	}
	if salesOrder.StoreMainMobilePhone == defaultSo.StoreMainMobilePhone {
		salesOrder.StoreMainMobilePhone = request.SalesOrder.StoreMainMobilePhone
	}
	if salesOrder.Store == defaultSo.Store {
		salesOrder.Store = request.SalesOrder.Store
	}
	if salesOrder.BrandID == defaultSo.BrandID {
		salesOrder.BrandID = request.SalesOrder.BrandID
	}
	if salesOrder.BrandName == defaultSo.BrandName {
		salesOrder.BrandName = request.SalesOrder.BrandName
	}
	if salesOrder.Brand == defaultSo.Brand {
		salesOrder.Brand = request.SalesOrder.Brand
	}
	if salesOrder.UserID == defaultSo.UserID {
		salesOrder.UserID = request.SalesOrder.UserID
	}
	if salesOrder.UserFirstName == defaultSo.UserFirstName {
		salesOrder.UserFirstName = request.SalesOrder.UserFirstName
	}
	if salesOrder.UserLastName == defaultSo.UserLastName {
		salesOrder.UserLastName = request.SalesOrder.UserLastName
	}
	if salesOrder.UserRoleID == defaultSo.UserRoleID {
		salesOrder.UserRoleID = request.SalesOrder.UserRoleID
	}
	if salesOrder.UserEmail == defaultSo.UserEmail {
		salesOrder.UserEmail = request.SalesOrder.UserEmail
	}
	if salesOrder.User == defaultSo.User {
		salesOrder.User = request.SalesOrder.User
	}
	if salesOrder.Salesman == defaultSo.Salesman {
		salesOrder.Salesman = request.SalesOrder.Salesman
	}
	if salesOrder.VisitationID == defaultSo.VisitationID {
		salesOrder.VisitationID = request.SalesOrder.VisitationID
	}
	if salesOrder.OrderSourceID == defaultSo.OrderSourceID {
		salesOrder.OrderSourceID = request.SalesOrder.OrderSourceID
	}
	if salesOrder.OrderSourceName == defaultSo.OrderSourceName {
		salesOrder.OrderSourceName = request.SalesOrder.OrderSourceName
	}
	if salesOrder.OrderSource == defaultSo.OrderSource {
		salesOrder.OrderSource = request.SalesOrder.OrderSource
	}
	if salesOrder.OrderStatusID == defaultSo.OrderStatusID {
		salesOrder.OrderStatusID = request.SalesOrder.OrderStatusID
	}
	if salesOrder.OrderStatus == defaultSo.OrderStatus {
		salesOrder.OrderStatus = request.SalesOrder.OrderStatus
	}
	if salesOrder.OrderStatusName == defaultSo.OrderStatusName {
		salesOrder.OrderStatusName = request.SalesOrder.OrderStatusName
	}
	if salesOrder.SoCode == defaultSo.SoCode {
		salesOrder.SoCode = request.SalesOrder.SoCode
	}
	if salesOrder.SoDate == defaultSo.SoDate {
		salesOrder.SoDate = request.SalesOrder.SoDate
	}
	if salesOrder.SoRefCode == defaultSo.SoRefCode {
		salesOrder.SoRefCode = request.SalesOrder.SoRefCode
	}
	if salesOrder.SoRefDate == defaultSo.SoRefDate {
		salesOrder.SoRefDate = request.SalesOrder.SoRefDate
	}
	if salesOrder.ReferralCode == defaultSo.ReferralCode {
		salesOrder.ReferralCode = request.SalesOrder.ReferralCode
	}
	if salesOrder.GLat == defaultSo.GLat {
		salesOrder.GLat = request.SalesOrder.GLat
	}
	if salesOrder.GLong == defaultSo.GLong {
		salesOrder.GLong = request.SalesOrder.GLong
	}
	if salesOrder.DeviceId == defaultSo.DeviceId {
		salesOrder.DeviceId = request.SalesOrder.DeviceId
	}
	if salesOrder.Note == defaultSo.Note {
		salesOrder.Note = request.SalesOrder.Note
	}
	if salesOrder.InternalComment == defaultSo.InternalComment {
		salesOrder.InternalComment = request.SalesOrder.InternalComment
	}
	if salesOrder.TotalAmount == defaultSo.TotalAmount {
		salesOrder.TotalAmount = request.SalesOrder.TotalAmount
	}
	if salesOrder.TotalTonase == defaultSo.TotalTonase {
		salesOrder.TotalTonase = request.SalesOrder.TotalTonase
	}
	if salesOrder.IsDoneSyncToEs == defaultSo.IsDoneSyncToEs {
		salesOrder.IsDoneSyncToEs = request.SalesOrder.IsDoneSyncToEs
	}
	if salesOrder.StartDateSyncToEs == defaultSo.StartDateSyncToEs {
		salesOrder.StartDateSyncToEs = request.SalesOrder.StartDateSyncToEs
	}
	if salesOrder.EndDateSyncToEs == defaultSo.EndDateSyncToEs {
		salesOrder.EndDateSyncToEs = request.SalesOrder.EndDateSyncToEs
	}
	if salesOrder.StartCreatedDate == defaultSo.StartCreatedDate {
		salesOrder.StartCreatedDate = request.SalesOrder.StartCreatedDate
	}
	if salesOrder.EndCreatedDate == defaultSo.EndCreatedDate {
		salesOrder.EndCreatedDate = request.SalesOrder.EndCreatedDate
	}
	if salesOrder.SalesOrderDetails == nil {
		salesOrder.SalesOrderDetails = request.SalesOrder.SalesOrderDetails
	}
	if salesOrder.DeliveryOrders == nil {
		salesOrder.DeliveryOrders = request.SalesOrder.DeliveryOrders
	}
	if salesOrder.SalesmanID == defaultSo.SalesmanID {
		salesOrder.SalesmanID = request.SalesOrder.SalesmanID
	}
	if salesOrder.SalesmanName == defaultSo.SalesmanName {
		salesOrder.SalesmanName = request.SalesOrder.SalesmanName
	}
	if salesOrder.SalesmanEmail == defaultSo.SalesmanEmail {
		salesOrder.SalesmanEmail = request.SalesOrder.SalesmanEmail
	}
	if salesOrder.SalesOrderLogID == defaultSo.SalesOrderLogID {
		salesOrder.SalesOrderLogID = request.SalesOrder.SalesOrderLogID
	}
	if salesOrder.CreatedBy == defaultSo.CreatedBy {
		salesOrder.CreatedBy = request.SalesOrder.CreatedBy
	}
	if salesOrder.LatestUpdatedBy == defaultSo.LatestUpdatedBy {
		salesOrder.LatestUpdatedBy = request.SalesOrder.LatestUpdatedBy
	}
	if salesOrder.CreatedAt == defaultSo.CreatedAt {
		salesOrder.CreatedAt = request.SalesOrder.CreatedAt
	}
	if salesOrder.UpdatedAt == defaultSo.UpdatedAt {
		salesOrder.UpdatedAt = request.SalesOrder.UpdatedAt
	}
	if salesOrder.DeletedAt == defaultSo.DeletedAt {
		salesOrder.DeletedAt = request.SalesOrder.DeletedAt
	}
	for k, v := range request.SalesOrder.SalesOrderDetails {
		for _, y := range salesOrder.SalesOrderDetails {
			if y.ID == v.ID {
				request.SalesOrder.SalesOrderDetails[k] = y
			}
		}
	}
	salesOrder.SalesOrderDetails = request.SalesOrder.SalesOrderDetails
	return
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
	salesOrder.StoreID = request.Store.ID
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
	result.ID = request.ID
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
	salesOrder.AgentProvinceName = request.AgentProvinceName
	salesOrder.AgentCityName = request.AgentCityName
	salesOrder.AgentDistrictName = request.AgentDistrictName
	salesOrder.AgentVillageName = request.AgentVillageName
	salesOrder.AgentAddress = request.AgentAddress
	salesOrder.AgentPhone = request.AgentPhone
	salesOrder.AgentMainMobilePhone = request.AgentMainMobilePhone

	salesOrder.StoreName = request.StoreName
	salesOrder.StoreCode = request.StoreCode
	salesOrder.StoreEmail = request.StoreEmail
	salesOrder.StoreName = request.StoreName
	salesOrder.StoreProvinceName = request.StoreProvinceName
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

func (salesOrder *SalesOrder) SalesOrderForDOMap(request *SalesOrder) {
	salesOrder.ID = request.ID
	salesOrder.AgentID = request.AgentID
	salesOrder.AgentName = request.AgentName
	salesOrder.AgentEmail = request.AgentEmail
	salesOrder.AgentName = request.AgentName
	salesOrder.AgentProvinceName = request.AgentProvinceName
	salesOrder.AgentProvinceID = request.AgentProvinceID
	salesOrder.AgentCityName = request.AgentCityName
	salesOrder.AgentCityID = request.AgentCityID
	salesOrder.AgentDistrictName = request.AgentDistrictName
	salesOrder.AgentDistrictID = request.AgentDistrictID
	salesOrder.AgentVillageName = request.AgentVillageName
	salesOrder.AgentVillageID = request.AgentVillageID
	salesOrder.AgentAddress = request.AgentAddress
	salesOrder.AgentPhone = request.AgentPhone
	salesOrder.AgentMainMobilePhone = request.AgentMainMobilePhone

	salesOrder.StoreID = request.StoreID
	salesOrder.StoreName = request.StoreName
	salesOrder.StoreCode = request.StoreCode
	salesOrder.StoreEmail = request.StoreEmail
	salesOrder.StoreName = request.StoreName
	salesOrder.StoreProvinceName = request.StoreProvinceName
	salesOrder.StoreProvinceID = request.StoreProvinceID
	salesOrder.StoreCityName = request.StoreCityName
	salesOrder.StoreCityID = request.StoreCityID
	salesOrder.StoreDistrictName = request.StoreDistrictName
	salesOrder.StoreDistrictID = request.StoreDistrictID
	salesOrder.StoreVillageName = request.StoreVillageName
	salesOrder.StoreVillageID = request.StoreVillageID
	salesOrder.StoreAddress = request.StoreAddress
	salesOrder.StorePhone = request.StorePhone
	salesOrder.StoreMainMobilePhone = request.StoreMainMobilePhone

	salesOrder.BrandID = request.BrandID
	salesOrder.BrandName = request.BrandName

	salesOrder.UserID = request.UserID
	salesOrder.UserFirstName = request.UserFirstName
	salesOrder.UserLastName = request.UserLastName
	salesOrder.UserEmail = request.UserEmail
	salesOrder.UserRoleID = request.UserRoleID

	salesOrder.OrderSourceID = request.OrderSourceID
	salesOrder.OrderSourceName = request.OrderSourceName
	salesOrder.OrderStatusID = request.OrderStatusID
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

	salesOrder.SalesmanID = request.SalesmanID
	salesOrder.SalesmanName = request.SalesmanName
	salesOrder.SalesmanEmail = request.SalesmanEmail
	salesOrder.CreatedAt = request.CreatedAt
	salesOrder.UpdatedAt = request.UpdatedAt

	return
}

func (salesOrderEventLog *SalesOrderEventLogResponse) SalesOrderEventLogResponseMap(request *GetSalesOrderLog, status string) {
	salesOrderEventLog.ID = request.ID
	salesOrderEventLog.RequestID = request.RequestID
	salesOrderEventLog.SoCode = request.SoCode
	salesOrderEventLog.Status = status
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

func (salesOrder *SalesOrder) SalesOrderUploadSOMap(request *UploadSOField, now time.Time) {

	salesOrder.AgentID = request.IDDistributor
	salesOrder.BrandID = request.KodeMerk
	salesOrder.SalesmanID = NullInt64{sql.NullInt64{Int64: int64(request.IDSalesman), Valid: true}}
	salesOrder.SoDate = request.TanggalOrder
	salesOrder.SoRefDate = NullString{sql.NullString{String: request.TanggalTokoOrder, Valid: true}}
	salesOrder.Note = NullString{sql.NullString{String: request.CatatanOrder, Valid: true}}
	salesOrder.InternalComment = NullString{sql.NullString{String: request.CatatanInternal, Valid: true}}
	salesOrder.IsDoneSyncToEs = "0"
	salesOrder.CreatedAt = &now
	salesOrder.StartDateSyncToEs = &now
	salesOrder.StartCreatedDate = &now
}

func (d *SalesOrderRequest) SalesOrderExportMap(r *SalesOrderExportRequest) {
	d.ID = r.ID
	d.PerPage = 0
	d.Page = 0
	d.SortField = ""
	d.SortValue = ""
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
	d.Keyword = ""
	d.SoDate = r.SoDate
	d.StoreProvinceID = r.StoreProvinceID
	d.StoreCityID = r.StoreCityID
	d.StoreDistrictID = r.StoreDistrictID
	d.StoreVillageID = r.StoreDistrictID
}

func (d *SalesOrderCsvResponse) DoDetailMap(r *SalesOrder) {

}

func (data *SalesOrder) MapToCsvRow() []string {
	return []string{
		data.OrderStatusName,
		data.OrderSourceName,
		data.ReferralCode.String,
		data.SoRefCode.String,
		data.SoCode,
		data.SoDate,
		strconv.Itoa(data.AgentID),
		data.AgentName.String,
		strconv.Itoa(int(data.SalesmanID.Int64)),
		data.SalesmanName.String,
		strconv.Itoa(data.StoreID),
		data.Store.AliasCode.String,
		data.StoreCode.String,
		data.StoreName.String,
		data.Store.DistrictID.String,
		data.StoreDistrictName.String,
		data.Store.CityID.String,
		data.StoreCityName.String,
		strconv.Itoa(data.StoreProvinceID),
		data.StoreProvinceName.String,
		strconv.Itoa(data.BrandID),
		data.BrandName,
		"soAmount",
		"doAmount",
		data.Note.String,
		data.InternalComment.String,
		"alasanCancel",
		"alasanReject",
		data.SoRefDate.String,
		data.CreatedAt.String(),
		data.UpdatedAt.String(),
		strconv.Itoa(data.CreatedBy),
		strconv.Itoa(data.LatestUpdatedBy),
	}
}
