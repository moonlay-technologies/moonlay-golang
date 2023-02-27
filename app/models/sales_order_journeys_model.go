package models

import (
	"order-service/global/utils/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SalesOrderJourneys struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SoId      int                `json:"so_id,omitempty" bson:"so_id,omitempty"`
	SoCode    string             `json:"so_code,omitempty" bson:"so_code,omitempty"`
	SoDate    string             `json:"so_date,omitempty" bson:"so_date,omitempty"`
	Status    string             `json:"status,omitempty" bson:"status,omitempty"`
	Remark    string             `json:"remark,omitempty" bson:"remark,omitempty"`
	Reason    string             `json:"reason,omitempty" bson:"reason,omitempty"`
	CreatedAt *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type SalesOrderJourneysChan struct {
	SalesOrderJourneys *SalesOrderJourneys
	Error              error
	ErrorLog           *model.ErrorLog
	Total              int64
	ID                 primitive.ObjectID `json:"_id" bson:"_id"`
}
