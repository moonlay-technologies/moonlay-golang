package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bxcodec/dbresolver"
	"github.com/go-redis/redis/v8"
	"poc-order-service/app/models"
	"poc-order-service/global/utils/helper"
	"poc-order-service/global/utils/redisdb"
	"time"
)

type UomRepositoryInterface interface {
	GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.UomChan)
}

type uom struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitUomRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) UomRepositoryInterface {
	return &uom{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *uom) GetByID(ID int, countOnly bool, ctx context.Context, resultChan chan *models.UomChan) {
	response := &models.UomChan{}
	var uom models.Uom
	var total int64

	uomRedisKey := fmt.Sprintf("%s:%d", "uom", ID)
	uomOnRedis, err := r.redisdb.Client().Get(ctx, uomRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM uoms WHERE deleted_at IS NULL AND id = ?", ID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("uom data not found")
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			uom = models.Uom{}
			err = r.db.QueryRow(""+
				"SELECT id, code, name FROM uoms "+
				"WHERE uoms.deleted_at IS NULL AND uoms.id = ?", ID).
				Scan(&uom.ID, &uom.Code, &uom.Name)

			if err != nil {
				errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			uomJson, _ := json.Marshal(uom)
			setUomOnRedis := r.redisdb.Client().Set(ctx, uomRedisKey, uomJson, 1*time.Hour)

			if setUomOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setUomOnRedis.Err(), 500, "Something went wrong, please try again later")
				response.Error = setUomOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.Uom = &uom
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
		_ = json.Unmarshal([]byte(uomOnRedis), &uom)
		response.Uom = &uom
		response.Total = total
		resultChan <- response
		return
	}
}
