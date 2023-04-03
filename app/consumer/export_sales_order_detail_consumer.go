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

type ExportSalesOrderDetailConsumerHandlerInterface interface {
	ProcessMessage()
}

type exportSalesOrderDetailConsumerHandler struct {
	kafkaClient               kafkadbo.KafkaClientInterface
	SalesOrderConsumerUseCase usecases.SalesOrderOpenSearchUseCaseInterface
	ctx                       context.Context
	args                      []interface{}
	db                        dbresolver.DB
}

func InitExportSalesOrderDetailConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, SalesOrderConsumerUseCase usecases.SalesOrderOpenSearchUseCaseInterface, db dbresolver.DB, ctx context.Context, args []interface{}) ExportSalesOrderDetailConsumerHandlerInterface {
	return &exportSalesOrderDetailConsumerHandler{
		kafkaClient:               kafkaClient,
		SalesOrderConsumerUseCase: SalesOrderConsumerUseCase,
		ctx:                       ctx,
		args:                      args,
		db:                        db,
	}
}

func (c *exportSalesOrderDetailConsumerHandler) ProcessMessage() {
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

		var salesOrderDetail models.SalesOrderDetailExportRequest
		err = json.Unmarshal(m.Value, &salesOrderDetail)

		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.CREATE_SALES_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
			fmt.Println(errorLogData)
			continue
		}

		errorLog := c.SalesOrderConsumerUseCase.GetDetails(&salesOrderDetail)

		if errorLog.Err != nil {
			errorLogData := helper.WriteLogConsumer(constants.CREATE_SALES_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), errorLog.Err, http.StatusInternalServerError, nil)
			fmt.Println(errorLogData)
			continue
		}
	}

	if err := reader.Close(); err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
}
