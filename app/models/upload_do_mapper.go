package models

func (result *UploadDOField) UploadDOFieldMap(request map[string]string) {
	result.IDDistributor = request["IDDistributor"]
	result.NoOrder = request["NoOrder"]
	result.TanggalSJ = request["TanggalSJ"]
	result.NoSJ = request["NoSJ"]
	result.Catatan = request["Catatan"]
	result.CatatanInternal = request["CatatanInternal"]
	result.NamaSupir = request["NamaSupir"]
	result.PlatNo = request["PlatNo"]
	result.KodeMerk = request["KodeMerk"]
	result.NamaMerk = request["NamaMerk"]
	result.KodeProduk = request["KodeProduk"]
	result.NamaProduk = request["NamaProduk"]
	result.QTYShip = request["QTYShip"]
	result.Unit = request["Unit"]
	result.KodeGudang = request["KodeGudang"]
}
