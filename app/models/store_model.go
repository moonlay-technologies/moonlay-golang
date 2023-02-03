package models

import (
	"order-service/global/utils/model"
	"time"
)

type Store struct {
	ID                    int        `json:"id,omitempty"`
	StoreCode             NullString `json:"store_code,omitempty" bson:"store_code,omitempty"`
	StoreCategory         NullString `json:"store_category,omitempty"`
	Name                  NullString `json:"name,omitempty"`
	Email                 NullString `json:"email,omitempty"`
	EmailVerified         int        `json:"email_verified,omitempty"`
	Description           NullString `json:"description,omitempty"`
	Address               NullString `json:"address,omitempty"`
	ProvinceID            NullString `json:"province_id,omitempty"`
	ProvinceName          NullString `json:"province_name,omitempty"`
	CityID                NullString `json:"city_id,omitempty"`
	CityName              NullString `json:"city_name,omitempty"`
	DistrictID            NullString `json:"district_id,omitempty"`
	DistrictName          NullString `json:"district_name,omitempty"`
	VillageID             NullString `json:"village_id,omitempty"`
	VillageName           NullString `json:"village_name,omitempty"`
	DataType              NullString `json:"data_type,omitempty"`
	PostalCode            NullString `json:"postal_code,omitempty"`
	GPlaceID              NullString `json:"g_place_id,omitempty"`
	GLat                  NullString `json:"g_lat,omitempty"`
	GLng                  NullString `json:"g_lng,omitempty"`
	ContactName           NullString `json:"contact_name,omitempty"`
	Website               NullString `json:"website,omitempty"`
	Phone                 NullString `json:"phone,omitempty"`
	MainMobilePhone       NullString `json:"main_mobile_phone,omitempty"`
	AlternatePhone1       NullString `json:"alternate_phone_1,omitempty"`
	AlternatePhone2       NullString `json:"alternate_phone_2,omitempty"`
	AlternaltePhone3      NullString `json:"alternalte_phone_3,omitempty"`
	YearEstablished       NullString `json:"year_established,omitempty"`
	Status                NullString `json:"status,omitempty"`
	ProofOfBusiness       NullString `json:"proof_of_business,omitempty"`
	IsBlacklisted         int        `json:"is_blacklisted,omitempty"`
	NoNpwp                NullString `json:"no_npwp,omitempty"`
	ImageNpwp             NullString `json:"image_npwp,omitempty"`
	NpSiup                NullString `json:"np_siup,omitempty"`
	ImageSiup             NullString `json:"image_siup,omitempty"`
	Image                 NullString `json:"image,omitempty"`
	AliasName             NullString `json:"alias_name,omitempty"`
	AliasCode             NullString `json:"alias_code,omitempty"`
	Uid                   NullString `json:"uid,omitempty"`
	Creator               NullString `json:"creator,omitempty"`
	DBOApprovalStatus     int        `json:"dbo_approval_status,omitempty"`
	AgentID               int        `json:"agent_id,omitempty"`
	ParentID              int        `json:"parent_id,omitempty"`
	HeadID                int        `json:"head_id,omitempty"`
	UseApps               int        `json:"use_apps,omitempty"`
	AgentReference        NullString `json:"agent_reference,omitempty"`
	UserIDCreated         int        `json:"user_id_created,omitempty"`
	UserIDUpdated         int        `json:"user_id_updated,omitempty"`
	StatusPengajuan       NullString `json:"status_pengajuan,omitempty"`
	DateSubmiited         *time.Time `json:"date_submiited,omitempty"`
	DateProcessed         *time.Time `json:"date_processed,omitempty"`
	ResubmitAllowed       *time.Time `json:"resubmit_allowed,omitempty"`
	Remarks               NullString `json:"remarks,omitempty"`
	VerifiedDBO           NullString `json:"verified_dbo,omitempty"`
	VerifiedDate          *time.Time `json:"verified_date,omitempty"`
	SalesmanReferralCode  NullString `json:"salesman_referral_code,omitempty"`
	ValidationStore       NullString `json:"validation_store,omitempty"`
	Channel               NullString `json:"channel,omitempty"`
	HookStatus            NullString `json:"hook_status,omitempty"`
	SjFromStoreOrderCount int        `json:"sj_from_store_order_count,omitempty"`
	CreatedAt             *time.Time `json:"created_at,omitempty"`
	UpdatedAt             *time.Time `json:"updated_at,omitempty"`
	DeletedAt             *time.Time `json:"deleted_at,omitempty"`
}

type StoreChan struct {
	Store    *Store
	Total    int64
	Error    error
	ErrorLog *model.ErrorLog
	ID       int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type StoresChan struct {
	Stores   []*Store
	Total    int64
	Error    error
	ErrorLog *model.ErrorLog
	ID       int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type StoreRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type Stores struct {
	Stores []*Store `json:"stores,omitempty"`
	Total  int64    `json:"total,omitempty"`
}
