package models

import (
	"database/sql"
	"time"
)

func (salesOrder *SalesOrder) SalesOrderRequestMap(request *SalesOrderStoreRequest, now time.Time) {
	salesOrder.CartID = request.CartID
	salesOrder.AgentID = request.AgentID
	salesOrder.StoreID = request.StoreID
	salesOrder.BrandID = request.BrandID
	salesOrder.UserID = request.UserID
	salesOrder.VisitationID = request.VisitationID
	salesOrder.OrderSourceID = request.OrderSourceID
	salesOrder.GLat = NullFloat64{NullFloat64: sql.NullFloat64{Float64: request.GLat, Valid: true}}
	salesOrder.GLong = NullFloat64{NullFloat64: sql.NullFloat64{Float64: request.GLong, Valid: true}}
	salesOrder.SoRefCode = NullString{NullString: sql.NullString{String: request.SoRefCode, Valid: true}}
	salesOrder.SoDate = now.Format("2006-01-02")
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

func (salesOrder *SalesOrder) OrderStatusChanMap(request *OrderStatusChan) {
	salesOrder.OrderStatus = request.OrderStatus
	salesOrder.OrderStatusID = request.OrderStatus.ID
	salesOrder.OrderStatusName = request.OrderStatus.Name
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
	salesOrder.SalesmanName = NullString{sql.NullString{String: request.Salesman.Name, Valid: true}}
	salesOrder.SalesmanEmail = request.Salesman.Email
	return
}

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
	return
}

func (result *SalesOrderResponse) SoResponseMap(request *SalesOrder) {
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
	result.StorePhone = request.StorePhone.String
	result.StoreMainMobilePhone = request.StoreMainMobilePhone.String
	result.BrandID = request.BrandID
	result.BrandName = request.BrandName
	result.UserID = request.UserID
	result.UserFirstName = request.UserFirstName.String
	result.UserLastName = request.UserLastName.String
	result.UserEmail = request.UserEmail.String
	result.OrderSourceID = request.OrderSourceID
	result.OrderSourceName = request.OrderSourceName
	result.OrderStatusID = request.OrderStatusID
	result.OrderStatusName = request.OrderStatusName
	result.SoCode = request.SoCode
	result.SoDate = request.SoDate
	result.SoRefCode = request.SoRefCode.String
	result.SoRefDate = request.SoRefDate.String
	result.GLat = request.GLat.Float64
	result.GLong = request.GLong.Float64
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

func (result *SalesOrderResponse) SoUpdateByIdResponseMap(request *SalesOrder) {
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
	result.StorePhone = request.StorePhone.String
	result.StoreMainMobilePhone = request.StoreMainMobilePhone.String
	result.BrandID = request.BrandID
	result.BrandName = request.BrandName
	result.UserID = request.UserID
	result.UserFirstName = request.UserFirstName.String
	result.UserLastName = request.UserLastName.String
	result.UserEmail = request.UserEmail.String
	result.OrderSourceID = request.OrderSourceID
	result.OrderSourceName = request.OrderSourceName
	result.OrderStatusID = request.OrderStatusID
	result.OrderStatusName = request.OrderStatusName
	result.SoCode = request.SoCode
	result.SoDate = request.SoDate
	result.SoRefCode = request.SoRefCode.String
	result.SoRefDate = request.SoRefDate.String
	result.GLat = request.GLat.Float64
	result.GLong = request.GLong.Float64
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
