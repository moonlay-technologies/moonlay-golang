package models

import "order-service/global/utils/model"

type UploadSOSJField struct {
	IDDistributor   int
	Status          string
	NoSuratJalan    string
	TglSuratJalan   string
	KodeTokoDBO     int
	IDMerk          int
	KodeProdukDBO   int
	Qty             int
	Unit            int
	NamaSupir       string
	PlatNo          string
	KodeGudang      int
	IDSalesman      int
	IDAlamat        string
	Catatan         string
	CatatanInternal string
}

type UploadSOSJFieldChan struct {
	UploadSOSJFields []*UploadSOSJField
	Total            int64
	ErrorLog         *model.ErrorLog
	Error            error
	ID               int64 `json:"id,omitempty" bson:"id,omitempty"`
}
