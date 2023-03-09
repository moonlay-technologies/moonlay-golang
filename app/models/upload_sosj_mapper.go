package models

import (
	"strconv"
	"strings"
)

func (result *UploadSOSJField) UploadSOSJFieldMap(request map[string]string, idDistributor int) {
	result.IDDistributor = idDistributor
	result.Status = request["_1"]
	result.TglSuratJalan = request["_3"]
	result.KodeTokoDBO, _ = strconv.Atoi(request["_4"])
	result.IDMerk, _ = strconv.Atoi(request["_5"])
	result.KodeProdukDBO, _ = strconv.Atoi(request["_6"])
	result.Qty, _ = strconv.Atoi(request["_7"])
	result.Unit, _ = strconv.Atoi(request["_8"])
	result.NamaSupir = request["_9"]
	result.PlatNo = request["_10"]
	result.KodeGudang, _ = strconv.Atoi(request["_11"])
	result.IDSalesman, _ = strconv.Atoi(request["_12"])
	result.IDAlamat = request["_13"]
	result.Catatan = request["_14"]
	result.CatatanInternal = strings.ReplaceAll(request["_15"], "\r", "")
}
