package repositories

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/global/utils/helper"
	"order-service/global/utils/model"
	"order-service/global/utils/mongodb"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeliveryOrderJourneyRepositoryInterface interface {
	InsertFromDO(request *models.DeliveryOrder, remarks string, now time.Time, ctx context.Context, result chan *models.DeliveryOrderJourneyChan)
	Get(request *models.DeliveryOrderJourneysRequest, countOnly bool, ctx context.Context, result chan *models.DeliveryOrderJourneysChan)
	GetByDoID(doID int, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderJourneysChan)
}

type deliveryOrderJourneyRepository struct {
	logger            log.Logger
	mongod            mongodb.MongoDBInterface
	collectionJourney string
}

func InitDeliveryOrderJourneyRepository(mongod mongodb.MongoDBInterface) DeliveryOrderJourneyRepositoryInterface {
	return &deliveryOrderJourneyRepository{
		mongod:            mongod,
		collectionJourney: constants.DELIVERY_ORDER_TABLE_JOURNEYS,
	}
}

func (r *deliveryOrderJourneyRepository) InsertFromDO(request *models.DeliveryOrder, remarks string, now time.Time, ctx context.Context, resultChan chan *models.DeliveryOrderJourneyChan) {
	deliveryOrderJourney := &models.DeliveryOrderJourney{
		DoId:      request.ID,
		DoCode:    request.DoCode,
		DoDate:    request.DoDate,
		Status:    helper.GetDOJourneyStatus(request.OrderStatusID),
		Remark:    remarks,
		Reason:    "",
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	r.Insert(deliveryOrderJourney, ctx, resultChan)
}

func (r *deliveryOrderJourneyRepository) Insert(request *models.DeliveryOrderJourney, ctx context.Context, resultChan chan *models.DeliveryOrderJourneyChan) {
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

func (r *deliveryOrderJourneyRepository) Get(request *models.DeliveryOrderJourneysRequest, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderJourneysChan) {
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
		createdAt, err := time.Parse(constants.DATE_FORMAT_COMMON, request.CreatedAt)
		if err != nil {
			errorLogData := helper.NewWriteLog(model.ErrorLog{
				Message:       "Ada kesalahan pada request data, silahkan dicek kembali",
				SystemMessage: helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_GET, fmt.Sprintf("field %s harus memiliki format yyyy-mm-dd", "created_at")),
				StatusCode:    http.StatusBadRequest,
				Err:           fmt.Errorf("invalid Process"),
			})
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}
		filter["created_at"] = bson.M{"$gte": createdAt, "$lte": createdAt.Add(24 * time.Hour)}
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
		err = helper.NewError(constants.ERROR_DATA_NOT_FOUND)
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

func (r *deliveryOrderJourneyRepository) GetByDoID(doID int, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderJourneysChan) {
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
		err = helper.NewError(constants.ERROR_DATA_NOT_FOUND)
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
