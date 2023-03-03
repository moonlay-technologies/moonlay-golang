package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/global/utils/helper"
	"order-service/global/utils/redisdb"
	"strings"
	"time"

	"github.com/bxcodec/dbresolver"
	"github.com/go-redis/redis/v8"
)

type SalesOrderDetailRepositoryInterface interface {
	GetBySalesOrderID(salesOrderID int, countOnly bool, ctx context.Context, result chan *models.SalesOrderDetailsChan)
	Insert(request *models.SalesOrderDetail, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.SalesOrderDetailChan)
	GetByID(salesOrderDetailID int, countOnly bool, ctx context.Context, result chan *models.SalesOrderDetailChan)
	UpdateByID(id int, request *models.SalesOrderDetail, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.SalesOrderDetailChan)
	RemoveCacheByID(id int, ctx context.Context, resultChan chan *models.SalesOrderDetailChan)
	DeleteByID(request *models.SalesOrderDetail, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.SalesOrderDetailChan)
}

type salesOrderDetail struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitSalesOrderDetailRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) SalesOrderDetailRepositoryInterface {
	return &salesOrderDetail{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *salesOrderDetail) GetBySalesOrderID(salesOrderID int, countOnly bool, ctx context.Context, resultChan chan *models.SalesOrderDetailsChan) {
	response := &models.SalesOrderDetailsChan{}
	var salesOrderDetails []*models.SalesOrderDetail
	var total int64

	salesOrderDetailRedisKey := fmt.Sprintf("%s:%d", constants.SALES_ORDER_DETAIL_BY_SALES_ORDER_ID, salesOrderID)
	salesOrderDetailsRedis, err := r.redisdb.Client().Get(ctx, salesOrderDetailRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM sales_order_details WHERE deleted_at IS NULL AND sales_order_id = ?", salesOrderID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("sales_order_detail data not found")
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			query, err := r.db.Query(""+
				"SELECT sod.id, sod.sales_order_id, sod.uom_id, sod.order_status_id, sod.so_detail_code, p.id as product_id, p.SKU as product_sku, p.productName as product_name, sod.qty, sod.sent_qty, sod.residual_qty, u.code as uom_code, sod.price, os.name as order_detail_status, sod.note, sod.is_done_sync_to_es, sod.start_date_sync_to_es, sod.end_date_sync_to_es, sod.created_at "+
				"FROM sales_order_details as sod "+
				"INNER JOIN products as p ON p.id = sod.product_id "+
				"INNER JOIN uoms as u ON u.id = sod.uom_id "+
				"INNER JOIN order_statuses as os ON os.id = sod.order_status_id "+
				"WHERE sod.deleted_at IS NULL AND sod.sales_order_id = ?", salesOrderID)

			if err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			for query.Next() {
				var salesOrderDetail models.SalesOrderDetail
				err = query.Scan(&salesOrderDetail.ID, &salesOrderDetail.SalesOrderID, &salesOrderDetail.UomID, &salesOrderDetail.OrderStatusID, &salesOrderDetail.SoDetailCode, &salesOrderDetail.ProductID, &salesOrderDetail.ProductSKU, &salesOrderDetail.ProductName, &salesOrderDetail.Qty, &salesOrderDetail.SentQty, &salesOrderDetail.ResidualQty, &salesOrderDetail.UomCode, &salesOrderDetail.Price, &salesOrderDetail.OrderStatusName, &salesOrderDetail.Note, &salesOrderDetail.IsDoneSyncToEs, &salesOrderDetail.StartDateSyncToEs, &salesOrderDetail.EndDateSyncToEs, &salesOrderDetail.CreatedAt)

				if err != nil {
					errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
					response.Error = err
					response.ErrorLog = errorLogData
					resultChan <- response
					return
				}

				salesOrderDetails = append(salesOrderDetails, &salesOrderDetail)
			}

			salesOrderDetailsJson, _ := json.Marshal(salesOrderDetails)
			setSalesOrderDetailsOnRedis := r.redisdb.Client().Set(ctx, salesOrderDetailsRedis, salesOrderDetailsJson, 1*time.Hour)

			if setSalesOrderDetailsOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setSalesOrderDetailsOnRedis.Err(), http.StatusInternalServerError, nil)
				response.Error = setSalesOrderDetailsOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.SalesOrderDetails = salesOrderDetails
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
		_ = json.Unmarshal([]byte(salesOrderDetailsRedis), &salesOrderDetails)
		response.SalesOrderDetails = salesOrderDetails
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *salesOrderDetail) Insert(request *models.SalesOrderDetail, sqlTransaction *sql.Tx, ctx context.Context, resultChan chan *models.SalesOrderDetailChan) {
	response := &models.SalesOrderDetailChan{}

	rawSqlFields := []string{}
	rawSqlDataTypes := []string{}
	rawSqlValues := []interface{}{}

	if request.SalesOrderID != 0 {
		rawSqlFields = append(rawSqlFields, "sales_order_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.SalesOrderID)
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

	if request.SoDetailCode != "" {
		rawSqlFields = append(rawSqlFields, "so_detail_code")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.SoDetailCode)
	}

	if request.Qty != 0 {
		rawSqlFields = append(rawSqlFields, "qty")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.Qty)
	}

	rawSqlFields = append(rawSqlFields, "sent_qty")
	rawSqlDataTypes = append(rawSqlDataTypes, "?")
	rawSqlValues = append(rawSqlValues, request.SentQty)

	rawSqlFields = append(rawSqlFields, "residual_qty")
	rawSqlDataTypes = append(rawSqlDataTypes, "?")
	rawSqlValues = append(rawSqlValues, request.ResidualQty)

	if request.Price != 0 {
		rawSqlFields = append(rawSqlFields, "price")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.Price)
	}

	if request.Note.String != "" {
		rawSqlFields = append(rawSqlFields, "note")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.Note)
	}

	if request.IsDoneSyncToEs != "" {
		rawSqlFields = append(rawSqlFields, "is_done_sync_to_es")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.IsDoneSyncToEs)
	}

	if request.StartDateSyncToEs != nil {
		rawSqlFields = append(rawSqlFields, "start_date_sync_to_es")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.StartDateSyncToEs.Format("2006-01-02 15:04:05"))
	}

	rawSqlFields = append(rawSqlFields, "created_at")
	rawSqlDataTypes = append(rawSqlDataTypes, "?")
	rawSqlValues = append(rawSqlValues, request.CreatedAt.Format("2006-01-02 15:04:05"))

	rawSqlFieldsJoin := strings.Join(rawSqlFields, ",")
	rawSqlDataTypesJoin := strings.Join(rawSqlDataTypes, ",")

	query := fmt.Sprintf("INSERT INTO sales_order_details (%s) VALUES (%v)", rawSqlFieldsJoin, rawSqlDataTypesJoin)
	result, err := sqlTransaction.ExecContext(ctx, query, rawSqlValues...)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	salesOrderDetailID, err := result.LastInsertId()

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	response.ID = salesOrderDetailID
	request.ID = int(salesOrderDetailID)
	response.SalesOrderDetail = request
	resultChan <- response
	return
}

func (r *salesOrderDetail) GetByID(id int, countOnly bool, ctx context.Context, resultChan chan *models.SalesOrderDetailChan) {
	response := &models.SalesOrderDetailChan{}
	var salesOrderDetail models.SalesOrderDetail
	var total int64

	salesOrderDetailRedisKey := fmt.Sprintf("%s:%d", constants.SALES_ORDER_DETAIL, id)
	salesOrderDetailOnRedis, err := r.redisdb.Client().Get(ctx, salesOrderDetailRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM sales_order_details WHERE deleted_at IS NULL AND id = ?", id).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("sales_order_detail data not found")
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			salesOrderDetail = models.SalesOrderDetail{}
			err = r.db.QueryRow(""+
				"SELECT id, product_id, uom_id, order_status_id, sales_order_id, qty, sent_qty, residual_qty, price, note, so_detail_code, created_at "+
				"FROM sales_order_details as sod "+
				"WHERE sod.deleted_at IS NULL AND sod.id = ?", id).
				Scan(&salesOrderDetail.ID, &salesOrderDetail.ProductID, &salesOrderDetail.UomID, &salesOrderDetail.OrderStatusID, &salesOrderDetail.SalesOrderID, &salesOrderDetail.Qty, &salesOrderDetail.SentQty, &salesOrderDetail.ResidualQty, &salesOrderDetail.Price, &salesOrderDetail.Note, &salesOrderDetail.SoDetailCode, &salesOrderDetail.CreatedAt)

			if err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			salesOrderDetailJson, _ := json.Marshal(salesOrderDetail)
			setSalesOrderDetailOnRedis := r.redisdb.Client().Set(ctx, salesOrderDetailRedisKey, salesOrderDetailJson, 1*time.Hour)

			if setSalesOrderDetailOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setSalesOrderDetailOnRedis.Err(), http.StatusInternalServerError, nil)
				response.Error = setSalesOrderDetailOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.SalesOrderDetail = &salesOrderDetail
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
		_ = json.Unmarshal([]byte(salesOrderDetailOnRedis), &salesOrderDetail)
		response.SalesOrderDetail = &salesOrderDetail
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *salesOrderDetail) UpdateByID(id int, request *models.SalesOrderDetail, sqlTransaction *sql.Tx, ctx context.Context, resultChan chan *models.SalesOrderDetailChan) {
	response := &models.SalesOrderDetailChan{}
	rawSqlQueries := []string{}

	if request.Qty != 0 {
		query := fmt.Sprintf("%s=%v", "product_id", request.ProductID)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.Qty != 0 {
		query := fmt.Sprintf("%s=%v", "uom_id", request.UomID)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.ResidualQty != 0 || request.SentQty != 0 {
		query := fmt.Sprintf("%s=%v", "residual_qty", request.ResidualQty)
		rawSqlQueries = append(rawSqlQueries, query)
		query = fmt.Sprintf("%s=%v", "sent_qty", request.SentQty)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.OrderStatusID != 0 {
		query := fmt.Sprintf("%s=%v", "order_status_id", request.OrderStatusID)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if len(request.Note.String) > 0 {
		query := fmt.Sprintf("%s='%v'", "note", request.Note.String)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if len(request.IsDoneSyncToEs) > 0 {
		query := fmt.Sprintf("%s='%v'", "is_done_sync_to_es", request.IsDoneSyncToEs)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.EndDateSyncToEs != nil {
		query := fmt.Sprintf("%s='%v'", "end_date_sync_to_es", request.EndDateSyncToEs.Format("2006-01-02 15:04:05"))
		rawSqlQueries = append(rawSqlQueries, query)
	}

	query := fmt.Sprintf("%s='%v'", "updated_at", request.UpdatedAt.Format("2006-01-02 15:04:05"))
	rawSqlQueries = append(rawSqlQueries, query)

	rawSqlQueriesJoin := strings.Join(rawSqlQueries, ",")

	updateQuery := fmt.Sprintf("UPDATE sales_order_details set %v WHERE id = ?", rawSqlQueriesJoin)
	result, err := sqlTransaction.ExecContext(ctx, updateQuery, id)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	salesOrderDetailID, err := result.LastInsertId()

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	response.ID = salesOrderDetailID
	response.SalesOrderDetail = request
	resultChan <- response
	return
}

func (r *salesOrderDetail) RemoveCacheByID(id int, ctx context.Context, resultChan chan *models.SalesOrderDetailChan) {
	response := &models.SalesOrderDetailChan{}
	salesOrderRedisKey := fmt.Sprintf("%s:%d", constants.SALES_ORDER_DETAIL, id)
	_, err := r.redisdb.Client().Del(ctx, salesOrderRedisKey).Result()

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	response.Error = nil
	resultChan <- response
	return
}

func (r *salesOrderDetail) DeleteByID(request *models.SalesOrderDetail, sqlTransaction *sql.Tx, ctx context.Context, resultChan chan *models.SalesOrderDetailChan) {
	now := time.Now()
	request.DeletedAt = &now
	request.UpdatedAt = &now
	response := &models.SalesOrderDetailChan{}
	rawSqlQueries := []string{}

	query := fmt.Sprintf("%s='%v'", "deleted_at", request.DeletedAt.Format("2006-01-02 15:04:05"))
	rawSqlQueries = append(rawSqlQueries, query)

	query = fmt.Sprintf("%s='%v'", "updated_at", request.UpdatedAt.Format("2006-01-02 15:04:05"))
	rawSqlQueries = append(rawSqlQueries, query)

	rawSqlQueriesJoin := strings.Join(rawSqlQueries, ",")

	updateQuery := fmt.Sprintf("UPDATE sales_order_details set %v WHERE id = ?", rawSqlQueriesJoin)
	result, err := sqlTransaction.ExecContext(ctx, updateQuery, request.ID)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	salesOrderDetailID, err := result.LastInsertId()

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	response.ID = salesOrderDetailID
	response.SalesOrderDetail = request
	resultChan <- response
	return
}
