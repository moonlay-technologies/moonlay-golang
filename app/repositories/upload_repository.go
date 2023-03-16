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
		Credentials:                   credentials.NewStaticCredentials("ASIAUHX63DBTCQN57MZS", "5K2+qV/g+mIbChmUHYidRD03k+fa31Uw6p5aeClm", "IQoJb3JpZ2luX2VjEEcaDmFwLXNvdXRoZWFzdC0xIkcwRQIga0lOYBYksz6SXGmezstpyiD67MH3d3xh+lWbh5aok6YCIQDYeX1oKoaj7jr4IxQZ5gXoK8utlTaqnqMc9g45XMVgVCr1AggQEAAaDDI5MTUxODQ4NjYzMCIME/yjuTO07RA9QSYFKtICZ/E0Pvn5RDOiTyNwxJ5/PObGXMv/1BMFK64o4mNijFlHusbBx6rtv9tbG07C5cxnF3J19wVtN0B3c1eTD3S7O5W74z6tdql8bFmUWfVKQi0Va2tLX1ffEKHu1/PwKTLj5te/fMEq7cT3eICGvNlT5LKVrWYsaKlgayQiRIp3rRQbksfmT6Z0ZUEgX2m59ssFbJNUTqDW4qZAl0zVX3dfiG0aGgbIjBE8DRGRgCBh2WLIGGHkqKLBGEOsnouf3gJqz0NfKztUy5FbU1fSICvctkE6/rFg8y1kl0UmBZnFDpCGWfPaFeeCm5eRLronvvOPFcvuGs15JA0bOh8XNY7Ed0WnnDOGkFPPJzt5JH+AKIFevnEmxoNUI/eR6Dk5FSWrEtfzYyv4XQ+O0XlD7IPRpxfkyX74txuZIWVIp/V1xeUAhihOIRrjiz5x/Gl9Fk7tZAIwyPrKoAY6pwE6Z/715CzS3/pWST3fWF7CF8d9pKUeLcaIxnSOTRRFD0lOPcQWPPsBTtYLJjr3e4/UOEqa3uPv1mdG4KaNsfyaMoebaPh48P3mafyZMXieYX+rWBUqN9LAIdXsBS+1TYgYMfzDcg4dJauRdKJx4gVVEtgRrmdreJ7EcxOT9xfglafU+hsESQ5y3/fOF57ss1dS2PwKdhSVhCyxZ3ryK9g0AM2Egsm55g=="),
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
