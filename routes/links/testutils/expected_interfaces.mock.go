// Code generated by MockGen. DO NOT EDIT.
// Source: routes/links/expected_interfaces.go

// Package testutils is a generated GoMock package.
package testutils

import (
	reflect "reflect"

	types "github.com/desmos-labs/caerus/types"
	gomock "go.uber.org/mock/gomock"
)

// MockDatabase is a mock of Database interface.
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase.
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance.
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return m.recorder
}

// GetAppDeepLinksCount mocks base method.
func (m *MockDatabase) GetAppDeepLinksCount(appID string) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAppDeepLinksCount", appID)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAppDeepLinksCount indicates an expected call of GetAppDeepLinksCount.
func (mr *MockDatabaseMockRecorder) GetAppDeepLinksCount(appID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAppDeepLinksCount", reflect.TypeOf((*MockDatabase)(nil).GetAppDeepLinksCount), appID)
}

// GetAppDeepLinksRateLimit mocks base method.
func (m *MockDatabase) GetAppDeepLinksRateLimit(appID string) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAppDeepLinksRateLimit", appID)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAppDeepLinksRateLimit indicates an expected call of GetAppDeepLinksRateLimit.
func (mr *MockDatabaseMockRecorder) GetAppDeepLinksRateLimit(appID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAppDeepLinksRateLimit", reflect.TypeOf((*MockDatabase)(nil).GetAppDeepLinksRateLimit), appID)
}

// GetDeepLinkConfig mocks base method.
func (m *MockDatabase) GetDeepLinkConfig(link string) (*types.LinkConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeepLinkConfig", link)
	ret0, _ := ret[0].(*types.LinkConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeepLinkConfig indicates an expected call of GetDeepLinkConfig.
func (mr *MockDatabaseMockRecorder) GetDeepLinkConfig(link interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeepLinkConfig", reflect.TypeOf((*MockDatabase)(nil).GetDeepLinkConfig), link)
}

// SaveCreatedDeepLink mocks base method.
func (m *MockDatabase) SaveCreatedDeepLink(link *types.CreatedDeepLink) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveCreatedDeepLink", link)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveCreatedDeepLink indicates an expected call of SaveCreatedDeepLink.
func (mr *MockDatabaseMockRecorder) SaveCreatedDeepLink(link interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveCreatedDeepLink", reflect.TypeOf((*MockDatabase)(nil).SaveCreatedDeepLink), link)
}

// MockDeepLinksClient is a mock of DeepLinksClient interface.
type MockDeepLinksClient struct {
	ctrl     *gomock.Controller
	recorder *MockDeepLinksClientMockRecorder
}

// MockDeepLinksClientMockRecorder is the mock recorder for MockDeepLinksClient.
type MockDeepLinksClientMockRecorder struct {
	mock *MockDeepLinksClient
}

// NewMockDeepLinksClient creates a new mock instance.
func NewMockDeepLinksClient(ctrl *gomock.Controller) *MockDeepLinksClient {
	mock := &MockDeepLinksClient{ctrl: ctrl}
	mock.recorder = &MockDeepLinksClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeepLinksClient) EXPECT() *MockDeepLinksClientMockRecorder {
	return m.recorder
}

// CreateDynamicLink mocks base method.
func (m *MockDeepLinksClient) CreateDynamicLink(apiKey string, linkConfig *types.LinkConfig) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDynamicLink", apiKey, linkConfig)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDynamicLink indicates an expected call of CreateDynamicLink.
func (mr *MockDeepLinksClientMockRecorder) CreateDynamicLink(apiKey, linkConfig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDynamicLink", reflect.TypeOf((*MockDeepLinksClient)(nil).CreateDynamicLink), apiKey, linkConfig)
}
