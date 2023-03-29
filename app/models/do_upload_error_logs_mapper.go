package models

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (result *DoUploadErrorLog) DoUploadErrorLogsMap(line int, doUploadHistoryId, requestId, bulkCode string, errors []string, now *time.Time) {
	result.RequestId = requestId
	result.DoUploadHistoryId, _ = primitive.ObjectIDFromHex(doUploadHistoryId)
	result.BulkCode = bulkCode
	result.ErrorRowLine = int64(line)
	result.ErrorMessage = strings.Join(errors, ";")
	result.CreatedAt = now
	result.UpdatedAt = now
}

func (result *RowDataDoUploadErrorLog) RowDataDoUploadErrorLogMap(request map[string]string, agentName string, warehouseName string) {
	result.AgentID = request["IDDistributor"]
	result.AgentName = &agentName
	result.OrderCode = request["NoOrder"]
	result.DoDate = request["TanggalSJ"]
	result.DoNumber = request["NoSJ"]
	result.Note = request["Catatan"]
	result.InternalNote = request["CatatanInternal"]
	result.DriverName = request["NamaSupir"]
	result.VehicleNo = request["PlatNo"]
	result.BrandID = request["KodeMerk"]
	result.BrandName = request["NamaMerk"]
	result.ProductCode = request["KodeProduk"]
	result.ProductName = request["NamaProduk"]
	result.DeliveryQty = request["QTYShip"]
	result.ProductUnit = request["Unit"]
	result.WarehouseCode = request["KodeGudang"]
	result.WarehouseName = &warehouseName
}
