package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/app/models"
	"order-service/global/utils/helper"
	"order-service/global/utils/redisdb"
	"time"

	"github.com/bxcodec/dbresolver"
	"github.com/go-redis/redis/v8"
)

type SalesmanRepositoryInterface interface {
	GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.SalesmanChan)
	GetByEmail(email string, countOnly bool, ctx context.Context, result chan *models.SalesmanChan)
}

type salesman struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitSalesmanRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) SalesmanRepositoryInterface {
	return &salesman{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *salesman) GetByID(ID int, countOnly bool, ctx context.Context, resultChan chan *models.SalesmanChan) {
	response := &models.SalesmanChan{}
	var salesman models.Salesman
	var total int64

	salesmanRedisKey := fmt.Sprintf("%s:%d", "salesman", ID)
	salesmanOnRedis, err := r.redisdb.Client().Get(ctx, salesmanRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM salesmans WHERE deleted_at IS NULL AND id = ?", ID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, 500, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("salesman data not found")
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			salesman = models.Salesman{}
			err = r.db.QueryRow(""+
				"SELECT id, name, email, phone_number FROM salesmans "+
				"WHERE salesmans.deleted_at IS NULL AND salesmans.id = ?", ID).
				Scan(&salesman.ID, &salesman.Name, &salesman.Email, &salesman.PhoneNumber)

			if err != nil {
				errorLogData := helper.WriteLog(err, 500, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			salesmanJson, _ := json.Marshal(salesman)
			setSalesmanOnRedis := r.redisdb.Client().Set(ctx, salesmanRedisKey, salesmanJson, 1*time.Hour)

			if setSalesmanOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setSalesmanOnRedis.Err(), 500, "Something went wrong, please try again later")
				response.Error = setSalesmanOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.Salesman = &salesman
			resultChan <- response
			return
		}

	} else if err != nil {
		errorLogData := helper.WriteLog(err, 500, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	} else {
		total = 1
		_ = json.Unmarshal([]byte(salesmanOnRedis), &salesman)
		response.Salesman = &salesman
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *salesman) GetByEmail(email string, countOnly bool, ctx context.Context, resultChan chan *models.SalesmanChan) {
	response := &models.SalesmanChan{}
	var salesman models.Salesman
	var total int64

	salesmanRedisKey := fmt.Sprintf("%s:%s", "salesman", email)
	salesmanOnRedis, err := r.redisdb.Client().Get(ctx, salesmanRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM salesmans WHERE deleted_at IS NULL AND email = ?", email).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, 500, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("salesman data not found")
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			salesman = models.Salesman{}
			err = r.db.QueryRow(""+
				"SELECT id, name, email, phone_number FROM salesmans "+
				"WHERE salesmans.deleted_at IS NULL AND salesmans.email = ?", email).
				Scan(&salesman.ID, &salesman.Name, &salesman.Email, &salesman.PhoneNumber)

			if err != nil {
				errorLogData := helper.WriteLog(err, 500, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			salesmanJson, _ := json.Marshal(salesman)
			setSalesmanOnRedis := r.redisdb.Client().Set(ctx, salesmanRedisKey, salesmanJson, 1*time.Hour)

			if setSalesmanOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setSalesmanOnRedis.Err(), 500, "Something went wrong, please try again later")
				response.Error = setSalesmanOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.Salesman = &salesman
			resultChan <- response
			return
		}

	} else if err != nil {
		errorLogData := helper.WriteLog(err, 500, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	} else {
		total = 1
		_ = json.Unmarshal([]byte(salesmanOnRedis), &salesman)
		response.Salesman = &salesman
		response.Total = total
		resultChan <- response
		return
	}
}
