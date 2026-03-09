package entity

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
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

type BroadcastRecipientFilter func(recipient *Character, referenceCharacter *Character) bool

type BroadcastOption struct {
	SendRaw            bool
	ExceptPlayerIDs    []uint32
	ReferenceCharacter *Character
	RecipientFilter    BroadcastRecipientFilter
}

func BroadcastVisibleByReference(recipient *Character, referenceCharacter *Character) bool {
	if referenceCharacter == nil || !referenceCharacter.IsHidden() {
		return true
	}
	return recipient.Role >= referenceCharacter.Role
}

func BroadcastRoleBelowReference(recipient *Character, referenceCharacter *Character) bool {
	if referenceCharacter == nil {
		return false
	}
	return recipient.Role < referenceCharacter.Role
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

	mapInstance.initializeNpcs()
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

func OnMobControllerChange(mob *Mob, before *Character, after *Character) {
}

func (m *Map) AddPlayer(playerID uint32, character *Character, spawnPoint uint8, init bool) error {
	if m.objects[types.OBJECT_TYPE_PLAYER] == nil {
		m.objects[types.OBJECT_TYPE_PLAYER] = make(map[uint32]interface{})
	}

	character.Map = m.ID
	character.spawnPoint = spawnPoint

	m.objects[types.OBJECT_TYPE_PLAYER][playerID] = character

	m.listener.OnPlayerAdded(m.ID, playerID, character, init)
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

func (m *Map) initializeNpcs() {
	if m.objects[types.OBJECT_TYPE_NPC] == nil {
		m.objects[types.OBJECT_TYPE_NPC] = make(map[uint32]interface{})
	}

	for _, wz := range m.model.NpcSpawns {
		oid := m.allocateOID()
		npc := &Npc{
			Object: Object{
				OID: oid,
				Position: types.Point[int16]{
					X: wz.BaseSpawn.Position.X,
					Y: wz.BaseSpawn.Position.Y,
				},
				Context: m.context,
			},
			Wz: &wz,
		}
		m.objects[types.OBJECT_TYPE_NPC][oid] = npc
	}
}

func (m *Map) initializeMobs() {
	for spawnId, mobSpawnSpec := range m.model.MobSpawns {
		m.MobSpawns[spawnId] = &MobSpawn{
			Wz:            &mobSpawnSpec,
			Spawned:       false,
			LastSpawnedAt: time.Time{},
		}
	}
}

// GetNpcs returns all NPCs on the map
func (m *Map) GetNpcs() map[uint32]interface{} {
	if m.objects[types.OBJECT_TYPE_NPC] == nil {
		return make(map[uint32]interface{})
	}
	return m.objects[types.OBJECT_TYPE_NPC]
}

func (m *Map) SpawnNpc(npcId uint32, position types.Point[int16]) (*Npc, error) {
	oid := m.allocateOID()

	footholdID := int16(0)
	foothold, ok := m.model.Footholds.Find(position)
	if ok {
		footholdID = foothold.ID
	}

	footholdPoint := m.model.FootholdPoint(position)
	spawnPosition := position
	if footholdPoint != nil {
		spawnPosition = *footholdPoint
	}

	baseSpawn := &wz.BaseSpawn{
		ID:              npcId,
		Position:        spawnPosition,
		RenderX0:        -50,
		RenderX1:        50,
		CollisionY:      spawnPosition.Y,
		Hide:            false,
		UseDay:          true,
		UseNight:        true,
		Foothold:        footholdID,
		FacingDirection: wz.FACING_DIRECTION_RIGHT,
		MobTime:         0,
		Info:            0,
		LimitedName:     "",
		NoFoothold:      false,
	}

	npcSpawn := wz.NpcSpawn{
		BaseSpawn: baseSpawn,
	}

	npc := &Npc{
		Object: Object{
			OID:      oid,
			Position: spawnPosition,
			Context:  m.context,
		},
		Wz: &npcSpawn,
	}

	if m.objects[types.OBJECT_TYPE_NPC] == nil {
		m.objects[types.OBJECT_TYPE_NPC] = make(map[uint32]interface{})
	}

	m.objects[types.OBJECT_TYPE_NPC][oid] = npc

	npcDTO := npc.ToDTO()
	spawnPacket := &response.SpawnNpc{
		NPC:     npcDTO,
		Visible: true,
	}
	controlPacket := &response.NpcControl{
		NPC:     npcDTO,
		MiniMap: true,
	}

	m.Broadcast(spawnPacket, nil)
	m.Broadcast(controlPacket, nil)

	return npc, nil
}

func (m *Map) SpawnMob(mobId uint32, position types.Point[int16], mobSpawn *MobSpawn) (*Mob, error) {
	oid := m.allocateOID()

	mobSpec, ok := m.context.GetResources().Monsters[mobId]
	if !ok {
		return nil, fmt.Errorf("mob model not found for ID: %d", mobId)
	}

	footholdID := int16(0)
	if mobSpawn != nil && mobSpawn.Wz != nil {
		footholdID = mobSpawn.Wz.Foothold
	}
	if footholdID == 0 {
		foothold, ok := m.model.Footholds.Find(position)
		if ok {
			footholdID = foothold.ID
		}
	}

	spawnPoint, ok := m.model.DropPoint(position)
	if !ok {
		spawnPoint = position
	}

	mob := &Mob{
		Life: Life{
			Object: Object{
				OID:      oid,
				Position: spawnPoint,
				Context:  m.context,
			},
			Hp:     uint16(mobSpec.MaxHP),
			BaseHp: uint16(mobSpec.MaxHP),
			Mp:     uint16(mobSpec.MaxMP),
			BaseMp: uint16(mobSpec.MaxMP),
			Stance: 5,
		},
		Foothold: footholdID,
		Wz:       mobSpec,
		MapID:    m.ID,
		Spawn:    mobSpawn,
	}

	if m.objects[types.OBJECT_TYPE_MONSTER] == nil {
		m.objects[types.OBJECT_TYPE_MONSTER] = make(map[uint32]interface{})
	}

	m.objects[types.OBJECT_TYPE_MONSTER][oid] = mob
	m.listener.OnMobSpawned(m.ID, oid, mob)
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

	if mob.Spawn != nil {
		mob.Spawn.Spawned = false
	}

	m.releaseOID(mobID)
	m.controllerTable.LeaveMob(mob)
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

// Broadcast sends a message to players on the map.
// When option is nil, defaults are used: encrypt policy, no except player, no filter.
func (m *Map) Broadcast(message types.Packet, option *BroadcastOption) {
	if m.objects[types.OBJECT_TYPE_PLAYER] == nil {
		return
	}

	policy := types.SEND_POLICY_ENCRYPT
	exceptPlayerIDs := []uint32(nil)
	var referenceCharacter *Character
	var recipientFilter BroadcastRecipientFilter

	if option != nil {
		if option.SendRaw {
			policy = types.SEND_POLICY_RAW
		}
		exceptPlayerIDs = option.ExceptPlayerIDs
		referenceCharacter = option.ReferenceCharacter
		recipientFilter = option.RecipientFilter
	}

	exceptSet := make(map[uint32]struct{}, len(exceptPlayerIDs))
	for _, playerID := range exceptPlayerIDs {
		exceptSet[playerID] = struct{}{}
	}

	for playerID, player := range m.objects[types.OBJECT_TYPE_PLAYER] {
		if _, excluded := exceptSet[playerID]; excluded {
			continue
		}

		character, ok := player.(*Character)
		if !ok {
			continue
		}

		if recipientFilter != nil && !recipientFilter(character, referenceCharacter) {
			continue
		}

		character.Send(message, policy)
	}
}

func (m *Map) SpawnItem(item Item, ownerID uint32, dropType constant.DropType) error {
	oid := m.allocateOID()

	drop := item.GetDrop()
	if drop == nil {
		return fmt.Errorf("item has no drop information")
	}

	dropPoint, ok := m.model.DropPoint(drop.Position)
	if !ok {
		dropPoint = drop.SpawnedPoint
	}

	drop.Position = dropPoint
	drop.OID = oid
	drop.Owner = ownerID
	drop.DropType = dropType
	drop.MapID = m.ID
	drop.RegisterExpire(constant.ITEM_EXPIRE_TIME)
	if dropType == constant.DROP_TYPE_OWNED || dropType == constant.DROP_TYPE_PARTY {
		drop.RegisterFFA(constant.ITEM_FFA_TIME)
	}

	if m.objects[types.OBJECT_TYPE_ITEM] == nil {
		m.objects[types.OBJECT_TYPE_ITEM] = make(map[uint32]interface{})
	}

	m.objects[types.OBJECT_TYPE_ITEM][oid] = item
	m.listener.OnItemSpawned(m.ID, oid, item, drop)

	return nil
}

func (m *Map) SpawnMeso(count int32, position types.Point[int16], ownerID uint32, dropType constant.DropType) error {
	oid := m.allocateOID()

	dropPoint, ok := m.model.DropPoint(position)
	if !ok {
		dropPoint = position
	}

	// Create meso entity
	meso := NewMeso(count, dropPoint, ownerID, dropType, oid, m.context, m.ID)

	// Register timers for meso
	drop := meso.GetDrop()
	if drop != nil {
		drop.RegisterExpire(constant.ITEM_EXPIRE_TIME)
		if dropType == constant.DROP_TYPE_OWNED || dropType == constant.DROP_TYPE_PARTY {
			drop.RegisterFFA(constant.ITEM_FFA_TIME)
		}
	}

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
		if drop.DropType == constant.DROP_TYPE_OWNED && drop.Owner != character.GetID() {
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
		if err := m.RemoveItem(itemID, constant.REMOVE_ITEM_TYPE_ANIMATED, character.GetID()); err != nil {
			return nil, constant.LOOT_FAILED_INVALID_ITEM
		}

		return item, constant.LOOT_SUCCESS

	case *Meso:
		drop := item.GetDrop()
		if drop == nil {
			return nil, constant.LOOT_FAILED_INVALID_ITEM
		}

		// Check ownership for owned drops
		if drop.DropType == constant.DROP_TYPE_OWNED && drop.Owner != character.GetID() {
			return nil, constant.LOOT_FAILED_NO_OWNERSHIP
		}

		// Check meso capacity before removing from map
		mesoCount := item.GetCount32()
		cap := math.MaxInt32 - character.Meso
		if int32(mesoCount) > cap {
			return nil, constant.LOOT_FAILED_MESO_FULL
		}

		// Remove meso from map using RemoveItem
		if err := m.RemoveItem(itemID, constant.REMOVE_ITEM_TYPE_ANIMATED, character.GetID()); err != nil {
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

// Luable interface implementation
func (m *Map) LuaTypeName() string {
	return "LuaMap"
}

func (m *Map) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"npcs": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			npcs := mapInstance.GetNpcs()
			tbl := L.NewTable()
			for _, npc := range npcs {
				if npcObj, ok := npc.(*Npc); ok {
					tbl.RawSetInt(int(npcObj.OID), luax.NewLuable(L, npcObj))
				}
			}
			L.Push(tbl)
			return 1
		},
		"mobs": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			mobs := mapInstance.GetMobs()
			tbl := L.NewTable()
			for _, mob := range mobs {
				if mobObj, ok := mob.(*Mob); ok {
					tbl.RawSetInt(int(mobObj.OID), luax.NewLuable(L, mobObj))
				}
			}
			L.Push(tbl)
			return 1
		},
		"characters": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			players := mapInstance.GetAllPlayers()
			tbl := L.NewTable()
			for _, player := range players {
				if char, ok := player.(*Character); ok {
					tbl.RawSetInt(int(char.GetID()), luax.NewLuable(L, char))
				}
			}
			L.Push(tbl)
			return 1
		},
		"items": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			items := mapInstance.GetItems()
			tbl := L.NewTable()
			for _, item := range items {
				if itemObj, ok := item.(Item); ok {
					drop := itemObj.GetDrop()
					if drop != nil && drop.Object != nil {
						// Item interface를 타입 어설션하여 실제 타입을 얻고 Luable로 변환
						switch v := itemObj.(type) {
						case *Equipment:
							tbl.RawSetInt(int(drop.OID), luax.NewLuable(L, v))
						case *Consume:
							tbl.RawSetInt(int(drop.OID), luax.NewLuable(L, v))
						case *CashItem:
							tbl.RawSetInt(int(drop.OID), luax.NewLuable(L, v))
						case *GeneralItem:
							tbl.RawSetInt(int(drop.OID), luax.NewLuable(L, v))
						case *Installation:
							tbl.RawSetInt(int(drop.OID), luax.NewLuable(L, v))
						case *Pet:
							tbl.RawSetInt(int(drop.OID), luax.NewLuable(L, v))
						default:
							// Fallback to Object
							tbl.RawSetInt(int(drop.OID), luax.NewLuable(L, drop.Object))
						}
					}
				}
			}
			L.Push(tbl)
			return 1
		},
		"wz": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			spec := mapInstance.GetSpec()
			if spec == nil {
				L.Push(lua.LNil)
				return 1
			}
			tbl := L.NewTable()
			tbl.RawSetString("id", lua.LNumber(spec.ID))
			tbl.RawSetString("name", lua.LString(spec.Name))
			tbl.RawSetString("return_map_id", lua.LNumber(spec.ReturnMapId))
			tbl.RawSetString("town", lua.LBool(spec.IsTown))
			L.Push(tbl)
			return 1
		},
	}
}

func (m *Map) String() string {
	return m.LuaTypeName()
}

func (m *Map) Type() lua.LValueType {
	return lua.LTUserData
}
