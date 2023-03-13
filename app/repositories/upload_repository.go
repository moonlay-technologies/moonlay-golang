package repositories

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"order-service/app/models"
	"order-service/global/utils/helper"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type UploadRepositoryInterface interface {
	UploadSOSJ(bucket, object, region string, user_id int, resultChan chan *models.UploadSOSJFieldChan)
	UploadDO(bucket, object, region string, resultChan chan *models.UploadDOFieldsChan)
	UploadSO(bucket, object, region string, resultChan chan *models.UploadSOFieldsChan)
}

type upload struct {
	requestValidationRepository RequestValidationRepositoryInterface
}

func InitUploadRepository(requestValidationRepository RequestValidationRepositoryInterface) UploadRepositoryInterface {
	return &upload{
		requestValidationRepository: requestValidationRepository,
	}
}

func (r *upload) UploadSOSJ(bucket, object, region string, user_id int, resultChan chan *models.UploadSOSJFieldChan) {
	response := &models.UploadSOSJFieldChan{}

	var errors []string

	var idDistributor int
	resultsWithHeader, err := r.ReadFile(bucket, object, region, s3.FileHeaderInfoUse)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	for _, v := range resultsWithHeader {
		for k2, v2 := range v {
			if v2 == "NoSuratJalan" {
				idDistributor, _ = strconv.Atoi(k2)
			}
		}
	}

	var uploadSOSJFields []*models.UploadSOSJField
	results, err := r.ReadFile(bucket, object, region, s3.FileHeaderInfoIgnore)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	noSuratJalan := []string{}
	for _, v := range results {
		if v["_1"] != "Status" {
			mandatoryError := uploadMandatoryValidation([]*models.TemplateRequest{
				{
					Field: "Status",
					Value: v["_1"],
				},
				{
					Field: "NoSuratJalan",
					Value: v["_2"],
				},
				{
					Field: "TglSuratJalan",
					Value: v["_3"],
				},
				{
					Field: "KodeTokoDBO",
					Value: v["_4"],
				},
				{
					Field: "IDMerk",
					Value: v["_5"],
				},
				{
					Field: "KodeProdukDBO",
					Value: v["_6"],
				},
				{
					Field: "Qty",
					Value: v["_7"],
				},
				{
					Field: "Unit",
					Value: v["_8"],
				},
			})
			if len(mandatoryError) > 1 {
				errors = append(errors, mandatoryError...)
				continue
			}

			intType := []*models.TemplateRequest{
				{
					Field: "KodeTokoDBO",
					Value: v["_4"],
				},
				{
					Field: "IDMerk",
					Value: v["_5"],
				},
				{
					Field: "KodeProdukDBO",
					Value: v["_6"],
				},
				{
					Field: "Qty",
					Value: v["_7"],
				},
				{
					Field: "Unit",
					Value: v["_8"],
				},
			}
			if v["_11"] != "" {
				intType = append(intType, &models.TemplateRequest{
					Field: "KodeGudang",
					Value: v["_11"],
				})
			}
			if v["_12"] != "" {
				intType = append(intType, &models.TemplateRequest{
					Field: "IDSalesman",
					Value: v["_12"],
				})
			}
			intTypeResult, intTypeError := uploadIntTypeValidation(intType)
			if len(intTypeError) > 1 {
				errors = append(errors, intTypeError...)
				continue
			}

			if intTypeResult["Qty"] < 1 {
				errors = append(errors, "Quantity harus lebih dari 0")
				continue
			}

			mustActiveField := []*models.MustActiveRequest{
				helper.GenerateMustActive("stores", "KodeTokoDBO", intTypeResult["KodeTokoDBO"], "active"),
				helper.GenerateMustActive("users", "user_id", user_id, "ACTIVE"),
				{
					Table:    "brands",
					ReqField: "IDMerk",
					Clause:   fmt.Sprintf("id = %d AND status_active = %d", intTypeResult["IDMerk"], 1),
					Id:       intTypeResult["IDMerk"],
				},
				{
					Table:    "products",
					ReqField: "KodeProdukDBO",
					Clause:   fmt.Sprintf("id = %d AND isActive = %d", intTypeResult["KodeProdukDBO"], 1),
					Id:       intTypeResult["KodeProdukDBO"],
				},
				{
					Table:    "uoms",
					ReqField: "Unit",
					Clause:   fmt.Sprintf("id = %d AND deleted_at IS NULL", intTypeResult["Unit"]),
					Id:       intTypeResult["Unit"],
				},
			}
			mustActiveError := r.mustActiveValidation(mustActiveField)
			if len(mustActiveError) > 1 {
				errors = append(errors, mustActiveError...)
				continue
			}

			if len(v["_12"]) > 0 {
				brandSalesman := make(chan *models.RequestIdValidationChan)
				go r.requestValidationRepository.BrandSalesmanValidation(intTypeResult["IDMerk"], intTypeResult["IDSalesman"], idDistributor, brandSalesman)
				brandSalesmanResult := <-brandSalesman

				if brandSalesmanResult.Total < 1 {
					errors = append(errors, fmt.Sprintf("Kode Merek = %d Tidak Terdaftar pada Distributor <nama_agent>. Silahkan gunakan Kode Merek yang lain.", intTypeResult["IDMerk"]))
					errors = append(errors, fmt.Sprintf("ID Salesman = %d Tidak Terdaftar pada Distributor <nama_agent>. Silahkan gunakan ID Salesman yang lain.", intTypeResult["IDSalesman"]))
					errors = append(errors, fmt.Sprintf("Salesman di Kode Toko = %d untuk Merek <Nama Merk> Tidak Terdaftar. Silahkan gunakan ID Salesman yang terdaftar.", intTypeResult["KodeTokoDBO"]))
					continue
				}
			}
			storeAddresses := make(chan *models.RequestIdValidationChan)
			go r.requestValidationRepository.StoreAddressesValidation(intTypeResult["KodeTokoDBO"], storeAddresses)
			storeAddressesResult := <-storeAddresses

			if storeAddressesResult.Total < 1 {
				errors = append(errors, fmt.Sprintf("Alamat Utama pada Kode Toko = %s Tidak Ditemukan. Silahkan gunakan Alamat Toko yang lain.", v["_4"]))
				continue
			}

			var uploadSOSJField models.UploadSOSJField
			uploadSOSJField.TglSuratJalan, err = helper.ParseDDYYMMtoYYYYMMDD(v["_3"])
			if err != nil {
				errors = append(errors, fmt.Sprintf("Format Tanggal Order = %s Salah, silahkan sesuaikan dengan format DD-MMM-YYYY, contoh 15/12/2021", v["_3"]))
				continue
			}
			uploadSOSJField.UploadSOSJFieldMap(v, idDistributor)

			checkIfNoSuratJalanExist := helper.InSliceString(noSuratJalan, v["_2"])
			if checkIfNoSuratJalanExist {

				for i := range uploadSOSJFields {
					brandId, _ := strconv.Atoi(v["_5"])
					if uploadSOSJFields[i].NoSuratJalan == v["_2"] && uploadSOSJFields[i].IDMerk != brandId {
						uploadSOSJFields[i].NoSuratJalan = uploadSOSJFields[i].NoSuratJalan + "-" + strconv.Itoa(uploadSOSJFields[i].IDMerk)
						uploadSOSJField.NoSuratJalan = v["_2"] + "-" + v["_5"]
						break
					} else {
						uploadSOSJField.NoSuratJalan = v["_2"]
					}
				}

			} else {
				uploadSOSJField.NoSuratJalan = v["_2"]
				noSuratJalan = append(noSuratJalan, v["_2"])
			}

			uploadSOSJFields = append(uploadSOSJFields, &uploadSOSJField)
		}
	}

	response.Total = int64(len(uploadSOSJFields))
	response.UploadSOSJFields = uploadSOSJFields
	resultChan <- response
	return

}

func (r *upload) UploadDO(bucket, object, region string, resultChan chan *models.UploadDOFieldsChan) {
	response := &models.UploadDOFieldsChan{}

	results, err := r.ReadFile(bucket, object, region, s3.FileHeaderInfoUse)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	var uploadDOFields []*models.UploadDOField
	for _, v := range results {
		// a, _ := json.Marshal(v)
		// fmt.Println(string(a))
		var uploadDOField models.UploadDOField
		uploadDOField.UploadDOFieldMap(v)

		uploadDOFields = append(uploadDOFields, &uploadDOField)
	}

	response.Total = int64(len(uploadDOFields))
	response.UploadDOFields = uploadDOFields
	resultChan <- response
	return
}

func (r *upload) UploadSO(bucket, object, region string, resultChan chan *models.UploadSOFieldsChan) {
	response := &models.UploadSOFieldsChan{}

	var errors []string

	results, err := r.ReadFile(bucket, object, region, s3.FileHeaderInfoUse)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		response.Error = err
		response.ErrorLog = errorLogData
		resultChan <- response
		return
	}

	brandIds := map[string][]map[string]string{}
	var uploadSOFields []*models.UploadSOField
	for _, v := range results {

		if brandIds[v["NoOrder"]] != nil {
			var error string

			if brandIds[v["NoOrder"]][0]["KodeMerk"] != v["KodeMerk"] {
				fmt.Println("Error satu file")
				response.Total = 0
				response.UploadSOFields = nil
				resultChan <- response
				return
			}

			for _, x := range brandIds[v["NoOrder"]] {
				if x["KodeProduk"] == v["KodeProduk"] && x["UnitProduk"] == v["UnitProduk"] {
					error = fmt.Sprintf("Duplikat row untuk No Order %s", v["NoOrder"])
					break
				}
			}

			if len(error) > 0 {
				errors = append(errors, error)
				continue
			}
		}

		brandIds[v["NoOrder"]] = append(brandIds[v["NoOrder"]], map[string]string{
			"KodeMerk":   v["KodeMerk"],
			"KodeProduk": v["KodeProduk"],
			"UnitProduk": v["UnitProduk"],
		})

		var uploadSOField models.UploadSOField
		uploadSOField.TanggalOrder, err = helper.ParseDDYYMMtoYYYYMMDD(v["TanggalOrder"])
		uploadSOField.TanggalTokoOrder, err = helper.ParseDDYYMMtoYYYYMMDD(v["TanggalTokoOrder"])
		uploadSOField.UploadSOFieldMap(v)

		uploadSOFields = append(uploadSOFields, &uploadSOField)
	}
	a, _ := json.Marshal(errors)
	fmt.Println("Errornya", string(a))

	response.Total = int64(len(uploadSOFields))
	response.UploadSOFields = uploadSOFields
	resultChan <- response
	return

}

func (r *upload) ReadFile(bucket, object, region, fileHeaderInfo string) ([]map[string]string, error) {
	sess := session.New(&aws.Config{
		Region:                        aws.String(region),
		CredentialsChainVerboseErrors: aws.Bool(true),
		Credentials:                   credentials.NewStaticCredentials("ASIAUHX63DBTJAQN5E7R", "gMlQhWzTQnVaWxg9b59FHlhbzZhxE/JMRtfBVI5l", "IQoJb3JpZ2luX2VjELL//////////wEaDmFwLXNvdXRoZWFzdC0xIkcwRQIhAIE4ECz8SUDtYnpt0AD8XyIYXSla6lqqgpWfcdyTMRTeAiAZ6+kl7lPmOZtESowR97YmDvTM5IysqCsgkRjjqUHc6yr1AghrEAAaDDI5MTUxODQ4NjYzMCIM7srqQ27ST0oAUGJCKtICVvZEFnan3ZI/Y/WSUhrFYWufHtMyEL1XC4u8ixlB9IsQUeY4QlzXgJHVETuv5u7QwVv6NvS/+I6IBBQHgRtIjQNIgJ5t1gqGRCQPXAITVUlvMbsd6n8xXF9DWcsz/y/rOhAe5j+nWwFU04esJgG/XxKjGjp2tRXtk64UY1o2woBYKE+iBL3cHju+CqoyKy7TsGLViuHzfbkCM3H6pFCK3xK6cYOn15svREsuTVhgfyGnwP7tZaalkwDCoo17+K8mw+cDKL0iHQq/XIZlG2j6vEe6rj7hKgTSpzzVZ3ILaXR+8Z1cLprM93x3IYevzJvGrvhLck/uGrGUKft4pDouwHr6uOuibpopKxQjauYNlrOtMhCbODv0TnXhv9foRIj3gDiIJ8ZHiSPCH0V8xT1PvADwp/vQYwYpzNSe7jMQqQ48zFp20zKS91qQQg+Qs6TKBU8wiqCqoAY6pwEADQBgh6s77Wg3ri9s7c6T1UWlwCMxYH3bm0Je7e45PZ3VRlburfCsQ7zx0dbHz9FoAQmN+AVNU8B2n4ghBrgFElhdEiI5rcdm7+eoBVPXF6C5uehlvS8p3YmAGaCjw9U6vDoXc07KrJZATpnLvQhw1zljU0ZJSyGUQjMBhwKE0XkdqWSGLvIyP23j4xX/2qYdNZU+OrURMgcNGwVakWs7Fnsma8gVBQ=="),
		HTTPClient:                    &http.Client{Timeout: 10 * time.Second},
	})
	svc := s3.New(sess)

	params := &s3.SelectObjectContentInput{
		Bucket:         aws.String(bucket),
		Key:            aws.String(object),
		ExpressionType: aws.String(s3.ExpressionTypeSql),
		Expression:     aws.String("SELECT * FROM S3Object"),
		InputSerialization: &s3.InputSerialization{
			CSV: &s3.CSVInput{
				FileHeaderInfo:  aws.String(fileHeaderInfo),
				RecordDelimiter: aws.String("\n"),
			},
		},
		OutputSerialization: &s3.OutputSerialization{
			JSON: &s3.JSONOutput{},
		},
	}

	resp, err := svc.SelectObjectContent(params)
	if err != nil {
		fmt.Println("error", err)
		return nil, err
	}
	defer resp.EventStream.Close()

	results, resultWriter := io.Pipe()

	go func() {
		defer resultWriter.Close()
		for event := range resp.EventStream.Events() {
			switch e := event.(type) {
			case *s3.RecordsEvent:
				resultWriter.Write(e.Payload)
				// fmt.Printf("Payload: %v\n", string(e.Payload))
			case *s3.StatsEvent:
				// fmt.Printf("Processed %d bytes\n", *e.Details.BytesProcessed)
			}
		}
	}()

	var records []map[string]string
	resReader := json.NewDecoder(results)
	for {
		var record map[string]string
		err := resReader.Decode(&record)

		if err == io.EOF {
			break
		}
		records = append(records, record)
	}
	if err := resp.EventStream.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "reading from event stream failed, %v\n", err)
		return nil, err
	}

	return records, nil
}

func uploadMandatoryValidation(request []*models.TemplateRequest) []string {
	errors := []string{}

	for _, value := range request {
		if len(value.Value) < 1 {
			error := fmt.Sprintf("Data %s tidak boleh kosong", value.Field)
			errors = append(errors, error)
		}
	}

	return errors
}

func uploadIntTypeValidation(request []*models.TemplateRequest) (map[string]int, []string) {
	result := map[string]int{}
	errors := []string{}

	for _, v := range request {
		parseInt, error := strconv.Atoi(v.Value)

		if error != nil {
			error := fmt.Sprintf("Data %s harus bertipe data integer", v.Value)
			errors = append(errors, error)
		} else {
			result[v.Field] = parseInt
		}
	}

	return result, errors
}

func (r *upload) mustActiveValidation(request []*models.MustActiveRequest) []string {

	errors := []string{}

	mustActive := make(chan *models.MustActiveRequestChan)
	go r.requestValidationRepository.MustActiveValidation(request, mustActive)
	mustActiveResult := <-mustActive

	for k, v := range mustActiveResult.Total {
		if v < 1 {

			var error string
			switch request[k].ReqField {
			case "KodeTokoDBO":
				error = fmt.Sprintf("Kode Toko = %d sudah Tidak Aktif. Silahkan gunakan Kode Toko yang lain", request[k].Id)
			case "KodeProdukDBO":
				error = fmt.Sprintf("Kode SKU = %d dengan Merek <Nama_Merk> sudah Tidak Aktif. Silahkan gunakan Kode SKU yang lain.", request[k].Id)
			case "IDMerk":
				error = fmt.Sprintf("Merk = %d sudah Tidak Aktif. Silahkan gunakan Kode yang lain.", request[k].Id)
			default:
				error = fmt.Sprintf("Kode = %d sudah Tidak Aktif. Silahkan gunakan Kode yang lain.", request[k].Id)
			}
			errors = append(errors, error)

		}
	}

	return errors
}
