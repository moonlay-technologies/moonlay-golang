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
	"order-service/global/utils/model"
	"time"

	"github.com/bxcodec/dbresolver"
)

type DeleteDeliveryOrderConsumerHandlerInterface interface {
	ProcessMessage()
}

type DeleteDeliveryOrderConsumerHandler struct {
	kafkaClient                    kafkadbo.KafkaClientInterface
	DeliveryOrderOpenSearchUseCase usecases.DeliveryOrderConsumerUseCaseInterface
	ctx                            context.Context
	args                           []interface{}
	db                             dbresolver.DB
	deliveryOrderLogRepository     mongoRepositories.DeliveryOrderLogRepositoryInterface
}

func InitDeleteDeliveryOrderConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, deliveryOrderLogRepository mongoRepositories.DeliveryOrderLogRepositoryInterface, DeliveryOrderOpenSearchUseCase usecases.DeliveryOrderConsumerUseCaseInterface, db dbresolver.DB, ctx context.Context, args []interface{}) DeleteDeliveryOrderConsumerHandlerInterface {
	return &DeleteDeliveryOrderConsumerHandler{
		kafkaClient:                    kafkaClient,
		DeliveryOrderOpenSearchUseCase: DeliveryOrderOpenSearchUseCase,
		ctx:                            ctx,
		args:                           args,
		db:                             db,
		deliveryOrderLogRepository:     deliveryOrderLogRepository,
	}
}

func (c *DeleteDeliveryOrderConsumerHandler) ProcessMessage() {
	fmt.Println("process ", constants.DELETE_DELIVERY_ORDER_TOPIC)
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
			Action:    constants.LOG_ACTION_MONGO_DELETE,
			Data:      string(m.Value),
			Status:    constants.LOG_STATUS_MONGO_ERROR,
			CreatedAt: &now,
		}

		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.DELETE_DELIVERY_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
			deliveryOrderLog.Error = errorLogData
			go c.deliveryOrderLogRepository.Insert(deliveryOrderLog, c.ctx, deliveryOrderLogResultChan)
			fmt.Println(deliveryOrderLog.Error)
			continue
		}
		go c.deliveryOrderLogRepository.GetByCode(deliveryOrder.DoCode, constants.LOG_STATUS_MONGO_DEFAULT, deliveryOrderLog.Action, false, c.ctx, deliveryOrderLogResultChan)
		deliveryOrderDetailResult := <-deliveryOrderLogResultChan
		if deliveryOrderDetailResult.Error != nil {
			errorLogData := helper.WriteLogConsumer(constants.DELETE_DELIVERY_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), deliveryOrderDetailResult.Error, http.StatusInternalServerError, nil)
			deliveryOrderLog.Error = errorLogData
			go c.deliveryOrderLogRepository.Insert(deliveryOrderLog, c.ctx, deliveryOrderLogResultChan)
			fmt.Println(deliveryOrderLog.Error)
			continue
		}
		deliveryOrderLog = deliveryOrderDetailResult.DeliveryOrderLog
		deliveryOrderLog.Status = constants.LOG_STATUS_MONGO_ERROR
		deliveryOrderLog.UpdatedAt = &now
		errorLog := &model.ErrorLog{}
		if deliveryOrder.DeliveryOrderDetails == nil {
			errorLog = c.DeliveryOrderOpenSearchUseCase.SyncToOpenSearchFromDeleteEvent(&deliveryOrder.ID, nil, c.ctx)
		} else {
			deliveryOrderDetailIds := []*int{}
			for _, v := range deliveryOrder.DeliveryOrderDetails {
				deliveryOrderDetailIds = append(deliveryOrderDetailIds, &v.ID)
			}
			errorLog = c.DeliveryOrderOpenSearchUseCase.SyncToOpenSearchFromDeleteEvent(&deliveryOrder.ID, deliveryOrderDetailIds, c.ctx)
		}

		if errorLog.Err != nil {
			dbTransaction.Rollback()
			errorLogData := helper.WriteLogConsumer(constants.DELETE_DELIVERY_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), errorLog.Err, http.StatusInternalServerError, nil)
			deliveryOrderLog.Error = errorLogData
			go c.deliveryOrderLogRepository.UpdateByID(deliveryOrderLog.ID.Hex(), deliveryOrderLog, c.ctx, deliveryOrderLogResultChan)
			fmt.Println(deliveryOrderLog.Error)
			continue
		}

		err = dbTransaction.Commit()
		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.DELETE_DELIVERY_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
			deliveryOrderLog.Error = errorLogData
			go c.deliveryOrderLogRepository.UpdateByID(deliveryOrderLog.ID.Hex(), deliveryOrderLog, c.ctx, deliveryOrderLogResultChan)
			fmt.Println(deliveryOrderLog.Error)
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
