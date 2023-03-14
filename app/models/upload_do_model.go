package models

import "order-service/global/utils/model"

type UploadDORequest struct {
	File string `json:"file,omitempty" binding:"required"`
}

type UploadDOField struct {
	IDDistributor   int
	NoOrder         string
	TanggalSJ       string
	NoSJ            string
	Catatan         string
	CatatanInternal string
	NamaSupir       string
	PlatNo          string
	KodeMerk        int
	NamaMerk        string
	KodeProduk      int
	NamaProduk      string
	QTYShip         int
	Unit            string
	KodeGudang      int
	IDUser          int
}

type UploadDOFieldsChan struct {
	UploadDOFields []*UploadDOField
	Total          int64
	ErrorLog       *model.ErrorLog
	Error          error
	ID             int64 `json:"id,omitempty" bson:"id,omitempty"`
}
