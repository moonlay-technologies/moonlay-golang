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

type AgentRepositoryInterface interface {
	GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.AgentChan)
}

type agent struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitAgentRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) AgentRepositoryInterface {
	return &agent{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *agent) GetByID(ID int, countOnly bool, ctx context.Context, resultChan chan *models.AgentChan) {
	response := &models.AgentChan{}
	var agent models.Agent
	var total int64

	agentRedisKey := fmt.Sprintf("%s:%d", "agent", ID)
	agentOnRedis, err := r.redisdb.Client().Get(ctx, agentRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM agents WHERE deleted_at IS NULL AND id = ?", ID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, 500, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("agents data not found")
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			agent = models.Agent{}
			err = r.db.QueryRow(""+
				"SELECT agents.id, agents.name, email, address, description, agents.province_id, agents.city_id, agents.district_id, agents.village_id, distributor_type, postal_code, g_lat, g_lng, contact_name, website, phone, main_mobile_phone, status, no_npwp, image_npwp, no_siup, provinces.name as province_name, cities.name as city_name, districts.name as district_name, villages.name as village_name FROM agents "+
				"INNER JOIN provinces on provinces.id = agents.province_id "+
				"INNER JOIN cities on cities.id = agents.city_id "+
				"INNER JOIN districts on districts.id = agents.district_id "+
				"INNEr JOIN villages on villages.id = agents.village_id "+
				"WHERE agents.deleted_at IS NULL AND agents.id = ?", ID).
				Scan(&agent.ID, &agent.Name, &agent.Email, &agent.Address, &agent.Description, &agent.ProvinceID, &agent.CityID, &agent.DistrictID, &agent.VillageID, &agent.DistributorType, &agent.PostalCode, &agent.GLat, &agent.GLng, &agent.ContactName, &agent.Website, &agent.Phone, &agent.MainMobilePhone, &agent.Status, &agent.NoNpwp, &agent.ImageNpwp, &agent.NoSiup, &agent.ProvinceName, &agent.CityName, &agent.DistrictName, &agent.VillageName)

			if err != nil {
				errorLogData := helper.WriteLog(err, 500, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			agentJson, _ := json.Marshal(agent)
			setAgentOnRedis := r.redisdb.Client().Set(ctx, agentRedisKey, agentJson, 1*time.Hour)

			if setAgentOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setAgentOnRedis.Err(), 500, "Something went wrong, please try again later")
				response.Error = setAgentOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.Agent = &agent
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
		_ = json.Unmarshal([]byte(agentOnRedis), &agent)
		response.Agent = &agent
		response.Total = total
		resultChan <- response
		return
	}
}
