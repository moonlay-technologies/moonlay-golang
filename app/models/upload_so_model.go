package models

import "order-service/global/utils/model"

type UploadSORequest struct {
	File string `json:"file,omitempty" binding:"required"`
}

type UploadSOField struct {
	IDDistributor     int
	KodeToko          int
	NamaToko          string
	IDSalesman        int
	NamaSalesman      string
	TanggalOrder      string
	NoOrder           string
	TanggalTokoOrder  string
	CatatanOrder      string
	CatatanInternal   string
	KodeMerk          int
	NamaMerk          string
	KodeProduk        string
	NamaProduk        string
	QTYOrder          int
	UnitProduk        string
	IDAlamat          int
	NamaAlamat        string
	IDUser            int
	SoUploadHistoryId string
	UploadType        string
}

type UploadSOFieldsChan struct {
	UploadSOFields []*UploadSOField
	Total          int64
	ErrorLog       *model.ErrorLog
	Error          error
	ID             int64 `json:"id,omitempty" bson:"id,omitempty"`
}
