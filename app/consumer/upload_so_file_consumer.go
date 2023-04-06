package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/app/middlewares"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/app/repositories"
	mongoRepositories "order-service/app/repositories/mongod"
	"order-service/global/utils/helper"
	kafkadbo "order-service/global/utils/kafka"
	"strconv"
	"strings"
	"time"

	"github.com/bxcodec/dbresolver"
)

type UploadSOFileConsumerHandlerInterface interface {
	ProcessMessage()
}

type uploadSOFileConsumerHandler struct {
	kafkaClient                 kafkadbo.KafkaClientInterface
	uploadRepository            repositories.UploadRepositoryInterface
	requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface
	requestValidationRepository repositories.RequestValidationRepositoryInterface
	soUploadHistoriesRepository mongoRepositories.SoUploadHistoriesRepositoryInterface
	soUploadErrorLogsRepository mongoRepositories.SoUploadErrorLogsRepositoryInterface
	salesOrderRepository        repositories.SalesOrderRepositoryInterface
	ctx                         context.Context
	args                        []interface{}
	db                          dbresolver.DB
}

func InitUploadSOFileConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, uploadRepository repositories.UploadRepositoryInterface, requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, requestValidationRepository repositories.RequestValidationRepositoryInterface, soUploadHistoriesRepository mongoRepositories.SoUploadHistoriesRepositoryInterface, soUploadErrorLogsRepository mongoRepositories.SoUploadErrorLogsRepositoryInterface, salesOrderRepository repositories.SalesOrderRepositoryInterface, db dbresolver.DB, ctx context.Context, args []interface{}) UploadSOFileConsumerHandlerInterface {
	return &uploadSOFileConsumerHandler{
		kafkaClient:                 kafkaClient,
		uploadRepository:            uploadRepository,
		requestValidationMiddleware: requestValidationMiddleware,
		requestValidationRepository: requestValidationRepository,
		soUploadHistoriesRepository: soUploadHistoriesRepository,
		soUploadErrorLogsRepository: soUploadErrorLogsRepository,
		salesOrderRepository:        salesOrderRepository,
		ctx:                         ctx,
		args:                        args,
		db:                          db,
	}
}

func (c *uploadSOFileConsumerHandler) ProcessMessage() {
	fmt.Println("process ", constants.UPLOAD_SO_FILE_TOPIC)
	topic := c.args[1].(string)
	groupID := c.args[2].(string)
	reader := c.kafkaClient.SetConsumerGroupReader(topic, groupID)

	for {
		m, err := reader.ReadMessage(c.ctx)
		if err != nil {
			break
		}

		fmt.Printf("message so at topic/partition/offset %v/%v/%v \n", m.Topic, m.Partition, m.Offset)
		now := time.Now()

		soUploadHistoryId := m.Value
		var key = string(m.Key[:])

		// var errors []string

		soUploadHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
		go c.soUploadHistoriesRepository.GetByID(string(soUploadHistoryId), false, c.ctx, soUploadHistoryJourneysResultChan)
		soUploadHistoryJourneysResult := <-soUploadHistoryJourneysResultChan
		if soUploadHistoryJourneysResult.Error != nil {
			fmt.Println(soUploadHistoryJourneysResult.Error.Error())
		}

		message := soUploadHistoryJourneysResult.SoUploadHistory
		file, err := c.uploadRepository.ReadFile(message.FilePath)

		if err != nil {
			c.updateSoUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
			continue
		}

		parseFile := string(file)
		data := strings.Split(parseFile, "\n")
		totalRows := int64(len(data) - 1)

		message.TotalRows = &totalRows
		c.updateSoUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_SUCCESS)

		agentIds := map[int]int{}
		results := []map[string]string{}
		for i, v := range data {
			if i == len(data)-1 {
				break
			}
			if i == 0 {
				continue
			}

			row := strings.Split(v, "\"")
			for j := 1; j < len(row); j = j + 2 {
				cell := strings.Split(row[j], ",")
				row[j] = strings.Join(cell, "")
			}
			v = strings.Join(row, "")

			headers := strings.Split(data[0], ",")
			line := strings.Split(v, ",")
			uploadSOField := map[string]string{}
			for j, y := range line {
				if j == 0 {
					agentId, err := strconv.Atoi(y)
					if err == nil {
						agentIds[agentId] = agentIds[agentId] + 1
					}
				}
				uploadSOField[strings.ReplaceAll(headers[j], "\r", "")] = strings.ReplaceAll(y, "\r", "")
			}
			results = append(results, uploadSOField)
		}

		if len(agentIds) > 1 {
			fmt.Println("Error multiple agentId")
			c.updateSoUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
			continue
		}

		brandIds := map[string][]map[string]string{}
		var uploadSOFields []*models.UploadSOField
		for i, v := range results {

			mandatoryError := c.requestValidationMiddleware.UploadMandatoryValidation([]*models.TemplateRequest{
				{
					Field: "IDDistributor",
					Value: v["IDDistributor"],
				},
				{
					Field: "KodeToko",
					Value: v["KodeToko"],
				},
				{
					Field: "IDSalesman",
					Value: v["IDSalesman"],
				},
				{
					Field: "TanggalOrder",
					Value: v["TanggalOrder"],
				},
				{
					Field: "NoOrder",
					Value: v["NoOrder"],
				},
				{
					Field: "TanggalTokoOrder",
					Value: v["TanggalTokoOrder"],
				},
				{
					Field: "KodeMerk",
					Value: v["KodeMerk"],
				},
				{
					Field: "KodeProduk",
					Value: v["KodeProduk"],
				},
				{
					Field: "QTYOrder",
					Value: v["QTYOrder"],
				},
				{
					Field: "UnitProduk",
					Value: v["UnitProduk"],
				},
			})
			if len(mandatoryError) >= 1 {

				if key == "retry" {
					c.updateSoUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := mandatoryError
					c.createSoUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, errors, &now, v)
					continue
				}
			}

			// Get Sales Order By Code
			getSalesOrderResultChan := make(chan *models.SalesOrderChan)
			go c.salesOrderRepository.GetBySoRefCode(v["NoOrder"], false, c.ctx, getSalesOrderResultChan)
			getSalesOrderResult := <-getSalesOrderResultChan

			if getSalesOrderResult.Error != nil || getSalesOrderResult.Total > 0 {

				if key == "retry" {
					c.updateSoUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errorMessage := fmt.Sprintf("No Order %s telah digunakan. Silahkan gunakan No Order lain.", v["NoOrder"])
					if getSalesOrderResult.Error != nil {
						errorMessage = getSalesOrderResult.Error.Error()
					}
					fmt.Println(errorMessage)
					errors := []string{errorMessage}

					c.createSoUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, errors, &now, v)

					continue
				}
			}

			if brandIds[v["NoOrder"]] != nil {

				var isBreak bool

				if brandIds[v["NoOrder"]][0]["KodeMerk"] != v["KodeMerk"] {
					fmt.Println("No Order " + v["NoOrder"] + " memiliki lebih dari 1 Kode Merk")

					c.updateSoUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)

					isBreak = true
					break
				}

				for _, x := range brandIds[v["NoOrder"]] {
					if x["KodeProduk"] == v["KodeProduk"] && x["UnitProduk"] == v["UnitProduk"] {
						if key == "retry" {
							c.updateSoUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)

							isBreak = true
							break
						} else {
							var errors []string
							errors = append(errors, fmt.Sprintf("Duplikat row untuk No Order %s", v["NoOrder"]))

							c.createSoUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, errors, &now, v)

							isBreak = false
							break
						}
					}
				}

				if isBreak {
					break
				} else {
					continue
				}
			}

			brandIds[v["NoOrder"]] = append(brandIds[v["NoOrder"]], map[string]string{
				"KodeMerk":   v["KodeMerk"],
				"KodeProduk": v["KodeProduk"],
				"UnitProduk": v["UnitProduk"],
			})

			intType := []*models.TemplateRequest{
				{
					Field: "IDDistributor",
					Value: v["IDDistributor"],
				},
				{
					Field: "IDSalesman",
					Value: v["IDSalesman"],
				},
				{
					Field: "KodeMerk",
					Value: v["KodeMerk"],
				},
				{
					Field: "QTYOrder",
					Value: v["QTYOrder"],
				},
			}
			if v["IDAlamat"] != "" {
				intType = append(intType, &models.TemplateRequest{
					Field: "IDAlamat",
					Value: v["IDAlamat"],
				})
			}
			intTypeResult, intTypeError := c.requestValidationMiddleware.UploadIntTypeValidation(intType)
			if len(intTypeError) >= 1 {

				if key == "retry" {
					c.updateSoUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)

					break
				} else {
					errors := intTypeError

					c.createSoUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, errors, &now, v)

					continue
				}
			}

			if intTypeResult["QTYOrder"] < 1 {
				if key == "retry" {
					c.updateSoUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)

					break
				} else {
					errors := []string{"Quantity harus lebih dari 0"}

					c.createSoUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, errors, &now, v)

					continue
				}
			}

			mustActiveField := []*models.MustActiveRequest{
				helper.GenerateMustActive("agents", "IDDistributor", intTypeResult["IDDistributor"], "active"),
				helper.GenerateMustActive("users", "user_id", 0, "ACTIVE"),
				{
					Table:    "salesmans",
					ReqField: "IDSalesman",
					Clause:   fmt.Sprintf(constants.CLAUSE_ID_VALIDATION, intTypeResult["IDSalesman"]),
					Id:       intTypeResult["IDSalesman"],
				},
				{
					Table:         "brands",
					ReqField:      "KodeMerk",
					Clause:        fmt.Sprintf("id = %d AND status_active = %d", intTypeResult["KodeMerk"], 1),
					Id:            intTypeResult["KodeMerk"],
					CustomMessage: fmt.Sprintf("Kode Merek = %d Tidak Terdaftar pada Distributor %s. Silahkan gunakan Kode Merek yang lain.", intTypeResult["KodeMerk"], soUploadHistoryJourneysResult.SoUploadHistory.AgentName),
				},
				{
					Table:         "products",
					ReqField:      "KodeProduk",
					Clause:        fmt.Sprintf("sku = '%s' AND isActive = %d", v["KodeProduk"], 1),
					Id:            v["KodeProduk"],
					CustomMessage: fmt.Sprintf("Kode SKU = %s dengan Merek %s sudah Tidak Aktif. Silahkan gunakan Kode SKU yang lain.", v["KodeProduk"], v["NamaMerk"]),
				},
				{
					Table:    "uoms",
					ReqField: "UnitProduk",
					Clause:   fmt.Sprintf("code = '%s' AND deleted_at IS NULL", v["UnitProduk"]),
					Id:       v["UnitProduk"],
				},
			}
			mustActiveError := c.requestValidationMiddleware.UploadMustActiveValidation(mustActiveField)
			if len(mustActiveError) >= 1 {
				if key == "retry" {
					c.updateSoUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)

					break
				} else {
					errors := mustActiveError

					c.createSoUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, errors, &now, v)

					continue
				}
			}

			brandSalesmanResultChan := make(chan *models.RequestIdValidationChan)
			go c.requestValidationRepository.BrandSalesmanValidation(intTypeResult["KodeMerk"], intTypeResult["IDSalesman"], intTypeResult["IDDistributor"], brandSalesmanResultChan)
			brandSalesmanResult := <-brandSalesmanResultChan

			if brandSalesmanResult.Total < 1 {
				if key == "retry" {
					c.updateSoUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)

					break
				} else {
					var errors []string

					errors = append(errors, fmt.Sprintf("Kode Merek = %d Tidak Terdaftar pada Distributor %s. Silahkan gunakan Kode Merek yang lain.", intTypeResult["KodeMerk"], soUploadHistoryJourneysResult.SoUploadHistory.AgentName))
					errors = append(errors, fmt.Sprintf("ID Salesman = %d Tidak Terdaftar pada Distributor %s. Silahkan gunakan ID Salesman yang lain.", intTypeResult["IDSalesman"], soUploadHistoryJourneysResult.SoUploadHistory.AgentName))

					c.createSoUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, errors, &now, v)

					continue
				}
			}

			storeAddressesResultChan := make(chan *models.RequestIdValidationChan)
			go c.requestValidationRepository.StoreAddressesValidation(v["KodeToko"], storeAddressesResultChan)
			storeAddressesResult := <-storeAddressesResultChan

			if storeAddressesResult.Total < 1 {
				if key == "retry" {
					c.updateSoUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)

					break
				} else {
					errors := []string{fmt.Sprintf("Alamat Utama pada Kode Toko = %s Tidak Ditemukan. Silahkan gunakan Alamat Toko yang lain.", v["KodeToko"])}

					c.createSoUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, errors, &now, v)

					continue
				}
			}

			tanggalOrder, err := helper.ParseDDYYMMtoYYYYMMDD(v["TanggalOrder"])
			if err != nil {
				if key == "retry" {
					c.updateSoUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)

					break
				} else {
					errors := []string{fmt.Sprintf("Format Tanggal Order = %s Salah, silahkan sesuaikan dengan format DD-MMM-YYYY, contoh 15/12/2021", v["TanggalOrder"])}

					c.createSoUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, errors, &now, v)

					continue
				}
			}

			nowWIB := time.Now().UTC().Add(7 * time.Hour)
			duration := time.Hour*time.Duration(nowWIB.Hour()) + time.Minute*time.Duration(nowWIB.Minute()) + time.Second*time.Duration(nowWIB.Second()) + time.Nanosecond*time.Duration(nowWIB.Nanosecond())

			parseTangalOrder, _ := time.Parse(constants.DATE_FORMAT_COMMON, tanggalOrder)
			if parseTangalOrder.Add(duration - 1*time.Minute).After(nowWIB) {
				if key == "retry" {
					c.updateSoUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)

					break
				} else {
					errors := []string{fmt.Sprintf("Tanggal Order = %s tidak boleh melebihi tanggal upload", v["TanggalOrder"])}

					c.createSoUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, errors, &now, v)

					continue
				}
			}

			tanggalTokoOrder, err := helper.ParseDDYYMMtoYYYYMMDD(v["TanggalTokoOrder"])
			if err != nil {
				if key == "retry" {
					c.updateSoUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)

					break
				} else {
					errors := []string{fmt.Sprintf("Format Tanggal Toko Order = %s Salah, silahkan sesuaikan dengan format DD-MMM-YYYY, contoh 15/12/2021", v["TanggalTokoOrder"])}

					c.createSoUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, errors, &now, v)

					continue
				}
			}

			parseTanggalTokoOrder, _ := time.Parse(constants.DATE_FORMAT_COMMON, tanggalTokoOrder)
			if parseTanggalTokoOrder.Add(duration - 1*time.Minute).After(parseTangalOrder.Add(duration)) {
				if key == "retry" {
					c.updateSoUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)

					break
				} else {
					errors := []string{fmt.Sprintf("Tanggal Toko Order = %s tidak boleh melebihi Tanggal Order", v["TanggalOrder"])}

					c.createSoUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, errors, &now, v)

					continue
				}
			}

			var uploadSOField models.UploadSOField
			uploadSOField.TanggalOrder = tanggalOrder
			uploadSOField.TanggalTokoOrder = tanggalTokoOrder

			uploadSOField.UploadSOFieldMap(v, int(*message.UploadedBy), message.ID.Hex())
			uploadSOField.BulkCode = message.BulkCode
			uploadSOField.ErrorLine = i + 2
			uploadSOField.UploadType = key
			uploadSOFields = append(uploadSOFields, &uploadSOField)
		}

		keyKafka := []byte(message.RequestId)
		messageKafka, _ := json.Marshal(uploadSOFields)

		err = c.kafkaClient.WriteToTopic(constants.UPLOAD_SO_ITEM_TOPIC, keyKafka, messageKafka)

		if err != nil {
			c.updateSoUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
			continue
		}

	}

	if err := reader.Close(); err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
}

func (c *uploadSOFileConsumerHandler) createSoUploadErrorLog(errorLine int, agentId, soUploadHistoryId, requestId, agentName, bulkCode string, errors []string, now *time.Time, item map[string]string) {
	soUploadErrorLog := &models.SoUploadErrorLog{}
	soUploadErrorLog.SoUploadErrorLogsMap(errorLine, soUploadHistoryId, requestId, bulkCode, errors, now)
	rowData := &models.RowDataSoUploadErrorLog{}
	rowData.RowDataSoUploadErrorLogMap(item, agentName)
	soUploadErrorLog.RowData = *rowData

	soUploadErrorLogResultChan := make(chan *models.SoUploadErrorLogChan)
	go c.soUploadErrorLogsRepository.Insert(soUploadErrorLog, c.ctx, soUploadErrorLogResultChan)
}

func (c *uploadSOFileConsumerHandler) updateSoUploadHistories(message *models.SoUploadHistory, status string) {
	message.UpdatedAt = time.Now()
	message.Status = status
	soUploadHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
	go c.soUploadHistoriesRepository.UpdateByID(message.ID.Hex(), message, c.ctx, soUploadHistoryJourneysResultChan)
}
