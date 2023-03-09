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

type DeliveryOrderDetailLog struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	RequestID    string             `json:"request_id,omitempty" bson:"request_id,omitempty"`
	DoDetailCode string             `json:"do_detail_code,omitempty" bson:"do_detail_code"`
	Data         interface{}        `json:"data,omitempty" bson:"data,omitempty"`
	Error        interface{}        `json:"error,omitempty" bson:"error,omitempty"`
	Action       string             `json:"action,omitempty" bson:"action,omitempty"`
	Status       string             `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt    *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type DeliveryOrderDetailLogChan struct {
	DeliveryOrderDetailLog *DeliveryOrderDetailLog
	Error                  error
	ErrorLog               *model.ErrorLog
	Total                  int64
	ID                     primitive.ObjectID `json:"_id" bson:"_id"`
}

type DeliveryOrderLogsChan struct {
	DeliveryOrderLogs []*DeliveryOrderLog
	Total             int64
	Error             error
	ErrorLog          *model.ErrorLog
	ID                primitive.ObjectID `json:"_id" bson:"_id"`
}

type DeliveryOrderJourney struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	DoId      int                `json:"do_id,omitempty" bson:"do_id,omitempty"`
	DoCode    string             `json:"do_code,omitempty" bson:"do_code,omitempty"`
	DoDate    string             `json:"do_date,omitempty" bson:"do_date,omitempty"`
	Status    string             `json:"status,omitempty" bson:"status,omitempty"`
	Remark    string             `json:"remark,omitempty" bson:"remark,omitempty"`
	Reason    string             `json:"reason,omitempty" bson:"reason,omitempty"`
	CreatedAt *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type DeliveryOrderJourneyChan struct {
	DeliveryOrderJourney *DeliveryOrderJourney
	Error                error
	ErrorLog             *model.ErrorLog
	Total                int64
	ID                   primitive.ObjectID `json:"_id" bson:"_id"`
}
