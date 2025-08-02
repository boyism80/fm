package entity

import (
	"fmt"
	"math"
	"time"

	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/constant"
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
	spec            *data.MapSpec // Map specification data
	sequence        uint32        // Sequence ID for generating unique object IDs
}

func NewMap(id uint32, listener MapListener, spec *data.MapSpec) *Map {
	if listener == nil {
		panic("MapListener cannot be nil")
	}
	if spec == nil {
		panic("MapSpec cannot be nil")
	}

	mapInstance := &Map{
		ID:              id,
		objects:         make(map[types.ObjectType]map[uint32]interface{}),
		controllerTable: NewControllerTable(nil),
		MobSpawns:       make(map[uint32]*MobSpawn),
		listener:        listener,
		spec:            spec,
		sequence:        0,
	}

	// Initialize NPCs from MapSpec (following old server pattern)
	mapInstance.initializeNpcs()

	return mapInstance
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

// GetSpec returns the map specification data
func (m *Map) GetSpec() *data.MapSpec {
	return m.spec
}

// FootholdPoint calculates the foothold position for a given point
func (m *Map) FootholdPoint(point types.Point[int16]) *types.Point[int16] {
	return m.spec.FootholdPoint(point)
}

// initializeNpcs initializes NPCs from MapSpec (following old server pattern)
func (m *Map) initializeNpcs() {
	if m.objects[types.OBJECT_TYPE_NPC] == nil {
		m.objects[types.OBJECT_TYPE_NPC] = make(map[uint32]interface{})
	}

	// Create NPCs from MapSpec (following old server pattern)
	for _, npcSpec := range m.spec.NpcSpawns {
		m.sequence++ // Generate unique OID for NPC
		npc := &Npc{
			Object: Object{
				OID: m.sequence, // Use sequence as OID
			},
			Spec: &npcSpec,
		}
		m.objects[types.OBJECT_TYPE_NPC][npc.OID] = npc
	}
}

// GetNpcs returns all NPCs on the map
func (m *Map) GetNpcs() map[uint32]interface{} {
	if m.objects[types.OBJECT_TYPE_NPC] == nil {
		return make(map[uint32]interface{})
	}
	return m.objects[types.OBJECT_TYPE_NPC]
}

// SpawnMob spawns a mob on the map (following old server pattern)
func (m *Map) SpawnMob(mobID uint32, foothold int16, position types.Point[int16]) (*Mob, error) {
	// Generate unique sequence ID for mob
	m.sequence++

	// Calculate spawn point
	spawnPoint, ok := m.spec.DropPoint(position)
	if !ok {
		spawnPoint = position
	}

	// Create mob entity
	mob := &Mob{
		Life: Life{
			Object: Object{
				OID:      m.sequence, // Use sequence as OID
				Position: spawnPoint,
			},
			Hp:     0, // Will be set by mob spec
			Mp:     0, // Will be set by mob spec
			MaxHp:  0, // Will be set by mob spec
			MaxMp:  0, // Will be set by mob spec
			Stance: 5,
		},
		Foothold: foothold,
	}

	// Initialize objects map for monsters if needed
	if m.objects[types.OBJECT_TYPE_MONSTER] == nil {
		m.objects[types.OBJECT_TYPE_MONSTER] = make(map[uint32]interface{})
	}

	// Add mob to map objects
	m.objects[types.OBJECT_TYPE_MONSTER][m.sequence] = mob

	// Notify listener about mob spawn
	m.listener.OnMobSpawned(m.ID, m.sequence, mob)

	return mob, nil
}

// RemoveMob removes a mob from the map
func (m *Map) RemoveMob(mobID uint32) error {
	if m.objects[types.OBJECT_TYPE_MONSTER] == nil {
		return fmt.Errorf("no monsters on map")
	}

	if _, exists := m.objects[types.OBJECT_TYPE_MONSTER][mobID]; !exists {
		return fmt.Errorf("mob %d not found on map", mobID)
	}

	delete(m.objects[types.OBJECT_TYPE_MONSTER], mobID)

	// Notify listener about mob removal
	m.listener.OnMobRemoved(m.ID, mobID)

	return nil
}

// GetMob retrieves a mob from the map
func (m *Map) GetMob(mobID uint32) *Mob {
	if m.objects[types.OBJECT_TYPE_MONSTER] == nil {
		return nil
	}

	if mob, ok := m.objects[types.OBJECT_TYPE_MONSTER][mobID].(*Mob); ok {
		return mob
	}
	return nil
}

// GetMobs returns all mobs on the map
func (m *Map) GetMobs() map[uint32]interface{} {
	if m.objects[types.OBJECT_TYPE_MONSTER] == nil {
		return make(map[uint32]interface{})
	}
	return m.objects[types.OBJECT_TYPE_MONSTER]
}

// BroadcastToPlayers sends a message to all players on the map (excluding specified player)
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

// BroadcastToAllPlayers sends a message to all players on the map (including sender)
func (m *Map) BroadcastToAllPlayers(message types.Packet, policy types.SendPolicy) {
	if m.objects[types.OBJECT_TYPE_PLAYER] == nil {
		return
	}

	for _, player := range m.objects[types.OBJECT_TYPE_PLAYER] {
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

	if player, ok := m.objects[types.OBJECT_TYPE_PLAYER][playerID].(*Character); ok {
		player.Send(message, policy)
	}
}

// SpawnItem spawns an item on the map (following old server pattern)
func (m *Map) SpawnItem(item Item, ownerID uint32, dropType constant.DropType) error {
	// Generate unique sequence ID
	m.sequence++

	// Get drop from item
	drop := item.GetDrop()
	if drop == nil {
		return fmt.Errorf("item has no drop information")
	}

	// Calculate drop point (following old server pattern)
	dropPoint, ok := m.spec.DropPoint(drop.Position)
	if !ok {
		dropPoint = drop.SpawnedPoint
	}

	// Update drop information
	drop.Position = dropPoint
	drop.OID = m.sequence
	drop.Owner = ownerID
	drop.DropType = dropType

	// Initialize objects map for items if needed
	if m.objects[types.OBJECT_TYPE_ITEM] == nil {
		m.objects[types.OBJECT_TYPE_ITEM] = make(map[uint32]interface{})
	}

	// Add item to map objects
	m.objects[types.OBJECT_TYPE_ITEM][m.sequence] = item

	// Notify listener about item spawn
	m.listener.OnItemSpawned(m.ID, m.sequence, item, drop)

	return nil
}

// SpawnMeso spawns meso on the map (following old server pattern)
func (m *Map) SpawnMeso(count int32, position types.Point[int16], ownerID uint32, dropType constant.DropType) error {
	// Generate unique sequence ID
	m.sequence++

	// Calculate drop point
	dropPoint, ok := m.spec.DropPoint(position)
	if !ok {
		dropPoint = position
	}

	// Create meso entity
	meso := &Meso{
		Drop: &Drop{
			Object: &Object{
				OID:      m.sequence,
				Position: dropPoint,
			},
			SpawnedPoint: position,
			DropType:     dropType,
			Owner:        ownerID,
		},
		Count: count,
	}

	// Initialize objects map for items if needed
	if m.objects[types.OBJECT_TYPE_ITEM] == nil {
		m.objects[types.OBJECT_TYPE_ITEM] = make(map[uint32]interface{})
	}

	// Add meso to map objects
	m.objects[types.OBJECT_TYPE_ITEM][m.sequence] = meso

	// Notify listener about meso spawn
	m.listener.OnMesoSpawned(m.ID, m.sequence, meso)

	return nil
}

// RemoveItem removes an item from the map
func (m *Map) RemoveItem(itemID uint32) error {
	if m.objects[types.OBJECT_TYPE_ITEM] == nil {
		return fmt.Errorf("no items on map")
	}

	if _, exists := m.objects[types.OBJECT_TYPE_ITEM][itemID]; !exists {
		return fmt.Errorf("item %d not found on map", itemID)
	}

	delete(m.objects[types.OBJECT_TYPE_ITEM], itemID)

	// Notify listener about item removal
	m.listener.OnItemRemoved(m.ID, itemID, 0, REMOVE_ITEM_TYPE_EXPIRED)

	return nil
}

// GetItem retrieves an item from the map
func (m *Map) GetItem(itemID uint32) Item {
	if m.objects[types.OBJECT_TYPE_ITEM] == nil {
		return nil
	}

	if item, ok := m.objects[types.OBJECT_TYPE_ITEM][itemID].(Item); ok {
		return item
	}
	return nil
}

// GetItems returns all items on the map
func (m *Map) GetItems() map[uint32]interface{} {
	if m.objects[types.OBJECT_TYPE_ITEM] == nil {
		return make(map[uint32]interface{})
	}
	return m.objects[types.OBJECT_TYPE_ITEM]
}

// LootItem attempts to loot an item from the map and returns the looted item and success status
// All capacity checks are performed before removing the item from the map
func (m *Map) LootItem(itemID uint32, character *Character, position types.Point[int16]) (interface{}, bool) {
	if m.objects[types.OBJECT_TYPE_ITEM] == nil {
		return nil, false
	}

	itemInterface, exists := m.objects[types.OBJECT_TYPE_ITEM][itemID]
	if !exists {
		return nil, false
	}

	switch item := itemInterface.(type) {
	case Item:
		drop := item.GetDrop()
		if drop == nil {
			return nil, false
		}

		// Check ownership for owned drops
		if drop.DropType == constant.DROP_TYPE_OWNED && drop.Owner != character.ID {
			return nil, false
		}

		// Check inventory capacity before removing from map
		invenType := item.GetInventoryType()
		inven := character.Inventory[invenType]
		spec := item.GetSpec()

		if !inven.IsFree(spec, item.GetCount()) {
			// Inventory is full, don't remove from map
			return nil, false
		}

		// Remove item from map immediately
		delete(m.objects[types.OBJECT_TYPE_ITEM], itemID)

		// Notify listener about item removal
		m.listener.OnItemRemoved(m.ID, itemID, character.ID, REMOVE_ITEM_TYPE_ANIMATED)

		return item, true

	case *Meso:
		drop := item.GetDrop()
		if drop == nil {
			return nil, false
		}

		// Check ownership for owned drops
		if drop.DropType == constant.DROP_TYPE_OWNED && drop.Owner != character.ID {
			return nil, false
		}

		// Check meso capacity before removing from map
		mesoCount := item.GetCount32()
		cap := math.MaxInt32 - character.Meso
		if int32(mesoCount) > cap {
			// Capacity exceeded, don't remove from map
			return nil, false
		}

		// Remove meso from map immediately
		delete(m.objects[types.OBJECT_TYPE_ITEM], itemID)

		// Notify listener about item removal
		m.listener.OnItemRemoved(m.ID, itemID, character.ID, REMOVE_ITEM_TYPE_ANIMATED)

		// Return the meso object
		return item, true

	default:
		return nil, false
	}
}
