package repositories

import (
	"encoding/json"
	"net/http"
	"poc-order-service/app/models"
	"poc-order-service/global/utils/helper"
	"poc-order-service/global/utils/model"
	"poc-order-service/global/utils/opensearch_dbo"
)

type DeliveryOrderOpenSearchRepositoryInterface interface {
	Create(request *models.DeliveryOrder, result chan *models.DeliveryOrderChan)
	Get(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan)
	GetByID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrderChan)
	GetBySalesOrderID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan)
	GetBySalesmanID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan)
	GetByStoreID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan)
	GetByAgentID(request *models.DeliveryOrderRequest, result chan *models.DeliveryOrdersChan)
	generateDeliveryOrderQueryOpenSearchResult(openSearchQueryJson []byte, withDeliveryOrderDetails bool) (*models.DeliveryOrders, *model.ErrorLog)
	generateDeliveryOrderQueryOpenSearchTermRequest(term_field string, term_value interface{}, request *models.DeliveryOrderRequest) []byte
}

type deliveryOrderOpenSearch struct {
	db opensearch_dbo.OpenSearchClientInterface
}

func InitDeliveryOrderOpenSearchRepository(db opensearch_dbo.OpenSearchClientInterface) DeliveryOrderOpenSearchRepositoryInterface {
	return &deliveryOrderOpenSearch{
		db: db,
	}
}

func (r *deliveryOrderOpenSearch) Create(request *models.DeliveryOrder, resultChan chan *models.DeliveryOrderChan) {
	response := &models.DeliveryOrderChan{}
	deliveryOrderJson, _ := json.Marshal(request)
	_, err := r.db.CreateDocument("delivery_orders", request.DoCode, deliveryOrderJson)

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
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
	openSearchDetailQuery := map[string]interface{}{}
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

	if len(request.StartSoDate) > 0 && len(request.EndSoDate) > 0 {
		filter := map[string]interface{}{
			"range": map[string]interface{}{
				"do_date": map[string]interface{}{
					"gte": request.StartSoDate,
					"lte": request.EndSoDate,
				},
			},
		}

		filters = append(filters, filter)
	}

	if len(request.SortField) > 0 && len(request.SortValue) > 0 {
		sortValue := map[string]interface{}{
			"order": request.SortValue,
		}

		if request.SortField == "created_at" {
			sortValue["unmapped_type"] = "date"
		}

		openSearchQuery["sort"] = []map[string]interface{}{
			{
				request.SortField: sortValue,
			},
		}
	}

	openSearchDetailBoolQuery["filter"] = filters
	openSearchDetailQuery["bool"] = openSearchDetailBoolQuery
	openSearchQuery["query"] = openSearchDetailQuery
	openSearchQueryJson, _ := json.Marshal(openSearchQuery)
	return openSearchQueryJson
}

func (r *deliveryOrderOpenSearch) generateDeliveryOrderQueryOpenSearchResult(openSearchQueryJson []byte, withDeliveryOrderDetails bool) (*models.DeliveryOrders, *model.ErrorLog) {
	openSearchQueryResult, err := r.db.Query("delivery_orders", openSearchQueryJson)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, "Ada kesalahan, silahkan coba lagi nanti")
		return &models.DeliveryOrders{}, errorLogData
	}

	if openSearchQueryResult.Hits.Total.Value == 0 {
		err = helper.NewError("delivery_orders_opensearch data not found")
		errorLogData := helper.WriteLog(err, http.StatusNotFound, "Data tidak ditemukan")
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
