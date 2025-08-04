package client

import (
	"net"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/crypt"
)

// LoginClient represents a login server client with file descriptor based thread assignment
type LoginClient struct {
	core.BaseClient
}

// Ensure LoginClient implements core.Client
var _ core.Client = (*LoginClient)(nil)

// GetThreadHash returns file descriptor for thread assignment
func (c *LoginClient) GetThreadHash() int {
	return c.GetFd()
}

// NewLoginClient creates a new LoginClient with file descriptor extraction and encryption
func NewLoginClient(conn net.Conn, clientID int) (*LoginClient, error) {
	fd, err := core.GetFileDescriptor(conn)
	if err != nil {
		return nil, err
	}
	ivSend := []byte{0x2F, 0xA3, 0x65, 0x43}
	ivRecv := []byte{0x65, 0x56, 0x12, 0xFD}
	se := crypt.NewEncryption(ivSend, -5)
	re := crypt.NewEncryption(ivRecv, 5)
	return &LoginClient{
		BaseClient: core.NewBaseClient(conn, clientID, fd, &se, &re),
	}, nil
}
