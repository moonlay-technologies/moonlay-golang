package commands

import (
	"context"
	"fmt"
	"net/http"
	"order-service/global/utils/helper"
	kafkadbo "order-service/global/utils/kafka"
)

type CreateKafkaTopicCommandHandlerInterface interface {
	CreateTopic(topic string, partition int) error
}

type createKafkaTopicCommandHandler struct {
	kafkaClient kafkadbo.KafkaClientInterface
	ctx         context.Context
}

func InitCreateKafkaTopicCommandHandler(kafkaClient kafkadbo.KafkaClientInterface, ctx context.Context) CreateKafkaTopicCommandHandlerInterface {
	return &createKafkaTopicCommandHandler{
		kafkaClient: kafkaClient,
		ctx:         ctx,
	}
}

func (c *createKafkaTopicCommandHandler) CreateTopic(topic string, partition int) error {
	brokers := c.kafkaClient.GetBrokers()
	totalBroker := len(brokers)
	replicationFactor := totalBroker - 1
	err := c.kafkaClient.CreateTopic(topic, partition, replicationFactor)

	if err != nil {
		errorLogData := helper.WriteLog(err, http.StatusInternalServerError, nil)
		fmt.Println(errorLogData)
		return err
	}

	return nil
}
