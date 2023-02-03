package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"order-service/global/utils/model"
	"time"
)

type SalesOrderLog struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	RequestID string             `json:"request_id,omitempty" bson:"request_id,omitempty"`
	SoCode    string             `json:"so_code,omitempty" bson:"so_code,omitempty"`
	Data      interface{}        `json:"data,omitempty" bson:"data,omitempty"`
	Status    string             `json:"status,omitempty" bson:"status,omitempty"`
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
