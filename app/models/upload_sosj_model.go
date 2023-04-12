package models

import "order-service/global/utils/model"

type UploadSOSJRequest struct {
	File string `json:"file,omitempty" binding:"required"`
}

type UploadSOSJField struct {
	IDDistributor       int
	Status              string
	NoSuratJalan        string
	TglSuratJalan       string
	KodeTokoDBO         string
	IDMerk              int
	KodeProdukDBO       string
	Qty                 int
	Unit                int
	NamaSupir           string
	PlatNo              string
	KodeGudang          int
	IDSalesman          int
	IDAlamat            string
	Catatan             string
	CatatanInternal     string
	IDUser              int
	SosjUploadHistoryId string
	BulkCode            string
	UploadType          string
	ErrorLine           int
	RowData             RowDataSosjUploadErrorLog
}

type UploadSOSJFieldChan struct {
	UploadSOSJFields []*UploadSOSJField
	Total            int64
	ErrorLog         *model.ErrorLog
	Error            error
	ID               int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type RetryUploadSOSJResponse struct {
	SosjUploadHistoryId string `json:"sosj_upload_history_id"`
	Message             string `json:"message"`
	Status              string `json:"status"`
}
