// Code generated by MockGen. DO NOT EDIT.
// Source: scheduler/expected_interfaces.go
//
// Generated by this command:
//
//	mockgen -source scheduler/expected_interfaces.go -destination scheduler/testutils/expected_interfaces.mock.go -package testutils
//
// Package testutils is a generated GoMock package.
package testutils

import (
	reflect "reflect"

	types "github.com/cosmos/cosmos-sdk/types"
	types0 "github.com/desmos-labs/caerus/types"
	types1 "github.com/desmos-labs/cosmos-go-wallet/types"
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

// GetApp mocks base method.
func (m *MockDatabase) GetApp(appID string) (*types0.Application, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApp", appID)
	ret0, _ := ret[0].(*types0.Application)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetApp indicates an expected call of GetApp.
func (mr *MockDatabaseMockRecorder) GetApp(appID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApp", reflect.TypeOf((*MockDatabase)(nil).GetApp), appID)
}

// GetNotGrantedFeeGrantRequests mocks base method.
func (m *MockDatabase) GetNotGrantedFeeGrantRequests(limit int) ([]types0.FeeGrantRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNotGrantedFeeGrantRequests", limit)
	ret0, _ := ret[0].([]types0.FeeGrantRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNotGrantedFeeGrantRequests indicates an expected call of GetNotGrantedFeeGrantRequests.
func (mr *MockDatabaseMockRecorder) GetNotGrantedFeeGrantRequests(limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNotGrantedFeeGrantRequests", reflect.TypeOf((*MockDatabase)(nil).GetNotGrantedFeeGrantRequests), limit)
}

// SetFeeGrantRequestsGranted mocks base method.
func (m *MockDatabase) SetFeeGrantRequestsGranted(ids []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetFeeGrantRequestsGranted", ids)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetFeeGrantRequestsGranted indicates an expected call of SetFeeGrantRequestsGranted.
func (mr *MockDatabaseMockRecorder) SetFeeGrantRequestsGranted(ids any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFeeGrantRequestsGranted", reflect.TypeOf((*MockDatabase)(nil).SetFeeGrantRequestsGranted), ids)
}

// MockFirebase is a mock of Firebase interface.
type MockFirebase struct {
	ctrl     *gomock.Controller
	recorder *MockFirebaseMockRecorder
}

// MockFirebaseMockRecorder is the mock recorder for MockFirebase.
type MockFirebaseMockRecorder struct {
	mock *MockFirebase
}

// NewMockFirebase creates a new mock instance.
func NewMockFirebase(ctrl *gomock.Controller) *MockFirebase {
	mock := &MockFirebase{ctrl: ctrl}
	mock.recorder = &MockFirebaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFirebase) EXPECT() *MockFirebaseMockRecorder {
	return m.recorder
}

// SendNotificationToApp mocks base method.
func (m *MockFirebase) SendNotificationToApp(appID string, notification types0.Notification) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendNotificationToApp", appID, notification)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendNotificationToApp indicates an expected call of SendNotificationToApp.
func (mr *MockFirebaseMockRecorder) SendNotificationToApp(appID, notification any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendNotificationToApp", reflect.TypeOf((*MockFirebase)(nil).SendNotificationToApp), appID, notification)
}

// MockChainClient is a mock of ChainClient interface.
type MockChainClient struct {
	ctrl     *gomock.Controller
	recorder *MockChainClientMockRecorder
}

// MockChainClientMockRecorder is the mock recorder for MockChainClient.
type MockChainClientMockRecorder struct {
	mock *MockChainClient
}

// NewMockChainClient creates a new mock instance.
func NewMockChainClient(ctrl *gomock.Controller) *MockChainClient {
	mock := &MockChainClient{ctrl: ctrl}
	mock.recorder = &MockChainClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChainClient) EXPECT() *MockChainClientMockRecorder {
	return m.recorder
}

// AccAddress mocks base method.
func (m *MockChainClient) AccAddress() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccAddress")
	ret0, _ := ret[0].(string)
	return ret0
}

// AccAddress indicates an expected call of AccAddress.
func (mr *MockChainClientMockRecorder) AccAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccAddress", reflect.TypeOf((*MockChainClient)(nil).AccAddress))
}

// BroadcastTxCommit mocks base method.
func (m *MockChainClient) BroadcastTxCommit(data *types1.TransactionData) (*types.TxResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BroadcastTxCommit", data)
	ret0, _ := ret[0].(*types.TxResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BroadcastTxCommit indicates an expected call of BroadcastTxCommit.
func (mr *MockChainClientMockRecorder) BroadcastTxCommit(data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BroadcastTxCommit", reflect.TypeOf((*MockChainClient)(nil).BroadcastTxCommit), data)
}

// HasFeeGrant mocks base method.
func (m *MockChainClient) HasFeeGrant(granteeAddress, granterAddress string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasFeeGrant", granteeAddress, granterAddress)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasFeeGrant indicates an expected call of HasFeeGrant.
func (mr *MockChainClientMockRecorder) HasFeeGrant(granteeAddress, granterAddress any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasFeeGrant", reflect.TypeOf((*MockChainClient)(nil).HasFeeGrant), granteeAddress, granterAddress)
}

// HasGrantedMsgGrantAllowanceAuthorization mocks base method.
func (m *MockChainClient) HasGrantedMsgGrantAllowanceAuthorization(appAddress string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasGrantedMsgGrantAllowanceAuthorization", appAddress)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasGrantedMsgGrantAllowanceAuthorization indicates an expected call of HasGrantedMsgGrantAllowanceAuthorization.
func (mr *MockChainClientMockRecorder) HasGrantedMsgGrantAllowanceAuthorization(appAddress any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasGrantedMsgGrantAllowanceAuthorization", reflect.TypeOf((*MockChainClient)(nil).HasGrantedMsgGrantAllowanceAuthorization), appAddress)
}
