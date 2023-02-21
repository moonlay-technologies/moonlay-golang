package constants

const (
	SALES_ORDER_TABLE_LOGS    = "sales_order_logs"
	DELIVERY_ORDER_TABLE_LOGS = "delivery_order_logs"

	COLUMN_SALES_ORDER_CODE    = "so_code"
	COLUMN_DELIVERY_ORDER_CODE = "do_code"
	COLUMN_STATUS              = "status"
	COLUMN_ACTION              = "action"

	LOG_STATUS_MONGO_DEFAULT = "0"
	LOG_STATUS_MONGO_SUCCESS = "1"
	LOG_STATUS_MONGO_ERROR   = "2"

	LOG_ACTION_MONGO_INSERT = "insert"
	LOG_ACTION_MONGO_UPDATE = "update"
	LOG_ACTION_MONGO_DELETE = "delete"
)
