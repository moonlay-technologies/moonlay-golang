package models

import (
	"order-service/global/utils/model"
	"time"
)

type DeliveryOrderDetail struct {
	ID                int               `json:"id,omitempty" bson:"id"`
	DeliveryOrderID   int               `json:"delivery_order_id,omitempty" bson:"delivery_order_id"`
	SoDetailID        int               `json:"so_detail_id,omitempty" bson:"so_detail_id"`
	SoDetail          *SalesOrderDetail `json:"so_detail,omitempty" bson:" so_detail"`
	BrandID           int               `json:"brand_id,omitempty" bson:"brand_id"`
	Brand             *Brand            `json:"brand,omitempty" bson:"brand"`
	BrandName         string            `json:"brand_name,omitempty"  bson:"brand_name"`
	ProductID         int               `json:"product_id,omitempty" bson:"product_id"`
	Product           *Product          `json:"product,omitempty" bson:"product"`
	ProductSKU        string            `json:"product_sku,omitempty" bson:"product_sku,omitempty"`
	ProductName       string            `json:"product_name,omitempty" bson:"product_name,omitempty"`
	UomID             int               `json:"uom_id,omitempty" bson:"uom_id"`
	Uom               *Uom              `json:"uom,omitempty" bson:"uom"`
	UomCode           string            `json:"uom_code,omitempty" bson:"uom_code,omitempty"`
	UomName           string            `json:"uom_name,omitempty" bson:"uom_name,omitempty"`
	OrderStatusID     int               `json:"order_status_id,omitempty" bson:"order_status_id"`
	OrderStatus       *OrderStatus      `json:"order_status,omitempty" bson:"order_status"`
	OrderStatusName   string            `json:"order_status_name,omitempty" bson:"order_status_name"`
	DoDetailCode      string            `json:"do_detail_code,omitempty" bson:"do_detail_code"`
	Qty               int               `json:"qty,omitempty" bson:"qty"`
	Note              NullString        `json:"note,omitempty" bson:"note"`
	IsDoneSyncToEs    string            `json:"is_done_sync_to_es,omitempty" bson:"is_done_sync_to_es"`
	StartDateSyncToEs *time.Time        `json:"start_date_sync_to_es,omitempty" bson:"start_date_sync_to_es"`
	EndDateSyncToEs   *time.Time        `json:"end_date_sync_to_es,omitempty" bson:"end_date_sync_to_es"`
	CreatedAt         *time.Time        `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt         *time.Time        `json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt         *time.Time        `json:"deleted_at,omitempty" bson:"deleted_at"`
}

type DeliveryOrderDetailStoreRequest struct {
	SoDetailID int    `json:"so_detail_id,omitempty"`
	Qty        int    `json:"qty,omitempty"`
	Note       string `json:"note,omitempty"`
}

type DeliveryOrderDetailChan struct {
	DeliveryOrderDetail *DeliveryOrderDetail
	Error               error
	ErrorLog            *model.ErrorLog
	ID                  int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type DeliveryOrderDetailsChan struct {
	DeliveryOrderDetails []*DeliveryOrderDetail
	Total                int64
	Error                error
	ErrorLog             *model.ErrorLog
	ID                   int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type DeliveryOrderDetailRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type DeliveryOrderDetails struct {
	DeliveryOrderDetails []*DeliveryOrderDetail `json:"delivery_order_details,omitempty"`
	Total                int64                  `json:"total,omitempty"`
}
