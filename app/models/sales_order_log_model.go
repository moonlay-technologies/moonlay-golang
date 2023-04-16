package models

import (
	"order-service/global/utils/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SalesOrderLog struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	RequestID string             `json:"request_id,omitempty" bson:"request_id,omitempty"`
	SoCode    string             `json:"so_code,omitempty" bson:"so_code,omitempty"`
	Data      interface{}        `json:"data,omitempty" bson:"data,omitempty"`
	Status    string             `json:"status,omitempty" bson:"status,omitempty"`
	Error     interface{}        `json:"error,omitempty" bson:"error,omitempty"`
	Action    string             `json:"action,omitempty" bson:"action,omitempty"`
	CreatedAt *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type SalesOrderLogChan struct {
	SalesOrderLog *SalesOrderLog
	Error         error
	ErrorLog      *model.ErrorLog
	Total         int64
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
}

type SalesOrderLogsChan struct {
	SalesOrderLogs []*SalesOrderLog
	Total          int64
	Error          error
	ErrorLog       *model.ErrorLog
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
}

type GetSalesOrderLog struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	RequestID string             `json:"request_id,omitempty" bson:"request_id,omitempty"`
	SoCode    string             `json:"so_code,omitempty" bson:"so_code,omitempty"`
	Data      *SalesOrder        `json:"data,omitempty" bson:"data,omitempty"`
	Status    string             `json:"status,omitempty" bson:"status,omitempty"`
	Error     interface{}        `json:"error,omitempty" bson:"error,omitempty"`
	Action    string             `json:"action,omitempty" bson:"action,omitempty"`
	CreatedAt *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type GetSalesOrderLogsChan struct {
	SalesOrderLogs []*GetSalesOrderLog
	Total          int64
	Error          error
	ErrorLog       *model.ErrorLog
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
}

type GetSalesOrderLogChan struct {
	SalesOrderLog *GetSalesOrderLog
	Total         int64
	Error         error
	ErrorLog      *model.ErrorLog
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
}

type SalesOrderEventLogRequest struct {
	Page              int    `json:"page,omitempty" bson:"page,omitempty"`
	PerPage           int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	SortField         string `json:"sort_field,omitempty" bson:"sort_field,omitempty"`
	SortValue         string `json:"sort_value,omitempty" bson:"sort_value,omitempty"`
	GlobalSearchValue string `json:"global_search_value,omitempty" bson:"global_search_value,omitempty"`
	RequestID         string `json:"request_id,omitempty" bson:"request_id,omitempty"`
	SoCode            string `json:"so_code,omitempty" bson:"so_code,omitempty"`
	Status            string `json:"status,omitempty" bson:"status,omitempty"`
	Action            string `json:"action,omitempty" bson:"action,omitempty"`
	AgentID           int    `json:"agent_id,omitempty" bson:"agent_id,omitempty"`
}

type SalesOrderEventLogResponse struct {
	ID        primitive.ObjectID      `json:"_id,omitempty" bson:"_id,omitempty"`
	RequestID *string                 `json:"request_id,omitempty" bson:"request_id,omitempty"`
	SoCode    string                  `json:"so_code,omitempty" bson:"so_code,omitempty"`
	Data      *DataSOEventLogResponse `json:"data,omitempty" bson:"data,omitempty"`
	Status    string                  `json:"status,omitempty" bson:"status,omitempty"`
	Action    string                  `json:"action,omitempty" bson:"action,omitempty"`
	CreatedAt *time.Time              `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time              `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type SODetailEventLogResponse struct {
	ID           int `json:"id,omitempty" bson:"id,omitempty"`
	SalesOrderID int `json:"sales_order_id,omitempty" bson:"sales_order_id,omitempty"`
	ProductID    int `json:"product_id,omitempty" bson:"product_id,omitempty"`
	OrderQty     int `json:"order_qty,omitempty" bson:"order_qty,omitempty"`
	UomID        int `json:"uom_id,omitempty" bson:"uom_id,omitempty"`
}

type DataSOEventLogResponse struct {
	AgentID           int                         `json:"agent_id,omitempty" bson:"agent_id,omitempty"`
	AgentName         string                      `json:"agent_name,omitempty" bson:"agent_name,omitempty"`
	StoreCode         string                      `json:"store_code,omitempty" bson:"store_code,omitempty"`
	StoreName         string                      `json:"store_name,omitempty" bson:"store_name,omitempty"`
	SalesID           NullInt64                   `json:"sales_id,omitempty" bson:"sales_id,omitempty"`
	SalesName         *string                     `json:"sales_name,omitempty" bson:"sales_name,omitempty"`
	OrderDate         *time.Time                  `json:"order_date,omitempty" bson:"order_date,omitempty"`
	StartOrderAt      *time.Time                  `json:"start_order_at,omitempty" bson:"start_order_at,omitempty"`
	OrderNote         NullString                  `json:"order_note,omitempty" bson:"order_note,omitempty"`
	InternalNote      NullString                  `json:"internal_note,omitempty" bson:"internal_note,omitempty"`
	BrandCode         int                         `json:"brand_code,omitempty" bson:"brand_code,omitempty"`
	BrandName         string                      `json:"brand_name,omitempty" bson:"brand_name,omitempty"`
	SalesOrderDetails []*SODetailEventLogResponse `json:"sales_order_details,omitempty" bson:"sales_order_details,omitempty"`
}
