package models

import "strconv"

func (result *UploadDOField) UploadDOFieldMap(request map[string]string, userId int) {
	result.IDDistributor, _ = strconv.Atoi(request["IDDistributor"])
	result.NoOrder = request["NoOrder"]
	result.TanggalSJ = request["TanggalSJ"]
	result.NoSJ = request["NoSJ"]
	result.Catatan = request["Catatan"]
	result.CatatanInternal = request["CatatanInternal"]
	result.NamaSupir = request["NamaSupir"]
	result.PlatNo = request["PlatNo"]
	result.KodeMerk, _ = strconv.Atoi(request["KodeMerk"])
	result.NamaMerk = request["NamaMerk"]
	result.KodeProduk, _ = strconv.Atoi(request["KodeProduk"])
	result.NamaProduk = request["NamaProduk"]
	result.QTYShip, _ = strconv.Atoi(request["QTYShip"])
	result.Unit = request["Unit"]
	result.KodeGudang, _ = strconv.Atoi(request["KodeGudang"])
	result.IDUser = userId
	return
}
