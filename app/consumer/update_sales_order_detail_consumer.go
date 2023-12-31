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

type UpdateSalesOrderDetailConsumerHandlerInterface interface {
	ProcessMessage()
}

type UpdateSalesOrderDetailConsumerHandler struct {
	kafkaClient                 kafkadbo.KafkaClientInterface
	salesOrderOpenSearchUseCase usecases.SalesOrderOpenSearchUseCaseInterface
	ctx                         context.Context
	args                        []interface{}
	db                          dbresolver.DB
	salesOrderLogRepository     mongoRepositories.SalesOrderLogRepositoryInterface
}

func InitUpdateSalesOrderDetailConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, salesOrderLogRepository mongoRepositories.SalesOrderLogRepositoryInterface, salesOrderOpenSearchUseCase usecases.SalesOrderOpenSearchUseCaseInterface, db dbresolver.DB, ctx context.Context, args []interface{}) UpdateSalesOrderDetailConsumerHandlerInterface {
	return &UpdateSalesOrderDetailConsumerHandler{
		kafkaClient:                 kafkaClient,
		salesOrderOpenSearchUseCase: salesOrderOpenSearchUseCase,
		ctx:                         ctx,
		args:                        args,
		db:                          db,
		salesOrderLogRepository:     salesOrderLogRepository,
	}
}

func (c *UpdateSalesOrderDetailConsumerHandler) ProcessMessage() {
	fmt.Println("process", constants.UPDATE_SALES_ORDER_DETAIL_TOPIC)
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
		key := string(m.Key[:])
		now := time.Now()
		salesOrderLogResultChan := make(chan *models.SalesOrderLogChan)

		dbTransaction, err := c.db.BeginTx(c.ctx, nil)
		salesOrderLog := &models.SalesOrderLog{
			RequestID: "",
			SoCode:    "",
			Data:      m.Value,
			Action:    constants.LOG_ACTION_MONGO_UPDATE,
			Status:    constants.LOG_STATUS_MONGO_ERROR,
			CreatedAt: &now,
		}

		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.UPDATE_SALES_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
			salesOrderLog.Error = errorLogData
			go c.salesOrderLogRepository.Insert(salesOrderLog, c.ctx, salesOrderLogResultChan)
			fmt.Println(errorLogData)
			continue
		}

		go c.salesOrderLogRepository.GetByCollumn(constants.COLUMN_SALES_ORDER_CODE, salesOrder.SoCode, false, c.ctx, salesOrderLogResultChan)
		salesOrderDetailResult := <-salesOrderLogResultChan
		if salesOrderDetailResult.Error != nil {
			salesOrderLog.Error = salesOrderDetailResult.Error
			go c.salesOrderLogRepository.Insert(salesOrderLog, c.ctx, salesOrderLogResultChan)
			fmt.Println(salesOrderDetailResult.Error)
			continue
		}
		logId := salesOrderDetailResult.ID.Hex()
		salesOrderLog = salesOrderDetailResult.SalesOrderLog
		salesOrderLog.Status = constants.LOG_STATUS_MONGO_ERROR
		salesOrderLog.UpdatedAt = &now

		errorLog := c.salesOrderOpenSearchUseCase.SyncToOpenSearchFromUpdateEvent(&salesOrder, c.ctx)
		if errorLog.Err != nil {
			dbTransaction.Rollback()
			errorLogData := helper.WriteLogConsumer(constants.UPDATE_SALES_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), errorLog.Err, http.StatusInternalServerError, nil)
			salesOrderLog.Error = errorLogData
			go c.salesOrderLogRepository.UpdateByID(logId, salesOrderLog, c.ctx, salesOrderLogResultChan)
			fmt.Println(errorLogData)
			continue
		}

		salesOrderDetail := &models.SalesOrderDetail{}
		for _, v := range salesOrder.SalesOrderDetails {
			if v.SoDetailCode == key {
				salesOrderDetail = v
				break
			}
		}
		errorLog = c.salesOrderOpenSearchUseCase.SyncDetailToOpenSearchFromUpdateEvent(&salesOrder, salesOrderDetail, dbTransaction, c.ctx)
		if errorLog.Err != nil {
			dbTransaction.Rollback()
			errorLogData := helper.WriteLogConsumer(constants.UPDATE_SALES_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), errorLog.Err, http.StatusInternalServerError, nil)
			salesOrderLog.Error = errorLogData
			go c.salesOrderLogRepository.UpdateByID(logId, salesOrderLog, c.ctx, salesOrderLogResultChan)
			fmt.Println(errorLogData)
			continue
		}

		err = dbTransaction.Commit()
		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.UPDATE_SALES_ORDER_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
			salesOrderLog.Error = errorLogData
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
