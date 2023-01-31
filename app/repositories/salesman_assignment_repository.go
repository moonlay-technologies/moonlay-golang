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

type SalesmanAssignmentRepositoryInterface interface {
	GetBySalesmanIDAndAgentID(agentID int, salesmanID int, countOnly bool, ctx context.Context, result chan *models.SalesmanAssignmentChan)
}

type salesmanAssignment struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitSalesmanAssignmentRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) SalesmanAssignmentRepositoryInterface {
	return &salesmanAssignment{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *salesmanAssignment) GetBySalesmanIDAndAgentID(agentID int, salesmanID int, countOnly bool, ctx context.Context, resultChan chan *models.SalesmanAssignmentChan) {
	response := &models.SalesmanAssignmentChan{}
	var salesmanAssignment models.SalesmanAssignment
	var total int64

	salesmanAssignmentRedisKey := fmt.Sprintf("%s:%d:%d", "salesmanAssignment", salesmanID, agentID)
	salesmanAssignmentOnRedis, err := r.redisdb.Client().Get(ctx, salesmanAssignmentRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM salesman_assignment as sa inner join salesmans as s ON s.id = sa.salesman_id WHERE sa.deleted_at IS NULL AND s.deleted_at IS NULL AND sa.salesman_id = ? AND s.agent_id = ?", salesmanID, agentID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("salesmanAssignments data not found")
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			salesmanAssignment = models.SalesmanAssignment{}
			err = r.db.QueryRow(""+
				"SELECT sa.id, sa.salesman_id, sa.brand_id, s.agent_id FROM salesman_assignment as sa "+
				"INNER JOIN salesmans as s "+
				"WHERE sa.deleted_at IS NULL AND s.deleted_at IS NULL AND sa.salesman_id = ? AND s.agent_id = ?", salesmanID, agentID).
				Scan(&salesmanAssignment.ID, &salesmanAssignment.SalesmanID, &salesmanAssignment.BrandID, &salesmanAssignment.AgentID)

			if err != nil {
				errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			salesmanAssignmentJson, _ := json.Marshal(salesmanAssignment)
			setSalesmanAssignmentOnRedis := r.redisdb.Client().Set(ctx, salesmanAssignmentRedisKey, salesmanAssignmentJson, 1*time.Hour)

			if setSalesmanAssignmentOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setSalesmanAssignmentOnRedis.Err(), 500, "Something went wrong, please try again later")
				response.Error = setSalesmanAssignmentOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.SalesmanAssignment = &salesmanAssignment
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
		_ = json.Unmarshal([]byte(salesmanAssignmentOnRedis), &salesmanAssignment)
		response.SalesmanAssignment = &salesmanAssignment
		response.Total = total
		resultChan <- response
		return
	}
}
