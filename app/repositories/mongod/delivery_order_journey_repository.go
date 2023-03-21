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
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeliveryOrderJourneysRepositoryInterface interface {
	Insert(request *models.DeliveryOrderJourney, ctx context.Context, result chan *models.DeliveryOrderJourneyChan)
	Get(request *models.DeliveryOrderJourneysRequest, countOnly bool, ctx context.Context, result chan *models.DeliveryOrderJourneysChan)
	GetByDoID(doID int, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderJourneysChan)
}

type deliveryOrderJourneysRepository struct {
	logger            log.Logger
	mongod            mongodb.MongoDBInterface
	collectionJourney string
}

func InitDeliveryOrderJourneysRepository(mongod mongodb.MongoDBInterface) DeliveryOrderJourneysRepositoryInterface {
	return &deliveryOrderJourneysRepository{
		mongod:            mongod,
		collectionJourney: constants.DELIVERY_ORDER_TABLE_JOURNEYS,
	}
}

func (r *deliveryOrderJourneysRepository) Insert(request *models.DeliveryOrderJourney, ctx context.Context, resultChan chan *models.DeliveryOrderJourneyChan) {
	response := &models.DeliveryOrderJourneyChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collectionJourney)
	result, err := collection.InsertOne(ctx, request)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	request.ID, _ = result.InsertedID.(primitive.ObjectID)
	response.DeliveryOrderJourney = request
	response.Error = nil
	resultChan <- response
	return
}

func (r *deliveryOrderJourneysRepository) Get(request *models.DeliveryOrderJourneysRequest, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderJourneysChan) {
	response := &models.DeliveryOrderJourneysChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collectionJourney)
	filter := bson.M{}

	if request.DoId > 0 {
		filter["do_id"] = request.DoId
	}

	if request.DoDate != "" {
		filter["do_date"] = request.DoDate
	}

	if request.Status != "" {
		filter["status"] = request.Status
	}

	if request.Remark != "" {
		filter["remark"] = request.Remark
	}

	if request.Reason != "" {
		filter["reason"] = request.Reason
	}

	if request.CreatedAt != "" {
		createdAt, _ := time.Parse("2006-01-02", request.CreatedAt)
		filter["created_at"] = &createdAt
	}

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

	if countOnly == false {
		deliveryOrderJourneys := []*models.DeliveryOrderJourney{}
		cursor, err := collection.Find(ctx, filter)

		if err != nil {
			response.Error = err
			resultChan <- response
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var deliveryOrderJourney *models.DeliveryOrderJourney
			if err := cursor.Decode(&deliveryOrderJourney); err != nil {
				response.Error = err
				resultChan <- response
				return
			}

			deliveryOrderJourneys = append(deliveryOrderJourneys, deliveryOrderJourney)
		}

		response.DeliveryOrderJourney = deliveryOrderJourneys
		response.Total = total
		response.Error = nil
		resultChan <- response
		return

	} else {
		response.DeliveryOrderJourney = nil
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *deliveryOrderJourneysRepository) GetByDoID(doID int, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderJourneysChan) {
	response := &models.DeliveryOrderJourneysChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collectionJourney)
	filter := bson.M{"do_id": doID}
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
		deliveryOrderJourneys := []*models.DeliveryOrderJourney{}
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			response.Error = err
			resultChan <- response
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			deliveryOrderJourney := &models.DeliveryOrderJourney{}
			if err := cursor.Decode(deliveryOrderJourney); err != nil {
				response.Error = err
				resultChan <- response
				return
			}
			deliveryOrderJourneys = append(deliveryOrderJourneys, deliveryOrderJourney)
		}
		response.DeliveryOrderJourney = deliveryOrderJourneys
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.DeliveryOrderJourney = nil
		response.Total = total
		resultChan <- response
		return
	}
}
