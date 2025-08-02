package core

import (
	"net"
	"time"
)

// MockObject represents a test object that can change its thread assignment
// Flow: Object creation -> Thread assignment -> Dynamic thread changes
// Purpose: Test thread reassignment functionality
// Error Conditions: Invalid thread assignment, state changes
type MockObject struct {
	id       int
	threadID int // Simulates changing thread assignment
}

// GetThreadHash returns the current thread assignment for this mock object
// Flow: Thread assignment check -> Hash generation
// Thread Assignment: Uses threadID for dynamic thread assignment
// Error Handling: Returns valid hash value
func (m *MockObject) GetThreadHash() int {
	return m.threadID
}

// Ensure MockObject implements ThreadAssignable
var _ ThreadAssignable = (*MockObject)(nil)

// MockConnection represents a mock network connection for testing
// Flow: Mock connection creation -> Socket testing
// Purpose: Test socket functionality without real network connections
// Error Conditions: Mock connection errors
type MockConnection struct {
	remoteAddr string
}

// Read implements net.Conn interface
func (m *MockConnection) Read(b []byte) (n int, err error) {
	return 0, nil
}

// Write implements net.Conn interface
func (m *MockConnection) Write(b []byte) (n int, err error) {
	return 0, nil
}

// Close implements net.Conn interface
func (m *MockConnection) Close() error {
	return nil
}

// LocalAddr implements net.Conn interface
func (m *MockConnection) LocalAddr() net.Addr {
	return &MockAddr{addr: "127.0.0.1:0"}
}

// RemoteAddr implements net.Conn interface
func (m *MockConnection) RemoteAddr() net.Addr {
	return &MockAddr{addr: m.remoteAddr}
}

// SetDeadline implements net.Conn interface
func (m *MockConnection) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline implements net.Conn interface
func (m *MockConnection) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline implements net.Conn interface
func (m *MockConnection) SetWriteDeadline(t time.Time) error {
	return nil
}

// MockAddr represents a mock network address
type MockAddr struct {
	addr string
}

// Network implements net.Addr interface
func (m *MockAddr) Network() string {
	return "tcp"
}

// String implements net.Addr interface
func (m *MockAddr) String() string {
	return m.addr
}
