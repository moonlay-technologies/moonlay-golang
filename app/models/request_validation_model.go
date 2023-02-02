package models

import "poc-order-service/global/utils/model"

type UniqueRequest struct {
	Table string
	Field string
	Value interface{}
}

type UniqueRequestChan struct {
	UniqueRequest *UniqueRequest
	Total         int64
	Error         error
	ErrorLog      *model.ErrorLog
	ID            int64 `json:"id,omitempty" bson:"id,omitempty"`
}
