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

const (
	sendQueueSize = 4096
	sendTimeout   = 10 * time.Second
)

type BaseClient struct {
	conn           net.Conn
	mu             sync.Mutex
	out            chan []byte
	done           chan struct{}
	closeOnce      sync.Once
	clientID       int
	sendEncryption *crypt.Encryption
	recvEncryption *crypt.Encryption
	lastPingAt     time.Time
	pongReceived   bool
}

func (c *BaseClient) GetConnection() net.Conn {
	return c.conn
}

func (c *BaseClient) GetRemoteIP() string {
	addr := c.conn.RemoteAddr().String()
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}

func (c *BaseClient) GetSendEncryption() *crypt.Encryption {
	return c.sendEncryption
}

func (c *BaseClient) GetRecvEncryption() *crypt.Encryption {
	return c.recvEncryption
}

func (c *BaseClient) GetClientID() int {
	return c.clientID
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

	switch {
	case policy == types.SEND_POLICY_RAW:
	case policy&types.SEND_POLICY_ENCRYPT != 0:
		writer = stream.NewStreamWriter(stream.LittleEndian)

		header := c.sendEncryption.GetPacketHeader(len(packetData))
		writer.Write(header)

		encryptedData := c.sendEncryption.Encrypt(packetData)
		writer.Write(encryptedData)

		packetData = writer.Bytes()
	default:
		return fmt.Errorf("unsupported send policy: %d", policy)
	}

	// Encrypted under c.mu and queued in the same order, so the client sees the IV sequence it expects.
	select {
	case <-c.done:
		return net.ErrClosed
	default:
	}
	select {
	case c.out <- packetData:
		return nil
	default:
		_ = c.conn.Close()
		return fmt.Errorf("send queue full")
	}
}

// Close stops the writer; the reader sees the closed connection and ends the session.
func (c *BaseClient) Close() {
	c.closeOnce.Do(func() {
		close(c.done)
	})
	_ = c.conn.Close()
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

func NewBaseClient(conn net.Conn, clientID int, sendEncryption, recvEncryption *crypt.Encryption) BaseClient {
	out := make(chan []byte, sendQueueSize)
	done := make(chan struct{})
	// A client that stops reading must not block the actor that sends to it; the writer owns the socket write and its deadline.
	go func() {
		for {
			select {
			case data := <-out:
				_ = conn.SetWriteDeadline(time.Now().Add(sendTimeout))
				if _, err := conn.Write(data); err != nil {
					_ = conn.Close()
					return
				}
			case <-done:
				return
			}
		}
	}()

	return BaseClient{
		conn:           conn,
		out:            out,
		done:           done,
		clientID:       clientID,
		sendEncryption: sendEncryption,
		recvEncryption: recvEncryption,
	}
}
