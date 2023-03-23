package models

import (
	"order-service/global/utils/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SosjUploadErrorLog struct {
	ID                  primitive.ObjectID        `json:"_id,omitempty" bson:"_id,omitempty"`
	RequestId           string                    `json:"request_id,omitempty" bson:"request_id,omitempty"`
	SosjUploadHistoryId primitive.ObjectID        `json:"sosj_upload_history_id,omitempty" bson:"sosj_upload_history_id,omitempty"`
	BulkCode            string                    `json:"bulk_code,omitempty" bson:"bulk_code,omitempty"`
	RowData             RowDataSosjUploadErrorLog `json:"row_data,omitempty" bson:"row_data,omitempty"`
	ErrorRowLine        int64                     `json:"error_row_line,omitempty" bson:"error_row_line,omitempty"`
	ErrorMessage        string                    `json:"error_message,omitempty" bson:"error_message,omitempty"`
	CreatedAt           *time.Time                `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt           *time.Time                `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type RowDataSosjUploadErrorLog struct {
	AgentId      string     `json:"agent_id,omitempty" bson:"agent_id,omitempty"`
	AgentName    NullString `json:"agent_name,omitempty" bson:"agent_name,omitempty"`
	SjStatus     string     `json:"sj_status,omitempty" bson:"sj_status,omitempty"`
	SjNo         string     `json:"sj_no,omitempty" bson:"sj_no,omitempty"`
	SjDate       string     `json:"sj_date,omitempty" bson:"sj_date,omitempty"`
	StoreCode    string     `json:"store_code,omitempty" bson:"store_code,omitempty"`
	StoreName    NullString `json:"store_name,omitempty" bson:"store_name,omitempty"`
	BrandId      string     `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	BrandName    NullString `json:"brand_name,omitempty" bson:"brand_name,omitempty"`
	ProductCode  string     `json:"product_code,omitempty" bson:"product_code,omitempty"`
	ProductName  NullString `json:"product_name,omitempty" bson:"product_name,omitempty"`
	DeliveryQty  string     `json:"deivery_qty,omitempty" bson:"deivery_qty,omitempty"`
	ProductUnit  string     `json:"product_unit,omitempty" bson:"product_unit,omitempty"`
	DriverName   string     `json:"driver_name,omitempty" bson:"driver_name,omitempty"`
	VehicleNo    string     `json:"vehicle_no,omitempty" bson:"vehicle_no,omitempty"`
	WhCode       string     `json:"wh_code,omitempty" bson:"wh_code,omitempty"`
	WhName       NullString `json:"wh_name,omitempty" bson:"wh_name,omitempty"`
	SalesmanId   string     `json:"salesman_id,omitempty" bson:"salesman_id,omitempty"`
	SalesName    NullString `json:"sales_name,omitempty" bson:"sales_name,omitempty"`
	AddresId     string     `json:"address_id,omitempty" bson:"address_id,omitempty"`
	Address      NullString `json:"address,omitempty" bson:"address,omitempty"`
	Note         string     `json:"note,omitempty" bson:"note,omitempty"`
	InternalNote string     `json:"internal_note,omitempty" bson:"internal_note,omitempty"`
}

type SosjUploadErrorLogChan struct {
	SosjUploadErrorLog *SosjUploadErrorLog
	Error              error
	ErrorLog           *model.ErrorLog
	Total              int64
	ID                 primitive.ObjectID `json:"_id" bson:"_id"`
}

type SosjUploadErrorLogsChan struct {
	SosjUploadErrorLogs []*SosjUploadErrorLog
	Error               error
	ErrorLog            *model.ErrorLog
	Total               int64
	ID                  primitive.ObjectID `json:"_id" bson:"_id"`
}

type RowDataSosjUploadErrorLogChan struct {
	RowDataSosjUploadErrorLog *RowDataSosjUploadErrorLog
	Error                     error
	ErrorLog                  *model.ErrorLog
	Total                     int64
	ID                        primitive.ObjectID `json:"_id" bson:"_id"`
}

type GetSosjUploadErrorLogsRequest struct {
	ID        string `json:"_id,omitempty" bson:"_id,omitempty"`
	Page      int    `json:"page,omitempty" bson:"page,omitempty"`
	PerPage   int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	SortField string `json:"sort_field,omitempty" bson:"sort_field,omitempty"`
	SortValue string `json:"sort_value,omitempty" bson:"sort_value,omitempty"`
	RequestID string `json:"request_id,omitempty" bson:"request_id,omitempty"`
	Status    string `json:"status,omitempty" bson:"status,omitempty"`
}

type GetSosjUploadErrorLogsResponse struct {
	SosjUploadErrorLogs []*SosjUploadErrorLog `json:"sosj_upload_error_logs,omitempty"`
	Total               int64                 `json:"total,omitempty"`
}
