package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/global/utils/helper"
	"order-service/global/utils/redisdb"
	"time"

	"github.com/bxcodec/dbresolver"
	"github.com/go-redis/redis/v8"
)

type OrderSourceRepositoryInterface interface {
	GetBySourceName(sourceName string, countOnly bool, ctx context.Context, result chan *models.OrderSourceChan)
	GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.OrderSourceChan)
}

type orderSource struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitOrderSourceRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) OrderSourceRepositoryInterface {
	return &orderSource{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *orderSource) GetBySourceName(sourceName string, countOnly bool, ctx context.Context, resultChan chan *models.OrderSourceChan) {
	response := &models.OrderSourceChan{}
	var orderSource models.OrderSource
	var total int64

	orderSourceRedisKey := fmt.Sprintf("%s:%s", "order-source", sourceName)
	orderSourceOnRedis, err := r.redisdb.Client().Get(ctx, orderSourceRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM order_sources WHERE deleted_at IS NULL AND source_name = ?", sourceName).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("order_source data not found")
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			orderSource = models.OrderSource{}
			err = r.db.QueryRow(""+
				"SELECT id, code, source_name  from order_sources as os"+
				"WHERE os.deleted_at IS NULL AND os.source_name = ?", sourceName).
				Scan(&orderSource.ID, &orderSource.Code, &orderSource.SourceName)

			if err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			orderSourceJson, _ := json.Marshal(orderSource)
			setOrderSourceOnRedis := r.redisdb.Client().Set(ctx, orderSourceRedisKey, orderSourceJson, 1*time.Hour)

			if setOrderSourceOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setOrderSourceOnRedis.Err(), http.StatusInternalServerError, nil)
				response.Error = setOrderSourceOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.OrderSource = &orderSource
			resultChan <- response
			return
		}

	} else if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	} else {
		total = 1
		_ = json.Unmarshal([]byte(orderSourceOnRedis), &orderSource)
		response.OrderSource = &orderSource
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *orderSource) GetByID(ID int, countOnly bool, ctx context.Context, resultChan chan *models.OrderSourceChan) {
	response := &models.OrderSourceChan{}
	var orderSource models.OrderSource
	var total int64

	orderSourceRedisKey := fmt.Sprintf("%s:%d", "order-source", ID)
	orderSourceOnRedis, err := r.redisdb.Client().Get(ctx, orderSourceRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM order_sources WHERE deleted_at IS NULL AND id = ?", ID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("order_source data not found")
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			orderSource = models.OrderSource{}
			err = r.db.QueryRow(""+
				"SELECT id, code, source_name  from order_sources as os "+
				"WHERE os.deleted_at IS NULL AND os.id = ?", ID).
				Scan(&orderSource.ID, &orderSource.Code, &orderSource.SourceName)

			if err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			orderSourceJson, _ := json.Marshal(orderSource)
			setOrderSourceOnRedis := r.redisdb.Client().Set(ctx, orderSourceRedisKey, orderSourceJson, 1*time.Hour)

			if setOrderSourceOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setOrderSourceOnRedis.Err(), http.StatusInternalServerError, nil)
				response.Error = setOrderSourceOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.OrderSource = &orderSource
			resultChan <- response
			return
		}

	} else if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	} else {
		total = 1
		_ = json.Unmarshal([]byte(orderSourceOnRedis), &orderSource)
		response.OrderSource = &orderSource
		response.Total = total
		resultChan <- response
		return
	}
}
