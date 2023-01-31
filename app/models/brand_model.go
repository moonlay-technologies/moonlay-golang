package models

import (
	"poc-order-service/global/utils/model"
	"time"
)

type Brand struct {
	ID           int        `json:"id,omitempty" bson:"id,omitempty"`
	Name         string     `json:"name,omitempty" bson:"name,omitempty"`
	PrincipleID  int        `json:"principle_id,omitempty" bson:"principle_id,omitempty"`
	Description  NullString `json:"description,omitempty" bson:"description,omitempty"`
	Uid          string     `json:"uid,omitempty" bson:"uid,omitempty"`
	Image        string     `json:"image,omitempty" bson:"image,omitempty"`
	DataType     string     `json:"data_type,omitempty" bson:"data_type,omitempty"`
	StartDate    *time.Time `json:"start_date,omitempty" bson:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty" bson:"end_date,omitempty"`
	StatusActive int        `json:"status_active,omitempty" bson:"status_active,omitempty"`
	Priority     int        `json:"priority,omitempty" bson:"priority,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type BrandChan struct {
	Brand    *Brand
	Total    int64
	Error    error
	ErrorLog *model.ErrorLog
	ID       int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type BrandsChan struct {
	Brands   []*Brand
	Total    int64
	Error    error
	ErrorLog *model.ErrorLog
	ID       int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type BrandRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type Brands struct {
	Brands []*Brand `json:"brands,omitempty"`
	Total  int64    `json:"total,omitempty"`
}
