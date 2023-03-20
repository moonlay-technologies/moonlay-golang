package repositories

import (
	"context"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/global/utils/helper"
	"order-service/global/utils/mongodb"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SosjUploadErrorLogsRepositoryInterface interface {
	Insert(request *models.SosjUploadErrorLog, ctx context.Context, resultChan chan *models.SosjUploadErrorLogChan)
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
