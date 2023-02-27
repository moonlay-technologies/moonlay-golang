package models

import (
	"order-service/global/utils/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SalesOrderDetailJourneys struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SoDetailId   int                `json:"so_detail_id,omitempty" bson:"so_detail_id,omitempty"`
	SoDetailCode string             `json:"so_detail_code,omitempty" bson:"so_detail_code,omitempty"`
	Status       string             `json:"status,omitempty" bson:"status,omitempty"`
	Remark       string             `json:"remark,omitempty" bson:"remark,omitempty"`
	Reason       string             `json:"reason,omitempty" bson:"reason,omitempty"`
	CreatedAt    *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type SalesOrderDetailJourneysChan struct {
	SalesOrderDetailJourneys *SalesOrderDetailJourneys
	Error                    error
	ErrorLog                 *model.ErrorLog
	Total                    int64
	ID                       primitive.ObjectID `json:"_id" bson:"_id"`
}
