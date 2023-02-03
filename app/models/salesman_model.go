package models

import (
	"order-service/global/utils/model"
	"time"
)

type Salesman struct {
	ID             int        `json:"id,omitempty"`
	ReferralCode   string     `json:"referral_code,omitempty"`
	Name           string     `json:"name,omitempty"`
	Email          NullString `json:"email,omitempty"`
	PhoneNumber    NullString `json:"phone_number,omitempty"`
	MobileNumber   NullString `json:"mobile_number,omitempty"`
	Address        string     `json:"address,omitempty"`
	ProvinceID     string     `json:"province_id,omitempty"`
	CityID         string     `json:"city_id,omitempty"`
	DistrictID     string     `json:"district_id,omitempty"`
	VillageID      string     `json:"village_id,omitempty"`
	GPlaceID       string     `json:"g_place_id,omitempty"`
	GLat           string     `json:"g_lat,omitempty"`
	GLng           string     `json:"g_lng,omitempty"`
	Province       string     `json:"province,omitempty"`
	Regency        string     `json:"regency,omitempty"`
	District       string     `json:"district,omitempty"`
	Village        string     `json:"village,omitempty"`
	NpKtp          string     `json:"np_ktp,omitempty"`
	NoKtpImage     string     `json:"no_ktp_image,omitempty"`
	NoSim          string     `json:"no_sim,omitempty"`
	NoSimImage     string     `json:"no_sim_image,omitempty"`
	AgentID        int        `json:"agent_id,omitempty"`
	SalesmanTypeID int        `json:"salesman_type_id,omitempty"`
	NameOnKtp      string     `json:"name_on_ktp,omitempty"`
	BirthDate      string     `json:"birth_date,omitempty"`
	BankName       string     `json:"bank_name,omitempty"`
	NameOnAccount  string     `json:"name_on_account,omitempty"`
	NoAccount      string     `json:"no_account,omitempty"`
	Uid            string     `json:"uid,omitempty"`
	IsVerified     int        `json:"is_verified,omitempty"`
	VerifiedBy     int        `json:"verified_by,omitempty"`
	VerifiedDate   *time.Time `json:"verified_date,omitempty"`
	SupervisorID   int        `json:"supervisor_id,omitempty"`
	IsDefault      int        `json:"is_default,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type SalesmanChan struct {
	Salesman *Salesman
	Total    int64
	Error    error
	ErrorLog *model.ErrorLog
	ID       int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type SalesmansChan struct {
	Salesmans []*Salesman
	Total     int64
	Error     error
	ErrorLog  *model.ErrorLog
	ID        int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type SalesmanRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type Salesmans struct {
	Salesmans []*Salesman `json:"salesmans,omitempty"`
	Total     int64       `json:"total,omitempty"`
}
