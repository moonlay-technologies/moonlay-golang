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
	"strconv"
	"strings"
	"time"

	"github.com/bxcodec/dbresolver"
)

type UploadSOSJFileConsumerHandlerInterface interface {
	ProcessMessage()
}

type uploadSOSJFileConsumerHandler struct {
	kafkaClient                   kafkadbo.KafkaClientInterface
	uploadRepository              repositories.UploadRepositoryInterface
	requestValidationMiddleware   middlewares.RequestValidationMiddlewareInterface
	requestValidationRepository   repositories.RequestValidationRepositoryInterface
	sosjUploadHistoriesRepository mongoRepositories.SOSJUploadHistoriesRepositoryInterface
	sosjUploadErrorLogsRepository mongoRepositories.SosjUploadErrorLogsRepositoryInterface
	salesOrderRepository          repositories.SalesOrderRepositoryInterface
	deliveryOrderRepository       repositories.DeliveryOrderRepositoryInterface
	warehouseRepository           repositories.WarehouseRepositoryInterface
	storeRepository               repositories.StoreRepositoryInterface
	productRepository             repositories.ProductRepositoryInterface
	uomRepository                 repositories.UomRepositoryInterface
	ctx                           context.Context
	args                          []interface{}
	db                            dbresolver.DB
}

func InitUploadSOSJFileConsumerHandlerInterface(kafkaClient kafkadbo.KafkaClientInterface, uploadRepository repositories.UploadRepositoryInterface, requestValidationMiddleware middlewares.RequestValidationMiddlewareInterface, requestValidationRepository repositories.RequestValidationRepositoryInterface, sosjUploadHistoriesRepository mongoRepositories.SOSJUploadHistoriesRepositoryInterface, sosjUploadErrorLogsRepository mongoRepositories.SosjUploadErrorLogsRepositoryInterface, salesOrderRepository repositories.SalesOrderRepositoryInterface, deliveryOrderRepository repositories.DeliveryOrderRepositoryInterface, warehouseRepository repositories.WarehouseRepositoryInterface, productRepository repositories.ProductRepositoryInterface, uomRepository repositories.UomRepositoryInterface, storeRepository repositories.StoreRepositoryInterface, db dbresolver.DB, ctx context.Context, args []interface{}) UploadSOFileConsumerHandlerInterface {
	return &uploadSOSJFileConsumerHandler{
		kafkaClient:                   kafkaClient,
		uploadRepository:              uploadRepository,
		requestValidationMiddleware:   requestValidationMiddleware,
		requestValidationRepository:   requestValidationRepository,
		sosjUploadHistoriesRepository: sosjUploadHistoriesRepository,
		sosjUploadErrorLogsRepository: sosjUploadErrorLogsRepository,
		salesOrderRepository:          salesOrderRepository,
		deliveryOrderRepository:       deliveryOrderRepository,
		warehouseRepository:           warehouseRepository,
		productRepository:             productRepository,
		uomRepository:                 uomRepository,
		storeRepository:               storeRepository,
		ctx:                           ctx,
		args:                          args,
		db:                            db,
	}
}

func (c *uploadSOSJFileConsumerHandler) ProcessMessage() {
	fmt.Println("process ", constants.UPLOAD_SOSJ_FILE_TOPIC)
	topic := c.args[1].(string)
	groupID := c.args[2].(string)
	reader := c.kafkaClient.SetConsumerGroupReader(topic, groupID)

	for {
		m, err := reader.ReadMessage(c.ctx)
		if err != nil {
			break
		}

		fmt.Printf("message so at topic/partition/offset %v/%v/%v \n", m.Topic, m.Partition, m.Offset)
		now := time.Now()

		sosjUploadHistoryId := m.Value
		var key = string(m.Key[:])
		var errors []string

		sosjUploadHistoryJourneysResultChan := make(chan *models.GetSosjUploadHistoryResponseChan)
		go c.sosjUploadHistoriesRepository.GetByID(string(sosjUploadHistoryId), false, c.ctx, sosjUploadHistoryJourneysResultChan)
		sosjUploadHistoryJourneysResult := <-sosjUploadHistoryJourneysResultChan
		if sosjUploadHistoryJourneysResult.Error != nil {
			fmt.Println(sosjUploadHistoryJourneysResult.Error.Error())
		}

		message := &sosjUploadHistoryJourneysResult.SosjUploadHistories.UploadHistory
		file, err := c.uploadRepository.ReadFile(message.FilePath)

		if err != nil {
			c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
			continue
		}

		parseFile := string(file)
		data := strings.Split(parseFile, "\n")

		if len(data) < 1 || strings.ReplaceAll(data[1], "\r", "") != "Status,NoSuratJalan,TglSuratJalan,KodeTokoDBO,IDMerk,KodeProdukDBO,Qty,Unit,NamaSupir,PlatNo,KodeGudang,IDSalesman,IDAlamat,Catatan,CatatanInternal" {
			c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
			continue
		}

		totalRows := int64(len(data) - 1)

		message.TotalRows = &totalRows
		c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_SUCCESS)

		results := []map[string]string{}
		var agentId string
		for i, v := range data {

			if i == len(data)-1 {
				break
			}

			line := strings.Split(v, ",")
			if i == 0 {
				agentId = line[1]
				continue
			}
			if i == 1 {
				continue
			}

			uploadSOSJField := map[string]string{}
			uploadSOSJField["_0"] = agentId
			for j, y := range line {
				uploadSOSJField[fmt.Sprintf("_%d", j+1)] = y
			}
			results = append(results, uploadSOSJField)
		}

		var uploadSOSJFields []*models.UploadSOSJField

		noSuratJalan := []string{}
		for i, v := range results {
			rowDataSosjUploadErrorLogResultChan := make(chan *models.RowDataSosjUploadErrorLogChan)
			go c.uploadRepository.GetSosjRowData(v["_0"], v["_4"], v["_5"], v["_6"], v["_11"], v["_12"], v["_13"], rowDataSosjUploadErrorLogResultChan)
			rowDataSosjUploadErrorLogResult := <-rowDataSosjUploadErrorLogResultChan
			rowData := &models.RowDataSosjUploadErrorLog{}
			rowData.RowDataSosjUploadErrorLogMap(*rowDataSosjUploadErrorLogResult.RowDataSosjUploadErrorLog, v)

			mandatoryError := c.requestValidationMiddleware.UploadMandatoryValidation([]*models.TemplateRequest{
				{
					Field: "IDDistributor",
					Value: rowData.AgentId,
				},
				{
					Field: "Status",
					Value: rowData.SjStatus,
				},
				{
					Field: "NoSuratJalan",
					Value: rowData.SjNo,
				},
				{
					Field: "TglSuratJalan",
					Value: rowData.SjDate,
				},
				{
					Field: "KodeTokoDBO",
					Value: rowData.StoreCode,
				},
				{
					Field: "IDMerk",
					Value: rowData.BrandId,
				},
				{
					Field: "KodeProdukDBO",
					Value: rowData.ProductCode,
				},
				{
					Field: "Qty",
					Value: rowData.DeliveryQty,
				},
				{
					Field: "Unit",
					Value: rowData.ProductUnit,
				},
			})
			if len(mandatoryError) >= 1 {
				if key == "retry" {

					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)

					break
				} else {
					errors := mandatoryError

					c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, rowData.AgentName.String, message.BulkCode, errors, &now, *rowData)
					continue
				}
			}

			intType := []*models.TemplateRequest{
				{
					Field: "IDDistributor",
					Value: rowData.AgentId,
				},
				{
					Field: "IDMerk",
					Value: rowData.BrandId,
				},
				{
					Field: "Qty",
					Value: rowData.DeliveryQty,
				},
				{
					Field: "Unit",
					Value: rowData.ProductUnit,
				},
			}
			if rowData.WhCode != "" {
				intType = append(intType, &models.TemplateRequest{
					Field: "KodeGudang",
					Value: rowData.WhCode,
				})
			}
			if rowData.SalesmanId != "" {
				intType = append(intType, &models.TemplateRequest{
					Field: "IDSalesman",
					Value: rowData.SalesmanId,
				})
			}
			intTypeResult, intTypeError := c.requestValidationMiddleware.UploadIntTypeValidation(intType)
			if len(intTypeError) >= 1 {
				if key == "retry" {

					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)

					break
				} else {
					errors := intTypeError

					c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, rowData.AgentName.String, message.BulkCode, errors, &now, *rowData)

					continue
				}
			}

			if intTypeResult["Qty"] < 1 {
				if key == "retry" {

					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors = []string{"Quantity harus lebih dari 0"}

					c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, rowData.AgentName.String, message.BulkCode, errors, &now, *rowData)
					continue
				}
			}

			getProductResultChan := make(chan *models.ProductChan)
			go c.productRepository.GetBySKU(rowData.ProductCode, false, c.ctx, getProductResultChan)
			getProductResult := <-getProductResultChan

			getUomResultChan := make(chan *models.UomChan)
			go c.uomRepository.GetByID(intTypeResult["Unit"], false, c.ctx, getUomResultChan)
			getUomResult := <-getUomResultChan

			if getUomResult.Error != nil || getProductResult.Error != nil {
				if key == "retry" {
					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{}
					if getUomResult.Error != nil {
						errors = append(errors, getUomResult.Error.Error())
					}
					if getProductResult.Error != nil {
						if getProductResult.ErrorLog.StatusCode == http.StatusNotFound {
							errors = append(errors, fmt.Sprintf("Kode SKU = %s dengan Merek %s Tidak Ditemukan. Silahkan gunakan Kode SKU yang lain", rowData.ProductCode, rowData.BrandName.String))
						} else {
							errors = append(errors, getProductResult.Error.Error())
						}
					}
					c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, rowData.AgentName.String, message.BulkCode, errors, &now, *rowData)
					continue
				}
			}

			mustActiveField := []*models.MustActiveRequest{
				helper.GenerateMustActive("agents", "IDDistributor", intTypeResult["IDDistributor"], "active"),
				helper.GenerateMustActive("users", "user_id", int(*message.UploadedBy), "ACTIVE"),
				{
					Table:         "brands",
					ReqField:      "IDMerk",
					Clause:        fmt.Sprintf("id = %d AND status_active = %d", intTypeResult["IDMerk"], 1),
					Id:            intTypeResult["IDMerk"],
					CustomMessage: fmt.Sprintf("Kode Merek = %s Tidak Terdaftar pada Distributor %s. Silahkan gunakan Kode Merek yang lain", rowData.BrandId, rowData.AgentName.String),
				},
				{
					Table:         "stores",
					ReqField:      "KodeTokoDBO",
					Clause:        fmt.Sprintf("IF((SELECT COUNT(store_code) FROM stores WHERE store_code = '%s'), stores.store_code = '%s', stores.alias_code = '%s') AND status = 'active'", rowData.StoreCode, rowData.StoreCode, rowData.StoreCode),
					Id:            rowData.StoreCode,
					CustomMessage: fmt.Sprintf("Kode Toko = %s sudah Tidak Aktif. Silahkan gunakan Kode Toko yang lain", rowData.StoreCode),
				},
				{
					Table:         "products",
					ReqField:      "KodeProdukDBO",
					Clause:        fmt.Sprintf("IF((SELECT COUNT(SKU) FROM products WHERE SKU = '%s'), products.SKU = '%s', products.aliasSku = '%s') AND isActive = %d", rowData.ProductCode, rowData.ProductCode, rowData.ProductCode, 1),
					Id:            rowData.ProductCode,
					CustomMessage: fmt.Sprintf("Kode SKU = %s dengan Merek %s sudah Tidak Aktif. Silahkan gunakan Kode SKU yang lain.", rowData.ProductCode, rowData.BrandName.String),
				},
				{
					Table:         "uoms",
					ReqField:      "Unit",
					Clause:        fmt.Sprintf(constants.CLAUSE_ID_VALIDATION, intTypeResult["Unit"]),
					Id:            intTypeResult["Unit"],
					CustomMessage: fmt.Sprintf("Unit Satuan = %s pada Kode SKU %s Tidak Sesuai. Silahkan gunakan unit satuan yang lain.", rowData.ProductUnit, rowData.ProductCode),
				},
			}
			mustActiveError := c.requestValidationMiddleware.UploadMustActiveValidation(mustActiveField)
			if len(mustActiveError) >= 1 {
				if key == "retry" {
					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := mustActiveError

					c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, rowData.AgentName.String, message.BulkCode, errors, &now, *rowData)

					continue
				}
			}

			var price float64
			if getUomResult.Uom.Code.String == getProductResult.Product.UnitMeasurementSmall.String {
				price = getProductResult.Product.PriceSmall
			} else if getUomResult.Uom.Code.String == getProductResult.Product.UnitMeasurementMedium.String {
				price = getProductResult.Product.PriceMedium
			} else if getUomResult.Uom.Code.String == getProductResult.Product.UnitMeasurementBig.String {
				price = getProductResult.Product.PriceBig
			} else {
				if key == "retry" {
					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{fmt.Sprintf("Unit Satuan = %s pada Kode SKU %s Tidak Sesuai. Silahkan gunakan unit satuan yang lain.", getUomResult.Uom.Code.String, rowData.ProductCode)}
					c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, rowData.AgentName.String, message.BulkCode, errors, &now, *rowData)
					continue
				}
			}

			if price < 1 {
				if key == "retry" {
					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{fmt.Sprintf("Produk dengan Kode SKU %s Belum Ada Harga atau Harga = 0. Silahkan gunakan Kode SKU Produk yang lain.", rowData.ProductCode)}
					c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, rowData.AgentName.String, message.BulkCode, errors, &now, *rowData)
					continue
				}
			}

			getStoreResultChan := make(chan *models.StoreChan)
			go c.storeRepository.GetIdByStoreCode(rowData.StoreCode, false, c.ctx, getStoreResultChan)
			getStoreResult := <-getStoreResultChan
			if getStoreResult.Error != nil {
				if key == "retry" {
					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{getStoreResult.Error.Error()}

					c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, rowData.AgentName.String, message.BulkCode, errors, &now, *rowData)
					continue
				}
			}

			storeIdValidationResultChan := make(chan *models.RequestIdValidationChan)
			go c.requestValidationRepository.StoreIdValidation(getStoreResult.Store.ID, intTypeResult["IDDistributor"], storeIdValidationResultChan)
			storeIdValidationResult := <-storeIdValidationResultChan

			if storeIdValidationResult.Total < 1 {
				if key == "retry" {
					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{fmt.Sprintf("KodeToko = %s Tidak Terdaftar pada Distributor %s. Silahkan gunakan Kode Toko yang lain.", rowData.StoreCode, rowData.AgentName.String)}

					c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, rowData.AgentName.String, message.BulkCode, errors, &now, *rowData)
					continue
				}
			}

			if rowData.WhCode != "" {
				getWarehouseResultChan := make(chan *models.WarehouseChan)
				go c.warehouseRepository.GetByID(intTypeResult["KodeGudang"], false, c.ctx, getWarehouseResultChan)
				getWarehouseResult := <-getWarehouseResultChan

				if getWarehouseResult.Error != nil {
					if key == "retry" {
						c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
						break
					} else {
						errors := []string{getWarehouseResult.Error.Error()}
						if getWarehouseResult.ErrorLog.StatusCode == http.StatusNotFound {
							errors = []string{fmt.Sprintf("Gudang dengan Kode %s Tidak Ditemukan pada Distributor %s. Silahkan gunakan Kode Gudang yang lain.", rowData.WhCode, rowData.AgentName.String)}
						}

						c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, rowData.AgentName.String, message.BulkCode, errors, &now, *rowData)
						continue
					}
				}

				if strconv.Itoa(getWarehouseResult.Warehouse.OwnerID) != agentId {
					if key == "retry" {
						c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
						break
					} else {

						errors := []string{fmt.Sprintf("Gudang Utama pada Distributor %s Tidak Ditemukan. Silahkan Periksa Kode Gudang Utama Anda atau gunakan Kode Gudang yang lain.", rowData.AgentName.String)}

						c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, rowData.AgentName.String, message.BulkCode, errors, &now, *rowData)
						continue
					}
				}

				if getWarehouseResult.Warehouse.Status != 1 {
					if key == "retry" {
						c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
						break
					} else {
						errors := []string{fmt.Sprintf("Gudang dengan Kode %s sudah Tidak Aktif pada Distributor %s. Silahkan gunakan Kode Gudang yang lain.", v["KodeGudang"], message.AgentName)}
						c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, rowData.AgentName.String, message.BulkCode, errors, &now, *rowData)
						continue
					}
				}
			}

			if len(v["_12"]) > 0 {
				brandSalesman := make(chan *models.RequestIdValidationChan)
				go c.requestValidationRepository.BrandSalesmanValidation(intTypeResult["IDMerk"], intTypeResult["IDSalesman"], intTypeResult["IDDistributor"], brandSalesman)
				brandSalesmanResult := <-brandSalesman

				if brandSalesmanResult.Total < 1 {

					if key == "retry" {
						c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
						break
					} else {
						errors := []string{}
						errors = append(errors, fmt.Sprintf("Kode Merek = %d Tidak Terdaftar pada Distributor %s. Silahkan gunakan Kode Merek yang lain.", intTypeResult["IDMerk"], rowData.AgentName.String))
						errors = append(errors, fmt.Sprintf("ID Salesman = %d Tidak Terdaftar pada Distributor %s. Silahkan gunakan ID Salesman yang lain.", intTypeResult["IDSalesman"], rowData.AgentName.String))
						errors = append(errors, fmt.Sprintf("Salesman di Kode Toko = %s untuk Merek %s Tidak Terdaftar. Silahkan gunakan ID Salesman yang terdaftar.", rowData.StoreCode, rowData.ProductName.String))

						c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, rowData.AgentName.String, message.BulkCode, errors, &now, *rowData)

						continue
					}
				}
			}
			storeAddresses := make(chan *models.RequestIdValidationChan)
			go c.requestValidationRepository.StoreAddressesValidation(rowData.StoreCode, storeAddresses)
			storeAddressesResult := <-storeAddresses

			if storeAddressesResult.Total < 1 {
				if key == "retry" {

					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors := []string{fmt.Sprintf("Alamat Utama pada Kode Toko = %s Tidak Ditemukan. Silahkan gunakan Alamat Toko yang lain.", rowData.StoreCode)}

					c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, rowData.AgentName.String, message.BulkCode, errors, &now, *rowData)

					continue
				}
			}

			var uploadSOSJField models.UploadSOSJField
			uploadSOSJField.TglSuratJalan, err = helper.ParseDDYYMMtoYYYYMMDD(rowData.SjDate)
			if err != nil {
				if key == "retry" {

					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors = []string{fmt.Sprintf("Format Tanggal Order = %s Salah, silahkan sesuaikan dengan format DD/MM/YYYY, contoh 15/12/2021", rowData.SjDate)}

					c.createSosjUploadErrorLog(i+3, rowData.AgentId, string(sosjUploadHistoryId), message.RequestId, rowData.AgentName.String, message.BulkCode, errors, &now, *rowData)

					continue
				}
			}
			uploadSOSJField.UploadSOSJFieldMap(v, intTypeResult["IDDistributor"], int(*message.UploadedBy))

			checkIfNoSuratJalanExist := helper.InSliceString(noSuratJalan, rowData.SjNo)
			if checkIfNoSuratJalanExist {

				for i := range uploadSOSJFields {
					brandId, _ := strconv.Atoi(rowData.BrandId)
					if uploadSOSJFields[i].NoSuratJalan == rowData.SjNo && uploadSOSJFields[i].IDMerk != brandId {
						uploadSOSJFields[i].NoSuratJalan = uploadSOSJFields[i].NoSuratJalan + "-" + strconv.Itoa(uploadSOSJFields[i].IDMerk)
						uploadSOSJField.NoSuratJalan = rowData.SjNo + "-" + rowData.BrandId
						break
					} else {
						uploadSOSJField.NoSuratJalan = rowData.SjNo
					}
				}

			} else {
				uploadSOSJField.NoSuratJalan = rowData.SjNo
				noSuratJalan = append(noSuratJalan, rowData.SjNo)
			}

			uploadSOSJField.BulkCode = message.BulkCode
			uploadSOSJField.SosjUploadHistoryId = message.ID.Hex()
			uploadSOSJField.ErrorLine = i + 3
			uploadSOSJField.UploadType = key
			uploadSOSJField.RowData = *rowData
			uploadSOSJFields = append(uploadSOSJFields, &uploadSOSJField)

		}

		var finalUploadSOSJFields []*models.UploadSOSJField
		for _, v := range uploadSOSJFields {

			salesOrderResultChan := make(chan *models.SalesOrderChan)
			go c.salesOrderRepository.GetBySoRefCode(v.NoSuratJalan, true, c.ctx, salesOrderResultChan)
			salesOrderResult := <-salesOrderResultChan
			deliveryOrderResultChan := make(chan *models.DeliveryOrderChan)
			go c.deliveryOrderRepository.GetByDoRefCode(v.NoSuratJalan, true, c.ctx, deliveryOrderResultChan)
			deliveryOrderResult := <-deliveryOrderResultChan

			if (deliveryOrderResult.Error != nil || salesOrderResult.Error != nil) && (deliveryOrderResult.ErrorLog.StatusCode != http.StatusNotFound || salesOrderResult.ErrorLog.StatusCode != http.StatusNotFound) {
				if key == "retry" {

					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {

					errors = []string{}
					if deliveryOrderResult.Error != nil {
						fmt.Println(deliveryOrderResult.Error.Error())
						errors = append(errors, deliveryOrderResult.Error.Error())
					}
					if salesOrderResult.Error != nil {
						fmt.Println(salesOrderResult.Error.Error())
						errors = append(errors, salesOrderResult.Error.Error())
					}

					c.createSosjUploadErrorLog(v.ErrorLine, v.RowData.AgentId, string(sosjUploadHistoryId), message.RequestId, v.RowData.AgentName.String, message.BulkCode, errors, &now, v.RowData)

					continue
				}
			}

			if deliveryOrderResult.Total > 0 || salesOrderResult.Total > 0 {
				if key == "retry" {

					c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
					break
				} else {
					errors = []string{}
					if deliveryOrderResult.Total > 0 {
						errors = append(errors, fmt.Sprintf("No. Surat Jalan = %s Sudah Terpakai pada Distributor %s, silahkan gunakan No. Surat Jalan lain.", v.NoSuratJalan, v.RowData.AgentName.String))
					}
					if salesOrderResult.Total > 0 {
						errors = append(errors, fmt.Sprintf("No. Order = %s Sudah Terpakai di Distributor %s. Silahkan gunakan No. Order yang lain", v.NoSuratJalan, v.RowData.AgentName.String))
					}

					c.createSosjUploadErrorLog(v.ErrorLine, v.RowData.AgentId, string(sosjUploadHistoryId), message.RequestId, v.RowData.AgentName.String, message.BulkCode, errors, &now, v.RowData)

					continue
				}
			}

			finalUploadSOSJFields = append(finalUploadSOSJFields, v)
		}

		if len(finalUploadSOSJFields) > 0 {
			keyKafka := []byte(message.RequestId)
			messageKafka, _ := json.Marshal(finalUploadSOSJFields)
			err = c.kafkaClient.WriteToTopic(constants.UPLOAD_SOSJ_ITEM_TOPIC, keyKafka, messageKafka)

			if err != nil {
				c.updateSosjUploadHistories(message, constants.UPLOAD_STATUS_HISTORY_FAILED)
				continue
			}
		}

	}

	if err := reader.Close(); err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
}

func (c *uploadSOSJFileConsumerHandler) createSosjUploadErrorLog(errorLine int, agentId, sosjUploadHistoryId, requestId, agentName, bulkCode string, errors []string, now *time.Time, rowData models.RowDataSosjUploadErrorLog) {
	sosjUploadErrorLog := &models.SosjUploadErrorLog{}
	sosjUploadErrorLog.SosjUploadErrorLogsMap(errorLine, agentId, sosjUploadHistoryId, requestId, agentName, bulkCode, errors, now)
	sosjUploadErrorLog.RowData = rowData

	sosjUploadErrorLogResultChan := make(chan *models.SosjUploadErrorLogChan)
	go c.sosjUploadErrorLogsRepository.Insert(sosjUploadErrorLog, c.ctx, sosjUploadErrorLogResultChan)
}

func (c *uploadSOSJFileConsumerHandler) updateSosjUploadHistories(message *models.UploadHistory, status string) {
	message.UpdatedAt = time.Now()
	message.Status = status
	sosjUploadHistoryJourneysResultChan := make(chan *models.UploadHistoryChan)
	go c.sosjUploadHistoriesRepository.UpdateByID(message.ID.Hex(), message, c.ctx, sosjUploadHistoryJourneysResultChan)
}
