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

// MapListener defines interface for map events. Prefer passing object pointers (*Map, *Character, *Mob) over IDs.
type MapListener interface {
	OnPlayerAdded(mapInstance *Map, character *Character, init bool)
	OnPlayerRemoved(mapInstance *Map, character *Character)
	OnPlayerMoved(mapInstance *Map, character *Character)
	OnPlayerMove(mapInstance *Map, character *Character, startPoint types.Vector2[int16], fragments []dto.MoveFragment)
	OnPlayerChat(mapInstance *Map, character *Character, message string)
	OnItemSpawned(mapInstance *Map, item Item, drop *Drop)
	OnMesoSpawned(mapInstance *Map, meso *Meso)
	OnItemRemoved(mapInstance *Map, itemID uint32, looterID uint32, mode constant.RemoveItemType)
	OnMobSpawned(mapInstance *Map, mob *Mob)
	OnMobRemoved(mapInstance *Map, mob *Mob, animationType constant.MobDieAnimationType)
	OnMobControllerChange(mob *Mob, before *Character, after *Character)
	OnMobMoved(mapInstance *Map, mob *Mob, isAggroed bool, centerSplit int8, skill1 uint8, skill2 uint8, skill3 uint8, skill4 uint8, startPoint types.Vector2[int16], movements []dto.MoveFragment)
	OnAttack(mapInstance *Map, character *Character, attackInfo dto.AttackInfo, skillLevel uint8)
	OnMobMobStatusApplied(mapInstance *Map, mob *Mob, debuff constant.MobStatus, value int32, skillID uint32, durationMs int64)
	OnMobMobStatusCancelled(mapInstance *Map, mob *Mob, debuff constant.MobStatus)
	OnMistSpawned(mapInstance *Map, mist *Mist)
	OnMistRemoved(mapInstance *Map, mist *Mist)
	OnDoorSpawned(mapInstance *Map, door *Door)
	OnDoorRemoved(mapInstance *Map, door *Door, animated bool)
}

type MobSpawn struct {
	Wz            *wz.MobSpawn
	Spawned       bool
	LastSpawnedAt time.Time
}

type Map struct {
	Wz *wz.Map // Map specification data

	id              uint32
	objects         map[constant.ObjectType]map[uint32]Object // Players, Mobs, Items, etc.
	controllerTable *ControllerTable
	MobSpawns       map[uint32]*MobSpawn
	listener        MapListener
	sequence        uint32      // Sequence ID for generating unique object IDs
	availableOIDs   []uint32    // Queue of available OIDs for reuse
	context         GameContext // GameContext for accessing resources
	actorPID        *actor.PID
	pidMutex        sync.RWMutex
}

type BroadcastRecipientFilter func(recipient *Character, reference Object) bool

type BroadcastOption struct {
	SendRaw         bool
	ExceptPlayerIDs []uint32
	Reference       Object
	RecipientFilter BroadcastRecipientFilter
}

func BroadcastVisibleByReference(recipient *Character, reference Object) bool {
	if reference == nil {
		return true
	}
	refCharacter, ok := reference.(*Character)
	if !ok {
		return true
	}
	if !refCharacter.IsHidden() {
		return true
	}
	return recipient.Role >= refCharacter.Role
}

func BroadcastRoleBelowReference(recipient *Character, reference Object) bool {
	if reference == nil {
		return false
	}
	refCharacter, ok := reference.(*Character)
	if !ok {
		return true
	}
	return recipient.Role < refCharacter.Role
}

func NewMap(id uint32, listener MapListener, mapId uint32, context GameContext) *Map {
	if listener == nil {
		panic("MapListener cannot be nil")
	}
	if context == nil {
		panic("GameContext cannot be nil")
	}

	// Get map model from resources
	wz, ok := context.GetResources().Maps[mapId]
	if !ok {
		panic(fmt.Sprintf("Wz not found for ID: %d", mapId))
	}

	mapInstance := &Map{
		id:              id,
		objects:         make(map[constant.ObjectType]map[uint32]Object),
		controllerTable: nil, // Will be set after mapInstance is created
		MobSpawns:       make(map[uint32]*MobSpawn),
		listener:        listener,
		Wz:              wz,
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
	if m.objects[constant.ObjectTypeCharacter] == nil {
		m.objects[constant.ObjectTypeCharacter] = make(map[uint32]Object)
	}

	character.Map = m
	character.spawnPoint = spawnPoint
	if pos, ok := m.Wz.GetSpawnPosition(spawnPoint); ok {
		character.Position = pos
	}

	m.objects[constant.ObjectTypeCharacter][playerID] = character

	m.listener.OnPlayerAdded(m, character, init)
	m.controllerTable.EnterPlayer(character)

	for _, summon := range character.GetSummons() {
		if summon == nil || summon.Owner != character {
			continue
		}
		summon.Position = character.Position
		m.AddSummon(summon)
	}

	return nil
}

func (m *Map) RemovePlayer(playerID uint32) error {
	if m.objects[constant.ObjectTypeCharacter] == nil {
		return fmt.Errorf("no players on map")
	}

	if _, exists := m.objects[constant.ObjectTypeCharacter][playerID]; !exists {
		return fmt.Errorf("player %d not found on map", playerID)
	}

	character := m.objects[constant.ObjectTypeCharacter][playerID].(*Character)

	for _, summon := range character.GetSummons() {
		if summon == nil || summon.Map != m || summon.OID == 0 {
			continue
		}
		m.RemoveSummon(summon.OID, false)
	}

	delete(m.objects[constant.ObjectTypeCharacter], playerID)

	character.SuspendTimers()
	character.Map = nil
	m.controllerTable.LeavePlayer(character)

	m.listener.OnPlayerRemoved(m, character)

	return nil
}

func (m *Map) GetPlayer(playerID uint32) *Character {
	if m.objects[constant.ObjectTypeCharacter] == nil {
		return nil
	}

	if player, ok := m.objects[constant.ObjectTypeCharacter][playerID].(*Character); ok {
		return player
	}
	return nil
}

func (m *Map) GetPlayerCount() int {
	if m.objects[constant.ObjectTypeCharacter] == nil {
		return 0
	}
	return len(m.objects[constant.ObjectTypeCharacter])
}

func (m *Map) GetAllPlayers() map[uint32]Object {
	if m.objects[constant.ObjectTypeCharacter] == nil {
		return make(map[uint32]Object)
	}
	return m.objects[constant.ObjectTypeCharacter]
}

func (m *Map) GetControllerTable() *ControllerTable {
	return m.controllerTable
}

func (m *Map) GetRecoveryRate() float32 {
	if m.Wz == nil {
		return 1.0
	}
	if m.Wz.RecoveryRate <= 0 {
		return 1.0
	}
	return m.Wz.RecoveryRate
}

func (m *Map) FootholdPoint(point types.Point[int16]) *types.Point[int16] {
	return m.Wz.FootholdPoint(point)
}

func (m *Map) initializeNpcs() {
	if m.objects[constant.ObjectTypeNpc] == nil {
		m.objects[constant.ObjectTypeNpc] = make(map[uint32]Object)
	}

	for _, wz := range m.Wz.NpcSpawns {
		oid := m.allocateOID()
		npc := &Npc{
			ObjectCore: ObjectCore{
				OID:      oid,
				Position: types.Point[int16]{X: wz.BaseSpawn.Position.X, Y: wz.BaseSpawn.Position.Y},
				Context:  m.context,
				Map:      m,
			},
			Wz: &wz,
		}
		m.objects[constant.ObjectTypeNpc][oid] = npc
	}
}

func (m *Map) AddSummon(s *Summon) {
	if s == nil {
		return
	}
	if s.Map == m && m.objects[constant.ObjectTypeSummon] != nil {
		if existing, ok := m.objects[constant.ObjectTypeSummon][s.OID]; ok && existing == s {
			return
		}
	}
	if s.Map != nil && s.Map != m {
		return
	}
	s.ObjectCore.Context = m.context
	if m.objects[constant.ObjectTypeSummon] == nil {
		m.objects[constant.ObjectTypeSummon] = make(map[uint32]Object)
	}
	if s.OID == 0 {
		s.OID = m.allocateOID()
	}
	if s.Map == nil {
		s.Map = m
	}
	m.objects[constant.ObjectTypeSummon][s.OID] = s

	if s.Owner != nil && s.Owner.Listener != nil {
		s.Owner.Listener.OnSummonSpawn(s.Owner, s)
	}
}

func (m *Map) RemoveSummon(oid uint32, animated bool) {
	if m.objects[constant.ObjectTypeSummon] == nil {
		return
	}
	obj, ok := m.objects[constant.ObjectTypeSummon][oid]
	if !ok {
		return
	}

	s, isSummon := obj.(*Summon)
	if !isSummon {
		delete(m.objects[constant.ObjectTypeSummon], oid)
		m.releaseOID(oid)
		return
	}

	owner := s.Owner
	if owner != nil && owner.Listener != nil {
		owner.Listener.OnSummonRemove(owner, s, animated)
	}

	delete(m.objects[constant.ObjectTypeSummon], oid)
	m.releaseOID(oid)
	if s.Map == m {
		s.Map = nil
		s.OID = 0
	}
}

func (m *Map) GetSummon(oid uint32) *Summon {
	if m.objects[constant.ObjectTypeSummon] == nil {
		return nil
	}
	if s, ok := m.objects[constant.ObjectTypeSummon][oid].(*Summon); ok {
		return s
	}
	return nil
}

func (m *Map) AddMist(mist *Mist) {
	if mist == nil {
		return
	}
	if mist.Map == m && m.objects[constant.ObjectTypeMist] != nil {
		if existing, ok := m.objects[constant.ObjectTypeMist][mist.OID]; ok && existing == mist {
			return
		}
	}
	if mist.Map != nil && mist.Map != m {
		return
	}
	mist.ObjectCore.Context = m.context
	if m.objects[constant.ObjectTypeMist] == nil {
		m.objects[constant.ObjectTypeMist] = make(map[uint32]Object)
	}
	if mist.OID == 0 {
		mist.OID = m.allocateOID()
	}
	if mist.Map == nil {
		mist.Map = m
	}
	m.objects[constant.ObjectTypeMist][mist.OID] = mist
	m.listener.OnMistSpawned(m, mist)
}

func (m *Map) RemoveMist(oid uint32) {
	if m.objects[constant.ObjectTypeMist] == nil {
		return
	}
	obj, ok := m.objects[constant.ObjectTypeMist][oid]
	if !ok {
		return
	}

	mi, isMist := obj.(*Mist)
	if !isMist {
		delete(m.objects[constant.ObjectTypeMist], oid)
		m.releaseOID(oid)
		return
	}

	m.listener.OnMistRemoved(m, mi)

	delete(m.objects[constant.ObjectTypeMist], oid)
	m.releaseOID(oid)
	if mi.Map == m {
		mi.Map = nil
		mi.OID = 0
	}
}

func (m *Map) GetMist(oid uint32) *Mist {
	if m.objects[constant.ObjectTypeMist] == nil {
		return nil
	}
	if mi, ok := m.objects[constant.ObjectTypeMist][oid].(*Mist); ok {
		return mi
	}
	return nil
}

func (m *Map) AddDoor(door *Door) {
	if door == nil {
		return
	}
	if door.Map == m && m.objects[constant.ObjectTypeDoor] != nil {
		if existing, ok := m.objects[constant.ObjectTypeDoor][door.OID]; ok && existing == door {
			return
		}
	}
	if door.Map != nil && door.Map != m {
		return
	}
	door.ObjectCore.Context = m.context
	if m.objects[constant.ObjectTypeDoor] == nil {
		m.objects[constant.ObjectTypeDoor] = make(map[uint32]Object)
	}
	if door.OID == 0 {
		door.OID = m.allocateOID()
	}
	if door.Map == nil {
		door.Map = m
	}
	m.objects[constant.ObjectTypeDoor][door.OID] = door
	m.listener.OnDoorSpawned(m, door)
}

func (m *Map) RemoveDoor(oid uint32, animated bool) {
	if m.objects[constant.ObjectTypeDoor] == nil {
		return
	}
	obj, ok := m.objects[constant.ObjectTypeDoor][oid]
	if !ok {
		return
	}

	door, isDoor := obj.(*Door)
	if !isDoor {
		delete(m.objects[constant.ObjectTypeDoor], oid)
		m.releaseOID(oid)
		return
	}

	m.listener.OnDoorRemoved(m, door, animated)

	delete(m.objects[constant.ObjectTypeDoor], oid)
	m.releaseOID(oid)
	if door.Map == m {
		door.Map = nil
		door.OID = 0
	}
}

func (m *Map) GetDoor(oid uint32) *Door {
	if m.objects[constant.ObjectTypeDoor] == nil {
		return nil
	}
	if door, ok := m.objects[constant.ObjectTypeDoor][oid].(*Door); ok {
		return door
	}
	return nil
}

func (m *Map) initializeMobs() {
	for spawnId, mobSpawnSpec := range m.Wz.MobSpawns {
		m.MobSpawns[spawnId] = &MobSpawn{
			Wz:            &mobSpawnSpec,
			Spawned:       false,
			LastSpawnedAt: time.Time{},
		}
	}
}

func (m *Map) GetNpcs() map[uint32]Object {
	if m.objects[constant.ObjectTypeNpc] == nil {
		return make(map[uint32]Object)
	}
	return m.objects[constant.ObjectTypeNpc]
}

func (m *Map) SpawnNpc(npcId uint32, position types.Point[int16]) (*Npc, error) {
	oid := m.allocateOID()

	footholdID := int16(0)
	foothold, ok := m.Wz.Footholds.Find(position)
	if ok {
		footholdID = foothold.ID
	}

	footholdPoint := m.Wz.FootholdPoint(position)
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
		ObjectCore: ObjectCore{
			OID:      oid,
			Position: spawnPosition,
			Context:  m.context,
			Map:      m,
		},
		Wz: &npcSpawn,
	}

	if m.objects[constant.ObjectTypeNpc] == nil {
		m.objects[constant.ObjectTypeNpc] = make(map[uint32]Object)
	}

	m.objects[constant.ObjectTypeNpc][oid] = npc

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
		foothold, ok := m.Wz.Footholds.Find(position)
		if ok {
			footholdID = foothold.ID
		}
	}

	spawnPoint, ok := m.Wz.DropPoint(position)
	if !ok {
		spawnPoint = position
	}

	mob := &Mob{
		LifeCore: LifeCore{
			ObjectCore: ObjectCore{
				OID:      oid,
				Position: spawnPoint,
				Context:  m.context,
				Map:      m,
			},
			Hp:     uint32(max(0, mobSpec.MaxHP)),
			Mp:     uint32(max(0, mobSpec.MaxMP)),
			BaseHp: uint32(max(0, mobSpec.MaxHP)),
			BaseMp: uint32(max(0, mobSpec.MaxMP)),
			Stance: 5,
		},
		Foothold: footholdID,
		Wz:       mobSpec,
		Spawn:    mobSpawn,
	}

	if m.objects[constant.ObjectTypeMob] == nil {
		m.objects[constant.ObjectTypeMob] = make(map[uint32]Object)
	}

	m.objects[constant.ObjectTypeMob][oid] = mob
	m.listener.OnMobSpawned(m, mob)
	m.controllerTable.EnterMob(mob)

	return mob, nil
}

func (m *Map) RemoveMob(mobID uint32, animationType constant.MobDieAnimationType) error {
	if m.objects[constant.ObjectTypeMob] == nil {
		return fmt.Errorf("no monsters on map")
	}

	if _, exists := m.objects[constant.ObjectTypeMob][mobID]; !exists {
		return fmt.Errorf("mob %d not found on map", mobID)
	}

	mob := m.objects[constant.ObjectTypeMob][mobID].(*Mob)
	delete(m.objects[constant.ObjectTypeMob], mobID)

	mob.ClearAllMobStatusTimers()
	if mob.Spawn != nil {
		mob.Spawn.Spawned = false
	}

	m.releaseOID(mobID)
	m.controllerTable.LeaveMob(mob)
	m.listener.OnMobRemoved(m, mob, animationType)

	return nil
}

func (m *Map) GetMob(mobID uint32) *Mob {
	if m.objects[constant.ObjectTypeMob] == nil {
		return nil
	}

	if mob, ok := m.objects[constant.ObjectTypeMob][mobID].(*Mob); ok {
		return mob
	}
	return nil
}

func (m *Map) GetMobs() map[uint32]Object {
	if m.objects[constant.ObjectTypeMob] == nil {
		return make(map[uint32]Object)
	}
	return m.objects[constant.ObjectTypeMob]
}

type ObjectsFilter struct {
	Area     *struct{ MinX, MinY, MaxX, MaxY int16 }
	Distance *struct {
		Dist int
		X, Y int16
	}
}

func (m *Map) GetObjects(filter constant.ObjectType, opts *ObjectsFilter) []Object {

	out := make([]Object, 0)
	for _, buckets := range m.objects {
		for _, obj := range buckets {
			if obj == nil {
				continue
			}
			if !obj.Is(filter) {
				continue
			}

			if opts != nil {
				pos := obj.GetPosition()
				if opts.Area != nil {
					a := opts.Area
					if pos.X < a.MinX || pos.X > a.MaxX || pos.Y < a.MinY || pos.Y > a.MaxY {
						continue
					}
				}
			}
			out = append(out, obj)
		}
	}
	return out
}

// Broadcast sends a message to players on the map.
// When option is nil, defaults are used: encrypt policy, no except player, no filter.
func (m *Map) Broadcast(message types.Packet, option *BroadcastOption) {
	if m.objects[constant.ObjectTypeCharacter] == nil {
		return
	}

	policy := types.SEND_POLICY_ENCRYPT
	exceptPlayerIDs := []uint32(nil)
	var reference Object
	var recipientFilter BroadcastRecipientFilter

	if option != nil {
		if option.SendRaw {
			policy = types.SEND_POLICY_RAW
		}
		exceptPlayerIDs = option.ExceptPlayerIDs
		reference = option.Reference
		recipientFilter = option.RecipientFilter
	}

	exceptSet := make(map[uint32]struct{}, len(exceptPlayerIDs))
	for _, playerID := range exceptPlayerIDs {
		exceptSet[playerID] = struct{}{}
	}

	for playerID, player := range m.objects[constant.ObjectTypeCharacter] {
		if _, excluded := exceptSet[playerID]; excluded {
			continue
		}

		character, ok := player.(*Character)
		if !ok {
			continue
		}

		if recipientFilter != nil && !recipientFilter(character, reference) {
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

	dropPoint, ok := m.Wz.DropPoint(drop.Position)
	if !ok {
		dropPoint = drop.SpawnedPoint
	}

	drop.Position = dropPoint
	drop.OID = oid
	drop.Owner = ownerID
	drop.DropType = dropType
	drop.Map = m
	drop.RegisterExpire(constant.ITEM_EXPIRE_TIME)
	if dropType == constant.DROP_TYPE_OWNED || dropType == constant.DROP_TYPE_PARTY {
		drop.RegisterFFA(constant.ITEM_FFA_TIME)
	}

	if m.objects[constant.ObjectTypeItem] == nil {
		m.objects[constant.ObjectTypeItem] = make(map[uint32]Object)
	}
	mapObj, ok := item.(Object)
	if !ok {
		return fmt.Errorf("item must implement Object")
	}
	m.objects[constant.ObjectTypeItem][oid] = mapObj
	m.listener.OnItemSpawned(m, item, drop)

	return nil
}

func (m *Map) SpawnMeso(count int32, position types.Point[int16], ownerID uint32, dropType constant.DropType) (*Meso, error) {
	oid := m.allocateOID()

	dropPoint, ok := m.Wz.DropPoint(position)
	if !ok {
		dropPoint = position
	}

	meso := NewMeso(count, dropPoint, ownerID, dropType, oid, m.context, m)

	// Register timers for meso
	drop := meso.GetDrop()
	if drop != nil {
		drop.RegisterExpire(constant.ITEM_EXPIRE_TIME)
		if dropType == constant.DROP_TYPE_OWNED || dropType == constant.DROP_TYPE_PARTY {
			drop.RegisterFFA(constant.ITEM_FFA_TIME)
		}
	}

	// Initialize objects map for items if needed
	if m.objects[constant.ObjectTypeItem] == nil {
		m.objects[constant.ObjectTypeItem] = make(map[uint32]Object)
	}

	// Add meso to map objects
	m.objects[constant.ObjectTypeItem][oid] = meso

	// Notify listener about meso spawn
	m.listener.OnMesoSpawned(m, meso)

	return meso, nil
}

// RemoveItem removes an item from the map
func (m *Map) RemoveItem(itemID uint32, removeType constant.RemoveItemType, playerID uint32) error {
	if m.objects[constant.ObjectTypeItem] == nil {
		return fmt.Errorf("no items on map")
	}

	if _, exists := m.objects[constant.ObjectTypeItem][itemID]; !exists {
		return fmt.Errorf("item %d not found on map", itemID)
	}

	delete(m.objects[constant.ObjectTypeItem], itemID)

	// Release OID for reuse
	m.releaseOID(itemID)

	// Notify listener about item removal
	m.listener.OnItemRemoved(m, itemID, playerID, removeType)

	return nil
}

// GetItem retrieves an item from the map
func (m *Map) GetItem(itemID uint32) Item {
	if m.objects[constant.ObjectTypeItem] == nil {
		return nil
	}

	if item, ok := m.objects[constant.ObjectTypeItem][itemID].(Item); ok {
		return item
	}
	return nil
}

// GetItems returns all items on the map
func (m *Map) GetItems() map[uint32]Object {
	if m.objects[constant.ObjectTypeItem] == nil {
		return make(map[uint32]Object)
	}
	return m.objects[constant.ObjectTypeItem]
}

// LootItem attempts to loot an item from the map and returns the looted item and reason
// All capacity checks are performed before removing the item from the map
func (m *Map) LootItem(itemID uint32, character *Character, position types.Point[int16]) (interface{}, constant.LootResult) {
	if m.objects[constant.ObjectTypeItem] == nil {
		return nil, constant.LOOT_FAILED_ITEM_NOT_FOUND
	}

	itemInterface, exists := m.objects[constant.ObjectTypeItem][itemID]
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
		"objects": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			filter := constant.ObjectTypeObject
			if L.GetTop() >= 2 {
				if lv := L.Get(2); lv.Type() == lua.LTNumber {
					filter = constant.ObjectType(lua.LVAsNumber(lv))
				}
			}
			var opts *ObjectsFilter
			if L.GetTop() >= 3 {
				tbl := L.CheckTable(3)
				opts = &ObjectsFilter{}
				readArea := func(t *lua.LTable) {
					minX := t.RawGetString("minX")
					minY := t.RawGetString("minY")
					maxX := t.RawGetString("maxX")
					maxY := t.RawGetString("maxY")
					if minX.Type() == lua.LTNumber && minY.Type() == lua.LTNumber && maxX.Type() == lua.LTNumber && maxY.Type() == lua.LTNumber {
						opts.Area = &struct{ MinX, MinY, MaxX, MaxY int16 }{
							MinX: int16(lua.LVAsNumber(minX)),
							MinY: int16(lua.LVAsNumber(minY)),
							MaxX: int16(lua.LVAsNumber(maxX)),
							MaxY: int16(lua.LVAsNumber(maxY)),
						}
					}
				}
				if area := tbl.RawGetString("area"); area.Type() == lua.LTTable {
					readArea(area.(*lua.LTable))
				} else if tbl.RawGetString("minX").Type() == lua.LTNumber {
					readArea(tbl)
				}
				if distLV := tbl.RawGetString("distance"); distLV.Type() == lua.LTNumber {
					xLV := tbl.RawGetString("x")
					yLV := tbl.RawGetString("y")
					x, y := int16(0), int16(0)
					if xLV.Type() == lua.LTNumber {
						x = int16(lua.LVAsNumber(xLV))
					}
					if yLV.Type() == lua.LTNumber {
						y = int16(lua.LVAsNumber(yLV))
					}
					opts.Distance = &struct {
						Dist int
						X, Y int16
					}{
						Dist: int(lua.LVAsNumber(distLV)),
						X:    x,
						Y:    y,
					}
				}
				if opts.Area == nil && opts.Distance == nil {
					opts = nil
				}
			}
			objs := mapInstance.GetObjects(filter, opts)
			result := L.NewTable()
			idx := 0
			for _, v := range objs {
				if luable, ok := v.(luax.Luable); ok {
					idx++
					result.RawSetInt(idx, luax.NewLuable(L, luable))
				}
			}
			L.Push(result)
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
					if drop != nil && drop.ObjectCore != nil {
						tbl.RawSetInt(int(drop.OID), luax.NewLuable(L, itemObj))
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
			spec := mapInstance.Wz
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
		"recovery_rate": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			L.Push(lua.LNumber(mapInstance.GetRecoveryRate()))
			return 1
		},
		"spawn_meso": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			count := int32(L.CheckInt(2))
			posTbl := L.CheckTable(3)
			var x, y int16
			if lx := posTbl.RawGetInt(1); lx != lua.LNil {
				x = int16(lua.LVAsNumber(lx))
			} else if lx := posTbl.RawGetString("x"); lx != lua.LNil {
				x = int16(lua.LVAsNumber(lx))
			}
			if ly := posTbl.RawGetInt(2); ly != lua.LNil {
				y = int16(lua.LVAsNumber(ly))
			} else if ly := posTbl.RawGetString("y"); ly != lua.LNil {
				y = int16(lua.LVAsNumber(ly))
			}
			var owner *Character
			if ownerLV := L.Get(4); ownerLV != lua.LNil {
				if ownerUd, ok := ownerLV.(*lua.LUserData); ok {
					if ch, ok := ownerUd.Value.(*Character); ok {
						owner = ch
					}
				}
			}
			if count <= 0 {
				return 0
			}
			pos := types.Point[int16]{X: x, Y: y}
			dropType := constant.DROP_TYPE_FFA
			ownerID := uint32(0)
			if owner != nil {
				dropType = constant.DROP_TYPE_OWNED
				ownerID = owner.GetID()
			}
			meso, err := mapInstance.SpawnMeso(count, pos, ownerID, dropType)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, meso))
			return 1
		},
		"spawn_item": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			if mapInstance.context == nil {
				return 0
			}
			resources := mapInstance.context.GetResources()
			if resources == nil {
				return 0
			}

			argc := L.GetTop()
			if argc < 3 {
				L.ArgError(2, "spawn_item(itemIdOrName, count, position [, owner]) requires at least 3 arguments")
				return 0
			}

			var itemId uint32
			switch lv := L.Get(2).(type) {
			case lua.LString:
				id, ok := resources.NameToItem(string(lv))
				if !ok {
					return 0
				}
				itemId = id
			case lua.LNumber:
				itemId = uint32(lv)
			default:
				L.ArgError(2, "item id (number) or item name (string) expected")
				return 0
			}

			if _, ok := resources.Items[itemId]; !ok {
				return 0
			}

			count := uint16(1)
			if n := L.CheckInt(3); n >= 1 {
				count = uint16(n)
			}

			posTbl := L.CheckTable(4)
			var x, y int16
			if lx := posTbl.RawGetInt(1); lx != lua.LNil {
				x = int16(lua.LVAsNumber(lx))
			} else if lx := posTbl.RawGetString("x"); lx != lua.LNil {
				x = int16(lua.LVAsNumber(lx))
			}
			if ly := posTbl.RawGetInt(2); ly != lua.LNil {
				y = int16(lua.LVAsNumber(ly))
			} else if ly := posTbl.RawGetString("y"); ly != lua.LNil {
				y = int16(lua.LVAsNumber(ly))
			}
			pos := types.Point[int16]{X: x, Y: y}

			var owner *Character
			if argc >= 5 {
				if ownerLV := L.Get(5); ownerLV != lua.LNil {
					if ownerUd, ok := ownerLV.(*lua.LUserData); ok {
						if ch, ok := ownerUd.Value.(*Character); ok {
							owner = ch
						}
					}
				}
			}

			item, err := NewItem(itemId, count, mapInstance.context)
			if err != nil {
				return 0
			}
			dropType := constant.DROP_TYPE_FFA
			ownerID := uint32(0)
			if owner != nil {
				dropType = constant.DROP_TYPE_OWNED
				ownerID = owner.GetID()
			}
			item.BindDrop(&Drop{
				ObjectCore:   &ObjectCore{Position: pos},
				Owner:        ownerID,
				SpawnedPoint: pos,
				DropType:     dropType,
			})

			if err := mapInstance.SpawnItem(item, ownerID, dropType); err != nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, item))
			return 1
		},
		"spawn_mob": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			if mapInstance.context == nil {
				L.RaiseError("spawn_mob: map has no context")
				return 0
			}
			resources := mapInstance.context.GetResources()
			if resources == nil {
				L.RaiseError("spawn_mob: no resources")
				return 0
			}
			var mobID uint32
			switch lv := L.Get(2).(type) {
			case lua.LString:
				id, ok := resources.NameToMob(string(lv))
				if !ok {
					L.RaiseError("spawn_mob: unknown mob name %q", string(lv))
					return 0
				}
				mobID = id
			case lua.LNumber:
				mobID = uint32(lv)
			default:
				L.ArgError(2, "mob id (number) or mob name (string) expected")
				return 0
			}
			x := int16(L.CheckInt(3))
			y := int16(L.CheckInt(4))
			pos := types.Point[int16]{X: x, Y: y}
			mob, err := mapInstance.SpawnMob(mobID, pos, nil)
			if err != nil {
				L.RaiseError("spawn_mob: %v", err)
				return 0
			}
			L.Push(luax.NewLuable(L, mob))
			return 1
		},
		"spawn_npc": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			if mapInstance.context == nil {
				L.RaiseError("spawn_npc: map has no context")
				return 0
			}
			resources := mapInstance.context.GetResources()
			if resources == nil {
				L.RaiseError("spawn_npc: no resources")
				return 0
			}
			var npcID uint32
			switch lv := L.Get(2).(type) {
			case lua.LString:
				id, ok := resources.NameToNpc(string(lv))
				if !ok {
					L.RaiseError("spawn_npc: unknown npc name %q", string(lv))
					return 0
				}
				npcID = id
			case lua.LNumber:
				npcID = uint32(lv)
			default:
				L.ArgError(2, "npc id (number) or npc name (string) expected")
				return 0
			}
			x := int16(L.CheckInt(3))
			y := int16(L.CheckInt(4))
			pos := types.Point[int16]{X: x, Y: y}
			npc, err := mapInstance.SpawnNpc(npcID, pos)
			if err != nil {
				L.RaiseError("spawn_npc: %v", err)
				return 0
			}
			L.Push(luax.NewLuable(L, npc))
			return 1
		},
		"remove_mob": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			oid := uint32(L.CheckInt(2))
			animType := constant.MOB_DIE_ANIMATION_TYPE_FADE_OUT
			if L.GetTop() >= 3 {
				animType = constant.MobDieAnimationType(L.CheckInt(3))
			}
			err := mapInstance.RemoveMob(oid, animType)
			if err != nil {
				L.RaiseError("remove_mob: %v", err)
				return 0
			}
			return 0
		},
		"remove_mist": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mapInstance, ok := ud.Value.(*Map)
			if !ok {
				L.ArgError(1, "Map expected")
				return 0
			}
			arg := L.Get(2)
			switch v := arg.(type) {
			case *lua.LUserData:
				mist, ok := v.Value.(*Mist)
				if !ok || mist == nil {
					L.ArgError(2, "Mist or mist OID expected")
					return 0
				}
				if mist.GetMap() != mapInstance || mist.OID == 0 {
					return 0
				}
				mapInstance.RemoveMist(mist.OID)
			case lua.LNumber:
				oid := uint32(v)
				if oid == 0 {
					return 0
				}
				mapInstance.RemoveMist(oid)
			default:
				L.ArgError(2, "Mist or mist OID expected")
				return 0
			}
			return 0
		},
	}
}

func (m *Map) String() string {
	return m.LuaTypeName()
}

func (m *Map) Type() lua.LValueType {
	return lua.LTUserData
}
