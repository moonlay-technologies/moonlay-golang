package repositories

import (
	"context"
	"fmt"
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

type DeliveryOrderDetailLogRepositoryInterface interface {
	Insert(request *models.DeliveryOrderDetailLog, ctx context.Context, result chan *models.DeliveryOrderDetailLogChan)
	GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderDetailLogChan)
	GetByCode(doDetailCode string, status string, action string, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderDetailLogChan)
	UpdateByID(ID string, request *models.DeliveryOrderDetailLog, ctx context.Context, result chan *models.DeliveryOrderDetailLogChan)
}

type deliveryOrderDetailLogRepository struct {
	logger     log.Logger
	mongod     mongodb.MongoDBInterface
	collection string
}

func InitDeliveryOrderDetailLogRepository(mongod mongodb.MongoDBInterface) DeliveryOrderDetailLogRepositoryInterface {
	return &deliveryOrderDetailLogRepository{
		mongod:     mongod,
		collection: constants.DELIVERY_ORDER_TABLE_LOGS,
	}
}

func (r *deliveryOrderDetailLogRepository) GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderDetailLogChan) {
	response := &models.DeliveryOrderDetailLogChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	objectID, _ := primitive.ObjectIDFromHex(ID)
	filter := bson.M{"_id": objectID}
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
		deliveryOrderDetailLog := &models.DeliveryOrderDetailLog{}
		err = collection.FindOne(ctx, filter).Decode(deliveryOrderDetailLog)

		if err != nil {
			fmt.Println(err.Error())
			response.Error = err
			resultChan <- response
			return
		}

		response.DeliveryOrderDetailLog = deliveryOrderDetailLog
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.DeliveryOrderDetailLog = nil
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *deliveryOrderDetailLogRepository) GetByCode(doDetailCode string, status string, action string, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderDetailLogChan) {
	response := &models.DeliveryOrderDetailLogChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	filter := bson.M{constants.COLUMN_DELIVERY_ORDER_CODE: doDetailCode, constants.COLUMN_STATUS: status, constants.COLUMN_ACTION: action}
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
		deliveryOrderDetailLog := &models.DeliveryOrderDetailLog{}
		err = collection.FindOne(ctx, filter).Decode(deliveryOrderDetailLog)

		if err != nil {
			fmt.Println(err.Error())
			response.Error = err
			resultChan <- response
			return
		}

		response.DeliveryOrderDetailLog = deliveryOrderDetailLog
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.DeliveryOrderDetailLog = nil
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *deliveryOrderDetailLogRepository) Insert(request *models.DeliveryOrderDetailLog, ctx context.Context, resultChan chan *models.DeliveryOrderDetailLogChan) {
	response := &models.DeliveryOrderDetailLogChan{}
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
	response.DeliveryOrderDetailLog = request
	response.Error = nil
	resultChan <- response
	return
}

func (r *deliveryOrderDetailLogRepository) UpdateByID(ID string, request *models.DeliveryOrderDetailLog, ctx context.Context, resultChan chan *models.DeliveryOrderDetailLogChan) {
	response := &models.DeliveryOrderDetailLogChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	objectID, _ := primitive.ObjectIDFromHex(ID)
	filter := bson.M{"_id": objectID}
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
		errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	updateData := bson.M{"$set": request}
	_, err = collection.UpdateOne(r.mongod.GetCtx(), filter, updateData)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	response.Error = nil
	resultChan <- response
	return
}
