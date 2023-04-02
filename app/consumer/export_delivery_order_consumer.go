package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/app/usecases"
	"order-service/global/utils/helper"
	kafkadbo "order-service/global/utils/kafka"

	"github.com/bxcodec/dbresolver"
)

type ExportDeliveryOrderConsumerHandlerInterface interface {
	ProcessMessage()
}

type exportDeliveryOrderConsumerHandler struct {
	kafkaClient                  kafkadbo.KafkaClientInterface
	DeliveryOrderConsumerUseCase usecases.DeliveryOrderConsumerUseCaseInterface
	ctx                          context.Context
	args                         []interface{}
	db                           dbresolver.DB
}

func InitExportDeliveryOrderConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, DeliveryOrderConsumerUseCase usecases.DeliveryOrderConsumerUseCaseInterface, db dbresolver.DB, ctx context.Context, args []interface{}) ExportDeliveryOrderConsumerHandlerInterface {
	return &exportDeliveryOrderConsumerHandler{
		kafkaClient:                  kafkaClient,
		DeliveryOrderConsumerUseCase: DeliveryOrderConsumerUseCase,
		ctx:                          ctx,
		args:                         args,
		db:                           db,
	}
}

func (c *exportDeliveryOrderConsumerHandler) ProcessMessage() {
	fmt.Println("process")
	topic := c.args[1].(string)
	groupID := c.args[2].(string)
	reader := c.kafkaClient.SetConsumerGroupReader(topic, groupID)

	for {
		m, err := reader.ReadMessage(c.ctx)
		if err != nil {
			break
		}

		fmt.Printf("message do at topic/partition/offset %v/%v/%v \n", m.Topic, m.Partition, m.Offset)

		var deliveryOrder models.DeliveryOrderExportRequest
		err = json.Unmarshal(m.Value, &deliveryOrder)

		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.CREATE_DELIVERY_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
			fmt.Println(errorLogData)
			continue
		}

		errorLog := c.DeliveryOrderConsumerUseCase.Get(&deliveryOrder)

		if errorLog.Err != nil {
			errorLogData := helper.WriteLogConsumer(constants.CREATE_DELIVERY_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), errorLog.Err, http.StatusInternalServerError, nil)
			fmt.Println(errorLogData)
			continue
		}
	}

	if err := reader.Close(); err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
}
