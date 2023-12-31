package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/app/repositories"
	"order-service/global/utils/helper"
	"order-service/global/utils/model"
	"order-service/global/utils/opensearch_dbo"
	"strings"
	"time"
)

type DeliveryOrderOpenSearchRepositoryInterface interface {
	Create(request *models.DeliveryOrder, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.DeliveryOrderChan)
	Get(request *models.DeliveryOrderRequest, isCountOnly bool, result chan *models.DeliveryOrdersChan)
	GetByID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrderChan)
	GetDetailsByDoID(request *models.DeliveryOrderDetailRequest, result chan *models.DeliveryOrderChan)
	GetDetailByID(doDetailID int, result chan *models.DeliveryOrderChan)
	GetBySalesOrderID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan)
	GetBySalesmanID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan)
	GetByStoreID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan)
	GetByAgentID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan)
	generateDeliveryOrderQueryOpenSearchResult(openSearchQueryJson []byte, withDeliveryOrderDetails bool, isCountOnly bool) (*models.DeliveryOrders, *model.ErrorLog)
	generateDeliveryOrderQueryOpenSearchTermRequest(term_field string, term_value interface{}, request *models.DeliveryOrderRequest) []byte
	generateDeliveryOrderDetailQueryOpenSearchTermRequest(term_field string, term_value interface{}, request *models.DeliveryOrderDetailRequest) []byte
}

type deliveryOrderOpenSearch struct {
	openSearch              opensearch_dbo.OpenSearchClientInterface
	deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface
}

func InitDeliveryOrderOpenSearchRepository(openSearch opensearch_dbo.OpenSearchClientInterface, deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface) DeliveryOrderOpenSearchRepositoryInterface {
	return &deliveryOrderOpenSearch{
		openSearch:              openSearch,
		deliveryOrderRepository: deliveryOrderRepository,
	}
}

func (r *deliveryOrderOpenSearch) Create(request *models.DeliveryOrder, sqlTransaction *sql.Tx, ctx context.Context, resultChan chan *models.DeliveryOrderChan) {
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

	updateDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
	go r.deliveryOrderRepository.UpdateByID(request.ID, request, "", false, sqlTransaction, ctx, updateDeliveryOrderResultChan)
	updateDeliveryOrderResult := <-updateDeliveryOrderResultChan

	if updateDeliveryOrderResult.Error != nil {
		errorLogData := helper.WriteLog(updateDeliveryOrderResult.Error, http.StatusInternalServerError, nil)
		response.Error = updateDeliveryOrderResult.Error
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	response.DeliveryOrder = request
	response.Total = 1
	resultChan <- response
	return
}

func (r *deliveryOrderOpenSearch) Get(request *models.DeliveryOrderRequest, isCountOnly bool, resultChan chan *models.DeliveryOrdersChan) {
	response := &models.DeliveryOrdersChan{}
	if isCountOnly {
		request.Page = 0
		request.PerPage = 0
		request.SortField = ""
		request.SortValue = ""
	}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchTermRequest("", "", request)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery, true, isCountOnly)

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
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery, true, false)

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

func (r *deliveryOrderOpenSearch) GetDetailsByDoID(request *models.DeliveryOrderDetailRequest, resultChan chan *models.DeliveryOrderChan) {
	response := &models.DeliveryOrderChan{}
	requestQuery := r.generateDeliveryOrderDetailQueryOpenSearchTermRequest("", "", request)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery, true, false)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrder = result.DeliveryOrders[0]
	response.Total = result.Total
	resultChan <- response
	return
}

func (r *deliveryOrderOpenSearch) GetDetailByID(doDetailID int, resultChan chan *models.DeliveryOrderChan) {
	response := &models.DeliveryOrderChan{}
	requestQuery := r.generateDeliveryOrderDetailQueryOpenSearchTermRequest("delivery_order_details.id", doDetailID, nil)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery, true, false)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.DeliveryOrder = result.DeliveryOrders[0]
	response.Total = result.Total
	resultChan <- response
	return
}

func (r *deliveryOrderOpenSearch) GetBySalesOrderID(request *models.DeliveryOrderRequest, resultChan chan *models.DeliveryOrdersChan) {
	response := &models.DeliveryOrdersChan{}
	requestQuery := r.generateDeliveryOrderQueryOpenSearchTermRequest("sales_order_id", request.SalesOrderID, request)
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery, true, false)

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
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery, true, false)

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
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery, true, false)

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
	result, err := r.generateDeliveryOrderQueryOpenSearchResult(requestQuery, true, false)

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
					"sales_order.brand_id": request.BrandID,
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

		if request.DoCode != "" {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"do_code.keyword": request.DoCode,
				},
			}

			filters = append(filters, filter)
		}

		if request.SoCode != "" {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"sales_order.so_code.keyword": request.SoCode,
				},
			}

			filters = append(filters, filter)
		}

		if request.DoRefCode != "" {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"do_ref_code.keyword": request.DoRefCode,
				},
			}

			filters = append(filters, filter)
		}

		if request.DoRefDate != "" {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"do_ref_date": request.DoRefDate,
				},
			}

			filters = append(filters, filter)
		}

		if request.ProductID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"delivery_order_details.product_id": request.ProductID,
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

		if request.SalesOrderID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"sales_order_id": request.SalesOrderID,
				},
			}

			filters = append(filters, filter)
		}

		if request.CategoryID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"delivery_order_details.product.category_id": request.CategoryID,
				},
			}

			filters = append(filters, filter)
		}

		if request.SalesmanID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"sales_order.salesman_id": request.SalesmanID,
				},
			}

			filters = append(filters, filter)
		}

		if request.ProvinceID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store.province_id": request.ProvinceID,
				},
			}

			filters = append(filters, filter)
		}

		if request.CityID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store.city_id": request.CityID,
				},
			}

			filters = append(filters, filter)
		}

		if request.DistrictID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store.district_id": request.DistrictID,
				},
			}

			filters = append(filters, filter)
		}

		if request.VillageID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store.village_id": request.VillageID,
				},
			}

			filters = append(filters, filter)
		}

		if request.UpdatedAt != "" {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"updated_at": request.UpdatedAt,
				},
			}

			filters = append(filters, filter)
		}

		if len(request.StartCreatedAt) > 0 && len(request.EndCreatedAt) > 0 {
			var startCreatedAt string
			var endCreatedAt string
			startTimeadj, err := time.Parse(constants.DATE_FORMAT_EXPORT_CREATED_AT, request.StartCreatedAt+constants.DATE_TIME_ZERO_HOUR_ADDITIONAL)
			if err == nil {
				startCreatedAt = startTimeadj.Format(time.RFC3339)
			} else {
				fmt.Println("error = ", err.Error())
			}
			endTimeadj, err := time.Parse(constants.DATE_FORMAT_EXPORT_CREATED_AT, request.EndCreatedAt+constants.DATE_TIME_ZERO_HOUR_ADDITIONAL)
			if err == nil {
				endTimeadj = endTimeadj.Add(time.Hour * +23).Add(time.Minute * 59).Add(time.Second * 59)
				endCreatedAt = endTimeadj.Format(time.RFC3339)
			} else {
				fmt.Println("error = ", err.Error())
			}
			filter := map[string]interface{}{
				"range": map[string]interface{}{
					"created_at": map[string]interface{}{
						"gte": startCreatedAt,
						"lte": endCreatedAt,
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
			if strings.Contains(request.GlobalSearchValue, "-") {
				globalSearchValue := strings.ReplaceAll(request.GlobalSearchValue, "-", " ")
				match := map[string]interface{}{
					"query_string": map[string]interface{}{
						"query":            "*" + globalSearchValue + "*",
						"fields":           []string{"do_code", "do_ref_code", "sales_order.so_code", "store.store_code", "store.name"},
						"default_operator": "AND",
					},
				}
				musts = append(musts, match)
			} else {
				match := map[string]interface{}{
					"query_string": map[string]interface{}{
						"query":            "*" + request.GlobalSearchValue + "*",
						"fields":           []string{"do_code", "do_ref_code", "sales_order.so_code", "store.store_code", "store.name"},
						"default_operator": "AND",
					},
				}

				musts = append(musts, match)
			}
		}

		if request.SortField != "" && request.SortValue != "" {
			sortValue := map[string]interface{}{
				"order": request.SortValue,
			}

			if helper.Contains(constants.UNMAPPED_TYPE_SORT_LIST(), request.SortField) {
				sortValue["unmapped_type"] = "date"
				openSearchQuery["sort"] = []map[string]interface{}{
					{
						request.SortField: sortValue,
					},
				}
			}

			if helper.Contains(constants.DELIVERY_ORDER_SORT_INT_LIST(), request.SortField) {
				openSearchQuery["sort"] = []map[string]interface{}{
					{
						request.SortField: sortValue,
					},
				}
			}

			if helper.Contains(constants.DELIVERY_ORDER_SORT_STRING_LIST(), request.SortField) {
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

func (r *deliveryOrderOpenSearch) generateDeliveryOrderDetailQueryOpenSearchTermRequest(term_field string, term_value interface{}, request *models.DeliveryOrderDetailRequest) []byte {
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

		if request.DoDetailID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"delivery_order_details.id": request.DoDetailID,
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
					"sales_order.brand_id": request.BrandID,
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

		if request.DoCode != "" {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"do_code.keyword": request.DoCode,
				},
			}

			filters = append(filters, filter)
		}

		if request.SoCode != "" {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"sales_order.so_code.keyword": request.SoCode,
				},
			}

			filters = append(filters, filter)
		}

		if request.DoRefCode != "" {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"do_ref_code.keyword": request.DoRefCode,
				},
			}

			filters = append(filters, filter)
		}

		if request.DoRefDate != "" {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"do_ref_date": request.DoRefDate,
				},
			}

			filters = append(filters, filter)
		}

		if request.ProductID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"delivery_order_details.product_id": request.ProductID,
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

		if request.SalesOrderID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"sales_order_id": request.SalesOrderID,
				},
			}

			filters = append(filters, filter)
		}

		if request.CategoryID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"delivery_order_details.product.category_id": request.CategoryID,
				},
			}

			filters = append(filters, filter)
		}

		if request.SalesmanID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"sales_order.salesman_id": request.SalesmanID,
				},
			}

			filters = append(filters, filter)
		}

		if request.ProvinceID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store.province_id": request.ProvinceID,
				},
			}

			filters = append(filters, filter)
		}

		if request.CityID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store.city_id": request.CityID,
				},
			}

			filters = append(filters, filter)
		}

		if request.DistrictID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store.district_id": request.DistrictID,
				},
			}

			filters = append(filters, filter)
		}

		if request.VillageID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store.village_id": request.VillageID,
				},
			}

			filters = append(filters, filter)
		}

		if request.Qty != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"delivery_order_details.qty": request.Qty,
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
					"fields":           []string{"do_code", "do_ref_code", "sales_order.so_code", "store.store_code", "store.name"},
					"default_operator": "AND",
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

func (r *deliveryOrderOpenSearch) generateDeliveryOrderQueryOpenSearchResult(openSearchQueryJson []byte, withDeliveryOrderDetails bool, isCountOnly bool) (*models.DeliveryOrders, *model.ErrorLog) {
	deliveryOrders := []*models.DeliveryOrder{}
	var total int64 = 0

	if isCountOnly {
		openSearchQueryResult, err := r.openSearch.Count(constants.DELIVERY_ORDERS_INDEX, openSearchQueryJson)
		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return &models.DeliveryOrders{}, errorLogData
		}

		if openSearchQueryResult <= 0 {
			err = helper.NewError(constants.ERROR_DATA_NOT_FOUND)
			errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
			return &models.DeliveryOrders{}, errorLogData
		}

		total = openSearchQueryResult
	} else {
		openSearchQueryResult, err := r.openSearch.Query(constants.DELIVERY_ORDERS_INDEX, openSearchQueryJson)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return &models.DeliveryOrders{}, errorLogData
		}

		total = int64(openSearchQueryResult.Hits.Total.Value)

		if int64(openSearchQueryResult.Hits.Total.Value) <= 0 {
			err = helper.NewError(constants.ERROR_DATA_NOT_FOUND)
			errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
			return &models.DeliveryOrders{}, errorLogData
		}

		loc, _ := time.LoadLocation("Asia/Jakarta")
		for _, v := range openSearchQueryResult.Hits.Hits {
			obj := v.Source.(map[string]interface{})
			deliveryOrder := models.DeliveryOrder{}
			objJson, _ := json.Marshal(obj)
			json.Unmarshal(objJson, &deliveryOrder)
			layout := time.RFC3339
			if obj["created_at"] != nil {
				createdAt, _ := time.ParseInLocation(layout, obj["created_at"].(string), loc)
				createdAt = createdAt.In(loc)
				deliveryOrder.CreatedAt = &createdAt
			}
			if obj["updated_at"] != nil {
				updatedAt, _ := time.ParseInLocation(layout, obj["created_at"].(string), loc)
				updatedAt = updatedAt.In(loc)
				deliveryOrder.UpdatedAt = &updatedAt
			}
			deliveryOrders = append(deliveryOrders, &deliveryOrder)
		}
	}

	result := &models.DeliveryOrders{
		Total:          total,
		DeliveryOrders: deliveryOrders,
	}
	return result, &model.ErrorLog{}
}
