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

type CreateDeliveryOrderConsumerHandlerInterface interface {
	ProcessMessage()
}

type createDeliveryOrderConsumerHandler struct {
	kafkaClient                 kafkadbo.KafkaClientInterface
	salesOrderUseCase           usecases.SalesOrderUseCaseInterface
	salesOrderOpenSearchUseCase usecases.SalesOrderOpenSearchUseCaseInterface
	deliveryOrderUseCase        usecases.DeliveryOrderUseCaseInterface
	ctx                         context.Context
	args                        []interface{}
	db                          dbresolver.DB
	deliveryOrderLogRepository  mongoRepositories.DeliveryOrderLogRepositoryInterface
}

func InitCreateDeliveryOrderConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, deliveryOrderLogRepository mongoRepositories.DeliveryOrderLogRepositoryInterface, salesOrderUseCase usecases.SalesOrderUseCaseInterface, salesOrderOpenSearchUseCase usecases.SalesOrderOpenSearchUseCaseInterface, deliveryOrderUseCase usecases.DeliveryOrderUseCaseInterface, db dbresolver.DB, ctx context.Context, args []interface{}) CreateDeliveryOrderConsumerHandlerInterface {
	return &createDeliveryOrderConsumerHandler{
		kafkaClient:                 kafkaClient,
		salesOrderUseCase:           salesOrderUseCase,
		salesOrderOpenSearchUseCase: salesOrderOpenSearchUseCase,
		deliveryOrderUseCase:        deliveryOrderUseCase,
		ctx:                         ctx,
		args:                        args,
		db:                          db,
		deliveryOrderLogRepository:  deliveryOrderLogRepository,
	}
}

func (c *createDeliveryOrderConsumerHandler) ProcessMessage() {
	fmt.Println("process")
	now := time.Now()
	topic := c.args[1].(string)
	groupID := c.args[2].(string)
	reader := c.kafkaClient.SetConsumerGroupReader(topic, groupID)
	deliveryOrderLogResultChan := make(chan *models.DeliveryOrderLogChan)

	for {
		m, err := reader.ReadMessage(c.ctx)
		if err != nil {
			break
		}

		fmt.Printf("message do at topic/partition/offset %v/%v/%v \n", m.Topic, m.Partition, m.Offset)

		var deliveryOrder models.DeliveryOrder
		err = json.Unmarshal(m.Value, &deliveryOrder)

		dbTransaction, err := c.db.BeginTx(c.ctx, nil)
		deliveryOrderLog := &models.DeliveryOrderLog{
			RequestID: "",
			DoCode:    "",
			Data:      m.Value,
			Action:    constants.LOG_ACTION_MONGO_INSERT,
			Status:    constants.LOG_STATUS_MONGO_ERROR,
			CreatedAt: &now,
		}

		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.CREATE_DELIVERY_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
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

		errorLog := c.deliveryOrderUseCase.SyncToOpenSearchFromCreateEvent(&deliveryOrder, c.salesOrderUseCase, dbTransaction, c.ctx)

		if errorLog.Err != nil {
			dbTransaction.Rollback()
			errorLogData := helper.WriteLogConsumer(constants.CREATE_DELIVERY_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), errorLog.Err, http.StatusInternalServerError, nil)
			go c.deliveryOrderLogRepository.UpdateByID(deliveryOrderLog.ID.Hex(), deliveryOrderLog, c.ctx, deliveryOrderLogResultChan)
			fmt.Println(errorLogData)
			continue
		}

		err = dbTransaction.Commit()
		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.CREATE_DELIVERY_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
			go c.deliveryOrderLogRepository.UpdateByID(deliveryOrderLog.ID.Hex(), deliveryOrderLog, c.ctx, deliveryOrderLogResultChan)
			fmt.Println(errorLogData)
			continue
		}

		salesOrderRequest := &models.SalesOrderRequest{
			ID:            deliveryOrder.SalesOrderID,
			OrderSourceID: deliveryOrder.OrderSourceID,
		}

		salesOrderWithDetail, errorLog := c.salesOrderUseCase.GetByIDWithDetail(salesOrderRequest, c.ctx)

		if errorLog.Err != nil {
			go c.deliveryOrderLogRepository.UpdateByID(deliveryOrderLog.ID.Hex(), deliveryOrderLog, c.ctx, deliveryOrderLogResultChan)
			fmt.Println(errorLog)
			continue
		}

		errorLog = c.salesOrderOpenSearchUseCase.SyncToOpenSearchFromUpdateEvent(salesOrderWithDetail, c.ctx)

		if errorLog.Err != nil {
			go c.deliveryOrderLogRepository.UpdateByID(deliveryOrderLog.ID.Hex(), deliveryOrderLog, c.ctx, deliveryOrderLogResultChan)
			fmt.Println(errorLog)
			continue
		}

		deliveryOrderRequest := &models.DeliveryOrderRequest{
			ID: deliveryOrder.ID,
		}

		deliveryOrderWithDetail, errorLog := c.deliveryOrderUseCase.GetByID(deliveryOrderRequest, c.ctx)

		if errorLog.Err != nil {
			go c.deliveryOrderLogRepository.UpdateByID(deliveryOrderLog.ID.Hex(), deliveryOrderLog, c.ctx, deliveryOrderLogResultChan)
			fmt.Println(errorLog)
			continue
		}

		salesOrderWithDetail.DeliveryOrders = nil
		deliveryOrderWithDetail.SalesOrder = salesOrderWithDetail
		errorLog = c.deliveryOrderUseCase.SyncToOpenSearchFromUpdateEvent(deliveryOrderWithDetail, c.ctx)

		if errorLog.Err != nil {
			go c.deliveryOrderLogRepository.UpdateByID(deliveryOrderLog.ID.Hex(), deliveryOrderLog, c.ctx, deliveryOrderLogResultChan)
			errorLogData := helper.WriteLogConsumer(constants.CREATE_DELIVERY_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), errorLog.Err, http.StatusInternalServerError, nil)
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
