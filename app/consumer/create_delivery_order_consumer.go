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

type CreateDeliveryOrderConsumerHandlerInterface interface {
	ProcessMessage()
}

type createDeliveryOrderConsumerHandler struct {
	kafkaClient          kafkadbo.KafkaClientInterface
	salesOrderUseCase    usecases.SalesOrderUseCaseInterface
	deliveryOrderUseCase usecases.DeliveryOrderUseCaseInterface
	ctx                  context.Context
	args                 []interface{}
	db                   dbresolver.DB
}

func InitCreateDeliveryOrderConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, salesOrderUseCase usecases.SalesOrderUseCaseInterface, deliveryOrderUseCase usecases.DeliveryOrderUseCaseInterface, db dbresolver.DB, ctx context.Context, args []interface{}) CreateDeliveryOrderConsumerHandlerInterface {
	return &createDeliveryOrderConsumerHandler{
		kafkaClient:          kafkaClient,
		salesOrderUseCase:    salesOrderUseCase,
		deliveryOrderUseCase: deliveryOrderUseCase,
		ctx:                  ctx,
		args:                 args,
		db:                   db,
	}
}

func (c *createDeliveryOrderConsumerHandler) ProcessMessage() {
	topic := c.args[1].(string)
	groupID := c.args[2].(string)
	reader := c.kafkaClient.SetConsumerGroupReader(topic, groupID)

	for {
		m, err := reader.ReadMessage(c.ctx)
		if err != nil {
			break
		}

		fmt.Printf("message do at topic/partition/offset %v/%v/%v \n", m.Topic, m.Partition, m.Offset)

		var deliveryOrder models.DeliveryOrder
		err = json.Unmarshal(m.Value, &deliveryOrder)

		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.CREATE_DELIVERY_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
			fmt.Println(errorLogData)
			continue
		}

		dbTransaction, err := c.db.BeginTx(c.ctx, nil)
		errorLog := c.deliveryOrderUseCase.SyncToOpenSearchFromCreateEvent(&deliveryOrder, c.salesOrderUseCase, dbTransaction, c.ctx)

		if errorLog.Err != nil {
			dbTransaction.Rollback()
			errorLogData := helper.WriteLogConsumer(constants.CREATE_DELIVERY_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), errorLog.Err, http.StatusInternalServerError, nil)
			fmt.Println(errorLogData)
			continue
		}

		err = dbTransaction.Commit()
		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.CREATE_DELIVERY_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
			fmt.Println(errorLogData)
			continue
		}

		salesOrderRequest := &models.SalesOrderRequest{
			ID:            deliveryOrder.SalesOrderID,
			OrderSourceID: deliveryOrder.OrderSourceID,
		}

		salesOrderWithDetail, errorLog := c.salesOrderUseCase.GetByID(salesOrderRequest, true, c.ctx)

		if errorLog.Err != nil {
			fmt.Println(errorLog)
			continue
		}

		errorLog = c.salesOrderUseCase.SyncToOpenSearchFromUpdateEvent(salesOrderWithDetail, c.ctx)

		if errorLog.Err != nil {
			fmt.Println(errorLog)
			continue
		}

		deliveryOrderRequest := &models.DeliveryOrderRequest{
			ID: deliveryOrder.ID,
		}

		deliveryOrderWithDetail, errorLog := c.deliveryOrderUseCase.GetByID(deliveryOrderRequest, c.ctx)

		if errorLog.Err != nil {
			fmt.Println(errorLog)
			continue
		}

		salesOrderWithDetail.DeliveryOrders = nil
		deliveryOrderWithDetail.SalesOrder = salesOrderWithDetail
		errorLog = c.deliveryOrderUseCase.SyncToOpenSearchFromUpdateEvent(deliveryOrderWithDetail, c.ctx)

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
