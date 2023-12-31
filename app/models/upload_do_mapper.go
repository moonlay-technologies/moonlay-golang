package models

import "strconv"

func (result *UploadDOField) UploadDOFieldMap(request map[string]string, userId int, sjUploadHistoryId string) {
	result.IDDistributor, _ = strconv.Atoi(request["IDDistributor"])
	result.NoOrder = request["NoOrder"]
	result.NoSJ = request["NoSJ"]
	result.Catatan = request["Catatan"]
	result.CatatanInternal = request["CatatanInternal"]
	result.NamaSupir = request["NamaSupir"]
	result.PlatNo = request["PlatNo"]
	result.KodeMerk, _ = strconv.Atoi(request["KodeMerk"])
	result.NamaMerk = request["NamaMerk"]
	result.KodeProduk = request["KodeProduk"]
	result.NamaProduk = request["NamaProduk"]
	result.QTYShip, _ = strconv.Atoi(request["QTYShip"])
	result.Unit = request["Unit"]
	result.KodeGudang = request["KodeGudang"]
	result.IDUser = userId
	result.SjUploadHistoryId = sjUploadHistoryId
	return
}
