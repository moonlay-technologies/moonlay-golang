package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/global/utils/helper"
	"order-service/global/utils/model"
	"order-service/global/utils/opensearch_dbo"
	"time"
)

type SalesOrderOpenSearchRepositoryInterface interface {
	Create(request *models.SalesOrder, result chan *models.SalesOrderChan)
	Get(request *models.SalesOrderRequest, result chan *models.SalesOrdersChan)
	GetByID(request *models.SalesOrderRequest, result chan *models.SalesOrderChan)
	GetBySalesmanID(request *models.SalesOrderRequest, result chan *models.SalesOrdersChan)
	GetByStoreID(request *models.SalesOrderRequest, result chan *models.SalesOrdersChan)
	GetByAgentID(request *models.SalesOrderRequest, result chan *models.SalesOrdersChan)
	GetByOrderStatusID(request *models.SalesOrderRequest, result chan *models.SalesOrdersChan)
	GetByOrderSourceID(request *models.SalesOrderRequest, result chan *models.SalesOrdersChan)
	generateSalesOrderQueryOpenSearchResult(openSearchQueryJson []byte, withSalesOrderDetails bool) (*models.SalesOrders, *model.ErrorLog)
	generateSalesOrderQueryOpenSearchTermRequest(term_field string, term_value interface{}, request *models.SalesOrderRequest) []byte
}

type salesOrderOpenSearch struct {
	db opensearch_dbo.OpenSearchClientInterface
}

func InitSalesOrderOpenSearchRepository(db opensearch_dbo.OpenSearchClientInterface) SalesOrderOpenSearchRepositoryInterface {
	return &salesOrderOpenSearch{
		db: db,
	}
}

func (r *salesOrderOpenSearch) Create(request *models.SalesOrder, resultChan chan *models.SalesOrderChan) {
	response := &models.SalesOrderChan{}
	salesOrderJson, _ := json.Marshal(request)
	st, err := r.db.CreateDocument("sales_orders", request.SoCode, salesOrderJson)

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, nil)
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

func (r *salesOrderOpenSearch) Get(request *models.SalesOrderRequest, resultChan chan *models.SalesOrdersChan) {
	response := &models.SalesOrdersChan{}
	requestQuery := r.generateSalesOrderQueryOpenSearchTermRequest("", "", request)
	result, err := r.generateSalesOrderQueryOpenSearchResult(requestQuery, true)

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
	result, err := r.generateSalesOrderQueryOpenSearchResult(requestQuery, true)

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
	result, err := r.generateSalesOrderQueryOpenSearchResult(requestQuery, true)

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
	result, err := r.generateSalesOrderQueryOpenSearchResult(requestQuery, true)

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
	result, err := r.generateSalesOrderQueryOpenSearchResult(requestQuery, true)

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
	result, err := r.generateSalesOrderQueryOpenSearchResult(requestQuery, true)

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
	result, err := r.generateSalesOrderQueryOpenSearchResult(requestQuery, true)

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
				"so_date": map[string]interface{}{
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

func (r *salesOrderOpenSearch) generateSalesOrderQueryOpenSearchResult(openSearchQueryJson []byte, withSalesOrderDetails bool) (*models.SalesOrders, *model.ErrorLog) {
	openSearchQueryResult, err := r.db.Query("sales_orders", openSearchQueryJson)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, "Ada kesalahan, silahkan coba lagi nanti")
		return &models.SalesOrders{}, errorLogData
	}

	if openSearchQueryResult.Hits.Total.Value == 0 {
		err = helper.NewError("data not found")
		errorLogData := helper.WriteLog(err, http.StatusNotFound, "Data tidak ditemukan")
		return &models.SalesOrders{}, errorLogData
	}

	salesOrders := []*models.SalesOrder{}

	if openSearchQueryResult.Hits.Total.Value > 0 {
		for _, v := range openSearchQueryResult.Hits.Hits {
			obj := v.Source.(map[string]interface{})
			storeObj := obj["store"].(map[string]interface{})
			agentObj := obj["agent"].(map[string]interface{})
			userObj := obj["user"].(map[string]interface{})
			orderStatusObj := obj["order_status"].(map[string]interface{})
			orderSourceObj := obj["order_source"].(map[string]interface{})
			brandObj := obj["brand"].(map[string]interface{})
			salesmanObj := map[string]interface{}{}

			if obj["salesman"] != nil {
				salesmanObj = obj["salesman"].(map[string]interface{})
			}

			objFloat := obj["id"].(float64)
			salesOrder := &models.SalesOrder{}
			salesOrder.ID = int(objFloat)
			salesOrder.SoCode = obj["so_code"].(string)
			salesOrder.SoDate = obj["so_date"].(string)

			if obj["so_ref_code"] != nil {
				salesOrder.SoRefCode = models.NullString{sql.NullString{String: obj["so_ref_code"].(string), Valid: true}}
			}

			if obj["so_ref_date"] != nil {
				salesOrder.SoRefDate = models.NullString{sql.NullString{String: obj["so_ref_date"].(string), Valid: true}}
			}

			if obj["g_long"] != nil {
				gLongString := obj["g_long"]
				gLong := gLongString.(float64)
				salesOrder.GLong = models.NullFloat64{sql.NullFloat64{Float64: gLong, Valid: true}}
			}

			if obj["g_lat"] != nil {
				gLatString := obj["g_lat"]
				gLat := gLatString.(float64)
				salesOrder.GLat = models.NullFloat64{sql.NullFloat64{Float64: gLat, Valid: true}}
			}

			if obj["note"] != nil {
				salesOrder.Note = models.NullString{sql.NullString{String: obj["note"].(string), Valid: true}}
			}

			if obj["internal_comment"] != nil {
				salesOrder.InternalComment = models.NullString{sql.NullString{String: obj["internal_comment"].(string), Valid: true}}
			}

			brandId := obj["brand_id"].(float64)
			salesOrder.BrandID = int(brandId)
			salesOrder.BrandName = brandObj["name"].(string)
			salesOrder.TotalTonase = obj["total_tonase"].(float64)
			salesOrder.TotalAmount = obj["total_amount"].(float64)
			salesOrder.AgentName = models.NullString{sql.NullString{String: agentObj["name"].(string), Valid: true}}
			salesOrder.AgentProvinceName = models.NullString{sql.NullString{String: agentObj["province_name"].(string), Valid: true}}
			salesOrder.AgentCityName = models.NullString{sql.NullString{String: agentObj["city_name"].(string), Valid: true}}
			salesOrder.AgentDistrictName = models.NullString{sql.NullString{String: agentObj["district_name"].(string), Valid: true}}
			salesOrder.AgentVillageName = models.NullString{sql.NullString{String: agentObj["village_name"].(string), Valid: true}}

			if agentObj["address"] != nil {
				salesOrder.AgentAddress = models.NullString{sql.NullString{String: agentObj["address"].(string), Valid: true}}
			}

			if agentObj["phone"] != nil {
				salesOrder.AgentPhone = models.NullString{sql.NullString{String: agentObj["phone"].(string), Valid: true}}
			}

			if agentObj["main_mobile_phone"] != nil {
				salesOrder.AgentMainMobilePhone = models.NullString{sql.NullString{String: agentObj["main_mobile_phone"].(string), Valid: true}}
			}

			salesOrder.StoreName = models.NullString{sql.NullString{String: storeObj["name"].(string), Valid: true}}
			salesOrder.StoreCode = models.NullString{sql.NullString{String: storeObj["store_code"].(string), Valid: true}}

			if storeObj["email"] != nil {
				salesOrder.StoreEmail = models.NullString{sql.NullString{String: storeObj["email"].(string), Valid: true}}
			}

			salesOrder.StoreProvinceName = models.NullString{sql.NullString{String: storeObj["province_name"].(string), Valid: true}}
			salesOrder.StoreCityName = models.NullString{sql.NullString{String: storeObj["city_name"].(string), Valid: true}}
			salesOrder.StoreDistrictName = models.NullString{sql.NullString{String: storeObj["district_name"].(string), Valid: true}}
			salesOrder.StoreVillageName = models.NullString{sql.NullString{String: storeObj["village_name"].(string), Valid: true}}

			if storeObj["address"] != nil {
				salesOrder.StoreAddress = models.NullString{sql.NullString{String: storeObj["address"].(string), Valid: true}}
			}

			if storeObj["phone"] != nil {
				salesOrder.StorePhone = models.NullString{sql.NullString{String: storeObj["phone"].(string), Valid: true}}
			}

			if storeObj["main_mobile_phone"] != nil {
				salesOrder.StoreMainMobilePhone = models.NullString{sql.NullString{String: storeObj["main_mobile_phone"].(string), Valid: true}}
			}

			if userObj["first_name"] != nil {
				salesOrder.UserFirstName = models.NullString{sql.NullString{String: userObj["first_name"].(string), Valid: true}}
			}

			if userObj["last_name"] != nil {
				salesOrder.UserLastName = models.NullString{sql.NullString{String: userObj["last_name"].(string), Valid: true}}
			}

			if userObj["email"] != nil {
				salesOrder.UserEmail = models.NullString{sql.NullString{String: userObj["email"].(string), Valid: true}}
			}

			salesOrder.OrderStatusName = orderStatusObj["name"].(string)
			salesOrder.OrderSourceName = orderSourceObj["source_name"].(string)

			if obj["salesman"] != nil {
				salesOrder.SalesmanName = models.NullString{sql.NullString{String: salesmanObj["name"].(string), Valid: true}}
				salesOrder.SalesmanEmail = models.NullString{sql.NullString{String: salesmanObj["email"].(string), Valid: true}}
			}

			salesOrderDetails := []*models.SalesOrderDetail{}
			if withSalesOrderDetails == true {
				salesOrderDetailsObj := obj["sales_order_details"].([]interface{})

				for _, salesOrderDetail := range salesOrderDetailsObj {
					salesOrderDetailJson, _ := json.Marshal(salesOrderDetail)
					salesOrderDetailObj := models.SalesOrderDetail{}
					_ = json.Unmarshal(salesOrderDetailJson, &salesOrderDetailObj)
					salesOrderDetails = append(salesOrderDetails, &salesOrderDetailObj)
				}

				salesOrder.SalesOrderDetails = salesOrderDetails

			}

			layout := "2006-01-02T15:04:05.000000+07:00"
			createdAt, _ := time.Parse(layout, obj["created_at"].(string))
			salesOrder.CreatedAt = &createdAt
			updatedAt, _ := time.Parse(layout, obj["updated_at"].(string))
			salesOrder.UpdatedAt = &updatedAt

			salesOrders = append(salesOrders, salesOrder)
		}
	}

	result := &models.SalesOrders{
		Total:       int64(openSearchQueryResult.Hits.Total.Value),
		SalesOrders: salesOrders,
	}

	return result, &model.ErrorLog{}
}
