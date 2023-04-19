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

	DO_STATUS_CANCEL = "SJCNCL"
	DO_STATUS_CLOSED = "SJCLS"
	DO_STATUS_OPEN   = "SJCR"

	ORDER_STATUS_OPEN      = "open"
	ORDER_STATUS_CANCELLED = "cancelled"
	ORDER_STATUS_CLOSED    = "closed"
	ORDER_STATUS_PARTIAL   = "partial"
	ORDER_STATUS_PENDING   = "pending"
	ORDER_STATUS_REJECTED  = "rejected"

	DATE_FORMAT_COMMON             = "2006-01-02"
	DATE_FORMAT_EXPORT             = "20060102-150405"
	DATE_FORMAT_EXPORT_CREATED_AT  = "2006-01-02 15:04:05"
	DATE_FORMAT_CODE_GENERATOR     = "20060102150405"
	DATE_TIME_FORMAT_COMON         = "2006-01-02 15:04:05"
	DATE_TIME_ZERO_HOUR_ADDITIONAL = " 00:00:00"

	CLAUSE_ID_VALIDATION = "id = %d AND deleted_at IS NULL"

	FILE_EXCEL_TYPE = "xlsx"
	FILE_CSV_TYPE   = "csv"
)

func UNMAPPED_TYPE_SORT_LIST() []string {
	return []string{
		"created_at",
		"updated_at",
		"do_date",
	}
}

func DELIVERY_ORDER_SORT_INT_LIST() []string {
	return []string{
		"order_status_id",
		"agent_id",
		"store_id",
		"product_id",
		"qty",
	}
}

func DELIVERY_ORDER_SORT_STRING_LIST() []string {
	return []string{
		"order_status.name",
		"do_ref_code",
		"do_code",
		"sales_order.store_code",
	}
}

func DELIVERY_ORDER_DETAIL_SORT_INT_LIST() []string {
	return []string{
		"order_status_id",
		"agent_id",
		"store_id",
		"product_id",
		"qty",
	}
}

func DELIVERY_ORDER_DETAIL_SORT_STRING_LIST() []string {
	return []string{
		"order_status.name",
		"do_code",
		"do_ref_code",
		"store_code",
	}
}

func SALES_ORDER_DETAIL_SORT_INT_LIST() []string {
	return []string{
		"order_status_id",
		"so_date",
		"store_id",
		"agent_id",
		"product_id",
	}
}

func SALES_ORDER_DETAIL_SORT_STRING_LIST() []string {
	return []string{
		"order_status.name",
		"so_ref_code",
		"so_code",
		"store_code",
		"store_name",
		"product_code",
	}
}

func SALES_ORDER_SORT_INT_LIST() []string {
	return []string{
		"order_status_id",
		"so_date",
		"store_id",
		"agent_id",
	}
}

func SALES_ORDER_SORT_STRING_LIST() []string {
	return []string{
		"order_status.name",
		"so_ref_code",
		"so_code",
		"store_code",
		"store_name",
	}
}

func DELIVERY_ORDER_EXPORT_SORT_INT_LIST() []string {
	return []string{
		"order_status_id",
		"do_date",
	}
}

func DELIVERY_ORDER_EXPORT_SORT_STRING_LIST() []string {
	return []string{
		"do_ref_code",
		"do_code",
		"do_ref_code",
		"store_code",
	}
}

func DELIVERY_ORDER_DETAIL_EXPORT_SORT_INT_LIST() []string {
	return []string{
		"order_status_id",
		"do_date",
	}
}

func DELIVERY_ORDER_DETAIL_EXPORT_SORT_STRING_LIST() []string {
	return []string{
		"do_code",
		"so_code",
		"do_ref_code",
		"store_code",
	}
}

func SALES_ORDER_DETAIL_EXPORT_SORT_INT_LIST() []string {
	return []string{
		"order_status_id",
		"product_id",
	}
}

func SALES_ORDER_DETAIL_EXPORT_SORT_STRING_LIST() []string {
	return []string{
		"so_code",
		"store_name",
	}
}

func SALES_ORDER_EXPORT_SORT_INT_LIST() []string {
	return []string{
		"order_status_id",
		"so_date",
	}
}

func SALES_ORDER_EXPORT_SORT_STRING_LIST() []string {
	return []string{
		"so_ref_code",
		"store_code",
	}
}
