package models

import "poc-order-service/global/utils/model"

type AgentBrand struct {
	AgentID int `json:"agent_id,omitempty"`
	BrandID int `json:"brand_id,omitempty"`
}

type AgentBrandChan struct {
	AgentBrand *AgentBrand
	Total      int64
	Error      error
	ErrorLog   *model.ErrorLog
	ID         int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type AgentBrandsChan struct {
	AgentBrands []*AgentBrand
	Total       int64
	Error       error
	ErrorLog    *model.ErrorLog
	ID          int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type AgentBrandRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type AgentBrands struct {
	AgentBrands []*AgentBrand `json:"agent_brands,omitempty"`
	Total       int64         `json:"total,omitempty"`
}
