package repositories

import (
	"context"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/global/utils/helper"
	"order-service/global/utils/mongodb"
	"os"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadSOSJHistoriesRepositoryInterface interface {
	Insert(request *models.UploadHistory, ctx context.Context, resultChan chan *models.UploadHistoryChan)
}

type uploadSOSJHistoriesRepository struct {
	logger     log.Logger
	mongod     mongodb.MongoDBInterface
	collection string
}

func InitUploadSOSJHistoriesRepositoryInterface(mongod mongodb.MongoDBInterface) UploadSOSJHistoriesRepositoryInterface {
	return &uploadSOSJHistoriesRepository{
		mongod:     mongod,
		collection: constants.UPLOAD_SOSJ_TABLE_HISTORIES,
	}
}

func (r *uploadSOSJHistoriesRepository) Insert(request *models.UploadHistory, ctx context.Context, resultChan chan *models.UploadHistoryChan) {
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
	return
}
