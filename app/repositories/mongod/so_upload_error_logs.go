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

type SoUploadErrorLogsRepositoryInterface interface {
	Insert(request *models.SoUploadErrorLog, ctx context.Context, resultChan chan *models.SoUploadErrorLogChan)
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
