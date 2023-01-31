package models

import "poc-order-service/global/utils/model"

type BrandStoreAgent struct {
	ID      int `json:"id,omitempty"`
	BrandID int `json:"brand_id,omitempty"`
	StoreID int `json:"store_id,omitempty"`
	AgentID int `json:"agent_id,omitempty"`
}

type BrandStoreAgentChan struct {
	BrandStoreAgent *BrandStoreAgent
	Error           error
	ErrorLog        *model.ErrorLog
	ID              int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type BrandStoreAgentsChan struct {
	BrandStoreAgents []*BrandStoreAgent
	Total            int64
	Error            error
	ErrorLog         *model.ErrorLog
	ID               int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type BrandStoreAgentRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type BrandStoreAgents struct {
	BrandStoreAgents []*BrandStoreAgent `json:"brand_store_agents,omitempty"`
	Total            int64              `json:"total,omitempty"`
}
