package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"order-service/app/models"
	"order-service/global/utils/helper"
	"order-service/global/utils/redisdb"
	"strings"
	"time"

	"github.com/bxcodec/dbresolver"
	"github.com/go-redis/redis/v8"
)

type CartRepositoryInterface interface {
	GetByUserID(userID int, cartStatusID int, countOnly bool, ctx context.Context, result chan *models.CartChan)
	Insert(request *models.Cart, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.CartChan)
}

type cart struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitCartRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) CartRepositoryInterface {
	return &cart{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *cart) GetByUserID(userID int, cartStatusID int, countOnly bool, ctx context.Context, resultChan chan *models.CartChan) {
	response := &models.CartChan{}
	var cart models.Cart
	var total int64

	cartRedisKey := fmt.Sprintf("%s:%d:%d", "cart-user-status", userID, cartStatusID)
	cartOnRedis, err := r.redisdb.Client().Get(ctx, cartRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM carts WHERE deleted_at IS NULL AND user_id = ? AND order_status_id= ?", userID, cartStatusID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, 500, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("cart data not found")
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			cart = models.Cart{}
			err = r.db.QueryRow(""+
				"SELECT id, agent_id, brand_id, visitation_id, user_id, store_id "+
				"FROM carts as c "+
				"WHERE c.deleted_at IS NULL AND c.user_id = ? AND c.order_status_id = ?", userID, cartStatusID).
				Scan(&cart.ID, &cart.AgentID, &cart.BrandID, &cart.VisitationID, &cart.UserID, &cart.StoreID)

			if err != nil {
				errorLogData := helper.WriteLog(err, 500, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			cartJson, _ := json.Marshal(cart)
			setCartOnRedis := r.redisdb.Client().Set(ctx, cartRedisKey, cartJson, 1*time.Hour)

			if setCartOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setCartOnRedis.Err(), 500, "Something went wrong, please try again later")
				response.Error = setCartOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.Cart = &cart
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
		_ = json.Unmarshal([]byte(cartOnRedis), &cart)
		response.Cart = &cart
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *cart) Insert(request *models.Cart, sqlTransaction *sql.Tx, ctx context.Context, resultChan chan *models.CartChan) {
	response := &models.CartChan{}
	rawSqlFields := []string{}
	rawSqlDataTypes := []string{}
	rawSqlValues := []interface{}{}

	if request.AgentID != 0 {
		rawSqlFields = append(rawSqlFields, "agent_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.AgentID)
	}

	if request.BrandID != 0 {
		rawSqlFields = append(rawSqlFields, "brand_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.BrandID)
	}

	if request.VisitationID != 0 {
		rawSqlFields = append(rawSqlFields, "visitation_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.VisitationID)
	}

	if request.UserID != 0 {
		rawSqlFields = append(rawSqlFields, "user_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.UserID)
	}

	if request.StoreID != 0 {
		rawSqlFields = append(rawSqlFields, "store_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.StoreID)
	}

	if request.OrderStatusID != 0 {
		rawSqlFields = append(rawSqlFields, "order_status_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.OrderStatusID)
	}

	if request.OrderSourceID != 0 {
		rawSqlFields = append(rawSqlFields, "order_source_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.OrderSourceID)
	}

	if request.TotalTonase != 0 {
		rawSqlFields = append(rawSqlFields, "total_tonase")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.TotalTonase)
	}

	if request.TotalAmount != 0 {
		rawSqlFields = append(rawSqlFields, "total_amount")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.TotalAmount)
	}

	if request.Note != "" {
		rawSqlFields = append(rawSqlFields, "note")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.Note)
	}

	rawSqlFields = append(rawSqlFields, "created_at")
	rawSqlDataTypes = append(rawSqlDataTypes, "?")
	rawSqlValues = append(rawSqlValues, request.CreatedAt.Format("2006-01-02 15:04:05"))

	rawSqlFieldsJoin := strings.Join(rawSqlFields, ",")
	rawSqlDataTypesJoin := strings.Join(rawSqlDataTypes, ",")

	query := fmt.Sprintf("INSERT INTO carts (%s) VALUES (%v)", rawSqlFieldsJoin, rawSqlDataTypesJoin)
	result, err := sqlTransaction.ExecContext(ctx, query, rawSqlValues...)

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	salesOrderID, err := result.LastInsertId()

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	response.ID = salesOrderID
	request.ID = int(salesOrderID)
	response.Cart = request
	resultChan <- response
	return
}
