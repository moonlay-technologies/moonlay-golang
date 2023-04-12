package helper

import (
	"fmt"
	"math/rand"
	"order-service/app/models"
	"order-service/app/models/constants"
	"strconv"
	"strings"
	"time"
)

func GenerateSORefCode(agentID int, soDate string) string {
	var code string
	rand.Seed(time.Now().UnixNano())
	randcode, _ := Generate(`[a-zA-Z0-9]{13}`)
	soDate = strings.ReplaceAll(soDate, "-", "")
	code = fmt.Sprintf("DBO-SOREF-AUTOGEN-%05d-%s-%s", agentID, soDate, randcode)
	return code
}

func GenerateSOCode(agentID int, orderSourceCode string) string {
	var code string
	var ym string

	acode := fmt.Sprintf("%04d", agentID)
	ym = time.Now().Format("060102")
	rand.Seed(time.Now().UnixNano())
	randcode, _ := Generate(`[A-Z0-9]{6}`)
	code = fmt.Sprintf("O%s%s%s%s", orderSourceCode, acode, ym, randcode)

	return code
}

func GenerateSODetailCode(soID int, agentID int, productID int, uomID int) (string, error) {
	var result string

	rand.Seed(time.Now().UTC().UnixNano())
	randoms, err := Generate(`[a-zA-Z0-9]{10}`)

	if err != nil {
		return result, err
	}

	tinow := time.Now().In(time.UTC).Format(constants.DATE_FORMAT_CODE_GENERATOR)

	result = fmt.Sprintf(
		"%d-%s-%s-%d-%s%d",
		soID, strconv.Itoa(int(agentID)),
		tinow, productID, randoms, uomID,
	)

	return result, nil
}

func GenerateUnprocessableErrorMessage(action_name, reason string) string {
	return fmt.Sprintf("Proses %s tidak dapat dilakukan karena %s", action_name, reason)
}

func CheckSalesOrderDetailStatus(soDetailId int, isNot bool, status string, soDetails []*models.SalesOrderDetail) int {
	total := 0

	for _, v := range soDetails {
		if isNot {
			if v.ID != soDetailId && v.OrderStatusName != status {
				total++
			}
		} else {
			if v.ID != soDetailId && v.OrderStatusName == status {
				total++
			}
		}
	}

	return total
}

func GetSOJourneyStatus(orderStatusId int) string {
	var journeyMap = map[int]string{
		5:  constants.SO_STATUS_APPV,
		6:  constants.SO_STATUS_PEND,
		12: constants.SO_STATUS_PEND,
		9:  constants.SO_STATUS_RJC,
		10: constants.SO_STATUS_CNCL,
		7:  constants.SO_STATUS_ORDPRT,
		8:  constants.SO_STATUS_ORDCLS,
		14: constants.SO_STATUS_CLS,
		11: constants.SO_STATUS_OPEN,
	}
	return journeyMap[orderStatusId]
}
