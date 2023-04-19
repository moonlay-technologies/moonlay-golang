// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	amqp "github.com/streadway/amqp"
	mock "github.com/stretchr/testify/mock"

	rabbitmq "order-service/global/utils/rabbitmq"
)

// RabbitMQ is an autogenerated mock type for the RabbitMQ type
type RabbitMQ struct {
	mock.Mock
}

// ChannelDone provides a mock function with given fields: target
func (_m *RabbitMQ) ChannelDone(target *rabbitmq.Channel) {
	_m.Called(target)
}

// ChannelPool provides a mock function with given fields:
func (_m *RabbitMQ) ChannelPool() rabbitmq.ChannelPool {
	ret := _m.Called()

	var r0 rabbitmq.ChannelPool
	if rf, ok := ret.Get(0).(func() rabbitmq.ChannelPool); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(rabbitmq.ChannelPool)
		}
	}

	return r0
}

// Connection provides a mock function with given fields:
func (_m *RabbitMQ) Connection() *rabbitmq.Connection {
	ret := _m.Called()

	var r0 *rabbitmq.Connection
	if rf, ok := ret.Get(0).(func() *rabbitmq.Connection); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rabbitmq.Connection)
		}
	}

	return r0
}

// DecodeMapType provides a mock function with given fields: input, output
func (_m *RabbitMQ) DecodeMapType(input interface{}, output interface{}) error {
	ret := _m.Called(input, output)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, interface{}) error); ok {
		r0 = rf(input, output)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Disconnect provides a mock function with given fields:
func (_m *RabbitMQ) Disconnect() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EncapsulateData provides a mock function with given fields: info, content
func (_m *RabbitMQ) EncapsulateData(info *rabbitmq.MessageInfo, content interface{}) rabbitmq.MessageBody {
	ret := _m.Called(info, content)

	var r0 rabbitmq.MessageBody
	if rf, ok := ret.Get(0).(func(*rabbitmq.MessageInfo, interface{}) rabbitmq.MessageBody); ok {
		r0 = rf(info, content)
	} else {
		r0 = ret.Get(0).(rabbitmq.MessageBody)
	}

	return r0
}

// ExtractMessageData provides a mock function with given fields: msg, out
func (_m *RabbitMQ) ExtractMessageData(msg amqp.Delivery, out interface{}) (*rabbitmq.MessageInfo, error) {
	ret := _m.Called(msg, out)

	var r0 *rabbitmq.MessageInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(amqp.Delivery, interface{}) (*rabbitmq.MessageInfo, error)); ok {
		return rf(msg, out)
	}
	if rf, ok := ret.Get(0).(func(amqp.Delivery, interface{}) *rabbitmq.MessageInfo); ok {
		r0 = rf(msg, out)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rabbitmq.MessageInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(amqp.Delivery, interface{}) error); ok {
		r1 = rf(msg, out)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChannel provides a mock function with given fields:
func (_m *RabbitMQ) GetChannel() (*rabbitmq.Channel, error) {
	ret := _m.Called()

	var r0 *rabbitmq.Channel
	var r1 error
	if rf, ok := ret.Get(0).(func() (*rabbitmq.Channel, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *rabbitmq.Channel); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rabbitmq.Channel)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InitQueue provides a mock function with given fields: queueName, exchangeName, routingKey, durable, autoDelete
func (_m *RabbitMQ) InitQueue(queueName string, exchangeName string, routingKey string, durable bool, autoDelete bool) (amqp.Queue, error) {
	ret := _m.Called(queueName, exchangeName, routingKey, durable, autoDelete)

	var r0 amqp.Queue
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, string, bool, bool) (amqp.Queue, error)); ok {
		return rf(queueName, exchangeName, routingKey, durable, autoDelete)
	}
	if rf, ok := ret.Get(0).(func(string, string, string, bool, bool) amqp.Queue); ok {
		r0 = rf(queueName, exchangeName, routingKey, durable, autoDelete)
	} else {
		r0 = ret.Get(0).(amqp.Queue)
	}

	if rf, ok := ret.Get(1).(func(string, string, string, bool, bool) error); ok {
		r1 = rf(queueName, exchangeName, routingKey, durable, autoDelete)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewChannel provides a mock function with given fields:
func (_m *RabbitMQ) NewChannel() (*rabbitmq.Channel, error) {
	ret := _m.Called()

	var r0 *rabbitmq.Channel
	var r1 error
	if rf, ok := ret.Get(0).(func() (*rabbitmq.Channel, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *rabbitmq.Channel); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rabbitmq.Channel)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PostProcessMessage provides a mock function with given fields: d, response
func (_m *RabbitMQ) PostProcessMessage(d amqp.Delivery, response interface{}) {
	_m.Called(d, response)
}

// PublishMessage provides a mock function with given fields: queueName, durable, autoDelete, exchange, routingKey, contentType, headers, body, msgInfo
func (_m *RabbitMQ) PublishMessage(queueName string, durable bool, autoDelete bool, exchange string, routingKey string, contentType string, headers map[string]interface{}, body interface{}, msgInfo *rabbitmq.MessageInfo) error {
	ret := _m.Called(queueName, durable, autoDelete, exchange, routingKey, contentType, headers, body, msgInfo)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, bool, bool, string, string, string, map[string]interface{}, interface{}, *rabbitmq.MessageInfo) error); ok {
		r0 = rf(queueName, durable, autoDelete, exchange, routingKey, contentType, headers, body, msgInfo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PublishMessageSync provides a mock function with given fields: exchange, routingKey, contentType, headers, body, msgInfo
func (_m *RabbitMQ) PublishMessageSync(exchange string, routingKey string, contentType string, headers map[string]interface{}, body interface{}, msgInfo *rabbitmq.MessageInfo) <-chan rabbitmq.Result {
	ret := _m.Called(exchange, routingKey, contentType, headers, body, msgInfo)

	var r0 <-chan rabbitmq.Result
	if rf, ok := ret.Get(0).(func(string, string, string, map[string]interface{}, interface{}, *rabbitmq.MessageInfo) <-chan rabbitmq.Result); ok {
		r0 = rf(exchange, routingKey, contentType, headers, body, msgInfo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan rabbitmq.Result)
		}
	}

	return r0
}

// ReadMessage provides a mock function with given fields: q
func (_m *RabbitMQ) ReadMessage(q amqp.Queue) (<-chan amqp.Delivery, *rabbitmq.Channel, error) {
	ret := _m.Called(q)

	var r0 <-chan amqp.Delivery
	var r1 *rabbitmq.Channel
	var r2 error
	if rf, ok := ret.Get(0).(func(amqp.Queue) (<-chan amqp.Delivery, *rabbitmq.Channel, error)); ok {
		return rf(q)
	}
	if rf, ok := ret.Get(0).(func(amqp.Queue) <-chan amqp.Delivery); ok {
		r0 = rf(q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan amqp.Delivery)
		}
	}

	if rf, ok := ret.Get(1).(func(amqp.Queue) *rabbitmq.Channel); ok {
		r1 = rf(q)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*rabbitmq.Channel)
		}
	}

	if rf, ok := ret.Get(2).(func(amqp.Queue) error); ok {
		r2 = rf(q)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// RegisterExchange provides a mock function with given fields: exchangeName, exchangeType
func (_m *RabbitMQ) RegisterExchange(exchangeName string, exchangeType string) error {
	ret := _m.Called(exchangeName, exchangeType)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(exchangeName, exchangeType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ServiceCode provides a mock function with given fields:
func (_m *RabbitMQ) ServiceCode() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SetLogLevel provides a mock function with given fields: logLevel
func (_m *RabbitMQ) SetLogLevel(logLevel string) {
	_m.Called(logLevel)
}

type mockConstructorTestingTNewRabbitMQ interface {
	mock.TestingT
	Cleanup(func())
}

// NewRabbitMQ creates a new instance of RabbitMQ. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRabbitMQ(t mockConstructorTestingTNewRabbitMQ) *RabbitMQ {
	mock := &RabbitMQ{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
