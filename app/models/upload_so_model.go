package models

import "order-service/global/utils/model"

type UploadSORequest struct {
	File string `json:"file,omitempty" binding:"required"`
}

type UploadSOField struct {
	IDDistributor     int
	KodeToko          string
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
	BulkCode          string
	UploadType        string
	ErrorLine         int
}

type UploadSOFieldsChan struct {
	UploadSOFields []*UploadSOField
	Total          int64
	ErrorLog       *model.ErrorLog
	Error          error
	ID             int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type RetryUploadSOResponse struct {
	SoUploadHistoryId string `json:"so_upload_history_id"`
	Message           string `json:"message"`
	Status            string `json:"status"`
}
