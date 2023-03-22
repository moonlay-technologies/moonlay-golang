package usecases

import (
	"context"
	"encoding/json"
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
	UploadDO(ctx context.Context) *model.ErrorLog
	UploadSO(request *models.UploadSORequest, ctx context.Context) *model.ErrorLog
	RetryUploadSO(soUploadHistoryId string, ctx context.Context) *model.ErrorLog
	RetryUploadSOSJ(soUploadHistoryId string, ctx context.Context) *model.ErrorLog
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
	uploadSOHistoriesRepository          mongoRepositories.UploadSOHistoriesRepositoryInterface
	sosjUploadHistoriesRepository        mongoRepositories.SOSJUploadHistoriesRepositoryInterface
	kafkaClient                          kafkadbo.KafkaClientInterface
	db                                   dbresolver.DB
	ctx                                  context.Context
}

func InitUploadUseCaseInterface(salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, orderSourceRepository repositories.OrderSourceRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, uomRepository repositories.UomRepositoryInterface, deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, deliveryOrderDetailRepository repositories.DeliveryOrderDetailRepositoryInterface, salesOrderLogRepository mongoRepositories.SalesOrderLogRepositoryInterface, salesOrderJourneysRepository mongoRepositories.SalesOrderJourneysRepositoryInterface, salesOrderDetailJourneysRepository mongoRepositories.SalesOrderDetailJourneysRepositoryInterface, userRepository repositories.UserRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, categoryRepository repositories.CategoryRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, salesOrderDetailOpenSearchRepository openSearchRepositories.SalesOrderDetailOpenSearchRepositoryInterface, uploadRepository repositories.UploadRepositoryInterface, warehouseRepository repositories.WarehouseRepositoryInterface, uploadSOHistoriesRepository mongoRepositories.UploadSOHistoriesRepositoryInterface, sosjUploadHistoriesRepository mongoRepositories.SOSJUploadHistoriesRepositoryInterface, kafkaClient kafkadbo.KafkaClientInterface, db dbresolver.DB, ctx context.Context) UploadUseCaseInterface {
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
		uploadSOHistoriesRepository:          uploadSOHistoriesRepository,
		sosjUploadHistoriesRepository:        sosjUploadHistoriesRepository,
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

	keyKafka := []byte(ctx.Value("RequestId").(string))
	messageKafka := []byte(sosjUploadHistoryJourneysResult.UploadHistory.ID.Hex())

	err := u.kafkaClient.WriteToTopic(constants.UPLOAD_SOSJ_FILE_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}

func (u *uploadUseCase) UploadDO(ctx context.Context) *model.ErrorLog {

	return nil
}

func (u *uploadUseCase) UploadSO(request *models.UploadSORequest, ctx context.Context) *model.ErrorLog {

	user := ctx.Value("user").(*models.UserClaims)

	// Check Agent By Id
	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(4, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	message := &models.SoUploadHistory{
		RequestId:       ctx.Value("RequestId").(string),
		FileName:        request.File,
		FilePath:        "upload-service/so/" + request.File,
		AgentId:         int64(user.AgentID),
		AgentName:       getAgentResult.Agent.Name,
		UploadedBy:      int64(user.UserID),
		UploadedByName:  user.FirstName + " " + user.LastName,
		UploadedByEmail: user.UserEmail,
		UpdatedBy:       int64(user.UserID),
		UpdatedByName:   user.FirstName + " " + user.LastName,
		UpdatedByEmail:  user.UserEmail,
	}

	keyKafka := []byte("upload")
	messageKafka, _ := json.Marshal(message)

	err := u.kafkaClient.WriteToTopic(constants.UPLOAD_SO_FILE_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}

func (u *uploadUseCase) RetryUploadSO(soUploadHistoryId string, ctx context.Context) *model.ErrorLog {

	user := ctx.Value("user").(*models.UserClaims)

	getUploadSOHistoriesResultChan := make(chan *models.SoUploadHistoryChan)
	go u.uploadSOHistoriesRepository.GetByID(soUploadHistoryId, false, ctx, getUploadSOHistoriesResultChan)
	getUploadSOHistoriesResult := <-getUploadSOHistoriesResultChan

	if getUploadSOHistoriesResult.Error != nil {
		return getUploadSOHistoriesResult.ErrorLog
	}

	getUploadSOHistoriesResult.SoUploadHistory.UpdatedBy = int64(user.UserID)
	getUploadSOHistoriesResult.SoUploadHistory.UpdatedByName = user.FirstName + " " + user.LastName
	getUploadSOHistoriesResult.SoUploadHistory.UpdatedByEmail = user.UserEmail

	keyKafka := []byte("retry")
	messageKafka, _ := json.Marshal(getUploadSOHistoriesResult.SoUploadHistory)

	err := u.kafkaClient.WriteToTopic(constants.UPLOAD_SO_FILE_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}

func (u *uploadUseCase) RetryUploadSOSJ(sosjUploadHistoryId string, ctx context.Context) *model.ErrorLog {

	user := ctx.Value("user").(*models.UserClaims)

	getSOSJUploadHistoriesResultChan := make(chan *models.UploadHistoryChan)
	go u.sosjUploadHistoriesRepository.GetByID(sosjUploadHistoryId, false, ctx, getSOSJUploadHistoriesResultChan)
	getUploadSOHistoriesResult := <-getSOSJUploadHistoriesResultChan

	if getUploadSOHistoriesResult.Error != nil {
		errorLogData := helper.WriteLog(getUploadSOHistoriesResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	userId := int64(user.UserID)
	getUploadSOHistoriesResult.UploadHistory.UpdatedBy = &userId
	getUploadSOHistoriesResult.UploadHistory.UpdatedByName = user.FirstName + " " + user.LastName
	getUploadSOHistoriesResult.UploadHistory.UpdatedByEmail = user.UserEmail
	getUploadSOHistoriesResult.UploadHistory.Status = constants.UPLOAD_STATUS_HISTORY_IN_PROGRESS

	sosjUploadHistoryJourneysResultChan := make(chan *models.UploadHistoryChan)
	go u.sosjUploadHistoriesRepository.UpdateByID(sosjUploadHistoryId, getUploadSOHistoriesResult.UploadHistory, ctx, sosjUploadHistoryJourneysResultChan)
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
