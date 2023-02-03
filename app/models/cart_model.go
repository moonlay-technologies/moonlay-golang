package models

import (
	"order-service/global/utils/model"
	"time"
)

type Cart struct {
	ID              int           `json:"id,omitempty"`
	AgentID         int           `json:"agent_id,omitempty"`
	BrandID         int           `json:"brand_id,omitempty"`
	VisitationID    int           `json:"visitation_id,omitempty"`
	UserID          int           `json:"user_id,omitempty"`
	StoreID         int           `json:"store_id,omitempty"`
	OrderStatusID   int           `json:"order_status_id,omitempty"`
	OrderSourceID   int           `json:"order_source_id,omitempty"`
	TotalTonase     float64       `json:"total_tonase,omitempty"`
	TotalAmount     float64       `json:"total_amount,omitempty"`
	Note            string        `json:"note,omitempty"`
	CartDetails     []*CartDetail `json:"cart_details,omitempty"`
	CreatedBy       int           `json:"created_by,omitempty" bson:"created_by,omitempty"`
	LatestUpdatedBy int           `json:"latest_updated_by,omitempty" bson:"latest_updated_by,omitempty"`
	CreatedAt       *time.Time    `json:"created_at,omitempty"`
	UpdatedAt       *time.Time    `json:"updated_at,omitempty"`
	DeletedAt       *time.Time    `json:"deleted_at,omitempty"`
}

type CartChan struct {
	Cart     *Cart
	Error    error
	ErrorLog *model.ErrorLog
	Total    int64
	ID       int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type CartsChan struct {
	Carts    []*Cart
	Total    int64
	Error    error
	ErrorLog *model.ErrorLog
	ID       int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type CartRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type Carts struct {
	Carts []*Cart `json:"carts,omitempty"`
	Total int64   `json:"total,omitempty"`
}
