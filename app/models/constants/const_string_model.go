package constants

const (
	CREATE_SALES_ORDER_CONSUMER        = "sales-order-consumer"
	UPDATE_SALES_ORDER_CONSUMER        = "update-sales-order-consumer"
	DELETE_SALES_ORDER_CONSUMER        = "delete-sales-order-consumer"
	UPDATE_SALES_ORDER_DETAIL_CONSUMER = "update-sales-order-detail-consumer"
	DELETE_SALES_ORDER_DETAIL_CONSUMER = "delete-sales-order-detail-consumer"

	CREATE_DELIVERY_ORDER_CONSUMER        = "create-delivery-order-consumer"
	UPDATE_DELIVERY_ORDER_CONSUMER        = "update-delivery-order-consumer"
	UPDATE_DELIVERY_ORDER_DETAIL_CONSUMER = "update-delivery-order-detail-consumer"
	DELETE_DELIVERY_ORDER_CONSUMER        = "delete-delivery-order-consumer"

	DELETE_DELIVERY_ORDER_DETAIL_CONSUMER = "delete-delivery-order-detail-consumer"

	SO_STATUS_APPV   = "APPV"
	SO_STATUS_REAPPV = "REAPPV"
	SO_STATUS_RJC    = "RJC"
	SO_STATUS_CNCL   = "CNCL"
	SO_STATUS_ORDPRT = "ORDPRT"
	SO_STATUS_ORDCLS = "ORDCLS"
	SO_STATUS_CLS    = "CLS"
	SO_STATUS_PEND   = "PEND"
	SO_STATUS_OPEN   = "OPEN"
	SO_STATUS_SJCR   = "SJCR"
	SO_STATUS_SJCLS  = "SJCLS"

	DO_STATUS_CNCL = "SJCNCL"
	DO_STATUS_CLS  = "SJCLS"
	DO_STATUS_OPEN = "SJCR"

	ORDER_STATUS_OPEN      = "open"
	ORDER_STATUS_CANCELLED = "cancelled"
	ORDER_STATUS_CLOSED    = "closed"
	ORDER_STATUS_PARTIAL   = "partial"
	ORDER_STATUS_PENDING   = "pending"
	ORDER_STATUS_REJECTED  = "rejected"

	DATE_FORMAT_COMMON         = "2006-01-02"
	DATE_FORMAT_EXPORT         = "20060102-150405"
	DATE_FORMAT_CODE_GENERATOR = "20060102150405"
	DATE_TIME_FORMAT_COMON     = "2006-01-02 15:04:05"

	CLAUSE_ID_VALIDATION = "id = %d AND deleted_at IS NULL"

	FILE_EXCEL_TYPE = "xlsx"
	FILE_CSV_TYPE   = "csv"
)
