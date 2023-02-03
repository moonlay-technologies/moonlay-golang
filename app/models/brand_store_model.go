package models

import "order-service/global/utils/model"

type BrandStore struct {
	ID      int `json:"id,omitempty"`
	BrandID int `json:"brand_id,omitempty"`
	StoreID int `json:"store_id,omitempty"`
}

type BrandStoreChan struct {
	BrandStore *BrandStore
	Error      error
	ErrorLog   *model.ErrorLog
	ID         int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type BrandStoresChan struct {
	BrandStores []*BrandStore
	Total       int64
	Error       error
	ErrorLog    *model.ErrorLog
	ID          int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type BrandStoreRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type BrandStores struct {
	BrandStores []*BrandStore `json:"brand_stores,omitempty"`
	Total       int64         `json:"total,omitempty"`
}
