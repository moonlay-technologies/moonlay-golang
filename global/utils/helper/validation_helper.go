package helper

import (
	"fmt"
	"order-service/app/models"
)

func GenerateMustActive(table string, reqField string, id int, status string) *models.MustActiveRequest {
	return &models.MustActiveRequest{
		Table:    table,
		ReqField: reqField,
		Clause:   fmt.Sprintf("id = %d AND status = '%s'", id, status),
		Id:       id,
	}
}
