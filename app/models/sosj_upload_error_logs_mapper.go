package models

import (
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (result *SosjUploadErrorLog) SosjUploadErrorLogsMap(line int, agentId, sosjUploadHistoryId, requestId, agentName, bulkCode string, errors []string, now *time.Time) {
	result.RequestId = requestId
	result.SosjUploadHistoryId, _ = primitive.ObjectIDFromHex(sosjUploadHistoryId)
	result.BulkCode = bulkCode
	result.ErrorRowLine = int64(line)
	result.ErrorMessage = strings.Join(errors, ";")
	result.CreatedAt = now
	result.UpdatedAt = now
}

func (result *RowDataSosjUploadErrorLog) RowDataSosjUploadErrorLogMap(rowData RowDataSosjUploadErrorLog, item map[string]string) {
	result.AgentId = item["_0"]
	result.AgentName = rowData.AgentName
	result.SjStatus = item["_1"]
	result.SjNo = item["_2"]
	result.SjDate = item["_3"]
	result.StoreCode = item["_4"]
	result.StoreName = rowData.StoreName
	result.BrandId = item["_5"]
	result.BrandName = rowData.BrandName
	result.ProductCode = item["_6"]
	result.ProductName = rowData.ProductName
	result.DeliveryQty = item["_7"]
	result.ProductUnit = item["_8"]
	result.DriverName = item["_9"]
	result.VehicleNo = item["_10"]
	result.WhCode = item["_11"]
	result.WhName = rowData.WhName
	result.SalesmanId = item["_12"]
	result.SalesName = rowData.SalesName
	result.AddresId = item["_13"]
	result.Address = rowData.Address
	result.Note = item["_14"]
	result.InternalNote = item["_15"]
}

func (result *RowDataSosjUploadErrorLog) RowDataSosjUploadErrorLogMap2(rowData RowDataSosjUploadErrorLog, item *UploadSOSJField) {
	result.AgentId = strconv.Itoa(item.IDDistributor)
	result.AgentName = rowData.AgentName
	result.SjStatus = item.Status
	result.SjNo = item.NoSuratJalan
	result.SjDate = item.TglSuratJalan
	result.StoreCode = item.KodeTokoDBO
	result.StoreName = rowData.StoreName
	result.BrandId = strconv.Itoa(item.IDMerk)
	result.BrandName = rowData.BrandName
	result.ProductCode = item.KodeProdukDBO
	result.ProductName = rowData.ProductName
	result.DeliveryQty = strconv.Itoa(item.Qty)
	result.ProductUnit = strconv.Itoa(item.Unit)
	result.DriverName = item.NamaSupir
	result.VehicleNo = item.PlatNo
	result.WhCode = strconv.Itoa(item.KodeGudang)
	result.WhName = rowData.WhName
	result.SalesmanId = strconv.Itoa(item.IDSalesman)
	result.SalesName = rowData.SalesName
	result.AddresId = item.IDAlamat
	result.Address = rowData.Address
	result.Note = item.Catatan
	result.InternalNote = item.CatatanInternal
}
