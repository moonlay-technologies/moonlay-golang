package constants

const (
	SALES_ORDERS_PATH                 = "sales-orders"
	SALES_ORDER_DETAILS_PATH          = "sales-order-details"
	DELIVERY_ORDERS_PATH              = "delivery-orders"
	AGENT_PATH                        = "agent"
	STORES_PATH                       = "stores"
	SALESMANS_PATH                    = "salesman"
	HEALTH_CHECK_PATH                 = "health-check"
	HOST_TO_HOST_PATH                 = "h2h"
	UPLOAD_SOSJ_PATH                  = "upload-sosj"
	UPLOAD_DO_PATH                    = "upload-delivery-orders"
	UPLOAD_SO_PATH                    = "upload-sales-orders"
	SOSJ_PATH                         = "sosj"
	UPLOAD_HISTORIES_PATH             = "upload-histories"
	DELIVERY_ORDER_EXPORT_PATH        = "https://lambda-upload-srv.s3.ap-southeast-1.amazonaws.com/" + S3_EXPORT_DO_PATH
	DELIVERY_ORDER_DETAIL_EXPORT_PATH = "https://lambda-upload-srv.s3.ap-southeast-1.amazonaws.com/" + S3_EXPORT_DO_DETAIL_PATH
	SALES_ORDER_EXPORT_PATH           = "https://lambda-upload-srv.s3.ap-southeast-1.amazonaws.com/" + S3_EXPORT_SO_PATH
	SALES_ORDER_DETAIL_EXPORT_PATH    = "https://lambda-upload-srv.s3.ap-southeast-1.amazonaws.com/" + S3_EXPORT_SO_DETAIL_PATH
)
