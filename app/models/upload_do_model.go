package models

import "order-service/global/utils/model"

type UploadDORequest struct {
	File string `json:"file,omitempty" binding:"required"`
}

type UploadDOField struct {
	IDDistributor     int
	NoOrder           string
	TanggalSJ         string
	NoSJ              string
	Catatan           string
	CatatanInternal   string
	NamaSupir         string
	PlatNo            string
	KodeMerk          int
	NamaMerk          string
	KodeProduk        string
	NamaProduk        string
	QTYShip           int
	Unit              string
	KodeGudang        string
	IDUser            int
	SjUploadHistoryId string
	BulkCode          string
	UploadType        string
	ErrorLine         int
}

type UploadDOFieldsChan struct {
	UploadDOFields []*UploadDOField
	Total          int64
	ErrorLog       *model.ErrorLog
	Error          error
	ID             int64 `json:"id,omitempty" bson:"id,omitempty"`
}
