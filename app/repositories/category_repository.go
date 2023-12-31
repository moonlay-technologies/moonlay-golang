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

type CategoryRepositoryInterface interface {
	GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.CategoryChan)
	GetByParentID(parentId int, countOnly bool, ctx context.Context, result chan *models.CategoryChan)
}

type category struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitCategoryRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) CategoryRepositoryInterface {
	return &category{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *category) GetByID(ID int, countOnly bool, ctx context.Context, resultChan chan *models.CategoryChan) {
	response := &models.CategoryChan{}
	var category models.Category
	var total int64

	categoryRedisKey := fmt.Sprintf("%s:%d", "category", ID)
	categoryOnRedis, err := r.redisdb.Client().Get(ctx, categoryRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM categories WHERE deleted_at IS NULL AND id = ?", ID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			errStr := fmt.Sprintf("category id %d data not found", ID)
			err = helper.NewError(errStr)
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			category = models.Category{}
			err = r.db.QueryRow(""+
				"SELECT id, name, parent_id from categories as c "+
				"WHERE c.deleted_at IS NULL AND c.id = ?", ID).
				Scan(&category.ID, &category.Name, &category.ParentID)

			if err != nil {
				errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			categoryJson, _ := json.Marshal(category)
			setCategoryOnRedis := r.redisdb.Client().Set(ctx, categoryRedisKey, categoryJson, 1*time.Hour)

			if setCategoryOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setCategoryOnRedis.Err(), 500, "Something went wrong, please try again later")
				response.Error = setCategoryOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.Category = &category
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
		_ = json.Unmarshal([]byte(categoryOnRedis), &category)
		response.Category = &category
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *category) GetByParentID(parentId int, countOnly bool, ctx context.Context, resultChan chan *models.CategoryChan) {
	response := &models.CategoryChan{}
	var category models.Category
	var total int64

	categoryRedisKey := fmt.Sprintf("%s:%d", "category-by-parent-id", parentId)
	categoryOnRedis, err := r.redisdb.Client().Get(ctx, categoryRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM categories WHERE deleted_at IS NULL AND parent_id = ?", parentId).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			errStr := fmt.Sprintf("category parent id %d data not found", parentId)
			err = helper.NewError(errStr)
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			category = models.Category{}
			err = r.db.QueryRow(""+
				"SELECT id, name from categories as c "+
				"WHERE c.deleted_at IS NULL AND c.parent_id = ?", parentId).
				Scan(&category.ID, &category.Name)

			if err != nil {
				errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			categoryJson, _ := json.Marshal(category)
			setCategoryOnRedis := r.redisdb.Client().Set(ctx, categoryRedisKey, categoryJson, 1*time.Hour)

			if setCategoryOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setCategoryOnRedis.Err(), 500, "Something went wrong, please try again later")
				response.Error = setCategoryOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.Category = &category
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
		_ = json.Unmarshal([]byte(categoryOnRedis), &category)
		response.Category = &category
		response.Total = total
		resultChan <- response
		return
	}
}
