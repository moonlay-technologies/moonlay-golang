package models

import "poc-order-service/global/utils/model"

type StoreAddress struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	StoreID     int    `json:"store_id,omitempty"`
	Email       string `json:"email,omitempty"`
	MobilePhone string `json:"mobile_phone,omitempty"`
	Phone       string `json:"phone,omitempty"`
	PicName     string `json:"pic_name,omitempty"`
	Description string `json:"description,omitempty"`
	IsMain      int    `json:"is_main,omitempty"`
	IsWarehouse int    `json:"is_warehouse,omitempty"`
	GLat        string `json:"g_lat,omitempty"`
	GLng        string `json:"g_lng,omitempty"`
	GPlaceID    string `json:"g_place_id,omitempty"`
	Address     string `json:"address,omitempty"`
	ProvinceID  int    `json:"province_id,omitempty"`
	CityID      int    `json:"city_id,omitempty"`
	DistrictID  int    `json:"district_id,omitempty"`
	VillageID   int    `json:"village_id,omitempty"`
}

type StoreAddressChan struct {
	StoreAddress *StoreAddress
	Error        error
	ErrorLog     *model.ErrorLog
	ID           int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type StoreAddressesChan struct {
	StoreAddresss []*StoreAddress
	Total         int64
	Error         error
	ErrorLog      *model.ErrorLog
	ID            int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type StoreAddressRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type StoreAddresses struct {
	StoreAddresses []*AgentBrand `json:"store_addresses,omitempty"`
	Total          int64         `json:"total,omitempty"`
}
