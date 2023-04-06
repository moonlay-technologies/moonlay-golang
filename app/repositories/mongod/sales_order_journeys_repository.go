package repositories

import (
	"context"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/global/utils/helper"
	"order-service/global/utils/model"
	"order-service/global/utils/mongodb"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SalesOrderJourneysRepositoryInterface interface {
	Insert(request *models.SalesOrderJourneys, ctx context.Context, resultChan chan *models.SalesOrderJourneysChan)
	Get(request *models.SalesOrderJourneyRequest, countOnly bool, ctx context.Context, result chan *models.SalesOrdersJourneysChan)
	GetBySoId(ID int, countOnly bool, ctx context.Context, resultChan chan *models.SalesOrdersJourneysChan)
}

type salesOrderJourneysRepository struct {
	logger     log.Logger
	mongod     mongodb.MongoDBInterface
	collection string
}

func InitSalesOrderJourneysRepository(mongod mongodb.MongoDBInterface) SalesOrderJourneysRepositoryInterface {
	return &salesOrderJourneysRepository{
		mongod:     mongod,
		collection: constants.SALES_ORDER_TABLE_JOURNEYS,
	}
}

func (r *salesOrderJourneysRepository) Insert(request *models.SalesOrderJourneys, ctx context.Context, resultChan chan *models.SalesOrderJourneysChan) {
	response := &models.SalesOrderJourneysChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	result, err := collection.InsertOne(ctx, request)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	request.ID, _ = result.InsertedID.(primitive.ObjectID)
	response.SalesOrderJourneys = request
	response.Error = nil
	resultChan <- response
	return
}

func (r *salesOrderJourneysRepository) Get(request *models.SalesOrderJourneyRequest, countOnly bool, ctx context.Context, resultChan chan *models.SalesOrdersJourneysChan) {
	response := &models.SalesOrdersJourneysChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	filter := bson.M{}
	sort := bson.M{}
	asc := 1
	desc := -1

	if request.SortField == "so_code" {
		if request.SortValue == "asc" {
			sort = bson.M{
				"so_code": asc,
			}
		} else if request.SortValue == "desc" {
			sort = bson.M{
				"so_code": desc,
			}
		}
	} else if request.SortField == "status" {
		if request.SortValue == "asc" {
			sort = bson.M{
				"status": asc,
			}
		} else if request.SortValue == "desc" {
			sort = bson.M{
				"status": desc,
			}
		}
	} else if request.SortField == "action" {
		if request.SortValue == "asc" {
			sort = bson.M{
				"action": asc,
			}
		} else if request.SortValue == "desc" {
			sort = bson.M{
				"action": desc,
			}
		}
	} else if request.SortField == "created_at" {
		if request.SortValue == "asc" {
			sort = bson.M{
				"created_at": asc,
			}
		} else if request.SortValue == "desc" {
			sort = bson.M{
				"created_at": desc,
			}
		}
	}

	if request.GlobalSearchValue != "" {
		soId, _ := strconv.ParseInt(request.GlobalSearchValue, 10, 32)
		createdAt, _ := time.Parse(constants.DATE_FORMAT_COMMON, request.GlobalSearchValue)
		filter = bson.M{
			"$or": []bson.M{
				{"status": bson.M{"$regex": request.GlobalSearchValue, "$options": "i"}},
				{"action": bson.M{"$regex": request.GlobalSearchValue, "$options": "i"}},
				{"created_at": bson.M{"$gte": createdAt, "$lte": createdAt.AddDate(0, 0, 1)}},
				{"so_id": soId},
			},
		}
	}

	if request.SoId > 0 {
		filter["so_id"] = request.SoId
	}

	if request.Status != "" {
		filter["status"] = request.Status
	}

	if request.Action != "" {
		filter["action"] = request.Action
	}

	if request.StartDate != "" && request.EndDate == "" {
		startDate, err := time.Parse(constants.DATE_FORMAT_COMMON, request.StartDate)
		if err != nil {
			errorLogData := helper.NewWriteLog(model.ErrorLog{
				Message:       "Ada kesalahan pada request data, silahkan dicek kembali",
				SystemMessage: helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_GET, fmt.Sprintf("field %s harus memiliki format yyyy-mm-dd", "start_date")),
				StatusCode:    http.StatusBadRequest,
				Err:           fmt.Errorf("invalid Process"),
			})
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}
		filter["created_at"] = bson.M{"$gte": startDate, "$lte": startDate.AddDate(0, 0, 1)}
	}

	if request.EndDate != "" && request.StartDate == "" {
		endDate, err := time.Parse(constants.DATE_FORMAT_COMMON, request.EndDate)
		if err != nil {
			errorLogData := helper.NewWriteLog(model.ErrorLog{
				Message:       "Ada kesalahan pada request data, silahkan dicek kembali",
				SystemMessage: helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_GET, fmt.Sprintf("field %s harus memiliki format yyyy-mm-dd", "end_date")),
				StatusCode:    http.StatusBadRequest,
				Err:           fmt.Errorf("invalid Process"),
			})
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}
		filter["created_at"] = bson.M{"$gte": endDate, "$lte": endDate.AddDate(0, 0, 1)}
	}

	if request.StartDate != "" && request.EndDate != "" {
		startDate, err := time.Parse(constants.DATE_FORMAT_COMMON, request.StartDate)
		if err != nil {
			errorLogData := helper.NewWriteLog(model.ErrorLog{
				Message:       "Ada kesalahan pada request data, silahkan dicek kembali",
				SystemMessage: helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_GET, fmt.Sprintf("field %s harus memiliki format yyyy-mm-dd", "start_date")),
				StatusCode:    http.StatusBadRequest,
				Err:           fmt.Errorf("invalid Process"),
			})
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		endDate, err := time.Parse(constants.DATE_FORMAT_COMMON, request.EndDate)
		if err != nil {
			errorLogData := helper.NewWriteLog(model.ErrorLog{
				Message:       "Ada kesalahan pada request data, silahkan dicek kembali",
				SystemMessage: helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_GET, fmt.Sprintf("field %s harus memiliki format yyyy-mm-dd", "end_date")),
				StatusCode:    http.StatusBadRequest,
				Err:           fmt.Errorf("invalid Process"),
			})
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		filter["created_at"] = bson.M{"$gte": startDate, "$lte": endDate.AddDate(0, 0, 1)}
	}

	option := options.Find().SetSkip(int64((request.Page - 1) * request.PerPage)).SetLimit(int64(request.PerPage)).SetSort(sort)
	total, err := collection.CountDocuments(ctx, filter)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	if total == 0 {
		err = helper.NewError(constants.ERROR_DATA_NOT_FOUND)
		errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	if countOnly == false {
		salesOrderJourneys := []*models.SalesOrderJourneys{}
		cursor, err := collection.Find(ctx, filter, option)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, "Terjadi Kesalahan, Silahkan Coba lagi Nanti")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var salesOrderJourney *models.SalesOrderJourneys
			if err := cursor.Decode(&salesOrderJourney); err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, "Terjadi Kesalahan, Silahkan Coba lagi Nanti")
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			salesOrderJourneys = append(salesOrderJourneys, salesOrderJourney)
		}

		response.SalesOrderJourneys = salesOrderJourneys
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.SalesOrderJourneys = nil
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *salesOrderJourneysRepository) GetBySoId(ID int, countOnly bool, ctx context.Context, resultChan chan *models.SalesOrdersJourneysChan) {
	response := &models.SalesOrdersJourneysChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	filter := bson.M{"so_id": ID}
	total, err := collection.CountDocuments(ctx, filter)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	if total == 0 {
		err = helper.NewError(constants.ERROR_DATA_NOT_FOUND)
		errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	if countOnly == false {
		salesOrderJourneys := []*models.SalesOrderJourneys{}
		cursor, err := collection.Find(ctx, filter)

		if err != nil {
			response.Error = err
			resultChan <- response
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var salesOrderJourney *models.SalesOrderJourneys
			if err := cursor.Decode(&salesOrderJourney); err != nil {
				response.Error = err
				resultChan <- response
				return
			}

			salesOrderJourneys = append(salesOrderJourneys, salesOrderJourney)
		}

		response.SalesOrderJourneys = salesOrderJourneys
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.SalesOrderJourneys = nil
		response.Total = total
		resultChan <- response
		return
	}
}
