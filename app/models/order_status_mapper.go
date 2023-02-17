package models

func (orderStatus *OrderStatusOpenSearchResponse) OrderStatusOpenSearchResponseMap(request *OrderStatus) {
	orderStatus.ID = request.ID
	orderStatus.Name = request.Name
	return
}
