// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	rabbitmq "order-service/global/utils/rabbitmq"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// ChannelPool is an autogenerated mock type for the ChannelPool type
type ChannelPool struct {
	mock.Mock
}

// CloseAll provides a mock function with given fields:
func (_m *ChannelPool) CloseAll() {
	_m.Called()
}

// Done provides a mock function with given fields: target
func (_m *ChannelPool) Done(target *rabbitmq.Channel) {
	_m.Called(target)
}

// GetChannel provides a mock function with given fields:
func (_m *ChannelPool) GetChannel() (*rabbitmq.Channel, error) {
	ret := _m.Called()

	var r0 *rabbitmq.Channel
	if rf, ok := ret.Get(0).(func() *rabbitmq.Channel); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rabbitmq.Channel)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetChannelIdleTimeout provides a mock function with given fields: timeout
func (_m *ChannelPool) SetChannelIdleTimeout(timeout time.Duration) {
	_m.Called(timeout)
}

// SetMaxPool provides a mock function with given fields: maxPool
func (_m *ChannelPool) SetMaxPool(maxPool int) {
	_m.Called(maxPool)
}

// SetMinIdleChannel provides a mock function with given fields: minIdleChannel
func (_m *ChannelPool) SetMinIdleChannel(minIdleChannel int) {
	_m.Called(minIdleChannel)
}

type mockConstructorTestingTNewChannelPool interface {
	mock.TestingT
	Cleanup(func())
}

// NewChannelPool creates a new instance of ChannelPool. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewChannelPool(t mockConstructorTestingTNewChannelPool) *ChannelPool {
	mock := &ChannelPool{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
