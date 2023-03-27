package models

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (result *SoUploadErrorLog) SoUploadErrorLogsMap(line int, soUploadHistoryId, requestId, bulkCode string, errors []string, now *time.Time) {
	result.RequestId = requestId
	result.SoUploadHistoryId, _ = primitive.ObjectIDFromHex(soUploadHistoryId)
	result.BulkCode = bulkCode
	result.ErrorRowLine = int64(line)
	result.ErrorMessage = strings.Join(errors, ";")
	result.CreatedAt = now
	result.UpdatedAt = now
}

func (result *RowDataSoUploadErrorLog) RowDataSoUploadErrorLogMap(item map[string]string) {
	namaToko := item["NamaToko"]
	namaSalesman := item["NamaSalesman"]
	namaMerk := item["NamaMerk"]
	namaAlamat := item["NamaAlamat"]

	result.AgentId = item["IDDistributor"]
	result.StoreCode = item["KodeToko"]
	result.StoreName = &namaToko
	result.SalesId = item["IDSalesman"]
	result.SalesName = &namaSalesman
	result.SoRefCode = item["NoOrder"]
	result.OrderDate = item["TanggalOrder"]
	result.OrderNote = item["CatatanOrder"]
	result.InternalNote = item["CatatanInternal"]
	result.BrandCode = item["KodeMerk"]
	result.BrandName = &namaMerk
	result.ProductCode = item["KodeProduk"]
	result.ProductName = item["NamaProduk"]
	result.OrderQty = item["QTYOrder"]
	result.ProductUnit = item["UnitProduk"]
	result.AddresId = item["IDAlamat"]
	result.Address = &namaAlamat
}
