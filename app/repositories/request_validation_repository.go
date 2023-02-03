package repositories

import (
	"fmt"
	"order-service/app/models"
	"order-service/global/utils/helper"

	"github.com/bxcodec/dbresolver"
)

type RequestValidationRepositoryInterface interface {
	UniqueValidation(value *models.UniqueRequest, resultChan chan *models.UniqueRequestChan)
	MustActiveValidation(value *models.MustActiveRequest, resultChan chan *models.MustActiveRequestChan)
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
