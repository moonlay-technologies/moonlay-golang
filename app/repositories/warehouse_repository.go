package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/global/utils/helper"
	"order-service/global/utils/redisdb"
	"time"

	"github.com/bxcodec/dbresolver"
	"github.com/go-redis/redis/v8"
)

type WarehouseRepositoryInterface interface {
	GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.WarehouseChan)
}

type warehouse struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitWarehouseRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) WarehouseRepositoryInterface {
	return &warehouse{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *warehouse) GetByID(ID int, countOnly bool, ctx context.Context, resultChan chan *models.WarehouseChan) {
	response := &models.WarehouseChan{}
	var warehouse models.Warehouse
	var total int64

	warehouseRedisKey := fmt.Sprintf("%s:%d", "warehouse", ID)
	warehouseOnRedis, err := r.redisdb.Client().Get(ctx, warehouseRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM warehouses WHERE id = ?", ID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("warehouse data not found")
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			warehouse = models.Warehouse{}
			err = r.db.QueryRow(""+
				"SELECT w.id, w.code, w.name, owner_id, w.province_id, p.name as province_name, w.city_id, c.name as city_name, w.district_id, d.name as district_name, w.village_id, v.name as village_name, w.address, w.phone, w.main_mobile_phone, w.email, w.pic_name, w.status, w.warehouse_type_id, w.is_main FROM warehouses as w "+
				"INNER JOIN provinces as p on p.id = w.province_id "+
				"INNER JOIN cities as c on c.id = w.city_id "+
				"INNER JOIN districts as d on d.id = w.district_id "+
				"INNEr JOIN villages as v on v.id = w.village_id "+
				"WHERE  w.id = ?", ID).
				Scan(&warehouse.ID, &warehouse.Code, &warehouse.Name, &warehouse.OwnerID, &warehouse.ProvinceID, &warehouse.ProvinceName, &warehouse.CityID, &warehouse.CityName, &warehouse.DistrictID, &warehouse.DistrictName, &warehouse.VillageID, &warehouse.VillageName, &warehouse.Address, &warehouse.Phone, &warehouse.MainMobilePhone, &warehouse.Email, &warehouse.PicName, &warehouse.Status, &warehouse.WarehouseTypeID, &warehouse.IsMain)

			if err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			warehouseJson, _ := json.Marshal(warehouse)
			setWarehouseOnRedis := r.redisdb.Client().Set(ctx, warehouseRedisKey, warehouseJson, 1*time.Hour)

			if setWarehouseOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setWarehouseOnRedis.Err(), http.StatusInternalServerError, nil)
				response.Error = setWarehouseOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.Warehouse = &warehouse
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
		_ = json.Unmarshal([]byte(warehouseOnRedis), &warehouse)
		response.Warehouse = &warehouse
		response.Total = total
		resultChan <- response
		return
	}
}
