package models

import (
	"order-service/global/utils/model"
	"time"
)

type DeliveryOrderDetail struct {
	ID                int               `json:"id,omitempty" bson:"id"`
	DeliveryOrderID   int               `json:"delivery_order_id,omitempty" bson:"delivery_order_id"`
	SoDetailID        int               `json:"so_detail_id,omitempty" bson:"so_detail_id"`
	SoDetail          *SalesOrderDetail `json:"so_detail,omitempty" bson:" so_detail"`
	BrandID           int               `json:"brand_id,omitempty" bson:"brand_id"`
	Brand             *Brand            `json:"brand,omitempty" bson:"brand"`
	BrandName         string            `json:"brand_name,omitempty"  bson:"brand_name"`
	ProductID         int               `json:"product_id,omitempty" bson:"product_id"`
	Product           *Product          `json:"product,omitempty" bson:"product"`
	ProductSKU        string            `json:"product_sku,omitempty" bson:"product_sku,omitempty"`
	ProductName       string            `json:"product_name,omitempty" bson:"product_name,omitempty"`
	UomID             int               `json:"uom_id,omitempty" bson:"uom_id"`
	Uom               *Uom              `json:"uom,omitempty" bson:"uom"`
	UomCode           string            `json:"uom_code,omitempty" bson:"uom_code,omitempty"`
	UomName           string            `json:"uom_name,omitempty" bson:"uom_name,omitempty"`
	OrderStatusID     int               `json:"order_status_id,omitempty" bson:"order_status_id"`
	OrderStatus       *OrderStatus      `json:"order_status,omitempty" bson:"order_status"`
	OrderStatusName   string            `json:"order_status_name,omitempty" bson:"order_status_name"`
	DoDetailCode      string            `json:"do_detail_code,omitempty" bson:"do_detail_code"`
	Qty               int               `json:"qty,omitempty" bson:"qty"`
	SentQty           int               `json:"sent_qty,omitempty" bson:"sent_qty,omitempty"`
	ResidualQty       int               `json:"residual_qty,omitempty" bson:"residual_qty,omitempty"`
	Price             float64           `json:"price,omitempty" bson:"price,omitempty"`
	Note              NullString        `json:"note,omitempty" bson:"note"`
	IsDoneSyncToEs    string            `json:"is_done_sync_to_es,omitempty" bson:"is_done_sync_to_es"`
	StartDateSyncToEs *time.Time        `json:"start_date_sync_to_es,omitempty" bson:"start_date_sync_to_es"`
	EndDateSyncToEs   *time.Time        `json:"end_date_sync_to_es,omitempty" bson:"end_date_sync_to_es"`
	CreatedAt         *time.Time        `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt         *time.Time        `json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt         *time.Time        `json:"deleted_at,omitempty" bson:"deleted_at"`
}

type DeliveryOrderDetailOpenSearch struct {
	ID                int               `json:"id,omitempty" bson:"id"`
	DeliveryOrderID   int               `json:"delivery_order_id,omitempty" bson:"delivery_order_id"`
	DoCode            string            `json:"do_code,omitempty"`
	DoDate            string            `json:"do_date,omitempty"`
	DoRefCode         string            `json:"do_ref_code,omitempty"`
	DoRefDate         string            `json:"do_ref_date,omitempty"`
	SalesOrderID      int               `json:"sales_order_id,omitempty"`
	SalesOrderCode    NullString        `json:"sales_order_code,omitempty" bson:"sales_order_code"`
	SalesOrderDate    NullString        `json:"sales_order_date,omitempty" bson:"sales_order_date"`
	SoRefDate         NullString        `json:"so_ref_date,omitempty" bson:"so_ref_date,omitempty"`
	SoDetailID        int               `json:"so_detail_id,omitempty" bson:"so_detail_id"`
	SoDetailCode      string            `json:"so_detail_code,omitempty" bson:"so_detail_code,omitempty"`
	SoDetail          *SalesOrderDetail `json:"so_detail,omitempty" bson:" so_detail"`
	AgentID           int               `json:"agentID,omitempty"`
	AgentName         string            `json:"agent_name,omitempty"`
	AgentAddress      NullString        `json:"agent_address,omitempty" bson:"agent_address,omitempty"`
	AgentPhone        NullString        `json:"agent_phone,omitempty" bson:"agent_phone,omitempty"`
	AgentProvinceID   int               `json:"agent_province_id,omitempty" bson:"agent_province_id,omitempty"`
	AgentProvinceName NullString        `json:"agent_province_name,omitempty" bson:"agent_province_name,omitempty"`
	AgentCityID       int               `json:"agent_city_id,omitempty" bson:"agent_city_id,omitempty"`
	AgentCityName     NullString        `json:"agent_city_name,omitempty" bson:"agent_city_name,omitempty"`
	AgentDistrictID   int               `json:"agent_district_id,omitempty" bson:"agent_district_id,omitempty"`
	AgentDistrictName NullString        `json:"agent_district_name,omitempty" bson:"agent_district_name,omitempty"`
	AgentVillageID    int               `json:"agent_village_id,omitempty" bson:"agent_village_id,omitempty"`
	AgentVillageName  NullString        `json:"agent_village_name,omitempty" bson:"agent_village_name,omitempty"`
	StoreID           int               `json:"store_id,omitempty" bson:"store_id,omitempty"`
	StoreName         NullString        `json:"store_name,omitempty" bson:"store_name,omitempty"`
	StoreCode         NullString        `json:"store_code,omitempty" bson:"store_code,omitempty"`
	StoreAddress      NullString        `json:"store_address,omitempty" bson:"store_address,omitempty"`
	StorePhone        NullString        `json:"store_phone,omitempty" bson:"store_phone,omitempty"`
	StoreProvinceID   int               `json:"store_province_id,omitempty" bson:"store_province_id,omitempty"`
	StoreProvinceName NullString        `json:"store_province_name,omitempty" bson:"store_province_name,omitempty"`
	StoreCityID       int               `json:"store_city_id,omitempty" bson:"store_city_id,omitempty"`
	StoreCityName     NullString        `json:"store_city_name,omitempty" bson:"store_city_name,omitempty"`
	StoreDistrictID   int               `json:"store_district_id,omitempty" bson:"store_district_id,omitempty"`
	StoreDistrictName NullString        `json:"store_district_name,omitempty" bson:"store_district_name,omitempty"`
	StoreVillageID    int               `json:"store_village_id,omitempty" bson:"store_village_id,omitempty"`
	StoreVillageName  NullString        `json:"store_village_name,omitempty" bson:"store_village_name,omitempty"`
	WarehouseID       int               `json:"warehouse_id,omitempty" bson:"warehouse_id"`
	WarehouseCode     string            `json:"warehouse_code,omitempty" bson:"warehouse_code"`
	WarehouseName     string            `json:"warehouse_name,omitempty" bson:"warehouse_name"`
	SalesmanID        int               `json:"salesman_id,omitempty" bson:"salesman_id,omitempty"`
	SalesmanName      string            `json:"salesman_name,omitempty" bson:"salesman_name,omitempty"`
	BrandID           int               `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	BrandName         string            `json:"brand_name,omitempty" bson:"brand_name,omitempty"`
	ProductID         int               `json:"product_id,omitempty" bson:"product_id"`
	ProductSKU        string            `json:"product_sku,omitempty" bson:"product_sku,omitempty"`
	ProductName       string            `json:"product_name,omitempty" bson:"product_name,omitempty"`
	Description       NullString        `json:"product_description,omitempty" bson:"product_description,omitempty"`
	UomID             int               `json:"uom_id,omitempty" bson:"uom_id"`
	UomCode           string            `json:"uom_code,omitempty" bson:"uom_code,omitempty"`
	UomName           string            `json:"uom_name,omitempty" bson:"uom_name,omitempty"`
	DoDetailCode      string            `json:"do_detail_code,omitempty" bson:"do_detail_code"`
	OrderStatusID     int               `json:"order_status_id,omitempty" bson:"order_status_id"`
	OrderStatusName   string            `json:"order_status_name,omitempty" bson:"order_status_name"`
	Qty               int               `json:"qty,omitempty" bson:"qty"`
	Note              NullString        `json:"note,omitempty" bson:"note"`
	CreatedAt         *time.Time        `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt         *time.Time        `json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt         *time.Time        `json:"deleted_at,omitempty" bson:"deleted_at"`
}

type DeliveryOrderDetailStoreRequest struct {
	SoDetailID int    `json:"so_detail_id,omitempty" bson:"so_detail_id,omitempty" binding:"required"`
	Qty        int    `json:"qty,omitempty" bson:"qty,omitempty"`
	Note       string `json:"note,omitempty" bson:"note,omitempty"`
}

type DeliveryOrderDetailStoreResponse struct {
	DeliveryOrderID int    `json:"delivery_order_id,omitempty" bson:"delivery_order_id,omitempty"`
	OrderStatusID   int    `json:"order_status_id,omitempty" bson:"order_status_id,omitempty"`
	SoDetailID      int    `json:"so_detail_id,omitempty" bson:"so_detail_id,omitempty"`
	ProductSku      string `json:"product_sku,omitempty" bson:"product_sku,omitempty"`
	ProductName     string `json:"product_name,omitempty" bson:"product_name,omitempty"`
	SalesOrderQty   int    `json:"sales_order_qty,omitempty" bson:"sales_order_qty,omitempty"`
	SentQty         int    `json:"sent_qty,omitempty" bson:"sent_qty,omitempty"`
	ResidualQty     int    `json:"residual_qty,omitempty" bson:"residual_qty,omitempty"`
	UomCode         string `json:"uom_code,omitempty" bson:"uom_code,omitempty"`
	Price           int    `json:"price,omitempty" bson:"price,omitempty"`
	Qty             int    `json:"qty,omitempty" bson:"qty,omitempty"`
	Note            string `json:"note,omitempty" bson:"note,omitempty"`
}

type DeliveryOrderDetailUpdateByIDRequest struct {
	RequestID string `json:"request_id,omitempty" bson:"request_id,omitempty"`
	ID        int    `json:"id,omitempty" bson:"id"`
	Qty       int    `json:"qty,omitempty" bson:"qty,omitempty"`
	Note      string `json:"note,omitempty" bson:"note,omitempty"`
}

type DeliveryOrderDetailUpdateByDeliveryOrderIDRequest struct {
	RequestID string `json:"request_id,omitempty" bson:"request_id,omitempty"`
	ID        int    `json:"id,omitempty" bson:"id,omitempty" binding:"required"`
	Qty       int    `json:"qty,omitempty" bson:"qty,omitempty" binding:"required"`
	Note      string `json:"note,omitempty" bson:"note,omitempty"`
}

type DeliveryOrderDetailUpdateByDeliveryOrdersIDRequest struct {
	DeliveryOrderDetailUpdateByDeliveryOrdersIDRequest []*DeliveryOrderDetailUpdateByDeliveryOrderIDRequest `json:"delivery_order_detail_update_by_do_id,omitempty"`
}

type DeliveryOrderDetailChan struct {
	DeliveryOrderDetail *DeliveryOrderDetail
	Error               error
	ErrorLog            *model.ErrorLog
	ID                  int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type DeliveryOrderDetailsChan struct {
	DeliveryOrderDetails []*DeliveryOrderDetail
	Total                int64
	Error                error
	ErrorLog             *model.ErrorLog
	ID                   int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type DeliveryOrderDetailRequest struct {
	ID                int     `json:"id,omitempty"`
	DoDetailID        int     `json:"do_detail_id,omitempty"`
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
	ProductID         int     `json:"product_id,omitempty"`
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
	Qty               int     `json:"qty,omitempty"`
	StartCreatedAt    string  `json:"start_created_at,omitempty"`
	EndCreatedAt      string  `json:"end_created_at,omitempty"`
	UpdatedAt         string  `json:"updated_at,omitempty"`
	StartDoDate       string  `json:"start_do_date,omitempty"`
	EndDoDate         string  `json:"end_do_date,omitempty"`
}

type DeliveryOrderDetails struct {
	DeliveryOrderDetails []*DeliveryOrderDetail `json:"delivery_order_details,omitempty"`
	Total                int64                  `json:"total,omitempty"`
}

type DeliveryOrderDetailOpenSearchDetailResponse struct {
	SoDetailID int       `json:"so_detail_id,omitempty" bson:"so_detail_id,omitempty"`
	Qty        NullInt64 `json:"qty,omitempty" bson:"qty,omitempty"`
}

type DeliveryOrderDetailOpenSearchResponse struct {
	ID              int                                     `json:"id,omitempty" bson:"id,omitempty"`
	DeliveryOrderID int                                     `json:"delivery_order_id,omitempty" bson:"delivery_order_id,omitempty"`
	ProductID       int                                     `json:"product_id,omitempty" bson:"product_id,omitempty"`
	Product         *ProductOpenSearchDeliveryOrderResponse `json:"product,omitempty" bson:"product,omitempty"`
	UomID           int                                     `json:"uom_id,omitempty" bson:"uom_id,omitempty"`
	UomName         string                                  `json:"uom_name,omitempty" bson:"uom_name,omitempty"`
	UomCode         string                                  `json:"uom_code,omitempty" bson:"uom_code,omitempty"`
	OrderStatusID   int                                     `json:"order_status_id,omitempty" bson:"order_status_id,omitempty"`
	OrderStatusName string                                  `json:"order_status_name,omitempty" bson:"order_status_name,omitempty"`
	DoDetailCode    string                                  `json:"do_detail_code,omitempty" bson:"do_detail_code,omitempty"`
	Qty             int                                     `json:"qty,omitempty" bson:"qty,omitempty"`
	SentQty         int                                     `json:"sent_qty,omitempty" bson:"sent_qty,omitempty"`
	ResidualQty     int                                     `json:"residual_qty,omitempty" bson:"residual_qty,omitempty"`
	Price           float64                                 `json:"price,omitempty" bson:"price,omitempty"`
	Note            NullString                              `json:"note,omitempty" bson:"note,omitempty"`
	CreatedAt       *time.Time                              `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type DeliveryOrderDetailOpenSearchRequest struct {
	ID                int               `json:"id,omitempty" bson:"id,omitempty"`
	PerPage           int               `json:"per_page,omitempty"`
	Page              int               `json:"page,omitempty"`
	SortField         string            `json:"sort_field,omitempty" bson:"sort_field,omitempty"`
	SortValue         string            `json:"sort_value,omitempty" bson:"sort_value,omitempty"`
	GlobalSearchValue string            `json:"global_search_value,omitempty" bson:"global_search_value,omitempty"`
	StartDoDate       string            `json:"start_so_date,omitempty"`
	EndDoDate         string            `json:"end_so_date,omitempty"`
	DeliveryOrderID   int               `json:"delivery_order_id,omitempty" bson:"delivery_order_id"`
	DoCode            string            `json:"do_code,omitempty"`
	DoDate            string            `json:"do_date,omitempty"`
	DoRefCode         string            `json:"do_ref_code,omitempty"`
	DoRefDate         string            `json:"do_ref_date,omitempty"`
	SalesOrderID      int               `json:"sales_order_id,omitempty"`
	SalesOrderCode    NullString        `json:"sales_order_code,omitempty" bson:"sales_order_code"`
	SalesOrderDate    NullString        `json:"sales_order_date,omitempty" bson:"sales_order_date"`
	SoRefDate         NullString        `json:"so_ref_date,omitempty" bson:"so_ref_date,omitempty"`
	SoDetailID        int               `json:"so_detail_id,omitempty" bson:"so_detail_id"`
	SoDetailCode      string            `json:"so_detail_code,omitempty" bson:"so_detail_code,omitempty"`
	SoDetail          *SalesOrderDetail `json:"so_detail,omitempty" bson:" so_detail"`
	AgentID           int               `json:"agentID,omitempty"`
	AgentName         string            `json:"agent_name,omitempty"`
	AgentAddress      NullString        `json:"agent_address,omitempty" bson:"agent_address,omitempty"`
	AgentPhone        NullString        `json:"agent_phone,omitempty" bson:"agent_phone,omitempty"`
	AgentProvinceID   int               `json:"agent_province_id,omitempty" bson:"agent_province_id,omitempty"`
	AgentProvinceName NullString        `json:"agent_province_name,omitempty" bson:"agent_province_name,omitempty"`
	AgentCityID       int               `json:"agent_city_id,omitempty" bson:"agent_city_id,omitempty"`
	AgentCityName     NullString        `json:"agent_city_name,omitempty" bson:"agent_city_name,omitempty"`
	AgentDistrictID   int               `json:"agent_district_id,omitempty" bson:"agent_district_id,omitempty"`
	AgentDistrictName NullString        `json:"agent_district_name,omitempty" bson:"agent_district_name,omitempty"`
	AgentVillageID    int               `json:"agent_village_id,omitempty" bson:"agent_village_id,omitempty"`
	AgentVillageName  NullString        `json:"agent_village_name,omitempty" bson:"agent_village_name,omitempty"`
	StoreID           int               `json:"store_id,omitempty" bson:"store_id,omitempty"`
	StoreName         NullString        `json:"store_name,omitempty" bson:"store_name,omitempty"`
	StoreCode         NullString        `json:"store_code,omitempty" bson:"store_code,omitempty"`
	StoreAddress      NullString        `json:"store_address,omitempty" bson:"store_address,omitempty"`
	StorePhone        NullString        `json:"store_phone,omitempty" bson:"store_phone,omitempty"`
	StoreProvinceID   int               `json:"store_province_id,omitempty" bson:"store_province_id,omitempty"`
	StoreProvinceName NullString        `json:"store_province_name,omitempty" bson:"store_province_name,omitempty"`
	StoreCityID       int               `json:"store_city_id,omitempty" bson:"store_city_id,omitempty"`
	StoreCityName     NullString        `json:"store_city_name,omitempty" bson:"store_city_name,omitempty"`
	StoreDistrictID   int               `json:"store_district_id,omitempty" bson:"store_district_id,omitempty"`
	StoreDistrictName NullString        `json:"store_district_name,omitempty" bson:"store_district_name,omitempty"`
	StoreVillageID    int               `json:"store_village_id,omitempty" bson:"store_village_id,omitempty"`
	StoreVillageName  NullString        `json:"store_village_name,omitempty" bson:"store_village_name,omitempty"`
	WarehouseID       int               `json:"warehouse_id,omitempty" bson:"warehouse_id"`
	WarehouseCode     string            `json:"warehouse_code,omitempty" bson:"warehouse_code"`
	WarehouseName     string            `json:"warehouse_name,omitempty" bson:"warehouse_name"`
	SalesmanID        int               `json:"salesman_id,omitempty" bson:"salesman_id,omitempty"`
	SalesmanName      string            `json:"salesman_name,omitempty" bson:"salesman_name,omitempty"`
	BrandID           int               `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	BrandName         string            `json:"brand_name,omitempty" bson:"brand_name,omitempty"`
	ProductID         int               `json:"product_id,omitempty" bson:"product_id"`
	ProductSKU        string            `json:"product_sku,omitempty" bson:"product_sku,omitempty"`
	ProductName       string            `json:"product_name,omitempty" bson:"product_name,omitempty"`
	Description       NullString        `json:"product_description,omitempty" bson:"product_description,omitempty"`
	UomID             int               `json:"uom_id,omitempty" bson:"uom_id"`
	UomCode           string            `json:"uom_code,omitempty" bson:"uom_code,omitempty"`
	UomName           string            `json:"uom_name,omitempty" bson:"uom_name,omitempty"`
	DoDetailCode      string            `json:"do_detail_code,omitempty" bson:"do_detail_code"`
	OrderStatusID     int               `json:"order_status_id,omitempty" bson:"order_status_id"`
	OrderStatusName   string            `json:"order_status_name,omitempty" bson:"order_status_name"`
	Qty               int               `json:"qty,omitempty" bson:"qty"`
	Note              NullString        `json:"note,omitempty" bson:"note"`
	StartCreatedAt    string            `json:"start_created_at,omitempty"`
	EndCreatedAt      string            `json:"end_created_at,omitempty"`
}

type DeliveryOrderDetailLogData struct {
	ID            int        `json:"id,omitempty" bson:"id"`
	AgentID       int        `json:"agent_id,omitempty" bson:"agent_id"`
	AgentName     string     `json:"agent_name,omitempty" bson:"agent_name"`
	DoRefCode     NullString `json:"do_ref_code,omitempty" bson:"do_ref_code"`
	DoDate        string     `json:"do_date,omitempty" bson:"do_date"`
	DoNumber      string     `json:"do_number,omitempty" bson:"do_number"`
	DoDetailCode  string     `json:"do_detail_code,omitempty" bson:"do_detail_code"`
	SoDetailID    int        `json:"so_detail_id,omitempty" bson:"so_detail_id"`
	Note          NullString `json:"note,omitempty" bson:"note"`
	InternalNote  NullString `json:"internal_note,omitempty" bson:"internal_note"`
	DriverName    NullString `json:"driver_name,omitempty" bson:"driver_name"`
	PlatNumber    NullString `json:"plat_number,omitempty" bson:"plat_number"`
	BrandID       int        `json:"brand_id,omitempty" bson:"brand_id"`
	BrandName     string     `json:"brand_name,omitempty"  bson:"brand_name"`
	ProductID     int        `json:"product_id,omitempty" bson:"product_id"`
	ProductName   string     `json:"product_name,omitempty" bson:"product_name,omitempty"`
	DeliveryQty   int        `json:"delivery_qty,omitempty" bson:"delivery_qty"`
	UomCode       string     `json:"uom_code,omitempty" bson:"uom_code,omitempty"`
	WarehouseID   int        `json:"warehouse_id,omitempty" bson:"warehouse_id"`
	WarehouseName string     `json:"warehouse_name,omitempty" bson:"warehouse_name"`
}
