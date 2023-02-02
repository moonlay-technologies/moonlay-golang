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

type AgentBrandRepositoryInterface interface {
	GetByAgentIDAndBrandID(agentID int, brandID int, countOnly bool, ctx context.Context, result chan *models.AgentBrandChan)
}

type agentBrand struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitAgentBrandRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) AgentBrandRepositoryInterface {
	return &agentBrand{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *agentBrand) GetByAgentIDAndBrandID(agentID int, brandID int, countOnly bool, ctx context.Context, resultChan chan *models.AgentBrandChan) {
	response := &models.AgentBrandChan{}
	var agentBrand models.AgentBrand
	var total int64

	agentBrandRedisKey := fmt.Sprintf("%s:%d:%d", "agentBrand", agentID, brandID)
	agentBrandOnRedis, err := r.redisdb.Client().Get(ctx, agentBrandRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM agent_brand WHERE agent_id = ? AND brand_id = ?", agentID, brandID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, 500, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("agent_brands data not found")
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			agentBrand = models.AgentBrand{}
			err = r.db.QueryRow(""+
				"SELECT agent_id, brand_id FROM agent_brand "+
				"WHERE agent_id = ? AND brand_id = ?", agentID, brandID).
				Scan(&agentBrand.AgentID, &agentBrand.BrandID)

			if err != nil {
				errorLogData := helper.WriteLog(err, 500, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			agentBrandJson, _ := json.Marshal(agentBrand)
			setAgentBrandOnRedis := r.redisdb.Client().Set(ctx, agentBrandRedisKey, agentBrandJson, 1*time.Hour)

			if setAgentBrandOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setAgentBrandOnRedis.Err(), 500, nil)
				response.Error = setAgentBrandOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.AgentBrand = &agentBrand
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
		_ = json.Unmarshal([]byte(agentBrandOnRedis), &agentBrand)
		response.AgentBrand = &agentBrand
		response.Total = total
		resultChan <- response
		return
	}
}
