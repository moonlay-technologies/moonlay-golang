package models

import (
	"poc-order-service/global/utils/model"
	"time"
)

type Warehouse struct {
	ID              int        `json:"id,omitempty" bson:"id,omitempty"`
	Code            string     `json:"code,omitempty" bson:"code,omitempty"`
	Name            string     `json:"name,omitempty" bson:"name,omitempty"`
	OwnerID         int        `json:"owner_id,omitempty" bson:"owner_id,omitempty"`
	ProvinceID      NullString `json:"province_id,omitempty" bson:"province_id,omitempty"`
	ProvinceName    NullString `json:"province_name,omitempty" bson:"province_name,omitempty"`
	CityID          NullString `json:"city_id,omitempty" bson:"city_id,omitempty"`
	CityName        NullString `json:"city_name,omitempty" bson:"city_name,omitempty"`
	DistrictID      NullString `json:"district_id,omitempty" bson:"district_id,omitempty"`
	DistrictName    NullString `json:"district_name,omitempty" bson:"district_name,omitempty"`
	VillageID       NullString `json:"village_id,omitempty" bson:"village_id,omitempty"`
	VillageName     NullString `json:"village_name,omitempty" bson:"village_name,omitempty"`
	Address         NullString `json:"address,omitempty" bson:"address,omitempty"`
	Phone           NullString `json:"phone,omitempty" bson:"phone,omitempty"`
	MainMobilePhone NullString `json:"main_mobile_phone,omitempty" bson:"main_mobile_phone,omitempty"`
	Email           NullString `json:"email,omitempty" bson:"email,omitempty"`
	PicName         NullString `json:"pic_name,omitempty" bson:"pic_name,omitempty"`
	Status          int        `json:"status,omitempty" bson:"status,omitempty"`
	WarehouseTypeID int        `json:"warehouse_type_id,omitempty" bson:"warehouse_type_id,omitempty"`
	IsMain          int        `json:"is_main,omitempty" bson:"is_main,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

type WarehouseChan struct {
	Warehouse *Warehouse
	Total     int64
	Error     error
	ErrorLog  *model.ErrorLog
	ID        int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type WarehousesChan struct {
	Warehouses []*Warehouse
	Total      int64
	Error      error
	ErrorLog   *model.ErrorLog
	ID         int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type WarehouseRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type Warehouses struct {
	Warehouses []*Warehouse `json:"warehouses,omitempty"`
	Total      int64        `json:"total,omitempty"`
}
