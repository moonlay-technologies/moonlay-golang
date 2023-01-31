package rabbitmq

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sync"
	"sync/atomic"
	"time"
)

// Connection amqp.Connection wrapper
type Connection struct {
	*amqp.Connection
	sync.Mutex
}

// Dial wrap amqp.Dial, Dial and get a reconnect connection
func NewConnection(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	connection := &Connection{
		Connection: conn,
	}

	go func() {
		for {
			reason, ok := <-connection.Connection.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok {
				logrus.Info("connection closed")
				break
			}
			logrus.Warnf("connection closed, reason: %v", reason)

			// reconnect if not closed by developer
			for {
				conn, err := amqp.Dial(url)
				if err == nil {
					connection.Lock()
					connection.Connection = conn
					connection.Unlock()
					logrus.Info("connection recovery success")
					break
				}

				logrus.Errorf("connection recovery failed, err: %v", err)
				// wait a moment for reconnect
				time.Sleep(delay)
			}
		}
	}()

	return connection, nil
}

// Channel amqp.Channel wapper
type Channel struct {
	*amqp.Channel
	ID     string
	closed int32
	sync.Mutex
	ServiceCode string
	lastUsed    time.Time
}

// Channel wrap amqp.Connection.Channel, get a auto reconnect chPool
func (c *Connection) Channel(serviceCode string) (*Channel, error) {
	c.Lock()
	ch, err := c.Connection.Channel()
	c.Unlock()
	if err != nil {
		return nil, err
	}

	channel := &Channel{
		ID:          uuid.New().String(),
		Channel:     ch,
		ServiceCode: serviceCode,
		lastUsed:    time.Now(),
	}
	go func() {
		for {
			// chPool.Lock()
			reason, ok := <-channel.Channel.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok || channel.IsClosed() {
				logrus.Infof("channel with id %s closed", channel.ID)
				channel.Close() // close again, ensure closed flag set when connection closed
				break
			}

			logrus.Warnf("channel id %s closed, reason: %v", channel.ID, reason)

			// reconnect if not closed by developer
			for {
				// wait 1s for connection reconnect
				c.Lock()
				ch, err := c.Connection.Channel()
				c.Unlock()
				if err == nil {
					channel.Lock()
					logrus.Infof("channel with id %s is successfuly recovered", channel.ID)
					channel.Channel = ch
					channel.Unlock()
					break
				}
				logrus.Errorf("failed to recover channel with id %s, err: %v\n", channel.ID, err)
				time.Sleep(delay)
			}
			// chPool.Unlock()
			// time.Sleep(5 * time.Second)
		}
	}()
	return channel, nil
}

func (ch *Channel) LastUsedNow() {
	ch.Lock()
	ch.lastUsed = time.Now()
	ch.Unlock()
}

// IsClosed indicate closed by developer
func (ch *Channel) IsClosed() bool {
	return atomic.LoadInt32(&ch.closed) == 1
}

// Close ensure closed flag set
func (ch *Channel) Close() error {
	if ch.IsClosed() {
		return amqp.ErrClosed
	}
	ch.Lock()
	atomic.StoreInt32(&ch.closed, 1)
	err := ch.Channel.Close()
	ch.Unlock()
	return err
}

// Consume warp amqp.Channel.Consume, the returned delivery will end only when channel closed by developer
func (ch *Channel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	deliveries := make(chan amqp.Delivery)

	go func() {
		for {
			ch.Lock()
			c := ch
			ch.Unlock()
			d, err := c.Channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
			if err != nil {
				logrus.Errorf("consume failed, err: %v", err)
				time.Sleep(delay)
				continue
			}

			for msg := range d {
				deliveries <- msg
			}

			// sleep before IsClose call. closed flag may not set before sleep.
			time.Sleep(delay)
			if ch.IsClosed() {
				ch.Unlock()
				break
			}
		}
	}()

	return deliveries, nil
}
