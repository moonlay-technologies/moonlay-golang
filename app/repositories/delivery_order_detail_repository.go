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

type DeliveryOrderDetailRepositoryInterface interface {
	GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.DeliveryOrderDetailChan)
	GetByDeliveryOrderID(deliveryOrderID int, countOnly bool, ctx context.Context, result chan *models.DeliveryOrderDetailsChan)
	GetBySalesOrderID(salesOrderID int, countOnly bool, ctx context.Context, result chan *models.SalesOrderDetailsChan)
	Insert(request *models.DeliveryOrderDetail, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.DeliveryOrderDetailChan)
	UpdateByID(id int, deliveryOrderDetail *models.DeliveryOrderDetail, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.DeliveryOrderDetailChan)
	DeleteByID(request *models.DeliveryOrderDetail, sqlTransaction *sql.Tx, ctx context.Context, resultChan chan *models.DeliveryOrderDetailChan)
}

type deliveryOrderDetail struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitDeliveryOrderDetailRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) DeliveryOrderDetailRepositoryInterface {
	return &deliveryOrderDetail{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *deliveryOrderDetail) GetByID(ID int, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderDetailChan) {
	response := &models.DeliveryOrderDetailChan{}
	var deliveryOrderDetail models.DeliveryOrderDetail
	var total int64

	deliveryOrderDetailRedisKey := fmt.Sprintf("%s:%d", constants.DELIVERY_ORDER_DETAIL_BY_ID, ID)
	deliveryOrderDetailsRedis, err := r.redisdb.Client().Get(ctx, deliveryOrderDetailRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM delivery_order_details WHERE deleted_at IS NULL AND id = ?", ID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("delivery_order_detail data not found")
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			query, err := r.db.Query(""+
				"SELECT dod.id, dod.delivery_order_id, dod.so_detail_id, dod.brand_id, dod.product_id, dod.uom_id, dod.order_status_id, dod.do_detail_code, dod.qty, dod.note, dod.is_done_sync_to_es, dod.start_date_sync_to_es, dod.end_date_sync_to_es, dod.created_at, dod.updated_at "+
				"FROM delivery_order_details as dod "+
				"INNER JOIN sales_order_details as sod ON sod.id = dod.so_detail_id "+
				"INNER JOIN brands as b ON b.id = dod.brand_id "+
				"INNER JOIN products as p ON p.id = dod.product_id "+
				"INNER JOIN uoms as u ON u.id = dod.uom_id "+
				"INNER JOIN order_statuses as os ON os.id = dod.order_status_id "+
				"WHERE dod.deleted_at IS NULL AND dod.id = ?", ID)

			if err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			for query.Next() {
				err = query.Scan(&deliveryOrderDetail.ID, &deliveryOrderDetail.DeliveryOrderID, &deliveryOrderDetail.SoDetailID, &deliveryOrderDetail.BrandID, &deliveryOrderDetail.ProductID, &deliveryOrderDetail.UomID, &deliveryOrderDetail.OrderStatusID, &deliveryOrderDetail.DoDetailCode, &deliveryOrderDetail.Qty, &deliveryOrderDetail.Note, &deliveryOrderDetail.IsDoneSyncToEs, &deliveryOrderDetail.StartDateSyncToEs, &deliveryOrderDetail.EndDateSyncToEs, &deliveryOrderDetail.CreatedAt, &deliveryOrderDetail.UpdatedAt)

				if err != nil {
					errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
					response.Error = err
					response.ErrorLog = errorLogData
					resultChan <- response
					return
				}
			}

			deliveryOrderDetailsJson, _ := json.Marshal(deliveryOrderDetail)
			setDeliveryOrderDetailsOnRedis := r.redisdb.Client().Set(ctx, deliveryOrderDetailsRedis, deliveryOrderDetailsJson, 1*time.Hour)

			if setDeliveryOrderDetailsOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setDeliveryOrderDetailsOnRedis.Err(), http.StatusInternalServerError, nil)
				response.Error = setDeliveryOrderDetailsOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.DeliveryOrderDetail = &deliveryOrderDetail
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
		_ = json.Unmarshal([]byte(deliveryOrderDetailsRedis), &deliveryOrderDetail)
		response.DeliveryOrderDetail = &deliveryOrderDetail
		resultChan <- response
		return
	}
}

func (r *deliveryOrderDetail) GetByDeliveryOrderID(deliveryOrderID int, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderDetailsChan) {
	response := &models.DeliveryOrderDetailsChan{}
	var deliveryOrderDetails []*models.DeliveryOrderDetail
	var total int64

	deliveryOrderDetailRedisKey := fmt.Sprintf("%s:%d", constants.DELIVERY_ORDER_DETAIL_BY_DELIVERY_ORDER_ID, deliveryOrderID)
	deliveryOrderDetailsRedis, err := r.redisdb.Client().Get(ctx, deliveryOrderDetailRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM delivery_order_details WHERE deleted_at IS NULL AND delivery_order_id = ?", deliveryOrderID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("delivery_order_detail data not found")
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			query, err := r.db.Query(""+
				"SELECT dod.id, dod.delivery_order_id, dod.so_detail_id, dod.brand_id, dod.product_id, dod.uom_id, dod.order_status_id, dod.do_detail_code, dod.qty, dod.note, dod.is_done_sync_to_es, dod.start_date_sync_to_es, dod.end_date_sync_to_es, dod.created_at, dod.updated_at "+
				"FROM delivery_order_details as dod "+
				"INNER JOIN sales_order_details as sod ON sod.id = dod.so_detail_id "+
				"INNER JOIN brands as b ON b.id = dod.brand_id "+
				"INNER JOIN products as p ON p.id = dod.product_id "+
				"INNER JOIN uoms as u ON u.id = dod.uom_id "+
				"INNER JOIN order_statuses as os ON os.id = dod.order_status_id "+
				"WHERE dod.deleted_at IS NULL AND dod.delivery_order_id = ?", deliveryOrderID)

			if err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			for query.Next() {
				var deliveryOrderDetail models.DeliveryOrderDetail
				err = query.Scan(&deliveryOrderDetail.ID, &deliveryOrderDetail.DeliveryOrderID, &deliveryOrderDetail.SoDetailID, &deliveryOrderDetail.BrandID, &deliveryOrderDetail.ProductID, &deliveryOrderDetail.UomID, &deliveryOrderDetail.OrderStatusID, &deliveryOrderDetail.DoDetailCode, &deliveryOrderDetail.Qty, &deliveryOrderDetail.Note, &deliveryOrderDetail.IsDoneSyncToEs, &deliveryOrderDetail.StartDateSyncToEs, &deliveryOrderDetail.EndDateSyncToEs, &deliveryOrderDetail.CreatedAt, &deliveryOrderDetail.UpdatedAt)

				if err != nil {
					errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
					response.Error = err
					response.ErrorLog = errorLogData
					resultChan <- response
					return
				}

				deliveryOrderDetails = append(deliveryOrderDetails, &deliveryOrderDetail)
			}

			deliveryOrderDetailsJson, _ := json.Marshal(deliveryOrderDetails)
			setDeliveryOrderDetailsOnRedis := r.redisdb.Client().Set(ctx, deliveryOrderDetailsRedis, deliveryOrderDetailsJson, 1*time.Hour)

			if setDeliveryOrderDetailsOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setDeliveryOrderDetailsOnRedis.Err(), http.StatusInternalServerError, nil)
				response.Error = setDeliveryOrderDetailsOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.DeliveryOrderDetails = deliveryOrderDetails
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
		_ = json.Unmarshal([]byte(deliveryOrderDetailsRedis), &deliveryOrderDetails)
		response.DeliveryOrderDetails = deliveryOrderDetails
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *deliveryOrderDetail) GetBySalesOrderID(salesOrderID int, countOnly bool, ctx context.Context, resultChan chan *models.SalesOrderDetailsChan) {
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
			err = helper.NewError("delivery_order data not found")
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			query, err := r.db.Query(""+
				"SELECT p.SKU as product_sku, p.productName as product_name, sod.qty, sod.sent_qty, sod.residual_qty, u.code as uom_code, sod.price, os.name as order_detail_status "+
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
				err = query.Scan(&salesOrderDetail.ProductSKU, &salesOrderDetail.ProductName, &salesOrderDetail.Qty, &salesOrderDetail.SentQty, &salesOrderDetail.ResidualQty, &salesOrderDetail.UomCode, &salesOrderDetail.Price, &salesOrderDetail.OrderStatusName)

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

func (r *deliveryOrderDetail) Insert(request *models.DeliveryOrderDetail, sqlTransaction *sql.Tx, ctx context.Context, resultChan chan *models.DeliveryOrderDetailChan) {
	response := &models.DeliveryOrderDetailChan{}

	rawSqlFields := []string{}
	rawSqlDataTypes := []string{}
	rawSqlValues := []interface{}{}

	if request.DeliveryOrderID != 0 {
		rawSqlFields = append(rawSqlFields, "delivery_order_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.DeliveryOrderID)
	}

	if request.SoDetailID != 0 {
		rawSqlFields = append(rawSqlFields, "so_detail_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.SoDetailID)
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

	if request.DoDetailCode != "" {
		rawSqlFields = append(rawSqlFields, "do_detail_code")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.DoDetailCode)
	}

	rawSqlFields = append(rawSqlFields, "qty")
	rawSqlDataTypes = append(rawSqlDataTypes, "?")
	rawSqlValues = append(rawSqlValues, request.Qty)

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
		rawSqlValues = append(rawSqlValues, request.StartDateSyncToEs.Format(constants.DATE_TIME_FORMAT_COMON))
	}

	if request.EndDateSyncToEs != nil {
		rawSqlFields = append(rawSqlFields, "end_date_sync_to_es")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.StartDateSyncToEs.Format(constants.DATE_TIME_FORMAT_COMON))
	}

	rawSqlFields = append(rawSqlFields, "created_at")
	rawSqlDataTypes = append(rawSqlDataTypes, "?")
	rawSqlValues = append(rawSqlValues, request.CreatedAt.Format(constants.DATE_TIME_FORMAT_COMON))

	rawSqlFields = append(rawSqlFields, "updated_at")
	rawSqlDataTypes = append(rawSqlDataTypes, "?")
	rawSqlValues = append(rawSqlValues, request.CreatedAt.Format(constants.DATE_TIME_FORMAT_COMON))

	rawSqlFieldsJoin := strings.Join(rawSqlFields, ",")
	rawSqlDataTypesJoin := strings.Join(rawSqlDataTypes, ",")

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%v)", constants.DELIVERY_ORDER_DETAILS_TABLE, rawSqlFieldsJoin, rawSqlDataTypesJoin)
	result, err := sqlTransaction.ExecContext(ctx, query, rawSqlValues...)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	deliveryOrderDetailID, err := result.LastInsertId()

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	response.ID = deliveryOrderDetailID
	request.ID = int(deliveryOrderDetailID)
	response.DeliveryOrderDetail = request
	resultChan <- response
	return
}

func (r *deliveryOrderDetail) UpdateByID(id int, request *models.DeliveryOrderDetail, sqlTransaction *sql.Tx, ctx context.Context, resultChan chan *models.DeliveryOrderDetailChan) {
	response := &models.DeliveryOrderDetailChan{}
	rawSqlQueries := []string{}

	if request.OrderStatusID != 0 {
		query := fmt.Sprintf("%s=%v", "order_status_id", request.OrderStatusID)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	query := fmt.Sprintf("%s=%v", "qty", request.Qty)
	rawSqlQueries = append(rawSqlQueries, query)

	if len(request.Note.String) > 0 {
		query := fmt.Sprintf("%s='%v'", "note", request.Note.String)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if len(request.IsDoneSyncToEs) > 0 {
		query := fmt.Sprintf("%s='%v'", "is_done_sync_to_es", request.IsDoneSyncToEs)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.StartDateSyncToEs != nil {
		query := fmt.Sprintf("%s='%v'", "start_date_sync_to_es", request.StartDateSyncToEs.Format(constants.DATE_TIME_FORMAT_COMON))
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.EndDateSyncToEs != nil {
		query := fmt.Sprintf("%s='%v'", "end_date_sync_to_es", request.EndDateSyncToEs.Format(constants.DATE_TIME_FORMAT_COMON))
		rawSqlQueries = append(rawSqlQueries, query)
	}

	query = fmt.Sprintf("%s='%v'", "updated_at", request.UpdatedAt.Format(constants.DATE_TIME_FORMAT_COMON))
	rawSqlQueries = append(rawSqlQueries, query)

	rawSqlQueriesJoin := strings.Join(rawSqlQueries, ",")

	updateQuery := fmt.Sprintf("UPDATE %s set %v WHERE id = ?", constants.DELIVERY_ORDER_DETAILS_TABLE, rawSqlQueriesJoin)
	result, err := sqlTransaction.ExecContext(ctx, updateQuery, id)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	deliveryOrderID, err := result.LastInsertId()

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	response.ID = deliveryOrderID
	response.DeliveryOrderDetail = request
	resultChan <- response
	return
}
func (r *deliveryOrderDetail) DeleteByID(request *models.DeliveryOrderDetail, sqlTransaction *sql.Tx, ctx context.Context, resultChan chan *models.DeliveryOrderDetailChan) {
	now := time.Now()
	request.DeletedAt = &now
	response := &models.DeliveryOrderDetailChan{}
	rawSqlQueries := []string{}

	query := fmt.Sprintf("%s=%v", "qty", "0")
	rawSqlQueries = append(rawSqlQueries, query)

	query = fmt.Sprintf("%s=%v", "order_status_id", "19")
	rawSqlQueries = append(rawSqlQueries, query)

	query = fmt.Sprintf("%s='%v'", "deleted_at", request.DeletedAt.Format(constants.DATE_TIME_FORMAT_COMON))
	rawSqlQueries = append(rawSqlQueries, query)

	query = fmt.Sprintf("%s='%v'", "is_done_sync_to_es", 0)
	rawSqlQueries = append(rawSqlQueries, query)

	rawSqlQueriesJoin := strings.Join(rawSqlQueries, ",")
	updateQuery := fmt.Sprintf("UPDATE %s set %v WHERE id = ?", constants.DELIVERY_ORDER_DETAILS_TABLE, rawSqlQueriesJoin)
	result, err := sqlTransaction.ExecContext(ctx, updateQuery, request.ID)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	deliveryOrderID, err := result.LastInsertId()

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	response.ID = deliveryOrderID
	response.DeliveryOrderDetail = request
	resultChan <- response
	return
}
