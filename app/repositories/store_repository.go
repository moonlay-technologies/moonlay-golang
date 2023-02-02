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

type StoreRepositoryInterface interface {
	GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.StoreChan)
}

type store struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitStoreRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) StoreRepositoryInterface {
	return &store{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *store) GetByID(ID int, countOnly bool, ctx context.Context, resultChan chan *models.StoreChan) {
	response := &models.StoreChan{}
	var store models.Store
	var total int64

	storeRedisKey := fmt.Sprintf("%s:%d", "store", ID)
	storeOnRedis, err := r.redisdb.Client().Get(ctx, storeRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM stores WHERE deleted_at IS NULL AND id = ?", ID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, 500, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			err = helper.NewError("store data not found")
			errorLogData := helper.WriteLog(err, 404, "data not found")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			store = models.Store{}
			err = r.db.QueryRow(""+
				"SELECT stores.id, store_code, stores_category, stores.name, email, email_verified, description, address, stores.province_id, stores.city_id, stores.district_id, stores.village_id, data_type, postal_code, g_lat, g_lng, contact_name, website, phone, main_mobile_phone, status, aliasName, DBOApprovalStatus, verified_dbo, verified_date, validation_store, channel, provinces.name as province_name, cities.name as city_name, districts.name as district_name, villages.name as village_name FROM stores  "+
				"INNER JOIN provinces on provinces.id = stores.province_id "+
				"INNER JOIN cities on cities.id = stores.city_id "+
				"INNER JOIN districts on districts.id = stores.district_id "+
				"INNEr JOIN villages on villages.id = stores.village_id "+
				"WHERE stores.deleted_at IS NULL AND stores.id = ?", ID).
				Scan(&store.ID, &store.StoreCode, &store.StoreCategory, &store.Name, &store.Email, &store.EmailVerified, &store.Description, &store.Address, &store.ProvinceID, &store.CityID, &store.DistrictID, &store.VillageID, &store.DataType, &store.PostalCode, &store.GLat, &store.GLng, &store.ContactName, &store.Website, &store.Phone, &store.MainMobilePhone, &store.Status, &store.AliasName, &store.DBOApprovalStatus, &store.VerifiedDBO, &store.VerifiedDate, &store.ValidationStore, &store.Channel, &store.ProvinceName, &store.CityName, &store.DistrictName, &store.VillageName)

			if err != nil {
				errorLogData := helper.WriteLog(err, 500, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			storeJson, _ := json.Marshal(store)
			setStoreOnRedis := r.redisdb.Client().Set(ctx, storeRedisKey, storeJson, 1*time.Hour)

			if setStoreOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setStoreOnRedis.Err(), 500, "Something went wrong, please try again later")
				response.Error = setStoreOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.Store = &store
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
		_ = json.Unmarshal([]byte(storeOnRedis), &store)
		response.Store = &store
		response.Total = total
		resultChan <- response
		return
	}
}
