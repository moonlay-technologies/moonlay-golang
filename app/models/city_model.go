package models

import "time"

type City struct {
	ID         string     `json:"id,omitempty" bson:"id"`
	ProvinceID string     `json:"province_id,omitempty" bson:"province_id"`
	Name       string     `json:"name,omitempty" bson:"name"`
	CreatedAt  *time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" bson:"deleted_at"`
}
