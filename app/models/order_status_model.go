package models

import (
	"order-service/global/utils/model"
	"time"
)

type OrderStatus struct {
	ID        int        `json:"id,omitempty" bson:"id"`
	Name      string     `json:"name,omitempty" bson:"name"`
	Sequence  int        `json:"sequence,omitempty" bson:"sequence"`
	OrderType string     `json:"order_type,omitempty" bson:"orderType"`
	CreatedAt *time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bson:"deleted_at"`
}

type OrderStatusChan struct {
	OrderStatus *OrderStatus
	Error       error
	ErrorLog    *model.ErrorLog
	Total       int64
	ID          int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type OrderStatusesChan struct {
	OrderStatuses []*OrderStatus
	Total         int64
	ErrorLog      *model.ErrorLog
	ID            int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type OrderStatusRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type OrderStatuses struct {
	OrderStatuses []*OrderStatus `json:"order_statuses,omitempty"`
	Total         int64          `json:"total,omitempty"`
}

type OrderStatusOpenSearchResponse struct {
	ID   int    `json:"id,omitempty" bson:"id"`
	Name string `json:"name,omitempty" bson:"name"`
}
