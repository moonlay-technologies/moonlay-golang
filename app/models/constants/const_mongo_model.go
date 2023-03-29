package constants

const (
	SALES_ORDER_TABLE_LOGS               = "sales_order_event_logs"
	SALES_ORDER_TABLE_DETAIL_LOGS        = "sales_order_detail_event_logs"
	SALES_ORDER_TABLE_JOURNEYS           = "sales_order_journeys"
	SALES_ORDER_DETAIL_TABLE_JOURNEYS    = "sales_order_detail_journeys"
	DELIVERY_ORDER_TABLE_LOGS            = "delivery_order_event_logs"
	DELIVERY_ORDER_TABLE_JOURNEYS        = "delivery_order_journeys"
	DELIVERY_ORDER_DETAIL_TABLE_JOURNEYS = "delivery_order_detail_journeys"
	SOSJ_UPLOAD_TABLE_HISTORIES          = "sosj_upload_histories"
	SOSJ_UPLOAD_ERROR_TABLE_LOGS         = "sosj_upload_error_logs"
	SO_UPLOAD_TABLE_HISTORIES            = "so_upload_histories"
	SO_UPLOAD_ERROR_TABLE_LOGS           = "so_upload_error_logs"
	SJ_UPLOAD_TABLE_HISTORIES            = "sj_upload_histories"
	SJ_UPLOAD_ERROR_TABLE_LOGS           = "sj_upload_error_logs"

	COLUMN_SALES_ORDER_CODE    = "so_code"
	COLUMN_DELIVERY_ORDER_CODE = "do_code"
	COLUMN_STATUS              = "status"
	COLUMN_ACTION              = "action"

	LOG_STATUS_MONGO_DEFAULT = "0"
	LOG_STATUS_MONGO_SUCCESS = "1"
	LOG_STATUS_MONGO_ERROR   = "2"

	EVENT_LOG_STATUS_0 = "in progress"
	EVENT_LOG_STATUS_1 = "success"
	EVENT_LOG_STATUS_2 = "failed"

	LOG_ACTION_MONGO_INSERT = "insert"
	LOG_ACTION_MONGO_UPDATE = "update"
	LOG_ACTION_MONGO_DELETE = "delete"

	UPLOAD_STATUS_HISTORY_SUCCESS     = "success"
	UPLOAD_STATUS_HISTORY_IN_PROGRESS = "in progress"
	UPLOAD_STATUS_HISTORY_FAILED      = "failed"
	UPLOAD_STATUS_HISTORY_UPLOADED    = "uploaded"
	UPLOAD_STATUS_HISTORY_ERR_UPLOAD  = "err_upload"
)
