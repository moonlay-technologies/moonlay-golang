package repositories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/global/utils/helper"
	"order-service/global/utils/model"
	"order-service/global/utils/opensearch_dbo"
	"strings"
	"time"
)

type SalesOrderOpenSearchRepositoryInterface interface {
	Create(request *models.SalesOrder, result chan *models.SalesOrderChan)
	Get(request *models.SalesOrderRequest, IsCountOnly bool, result chan *models.SalesOrdersChan)
	GetByID(request *models.SalesOrderRequest, result chan *models.SalesOrderChan)
	GetDetailByID(id int, result chan *models.SalesOrderChan)
	GetBySalesmanID(request *models.SalesOrderRequest, result chan *models.SalesOrdersChan)
	GetByStoreID(request *models.SalesOrderRequest, result chan *models.SalesOrdersChan)
	GetByAgentID(request *models.SalesOrderRequest, result chan *models.SalesOrdersChan)
	GetByOrderStatusID(request *models.SalesOrderRequest, result chan *models.SalesOrdersChan)
	GetByOrderSourceID(request *models.SalesOrderRequest, result chan *models.SalesOrdersChan)
	generateSalesOrderQueryOpenSearchResult(openSearchQueryJson []byte, withSalesOrderDetails bool, isCountOnly bool) (*models.SalesOrders, *model.ErrorLog)
	generateSalesOrderQueryOpenSearchTermRequest(term_field string, term_value interface{}, request *models.SalesOrderRequest) []byte
}

type salesOrderOpenSearch struct {
	openSearch opensearch_dbo.OpenSearchClientInterface
}

func InitSalesOrderOpenSearchRepository(openSearch opensearch_dbo.OpenSearchClientInterface) SalesOrderOpenSearchRepositoryInterface {
	return &salesOrderOpenSearch{
		openSearch: openSearch,
	}
}

func (r *salesOrderOpenSearch) Create(request *models.SalesOrder, resultChan chan *models.SalesOrderChan) {
	response := &models.SalesOrderChan{}
	salesOrderJson, _ := json.Marshal(request)
	st, err := r.openSearch.CreateDocument(constants.SALES_ORDERS_INDEX, request.SoCode, salesOrderJson)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}
	fmt.Println("hasilnya")
	fmt.Println(st)
	response.Error = nil
	response.ErrorLog = &model.ErrorLog{}
	resultChan <- response
	return
}

func (r *salesOrderOpenSearch) Get(request *models.SalesOrderRequest, isCountOnly bool, resultChan chan *models.SalesOrdersChan) {
	response := &models.SalesOrdersChan{}
	if isCountOnly {
		request.Page = 0
		request.PerPage = 0
		request.SortField = ""
		request.SortValue = ""
	}
	requestQuery := r.generateSalesOrderQueryOpenSearchTermRequest("", "", request)
	result, err := r.generateSalesOrderQueryOpenSearchResult(requestQuery, true, isCountOnly)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.SalesOrders = result.SalesOrders
	response.Total = result.Total
	resultChan <- response
	return
}

func (r *salesOrderOpenSearch) GetByID(request *models.SalesOrderRequest, resultChan chan *models.SalesOrderChan) {
	response := &models.SalesOrderChan{}
	requestQuery := r.generateSalesOrderQueryOpenSearchTermRequest("id", request.ID, request)
	result, err := r.generateSalesOrderQueryOpenSearchResult(requestQuery, true, false)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.SalesOrder = result.SalesOrders[0]
	resultChan <- response
	return
}

func (r *salesOrderOpenSearch) GetDetailByID(id int, resultChan chan *models.SalesOrderChan) {
	response := &models.SalesOrderChan{}
	requestQuery := r.generateSalesOrderQueryOpenSearchTermRequest("sales_order_details.id", id, nil)
	result, err := r.generateSalesOrderQueryOpenSearchResult(requestQuery, true, false)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.SalesOrder = result.SalesOrders[0]
	resultChan <- response
	return
}

func (r *salesOrderOpenSearch) GetBySalesmanID(request *models.SalesOrderRequest, resultChan chan *models.SalesOrdersChan) {
	response := &models.SalesOrdersChan{}
	requestQuery := r.generateSalesOrderQueryOpenSearchTermRequest("user_id", request.SalesmanID, request)
	result, err := r.generateSalesOrderQueryOpenSearchResult(requestQuery, true, false)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.SalesOrders = result.SalesOrders
	resultChan <- response
	return
}

func (r *salesOrderOpenSearch) GetByStoreID(request *models.SalesOrderRequest, resultChan chan *models.SalesOrdersChan) {
	response := &models.SalesOrdersChan{}
	requestQuery := r.generateSalesOrderQueryOpenSearchTermRequest("store_id", request.StoreID, request)
	result, err := r.generateSalesOrderQueryOpenSearchResult(requestQuery, true, false)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.SalesOrders = result.SalesOrders
	resultChan <- response
	return
}

func (r *salesOrderOpenSearch) GetByAgentID(request *models.SalesOrderRequest, resultChan chan *models.SalesOrdersChan) {
	response := &models.SalesOrdersChan{}
	requestQuery := r.generateSalesOrderQueryOpenSearchTermRequest("agent_id", request.AgentID, request)
	result, err := r.generateSalesOrderQueryOpenSearchResult(requestQuery, true, false)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.SalesOrders = result.SalesOrders
	resultChan <- response
	return
}

func (r *salesOrderOpenSearch) GetByOrderStatusID(request *models.SalesOrderRequest, resultChan chan *models.SalesOrdersChan) {
	response := &models.SalesOrdersChan{}
	requestQuery := r.generateSalesOrderQueryOpenSearchTermRequest("order_status_id", request.OrderStatusID, request)
	result, err := r.generateSalesOrderQueryOpenSearchResult(requestQuery, true, false)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.SalesOrders = result.SalesOrders
	resultChan <- response
	return
}

func (r *salesOrderOpenSearch) GetByOrderSourceID(request *models.SalesOrderRequest, resultChan chan *models.SalesOrdersChan) {
	response := &models.SalesOrdersChan{}
	requestQuery := r.generateSalesOrderQueryOpenSearchTermRequest("order_source_id", request.OrderSourceID, request)
	result, err := r.generateSalesOrderQueryOpenSearchResult(requestQuery, true, false)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.SalesOrders = result.SalesOrders
	resultChan <- response
	return
}

func (r *salesOrderOpenSearch) generateSalesOrderQueryOpenSearchTermRequest(term_field string, term_value interface{}, request *models.SalesOrderRequest) []byte {
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

		if request.SalesmanID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"salesman.id": request.SalesmanID,
				},
			}

			filters = append(filters, filter)
		}

		if request.AgentID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"agent.id": request.AgentID,
				},
			}

			filters = append(filters, filter)
		}

		if request.StoreID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store_id": request.StoreID,
				},
			}

			filters = append(filters, filter)
		}

		if request.ID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"id": request.ID,
				},
			}

			filters = append(filters, filter)
		}

		if request.ProductID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"sales_order_details.product_id": request.ProductID,
				},
			}

			filters = append(filters, filter)
		}

		if request.CategoryID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"sales_order_details.product.category_id": request.CategoryID,
				},
			}

			filters = append(filters, filter)
		}

		if request.BrandID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"brand_id": request.BrandID,
				},
			}

			filters = append(filters, filter)
		}

		if request.OrderSourceID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"order_source_id": request.OrderSourceID,
				},
			}

			filters = append(filters, filter)
		}

		if request.OrderStatusID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"order_status_id": request.OrderStatusID,
				},
			}

			filters = append(filters, filter)
		}

		if request.ProvinceID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store.province_id": request.ProvinceID,
				},
			}

			filters = append(filters, filter)
		}

		if request.CityID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store.city_id": request.CityID,
				},
			}

			filters = append(filters, filter)
		}

		if request.DistrictID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store.district_id": request.DistrictID,
				},
			}

			filters = append(filters, filter)
		}

		if request.VillageID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store.village_id": request.VillageID,
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

		if len(request.StartSoDate) > 0 && len(request.EndSoDate) > 0 {
			filter := map[string]interface{}{
				"range": map[string]interface{}{
					"so_date": map[string]interface{}{
						"gte": request.StartSoDate,
						"lte": request.EndSoDate,
					},
				},
			}

			filters = append(filters, filter)
		}

		if request.GlobalSearchValue != "" {
			globalSearchValue := strings.ReplaceAll(request.GlobalSearchValue, "-", " ")
			match := map[string]interface{}{
				"query_string": map[string]interface{}{
					"query":            "*" + globalSearchValue + "*",
					"fields":           []string{"store_code", "store_name", "so_code", "so_ref_code"},
					"default_operator": "AND",
				},
			}

			musts = append(musts, match)
		}

		if len(request.SortField) > 0 && len(request.SortValue) > 0 {
			sortValue := map[string]interface{}{
				"order": request.SortValue,
			}

			if helper.Contains(constants.UNMAPPED_TYPE_SORT_LIST(), request.SortField) {
				sortValue["unmapped_type"] = "date"
			}

			if helper.Contains(constants.SALES_ORDER_SORT_INT_LIST(), request.SortField) {
				openSearchQuery["sort"] = []map[string]interface{}{
					{
						request.SortField: sortValue,
					},
				}
			}

			if helper.Contains(constants.SALES_ORDER_SORT_STRING_LIST(), request.SortField) {
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

func (r *salesOrderOpenSearch) generateSalesOrderQueryOpenSearchResult(openSearchQueryJson []byte, withSalesOrderDetails bool, isCountOnly bool) (*models.SalesOrders, *model.ErrorLog) {
	salesOrders := []*models.SalesOrder{}
	var total int64 = 0

	if isCountOnly {
		openSearchQueryResult, err := r.openSearch.Count(constants.SALES_ORDERS_INDEX, openSearchQueryJson)
		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return &models.SalesOrders{}, errorLogData
		}

		if openSearchQueryResult <= 0 {
			err = helper.NewError("sales_orders_opensearch data not found")
			errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
			return &models.SalesOrders{}, errorLogData
		}

		total = openSearchQueryResult
	} else {
		openSearchQueryResult, err := r.openSearch.Query(constants.SALES_ORDERS_INDEX, openSearchQueryJson)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return &models.SalesOrders{}, errorLogData
		}

		if openSearchQueryResult.Hits.Total.Value == 0 {
			err = helper.NewError(helper.DefaultStatusText[http.StatusNotFound])
			errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
			return &models.SalesOrders{}, errorLogData
		}

		total = int64(openSearchQueryResult.Hits.Total.Value)

		if openSearchQueryResult.Hits.Total.Value > 0 {
			for _, v := range openSearchQueryResult.Hits.Hits {
				obj := v.Source.(map[string]interface{})
				salesOrder := &models.SalesOrder{}
				jsonStr, err := json.Marshal(v.Source)
				if err != nil {
					fmt.Println(err)
				}
				if err := json.Unmarshal(jsonStr, &salesOrder); err != nil {
					fmt.Println(err)
				}

				layout := time.RFC3339
				createdAt, _ := time.Parse(layout, obj["created_at"].(string))
				salesOrder.CreatedAt = &createdAt
				updatedAt, _ := time.Parse(layout, obj["updated_at"].(string))
				salesOrder.UpdatedAt = &updatedAt

				salesOrders = append(salesOrders, salesOrder)
			}
		}
	}

	result := &models.SalesOrders{
		Total:       total,
		SalesOrders: salesOrders,
	}

	return result, &model.ErrorLog{}
}
