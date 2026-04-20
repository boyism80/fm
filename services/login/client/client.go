package client

import (
	"net"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/crypt"
)

type LoginClient struct {
	core.BaseClient
	logicActorPID      *actor.PID
	pidMutex           sync.RWMutex
	accountId          uint32
	worldId            uint32
	channelId          uint8
	accountMutex       sync.RWMutex
	transferDisconnect bool
	transferMu         sync.Mutex
}

var _ core.Client = (*LoginClient)(nil)

func (c *LoginClient) GetLogicActorPID() *actor.PID {
	c.pidMutex.RLock()
	defer c.pidMutex.RUnlock()
	return c.logicActorPID
}

func (c *LoginClient) SetLogicActorPID(pid *actor.PID) {
	c.pidMutex.Lock()
	defer c.pidMutex.Unlock()
	c.logicActorPID = pid
}

func (c *LoginClient) GetAccountId() uint32 {
	c.accountMutex.RLock()
	defer c.accountMutex.RUnlock()
	return c.accountId
}

func (c *LoginClient) SetAccountId(id uint32) {
	c.accountMutex.Lock()
	defer c.accountMutex.Unlock()
	c.accountId = id
}

func (c *LoginClient) GetWorldId() uint32 {
	c.accountMutex.RLock()
	defer c.accountMutex.RUnlock()
	return c.worldId
}

func (c *LoginClient) SetWorldId(id uint32) {
	c.accountMutex.Lock()
	defer c.accountMutex.Unlock()
	c.worldId = id
}

func (c *LoginClient) GetChannelId() uint8 {
	c.accountMutex.RLock()
	defer c.accountMutex.RUnlock()
	return c.channelId
}

func (c *LoginClient) SetChannelId(id uint8) {
	c.accountMutex.Lock()
	defer c.accountMutex.Unlock()
	c.channelId = id
}

func (c *LoginClient) SetTransferDisconnect(v bool) {
	c.transferMu.Lock()
	defer c.transferMu.Unlock()
	c.transferDisconnect = v
}

func (c *LoginClient) TakeTransferDisconnect() bool {
	c.transferMu.Lock()
	defer c.transferMu.Unlock()
	v := c.transferDisconnect
	c.transferDisconnect = false
	return v
}

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
