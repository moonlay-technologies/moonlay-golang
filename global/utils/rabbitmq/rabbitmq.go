package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	maps "github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

const (
	DefaultExchange string = ""
	delay                  = 300 * time.Millisecond // reconnect after delay seconds
)

type MessageBody struct {
	Info    *MessageInfo `json:"info"`
	Content interface{}  `json:"content"`
}

type MessageInfo struct {
	RequestID string `json:"request_id"`
	Origin    string `json:"origin"`
	Source    string `json:"source"`
}

//Result data structure
type Result struct {
	Data  interface{}
	Error error
}

//RabbitMQ :
type RabbitMQ interface {
	Connection() *Connection
	ChannelPool() ChannelPool
	ServiceCode() string
	NewChannel() (*Channel, error)
	GetChannel() (*Channel, error)
	ChannelDone(target *Channel)
	Disconnect() error
	SetLogLevel(logLevel string)
	RegisterExchange(exchangeName string, exchangeType string) error
	InitQueue(queueName string, exchangeName string, routingKey string, durable, autoDelete bool) (amqp.Queue, error)

	EncapsulateData(info *MessageInfo, content interface{}) MessageBody
	DecodeMapType(input interface{}, output interface{}) error
	ExtractMessageData(msg amqp.Delivery, out interface{}) (*MessageInfo, error)

	ReadMessage(q amqp.Queue) (<-chan amqp.Delivery, *Channel, error)
	PublishMessage(queueName string, durable bool, autoDelete bool, exchange string, routingKey string, contentType string, headers map[string]interface{}, body interface{}, msgInfo *MessageInfo) error
	PublishMessageSync(exchange string, routingKey string, contentType string, headers map[string]interface{}, body interface{}, msgInfo *MessageInfo) <-chan Result
	PostProcessMessage(d amqp.Delivery, response interface{})
}
type rabbitMQ struct {
	mConnection sync.Mutex
	connection  *Connection
	mChannel    sync.Mutex
	chPool      ChannelPool
	serviceCode string
}

//NewRabbitMQ : Create New RabbitMQ connection (connection and chPool)
func NewRabbitMQ(connString string, serviceCode string) (RabbitMQ, error) {
	conn, err := NewConnection(connString)
	if err != nil {
		return nil, err
	}
	c, e := conn.NewChannelPool(serviceCode)
	if e != nil {
		return nil, err
	}
	return &rabbitMQ{
		connection:  conn,
		chPool:      c,
		serviceCode: serviceCode,
	}, nil
}

func (rmq *rabbitMQ) Connection() *Connection {
	rmq.mConnection.Lock()
	defer rmq.mConnection.Unlock()
	return rmq.connection
}

func (rmq *rabbitMQ) NewChannel() (*Channel, error) {
	rmq.mConnection.Lock()
	defer rmq.mConnection.Unlock()
	nChannel, err := rmq.connection.Channel(rmq.serviceCode)
	if err != nil {
		return nil, err
	}
	return nChannel, nil
}

func (rmq *rabbitMQ) Disconnect() error {
	rmq.mConnection.Lock()
	defer rmq.mConnection.Unlock()
	rmq.chPool.CloseAll()

	err := rmq.connection.Close()
	if err != nil {
		return err
	}

	return nil
}

func (rmq *rabbitMQ) SetLogLevel(logLevel string) {
	switch strings.ToLower(logLevel) {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.ErrorLevel)
	}

}

func (rmq *rabbitMQ) ChannelPool() ChannelPool {
	rmq.mChannel.Lock()
	defer rmq.mChannel.Unlock()
	return rmq.chPool
}

func (rmq *rabbitMQ) GetChannel() (*Channel, error) {
	return rmq.chPool.GetChannel()
}

func (rmq *rabbitMQ) ChannelDone(target *Channel) {
	rmq.chPool.Done(target)
}

func (rmq *rabbitMQ) ServiceCode() string {
	rmq.mConnection.Lock()
	defer rmq.mConnection.Unlock()
	return rmq.serviceCode
}

//DeclareExchange : Declaring exchange
func (rmq *rabbitMQ) RegisterExchange(exchangeName string, exchangeType string) error {
	ch, _ := rmq.GetChannel()
	defer rmq.ChannelDone(ch)
	return ch.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,         // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}

//InitQueue : Declare and binding queue,  queue name will be combination of is q.ServiceCode.routingKey
func (rmq *rabbitMQ) InitQueue(queueName string, exchangeName string, routingKey string, durable, autoDelete bool) (amqp.Queue, error) {
	var err error
	var queue amqp.Queue
	if queue, err = rmq.declareQueue(queueName, durable, autoDelete); err != nil {
		log.Error(err)
	} else {
		if err = rmq.bindQueue(queue, exchangeName, routingKey); err != nil {
			log.Error(err)
		}
	}

	return queue, err
}

// declareQueue :
func (rmq *rabbitMQ) declareQueue(name string, durable, autoDelete bool) (amqp.Queue, error) {
	ch, _ := rmq.GetChannel()
	defer rmq.ChannelDone(ch)
	q, err := ch.QueueDeclare(
		name,       // name
		durable,    // durable
		autoDelete, // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	rmq.failOnError("error when declaring a queue", err)
	return q, err
}

//bindQueue :
func (rmq *rabbitMQ) bindQueue(q amqp.Queue, exchangeName, routingKey string) error {
	ch, _ := rmq.GetChannel()
	defer rmq.ChannelDone(ch)
	err := ch.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	rmq.failOnError("error when binding a queue", err)
	return err
}

func (rmq *rabbitMQ) failOnError(msg string, err error) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

func (rmq *rabbitMQ) EncapsulateData(info *MessageInfo, content interface{}) MessageBody {
	var body MessageBody
	if info == nil {
		info = &MessageInfo{
			RequestID: uuid.New().String(),
			Origin:    rmq.ServiceCode(),
			Source:    rmq.ServiceCode(),
		}
	} else {
		info.Source = rmq.ServiceCode()
	}
	switch content.(type) {
	case []byte:
		body = MessageBody{
			Info:    info,
			Content: string(content.([]byte)),
		}
	default:
		body = MessageBody{
			Info:    info,
			Content: content,
		}
	}
	return body
}

func (rmq *rabbitMQ) DecodeMapType(input interface{}, output interface{}) error {
	if input == nil {
		return fmt.Errorf("error when decoding message : nil data input")
	}
	config := &maps.DecoderConfig{
		Metadata:   nil,
		Result:     output,
		TagName:    "json",
		DecodeHook: maps.ComposeDecodeHookFunc(rmq.toTimeHookFunc()),
	}

	decoder, err := maps.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func (rmq *rabbitMQ) toTimeHookFunc() maps.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
		// Convert it by parsing
	}
}

func (rmq *rabbitMQ) ExtractMessageData(msg amqp.Delivery, out interface{}) (*MessageInfo, error) {
	var resp MessageBody
	if err := json.Unmarshal(msg.Body, &resp); err != nil {
		return nil, fmt.Errorf("error when unmarshaling data %s", err.Error())
	}
	if err := rmq.DecodeMapType(resp.Content, &out); err != nil {
		return nil, fmt.Errorf("error when decoding data %s", err.Error())
	} else {
		return resp.Info, nil
	}
}

//ReadMessage :
func (rmq *rabbitMQ) ReadMessage(q amqp.Queue) (<-chan amqp.Delivery, *Channel, error) {
	ch, _ := rmq.NewChannel()
	//defer func() {
	//	fmt.Println("start defer ReadMessage")
	//	for {
	//		select {
	//		case <-ctx.Done(): // if cancel() execute
	//			ch.Close()
	//			break
	//		}
	//		time.Sleep(100 * time.Millisecond)
	//	}
	//	fmt.Println("defer ReadMessage")
	//}()

	msg, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	rmq.failOnError("error when reading the messages", err)
	return msg, ch, err
}

func (rmq *rabbitMQ) PublishMessage(queueName string, durable bool, autoDelete bool,exchange string, routingKey string, contentType string, headers map[string]interface{}, body interface{}, msgInfo *MessageInfo) error {
	ch, _ := rmq.NewChannel()
	defer rmq.ChannelDone(ch)
	b, _ := json.Marshal(body)
	log.Tracef("Trying to publish message to %s : %s\n", routingKey, string(b))

	queue, err := rmq.declareQueue(queueName, durable, autoDelete)

	if err != nil {
		log.Error(err)
		os.Exit(0)

	}

	err = rmq.bindQueue(queue, exchange, routingKey)

	if err != nil {
		log.Error(err)
		os.Exit(0)
	}

	err = ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{ContentType: contentType, Body: b, DeliveryMode: amqp.Persistent},
	)

	if err != nil {
		fmt.Println("error publish " + err.Error())
	}

	return err
}

//PublishMessageSync : Publish a message and waiting for response (this function implement request reply pattern)
func (rmq *rabbitMQ) PublishMessageSync(exchange string, routingKey string, contentType string, headers map[string]interface{}, body interface{}, msgInfo *MessageInfo) <-chan Result {
	output := make(chan Result)
	go func() {
		defer close(output)
		chPub, _ := rmq.GetChannel()
		defer rmq.ChannelDone(chPub)

		msgBody := rmq.EncapsulateData(msgInfo, body)

		corrId := uuid.New().String()
		r := fmt.Sprintf("%s.%s", routingKey, corrId)

		q, err := chPub.QueueDeclare(
			r,     // name
			false, // durable
			false, // delete when unused
			true,  // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			output <- Result{Error: fmt.Errorf("error when initalize response queue %s", err.Error())}
			return
		}

		b, _ := json.Marshal(msgBody)
		log.Tracef("Trying to publish message to %s : %s", routingKey, string(b))
		//msgReturn := make(chan amqp.Return)
		err = chPub.Publish(
			exchange,   // exchange
			routingKey, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{ContentType: contentType, Body: b,
				DeliveryMode:  amqp.Persistent,
				CorrelationId: corrId,
				ReplyTo:       q.Name,
			},
		)
		if err != nil {
			log.Error("error when publishing data ", err)
			output <- Result{Error: fmt.Errorf("error when publishing data %s", err.Error())}
			return
		}

		var resp MessageBody
		msg, ch, err := rmq.ReadMessage(q)
		defer ch.Close()
		if err != nil {
			output <- Result{Error: fmt.Errorf("error when reading the message %s", err.Error())}
			return
		}
		for d := range msg {
			if d.CorrelationId == corrId {
				if err := json.Unmarshal(d.Body, &resp); err != nil {
					output <- Result{Error: fmt.Errorf("error when unmarshaling data %s", err.Error())}
					return
				}
				break
			}
		}
		var result Result
		if err := rmq.DecodeMapType(resp.Content, &result); err != nil {
			output <- Result{Error: fmt.Errorf("error when decoding data %s", err.Error())}
			return
		}
		output <- result
	}()
	return output
}

func (rmq *rabbitMQ) PostProcessMessage(d amqp.Delivery, response interface{}) {
	go func() {
		ch, _ := rmq.GetChannel()
		defer rmq.ChannelDone(ch)
		if d.CorrelationId != "" && d.ReplyTo != "" {
			res := rmq.EncapsulateData(nil, response)
			b, _ := json.Marshal(res)
			if err := ch.Publish(DefaultExchange,
				d.ReplyTo,
				false,
				false,
				amqp.Publishing{
					ContentType:   d.ContentType,
					CorrelationId: d.CorrelationId,
					Body:          b,
				},
			); err != nil {
				log.Error("error when publishing response")
			}
		}
	}()
}
