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
