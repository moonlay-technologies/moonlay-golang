package models

import "order-service/global/utils/model"

type SalesmanStore struct {
	SalesmanID    int    `json:"salesman_id,omitempty"`
	StoreID       int    `json:"store_id,omitempty"`
	AgentID       int    `json:"agent_id,omitempty"`
	IsBlackListed int    `json:"is_black_listed,omitempty"`
	BrandID       int    `json:"brand_id,omitempty"`
	EventSource   string `json:"event_source,omitempty"`
}

type SalesmanStoreChan struct {
	SalesmanStore *SalesmanStore
	Error         error
	ErrorLog      *model.ErrorLog
	ID            int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type SalesmanStoresChan struct {
	SalesmanStores []*SalesmanStore
	Total          int64
	Error          error
	ErrorLog       *model.ErrorLog
	ID             int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type SalesmanStoreRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type SalesmanStores struct {
	SalesmanStores []*SalesmanStore `json:"salesman_stores,omitempty"`
	Total          int64            `json:"total,omitempty"`
}
