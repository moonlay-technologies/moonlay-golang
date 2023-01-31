package models

import "poc-order-service/global/utils/model"

type SalesmanBrand struct {
	SalesmanID int `json:"salesman_id,omitempty"`
	BrandID    int `json:"brand_id,omitempty"`
}

type SalesmanBrandChan struct {
	SalesmanBrand *SalesmanBrand
	Error         error
	ErrorLog      *model.ErrorLog
	ID            int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type SalesmanBrandsChan struct {
	SalesmanBrands []*SalesmanBrand
	Total          int64
	Error          error
	ErrorLog       *model.ErrorLog
	ID             int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type SalesmanBrandRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type SalesmanBrands struct {
	SalesmanBrands []*SalesmanBrand `json:"salesman_brands,omitempty"`
	Total          int64            `json:"total,omitempty"`
}
