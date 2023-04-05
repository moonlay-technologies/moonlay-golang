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

type DeliveryOrderLogRepositoryInterface interface {
	Insert(request *models.DeliveryOrderLog, ctx context.Context, result chan *models.DeliveryOrderLogChan)
	Get(request *models.DeliveryOrderEventLogRequest, countOnly bool, ctx context.Context, resultChan chan *models.GetDeliveryOrderLogsChan)
	GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.GetDeliveryOrderLogChan)
	GetByCode(doCode string, status string, action string, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderLogChan)
	UpdateByID(ID string, request *models.DeliveryOrderLog, ctx context.Context, result chan *models.DeliveryOrderLogChan)
	InsertJourney(request *models.DeliveryOrderJourney, ctx context.Context, result chan *models.DeliveryOrderJourneyChan)
}

type deliveryOrderLogRepository struct {
	logger            log.Logger
	mongod            mongodb.MongoDBInterface
	collectionLog     string
	collectionJourney string
}

func InitDeliveryOrderLogRepository(mongod mongodb.MongoDBInterface) DeliveryOrderLogRepositoryInterface {
	return &deliveryOrderLogRepository{
		mongod:            mongod,
		collectionLog:     constants.DELIVERY_ORDER_TABLE_LOGS,
		collectionJourney: constants.DELIVERY_ORDER_TABLE_JOURNEYS,
	}
}

func (r *deliveryOrderLogRepository) Insert(request *models.DeliveryOrderLog, ctx context.Context, resultChan chan *models.DeliveryOrderLogChan) {
	response := &models.DeliveryOrderLogChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collectionLog)
	result, err := collection.InsertOne(ctx, request)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
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

func (r *deliveryOrderLogRepository) Get(request *models.DeliveryOrderEventLogRequest, countOnly bool, ctx context.Context, resultChan chan *models.GetDeliveryOrderLogsChan) {
	response := &models.GetDeliveryOrderLogsChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collectionLog)
	filter := bson.M{}
	sort := bson.M{}
	asc := 1
	desc := -1

	if request.SortField == "do_code" {
		if request.SortValue == "asc" {
			sort = bson.M{
				"do_code": asc,
			}
		} else if request.SortValue == "desc" {
			sort = bson.M{
				"do_code": desc,
			}
		}
	} else if request.SortField == "status" {
		if request.SortValue == "asc" {
			sort = bson.M{
				"status": asc,
			}
		} else if request.SortValue == "desc" {
			sort = bson.M{
				"status": desc,
			}
		}
	} else if request.SortField == "agent_name" {
		if request.SortValue == "asc" {
			sort = bson.M{
				"data.agent_name": asc,
			}
		} else if request.SortValue == "desc" {
			sort = bson.M{
				"data.agent_name": desc,
			}
		}
	} else if request.SortField == "created_at" {
		if request.SortValue == "asc" {
			sort = bson.M{
				"created_at": asc,
			}
		} else if request.SortValue == "desc" {
			sort = bson.M{
				"created_at": desc,
			}
		}
	}

	if request.GlobalSearchValue != "" {
		var status string
		switch request.GlobalSearchValue {
		case constants.EVENT_LOG_STATUS_0:
			status = "0"
			filter = bson.M{
				"$or": []bson.M{
					{"do_code": bson.M{"$regex": request.GlobalSearchValue, "$options": "i"}},
					{"status": bson.M{"$regex": status, "$options": "i"}},
					{"data.agent_name": bson.M{"$regex": request.GlobalSearchValue, "$options": "i"}},
				},
			}
		case constants.EVENT_LOG_STATUS_1:
			status = "1"
			filter = bson.M{
				"$or": []bson.M{
					{"do_code": bson.M{"$regex": request.GlobalSearchValue, "$options": "i"}},
					{"status": bson.M{"$regex": status, "$options": "i"}},
					{"data.agent_name": bson.M{"$regex": request.GlobalSearchValue, "$options": "i"}},
				},
			}
		case constants.EVENT_LOG_STATUS_2:
			status = "2"
			filter = bson.M{
				"$or": []bson.M{
					{"do_code": bson.M{"$regex": request.GlobalSearchValue, "$options": "i"}},
					{"status": bson.M{"$regex": status, "$options": "i"}},
					{"data.agent_name": bson.M{"$regex": request.GlobalSearchValue, "$options": "i"}},
				},
			}
		default:
			filter = bson.M{
				"$or": []bson.M{
					{"do_code": bson.M{"$regex": request.GlobalSearchValue, "$options": "i"}},
					{"status": bson.M{"$regex": request.GlobalSearchValue, "$options": "i"}},
					{"data.agent_name": bson.M{"$regex": request.GlobalSearchValue, "$options": "i"}},
				},
			}
		}
	}

	if request.ID != "" {
		id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusBadRequest, "Ada kesalahan pada request data, silahkan dicek kembali")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}
		filter["_id"] = id
	}

	if request.RequestID != "" {
		filter["request_id"] = request.RequestID
	}

	if request.AgentID > 0 {
		filter["data.agent_id"] = request.AgentID
	}

	if request.Status != "" {
		var status string
		switch request.Status {
		case constants.EVENT_LOG_STATUS_0:
			status = "0"
		case constants.EVENT_LOG_STATUS_1:
			status = "1"
		case constants.EVENT_LOG_STATUS_2:
			status = "2"
		default:
			status = ""
		}
		filter["status"] = status
	}

	option := options.Find().SetSkip(int64((request.Page - 1) * request.PerPage)).SetLimit(int64(request.PerPage)).SetSort(sort)
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
		deliveryOrderLogs := []*models.GetDeliveryOrderLog{}
		cursor, err := collection.Find(ctx, filter, option)

		if err != nil {
			errorLogData := helper.WriteLog(err, http.StatusInternalServerError, "Terjadi Kesalahan, Silahkan Coba lagi Nanti")
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var deliveryOrderLog *models.GetDeliveryOrderLog
			if err := cursor.Decode(&deliveryOrderLog); err != nil {
				errorLogData := helper.WriteLog(err, http.StatusInternalServerError, "Terjadi Kesalahan, Silahkan Coba lagi Nanti")
				response.Error = err
				response.ErrorLog = errorLogData
				resultChan <- response
				return
			}

			deliveryOrderLogs = append(deliveryOrderLogs, deliveryOrderLog)
		}

		response.DeliveryOrderLog = deliveryOrderLogs
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

func (r *deliveryOrderLogRepository) GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.GetDeliveryOrderLogChan) {
	response := &models.GetDeliveryOrderLogChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collectionLog)
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

	if countOnly == false {
		deliveryOrderLog := &models.GetDeliveryOrderLog{}
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

func (r *deliveryOrderLogRepository) GetByCode(doCode string, status string, action string, countOnly bool, ctx context.Context, resultChan chan *models.DeliveryOrderLogChan) {
	response := &models.DeliveryOrderLogChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collectionLog)
	filter := bson.M{constants.COLUMN_DELIVERY_ORDER_CODE: doCode, constants.COLUMN_STATUS: status, constants.COLUMN_ACTION: action}
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

func (r *deliveryOrderLogRepository) UpdateByID(ID string, request *models.DeliveryOrderLog, ctx context.Context, resultChan chan *models.DeliveryOrderLogChan) {
	response := &models.DeliveryOrderLogChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collectionLog)
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

func (r *deliveryOrderLogRepository) InsertJourney(request *models.DeliveryOrderJourney, ctx context.Context, resultChan chan *models.DeliveryOrderJourneyChan) {
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
