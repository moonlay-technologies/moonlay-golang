package constants

func ExcelCollumnMapper() []string {
	return []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ",
		"BA", "BB", "BC", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BK", "BL", "BM", "BN", "BO", "BP", "BQ", "BR", "BS", "BT", "BU", "BV", "BW", "BX", "BY", "BZ",
		"CA", "CB", "CC", "CD", "CE", "CF", "CG", "CH", "CI", "CJ", "CK", "CL", "CM", "CN", "CO", "CP", "CQ", "CR", "CS", "CT", "CU", "CV", "CW", "CX", "CY", "CZ",
	}
}

func DELIVERY_ORDER_EXPORT_HEADER() []interface{} {
	return []interface{}{
		"DOStatus",
		"DODate",
		"SJNo",
		"DONO",
		"OrderNo",
		"SODate",
		"SONO",
		"SOSource",
		"AgentID",
		"AgentName",
		"GudangID",
		"GudangName",
		"BrandID",
		"BrandName",
		"KodeSalesman",
		"Salesman",
		"KategoriToko",
		"KodeTokoDBO",
		"KodeToko",
		"NamaToko",
		"KodeKecamatan",
		"Kecamatan",
		"Kode City",
		"City",
		"KodeProvince",
		"Province",
		"DOAmount",
		"NamaSupir",
		"PlatNo",
		"Catatan",
		"CreatedDate",
		"LastUpdate",
		"UserIDCreated",
		"UserIDModified"}
}

func DELIVERY_ORDER_DETAIL_EXPORT_HEADER() []interface{} {
	return []interface{}{
		"DOStatus",
		"DODate",
		"SJNo",
		"DONO",
		"SODate",
		"SONO",
		"SOSource",
		"AgentID",
		"AgentName",
		"GudangID",
		"GudangName",
		"BrandID",
		"BrandName",
		"KodeSalesman",
		"Salesman",
		"KategoriToko",
		"KodeTokoDBO",
		"KodeToko",
		"NamaToko",
		"KodeKecamatan",
		"Kecamatan",
		"Kode City",
		"City",
		"KodeProvince",
		"Province",
		"BrandID",
		"BrandName",
		"KategoriID_L1",
		"KategoriName_L1",
		"KategoriID_Last",
		"KategoriName_Last",
		"ItemCode",
		"ItemName",
		"UnitSatuan",
		"Price",
		"SOQty",
		"SOQty_Sisa",
		"DOQty",
		"DOAmount",
		"CreatedDate",
		"LastUpdate",
		"UserIDCreated",
		"UserIDModified"}
}

func SALES_ORDER_EXPORT_HEADER() []interface{} {
	return []interface{}{
		"SOStatus",
		"SOSource",
		"KodeReferralOrder",
		"OrderNo",
		"SONO",
		"SODate",
		"DistributorID",
		"DistributorName",
		"KodeSalesman",
		"Salesman",
		"TokoType",
		"KodeTokoDBO",
		"KodeToko",
		"TokoName",
		"KodeKecamatan",
		"Kecamatan",
		"KodeCity",
		"City",
		"KodeProvince",
		"Province",
		"BrandID",
		"BrandName",
		"SOAmount",
		"DOAmount",
		"OrderNotes",
		"InternalNotes",
		"AlasanCancel",
		"AlasanReject",
		"SORefDate",
		"CreatedDate",
		"LastUpdate",
		"UserIDCreated",
		"UserIDModified"}
}

func SALES_ORDER_DETAIL_EXPORT_HEADER() []interface{} {
	return []interface{}{
		"SOStatus",
		"SOSource",
		"KodeReferralOrder",
		"OrderNo",
		"SONO",
		"SODate",
		"DistributorID",
		"DistributorName",
		"KodeSalesman",
		"Salesman",
		"TokoType",
		"KodeTokoDBO",
		"KodeToko",
		"TokoName",
		"KodeKecamatan",
		"Kecamatan",
		"KodeCity",
		"City",
		"KodeProvince",
		"Province",
		"BrandID",
		"BrandName",
		"KategoriID_L1",
		"KategoriName_L1",
		"KategoriID_Last",
		"KategoriName_Last",
		"ItemCode",
		"ItemName",
		"UnitSatuan",
		"Price",
		"SOQty",
		"SOAmount",
		"DOQty",
		"DOAmount",
		"SOItemStatus",
		"AlasanCancel",
		"CreatedDate",
		"LastUpdate",
		"UserIDCreated",
		"UserIDModified",
		"OrderNotes",
		"InternalNotes",}
}
