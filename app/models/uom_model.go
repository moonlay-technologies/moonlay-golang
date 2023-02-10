package models

import (
	"order-service/global/utils/model"
	"time"
)

type Uom struct {
	ID        int        `json:"id,omitempty"`
	Name      NullString `json:"name,omitempty"`
	Code      NullString `json:"code,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type UomChan struct {
	Uom      *Uom
	Error    error
	ErrorLog *model.ErrorLog
	Total    int64
	ID       int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type UomsChan struct {
	Uoms     []*Uom
	Error    error
	ErrorLog *model.ErrorLog
	Total    int64
	ID       int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type UomRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type Uoms struct {
	Uoms  []*Uom `json:"uoms,omitempty"`
	Total int64  `json:"total,omitempty"`
}

type UomOpenSearchResponse struct {
	Name NullString `json:"name,omitempty"`
	Code NullString `json:"code,omitempty"`
}
