package models

import (
	"order-service/global/utils/model"
	"time"
)

type Agent struct {
	ID                   int        `json:"id,omitempty" bson:"id,omitempty"`
	Name                 string     `json:"name,omitempty" bson:"name,omitempty"`
	UID                  string     `json:"uid,omitempty" bson:"uid,omitempty"`
	ParentID             int        `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
	Email                NullString `json:"email,omitempty" bson:"email,omitempty"`
	ChangePasswordStatus int        `json:"change_password_status,omitempty" bson:"change_password_status,omitempty"`
	Description          NullString `json:"description,omitempty" bson:"description,omitempty"`
	Address              NullString `json:"address,omitempty" bson:"address,omitempty"`
	ProvinceID           NullString `json:"province_id,omitempty" bson:"province_id,omitempty"`
	ProvinceName         NullString `json:"province_name,omitempty" bson:"province_name,omitempty"`
	CityID               NullString `json:"city_id,omitempty" bson:"city_id,omitempty"`
	CityName             NullString `json:"city_name,omitempty" bson:"city_name,omitempty"`
	DistrictID           NullString `json:"district_id,omitempty" bson:"district_id,omitempty"`
	DistrictName         NullString `json:"district_name,omitempty" bson:"district_name,omitempty"`
	VillageID            NullString `json:"village_id,omitempty" bson:"village_id,omitempty"`
	VillageName          NullString `json:"village_name,omitempty" bson:"village_name,omitempty"`
	DataType             NullString `json:"data_type,omitempty" bson:"data_type,omitempty"`
	DistributorType      NullString `json:"distributor_type,omitempty" bson:"distributor_type,omitempty"`
	PostalCode           NullString `json:"postal_code,omitempty" bson:"postal_code,omitempty"`
	GPlaceID             NullString `json:"g_place_id,omitempty" bson:"g_place_id,omitempty"`
	GLat                 NullString `json:"g_lat,omitempty" bson:"g_lat,omitempty"`
	GLng                 NullString `json:"g_lng,omitempty" bson:"g_lng,omitempty"`
	ContactName          NullString `json:"contact_name,omitempty" bson:"contact_name,omitempty"`
	Website              NullString `json:"website,omitempty" bson:"website,omitempty"`
	Phone                NullString `json:"phone,omitempty" bson:"phone,omitempty"`
	MainMobilePhone      NullString `json:"main_mobile_phone,omitempty" bson:"main_mobile_phone,omitempty"`
	AlternatePhone1      NullString `json:"alternate_phone_1,omitempty" bson:"alternate_phone_1,omitempty"`
	AlternatePhone2      NullString `json:"alternate_phone_2,omitempty" bson:"alternate_phone_2,omitempty"`
	AlternatePhone3      NullString `json:"alternate_phone_3,omitempty" bson:"alternate_phone_3,omitempty"`
	YearEstablished      NullString `json:"year_established,omitempty" bson:"year_established,omitempty"`
	Status               NullString `json:"status,omitempty" bson:"status,omitempty"`
	NoNpwp               NullString `json:"no_npwp,omitempty" bson:"no_npwp,omitempty"`
	ImageNpwp            NullString `json:"image_npwp,omitempty" bson:"image_npwp,omitempty"`
	NoSiup               NullString `json:"no_siup,omitempty" bson:"no_siup,omitempty"`
	ImageSiup            NullString `json:"image_siup,omitempty" bson:"image_siup,omitempty"`
	Image                NullString `json:"image,omitempty" bson:"image,omitempty"`
	UserIDCreated        int        `json:"user_id_created,omitempty" bson:"user_id_created,omitempty"`
	UserIDUpdated        int        `json:"user_id_updated,omitempty" bson:"user_id_updated,omitempty"`
	CustomerCode         NullString `json:"customer_code,omitempty" bson:"customer_code,omitempty"`
	ApiUrl               NullString `json:"api_url,omitempty" bson:"api_url,omitempty"`
	IntegrasiApi         NullString `json:"integrasi_api,omitempty" bson:"integrasi_api,omitempty"`
	ApiUrlOrder          NullString `json:"api_url_order,omitempty" bson:"api_url_order,omitempty"`
	ApiEnvOrder          NullString `json:"api_env_order,omitempty" bson:"api_env_order,omitempty"`
	ResponseOrderApi     NullString `json:"response_order_api,omitempty" bson:"response_order_api,omitempty"`
	PicCollectionID      int        `json:"pic_collection_id,omitempty" bson:"pic_collection_id,omitempty"`
	PicCollectionName    NullString `json:"pic_collection_name,omitempty" bson:"pic_collection_name,omitempty"`
	PicCollectionTelp    NullString `json:"pic_collection_telp,omitempty" bson:"pic_collection_telp,omitempty"`
	Initial              NullString `json:"initial,omitempty" bson:"initial,omitempty"`
	CreatedAt            *time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt            *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type AgentChan struct {
	Agent    *Agent
	Total    int64
	Error    error
	ErrorLog *model.ErrorLog
	ID       int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type AgentsChan struct {
	Agents   []*Agent
	Total    int64
	Error    error
	ErrorLog *model.ErrorLog
	ID       int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type AgentRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type Agents struct {
	Agents []*Agent `json:"agents,omitempty"`
	Total  int64    `json:"total,omitempty"`
}
