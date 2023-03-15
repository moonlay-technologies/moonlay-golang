package models

import (
	"strings"
	"time"
)

func (result *SoUploadErrorLog) SoUploadErrorLogsMap(request map[string]string, agentName string, line int, errors []string, now *time.Time) {
	result.Address = strings.ReplaceAll(request["NamaAlamat"], "\r", "")
	result.AddressId = request["IDAlamat"]
	result.AgentId = request["IDDistributor"]
	result.AgentName = agentName
	result.BrandCode = request["KodeMerk"]
	result.BrandName = request["NamaMerk"]
	result.SalesId = request["IDSalesman"]
	result.SalesName = request["NamaSalesman"]
	result.StoreCode = request["KodeToko"]
	result.StoreName = request["NamaToko"]
	result.StoreOrderAt = request["TanggalTokoOrder"]
	result.InternalNote = request["CatatanInternal"]
	result.OrderCode = request["NoOrder"]
	result.OrderDate = request["TanggalOrder"]
	result.OrderNote = request["CatatanOrder"]
	result.Line = line
	result.ErrorLog = strings.Join(errors, ";")
	result.CreatedAt = now
	result.UpdatedAt = now
}
