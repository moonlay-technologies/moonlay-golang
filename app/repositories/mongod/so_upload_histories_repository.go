package repositories

import (
	"context"
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

type SoUploadHistoriesRepositoryInterface interface {
	Insert(request *models.SoUploadHistory, ctx context.Context, resultChan chan *models.SoUploadHistoryChan)
	Get(request *models.GetSoUploadHistoriesRequest, countOnly bool, ctx context.Context, resultChan chan *models.SoUploadHistoriesChan)
	GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.SoUploadHistoryChan)
	GetByHistoryID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.GetSoUploadHistoryResponseChan)
	UpdateByID(ID string, request *models.SoUploadHistory, ctx context.Context, resultChan chan *models.SoUploadHistoryChan)
}

type soUploadHistoriesRepository struct {
	logger     log.Logger
	mongod     mongodb.MongoDBInterface
	collection string
}

func InitSoUploadHistoriesRepositoryInterface(mongod mongodb.MongoDBInterface) SoUploadHistoriesRepositoryInterface {
	return &soUploadHistoriesRepository{
		mongod:     mongod,
		collection: constants.SO_UPLOAD_TABLE_HISTORIES,
	}
}

func (r *soUploadHistoriesRepository) Insert(request *models.SoUploadHistory, ctx context.Context, resultChan chan *models.SoUploadHistoryChan) {
	response := &models.SoUploadHistoryChan{}
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
	response.SoUploadHistory = request
	response.Error = nil
	resultChan <- response
	return
}

func (r *soUploadHistoriesRepository) Get(request *models.GetSoUploadHistoriesRequest, countOnly bool, ctx context.Context, resultChan chan *models.SoUploadHistoriesChan) {
	response := &models.SoUploadHistoriesChan{}
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
		err = helper.NewError(helper.DefaultStatusText[http.StatusNotFound])
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	if !countOnly {
		SoUploadHistories := []*models.SoUploadHistory{}
		cursor, err := collection.Find(ctx, filter, option)

		if err != nil {
			fmt.Println(err.Error())
			response.Error = err
			resultChan <- response
			return
		}
		defer cursor.Close(ctx)

		response.SoUploadHistories = SoUploadHistories
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.SoUploadHistories = nil
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *soUploadHistoriesRepository) GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.SoUploadHistoryChan) {
	response := &models.SoUploadHistoryChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	objectID, _ := primitive.ObjectIDFromHex(ID)
	filter := bson.M{"_id": objectID}
	total, err := collection.CountDocuments(ctx, filter)

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

	if !countOnly {
		soUploadHistory := &models.SoUploadHistory{}
		err = collection.FindOne(ctx, filter).Decode(soUploadHistory)

		if err != nil {
			fmt.Println(err.Error())
			response.Error = err
			resultChan <- response
			return
		}

		response.SoUploadHistory = soUploadHistory
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.SoUploadHistory = nil
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *soUploadHistoriesRepository) GetByHistoryID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.GetSoUploadHistoryResponseChan) {
	response := &models.GetSoUploadHistoryResponseChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, "Ada kesalahan pada request data, silahkan dicek kembali")
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	filter := bson.M{"_id": objectID}
	total, err := collection.CountDocuments(ctx, filter)

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

	if !countOnly {
		cursor, err := collection.Aggregate(ctx, bson.A{
			bson.D{{"$match", bson.D{{"_id", objectID}}}},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "so_upload_error_logs"},
						{"localField", "_id"},
						{"foreignField", "so_upload_history_id"},
						{"as", "so_upload_error_logs"},
					},
				},
			},
		})

		if err != nil {
			fmt.Println(err.Error())
			response.Error = err
			resultChan <- response
			return
		}

		defer cursor.Close(ctx)
		soUploadHistoryResponse := &models.GetSoUploadHistoryResponse{}
		soUploadHistory := &models.SoUploadHistory{}
		for cursor.Next(ctx) {
			if err := cursor.Decode(&soUploadHistoryResponse); err != nil {
				response.Error = err
				resultChan <- response
				return
			}

			if err := cursor.Decode(&soUploadHistory); err != nil {
				response.Error = err
				resultChan <- response
				return
			}
		}

		soUploadHistoryResponse.GetSoUploadHistoryResponseMap(soUploadHistory)

		response.SoUploadHistories = soUploadHistoryResponse
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.SoUploadHistories = nil
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *soUploadHistoriesRepository) UpdateByID(ID string, request *models.SoUploadHistory, ctx context.Context, resultChan chan *models.SoUploadHistoryChan) {
	response := &models.SoUploadHistoryChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	objectID, _ := primitive.ObjectIDFromHex(ID)
	filter := bson.M{"_id": objectID}
	total, err := collection.CountDocuments(ctx, filter)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	if total == 0 {
		err = helper.NewError(helper.DefaultStatusText[http.StatusNotFound])
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	updateData := bson.M{"$set": request}
	_, err = collection.UpdateOne(r.mongod.GetCtx(), filter, updateData)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	response.Error = nil
	response.SoUploadHistory = request
	resultChan <- response
}
