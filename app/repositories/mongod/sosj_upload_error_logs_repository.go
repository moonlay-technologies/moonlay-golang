package repositories

import (
	"context"
	"fmt"
	"math"
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

type SosjUploadErrorLogsRepositoryInterface interface {
	Insert(request *models.SosjUploadErrorLog, ctx context.Context, resultChan chan *models.SosjUploadErrorLogChan)
	Get(request *models.GetSosjUploadErrorLogsRequest, countOnly bool, ctx context.Context, resultChan chan *models.SosjUploadErrorLogsChan)
}

type sosjUploadErrorLogsRepository struct {
	mongod     mongodb.MongoDBInterface
	collection string
}

func InitSOSJUploadErrorLogsRepositoryInterface(mongod mongodb.MongoDBInterface) SosjUploadErrorLogsRepositoryInterface {
	return &sosjUploadErrorLogsRepository{
		mongod:     mongod,
		collection: constants.SOSJ_UPLOAD_ERROR_TABLE_LOGS,
	}
}

func (r *sosjUploadErrorLogsRepository) Insert(request *models.SosjUploadErrorLog, ctx context.Context, resultChan chan *models.SosjUploadErrorLogChan) {
	response := &models.SosjUploadErrorLogChan{}
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
	response.SosjUploadErrorLog = request
	response.Error = nil
	resultChan <- response
}

func (r *sosjUploadErrorLogsRepository) Get(request *models.GetSosjUploadErrorLogsRequest, countOnly bool, ctx context.Context, resultChan chan *models.SosjUploadErrorLogsChan) {
	response := &models.SosjUploadErrorLogsChan{}
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
		perPage = math.MaxInt
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

	if request.SoSjUploadHistoryID != "" {
		sosjUploadHistoryID, err := primitive.ObjectIDFromHex(request.SoSjUploadHistoryID)
		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusBadRequest, "Ada kesalahan pada request data, silahkan dicek kembali")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}
		filter["sosj_upload_history_id"] = sosjUploadHistoryID
	}

	option := options.Find().SetSkip(int64((page - 1) * perPage)).SetLimit(int64(perPage)).SetSort(sort)
	total, err := collection.CountDocuments(ctx, filter)
	fmt.Println("filter", filter)
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
		sosjUploadErorLogs := []*models.SosjUploadErrorLog{}
		cursor, errs := collection.Find(ctx, filter, option)

		if errs != nil {
			response.Error = err
			resultChan <- response
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var sosjUploadErrorLog *models.SosjUploadErrorLog
			if err := cursor.Decode(&sosjUploadErrorLog); err != nil {
				response.Error = err
				resultChan <- response
				return
			}

			sosjUploadErorLogs = append(sosjUploadErorLogs, sosjUploadErrorLog)
		}

		response.SosjUploadErrorLogs = sosjUploadErorLogs
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.SosjUploadErrorLogs = nil
		response.Total = total
		resultChan <- response
		return
	}
}
