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

type DeliveryOrderOpenSearchRepositoryInterface interface {
	Create(request *models.DeliveryOrder, result chan *models.DeliveryOrderChan)
	Get(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan)
	GetByID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrderChan)
	GetBySalesOrderID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan)
	GetBySalesmanID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan)
	GetBySalesmansID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan)
	GetByStoreID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan)
	GetByAgentID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan)
	generateDeliveryOrderQueryOpenSearchResult(openSearchQueryJson []byte, withDeliveryOrderDetails bool) (*models.DeliveryOrders, *model.ErrorLog)
	generateDeliveryOrderQueryOpenSearchTermRequest(term_field string, term_value interface{}, request *models.DeliveryOrderRequest) []byte
	generateDeliveryOrderQueryOpenSearchByQueryParamTermRequest(term_field string, term_value interface{}, request *models.DeliveryOrderRequest) []byte
}

type deliveryOrderOpenSearch struct {
	openSearch opensearch_dbo.OpenSearchClientInterface
}

func InitDeliveryOrderOpenSearchRepository(openSearch opensearch_dbo.OpenSearchClientInterface) DeliveryOrderOpenSearchRepositoryInterface {
	return &deliveryOrderOpenSearch{
		openSearch: openSearch,
	}
}

func (r *deliveryOrderOpenSearch) Create(request *models.DeliveryOrder, resultChan chan *models.DeliveryOrderChan) {
	response := &models.DeliveryOrderChan{}
	deliveryOrderJson, _ := json.Marshal(request)
	_, err := r.openSearch.CreateDocument(constants.DELIVERY_ORDERS_INDEX, request.DoCode, deliveryOrderJson)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	response.DeliveryOrder = request
	response.Total = 1
	resultChan <- response
	return
}

func (r *deliveryOrderOpenSearch) Get(request *models.DeliveryOrderRequest, resultChan chan *models.DeliveryOrdersChan) {
	response := &models.DeliveryOrdersChan{}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchTermRequest("", "", request)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery, true)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrders = result.DeliveryOrders
	response.Total = result.Total
	resultChan <- response
	return
}

func (r *deliveryOrderOpenSearch) GetByID(request *models.DeliveryOrderRequest, resultChan chan *models.DeliveryOrderChan) {
	response := &models.DeliveryOrderChan{}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchTermRequest("id", request.ID, request)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery, true)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrder = result.DeliveryOrders[0]
	resultChan <- response
	return
}

func (r *deliveryOrderOpenSearch) GetBySalesOrderID(request *models.DeliveryOrderRequest, resultChan chan *models.DeliveryOrdersChan) {
	response := &models.DeliveryOrdersChan{}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchTermRequest("sales_order_id", request.SalesOrderID, request)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery, true)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrders = result.DeliveryOrders
	resultChan <- response
	return
}

func (r *deliveryOrderOpenSearch) GetBySalesmanID(request *models.DeliveryOrderRequest, resultChan chan *models.DeliveryOrdersChan) {
	response := &models.DeliveryOrdersChan{}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchTermRequest("sales_order.user_id", request.SalesmanID, request)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery, true)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrders = result.DeliveryOrders
	resultChan <- response
	return
}

func (r *deliveryOrderOpenSearch) GetBySalesmansID(request *models.DeliveryOrderRequest, resultChan chan *models.DeliveryOrdersChan) {
	response := &models.DeliveryOrdersChan{}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchByQueryParamTermRequest("", "", request)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery, true)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrders = result.DeliveryOrders
	resultChan <- response
	return
}

func (r *deliveryOrderOpenSearch) GetByStoreID(request *models.DeliveryOrderRequest, resultChan chan *models.DeliveryOrdersChan) {
	response := &models.DeliveryOrdersChan{}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchTermRequest("store_id", request.StoreID, request)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery, true)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrders = result.DeliveryOrders
	resultChan <- response
	return
}

func (r *deliveryOrderOpenSearch) GetByAgentID(request *models.DeliveryOrderRequest, resultChan chan *models.DeliveryOrdersChan) {
	response := &models.DeliveryOrdersChan{}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchTermRequest("agent_id", request.AgentID, request)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery, true)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrders = result.DeliveryOrders
	resultChan <- response
	return
}

func (r *deliveryOrderOpenSearch) generateDeliveryOrderQueryOpenSearchTermRequest(term_field string, term_value interface{}, request *models.DeliveryOrderRequest) []byte {
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

	if request.CategoryID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"category_id": request.CategoryID,
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

	if request.ProvinceID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"store.province_id": request.ProvinceID,
			},
		}

		musts = append(musts, match)
	}

	if request.CityID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"store.city_id": request.CityID,
			},
		}

		musts = append(musts, match)
	}

	if request.DistrictID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"store.district_id": request.DistrictID,
			},
		}

		musts = append(musts, match)
	}

	if request.VillageID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"store.village_id": request.VillageID,
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

func (r *deliveryOrderOpenSearch) generateDeliveryOrderQueryOpenSearchByQueryParamTermRequest(term_field string, term_value interface{}, request *models.DeliveryOrderRequest) []byte {
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

	if request.CategoryID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"category_id": request.CategoryID,
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

	if request.ProvinceID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"store.province_id": request.ProvinceID,
			},
		}

		musts = append(musts, match)
	}

	if request.CityID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"store.city_id": request.CityID,
			},
		}

		musts = append(musts, match)
	}

	if request.DistrictID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"store.district_id": request.DistrictID,
			},
		}

		musts = append(musts, match)
	}

	if request.VillageID != 0 {
		match := map[string]interface{}{
			"match": map[string]interface{}{
				"store.village_id": request.VillageID,
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

func (r *deliveryOrderOpenSearch) generateDeliveryOrderQueryOpenSearchResult(openSearchQueryJson []byte, withDeliveryOrderDetails bool) (*models.DeliveryOrders, *model.ErrorLog) {
	openSearchQueryResult, err := r.openSearch.Query(constants.DELIVERY_ORDERS_INDEX, openSearchQueryJson)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return &models.DeliveryOrders{}, errorLogData
	}

	if openSearchQueryResult.Hits.Total.Value == 0 {
		err = helper.NewError("delivery_orders_opensearch data not found")
		errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
		return &models.DeliveryOrders{}, errorLogData
	}

	deliveryOrders := []*models.DeliveryOrder{}

	if openSearchQueryResult.Hits.Total.Value > 0 {
		for _, v := range openSearchQueryResult.Hits.Hits {
			obj := v.Source.(map[string]interface{})
			deliveryOrder := models.DeliveryOrder{}
			objJson, _ := json.Marshal(obj)
			json.Unmarshal(objJson, &deliveryOrder)
			deliveryOrders = append(deliveryOrders, &deliveryOrder)
		}
	}

	result := &models.DeliveryOrders{
		Total:          int64(openSearchQueryResult.Hits.Total.Value),
		DeliveryOrders: deliveryOrders,
	}

	return result, &model.ErrorLog{}
}
