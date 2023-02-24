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
	MustActiveValidation(value *models.MustActiveRequest, resultChan chan *models.MustActiveRequestChan)
	MustEmptyValidation(value *models.MustEmptyValidationRequest, resultChan chan *models.MustEmptyValidationRequestChan)
}

type requestValidationRepository struct {
	db dbresolver.DB
}

func InitUniqueRequestValidationRepository(db dbresolver.DB) RequestValidationRepositoryInterface {
	return &requestValidationRepository{
		db: db,
	}
}

func (r *requestValidationRepository) UniqueValidation(value *models.UniqueRequest, resultChan chan *models.UniqueRequestChan) {
	response := &models.UniqueRequestChan{}
	var total int64

	query := fmt.Sprintf("SELECT COUNT(*) as total FROM %s WHERE %s = ?", value.Table, value.Field)
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

func (r *requestValidationRepository) MustActiveValidation(value *models.MustActiveRequest, resultChan chan *models.MustActiveRequestChan) {
	response := &models.MustActiveRequestChan{}
	var total int64

	query := fmt.Sprintf("SELECT COUNT(*) as total FROM %s WHERE %s ", value.Table, value.Clause)
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

func (r *requestValidationRepository) MustEmptyValidation(value *models.MustEmptyValidationRequest, resultChan chan *models.MustEmptyValidationRequestChan) {
	response := &models.MustEmptyValidationRequestChan{}
	query := fmt.Sprintf("SELECT %[4]s as resultQuery FROM %[1]s JOIN %[2]s ON %[2]s.id = %[1]s.%[3]s WHERE %[5]s", value.Table, value.TableJoin, value.ForeignKey, value.SelectedCollumn, value.Clause)
	q, err := r.db.Query(query)
	if err != nil {
		response.Result = false
		errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
		response.Error = err
		response.ErrorLog = errorLogData
	} else {
		result := ""
		var resultQuery string
		for q.Next() {
			q.Scan(&resultQuery)
			result += resultQuery + ", "
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
