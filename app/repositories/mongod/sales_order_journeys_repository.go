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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SalesOrderJourneysRepositoryInterface interface {
	Insert(request *models.SalesOrderJourneys, ctx context.Context, resultChan chan *models.SalesOrderJourneysChan)
	GetBySoId(ID int, countOnly bool, ctx context.Context, resultChan chan *models.SalesOrdersJourneysChan)
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

func (r *salesOrderJourneysRepository) GetBySoId(ID int, countOnly bool, ctx context.Context, resultChan chan *models.SalesOrdersJourneysChan) {
	response := &models.SalesOrdersJourneysChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	filter := bson.M{"so_id": ID}
	total, err := collection.CountDocuments(ctx, filter)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	if total == 0 {
		err = helper.NewError("data not found")
		errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	if countOnly == false {
		salesOrderJourneys := []*models.SalesOrderJourneys{}
		cursor, err := collection.Find(ctx, filter)

		if err != nil {
			response.Error = err
			resultChan <- response
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var salesOrderJourney *models.SalesOrderJourneys
			if err := cursor.Decode(&salesOrderJourney); err != nil {
				response.Error = err
				resultChan <- response
				return
			}

			salesOrderJourneys = append(salesOrderJourneys, salesOrderJourney)
		}

		response.SalesOrderJourneys = salesOrderJourneys
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.SalesOrderJourneys = nil
		response.Total = total
		resultChan <- response
		return
	}
}
