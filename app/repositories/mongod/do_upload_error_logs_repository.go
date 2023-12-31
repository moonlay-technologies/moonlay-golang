package repositories

import (
	"context"
	"log"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/global/utils/helper"
	"order-service/global/utils/mongodb"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DoUploadErrorLogsRepositoryInterface interface {
	Insert(request *models.DoUploadErrorLog, ctx context.Context, resultChan chan *models.DoUploadErrorLogChan)
	Get(request *models.GetDoUploadErrorLogsRequest, countOnly bool, ctx context.Context, resultChan chan *models.DoUploadErrorLogsChan)
}

type doUploadErrorLogsRepository struct {
	logger     log.Logger
	mongod     mongodb.MongoDBInterface
	collection string
}

func InitDoUploadErrorLogsRepositoryInterface(mongod mongodb.MongoDBInterface) DoUploadErrorLogsRepositoryInterface {
	return &doUploadErrorLogsRepository{
		mongod:     mongod,
		collection: constants.SJ_UPLOAD_ERROR_TABLE_LOGS,
	}
}

func (r *doUploadErrorLogsRepository) Insert(request *models.DoUploadErrorLog, ctx context.Context, resultChan chan *models.DoUploadErrorLogChan) {
	response := &models.DoUploadErrorLogChan{}
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
	response.DoUploadErrorLog = request
	response.Error = nil
	resultChan <- response
	return
}

func (r *doUploadErrorLogsRepository) Get(request *models.GetDoUploadErrorLogsRequest, countOnly bool, ctx context.Context, resultChan chan *models.DoUploadErrorLogsChan) {
	response := &models.DoUploadErrorLogsChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	filter := bson.M{}
	sort := bson.M{}
	asc := 1
	desc := -1

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

	if request.DoUploadHistoryID != "" {
		doUploadHistoryID, err := primitive.ObjectIDFromHex(request.DoUploadHistoryID)
		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusBadRequest, "Ada kesalahan pada request data, silahkan dicek kembali")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		filter["do_upload_history_id"] = doUploadHistoryID
	}

	option := options.Find().SetSort(sort)
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
		doUploadErorLogs := []*models.DoUploadErrorLog{}
		cursor, errs := collection.Find(ctx, filter, option)

		if errs != nil {
			response.Error = err
			resultChan <- response
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var doUploadErrorLog *models.DoUploadErrorLog
			if err := cursor.Decode(&doUploadErrorLog); err != nil {
				response.Error = err
				resultChan <- response
				return
			}

			doUploadErorLogs = append(doUploadErorLogs, doUploadErrorLog)
		}

		response.DoUploadErrorLogs = doUploadErorLogs
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.DoUploadErrorLogs = nil
		response.Total = total
		resultChan <- response
		return
	}
}
