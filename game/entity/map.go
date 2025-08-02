package entity

import (
	"fmt"
	"time"

	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/data"
)

type MobSpawn struct {
	Spec          *data.MobSpawnSpec
	Spawned       bool
	LastSpawnedAt time.Time
}

type Map struct {
	ID              uint32
	objects         map[types.ObjectType]map[uint32]interface{} // Players, Mobs, Items, etc.
	controllerTable *ControllerTable
	MobSpawns       map[uint32]*MobSpawn
	listener        MapListener
}

func NewMap(id uint32, listener MapListener) *Map {
	if listener == nil {
		panic("MapListener cannot be nil")
	}

	return &Map{
		ID:              id,
		objects:         make(map[types.ObjectType]map[uint32]interface{}),
		controllerTable: NewControllerTable(nil),
		MobSpawns:       make(map[uint32]*MobSpawn),
		listener:        listener,
	}
}

func (m *Map) AddPlayer(playerID uint32, character *Character, init bool) error {
	if m.objects[types.OBJECT_TYPE_PLAYER] == nil {
		m.objects[types.OBJECT_TYPE_PLAYER] = make(map[uint32]interface{})
	}

	m.objects[types.OBJECT_TYPE_PLAYER][playerID] = character

	// Notify listener about player addition
	m.listener.OnPlayerAdded(m.ID, playerID, character, init)

	return nil
}

func (m *Map) RemovePlayer(playerID uint32) error {
	if m.objects[types.OBJECT_TYPE_PLAYER] == nil {
		return fmt.Errorf("no players on map")
	}

	if _, exists := m.objects[types.OBJECT_TYPE_PLAYER][playerID]; !exists {
		return fmt.Errorf("player %d not found on map", playerID)
	}

	delete(m.objects[types.OBJECT_TYPE_PLAYER], playerID)

	// Notify listener about player removal
	m.listener.OnPlayerRemoved(m.ID, playerID)

	return nil
}

func (m *Map) GetPlayer(playerID uint32) *Character {
	if m.objects[types.OBJECT_TYPE_PLAYER] == nil {
		return nil
	}

	if player, ok := m.objects[types.OBJECT_TYPE_PLAYER][playerID].(*Character); ok {
		return player
	}
	return nil
}

func (m *Map) GetPlayerCount() int {
	if m.objects[types.OBJECT_TYPE_PLAYER] == nil {
		return 0
	}
	return len(m.objects[types.OBJECT_TYPE_PLAYER])
}

// GetAllPlayers returns all players on the map
func (m *Map) GetAllPlayers() map[uint32]interface{} {
	if m.objects[types.OBJECT_TYPE_PLAYER] == nil {
		return make(map[uint32]interface{})
	}
	return m.objects[types.OBJECT_TYPE_PLAYER]
}

func (m *Map) GetControllerTable() *ControllerTable {
	return m.controllerTable
}

// BroadcastToPlayers sends a message to all players on the map
func (m *Map) BroadcastToPlayers(message types.Packet, policy types.SendPolicy, exceptPlayerID uint32) {
	if m.objects[types.OBJECT_TYPE_PLAYER] == nil {
		return
	}

	for playerID, player := range m.objects[types.OBJECT_TYPE_PLAYER] {
		if playerID == exceptPlayerID {
			continue
		}

		if character, ok := player.(*Character); ok {
			character.Send(message, policy)
		}
	}
}

// BroadcastToPlayer sends a message to a specific player on the map
func (m *Map) BroadcastToPlayer(playerID uint32, message types.Packet, policy types.SendPolicy) {
	if m.objects[types.OBJECT_TYPE_PLAYER] == nil {
		return
	}

	if player, ok := m.objects[types.OBJECT_TYPE_PLAYER][playerID]; ok {
		if character, ok := player.(*Character); ok {
			character.Send(message, policy)
		}
	}
}
