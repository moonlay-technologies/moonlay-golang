package opensearch_dbo

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/global/utils/helper"
	"order-service/global/utils/model"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

type openSearchClient struct {
	addresses     []string
	username      string
	password      string
	maxRetries    int
	retryOnStatus []int
	ctx           context.Context
	client        *opensearch.Client
}

type OpenSearchClientInterface interface {
	CreateDocument(index string, documentID string, document []byte) (string, error)
	GetByID(index string, documentID string) (*model.OpenSearchGetResponse, error)
	Query(index string, query []byte) (*model.OpenSearchQueryResponse, error)
	Count(index string, query []byte) (int64, error)
	GetConnection() *opensearch.Client
}

func InitOpenSearchClientInterface(addresses []string, username string, password string, ctx context.Context) OpenSearchClientInterface {
	osClient := &openSearchClient{
		addresses:     addresses,
		username:      username,
		password:      password,
		maxRetries:    3,
		retryOnStatus: []int{502, 503, 504},
		ctx:           ctx,
	}

	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses:     osClient.addresses,
		Username:      osClient.username,
		Password:      osClient.password,
		MaxRetries:    osClient.maxRetries,
		RetryOnStatus: osClient.retryOnStatus,
	})

	if err != nil {
		fmt.Println("Error Connect to OpenSearch : ", err.Error())
		panic(err)
	}

	osClient.client = client
	return osClient
}

func (os *openSearchClient) CreateDocument(index string, documentID string, document []byte) (string, error) {
	req := opensearchapi.IndexRequest{
		Index: index,
		Body:  bytes.NewReader(document),
	}

	if len(documentID) > 0 {
		req.DocumentID = documentID
	}

	response, err := req.Do(os.ctx, os.client)

	if err != nil {
		fmt.Println("Error Insert Index Open Search : ", err.Error())
		return "", err
	}

	defer response.Body.Close()

	if response.IsError() {
		errStr := fmt.Sprintf("Error failed Insert Document to OpenSearch")
		fmt.Println(errStr)
		var errorRes interface{}
		err = json.NewDecoder(response.Body).Decode(&errorRes)
		fmt.Println(errorRes)
		fmt.Println(response.Status())
		return "", err
	}

	insertResponse := model.OpenSearchInsertResponse{}
	err = json.NewDecoder(response.Body).Decode(&insertResponse)

	if err != nil {
		errStr := fmt.Sprintf("Error decode insert document result")
		fmt.Println(errStr)
		fmt.Println(err)
		return "", err
	}

	return insertResponse.ID, nil
}

func (os *openSearchClient) GetByID(index string, documentID string) (*model.OpenSearchGetResponse, error) {
	req := opensearchapi.GetRequest{
		Index:      index,
		DocumentID: documentID,
	}

	response, err := req.Do(os.ctx, os.client)

	if err != nil {
		fmt.Println("Error Get Index Open Search : ", err.Error())
		return &model.OpenSearchGetResponse{}, err
	}

	defer response.Body.Close()

	if response.IsError() {
		errStr := fmt.Sprintf("Error failed Get Document to elasticsearch")
		fmt.Println(errStr)
		fmt.Println(response.Status())
		return &model.OpenSearchGetResponse{}, err
	}

	getResponse := model.OpenSearchGetResponse{}
	err = json.NewDecoder(response.Body).Decode(&getResponse)

	if err != nil {
		errStr := fmt.Sprintf("Error decode insert document result")
		fmt.Println(errStr)
		fmt.Println(err)
		return &model.OpenSearchGetResponse{}, err
	}

	return &getResponse, nil
}

func (os *openSearchClient) Query(index string, query []byte) (*model.OpenSearchQueryResponse, error) {
	req := opensearchapi.SearchRequest{
		Index: []string{index},
		Body:  bytes.NewReader(query),
	}

	response, err := req.Do(os.ctx, os.client)

	if err != nil {
		fmt.Println("Error Query Index Open Search : ", err.Error())
		return &model.OpenSearchQueryResponse{}, err
	}

	if response.IsError() == true {
		errStr := fmt.Sprintf("Error failed Query Document to open_search")
		err = helper.NewError(errStr)
		fmt.Println(errStr)
		fmt.Println(response.Status())
		return &model.OpenSearchQueryResponse{}, err
	}

	searchResponse := &model.OpenSearchQueryResponse{}

	err = json.NewDecoder(response.Body).Decode(searchResponse)

	if err != nil {
		errStr := fmt.Sprintf("Error decode Query document result")
		fmt.Println(errStr)
		fmt.Println(err)
		return &model.OpenSearchQueryResponse{}, err
	}

	return searchResponse, nil
}

func (os *openSearchClient) Count(index string, query []byte) (int64, error) {
	req := opensearchapi.CountRequest{
		Index: []string{index},
		Body:  bytes.NewReader(query),
	}

	response, err := req.Do(os.ctx, os.client)

	if err != nil {
		fmt.Println("Error Count Index Open Search : ", err.Error())
		return 0, err
	}

	if response.IsError() == true {
		errStr := fmt.Sprintf("Error failed Count Document to open_search")
		err = helper.NewError(errStr)
		fmt.Println(errStr)
		fmt.Println(response.Status())
		return 0, err
	}

	searchResponse := &model.OpenSearchCountResponse{}

	err = json.NewDecoder(response.Body).Decode(searchResponse)

	if err != nil {
		errStr := fmt.Sprintf("Error decode Query document result")
		fmt.Println(errStr)
		fmt.Println(err)
		return 0, err
	}

	return searchResponse.Count, nil
}

func (os *openSearchClient) GetConnection() *opensearch.Client {
	return os.client
}
