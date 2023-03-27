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

type UploadDOFileConsumerHandlerInterface interface {
	ProcessMessage()
}

type uploadDOFileConsumerHandler struct {
	kafkaClient                 kafkadbo.KafkaClientInterface
	uploadRepository            repositories.UploadRepositoryInterface
	requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface
	requestValidationRepository repositories.RequestValidationRepositoryInterface
	uploadSJHistoriesRepository mongoRepositories.UploadSJHistoriesRepositoryInterface
	ctx                         context.Context
	args                        []interface{}
	db                          dbresolver.DB
}

func InitUploadDOFileConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, uploadRepository repositories.UploadRepositoryInterface, requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, requestValidationRepository repositories.RequestValidationRepositoryInterface, uploadSJHistoriesRepository mongoRepositories.UploadSJHistoriesRepositoryInterface, ctx context.Context, args []interface{}, db dbresolver.DB) UploadDOFileConsumerHandlerInterface {
	return &uploadDOFileConsumerHandler{
		kafkaClient:                 kafkaClient,
		uploadRepository:            uploadRepository,
		requestValidationMiddleware: requestValidationMiddleware,
		requestValidationRepository: requestValidationRepository,
		uploadSJHistoriesRepository: uploadSJHistoriesRepository,
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

		var message models.UploadSJHistory
		err = json.Unmarshal(m.Value, &message)
		message.CreatedAt = now
		message.UpdatedAt = now

		var errors []string

		results, err := c.uploadRepository.ReadFile("be-so-service", message.FilePath, "ap-southeast-1", s3.FileHeaderInfoUse)

		if err != nil {
			fmt.Println(err.Error())
			message.Status = "Error"
			uploadSJHistoryJourneyResultChan := make(chan *models.UploadSJHistoryChan)
			go c.uploadSJHistoriesRepository.Insert(&message, c.ctx, uploadSJHistoryJourneyResultChan)
			continue
		}

		var uploadDOFields []*models.UploadDOField
		for _, v := range results {

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
				errors = append(errors, mandatoryError...)
				continue
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
				// {
				// 	Field: "KodeProduk",
				// 	Value: v["KodeProduk"],
				// },
				{
					Field: "QTYShip",
					Value: v["QTYShip"],
				},
			}
			intTypeResult, intTypeError := c.requestValidationMiddleware.UploadIntTypeValidation(intType)
			if len(intTypeError) > 1 {
				errors = append(errors, intTypeError...)
				continue
			}

			if intTypeResult["QTYShip"] < 1 {
				errors = append(errors, "Quantity harus lebih dari 0")
				continue
			}

			mustActiveField := []*models.MustActiveRequest{
				helper.GenerateMustActive("agents", "IDDistributor", intTypeResult["IDDistributor"], "active"),
			}
			mustActiveError := c.requestValidationMiddleware.UploadMustActiveValidation(mustActiveField)
			if len(mustActiveError) > 1 {
				errors = append(errors, mustActiveError...)
				continue
			}

			var uploadDOField models.UploadDOField
			uploadDOField.TanggalSJ, err = helper.ParseDDYYMMtoYYYYMMDD(v["TanggalSJ"])
			if err != nil {
				errors = append(errors, fmt.Sprintf("Format Tanggal Order = %s Salah, silahkan sesuaikan dengan format DD-MMM-YYYY, contoh 15/12/2021", v["TanggalSJ"]))
				continue
			}
			uploadDOField.UploadDOFieldMap(v, int(message.UploadedBy))

			uploadDOFields = append(uploadDOFields, &uploadDOField)
		}

		keyKafka := []byte(message.RequestId)
		messageKafka, _ := json.Marshal(uploadDOFields)

		err = c.kafkaClient.WriteToTopic(constants.UPLOAD_DO_ITEM_TOPIC, keyKafka, messageKafka)

		if err != nil {
			message.Status = "Error"
			uploadSJHistoryJourneyResultChan := make(chan *models.UploadSJHistoryChan)
			go c.uploadSJHistoriesRepository.Insert(&message, c.ctx, uploadSJHistoryJourneyResultChan)
			continue
		}

		message.Status = "Uploaded"
		uploadSJHistoryJourneyResultChan := make(chan *models.UploadSJHistoryChan)
		go c.uploadSJHistoriesRepository.Insert(&message, c.ctx, uploadSJHistoryJourneyResultChan)
	}

	if err := reader.Close(); err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
}
