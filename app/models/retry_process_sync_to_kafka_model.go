package models

type SORetryProcessSyncToKafkaResponse struct {
	SalesOrderLogEventId string `json:"sales_order_event_logs_id,omitempty" bson:"sales_order_event_logs_id,omitempty"`
	Status               string `json:"status,omitempty" bson:"status,omitempty"`
	Message              string `json:"message,omitempty" bson:"message,omitempty"`
}

type DORetryProcessSyncToKafkaResponse struct {
	DeliveryOrderLogEventId string `json:"delivery_order_event_logs_id,omitempty" bson:"delivery_order_event_logs_id,omitempty"`
	Status                  string `json:"status,omitempty" bson:"status,omitempty"`
	Message                 string `json:"message,omitempty" bson:"message,omitempty"`
}
