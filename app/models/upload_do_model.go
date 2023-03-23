package models

import "order-service/global/utils/model"

type UploadDOField struct {
	IDDistributor   string
	NoOrder         string
	TanggalSJ       string
	NoSJ            string
	Catatan         string
	CatatanInternal string
	NamaSupir       string
	PlatNo          string
	KodeMerk        string
	NamaMerk        string
	KodeProduk      string
	NamaProduk      string
	QTYShip         string
	Unit            string
	KodeGudang      string
}

type UploadDOFieldsChan struct {
	UploadDOFields []*UploadDOField
	Total          int64
	ErrorLog       *model.ErrorLog
	Error          error
	ID             int64 `json:"id,omitempty" bson:"id,omitempty"`
}
