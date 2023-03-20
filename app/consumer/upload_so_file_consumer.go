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
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
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
	uploadSOHistoriesRepository mongoRepositories.UploadSOHistoriesRepositoryInterface
	soUploadErrorLogsRepository mongoRepositories.SoUploadErrorLogsRepositoryInterface
	ctx                         context.Context
	args                        []interface{}
	db                          dbresolver.DB
}

func InitUploadSOFileConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, uploadRepository repositories.UploadRepositoryInterface, requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, requestValidationRepository repositories.RequestValidationRepositoryInterface, uploadSOHistoriesRepository mongoRepositories.UploadSOHistoriesRepositoryInterface, soUploadErrorLogsRepository mongoRepositories.SoUploadErrorLogsRepositoryInterface, db dbresolver.DB, ctx context.Context, args []interface{}) UploadSOFileConsumerHandlerInterface {
	return &uploadSOFileConsumerHandler{
		kafkaClient:                 kafkaClient,
		uploadRepository:            uploadRepository,
		requestValidationMiddleware: requestValidationMiddleware,
		requestValidationRepository: requestValidationRepository,
		uploadSOHistoriesRepository: uploadSOHistoriesRepository,
		soUploadErrorLogsRepository: soUploadErrorLogsRepository,
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

		var message models.SoUploadHistory
		err = json.Unmarshal(m.Value, &message)
		var key = string(m.Key[:])

		if key == "retry" {
			message.UpdatedAt = now
		} else {
			message.CreatedAt = now
			message.UpdatedAt = now
		}

		results, err := c.uploadRepository.ReadFile("be-so-service", message.FilePath, "ap-southeast-1", s3.FileHeaderInfoUse)

		if err != nil {

			message.Status = constants.UPLOAD_STATUS_HISTORY_ERR_UPLOAD
			uploadSOHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
			if key == "retry" {
				go c.uploadSOHistoriesRepository.UpdateByID(message.ID.Hex(), &message, c.ctx, uploadSOHistoryJourneysResultChan)
			} else {
				go c.uploadSOHistoriesRepository.Insert(&message, c.ctx, uploadSOHistoryJourneysResultChan)
			}
			continue

		}

		message.Status = constants.UPLOAD_STATUS_HISTORY_UPLOADED
		uploadSOHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
		uploadSOHistoryJourneysResult := &models.SoUploadHistoryChan{}
		if key == "retry" {
			go c.uploadSOHistoriesRepository.UpdateByID(message.ID.Hex(), &message, c.ctx, uploadSOHistoryJourneysResultChan)
			uploadSOHistoryJourneysResult = <-uploadSOHistoryJourneysResultChan
		} else {
			go c.uploadSOHistoriesRepository.Insert(&message, c.ctx, uploadSOHistoryJourneysResultChan)
			uploadSOHistoryJourneysResult = <-uploadSOHistoryJourneysResultChan
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
			if len(mandatoryError) > 1 {

				if key == "retry" {
					message.Status = constants.UPLOAD_STATUS_HISTORY_ERR_UPLOAD
					uploadSOHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
					go c.uploadSOHistoriesRepository.UpdateByID(message.ID.Hex(), &message, c.ctx, uploadSOHistoryJourneysResultChan)
					break
				} else {
					errors := mandatoryError

					soUploadErrorLogs := &models.SoUploadErrorLog{}
					soUploadErrorLogs.SoUploadErrorLogsMap(v, message.AgentName, i+2, errors, &now)
					soUploadErrorLogs.SoUploadHistoryId = uploadSOHistoryJourneysResult.SoUploadHistory.ID

					soUploadErrorLogssResultChan := make(chan *models.SoUploadErrorLogChan)
					go c.soUploadErrorLogsRepository.Insert(soUploadErrorLogs, c.ctx, soUploadErrorLogssResultChan)
					continue
				}
			}

			if brandIds[v["NoOrder"]] != nil {

				var isBreak bool

				if brandIds[v["NoOrder"]][0]["KodeMerk"] != v["KodeMerk"] {
					message.Status = constants.UPLOAD_STATUS_HISTORY_ERR_UPLOAD
					uploadSOHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
					if key == "retry" {
						go c.uploadSOHistoriesRepository.UpdateByID(message.ID.Hex(), &message, c.ctx, uploadSOHistoryJourneysResultChan)
					} else {
						go c.uploadSOHistoriesRepository.Insert(&message, c.ctx, uploadSOHistoryJourneysResultChan)
					}
					isBreak = true
					break
				}

				for _, x := range brandIds[v["NoOrder"]] {
					if x["KodeProduk"] == v["KodeProduk"] && x["UnitProduk"] == v["UnitProduk"] {
						if key == "retry" {
							message.Status = constants.UPLOAD_STATUS_HISTORY_ERR_UPLOAD
							uploadSOHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
							go c.uploadSOHistoriesRepository.UpdateByID(message.ID.Hex(), &message, c.ctx, uploadSOHistoryJourneysResultChan)

							isBreak = true
							break
						} else {
							var errors []string
							errors = append(errors, fmt.Sprintf("Duplikat row untuk No Order %s", v["NoOrder"]))

							soUploadErrorLogs := &models.SoUploadErrorLog{}
							soUploadErrorLogs.SoUploadErrorLogsMap(v, message.AgentName, i+2, errors, &now)
							soUploadErrorLogs.SoUploadHistoryId = uploadSOHistoryJourneysResult.SoUploadHistory.ID

							soUploadErrorLogssResultChan := make(chan *models.SoUploadErrorLogChan)
							go c.soUploadErrorLogsRepository.Insert(soUploadErrorLogs, c.ctx, soUploadErrorLogssResultChan)

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
					Field: "KodeToko",
					Value: v["KodeToko"],
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
			if len(intTypeError) > 1 {

				if key == "retry" {
					message.Status = constants.UPLOAD_STATUS_HISTORY_ERR_UPLOAD
					uploadSOHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
					go c.uploadSOHistoriesRepository.UpdateByID(message.ID.Hex(), &message, c.ctx, uploadSOHistoryJourneysResultChan)

					break
				} else {
					errors := intTypeError

					soUploadErrorLogs := &models.SoUploadErrorLog{}
					soUploadErrorLogs.SoUploadErrorLogsMap(v, message.AgentName, i+2, errors, &now)
					soUploadErrorLogs.SoUploadHistoryId = uploadSOHistoryJourneysResult.SoUploadHistory.ID

					soUploadErrorLogssResultChan := make(chan *models.SoUploadErrorLogChan)
					go c.soUploadErrorLogsRepository.Insert(soUploadErrorLogs, c.ctx, soUploadErrorLogssResultChan)

					continue
				}
			}

			if intTypeResult["QTYOrder"] < 1 {
				if key == "retry" {
					message.Status = constants.UPLOAD_STATUS_HISTORY_ERR_UPLOAD
					uploadSOHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
					go c.uploadSOHistoriesRepository.UpdateByID(message.ID.Hex(), &message, c.ctx, uploadSOHistoryJourneysResultChan)

					break
				} else {
					errors := []string{"Quantity harus lebih dari 0"}

					soUploadErrorLogs := &models.SoUploadErrorLog{}
					soUploadErrorLogs.SoUploadErrorLogsMap(v, message.AgentName, i+2, errors, &now)
					soUploadErrorLogs.SoUploadHistoryId = uploadSOHistoryJourneysResult.SoUploadHistory.ID

					soUploadErrorLogssResultChan := make(chan *models.SoUploadErrorLogChan)
					go c.soUploadErrorLogsRepository.Insert(soUploadErrorLogs, c.ctx, soUploadErrorLogssResultChan)

					continue
				}
			}

			mustActiveField := []*models.MustActiveRequest{
				helper.GenerateMustActive("agents", "IDDistributor", intTypeResult["IDDistributor"], "active"),
				helper.GenerateMustActive("users", "user_id", int(message.UploadedBy), "ACTIVE"),
				{
					Table:    "salesmans",
					ReqField: "IDSalesman",
					Clause:   fmt.Sprintf("id = %d AND deleted_at IS NULL", intTypeResult["IDSalesman"]),
					Id:       intTypeResult["IDSalesman"],
				},
				{
					Table:         "brands",
					ReqField:      "KodeMerk",
					Clause:        fmt.Sprintf("id = %d AND status_active = %d", intTypeResult["KodeMerk"], 1),
					Id:            intTypeResult["KodeMerk"],
					CustomMessage: fmt.Sprintf("Kode Merek = %d Tidak Terdaftar pada Distributor <nama_agent>. Silahkan gunakan Kode Merek yang lain.", intTypeResult["KodeMerk"]),
				},
				{
					Table:         "products",
					ReqField:      "KodeProduk",
					Clause:        fmt.Sprintf("sku = '%s' AND isActive = %d", v["KodeProduk"], 1),
					Id:            v["KodeProduk"],
					CustomMessage: fmt.Sprintf("Kode SKU = %s dengan Merek %ssudah Tidak Aktif. Silahkan gunakan Kode SKU yang lain.", v["KodeProduk"], v["NamaMerk"]),
				},
				{
					Table:    "uoms",
					ReqField: "UnitProduk",
					Clause:   fmt.Sprintf("code = '%s' AND deleted_at IS NULL", v["UnitProduk"]),
					Id:       v["UnitProduk"],
				},
			}
			mustActiveError := c.requestValidationMiddleware.UploadMustActiveValidation(mustActiveField)
			if len(mustActiveError) > 1 {
				if key == "retry" {
					message.Status = constants.UPLOAD_STATUS_HISTORY_ERR_UPLOAD
					uploadSOHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
					go c.uploadSOHistoriesRepository.UpdateByID(message.ID.Hex(), &message, c.ctx, uploadSOHistoryJourneysResultChan)

					break
				} else {
					errors := mustActiveError

					soUploadErrorLogs := &models.SoUploadErrorLog{}
					soUploadErrorLogs.SoUploadErrorLogsMap(v, message.AgentName, i+2, errors, &now)
					soUploadErrorLogs.SoUploadHistoryId = uploadSOHistoryJourneysResult.SoUploadHistory.ID

					soUploadErrorLogssResultChan := make(chan *models.SoUploadErrorLogChan)
					go c.soUploadErrorLogsRepository.Insert(soUploadErrorLogs, c.ctx, soUploadErrorLogssResultChan)

					continue
				}
			}

			brandSalesmanResultChan := make(chan *models.RequestIdValidationChan)
			go c.requestValidationRepository.BrandSalesmanValidation(intTypeResult["KodeMerk"], intTypeResult["IDSalesman"], intTypeResult["IDDistributor"], brandSalesmanResultChan)
			brandSalesmanResult := <-brandSalesmanResultChan

			if brandSalesmanResult.Total < 1 {
				if key == "retry" {
					message.Status = constants.UPLOAD_STATUS_HISTORY_ERR_UPLOAD
					uploadSOHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
					go c.uploadSOHistoriesRepository.UpdateByID(message.ID.Hex(), &message, c.ctx, uploadSOHistoryJourneysResultChan)

					break
				} else {
					var errors []string

					errors = append(errors, fmt.Sprintf("Kode Merek = %d Tidak Terdaftar pada Distributor <nama_agent>. Silahkan gunakan Kode Merek yang lain.", intTypeResult["KodeMerk"]))
					errors = append(errors, fmt.Sprintf("ID Salesman = %d Tidak Terdaftar pada Distributor <nama_agent>. Silahkan gunakan ID Salesman yang lain.", intTypeResult["IDSalesman"]))

					soUploadErrorLogs := &models.SoUploadErrorLog{}
					soUploadErrorLogs.SoUploadErrorLogsMap(v, message.AgentName, i+2, errors, &now)

					soUploadErrorLogssResultChan := make(chan *models.SoUploadErrorLogChan)
					go c.soUploadErrorLogsRepository.Insert(soUploadErrorLogs, c.ctx, soUploadErrorLogssResultChan)

					continue
				}
			}

			storeAddressesResultChan := make(chan *models.RequestIdValidationChan)
			go c.requestValidationRepository.StoreAddressesValidation(intTypeResult["KodeToko"], storeAddressesResultChan)
			storeAddressesResult := <-storeAddressesResultChan

			if storeAddressesResult.Total < 1 {
				if key == "retry" {
					message.Status = constants.UPLOAD_STATUS_HISTORY_ERR_UPLOAD
					uploadSOHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
					go c.uploadSOHistoriesRepository.UpdateByID(message.ID.Hex(), &message, c.ctx, uploadSOHistoryJourneysResultChan)

					break
				} else {
					errors := []string{fmt.Sprintf("Alamat Utama pada Kode Toko = %s Tidak Ditemukan. Silahkan gunakan Alamat Toko yang lain.", v["KodeToko"])}

					soUploadErrorLogs := &models.SoUploadErrorLog{}
					soUploadErrorLogs.SoUploadErrorLogsMap(v, message.AgentName, i+2, errors, &now)
					soUploadErrorLogs.SoUploadHistoryId = uploadSOHistoryJourneysResult.SoUploadHistory.ID

					soUploadErrorLogssResultChan := make(chan *models.SoUploadErrorLogChan)
					go c.soUploadErrorLogsRepository.Insert(soUploadErrorLogs, c.ctx, soUploadErrorLogssResultChan)

					continue
				}
			}

			var uploadSOField models.UploadSOField
			uploadSOField.TanggalOrder, err = helper.ParseDDYYMMtoYYYYMMDD(v["TanggalOrder"])
			if err != nil {
				if key == "retry" {
					message.Status = constants.UPLOAD_STATUS_HISTORY_ERR_UPLOAD
					uploadSOHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
					go c.uploadSOHistoriesRepository.UpdateByID(message.ID.Hex(), &message, c.ctx, uploadSOHistoryJourneysResultChan)

					break
				} else {
					errors := []string{fmt.Sprintf("Format Tanggal Order = %s Salah, silahkan sesuaikan dengan format DD-MMM-YYYY, contoh 15/12/2021", v["TanggalOrder"])}

					soUploadErrorLogs := &models.SoUploadErrorLog{}
					soUploadErrorLogs.SoUploadErrorLogsMap(v, message.AgentName, i+2, errors, &now)
					soUploadErrorLogs.SoUploadHistoryId = uploadSOHistoryJourneysResult.SoUploadHistory.ID

					soUploadErrorLogssResultChan := make(chan *models.SoUploadErrorLogChan)
					go c.soUploadErrorLogsRepository.Insert(soUploadErrorLogs, c.ctx, soUploadErrorLogssResultChan)

					continue
				}
			}

			uploadSOField.TanggalTokoOrder, err = helper.ParseDDYYMMtoYYYYMMDD(v["TanggalTokoOrder"])
			if err != nil {
				if key == "retry" {
					message.Status = constants.UPLOAD_STATUS_HISTORY_ERR_UPLOAD
					uploadSOHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
					go c.uploadSOHistoriesRepository.UpdateByID(message.ID.Hex(), &message, c.ctx, uploadSOHistoryJourneysResultChan)

					break
				} else {
					errors := []string{fmt.Sprintf("Format Tanggal Toko Order = %s Salah, silahkan sesuaikan dengan format DD-MMM-YYYY, contoh 15/12/2021", v["TanggalTokoOrder"])}

					soUploadErrorLogs := &models.SoUploadErrorLog{}
					soUploadErrorLogs.SoUploadErrorLogsMap(v, message.AgentName, i+2, errors, &now)
					soUploadErrorLogs.SoUploadHistoryId = uploadSOHistoryJourneysResult.SoUploadHistory.ID

					soUploadErrorLogssResultChan := make(chan *models.SoUploadErrorLogChan)
					go c.soUploadErrorLogsRepository.Insert(soUploadErrorLogs, c.ctx, soUploadErrorLogssResultChan)

					continue
				}
			}

			if key == "retry" {
				uploadSOField.UploadType = "retry"
			} else {
				uploadSOField.UploadType = "upload"
			}

			uploadSOField.UploadSOFieldMap(v, int(message.UploadedBy), uploadSOHistoryJourneysResult.SoUploadHistory.ID.Hex())

			uploadSOFields = append(uploadSOFields, &uploadSOField)
		}

		keyKafka := []byte(message.RequestId)
		messageKafka, _ := json.Marshal(uploadSOFields)

		err = c.kafkaClient.WriteToTopic(constants.UPLOAD_SO_ITEM_TOPIC, keyKafka, messageKafka)

		if err != nil {
			message.Status = constants.UPLOAD_STATUS_HISTORY_UPLOADED
			uploadSOHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
			go c.uploadSOHistoriesRepository.Insert(&message, c.ctx, uploadSOHistoryJourneysResultChan)
			continue
		}

	}

	if err := reader.Close(); err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
}
