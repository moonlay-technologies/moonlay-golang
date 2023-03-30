package repositories

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"order-service/app/models"
	"order-service/global/utils/helper"
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
	UploadFile(data *bytes.Buffer, date string, fileType string) error
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

func (r *uploadRepository) UploadFile(data *bytes.Buffer, date string, fileType string) error {

	session, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(
			"ASIAUHX63DBTOXGQIU4F",
			"ADYRZWgcYMl3j9U7rmX/Sres0UOC+z+R7WJEUGtX",
			"IQoJb3JpZ2luX2VjEGIaDmFwLXNvdXRoZWFzdC0xIkYwRAIgbb9Cu0hH/ZFIiApK7GKvR4Qsz4+dgKXWSfH34A58kZ0CIB8zYxOQ2gtlEn7A+SoV7UT4lwsjpwR/C/Y1mVKpwinXKvcCCDsQABoMMjkxNTE4NDg2NjMwIgw3hSmJLllh3sToJOgq1AJf3aus62JEYsyutdpVTt5sbx871q2E+ahnx6apIB7mwlEAeJ6pJXXtT3Ud3OF3XVVf9jg7xPMAmJ+s7ZskNAtzAyk8vTr5oT5b8lwi4Hpw4P/Xuy574E0lOHrsY6dbOgBHJmPSDHwLwjmMgdSu3xpn0SjkxzCrMvfYSWZ137nENcu/u2WWxYjlaoSpvAxu2y14cNYDSqT68CjfIehSrjdGm7QT2HOWxUy1uihRZQAqsignj17xw3OY3hRfLP3gSlRhSmCJKXL9unQScpLzd+kusmivLRfMcmS7PwRskuSkZVYxKyTB6GTuNHbvj/BeBzVnDOwK4JX71DGQkd+aNBIUjiAIdM7aEGW01Ydw0ajNYS9JuNmZ/S2wkk36HMXmSlgnGm1M04oueiAcTjhN918gnHXUC9tX2CluGwK9Y4DGJAQLh1QVj6VUsR2+NrkdtpEb76m7MPeOiaEGOqgB2EZc9o/uCCOO5e1HYIzquqabkduUKZBpiPxD5e3bCKzR39qLSMHSel7ftinCRUq8sPKsLriEFIulfXs9o+ukYcTeMDMdKCnFY5Sz159XHf0LLPhCehy+RaTYmqEVqkYF3d5xYcKrRzpBzra+zO0Q1MWm6+/KZ3ufhBRIxODpzx9JJF7JBizp7PKJh/Q1uGNVjvHMotd441Mr7rP5a9Y2XwCfgyYLiQvV",
		)})
	if err != nil {
		log.Fatal(err)
	}

	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String("lambda-upload-srv"),
		Key:                aws.String(fmt.Sprintf("lambda-upload-srv/order-service/export-delivery-orders/delivery_order_%s.%s", date, fileType)),
		ACL:                aws.String("public-read"),
		Body:               bytes.NewReader(data.Bytes()),
		ContentLength:      aws.Int64(int64(len(data.Bytes()))),
		ContentType:        aws.String("csv"),
		ContentDisposition: aws.String("attachment"),
	})
	return err
}
