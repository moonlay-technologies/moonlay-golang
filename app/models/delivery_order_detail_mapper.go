package models

func (deliveryOrderDetail *DeliveryOrderDetailOpenSearchDetailResponse) DeliveryOrderDetailOpenSearchResponseMap(request *DeliveryOrderDetail) {
	deliveryOrderDetail.SoDetailID = request.SoDetailID
	deliveryOrderDetail.Qty = request.Qty
	return
}
