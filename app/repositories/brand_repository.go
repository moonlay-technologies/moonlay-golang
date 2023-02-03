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

type BrandRepositoryInterface interface {
	GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.BrandChan)
}

type brand struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitBrandRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) BrandRepositoryInterface {
	return &brand{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *brand) GetByID(ID int, countOnly bool, ctx context.Context, resultChan chan *models.BrandChan) {
	response := &models.BrandChan{}
	var brand models.Brand
	var total int64

	brandRedisKey := fmt.Sprintf("%s:%d", "brand", ID)
	brandOnRedis, err := r.redisdb.Client().Get(ctx, brandRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM brands WHERE deleted_at IS NULL AND id = ?", ID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("brands data not found")
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			brand = models.Brand{}
			err = r.db.QueryRow(""+
				"SELECT id, name, principle_id, description, image, status_active FROM brands "+
				"WHERE brands.deleted_at IS NULL AND brands.id = ?", ID).
				Scan(&brand.ID, &brand.Name, &brand.PrincipleID, &brand.Description, &brand.Image, &brand.StatusActive)

			if err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			brandJson, _ := json.Marshal(brand)
			setBrandOnRedis := r.redisdb.Client().Set(ctx, brandRedisKey, brandJson, 1*time.Hour)

			if setBrandOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setBrandOnRedis.Err(), http.StatusInternalServerError, nil)
				response.Error = setBrandOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.Brand = &brand
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
		_ = json.Unmarshal([]byte(brandOnRedis), &brand)
		response.Brand = &brand
		response.Total = total
		resultChan <- response
		return
	}
}
