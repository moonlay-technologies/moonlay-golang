package usecases

import (
	"context"
	"fmt"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/app/repositories"
	mongoRepositories "order-service/app/repositories/mongod"
	openSearchRepositories "order-service/app/repositories/open_search"
	"order-service/global/utils/helper"
	kafkadbo "order-service/global/utils/kafka"
	"order-service/global/utils/model"
	"strconv"
	"strings"
	"time"

	"github.com/bxcodec/dbresolver"
)

type UploadUseCaseInterface interface {
	UploadSOSJ(request *models.UploadSOSJRequest, ctx context.Context) *model.ErrorLog
	UploadDO(request *models.UploadDORequest, ctx context.Context) *model.ErrorLog
	UploadSO(request *models.UploadSORequest, ctx context.Context) *model.ErrorLog
	RetryUploadSO(soUploadHistoryId string, ctx context.Context) *model.ErrorLog
	RetryUploadDO(sjUploadHistoryId string, ctx context.Context) *model.ErrorLog
	RetryUploadSOSJ(sosjUploadHistoryId string, ctx context.Context) *model.ErrorLog
	GetSosjUploadErrorLogs(request *models.GetSosjUploadErrorLogsRequest, ctx context.Context) (*models.GetSosjUploadErrorLogsResponse, *model.ErrorLog)
	GetSosjUploadHistoryById(id string, ctx context.Context) (*models.GetSosjUploadHistoryResponse, *model.ErrorLog)
	GetSosjUploadErrorLogBySosjUploadHistoryId(request *models.GetSosjUploadErrorLogsRequest, ctx context.Context) (*models.GetSosjUploadErrorLogsResponse, *model.ErrorLog)
}

type uploadUseCase struct {
	uploadRepository                     repositories.UploadRepositoryInterface
	salesOrderRepository                 repositories.SalesOrderRepositoryInterface
	salesOrderDetailRepository           repositories.SalesOrderDetailRepositoryInterface
	orderStatusRepository                repositories.OrderStatusRepositoryInterface
	orderSourceRepository                repositories.OrderSourceRepositoryInterface
	agentRepository                      repositories.AgentRepositoryInterface
	brandRepository                      repositories.BrandRepositoryInterface
	storeRepository                      repositories.StoreRepositoryInterface
	productRepository                    repositories.ProductRepositoryInterface
	uomRepository                        repositories.UomRepositoryInterface
	deliveryOrderRepository              repositories.DeliveryOrderRepositoryInterface
	deliveryOrderDetailRepository        repositories.DeliveryOrderDetailRepositoryInterface
	salesOrderLogRepository              mongoRepositories.SalesOrderLogRepositoryInterface
	salesOrderJourneysRepository         mongoRepositories.SalesOrderJourneysRepositoryInterface
	salesOrderDetailJourneysRepository   mongoRepositories.SalesOrderDetailJourneysRepositoryInterface
	userRepository                       repositories.UserRepositoryInterface
	salesmanRepository                   repositories.SalesmanRepositoryInterface
	categoryRepository                   repositories.CategoryRepositoryInterface
	salesOrderOpenSearchRepository       openSearchRepositories.SalesOrderOpenSearchRepositoryInterface
	salesOrderDetailOpenSearchRepository openSearchRepositories.SalesOrderDetailOpenSearchRepositoryInterface
	warehouseRepository                  repositories.WarehouseRepositoryInterface
	soUploadHistoriesRepository          mongoRepositories.SoUploadHistoriesRepositoryInterface
	soUploadErrorLogsRepository          mongoRepositories.SoUploadErrorLogsRepositoryInterface
	sosjUploadHistoriesRepository        mongoRepositories.SOSJUploadHistoriesRepositoryInterface
	sosjUploadErrorLogsRepository        mongoRepositories.SosjUploadErrorLogsRepositoryInterface
	doUploadHistoriesRepository          mongoRepositories.DoUploadHistoriesRepositoryInterface
	kafkaClient                          kafkadbo.KafkaClientInterface
	db                                   dbresolver.DB
	ctx                                  context.Context
}

func InitUploadUseCaseInterface(salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, orderSourceRepository repositories.OrderSourceRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, uomRepository repositories.UomRepositoryInterface, deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, deliveryOrderDetailRepository repositories.DeliveryOrderDetailRepositoryInterface, salesOrderLogRepository mongoRepositories.SalesOrderLogRepositoryInterface, salesOrderJourneysRepository mongoRepositories.SalesOrderJourneysRepositoryInterface, salesOrderDetailJourneysRepository mongoRepositories.SalesOrderDetailJourneysRepositoryInterface, userRepository repositories.UserRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, categoryRepository repositories.CategoryRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, salesOrderDetailOpenSearchRepository openSearchRepositories.SalesOrderDetailOpenSearchRepositoryInterface, uploadRepository repositories.UploadRepositoryInterface, warehouseRepository repositories.WarehouseRepositoryInterface, soUploadHistoriesRepository mongoRepositories.SoUploadHistoriesRepositoryInterface, soUploadErrorLogsRepository mongoRepositories.SoUploadErrorLogsRepositoryInterface, sosjUploadHistoriesRepository mongoRepositories.SOSJUploadHistoriesRepositoryInterface, sosjUploadErrorLogsRepository mongoRepositories.SosjUploadErrorLogsRepositoryInterface, doUploadHistoriesRepository mongoRepositories.DoUploadHistoriesRepositoryInterface, kafkaClient kafkadbo.KafkaClientInterface, db dbresolver.DB, ctx context.Context) UploadUseCaseInterface {
	return &uploadUseCase{
		salesOrderRepository:                 salesOrderRepository,
		salesOrderDetailRepository:           salesOrderDetailRepository,
		orderStatusRepository:                orderStatusRepository,
		orderSourceRepository:                orderSourceRepository,
		agentRepository:                      agentRepository,
		brandRepository:                      brandRepository,
		storeRepository:                      storeRepository,
		productRepository:                    productRepository,
		uomRepository:                        uomRepository,
		deliveryOrderRepository:              deliveryOrderRepository,
		deliveryOrderDetailRepository:        deliveryOrderDetailRepository,
		salesOrderLogRepository:              salesOrderLogRepository,
		salesOrderJourneysRepository:         salesOrderJourneysRepository,
		salesOrderDetailJourneysRepository:   salesOrderDetailJourneysRepository,
		userRepository:                       userRepository,
		salesmanRepository:                   salesmanRepository,
		categoryRepository:                   categoryRepository,
		salesOrderOpenSearchRepository:       salesOrderOpenSearchRepository,
		salesOrderDetailOpenSearchRepository: salesOrderDetailOpenSearchRepository,
		uploadRepository:                     uploadRepository,
		warehouseRepository:                  warehouseRepository,
		soUploadHistoriesRepository:          soUploadHistoriesRepository,
		soUploadErrorLogsRepository:          soUploadErrorLogsRepository,
		sosjUploadHistoriesRepository:        sosjUploadHistoriesRepository,
		sosjUploadErrorLogsRepository:        sosjUploadErrorLogsRepository,
		doUploadHistoriesRepository:          doUploadHistoriesRepository,
		kafkaClient:                          kafkaClient,
		db:                                   db,
		ctx:                                  ctx,
	}
}

func (u *uploadUseCase) UploadSOSJ(request *models.UploadSOSJRequest, ctx context.Context) *model.ErrorLog {

	now := time.Now()
	user := ctx.Value("user").(*models.UserClaims)

	// Check Agent By Id
	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(4, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	agentId := int64(user.AgentID)
	userId := int64(user.UserID)
	url := strings.Split(request.File, "/")
	message := &models.UploadHistory{
		RequestId:       ctx.Value("RequestId").(string),
		BulkCode:        "SOSJ-" + strconv.Itoa(user.AgentID) + "-" + fmt.Sprint(now.Unix()),
		FileName:        url[len(url)-1],
		FilePath:        request.File,
		AgentId:         &agentId,
		AgentName:       getAgentResult.Agent.Name,
		UploadedBy:      &userId,
		UploadedByName:  user.FirstName + " " + user.LastName,
		UploadedByEmail: user.UserEmail,
		Status:          constants.UPLOAD_STATUS_HISTORY_IN_PROGRESS,
		UpdatedBy:       &userId,
		UpdatedByName:   user.FirstName + " " + user.LastName,
		UpdatedByEmail:  user.UserEmail,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	sosjUploadHistoryJourneysResultChan := make(chan *models.UploadHistoryChan)
	go u.sosjUploadHistoriesRepository.Insert(message, ctx, sosjUploadHistoryJourneysResultChan)
	sosjUploadHistoryJourneysResult := <-sosjUploadHistoryJourneysResultChan

	if sosjUploadHistoryJourneysResult.Error != nil {
		errorLogData := helper.WriteLog(sosjUploadHistoryJourneysResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	keyKafka := []byte("upload")
	messageKafka := []byte(sosjUploadHistoryJourneysResult.UploadHistory.ID.Hex())

	err := u.kafkaClient.WriteToTopic(constants.UPLOAD_SOSJ_FILE_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}

func (u *uploadUseCase) UploadDO(request *models.UploadDORequest, ctx context.Context) *model.ErrorLog {

	now := time.Now()
	user := ctx.Value("user").(*models.UserClaims)

	// Check Agent By Id
	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(4, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	agentId := int64(user.AgentID)
	userId := int64(user.UserID)
	url := strings.Split(request.File, "/")
	message := &models.DoUploadHistory{
		RequestId:       ctx.Value("RequestId").(string),
		BulkCode:        "SJ-" + strconv.Itoa(user.AgentID) + "-" + fmt.Sprint(now.Unix()),
		FileName:        url[len(url)-1],
		FilePath:        request.File,
		AgentId:         &agentId,
		AgentName:       getAgentResult.Agent.Name,
		UploadedBy:      &userId,
		UploadedByName:  user.FirstName + " " + user.LastName,
		UploadedByEmail: user.UserEmail,
		Status:          constants.UPLOAD_STATUS_HISTORY_IN_PROGRESS,
		UpdatedBy:       &userId,
		UpdatedByName:   user.FirstName + " " + user.LastName,
		UpdatedByEmail:  user.UserEmail,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	sjUploadHistoryJourneysResultChan := make(chan *models.DoUploadHistoryChan)
	go u.doUploadHistoriesRepository.Insert(message, ctx, sjUploadHistoryJourneysResultChan)
	sjUploadHistoryJourneysResult := <-sjUploadHistoryJourneysResultChan

	if sjUploadHistoryJourneysResult.Error != nil {
		errorLogData := helper.WriteLog(sjUploadHistoryJourneysResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	keyKafka := []byte("upload")
	messageKafka := []byte(sjUploadHistoryJourneysResult.DoUploadHistory.ID.Hex())

	err := u.kafkaClient.WriteToTopic(constants.UPLOAD_DO_FILE_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}

func (u *uploadUseCase) UploadSO(request *models.UploadSORequest, ctx context.Context) *model.ErrorLog {

	now := time.Now()
	user := ctx.Value("user").(*models.UserClaims)

	// Check Agent By Id
	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(4, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	agentId := int64(user.AgentID)
	userId := int64(user.UserID)
	url := strings.Split(request.File, "/")
	message := &models.SoUploadHistory{
		RequestId:       ctx.Value("RequestId").(string),
		BulkCode:        "SO-" + strconv.Itoa(user.AgentID) + "-" + fmt.Sprint(now.Unix()),
		FileName:        url[len(url)-1],
		FilePath:        request.File,
		AgentId:         &agentId,
		AgentName:       getAgentResult.Agent.Name,
		UploadedBy:      &userId,
		UploadedByName:  user.FirstName + " " + user.LastName,
		UploadedByEmail: user.UserEmail,
		Status:          constants.UPLOAD_STATUS_HISTORY_IN_PROGRESS,
		UpdatedBy:       &userId,
		UpdatedByName:   user.FirstName + " " + user.LastName,
		UpdatedByEmail:  user.UserEmail,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	soUploadHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
	go u.soUploadHistoriesRepository.Insert(message, ctx, soUploadHistoryJourneysResultChan)
	soUploadHistoryJourneysResult := <-soUploadHistoryJourneysResultChan

	if soUploadHistoryJourneysResult.Error != nil {
		errorLogData := helper.WriteLog(soUploadHistoryJourneysResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	keyKafka := []byte("upload")
	messageKafka := []byte(soUploadHistoryJourneysResult.SoUploadHistory.ID.Hex())

	err := u.kafkaClient.WriteToTopic(constants.UPLOAD_SO_FILE_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}

func (u *uploadUseCase) RetryUploadSO(soUploadHistoryId string, ctx context.Context) *model.ErrorLog {

	user := ctx.Value("user").(*models.UserClaims)

	getSoUploadHistoriesResultChan := make(chan *models.SoUploadHistoryChan)
	go u.soUploadHistoriesRepository.GetByID(soUploadHistoryId, false, ctx, getSoUploadHistoriesResultChan)
	getSoUploadHistoriesResult := <-getSoUploadHistoriesResultChan

	if getSoUploadHistoriesResult.Error != nil {
		errorLogData := helper.WriteLog(getSoUploadHistoriesResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	userId := int64(user.UserID)
	getSoUploadHistoriesResult.SoUploadHistory.UpdatedBy = &userId
	getSoUploadHistoriesResult.SoUploadHistory.UpdatedByName = user.FirstName + " " + user.LastName
	getSoUploadHistoriesResult.SoUploadHistory.UpdatedByEmail = user.UserEmail
	getSoUploadHistoriesResult.SoUploadHistory.Status = constants.UPLOAD_STATUS_HISTORY_IN_PROGRESS

	soUploadHistoryJourneysResultChan := make(chan *models.SoUploadHistoryChan)
	go u.soUploadHistoriesRepository.UpdateByID(soUploadHistoryId, getSoUploadHistoriesResult.SoUploadHistory, ctx, soUploadHistoryJourneysResultChan)
	sosjUploadHistoryJourneysResult := <-soUploadHistoryJourneysResultChan

	if sosjUploadHistoryJourneysResult.Error != nil {
		errorLogData := helper.WriteLog(sosjUploadHistoryJourneysResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	keyKafka := []byte("retry")
	messageKafka := []byte(soUploadHistoryId)

	err := u.kafkaClient.WriteToTopic(constants.UPLOAD_SO_FILE_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}

func (u *uploadUseCase) RetryUploadDO(sjUploadHistoryId string, ctx context.Context) *model.ErrorLog {

	user := ctx.Value("user").(*models.UserClaims)

	getSoUploadHistoriesResultChan := make(chan *models.DoUploadHistoryChan)
	go u.doUploadHistoriesRepository.GetByID(sjUploadHistoryId, false, ctx, getSoUploadHistoriesResultChan)
	getSoUploadHistoriesResult := <-getSoUploadHistoriesResultChan

	if getSoUploadHistoriesResult.Error != nil {
		errorLogData := helper.WriteLog(getSoUploadHistoriesResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	userId := int64(user.UserID)
	getSoUploadHistoriesResult.DoUploadHistory.UpdatedBy = &userId
	getSoUploadHistoriesResult.DoUploadHistory.UpdatedByName = user.FirstName + " " + user.LastName
	getSoUploadHistoriesResult.DoUploadHistory.UpdatedByEmail = user.UserEmail
	getSoUploadHistoriesResult.DoUploadHistory.Status = constants.UPLOAD_STATUS_HISTORY_IN_PROGRESS

	sjUploadHistoryJourneysResultChan := make(chan *models.DoUploadHistoryChan)
	go u.doUploadHistoriesRepository.UpdateByID(sjUploadHistoryId, getSoUploadHistoriesResult.DoUploadHistory, ctx, sjUploadHistoryJourneysResultChan)
	sjsjUploadHistoryJourneysResult := <-sjUploadHistoryJourneysResultChan

	if sjsjUploadHistoryJourneysResult.Error != nil {
		errorLogData := helper.WriteLog(sjsjUploadHistoryJourneysResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	keyKafka := []byte("retry")
	messageKafka := []byte(sjUploadHistoryId)

	err := u.kafkaClient.WriteToTopic(constants.UPLOAD_DO_FILE_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}

func (u *uploadUseCase) RetryUploadSOSJ(sosjUploadHistoryId string, ctx context.Context) *model.ErrorLog {

	user := ctx.Value("user").(*models.UserClaims)

	getSOSJUploadHistoriesResultChan := make(chan *models.GetSosjUploadHistoryResponseChan)
	go u.sosjUploadHistoriesRepository.GetByID(sosjUploadHistoryId, false, ctx, getSOSJUploadHistoriesResultChan)
	getSosjUploadHistoriesResult := <-getSOSJUploadHistoriesResultChan

	if getSosjUploadHistoriesResult.Error != nil {
		errorLogData := helper.WriteLog(getSosjUploadHistoriesResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	userId := int64(user.UserID)
	getSosjUploadHistoriesResult.SosjUploadHistories.UpdatedBy = &userId
	getSosjUploadHistoriesResult.SosjUploadHistories.UpdatedByName = user.FirstName + " " + user.LastName
	getSosjUploadHistoriesResult.SosjUploadHistories.UpdatedByEmail = user.UserEmail
	getSosjUploadHistoriesResult.SosjUploadHistories.Status = constants.UPLOAD_STATUS_HISTORY_IN_PROGRESS

	sosjUploadHistoryJourneysResultChan := make(chan *models.UploadHistoryChan)
	go u.sosjUploadHistoriesRepository.UpdateByID(sosjUploadHistoryId, &getSosjUploadHistoriesResult.SosjUploadHistories.UploadHistory, ctx, sosjUploadHistoryJourneysResultChan)
	sosjUploadHistoryJourneysResult := <-sosjUploadHistoryJourneysResultChan

	if sosjUploadHistoryJourneysResult.Error != nil {
		errorLogData := helper.WriteLog(sosjUploadHistoryJourneysResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	keyKafka := []byte("retry")
	messageKafka := []byte(sosjUploadHistoryId)

	err := u.kafkaClient.WriteToTopic(constants.UPLOAD_SOSJ_FILE_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}

func (u *uploadUseCase) GetSosjUploadErrorLogs(request *models.GetSosjUploadErrorLogsRequest, ctx context.Context) (*models.GetSosjUploadErrorLogsResponse, *model.ErrorLog) {

	getSosjUploadErrorLogsResultChan := make(chan *models.SosjUploadErrorLogsChan)
	go u.sosjUploadErrorLogsRepository.Get(request, false, ctx, getSosjUploadErrorLogsResultChan)
	getSosjUploadHistoriesResult := <-getSosjUploadErrorLogsResultChan

	if getSosjUploadHistoriesResult.Error != nil {
		return &models.GetSosjUploadErrorLogsResponse{}, getSosjUploadHistoriesResult.ErrorLog
	}

	result := models.GetSosjUploadErrorLogsResponse{
		SosjUploadErrorLogs: getSosjUploadHistoriesResult.SosjUploadErrorLogs,
		Total:               getSosjUploadHistoriesResult.Total,
	}

	return &result, nil

}

func (u *uploadUseCase) GetSosjUploadHistoryById(id string, ctx context.Context) (*models.GetSosjUploadHistoryResponse, *model.ErrorLog) {

	getSosjUploadHistoryByIdResultChan := make(chan *models.GetSosjUploadHistoryResponseChan)
	go u.sosjUploadHistoriesRepository.GetByID(id, false, ctx, getSosjUploadHistoryByIdResultChan)
	getSosjUploadHistoryByIdResult := <-getSosjUploadHistoryByIdResultChan

	if getSosjUploadHistoryByIdResult.Error != nil {
		return &models.GetSosjUploadHistoryResponse{}, getSosjUploadHistoryByIdResult.ErrorLog
	}

	return getSosjUploadHistoryByIdResult.SosjUploadHistories, nil

}

func (u *uploadUseCase) GetSosjUploadErrorLogBySosjUploadHistoryId(request *models.GetSosjUploadErrorLogsRequest, ctx context.Context) (*models.GetSosjUploadErrorLogsResponse, *model.ErrorLog) {

	getSosjUploadErrorLogsResultChan := make(chan *models.SosjUploadErrorLogsChan)
	go u.sosjUploadErrorLogsRepository.Get(request, false, ctx, getSosjUploadErrorLogsResultChan)
	getSosjUploadErrorLogsResult := <-getSosjUploadErrorLogsResultChan

	if getSosjUploadErrorLogsResult.Error != nil {
		return &models.GetSosjUploadErrorLogsResponse{}, getSosjUploadErrorLogsResult.ErrorLog
	}

	result := models.GetSosjUploadErrorLogsResponse{
		SosjUploadErrorLogs: getSosjUploadErrorLogsResult.SosjUploadErrorLogs,
		Total:               getSosjUploadErrorLogsResult.Total,
	}

	return &result, nil
}
