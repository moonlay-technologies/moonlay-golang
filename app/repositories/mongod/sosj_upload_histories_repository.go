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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SOSJUploadHistoriesRepositoryInterface interface {
	Insert(request *models.UploadHistory, ctx context.Context, resultChan chan *models.UploadHistoryChan)
	GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.GetSosjUploadHistoryResponseChan)
	UpdateByID(ID string, request *models.UploadHistory, ctx context.Context, resultChan chan *models.UploadHistoryChan)
}

type sosjUploadHistoriesRepository struct {
	mongod     mongodb.MongoDBInterface
	collection string
}

func InitSOSJUploadHistoriesRepositoryInterface(mongod mongodb.MongoDBInterface) SOSJUploadHistoriesRepositoryInterface {
	return &sosjUploadHistoriesRepository{
		mongod:     mongod,
		collection: constants.SOSJ_UPLOAD_TABLE_HISTORIES,
	}
}

func (r *sosjUploadHistoriesRepository) Insert(request *models.UploadHistory, ctx context.Context, resultChan chan *models.UploadHistoryChan) {
	response := &models.UploadHistoryChan{}
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
	response.UploadHistory = request
	response.Error = nil
	resultChan <- response
}

func (r *sosjUploadHistoriesRepository) GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.GetSosjUploadHistoryResponseChan) {

	response := &models.GetSosjUploadHistoryResponseChan{}
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
		cursor, err := collection.Aggregate(ctx, bson.A{
			bson.D{{"$match", bson.D{{"_id", objectID}}}},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "sosj_upload_error_logs"},
						{"localField", "_id"},
						{"foreignField", "sosj_upload_history_id"},
						{"as", "sosj_upload_error_logs"},
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

		sosjUploadHistory := &models.GetSosjUploadHistoryResponse{}
		uploadHistory := &models.UploadHistory{}
		for cursor.Next(ctx) {
			if err := cursor.Decode(&sosjUploadHistory); err != nil {
				response.Error = err
				resultChan <- response
				return
			}
			if err := cursor.Decode(&uploadHistory); err != nil {
				response.Error = err
				resultChan <- response
				return
			}
		}
		sosjUploadHistory.GetSosjUploadHistoryResponseMap(uploadHistory)

		response.SosjUploadHistories = sosjUploadHistory
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.SosjUploadHistories = nil
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *sosjUploadHistoriesRepository) UpdateByID(ID string, request *models.UploadHistory, ctx context.Context, resultChan chan *models.UploadHistoryChan) {
	response := &models.UploadHistoryChan{}
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
	response.UploadHistory = request
	resultChan <- response
}
