package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	mongoRepositories "order-service/app/repositories/mongod"
	"order-service/app/usecases"
	"order-service/global/utils/helper"
	kafkadbo "order-service/global/utils/kafka"
	"time"

	"github.com/bxcodec/dbresolver"
)

type UpdateDeliveryOrderConsumerHandlerInterface interface {
	ProcessMessage()
}

type UpdateDeliveryOrderConsumerHandler struct {
	kafkaClient                    kafkadbo.KafkaClientInterface
	DeliveryOrderOpenSearchUseCase usecases.DeliveryOrderOpenSearchUseCaseInterface
	ctx                            context.Context
	args                           []interface{}
	db                             dbresolver.DB
	deliveryOrderLogRepository     mongoRepositories.DeliveryOrderLogRepositoryInterface
}

func InitUpdateDeliveryOrderConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, deliveryOrderLogRepository mongoRepositories.DeliveryOrderLogRepositoryInterface, DeliveryOrderOpenSearchUseCase usecases.DeliveryOrderOpenSearchUseCaseInterface, db dbresolver.DB, ctx context.Context, args []interface{}) UpdateDeliveryOrderConsumerHandlerInterface {
	return &UpdateDeliveryOrderConsumerHandler{
		kafkaClient:                    kafkaClient,
		DeliveryOrderOpenSearchUseCase: DeliveryOrderOpenSearchUseCase,
		ctx:                            ctx,
		args:                           args,
		db:                             db,
		deliveryOrderLogRepository:     deliveryOrderLogRepository,
	}
}

func (c *UpdateDeliveryOrderConsumerHandler) ProcessMessage() {
	fmt.Println("process ", constants.UPDATE_DELIVERY_ORDER_TOPIC)
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
		now := time.Now()
		deliveryOrderLogResultChan := make(chan *models.DeliveryOrderLogChan)

		dbTransaction, err := c.db.BeginTx(c.ctx, nil)
		deliveryOrderLog := &models.DeliveryOrderLog{
			RequestID: "",
			DoCode:    "",
			Data:      m.Value,
			Status:    constants.LOG_STATUS_MONGO_ERROR,
			Action:    constants.LOG_ACTION_MONGO_UPDATE,
			CreatedAt: &now,
		}

		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.UPDATE_DELIVERY_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
			go c.deliveryOrderLogRepository.Insert(deliveryOrderLog, c.ctx, deliveryOrderLogResultChan)
			fmt.Println(errorLogData)
			continue
		}
		go c.deliveryOrderLogRepository.GetByCode(deliveryOrder.DoCode, constants.LOG_STATUS_MONGO_DEFAULT, deliveryOrderLog.Action, false, c.ctx, deliveryOrderLogResultChan)
		deliveryOrderDetailResult := <-deliveryOrderLogResultChan
		if deliveryOrderDetailResult.Error != nil {
			go c.deliveryOrderLogRepository.Insert(deliveryOrderLog, c.ctx, deliveryOrderLogResultChan)
			fmt.Println(deliveryOrderDetailResult.Error)
			continue
		}
		deliveryOrderLog = deliveryOrderDetailResult.DeliveryOrderLog
		deliveryOrderLog.Status = constants.LOG_STATUS_MONGO_ERROR
		deliveryOrderLog.UpdatedAt = &now
		errorLog := c.DeliveryOrderOpenSearchUseCase.SyncToOpenSearchFromUpdateEvent(&deliveryOrder, c.ctx)

		if errorLog.Err != nil {
			dbTransaction.Rollback()
			errorLogData := helper.WriteLogConsumer(constants.UPDATE_DELIVERY_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), errorLog.Err, http.StatusInternalServerError, nil)
			go c.deliveryOrderLogRepository.UpdateByID(deliveryOrderLog.ID.Hex(), deliveryOrderLog, c.ctx, deliveryOrderLogResultChan)
			fmt.Println(errorLogData)
			continue
		}
		err = dbTransaction.Commit()
		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.UPDATE_DELIVERY_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
			go c.deliveryOrderLogRepository.UpdateByID(deliveryOrderLog.ID.Hex(), deliveryOrderLog, c.ctx, deliveryOrderLogResultChan)
			fmt.Println(errorLogData)
			continue
		}

		deliveryOrderLog.Status = constants.LOG_STATUS_MONGO_SUCCESS
		go c.deliveryOrderLogRepository.UpdateByID(deliveryOrderLog.ID.Hex(), deliveryOrderLog, c.ctx, deliveryOrderLogResultChan)
	}

	if err := reader.Close(); err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
}
