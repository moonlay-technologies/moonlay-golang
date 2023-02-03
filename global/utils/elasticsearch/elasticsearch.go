package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"order-service/global/utils/model"
)

type elasticSearchClient struct {
	connection *elasticsearch8.Client
	ctx        context.Context
}

type ElasticSearchClientInterface interface {
	GetConnection() *elasticsearch8.Client
	InsertDocument(index string, body []byte) (string, error)
	Query(index string, query map[string]interface{}) (*model.ElasticSearchQueryResponse, error)
}

func InitElasticSearchClientInterface(context context.Context, hosts []string) ElasticSearchClientInterface {
	cfg := elasticsearch8.Config{
		Addresses: hosts,
	}

	es, err := elasticsearch8.NewClient(cfg)

	if err != nil {
		errStr := fmt.Sprintf("Error failed connect to elasticsearch")
		fmt.Println(errStr)
		fmt.Println(err)
		panic(err)
	}

	esClient := &elasticSearchClient{
		connection: es,
		ctx:        context,
	}

	return esClient
}

func (es *elasticSearchClient) GetConnection() *elasticsearch8.Client {
	return es.connection
}

func (es *elasticSearchClient) InsertDocument(index string, body []byte) (string, error) {
	request := esapi.IndexRequest{
		Index:   index,
		Body:    bytes.NewReader(body),
		Refresh: "true",
	}

	response, err := request.Do(es.ctx, es.connection)

	if err != nil {
		errStr := fmt.Sprintf("Error failed request document to elasticsearch")
		fmt.Println(errStr)
		fmt.Println(err)
		return "", err
	}

	defer response.Body.Close()

	if response.IsError() {
		errStr := fmt.Sprintf("Error failed insert document to elasticsearch")
		fmt.Println(errStr)
		fmt.Println(response.Status())
		return "", err
	}

	elasticInsertResponse := model.ElasticSearchInsertResponse{}
	err = json.NewDecoder(response.Body).Decode(&elasticInsertResponse)

	if err != nil {
		errStr := fmt.Sprintf("Error decode insert document result")
		fmt.Println(errStr)
		fmt.Println(err)
		return "", err
	}

	return elasticInsertResponse.ID, nil
}

func (es *elasticSearchClient) Query(index string, query map[string]interface{}) (*model.ElasticSearchQueryResponse, error) {
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(query)

	if err != nil {
		errStr := fmt.Sprintf("Error decode query search")
		fmt.Println(errStr)
		fmt.Println(err)
		return &model.ElasticSearchQueryResponse{}, err
	}

	response, err := es.connection.Search(
		es.connection.Search.WithContext(es.ctx),
		es.connection.Search.WithIndex(index),
		es.connection.Search.WithBody(&buf),
		es.connection.Search.WithTrackTotalHits(true),
		es.connection.Search.WithPretty(),
	)

	if err != nil {
		errStr := fmt.Sprintf("Error get query search to elasticsearch")
		fmt.Println(errStr)
		fmt.Println(err)
		return &model.ElasticSearchQueryResponse{}, err
	}

	defer response.Body.Close()

	if response.IsError() {
		errStr := fmt.Sprintf("Error failed get query search to elasticsearch")
		fmt.Println(errStr)
		fmt.Println(response.Status())
		return &model.ElasticSearchQueryResponse{}, err
	}

	elasticSearchQueryResponse := model.ElasticSearchQueryResponse{}
	err = json.NewDecoder(response.Body).Decode(&elasticSearchQueryResponse)

	if err != nil {
		errStr := fmt.Sprintf("Error decode search query result")
		fmt.Println(errStr)
		fmt.Println(err)
		return &model.ElasticSearchQueryResponse{}, err
	}

	return &elasticSearchQueryResponse, nil
}
