package entity

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
)

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
	OnMobHomingRemoved(mapInstance *Map, mob *Mob, removed *Homing, causer *Character)
	OnMobHomingSet(mapInstance *Map, mob *Mob, homing *Homing, causer *Character)
	OnMobControllerChange(mob *Mob, before *Character, after *Character)
	OnMobMoved(mapInstance *Map, mob *Mob, isAggroed bool, centerSplit int8, skill1 uint8, skill2 uint8, skill3 uint8, skill4 uint8, startPoint types.Vector2[int16], movements []dto.MoveFragment)
	OnAttack(mapInstance *Map, character *Character, attackPayload dto.AttackPayload, skillLevel uint8)
	OnRangedAttack(mapInstance *Map, character *Character, attackPayload dto.AttackPayload, skillLevel uint8)
	OnMagicAttack(mapInstance *Map, character *Character, attackPayload dto.AttackPayload, skillLevel uint8)
	OnMobMobBuffApplied(mapInstance *Map, mob *Mob, buff constant.MobBuffFlag, value int32, skillID uint32, durationMs int64)
	OnMobMobBuffCancelled(mapInstance *Map, mob *Mob, buff constant.MobBuffFlag)
	OnMistSpawned(mapInstance *Map, mist *Mist)
	OnMistRemoved(mapInstance *Map, mist *Mist)
	OnDoorRemoved(mapInstance *Map, door *Door, animated bool)
}

type MobSpawn struct {
	Wz            *wz.MobSpawn
	Spawned       bool
	LastSpawnedAt time.Time
}

type Map struct {
	Wz                *wz.Map
	id                uint32
	objects           map[constant.ObjectType]map[uint32]Object
	controllerTable   *ControllerTable
	MobSpawns         map[uint32]*MobSpawn
	listener          MapListener
	sequence          uint32
	availableOIDs     []uint32
	Context           GameContext
	actorPID          *actor.PID
	pidMutex          sync.RWMutex
	UsedDoorPortalIDs map[uint8]struct{}
}

type BroadcastOption struct {
	SendRaw bool
}

func NewMap(id uint32, listener MapListener, mapId uint32, context GameContext) *Map {
	if listener == nil {
		panic("MapListener cannot be nil")
	}
	if context == nil {
		panic("GameContext cannot be nil")
	}

	wz, ok := context.GetResources().Maps[mapId]
	if !ok {
		panic(fmt.Sprintf("Wz not found for ID: %d", mapId))
	}

	mapInstance := &Map{
		id:              id,
		objects:         make(map[constant.ObjectType]map[uint32]Object),
		controllerTable: nil,
		MobSpawns:       make(map[uint32]*MobSpawn),
		listener:        listener,
		Wz:              wz,
		sequence:        0,
		availableOIDs:   make([]uint32, 0),
		Context:         context,
	}

	mapInstance.controllerTable = NewControllerTable(mapInstance.onMobControllerChange)

	mapInstance.initializeNpcs()
	mapInstance.initializeMobs()

	return mapInstance
}

func (m *Map) onMobControllerChange(mob *Mob, before *Character, after *Character) {

	m.listener.OnMobControllerChange(mob, before, after)
}

func (m *Map) allocateOID() uint32 {
	if len(m.availableOIDs) > 0 {

		oid := m.availableOIDs[0]
		m.availableOIDs = m.availableOIDs[1:]
		return oid
	}

	m.sequence++
	return m.sequence
}

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
	character.Stance = constant.StanceDefaultValue

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

	m.callMapLifecycleScript(character, "on_map_enter")

	return nil
}

func (m *Map) callMapLifecycleScript(character *Character, hook string) {
	if m == nil || character == nil || character.Context == nil {
		return
	}
	if character.GetMap() != m {
		return
	}
	pid := m.GetActorPID()
	if pid == nil {
		return
	}
	root := luax.GetRootLuaState(pid.String())
	if root == nil {
		return
	}
	_, thread, _ := luax.Call(root, "script/script.lua", hook, character, m)
	if thread != nil {
		thread.Close()
	}
}

func (m *Map) RemovePlayer(playerID uint32) error {
	if m.objects[constant.ObjectTypeCharacter] == nil {
		return fmt.Errorf("no players on map")
	}

	if _, exists := m.objects[constant.ObjectTypeCharacter][playerID]; !exists {
		return fmt.Errorf("player %d not found on map", playerID)
	}

	character := m.objects[constant.ObjectTypeCharacter][playerID].(*Character)

	m.callMapLifecycleScript(character, "on_map_leave")

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

func (m *Map) GetMapID() uint32 { return m.id }

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
	if m == nil || m.Wz == nil {
		return nil
	}
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
				Context:  m.Context,
				Map:      m,
			},
			Wz: &wz,
		}
		npc.ObjectCore.self = npc
		m.objects[constant.ObjectTypeNpc][oid] = npc
	}
}

func (m *Map) AddSummon(s *Summon) {
	if s == nil {
		return
	}
	if s.Owner == nil {
		panic("AddSummon: summon.Owner must not be nil")
	}
	if s.LifeCore.ObjectCore.self == nil {
		s.LifeCore.ObjectCore.self = s
	}
	if s.Map == m && m.objects[constant.ObjectTypeSummon] != nil {
		if existing, ok := m.objects[constant.ObjectTypeSummon][s.OID]; ok && existing == s {
			return
		}
	}
	if s.Map != nil && s.Map != m {
		return
	}
	s.ObjectCore.Context = m.Context
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

	s.Owner.Listener.OnSummonSpawn(s.Owner, s)
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

	s.Owner.Listener.OnSummonRemove(s.Owner, s, animated)

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
	if mist.ObjectCore.self == nil {
		mist.ObjectCore.self = mist
	}
	if mist.Map == m && m.objects[constant.ObjectTypeMist] != nil {
		if existing, ok := m.objects[constant.ObjectTypeMist][mist.OID]; ok && existing == mist {
			return
		}
	}
	if mist.Map != nil && mist.Map != m {
		return
	}
	mist.ObjectCore.Context = m.Context
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
	if door.ObjectCore.self == nil {
		door.ObjectCore.self = door
	}
	if door.Map == m && m.objects[constant.ObjectTypeDoor] != nil {
		if existing, ok := m.objects[constant.ObjectTypeDoor][door.OID]; ok && existing == door {
			return
		}
	}
	if door.Map != nil && door.Map != m {
		return
	}
	door.ObjectCore.Context = m.Context
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
	door.BroadcastCall(func(obj Object) {
		ch, ok := obj.(*Character)
		if !ok || ch == nil {
			return
		}
		door.SendSpawnSyncToViewer(ch)
	}, nil)
}

func (m *Map) RemoveDoor(oid uint32, animated bool) {
	m.removeDoorInternal(oid, animated, true)
}

func (m *Map) RemoveDoorByOwnerSkill(ownerID uint32, skillID constant.SkillID, animated bool) {
	if m == nil {
		return
	}
	for _, object := range m.GetObjects(constant.ObjectTypeDoor) {
		door, ok := object.(*Door)
		if !ok || door == nil || door.OID == 0 {
			continue
		}
		if door.OwnerID != ownerID || door.SkillID != skillID {
			continue
		}
		m.removeDoorInternal(door.OID, animated, false)
		return
	}
}

func (m *Map) removeDoorInternal(oid uint32, animated bool, notifyMysticCounterpart bool) {
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

	if ch := m.GetPlayer(door.OwnerID); ch != nil {
		ch.forgetDoorRegistrationIfSame(door)
	}

	var counterpartMapWZID uint32
	switch m.Wz.ID {
	case door.FieldMapID:
		counterpartMapWZID = door.ReturnMapID
	case door.ReturnMapID:
		counterpartMapWZID = door.FieldMapID
	}
	ownerID := door.OwnerID
	skillID := door.SkillID

	m.listener.OnDoorRemoved(m, door, animated)

	if door.SkillID == constant.SkillMysticDoor && m.Wz != nil && uint32(m.Wz.ID) == door.ReturnMapID {
		m.ReleaseMysticReturnPortal(door.ReturnPortalID)
	}

	delete(m.objects[constant.ObjectTypeDoor], oid)
	m.releaseOID(oid)
	if door.Map == m {
		door.Map = nil
		door.OID = 0
	}

	if notifyMysticCounterpart && skillID == constant.SkillMysticDoor && counterpartMapWZID != 0 && m.Context != nil {
		m.Context.NotifyDoorRemove(ownerID, uint32(skillID), counterpartMapWZID)
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

func (m *Map) FindDoorByOwner(ownerID uint32) *Door {
	if m == nil || ownerID == 0 {
		return nil
	}
	for _, o := range m.GetObjects(constant.ObjectTypeDoor) {
		d, ok := o.(*Door)
		if !ok || d == nil {
			continue
		}
		if d.OwnerID == ownerID {
			return d
		}
	}
	return nil
}

func (m *Map) ApplyPartyLeaveDoorSync(leaverID uint32) {
	if m == nil || leaverID == 0 {
		return
	}
	ch := m.GetPlayer(leaverID)
	if ch == nil {
		return
	}
	m.resyncOwnerDoorPortals(ch)
}

func (m *Map) ApplyPartyDisbandDoorSync(formerMemberIDs []uint32) {
	if m == nil || len(formerMemberIDs) == 0 {
		return
	}
	former := make(map[uint32]struct{}, len(formerMemberIDs))
	for _, id := range formerMemberIDs {
		if id != 0 {
			former[id] = struct{}{}
		}
	}
	for _, obj := range m.GetAllPlayers() {
		ch, ok := obj.(*Character)
		if !ok || ch == nil {
			continue
		}
		if _, in := former[ch.GetID()]; !in {
			continue
		}
		m.resyncOwnerDoorPortals(ch)
	}
}

func (m *Map) resyncOwnerDoorPortals(owner *Character) {
	if m == nil || owner == nil {
		return
	}
	oid := owner.GetID()
	for _, obj := range m.GetObjects(constant.ObjectTypeDoor) {
		d, ok := obj.(*Door)
		if !ok || d == nil || d.OwnerID != oid {
			continue
		}
		d.SendOwnerPortalResync(owner)
	}
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
			Context:  m.Context,
			Map:      m,
		},
		Wz: &npcSpawn,
	}
	npc.ObjectCore.self = npc

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

	mobSpec, ok := m.Context.GetResources().Monsters[mobId]
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
				Context:  m.Context,
				Map:      m,
			},
			hp:     uint32(max(0, mobSpec.MaxHP)),
			mp:     uint32(max(0, mobSpec.MaxMP)),
			BaseHp: uint32(max(0, mobSpec.MaxHP)),
			BaseMp: uint32(max(0, mobSpec.MaxMP)),
			Stance: 5,
		},
		Foothold: footholdID,
		Wz:       mobSpec,
		Spawn:    mobSpawn,
		ExpRate:  100,
		DropRate: 100,
		Homing:   make(map[uint32]*Homing),
	}
	mob.LifeCore.ObjectCore.self = mob

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
	mob.ClearAllHoming()
	delete(m.objects[constant.ObjectTypeMob], mobID)

	mob.ClearAllMobBuffTimers()
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

func (m *Map) GetObjects(filter constant.ObjectType) []Object {
	out := make([]Object, 0)
	for _, buckets := range m.objects {
		for _, obj := range buckets {
			if obj == nil {
				continue
			}
			if !obj.Is(filter) {
				continue
			}
			out = append(out, obj)
		}
	}
	return out
}

func (m *Map) Broadcast(message types.Packet, option *BroadcastOption) {
	if m.objects[constant.ObjectTypeCharacter] == nil {
		return
	}

	policy := types.SEND_POLICY_ENCRYPT
	if option != nil && option.SendRaw {
		policy = types.SEND_POLICY_RAW
	}

	for _, player := range m.objects[constant.ObjectTypeCharacter] {
		character, ok := player.(*Character)
		if !ok {
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
	if drop.ObjectCore != nil {
		drop.ObjectCore.self = mapObj
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

	meso := NewMeso(count, dropPoint, ownerID, dropType, oid, m.Context, m)

	drop := meso.GetDrop()
	if drop != nil {
		drop.RegisterExpire(constant.ITEM_EXPIRE_TIME)
		if dropType == constant.DROP_TYPE_OWNED || dropType == constant.DROP_TYPE_PARTY {
			drop.RegisterFFA(constant.ITEM_FFA_TIME)
		}
	}

	if m.objects[constant.ObjectTypeItem] == nil {
		m.objects[constant.ObjectTypeItem] = make(map[uint32]Object)
	}

	m.objects[constant.ObjectTypeItem][oid] = meso

	m.listener.OnMesoSpawned(m, meso)

	return meso, nil
}

func (m *Map) RemoveItem(itemID uint32, removeType constant.RemoveItemType, playerID uint32) error {
	if m.objects[constant.ObjectTypeItem] == nil {
		return fmt.Errorf("no items on map")
	}

	if _, exists := m.objects[constant.ObjectTypeItem][itemID]; !exists {
		return fmt.Errorf("item %d not found on map", itemID)
	}

	delete(m.objects[constant.ObjectTypeItem], itemID)

	m.releaseOID(itemID)

	m.listener.OnItemRemoved(m, itemID, playerID, removeType)

	return nil
}

func (m *Map) GetItem(itemID uint32) Item {
	if m.objects[constant.ObjectTypeItem] == nil {
		return nil
	}

	if item, ok := m.objects[constant.ObjectTypeItem][itemID].(Item); ok {
		return item
	}
	return nil
}

func (m *Map) GetItems() map[uint32]Object {
	if m.objects[constant.ObjectTypeItem] == nil {
		return make(map[uint32]Object)
	}
	return m.objects[constant.ObjectTypeItem]
}

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

		if drop.DropType == constant.DROP_TYPE_OWNED && drop.Owner != character.GetID() {
			return nil, constant.LOOT_FAILED_NO_OWNERSHIP
		}

		invenType := item.GetInventoryType()
		inven := character.Inventory[invenType]
		model := item.GetModel()

		if !inven.IsFree(model, item.GetCount()) {
			return nil, constant.LOOT_FAILED_INVENTORY_FULL
		}

		if err := m.RemoveItem(itemID, constant.REMOVE_ITEM_TYPE_ANIMATED, character.GetID()); err != nil {
			return nil, constant.LOOT_FAILED_INVALID_ITEM
		}

		return item, constant.LOOT_SUCCESS

	case *Meso:
		drop := item.GetDrop()
		if drop == nil {
			return nil, constant.LOOT_FAILED_INVALID_ITEM
		}

		if drop.DropType == constant.DROP_TYPE_OWNED && drop.Owner != character.GetID() {
			return nil, constant.LOOT_FAILED_NO_OWNERSHIP
		}

		mesoCount := item.GetCount32()
		cap := math.MaxInt32 - character.Meso
		if int32(mesoCount) > cap {
			return nil, constant.LOOT_FAILED_MESO_FULL
		}

		if err := m.RemoveItem(itemID, constant.REMOVE_ITEM_TYPE_ANIMATED, character.GetID()); err != nil {
			return nil, constant.LOOT_FAILED_INVALID_ITEM
		}

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
