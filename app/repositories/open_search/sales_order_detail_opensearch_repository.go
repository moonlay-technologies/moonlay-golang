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

type SalesOrderDetailOpenSearchRepositoryInterface interface {
	Create(request *models.SalesOrderDetailOpenSearch, resultChan chan *models.SalesOrderDetailOpenSearchChan)
	Get(request *models.GetSalesOrderDetailRequest, isCountOnly bool, result chan *models.SalesOrderDetailsOpenSearchChan)
	generateSalesOrderDetailQueryOpenSearchTermRequest(term_field string, term_value interface{}, request *models.GetSalesOrderDetailRequest) []byte
	generateSalesOrderDetailQueryOpenSearchResult(openSearchQueryJson []byte, isCountOnly bool) (*models.SalesOrderDetailsOpenSearch, *model.ErrorLog)
}

type salesOrderDetailOpenSearch struct {
	openSearch opensearch_dbo.OpenSearchClientInterface
}

func InitSalesOrderDetailOpenSearchRepository(openSearch opensearch_dbo.OpenSearchClientInterface) SalesOrderDetailOpenSearchRepositoryInterface {
	return &salesOrderDetailOpenSearch{
		openSearch: openSearch,
	}
}

func (r *salesOrderDetailOpenSearch) Create(request *models.SalesOrderDetailOpenSearch, resultChan chan *models.SalesOrderDetailOpenSearchChan) {
	response := &models.SalesOrderDetailOpenSearchChan{}
	salesOrderJson, _ := json.Marshal(request)
	st, err := r.openSearch.CreateDocument(constants.SALES_ORDER_DETAILS_INDEX, request.SoDetailCode, salesOrderJson)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}
	fmt.Println("so detail code : ", st)
	response.Error = nil
	response.ErrorLog = &model.ErrorLog{}
	resultChan <- response
	return
}

func (r *salesOrderDetailOpenSearch) Get(request *models.GetSalesOrderDetailRequest, isCountOnly bool, resultChan chan *models.SalesOrderDetailsOpenSearchChan) {
	response := &models.SalesOrderDetailsOpenSearchChan{}
	if isCountOnly {
		request.Page = 0
		request.PerPage = 0
		request.SortField = ""
		request.SortValue = ""
	}
	requestQuery := r.generateSalesOrderDetailQueryOpenSearchTermRequest("", "", request)
	result, err := r.generateSalesOrderDetailQueryOpenSearchResult(requestQuery, isCountOnly)

	if err.Err != nil {
		response.Error = err.Err
		response.ErrorLog = err
		resultChan <- response
		return
	}

	response.SalesOrderDetails = result.SalesOrderDetails
	response.Total = result.Total
	resultChan <- response
	return
}

func (r *salesOrderDetailOpenSearch) generateSalesOrderDetailQueryOpenSearchTermRequest(term_field string, term_value interface{}, request *models.GetSalesOrderDetailRequest) []byte {
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
					"salesman_id": request.SalesmanID,
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

		if request.SoID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"so_id": request.SoID,
				},
			}

			filters = append(filters, filter)
		}

		if request.ProductID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"product_id": request.ProductID,
				},
			}

			filters = append(filters, filter)
		}

		if request.CategoryID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"category_id": request.CategoryID,
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
					"order_status.id": request.OrderStatusID,
				},
			}

			filters = append(filters, filter)
		}

		if request.ProvinceID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store_province_id": request.ProvinceID,
				},
			}

			filters = append(filters, filter)
		}

		if request.CityID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store_city_id": request.CityID,
				},
			}

			filters = append(filters, filter)
		}

		if request.DistrictID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store_district_id": request.DistrictID,
				},
			}

			filters = append(filters, filter)
		}

		if request.VillageID > 0 {
			filter := map[string]interface{}{
				"term": map[string]interface{}{
					"store_village_id": request.VillageID,
				},
			}

			filters = append(filters, filter)
		}

		if len(request.StartCreatedAt) > 0 && len(request.EndCreatedAt) > 0 {
			startTimeadj, err := time.Parse(constants.DATE_FORMAT_EXPORT_CREATED_AT, request.StartCreatedAt+constants.DATE_TIME_ZERO_HOUR_ADDITIONAL)
			if err == nil {
				startTimeadj = startTimeadj.Add(time.Hour * -7)
				request.StartCreatedAt = startTimeadj.Format(time.RFC3339)
			} else {
				fmt.Println("error = ", err.Error())
			}
			endTimeadj, err := time.Parse(constants.DATE_FORMAT_EXPORT_CREATED_AT, request.EndCreatedAt+constants.DATE_TIME_ZERO_HOUR_ADDITIONAL)
			if err == nil {
				endTimeadj = endTimeadj.Add(time.Hour * +16).Add(time.Minute * 59).Add(time.Second * 59)
				request.EndCreatedAt = endTimeadj.Format(time.RFC3339)
			} else {
				fmt.Println("error = ", err.Error())
			}
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
				openSearchQuery["sort"] = []map[string]interface{}{
					{
						request.SortField: sortValue,
					},
				}
			}

			if helper.Contains(constants.SALES_ORDER_DETAIL_SORT_INT_LIST(), request.SortField) {
				openSearchQuery["sort"] = []map[string]interface{}{
					{
						request.SortField: sortValue,
					},
				}
			}

			if helper.Contains(constants.SALES_ORDER_DETAIL_SORT_STRING_LIST(), request.SortField) {
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

func (r *salesOrderDetailOpenSearch) generateSalesOrderDetailQueryOpenSearchResult(openSearchQueryJson []byte, isCountOnly bool) (*models.SalesOrderDetailsOpenSearch, *model.ErrorLog) {
	salesOrderDetails := []*models.SalesOrderDetailOpenSearch{}
	var total int64 = 0

	if isCountOnly {
		openSearchQueryResult, err := r.openSearch.Count(constants.SALES_ORDER_DETAILS_INDEX, openSearchQueryJson)
		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return &models.SalesOrderDetailsOpenSearch{}, errorLogData
		}

		if openSearchQueryResult <= 0 {
			err = helper.NewError("sales_order_details_opensearch data not found")
			errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
			return &models.SalesOrderDetailsOpenSearch{}, errorLogData
		}

		total = openSearchQueryResult
	} else {
		openSearchQueryResult, err := r.openSearch.Query(constants.SALES_ORDER_DETAILS_INDEX, openSearchQueryJson)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			return &models.SalesOrderDetailsOpenSearch{}, errorLogData
		}

		if openSearchQueryResult.Hits.Total.Value == 0 {
			err = helper.NewError(helper.DefaultStatusText[http.StatusNotFound])
			errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
			return &models.SalesOrderDetailsOpenSearch{}, errorLogData
		}

		total = int64(openSearchQueryResult.Hits.Total.Value)

		if openSearchQueryResult.Hits.Total.Value > 0 {
			loc, _ := time.LoadLocation("Asia/Jakarta")
			for _, v := range openSearchQueryResult.Hits.Hits {
				obj := v.Source.(map[string]interface{})

				salesOrderDetail := &models.SalesOrderDetailOpenSearch{}
				jsonStr, err := json.Marshal(v.Source)
				if err != nil {
					fmt.Println(err)
				}
				if err := json.Unmarshal(jsonStr, &salesOrderDetail); err != nil {
					fmt.Println(err)
				}

				layout := time.RFC3339
				if obj["created_at"] != nil {
					createdAt, _ := time.Parse(layout, obj["created_at"].(string))
					createdAt = createdAt.In(loc)
					salesOrderDetail.CreatedAt = &createdAt
				}
				if obj["updated_at"] != nil {
					updatedAt, _ := time.Parse(layout, obj["updated_at"].(string))
					updatedAt = updatedAt.In(loc)
					salesOrderDetail.UpdatedAt = &updatedAt
				}

				salesOrderDetails = append(salesOrderDetails, salesOrderDetail)
			}
		}
	}

	result := &models.SalesOrderDetailsOpenSearch{
		Total:             total,
		SalesOrderDetails: salesOrderDetails,
	}

	return result, &model.ErrorLog{}
}
