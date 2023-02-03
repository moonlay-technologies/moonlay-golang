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

type DeliveryOrderRepositoryInterface interface {
	GetBySalesOrderID(deliveryOrderID int, countOnly bool, ctx context.Context, result chan *models.DeliveryOrdersChan)
	Insert(request *models.DeliveryOrder, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.DeliveryOrderChan)
	GetByID(id int, countOnly bool, ctx context.Context, result chan *models.DeliveryOrderChan)
	//GetByAgentID(id int, countOnly bool, ctx context.Context, result chan *models.DeliveryOrderChan)
	UpdateByID(id int, deliveryOrder *models.DeliveryOrder, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.DeliveryOrderChan)
}

type deliveryOrder struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitDeliveryRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) DeliveryOrderRepositoryInterface {
	return &deliveryOrder{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *deliveryOrder) GetByID(id int, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderChan) {
	response := &models.DeliveryOrderChan{}
	var deliveryOrder models.DeliveryOrder
	var total int64

	deliveryOrderRedisKey := fmt.Sprintf("%s:%d", constants.DELIVERY_ORDER, id)
	deliveryOrderOnRedis, err := r.redisdb.Client().Get(ctx, deliveryOrderRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM delivery_orders WHERE deleted_at IS NULL AND id = ?", id).Scan(&total)

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
			deliveryOrder = models.DeliveryOrder{}
			err = r.db.QueryRow(""+
				"SELECT do.id, sales_order_id, warehouse_id, order_status_id, order_source_id, do_code, do_date, do_ref_code, do_ref_date, driver_name, plat_number, note, created_at, so.so_code, so.so_date, w.code, w.name, w.province_id, w.city_id, w.district_id, w.village_id, provinces.name as province_name, cities.name as city_name, districts.name as district_name, villages.name as village_name, order_statuses.name as order_status_name, order_sources.source_name as order_source_name "+
				"FROM delivery_orders as do "+
				"INNER JOIN "+constants.SALES_ORDERS_TABLE+" as so ON so.id = do.sales_order_id "+
				"INNER JOIN warehouses as w ON w.id = do.warehouse_id "+
				"INNER JOIN provinces ON provinces.id = w.province_id "+
				"INNER JOIN cities ON cities.province_id = provinces.id "+
				"INNER JOIN districts ON districts.city_id = cities.id "+
				"INNER JOIN villages ON villages.district_id = districts.id "+
				"INNER JOIN order_statuses ON order_statuses.id = do.order_status_id "+
				"INNER JOIN order_sources  ON order_sources.id = do.order_status_id "+
				"WHERE do.deleted_at IS NULL AND do.id = ?", id).
				Scan(&deliveryOrder.ID, &deliveryOrder.SalesOrderID, &deliveryOrder.WarehouseID, &deliveryOrder.OrderStatusID, &deliveryOrder.OrderSourceID, &deliveryOrder.DoCode, &deliveryOrder.DoDate, &deliveryOrder.DoRefCode, &deliveryOrder.DoRefDate, &deliveryOrder.DriverName, &deliveryOrder.PlatNumber, &deliveryOrder.Note, &deliveryOrder.CreatedAt, &deliveryOrder.SalesOrderCode, &deliveryOrder.SalesOrderDate, &deliveryOrder.WarehouseCode, &deliveryOrder.WarehouseName, &deliveryOrder.WarehouseProvinceID, &deliveryOrder.WarehouseCityID, &deliveryOrder.WarehouseDistrictID, &deliveryOrder.WarehouseVillageID, &deliveryOrder.WarehouseProvinceName, &deliveryOrder.WarehouseCityName, &deliveryOrder.WarehouseDistrictName, &deliveryOrder.WarehouseVillageName, &deliveryOrder.OrderStatusName, &deliveryOrder.OrderSourceName)

			if err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			deliveryOrderJson, _ := json.Marshal(deliveryOrder)
			setDeliveryOrderOnRedis := r.redisdb.Client().Set(ctx, deliveryOrderRedisKey, deliveryOrderJson, 1*time.Hour)

			if setDeliveryOrderOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setDeliveryOrderOnRedis.Err(), http.StatusInternalServerError, nil)
				response.Error = setDeliveryOrderOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.DeliveryOrder = &deliveryOrder
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
		_ = json.Unmarshal([]byte(deliveryOrderOnRedis), &deliveryOrder)
		response.DeliveryOrder = &deliveryOrder
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *deliveryOrder) GetBySalesOrderID(salesOrderID int, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrdersChan) {
	response := &models.DeliveryOrdersChan{}
	var deliveryOrdersResult *models.DeliveryOrders
	var total int64

	deliveryOrderRedisKey := fmt.Sprintf("%s:%d", constants.DELIVERY_ORDER_BY_SALES_ORDER, salesOrderID)
	deliveryOrderOnRedis, err := r.redisdb.Client().Get(ctx, deliveryOrderRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM delivery_orders WHERE deleted_at IS NULL AND sales_order_id = ?", salesOrderID).Scan(&total)

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
				"SELECT do.id, sales_order_id, warehouse_id, order_status_id, order_source_id, do_code, do_date, do_ref_code, do_ref_date, driver_name, plat_number, note, created_at, so.so_code, so.so_date, w.code, w.name, w.province_id, w.city_id, w.district_id, w.village_id, provinces.name as province_name, cities.name as city_name, districts.name as district_name, villages.name as village_name, order_statuses.name as order_status_name, order_sources.source_name as order_source_name "+
				"FROM delivery_orders as do "+
				"INNER JOIN "+constants.SALES_ORDERS_TABLE+" as so ON so.id = do.sales_order_id "+
				"INNER JOIN warehouses as w ON w.id = do.warehouse_id "+
				"INNER JOIN provinces ON provinces.id = w.province_id "+
				"INNER JOIN cities ON cities.province_id = provinces.id "+
				"INNER JOIN districts ON districts.city_id = cities.id "+
				"INNER JOIN villages ON villages.district_id = districts.id "+
				"INNER JOIN order_statuses ON order_statuses.id = do.order_status_id "+
				"INNER JOIN order_sources  ON order_sources.id = do.order_status_id "+
				"WHERE do.deleted_at IS NULL AND do.sales_order_id = ?", salesOrderID)

			if err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			deliveryOrders := []*models.DeliveryOrder{}

			for query.Next() {
				deliveryOrder := models.DeliveryOrder{}
				err = query.Scan(&deliveryOrder.ID, &deliveryOrder.SalesOrderID, &deliveryOrder.WarehouseID, &deliveryOrder.OrderStatusID, &deliveryOrder.OrderSourceID, &deliveryOrder.DoCode, &deliveryOrder.DoDate, &deliveryOrder.DoRefCode, &deliveryOrder.DoRefDate, &deliveryOrder.DriverName, &deliveryOrder.PlatNumber, &deliveryOrder.Note, &deliveryOrder.CreatedAt, &deliveryOrder.SalesOrderCode, &deliveryOrder.SalesOrderDate, &deliveryOrder.WarehouseCode, &deliveryOrder.WarehouseName, &deliveryOrder.WarehouseProvinceID, &deliveryOrder.WarehouseCityID, &deliveryOrder.WarehouseDistrictID, &deliveryOrder.WarehouseVillageID, &deliveryOrder.WarehouseProvinceName, &deliveryOrder.WarehouseCityName, &deliveryOrder.WarehouseDistrictName, &deliveryOrder.WarehouseVillageName, &deliveryOrder.OrderStatusName, &deliveryOrder.OrderSourceName)

				if err != nil {
					errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
					response.Error = err
					response.ErrorLog = errorLogData
					resultChan <- response
					return
				}

				deliveryOrders = append(deliveryOrders, &deliveryOrder)
			}

			deliveryOrdersResult = &models.DeliveryOrders{
				DeliveryOrders: deliveryOrders,
				Total:          total,
			}

			deliveryOrderJson, _ := json.Marshal(deliveryOrdersResult)
			setDeliveryOrderOnRedis := r.redisdb.Client().Set(ctx, deliveryOrderRedisKey, deliveryOrderJson, 1*time.Hour)

			if setDeliveryOrderOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setDeliveryOrderOnRedis.Err(), http.StatusInternalServerError, nil)
				response.Error = setDeliveryOrderOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.DeliveryOrders = deliveryOrders
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
		_ = json.Unmarshal([]byte(deliveryOrderOnRedis), &deliveryOrdersResult)
		response.DeliveryOrders = deliveryOrdersResult.DeliveryOrders
		response.Total = deliveryOrdersResult.Total
		resultChan <- response
		return
	}
}

func (r *deliveryOrder) Insert(request *models.DeliveryOrder, sqlTransaction *sql.Tx, ctx context.Context, resultChan chan *models.DeliveryOrderChan) {
	response := &models.DeliveryOrderChan{}
	rawSqlFields := []string{}
	rawSqlDataTypes := []string{}
	rawSqlValues := []interface{}{}

	if request.SalesOrderID != 0 {
		rawSqlFields = append(rawSqlFields, "sales_order_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.SalesOrderID)
	}

	if request.WarehouseID != 0 {
		rawSqlFields = append(rawSqlFields, "warehouse_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.WarehouseID)
	}

	if request.OrderStatusID != 0 {
		rawSqlFields = append(rawSqlFields, "order_status_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.OrderStatusID)
	}

	if request.AgentID != 0 {
		rawSqlFields = append(rawSqlFields, "agent_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.AgentID)
	}

	if request.StoreID != 0 {
		rawSqlFields = append(rawSqlFields, "store_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.StoreID)
	}

	if request.OrderSourceID != 0 {
		rawSqlFields = append(rawSqlFields, "order_source_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.OrderSourceID)
	}

	if request.DoCode != "" {
		rawSqlFields = append(rawSqlFields, "do_code")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.DoCode)
	}

	if request.DoDate != "" {
		rawSqlFields = append(rawSqlFields, "do_date")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.DoDate)
	}

	if request.DoRefDate.String != "" {
		rawSqlFields = append(rawSqlFields, "do_ref_code")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.DoRefCode)
	}

	if request.DoRefDate.String != "" {
		rawSqlFields = append(rawSqlFields, "do_ref_date")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.DoRefDate)
	}

	if request.PlatNumber.String != "" {
		rawSqlFields = append(rawSqlFields, "plat_number")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.PlatNumber)
	}

	if request.Note.String != "" {
		rawSqlFields = append(rawSqlFields, "note")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.Note.String)
	}

	if request.DriverName.String != "" {
		rawSqlFields = append(rawSqlFields, "driver_name")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.DriverName.String)
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

	query := fmt.Sprintf("INSERT INTO delivery_orders (%s) VALUES (%v)", rawSqlFieldsJoin, rawSqlDataTypesJoin)
	result, err := sqlTransaction.ExecContext(ctx, query, rawSqlValues...)

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
	request.ID = int(deliveryOrderID)
	response.DeliveryOrder = request
	resultChan <- response
	return
}

func (r *deliveryOrder) UpdateByID(id int, request *models.DeliveryOrder, sqlTransaction *sql.Tx, ctx context.Context, resultChan chan *models.DeliveryOrderChan) {
	response := &models.DeliveryOrderChan{}
	rawSqlQueries := []string{}

	if request.WarehouseID != 0 {
		query := fmt.Sprintf("%s=%v", "warehouse_id", request.WarehouseID)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.OrderStatusID != 0 {
		query := fmt.Sprintf("%s=%v", "order_status_id", request.OrderStatusID)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.OrderSourceID != 0 {
		query := fmt.Sprintf("%s=%v", "order_source_id", request.OrderSourceID)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.AgentID != 0 {
		query := fmt.Sprintf("%s=%v", "agent_id", request.StoreID)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.StoreID != 0 {
		query := fmt.Sprintf("%s=%v", "store_id", request.AgentID)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.DoCode != "" {
		query := fmt.Sprintf("%s='%v'", "do_code", request.DoCode)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.DoDate != "" {
		query := fmt.Sprintf("%s='%v'", "so_date", request.DoDate)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.DoRefCode.String != "" {
		query := fmt.Sprintf("%s='%v'", "do_ref_code", request.DoRefCode.String)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.DoRefDate.String != "" {
		query := fmt.Sprintf("%s='%v'", "do_ref_date", request.DoRefDate.String)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.DriverName.String != "" {
		query := fmt.Sprintf("%s='%v'", "driver_name", request.DriverName.String)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.PlatNumber.String != "" {
		query := fmt.Sprintf("%s='%v'", "plat_number", request.PlatNumber.String)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.Note.String != "" {
		query := fmt.Sprintf("%s='%v'", "note", request.Note.String)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.IsDoneSyncToEs != "" {
		query := fmt.Sprintf("%s='%v'", "is_done_sync_to_es", request.IsDoneSyncToEs)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.StartDateSyncToEs != nil {
		query := fmt.Sprintf("%s='%v'", "start_date_sync_to_es", request.StartDateSyncToEs.Format("2006-01-02 15:04:05"))
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.EndDateSyncToEs != nil {
		query := fmt.Sprintf("%s='%v'", "end_date_sync_to_es", request.EndDateSyncToEs.Format("2006-01-02 15:04:05"))
		rawSqlQueries = append(rawSqlQueries, query)
	}

	query := fmt.Sprintf("%s='%v'", "updated_at", request.UpdatedAt.Format("2006-01-02 15:04:05"))
	rawSqlQueries = append(rawSqlQueries, query)

	rawSqlQueriesJoin := strings.Join(rawSqlQueries, ",")
	fmt.Println(rawSqlQueriesJoin)
	updateQuery := fmt.Sprintf("UPDATE delivery_orders set %v WHERE id = ?", rawSqlQueriesJoin)
	result, err := sqlTransaction.ExecContext(ctx, updateQuery, id)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	salesOrderID, err := result.LastInsertId()

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	deliveryOrderRedisKey := fmt.Sprintf("%s", constants.DELIVERY_ORDER+"*")
	_, err = r.redisdb.Client().Del(ctx, deliveryOrderRedisKey).Result()

	response.ID = salesOrderID
	request.ID = int(salesOrderID)
	response.DeliveryOrder = request
	resultChan <- response
	return
}
