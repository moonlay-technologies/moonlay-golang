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

type DeleteDeliveryOrderDetailConsumerHandlerInterface interface {
	ProcessMessage()
}

type DeleteDeliveryOrderDetailConsumerHandler struct {
	kafkaClient                    kafkadbo.KafkaClientInterface
	DeliveryOrderOpenSearchUseCase usecases.DeliveryOrderOpenSearchUseCaseInterface
	ctx                            context.Context
	args                           []interface{}
	db                             dbresolver.DB
	deliveryOrderLogRepository     mongoRepositories.DeliveryOrderLogRepositoryInterface
}

func InitDeleteDeliveryOrderDetailConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, deliveryOrderLogRepository mongoRepositories.DeliveryOrderLogRepositoryInterface, DeliveryOrderOpenSearchUseCase usecases.DeliveryOrderOpenSearchUseCaseInterface, db dbresolver.DB, ctx context.Context, args []interface{}) DeleteDeliveryOrderDetailConsumerHandlerInterface {
	return &DeleteDeliveryOrderDetailConsumerHandler{
		kafkaClient:                    kafkaClient,
		DeliveryOrderOpenSearchUseCase: DeliveryOrderOpenSearchUseCase,
		ctx:                            ctx,
		args:                           args,
		db:                             db,
		deliveryOrderLogRepository:     deliveryOrderLogRepository,
	}
}

func (c *DeleteDeliveryOrderDetailConsumerHandler) ProcessMessage() {
	fmt.Println("process ", constants.DELETE_DELIVERY_ORDER_DETAIL_TOPIC)
	topic := c.args[1].(string)
	groupID := c.args[2].(string)
	reader := c.kafkaClient.SetConsumerGroupReader(topic, groupID)

	for {
		m, err := reader.ReadMessage(c.ctx)
		if err != nil {
			break
		}

		fmt.Printf("message do at topic/partition/offset %v/%v/%v \n", m.Topic, m.Partition, m.Offset)

		var deliveryOrderDetail models.DeliveryOrderDetail
		err = json.Unmarshal(m.Value, &deliveryOrderDetail)
		now := time.Now()
		deliveryOrderDetailLogResultChan := make(chan *models.DeliveryOrderLogChan)

		dbTransaction, err := c.db.BeginTx(c.ctx, nil)
		deliveryOrderDetailLog := &models.DeliveryOrderLog{
			RequestID: "",
			DoCode:    "",
			Action:    constants.LOG_ACTION_MONGO_DELETE,
			Data:      m.Value,
			Status:    constants.LOG_STATUS_MONGO_ERROR,
			CreatedAt: &now,
		}

		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.DELETE_DELIVERY_ORDER_DETAIL_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
			deliveryOrderDetailLog.Error = errorLogData
			go c.deliveryOrderLogRepository.Insert(deliveryOrderDetailLog, c.ctx, deliveryOrderDetailLogResultChan)
			fmt.Println(deliveryOrderDetailLog.Error)
			continue
		}
		go c.deliveryOrderLogRepository.GetByCode(deliveryOrderDetail.DoDetailCode, constants.LOG_STATUS_MONGO_DEFAULT, deliveryOrderDetailLog.Action, false, c.ctx, deliveryOrderDetailLogResultChan)
		deliveryOrderDetailResult := <-deliveryOrderDetailLogResultChan
		if deliveryOrderDetailResult.Error != nil {
			errorLogData := helper.WriteLogConsumer(constants.DELETE_DELIVERY_ORDER_DETAIL_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), deliveryOrderDetailResult.Error, http.StatusInternalServerError, nil)
			deliveryOrderDetailLog.Error = errorLogData
			go c.deliveryOrderLogRepository.Insert(deliveryOrderDetailLog, c.ctx, deliveryOrderDetailLogResultChan)
			fmt.Println(deliveryOrderDetailLog.Error)
			continue
		}
		deliveryOrderDetailLog = deliveryOrderDetailResult.DeliveryOrderLog
		deliveryOrderDetailLog.Status = constants.LOG_STATUS_MONGO_ERROR
		deliveryOrderDetailLog.UpdatedAt = &now

		errorLog := c.DeliveryOrderOpenSearchUseCase.SyncToOpenSearchFromDeleteEvent(&deliveryOrderDetail.DeliveryOrderID, []*int{&deliveryOrderDetail.ID}, c.ctx)

		if errorLog.Err != nil {
			dbTransaction.Rollback()
			errorLogData := helper.WriteLogConsumer(constants.DELETE_DELIVERY_ORDER_DETAIL_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), errorLog.Err, http.StatusInternalServerError, nil)
			deliveryOrderDetailLog.Error = errorLogData
			go c.deliveryOrderLogRepository.UpdateByID(deliveryOrderDetailLog.ID.Hex(), deliveryOrderDetailLog, c.ctx, deliveryOrderDetailLogResultChan)
			fmt.Println(deliveryOrderDetailLog.Error)
			continue
		}

		err = dbTransaction.Commit()
		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.DELETE_DELIVERY_ORDER_DETAIL_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
			deliveryOrderDetailLog.Error = errorLogData
			go c.deliveryOrderLogRepository.UpdateByID(deliveryOrderDetailLog.ID.Hex(), deliveryOrderDetailLog, c.ctx, deliveryOrderDetailLogResultChan)
			fmt.Println(deliveryOrderDetailLog.Error)
			continue
		}

		deliveryOrderDetailLog.Status = constants.LOG_STATUS_MONGO_SUCCESS
		go c.deliveryOrderLogRepository.UpdateByID(deliveryOrderDetailLog.ID.Hex(), deliveryOrderDetailLog, c.ctx, deliveryOrderDetailLogResultChan)
	}

	if err := reader.Close(); err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
}
