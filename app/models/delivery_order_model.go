package models

import (
	"order-service/global/utils/model"
	"time"
)

type DeliveryOrder struct {
	ID                    int                    `json:"id,omitempty" bson:"id"`
	SalesOrderID          int                    `json:"sales_order_id,omitempty" bson:"sales_order_id"`
	SalesOrder            *SalesOrder            `json:"sales_order,omitempty" bson:"sales_order"`
	SalesOrderCode        NullString             `json:"sales_order_code,omitempty" bson:"sales_order_code"`
	SalesOrderDate        NullString             `json:"sales_order_date,omitempty" bson:"sales_order_date"`
	Salesman              *Salesman              `json:"salesman,omitempty" bson:"salesman"`
	WarehouseID           int                    `json:"warehouse_id,omitempty" bson:"warehouse_id"`
	Warehouse             *Warehouse             `json:"warehouse,omitempty" bson:"warehouse"`
	WarehouseName         string                 `json:"warehouse_name,omitempty" bson:"warehouse_name"`
	WarehouseCode         string                 `json:"warehouse_code,omitempty" bson:"warehouse_code"`
	WarehouseProvinceID   string                 `json:"warehouse_province_id,omitempty" bson:"warehouse_province_id"`
	WarehouseProvinceName NullString             `json:"warehouse_province_name,omitempty" bson:"warehouse_province_name"`
	WarehouseCityID       string                 `json:"warehouse_city_id,omitempty" bson:"warehouse_city_id"`
	WarehouseCityName     NullString             `json:"warehouse_city_name,omitempty" bson:"warehouse_city_name"`
	WarehouseDistrictID   string                 `json:"warehouse_district_id,omitempty" bson:"warehouse_district_id"`
	WarehouseDistrictName NullString             `json:"warehouse_district_name,omitempty" bson:"warehouse_district_name"`
	WarehouseVillageID    string                 `json:"warehouse_village_id,omitempty" bson:"warehouse_village_id"`
	WarehouseVillageName  NullString             `json:"warehouse_village_name,omitempty" bson:"warehouse_village_name"`
	OrderStatusID         int                    `json:"order_status_id,omitempty" bson:"order_status_id"`
	OrderStatus           *OrderStatus           `json:"order_status,omitempty" bson:"order_status"`
	OrderStatusName       NullString             `json:"order_status_name,omitempty" bson:"order_status_name"`
	OrderSourceID         int                    `json:"order_source_id,omitempty" bson:"order_source_id"`
	OrderSource           *OrderSource           `json:"order_source,omitempty" bson:"order_source"`
	OrderSourceName       NullString             `json:"order_source_name,omitempty" bson:"order_source_name"`
	AgentID               int                    `json:"agent_id,omitempty" bson:"agent_id"`
	Agent                 *Agent                 `json:"agent,omitempty" bson:"agent"`
	StoreID               int                    `json:"store_id,omitempty" bson:"store_id"`
	Store                 *Store                 `json:"store,omitempty" bson:"store"`
	DoCode                string                 `json:"do_code,omitempty" bson:"do_code"`
	DoDate                string                 `json:"do_date,omitempty" bson:"do_date"`
	DoRefCode             NullString             `json:"do_ref_code,omitempty" bson:"do_ref_code"`
	DoRefDate             NullString             `json:"do_ref_date,omitempty" bson:"do_ref_date"`
	DriverName            NullString             `json:"driver_name,omitempty" bson:"driver_name"`
	PlatNumber            NullString             `json:"plat_number,omitempty" bson:"plat_number"`
	Note                  NullString             `json:"note,omitempty" bson:"note"`
	DeliveryOrderDetails  []*DeliveryOrderDetail `json:"delivery_order_details,omitempty" bson:"delivery_order_details"`
	IsDoneSyncToEs        string                 `json:"is_done_sync_to_es,omitempty" bson:"is_done_sync_to_es"`
	StartDateSyncToEs     *time.Time             `json:"start_date_sync_to_es,omitempty" bson:"start_date_sync_to_es"`
	EndDateSyncToEs       *time.Time             `json:"end_date_sync_to_es,omitempty" bson:"end_date_sync_to_es"`
	CreatedBy             int                    `json:"created_by,omitempty" bson:"created_by"`
	LatestUpdatedBy       int                    `json:"latest_updated_by" bson:"latest_updated_by"`
	StartCreatedDate      *time.Time             `json:"start_created_date,omitempty" bson:"start_created_date,omitempty"`
	EndCreatedDate        *time.Time             `json:"end_created_date,omitempty" bson:"end_created_date,omitempty"`
	CreatedAt             *time.Time             `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt             *time.Time             `json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt             *time.Time             `json:"deleted_at,omitempty" bson:"deleted_at"`
}

type DeliveryOrderStoreRequest struct {
	RequestID            string                             `json:"request_id,omitempty" bson:"request_id,omitempty"`
	AgentID              int                                `json:"agent_id,omitempty" bson:"agent_id,omitempty" binding:"required"`
	StoreID              int                                `json:"store_id,omitempty" bson:"store_id,omitempty" binding:"required"`
	SalesOrderID         int                                `json:"sales_order_id,omitempty" bson:"sales_order_id,omitempty" binding:"required"`
	WarehouseID          int                                `json:"warehouse_id,omitempty" bson:"warehouse_id,omitempty" binding:"required"`
	OrderSourceID        int                                `json:"order_source_id,omitempty" bson:"order_source_id,omitempty" binding:"required"`
	OrderStatusID        int                                `json:"order_status_id,omitempty" bson:"order_status_id,omitempty" binding:"required"`
	DoCode               string                             `json:"do_code,omitempty" bson:"do_code,omitempty" binding:"required"`
	DoDate               string                             `json:"do_date,omitempty" bson:"do_date,omitempty"`
	DoRefCode            string                             `json:"do_ref_code,omitempty" bson:"do_ref_code,omitempty" binding:"required"`
	DoRefDate            string                             `json:"do_ref_date,omitempty" bson:"do_ref_date,omitempty" binding:"required"`
	DriverName           string                             `json:"driver_name,omitempty" bson:"driver_name,omitempty"`
	PlatNumber           string                             `json:"plat_number,omitempty" bson:"plat_number,omitempty"`
	Note                 string                             `json:"note,omitempty" bson:"note,omitempty"`
	DeliveryOrderDetails []*DeliveryOrderDetailStoreRequest `json:"delivery_order_details,omitempty" bson:"delivery_order_details,omitempty" binding:"required,dive,required"`
}

type DeliveryOrderStoreResponse struct {
	SalesOrderID              int                                 `json:"sales_order_id,omitempty" bson:"sales_order_id,omitempty"`
	SalesOrderSoCode          string                              `json:"sales_order_so_code,omitempty" bson:"sales_order_so_code,omitempty"`
	SalesOrderSoDate          string                              `json:"sales_order_so_date,omitempty" bson:"sales_order_so_date,omitempty"`
	SalesOrderReferralCode    string                              `json:"sales_order_refferal_code,omitempty" bson:"sales_order_refferral_code,omitempty"`
	SalesOrderNote            string                              `json:"sales_order_note,omitempty" bson:"sales_order_note,omitempty"`
	SalesOrderInternalComment string                              `json:"sales_order_internal_comment,omitempty" bson:"sales_order_internal_comment,omitempty"`
	SalesmanName              string                              `json:"salesman_name,omitempty" bson:"salesman_name,omitempty"`
	StoreName                 string                              `json:"store_name,omitempty" bson:"store_name,omitempty"`
	StoreCityName             string                              `json:"store_city_name,omitempty" bson:"store_city_name,omitempty"`
	StoreProvinceName         string                              `json:"store_province_name,omitempty" bson:"store_province_name,omitempty"`
	TotalAmount               int                                 `json:"total_amount,omitempty" bson:"total_amount,omitempty"`
	WarehouseID               int                                 `json:"warehouse_id,omitempty" bson:"warehouse_id,omitempty"`
	WarehouseAddress          string                              `json:"warehouse_address,omitempty" bson:"warehouse_address,omitempty"`
	OrderSourceID             int                                 `json:"order_source_id,omitempty" bson:"order_source_id,omitempty"`
	OrderStatusID             int                                 `json:"order_status_id,omitempty" bson:"order_status_id,omitempty"`
	AgentID                   int                                 `json:"agent_id,omitempty" bson:"agent_id,omitempty"`
	StoreID                   int                                 `json:"store_id,omitempty" bson:"store_id,omitempty"`
	DoCode                    string                              `json:"do_code,omitempty" bson:"do_code,omitempty"`
	DoDate                    string                              `json:"do_date,omitempty" bson:"do_date,omitempty"`
	DoRefCode                 string                              `json:"do_ref_code,omitempty" bson:"do_ref_code,omitempty"`
	DoRefDate                 string                              `json:"do_ref_date,omitempty" bson:"do_ref_date,omitempty"`
	DriverName                string                              `json:"driver_name,omitempty" bson:"driver_name,omitempty"`
	PlatNumber                string                              `json:"plat_number,omitempty" bson:"plat_number,omitempty"`
	Note                      string                              `json:"note,omitempty" bson:"note,omitempty"`
	DeliveryOrderDetails      []*DeliveryOrderDetailStoreResponse `json:"delivery_order_details,omitempty" bson:"delivery_order_details,omitempty"`
}

type DeliveryOrderUpdateByIDRequest struct {
	RequestID            string                                  `json:"request_id,omitempty" bson:"request_id,omitempty"`
	WarehouseID          int                                     `json:"warehouse_id,omitempty" bson:"warehouse_id,omitempty" binding:"required"`
	OrderSourceID        int                                     `json:"order_source_id,omitempty" bson:"order_source_id,omitempty" binding:"required"`
	OrderStatusID        int                                     `json:"order_status_id,omitempty" bson:"order_status_id,omitempty" binding:"required"`
	DoRefCode            string                                  `json:"do_ref_code,omitempty" bson:"do_ref_code,omitempty" binding:"required"`
	DoRefDate            string                                  `json:"do_ref_date,omitempty" bson:"do_ref_date,omitempty" binding:"required"`
	DriverName           string                                  `json:"driver_name,omitempty" bson:"driver_name,omitempty"`
	PlatNumber           string                                  `json:"plat_number,omitempty" bson:"plat_number,omitempty"`
	Note                 string                                  `json:"note,omitempty" bson:"note,omitempty"`
	DeliveryOrderDetails []*DeliveryOrderDetailUpdateByIDRequest `json:"delivery_order_details,omitempty" bson:"delivery_order_details,omitempty" binding:"required,dive,required"`
}

type DeliveryOrderChan struct {
	DeliveryOrder *DeliveryOrder
	Total         int64
	Error         error
	ErrorLog      *model.ErrorLog
	ID            int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type DeliveryOrdersChan struct {
	DeliveryOrders []*DeliveryOrder
	Total          int64
	Error          error
	ErrorLog       *model.ErrorLog
	ID             int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type DeliveryOrderRequest struct {
	ID                int     `json:"id,omitempty"`
	PerPage           int     `json:"per_page,omitempty"`
	Page              int     `json:"page,omitempty"`
	SortField         string  `json:"sort_field,omitempty" bson:"sort_field,omitempty"`
	SortValue         string  `json:"sort_value,omitempty" bson:"sort_value,omitempty"`
	GlobalSearchValue string  `json:"global_search_value,omitempty" bson:"global_search_value,omitempty"`
	Keyword           string  `json:"keyword,omitempty"`
	AgentID           int     `json:"agentID,omitempty"`
	AgentName         string  `json:"agent_name,omitempty"`
	StoreID           int     `json:"storeID,omitempty"`
	StoreName         string  `json:"store_name,omitempty"`
	BrandID           int     `json:"brand_id,omitempty"`
	BrandName         string  `json:"brand_name,omitempty"`
	OrderSourceID     int     `json:"order_source_id,omitempty"`
	OrderStatusID     int     `json:"order_status_id,omitempty"`
	SalesOrderID      int     `json:"sales_order_id,omitempty"`
	SoCode            string  `json:"so_code,omitempty"`
	WarehouseID       int     `json:"warehouse_id,omitempty"`
	WarehouseCode     string  `json:"warehouse_code,omitempty"`
	DoCode            string  `json:"do_code,omitempty"`
	DoDate            string  `json:"do_date,omitempty"`
	DoRefCode         string  `json:"do_ref_code,omitempty"`
	DoRefDate         string  `json:"do_ref_date,omitempty"`
	DoRefferalCode    string  `json:"do_refferal_code,omitempty"`
	TotalAmount       float64 `json:"total_amount,omitempty"`
	TotalTonase       float64 `json:"total_tonase,omitempty"`
	ProductSKU        string  `json:"product_sku,omitempty"`
	ProductCode       string  `json:"product_code,omitempty"`
	ProductName       string  `json:"product_name,omitempty"`
	CategoryID        int     `json:"category_id,omitempty"`
	SalesmanID        int     `json:"salesman_id,omitempty"`
	ProvinceID        int     `json:"province_id,omitempty"`
	CityID            int     `json:"city_id,omitempty"`
	DistrictID        int     `json:"district_id,omitempty"`
	VillageID         int     `json:"village_id,omitempty"`
	StoreProvinceID   int     `json:"store_province_id,omitempty"`
	StoreCityID       int     `json:"store_city_id,omitempty"`
	StoreDistrictID   int     `json:"store_district_id,omitempty"`
	StoreVillageID    int     `json:"store_village_id,omitempty"`
	StoreCode         string  `json:"store_code,omitempty"`
	StartCreatedAt    string  `json:"start_created_at,omitempty"`
	EndCreatedAt      string  `json:"end_created_at,omitempty"`
	StartDoDate       string  `json:"start_do_date,omitempty"`
	EndDoDate         string  `json:"end_do_date,omitempty"`
}

type DeliveryOrderResponse struct {
	SalesOrderID        int                            `json:"sales_order_id,omitempty"`
	WarehouseID         int                            `json:"warehouse_id,omitempty"`
	OrderSourceID       int                            `json:"order_source_id,omitempty"`
	AgentID             int                            `json:"agentID,omitempty"`
	StoreID             int                            `json:"storeID,omitempty"`
	DoCode              string                         `json:"do_code,omitempty"`
	DoDate              string                         `json:"do_date,omitempty"`
	DoRefCode           string                         `json:"do_ref_code,omitempty"`
	DoRefDate           string                         `json:"do_ref_date,omitempty"`
	DriverName          string                         `json:"driver_name,omitempty"`
	PlatNumber          string                         `json:"plat_number,omitempty"`
	Note                string                         `json:"note,omitempty"`
	DeliveryOrderDetail []*DeliveryOrderDetailResponse `json:"delivery_order_details,omitempty"`
}

type DeliveryOrders struct {
	DeliveryOrders []*DeliveryOrder `json:"delivery_orders,omitempty"`
	Total          int64            `json:"total,omitempty"`
}
