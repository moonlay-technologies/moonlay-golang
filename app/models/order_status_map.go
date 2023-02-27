package models

func (salesOrder *SalesOrder) SalesOrderStatusChanMap(request *OrderStatusChan) {
	salesOrder.OrderStatus = request.OrderStatus
	salesOrder.OrderStatusID = request.OrderStatus.ID
	salesOrder.OrderStatusName = request.OrderStatus.Name
	return
}

func (salesOrderDetail *SalesOrderDetail) SalesOrderDetailStatusChanMap(request *OrderStatusChan) {
	salesOrderDetail.OrderStatus = request.OrderStatus
	salesOrderDetail.OrderStatusID = request.OrderStatus.ID
	salesOrderDetail.OrderStatusName = request.OrderStatus.Name
	return
}
