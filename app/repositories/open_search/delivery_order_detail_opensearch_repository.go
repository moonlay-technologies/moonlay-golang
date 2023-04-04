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
	Create(request *models.DeliveryOrderDetailOpenSearch, result chan *models.DeliveryOrderDetailOpenSearchChan)
	Get(request *models.DeliveryOrderDetailOpenSearchRequest, isCountOnly bool, result chan *models.DeliveryOrderDetailsOpenSearchChan)
	GetByID(request *models.DeliveryOrderRequest, resultChan chan *models.DeliveryOrderDetailOpenSearchChan)
	generateDeliveryOrderQueryOpenSearchResult(openSearchQueryJson []byte, isCountOnly bool) (*models.DeliveryOrderDetailsOpenSearch, *model.ErrorLog)
	generateDeliveryOrderQueryOpenSearchTermRequest(term_field string, term_value interface{}, request *models.DeliveryOrderDetailOpenSearchRequest) []byte
}

type deliveryOrderDetailOpenSearch struct {
	openSearch opensearch_dbo.OpenSearchClientInterface
}

func InitDeliveryOrderDetailOpenSearchRepository(openSearch opensearch_dbo.OpenSearchClientInterface) DeliveryOrderDetailOpenSearchRepositoryInterface {
	return &deliveryOrderDetailOpenSearch{
		openSearch: openSearch,
	}
}

func (r *deliveryOrderDetailOpenSearch) Create(request *models.DeliveryOrderDetailOpenSearch, resultChan chan *models.DeliveryOrderDetailOpenSearchChan) {
	response := &models.DeliveryOrderDetailOpenSearchChan{}
	deliveryOrderDetailJson, _ := json.Marshal(request)
	_, err := r.openSearch.CreateDocument(constants.DELIVERY_ORDER_DETAILS_INDEX, request.DoDetailCode, deliveryOrderDetailJson)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	response.DeliveryOrderDetailOpenSearch = request
	resultChan <- response
	return
}

func (r *deliveryOrderDetailOpenSearch) Get(request *models.DeliveryOrderDetailOpenSearchRequest, isCountOnly bool, resultChan chan *models.DeliveryOrderDetailsOpenSearchChan) {
	response := &models.DeliveryOrderDetailsOpenSearchChan{}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchTermRequest("", "", request)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery, isCountOnly)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrderDetailOpenSearch = result.DeliveryOrderDetails
	response.Total = result.Total
	resultChan <- response
	return
}

func (r *deliveryOrderDetailOpenSearch) GetByID(request *models.DeliveryOrderRequest, resultChan chan *models.DeliveryOrderDetailOpenSearchChan) {
	response := &models.DeliveryOrderDetailOpenSearchChan{}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchTermRequest("id", request.ID, nil)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery, false)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrderDetailOpenSearch = result.DeliveryOrderDetails[0]
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

		if len(request.StartDoDate) > 0 && len(request.EndDoDate) == 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"do_date": request.StartDoDate,
				},
			}

			filters = append(filters, filter)
		}

		if len(request.EndDoDate) > 0 && len(request.StartDoDate) == 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"do_date": request.EndDoDate,
				},
			}

			filters = append(filters, filter)
		}

		if len(request.StartCreatedAt) > 0 && len(request.EndCreatedAt) == 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"created_at": request.StartCreatedAt,
				},
			}

			filters = append(filters, filter)
		}

		if len(request.EndCreatedAt) > 0 && len(request.StartCreatedAt) == 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"created_at": request.EndCreatedAt,
				},
			}

			filters = append(filters, filter)
		}

		if request.ID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"id": request.ID,
				},
			}

			filters = append(filters, filter)
		}

		if request.DeliveryOrderID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"delivery_order_id": request.DeliveryOrderID,
				},
			}

			filters = append(filters, filter)
		}

		if request.DoCode != "" {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"do_code.keyword": request.DoCode,
				},
			}

			filters = append(filters, filter)
		}

		if request.SalesOrderID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"so_detail.sales_order_id": request.SalesOrderID,
				},
			}

			filters = append(filters, filter)
		}

		if request.SoCode != "" {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"so_code.keyword": request.SoCode,
				},
			}

			filters = append(filters, filter)
		}

		if request.AgentID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"agent_id": request.AgentID,
				},
			}

			filters = append(filters, filter)
		}

		if request.StoreID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store_id": request.StoreID,
				},
			}

			filters = append(filters, filter)
		}

		if request.BrandID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"brand_id": request.BrandID,
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

		if request.ProductID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"product_id": request.ProductID,
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

		if request.GlobalSearchValue != "" {
			match := map[string]interface{}{
				"query_string": map[string]interface{}{
					"query":            "*" + request.GlobalSearchValue + "*",
					"fields":           []string{"do_code", "so_code", "order_status.name^0.5", "qty^3"},
					"type":             "best_fields",
					"default_operator": "AND",
					"lenient":          true,
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

			if request.SortField == "order_status_id" || request.SortField == "agent_id" || request.SortField == "store_id" || request.SortField == "product_id" || request.SortField == "qty" {
				openSearchQuery["sort"] = []map[string]interface{}{
					{
						request.SortField: sortValue,
					},
				}
			}

			if request.SortField == "do_code" || request.SortField == "so_code" {
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

func (r *deliveryOrderDetailOpenSearch) generateDeliveryOrderQueryOpenSearchResult(openSearchQueryJson []byte, isCountOnly bool) (*models.DeliveryOrderDetailsOpenSearch, *model.ErrorLog) {
	deliveryOrderDetails := []*models.DeliveryOrderDetailOpenSearch{}
	var total int64 = 0

	if isCountOnly {
		openSearchQueryResult, err := r.openSearch.Count(constants.DELIVERY_ORDERS_INDEX, openSearchQueryJson)
		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return &models.DeliveryOrderDetailsOpenSearch{}, errorLogData
		}

		if openSearchQueryResult <= 0 {
			err = helper.NewError("delivery_orders_opensearch data not found")
			errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
			return &models.DeliveryOrderDetailsOpenSearch{}, errorLogData
		}

		total = openSearchQueryResult
	} else {
		openSearchQueryResult, err := r.openSearch.Query(constants.DELIVERY_ORDER_DETAILS_INDEX, openSearchQueryJson)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return &models.DeliveryOrderDetailsOpenSearch{}, errorLogData
		}

		if openSearchQueryResult.Hits.Total.Value == 0 {
			err = helper.NewError("delivery_orders_opensearch data not found")
			errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
			return &models.DeliveryOrderDetailsOpenSearch{}, errorLogData
		}

		total = int64(openSearchQueryResult.Hits.Total.Value)

		if openSearchQueryResult.Hits.Total.Value > 0 {
			for _, v := range openSearchQueryResult.Hits.Hits {
				obj := v.Source.(map[string]interface{})
				deliveryOrderDetail := models.DeliveryOrderDetailOpenSearch{}
				objJson, _ := json.Marshal(obj)
				json.Unmarshal(objJson, &deliveryOrderDetail)
				deliveryOrderDetails = append(deliveryOrderDetails, &deliveryOrderDetail)
			}
		}
	}

	result := &models.DeliveryOrderDetailsOpenSearch{
		Total:                total,
		DeliveryOrderDetails: deliveryOrderDetails,
	}

	return result, &model.ErrorLog{}
}
