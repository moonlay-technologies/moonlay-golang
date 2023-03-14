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
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
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
	uploadSOSJHistoriesRepository mongoRepositories.UploadSOSJHistoriesRepositoryInterface
	ctx                           context.Context
	args                          []interface{}
	db                            dbresolver.DB
}

func InitUploadSOSJFileConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, uploadRepository repositories.UploadRepositoryInterface, requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, requestValidationRepository repositories.RequestValidationRepositoryInterface, uploadSOSJHistoriesRepository mongoRepositories.UploadSOSJHistoriesRepositoryInterface, db dbresolver.DB, ctx context.Context, args []interface{}) UploadSOFileConsumerHandlerInterface {
	return &uploadSOSJFileConsumerHandler{
		kafkaClient:                   kafkaClient,
		uploadRepository:              uploadRepository,
		requestValidationMiddleware:   requestValidationMiddleware,
		requestValidationRepository:   requestValidationRepository,
		uploadSOSJHistoriesRepository: uploadSOSJHistoriesRepository,
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

		var message models.UploadHistory
		err = json.Unmarshal(m.Value, &message)
		message.CreatedDate = &now
		message.UpdatedAt = &now

		var errors []string

		var idDistributor int
		resultsWithHeader, err := c.uploadRepository.ReadFile("be-so-service", message.FPath, "ap-southeast-1", s3.FileHeaderInfoUse)

		if err != nil {
			fmt.Println("Ini error", err.Error())
			message.UploadStatus = constants.UPLOAD_STATUS_HISTORY_ERR_UPLOAD
			uploadSOSJHistoryJourneysResultChan := make(chan *models.UploadHistoryChan)
			go c.uploadSOSJHistoriesRepository.Insert(&message, c.ctx, uploadSOSJHistoryJourneysResultChan)
			continue
		}

		for _, v := range resultsWithHeader {
			for k2, v2 := range v {
				if v2 == "NoSuratJalan" {
					idDistributor, _ = strconv.Atoi(k2)
				}
			}
		}

		var uploadSOSJFields []*models.UploadSOSJField
		results, err := c.uploadRepository.ReadFile("be-so-service", message.FPath, "ap-southeast-1", s3.FileHeaderInfoIgnore)

		if err != nil {
			fmt.Println("Ini error 2", err.Error())
			message.UploadStatus = constants.UPLOAD_STATUS_HISTORY_ERR_UPLOAD
			uploadSOSJHistoryJourneysResultChan := make(chan *models.UploadHistoryChan)
			go c.uploadSOSJHistoriesRepository.Insert(&message, c.ctx, uploadSOSJHistoryJourneysResultChan)
			continue
		}

		noSuratJalan := []string{}
		for _, v := range results {
			if v["_1"] != "Status" {
				mandatoryError := c.requestValidationMiddleware.UploadMandatoryValidation([]*models.TemplateRequest{
					{
						Field: "Status",
						Value: v["_1"],
					},
					{
						Field: "NoSuratJalan",
						Value: v["_2"],
					},
					{
						Field: "TglSuratJalan",
						Value: v["_3"],
					},
					{
						Field: "KodeTokoDBO",
						Value: v["_4"],
					},
					{
						Field: "IDMerk",
						Value: v["_5"],
					},
					{
						Field: "KodeProdukDBO",
						Value: v["_6"],
					},
					{
						Field: "Qty",
						Value: v["_7"],
					},
					{
						Field: "Unit",
						Value: v["_8"],
					},
				})
				if len(mandatoryError) > 1 {
					errors = append(errors, mandatoryError...)
					continue
				}

				intType := []*models.TemplateRequest{
					{
						Field: "KodeTokoDBO",
						Value: v["_4"],
					},
					{
						Field: "IDMerk",
						Value: v["_5"],
					},
					{
						Field: "KodeProdukDBO",
						Value: v["_6"],
					},
					{
						Field: "Qty",
						Value: v["_7"],
					},
					{
						Field: "Unit",
						Value: v["_8"],
					},
				}
				if v["_11"] != "" {
					intType = append(intType, &models.TemplateRequest{
						Field: "KodeGudang",
						Value: v["_11"],
					})
				}
				if v["_12"] != "" {
					intType = append(intType, &models.TemplateRequest{
						Field: "IDSalesman",
						Value: v["_12"],
					})
				}
				intTypeResult, intTypeError := c.requestValidationMiddleware.UploadIntTypeValidation(intType)
				if len(intTypeError) > 1 {
					errors = append(errors, intTypeError...)
					continue
				}

				if intTypeResult["Qty"] < 1 {
					errors = append(errors, "Quantity harus lebih dari 0")
					continue
				}

				mustActiveField := []*models.MustActiveRequest{
					helper.GenerateMustActive("stores", "KodeTokoDBO", intTypeResult["KodeTokoDBO"], "active"),
					helper.GenerateMustActive("users", "user_id", message.UploadById, "ACTIVE"),
					{
						Table:    "brands",
						ReqField: "IDMerk",
						Clause:   fmt.Sprintf("id = %d AND status_active = %d", intTypeResult["IDMerk"], 1),
						Id:       intTypeResult["IDMerk"],
					},
					{
						Table:    "products",
						ReqField: "KodeProdukDBO",
						Clause:   fmt.Sprintf("id = %d AND isActive = %d", intTypeResult["KodeProdukDBO"], 1),
						Id:       intTypeResult["KodeProdukDBO"],
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
					errors = append(errors, mustActiveError...)
					continue
				}

				if len(v["_12"]) > 0 {
					brandSalesman := make(chan *models.RequestIdValidationChan)
					go c.requestValidationRepository.BrandSalesmanValidation(intTypeResult["IDMerk"], intTypeResult["IDSalesman"], idDistributor, brandSalesman)
					brandSalesmanResult := <-brandSalesman

					if brandSalesmanResult.Total < 1 {
						errors = append(errors, fmt.Sprintf("Kode Merek = %d Tidak Terdaftar pada Distributor <nama_agent>. Silahkan gunakan Kode Merek yang lain.", intTypeResult["IDMerk"]))
						errors = append(errors, fmt.Sprintf("ID Salesman = %d Tidak Terdaftar pada Distributor <nama_agent>. Silahkan gunakan ID Salesman yang lain.", intTypeResult["IDSalesman"]))
						errors = append(errors, fmt.Sprintf("Salesman di Kode Toko = %d untuk Merek <Nama Merk> Tidak Terdaftar. Silahkan gunakan ID Salesman yang terdaftar.", intTypeResult["KodeTokoDBO"]))
						continue
					}
				}
				storeAddresses := make(chan *models.RequestIdValidationChan)
				go c.requestValidationRepository.StoreAddressesValidation(intTypeResult["KodeTokoDBO"], storeAddresses)
				storeAddressesResult := <-storeAddresses

				if storeAddressesResult.Total < 1 {
					errors = append(errors, fmt.Sprintf("Alamat Utama pada Kode Toko = %s Tidak Ditemukan. Silahkan gunakan Alamat Toko yang lain.", v["_4"]))
					continue
				}

				var uploadSOSJField models.UploadSOSJField
				uploadSOSJField.TglSuratJalan, err = helper.ParseDDYYMMtoYYYYMMDD(v["_3"])
				if err != nil {
					errors = append(errors, fmt.Sprintf("Format Tanggal Order = %s Salah, silahkan sesuaikan dengan format DD-MMM-YYYY, contoh 15/12/2021", v["_3"]))
					continue
				}
				uploadSOSJField.UploadSOSJFieldMap(v, idDistributor)

				checkIfNoSuratJalanExist := helper.InSliceString(noSuratJalan, v["_2"])
				if checkIfNoSuratJalanExist {

					for i := range uploadSOSJFields {
						brandId, _ := strconv.Atoi(v["_5"])
						if uploadSOSJFields[i].NoSuratJalan == v["_2"] && uploadSOSJFields[i].IDMerk != brandId {
							uploadSOSJFields[i].NoSuratJalan = uploadSOSJFields[i].NoSuratJalan + "-" + strconv.Itoa(uploadSOSJFields[i].IDMerk)
							uploadSOSJField.NoSuratJalan = v["_2"] + "-" + v["_5"]
							break
						} else {
							uploadSOSJField.NoSuratJalan = v["_2"]
						}
					}

				} else {
					uploadSOSJField.NoSuratJalan = v["_2"]
					noSuratJalan = append(noSuratJalan, v["_2"])
				}

				uploadSOSJFields = append(uploadSOSJFields, &uploadSOSJField)
			}
		}

		keyKafka := []byte(message.RequestId)
		messageKafka, _ := json.Marshal(uploadSOSJFields)

		err = c.kafkaClient.WriteToTopic(constants.UPLOAD_SOSJ_ITEM_TOPIC, keyKafka, messageKafka)

		if err != nil {
			message.UploadStatus = constants.UPLOAD_STATUS_HISTORY_ERR_UPLOAD
			uploadSOSJHistoryJourneysResultChan := make(chan *models.UploadHistoryChan)
			go c.uploadSOSJHistoriesRepository.Insert(&message, c.ctx, uploadSOSJHistoryJourneysResultChan)
			continue
		}

		message.UploadStatus = constants.UPLOAD_STATUS_HISTORY_UPLOADED
		uploadSOSJHistoryJourneysResultChan := make(chan *models.UploadHistoryChan)
		go c.uploadSOSJHistoriesRepository.Insert(&message, c.ctx, uploadSOSJHistoryJourneysResultChan)
	}

	if err := reader.Close(); err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
}
