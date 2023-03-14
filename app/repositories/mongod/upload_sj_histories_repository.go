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

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadSJHistoriesRepositoryInterface interface {
	Insert(request *models.UploadSJHistory, ctx context.Context, resultChan chan *models.UploadSJHistoryChan)
}

type uploadSJHistoriesRepository struct {
	logger     log.Logger
	mongod     mongodb.MongoDBInterface
	collection string
}

func InitUploadSJHistoriesRepositoryInterface(mongod mongodb.MongoDBInterface) UploadSJHistoriesRepositoryInterface {
	return &uploadSJHistoriesRepository{
		mongod:     mongod,
		collection: constants.UPLOAD_DO_TABLE_HISTORIES,
	}
}

func (r *uploadSJHistoriesRepository) Insert(request *models.UploadSJHistory, ctx context.Context, resultChan chan *models.UploadSJHistoryChan) {
	response := &models.UploadSJHistoryChan{}
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
	response.UploadSJHistory = request
	response.Error = nil
	resultChan <- response
	return
}
