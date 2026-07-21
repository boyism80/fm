package client

import (
	"net"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/crypt"
	"github.com/boyism80/fm/services/game/entity"
)

type GameClient struct {
	core.BaseClient
	character          *entity.Character
	transferDisconnect bool
	mu                 sync.Mutex
}

var _ core.Client = (*GameClient)(nil)

func NewGameClient(conn net.Conn, clientID int) (*GameClient, error) {
	fd, err := core.GetFileDescriptor(conn)
	if err != nil {
		return nil, err
	}
	ivSend := []byte{0x2F, 0xA3, 0x65, 0x43}
	ivRecv := []byte{0x65, 0x56, 0x12, 0xFD}
	se := crypt.NewEncryption(ivSend, -5)
	re := crypt.NewEncryption(ivRecv, 5)
	return &GameClient{
		BaseClient: core.NewBaseClient(conn, clientID, fd, &se, &re),
	}, nil
}

func (c *GameClient) SetCharacter(character *entity.Character) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.character = character
}

func (c *GameClient) GetCharacter() *entity.Character {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.character
}

// GetLogicActorPID is the packet mailbox for the map the character stands on;
// nil falls back to core.Server's nil LogicActor (login / between-maps).
func (c *GameClient) GetLogicActorPID() *actor.PID {
	c.mu.Lock()
	ch := c.character
	c.mu.Unlock()
	if ch == nil {
		return nil
	}
	m := ch.GetMap()
	if m == nil {
		return nil
	}
	return m.LogicActorPID()
}

func (c *GameClient) SetTransferDisconnect(v bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.transferDisconnect = v
}

func (c *GameClient) TakeTransferDisconnect() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	v := c.transferDisconnect
	c.transferDisconnect = false
	return v
}
