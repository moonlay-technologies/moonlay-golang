package models

import (
	"order-service/global/utils/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SoUploadErrorLog struct {
	ID                primitive.ObjectID      `json:"_id,omitempty" bson:"_id,omitempty"`
	RequestId         string                  `json:"request_id,omitempty" bson:"request_id,omitempty"`
	SoUploadHistoryId primitive.ObjectID      `json:"so_upload_history_id,omitempty" bson:"so_upload_history_id,omitempty"`
	BulkCode          string                  `json:"bulk_code,omitempty" bson:"bulk_code,omitempty"`
	RowData           RowDataSoUploadErrorLog `json:"row_data,omitempty" bson:"row_data,omitempty"`
	ErrorRowLine      int64                   `json:"error_row_line,omitempty" bson:"error_row_line,omitempty"`
	ErrorMessage      string                  `json:"error_message,omitempty" bson:"error_message,omitempty"`
	CreatedAt         *time.Time              `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt         *time.Time              `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type RowDataSoUploadErrorLog struct {
	AgentId      string  `json:"agent_id,omitempty" bson:"agent_id,omitempty"`
	AgentName    *string `json:"agent_name,omitempty" bson:"agent_name,omitempty"`
	StoreCode    string  `json:"store_code,omitempty" bson:"store_code,omitempty"`
	StoreName    *string `json:"store_name,omitempty" bson:"store_name,omitempty"`
	SalesId      string  `json:"sales_id,omitempty" bson:"sales_id,omitempty"`
	SalesName    *string `json:"sales_name,omitempty" bson:"sales_name,omitempty"`
	SoRefCode    string  `json:"so_ref_code,omitempty" bson:"so_ref_code,omitempty"`
	OrderDate    string  `json:"order_date,omitempty" bson:"order_date,omitempty"`
	OrderNote    string  `json:"order_note,omitempty" bson:"order_note,omitempty"`
	InternalNote string  `json:"internal_note,omitempty" bson:"internal_note,omitempty"`
	BrandCode    string  `json:"brand_code,omitempty" bson:"brand_code,omitempty"`
	BrandName    *string `json:"brand_name,omitempty" bson:"brand_name,omitempty"`
	ProductCode  string  `json:"product_code,omitempty" bson:"product_code,omitempty"`
	ProductName  string  `json:"product_name,omitempty" bson:"product_name,omitempty"`
	OrderQty     string  `json:"order_qty,omitempty" bson:"order_qty,omitempty"`
	ProductUnit  string  `json:"product_unit,omitempty" bson:"product_unit,omitempty"`
	AddresId     string  `json:"address_id,omitempty" bson:"address_id,omitempty"`
	Address      *string `json:"address,omitempty" bson:"address,omitempty"`
}

type SoUploadErrorLogChan struct {
	SoUploadErrorLog *SoUploadErrorLog
	Error            error
	ErrorLog         *model.ErrorLog
	Total            int64
	ID               primitive.ObjectID `json:"_id" bson:"_id"`
}

type SoUploadErrorLogsChan struct {
	SoUploadErrorLogs []*SoUploadErrorLog
	Error             error
	ErrorLog          *model.ErrorLog
	Total             int64
	ID                primitive.ObjectID `json:"_id" bson:"_id"`
}

type GetSoUploadErrorLogsRequest struct {
	ID        string `json:"_id,omitempty" bson:"_id,omitempty"`
	Page      int    `json:"page,omitempty" bson:"page,omitempty"`
	PerPage   int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	SortField string `json:"sort_field,omitempty" bson:"sort_field,omitempty"`
	SortValue string `json:"sort_value,omitempty" bson:"sort_value,omitempty"`
	RequestID string `json:"request_id,omitempty" bson:"request_id,omitempty"`
	Status    string `json:"status,omitempty" bson:"status,omitempty"`
}

type GetSoUploadErrorLogsResponse struct {
	SoUploadErrosLogs []*SoUploadErrorLog `json:"so_upload_error_logs,omitempty"`
	Total             int64               `json:"total,omitempty"`
}
