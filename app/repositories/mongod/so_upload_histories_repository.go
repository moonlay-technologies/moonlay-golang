package repositories

import (
	"context"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/global/utils/helper"
	"order-service/global/utils/model"
	"order-service/global/utils/mongodb"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SoUploadHistoriesRepositoryInterface interface {
	Insert(request *models.SoUploadHistory, ctx context.Context, resultChan chan *models.SoUploadHistoryChan)
	Get(request *models.GetSoUploadHistoriesRequest, countOnly bool, ctx context.Context, resultChan chan *models.SoUploadHistoriesChan)
	GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.SoUploadHistoryChan)
	GetByHistoryID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.GetSoUploadHistoryResponseChan)
	UpdateByID(ID string, request *models.SoUploadHistory, ctx context.Context, resultChan chan *models.SoUploadHistoryChan)
}

type soUploadHistoriesRepository struct {
	logger     log.Logger
	mongod     mongodb.MongoDBInterface
	collection string
}

func InitSoUploadHistoriesRepositoryInterface(mongod mongodb.MongoDBInterface) SoUploadHistoriesRepositoryInterface {
	return &soUploadHistoriesRepository{
		mongod:     mongod,
		collection: constants.SO_UPLOAD_TABLE_HISTORIES,
	}
}

func (r *soUploadHistoriesRepository) Insert(request *models.SoUploadHistory, ctx context.Context, resultChan chan *models.SoUploadHistoryChan) {
	response := &models.SoUploadHistoryChan{}
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
	response.SoUploadHistory = request
	response.Error = nil
	resultChan <- response
	return
}

func (r *soUploadHistoriesRepository) Get(request *models.GetSoUploadHistoriesRequest, countOnly bool, ctx context.Context, resultChan chan *models.SoUploadHistoriesChan) {
	response := &models.SoUploadHistoriesChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	filter := bson.M{}
	sort := bson.M{}
	asc := 1
	desc := -1

	if request.SortField == "agent_name" {
		if request.SortValue == "asc" {
			sort = bson.M{
				"agent_name": asc,
			}
		} else if request.SortValue == "desc" {
			sort = bson.M{
				"agent_name": desc,
			}
		}
	} else if request.SortField == "file_name" {
		if request.SortValue == "asc" {
			sort = bson.M{
				"file_name": asc,
			}
		} else if request.SortValue == "desc" {
			sort = bson.M{
				"file_name": desc,
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
		uploadedBy, _ := strconv.ParseInt(request.GlobalSearchValue, 10, 64)
		createdAt, _ := time.Parse("2006-01-02", request.GlobalSearchValue)

		filter = bson.M{
			"$or": []bson.M{
				{"bulk_code": bson.M{"$regex": ".*" + request.GlobalSearchValue + ".*", "$options": "i"}},
				{"agent_name": bson.M{"$regex": ".*" + request.GlobalSearchValue + ".*", "$options": "i"}},
				{"file_name": bson.M{"$regex": ".*" + request.GlobalSearchValue + ".*", "$options": "i"}},
				{"status": bson.M{"$regex": ".*" + request.GlobalSearchValue + ".*", "$options": "i"}},
				{"uploaded_by": uploadedBy},
				{"created_at": bson.M{"$gte": createdAt, "$lte": createdAt.AddDate(0, 0, 1)}},
			},
		}
	}

	if request.RequestID != "" {
		filter["request_id"] = request.RequestID
	}

	if request.FileName != "" {
		filter["file_name"] = request.FileName
	}

	if request.BulkCode != "" {
		filter["bulk_code"] = request.BulkCode
	}

	if request.AgentID > 0 {
		filter["agent_id"] = request.AgentID
	}

	if request.Status != "" {
		filter["status"] = request.Status
	}

	if request.UploadedBy > 0 {
		filter["uploaded_by"] = request.UploadedBy
	}

	if request.StartUploadAt != "" && request.EndUploadAt == "" {
		startUploadAt, err := time.Parse("2006-01-02", request.StartUploadAt)
		if err != nil {
			errorLogData := helper.NewWriteLog(model.ErrorLog{
				Message:       "Ada kesalahan pada request data, silahkan dicek kembali",
				SystemMessage: helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_GET, fmt.Sprintf("field %s harus memiliki format yyyy-mm-dd", "start_upload_at")),
				StatusCode:    http.StatusBadRequest,
				Err:           fmt.Errorf("invalid Process"),
			})
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}
		filter["created_at"] = bson.M{"$gte": startUploadAt, "$lte": startUploadAt.AddDate(0, 0, 1)}
	}

	if request.EndUploadAt != "" && request.StartUploadAt == "" {
		endUploadAt, err := time.Parse("2006-01-02", request.EndUploadAt)
		if err != nil {
			errorLogData := helper.NewWriteLog(model.ErrorLog{
				Message:       "Ada kesalahan pada request data, silahkan dicek kembali",
				SystemMessage: helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_GET, fmt.Sprintf("field %s harus memiliki format yyyy-mm-dd", "end_upload_at")),
				StatusCode:    http.StatusBadRequest,
				Err:           fmt.Errorf("invalid Process"),
			})
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}
		filter["created_at"] = bson.M{"$gte": endUploadAt, "$lte": endUploadAt.AddDate(0, 0, 1)}
	}

	if request.FinishProcessDateStart != "" && request.FinishProcessDateEnd == "" {
		finishProcessDateStart, err := time.Parse("2006-01-02", request.FinishProcessDateStart)
		if err != nil {
			errorLogData := helper.NewWriteLog(model.ErrorLog{
				Message:       "Ada kesalahan pada request data, silahkan dicek kembali",
				SystemMessage: helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_GET, fmt.Sprintf("field %s harus memiliki format yyyy-mm-dd", "finish_process_date_start")),
				StatusCode:    http.StatusBadRequest,
				Err:           fmt.Errorf("invalid Process"),
			})
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}
		filter["updated_at"] = bson.M{"$gte": finishProcessDateStart, "$lte": finishProcessDateStart.AddDate(0, 0, 1)}
	}

	if request.FinishProcessDateEnd != "" && request.FinishProcessDateStart == "" {
		finishProcessDateEnd, err := time.Parse("2006-01-02", request.FinishProcessDateEnd)
		if err != nil {
			errorLogData := helper.NewWriteLog(model.ErrorLog{
				Message:       "Ada kesalahan pada request data, silahkan dicek kembali",
				SystemMessage: helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_GET, fmt.Sprintf("field %s harus memiliki format yyyy-mm-dd", "finish_process_date_end")),
				StatusCode:    http.StatusBadRequest,
				Err:           fmt.Errorf("invalid Process"),
			})
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}
		filter["updated_at"] = bson.M{"$gte": finishProcessDateEnd, "$lte": finishProcessDateEnd.AddDate(0, 0, 1)}
	}

	if request.StartUploadAt != "" && request.EndUploadAt != "" {
		startUploadAt, err := time.Parse("2006-01-02", request.StartUploadAt)
		if err != nil {
			errorLogData := helper.NewWriteLog(model.ErrorLog{
				Message:       "Ada kesalahan pada request data, silahkan dicek kembali",
				SystemMessage: helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_GET, fmt.Sprintf("field %s harus memiliki format yyyy-mm-dd", "start_upload_at")),
				StatusCode:    http.StatusBadRequest,
				Err:           fmt.Errorf("invalid Process"),
			})
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		endUploadAt, err := time.Parse("2006-01-02", request.EndUploadAt)
		if err != nil {
			errorLogData := helper.NewWriteLog(model.ErrorLog{
				Message:       "Ada kesalahan pada request data, silahkan dicek kembali",
				SystemMessage: helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_GET, fmt.Sprintf("field %s harus memiliki format yyyy-mm-dd", "end_upload_at")),
				StatusCode:    http.StatusBadRequest,
				Err:           fmt.Errorf("invalid Process"),
			})
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		filter["created_at"] = bson.M{"$gte": startUploadAt, "$lte": endUploadAt.AddDate(0, 0, 1)}
	}

	if request.FinishProcessDateStart != "" && request.FinishProcessDateEnd != "" {
		finishProcessDateStart, err := time.Parse("2006-01-02", request.FinishProcessDateStart)
		if err != nil {
			errorLogData := helper.NewWriteLog(model.ErrorLog{
				Message:       "Ada kesalahan pada request data, silahkan dicek kembali",
				SystemMessage: helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_GET, fmt.Sprintf("field %s harus memiliki format yyyy-mm-dd", "finish_process_date_start")),
				StatusCode:    http.StatusBadRequest,
				Err:           fmt.Errorf("invalid Process"),
			})
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		finishProcessDateEnd, err := time.Parse("2006-01-02", request.FinishProcessDateEnd)
		if err != nil {
			errorLogData := helper.NewWriteLog(model.ErrorLog{
				Message:       "Ada kesalahan pada request data, silahkan dicek kembali",
				SystemMessage: helper.GenerateUnprocessableErrorMessage(constants.ERROR_ACTION_NAME_GET, fmt.Sprintf("field %s harus memiliki format yyyy-mm-dd", "finish_process_date_end")),
				StatusCode:    http.StatusBadRequest,
				Err:           fmt.Errorf("invalid Process"),
			})
			response.Error = err
			response.ErrorLog = errorLogData
			resultChan <- response
			return
		}

		filter["updated_at"] = bson.M{"$gte": finishProcessDateStart, "$lte": finishProcessDateEnd.AddDate(0, 0, 1)}
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

	if !countOnly {
		SoUploadHistories := []*models.SoUploadHistory{}
		cursor, err := collection.Find(ctx, filter, option)

		if err != nil {
			fmt.Println(err.Error())
			response.Error = err
			resultChan <- response
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var soUploadHistory *models.SoUploadHistory
			if err := cursor.Decode(&soUploadHistory); err != nil {
				response.Error = err
				resultChan <- response
				return
			}

			SoUploadHistories = append(SoUploadHistories, soUploadHistory)
		}

		response.SoUploadHistories = SoUploadHistories
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.SoUploadHistories = nil
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *soUploadHistoriesRepository) GetByID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.SoUploadHistoryChan) {
	response := &models.SoUploadHistoryChan{}
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
		err = helper.NewError(constants.ERROR_DATA_NOT_FOUND)
		errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	if !countOnly {
		soUploadHistory := &models.SoUploadHistory{}
		err = collection.FindOne(ctx, filter).Decode(soUploadHistory)

		if err != nil {
			fmt.Println(err.Error())
			response.Error = err
			resultChan <- response
			return
		}

		response.SoUploadHistory = soUploadHistory
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.SoUploadHistory = nil
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *soUploadHistoriesRepository) GetByHistoryID(ID string, countOnly bool, ctx context.Context, resultChan chan *models.GetSoUploadHistoryResponseChan) {
	response := &models.GetSoUploadHistoryResponseChan{}
	collection := r.mongod.Client().Database(os.Getenv("MONGO_DATABASE")).Collection(r.collection)
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusBadRequest, "Ada kesalahan pada request data, silahkan dicek kembali")
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

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
		err = helper.NewError(constants.ERROR_DATA_NOT_FOUND)
		errorLogData := helper.WriteLog(err, http.StatusNotFound, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	if !countOnly {
		cursor, err := collection.Aggregate(ctx, bson.A{
			bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: objectID}}}},
			bson.D{
				{Key: "$lookup",
					Value: bson.D{
						{Key: "from", Value: "so_upload_error_logs"},
						{Key: "localField", Value: "_id"},
						{Key: "foreignField", Value: "so_upload_history_id"},
						{Key: "as", Value: "so_upload_error_logs"},
					},
				},
			},
		})

		if err != nil {
			fmt.Println(err.Error())
			response.Error = err
			resultChan <- response
			return
		}

		defer cursor.Close(ctx)
		soUploadHistoryResponse := &models.GetSoUploadHistoryResponse{}
		soUploadHistory := &models.SoUploadHistory{}
		for cursor.Next(ctx) {
			if err := cursor.Decode(&soUploadHistoryResponse); err != nil {
				response.Error = err
				resultChan <- response
				return
			}

			if err := cursor.Decode(&soUploadHistory); err != nil {
				response.Error = err
				resultChan <- response
				return
			}
		}

		soUploadHistoryResponse.GetSoUploadHistoryResponseMap(soUploadHistory)

		response.SoUploadHistories = soUploadHistoryResponse
		response.Total = total
		response.Error = nil
		resultChan <- response
		return
	} else {
		response.SoUploadHistories = nil
		response.Total = total
		resultChan <- response
		return
	}
}

func (r *soUploadHistoriesRepository) UpdateByID(ID string, request *models.SoUploadHistory, ctx context.Context, resultChan chan *models.SoUploadHistoryChan) {
	response := &models.SoUploadHistoryChan{}
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
	response.SoUploadHistory = request
	resultChan <- response
}
