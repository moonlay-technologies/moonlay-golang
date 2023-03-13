package models

import "order-service/global/utils/model"

type UploadSOField struct {
	IDDistributor    int
	KodeToko         int
	NamaToko         string
	IDSalesman       int
	NamaSalesman     string
	TanggalOrder     string
	NoOrder          string
	TanggalTokoOrder string
	CatatanOrder     string
	CatatanInternal  string
	KodeMerk         int
	NamaMerk         string
	KodeProduk       string
	NamaProduk       string
	QTYOrder         int
	UnitProduk       string
	IDAlamat         int
	NamaAlamat       string
}

type UploadSOFieldsChan struct {
	UploadSOFields []*UploadSOField
	Total          int64
	ErrorLog       *model.ErrorLog
	Error          error
	ID             int64 `json:"id,omitempty" bson:"id,omitempty"`
}
