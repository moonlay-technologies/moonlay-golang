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

type UploadSOHistoriesRepositoryInterface interface {
	Insert(request *models.UploadSOHistory, ctx context.Context, resultChan chan *models.UploadSOHistoryChan)
}

type uploadSOHistoriesRepository struct {
	logger     log.Logger
	mongod     mongodb.MongoDBInterface
	collection string
}

func InitUploadSOHistoriesRepositoryInterface(mongod mongodb.MongoDBInterface) UploadSOHistoriesRepositoryInterface {
	return &uploadSOHistoriesRepository{
		mongod:     mongod,
		collection: constants.UPLOAD_SO_TABLE_HISTORIES,
	}
}

func (r *uploadSOHistoriesRepository) Insert(request *models.UploadSOHistory, ctx context.Context, resultChan chan *models.UploadSOHistoryChan) {
	response := &models.UploadSOHistoryChan{}
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
	response.UploadSOHistory = request
	response.Error = nil
	resultChan <- response
	return
}
