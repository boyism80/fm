package core

import (
	"fmt"
	"net"
	"sync"

	"github.com/boyism80/fm/common/crypt"
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
)

// Client represents a connected game client with socket functionality
// Flow: Network connection -> IO thread processing -> Logic thread game state
// State: Connection, authentication status, game session, file descriptor, encryption
// Error Conditions: Connection lost, invalid packets, authentication failure
type Client[T any] struct {
	conn           net.Conn
	mu             sync.Mutex
	clientID       int              // Unique client identifier for thread assignment
	fd             int              // File descriptor for thread assignment (extracted from connection)
	sendEncryption crypt.Encryption // Encryption for outgoing packets
	recvEncryption crypt.Encryption // Encryption for incoming packets
	data           T                // Generic data for server-specific information
}

// Ensure Client implements ThreadAssignable
var _ ThreadAssignable = (*Client[any])(nil)

// GetThreadHash returns a hash value for thread assignment based on file descriptor
// Flow: File descriptor -> Hash generation -> Thread assignment
// Thread Assignment: Uses file descriptor for consistent thread assignment
// Error Handling: Returns valid hash value
func (c *Client[T]) GetThreadHash() int {
	return c.fd
}

// GetConnection returns the underlying network connection
func (c *Client[T]) GetConnection() net.Conn {
	return c.conn
}

// GetFileDescriptor returns the file descriptor of this client
func (c *Client[T]) GetFileDescriptor() int {
	return c.fd
}

// NewClient creates a new Client with file descriptor extraction and encryption
func NewClient[T any](conn net.Conn, clientID int, data T) (*Client[T], error) {
	fd, err := getFileDescriptor(conn)
	if err != nil {
		return nil, err
	}

	// Initialize encryption with default IVs
	ivSend := []byte{0x2F, 0xA3, 0x65, 0x43}
	ivRecv := []byte{0x65, 0x56, 0x12, 0xFD}

	return &Client[T]{
		conn:           conn,
		clientID:       clientID,
		fd:             fd,
		sendEncryption: crypt.NewEncryption(ivSend, -5),
		recvEncryption: crypt.NewEncryption(ivRecv, 5),
		data:           data,
	}, nil
}

// getFileDescriptor extracts the file descriptor from a net.Conn
func getFileDescriptor(conn net.Conn) (int, error) {
	remoteAddr := conn.RemoteAddr().String()
	fd := 0
	for _, char := range remoteAddr {
		fd = fd*31 + int(char)
	}
	if fd < 0 {
		fd = -fd
	}
	return fd, nil
}

// Send sends a packet to the client with optional encryption
func (c *Client[T]) Send(packet types.Packet, policy types.SendPolicy) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Serialize packet with opcode
	writer := stream.NewStreamWriter(stream.LittleEndian)
	writer.WriteU16(uint16(packet.Opcode()))
	if err := packet.Serialize(writer); err != nil {
		return fmt.Errorf("failed to serialize packet: %w", err)
	}
	packetData := writer.Bytes()

	// Handle raw policy (no encryption)
	if policy == types.SEND_POLICY_RAW {
		_, err := c.conn.Write(packetData)
		return err
	}

	// Handle encrypted policy
	if policy&types.SEND_POLICY_ENCRYPT != 0 {
		// Create new writer for encrypted packet
		writer = stream.NewStreamWriter(stream.LittleEndian)

		// Add packet header
		header := c.sendEncryption.GetPacketHeader(len(packetData))
		writer.Write(header)

		// Encrypt packet data
		encryptedData := c.sendEncryption.Encrypt(packetData)
		writer.Write(encryptedData)

		// Send encrypted packet
		_, err := c.conn.Write(writer.Bytes())
		return err
	}

	return fmt.Errorf("unsupported send policy: %d", policy)
}

// GetSendEncryption returns the send encryption object
func (c *Client[T]) GetSendEncryption() crypt.Encryption {
	return c.sendEncryption
}

// GetRecvEncryption returns the receive encryption object
func (c *Client[T]) GetRecvEncryption() crypt.Encryption {
	return c.recvEncryption
}

// GetData returns the generic data for this client
func (c *Client[T]) GetData() T {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.data
}

// SetData sets the generic data for this client
func (c *Client[T]) SetData(data T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = data
}
