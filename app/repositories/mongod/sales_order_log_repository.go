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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SalesOrderLogRepositoryInterface interface {
	Insert(request *models.SalesOrderLog, ctx context.Context, result chan *models.SalesOrderLogChan)
	GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.SalesOrderLogChan)
	GetByCollumn(collumnName string, value string, countOnly bool, ctx context.Context, resultChan chan *models.SalesOrderLogChan)
	UpdateByID(ID string, request *models.SalesOrderLog, ctx context.Context, result chan *models.SalesOrderLogChan)
}

type salesOrderLogRepository struct {
	logger     log.Logger
	mongod     mongodb.MongoDBInterface
	collection string
}

func InitSalesOrderLogRepository(mongod mongodb.MongoDBInterface) SalesOrderLogRepositoryInterface {
	return &salesOrderLogRepository{
		mongod:     mongod,
		collection: constants.SALES_ORDER_TABLE_LOGS,
	}
}

func (r *salesOrderLogRepository) GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.SalesOrderLogChan) {
	response := &models.SalesOrderLogChan{}
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
		salesOrderLog := &models.SalesOrderLog{}
		err = collection.FindOne(ctx, filter).Decode(salesOrderLog)

		if err != nil {
			fmt.Println(err.Error())
			response.Error = err
			resultChan <- response
			return
		}

		response.SalesOrderLog = salesOrderLog
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.SalesOrderLog = nil
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *salesOrderLogRepository) GetByCollumn(collumnName string, value string, countOnly bool, ctx context.Context, resultChan chan *models.SalesOrderLogChan) {
	response := &models.SalesOrderLogChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	fmt.Println("collumn = ", collumnName)
	fmt.Println("value = ", value)
	filter := bson.M{collumnName: value}
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
		salesOrderLog := &models.SalesOrderLog{}
		opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
		err = collection.FindOne(ctx, filter, opts).Decode(salesOrderLog)

		if err != nil {
			fmt.Println(err.Error())
			response.Error = err
			resultChan <- response
			return
		}

		response.SalesOrderLog = salesOrderLog
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.SalesOrderLog = nil
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *salesOrderLogRepository) Insert(request *models.SalesOrderLog, ctx context.Context, resultChan chan *models.SalesOrderLogChan) {
	response := &models.SalesOrderLogChan{}
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
	response.SalesOrderLog = request
	response.Error = nil
	resultChan <- response
	return
}

func (r *salesOrderLogRepository) UpdateByID(ID string, request *models.SalesOrderLog, ctx context.Context, resultChan chan *models.SalesOrderLogChan) {
	response := &models.SalesOrderLogChan{}
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
