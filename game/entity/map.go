package entity

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/types"
)

// MapListener defines interface for map events
type MapListener interface {
	OnPlayerAdded(mapID uint32, playerID uint32, character *Character, init bool)
	OnPlayerRemoved(mapID uint32, playerID uint32)
	OnPlayerMoved(mapID uint32, playerID uint32, character *Character)
	OnPlayerMove(mapID uint32, playerID uint32, character *Character, startPoint types.Vector2[int16], fragments []dto.MoveFragment)
	OnPlayerChat(mapID uint32, playerID uint32, message string)
	OnItemSpawned(mapID uint32, itemID uint32, item Item, drop *Drop)
	OnMesoSpawned(mapID uint32, itemID uint32, meso *Meso)
	OnItemRemoved(mapID uint32, itemID uint32, characterID uint32, mode constant.RemoveItemType)
	OnMobSpawned(mapID uint32, mobID uint32, mob *Mob)
	OnMobRemoved(mapID uint32, mobID uint32, animationType constant.MobDieAnimationType)
	OnMobControllerChange(mob *Mob, before *Character, after *Character)
	OnMobMoved(mapID uint32, mobID uint32, isAggroed bool, centerSplit int8, skill1 uint8, skill2 uint8, skill3 uint8, skill4 uint8, startPoint types.Vector2[int16], movements []dto.MoveFragment)
	OnAttack(mapID uint32, characterID uint32, attackInfo dto.AttackInfo, skillLevel uint8)
	OnWarpCharacter(character *Character, targetMapID uint32, portal uint8)
}

type MobSpawn struct {
	Wz            *wz.MobSpawn
	Spawned       bool
	LastSpawnedAt time.Time
}

type Map struct {
	ID              uint32
	objects         map[types.ObjectType]map[uint32]interface{} // Players, Mobs, Items, etc.
	controllerTable *ControllerTable
	MobSpawns       map[uint32]*MobSpawn
	listener        MapListener
	model           *wz.Map     // Map specification data
	sequence        uint32      // Sequence ID for generating unique object IDs
	availableOIDs   []uint32    // Queue of available OIDs for reuse
	context         GameContext // GameContext for accessing resources
	actorPID        *actor.PID
	pidMutex        sync.RWMutex
}

func NewMap(id uint32, listener MapListener, mapId uint32, context GameContext) *Map {
	if listener == nil {
		panic("MapListener cannot be nil")
	}
	if context == nil {
		panic("GameContext cannot be nil")
	}

	// Get map model from resources
	mapSpec, ok := context.GetResources().Maps[mapId]
	if !ok {
		panic(fmt.Sprintf("MapSpec not found for ID: %d", mapId))
	}

	mapInstance := &Map{
		ID:              id,
		objects:         make(map[types.ObjectType]map[uint32]interface{}),
		controllerTable: nil, // Will be set after mapInstance is created
		MobSpawns:       make(map[uint32]*MobSpawn),
		listener:        listener,
		model:           mapSpec,
		sequence:        0,
		availableOIDs:   make([]uint32, 0),
		context:         context,
	}

	// Set up controller table with Map-specific callback
	mapInstance.controllerTable = NewControllerTable(mapInstance.onMobControllerChange)

	// Initialize NPCs from MapSpec (following old server pattern)
	mapInstance.initializeNpcs()

	// Initialize mobs from MapSpec (following old server pattern)
	mapInstance.initializeMobs()

	return mapInstance
}

// onMobControllerChange is the Map-specific callback for mob controller changes
func (m *Map) onMobControllerChange(mob *Mob, before *Character, after *Character) {
	// Delegate to MapListener to handle the mob controller change
	m.listener.OnMobControllerChange(mob, before, after)
}

// allocateOID allocates a new OID, reusing from queue if available
func (m *Map) allocateOID() uint32 {
	if len(m.availableOIDs) > 0 {
		// Reuse OID from queue
		oid := m.availableOIDs[0]
		m.availableOIDs = m.availableOIDs[1:]
		return oid
	}
	// Generate new OID
	m.sequence++
	return m.sequence
}

// releaseOID releases an OID back to the queue for reuse
func (m *Map) releaseOID(oid uint32) {
	m.availableOIDs = append(m.availableOIDs, oid)
}

// OnMobControllerChange handles mob controller changes (following old server pattern)
// This is a global callback that will be replaced with Map-specific callback
func OnMobControllerChange(mob *Mob, before *Character, after *Character) {
	// This callback is called when a mob's controller changes
	// In the new Entity-based architecture, we'll handle this through the MapListener
	// The actual mob AI logic will be implemented in the MapListener
}

func (m *Map) AddPlayer(playerID uint32, character *Character, init bool) error {
	if m.objects[types.OBJECT_TYPE_PLAYER] == nil {
		m.objects[types.OBJECT_TYPE_PLAYER] = make(map[uint32]interface{})
	}

	m.objects[types.OBJECT_TYPE_PLAYER][playerID] = character

	// Notify listener about player addition first (sends Warp and SpawnMob packets)
	// Controller assignment will be done after SpawnMob packets are sent
	m.listener.OnPlayerAdded(m.ID, playerID, character, init)

	// Add player to controller table for mob AI (following old server pattern)
	// This will trigger StartControlMob packets after SpawnMob packets
	m.controllerTable.EnterPlayer(character)

	return nil
}

func (m *Map) RemovePlayer(playerID uint32) error {
	if m.objects[types.OBJECT_TYPE_PLAYER] == nil {
		return fmt.Errorf("no players on map")
	}

	if _, exists := m.objects[types.OBJECT_TYPE_PLAYER][playerID]; !exists {
		return fmt.Errorf("player %d not found on map", playerID)
	}

	character := m.objects[types.OBJECT_TYPE_PLAYER][playerID].(*Character)
	delete(m.objects[types.OBJECT_TYPE_PLAYER], playerID)

	// Remove player from controller table for mob AI (following old server pattern)
	m.controllerTable.LeavePlayer(character)

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
func (m *Map) GetSpec() *wz.Map {
	return m.model
}

// FootholdPoint calculates the foothold position for a given point
func (m *Map) FootholdPoint(point types.Point[int16]) *types.Point[int16] {
	return m.model.FootholdPoint(point)
}

// initializeNpcs initializes NPCs from MapSpec (following old server pattern)
func (m *Map) initializeNpcs() {
	if m.objects[types.OBJECT_TYPE_NPC] == nil {
		m.objects[types.OBJECT_TYPE_NPC] = make(map[uint32]interface{})
	}

	// Create NPCs from MapSpec (following old server pattern)
	for _, npcSpec := range m.model.NpcSpawns {
		oid := m.allocateOID()
		npc := &Npc{
			Object: Object{
				OID: oid,
			},
			Wz: &npcSpec,
		}
		m.objects[types.OBJECT_TYPE_NPC][npc.OID] = npc
	}
}

// initializeMobs initializes mobs from MapSpec MobSpawns (following old server pattern)
func (m *Map) initializeMobs() {
	// Initialize MobSpawns map from MapSpec
	for spawnId, mobSpawnSpec := range m.model.MobSpawns {
		// Create MobSpawn entry
		m.MobSpawns[spawnId] = &MobSpawn{
			Wz:            &mobSpawnSpec,
			Spawned:       false,
			LastSpawnedAt: time.Time{},
		}

		// Spawn the mob immediately
		position := types.Point[int16]{
			X: mobSpawnSpec.Position.X,
			Y: mobSpawnSpec.Position.Y,
		}

		_, err := m.SpawnMob(mobSpawnSpec.ID, position)
		if err != nil {
			// Log error but continue with other mobs
			fmt.Printf("Failed to spawn mob %d at spawn point %d: %v\n", mobSpawnSpec.ID, spawnId, err)
			continue
		}

		// Mark as spawned
		m.MobSpawns[spawnId].Spawned = true
		m.MobSpawns[spawnId].LastSpawnedAt = time.Now()
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
func (m *Map) SpawnMob(mobId uint32, position types.Point[int16]) (*Mob, error) {
	// Allocate OID for mob
	oid := m.allocateOID()

	// Get mob model from resources
	mobSpec, ok := m.context.GetResources().Monsters[mobId]
	if !ok {
		return nil, fmt.Errorf("mob model not found for ID: %d", mobId)
	}

	// Find foothold for the position (following old server pattern)
	foothold, ok := m.model.Footholds.Find(position)
	if !ok {
		return nil, fmt.Errorf("no valid foothold found at position: %v", position)
	}

	// Calculate spawn point
	spawnPoint, ok := m.model.DropPoint(position)
	if !ok {
		spawnPoint = position
	}

	// Create mob entity
	mob := &Mob{
		Life: Life{
			Object: Object{
				OID:      oid,
				Position: spawnPoint,
				Context:  m.context, // Pass GameContext to Object
			},
			Hp:     uint16(mobSpec.MaxHP), // Set from mob model
			Mp:     uint16(mobSpec.MaxMP), // Set from mob model
			MaxHp:  uint16(mobSpec.MaxHP), // Set from mob model
			MaxMp:  uint16(mobSpec.MaxMP), // Set from mob model
			Stance: 5,
		},
		Foothold: foothold.ID,
		Wz:       mobSpec, // Pass MobSpec directly
		MapID:    m.ID,    // Set the map ID where this mob is spawned
	}

	// Initialize objects map for monsters if needed
	if m.objects[types.OBJECT_TYPE_MONSTER] == nil {
		m.objects[types.OBJECT_TYPE_MONSTER] = make(map[uint32]interface{})
	}

	// Add mob to map objects
	m.objects[types.OBJECT_TYPE_MONSTER][oid] = mob

	// Notify listener about mob spawn
	m.listener.OnMobSpawned(m.ID, oid, mob)

	// Add mob to controller table for mob AI (following old server pattern)
	m.controllerTable.EnterMob(mob)

	return mob, nil
}

// RemoveMob removes a mob from the map
func (m *Map) RemoveMob(mobID uint32, animationType constant.MobDieAnimationType) error {
	if m.objects[types.OBJECT_TYPE_MONSTER] == nil {
		return fmt.Errorf("no monsters on map")
	}

	if _, exists := m.objects[types.OBJECT_TYPE_MONSTER][mobID]; !exists {
		return fmt.Errorf("mob %d not found on map", mobID)
	}

	mob := m.objects[types.OBJECT_TYPE_MONSTER][mobID].(*Mob)
	delete(m.objects[types.OBJECT_TYPE_MONSTER], mobID)

	// Release OID for reuse
	m.releaseOID(mobID)

	// Remove mob from controller table for mob AI (following old server pattern)
	m.controllerTable.LeaveMob(mob)

	// Notify listener about mob removal
	m.listener.OnMobRemoved(m.ID, mobID, animationType)

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
	// Allocate OID for item
	oid := m.allocateOID()

	// Get drop from item
	drop := item.GetDrop()
	if drop == nil {
		return fmt.Errorf("item has no drop information")
	}

	// Calculate drop point (following old server pattern)
	dropPoint, ok := m.model.DropPoint(drop.Position)
	if !ok {
		dropPoint = drop.SpawnedPoint
	}

	// Update drop information
	drop.Position = dropPoint
	drop.OID = oid
	drop.Owner = ownerID
	drop.DropType = dropType
	drop.MapID = m.ID
	drop.setupDropTimers()

	// Initialize objects map for items if needed
	if m.objects[types.OBJECT_TYPE_ITEM] == nil {
		m.objects[types.OBJECT_TYPE_ITEM] = make(map[uint32]interface{})
	}

	// Add item to map objects
	m.objects[types.OBJECT_TYPE_ITEM][oid] = item

	// Notify listener about item spawn
	m.listener.OnItemSpawned(m.ID, oid, item, drop)

	return nil
}

// SpawnMeso spawns meso on the map (following old server pattern)
func (m *Map) SpawnMeso(count int32, position types.Point[int16], ownerID uint32, dropType constant.DropType) error {
	// Allocate OID for meso
	oid := m.allocateOID()

	// Calculate drop point
	dropPoint, ok := m.model.DropPoint(position)
	if !ok {
		dropPoint = position
	}

	// Create meso entity
	meso := NewMeso(count, dropPoint, ownerID, dropType, oid, m.context, m.ID)

	// Initialize objects map for items if needed
	if m.objects[types.OBJECT_TYPE_ITEM] == nil {
		m.objects[types.OBJECT_TYPE_ITEM] = make(map[uint32]interface{})
	}

	// Add meso to map objects
	m.objects[types.OBJECT_TYPE_ITEM][oid] = meso

	// Notify listener about meso spawn
	m.listener.OnMesoSpawned(m.ID, oid, meso)

	return nil
}

// RemoveItem removes an item from the map
func (m *Map) RemoveItem(itemID uint32, removeType constant.RemoveItemType, playerID uint32) error {
	if m.objects[types.OBJECT_TYPE_ITEM] == nil {
		return fmt.Errorf("no items on map")
	}

	if _, exists := m.objects[types.OBJECT_TYPE_ITEM][itemID]; !exists {
		return fmt.Errorf("item %d not found on map", itemID)
	}

	// Get the item before removing it to cancel its timers
	itemInterface := m.objects[types.OBJECT_TYPE_ITEM][itemID]

	// Cancel timers if it's a dropable item
	if dropable, ok := itemInterface.(Dropable); ok {
		if drop := dropable.GetDrop(); drop != nil {
			drop.cancelTimers()
		}
	}

	delete(m.objects[types.OBJECT_TYPE_ITEM], itemID)

	// Release OID for reuse
	m.releaseOID(itemID)

	// Notify listener about item removal
	m.listener.OnItemRemoved(m.ID, itemID, playerID, removeType)

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

// LootItem attempts to loot an item from the map and returns the looted item and reason
// All capacity checks are performed before removing the item from the map
func (m *Map) LootItem(itemID uint32, character *Character, position types.Point[int16]) (interface{}, constant.LootResult) {
	if m.objects[types.OBJECT_TYPE_ITEM] == nil {
		return nil, constant.LOOT_FAILED_ITEM_NOT_FOUND
	}

	itemInterface, exists := m.objects[types.OBJECT_TYPE_ITEM][itemID]
	if !exists {
		return nil, constant.LOOT_FAILED_ITEM_NOT_FOUND
	}

	switch item := itemInterface.(type) {
	case Item:
		drop := item.GetDrop()
		if drop == nil {
			return nil, constant.LOOT_FAILED_INVALID_ITEM
		}

		// Check ownership for owned drops
		if drop.DropType == constant.DROP_TYPE_OWNED && drop.Owner != character.ID {
			return nil, constant.LOOT_FAILED_NO_OWNERSHIP
		}

		// Check inventory capacity before removing from map
		invenType := item.GetInventoryType()
		inven := character.Inventory[invenType]
		model := item.GetModel()

		if !inven.IsFree(model, item.GetCount()) {
			return nil, constant.LOOT_FAILED_INVENTORY_FULL
		}

		// Remove item from map using RemoveItem
		if err := m.RemoveItem(itemID, constant.REMOVE_ITEM_TYPE_ANIMATED, character.ID); err != nil {
			return nil, constant.LOOT_FAILED_INVALID_ITEM
		}

		return item, constant.LOOT_SUCCESS

	case *Meso:
		drop := item.GetDrop()
		if drop == nil {
			return nil, constant.LOOT_FAILED_INVALID_ITEM
		}

		// Check ownership for owned drops
		if drop.DropType == constant.DROP_TYPE_OWNED && drop.Owner != character.ID {
			return nil, constant.LOOT_FAILED_NO_OWNERSHIP
		}

		// Check meso capacity before removing from map
		mesoCount := item.GetCount32()
		cap := math.MaxInt32 - character.Meso
		if int32(mesoCount) > cap {
			return nil, constant.LOOT_FAILED_MESO_FULL
		}

		// Remove meso from map using RemoveItem
		if err := m.RemoveItem(itemID, constant.REMOVE_ITEM_TYPE_ANIMATED, character.ID); err != nil {
			return nil, constant.LOOT_FAILED_INVALID_ITEM
		}

		// Return the meso object
		return item, constant.LOOT_SUCCESS

	default:
		return nil, constant.LOOT_FAILED_INVALID_ITEM
	}
}

func (m *Map) GetActorPID() *actor.PID {
	m.pidMutex.RLock()
	defer m.pidMutex.RUnlock()
	return m.actorPID
}

func (m *Map) SetActorPID(pid *actor.PID) {
	m.pidMutex.Lock()
	defer m.pidMutex.Unlock()
	m.actorPID = pid
}

func (m *Map) WarpCharacter(character *Character, targetMapID uint32, portal uint8) {
	if m != nil {
		if err := m.RemovePlayer(character.ID); err != nil {
			log.Fatalf("Failed to remove player %d from map: %v", character.ID, err)
			return
		}
	}
	character.Map = targetMapID
	character.SpawnPoint = portal
	m.listener.OnWarpCharacter(character, targetMapID, portal)
}
