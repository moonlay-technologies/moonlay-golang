package models

import "time"

type District struct {
	ID        string     `json:"id,omitempty" bson:"id"`
	CityID    string     `json:"city_id,omitempty" bson:"city_id"`
	Name      string     `json:"name,omitempty" bson:"name"`
	CreatedAt *time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bson:"deleted_at"`
}
