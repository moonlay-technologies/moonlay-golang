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

type SalesOrderJourneysRepositoryInterface interface {
	Insert(request *models.SalesOrderJourneys, ctx context.Context, resultChan chan *models.SalesOrderJourneysChan)
}

type salesOrderJourneysRepository struct {
	logger     log.Logger
	mongod     mongodb.MongoDBInterface
	collection string
}

func InitSalesOrderJourneysRepository(mongod mongodb.MongoDBInterface) SalesOrderJourneysRepositoryInterface {
	return &salesOrderJourneysRepository{
		mongod:     mongod,
		collection: constants.SALES_ORDER_TABLE_JOURNEYS,
	}
}

func (r *salesOrderJourneysRepository) Insert(request *models.SalesOrderJourneys, ctx context.Context, resultChan chan *models.SalesOrderJourneysChan) {
	response := &models.SalesOrderJourneysChan{}
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
	response.SalesOrderJourneys = request
	response.Error = nil
	resultChan <- response
	return
}
