// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// TelegramBotInterface is an autogenerated mock type for the TelegramBotInterface type
type TelegramBotInterface struct {
	mock.Mock
}

// SendMessage provides a mock function with given fields: messages, result
func (_m *TelegramBotInterface) SendMessage(messages string, result chan error) {
	_m.Called(messages, result)
}

// SetChatID provides a mock function with given fields: ChatID
func (_m *TelegramBotInterface) SetChatID(ChatID string) {
	_m.Called(ChatID)
}

// SetToken provides a mock function with given fields: token
func (_m *TelegramBotInterface) SetToken(token string) {
	_m.Called(token)
}

type mockConstructorTestingTNewTelegramBotInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewTelegramBotInterface creates a new instance of TelegramBotInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTelegramBotInterface(t mockConstructorTestingTNewTelegramBotInterface) *TelegramBotInterface {
	mock := &TelegramBotInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
