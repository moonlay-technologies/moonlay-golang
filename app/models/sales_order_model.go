package models

import (
	"order-service/global/utils/model"
	"time"
)

type SalesOrder struct {
	ID                   int                 `json:"id,omitempty" bson:"id,omitempty"`
	CartID               int                 `json:"cart_id,omitempty" bson:"cart_id,omitempty"`
	AgentID              int                 `json:"agent_id,omitempty" bson:"agent_id,omitempty"`
	AgentName            NullString          `json:"agent_name,omitempty" bson:"agent_name,omitempty"`
	AgentEmail           NullString          `json:"agent_email,omitempty" bson:"agent_email,omitempty"`
	AgentProvinceID      int                 `json:"agent_province_id,omitempty" bson:"agent_province_id,omitempty"`
	AgentProvinceName    NullString          `json:"agent_province_name,omitempty" bson:"agent_province_name,omitempty"`
	AgentCityID          int                 `json:"agent_city_id,omitempty" bson:"agent_city_id,omitempty"`
	AgentCityName        NullString          `json:"agent_city_name,omitempty" bson:"agent_city_name,omitempty"`
	AgentDistrictID      int                 `json:"agent_district_id,omitempty" bson:"agent_district_id,omitempty"`
	AgentDistrictName    NullString          `json:"agent_district_name,omitempty" bson:"agent_district_name,omitempty"`
	AgentVillageID       int                 `json:"agent_village_id,omitempty" bson:"agent_village_id,omitempty"`
	AgentVillageName     NullString          `json:"agent_village_name,omitempty" bson:"agent_village_name,omitempty"`
	AgentAddress         NullString          `json:"agent_address,omitempty" bson:"agent_address,omitempty"`
	AgentPhone           NullString          `json:"agent_phone,omitempty" bson:"agent_phone,omitempty"`
	AgentMainMobilePhone NullString          `json:"agent_main_mobile_phone,omitempty" bson:"agent_main_mobile_phone,omitempty"`
	Agent                *Agent              `json:"agent,omitempty" bson:"agent,omitempty"`
	StoreID              int                 `json:"store_id,omitempty" bson:"store_id,omitempty"`
	StoreName            NullString          `json:"store_name,omitempty" bson:"store_name,omitempty"`
	StoreCode            NullString          `json:"store_code,omitempty" bson:"store_code,omitempty"`
	StoreEmail           NullString          `json:"store_email,omitempty" bson:"store_email,omitempty"`
	StoreProvinceID      int                 `json:"store_province_id,omitempty" bson:"store_province_id,omitempty"`
	StoreProvinceName    NullString          `json:"store_province_name,omitempty" bson:"store_province_name,omitempty"`
	StoreCityID          int                 `json:"store_city_id,omitempty" bson:"store_city_id,omitempty"`
	StoreCityName        NullString          `json:"store_city_name,omitempty" bson:"store_city_name,omitempty"`
	StoreDistrictID      int                 `json:"store_district_id,omitempty" bson:"store_district_id,omitempty"`
	StoreDistrictName    NullString          `json:"store_district_name,omitempty" bson:"store_district_name,omitempty"`
	StoreVillageID       int                 `json:"store_village_id,omitempty" bson:"store_village_id,omitempty"`
	StoreVillageName     NullString          `json:"store_village_name,omitempty" bson:"store_village_name,omitempty"`
	StoreAddress         NullString          `json:"store_address,omitempty" bson:"store_address,omitempty"`
	StorePhone           NullString          `json:"store_phone,omitempty" bson:"store_phone,omitempty"`
	StoreMainMobilePhone NullString          `json:"store_main_mobile_phone,omitempty" bson:"store_main_mobile_phone,omitempty"`
	Store                *Store              `json:"store,omitempty" bson:"store,omitempty"`
	BrandID              int                 `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	BrandName            string              `json:"brand_name,omitempty" bson:"brand_name,omitempty"`
	Brand                *Brand              `json:"brand,omitempty" bson:"brand,omitempty"`
	UserID               int                 `json:"user_id,omitempty" bson:"user_id,omitempty"`
	UserFirstName        NullString          `json:"user_first_name,omitempty" bson:"user_first_name,omitempty"`
	UserLastName         NullString          `json:"user_last_name,omitempty" bson:"user_last_name,omitempty"`
	UserRoleID           int                 `json:"user_role_id,omitempty" bson:"user_role_id,omitempty"`
	UserEmail            NullString          `json:"user_email,omitempty" bson:"user_email,omitempty"`
	User                 *User               `json:"user,omitempty" bson:"user,omitempty"`
	Salesman             *Salesman           `json:"salesman,omitempty" bson:"salesman,omitempty"`
	VisitationID         int                 `json:"visitation_id,omitempty" bson:"visitation_id,omitempty"`
	OrderSourceID        int                 `json:"order_source_id,omitempty" bson:"order_source_id,omitempty"`
	OrderSourceName      string              `json:"order_source_name,omitempty" bson:"order_source_name,omitempty"`
	OrderSource          *OrderSource        `json:"order_source,omitempty" bson:"order_source,omitempty"`
	OrderStatusID        int                 `json:"order_status_id,omitempty" bson:"order_status_id,omitempty"`
	OrderStatus          *OrderStatus        `json:"order_status,omitempty" bson:"order_status,omitempty"`
	OrderStatusName      string              `json:"order_status_name,omitempty" bson:"order_status_name,omitempty"`
	SoCode               string              `json:"so_code,omitempty" bson:"so_code,omitempty"`
	SoDate               string              `json:"so_date,omitempty" bson:"so_date,omitempty"`
	SoRefCode            NullString          `json:"so_ref_code,omitempty" bson:"so_ref_code,omitempty"`
	SoRefDate            NullString          `json:"so_ref_date,omitempty" bson:"so_ref_date,omitempty"`
	ReferralCode         NullString          `json:"referral_code,omitempty" bson:"referral_code,omitempty"`
	GLat                 NullFloat64         `json:"g_lat,omitempty" bson:"g_lat,omitempty"`
	GLong                NullFloat64         `json:"g_long,omitempty" bson:"g_long,omitempty"`
	DeviceId             NullString          `json:"device_id,omitempty" bson:"device_id,omitempty"`
	Note                 NullString          `json:"note,omitempty" bson:"note,omitempty"`
	InternalComment      NullString          `json:"internal_comment,omitempty" bson:"internal_comment,omitempty"`
	TotalAmount          float64             `json:"total_amount,omitempty" bson:"total_amount,omitempty"`
	TotalTonase          float64             `json:"total_tonase,omitempty" bson:"total_tonase,omitempty"`
	IsDoneSyncToEs       string              `json:"is_done_sync_to_es,omitempty" bson:"is_done_sync_to_es,omitempty"`
	StartDateSyncToEs    *time.Time          `json:"start_date_sync_to_es,omitempty" bson:"start_date_sync_to_es,omitempty"`
	EndDateSyncToEs      *time.Time          `json:"end_date_sync_to_es,omitempty" bson:"end_date_sync_to_es,omitempty"`
	StartCreatedDate     *time.Time          `json:"start_created_date,omitempty" bson:"start_created_date,omitempty"`
	EndCreatedDate       *time.Time          `json:"end_created_date,omitempty" bson:"end_created_date,omitempty"`
	SalesOrderDetails    []*SalesOrderDetail `json:"sales_order_details,omitempty" bson:"sales_order_details,omitempty"`
	DeliveryOrders       []*DeliveryOrder    `json:"delivery_orders,omitempty" bson:"delivery_orders,omitempty"`
	SalesmanID           NullInt64           `json:"salesman_id,omitempty" bson:"salesman_id,omitempty"`
	SalesmanName         NullString          `json:"salesman_name,omitempty" bson:"salesman_name,omitempty"`
	SalesmanEmail        NullString          `json:"salesman_email,omitempty" bson:"salesman_email,omitempty"`
	SalesOrderLogID      string              `json:"sales_order_log_id,omitempty" bson:"sales_order_log_id,omitempty"`
	CreatedBy            int                 `json:"created_by,omitempty" bson:"created_by,omitempty"`
	LatestUpdatedBy      int                 `json:"latest_updated_by,omitempty" bson:"latest_updated_by,omitempty"`
	CreatedAt            *time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt            *time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt            *time.Time          `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type SalesOrderTemplate struct {
	AgentID         int     `json:"agent_id,omitempty" bson:"agent_id,omitempty" binding:"required"`
	StoreID         int     `json:"store_id,omitempty" bson:"store_id,omitempty" binding:"required"`
	UserID          int     `json:"user_id,omitempty" bson:"user_id,omitempty" binding:"required"`
	OrderSourceID   int     `json:"order_source_id,omitempty" bson:"order_source_id,omitempty" binding:"required"`
	SoRefCode       string  `json:"so_ref_code,omitempty" bson:"so_ref_code,omitempty"`
	GLong           float64 `json:"g_long,omitempty" bson:"g_long,omitempty"`
	GLat            float64 `json:"g_lat,omitempty" bson:"g_lat,omitempty"`
	Note            string  `json:"note,omitempty" bson:"note,omitempty"`
	InternalComment string  `json:"internal_comment,omitempty" bson:"internal_comment,omitempty"`
	TotalAmount     float64 `json:"total_amount,omitempty" bson:"total_amount,omitempty" binding:"required"`
	TotalTonase     float64 `json:"total_tonase,omitempty" bson:"total_tonase,omitempty" binding:"required"`
	DeviceId        string  `json:"device_id,omitempty" bson:"device_id,omitempty"`
	ReferralCode    string  `json:"referral_code,omitempty" bson:"referral_code,omitempty"`
}

type SalesOrderStoreRequest struct {
	CartID int `json:"cart_id,omitempty" bson:"cart_id,omitempty"`
	SalesOrderTemplate
	BrandID           int                             `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	SalesmanID        int                             `json:"salesman_id,omitempty" bson:"salesman_id,omitempty"`
	VisitationID      int                             `json:"visitation_id,omitempty" bson:"visitation_id,omitempty"`
	SoDate            string                          `json:"so_date,omitempty" bson:"so_date,omitempty" binding:"required"`
	SoRefDate         string                          `json:"so_ref_date,omitempty" bson:"so_ref_date,omitempty" binding:"required"`
	OrderStatusID     int                             `json:"order_status_id,omitempty" bson:"order_status_id,omitempty"`
	RequestID         string                          `json:"request_id,omitempty" bson:"request_id,omitempty"`
	SalesOrderDetails []*SalesOrderDetailStoreRequest `json:"sales_order_details" bson:"sales_order_details" binding:"required,dive,required"`
}

type SalesOrderResponse struct {
	ID                   int                              `json:"id,omitempty" bson:"id,omitempty"`
	CartID               int                              `json:"cart_id,omitempty" bson:"cart_id,omitempty"`
	AgentID              int                              `json:"agent_id,omitempty" bson:"agent_id,omitempty" binding:"required"`
	AgentName            string                           `json:"agent_name,omitempty" bson:"agent_name,omitempty"`
	AgentEmail           string                           `json:"agent_email,omitempty" bson:"agent_email,omitempty"`
	AgentProvinceName    string                           `json:"agent_province_name,omitempty" bson:"agent_province_name,omitempty"`
	AgentCityName        string                           `json:"agent_city_name,omitempty" bson:"agent_city_name,omitempty"`
	AgentDistrictName    string                           `json:"agent_district_name,omitempty" bson:"agent_district_name,omitempty"`
	AgentVillageName     string                           `json:"agent_village_name,omitempty" bson:"agent_village_name,omitempty"`
	AgentAddress         string                           `json:"agent_address,omitempty" bson:"agent_address,omitempty"`
	AgentPhone           string                           `json:"agent_phone,omitempty" bson:"agent_phone,omitempty"`
	AgentMainMobilePhone string                           `json:"agent_main_mobile_phone,omitempty" bson:"agent_main_mobile_phone,omitempty"`
	StoreID              int                              `json:"store_id,omitempty" bson:"store_id,omitempty" binding:"required"`
	StoreCode            string                           `json:"store_code,omitempty" bson:"store_code,omitempty"`
	StoreName            string                           `json:"store_name,omitempty" bson:"store_name,omitempty"`
	StoreStatus          string                           `json:"store_status,omitempty" bson:"store_status,omitempty"`
	StoreOwner           string                           `json:"store_owner,omitempty" bson:"store_name,omitempty"`
	StoreEmail           string                           `json:"store_email,omitempty" bson:"store_email,omitempty"`
	StoreProvinceId      int                              `json:"store_province_id,omitempty" bson:"store_province_id,omitempty"`
	StoreProvinceName    string                           `json:"store_province_name,omitempty" bson:"store_province_name,omitempty"`
	StoreCityId          int                              `json:"store_city_id,omitempty" bson:"store_city_id,omitempty"`
	StoreCityName        string                           `json:"store_city_name,omitempty" bson:"store_city_name,omitempty"`
	StoreDistrictId      int                              `json:"store_district_id,omitempty" bson:"store_district_id,omitempty"`
	StoreDistrictName    string                           `json:"store_district_name,omitempty" bson:"store_district_name,omitempty"`
	StoreVillageName     string                           `json:"store_village_name,omitempty" bson:"store_village_name,omitempty"`
	StoreAddress         string                           `json:"store_address,omitempty" bson:"store_address,omitempty"`
	StorePhone           NullString                       `json:"store_phone,omitempty" bson:"store_phone,omitempty"`
	StoreMainMobilePhone string                           `json:"store_main_mobile_phone,omitempty" bson:"store_main_mobile_phone,omitempty"`
	BrandID              int                              `json:"brand_id,omitempty" bson:"brand_id,omitempty" binding:"required"`
	BrandName            string                           `json:"brand_name,omitempty" bson:"brand_name,omitempty"`
	UserID               int                              `json:"user_id,omitempty" bson:"user_id,omitempty" binding:"required"`
	UserFirstName        string                           `json:"user_first_name,omitempty" bson:"user_first_name,omitempty"`
	UserLastName         string                           `json:"user_last_name,omitempty" bson:"user_last_name,omitempty"`
	UserEmail            string                           `json:"user_email,omitempty" bson:"user_email,omitempty"`
	SalesmanID           int                              `json:"salesman_id,omitempty" bson:"salesman_id,omitempty"`
	SalesmanName         string                           `json:"salesman_name,omitempty" bson:"salesman_name,omitempty"`
	SalesmanEmail        string                           `json:"salesman_email,omitempty" bson:"salesman_email,omitempty"`
	VisitationID         int                              `json:"visitation_id,omitempty" bson:"visitation_id,omitempty"`
	OrderSourceID        int                              `json:"order_source_id,omitempty" bson:"order_source_id,omitempty" binding:"required"`
	OrderSourceName      string                           `json:"order_source_name,omitempty" bson:"order_source_name,omitempty"`
	OrderStatusID        int                              `json:"order_status_id,omitempty" bson:"order_status_id,omitempty"`
	OrderStatusName      string                           `json:"order_status_name,omitempty" bson:"order_status_name,omitempty"`
	SoCode               string                           `json:"so_code,omitempty" bson:"so_code,omitempty"`
	SoDate               string                           `json:"so_date,omitempty" bson:"so_date,omitempty" binding:"required"`
	SoRefCode            string                           `json:"so_ref_code,omitempty" bson:"so_ref_code,omitempty"`
	SoRefDate            string                           `json:"so_ref_date,omitempty" bson:"so_ref_date,omitempty" binding:"required"`
	GLong                float64                          `json:"g_long,omitempty" bson:"g_long,omitempty"`
	GLat                 float64                          `json:"g_lat,omitempty" bson:"g_lat,omitempty"`
	Note                 string                           `json:"note,omitempty" bson:"note,omitempty"`
	InternalComment      string                           `json:"internal_comment,omitempty" bson:"internal_comment,omitempty"`
	TotalAmount          float64                          `json:"total_amount,omitempty" bson:"total_amount,omitempty" binding:"required"`
	TotalTonase          float64                          `json:"total_tonase,omitempty" bson:"total_tonase,omitempty" binding:"required"`
	DeviceId             string                           `json:"device_id,omitempty" bson:"device_id,omitempty" binding:"required"`
	ReferralCode         string                           `json:"referral_code,omitempty" bson:"referral_code,omitempty"`
	StartCreatedDate     *time.Time                       `json:"start_created_date,omitempty" bson:"start_created_date,omitempty"`
	SalesOrderDetails    []*SalesOrderDetailStoreResponse `json:"sales_order_details" bson:"sales_order_details,omitempty"`
	CreatedAt            *time.Time                       `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type SalesOrderUpdateRequest struct {
	SoDate            string                           `json:"so_date,omitempty" bson:"so_date,omitempty"`
	SoCode            string                           `json:"so_code,omitempty" bson:"so_code,omitempty"`
	OrderStatusID     int                              `json:"order_status_id,omitempty" bson:"order_status_id,omitempty" binding:"required"`
	Reason            string                           `json:"reason,omitempty" bson:"reason,omitempty"`
	SalesOrderDetails []*SalesOrderDetailUpdateRequest `json:"sales_order_details" bson:"sales_order_details" binding:"required,dive,required"`
}

type SalesOrderChan struct {
	SalesOrder *SalesOrder
	Total      int64
	Error      error
	ErrorLog   *model.ErrorLog
	ID         int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type SalesOrdersChan struct {
	SalesOrders []*SalesOrder
	Total       int64
	ErrorLog    *model.ErrorLog
	Error       error
	ID          int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type SalesOrderRequest struct {
	ID                int    `json:"id,omitempty"`
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
	Keyword           string `json:"keyword,omitempty"`
	SoDate            string `json:"so_date,omitempty"`
	StoreProvinceID   int    `json:"store_province_id,omitempty"`
	StoreCityID       int    `json:"store_city_id,omitempty"`
	StoreDistrictID   int    `json:"store_district_id,omitempty"`
	StoreVillageID    int    `json:"store_village_id,omitempty"`
}

type SalesOrders struct {
	SalesOrders []*SalesOrder `json:"sales_orders,omitempty"`
	Total       int64         `json:"total,omitempty"`
}

type SalesOrderOpenSearchResponse struct {
	ID                   int                                   `json:"id,omitempty" bson:"id,omitempty"`
	AgentName            NullString                            `json:"agent_name,omitempty" bson:"agent_name,omitempty"`
	AgentEmail           NullString                            `json:"agent_email,omitempty" bson:"agent_email,omitempty"`
	AgentProvinceName    NullString                            `json:"agent_province_name,omitempty" bson:"agent_province_name,omitempty"`
	AgentCityName        NullString                            `json:"agent_city_name,omitempty" bson:"agent_city_name,omitempty"`
	AgentDistrictName    NullString                            `json:"agent_district_name,omitempty" bson:"agent_district_name,omitempty"`
	AgentVillageName     NullString                            `json:"agent_village_name,omitempty" bson:"agent_village_name,omitempty"`
	AgentAddress         NullString                            `json:"agent_address,omitempty" bson:"agent_address,omitempty"`
	AgentPhone           NullString                            `json:"agent_phone,omitempty" bson:"agent_phone,omitempty"`
	AgentMainMobilePhone NullString                            `json:"agent_main_mobile_phone,omitempty" bson:"agent_main_mobile_phone,omitempty"`
	StoreName            NullString                            `json:"store_name,omitempty" bson:"store_name,omitempty"`
	StoreCode            NullString                            `json:"store_code,omitempty" bson:"store_code,omitempty"`
	StoreEmail           NullString                            `json:"store_email,omitempty" bson:"store_email,omitempty"`
	StoreProvinceName    NullString                            `json:"store_province_name,omitempty" bson:"store_province_name,omitempty"`
	StoreCityName        NullString                            `json:"store_city_name,omitempty" bson:"store_city_name,omitempty"`
	StoreDistrictName    NullString                            `json:"store_district_name,omitempty" bson:"store_district_name,omitempty"`
	StoreVillageName     NullString                            `json:"store_village_name,omitempty" bson:"store_village_name,omitempty"`
	StoreAddress         NullString                            `json:"store_address,omitempty" bson:"store_address,omitempty"`
	StorePhone           NullString                            `json:"store_phone,omitempty" bson:"store_phone,omitempty"`
	StoreMainMobilePhone NullString                            `json:"store_main_mobile_phone,omitempty" bson:"store_main_mobile_phone,omitempty"`
	BrandID              int                                   `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	BrandName            string                                `json:"brand_name,omitempty" bson:"brand_name,omitempty"`
	UserFirstName        NullString                            `json:"user_first_name,omitempty" bson:"user_first_name,omitempty"`
	UserLastName         NullString                            `json:"user_last_name,omitempty" bson:"user_last_name,omitempty"`
	UserEmail            NullString                            `json:"user_email,omitempty" bson:"user_email,omitempty"`
	OrderSourceName      string                                `json:"order_source_name,omitempty" bson:"order_source_name,omitempty"`
	OrderStatusName      string                                `json:"order_status_name,omitempty" bson:"order_status_name,omitempty"`
	SoCode               string                                `json:"so_code,omitempty" bson:"so_code,omitempty"`
	SoDate               string                                `json:"so_date,omitempty" bson:"so_date,omitempty"`
	SoRefCode            NullString                            `json:"so_ref_code,omitempty" bson:"so_ref_code,omitempty"`
	SoRefDate            NullString                            `json:"so_ref_date,omitempty" bson:"so_ref_date,omitempty"`
	ReferralCode         NullString                            `json:"referral_code,omitempty" bson:"referral_code,omitempty"`
	GLat                 NullFloat64                           `json:"g_lat,omitempty" bson:"g_lat,omitempty"`
	GLong                NullFloat64                           `json:"g_long,omitempty" bson:"g_long,omitempty"`
	Note                 NullString                            `json:"note,omitempty" bson:"note,omitempty"`
	InternalComment      NullString                            `json:"internal_comment,omitempty" bson:"internal_comment,omitempty"`
	TotalAmount          float64                               `json:"total_amount,omitempty" bson:"total_amount,omitempty"`
	TotalTonase          float64                               `json:"total_tonase,omitempty" bson:"total_tonase,omitempty"`
	SalesOrderDetails    []*SalesOrderDetailOpenSearchResponse `json:"sales_order_details,omitempty" bson:"sales_order_details,omitempty"`
	SalesmanName         NullString                            `json:"salesman_name,omitempty" bson:"salesman_name,omitempty"`
	SalesmanEmail        NullString                            `json:"salesman_email,omitempty" bson:"salesman_email,omitempty"`
	CreatedAt            *time.Time                            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt            *time.Time                            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type SalesOrdersOpenSearchResponse struct {
	SalesOrders []*SalesOrderOpenSearchResponse `json:"sales_orders,omitempty"`
	Total       int64                           `json:"total,omitempty"`
}
type SalesOrderExportRequest struct {
	ID                int    `json:"id,omitempty"`
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
	SoRefCode         string `json:"so_ref_code,omitempty" bson:"so_ref_code,omitempty"`
	ProductID         int    `json:"product_id,omitempty"`
	CategoryID        int    `json:"category_id,omitempty"`
	SalesmanID        int    `json:"salesman_id,omitempty"`
	ProvinceID        int    `json:"province_id,omitempty"`
	CityID            int    `json:"city_id,omitempty"`
	DistrictID        int    `json:"district_id,omitempty"`
	VillageID         int    `json:"village_id,omitempty"`
	SoDate            string `json:"so_date,omitempty"`
	StoreProvinceID   int    `json:"store_province_id,omitempty"`
	StoreCityID       int    `json:"store_city_id,omitempty"`
	StoreDistrictID   int    `json:"store_district_id,omitempty"`
	StoreVillageID    int    `json:"store_village_id,omitempty"`
	StartCreatedAt    string `json:"start_created_at,omitempty"`
	EndCreatedAt      string `json:"end_created_at,omitempty"`
	UpdatedAt         string `json:"updated_at,omitempty"`
	UserID            int    `json:"user_id,omitempty"`
}

type SalesOrderCsvResponse struct {
	SOStatus          string
	SOSource          string
	KodeReferralOrder string
	OrderNo           string
	SONO              string
	SODate            string
	DistributorID     string
	DistributorName   string
	KodeSalesman      string
	Salesman          string
	TokoType          string
	KodeTokoDBO       string
	KodeToko          string
	TokoName          string
	KodeKecamatan     string
	Kecamatan         string
	KodeCity          string
	City              string
	KodeProvince      string
	Province          string
	BrandID           string
	BrandName         string
	SOAmount          string
	DOAmount          string
	OrderNotes        string
	InternalNotes     string
	AlasanCancel      string
	AlasanReject      string
	SORefDate         string
	CreatedDate       string
	LastUpdate        string
	UserIDCreated     string
	UserIDModified    string
}
type SalesOrderExportResponse struct {
	StatusCode int    `json:"status_code,omitempty"`
	UrlFile    string `json:"url_file,omitempty"`
}
