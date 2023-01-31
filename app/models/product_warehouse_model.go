package models

import "poc-order-service/global/utils/model"

type ProductWarehouse struct {
	ID             int `json:"id,omitempty"`
	ProductAgentID int `json:"product_agent_id,omitempty"`
	WarehouseID    int `json:"warehouse_id,omitempty"`
	StockBig       int `json:"stock_big,omitempty"`
	StockMedium    int `json:"stock_medium,omitempty"`
	StockSmall     int `json:"stock_small,omitempty"`
}

type ProductWarehouseChan struct {
	AgentBrand *AgentBrand
	Error      error
	ErrorLog   *model.ErrorLog
	ID         int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type ProductWarehousesChan struct {
	ProductWarehouses []*ProductWarehouse
	Total             int64
	Error             error
	ErrorLog          *model.ErrorLog
	ID                int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type ProductWarehouseRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type ProductWarehouses struct {
	ProductWarehouses []*ProductWarehouse `json:"product_warehouses,omitempty"`
	Total             int64               `json:"total,omitempty"`
}
