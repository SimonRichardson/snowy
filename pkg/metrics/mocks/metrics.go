// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/trussle/snowy/pkg/metrics (interfaces: Gauge,HistogramVec,Counter)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
	prometheus "github.com/prometheus/client_golang/prometheus"
	reflect "reflect"
)

// MockGauge is a mock of Gauge interface
type MockGauge struct {
	ctrl     *gomock.Controller
	recorder *MockGaugeMockRecorder
}

// MockGaugeMockRecorder is the mock recorder for MockGauge
type MockGaugeMockRecorder struct {
	mock *MockGauge
}

// NewMockGauge creates a new mock instance
func NewMockGauge(ctrl *gomock.Controller) *MockGauge {
	mock := &MockGauge{ctrl: ctrl}
	mock.recorder = &MockGaugeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockGauge) EXPECT() *MockGaugeMockRecorder {
	return _m.recorder
}

// Dec mocks base method
func (_m *MockGauge) Dec() {
	_m.ctrl.Call(_m, "Dec")
}

// Dec indicates an expected call of Dec
func (_mr *MockGaugeMockRecorder) Dec() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Dec", reflect.TypeOf((*MockGauge)(nil).Dec))
}

// Inc mocks base method
func (_m *MockGauge) Inc() {
	_m.ctrl.Call(_m, "Inc")
}

// Inc indicates an expected call of Inc
func (_mr *MockGaugeMockRecorder) Inc() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Inc", reflect.TypeOf((*MockGauge)(nil).Inc))
}

// MockHistogramVec is a mock of HistogramVec interface
type MockHistogramVec struct {
	ctrl     *gomock.Controller
	recorder *MockHistogramVecMockRecorder
}

// MockHistogramVecMockRecorder is the mock recorder for MockHistogramVec
type MockHistogramVecMockRecorder struct {
	mock *MockHistogramVec
}

// NewMockHistogramVec creates a new mock instance
func NewMockHistogramVec(ctrl *gomock.Controller) *MockHistogramVec {
	mock := &MockHistogramVec{ctrl: ctrl}
	mock.recorder = &MockHistogramVecMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockHistogramVec) EXPECT() *MockHistogramVecMockRecorder {
	return _m.recorder
}

// WithLabelValues mocks base method
func (_m *MockHistogramVec) WithLabelValues(_param0 ...string) prometheus.Observer {
	_s := []interface{}{}
	for _, _x := range _param0 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "WithLabelValues", _s...)
	ret0, _ := ret[0].(prometheus.Observer)
	return ret0
}

// WithLabelValues indicates an expected call of WithLabelValues
func (_mr *MockHistogramVecMockRecorder) WithLabelValues(arg0 ...interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "WithLabelValues", reflect.TypeOf((*MockHistogramVec)(nil).WithLabelValues), arg0...)
}

// MockCounter is a mock of Counter interface
type MockCounter struct {
	ctrl     *gomock.Controller
	recorder *MockCounterMockRecorder
}

// MockCounterMockRecorder is the mock recorder for MockCounter
type MockCounterMockRecorder struct {
	mock *MockCounter
}

// NewMockCounter creates a new mock instance
func NewMockCounter(ctrl *gomock.Controller) *MockCounter {
	mock := &MockCounter{ctrl: ctrl}
	mock.recorder = &MockCounterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockCounter) EXPECT() *MockCounterMockRecorder {
	return _m.recorder
}

// Add mocks base method
func (_m *MockCounter) Add(_param0 float64) {
	_m.ctrl.Call(_m, "Add", _param0)
}

// Add indicates an expected call of Add
func (_mr *MockCounterMockRecorder) Add(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Add", reflect.TypeOf((*MockCounter)(nil).Add), arg0)
}

// Inc mocks base method
func (_m *MockCounter) Inc() {
	_m.ctrl.Call(_m, "Inc")
}

// Inc indicates an expected call of Inc
func (_mr *MockCounterMockRecorder) Inc() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Inc", reflect.TypeOf((*MockCounter)(nil).Inc))
}
