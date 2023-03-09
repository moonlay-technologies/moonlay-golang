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
	UploadSOSJ(bucket, object, region string, resultChan chan *models.UploadSOSJFieldChan)
}

type upload struct {
}

func InitUploadRepository() UploadRepositoryInterface {
	return &upload{}
}

func (r *upload) UploadSOSJ(bucket, object, region string, resultChan chan *models.UploadSOSJFieldChan) {
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
		intType := []*models.TemplateRequest{
			{
				Field: "KodeTokoDBO",
				Value: v["_4"],
			},
			{
				Field: "IDMerek",
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
			{
				Field: "KodeGudang",
				Value: v["_11"],
			},
		}
		if v["_12"] != "" {
			intType = append(intType, &models.TemplateRequest{
				Field: "IDSalesman",
				Value: v["_12"],
			})
		}
		intTypeError := uploadIntTypeValidation(intType)
		if len(intTypeError) > 1 {
			errors = append(errors, intTypeError...)
			continue
		}

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

		// mustActiveField := []*models.MustActiveRequest{
		// 	helper.GenerateMustActive("stores", "store_id", v["_4"], "active"),
		// 	helper.GenerateMustActive("users", "user_id", insertRequest.UserID, "ACTIVE"),
		// }

		if v["_1"] != "Status" {

			var uploadSOSJField models.UploadSOSJField
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
			uploadSOSJField.UploadSOSJFieldMap(v, idDistributor)
			uploadSOSJField.TglSuratJalan, _ = helper.ParseDDYYMMtoYYYYMMDD(uploadSOSJField.TglSuratJalan)
			uploadSOSJFields = append(uploadSOSJFields, &uploadSOSJField)
		}
	}

	response.Total = int64(len(uploadSOSJFields))
	response.UploadSOSJFields = uploadSOSJFields
	resultChan <- response
	return

}

func (r *upload) ReadFile(bucket, object, region, fileHeaderInfo string) ([]map[string]string, error) {
	sess := session.New(&aws.Config{
		Region:                        aws.String(region),
		CredentialsChainVerboseErrors: aws.Bool(true),
		Credentials:                   credentials.NewStaticCredentials("ASIAUHX63DBTJSNFFWFY", "OWoH8MJin2iBegS32dT8HDGx0ilgY78ih8y8kLSM", "IQoJb3JpZ2luX2VjEJr//////////wEaDmFwLXNvdXRoZWFzdC0xIkcwRQIhALE9XsKawi1OxiCjJ88z4OrbQdMeAJKXm8DO2rxjmQEpAiBAf+9+oMN1q/gfe/u8nFebukjhqMjhRPCESkoBtVTYJyr1AghTEAAaDDI5MTUxODQ4NjYzMCIMQb3/yCle1EhmcFtJKtICgdc0v+q9l2yWvuPzSzl4zUIqTb/r6QSkT+I5XnSJNoCvwkCQVphLhWB6CwaoQKoztHdbAt+rGQuCPMkhnnXV+rRD8nte/hBr83p2qXxKyeCotVZlA6+6QloGwbNp23+A/cD43/0qHgrQnXQXdVlBfT4f9k9TJZ92cRNgesb6+5JWDOb+epPMbomuBn9HuKlnjzOHEqwXODDTbyooU0l0HX+gnf4+uWqs3z0MlA1U1dInsWAzx78HWs0ggdMkrZvISXi1I7B8k5tNRfYV+PFGQq7TAuBCZ1bm2kxbeUU1WzVf1S77alr1hoI+gWgpRSql+1O2nYMHHdVppPlv1WaIHli1uxsSqKSK3rmC/bFO8zXOpUvzEFjQEfYql4pozrRsWo97Rt3XhrzmBQ9sxeAOo60ImWrvPuOWhK4YHUuLkwwZTqk5X9T3i7W5W5Ef4zq2VbQwhfakoAY6pwGmXuCOpI76IRaCiP17X7aItZiE8H4dV15CIk28bHHGXzGSKNf20Zn2elA2kLy/1xfZP0+jq99N1keK/qwXevEBiDY/LaON7Ng51XKVJok856s7Sygr5wkcr6+RzwORGMl516hMS3PZTN2yfZMIq+m4nfJMI639CyTwfGMCQ5lvI4ok3rwUayhqjUfkccmO3+LSbjNt5qE/lrUWxwYbv7KI2QEdVaTinw=="),
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

func uploadIntTypeValidation(request []*models.TemplateRequest) []string {
	errors := []string{}

	for _, v := range request {
		_, error := strconv.Atoi(v.Value)
		if error != nil {
			error := fmt.Sprintf("Data %s harus bertipe data integer", v.Value)
			errors = append(errors, error)
		}
	}

	return errors
}
