package models

func (uom *UomOpenSearchResponse) UomOpenSearchResponseMap(request *Uom) {
	uom.Name = request.Name
	uom.Code = request.Code
	return
}
