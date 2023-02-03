package commands

import (
	"context"
	kafkadbo "order-service/global/utils/kafka"
)

func InitCreateKafkaTopicCommand(kafkaClient kafkadbo.KafkaClientInterface, ctx context.Context) CreateKafkaTopicCommandHandlerInterface {
	handler := InitCreateKafkaTopicCommandHandler(kafkaClient, ctx)
	return handler
}
