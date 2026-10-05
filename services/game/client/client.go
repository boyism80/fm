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
	character       *entity.Character
	changingChannel bool
	loggedOut       bool
	sessionLost     bool
	mu              sync.Mutex
}

var _ core.Client = (*GameClient)(nil)

func NewGameClient(conn net.Conn, clientID int) (*GameClient, error) {
	ivSend := []byte{0x2F, 0xA3, 0x65, 0x43}
	ivRecv := []byte{0x65, 0x56, 0x12, 0xFD}
	se := crypt.NewEncryption(ivSend, -5)
	re := crypt.NewEncryption(ivRecv, 5)
	return &GameClient{
		BaseClient: core.NewBaseClient(conn, clientID, &se, &re),
	}, nil
}

// SetCharacter fails once the client has logged out, so a login that finishes after the disconnect never enters the game.
func (c *GameClient) SetCharacter(character *entity.Character) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.loggedOut {
		return false
	}
	c.character = character
	return true
}

// Logout takes the character off the client; later packets and the disconnect find no character.
func (c *GameClient) Logout() *entity.Character {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.loggedOut = true
	character := c.character
	c.character = nil
	return character
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
		m = ch.Destination
	}
	if m == nil {
		return nil
	}
	return m.LogicActorPID()
}

func (c *GameClient) SetChangingChannel(v bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.changingChannel = v
}

func (c *GameClient) ChangingChannel() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.changingChannel
}

func (c *GameClient) LoseSession() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sessionLost = true
}

func (c *GameClient) SessionLost() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sessionLost
}
