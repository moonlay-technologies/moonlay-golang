package models

import (
	"order-service/global/utils/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DoUploadErrorLog struct {
	ID                primitive.ObjectID      `json:"_id,omitempty" bson:"_id,omitempty"`
	RequestId         string                  `json:"request_id,omitempty" bson:"request_id,omitempty"`
	DoUploadHistoryId primitive.ObjectID      `json:"do_upload_history_id,omitempty" bson:"do_upload_history_id,omitempty"`
	BulkCode          string                  `json:"bulk_code,omitempty" bson:"bulk_code,omitempty"`
	RowData           RowDataDoUploadErrorLog `json:"row_data,omitempty" bson:"row_data,omitempty"`
	ErrorRowLine      int64                   `json:"error_row_line,omitempty" bson:"error_row_line,omitempty"`
	ErrorMessage      string                  `json:"error_message,omitempty" bson:"error_message,omitempty"`
	CreatedAt         *time.Time              `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt         *time.Time              `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type RowDataDoUploadErrorLog struct {
	AgentID       string  `json:"agent_id,omitempty" bson:"agent_id,omitempty"`
	AgentName     *string `json:"agent_name,omitempty" bson:"agent_name,omitempty"`
	OrderCode     string  `json:"order_code,omitempty" bson:"order_code,omitempty"`
	DoDate        string  `json:"do_date,omitempty" bson:"do_date,omitempty"`
	DoNumber      string  `json:"do_number,omitempty" bson:"do_number,omitempty"`
	Note          *string `json:"note,omitempty" bson:"note,omitempty"`
	InternalNote  *string `json:"internal_note,omitempty" bson:"internal_note,omitempty"`
	DriverName    *string `json:"driver_name,omitempty" bson:"driver_name,omitempty"`
	VehicleNo     *string `json:"vehicle_no,omitempty" bson:"vehicle_no,omitempty"`
	BrandID       string  `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	BrandName     *string `json:"brand_name,omitempty" bson:"brand_name,omitempty"`
	ProductCode   string  `json:"product_code,omitempty" bson:"product_code,omitempty"`
	ProductName   *string `json:"product_name,omitempty" bson:"product_name,omitempty"`
	DeliveryQty   *string `json:"delivery_qty,omitempty" bson:"delivery_qty,omitempty"`
	ProductUnit   string  `json:"product_unit,omitempty" bson:"product_unit,omitempty"`
	WarehouseCode string  `json:"warehouse_code,omitempty" bson:"warehouse_code,omitempty"`
	WarehouseName *string `json:"warehouse_name,omitempty" bson:"warehouse_name,omitempty"`
}

type DoUploadErrorLogChan struct {
	DoUploadErrorLog *DoUploadErrorLog
	Error            error
	ErrorLog         *model.ErrorLog
	Total            int64
	ID               primitive.ObjectID `json:"_id" bson:"_id"`
}

type DoUploadErrorLogsChan struct {
	DoUploadErrorLogs []*DoUploadErrorLog
	Error             error
	ErrorLog          *model.ErrorLog
	Total             int64
	ID                primitive.ObjectID `json:"_id" bson:"_id"`
}

type GetDoUploadErrorLogsRequest struct {
	ID                string `json:"_id,omitempty" bson:"_id,omitempty"`
	Page              int    `json:"page,omitempty" bson:"page,omitempty"`
	PerPage           int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	SortField         string `json:"sort_field,omitempty" bson:"sort_field,omitempty"`
	SortValue         string `json:"sort_value,omitempty" bson:"sort_value,omitempty"`
	RequestID         string `json:"request_id,omitempty" bson:"request_id,omitempty"`
	DoUploadHistoryID string `json:"do_upload_history_id,omitempty" bson:"do_upload_history_id,omitempty"`
	Status            string `json:"status,omitempty" bson:"status,omitempty"`
}

type GetDoUploadErrorLogsResponse struct {
	DoUploadErrorLogs []*DoUploadErrorLog `json:"do_upload_error_logs,omitempty"`
	Total             int64               `json:"total,omitempty"`
}
