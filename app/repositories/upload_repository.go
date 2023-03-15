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
	ReadFile(bucket, object, region, fileHeaderInfo string) ([]map[string]string, error)
	UploadDO(bucket, object, region string, resultChan chan *models.UploadDOFieldsChan)
}

type upload struct {
	requestValidationRepository RequestValidationRepositoryInterface
}

func InitUploadRepository(requestValidationRepository RequestValidationRepositoryInterface) UploadRepositoryInterface {
	return &upload{
		requestValidationRepository: requestValidationRepository,
	}
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
		Credentials:                   credentials.NewStaticCredentials("ASIAUHX63DBTNIVERB5G", "I6sZyfIFZIrUndQWUSchCmQ6FPT5RWN4lGrMxMXR", "IQoJb3JpZ2luX2VjEDYaDmFwLXNvdXRoZWFzdC0xIkgwRgIhAK5S809Aka25Kw2LPZRYhNX3c3npXLiV9h6pLmKcuHCIAiEAoYisa4lW9r/vGITaslY744XxfiHRyUXcaOucyXdXQXwq/gII7///////////ARAAGgwyOTE1MTg0ODY2MzAiDHj2t9wxBqBt6c/21irSAhDhUh+GgS3FEIEVM3exRVxwi/WlBpenXQx1eJBo8L62ZpI5tW/tT7tpeOMqqrZ8j5Sxhzxp4MPjtl7jCLFflNDb1WaFUqwcCTT9n32TDPXUdLTdhIMjrrNnwdAkf0UPTBH+xzl3bhn7dXBVDOR2Dcs37K0ZlydDQ+gBtDrGmD7CdVUt/o21hfkUmeSfs9Q0LIYhfGoGNte8M8sgAgBV0ZU7vZSrHlHfgoWHbfTEzYpntpl3jHcUE7/bdYNW3YdLoC8ZW8/wgDBKqVt6MkilhkyhcVTzwWQJPBJz+daHukoCw6aHaFbhzQoLs5ED1eyAxpWJSD9JGjd2WLNBRYOZ/tzcWmNk8vKbEUXGlC+I2k7SwRDFDFHPmNSyNyi7kbl0Yn5+raLNCFKEAbirJUXM+/zmhbMXhbYsj2DX2S0iDG9giv05Um6KTXAMJ59zu3QwGAEKMPuZx6AGOqYBOdK5d4hnjguPsA/F4MR1uorY5s0Gk1e3yP1D8TycGpcDIfeeowU9YrPrNMNOucjIBJIDnAYNcdqDrjfrdzPEsIaTAv7TpJWvUcrnGg+Ll9jczfm1g6vDftKWGzzBL8PSdAyVLU5Hn1/xloFL3JpRv4aUZZvV0s+amR7tPR6BFz9HjQrVCT+WK4yj2uMRvfztkLXrqG5VR0aDBKDHCwR/iflgwIOzcA=="),
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
