package models

import (
	"order-service/global/utils/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadHistory struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BulkCode       string             `json:"bulk_code,omitempty" bson:"bulk_code,omitempty"`
	FName          string             `json:"f_name,omitempty" bson:"f_name,omitempty"`
	FPath          string             `json:"f_path,omitempty" bson:"f_path,omitempty"`
	AgentId        int                `json:"agent_id,omitempty" bson:"agent_id,omitempty"`
	AgentName      string             `json:"agent_name,omitempty" bson:"agent_name,omitempty"`
	UploadById     int                `json:"upload_by_id,omitempty" bson:"upload_by_id,omitempty"`
	UploadByEmail  string             `json:"upload_by_email,omitempty" bson:"upload_by_email,omitempty"`
	UploadStatus   string             `json:"upload_status,omitempty" bson:"upload_status,omitempty"`
	CurrentProcess string             `json:"current_process,omitempty" bson:"current_process,omitempty"`
	CreatedDate    *time.Time         `json:"created_date,omitempty" bson:"created_date,omitempty"`
	DoneDate       *time.Time         `json:"done_date,omitempty" bson:"done_date,omitempty"`
	TotalRow       string             `json:"total_row,omitempty" bson:"total_row,omitempty"`
	UseMater       string             `json:"use_master,omitempty" bson:"use_master,omitempty"`
	RequestId      string             `json:"request_id,omitempty" bson:"request_id,omitempty"`
	UploadedByName string             `json:"uploaded_by_name,omitempty" bson:"uploaded_by_name,omitempty"`
	UpdatedAt      *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UploadHistoryChan struct {
	UploadHistory *UploadHistory
	Error         error
	ErrorLog      *model.ErrorLog
	Total         int64
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
}
