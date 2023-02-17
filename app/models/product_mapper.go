package models

func (product *ProductOpenSearchResponse) ProductOpenSearchResponseMap(request *Product) {
	product.ID = request.ID
	product.Sku = request.Sku
	product.AliasSku = request.AliasSku
	product.ProductName = request.ProductName
	product.Description = request.Description
	product.CategoryID = request.CategoryID
	return
}
