package models

import (
	"order-service/global/utils/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type DeliveryOrderJourneys struct {
	DeliveryOrderJourneys []*DeliveryOrderJourney `json:"delivery_order_journeys,omitempty" bson:"delivery_order_journeys,omitempty"`
	Total                 int64                   `json:"total,omitempty" bson:"total,omitempty"`
}

type DeliveryOrderJourneyChan struct {
	DeliveryOrderJourney *DeliveryOrderJourney
	Error                error
	ErrorLog             *model.ErrorLog
	Total                int64
	ID                   primitive.ObjectID `json:"_id" bson:"_id"`
}

type DeliveryOrderJourneysChan struct {
	DeliveryOrderJourney []*DeliveryOrderJourney
	Error                error
	ErrorLog             *model.ErrorLog
	Total                int64
	ID                   primitive.ObjectID `json:"_id" bson:"_id"`
}

type DeliveryOrderJourneysResponse struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	DoId      int                `json:"do_id,omitempty" bson:"do_id,omitempty"`
	DoCode    string             `json:"do_code,omitempty" bson:"do_code,omitempty"`
	DoDate    string             `json:"do_date,omitempty" bson:"do_date,omitempty"`
	Status    string             `json:"status,omitempty" bson:"status,omitempty"`
	Remark    NullString         `json:"remark,omitempty" bson:"remark,omitempty"`
	Reason    NullString         `json:"reason,omitempty" bson:"reason,omitempty"`
	CreatedAt *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type DeliveryOrderJourneysResponses struct {
	DeliveryOrderJourneys []*DeliveryOrderJourneysResponse `json:"delivery_order_journeys,omitempty"`
	Total                 int64                            `json:"total,omitempty"`
}
