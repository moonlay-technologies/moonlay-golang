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

type DeleteSalesOrderConsumerHandlerInterface interface {
	ProcessMessage()
}

type DeleteSalesOrderConsumerHandler struct {
	kafkaClient             kafkadbo.KafkaClientInterface
	salesOrderUseCase       usecases.SalesOrderUseCaseInterface
	ctx                     context.Context
	args                    []interface{}
	db                      dbresolver.DB
	salesOrderLogRepository mongoRepositories.SalesOrderLogRepositoryInterface
}

func InitDeleteSalesOrderConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, salesOrderLogRepository mongoRepositories.SalesOrderLogRepositoryInterface, salesOrderUseCase usecases.SalesOrderUseCaseInterface, db dbresolver.DB, ctx context.Context, args []interface{}) UpdateSalesOrderConsumerHandlerInterface {
	return &updateSalesOrderConsumerHandler{
		kafkaClient:             kafkaClient,
		salesOrderUseCase:       salesOrderUseCase,
		ctx:                     ctx,
		args:                    args,
		db:                      db,
		salesOrderLogRepository: salesOrderLogRepository,
	}
}

func (c *DeleteSalesOrderConsumerHandler) ProcessMessage() {
	fmt.Println("process", constants.UPDATE_SALES_ORDER_TOPIC)
	now := time.Now()
	topic := c.args[1].(string)
	groupID := c.args[2].(string)
	reader := c.kafkaClient.SetConsumerGroupReader(topic, groupID)
	salesOrderLogResultChan := make(chan *models.SalesOrderLogChan)

	for {
		m, err := reader.ReadMessage(c.ctx)
		if err != nil {
			break
		}

		fmt.Printf("message at topic/partition/offset %v/%v/%v \n", m.Topic, m.Partition, m.Offset)

		var salesOrder models.SalesOrder
		err = json.Unmarshal(m.Value, &salesOrder)

		dbTransaction, err := c.db.BeginTx(c.ctx, nil)
		salesOrderLog := &models.SalesOrderLog{
			RequestID: "",
			SoCode:    "",
			Data:      m.Value,
			Status:    constants.LOG_STATUS_MONGO_ERROR,
			CreatedAt: &now,
		}

		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.UPDATE_SALES_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
			go c.salesOrderLogRepository.Insert(salesOrderLog, c.ctx, salesOrderLogResultChan)
			fmt.Println(errorLogData)
			continue
		}

		go c.salesOrderLogRepository.GetByCollumn(constants.SALES_ORDER_CODE_COLLUMN, salesOrder.SoCode, false, c.ctx, salesOrderLogResultChan)
		salesOrderDetailResult := <-salesOrderLogResultChan
		if salesOrderDetailResult.Error != nil {
			go c.salesOrderLogRepository.Insert(salesOrderLog, c.ctx, salesOrderLogResultChan)
			fmt.Println(salesOrderDetailResult.Error)
			continue
		}
		logId := salesOrderDetailResult.ID.Hex()
		salesOrderLog = salesOrderDetailResult.SalesOrderLog
		salesOrderLog.Status = constants.LOG_STATUS_MONGO_ERROR
		salesOrderLog.UpdatedAt = &now
		errorLog := c.salesOrderUseCase.SyncToOpenSearchFromUpdateEvent(&salesOrder, c.ctx)
		if errorLog.Err != nil {
			dbTransaction.Rollback()
			errorLogData := helper.WriteLogConsumer(constants.UPDATE_SALES_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), errorLog.Err, http.StatusInternalServerError, nil)
			go c.salesOrderLogRepository.UpdateByID(logId, salesOrderLog, c.ctx, salesOrderLogResultChan)
			fmt.Println(errorLogData)
			continue
		}

		err = dbTransaction.Commit()
		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.UPDATE_SALES_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
			go c.salesOrderLogRepository.UpdateByID(logId, salesOrderLog, c.ctx, salesOrderLogResultChan)
			fmt.Println(errorLogData)
			continue
		}

		salesOrderLog.Status = constants.LOG_STATUS_MONGO_SUCCESS
		go c.salesOrderLogRepository.UpdateByID(logId, salesOrderLog, c.ctx, salesOrderLogResultChan)
	}

	if err := reader.Close(); err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
}