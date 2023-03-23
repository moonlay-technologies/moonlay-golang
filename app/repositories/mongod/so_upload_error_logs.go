package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/global/utils/helper"
	"order-service/global/utils/mongodb"
	"os"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SoUploadErrorLogsRepositoryInterface interface {
	Insert(request *models.SoUploadErrorLog, ctx context.Context, resultChan chan *models.SoUploadErrorLogChan)
	Get(request *models.GetSoUploadErrorLogsRequest, countOnly bool, ctx context.Context, resultChan chan *models.SoUploadErrorLogsChan)
}

type soUploadErrorLogsRepository struct {
	logger     log.Logger
	mongod     mongodb.MongoDBInterface
	collection string
}

func InitSoUploadErrorLogsRepositoryInterface(mongod mongodb.MongoDBInterface) SoUploadErrorLogsRepositoryInterface {
	return &soUploadErrorLogsRepository{
		mongod:     mongod,
		collection: constants.SO_UPLOAD_ERROR_TABLE_LOGS,
	}
}

func (r *soUploadErrorLogsRepository) Insert(request *models.SoUploadErrorLog, ctx context.Context, resultChan chan *models.SoUploadErrorLogChan) {
	response := &models.SoUploadErrorLogChan{}
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
	response.SoUploadErrorLog = request
	response.Error = nil
	resultChan <- response
	return
}

func (r *soUploadErrorLogsRepository) Get(request *models.GetSoUploadErrorLogsRequest, countOnly bool, ctx context.Context, resultChan chan *models.SoUploadErrorLogsChan) {
	response := &models.SoUploadErrorLogsChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	filter := bson.M{}
	sort := bson.M{}
	asc := 1
	desc := -1

	page := request.Page
	if page == 0 {
		page = 1
	}

	perPage := request.PerPage
	if perPage == 0 {
		perPage = 10
	}

	if request.SortField == "updated_at" {
		if request.SortValue == "asc" {
			sort = bson.M{
				"updated_at": asc,
			}
		} else {
			sort = bson.M{
				"updated_at": desc,
			}
		}
	} else {
		if request.SortValue == "asc" {
			sort = bson.M{
				"created_at": asc,
			}
		} else {
			sort = bson.M{
				"created_at": desc,
			}
		}
	}

	if request.Status != "" {
		filter["status"] = request.Status
	}

	if request.RequestID != "" {
		filter["request_id"] = request.RequestID
	}

	option := options.Find().SetSkip(int64((page - 1) * perPage)).SetLimit(int64(perPage)).SetSort(sort)
	total, err := collection.CountDocuments(ctx, filter)

	fmt.Println("total", total)
	a, _ := json.Marshal(filter)
	fmt.Println(string(a))

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	if total == 0 {
		err = helper.NewError(helper.DefaultStatusText[http.StatusNotFound])
		errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	if countOnly == false {
		soUploadErorLogs := []*models.SoUploadErrorLog{}
		cursor, errs := collection.Find(ctx, filter, option)

		if errs != nil {
			response.Error = err
			resultChan <- response
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var salesOrderLog *models.SoUploadErrorLog
			if err := cursor.Decode(&salesOrderLog); err != nil {
				response.Error = err
				resultChan <- response
				return
			}

			soUploadErorLogs = append(soUploadErorLogs, salesOrderLog)
		}

		response.SoUploadErrorLogs = soUploadErorLogs
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.SoUploadErrorLogs = nil
		response.Total = total
		resultChan <- response
		return
	}
}
