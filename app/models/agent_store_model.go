package models

import "order-service/global/utils/model"

type AgentStore struct {
	ID            int `json:"id,omitempty"`
	AgentID       int `json:"agent_id,omitempty"`
	StoreID       int `json:"store_id,omitempty"`
	SalesmanID    int `json:"salesman_id,omitempty"`
	IsMyStore     int `json:"is_my_store,omitempty"`
	IsBlackListed int `json:"is_black_listed,omitempty"`
}

type AgentStoreChan struct {
	AgentStore *AgentStore
	Error      error
	ErrorLog   *model.ErrorLog
	ID         int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type AgentStoresChan struct {
	AgentStores []*AgentStore
	Total       int64
	Error       error
	ErrorLog    *model.ErrorLog
	ID          int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type AgentStoreRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type AgentStores struct {
	AgentStores []*AgentStore `json:"agent_stores,omitempty"`
	Total       int64         `json:"total,omitempty"`
}
