package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/app/repositories"
	mongoRepositories "order-service/app/repositories/mongod"
	"order-service/app/usecases"
	"order-service/global/utils/helper"
	kafkadbo "order-service/global/utils/kafka"
	"time"

	"github.com/bxcodec/dbresolver"
)

type DeleteSalesOrderDetailConsumerHandlerInterface interface {
	ProcessMessage()
}

type DeleteSalesOrderDetailConsumerHandler struct {
	kafkaClient                 kafkadbo.KafkaClientInterface
	salesOrderOpenSearchUseCase usecases.SalesOrderOpenSearchUseCaseInterface
	salesOrderRepository        repositories.SalesOrderRepositoryInterface
	ctx                         context.Context
	args                        []interface{}
	db                          dbresolver.DB
	salesOrderLogRepository     mongoRepositories.SalesOrderLogRepositoryInterface
}

func InitDeleteSalesOrderDetailConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, salesOrderLogRepository mongoRepositories.SalesOrderLogRepositoryInterface, salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderOpenSearchUseCase usecases.SalesOrderOpenSearchUseCaseInterface, db dbresolver.DB, ctx context.Context, args []interface{}) DeleteSalesOrderDetailConsumerHandlerInterface {
	return &DeleteSalesOrderDetailConsumerHandler{
		kafkaClient:                 kafkaClient,
		salesOrderOpenSearchUseCase: salesOrderOpenSearchUseCase,
		salesOrderRepository:        salesOrderRepository,
		ctx:                         ctx,
		args:                        args,
		db:                          db,
		salesOrderLogRepository:     salesOrderLogRepository,
	}
}

func (c *DeleteSalesOrderDetailConsumerHandler) ProcessMessage() {
	fmt.Println("process", constants.DELETE_SALES_ORDER_DETAIL_TOPIC)
	topic := c.args[1].(string)
	groupID := c.args[2].(string)
	reader := c.kafkaClient.SetConsumerGroupReader(topic, groupID)

	for {
		m, err := reader.ReadMessage(c.ctx)
		if err != nil {
			break
		}

		fmt.Printf("message at topic/partition/offset %v/%v/%v \n", m.Topic, m.Partition, m.Offset)

		var salesOrderDetail models.SalesOrderDetail
		err = json.Unmarshal(m.Value, &salesOrderDetail)
		now := time.Now()
		salesOrderLogResultChan := make(chan *models.SalesOrderLogChan)

		dbTransaction, err := c.db.BeginTx(c.ctx, nil)
		salesOrderLog := &models.SalesOrderLog{
			RequestID: "",
			SoCode:    "",
			Data:      m.Value,
			Action:    constants.LOG_ACTION_MONGO_DELETE,
			Status:    constants.LOG_STATUS_MONGO_ERROR,
			CreatedAt: &now,
		}

		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.DELETE_SALES_ORDER_DETAIL_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
			salesOrderLog.Error = errorLogData
			go c.salesOrderLogRepository.Insert(salesOrderLog, c.ctx, salesOrderLogResultChan)
			fmt.Println(errorLogData)
			continue
		}

		getSalesOrderByIDResultChan := make(chan *models.SalesOrderChan)
		go c.salesOrderRepository.GetByID(salesOrderDetail.SalesOrderID, false, c.ctx, getSalesOrderByIDResultChan)
		getSalesOrderByIDResult := <-getSalesOrderByIDResultChan

		if getSalesOrderByIDResult.Error != nil {
			salesOrderLog.Error = getSalesOrderByIDResult.Error
			go c.salesOrderLogRepository.Insert(salesOrderLog, c.ctx, salesOrderLogResultChan)
			fmt.Println(getSalesOrderByIDResult.ErrorLog)
			continue
		}

		go c.salesOrderLogRepository.GetByCollumn(constants.COLUMN_SALES_ORDER_CODE, getSalesOrderByIDResult.SalesOrder.SoCode, false, c.ctx, salesOrderLogResultChan)
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
		errorLog := c.salesOrderOpenSearchUseCase.SyncDetailToOpenSearchFromDeleteEvent(&salesOrderDetail, c.ctx)
		if errorLog.Err != nil {
			dbTransaction.Rollback()
			errorLogData := helper.WriteLogConsumer(constants.DELETE_SALES_ORDER_DETAIL_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), errorLog.Err, http.StatusInternalServerError, nil)
			salesOrderLog.Error = errorLogData
			go c.salesOrderLogRepository.UpdateByID(logId, salesOrderLog, c.ctx, salesOrderLogResultChan)
			fmt.Println(errorLogData)
			continue
		}

		err = dbTransaction.Commit()
		if err != nil {
			errorLogData := helper.WriteLogConsumer(constants.DELETE_SALES_ORDER_DETAIL_CONSUMER, m.Topic, m.Partition, m.Offset, string(m.Key), err, http.StatusInternalServerError, nil)
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
