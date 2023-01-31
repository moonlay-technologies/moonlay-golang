package models

import "time"

type Village struct {
	ID         string     `json:"id,omitempty"`
	DistrictID string     `json:"district_id,omitempty"`
	Name       string     `json:"name,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}
