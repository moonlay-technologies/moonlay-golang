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

type ProductRepositoryInterface interface {
	GetByID(ID int, countOnly bool, ctx context.Context, result chan *models.ProductChan)
	GetBySKU(SKU string, countOnly bool, ctx context.Context, resultChan chan *models.ProductChan)
}

type product struct {
	db      dbresolver.DB
	redisdb redisdb.RedisInterface
}

func InitProductRepository(db dbresolver.DB, redisdb redisdb.RedisInterface) ProductRepositoryInterface {
	return &product{
		db:      db,
		redisdb: redisdb,
	}
}

func (r *product) GetByID(ID int, countOnly bool, ctx context.Context, resultChan chan *models.ProductChan) {
	response := &models.ProductChan{}
	var product models.Product
	var total int64

	productRedisKey := fmt.Sprintf("%s:%d", "product", ID)
	productOnRedis, err := r.redisdb.Client().Get(ctx, productRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM products WHERE deleted_at IS NULL AND id = ?", ID).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			errStr := fmt.Sprintf("product id %d data not found", ID)
			err = helper.NewError(errStr)
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			product = models.Product{}
			err = r.db.QueryRow(""+
				"SELECT id, SKU, aliasSKU, productName, description, category_id, unitMeasurementSmall, unitMeasurementMedium, unitMeasurementBig, priceSmall, priceMedium, priceBig, isActive, nettWeight  from products as p "+
				"WHERE p.deleted_at IS NULL AND p.id = ?", ID).
				Scan(&product.ID, &product.Sku, &product.AliasSku, &product.ProductName, &product.Description, &product.CategoryID, &product.UnitMeasurementSmall, &product.UnitMeasurementMedium, &product.UnitMeasurementBig, &product.PriceSmall, &product.PriceMedium, &product.PriceBig, &product.IsActive, &product.NettWeight)

			if err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			productJson, _ := json.Marshal(product)
			setProductOnRedis := r.redisdb.Client().Set(ctx, productRedisKey, productJson, 1*time.Hour)

			if setProductOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setProductOnRedis.Err(), http.StatusInternalServerError, nil)
				response.Error = setProductOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.Product = &product
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
		_ = json.Unmarshal([]byte(productOnRedis), &product)
		response.Product = &product
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *product) GetBySKU(SKU string, countOnly bool, ctx context.Context, resultChan chan *models.ProductChan) {
	response := &models.ProductChan{}
	var product models.Product
	var total int64

	productRedisKey := fmt.Sprintf("%s:%s", "product", SKU)
	productOnRedis, err := r.redisdb.Client().Get(ctx, productRedisKey).Result()

	if err == redis.Nil {
		err = r.db.QueryRow("SELECT COUNT(*) as total FROM products WHERE deleted_at IS NULL AND SKU = ?", SKU).Scan(&total)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if total == 0 {
			errStr := fmt.Sprintf("product id %s data not found", SKU)
			err = helper.NewError(errStr)
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		if countOnly == false {
			product = models.Product{}
			query := fmt.Sprintf("SELECT id, IF(SKU = '%s', SKU, aliasSKU) AS SKU, aliasSKU, IF(SKU = '%s', productName, aliasName) AS productName, description, category_id, unitMeasurementSmall, unitMeasurementMedium, unitMeasurementBig, priceSmall, priceMedium, priceBig, isActive, nettWeight FROM products AS p WHERE p.deleted_at IS NULL AND IF((SELECT COUNT(SKU) FROM products WHERE SKU = '%s'), p.SKU = '%s', p.aliasSku = '%s')", SKU, SKU, SKU, SKU, SKU)
			err = r.db.QueryRow(query).
				Scan(&product.ID, &product.Sku, &product.AliasSku, &product.ProductName, &product.Description, &product.CategoryID, &product.UnitMeasurementSmall, &product.UnitMeasurementMedium, &product.UnitMeasurementBig, &product.PriceSmall, &product.PriceMedium, &product.PriceBig, &product.IsActive, &product.NettWeight)

			if err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			productJson, _ := json.Marshal(product)
			setProductOnRedis := r.redisdb.Client().Set(ctx, productRedisKey, productJson, 1*time.Hour)

			if setProductOnRedis.Err() != nil {
				errorLogData := helper.WriteLog(setProductOnRedis.Err(), http.StatusInternalServerError, nil)
				response.Error = setProductOnRedis.Err()
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			response.Total = total
			response.Product = &product
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
		_ = json.Unmarshal([]byte(productOnRedis), &product)
		response.Product = &product
		response.Total = total
		resultChan <- response
		return
	}
}
