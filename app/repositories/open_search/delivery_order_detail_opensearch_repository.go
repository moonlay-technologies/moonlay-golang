package repositories

import (
	"encoding/json"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/global/utils/helper"
	"order-service/global/utils/model"
	"order-service/global/utils/opensearch_dbo"
)

type DeliveryOrderDetailOpenSearchRepositoryInterface interface {
	Create(request *models.DeliveryOrderDetail, result chan *models.DeliveryOrderDetailChan)
	Get(request *models.DeliveryOrderDetailOpenSearchRequest, result chan *models.DeliveryOrderDetailsChan)
	GetByDoID(request *models.DeliveryOrderRequest, resultChan chan *models.DeliveryOrderDetailsChan)
	GetBySalesOrderID(request *models.DeliveryOrderDetailOpenSearchRequest, resultChan chan *models.DeliveryOrderDetailsChan)
	GetBySalesmanID(request *models.DeliveryOrderDetailOpenSearchRequest, resultChan chan *models.DeliveryOrderDetailsChan)
	GetBySalesmansID(request *models.DeliveryOrderDetailOpenSearchRequest, resultChan chan *models.DeliveryOrderDetailsChan)
	GetByStoreID(request *models.DeliveryOrderDetailOpenSearchRequest, resultChan chan *models.DeliveryOrderDetailsChan)
	GetByAgentID(request *models.DeliveryOrderDetailOpenSearchRequest, resultChan chan *models.DeliveryOrderDetailsChan)
	generateDeliveryOrderQueryOpenSearchResult(openSearchQueryJson []byte) (*models.DeliveryOrderDetails, *model.ErrorLog)
	generateDeliveryOrderQueryOpenSearchTermRequest(term_field string, term_value interface{}, request *models.DeliveryOrderDetailOpenSearchRequest) []byte
	generateDeliveryOrderQueryOpenSearchByQueryParamTermRequest(term_field string, term_value interface{}, request *models.DeliveryOrderDetailOpenSearchRequest) []byte
}

type deliveryOrderDetailOpenSearch struct {
	openSearch opensearch_dbo.OpenSearchClientInterface
}

func InitDeliveryOrderDetailOpenSearchRepository(openSearch opensearch_dbo.OpenSearchClientInterface) DeliveryOrderDetailOpenSearchRepositoryInterface {
	return &deliveryOrderDetailOpenSearch{
		openSearch: openSearch,
	}
}

func (r *deliveryOrderDetailOpenSearch) Create(request *models.DeliveryOrderDetail, resultChan chan *models.DeliveryOrderDetailChan) {
	response := &models.DeliveryOrderDetailChan{}
	deliveryOrderDetailJson, _ := json.Marshal(request)
	_, err := r.openSearch.CreateDocument(constants.DELIVERY_ORDER_DETAILS_INDEX, request.DoDetailCode, deliveryOrderDetailJson)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	response.DeliveryOrderDetail = request
	resultChan <- response
	return
}

func (r *deliveryOrderDetailOpenSearch) Get(request *models.DeliveryOrderDetailOpenSearchRequest, resultChan chan *models.DeliveryOrderDetailsChan) {
	response := &models.DeliveryOrderDetailsChan{}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchTermRequest("", "", request)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrderDetails = result.DeliveryOrderDetails
	response.Total = result.Total
	resultChan <- response
	return
}

func (r *deliveryOrderDetailOpenSearch) GetByDoID(request *models.DeliveryOrderRequest, resultChan chan *models.DeliveryOrderDetailsChan) {
	response := &models.DeliveryOrderDetailsChan{}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchTermRequest("id", request.ID, nil)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrderDetails = result.DeliveryOrderDetails
	resultChan <- response
	return
}

func (r *deliveryOrderDetailOpenSearch) GetBySalesOrderID(request *models.DeliveryOrderDetailOpenSearchRequest, resultChan chan *models.DeliveryOrderDetailsChan) {
	response := &models.DeliveryOrderDetailsChan{}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchTermRequest("sales_order_id", request.SalesOrderID, request)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrderDetails = result.DeliveryOrderDetails
	resultChan <- response
	return
}

func (r *deliveryOrderDetailOpenSearch) GetBySalesmanID(request *models.DeliveryOrderDetailOpenSearchRequest, resultChan chan *models.DeliveryOrderDetailsChan) {
	response := &models.DeliveryOrderDetailsChan{}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchTermRequest("sales_order.user_id", request.SalesmanID, request)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrderDetails = result.DeliveryOrderDetails
	resultChan <- response
	return
}

func (r *deliveryOrderDetailOpenSearch) GetBySalesmansID(request *models.DeliveryOrderDetailOpenSearchRequest, resultChan chan *models.DeliveryOrderDetailsChan) {
	response := &models.DeliveryOrderDetailsChan{}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchByQueryParamTermRequest("", "", request)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrderDetails = result.DeliveryOrderDetails
	resultChan <- response
	return
}

func (r *deliveryOrderDetailOpenSearch) GetByStoreID(request *models.DeliveryOrderDetailOpenSearchRequest, resultChan chan *models.DeliveryOrderDetailsChan) {
	response := &models.DeliveryOrderDetailsChan{}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchTermRequest("store_id", request.StoreID, request)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrderDetails = result.DeliveryOrderDetails
	resultChan <- response
	return
}

func (r *deliveryOrderDetailOpenSearch) GetByAgentID(request *models.DeliveryOrderDetailOpenSearchRequest, resultChan chan *models.DeliveryOrderDetailsChan) {
	response := &models.DeliveryOrderDetailsChan{}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchTermRequest("agent_id", request.AgentID, request)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrderDetails = result.DeliveryOrderDetails
	resultChan <- response
	return
}

func (r *deliveryOrderDetailOpenSearch) generateDeliveryOrderQueryOpenSearchTermRequest(term_field string, term_value interface{}, request *models.DeliveryOrderDetailOpenSearchRequest) []byte {
	openSearchQuery := map[string]interface{}{}
	openSearchDetailQuery := map[string]interface{}{}
	openSearchDetailBoolQuery := map[string]interface{}{}

	filters := []map[string]interface{}{}
	musts := []map[string]interface{}{}

	if len(term_field) > 0 && term_value != nil {
		filter := map[string]interface{}{
			"term": map[string]interface{}{
				term_field: term_value,
			},
		}

		filters = append(filters, filter)
	}

	if request != nil {

		if request.Page > 0 {
			page := request.PerPage * (request.Page - 1)
			openSearchQuery["size"] = request.PerPage
			openSearchQuery["from"] = page
		}

		if request.DeliveryOrderID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"delivery_order_id": request.DeliveryOrderID,
				},
			}

			filters = append(filters, filter)
		}

		if request.ProductID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"product_id": request.ProductID,
				},
			}

			filters = append(filters, filter)
		}

		if request.UomID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"uom_id": request.UomID,
				},
			}

			filters = append(filters, filter)
		}

		if request.UomName != "" {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"uom_name": request.UomName,
				},
			}

			filters = append(filters, filter)
		}

		if request.UomCode != "" {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"uom_code": request.UomCode,
				},
			}

			filters = append(filters, filter)
		}

		if request.OrderStatusID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"order_status_id": request.OrderStatusID,
				},
			}

			filters = append(filters, filter)
		}

		if request.OrderStatusName != "" {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"order_status_name": request.UomCode,
				},
			}

			filters = append(filters, filter)
		}

		if request.DoDetailCode != "" {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"do_detail_code": request.DoDetailCode,
				},
			}

			filters = append(filters, filter)
		}

		if request.Qty != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"qty": request.Qty,
				},
			}

			filters = append(filters, filter)
		}

		if request.SoDetail.SentQty != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"sent_qty": request.SoDetail.SentQty,
				},
			}

			filters = append(filters, filter)
		}

		if request.SoDetail.ResidualQty != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"residual_qty": request.SoDetail.ResidualQty,
				},
			}

			filters = append(filters, filter)
		}

		if request.SoDetail.Price != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"price": request.SoDetail.Price,
				},
			}

			filters = append(filters, filter)
		}

		if len(request.StartCreatedAt) > 0 && len(request.EndCreatedAt) > 0 {
			filter := map[string]interface{}{
				"range": map[string]interface{}{
					"created_at": map[string]interface{}{
						"gte": request.StartCreatedAt,
						"lte": request.EndCreatedAt,
					},
				},
			}

			filters = append(filters, filter)
		}

		if len(request.StartDoDate) > 0 && len(request.EndDoDate) > 0 {
			filter := map[string]interface{}{
				"range": map[string]interface{}{
					"so_date": map[string]interface{}{
						"gte": request.StartDoDate,
						"lte": request.EndDoDate,
					},
				},
			}

			filters = append(filters, filter)
		}

		if request.GlobalSearchValue != "" {
			match := map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":  request.GlobalSearchValue,
					"fields": []string{"do_code", "do_ref_code", "sales_order.so_code", "store.store_code", "store.name"},
				},
			}

			musts = append(musts, match)
		}

		if request.SortField != "" && request.SortValue != "" {
			sortValue := map[string]interface{}{
				"order": request.SortValue,
			}

			if request.SortField == "created_at" {
				sortValue["unmapped_type"] = "date"
			}

			if request.SortField == "updated_at" {
				sortValue["unmapped_type"] = "date"
			}

			if request.SortField == "do_date" || request.SortField == "order_status_id" || request.SortField == "created_at" || request.SortField == "updated_at" {
				openSearchQuery["sort"] = []map[string]interface{}{
					{
						request.SortField: sortValue,
					},
				}
			}

			if request.SortField == "do_ref_code" {
				openSearchQuery["sort"] = []map[string]interface{}{
					{
						request.SortField + ".keyword": sortValue,
					},
				}
			}
		}
	}

	openSearchDetailBoolQuery["filter"] = filters
	openSearchDetailBoolQuery["must"] = musts
	openSearchDetailQuery["bool"] = openSearchDetailBoolQuery
	openSearchQuery["query"] = openSearchDetailQuery
	openSearchQueryJson, _ := json.Marshal(openSearchQuery)

	return openSearchQueryJson
}

func (r *deliveryOrderDetailOpenSearch) generateDeliveryOrderQueryOpenSearchByQueryParamTermRequest(term_field string, term_value interface{}, request *models.DeliveryOrderDetailOpenSearchRequest) []byte {
	openSearchQuery := map[string]interface{}{}
	// openSearchDetailQuery := map[string]interface{}{}
	openSearchDetailBoolQuery := map[string]interface{}{}

	if request.Page > 0 {
		page := request.PerPage * (request.Page - 1)
		openSearchQuery["size"] = request.PerPage
		openSearchQuery["from"] = page
	}

	filters := []map[string]interface{}{}

	if len(term_field) > 0 && term_value != nil {
		filter := map[string]interface{}{
			"term": map[string]interface{}{
				term_field: term_value,
			},
		}

		filters = append(filters, filter)
	}

	if len(request.StartCreatedAt) > 0 && len(request.EndCreatedAt) > 0 {
		filter := map[string]interface{}{
			"range": map[string]interface{}{
				"created_at": map[string]interface{}{
					"gte": request.StartCreatedAt,
					"lte": request.EndCreatedAt,
				},
			},
		}

		filters = append(filters, filter)
	}

	if len(request.StartDoDate) > 0 && len(request.EndDoDate) > 0 {
		filter := map[string]interface{}{
			"range": map[string]interface{}{
				"do_date": map[string]interface{}{
					"gte": request.StartDoDate,
					"lte": request.EndDoDate,
				},
			},
		}

		filters = append(filters, filter)
	}

	musts := []map[string]interface{}{}
	multiMatch := map[string]interface{}{}

	if request.GlobalSearchValue != "" {
		multiMatch = map[string]interface{}{
			"query":  request.GlobalSearchValue,
			"fields": []string{"do_code", "do_ref_code", "sales_order.so_code", "store.store_code", "store.name"},
		}
	}

	if request.StartDoDate != "" {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"do_date": request.StartDoDate,
			},
		}

		musts = append(musts, match)
	}

	if request.EndDoDate != "" {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"do_date": request.EndDoDate,
			},
		}

		musts = append(musts, match)
	}

	if request.StartCreatedAt != "" {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"created_at": request.StartCreatedAt,
			},
		}

		musts = append(musts, match)
	}

	if request.EndCreatedAt != "" {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"created_at": request.EndCreatedAt,
			},
		}

		musts = append(musts, match)
	}

	if request.AgentID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"agent_id": request.AgentID,
			},
		}

		musts = append(musts, match)
	}

	if request.StoreID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"store_id": request.StoreID,
			},
		}

		musts = append(musts, match)
	}

	if request.BrandID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"sales_order.brand_id": request.BrandID,
			},
		}

		musts = append(musts, match)
	}

	if request.OrderStatusID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"order_status_id": request.OrderStatusID,
			},
		}

		musts = append(musts, match)
	}

	if request.DoCode != "" {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"do_code": request.DoCode,
			},
		}

		musts = append(musts, match)
	}

	if request.DoRefCode != "" {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"do_ref_code": request.DoRefCode,
			},
		}

		musts = append(musts, match)
	}

	if request.DoRefDate != "" {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"do_ref_date": request.DoRefDate,
			},
		}

		musts = append(musts, match)
	}

	if request.ProductID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"delivery_order_details.product_id": request.ProductID,
			},
		}

		musts = append(musts, match)
	}

	if request.ID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"id": request.ID,
			},
		}

		musts = append(musts, match)
	}

	if request.SalesOrderID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"sales_order_id": request.SalesOrderID,
			},
		}

		musts = append(musts, match)
	}

	if request.SalesmanID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"sales_order.salesman_id": request.SalesmanID,
			},
		}

		musts = append(musts, match)
	}

	if request.StoreProvinceID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"store_province_id": request.StoreProvinceID,
			},
		}

		musts = append(musts, match)
	}

	if request.StoreCityID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"store_city_id": request.StoreCityID,
			},
		}

		musts = append(musts, match)
	}

	if request.StoreDistrictID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"store_district_id": request.StoreDistrictID,
			},
		}

		musts = append(musts, match)
	}

	if request.StoreVillageID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"store_village_id": request.StoreVillageID,
			},
		}

		musts = append(musts, match)
	}

	if request.SortField != "" && request.SortValue != "" {
		sortValue := map[string]interface{}{
			"order": request.SortValue,
		}

		if request.SortField == "created_at" {
			sortValue["unmapped_type"] = "date"
		}

		if request.SortField == "updated_at" {
			sortValue["unmapped_type"] = "date"
		}

		if request.SortField == "do_date" || request.SortField == "order_status_id" || request.SortField == "created_at" || request.SortField == "updated_at" {
			openSearchQuery["sort"] = []map[string]interface{}{
				{
					request.SortField: sortValue,
				},
			}
		}

		if request.SortField == "do_ref_code" {
			openSearchQuery["sort"] = []map[string]interface{}{
				{
					request.SortField + ".keyword": sortValue,
				},
			}
		}

	}

	openSearchDetailQueryString := map[string]interface{}{}
	openSearchDetailQueryBool := map[string]interface{}{}
	openSearchDetailBoolQuery["filter"] = filters
	openSearchDetailBoolQuery["must"] = musts
	openSearchDetailQueryBool["bool"] = openSearchDetailBoolQuery
	openSearchDetailQueryString["multi_match"] = multiMatch
	if request.GlobalSearchValue != "" {
		openSearchQuery["query"] = openSearchDetailQueryString
		openSearchQueryJson, _ := json.Marshal(openSearchQuery)
		return openSearchQueryJson
	}
	openSearchQuery["query"] = openSearchDetailQueryBool
	openSearchQueryJson, _ := json.Marshal(openSearchQuery)
	return openSearchQueryJson
}

func (r *deliveryOrderDetailOpenSearch) generateDeliveryOrderQueryOpenSearchResult(openSearchQueryJson []byte) (*models.DeliveryOrderDetails, *model.ErrorLog) {
	openSearchQueryResult, err := r.openSearch.Query(constants.DELIVERY_ORDER_DETAILS_INDEX, openSearchQueryJson)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.DeliveryOrderDetails{}, errorLogData
	}

	if openSearchQueryResult.Hits.Total.Value == 0 {
		err = helper.NewError("delivery_orders_opensearch data not found")
		errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
		return &models.DeliveryOrderDetails{}, errorLogData
	}

	deliveryOrderDetails := []*models.DeliveryOrderDetail{}

	if openSearchQueryResult.Hits.Total.Value > 0 {
		for _, v := range openSearchQueryResult.Hits.Hits {
			obj := v.Source.(map[string]interface{})
			deliveryOrderDetail := models.DeliveryOrderDetail{}
			objJson, _ := json.Marshal(obj)
			json.Unmarshal(objJson, &deliveryOrderDetail)
			deliveryOrderDetails = append(deliveryOrderDetails, &deliveryOrderDetail)
		}
	}

	result := &models.DeliveryOrderDetails{
		Total:                int64(openSearchQueryResult.Hits.Total.Value),
		DeliveryOrderDetails: deliveryOrderDetails,
	}

	return result, &model.ErrorLog{}
}
