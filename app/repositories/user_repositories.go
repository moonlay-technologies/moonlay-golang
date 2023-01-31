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

type UserRepositoryInterface interface {
	GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.UserChan)
}

type user struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitUserRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) UserRepositoryInterface {
	return &user{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *user) GetByID(ID int, countOnly bool, ctx context.Context, resultChan chan *models.UserChan) {
	response := &models.UserChan{}
	var user models.User
	var total int64

	userRedisKey := fmt.Sprintf("%s:%d", "user", ID)
	userOnRedis, err := r.redisdb.Client().Get(ctx, userRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM users WHERE deleted_at IS NULL AND id = ?", ID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("user data not found")
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			user = models.User{}
			err = r.db.QueryRow(""+
				"SELECT id, email, first_name, last_name, role_id, status, is_admin, is_owner FROM users "+
				"WHERE users.deleted_at IS NULL AND users.id = ?", ID).
				Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.RoleID, &user.Status, &user.IsAdmin, &user.IsOwner)

			if err != nil {
				errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			userJson, _ := json.Marshal(user)
			setUserOnRedis := r.redisdb.Client().Set(ctx, userRedisKey, userJson, 1*time.Hour)

			if setUserOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setUserOnRedis.Err(), 500, "Something went wrong, please try again later")
				response.Error = setUserOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.User = &user
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
		_ = json.Unmarshal([]byte(userOnRedis), &user)
		response.User = &user
		response.Total = total
		resultChan <- response
		return
	}
}
