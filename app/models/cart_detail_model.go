package models

import (
	"order-service/global/utils/model"
	"time"
)

type CartDetail struct {
	ID            int        `json:"id,omitempty"`
	CartID        int        `json:"cart_id,omitempty"`
	BrandID       int        `json:"brand_id,omitempty"`
	ProductID     int        `json:"product_id,omitempty"`
	UomID         int        `json:"uom_id,omitempty"`
	OrderStatusID int        `json:"order_status_id,omitempty"`
	Qty           int        `json:"qty,omitempty"`
	Price         float64    `json:"price,omitempty"`
	Note          string     `json:"note,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

type CartDetailChan struct {
	CartDetail *CartDetail
	Error      error
	ErrorLog   *model.ErrorLog
	ID         int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type CartDetailsChan struct {
	CartDetails []*CartDetail
	Total       int64
	Error       error
	ErrorLog    *model.ErrorLog
	ID          int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type CartDetailRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type CartDetails struct {
	CartDetails []*CartDetail `json:"cart_details,omitempty"`
	Total       int64         `json:"total,omitempty"`
}
