package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/app/middlewares"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/app/repositories"
	mongoRepositories "order-service/app/repositories/mongod"
	"order-service/global/utils/helper"
	kafkadbo "order-service/global/utils/kafka"
	"order-service/global/utils/redisdb"
	"strconv"
	"strings"
	"time"

	"github.com/bxcodec/dbresolver"
)

type UploadDOFileConsumerHandlerInterface interface {
	ProcessMessage()
}

type uploadDOFileConsumerHandler struct {
	kafkaClient                 kafkadbo.KafkaClientInterface
	uploadRepository            repositories.UploadRepositoryInterface
	requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface
	requestValidationRepository repositories.RequestValidationRepositoryInterface
	sjUploadHistoriesRepository mongoRepositories.DoUploadHistoriesRepositoryInterface
	sjUploadErrorLogsRepository mongoRepositories.DoUploadErrorLogsRepositoryInterface
	warehouseRepository         repositories.WarehouseRepositoryInterface
	salesOrderRepository        repositories.SalesOrderRepositoryInterface
	salesOrderDetailRepository  repositories.SalesOrderDetailRepositoryInterface
	deliveryOrderRepository     repositories.DeliveryOrderRepositoryInterface
	orderStatusRepository       repositories.OrderStatusRepositoryInterface
	ctx                         context.Context
	redisdb                     redisdb.RedisInterface
	args                        []interface{}
	db                          dbresolver.DB
}

func InitUploadDOFileConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, uploadRepository repositories.UploadRepositoryInterface, requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, requestValidationRepository repositories.RequestValidationRepositoryInterface, sjUploadHistoriesRepository mongoRepositories.DoUploadHistoriesRepositoryInterface, sjUploadErrorLogsRepository mongoRepositories.DoUploadErrorLogsRepositoryInterface, warehouseRepository repositories.WarehouseRepositoryInterface, salesOrderRepository repositories.SalesOrderRepositoryInterface, salesOrderDetailRepository repositories.SalesOrderDetailRepositoryInterface, deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, orderStatusRepository repositories.OrderStatusRepositoryInterface, ctx context.Context, redisdb redisdb.RedisInterface, args []interface{}, db dbresolver.DB) UploadDOFileConsumerHandlerInterface {
	return &uploadDOFileConsumerHandler{
		kafkaClient:                 kafkaClient,
		uploadRepository:            uploadRepository,
		requestValidationMiddleware: requestValidationMiddleware,
		requestValidationRepository: requestValidationRepository,
		sjUploadHistoriesRepository: sjUploadHistoriesRepository,
		sjUploadErrorLogsRepository: sjUploadErrorLogsRepository,
		warehouseRepository:         warehouseRepository,
		salesOrderRepository:        salesOrderRepository,
		salesOrderDetailRepository:  salesOrderDetailRepository,
		deliveryOrderRepository:     deliveryOrderRepository,
		orderStatusRepository:       orderStatusRepository,
		ctx:                         ctx,
		redisdb:                     redisdb,
		args:                        args,
		db:                          db,
	}
}

func (c *uploadDOFileConsumerHandler) ProcessMessage() {
	fmt.Println("process ", constants.UPLOAD_DO_FILE_TOPIC)
	topic := c.args[1].(string)
	groupID := c.args[2].(string)
	reader := c.kafkaClient.SetConsumerGroupReader(topic, groupID)

	for {
		m, err := reader.ReadMessage(c.ctx)
		if err != nil {
			break
		}

		fmt.Printf("message do at topic/partition/offset %v/%v/%v \n", m.Topic, m.Partition, m.Offset)
		now := time.Now()

		sjUploadHistoryId := m.Value
		var key = string(m.Key[:])

		sjUploadHistoryJourneysResultChan := make(chan *models.DoUploadHistoryChan)
		go c.sjUploadHistoriesRepository.GetByID(string(sjUploadHistoryId), false, c.ctx, sjUploadHistoryJourneysResultChan)
		sjUploadHistoryJourneysResult := <-sjUploadHistoryJourneysResultChan
		if sjUploadHistoryJourneysResult.Error != nil {
			fmt.Println(sjUploadHistoryJourneysResult.Error.Error())
		}

		message := sjUploadHistoryJourneysResult.DoUploadHistory
		file, err := c.uploadRepository.ReadFile(message.FilePath)

		if err != nil {
			c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
			continue
		}

		parseFile := string(file)
		data := strings.Split(parseFile, "\n")
		if len(data) < 1 || strings.ReplaceAll(data[0], "\r", "") != "IDDistributor,NoOrder,TanggalSJ,NoSJ,Catatan,CatatanInternal,NamaSupir,PlatNo,KodeMerk,NamaMerk,KodeProduk,NamaProduk,QTYShip,Unit,KodeGudang" {
			c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
			continue
		}
		totalRows := int64(len(data) - 1)

		message.TotalRows = &totalRows
		c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_SUCCESS)

		results := []map[string]string{}
		for i, v := range data {
			if i == len(data)-1 {
				break
			}
			if i == 0 {
				continue
			}

			row := strings.Split(v, "\"")
			for j := 1; j < len(row); j = j + 2 {
				cell := strings.Split(row[j], ",")
				row[j] = strings.Join(cell, "")
			}
			v = strings.Join(row, "")

			headers := strings.Split(data[0], ",")
			line := strings.Split(v, ",")
			uploadSjField := map[string]string{}
			for j, y := range line {
				uploadSjField[strings.ReplaceAll(headers[j], "\r", "")] = strings.ReplaceAll(y, "\r", "")
			}
			results = append(results, uploadSjField)
		}

		var uploadDOFields []*models.UploadDOField
		for i, v := range results {
			mandatoryError := c.requestValidationMiddleware.UploadMandatoryValidation([]*models.TemplateRequest{
				{
					Field: "IDDistributor",
					Value: v["IDDistributor"],
				},
				{
					Field: "NoOrder",
					Value: v["NoOrder"],
				},
				{
					Field: "TanggalSJ",
					Value: v["TanggalSJ"],
				},
				{
					Field: "NoSJ",
					Value: v["NoSJ"],
				},
				{
					Field: "KodeMerk",
					Value: v["KodeMerk"],
				},
				{
					Field: "KodeProduk",
					Value: v["KodeProduk"],
				},
				{
					Field: "QTYShip",
					Value: v["QTYShip"],
				},
				{
					Field: "Unit",
					Value: v["Unit"],
				},
			})
			if len(mandatoryError) >= 1 {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := mandatoryError
					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, "", errors, &now, v)
					continue
				}
			}

			intType := []*models.TemplateRequest{
				{
					Field: "IDDistributor",
					Value: v["IDDistributor"],
				},
				{
					Field: "KodeMerk",
					Value: v["KodeMerk"],
				},
				{
					Field: "QTYShip",
					Value: v["QTYShip"],
				},
			}
			intTypeResult, intTypeError := c.requestValidationMiddleware.UploadIntTypeValidation(intType)
			if len(intTypeError) >= 1 {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := intTypeError
					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, "", errors, &now, v)
					continue
				}
			}

			agentId, _ := strconv.Atoi(v["IDDistributor"])
			getWarehouseResultChan := make(chan *models.WarehouseChan)
			if len(v["KodeGudang"]) > 0 {
				go c.warehouseRepository.GetByCode(v["KodeGudang"], false, c.ctx, getWarehouseResultChan)
			} else {
				go c.warehouseRepository.GetByAgentID(agentId, false, c.ctx, getWarehouseResultChan)
			}
			getWarehouseResult := <-getWarehouseResultChan
			if getWarehouseResult.Error != nil && getWarehouseResult.ErrorLog.StatusCode == http.StatusNotFound {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{getWarehouseResult.Error.Error()}
					if getWarehouseResult.ErrorLog.StatusCode == http.StatusNotFound {
						errors = []string{fmt.Sprintf("Gudang dengan Kode %s Tidak Ditemukan pada Distributor %s. Silahkan gunakan Kode Gudang yang lain.", v["KodeGudang"], message.AgentName)}
					}
					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, "", errors, &now, v)
					continue
				}
			}
			warehouseName := getWarehouseResult.Warehouse.Name

			if getWarehouseResult.Warehouse.OwnerID != agentId {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {

					errors := []string{fmt.Sprintf("Gudang Utama pada Distributor %s Tidak Ditemukan. Silahkan Periksa Kode Gudang Utama Anda atau gunakan Kode Gudang yang lain.", message.AgentName)}

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			if getWarehouseResult.Warehouse.Status != 1 {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{fmt.Sprintf("Gudang dengan Kode %s sudah Tidak Aktif pada Distributor %s. Silahkan gunakan Kode Gudang yang lain.", v["KodeGudang"], message.AgentName)}
					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, "", errors, &now, v)
					continue
				}
			}

			if intTypeResult["QTYShip"] < 1 {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{"Quantity harus lebih dari 0"}

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			mustActiveField := []*models.MustActiveRequest{
				helper.GenerateMustActive("agents", "IDDistributor", intTypeResult["IDDistributor"], "active"),
			}
			mustActiveError := c.requestValidationMiddleware.UploadMustActiveValidation(mustActiveField)
			if len(mustActiveError) >= 1 {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := mustActiveError

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			salesOrderRedisKey := fmt.Sprintf("%s:%s", constants.SALES_ORDER_BY_CODE, v["NoOrder"])
			_, err := c.redisdb.Client().Del(c.ctx, salesOrderRedisKey).Result()
			if err != nil {
				fmt.Println(err)
			}

			getSalesOrderResultChan := make(chan *models.SalesOrderChan)
			go c.salesOrderRepository.GetByCode(v["NoOrder"], false, c.ctx, getSalesOrderResultChan)
			getSalesOrderResult := <-getSalesOrderResultChan
			if getSalesOrderResult.Error != nil {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{getSalesOrderResult.Error.Error()}
					if getSalesOrderResult.ErrorLog.StatusCode == http.StatusNotFound {
						errors = []string{fmt.Sprintf("Nomer Order = %s Tidak Ditemukan. Silahkan gunakan Nomer Order yang lain.", v["NoOrder"])}
					}
					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			getOrderStatusSOResultChan := make(chan *models.OrderStatusChan)
			go c.orderStatusRepository.GetByID(getSalesOrderResult.SalesOrder.OrderStatusID, false, c.ctx, getOrderStatusSOResultChan)
			getOrderStatusResult := <-getOrderStatusSOResultChan

			if getOrderStatusResult.Error != nil {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{getSalesOrderResult.Error.Error()}

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			if getOrderStatusResult.OrderStatus.Name != "open" && getOrderStatusResult.OrderStatus.Name != "partial" {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{fmt.Sprintf("Status Sales Order %s. Mohon sesuaikan kembali.", getOrderStatusResult.OrderStatus.Name)}

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			brandId, _ := strconv.Atoi(v["KodeMerk"])
			if getSalesOrderResult.SalesOrder.BrandID != brandId {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{fmt.Sprintf("Kode Merek = %s di Surat Jalan BERBEDA dengan Kode Merek yang terdapat pada No Order = %s. Silahkan menggunakan Kode Merek yang sama untuk Surat Jalan dan Order.", v["KodeMerk"], v["NoOrder"])}

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			salesOrderDetailRedisKey := fmt.Sprintf("%s:%d", constants.SALES_ORDER_DETAIL_BY_SOID_SKU, getSalesOrderResult.SalesOrder.ID)
			_, err = c.redisdb.Client().Del(c.ctx, salesOrderDetailRedisKey).Result()
			if err != nil {
				fmt.Println(err)
			}

			getSODetailBySoIdAndSkuResultChan := make(chan *models.SalesOrderDetailsChan)
			go c.salesOrderDetailRepository.GetBySOIDAndSku(getSalesOrderResult.SalesOrder.ID, v["KodeProduk"], false, c.ctx, getSODetailBySoIdAndSkuResultChan)
			getSODetailBySoIdAndSkuResult := <-getSODetailBySoIdAndSkuResultChan

			if getSODetailBySoIdAndSkuResult.Error != nil {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break

				} else {
					errors := []string{getSODetailBySoIdAndSkuResult.Error.Error()}
					if getSODetailBySoIdAndSkuResult.ErrorLog.StatusCode == http.StatusNotFound {
						errors = []string{fmt.Sprintf("Kode Produk = %s pada data SJ Tidak ditemukan di No. Order = %s. Silahkan gunakan kode product yang lain.", v["KodeProduk"], v["NoOrder"])}
					}

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			isUomExist := false
			soDetail := &models.SalesOrderDetail{}
			for _, y := range getSODetailBySoIdAndSkuResult.SalesOrderDetails {
				if y.UomCode == v["Unit"] {
					isUomExist = true
					soDetail = y
					break
				}
			}

			if !isUomExist {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break

				} else {
					errors := []string{fmt.Sprintf("Satuan produk untuk SKU %s di data order dan file upload surat jalan tidak sesuai. Mohon sesuaikan kembali", v["KodeProduk"])}

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			qtyShip, _ := strconv.Atoi(v["QTYShip"])
			if qtyShip > soDetail.Qty {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break

				} else {
					errors := []string{fmt.Sprintf("QTY SJ untuk Kode Produk = %s Lebih Besar dari QTY Order %s. Silahkan QTY SJ disesuaikan kembali.", v["KodeProduk"], v["NoOrder"])}

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			if qtyShip > soDetail.ResidualQty {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break

				} else {
					errors := []string{fmt.Sprintf("QTY SJ untuk Kode Produk = %s Melebihi Sisa QTY Order %s. Silahkan menyesuaikan kembali QTY produk.", v["KodeProduk"], v["NoOrder"])}

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			deliveryOrderRedisKey := fmt.Sprintf("%s:%s", constants.DELIVERY_ORDER_BY_DO_REF_CODE, v["NoSJ"])
			_, err = c.redisdb.Client().Del(c.ctx, deliveryOrderRedisKey).Result()
			if err != nil {
				fmt.Println(err)
			}

			getDeliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
			go c.deliveryOrderRepository.GetByDoRefCode(v["NoSJ"], false, c.ctx, getDeliveryOrderResultChan)
			getDeliveryOrderResult := <-getDeliveryOrderResultChan
			if getDeliveryOrderResult.Error != nil && getDeliveryOrderResult.ErrorLog.StatusCode != http.StatusNotFound {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{getDeliveryOrderResult.Error.Error()}

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			if getDeliveryOrderResult.Total > 0 {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{fmt.Sprintf("No. Surat Jalan = %s Sudah Terpakai pada Distributor %s, silahkan gunakan No. Surat Jalan lain.", v["NoSJ"], getDeliveryOrderResult.DeliveryOrder.AgentName)}

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			tanggalSJ, err := helper.ParseDDYYMMtoYYYYMMDD(v["TanggalSJ"])
			if err != nil {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{fmt.Sprintf("Format Tanggal Order = %s Salah, silahkan sesuaikan dengan format DD-MMM-YYYY, contoh 15/12/2021", v["TanggalSJ"])}

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)
					continue
				}
			}

			nowWIB := time.Now().UTC().Add(7 * time.Hour)
			duration := time.Hour*time.Duration(nowWIB.Hour()) + time.Minute*time.Duration(nowWIB.Minute()) + time.Second*time.Duration(nowWIB.Second()) + time.Nanosecond*time.Duration(nowWIB.Nanosecond())
			parseTangalSJ, _ := time.Parse(constants.DATE_FORMAT_COMMON, tanggalSJ)
			tanggalOrder, _ := time.Parse(time.RFC3339, getSalesOrderResult.SalesOrder.SoDate)

			if parseTangalSJ.Add(duration + 1*time.Minute).Before(tanggalOrder.Add(duration)) {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)

					break
				} else {
					errors := []string{"Tanggal Surat Jalan TIDAK BOLEH LEBIH AWAL dari Tanggal Order. Silahkan disesuaikan kembali"}

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)

					continue
				}
			}

			if parseTangalSJ.Add(duration - 1*time.Minute).After(nowWIB) {
				if key == "retry" {
					c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)

					break
				} else {
					errors := []string{"Tanggal Surat Jalan TIDAK BOLEH MELEBIHI dari Tanggal Pembuatan Surat Jalan. Silahkan disesuaikan kembali"}

					c.createSjUploadErrorLog(i+2, v["IDDistributor"], message.ID.Hex(), message.RequestId, message.AgentName, message.BulkCode, warehouseName, errors, &now, v)

					continue
				}
			}

			var uploadDOField models.UploadDOField
			uploadDOField.TanggalSJ = tanggalSJ
			uploadDOField.UploadDOFieldMap(v, int(*message.UploadedBy), message.ID.Hex())
			uploadDOField.BulkCode = message.BulkCode
			uploadDOField.ErrorLine = i + 2
			uploadDOField.UploadType = key

			uploadDOFields = append(uploadDOFields, &uploadDOField)
		}

		keyKafka := []byte(message.RequestId)
		messageKafka, _ := json.Marshal(uploadDOFields)

		err = c.kafkaClient.WriteToTopic(constants.UPLOAD_DO_ITEM_TOPIC, keyKafka, messageKafka)

		if err != nil {
			c.updateSjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
			continue
		}

	}

	if err := reader.Close(); err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
}

func (c *uploadDOFileConsumerHandler) createSjUploadErrorLog(errorLine int, agentId, sjUploadHistoryId, requestId, agentName, bulkCode, warehouseName string, errors []string, now *time.Time, item map[string]string) {
	sjUploadErrorLog := &models.DoUploadErrorLog{}
	sjUploadErrorLog.DoUploadErrorLogsMap(errorLine, sjUploadHistoryId, requestId, bulkCode, errors, now)
	rowData := &models.RowDataDoUploadErrorLog{}
	rowData.RowDataDoUploadErrorLogMap(item, agentName, warehouseName)
	sjUploadErrorLog.RowData = *rowData

	sjUploadErrorLogResultChan := make(chan *models.DoUploadErrorLogChan)
	go c.sjUploadErrorLogsRepository.Insert(sjUploadErrorLog, c.ctx, sjUploadErrorLogResultChan)
}

func (c *uploadDOFileConsumerHandler) updateSjUploadHistories(message *models.DoUploadHistory, status string) {
	message.UpdatedAt = time.Now()
	message.Status = status
	sjUploadHistoryJourneysResultChan := make(chan *models.DoUploadHistoryChan)
	go c.sjUploadHistoriesRepository.UpdateByID(message.ID.Hex(), message, c.ctx, sjUploadHistoryJourneysResultChan)
}
