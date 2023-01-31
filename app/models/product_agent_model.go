package models

import (
	"poc-order-service/global/utils/model"
	"time"
)

type ProductAgent struct {
	ID              int        `json:"id,omitempty"`
	MasterProductID int        `json:"master_product_id,omitempty"`
	AgentID         int        `json:"agent_id,omitempty"`
	PriceBig        float64    `json:"price_big,omitempty"`
	PriceMedium     float64    `json:"price_medium,omitempty"`
	PriceSmall      float64    `json:"price_small,omitempty"`
	AliasName       string     `json:"alias_name,omitempty"`
	AliasSKU        string     `json:"alias_sku,omitempty"`
	IsActive        int        `json:"is_active,omitempty"`
	StartDate       string     `json:"start_date,omitempty"`
	EndDate         string     `json:"end_date,omitempty"`
	PieceInBox      int        `json:"piece_in_box,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

type ProductAgentChan struct {
	ProductAgent *ProductAgent
	Error        error
	ErrorLog     *model.ErrorLog
	ID           int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type ProductAgentsChan struct {
	ProductAgents []*ProductAgent
	Total         int64
	Error         error
	ErrorLog      *model.ErrorLog
	ID            int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type ProductAgentRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type ProductAgents struct {
	ProductAgents []*ProductAgent `json:"product_agents,omitempty"`
	Total         int64           `json:"total,omitempty"`
}
