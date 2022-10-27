package mocks

import (
	gomock "github.com/golang/mock/gomock"
	acceptor "github.com/gotechbook/gotechbook-framework-acceptor"
	net "net"
	reflect "reflect"
	time "time"
)

// MockAcceptor is a mock of Acceptor interface.
type MockAcceptor struct {
	ctrl     *gomock.Controller
	recorder *MockAcceptorMockRecorder
}

// MockAcceptorMockRecorder is the mock recorder for MockAcceptor.
type MockAcceptorMockRecorder struct {
	mock *MockAcceptor
}

// NewMockAcceptor creates a new mock instance.
func NewMockAcceptor(ctrl *gomock.Controller) *MockAcceptor {
	mock := &MockAcceptor{ctrl: ctrl}
	mock.recorder = &MockAcceptorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAcceptor) EXPECT() *MockAcceptorMockRecorder {
	return m.recorder
}

// GetAddr mocks base method.
func (m *MockAcceptor) GetAddr() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAddr")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetAddr indicates an expected call of GetAddr.
func (mr *MockAcceptorMockRecorder) GetAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAddr", reflect.TypeOf((*MockAcceptor)(nil).GetAddr))
}

// GetConnChan mocks base method.
func (m *MockAcceptor) GetConnChan() chan acceptor.Conn {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnChan")
	ret0, _ := ret[0].(chan acceptor.Conn)
	return ret0
}

// GetConnChan indicates an expected call of GetConnChan.
func (mr *MockAcceptorMockRecorder) GetConnChan() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnChan", reflect.TypeOf((*MockAcceptor)(nil).GetConnChan))
}

// ListenAndServe mocks base method.
func (m *MockAcceptor) ListenAndServe() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ListenAndServe")
}

// ListenAndServe indicates an expected call of ListenAndServe.
func (mr *MockAcceptorMockRecorder) ListenAndServe() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenAndServe", reflect.TypeOf((*MockAcceptor)(nil).ListenAndServe))
}

// Stop mocks base method.
func (m *MockAcceptor) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockAcceptorMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockAcceptor)(nil).Stop))
}

// MockConn is a mock of Conn interface.
type MockConn struct {
	ctrl     *gomock.Controller
	recorder *MockConnMockRecorder
}

// MockConnMockRecorder is the mock recorder for MockConn.
type MockConnMockRecorder struct {
	mock *MockConn
}

// NewMockConn creates a new mock instance.
func NewMockConn(ctrl *gomock.Controller) *MockConn {
	mock := &MockConn{ctrl: ctrl}
	mock.recorder = &MockConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConn) EXPECT() *MockConnMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockConn) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockConnMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConn)(nil).Close))
}

// GetNextMessage mocks base method.
func (m *MockConn) GetNextMessage() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNextMessage")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNextMessage indicates an expected call of GetNextMessage.
func (mr *MockConnMockRecorder) GetNextMessage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNextMessage", reflect.TypeOf((*MockConn)(nil).GetNextMessage))
}

// LocalAddr mocks base method.
func (m *MockConn) LocalAddr() net.Addr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LocalAddr")
	ret0, _ := ret[0].(net.Addr)
	return ret0
}

// LocalAddr indicates an expected call of LocalAddr.
func (mr *MockConnMockRecorder) LocalAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LocalAddr", reflect.TypeOf((*MockConn)(nil).LocalAddr))
}

// Read mocks base method.
func (m *MockConn) Read(b []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", b)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockConnMockRecorder) Read(b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockConn)(nil).Read), b)
}

// RemoteAddr mocks base method.
func (m *MockConn) RemoteAddr() net.Addr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoteAddr")
	ret0, _ := ret[0].(net.Addr)
	return ret0
}

// RemoteAddr indicates an expected call of RemoteAddr.
func (mr *MockConnMockRecorder) RemoteAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoteAddr", reflect.TypeOf((*MockConn)(nil).RemoteAddr))
}

// SetDeadline mocks base method.
func (m *MockConn) SetDeadline(t time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetDeadline", t)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDeadline indicates an expected call of SetDeadline.
func (mr *MockConnMockRecorder) SetDeadline(t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDeadline", reflect.TypeOf((*MockConn)(nil).SetDeadline), t)
}

// SetReadDeadline mocks base method.
func (m *MockConn) SetReadDeadline(t time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetReadDeadline", t)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetReadDeadline indicates an expected call of SetReadDeadline.
func (mr *MockConnMockRecorder) SetReadDeadline(t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetReadDeadline", reflect.TypeOf((*MockConn)(nil).SetReadDeadline), t)
}

// SetWriteDeadline mocks base method.
func (m *MockConn) SetWriteDeadline(t time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetWriteDeadline", t)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetWriteDeadline indicates an expected call of SetWriteDeadline.
func (mr *MockConnMockRecorder) SetWriteDeadline(t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWriteDeadline", reflect.TypeOf((*MockConn)(nil).SetWriteDeadline), t)
}

// Write mocks base method.
func (m *MockConn) Write(b []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", b)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockConnMockRecorder) Write(b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockConn)(nil).Write), b)
}

// MockCodec is a mock of Codec interface.
type MockCodec struct {
	ctrl     *gomock.Controller
	recorder *MockCodecMockRecorder
}

// MockCodecMockRecorder is the mock recorder for MockCodec.
type MockCodecMockRecorder struct {
	mock *MockCodec
}

// NewMockCodec creates a new mock instance.
func NewMockCodec(ctrl *gomock.Controller) *MockCodec {
	mock := &MockCodec{ctrl: ctrl}
	mock.recorder = &MockCodecMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCodec) EXPECT() *MockCodecMockRecorder {
	return m.recorder
}

// Decode mocks base method.
func (m *MockCodec) Decode(data []byte) ([]*acceptor.Packet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decode", data)
	ret0, _ := ret[0].([]*acceptor.Packet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Decode indicates an expected call of Decode.
func (mr *MockCodecMockRecorder) Decode(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decode", reflect.TypeOf((*MockCodec)(nil).Decode), data)
}

// Encode mocks base method.
func (m *MockCodec) Encode(typType, data []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Encode", typType, data)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Encode indicates an expected call of Encode.
func (mr *MockCodecMockRecorder) Encode(typType, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Encode", reflect.TypeOf((*MockCodec)(nil).Encode), typType, data)
}
