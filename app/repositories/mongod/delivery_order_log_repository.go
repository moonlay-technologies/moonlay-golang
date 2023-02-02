package repositories

import (
	"context"
	"fmt"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/global/utils/helper"
	"order-service/global/utils/mongodb"
	"os"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeliveryOrderLogRepositoryInterface interface {
	Insert(request *models.DeliveryOrderLog, ctx context.Context, result chan *models.DeliveryOrderLogChan)
	GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderLogChan)
	UpdateByID(ID string, request *models.DeliveryOrderLog, ctx context.Context, result chan *models.DeliveryOrderLogChan)
}

type deliveryOrderLogRepository struct {
	logger     log.Logger
	mongod     mongodb.MongoDBInterface
	collection string
}

func InitDeliveryOrderLogRepository(mongod mongodb.MongoDBInterface) DeliveryOrderLogRepositoryInterface {
	return &deliveryOrderLogRepository{
		mongod:     mongod,
		collection: "delivery_order_logs",
	}
}

func (r *deliveryOrderLogRepository) GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderLogChan) {
	response := &models.DeliveryOrderLogChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	objectID, _ := primitive.ObjectIDFromHex(ID)
	filter := bson.M{"_id": objectID}
	total, err := collection.CountDocuments(ctx, filter)

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	if total == 0 {
		err = helper.NewError("data not found")
		errorLogData := helper.WriteLog(err, 404, "data not found")
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	if countOnly == false {
		deliveryOrderLog := &models.DeliveryOrderLog{}
		err = collection.FindOne(ctx, filter).Decode(deliveryOrderLog)

		if err != nil {
			fmt.Println(err.Error())
			response.Error = err
			resultChan <- response
			return
		}

		response.DeliveryOrderLog = deliveryOrderLog
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.DeliveryOrderLog = nil
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *deliveryOrderLogRepository) Insert(request *models.DeliveryOrderLog, ctx context.Context, resultChan chan *models.DeliveryOrderLogChan) {
	response := &models.DeliveryOrderLogChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	result, err := collection.InsertOne(ctx, request)

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	request.ID, _ = result.InsertedID.(primitive.ObjectID)
	response.DeliveryOrderLog = request
	response.Error = nil
	resultChan <- response
	return
}

func (r *deliveryOrderLogRepository) UpdateByID(ID string, request *models.DeliveryOrderLog, ctx context.Context, resultChan chan *models.DeliveryOrderLogChan) {
	response := &models.DeliveryOrderLogChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	objectID, _ := primitive.ObjectIDFromHex(ID)
	filter := bson.M{"_id": objectID}
	total, err := collection.CountDocuments(ctx, filter)

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	if total == 0 {
		err = helper.NewError("data not found")
		errorLogData := helper.WriteLog(err, 404, "data not found")
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	updateData := bson.M{"$set": request}
	_, err = collection.UpdateOne(r.mongod.GetCtx(), filter, updateData)

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	response.Error = nil
	resultChan <- response
	return
}
