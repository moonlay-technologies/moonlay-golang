package models

import (
	"poc-order-service/global/utils/model"
	"time"
)

type Product struct {
	ID                    int        `json:"id,omitempty" bson:"id,omitempty"`
	Sku                   NullString `json:"sku,omitempty" bson:"sku,omitempty"`
	AliasSku              NullString `json:"alias_sku,omitempty" bson:"alias_sku,omitempty"`
	ProductName           NullString `json:"product_name,omitempty" bson:"product_name,omitempty"`
	Description           NullString `json:"description,omitempty" bson:"description,omitempty"`
	Uid                   NullString `json:"uid,omitempty" bson:"uid,omitempty"`
	UnitMeasurementSmall  NullString `json:"unit_measurement_small,omitempty" bson:"unit_measurement_small,omitempty"`
	UnitMeasurementMedium NullString `json:"unit_measurement_medium,omitempty" bson:"unit_measurement_medium,omitempty"`
	UnitMeasurementBig    NullString `json:"unit_measurement_big,omitempty" bson:"unit_measurement_big,omitempty"`
	Ukuran                NullString `json:"ukuran,omitempty" bson:"ukuran,omitempty"`
	NettWeight            float64    `json:"nett_weight,omitempty" bson:"nett_weight,omitempty"`
	NettWeightUm          NullString `json:"nett_weight_um,omitempty" bson:"nett_weight_um,omitempty"`
	Volume                float64    `json:"volume,omitempty" bson:"volume,omitempty"`
	SmallInMediumAmount   int        `json:"small_in_medium_amount,omitempty" bson:"small_in_medium_amount,omitempty"`
	MediumInBigAmount     int        `json:"medium_in_big_amount,omitempty" bson:"medium_in_big_amount,omitempty"`
	PriceBig              float64    `json:"price_big,omitempty" bson:"price_big,omitempty"`
	PriceSmall            float64    `json:"price_small,omitempty" bson:"price_small,omitempty"`
	PriceMedium           float64    `json:"price_medium,omitempty" bson:"price_medium,omitempty"`
	Currency              NullString `json:"currency,omitempty" bson:"currency,omitempty"`
	Priority              int        `json:"priority,omitempty" bson:"priority,omitempty"`
	IsActive              int        `json:"is_active,omitempty" bson:"is_active,omitempty"`
	DataType              NullString `json:"data_type,omitempty" bson:"data_type,omitempty"`
	StartDate             *time.Time `json:"start_date,omitempty" bson:"start_date,omitempty"`
	EndDate               *time.Time `json:"end_date,omitempty" bson:"end_date,omitempty"`
	Stock                 float64    `json:"stock,omitempty" bson:"stock,omitempty"`
	Image                 NullString `json:"image,omitempty" bson:"image,omitempty"`
	CategoryID            int        `json:"category_id" bson:"category_id,omitempty"`
	CreatedBy             int        `json:"created_by,omitempty" bson:"created_by,omitempty"`
	UpdatedBy             int        `json:"updated_by,omitempty" bson:"updated_by,omitempty"`
	CreatedAt             *time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt             *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt             *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type ProductChan struct {
	Product  *Product
	Error    error
	ErrorLog *model.ErrorLog
	Total    int64
	ID       int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type ProductsChan struct {
	Products []*Product
	Total    int64
	Error    error
	ErrorLog *model.ErrorLog
	ID       int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type ProductRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type Products struct {
	Products []*Product `json:"products,omitempty"`
	Total    int64      `json:"total,omitempty"`
}
