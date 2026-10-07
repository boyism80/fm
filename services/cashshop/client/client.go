package client

import (
	"net"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/crypt"
	"github.com/boyism80/fm/services/cashshop/entity"
)

type CashShopClient struct {
	core.BaseClient
	mu            sync.RWMutex
	logicActorPID *actor.PID
	character     *entity.Character
	disconnected  bool
	leaving       bool
}

var _ core.Client = (*CashShopClient)(nil)

func (c *CashShopClient) GetLogicActorPID() *actor.PID {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.logicActorPID
}

func (c *CashShopClient) SetLogicActorPID(pid *actor.PID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logicActorPID = pid
}

func (c *CashShopClient) GetCharacter() *entity.Character {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.character
}

func (c *CashShopClient) SetCharacter(character *entity.Character) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.disconnected {
		return false
	}
	c.character = character
	return true
}

func (c *CashShopClient) Disconnect() *entity.Character {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.disconnected = true
	if c.leaving {
		return nil
	}
	return c.character
}

func (c *CashShopClient) Leave() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.leaving = true
}

func (c *CashShopClient) Leaving() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.leaving
}

func NewCashShopClient(conn net.Conn, clientID int) (*CashShopClient, error) {
	ivSend := []byte{0x2F, 0xA3, 0x65, 0x43}
	ivRecv := []byte{0x65, 0x56, 0x12, 0xFD}
	se := crypt.NewEncryption(ivSend, -5)
	re := crypt.NewEncryption(ivRecv, 5)
	return &CashShopClient{
		BaseClient: core.NewBaseClient(conn, clientID, &se, &re),
	}, nil
}
