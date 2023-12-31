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
	"time"
)

type DeliveryOrderDetailOpenSearchRepositoryInterface interface {
	Create(request *models.DeliveryOrderDetailOpenSearch, doDetail *models.DeliveryOrderDetail, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.DeliveryOrderDetailOpenSearchChan)
	Get(request *models.DeliveryOrderDetailOpenSearchRequest, isCountOnly bool, result chan *models.DeliveryOrderDetailsOpenSearchChan)
	GetByID(request *models.DeliveryOrderRequest, resultChan chan *models.DeliveryOrderDetailOpenSearchChan)
	generateDeliveryOrderQueryOpenSearchResult(openSearchQueryJson []byte, isCountOnly bool) (*models.DeliveryOrderDetailsOpenSearch, *model.ErrorLog)
	generateDeliveryOrderQueryOpenSearchTermRequest(term_field string, term_value interface{}, request *models.DeliveryOrderDetailOpenSearchRequest) []byte
}

type deliveryOrderDetailOpenSearch struct {
	openSearch                    opensearch_dbo.OpenSearchClientInterface
	deliveryOrderDetailRepository repositories.DeliveryOrderDetailRepositoryInterface
}

func InitDeliveryOrderDetailOpenSearchRepository(openSearch opensearch_dbo.OpenSearchClientInterface, deliveryOrderDetailRepository repositories.DeliveryOrderDetailRepositoryInterface) DeliveryOrderDetailOpenSearchRepositoryInterface {
	return &deliveryOrderDetailOpenSearch{
		openSearch:                    openSearch,
		deliveryOrderDetailRepository: deliveryOrderDetailRepository,
	}
}

func (r *deliveryOrderDetailOpenSearch) Create(request *models.DeliveryOrderDetailOpenSearch, doDetail *models.DeliveryOrderDetail, sqlTransaction *sql.Tx, ctx context.Context, resultChan chan *models.DeliveryOrderDetailOpenSearchChan) {
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

	updateDeliveryOrderDetailResultChan := make(chan *models.DeliveryOrderDetailChan)
	go r.deliveryOrderDetailRepository.UpdateByID(doDetail.ID, doDetail, sqlTransaction, ctx, updateDeliveryOrderDetailResultChan)
	updateDeliveryOrderDetailResult := <-updateDeliveryOrderDetailResultChan

	if updateDeliveryOrderDetailResult.Error != nil {
		errorLogData := helper.WriteLog(updateDeliveryOrderDetailResult.Error, http.StatusInternalServerError, nil)
		response.Error = updateDeliveryOrderDetailResult.Error
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
	page := request.Page
	perPage := request.PerPage
	sortField := request.SortField
	sortValue := request.SortValue
	if isCountOnly {
		request.Page = 0
		request.PerPage = 0
		request.SortField = ""
		request.SortValue = ""

	}
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
	if isCountOnly {
		request.Page = page
		request.PerPage = perPage
		request.SortField = sortField
		request.SortValue = sortValue
	}
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

		if len(request.UpdatedAt) > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"updated_at": request.UpdatedAt,
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

		if request.DoDetailID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"id": request.DoDetailID,
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

		if request.SalesOrderID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"sales_order_id": request.SalesOrderID,
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

		if request.SalesmanID != 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"salesman_id": request.SalesmanID,
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
					"qty": request.Qty,
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
			match := map[string]interface{}{
				"query_string": map[string]interface{}{
					"query":            "*" + request.GlobalSearchValue + "*",
					"fields":           []string{"do_code", "do_ref_code", "so_code", "store.name", "store.store_code", "order_status.name^0.5", "qty^3"},
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

			if helper.Contains(constants.UNMAPPED_TYPE_SORT_LIST(), request.SortField) {
				sortValue["unmapped_type"] = "date"
				openSearchQuery["sort"] = []map[string]interface{}{
					{
						request.SortField: sortValue,
					},
				}
			}

			if helper.Contains(constants.DELIVERY_ORDER_DETAIL_SORT_INT_LIST(), request.SortField) {
				openSearchQuery["sort"] = []map[string]interface{}{
					{
						request.SortField: sortValue,
					},
				}
			}

			if helper.Contains(constants.DELIVERY_ORDER_DETAIL_SORT_STRING_LIST(), request.SortField) {
				var field string
				if request.SortField == "store_code" {
					field = "store.store_code"
				} else {
					field = request.SortField
				}
				openSearchQuery["sort"] = []map[string]interface{}{
					{
						field + ".keyword": sortValue,
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
		openSearchQueryResult, err := r.openSearch.Count(constants.DELIVERY_ORDER_DETAILS_INDEX, openSearchQueryJson)
		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return &models.DeliveryOrderDetailsOpenSearch{}, errorLogData
		}

		if openSearchQueryResult <= 0 {
			err = helper.NewError(constants.ERROR_DATA_NOT_FOUND)
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
			err = helper.NewError(constants.ERROR_DATA_NOT_FOUND)
			errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
			return &models.DeliveryOrderDetailsOpenSearch{}, errorLogData
		}

		total = int64(openSearchQueryResult.Hits.Total.Value)

		if openSearchQueryResult.Hits.Total.Value > 0 {
			loc, _ := time.LoadLocation("Asia/Jakarta")
			for _, v := range openSearchQueryResult.Hits.Hits {
				obj := v.Source.(map[string]interface{})
				deliveryOrderDetail := models.DeliveryOrderDetailOpenSearch{}
				objJson, _ := json.Marshal(obj)
				json.Unmarshal(objJson, &deliveryOrderDetail)
				layout := time.RFC3339
				if obj["created_at"] != nil {
					createdAt, _ := time.ParseInLocation(layout, obj["created_at"].(string), loc)
					createdAt = createdAt.In(loc)
					deliveryOrderDetail.CreatedAt = &createdAt
				}
				if obj["updated_at"] != nil {
					updatedAt, _ := time.ParseInLocation(layout, obj["created_at"].(string), loc)
					updatedAt = updatedAt.In(loc)
					deliveryOrderDetail.UpdatedAt = &updatedAt
				}
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
