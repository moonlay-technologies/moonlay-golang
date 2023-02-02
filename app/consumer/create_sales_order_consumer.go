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

type CreateSalesOrderConsumerHandlerInterface interface {
	ProcessMessage()
}

type createSalesOrderConsumerHandler struct {
	kafkaClient       kafkadbo.KafkaClientInterface
	salesOrderUseCase usecases.SalesOrderUseCaseInterface
	ctx               context.Context
	args              []interface{}
	db                dbresolver.DB
}

func InitCreateSalesOrderConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, salesOrderUseCase usecases.SalesOrderUseCaseInterface, db dbresolver.DB, ctx context.Context, args []interface{}) CreateSalesOrderConsumerHandlerInterface {
	return &createSalesOrderConsumerHandler{
		kafkaClient:       kafkaClient,
		salesOrderUseCase: salesOrderUseCase,
		ctx:               ctx,
		args:              args,
		db:                db,
	}
}

func (c *createSalesOrderConsumerHandler) ProcessMessage() {
	fmt.Println("process")
	topic := c.args[1].(string)
	groupID := c.args[2].(string)
	reader := c.kafkaClient.SetConsumerGroupReader(topic, groupID)

	for {
		m, err := reader.ReadMessage(c.ctx)
		if err != nil {
			break
		}

		fmt.Printf("message at topic/partition/offset %v/%v/%v \n", m.Topic, m.Partition, m.Offset)

		var salesOrder models.SalesOrder
		err = json.Unmarshal(m.Value, &salesOrder)

		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.SALES_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, "Ada kesalahan, silahkan coba lagi nanti")
			fmt.Println(errorLogData)
			continue
		}

		dbTransaction, err := c.db.BeginTx(c.ctx, nil)
		errorLog := c.salesOrderUseCase.SyncToOpenSearchFromCreateEvent(&salesOrder, dbTransaction, c.ctx)

		if errorLog.Err != nil {
			dbTransaction.Rollback()
			errorLogData := helper.WriteLogConsumer(constants.SALES_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), errorLog.Err, http.StatusInternalServerError, "Ada kesalahan, silahkan coba lagi nanti")
			fmt.Println(errorLogData)
			continue
		}

		err = dbTransaction.Commit()
		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.SALES_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, "Ada kesalahan, silahkan coba lagi nanti")
			fmt.Println(errorLogData)
			continue
		}
	}

	if err := reader.Close(); err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
}
