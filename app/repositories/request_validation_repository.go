package repositories

import (
	"fmt"
	"order-service/app/models"
	"order-service/global/utils/helper"
	"strings"

	"github.com/bxcodec/dbresolver"
)

type RequestValidationRepositoryInterface interface {
	UniqueValidation(value *models.UniqueRequest, resultChan chan *models.UniqueRequestChan)
	MustActiveValidation(value []*models.MustActiveRequest, resultChan chan *models.MustActiveRequestChan)
	MustEmptyValidation(value *models.MustEmptyValidationRequest, resultChan chan *models.MustEmptyValidationRequestChan)
	AgentIdValidation(agentId, userId int, resultChan chan *models.RequestIdValidationChan)
	StoreIdValidation(storeId, agentId int, resultChan chan *models.RequestIdValidationChan)
	SalesmanIdValidation(salesmanId, agentId int, resultChan chan *models.RequestIdValidationChan)
	BrandIdValidation(brandId, agentId int, resultChan chan *models.RequestIdValidationChan)
	BrandSalesmanValidation(brandId, salesmanId, agentId int, resultChan chan *models.RequestIdValidationChan)
	StoreAddressesValidation(storeCode string, resultChan chan *models.RequestIdValidationChan)
}

type requestValidationRepository struct {
	db dbresolver.DB
}

func InitRequestValidationRepository(db dbresolver.DB) RequestValidationRepositoryInterface {
	return &requestValidationRepository{
		db: db,
	}
}

func (r *requestValidationRepository) UniqueValidation(value *models.UniqueRequest, resultChan chan *models.UniqueRequestChan) {
	response := &models.UniqueRequestChan{}
	var total int64

	query := fmt.Sprintf("SELECT COUNT(*) as total FROM %s WHERE %s = ? AND deleted_at IS NULL", value.Table, value.Field)
	err := r.db.QueryRow(query, value.Value).Scan(&total)

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	} else {
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *requestValidationRepository) MustActiveValidation(values []*models.MustActiveRequest, resultChan chan *models.MustActiveRequestChan) {
	response := &models.MustActiveRequestChan{}
	var total int64
	var query = ""
	fmt.Println(query)
	for k, value := range values {
		if k == len(values)-1 {
			query += fmt.Sprintf("SELECT COUNT(*) as total FROM %s WHERE %s ", value.Table, value.Clause)
		} else {
			query += fmt.Sprintf("SELECT COUNT(*) as total FROM %s WHERE %s UNION ALL ", value.Table, value.Clause)
		}
	}

	q, err := r.db.Query(query)

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
		response.Error = err
		response.ErrorLog = errorLogData
	}

	for q.Next() {
		q.Scan(&total)
		response.Total = append(response.Total, total)
	}
	resultChan <- response
	return
}

func (r *requestValidationRepository) MustEmptyValidation(value *models.MustEmptyValidationRequest, resultChan chan *models.MustEmptyValidationRequestChan) {
	response := &models.MustEmptyValidationRequestChan{}
	query := fmt.Sprintf("SELECT %[2]s as resultQuery FROM %[1]s WHERE %[3]s", value.Table, value.SelectedCollumn, value.Clause)
	q, err := r.db.Query(query)
	if err != nil {
		response.Result = false
		errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
		response.Error = err
		response.ErrorLog = errorLogData
	} else {
		result := ""
		var resultQuery string
		isNotLast := q.Next()
		for isNotLast {
			q.Scan(&resultQuery)
			isNotLast = q.Next()
			result += resultQuery
			if isNotLast {
				result += ", "
			}
		}
		if result == "" {
			response.Result = true
		} else {
			response.Result = false
			response.Message = strings.Replace(value.MessageFormat, "<result>", result, 1)
		}
	}
	resultChan <- response
	return
}

func (r *requestValidationRepository) AgentIdValidation(agentId, userId int, resultChan chan *models.RequestIdValidationChan) {
	response := &models.RequestIdValidationChan{}
	var total int64

	query := fmt.Sprintf("SELECT COUNT(*) as total FROM agents WHERE id = %d AND user_id_updated = %d AND deleted_at IS NULL", agentId, userId)
	err := r.db.QueryRow(query).Scan(&total)

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	} else {
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *requestValidationRepository) StoreIdValidation(storeId, agentId int, resultChan chan *models.RequestIdValidationChan) {
	response := &models.RequestIdValidationChan{}
	var total int64

	query := fmt.Sprintf("SELECT COUNT(*) AS total FROM stores JOIN agent_store ON stores.id = agent_store.store_id JOIN agents ON agent_store.agent_id = agents.id WHERE stores.id = %d AND agents.id = %d AND stores.deleted_at IS NULL", storeId, agentId)
	err := r.db.QueryRow(query).Scan(&total)

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	} else {
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *requestValidationRepository) SalesmanIdValidation(salesmanId, agentId int, resultChan chan *models.RequestIdValidationChan) {
	response := &models.RequestIdValidationChan{}
	var total int64

	query := fmt.Sprintf("SELECT COUNT(*) AS total FROM salesmans JOIN brand_salesman ON salesmans.id = brand_salesman.salesman_id JOIN agents ON brand_salesman.agent_id = agents.id WHERE salesmans.id = %d AND agents.id = %d AND salesmans.deleted_at IS NULL", salesmanId, agentId)
	err := r.db.QueryRow(query).Scan(&total)

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	} else {
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *requestValidationRepository) BrandIdValidation(brandId, agentId int, resultChan chan *models.RequestIdValidationChan) {
	response := &models.RequestIdValidationChan{}
	var total int64

	query := fmt.Sprintf("SELECT COUNT(*) AS total FROM brands JOIN brand_salesman ON brands.id = brand_salesman.brand_id JOIN agents ON brand_salesman.agent_id = agents.id WHERE brands.id = %d AND agents.id = %d AND brands.deleted_at IS NULL", brandId, agentId)
	err := r.db.QueryRow(query).Scan(&total)

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	} else {
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *requestValidationRepository) BrandSalesmanValidation(brandId, salesmanId, agentId int, resultChan chan *models.RequestIdValidationChan) {
	response := &models.RequestIdValidationChan{}
	var total int64

	query := fmt.Sprintf("SELECT COUNT(*) AS total FROM brand_salesman WHERE brand_id = %d AND salesman_id = %d AND agent_id = %d", brandId, salesmanId, agentId)
	err := r.db.QueryRow(query).Scan(&total)

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	} else {
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *requestValidationRepository) StoreAddressesValidation(storeCode string, resultChan chan *models.RequestIdValidationChan) {
	response := &models.RequestIdValidationChan{}
	var total int64

	query := fmt.Sprintf("SELECT COUNT(*) AS total FROM store_addresses JOIN stores ON store_addresses.store_id = stores.id WHERE stores.store_code = '%s'", storeCode)
	err := r.db.QueryRow(query).Scan(&total)

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	} else {
		response.Total = total
		resultChan <- response
		return
	}
}
