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
	SentQty     int    `json:"sent_qty,omitempty"`
	ResidualQty int    `json:"residual_qty,omitempty"`
	Note        string `json:"note,omitempty"`
}

type SalesOrderDetailStoreRequest struct {
	SalesOrderDetailTemplate
	SalesOrderId  int     `json:"sales_order_id,omitempty"`
	OrderStatusId int     `json:"order_status_id,omitempty"`
	SoDetailCode  string  `json:"so_detail_code,omitempty"`
	Price         float64 `json:"price,omitempty" binding:"required"`
}

type SalesOrderDetailStoreResponse struct {
	ID int `json:"id,omitempty" bson:"id,omitempty"`
	SalesOrderDetailStoreRequest
	ProductSKU   string     `json:"product_sku,omitempty"`
	ProductName  string     `json:"product_name,omitempty"`
	CategoryName string     `json:"category_name,omitempty"`
	UomCode      string     `json:"uom_code,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
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

type SalesOrderDetailOpenSearch struct {
	ID                int          `json:"id,omitempty" bson:"id,omitempty"`
	SoDetailCode      string       `json:"so_detail_code,omitempty" bson:"so_detail_code,omitempty"`
	Qty               int          `json:"qty,omitempty" bson:"qty,omitempty"`
	SentQty           int          `json:"sent_qty,omitempty" bson:"sent_qty,omitempty"`
	ResidualQty       int          `json:"residual_qty,omitempty" bson:"residual_qty,omitempty"`
	Price             float64      `json:"price,omitempty" bson:"price,omitempty"`
	Subtotal          float64      `json:"subtotal,omitempty" bson:"subtotal,omitempty"`
	Note              NullString   `json:"note,omitempty" bson:"note,omitempty"`
	SalesOrderID      int          `json:"sales_order_id,omitempty" bson:"sales_order_id,omitempty"`
	SoCode            string       `json:"so_code,omitempty" bson:"so_code,omitempty"`
	SoDate            string       `json:"so_date,omitempty" bson:"so_date,omitempty"`
	SoRefCode         string       `json:"so_ref_code,omitempty" bson:"so_ref_code,omitempty"`
	SoRefDate         NullString   `json:"so_ref_date,omitempty" bson:"so_ref_date,omitempty"`
	SoReferralCode    NullString   `json:"so_referral_code,omitempty" bson:"so_referral_code,omitempty"`
	AgentId           int          `json:"agent_id,omitempty" bson:"agent_id,omitempty"`
	Agent             *Agent       `json:"agent,omitempty" bson:"agent,omitempty"`
	StoreID           int          `json:"store_id,omitempty" bson:"store_id,omitempty"`
	Store             *Store       `json:"store,omitempty" bson:"store,omitempty"`
	BrandID           int          `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	BrandName         string       `json:"brand_name,omitempty" bson:"brand_name,omitempty"`
	Brand             *Brand       `json:"brand,omitempty" bson:"brand,omitempty"`
	UserID            int          `json:"user_id,omitempty" bson:"user_id,omitempty"`
	User              *User        `json:"user,omitempty" bson:"user,omitempty"`
	SalesmanID        int          `json:"salesman_id,omitempty" bson:"salesman_id,omitempty"`
	Salesman          *Salesman    `json:"salesman,omitempty" bson:"salesman,omitempty"`
	OrderSourceID     int          `json:"order_source_id,omitempty" bson:"order_source_id,omitempty"`
	OrderSource       *OrderSource `json:"order_source,omitempty" bson:"order_source,omitempty"`
	OrderStatusID     int          `json:"order_status_id,omitempty" bson:"order_status_id,omitempty"`
	OrderStatus       *OrderStatus `json:"order_status,omitempty" bson:"order_status,omitempty"`
	GLat              NullFloat64  `json:"g_lat,omitempty" bson:"g_lat,omitempty"`
	GLong             NullFloat64  `json:"g_long,omitempty" bson:"g_long,omitempty"`
	ProductID         int          `json:"product_id,omitempty" bson:"product_id,omitempty" `
	Product           *Product     `json:"product,omitempty" bson:"product,omitempty"`
	UomID             int          `json:"uom_id,omitempty" bson:"uom_id,omitempty"`
	Uom               *Uom         `json:"uom,omitempty" bson:"uom,omitempty"`
	IsDoneSyncToEs    string       `json:"is_done_sync_to_es,omitempty" bson:"is_done_sync_to_es,omitempty"`
	StartDateSyncToEs *time.Time   `json:"start_date_sync_to_es,omitempty" bson:"start_date_sync_to_es,omitempty"`
	EndDateSyncToEs   *time.Time   `json:"end_date_sync_to_es,omitempty" bson:"end_date_sync_to_es,omitempty"`
	CreatedAt         *time.Time   `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt         *time.Time   `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt         *time.Time   `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type SalesOrderDetailOpenSearchResponse struct {
	ID            int                            `json:"id,omitempty" bson:"id,omitempty"`
	SalesOrderID  int                            `json:"sales_order_id,omitempty" bson:"sales_order_id,omitempty"`
	ProductID     int                            `json:"product_id,omitempty" bson:"product_id,omitempty" `
	BrandID       int                            `json:"brand_id,omitempty" bson:"brand_id,omitempty" `
	Product       *ProductOpenSearchResponse     `json:"product,omitempty" bson:"product,omitempty"`
	UomID         int                            `json:"uom_id,omitempty" bson:"uom_id,omitempty"`
	Uom           *UomOpenSearchResponse         `json:"uom,omitempty" bson:"uom,omitempty"`
	OrderStatusID int                            `json:"order_status_id,omitempty" bson:"order_status_id,omitempty"`
	OrderStatus   *OrderStatusOpenSearchResponse `json:"order_status,omitempty" bson:"order_status,omitempty"`
	SoDetailCode  string                         `json:"so_detail_code,omitempty" bson:"so_detail_code,omitempty"`
	Qty           NullInt64                      `json:"qty,omitempty" bson:"qty,omitempty"`
	SentQty       NullInt64                      `json:"sent_qty,omitempty" bson:"sent_qty,omitempty"`
	ResidualQty   NullInt64                      `json:"residual_qty,omitempty" bson:"residual_qty,omitempty"`
	Price         float64                        `json:"price,omitempty" bson:"price,omitempty"`
	Note          NullString                     `json:"note,omitempty" bson:"note,omitempty"`
	CreatedAt     *time.Time                     `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type SalesOrderDetailsOpenSearchResponse struct {
	SalesOrderDetails []*SalesOrderDetailOpenSearchResponse `json:"sales_order_details,omitempty"`
	Total             int64                                 `json:"total,omitempty"`
}
