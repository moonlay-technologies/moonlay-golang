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
)

type SalesOrderDetailOpenSearchRepositoryInterface interface {
	Create(request *models.SalesOrderDetailOpenSearch, resultChan chan *models.SalesOrderDetailOpenSearchChan)
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
