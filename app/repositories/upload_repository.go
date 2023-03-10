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
	UploadDO(bucket, object, region string, resultChan chan *models.UploadDOFieldsChan)
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

func (r *upload) ReadFile(bucket, object, region, fileHeaderInfo string) ([]map[string]string, error) {
	sess := session.New(&aws.Config{
		Region:                        aws.String(region),
		CredentialsChainVerboseErrors: aws.Bool(true),
		Credentials:                   credentials.NewStaticCredentials("ASIAUHX63DBTMP5NRTFF", "KU2edn3ocS0YzFojWFl73S+9TlXAYwE1GqrKVmAG", "IQoJb3JpZ2luX2VjEJ3//////////wEaDmFwLXNvdXRoZWFzdC0xIkcwRQIhAP5c2AqDnMDdUMkjHHui/fqbZ7TyvnaSrEZCeZpUachPAiAeWBT4uuvvUOosFTlio+S8T/6wZZryFQNMXfzchZZBBSr3AghWEAAaDDI5MTUxODQ4NjYzMCIMJHlP5ptN5g9AkgJEKtQCNISePFrLtTof+txAe9VxibuYOQBGXTPPk8VKdgC1QZJ1lmtoU+5bX5gR7Kc9YtVfZHvr9/QVaIUj9I0UtxPuZhWtKndvJ5Oa73NTJ9opp+Pm8JtnBB7fx6lA4NS4frvCYZnP8Vjx9BcN5+VNhaVBv2HZSQTViS3pClNWiGc3iY0Pw2A5PSYC2L+++jUeZv5MLLkcZta+d7OPC0uQ18aAjnkFa+l2wa4j4G/+EnFYbbbKSKprtg+nsttYqxoMmfDS/2P2rsNyLLIDc2XSqLjbNEGbenXXzFEAfOkTruLTvMV4KfZGGapszfZJrUMSVEBZkUmPvMejDDtAzxfpMfpZ0swPzA0Tfr5G2fu36EGjziaO1RpSoAORqPfEmLGhod/uL+1WmJEhf4XayHwTCcdez6rpp2T9BUgHP2Lxa9Yr80ZK8aEAM7GAKSLmapHCyv/0xTVGyDCgwqWgBjqnAZqtxL7g3kQ3FAY+c8P7HD1AT23DVoqc9jUSyWoR3ko1BQ7dxdxCwJxRCuUrFgg0wNNNsOPSQxTspnSAR0luI+mbP3YNLEDfQIFhqTZ92dubzQm5A9m44Yeecug2CoxA/MaKMfgRJP34V8/rKI/+XVKk4wmrrARfrjG/5sCAcSHkXtaiXD2PfgER4ZNKw8bIG+3z90wPnldBA7XWHGg2zeKIIMxuKKGV"),
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
