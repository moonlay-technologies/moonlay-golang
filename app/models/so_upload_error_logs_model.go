package models

import (
	"order-service/global/utils/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SoUploadErrorLog struct {
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SoUploadHistoryId primitive.ObjectID `json:"so_upload_history_id,omitempty" bson:"so_upload_history_id,omitempty"`
	BulkCode          string             `json:"bulk_code,omitempty" bson:"bulk_code,omitempty"`
	Address           string             `json:"address,omitempty" bson:"address,omitempty"`
	AddressId         string             `json:"address_id,omitempty" bson:"address_id,omitempty"`
	AgentName         string             `json:"agent_name,omitempty" bson:"agent_name,omitempty"`
	AgentId           string             `json:"agent_id,omitempty" bson:"agent_id,omitempty"`
	BrandCode         string             `json:"brand_code,omitempty" bson:"brand_code,omitempty"`
	BrandName         string             `json:"brand_name,omitempty" bson:"brand_name,omitempty"`
	SalesId           string             `json:"sales_id,omitempty" bson:"sales_id,omitempty"`
	SalesName         string             `json:"sales_name,omitempty" bson:"sales_name,omitempty"`
	StoreCode         string             `json:"store_code,omitempty" bson:"store_code,omitempty"`
	StoreName         string             `json:"store_name,omitempty" bson:"store_name,omitempty"`
	StoreOrderAt      string             `json:"store_order_at,omitempty" bson:"store_order_at,omitempty"`
	InternalNote      string             `json:"internal_note,omitempty" bson:"internal_note,omitempty"`
	OrderCode         string             `json:"order_code,omitempty" bson:"order_code,omitempty"`
	OrderDate         string             `json:"order_date,omitempty" bson:"order_date,omitempty"`
	OrderNote         string             `json:"order_note,omitempty" bson:"order_note,omitempty"`
	Line              int                `json:"line,omitempty" bson:"line,omitempty"`
	ErrorLog          string             `json:"error_log,omitempty" bson:"error_log,omitempty"`
	CreatedAt         *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt         *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type SoUploadErrorLogChan struct {
	SoUploadErrorLog *SoUploadErrorLog
	Error            error
	ErrorLog         *model.ErrorLog
	Total            int64
	ID               primitive.ObjectID `json:"_id" bson:"_id"`
}
