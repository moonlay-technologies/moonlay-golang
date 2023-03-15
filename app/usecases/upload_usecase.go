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

	"github.com/bxcodec/dbresolver"
)

type UploadUseCaseInterface interface {
	UploadSOSJ(request *models.UploadSOSJRequest, ctx context.Context) *model.ErrorLog
	UploadDO(ctx context.Context) *model.ErrorLog
	UploadSO(request *models.UploadSORequest, ctx context.Context) *model.ErrorLog
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
	kafkaClient                          kafkadbo.KafkaClientInterface
	db                                   dbresolver.DB
	ctx                                  context.Context
}

func InitUploadUseCaseInterface(salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, orderSourceRepository repositories.OrderSourceRepositoryInterface, agentRepository repositories.AgentRepositoryInterface, brandRepository repositories.BrandRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, productRepository repositories.ProductRepositoryInterface, uomRepository repositories.UomRepositoryInterface, deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, deliveryOrderDetailRepository repositories.DeliveryOrderDetailRepositoryInterface, salesOrderLogRepository mongoRepositories.SalesOrderLogRepositoryInterface, salesOrderJourneysRepository mongoRepositories.SalesOrderJourneysRepositoryInterface, salesOrderDetailJourneysRepository mongoRepositories.SalesOrderDetailJourneysRepositoryInterface, userRepository repositories.UserRepositoryInterface, salesmanRepository repositories.SalesmanRepositoryInterface, categoryRepository repositories.CategoryRepositoryInterface, salesOrderOpenSearchRepository openSearchRepositories.SalesOrderOpenSearchRepositoryInterface, salesOrderDetailOpenSearchRepository openSearchRepositories.SalesOrderDetailOpenSearchRepositoryInterface, uploadRepository repositories.UploadRepositoryInterface, warehouseRepository repositories.WarehouseRepositoryInterface, kafkaClient kafkadbo.KafkaClientInterface, db dbresolver.DB, ctx context.Context) UploadUseCaseInterface {
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
		kafkaClient:                          kafkaClient,
		db:                                   db,
		ctx:                                  ctx,
	}
}

func (u *uploadUseCase) UploadSOSJ(request *models.UploadSOSJRequest, ctx context.Context) *model.ErrorLog {

	user := ctx.Value("user").(*models.UserClaims)

	// Check Agent By Id
	getAgentResultChan := make(chan *models.AgentChan)
	go u.agentRepository.GetByID(4, false, ctx, getAgentResultChan)
	getAgentResult := <-getAgentResultChan

	if getAgentResult.Error != nil {
		errorLogData := helper.WriteLog(getAgentResult.Error, http.StatusInternalServerError, nil)
		return errorLogData
	}

	message := &models.UploadHistory{
		RequestId:      ctx.Value("RequestId").(string),
		FName:          request.File,
		FPath:          "upload-service/sosj/" + request.File,
		AgentId:        user.AgentID,
		AgentName:      getAgentResult.Agent.Name,
		UploadById:     user.UserID,
		UploadByEmail:  user.UserEmail,
		UploadedByName: user.FirstName + " " + user.LastName,
	}

	keyKafka := []byte(ctx.Value("RequestId").(string))
	messageKafka, _ := json.Marshal(message)

	err := u.kafkaClient.WriteToTopic(constants.UPLOAD_SOSJ_FILE_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}

func (u *uploadUseCase) UploadDO(ctx context.Context) *model.ErrorLog {

	uploadDOResultChan := make(chan *models.UploadDOFieldsChan)
	go u.uploadRepository.UploadDO("be-so-service", "upload-service/do/format-file-upload-data-DO-V2 (1).csv", "ap-southeast-1", uploadDOResultChan)
	uploadDOResult := <-uploadDOResultChan

	for _, v := range uploadDOResult.UploadDOFields {
		a, _ := json.Marshal(v)
		fmt.Println(string(a))
	}
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

	message := &models.UploadSOHistory{
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

	keyKafka := []byte(ctx.Value("RequestId").(string))
	messageKafka, _ := json.Marshal(message)

	err := u.kafkaClient.WriteToTopic(constants.UPLOAD_SO_FILE_TOPIC, keyKafka, messageKafka)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		return errorLogData
	}

	return nil
}
