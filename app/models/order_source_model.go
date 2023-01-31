package models

import (
	"poc-order-service/global/utils/model"
	"time"
)

type OrderSource struct {
	ID         int        `json:"id,omitempty" bson:"id"`
	SourceName string     `json:"source_name,omitempty" bson:"source_name"`
	Code       string     `json:"code,omitempty" bson:"code"`
	ParentID   int        `json:"parent_id,omitempty" bson:"parent_id"`
	CreatedAt  *time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" bson:"deleted_at"`
}

type OrderSourceChan struct {
	OrderSource *OrderSource
	Total       int64
	Error       error
	ErrorLog    *model.ErrorLog
	ID          int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type OrderSourcesChan struct {
	OrderSources []*OrderSource
	Total        int64
	Error        error
	ErrorLog     *model.ErrorLog
	ID           int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type OrderSourceRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type OrderSources struct {
	OrderSources []*OrderSource `json:"order_sources,omitempty"`
	Total        int64          `json:"total,omitempty"`
}
