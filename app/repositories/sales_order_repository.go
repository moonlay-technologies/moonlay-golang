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

type SalesOrderRepositoryInterface interface {
	GetByID(id int, countOnly bool, ctx context.Context, result chan *models.SalesOrderChan)
	GetByCode(soCode string, countOnly bool, ctx context.Context, result chan *models.SalesOrderChan)
	GetByAgentRefCode(soRefCode string, agentID int, countOnly bool, ctx context.Context, result chan *models.SalesOrderChan)
	Insert(request *models.SalesOrder, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.SalesOrderChan)
	UpdateByID(id int, salesOrder *models.SalesOrder, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.SalesOrderChan)
	RemoveCacheByID(id int, ctx context.Context, result chan *models.SalesOrderChan)
	DeleteByID(salesOrder *models.SalesOrder, sqlTransaction *sql.Tx, ctx context.Context, result chan *models.SalesOrderChan)
}

type salesOrder struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitSalesOrderRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) SalesOrderRepositoryInterface {
	return &salesOrder{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *salesOrder) GetByID(id int, countOnly bool, ctx context.Context, resultChan chan *models.SalesOrderChan) {
	response := &models.SalesOrderChan{}
	var salesOrder models.SalesOrder
	var total int64

	salesOrderRedisKey := fmt.Sprintf("%s:%d", constants.SALES_ORDER, id)
	salesOrderOnRedis, err := r.redisdb.Client().Get(ctx, salesOrderRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM "+constants.SALES_ORDERS_TABLE+" WHERE deleted_at IS NULL AND id = ?", id).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("sales_order data not found")
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			salesOrder = models.SalesOrder{}
			err = r.db.QueryRow(""+
				"SELECT so.id, so.agent_id, a.name as agent_name, a.email as agent_email, a.province_id as agent_province_id, ap.name agent_province_name, a.city_id as agent_city_id, ac.name as agent_city_name, a.district_id as agent_district_id, ad.name as agent_district_name, a.village_id as agent_village_id, av.name as agent_village_name, a.address as agent_address, a.phone as agent_phone, a.main_mobile_phone as agent_mobile_phone, so.store_id, s.name as store_name, s.store_code as store_code, s.email as store_email, s.province_id as store_province_id, sp.name as store_province_name, s.city_id, sc.name as store_city_name, s.district_id, sd.name as store_district_name, s.village_id as store_village_id, sv.name as store_village_name, s.address as store_address, s.phone as store_phone, s.main_mobile_phone as store_mobile_phone, so.brand_id, b.name as brand_name, so.user_id, u.email as user_email, u.first_name as user_first_name, u.last_name as user_last_name, u.role_id as user_role_id, so.order_source_id, os.source_name as order_source_name, so.order_status_id, ost.name as order_status_name, IFNULL(salesmans.id,0) as salesman_id, salesmans.name as salesman_name, so.so_code, so.so_date, so_ref_code, so_ref_date, so.g_long, so.g_lat, so.note, so.internal_comment, so.total_amount, so.total_tonase, so.start_created_date, so.created_at, so.updated_at "+
				"FROM "+constants.SALES_ORDERS_TABLE+" as so "+
				"INNER JOIN agents as a ON a.id = so.agent_id "+
				"INNER JOIN stores as s ON s.id = so.store_id "+
				"INNER JOIN order_sources as os ON os.id = so.order_source_id "+
				"INNER JOIN order_statuses as ost ON ost.id = so.order_status_id "+
				"INNER JOIN brands as b ON b.id = so.brand_id "+
				"INNER JOIN users as u ON u.id = so.user_id "+
				"INNER JOIN provinces as ap ON ap.id = a.province_id "+
				"INNER JOIN cities as ac ON ac.id = a.city_id "+
				"INNER JOIN districts as ad ON ad.id = a.district_id "+
				"INNER JOIN villages as av ON av.id = a.village_id "+
				"INNER JOIN provinces as sp ON sp.id = s.province_id "+
				"INNER JOIN cities as sc ON sc.id = s.city_id "+
				"INNER JOIN districts as sd ON sd.id = s.district_id "+
				"INNER JOIN villages as sv ON sv.id = s.village_id "+
				"LEFT JOIN salesmans on salesmans.email = u.email "+
				"WHERE so.deleted_at IS NULL AND so.id = ?", id).
				Scan(&salesOrder.ID, &salesOrder.AgentID, &salesOrder.AgentName, &salesOrder.AgentEmail, &salesOrder.AgentProvinceID, &salesOrder.AgentProvinceName, &salesOrder.AgentCityID, &salesOrder.AgentCityName, &salesOrder.AgentDistrictID, &salesOrder.AgentDistrictName, &salesOrder.AgentVillageID, &salesOrder.AgentVillageName, &salesOrder.AgentAddress, &salesOrder.AgentPhone, &salesOrder.AgentMainMobilePhone, &salesOrder.StoreID, &salesOrder.StoreName, &salesOrder.StoreCode, &salesOrder.StoreEmail, &salesOrder.StoreProvinceID, &salesOrder.StoreProvinceName, &salesOrder.StoreCityID, &salesOrder.StoreCityName, &salesOrder.StoreDistrictID, &salesOrder.StoreDistrictName, &salesOrder.StoreVillageID, &salesOrder.StoreVillageName, &salesOrder.StoreAddress, &salesOrder.StorePhone, &salesOrder.StoreMainMobilePhone, &salesOrder.BrandID, &salesOrder.BrandName, &salesOrder.UserID, &salesOrder.UserEmail, &salesOrder.UserFirstName, &salesOrder.UserLastName, &salesOrder.UserRoleID, &salesOrder.OrderSourceID, &salesOrder.OrderSourceName, &salesOrder.OrderStatusID, &salesOrder.OrderStatusName, &salesOrder.SalesmanID, &salesOrder.SalesmanName, &salesOrder.SoCode, &salesOrder.SoDate, &salesOrder.SoRefCode, &salesOrder.SoRefDate, &salesOrder.GLong, &salesOrder.GLat, &salesOrder.Note, &salesOrder.InternalComment, &salesOrder.TotalAmount, &salesOrder.TotalTonase, &salesOrder.StartCreatedDate, &salesOrder.CreatedAt, &salesOrder.UpdatedAt)

			if err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			salesOrderJson, _ := json.Marshal(salesOrder)
			setSalesOrderOnRedis := r.redisdb.Client().Set(ctx, salesOrderRedisKey, salesOrderJson, 1*time.Hour)

			if setSalesOrderOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setSalesOrderOnRedis.Err(), http.StatusInternalServerError, nil)
				response.Error = setSalesOrderOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.SalesOrder = &salesOrder
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
		_ = json.Unmarshal([]byte(salesOrderOnRedis), &salesOrder)
		response.SalesOrder = &salesOrder
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *salesOrder) GetByCode(soCode string, countOnly bool, ctx context.Context, resultChan chan *models.SalesOrderChan) {
	response := &models.SalesOrderChan{}
	var salesOrder models.SalesOrder
	var total int64

	salesOrderRedisKey := fmt.Sprintf("%s:%s", constants.SALES_ORDER_BY_CODE, soCode)
	salesOrderOnRedis, err := r.redisdb.Client().Get(ctx, salesOrderRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM "+constants.SALES_ORDERS_TABLE+" WHERE deleted_at IS NULL AND so_code = ?", soCode).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("sales_order data not found")
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			salesOrder = models.SalesOrder{}
			err = r.db.QueryRow(""+
				"SELECT so.id, so.agent_id, a.name as agent_name, a.email as agent_email, a.province_id as agent_province_id, ap.name agent_province_name, a.city_id as agent_city_id, ac.name as agent_city_name, a.district_id as agent_district_id, ad.name as agent_district_name, a.village_id as agent_village_id, av.name as agent_village_name, a.address as agent_address, a.phone as agent_phone, a.main_mobile_phone as agent_mobile_phone, so.store_id, s.name as store_name, s.store_code as store_code, s.email as store_email, s.province_id as store_province_id, sp.name as store_province_name, s.city_id, sc.name as store_city_name, s.district_id, sd.name as store_district_name, s.village_id as store_village_id, sv.name as store_village_name, s.address as store_address, s.phone as store_phone, s.main_mobile_phone as store_mobile_phone, so.brand_id, b.name as brand_name, so.user_id, u.email as user_email, u.first_name as user_first_name, u.last_name as user_last_name, u.role_id as user_role_id, so.order_source_id, os.source_name as order_source_name, so.order_status_id, ost.name as order_status_name, IFNULL(salesmans.id,0) as salesman_id, salesmans.name as salesman_name, so.so_code, so.so_date, so_ref_code, so_ref_date, so.g_long, so.g_lat, so.note, so.internal_comment, so.total_amount, so.total_tonase, so.created_at, so.updated_at "+
				"FROM "+constants.SALES_ORDERS_TABLE+" as so "+
				"INNER JOIN agents as a ON a.id = so.agent_id "+
				"INNER JOIN stores as s ON s.id = so.store_id "+
				"INNER JOIN order_sources as os ON os.id = so.order_source_id "+
				"INNER JOIN order_statuses as ost ON ost.id = so.order_status_id "+
				"INNER JOIN brands as b ON b.id = so.brand_id "+
				"INNER JOIN users as u ON u.id = so.user_id "+
				"INNER JOIN provinces as ap ON ap.id = a.province_id "+
				"INNER JOIN cities as ac ON ac.id = a.city_id "+
				"INNER JOIN districts as ad ON ad.id = a.district_id "+
				"INNER JOIN villages as av ON av.id = a.village_id "+
				"INNER JOIN provinces as sp ON sp.id = s.province_id "+
				"INNER JOIN cities as sc ON sc.id = s.city_id "+
				"INNER JOIN districts as sd ON sd.id = s.district_id "+
				"INNER JOIN villages as sv ON sv.id = s.village_id "+
				"LEFT JOIN salesmans on salesmans.email = u.email "+
				"WHERE so.deleted_at IS NULL AND so.so_code = ?", soCode).
				Scan(&salesOrder.ID, &salesOrder.AgentID, &salesOrder.AgentName, &salesOrder.AgentEmail, &salesOrder.AgentProvinceID, &salesOrder.AgentProvinceName, &salesOrder.AgentCityID, &salesOrder.AgentCityName, &salesOrder.AgentDistrictID, &salesOrder.AgentDistrictName, &salesOrder.AgentVillageID, &salesOrder.AgentVillageName, &salesOrder.AgentAddress, &salesOrder.AgentPhone, &salesOrder.AgentMainMobilePhone, &salesOrder.StoreID, &salesOrder.StoreName, &salesOrder.StoreCode, &salesOrder.StoreEmail, &salesOrder.StoreProvinceID, &salesOrder.StoreProvinceName, &salesOrder.StoreCityID, &salesOrder.StoreCityName, &salesOrder.StoreDistrictID, &salesOrder.StoreDistrictName, &salesOrder.StoreVillageID, &salesOrder.StoreVillageName, &salesOrder.StoreAddress, &salesOrder.StorePhone, &salesOrder.StoreMainMobilePhone, &salesOrder.BrandID, &salesOrder.BrandName, &salesOrder.UserID, &salesOrder.UserEmail, &salesOrder.UserFirstName, &salesOrder.UserLastName, &salesOrder.UserRoleID, &salesOrder.OrderSourceID, &salesOrder.OrderSourceName, &salesOrder.OrderStatusID, &salesOrder.OrderStatusName, &salesOrder.SalesmanID, &salesOrder.SalesmanName, &salesOrder.SoCode, &salesOrder.SoDate, &salesOrder.SoRefCode, &salesOrder.SoRefDate, &salesOrder.GLong, &salesOrder.GLat, &salesOrder.Note, &salesOrder.InternalComment, &salesOrder.TotalAmount, &salesOrder.TotalTonase, &salesOrder.CreatedAt, &salesOrder.UpdatedAt)

			if err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			salesOrderJson, _ := json.Marshal(salesOrder)
			setSalesOrderOnRedis := r.redisdb.Client().Set(ctx, salesOrderRedisKey, salesOrderJson, 1*time.Hour)

			if setSalesOrderOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setSalesOrderOnRedis.Err(), http.StatusInternalServerError, nil)
				response.Error = setSalesOrderOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.SalesOrder = &salesOrder
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
		_ = json.Unmarshal([]byte(salesOrderOnRedis), &salesOrder)
		response.SalesOrder = &salesOrder
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *salesOrder) GetByAgentRefCode(soRefCode string, agentID int, countOnly bool, ctx context.Context, resultChan chan *models.SalesOrderChan) {
	response := &models.SalesOrderChan{}
	var salesOrder models.SalesOrder
	var total int64

	salesOrderRedisKey := fmt.Sprintf("%s:%s:%d", constants.SALES_ORDER_BY_AGENT_REF_CODE, soRefCode, agentID)
	salesOrderOnRedis, err := r.redisdb.Client().Get(ctx, salesOrderRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM "+constants.SALES_ORDERS_TABLE+" WHERE deleted_at IS NULL AND so_ref_code = ? AND agent_id= ?", soRefCode, agentID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("sales_order data not found")
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			salesOrder = models.SalesOrder{}
			err = r.db.QueryRow(""+
				"SELECT so.id, so.agent_id, a.name as agent_name, a.email as agent_email, a.province_id as agent_province_id, ap.name agent_province_name, a.city_id as agent_city_id, ac.name as agent_city_name, a.district_id as agent_district_id, ad.name as agent_district_name, a.village_id as agent_village_id, av.name as agent_village_name, a.address as agent_address, a.phone as agent_phone, a.main_mobile_phone as agent_mobile_phone, so.store_id, s.name as store_name, s.store_code as store_code, s.email as store_email, s.province_id as store_province_id, sp.name as store_province_name, s.city_id, sc.name as store_city_name, s.district_id, sd.name as store_district_name, s.village_id as store_village_id, sv.name as store_village_name, s.address as store_address, s.phone as store_phone, s.main_mobile_phone as store_mobile_phone, so.brand_id, b.name as brand_name, so.user_id, u.email as user_email, u.first_name as user_first_name, u.last_name as user_last_name, u.role_id as user_role_id, so.order_source_id, os.source_name as order_source_name, so.order_status_id, ost.name as order_status_name, salesmans.id as salesman_id, salesmans.name as salesman_name, so.so_code, so.so_date, so_ref_code, so_ref_date, so.g_long, so.g_lat, so.note, so.internal_comment, so.total_amount, so.total_tonase, so.created_at, so.updated_at "+
				"FROM "+constants.SALES_ORDERS_TABLE+" as so "+
				"INNER JOIN agents as a ON a.id = so.agent_id "+
				"INNER JOIN stores as s ON s.id = so.store_id "+
				"INNER JOIN order_sources as os ON os.id = so.order_source_id "+
				"INNER JOIN order_statuses as ost ON ost.id = so.order_status_id "+
				"INNER JOIN brands as b ON b.id = so.brand_id "+
				"INNER JOIN users as u ON u.id = so.user_id "+
				"INNER JOIN provinces as ap ON ap.id = a.province_id "+
				"INNER JOIN cities as ac ON ac.id = a.city_id "+
				"INNER JOIN districts as ad ON ad.id = a.district_id "+
				"INNER JOIN villages as av ON av.id = a.village_id "+
				"INNER JOIN provinces as sp ON sp.id = s.province_id "+
				"INNER JOIN cities as sc ON sc.id = s.city_id "+
				"INNER JOIN districts as sd ON sd.id = s.district_id "+
				"INNER JOIN villages as sv ON sv.id = s.village_id "+
				"LEFT JOIN salesmans on salesmans.email = u.email "+
				"WHERE so.deleted_at IS NULL AND so.so_ref_code = ? AND so.agent_id= ?", soRefCode, agentID).
				Scan(&salesOrder.ID, &salesOrder.AgentID, &salesOrder.AgentName, &salesOrder.AgentEmail, &salesOrder.AgentProvinceID, &salesOrder.AgentProvinceName, &salesOrder.AgentCityID, &salesOrder.AgentCityName, &salesOrder.AgentDistrictID, &salesOrder.AgentDistrictName, &salesOrder.AgentVillageID, &salesOrder.AgentVillageName, &salesOrder.AgentAddress, &salesOrder.AgentPhone, &salesOrder.AgentMainMobilePhone, &salesOrder.StoreID, &salesOrder.StoreName, &salesOrder.StoreCode, &salesOrder.StoreEmail, &salesOrder.StoreProvinceID, &salesOrder.StoreProvinceName, &salesOrder.StoreCityID, &salesOrder.StoreCityName, &salesOrder.StoreDistrictID, &salesOrder.StoreDistrictName, &salesOrder.StoreVillageID, &salesOrder.StoreVillageName, &salesOrder.StoreAddress, &salesOrder.StorePhone, &salesOrder.StoreMainMobilePhone, &salesOrder.BrandID, &salesOrder.BrandName, &salesOrder.UserID, &salesOrder.UserEmail, &salesOrder.UserFirstName, &salesOrder.UserLastName, &salesOrder.UserRoleID, &salesOrder.OrderSourceID, &salesOrder.OrderSourceName, &salesOrder.OrderStatusID, &salesOrder.OrderStatusName, &salesOrder.SalesmanID, &salesOrder.SalesmanName, &salesOrder.SoCode, &salesOrder.SoDate, &salesOrder.SoRefCode, &salesOrder.SoRefDate, &salesOrder.GLong, &salesOrder.GLat, &salesOrder.Note, &salesOrder.InternalComment, &salesOrder.TotalAmount, &salesOrder.TotalTonase, &salesOrder.CreatedAt, &salesOrder.UpdatedAt)

			if err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			salesOrderJson, _ := json.Marshal(salesOrder)
			setSalesOrderOnRedis := r.redisdb.Client().Set(ctx, salesOrderRedisKey, salesOrderJson, 1*time.Hour)

			if setSalesOrderOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setSalesOrderOnRedis.Err(), http.StatusInternalServerError, nil)
				response.Error = setSalesOrderOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.SalesOrder = &salesOrder
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
		_ = json.Unmarshal([]byte(salesOrderOnRedis), &salesOrder)
		response.SalesOrder = &salesOrder
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *salesOrder) Insert(request *models.SalesOrder, sqlTransaction *sql.Tx, ctx context.Context, resultChan chan *models.SalesOrderChan) {
	response := &models.SalesOrderChan{}
	rawSqlFields := []string{}
	rawSqlDataTypes := []string{}
	rawSqlValues := []interface{}{}

	if request.CartID != 0 {
		rawSqlFields = append(rawSqlFields, "cart_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.CartID)
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

	if request.BrandID != 0 {
		rawSqlFields = append(rawSqlFields, "brand_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.BrandID)
	}

	if request.UserID != 0 {
		rawSqlFields = append(rawSqlFields, "user_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.UserID)
	}

	if request.VisitationID != 0 {
		rawSqlFields = append(rawSqlFields, "visitation_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.VisitationID)
	}

	if request.OrderSourceID != 0 {
		rawSqlFields = append(rawSqlFields, "order_source_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.OrderSourceID)
	}

	if request.OrderStatusID != 0 {
		rawSqlFields = append(rawSqlFields, "order_status_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.OrderStatusID)
	}

	if request.SoCode != "" {
		rawSqlFields = append(rawSqlFields, "so_code")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.SoCode)
	}

	if request.SoDate != "" {
		rawSqlFields = append(rawSqlFields, "so_date")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.SoDate)
	}

	if request.SoRefCode.String != "" {
		rawSqlFields = append(rawSqlFields, "so_ref_code")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.SoRefCode.String)
	}

	if request.SoRefDate.String != "" {
		rawSqlFields = append(rawSqlFields, "so_ref_date")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.SoRefDate.String)
	}

	if request.ReferralCode.String != "" {
		rawSqlFields = append(rawSqlFields, "referral_code")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.ReferralCode.String)
	}

	if request.GLat.Float64 != 0 {
		rawSqlFields = append(rawSqlFields, "g_lat")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.GLat)
	}

	if request.GLong.Float64 != 0 {
		rawSqlFields = append(rawSqlFields, "g_long")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.GLong)
	}

	if request.DeviceId.String != "" {
		rawSqlFields = append(rawSqlFields, "device_id")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.DeviceId.String)
	}

	if request.Note.String != "" {
		rawSqlFields = append(rawSqlFields, "note")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.Note.String)
	}

	if request.InternalComment.String != "" {
		rawSqlFields = append(rawSqlFields, "internal_comment")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.InternalComment.String)
	}

	if request.TotalAmount != 0 {
		rawSqlFields = append(rawSqlFields, "total_amount")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.TotalAmount)
	}

	if request.TotalTonase != 0 {
		rawSqlFields = append(rawSqlFields, "total_tonase")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.TotalTonase)
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

	if request.StartCreatedDate != nil {
		rawSqlFields = append(rawSqlFields, "start_created_date")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.StartCreatedDate.Format("2006-01-02 15:04:05"))
	}

	if request.CreatedBy != 0 {
		rawSqlFields = append(rawSqlFields, "created_by")
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		rawSqlValues = append(rawSqlValues, request.CreatedBy)
	}

	rawSqlFields = append(rawSqlFields, "created_at")
	rawSqlDataTypes = append(rawSqlDataTypes, "?")
	rawSqlValues = append(rawSqlValues, request.CreatedAt.Format("2006-01-02 15:04:05"))

	rawSqlFieldsJoin := strings.Join(rawSqlFields, ",")
	rawSqlDataTypesJoin := strings.Join(rawSqlDataTypes, ",")

	query := fmt.Sprintf("INSERT INTO "+constants.SALES_ORDERS_TABLE+" (%s) VALUES (%v)", rawSqlFieldsJoin, rawSqlDataTypesJoin)
	result, err := sqlTransaction.ExecContext(ctx, query, rawSqlValues...)

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

	response.ID = salesOrderID
	request.ID = int(salesOrderID)
	response.SalesOrder = request
	resultChan <- response
	return
}

func (r *salesOrder) UpdateByID(id int, request *models.SalesOrder, sqlTransaction *sql.Tx, ctx context.Context, resultChan chan *models.SalesOrderChan) {
	response := &models.SalesOrderChan{}
	rawSqlQueries := []string{}

	if request.AgentID != 0 {
		query := fmt.Sprintf("%s=%v", "agent_id", request.AgentID)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.UserID != 0 {
		query := fmt.Sprintf("%s=%v", "user_id", request.UserID)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.StoreID != 0 {
		query := fmt.Sprintf("%s=%v", "store_id", request.StoreID)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.OrderSourceID != 0 {
		query := fmt.Sprintf("%s=%v", "order_source_id", request.OrderSourceID)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.GLong.Float64 != 0 {
		query := fmt.Sprintf("%s=%v", "g_long", request.GLong.Float64)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.GLat.Float64 != 0 {
		query := fmt.Sprintf("%s=%v", "g_lat", request.GLat.Float64)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.Note.String != "" {
		query := fmt.Sprintf("%s='%v'", "note", request.Note.String)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.InternalComment.String != "" {
		query := fmt.Sprintf("%s='%v'", "internal_comment", request.InternalComment.String)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.SoDate != "" {
		query := fmt.Sprintf("%s='%v'", "so_date", request.SoDate)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.SoRefCode.String != "" {
		query := fmt.Sprintf("%s='%v'", "so_ref_code", request.SoRefCode.String)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.SoRefDate.String != "" {
		query := fmt.Sprintf("%s='%v'", "so_ref_date", request.SoRefDate.String)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.TotalAmount != 0 {
		query := fmt.Sprintf("%s=%v", "total_amount", request.TotalAmount)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.TotalTonase != 0 {
		query := fmt.Sprintf("%s=%v", "total_tonase", request.TotalTonase)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.DeviceId.String != "" {
		query := fmt.Sprintf("%s='%v'", "device_id", request.DeviceId.String)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.ReferralCode.String != "" {
		query := fmt.Sprintf("%s='%v'", "referral_code", request.ReferralCode.String)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.CartID != 0 {
		query := fmt.Sprintf("%s=%v", "cart_id", request.CartID)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.BrandID != 0 {
		query := fmt.Sprintf("%s=%v", "brand_id", request.BrandID)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.VisitationID != 0 {
		query := fmt.Sprintf("%s=%v", "visitation_id", request.VisitationID)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.OrderStatusID != 0 {
		query := fmt.Sprintf("%s=%v", "order_status_id", request.OrderStatusID)
		rawSqlQueries = append(rawSqlQueries, query)
	}

	if request.SoCode != "" {
		query := fmt.Sprintf("%s='%v'", "so_code", request.SoCode)
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

	if request.EndCreatedDate != nil {
		query := fmt.Sprintf("%s='%v'", "end_created_date", request.EndCreatedDate.Format("2006-01-02 15:04:05"))
		rawSqlQueries = append(rawSqlQueries, query)
	}

	query := fmt.Sprintf("%s='%v'", "updated_at", request.UpdatedAt.Format("2006-01-02 15:04:05"))
	rawSqlQueries = append(rawSqlQueries, query)

	rawSqlQueriesJoin := strings.Join(rawSqlQueries, ",")
	updateQuery := fmt.Sprintf("UPDATE "+constants.SALES_ORDERS_TABLE+" set %v WHERE id = ?", rawSqlQueriesJoin)
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

	response.ID = salesOrderID
	request.ID = int(salesOrderID)
	response.SalesOrder = request
	resultChan <- response
	return
}

func (r *salesOrder) RemoveCacheByID(id int, ctx context.Context, resultChan chan *models.SalesOrderChan) {
	response := &models.SalesOrderChan{}
	salesOrderRedisKey := fmt.Sprintf("%s:%d", constants.SALES_ORDER, id)
	result, err := r.redisdb.Client().Del(ctx, salesOrderRedisKey).Result()

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	fmt.Println(result)

	salesOrderDetailRedisKey := fmt.Sprintf("%s:%d", constants.SALES_ORDER_DETAIL_BY_SALES_ORDER_ID, id)
	result, err = r.redisdb.Client().Del(ctx, salesOrderDetailRedisKey).Result()

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	fmt.Println(result)

	response.Error = nil
	resultChan <- response
	return
}

func (r *salesOrder) DeleteByID(request *models.SalesOrder, sqlTransaction *sql.Tx, ctx context.Context, resultChan chan *models.SalesOrderChan) {
	now := time.Now()
	request.DeletedAt = &now
	request.UpdatedAt = &now
	response := &models.SalesOrderChan{}
	rawSqlQueries := []string{}

	query := fmt.Sprintf("%s='%v'", "deleted_at", request.DeletedAt.Format("2006-01-02 15:04:05"))
	rawSqlQueries = append(rawSqlQueries, query)

	query = fmt.Sprintf("%s='%v'", "updated_at", request.UpdatedAt.Format("2006-01-02 15:04:05"))
	rawSqlQueries = append(rawSqlQueries, query)

	query = fmt.Sprintf("%s='%v'", "is_done_sync_to_es", 0)
	rawSqlQueries = append(rawSqlQueries, query)

	rawSqlQueriesJoin := strings.Join(rawSqlQueries, ",")
	updateQuery := fmt.Sprintf("UPDATE "+constants.SALES_ORDERS_TABLE+" set %v WHERE id = ?", rawSqlQueriesJoin)
	result, err := sqlTransaction.ExecContext(ctx, updateQuery, request.ID)

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

	response.ID = salesOrderID
	response.SalesOrder = request
	resultChan <- response
	return
}
