package models

import (
	"order-service/global/utils/model"
	"time"
)

type SalesOrderDetail struct {
	ID                int          `json:"id,omitempty" bson:"id,omitempty"`
	SalesOrderID      int          `json:"sales_order_id,omitempty" bson:"sales_order_id,omitempty"`
	ProductID         int          `json:"product_id,omitempty" bson:"product_id,omitempty" `
	ProductSKU        string       `json:"product_sku,omitempty" bson:"product_sku,omitempty"`
	ProductName       string       `json:"product_name,omitempty" bson:"product_name,omitempty"`
	Product           *Product     `json:"product,omitempty" bson:"product,omitempty"`
	BrandID           int          `json:"brand_id,omitempty" bson:"brand_id,omitempty" `
	UomID             int          `json:"uom_id,omitempty" bson:"uom_id,omitempty"`
	UomCode           string       `json:"uom_code,omitempty" bson:"uom_code,omitempty"`
	UomName           string       `json:"uom_name,omitempty" bson:"uom_name,omitempty"`
	Uom               *Uom         `json:"uom,omitempty" bson:"uom,omitempty"`
	OrderStatusID     int          `json:"order_status_id,omitempty" bson:"order_status_id,omitempty"`
	OrderStatusName   string       `json:"order_status_name,omitempty" bson:"order_status_name,omitempty"`
	OrderStatus       *OrderStatus `json:"order_status,omitempty" bson:"order_status,omitempty"`
	SoDetailCode      string       `json:"so_detail_code,omitempty" bson:"so_detail_code,omitempty"`
	Qty               int          `json:"qty,omitempty" bson:"qty,omitempty"`
	SentQty           int          `json:"sent_qty,omitempty" bson:"sent_qty,omitempty"`
	ResidualQty       int          `json:"residual_qty,omitempty" bson:"residual_qty,omitempty"`
	Price             float64      `json:"price,omitempty" bson:"price,omitempty"`
	Note              NullString   `json:"note,omitempty" bson:"note,omitempty"`
	IsDoneSyncToEs    string       `json:"is_done_sync_to_es,omitempty" bson:"is_done_sync_to_es,omitempty"`
	StartDateSyncToEs *time.Time   `json:"start_date_sync_to_es,omitempty" bson:"start_date_sync_to_es,omitempty"`
	EndDateSyncToEs   *time.Time   `json:"end_date_sync_to_es,omitempty" bson:"end_date_sync_to_es,omitempty"`
	UomType           string       `json:"uom_type,omitempty" bson:"uom_type,omitempty"`
	Subtotal          float64      `json:"subtotal,omitempty" bson:"subtotal,omitempty"`
	FirstCategoryId   int          `json:"first_category_id,omitempty" bson:"first_category_id,omitempty"`
	FirstCategoryName *string      `json:"first_category_name,omitempty" bson:"first_category_name,omitempty"`
	LastCategoryId    int          `json:"last_category_id,omitempty" bson:"last_category_id,omitempty"`
	LastCategoryName  *string      `json:"last_category_name,omitempty" bson:"last_category_name,omitempty"`
	CreatedBy         int          `json:"created_by,omitempty" bson:"created_by,omitempty"`
	LatestUpdatedBy   int          `json:"latest_updated_by,omitempty" bson:"latest_updated_by,omitempty"`
	DeletedBy         int          `json:"deleted_by,omitempty" bson:"deleted_by,omitempty"`
	CreatedAt         *time.Time   `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt         *time.Time   `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt         *time.Time   `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type SalesOrderDetailTemplate struct {
	ProductID int    `json:"product_id,omitempty" binding:"required"`
	UomID     int    `json:"uom_id,omitempty" binding:"required"`
	Qty       int    `json:"qty,omitempty" binding:"required"`
	Note      string `json:"note,omitempty"`
}

type SalesOrderDetailStoreRequest struct {
	SalesOrderDetailTemplate
	BrandID int     `json:"brand_id,omitempty" bson:"brand_id,omitempty" binding:"required"`
	Price   float64 `json:"price,omitempty" binding:"required"`
}

type SalesOrderDetailStoreResponse struct {
	ID                    int        `json:"id,omitempty" bson:"id,omitempty"`
	SalesOrderId          int        `json:"sales_order_id,omitempty"`
	BrandID               int        `json:"brand_id,omitempty" bson:"brand_id,omitempty" binding:"required"`
	ProductID             int        `json:"product_id,omitempty" binding:"required"`
	ProductSKU            string     `json:"product_sku,omitempty"`
	ProductName           string     `json:"product_name,omitempty"`
	CategoryId            int        `json:"category_id,omitempty"`
	CategoryName          string     `json:"category_name,omitempty"`
	UnitMeasurementSmall  string     `json:"unit_measurement_small,omitempty" bson:"unit_measurement_small,omitempty"`
	UnitMeasurementMedium string     `json:"unit_measurement_medium,omitempty" bson:"unit_measurement_medium,omitempty"`
	UnitMeasurementBig    string     `json:"unit_measurement_big,omitempty" bson:"unit_measurement_big,omitempty"`
	UomID                 int        `json:"uom_id,omitempty" binding:"required"`
	UomCode               string     `json:"uom_code,omitempty"`
	OrderStatusId         int        `json:"order_status_id,omitempty"`
	OrderStatusName       string     `json:"order_status_name,omitempty"`
	SoDetailCode          string     `json:"so_detail_code,omitempty"`
	Qty                   int        `json:"qty,omitempty" binding:"required"`
	SentQty               NullInt64  `json:"sent_qty,omitempty"`
	ResidualQty           NullInt64  `json:"residual_qty,omitempty"`
	Price                 float64    `json:"price,omitempty" binding:"required"`
	Note                  string     `json:"note,omitempty"`
	CreatedAt             *time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type SalesOrderDetailUpdateRequest struct {
	ID            int    `json:"id,omitempty" bson:"id,omitempty" binding:"required"`
	SoDetailCode  string `json:"so_detail_code,omitempty" bson:"so_detail_code,omitempty"`
	OrderStatusID int    `json:"order_status_id,omitempty" bson:"order_status_id,omitempty" binding:"required"`
	Reason        string `json:"reason,omitempty" bson:"reason,omitempty" binding:"required"`
}

type SalesOrderDetailUpdateByIdRequest struct {
	SoDetailCode  string `json:"so_detail_code,omitempty" bson:"so_detail_code,omitempty"`
	OrderStatusID int    `json:"order_status_id,omitempty" bson:"order_status_id,omitempty" binding:"required"`
	Reason        string `json:"reason,omitempty" bson:"reason,omitempty" binding:"required"`
}

type SalesOrderDetailChan struct {
	SalesOrderDetail *SalesOrderDetail
	Total            int64
	Error            error
	ErrorLog         *model.ErrorLog
	ID               int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type SalesOrderDetailsChan struct {
	SalesOrderDetails []*SalesOrderDetail
	ErrorLog          *model.ErrorLog
	Total             int64
	Error             error
	ID                int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type SalesOrderDetailRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type SalesOrderDetails struct {
	SalesOrderDetails []*SalesOrderDetail `json:"sales_order_details,omitempty"`
	Total             int64               `json:"total,omitempty"`
}

type SalesOrderDetailOpenSearch struct {
	ID                   int          `json:"id,omitempty" bson:"id,omitempty"`
	SoID                 int          `json:"so_id,omitempty" bson:"so_id,omitempty"`
	SoDetailCode         string       `json:"so_detail_code,omitempty" bson:"so_detail_code,omitempty"`
	Qty                  int          `json:"qty,omitempty" bson:"qty,omitempty"`
	SentQty              int          `json:"sent_qty,omitempty" bson:"sent_qty,omitempty"`
	ResidualQty          int          `json:"residual_qty,omitempty" bson:"residual_qty,omitempty"`
	Price                float64      `json:"price,omitempty" bson:"price,omitempty"`
	Subtotal             float64      `json:"subtotal,omitempty" bson:"subtotal,omitempty"`
	Note                 NullString   `json:"note,omitempty" bson:"note,omitempty"`
	SalesOrderID         int          `json:"sales_order_id,omitempty" bson:"sales_order_id,omitempty"`
	SoCode               string       `json:"so_code,omitempty" bson:"so_code,omitempty"`
	SoDate               string       `json:"so_date,omitempty" bson:"so_date,omitempty"`
	SoRefCode            string       `json:"so_ref_code,omitempty" bson:"so_ref_code,omitempty"`
	SoRefDate            NullString   `json:"so_ref_date,omitempty" bson:"so_ref_date,omitempty"`
	ReferralCode         NullString   `json:"referral_code,omitempty" bson:"referral_code,omitempty"`
	SoReferralCode       NullString   `json:"so_referral_code,omitempty" bson:"so_referral_code,omitempty"`
	AgentId              int          `json:"agent_id,omitempty" bson:"agent_id,omitempty"`
	AgentName            string       `json:"agent_name,omitempty" bson:"agent_name,omitempty"`
	AgentProvinceID      int          `json:"agent_province_id,omitempty" bson:"agent_province_id,omitempty"`
	AgentProvinceName    string       `json:"agent_province_name,omitempty" bson:"agent_province_name,omitempty"`
	AgentCityID          int          `json:"agent_city_id,omitempty" bson:"agent_city_id,omitempty"`
	AgentCityName        string       `json:"agent_city_name,omitempty" bson:"agent_city_name,omitempty"`
	AgentDistrictID      int          `json:"agent_district_id,omitempty" bson:"agent_district_id,omitempty"`
	AgentDistrictName    string       `json:"agent_district_name,omitempty" bson:"agent_district_name,omitempty"`
	AgentVillageID       int          `json:"agent_village_id,omitempty" bson:"agent_village_id,omitempty"`
	AgentVillageName     string       `json:"agent_village_name,omitempty" bson:"agent_village_name,omitempty"`
	AgentPhone           string       `json:"agent_phone,omitempty" bson:"agent_phone,omitempty"`
	AgentAddress         string       `json:"agent_address,omitempty" bson:"agent_address,omitempty"`
	Agent                *Agent       `json:"agent,omitempty" bson:"agent,omitempty"`
	StoreID              int          `json:"store_id,omitempty" bson:"store_id,omitempty"`
	StoreCode            string       `json:"store_code,omitempty" bson:"store_code,omitempty"`
	StoreName            string       `json:"store_name,omitempty" bson:"store_name,omitempty"`
	StoreProvinceID      int          `json:"store_province_id,omitempty" bson:"store_province_id,omitempty"`
	StoreProvinceName    string       `json:"store_province_name,omitempty" bson:"store_province_name,omitempty"`
	StoreCityID          int          `json:"store_city_id,omitempty" bson:"store_city_id,omitempty"`
	StoreCityName        string       `json:"store_city_name,omitempty" bson:"store_city_name,omitempty"`
	StoreDistrictID      int          `json:"store_district_id,omitempty" bson:"store_district_id,omitempty"`
	StoreDistrictName    string       `json:"store_district_name,omitempty" bson:"store_district_name,omitempty"`
	StoreVillageID       int          `json:"store_village_id,omitempty" bson:"store_village_id,omitempty"`
	StoreVillageName     string       `json:"store_village_name,omitempty" bson:"store_village_name,omitempty"`
	StoreAddress         string       `json:"store_address,omitempty" bson:"store_address,omitempty"`
	StorePhone           string       `json:"store_phone,omitempty" bson:"store_phone,omitempty"`
	StoreMainMobilePhone string       `json:"store_main_mobile_phone,omitempty" bson:"store_main_mobile_phone,omitempty"`
	Store                *Store       `json:"store,omitempty" bson:"store,omitempty"`
	BrandID              int          `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	BrandName            string       `json:"brand_name,omitempty" bson:"brand_name,omitempty"`
	Brand                *Brand       `json:"brand,omitempty" bson:"brand,omitempty"`
	UserID               int          `json:"user_id,omitempty" bson:"user_id,omitempty"`
	UserFirstName        string       `json:"user_first_name,omitempty" bson:"user_first_name,omitempty"`
	UserLastName         string       `json:"user_last_name,omitempty" bson:"user_last_name,omitempty"`
	UserRoleID           int          `json:"user_role_id,omitempty" bson:"user_role_id,omitempty"`
	UserEmail            string       `json:"user_email,omitempty" bson:"user_email,omitempty"`
	User                 *User        `json:"user,omitempty" bson:"user,omitempty"`
	SalesmanID           int          `json:"salesman_id,omitempty" bson:"salesman_id,omitempty"`
	SalesmanName         string       `json:"salesman_name,omitempty" bson:"salesman_name,omitempty"`
	SalesmanEmail        string       `json:"salesman_email,omitempty" bson:"salesman_email,omitempty"`
	Salesman             *Salesman    `json:"salesman,omitempty" bson:"salesman,omitempty"`
	OrderSourceID        int          `json:"order_source_id,omitempty" bson:"order_source_id,omitempty"`
	OrderSourceName      string       `json:"order_source_name,omitempty" bson:"order_source_name,omitempty"`
	OrderSource          *OrderSource `json:"order_source,omitempty" bson:"order_source,omitempty"`
	OrderStatusID        int          `json:"order_status_id,omitempty" bson:"order_status_id,omitempty"`
	OrderStatusName      string       `json:"order_status_name,omitempty" bson:"order_status_name,omitempty"`
	OrderStatus          *OrderStatus `json:"order_status,omitempty" bson:"order_status,omitempty"`
	GLat                 NullFloat64  `json:"g_lat,omitempty" bson:"g_lat,omitempty"`
	GLong                NullFloat64  `json:"g_long,omitempty" bson:"g_long,omitempty"`
	ProductID            int          `json:"product_id,omitempty" bson:"product_id,omitempty" `
	ProductSKU           string       `json:"product_sku,omitempty" bson:"product_sku,omitempty"`
	ProductName          string       `json:"product_name,omitempty" bson:"product_name,omitempty"`
	ProductDescription   string       `json:"product_description,omitempty" bson:"product_description,omitempty"`
	CategoryID           int          `json:"category_id,omitempty" bson:"category_id,omitempty" `
	Product              *Product     `json:"product,omitempty" bson:"product,omitempty"`
	UomID                int          `json:"uom_id,omitempty" bson:"uom_id,omitempty"`
	UomCode              string       `json:"uom_code,omitempty" bson:"uom_code,omitempty"`
	UomName              string       `json:"uom_name,omitempty" bson:"uom_name,omitempty"`
	Uom                  *Uom         `json:"uom,omitempty" bson:"uom,omitempty"`
	IsDoneSyncToEs       string       `json:"is_done_sync_to_es,omitempty" bson:"is_done_sync_to_es,omitempty"`
	StartDateSyncToEs    *time.Time   `json:"start_date_sync_to_es,omitempty" bson:"start_date_sync_to_es,omitempty"`
	EndDateSyncToEs      *time.Time   `json:"end_date_sync_to_es,omitempty" bson:"end_date_sync_to_es,omitempty"`
	FirstCategoryId      int          `json:"first_category_id,omitempty" bson:"first_category_id,omitempty"`
	FirstCategoryName    *string      `json:"first_category_name,omitempty" bson:"first_category_name,omitempty"`
	LastCategoryId       int          `json:"last_category_id,omitempty" bson:"last_category_id,omitempty"`
	LastCategoryName     *string      `json:"last_category_name,omitempty" bson:"last_category_name,omitempty"`
	CreatedBy            int          `json:"created_by,omitempty" bson:"created_by,omitempty"`
	UpdatedBy            int          `json:"updated_by,omitempty" bson:"updated_by,omitempty"`
	DeletedBy            int          `json:"deleted_by,omitempty" bson:"deleted_by,omitempty"`
	CreatedAt            *time.Time   `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt            *time.Time   `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt            *time.Time   `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type SalesOrderDetailOpenSearchChan struct {
	SalesOrderDetail *SalesOrderDetailOpenSearch
	Total            int64
	Error            error
	ErrorLog         *model.ErrorLog
	ID               int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type SalesOrderDetailsOpenSearchChan struct {
	SalesOrderDetails []*SalesOrderDetailOpenSearch
	Total             int64
	Error             error
	ErrorLog          *model.ErrorLog
	ID                int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type SalesOrderDetailsOpenSearch struct {
	SalesOrderDetails []*SalesOrderDetailOpenSearch `json:"sales_order_details,omitempty"`
	Total             int64                         `json:"total,omitempty"`
}

type SalesOrderDetailOpenSearchResponse struct {
	ID            int                            `json:"id,omitempty" bson:"id,omitempty"`
	SalesOrderID  int                            `json:"sales_order_id,omitempty" bson:"sales_order_id,omitempty"`
	ProductID     int                            `json:"product_id,omitempty" bson:"product_id,omitempty" `
	BrandID       int                            `json:"brand_id,omitempty" bson:"brand_id,omitempty" `
	Product       *ProductOpenSearchResponse     `json:"product,omitempty" bson:"product,omitempty"`
	UomID         int                            `json:"uom_id,omitempty" bson:"uom_id,omitempty"`
	Uom           *UomOpenSearchResponse         `json:"uom,omitempty" bson:"uom,omitempty"`
	OrderStatusID int                            `json:"order_status_id,omitempty" bson:"order_status_id,omitempty"`
	OrderStatus   *OrderStatusOpenSearchResponse `json:"order_status,omitempty" bson:"order_status,omitempty"`
	SoDetailCode  string                         `json:"so_detail_code,omitempty" bson:"so_detail_code,omitempty"`
	Qty           NullInt64                      `json:"qty,omitempty" bson:"qty,omitempty"`
	SentQty       NullInt64                      `json:"sent_qty,omitempty" bson:"sent_qty,omitempty"`
	ResidualQty   NullInt64                      `json:"residual_qty,omitempty" bson:"residual_qty,omitempty"`
	Price         float64                        `json:"price,omitempty" bson:"price,omitempty"`
	Note          NullString                     `json:"note,omitempty" bson:"note,omitempty"`
	CreatedAt     *time.Time                     `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type SalesOrderDetailsOpenSearchResponse struct {
	SalesOrderDetails []*SalesOrderDetailOpenSearchResponse `json:"sales_order_details,omitempty"`
	Total             int64                                 `json:"total,omitempty"`
}

type GetSalesOrderDetailRequest struct {
	ID                int    `json:"id,omitempty"`
	SoID              int    `json:"so_id,omitempty"`
	PerPage           int    `json:"per_page,omitempty"`
	Page              int    `json:"page,omitempty"`
	SortField         string `json:"sort_field,omitempty" bson:"sort_field,omitempty"`
	SortValue         string `json:"sort_value,omitempty" bson:"sort_value,omitempty"`
	GlobalSearchValue string `json:"global_search_value,omitempty" bson:"global_search_value,omitempty"`
	AgentID           int    `json:"agent_id,omitempty"`
	StoreID           int    `json:"store_id,omitempty"`
	BrandID           int    `json:"brand_id,omitempty"`
	OrderSourceID     int    `json:"order_source_id,omitempty"`
	OrderStatusID     int    `json:"order_status_id,omitempty"`
	StartSoDate       string `json:"start_so_date,omitempty"`
	EndSoDate         string `json:"end_so_date,omitempty"`
	ProductID         int    `json:"product_id,omitempty"`
	CategoryID        int    `json:"category_id,omitempty"`
	SalesmanID        int    `json:"salesman_id,omitempty"`
	ProvinceID        int    `json:"province_id,omitempty"`
	CityID            int    `json:"city_id,omitempty"`
	DistrictID        int    `json:"district_id,omitempty"`
	VillageID         int    `json:"village_id,omitempty"`
	StartCreatedAt    string `json:"start_created_at,omitempty"`
	EndCreatedAt      string `json:"end_created_at,omitempty"`
}

type SalesOrderDetailExportRequest struct {
	ID                int    `json:"id,omitempty"`
	SoID              int    `json:"so_id,omitempty"`
	SortField         string `json:"sort_field,omitempty" bson:"sort_field,omitempty"`
	SortValue         string `json:"sort_value,omitempty" bson:"sort_value,omitempty"`
	GlobalSearchValue string `json:"global_search_value,omitempty" bson:"global_search_value,omitempty"`
	FileType          string `json:"file_type,omitempty"`
	FileName          string `json:"file_name,omitempty"`
	AgentID           int    `json:"agent_id,omitempty"`
	StoreID           int    `json:"store_id,omitempty"`
	BrandID           int    `json:"brand_id,omitempty"`
	OrderSourceID     int    `json:"order_source_id,omitempty"`
	OrderStatusID     int    `json:"order_status_id,omitempty"`
	StartSoDate       string `json:"start_so_date,omitempty"`
	EndSoDate         string `json:"end_so_date,omitempty"`
	ProductID         int    `json:"product_id,omitempty"`
	CategoryID        int    `json:"category_id,omitempty"`
	SalesmanID        int    `json:"salesman_id,omitempty"`
	ProvinceID        int    `json:"province_id,omitempty"`
	CityID            int    `json:"city_id,omitempty"`
	DistrictID        int    `json:"district_id,omitempty"`
	VillageID         int    `json:"village_id,omitempty"`
	StartCreatedAt    string `json:"start_created_at,omitempty"`
	EndCreatedAt      string `json:"end_created_at,omitempty"`
}
