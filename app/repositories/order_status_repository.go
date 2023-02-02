package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bxcodec/dbresolver"
	"github.com/go-redis/redis/v8"
	"order-service/app/models"
	"order-service/global/utils/helper"
	"order-service/global/utils/redisdb"
	"time"
)

type OrderStatusRepositoryInterface interface {
	GetByNameAndType(name string, statusType string, countOnly bool, ctx context.Context, result chan *models.OrderStatusChan)
	GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.OrderStatusChan)
}

type orderStatus struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitOrderStatusRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) OrderStatusRepositoryInterface {
	return &orderStatus{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *orderStatus) GetByNameAndType(name string, statusType string, countOnly bool, ctx context.Context, resultChan chan *models.OrderStatusChan) {
	response := &models.OrderStatusChan{}
	var orderStatus models.OrderStatus
	var total int64

	orderStatusRedisKey := fmt.Sprintf("%s:%s:%s", "order-status", name, statusType)
	orderStatusOnRedis, err := r.redisdb.Client().Get(ctx, orderStatusRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM order_statuses WHERE deleted_at IS NULL AND name = ? AND order_type = ?", name, statusType).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("order_status data not found")
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			orderStatus = models.OrderStatus{}
			err = r.db.QueryRow(""+
				"SELECT id, name, order_type "+
				"FROM order_statuses as os "+
				"WHERE os.deleted_at IS NULL AND os.name = ? AND os.order_type = ?", name, statusType).
				Scan(&orderStatus.ID, &orderStatus.Name, &orderStatus.OrderType)

			if err != nil {
				errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			orderStatusJson, err := json.Marshal(orderStatus)

			if err != nil {
				errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			_, err = r.redisdb.Client().Set(ctx, orderStatusRedisKey, orderStatusJson, 1*time.Hour).Result()

			if err != nil {
				errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.OrderStatus = &orderStatus
			resultChan <- response
			return
		}

	} else if err != nil {
		errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	} else {
		total = 1
		_ = json.Unmarshal([]byte(orderStatusOnRedis), &orderStatus)
		response.OrderStatus = &orderStatus
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *orderStatus) GetByID(ID int, countOnly bool, ctx context.Context, resultChan chan *models.OrderStatusChan) {
	response := &models.OrderStatusChan{}
	var orderStatus models.OrderStatus
	var total int64

	orderStatusRedisKey := fmt.Sprintf("%s:%d", "order-status", ID)
	orderStatusOnRedis, err := r.redisdb.Client().Get(ctx, orderStatusRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM order_statuses WHERE deleted_at IS NULL AND id = ?", ID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("order_status data not found")
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			orderStatus = models.OrderStatus{}
			err = r.db.QueryRow(""+
				"SELECT id, name, sequence, order_type  from order_statuses as os "+
				"WHERE os.deleted_at IS NULL AND os.id = ?", ID).
				Scan(&orderStatus.ID, &orderStatus.Name, &orderStatus.Sequence, &orderStatus.OrderType)

			if err != nil {
				errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			orderStatusJson, _ := json.Marshal(orderStatus)
			setOrderStatusOnRedis := r.redisdb.Client().Set(ctx, orderStatusRedisKey, orderStatusJson, 1*time.Hour)

			if setOrderStatusOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setOrderStatusOnRedis.Err(), 500, "Something went wrong, please try again later")
				response.Error = setOrderStatusOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.OrderStatus = &orderStatus
			resultChan <- response
			return
		}

	} else if err != nil {
		errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	} else {
		total = 1
		_ = json.Unmarshal([]byte(orderStatusOnRedis), &orderStatus)
		response.OrderStatus = &orderStatus
		response.Total = total
		resultChan <- response
		return
	}
}
