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

type SalesOrdersJourneysChan struct {
	SalesOrderJourneys []*SalesOrderJourneys
	Error              error
	ErrorLog           *model.ErrorLog
	Total              int64
	ID                 primitive.ObjectID `json:"_id" bson:"_id"`
}

type SalesOrderJourneyRequest struct {
	Page              int    `json:"page,omitempty" bson:"page,omitempty"`
	PerPage           int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	SortField         string `json:"sort_field,omitempty" bson:"sort_field,omitempty"`
	SortValue         string `json:"sort_value,omitempty" bson:"sort_value,omitempty"`
	GlobalSearchValue string `json:"global_search_value,omitempty" bson:"global_search_value,omitempty"`
	SoId              int    `json:"so_id,omitempty" bson:"so_id,omitempty"`
	SoCode            string `json:"so_code,omitempty" bson:"so_code,omitempty"`
	Status            string `json:"status,omitempty" bson:"status,omitempty"`
	Action            string `json:"action,omitempty" bson:"action,omitempty"`
	StartDate         string `json:"start_date,omitempty" bson:"start_date,omitempty"`
	EndDate           string `json:"end_date,omitempty" bson:"end_date,omitempty"`
}

type SalesOrderJourneyResponse struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	SoId            *int               `json:"so_id,omitempty" bson:"so_id,omitempty"`
	SoCode          *string            `json:"so_code,omitempty" bson:"so_code,omitempty"`
	SoDate          *string            `json:"so_date,omitempty" bson:"so_date,omitempty"`
	OrderStatusID   *int               `json:"order_status_id,omitempty" bson:"order_status_id,omitempty"`
	OrderStatusName *string            `json:"order_status_name,omitempty" bson:"order_status_name,omitempty"`
	Remark          NullString         `json:"remark,omitempty" bson:"remark,omitempty"`
	Reason          NullString         `json:"reason,omitempty" bson:"reason,omitempty"`
	CreatedAt       *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt       *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type SalesOrderJourneyResponses struct {
	SalesOrderJourneys []*SalesOrderJourneyResponse `json:"sales_order_journeys,omitempty"`
	Total              int64                        `json:"total,omitempty"`
}
