package handlers

import (
	"context"
	"fmt"
	"order-service/app/commands"
	kafkadbo "order-service/global/utils/kafka"
	"os"
	"strconv"
)

func MainCommandHandler(kafkaClient kafkadbo.KafkaClientInterface, ctx context.Context, args []interface{}) {
	switch args[1] {
	case "create-kafka-topic":
		createKafkaTopicCommand := commands.InitCreateKafkaTopicCommand(kafkaClient, ctx)
		totalPartitionStr := args[3].(string)
		totalPartitionInt, _ := strconv.Atoi(totalPartitionStr)
		err := createKafkaTopicCommand.CreateTopic(args[2].(string), totalPartitionInt)
		if err != nil {
			panic(err)
		}
		break

	default:
		fmt.Println("Choose Command Type You Want")
		os.Exit(0)
	}
}
