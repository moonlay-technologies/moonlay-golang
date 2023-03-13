package models

import (
	"strconv"
	"strings"
)

func (result *UploadSOField) UploadSOFieldMap(request map[string]string) {
	result.IDDistributor, _ = strconv.Atoi(request["IDDistributor"])
	result.KodeToko, _ = strconv.Atoi(request["KodeToko"])
	result.NamaToko = request["NamaToko"]
	result.IDSalesman, _ = strconv.Atoi(request["IDSalesman"])
	result.NamaSalesman = request["NamaSalesman"]
	result.NoOrder = request["NoOrder"]
	result.CatatanOrder = request["CatatanOrder"]
	result.CatatanInternal = request["CatatanInternal"]
	result.KodeMerk, _ = strconv.Atoi(request["KodeMerk"])
	result.NamaMerk = request["NamaMerk"]
	result.KodeProduk = request["KodeProduk"]
	result.NamaProduk = request["NamaProduk"]
	result.QTYOrder, _ = strconv.Atoi(request["QTYOrder"])
	result.UnitProduk = request["UnitProduk"]
	result.IDAlamat, _ = strconv.Atoi(request["IDAlamat"])
	result.NamaAlamat = strings.ReplaceAll(request["NamaAlamat"], "\r", "")
}
