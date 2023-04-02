package repositories

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"order-service/app/models"
	"order-service/app/models/constants"
	"order-service/global/utils/helper"
	"os"
	"strconv"

	"github.com/bxcodec/dbresolver"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type UploadRepositoryInterface interface {
	ReadFile(url string) ([]byte, error)
	GetSosjRowData(agentId, storeCode, brandId, productSku, warehouseId, salesmanId, addressId string, resultChan chan *models.RowDataSosjUploadErrorLogChan)
	UploadFile(data *bytes.Buffer, fileName string, fileType string) error
}

type uploadRepository struct {
	db dbresolver.DB
}

func InitUploadRepository(db dbresolver.DB) UploadRepositoryInterface {
	return &uploadRepository{
		db: db,
	}
}

func (r *uploadRepository) ReadFile(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return body, nil
}

func (r *uploadRepository) GetSosjRowData(agentId, storeCode, brandId, productSku, warehouseId, salesmanId, addressId string, resultChan chan *models.RowDataSosjUploadErrorLogChan) {

	rowData := models.RowDataSosjUploadErrorLog{}
	response := &models.RowDataSosjUploadErrorLogChan{}

	agentIdInt, _ := strconv.Atoi(agentId)
	brandIdInt, _ := strconv.Atoi(brandId)
	warehouseIdInt, _ := strconv.Atoi(warehouseId)
	salesmanIdInt, _ := strconv.Atoi(salesmanId)
	addressIdInt, _ := strconv.Atoi(addressId)

	err := r.db.QueryRow("SELECT (SELECT agents.name FROM agents WHERE id = ?)  AS agent_name,(SELECT stores.name FROM stores WHERE IF((SELECT COUNT(store_code) FROM stores WHERE stores.store_code = ?), stores.store_code = ?, stores.alias_code = ?)) AS store_name, (SELECT brands.name FROM brands WHERE id = ?) AS brand_name, (SELECT products.productName FROM products WHERE IF((SELECT COUNT(SKU) FROM products WHERE SKU = ?), products.SKU = ?, products.aliasSku = ?)) AS product_name, (SELECT warehouses.name FROM warehouses WHERE id = ?) AS wh_name, (SELECT salesmans.name FROM salesmans WHERE id = ?) AS sales_name, (SELECT store_addresses.address FROM store_addresses WHERE id = ?) AS address", agentIdInt, storeCode, storeCode, storeCode, brandIdInt, productSku, productSku, productSku, warehouseIdInt, salesmanIdInt, addressIdInt).Scan(&rowData.AgentName, &rowData.StoreName, &rowData.BrandName, &rowData.ProductName, &rowData.WhName, &rowData.SalesName, &rowData.Address)

	if err != nil {
		errorLogData := helper.WriteLog(err, 500, "Something went wrong, please try again later")
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	} else {
		response.RowDataSosjUploadErrorLog = &rowData
		resultChan <- response
		return
	}
}

func (r *uploadRepository) UploadFile(data *bytes.Buffer, fileName string, fileType string) error {

	AccessKeyID := os.Getenv("S3_ACCESS_KEY_ID")
	SecretAccessKey := os.Getenv("S3_SECRET_ACCESS_KEY")
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(
			AccessKeyID,
			SecretAccessKey,
			"",
		)})
	if err != nil {
		log.Fatal(err)
	}

	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(constants.S3_EXPORT_BUCKET),
		Key:                aws.String(fmt.Sprintf("%s/%s.%s", constants.S3_EXPORT_DO_PATH, fileName, fileType)),
		ACL:                aws.String(constants.S3_EXPORT_ACL),
		Body:               bytes.NewReader(data.Bytes()),
		ContentLength:      aws.Int64(int64(len(data.Bytes()))),
		ContentType:        aws.String("csv"),
		ContentDisposition: aws.String(constants.S3_EXPORT_CONTENT_DISPOSISTION),
	})
	return err
}
