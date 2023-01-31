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

type ProductRepositoryInterface interface {
	GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.ProductChan)
}

type product struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitProductRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) ProductRepositoryInterface {
	return &product{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *product) GetByID(ID int, countOnly bool, ctx context.Context, resultChan chan *models.ProductChan) {
	response := &models.ProductChan{}
	var product models.Product
	var total int64

	productRedisKey := fmt.Sprintf("%s:%d", "product", ID)
	productOnRedis, err := r.redisdb.Client().Get(ctx, productRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM products WHERE deleted_at IS NULL AND id = ?", ID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			errStr := fmt.Sprintf("product id %d data not found", ID)
			err = helper.NewError(errStr)
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			product = models.Product{}
			err = r.db.QueryRow(""+
				"SELECT id, SKU, productName, unitMeasurementSmall, unitMeasurementMedium, isActive  from products as p "+
				"WHERE p.deleted_at IS NULL AND p.id = ?", ID).
				Scan(&product.ID, &product.Sku, &product.ProductName, &product.UnitMeasurementSmall, &product.UnitMeasurementMedium, &product.IsActive)

			if err != nil {
				errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			productJson, _ := json.Marshal(product)
			setProductOnRedis := r.redisdb.Client().Set(ctx, productRedisKey, productJson, 1*time.Hour)

			if setProductOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setProductOnRedis.Err(), 500, "Something went wrong, please try again later")
				response.Error = setProductOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.Product = &product
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
		_ = json.Unmarshal([]byte(productOnRedis), &product)
		response.Product = &product
		response.Total = total
		resultChan <- response
		return
	}
}
