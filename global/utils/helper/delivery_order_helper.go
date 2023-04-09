package helper

import (
	"fmt"
	"math/rand"
	"order-service/app/models/constants"
	"strconv"
	"time"
)

func GenerateDOCode(agentID int, orderSourceCode string) string {
	var code string
	var ym string

	acode := fmt.Sprintf("%04d", agentID)
	ym = time.Now().Format("060102")
	rand.Seed(time.Now().UnixNano())
	randcode, _ := Generate(`[A-Z0-9]{6}`)
	code = fmt.Sprintf("D%s%s%s%s", orderSourceCode, acode, ym, randcode)

	return code
}

func GenerateDODetailCode(doID int, agentID int, productID int, uomID int) (string, error) {
	var result string

	rand.Seed(time.Now().UTC().UnixNano())
	randoms, err := Generate(`[a-zA-Z0-9]{10}`)

	if err != nil {
		return result, err
	}

	tinow := time.Now().In(time.UTC).Format(constants.DATE_FORMAT_CODE_GENERATOR)

	result = fmt.Sprintf(
		"%d-%s-%s-%d-%s%d",
		doID, strconv.Itoa(int(agentID)),
		tinow, productID, randoms, uomID,
	)

	return result, nil
}
