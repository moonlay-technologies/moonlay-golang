package models

import (
	"order-service/global/utils/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeliveryOrderLog struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	RequestID string             `json:"request_id,omitempty" bson:"request_id,omitempty"`
	DoCode    string             `json:"do_code,omitempty" bson:"do_code,omitempty"`
	Data      interface{}        `json:"data,omitempty" bson:"data,omitempty"`
	Error     interface{}        `json:"error,omitempty" bson:"error,omitempty"`
	Action    string             `json:"action,omitempty" bson:"action,omitempty"`
	Status    string             `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type DeliveryOrderLogChan struct {
	DeliveryOrderLog *DeliveryOrderLog
	Error            error
	ErrorLog         *model.ErrorLog
	Total            int64
	ID               primitive.ObjectID `json:"_id" bson:"_id"`
}

type DeliveryOrderLogsChan struct {
	DeliveryOrderLogs []*DeliveryOrderLog
	Total             int64
	Error             error
	ErrorLog          *model.ErrorLog
	ID                primitive.ObjectID `json:"_id" bson:"_id"`
}

type GetDeliveryOrderLog struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	RequestID string             `json:"request_id,omitempty" bson:"request_id,omitempty"`
	DoCode    string             `json:"do_code,omitempty" bson:"do_code,omitempty"`
	Data      *DeliveryOrder     `json:"data,omitempty" bson:"data,omitempty"`
	Error     interface{}        `json:"error,omitempty" bson:"error,omitempty"`
	Action    string             `json:"action,omitempty" bson:"action,omitempty"`
	Status    string             `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type GetDeliveryOrderLogChan struct {
	DeliveryOrderLog *GetDeliveryOrderLog
	Total            int64
	Error            error
	ErrorLog         *model.ErrorLog
	ID               primitive.ObjectID `json:"_id" bson:"_id"`
}

type GetDeliveryOrderLogsChan struct {
	DeliveryOrderLog []*GetDeliveryOrderLog
	Total            int64
	Error            error
	ErrorLog         *model.ErrorLog
	ID               primitive.ObjectID `json:"_id" bson:"_id"`
}

type DeliveryOrderEventLogRequest struct {
	Page              int    `json:"page,omitempty" bson:"page,omitempty"`
	PerPage           int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	SortField         string `json:"sort_field,omitempty" bson:"sort_field,omitempty"`
	SortValue         string `json:"sort_value,omitempty" bson:"sort_value,omitempty"`
	GlobalSearchValue string `json:"global_search_value,omitempty" bson:"global_search_value,omitempty"`
	ID                string `json:"id,omitempty" bson:"id,omitempty"`
	RequestID         string `json:"request_id,omitempty" bson:"request_id,omitempty"`
	AgentID           int    `json:"agent_id,omitempty" bson:"agent_id,omitempty"`
	Status            string `json:"status,omitempty" bson:"status,omitempty"`
}

type DeliveryOrderEventLogResponse struct {
	ID        primitive.ObjectID      `json:"_id,omitempty" bson:"_id,omitempty"`
	RequestID string                  `json:"request_id,omitempty" bson:"request_id,omitempty"`
	DoID      int                     `json:"do_id,omitempty" bson:"do_id,omitempty"`
	DoCode    string                  `json:"do_code,omitempty" bson:"do_code,omitempty"`
	Data      *DataDOEventLogResponse `json:"data,omitempty" bson:"data,omitempty"`
	Status    string                  `json:"status,omitempty" bson:"status,omitempty"`
	Action    string                  `json:"action,omitempty" bson:"action,omitempty"`
	CreatedAt *time.Time              `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time              `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type DataDOEventLogResponse struct {
	AgentID              int                         `json:"agent_id,omitempty" bson:"agent_id,omitempty"`
	AgentName            string                      `json:"agent_name,omitempty" bson:"agent_name,omitempty"`
	SoCode               string                      `json:"so_code,omitempty" bson:"so_code,omitempty"`
	DoDate               string                      `json:"do_date,omitempty" bson:"do_date,omitempty"`
	DoRefCode            string                      `json:"do_ref_code,omitempty" bson:"do_ref_code,omitempty"`
	Note                 NullString                  `json:"note,omitempty" bson:"note,omitempty"`
	InternalComment      NullString                  `json:"internal_comment,omitempty" bson:"internal_comment,omitempty"`
	DriverName           NullString                  `json:"driver_name,omitempty" bson:"driver_name,omitempty"`
	PlatNumber           NullString                  `json:"plat_number,omitempty" bson:"plat_number,omitempty"`
	BrandID              int                         `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	BrandName            string                      `json:"brand_name,omitempty" bson:"brand_name,omitempty"`
	WarehouseCode        string                      `json:"warehouse_code,omitempty" bson:"warehouse_code,omitempty"`
	WarehouseName        string                      `json:"warehouse_name,omitempty" bson:"warehouse_name,omitempty"`
	DeliveryOrderDetails []*DODetailEventLogResponse `json:"delivery_order_details,omitempty" bson:"delivery_order_details,omitempty"`
}

type DODetailEventLogResponse struct {
	ID          int        `json:"id,omitempty" bson:"id,omitempty"`
	ProductCode string     `json:"product_code,omitempty" bson:"product_code,omitempty"`
	ProductName string     `json:"product_name,omitempty" bson:"product_name,omitempty"`
	DeliveryQty NullInt64  `json:"delivery_qty,omitempty" bson:"delivery_qty,omitempty"`
	ProductUnit NullString `json:"product_unit,omitempty" bson:"product_unit,omitempty"`
}
