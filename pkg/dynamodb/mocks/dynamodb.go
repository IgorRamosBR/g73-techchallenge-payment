// Code generated by MockGen. DO NOT EDIT.
// Source: dynamodb.go
//
// Generated by this command:
//
//	mockgen -source=dynamodb.go -destination=mocks/dynamodb.go
//
// Package mock_dynamodb is a generated GoMock package.
package mock_dynamodb

import (
	reflect "reflect"

	expression "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	types "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	gomock "go.uber.org/mock/gomock"
)

// MockDynamoDBClient is a mock of DynamoDBClient interface.
type MockDynamoDBClient struct {
	ctrl     *gomock.Controller
	recorder *MockDynamoDBClientMockRecorder
}

// MockDynamoDBClientMockRecorder is the mock recorder for MockDynamoDBClient.
type MockDynamoDBClientMockRecorder struct {
	mock *MockDynamoDBClient
}

// NewMockDynamoDBClient creates a new mock instance.
func NewMockDynamoDBClient(ctrl *gomock.Controller) *MockDynamoDBClient {
	mock := &MockDynamoDBClient{ctrl: ctrl}
	mock.recorder = &MockDynamoDBClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDynamoDBClient) EXPECT() *MockDynamoDBClientMockRecorder {
	return m.recorder
}

// GetItem mocks base method.
func (m *MockDynamoDBClient) GetItem(tableName string, key map[string]types.AttributeValue) (map[string]types.AttributeValue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItem", tableName, key)
	ret0, _ := ret[0].(map[string]types.AttributeValue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItem indicates an expected call of GetItem.
func (mr *MockDynamoDBClientMockRecorder) GetItem(tableName, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItem", reflect.TypeOf((*MockDynamoDBClient)(nil).GetItem), tableName, key)
}

// PutItem mocks base method.
func (m *MockDynamoDBClient) PutItem(tableName string, item map[string]types.AttributeValue) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutItem", tableName, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutItem indicates an expected call of PutItem.
func (mr *MockDynamoDBClientMockRecorder) PutItem(tableName, item any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutItem", reflect.TypeOf((*MockDynamoDBClient)(nil).PutItem), tableName, item)
}

// QueryItem mocks base method.
func (m *MockDynamoDBClient) QueryItem(tableName string, expr expression.Expression, indexName string) ([]map[string]types.AttributeValue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryItem", tableName, expr, indexName)
	ret0, _ := ret[0].([]map[string]types.AttributeValue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryItem indicates an expected call of QueryItem.
func (mr *MockDynamoDBClientMockRecorder) QueryItem(tableName, expr, indexName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryItem", reflect.TypeOf((*MockDynamoDBClient)(nil).QueryItem), tableName, expr, indexName)
}

// UpdateItem mocks base method.
func (m *MockDynamoDBClient) UpdateItem(tableName string, key map[string]types.AttributeValue, expr expression.Expression) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateItem", tableName, key, expr)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateItem indicates an expected call of UpdateItem.
func (mr *MockDynamoDBClientMockRecorder) UpdateItem(tableName, key, expr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateItem", reflect.TypeOf((*MockDynamoDBClient)(nil).UpdateItem), tableName, key, expr)
}
