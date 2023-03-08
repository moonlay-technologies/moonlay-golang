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
	Table              string
	ReqField           string
	Clause             string
	CustomMessage      string
	CustomResponseCode int
}

type MustActiveRequestChan struct {
	UniqueRequest []*MustActiveRequest
	Total         []int64
	Error         error
	ErrorLog      *model.ErrorLog
	ID            int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type MustEmptyValidationRequest struct {
	Table           string
	SelectedCollumn string
	Clause          string
	MessageFormat   string
}

type MustEmptyValidationRequestChan struct {
	UniqueRequest *MustEmptyValidationRequest
	Result        bool
	Message       string
	Error         error
	ErrorLog      *model.ErrorLog
	ID            int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type DateInputRequest struct {
	Field string
	Value string
}

type DateInputRequestChan struct {
	UniqueRequest *DateInputRequest
	Total         int64
	Error         error
	ErrorLog      *model.ErrorLog
	ID            int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type RequestIdValidationChan struct {
	Total    int64
	Error    error
	ErrorLog *model.ErrorLog
	ID       int64 `json:"id,omitempty" bson:"id,omitempty"`
}
