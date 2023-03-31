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
	"strings"
	"time"

	"github.com/bxcodec/dbresolver"
)

type UploadDOFileConsumerHandlerInterface interface {
	ProcessMessage()
}

type uploadDOFileConsumerHandler struct {
	kafkaClient                 kafkadbo.KafkaClientInterface
	uploadRepository            repositories.UploadRepositoryInterface
	requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface
	requestValidationRepository repositories.RequestValidationRepositoryInterface
	sjUploadHistoriesRepository mongoRepositories.DoUploadHistoriesRepositoryInterface
	sjUploadErrorLogsRepository mongoRepositories.DoUploadErrorLogsRepositoryInterface
	warehouseRepository         repositories.WarehouseRepositoryInterface
	ctx                         context.Context
	args                        []interface{}
	db                          dbresolver.DB
}

func InitUploadDOFileConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, uploadRepository repositories.UploadRepositoryInterface, requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, requestValidationRepository repositories.RequestValidationRepositoryInterface, sjUploadHistoriesRepository mongoRepositories.DoUploadHistoriesRepositoryInterface, sjUploadErrorLogsRepository mongoRepositories.DoUploadErrorLogsRepositoryInterface, warehouseRepository repositories.WarehouseRepositoryInterface, ctx context.Context, args []interface{}, db dbresolver.DB) UploadDOFileConsumerHandlerInterface {
	return &uploadDOFileConsumerHandler{
		kafkaClient:                 kafkaClient,
		uploadRepository:            uploadRepository,
		requestValidationMiddleware: requestValidationMiddleware,
		requestValidationRepository: requestValidationRepository,
		sjUploadHistoriesRepository: sjUploadHistoriesRepository,
		sjUploadErrorLogsRepository: sjUploadErrorLogsRepository,
		warehouseRepository:         warehouseRepository,
		ctx:                         ctx,
		args:                        args,
		db:                          db,
	}
}

func (c *uploadDOFileConsumerHandler) ProcessMessage() {
	fmt.Println("process ", constants.UPLOAD_DO_FILE_TOPIC)
	topic := c.args[1].(string)
	groupID := c.args[2].(string)
	reader := c.kafkaClient.SetConsumerGroupReader(topic, groupID)

	for {
		m, err := reader.ReadMessage(c.ctx)
		if err != nil {
			break
		}

		fmt.Printf("message do at topic/partition/offset %v/%v/%v \n", m.Topic, m.Partition, m.Offset)
		now := time.Now()

		sjUploadHistoryId := m.Value
		var key = string(m.Key[:])

		// var errors []string

		sjUploadHistoryJourneysResultChan := make(chan *models.DoUploadHistoryChan)
		go c.sjUploadHistoriesRepository.GetByID(string(sjUploadHistoryId), false, c.ctx, sjUploadHistoryJourneysResultChan)
		sjUploadHistoryJourneysResult := <-sjUploadHistoryJourneysResultChan
		if sjUploadHistoryJourneysResult.Error != nil {
			fmt.Println(sjUploadHistoryJourneysResult.Error.Error())
		}

		message := sjUploadHistoryJourneysResult.DoUploadHistory
		file, err := c.uploadRepository.ReadFile(message.FilePath)

		if err != nil {
			c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
			continue
		}

		parseFile := string(file)
		data := strings.Split(parseFile, "\n")
		totalRows := int64(len(data) - 1)

		message.TotalRows = &totalRows
		c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_SUCCESS)

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
			uploadSjField := map[string]string{}
			for j, y := range line {
				uploadSjField[strings.ReplaceAll(headers[j], "\r", "")] = strings.ReplaceAll(y, "\r", "")
			}
			results = append(results, uploadSjField)
		}

		var uploadDOFields []*models.UploadDOField
		for i, v := range results {
			warehouseName := ""
			if len(v["KodeGudang"]) > 0 {
				getWarehouseResultChan := make(chan *models.WarehouseChan)
				go c.warehouseRepository.GetByCode(v["KodeGudang"], false, c.ctx, getWarehouseResultChan)
				getWarehouseResult := <-getWarehouseResultChan
				if getWarehouseResult.Error != nil {
					if key == "retry" {
						c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
						break
					} else {
						errors := []string{getWarehouseResult.Error.Error()}
						c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
						continue
					}
				}
				warehouseName = getWarehouseResult.Warehouse.Name
			}
			mandatoryError := c.requestValidationMiddleware.UploadMandatoryValidation([]*models.TemplateRequest{
				{
					Field: "IDDistributor",
					Value: v["IDDistributor"],
				},
				{
					Field: "NoOrder",
					Value: v["NoOrder"],
				},
				{
					Field: "TanggalSJ",
					Value: v["TanggalSJ"],
				},
				{
					Field: "NoSJ",
					Value: v["NoSJ"],
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
					Field: "QTYShip",
					Value: v["QTYShip"],
				},
				{
					Field: "Unit",
					Value: v["Unit"],
				},
			})
			if len(mandatoryError) > 1 {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := mandatoryError
					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			intType := []*models.TemplateRequest{
				{
					Field: "IDDistributor",
					Value: v["IDDistributor"],
				},
				{
					Field: "KodeMerk",
					Value: v["KodeMerk"],
				},
				{
					Field: "QTYShip",
					Value: v["QTYShip"],
				},
			}
			intTypeResult, intTypeError := c.requestValidationMiddleware.UploadIntTypeValidation(intType)
			if len(intTypeError) > 1 {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := intTypeError
					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			if intTypeResult["QTYShip"] < 1 {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{"Quantity harus lebih dari 0"}

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			mustActiveField := []*models.MustActiveRequest{
				helper.GenerateMustActive("agents", "IDDistributor", intTypeResult["IDDistributor"], "active"),
			}
			mustActiveError := c.requestValidationMiddleware.UploadMustActiveValidation(mustActiveField)
			if len(mustActiveError) > 1 {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := mustActiveError

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			tanggalSJ, err := helper.ParseDDYYMMtoYYYYMMDD(v["TanggalSJ"])
			if err != nil {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{fmt.Sprintf("Format Tanggal Order = %s Salah, silahkan sesuaikan dengan format DD-MMM-YYYY, contoh 15/12/2021", v["TanggalSJ"])}

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			var uploadDOField models.UploadDOField
			uploadDOField.TanggalSJ = tanggalSJ
			uploadDOField.UploadDOFieldMap(v, int(*message.UploadedBy), message.ID.Hex())
			uploadDOField.BulkCode = message.BulkCode
			uploadDOField.ErrorLine = i + 2
			uploadDOField.UploadType = key

			uploadDOFields = append(uploadDOFields, &uploadDOField)
		}

		keyKafka := []byte(message.RequestId)
		messageKafka, _ := json.Marshal(uploadDOFields)

		err = c.kafkaClient.WriteToTopic(constants.UPLOAD_DO_ITEM_TOPIC, keyKafka, messageKafka)

		if err != nil {
			c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
			continue
		}

	}

	if err := reader.Close(); err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
}

func (c *uploadDOFileConsumerHandler) createSjUploadErrorLog(errorLine int, agentId, sjUploadHistoryId, requestId, agentName, bulkCode, warehouseName string, errors []string, now *time.Time, item map[string]string) {
	sjUploadErrorLog := &models.DoUploadErrorLog{}
	sjUploadErrorLog.DoUploadErrorLogsMap(errorLine, sjUploadHistoryId, requestId, bulkCode, errors, now)
	rowData := &models.RowDataDoUploadErrorLog{}
	rowData.RowDataDoUploadErrorLogMap(item, agentName, warehouseName)
	sjUploadErrorLog.RowData = *rowData

	sjUploadErrorLogResultChan := make(chan *models.DoUploadErrorLogChan)
	go c.sjUploadErrorLogsRepository.Insert(sjUploadErrorLog, c.ctx, sjUploadErrorLogResultChan)
}

func (c *uploadDOFileConsumerHandler) updateSjUploadHistories(message *models.DoUploadHistory, status string) {
	message.UpdatedAt = time.Now()
	message.Status = status
	sjUploadHistoryJourneysResultChan := make(chan *models.DoUploadHistoryChan)
	go c.sjUploadHistoriesRepository.UpdateByID(message.ID.Hex(), message, c.ctx, sjUploadHistoryJourneysResultChan)
}
