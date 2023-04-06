package constants

const (
	ERROR_ACTION_NAME_CREATE = "create"
	ERROR_ACTION_NAME_GET    = "get"
	ERROR_ACTION_NAME_UPDATE = "update"
	ERROR_ACTION_NAME_DELETE = "delete"
	ERROR_ACTION_NAME_UPLOAD = "upload"

	ERROR_INVALID_PROCESS              = "Invalid Process"
	ERROR_BAD_REQUEST_INT_ID_PARAMS    = "Parameter 'id' harus bernilai integer"
	ERROR_BAD_REQUEST_INT_SO_ID_PARAMS = "sales order id harus bernilai integer"
	ERROR_DATA_NOT_FOUND               = "data not found"
	ERROR_INTERNAL_SERVER_1            = "Ada kesalahan, silahkan coba lagi nanti"

	ERROR_UPDATE_SO_MESSAGE = "Sales Order Has Delivery Order <result>, Please Delete it First"

	ERROR_QTY_CANT_NEGATIVE = "qty must be higher or equal 0"
)
