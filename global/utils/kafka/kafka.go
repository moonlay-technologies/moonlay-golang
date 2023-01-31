package kafkadbo

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"
)

type kafkaClient struct {
	connection          *kafka.Conn
	brokers             []string
	writer              *kafka.Writer
	reader              *kafka.Reader
	consumerGroupReader *kafka.Reader
	controller          *kafka.Conn
	ctx                 context.Context
}

type KafkaClientInterface interface {
	GetConnection() *kafka.Conn
	GetController() *kafka.Conn
	CreateTopic(topic string, totalPartition int, totalReplicationFactor int) error
	WriteToTopic(topic string, key []byte, message []byte) error
	SetWriter(topic string)
	SetReader(topic string, partition int, offset int64) *kafka.Reader
	SetConsumerGroupReader(topic string, groupID string) *kafka.Reader
	GetBrokers() []string
}

func InitKafkaClientInterface(context context.Context, brokers []string) KafkaClientInterface {
	client := &kafkaClient{
		brokers: brokers,
		ctx:     context,
	}

	return client
}

func (k *kafkaClient) GetBrokers() []string {
	return k.brokers
}

func (k *kafkaClient) GetConnection() *kafka.Conn {
	return k.connection
}

func (k *kafkaClient) GetController() *kafka.Conn {
	return k.controller
}

func (k *kafkaClient) CreateTopic(topic string, totalPartition int, totalReplicationFactor int) error {
	conn, err := kafka.Dial("tcp", k.brokers[0])

	if err != nil {
		errStr := fmt.Sprintf("Error failed connect to kafka")
		fmt.Println(errStr)
		fmt.Println(err)
		panic(err)
	}

	controller, err := conn.Controller()

	if err != nil {
		errStr := fmt.Sprintf("Error failed get to controller")
		fmt.Println(errStr)
		fmt.Println(err)
		panic(err)
	}

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))

	if err != nil {
		errStr := fmt.Sprintf("Error failed connect to controller")
		fmt.Println(errStr)
		fmt.Println(err)
		panic(err)
	}

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     totalPartition,
			ReplicationFactor: totalReplicationFactor,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)

	if err != nil {
		if errors.Is(err, syscall.EPIPE) {
			fmt.Println("broken pipe")
		} else {
			errStr := fmt.Sprintf("Error failed create topic %s", topic)
			fmt.Println(errStr)
			fmt.Println(err)
			return err
		}
	}

	k.controller = controllerConn
	k.connection = conn

	return nil
}

func (k *kafkaClient) SetWriter(topic string) {
	var writer *kafka.Writer

	if os.Getenv("ENVIRONMENT") == "local" {
		writer = &kafka.Writer{
			Addr:         kafka.TCP(k.brokers[0], k.brokers[1], k.brokers[2]),
			Topic:        topic,
			Balancer:     kafka.CRC32Balancer{},
			BatchTimeout: 5 * time.Millisecond,
		}
	} else {
		writer = &kafka.Writer{
			Addr:         kafka.TCP(k.brokers[0], k.brokers[1]),
			Topic:        topic,
			Balancer:     kafka.CRC32Balancer{},
			BatchTimeout: 5 * time.Millisecond,
		}
	}

	k.writer = writer
}

func (k *kafkaClient) WriteToTopic(topic string, key []byte, message []byte) error {
	k.SetWriter(topic)

	err := k.writer.WriteMessages(k.ctx,
		kafka.Message{
			Value: message,
		},
	)

	if err != nil {
		errStr := fmt.Sprintf("Error failed insert to topic %s", topic)
		fmt.Println(errStr)
		fmt.Println(err)
		return err
	}

	defer k.writer.Close()
	return nil
}

func (k *kafkaClient) SetReader(topic string, partition int, offset int64) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   k.brokers,
		Topic:     topic,
		Partition: partition,
		MinBytes:  10e3,
		MaxWait:   10e6,
	})

	if offset > 0 {
		reader.SetOffset(offset)
	}

	k.reader = reader
	return reader
}

func (k *kafkaClient) SetConsumerGroupReader(topic string, groupID string) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: k.brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	k.consumerGroupReader = reader
	return reader
}
