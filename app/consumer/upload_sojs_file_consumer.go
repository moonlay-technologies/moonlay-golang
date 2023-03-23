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

type UploadSOSJFileConsumerHandlerInterface interface {
	ProcessMessage()
}

type uploadSOSJFileConsumerHandler struct {
	kafkaClient                   kafkadbo.KafkaClientInterface
	uploadRepository              repositories.UploadRepositoryInterface
	requestValidationMiddleware   middlewares.RequestValidationMiddlewareInterface
	requestValidationRepository   repositories.RequestValidationRepositoryInterface
	sosjUploadHistoriesRepository mongoRepositories.SOSJUploadHistoriesRepositoryInterface
	sosjUploadErrorLogsRepository mongoRepositories.SosjUploadErrorLogsRepositoryInterface
	ctx                           context.Context
	args                          []interface{}
	db                            dbresolver.DB
}

func InitUploadSOSJFileConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, uploadRepository repositories.UploadRepositoryInterface, requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, requestValidationRepository repositories.RequestValidationRepositoryInterface, sosjUploadHistoriesRepository mongoRepositories.SOSJUploadHistoriesRepositoryInterface, sosjUploadErrorLogsRepository mongoRepositories.SosjUploadErrorLogsRepositoryInterface, db dbresolver.DB, ctx context.Context, args []interface{}) UploadSOFileConsumerHandlerInterface {
	return &uploadSOSJFileConsumerHandler{
		kafkaClient:                   kafkaClient,
		uploadRepository:              uploadRepository,
		requestValidationMiddleware:   requestValidationMiddleware,
		requestValidationRepository:   requestValidationRepository,
		sosjUploadHistoriesRepository: sosjUploadHistoriesRepository,
		sosjUploadErrorLogsRepository: sosjUploadErrorLogsRepository,
		ctx:                           ctx,
		args:                          args,
		db:                            db,
	}
}

func (c *uploadSOSJFileConsumerHandler) ProcessMessage() {
	fmt.Println("process ", constants.UPLOAD_SOSJ_FILE_TOPIC)
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

		sosjUploadHistoryId := m.Value
		var key = string(m.Key[:])
		var errors []string

		sosjUploadHistoryJourneysResultChan := make(chan *models.GetSosjUploadHistoryResponseChan)
		go c.sosjUploadHistoriesRepository.GetByID(string(sosjUploadHistoryId), false, c.ctx, sosjUploadHistoryJourneysResultChan)
		sosjUploadHistoryJourneysResult := <-sosjUploadHistoryJourneysResultChan
		if sosjUploadHistoryJourneysResult.Error != nil {
			fmt.Println(sosjUploadHistoryJourneysResult.Error.Error())
		}

		message := &sosjUploadHistoryJourneysResult.SosjUploadHistories.UploadHistory
		file, err := c.uploadRepository.ReadFile(message.FilePath)

		if err != nil {
			c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
			continue
		}

		parseFile := string(file)
		data := strings.Split(parseFile, "\n")
		totalRows := int64(len(data) - 1)

		message.TotalRows = &totalRows
		c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_SUCCESS)

		results := []map[string]string{}
		var agentId string
		for i, v := range data {

			if i == len(data)-1 {
				break
			}

			line := strings.Split(v, ",")
			if i == 0 {
				agentId = line[1]
				continue
			}
			if i == 1 {
				continue
			}

			uploadSOSJField := map[string]string{}
			uploadSOSJField["_0"] = agentId
			for j, y := range line {
				uploadSOSJField[fmt.Sprintf("_%d", j+1)] = y
			}
			results = append(results, uploadSOSJField)
		}

		var uploadSOSJFields []*models.UploadSOSJField

		noSuratJalan := []string{}
		for i, v := range results {
			rowDataSosjUploadErrorLogResultChan := make(chan *models.RowDataSosjUploadErrorLogChan)
			go c.uploadRepository.GetSosjRowData(v["_0"], v["_4"], v["_5"], v["_6"], v["_11"], v["_12"], v["_13"], rowDataSosjUploadErrorLogResultChan)
			rowDataSosjUploadErrorLogResult := <-rowDataSosjUploadErrorLogResultChan
			rowData := &models.RowDataSosjUploadErrorLog{}
			rowData.RowDataSosjUploadErrorLogMap(*rowDataSosjUploadErrorLogResult.RowDataSosjUploadErrorLog, v)

			mandatoryError := c.requestValidationMiddleware.UploadMandatoryValidation([]*models.TemplateRequest{
				{
					Field: "IDDistributor",
					Value: rowData.AgentId,
				},
				{
					Field: "Status",
					Value: rowData.SjStatus,
				},
				{
					Field: "NoSuratJalan",
					Value: rowData.SjNo,
				},
				{
					Field: "TglSuratJalan",
					Value: rowData.SjDate,
				},
				{
					Field: "KodeTokoDBO",
					Value: rowData.StoreCode,
				},
				{
					Field: "IDMerk",
					Value: rowData.BrandId,
				},
				{
					Field: "KodeProdukDBO",
					Value: rowData.ProductCode,
				},
				{
					Field: "Qty",
					Value: rowData.DeliveryQty,
				},
				{
					Field: "Unit",
					Value: rowData.ProductUnit,
				},
			})
			if len(mandatoryError) > 1 {
				if key == "retry" {

					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)

					break
				} else {
					errors := mandatoryError

					c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, message.AgentName, message.BulkCode, errors, &now, *rowData)
					continue
				}
			}

			intType := []*models.TemplateRequest{
				{
					Field: "IDDistributor",
					Value: rowData.AgentId,
				},
				{
					Field: "IDMerk",
					Value: rowData.BrandId,
				},
				{
					Field: "Qty",
					Value: rowData.DeliveryQty,
				},
				{
					Field: "Unit",
					Value: rowData.ProductUnit,
				},
			}
			if rowData.WhCode != "" {
				intType = append(intType, &models.TemplateRequest{
					Field: "KodeGudang",
					Value: rowData.WhCode,
				})
			}
			if rowData.SalesmanId != "" {
				intType = append(intType, &models.TemplateRequest{
					Field: "IDSalesman",
					Value: rowData.SalesmanId,
				})
			}
			intTypeResult, intTypeError := c.requestValidationMiddleware.UploadIntTypeValidation(intType)
			if len(intTypeError) > 1 {
				if key == "retry" {

					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)

					break
				} else {
					errors := intTypeError

					c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, message.AgentName, message.BulkCode, errors, &now, *rowData)

					continue
				}
			}

			if intTypeResult["Qty"] < 1 {
				if key == "retry" {

					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors = []string{"Quantity harus lebih dari 0"}

					c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, message.AgentName, message.BulkCode, errors, &now, *rowData)
					continue
				}
			}

			mustActiveField := []*models.MustActiveRequest{
				helper.GenerateMustActive("agents", "IDDistributor", intTypeResult["IDDistributor"], "active"),
				helper.GenerateMustActive("users", "user_id", int(*message.UploadedBy), "ACTIVE"),
				{
					Table:    "brands",
					ReqField: "IDMerk",
					Clause:   fmt.Sprintf("id = %d AND status_active = %d", intTypeResult["IDMerk"], 1),
					Id:       intTypeResult["IDMerk"],
				},
				{
					Table:    "stores",
					ReqField: "KodeTokoDBO",
					Clause:   fmt.Sprintf("IF((SELECT COUNT(store_code) FROM stores WHERE store_code = '%s'), stores.store_code = '%s', stores.alias_code = '%s') AND status = 'active'", rowData.StoreCode, rowData.StoreCode, rowData.StoreCode),
					Id:       rowData.StoreCode,
				},
				{
					Table:    "products",
					ReqField: "KodeProdukDBO",
					Clause:   fmt.Sprintf("IF((SELECT COUNT(SKU) FROM products WHERE SKU = '%s'), products.SKU = '%s', products.aliasSku = '%s') AND isActive = %d", rowData.ProductCode, rowData.ProductCode, rowData.ProductCode, 1),
					Id:       rowData.ProductCode,
				},
				{
					Table:    "uoms",
					ReqField: "Unit",
					Clause:   fmt.Sprintf("id = %d AND deleted_at IS NULL", intTypeResult["Unit"]),
					Id:       intTypeResult["Unit"],
				},
			}
			mustActiveError := c.requestValidationMiddleware.UploadMustActiveValidation(mustActiveField)
			if len(mustActiveError) > 1 {
				if key == "retry" {
					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := mustActiveError

					c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, message.AgentName, message.BulkCode, errors, &now, *rowData)

					continue
				}
			}

			if len(v["_12"]) > 0 {
				brandSalesman := make(chan *models.RequestIdValidationChan)
				go c.requestValidationRepository.BrandSalesmanValidation(intTypeResult["IDMerk"], intTypeResult["IDSalesman"], intTypeResult["IDDistributor"], brandSalesman)
				brandSalesmanResult := <-brandSalesman

				if brandSalesmanResult.Total < 1 {

					if key == "retry" {
						c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
						break
					} else {
						errors := []string{}
						errors = append(errors, fmt.Sprintf("Kode Merek = %d Tidak Terdaftar pada Distributor %s. Silahkan gunakan Kode Merek yang lain.", intTypeResult["IDMerk"], message.AgentName))
						errors = append(errors, fmt.Sprintf("ID Salesman = %d Tidak Terdaftar pada Distributor %s. Silahkan gunakan ID Salesman yang lain.", intTypeResult["IDSalesman"], message.AgentName))
						errors = append(errors, fmt.Sprintf("Salesman di Kode Toko = %s untuk Merek %s Tidak Terdaftar. Silahkan gunakan ID Salesman yang terdaftar.", rowData.StoreCode, rowData.ProductName.String))

						c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, message.AgentName, message.BulkCode, errors, &now, *rowData)

						continue
					}
				}
			}
			storeAddresses := make(chan *models.RequestIdValidationChan)
			go c.requestValidationRepository.StoreAddressesValidation(rowData.StoreCode, storeAddresses)
			storeAddressesResult := <-storeAddresses

			if storeAddressesResult.Total < 1 {
				if key == "retry" {

					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{fmt.Sprintf("Alamat Utama pada Kode Toko = %s Tidak Ditemukan. Silahkan gunakan Alamat Toko yang lain.", rowData.StoreCode)}

					c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, message.AgentName, message.BulkCode, errors, &now, *rowData)

					continue
				}
			}

			var uploadSOSJField models.UploadSOSJField
			uploadSOSJField.TglSuratJalan, err = helper.ParseDDYYMMtoYYYYMMDD(rowData.SjDate)
			if err != nil {
				if key == "retry" {

					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors = []string{fmt.Sprintf("Format Tanggal Order = %s Salah, silahkan sesuaikan dengan format DD-MMM-YYYY, contoh 15/12/2021", rowData.SjDate)}

					c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, message.AgentName, message.BulkCode, errors, &now, *rowData)

					continue
				}
			}
			uploadSOSJField.UploadSOSJFieldMap(v, intTypeResult["IDDistributor"], int(*message.UploadedBy))

			checkIfNoSuratJalanExist := helper.InSliceString(noSuratJalan, rowData.SjNo)
			if checkIfNoSuratJalanExist {

				for i := range uploadSOSJFields {
					brandId, _ := strconv.Atoi(rowData.BrandId)
					if uploadSOSJFields[i].NoSuratJalan == rowData.SjNo && uploadSOSJFields[i].IDMerk != brandId {
						uploadSOSJFields[i].NoSuratJalan = uploadSOSJFields[i].NoSuratJalan + "-" + strconv.Itoa(uploadSOSJFields[i].IDMerk)
						uploadSOSJField.NoSuratJalan = rowData.SjNo + "-" + rowData.BrandId
						break
					} else {
						uploadSOSJField.NoSuratJalan = rowData.SjNo
					}
				}

			} else {
				uploadSOSJField.NoSuratJalan = rowData.SjNo
				noSuratJalan = append(noSuratJalan, rowData.SjNo)
			}

			uploadSOSJField.BulkCode = message.BulkCode
			uploadSOSJField.SosjUploadHistoryId = message.ID.Hex()
			uploadSOSJField.ErrorLine = i + 3
			uploadSOSJField.UploadType = key
			uploadSOSJFields = append(uploadSOSJFields, &uploadSOSJField)

		}

		keyKafka := []byte(message.RequestId)
		messageKafka, _ := json.Marshal(uploadSOSJFields)

		err = c.kafkaClient.WriteToTopic(constants.UPLOAD_SOSJ_ITEM_TOPIC, keyKafka, messageKafka)

		if err != nil {
			c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
			continue
		}

	}

	if err := reader.Close(); err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
}

func (c *uploadSOSJFileConsumerHandler) createSosjUploadErrorLog(errorLine int, agentId, sosjUploadHistoryId, requestId, agentName, bulkCode string, errors []string, now *time.Time, rowData models.RowDataSosjUploadErrorLog) {
	sosjUploadErrorLog := &models.SosjUploadErrorLog{}
	sosjUploadErrorLog.SosjUploadErrorLogsMap(errorLine, agentId, sosjUploadHistoryId, requestId, agentName, bulkCode, errors, now)
	sosjUploadErrorLog.RowData = rowData

	sosjUploadErrorLogResultChan := make(chan *models.SosjUploadErrorLogChan)
	go c.sosjUploadErrorLogsRepository.Insert(sosjUploadErrorLog, c.ctx, sosjUploadErrorLogResultChan)
}

func (c *uploadSOSJFileConsumerHandler) updateSosjUploadHistories(message *models.UploadHistory, status string) {
	message.UpdatedAt = time.Now()
	message.Status = status
	sosjUploadHistoryJourneysResultChan := make(chan *models.UploadHistoryChan)
	go c.sosjUploadHistoriesRepository.UpdateByID(message.ID.Hex(), message, c.ctx, sosjUploadHistoryJourneysResultChan)
}
