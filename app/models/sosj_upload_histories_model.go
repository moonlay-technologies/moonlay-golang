package models

import (
	"order-service/global/utils/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadHistory struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	RequestId       string             `json:"request_id,omitempty" bson:"request_id,omitempty"`
	BulkCode        string             `json:"bulk_code,omitempty" bson:"bulk_code,omitempty"`
	FileName        string             `json:"file_name,omitempty" bson:"file_name,omitempty"`
	FilePath        string             `json:"file_path,omitempty" bson:"file_path,omitempty"`
	AgentId         *int64             `json:"agent_id,omitempty" bson:"agent_id,omitempty"`
	AgentName       string             `json:"agent_name,omitempty" bson:"agent_name,omitempty"`
	UploadedBy      *int64             `json:"uploaded_by,omitempty" bson:"uploaded_by,omitempty"`
	UploadedByName  string             `json:"uploaded_by_name,omitempty" bson:"uploaded_by_name,omitempty"`
	UploadedByEmail string             `json:"uploaded_by_email,omitempty" bson:"uploaded_by_email,omitempty"`
	UpdatedBy       *int64             `json:"updated_by,omitempty" bson:"updated_by,omitempty"`
	UpdatedByName   string             `json:"updated_by_name,omitempty" bson:"updated_by_name,omitempty"`
	UpdatedByEmail  string             `json:"updated_by_email,omitempty" bson:"updated_by_email,omitempty"`
	Status          string             `json:"status,omitempty" bson:"status,omitempty"`
	TotalRows       *int64             `json:"total_rows,omitempty" bson:"total_rows,omitempty"`
	CreatedAt       time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt       time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UploadHistoryChan struct {
	UploadHistory *UploadHistory
	Error         error
	ErrorLog      *model.ErrorLog
	Total         int64
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
}
