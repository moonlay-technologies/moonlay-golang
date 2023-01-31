package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bxcodec/dbresolver"
	"github.com/go-redis/redis/v8"
	"poc-order-service/app/models"
	"poc-order-service/global/utils/helper"
	"poc-order-service/global/utils/redisdb"
	"strings"
	"time"
)

type CartDetailRepositoryInterface interface {
	GetByCartID(cartID, countOnly bool, ctx context.Context, result chan *models.CartDetailsChan)
	Insert(request *models.CartDetail, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.CartDetailChan)
}

type cartDetail struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitCartDetailRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) CartDetailRepositoryInterface {
	return &cartDetail{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *cartDetail) GetByCartID(cartID, countOnly bool, ctx context.Context, resultChan chan *models.CartDetailsChan) {
	response := &models.CartDetailsChan{}
	var cartDetails []*models.CartDetail
	var total int64

	cartDetailsRedisKey := fmt.Sprintf("%s:%d", "cart-details-by-cart-id", cartID)
	cartDetailsOnRedis, err := r.redisdb.Client().Get(ctx, cartDetailsRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM cart_details WHERE deleted_at IS NULL AND cart_id = ? ", cartID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("cart_detail data not found")
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			query, err := r.db.Query(""+
				"SELECT id, cart_id, brand_id, product_id, uom_id, order_status_id, qty, price "+
				"FROM cart_details as cd "+
				"WHERE cd.deleted_at IS NULL AND cd.cart_id = ?", cartID)

			if err != nil {
				errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			for query.Next() {
				var cartDetail models.CartDetail
				err = query.Scan(&cartDetail.ID, &cartDetail.CartID, &cartDetail.BrandID, &cartDetail.ProductID, &cartDetail.UomID, &cartDetail.OrderStatusID, &cartDetail.Qty, &cartDetail.Price)

				if err != nil {
					errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
					response.Error = err
					response.ErrorLog = errorLogData
					resultChan <- response
					return
				}

				cartDetails = append(cartDetails, &cartDetail)
			}

			cartDetailsJson, _ := json.Marshal(cartDetails)
			setCartDetailsOnRedis := r.redisdb.Client().Set(ctx, cartDetailsRedisKey, cartDetailsJson, 1*time.Hour)

			if setCartDetailsOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setCartDetailsOnRedis.Err(), 500, "Something went wrong, please try again later")
				response.Error = setCartDetailsOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.CartDetails = cartDetails
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
		_ = json.Unmarshal([]byte(cartDetailsOnRedis), &cartDetails)
		response.CartDetails = cartDetails
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *cartDetail) Insert(request *models.CartDetail, sqlTransaction *sql.Tx, ctx context.Context, resultChan chan *models.CartDetailChan) {
	response := &models.CartDetailChan{}
	rawSqlFields := []string{}
	rawSqlDataTypes := []string{}
	rawSqlValues := []interface{}{}

	if request.CartID != 0 {
		rawSqlFields = append(rawSqlFields, "cart_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.CartID)
	}

	if request.BrandID != 0 {
		rawSqlFields = append(rawSqlFields, "brand_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.BrandID)
	}

	if request.ProductID != 0 {
		rawSqlFields = append(rawSqlFields, "product_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.ProductID)
	}

	if request.UomID != 0 {
		rawSqlFields = append(rawSqlFields, "uom_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.UomID)
	}

	if request.OrderStatusID != 0 {
		rawSqlFields = append(rawSqlFields, "order_status_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.OrderStatusID)
	}

	if request.Qty != 0 {
		rawSqlFields = append(rawSqlFields, "qty")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.Qty)
	}

	if request.Price != 0 {
		rawSqlFields = append(rawSqlFields, "price")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.Price)
	}

	rawSqlFields = append(rawSqlFields, "created_at")
	rawSqlDataTypes = append(rawSqlDataTypes, "?")
	rawSqlValues = append(rawSqlValues, request.CreatedAt.Format("2006-01-02 15:04:05"))

	rawSqlFieldsJoin := strings.Join(rawSqlFields, ",")
	rawSqlDataTypesJoin := strings.Join(rawSqlDataTypes, ",")

	query := fmt.Sprintf("INSERT INTO cart_details (%s) VALUES (%v)", rawSqlFieldsJoin, rawSqlDataTypesJoin)
	result, err := sqlTransaction.ExecContext(ctx, query, rawSqlValues...)

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	cartDetailID, err := result.LastInsertId()

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	response.ID = cartDetailID
	request.ID = int(cartDetailID)
	response.CartDetail = request
	resultChan <- response
	return
}
