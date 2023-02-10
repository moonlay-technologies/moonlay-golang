package models

import "order-service/global/utils/model"

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

type MustActiveRequest struct {
	Table    string
	ReqField string
	Clause   string
	// Column1  string
	// Value1   interface{}
	// Column2  string
	// Value2   interface{}
}

type MustActiveRequestChan struct {
	UniqueRequest *MustActiveRequest
	Total         int64
	Error         error
	ErrorLog      *model.ErrorLog
	ID            int64 `json:"id,omitempty" bson:"id,omitempty"`
}