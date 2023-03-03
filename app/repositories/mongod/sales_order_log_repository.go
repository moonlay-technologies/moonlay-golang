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
	Get(request *models.SalesOrderEventLogRequest, countOnly bool, ctx context.Context, resultChan chan *models.GetSalesOrderLogsChan)
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

func (r *salesOrderLogRepository) Get(request *models.SalesOrderEventLogRequest, countOnly bool, ctx context.Context, resultChan chan *models.GetSalesOrderLogsChan) {
	response := &models.GetSalesOrderLogsChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	filter := bson.M{}
	sort := bson.M{}
	asc := 1
	desc := -1

	if request.SortField == "created_at" && request.SortValue == "asc" {
		sort = bson.M{
			"created_at": asc,
		}
	} else if request.SortField == "updated_at" && request.SortValue == "asc" {
		sort = bson.M{
			"updated_at": asc,
		}
	} else if request.SortField == "created_at" && request.SortValue == "desc" {
		sort = bson.M{
			"created_at": desc,
		}
	} else if request.SortField == "updated_at" && request.SortValue == "desc" {
		sort = bson.M{
			"updated_at": desc,
		}
	} else {
		sort = bson.M{
			"created_at": desc,
		}
	}

	// if request.GlobalSearchValue != "" {
	// 	filter = bson.M{
	// 		"$or": []bson.M{
	// 			bson.M{"$in":},
	// 			bson.M{},
	// 		},
	// }

	if request.RequestID != "" {
		filter = bson.M{"request_id": request.RequestID}
	}

	if request.SoCode != "" {
		filter["so_code"] = request.SoCode
	}

	if request.Status != "" {
		filter["status"] = request.Status
	}

	if request.Action != "" {
		filter["action"] = request.Action
	}

	if request.AgentID > 0 {
		filter["data.agent_id"] = request.AgentID
	}

	option := options.Find().SetSkip(int64((request.Page - 1) * request.PerPage)).SetLimit(int64(request.PerPage)).SetSort(sort)
	total, err := collection.CountDocuments(ctx, filter)
	fmt.Println("errs", err)
	fmt.Println("total", total)
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
		salesOrderLogs := []*models.GetSalesOrderLog{}
		cursor, errs := collection.Find(ctx, filter, option)

		if errs != nil {
			response.Error = err
			resultChan <- response
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var salesOrderLog *models.GetSalesOrderLog
			if err := cursor.Decode(&salesOrderLog); err != nil {
				response.Error = err
				resultChan <- response
				return
			}

			salesOrderLogs = append(salesOrderLogs, salesOrderLog)
		}

		response.SalesOrderLogs = salesOrderLogs
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.SalesOrderLogs = nil
		response.Total = total
		resultChan <- response
		return
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
