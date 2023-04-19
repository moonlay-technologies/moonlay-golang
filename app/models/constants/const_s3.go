package constants

const (
	S3_EXPORT_BUCKET               = "lambda-upload-srv"
	S3_EXPORT_ACL                  = "public-read"
	S3_EXPORT_CONTENT_DISPOSISTION = "attachment"
	S3_EXPORT_DO_PATH              = "lambda-upload-srv/order-service/export-delivery-orders"
	S3_EXPORT_DO_DETAIL_PATH       = "lambda-upload-srv/order-service/export-delivery-order-details"
	S3_EXPORT_SO_PATH              = "lambda-upload-srv/order-service/export-sales-orders"
	S3_EXPORT_SO_DETAIL_PATH       = "lambda-upload-srv/order-service/export-sales-orders-details"

	EXPORT_PARTIAL_DEFAULT = 50
)
