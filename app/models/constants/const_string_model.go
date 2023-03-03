package constants

const (
	CREATE_SALES_ORDER_CONSUMER        = "sales-order-consumer"
	UPDATE_SALES_ORDER_CONSUMER        = "update-sales-order-consumer"
	DELETE_SALES_ORDER_CONSUMER        = "delete-sales-order-consumer"
	UPDATE_SALES_ORDER_DETAIL_CONSUMER = "update-sales-order-detail-consumer"

	CREATE_DELIVERY_ORDER_CONSUMER = "create-delivery-order-consumer"
	UPDATE_DELIVERY_ORDER_CONSUMER = "update-delivery-order-consumer"
	DELETE_DELIVERY_ORDER_CONSUMER = "delete-delivery-order-consumer"

	UPDATE_SO_STATUS_APPV = "APPV"
	UPDATE_SO_STATUS_RJC  = "RJC"
	UPDATE_SO_STATUS_CNCL = "CNCL"

	ORDER_STATUS_OPEN      = "open"
	ORDER_STATUS_CANCELLED = "cancelled"
	ORDER_STATUS_CLOSED    = "closed"
	ORDER_STATUS_PARTIAL   = "partial"
	ORDER_STATUS_PENDING   = "pending"
	ORDER_STATUS_REJECTED  = "rejected"
)
