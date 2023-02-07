package models

import (
	"order-service/global/utils/model"
	"time"
)

type SalesOrderDetail struct {
	ID                int          `json:"id,omitempty" bson:"id,omitempty"`
	SalesOrderID      int          `json:"sales_order_id,omitempty" bson:"sales_order_id,omitempty"`
	ProductID         int          `json:"product_id,omitempty" bson:"product_id,omitempty" `
	BrandID           int          `json:"brand_id,omitempty" bson:"brand_id,omitempty" `
	Product           *Product     `json:"product,omitempty" bson:"product,omitempty"`
	UomID             int          `json:"uom_id,omitempty" bson:"uom_id,omitempty"`
	Uom               *Uom         `json:"uom,omitempty" bson:"uom,omitempty"`
	OrderStatusID     int          `json:"order_status_id,omitempty" bson:"order_status_id,omitempty"`
	OrderStatus       *OrderStatus `json:"order_status,omitempty" bson:"order_status,omitempty"`
	SoDetailCode      string       `json:"so_detail_code,omitempty" bson:"so_detail_code,omitempty"`
	Qty               int          `json:"qty,omitempty" bson:"qty,omitempty"`
	SentQty           int          `json:"sent_qty,omitempty" bson:"sent_qty,omitempty"`
	ResidualQty       int          `json:"residual_qty,omitempty" bson:"residual_qty,omitempty"`
	Price             float64      `json:"price,omitempty" bson:"price,omitempty"`
	Note              NullString   `json:"note,omitempty" bson:"note,omitempty"`
	IsDoneSyncToEs    string       `json:"is_done_sync_to_es,omitempty" bson:"is_done_sync_to_es,omitempty"`
	StartDateSyncToEs *time.Time   `json:"start_date_sync_to_es,omitempty" bson:"start_date_sync_to_es,omitempty"`
	EndDateSyncToEs   *time.Time   `json:"end_date_sync_to_es,omitempty" bson:"end_date_sync_to_es,omitempty"`
	ProductSKU        string       `json:"product_sku,omitempty" bson:"product_sku,omitempty"`
	ProductName       string       `json:"product_name,omitempty" bson:"product_name,omitempty"`
	UomCode           string       `json:"uom_code,omitempty" bson:"uom_code,omitempty"`
	UomName           string       `json:"uom_name,omitempty" bson:"uom_name,omitempty"`
	UomType           string       `json:"uom_type,omitempty" bson:"uom_type,omitempty"`
	Subtotal          float64      `json:"subtotal,omitempty" bson:"subtotal,omitempty"`
	OrderStatusName   string       `json:"order_status_name,omitempty" bson:"order_status_name,omitempty"`
	CreatedAt         *time.Time   `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt         *time.Time   `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt         *time.Time   `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type SalesOrderDetailTemplate struct {
	ProductID   int    `json:"product_id,omitempty" binding:"required"`
	UomID       int    `json:"uom_id,omitempty" binding:"required"`
	Qty         int    `json:"qty,omitempty" binding:"required"`
	SentQty     int    `json:"sent_qty,omitempty" binding:"required"`
	ResidualQty int    `json:"residual_qty,omitempty" binding:"required"`
	Note        string `json:"note,omitempty"`
}

type SalesOrderDetailStoreRequest struct {
	SalesOrderDetailTemplate
	SalesOrderId  int     `json:"sales_order_id,omitempty" binding:"required"`
	OrderStatusId int     `json:"order_status_id,omitempty" binding:"required"`
	SoDetailCode  string  `json:"so_detail_code,omitempty" binding:"required"`
	Price         float64 `json:"price,omitempty" binding:"required"`
}

type SalesOrderDetailStoreResponse struct {
	SalesOrderDetailStoreRequest
	ProductSKU   string `json:"product_sku,omitempty"`
	ProductName  string `json:"product_name,omitempty"`
	CategoryName string `json:"category_name,omitempty"`
	UomCode      string `json:"uom_code,omitempty"`
}

type SalesOrderDetailUpdateRequest struct {
	SalesOrderDetailTemplate
	ID      int     `json:"id,omitempty" bson:"id,omitempty" binding:"required"`
	BrandID int     `json:"brand_id,omitempty" binding:"required"`
	Price   float64 `json:"price,omitempty"`
}

type SalesOrderDetailChan struct {
	SalesOrderDetail *SalesOrderDetail
	Total            int64
	Error            error
	ErrorLog         *model.ErrorLog
	ID               int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type SalesOrderDetailsChan struct {
	SalesOrderDetails []*SalesOrderDetail
	ErrorLog          *model.ErrorLog
	Total             int64
	Error             error
	ID                int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type SalesOrderDetailRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type SalesOrderDetails struct {
	SalesOrderDetails []*SalesOrderDetail `json:"sales_order_details,omitempty"`
	Total             int64               `json:"total,omitempty"`
}
