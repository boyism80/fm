package core

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/boyism80/fm/core/client"
	"github.com/boyism80/fm/core/crypt"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type Client = client.Client

type BaseClient struct {
	conn           net.Conn
	mu             sync.Mutex
	clientID       int
	fd             int
	sendEncryption *crypt.Encryption
	recvEncryption *crypt.Encryption
	lastPingAt     time.Time
	pongReceived   bool
}

func (c *BaseClient) GetConnection() net.Conn {
	return c.conn
}

func (c *BaseClient) GetSendEncryption() *crypt.Encryption {
	return c.sendEncryption
}

func (c *BaseClient) GetRecvEncryption() *crypt.Encryption {
	return c.recvEncryption
}

func (c *BaseClient) GetFd() int {
	return c.fd
}

func (c *BaseClient) Send(packet types.Packet, policy types.SendPolicy) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	writer := stream.NewStreamWriter(stream.LittleEndian)
	writer.WriteU16(uint16(packet.Opcode()))
	if err := packet.Serialize(writer); err != nil {
		return fmt.Errorf("failed to serialize packet: %w", err)
	}
	packetData := writer.Bytes()

	if policy == types.SEND_POLICY_RAW {
		_, err := c.conn.Write(packetData)
		return err
	}

	if policy&types.SEND_POLICY_ENCRYPT != 0 {

		writer = stream.NewStreamWriter(stream.LittleEndian)

		header := c.sendEncryption.GetPacketHeader(len(packetData))
		writer.Write(header)

		encryptedData := c.sendEncryption.Encrypt(packetData)
		writer.Write(encryptedData)

		_, err := c.conn.Write(writer.Bytes())
		return err
	}

	return fmt.Errorf("unsupported send policy: %d", policy)
}

func (c *BaseClient) MarkPongReceived() {
	c.pongReceived = true
}

func (c *BaseClient) NextPingAction(now time.Time, interval time.Duration) (sendPing bool, disconnect bool) {
	if c.lastPingAt.IsZero() {
		c.lastPingAt = now
		c.pongReceived = false
		return true, false
	}

	if now.Sub(c.lastPingAt) < interval {
		return false, false
	}

	if !c.pongReceived {
		return false, true
	}

	c.lastPingAt = now
	c.pongReceived = false
	return true, false
}

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

func NewBaseClient(conn net.Conn, clientID int, fd int, sendEncryption, recvEncryption *crypt.Encryption) BaseClient {
	return BaseClient{
		conn:           conn,
		clientID:       clientID,
		fd:             fd,
		sendEncryption: sendEncryption,
		recvEncryption: recvEncryption,
	}
}
