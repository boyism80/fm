package core

import (
	"fmt"
	"net"
	"sync"

	"github.com/boyism80/fm/core/client"
	"github.com/boyism80/fm/core/crypt"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
	"github.com/boyism80/fm/util"
)

type Client = client.Client

// BaseClient contains common fields and methods for all client types
type BaseClient struct {
	conn           net.Conn
	mu             sync.Mutex
	clientID       int
	fd             int
	sendEncryption *crypt.Encryption
	recvEncryption *crypt.Encryption
}

// GetConnection returns the underlying network connection
func (c *BaseClient) GetConnection() net.Conn {
	return c.conn
}

// GetSendEncryption returns the send encryption object
func (c *BaseClient) GetSendEncryption() *crypt.Encryption {
	return c.sendEncryption
}

// GetRecvEncryption returns the receive encryption object
func (c *BaseClient) GetRecvEncryption() *crypt.Encryption {
	return c.recvEncryption
}

// GetFd returns the file descriptor
func (c *BaseClient) GetFd() int {
	return c.fd
}

// Send sends a packet to the client with optional encryption
func (c *BaseClient) Send(packet types.Packet, policy types.SendPolicy) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Serialize packet with opcode
	writer := stream.NewStreamWriter(stream.LittleEndian)
	writer.WriteU16(uint16(packet.Opcode()))
	if err := packet.Serialize(writer); err != nil {
		return fmt.Errorf("failed to serialize packet: %w", err)
	}
	packetData := writer.Bytes()
	fmt.Println(util.ToHexString(packetData))

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

// GetFileDescriptor extracts the file descriptor from a net.Conn
func GetFileDescriptor(conn net.Conn) (int, error) {
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

// NewBaseClient creates a new BaseClient with the given parameters
func NewBaseClient(conn net.Conn, clientID int, fd int, sendEncryption, recvEncryption *crypt.Encryption) BaseClient {
	return BaseClient{
		conn:           conn,
		clientID:       clientID,
		fd:             fd,
		sendEncryption: sendEncryption,
		recvEncryption: recvEncryption,
	}
}
