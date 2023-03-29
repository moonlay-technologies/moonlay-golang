package repositories

import (
	"context"
	"fmt"
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

type DoUploadHistoriesRepositoryInterface interface {
	Insert(request *models.DoUploadHistory, ctx context.Context, resultChan chan *models.DoUploadHistoryChan)
	Get(request *models.GetDoUploadHistoriesRequest, countOnly bool, ctx context.Context, resultChan chan *models.DoUploadHistoriesChan)
	GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.DoUploadHistoryChan)
	GetByHistoryID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.GetDoUploadHistoryResponseChan)
	UpdateByID(ID string, request *models.DoUploadHistory, ctx context.Context, resultChan chan *models.DoUploadHistoryChan)
}

type doUploadHistoriesRepository struct {
	logger     log.Logger
	mongod     mongodb.MongoDBInterface
	collection string
}

func InitDoUploadHistoriesRepositoryInterface(mongod mongodb.MongoDBInterface) DoUploadHistoriesRepositoryInterface {
	return &doUploadHistoriesRepository{
		mongod:     mongod,
		collection: constants.SJ_UPLOAD_TABLE_HISTORIES,
	}
}

func (r *doUploadHistoriesRepository) Insert(request *models.DoUploadHistory, ctx context.Context, resultChan chan *models.DoUploadHistoryChan) {
	response := &models.DoUploadHistoryChan{}
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
	response.DoUploadHistory = request
	response.Error = nil
	resultChan <- response
	return
}

func (r *doUploadHistoriesRepository) Get(request *models.GetDoUploadHistoriesRequest, countOnly bool, ctx context.Context, resultChan chan *models.DoUploadHistoriesChan) {
	response := &models.DoUploadHistoriesChan{}
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
		doUploadHistories := []*models.DoUploadHistory{}
		cursor, err := collection.Find(ctx, filter, option)

		if err != nil {
			fmt.Println(err.Error())
			response.Error = err
			resultChan <- response
			return
		}
		defer cursor.Close(ctx)

		response.DoUploadHistories = doUploadHistories
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.DoUploadHistories = nil
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *doUploadHistoriesRepository) GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.DoUploadHistoryChan) {
	response := &models.DoUploadHistoryChan{}
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
		doUploadHistory := &models.DoUploadHistory{}
		err = collection.FindOne(ctx, filter).Decode(doUploadHistory)

		if err != nil {
			fmt.Println(err.Error())
			response.Error = err
			resultChan <- response
			return
		}

		response.DoUploadHistory = doUploadHistory
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.DoUploadHistory = nil
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *doUploadHistoriesRepository) GetByHistoryID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.GetDoUploadHistoryResponseChan) {
	response := &models.GetDoUploadHistoryResponseChan{}
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
						{"from", "do_upload_error_logs"},
						{"localField", "_id"},
						{"foreignField", "do_upload_history_id"},
						{"as", "do_upload_error_logs"},
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
		doUploadHistoryResponse := &models.GetDoUploadHistoryResponse{}
		doUploadHistory := &models.DoUploadHistory{}
		for cursor.Next(ctx) {
			if err := cursor.Decode(&doUploadHistoryResponse); err != nil {
				response.Error = err
				resultChan <- response
				return
			}

			if err := cursor.Decode(&doUploadHistory); err != nil {
				response.Error = err
				resultChan <- response
				return
			}
		}

		doUploadHistoryResponse.GetDoUploadHistoryResponseMap(doUploadHistory)

		response.DoUploadHistories = doUploadHistoryResponse
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.DoUploadHistories = nil
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *doUploadHistoriesRepository) UpdateByID(ID string, request *models.DoUploadHistory, ctx context.Context, resultChan chan *models.DoUploadHistoryChan) {
	response := &models.DoUploadHistoryChan{}
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
	response.DoUploadHistory = request
	resultChan <- response
}
