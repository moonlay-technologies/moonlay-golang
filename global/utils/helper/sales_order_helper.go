package helper

import (
	"fmt"
	"math/rand"
	"order-service/app/models"
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

	tinow := time.Now().In(time.UTC).Format("20060102150405")

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
