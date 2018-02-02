// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/trussle/snowy/pkg/repository (interfaces: Repository)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/trussle/snowy/pkg/models"
	repository "github.com/trussle/snowy/pkg/repository"
	uuid "github.com/trussle/uuid"
	reflect "reflect"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AppendLedger mocks base method
func (m *MockRepository) AppendLedger(arg0 uuid.UUID, arg1 models.Ledger) (models.Ledger, error) {
	ret := m.ctrl.Call(m, "AppendLedger", arg0, arg1)
	ret0, _ := ret[0].(models.Ledger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AppendLedger indicates an expected call of AppendLedger
func (mr *MockRepositoryMockRecorder) AppendLedger(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendLedger", reflect.TypeOf((*MockRepository)(nil).AppendLedger), arg0, arg1)
}

// Close mocks base method
func (m *MockRepository) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockRepositoryMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRepository)(nil).Close))
}

// ForkLedger mocks base method
func (m *MockRepository) ForkLedger(arg0 uuid.UUID, arg1 models.Ledger) (models.Ledger, error) {
	ret := m.ctrl.Call(m, "ForkLedger", arg0, arg1)
	ret0, _ := ret[0].(models.Ledger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ForkLedger indicates an expected call of ForkLedger
func (mr *MockRepositoryMockRecorder) ForkLedger(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForkLedger", reflect.TypeOf((*MockRepository)(nil).ForkLedger), arg0, arg1)
}

// InsertLedger mocks base method
func (m *MockRepository) InsertLedger(arg0 models.Ledger) (models.Ledger, error) {
	ret := m.ctrl.Call(m, "InsertLedger", arg0)
	ret0, _ := ret[0].(models.Ledger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertLedger indicates an expected call of InsertLedger
func (mr *MockRepositoryMockRecorder) InsertLedger(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertLedger", reflect.TypeOf((*MockRepository)(nil).InsertLedger), arg0)
}

// LedgerStatistics mocks base method
func (m *MockRepository) LedgerStatistics() (models.LedgerStatistics, error) {
	ret := m.ctrl.Call(m, "LedgerStatistics")
	ret0, _ := ret[0].(models.LedgerStatistics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LedgerStatistics indicates an expected call of LedgerStatistics
func (mr *MockRepositoryMockRecorder) LedgerStatistics() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LedgerStatistics", reflect.TypeOf((*MockRepository)(nil).LedgerStatistics))
}

// PutContent mocks base method
func (m *MockRepository) PutContent(arg0 models.Content) (models.Content, error) {
	ret := m.ctrl.Call(m, "PutContent", arg0)
	ret0, _ := ret[0].(models.Content)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutContent indicates an expected call of PutContent
func (mr *MockRepositoryMockRecorder) PutContent(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutContent", reflect.TypeOf((*MockRepository)(nil).PutContent), arg0)
}

// SelectContent mocks base method
func (m *MockRepository) SelectContent(arg0 uuid.UUID, arg1 repository.Query) (models.Content, error) {
	ret := m.ctrl.Call(m, "SelectContent", arg0, arg1)
	ret0, _ := ret[0].(models.Content)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectContent indicates an expected call of SelectContent
func (mr *MockRepositoryMockRecorder) SelectContent(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectContent", reflect.TypeOf((*MockRepository)(nil).SelectContent), arg0, arg1)
}

// SelectContents mocks base method
func (m *MockRepository) SelectContents(arg0 uuid.UUID, arg1 repository.Query) ([]models.Content, error) {
	ret := m.ctrl.Call(m, "SelectContents", arg0, arg1)
	ret0, _ := ret[0].([]models.Content)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectContents indicates an expected call of SelectContents
func (mr *MockRepositoryMockRecorder) SelectContents(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectContents", reflect.TypeOf((*MockRepository)(nil).SelectContents), arg0, arg1)
}

// SelectForkLedgers mocks base method
func (m *MockRepository) SelectForkLedgers(arg0 uuid.UUID) ([]models.Ledger, error) {
	ret := m.ctrl.Call(m, "SelectForkLedgers", arg0)
	ret0, _ := ret[0].([]models.Ledger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectForkLedgers indicates an expected call of SelectForkLedgers
func (mr *MockRepositoryMockRecorder) SelectForkLedgers(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectForkLedgers", reflect.TypeOf((*MockRepository)(nil).SelectForkLedgers), arg0)
}

// SelectLedger mocks base method
func (m *MockRepository) SelectLedger(arg0 uuid.UUID, arg1 repository.Query) (models.Ledger, error) {
	ret := m.ctrl.Call(m, "SelectLedger", arg0, arg1)
	ret0, _ := ret[0].(models.Ledger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectLedger indicates an expected call of SelectLedger
func (mr *MockRepositoryMockRecorder) SelectLedger(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectLedger", reflect.TypeOf((*MockRepository)(nil).SelectLedger), arg0, arg1)
}

// SelectLedgers mocks base method
func (m *MockRepository) SelectLedgers(arg0 uuid.UUID, arg1 repository.Query) ([]models.Ledger, error) {
	ret := m.ctrl.Call(m, "SelectLedgers", arg0, arg1)
	ret0, _ := ret[0].([]models.Ledger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectLedgers indicates an expected call of SelectLedgers
func (mr *MockRepositoryMockRecorder) SelectLedgers(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectLedgers", reflect.TypeOf((*MockRepository)(nil).SelectLedgers), arg0, arg1)
}
