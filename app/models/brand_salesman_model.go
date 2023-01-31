package models

import "poc-order-service/global/utils/model"

type BrandSalesman struct {
	ID         int `json:"id,omitempty"`
	BrandID    int `json:"brand_id,omitempty"`
	SalesmanID int `json:"salesman_id,omitempty"`
	AgentID    int `json:"agent_id,omitempty"`
}

type BrandSalesmanChan struct {
	BrandSalesman *BrandSalesman
	Error         error
	ErrorLog      *model.ErrorLog
	ID            int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type BrandSalesmansChan struct {
	BrandSalesmans []*BrandSalesman
	Total          int64
	Error          error
	ErrorLog       *model.ErrorLog
	ID             int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type BrandSalesmanRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type BrandSalesmans struct {
	BrandSalesmans []*BrandSalesman `json:"brand_salesmans,omitempty"`
	Total          int64            `json:"total,omitempty"`
}
