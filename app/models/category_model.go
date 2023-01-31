package models

import "time"

type Category struct {
	ID          int        `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Slug        string     `json:"slug,omitempty"`
	Description string     `json:"description,omitempty"`
	IsActive    int        `json:"is_active,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Image       string     `json:"image,omitempty"`
	ParentID    int        `json:"parent_id,omitempty"`
	BrandID     int        `json:"brand_id,omitempty"`
	Order       int        `json:"order,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
