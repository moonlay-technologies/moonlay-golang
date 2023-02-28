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

type SalesOrderDetailJourneysRepositoryInterface interface {
	Insert(request *models.SalesOrderDetailJourneys, ctx context.Context, resultChan chan *models.SalesOrderDetailJourneysChan)
}

type salesOrderDetailJourneysRepository struct {
	logger     log.Logger
	mongod     mongodb.MongoDBInterface
	collection string
}

func InitSalesOrderDetailJourneysRepository(mongod mongodb.MongoDBInterface) SalesOrderDetailJourneysRepositoryInterface {
	return &salesOrderDetailJourneysRepository{
		mongod:     mongod,
		collection: constants.SALES_ORDER_DETAIL_TABLE_JOURNEYS,
	}
}

func (r *salesOrderDetailJourneysRepository) Insert(request *models.SalesOrderDetailJourneys, ctx context.Context, resultChan chan *models.SalesOrderDetailJourneysChan) {
	response := &models.SalesOrderDetailJourneysChan{}
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
	response.SalesOrderDetailJourneys = request
	response.Error = nil
	resultChan <- response
	return
}
