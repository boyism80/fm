package entity

import (
	"fmt"
	"log"
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
	lua "github.com/yuin/gopher-lua"
)

type MapListener interface {
	OnPlayerAdded(ctx actor.Context, mapInstance *Map, character *Character, init bool)
	OnPlayerRemoved(mapInstance *Map, character *Character)
	OnPlayerMoved(mapInstance *Map, character *Character)
	OnPlayerMove(mapInstance *Map, character *Character, startPoint types.Vector2[int16], fragments []dto.MoveFragment)
	OnPlayerChat(mapInstance *Map, character *Character, message string)
	OnItemSpawned(mapInstance *Map, item Item, placement *FieldPlacement)
	OnMesoSpawned(mapInstance *Map, meso *Meso)
	OnItemRemoved(mapInstance *Map, itemID uint32, looterID uint32, mode constant.RemoveItemType)
	OnMobSpawned(mapInstance *Map, mob *Mob, spawnType constant.MobSpawnType, link uint32)
	OnMobRemoved(mapInstance *Map, mob *Mob, animationType constant.MobDieAnimationType)
	OnReactorSpawned(mapInstance *Map, reactor *Reactor)
	OnReactorRemoved(mapInstance *Map, reactor *Reactor)
	OnReactorTriggered(mapInstance *Map, reactor *Reactor, stance int32)
	OnMusicChanged(mapInstance *Map, song string)
	OnMapMessage(mapInstance *Map, messageType constant.ServerMessageType, message string)
	OnMobHomingRemoved(mapInstance *Map, mob *Mob, removed *Homing, causer *Character)
	OnMobHomingSet(mapInstance *Map, mob *Mob, homing *Homing, causer *Character)
	OnMobControllerChange(mob *Mob, before *Character, after *Character, aggro bool)
	OnMobMoved(mapInstance *Map, mob *Mob, isAggroed bool, centerSplit int8, skill1 uint8, skill2 uint8, skill3 uint8, skill4 uint8, startPoint types.Vector2[int16], movements []dto.MoveFragment)
	OnAttack(mapInstance *Map, character *Character, attackPayload dto.AttackPayload, skillLevel uint8)
	OnRangedAttack(mapInstance *Map, character *Character, attackPayload dto.AttackPayload, skillLevel uint8)
	OnMagicAttack(mapInstance *Map, character *Character, attackPayload dto.AttackPayload, skillLevel uint8)
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
	ReactorSpawns     map[uint32]*ReactorSpawn
	listener          MapListener
	mobListener       MobListener
	sequence          uint32
	availableOIDs     []uint32
	GameWorld         GameWorld
	actorPID          *actor.PID
	luaRoot           *lua.LState
	pidMutex          sync.RWMutex
	UsedDoorPortalIDs map[uint8]struct{}
}

type BroadcastOption struct {
	SendRaw bool
}

func NewMap(id uint32, listener MapListener, mobListener MobListener, mapId uint32, gw GameWorld) *Map {
	if listener == nil {
		panic("MapListener cannot be nil")
	}
	if mobListener == nil {
		panic("MobListener cannot be nil")
	}
	if gw == nil {
		panic("GameWorld cannot be nil")
	}

	wz, ok := gw.GetResources().Maps[mapId]
	if !ok {
		panic(fmt.Sprintf("Wz not found for ID: %d", mapId))
	}

	mapInstance := &Map{
		id:              id,
		objects:         make(map[constant.ObjectType]map[uint32]Object),
		controllerTable: nil,
		MobSpawns:       make(map[uint32]*MobSpawn),
		listener:        listener,
		mobListener:     mobListener,
		Wz:              wz,
		sequence:        0,
		availableOIDs:   make([]uint32, 0),
		GameWorld:       gw,
	}

	mapInstance.controllerTable = NewControllerTable(mapInstance.onMobControllerChange)

	mapInstance.initNpcs()
	mapInstance.initReactors()
	mapInstance.initMobs()

	return mapInstance
}

func (m *Map) GetLuaRoot() *lua.LState {
	if m == nil {
		return nil
	}
	return m.luaRoot
}

func (m *Map) EnsureLuaRoot(ctx actor.Context) *lua.LState {
	if m == nil {
		return nil
	}
	if m.luaRoot == nil {
		m.luaRoot = luax.NewState()
	}
	root := m.luaRoot

	if root == nil {
		return nil
	}
	return root
}

func (m *Map) ClearLuaRoot() {
	if m == nil {
		return
	}
	m.luaRoot = nil
}

func (m *Map) onMobControllerChange(mob *Mob, before *Character, after *Character, aggro bool) {

	m.listener.OnMobControllerChange(mob, before, after, aggro)
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

func OnMobControllerChange(mob *Mob, before *Character, after *Character, aggro bool) {
}

func (m *Map) AddPlayer(ctx actor.Context, playerID uint32, character *Character, spawnPoint uint8, init bool) error {
	if m.objects[constant.ObjectTypeCharacter] == nil {
		m.objects[constant.ObjectTypeCharacter] = make(map[uint32]Object)
	}

	character.Map = m
	if pos, ok := m.Wz.GetSpawnPosition(spawnPoint); ok {
		character.Position = pos
	}
	character.Stance = constant.StanceDefaultValue

	m.objects[constant.ObjectTypeCharacter][playerID] = character
	m.EnsureLuaRoot(ctx)

	m.listener.OnPlayerAdded(ctx, m, character, init)
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
	if m == nil || character == nil || character.GameWorld == nil {
		return
	}
	if character.GetMap() != m {
		return
	}
	root := m.GetLuaRoot()
	if root == nil {
		return
	}
	thread, err := luax.NewThread(root, constant.CharacterHookScriptPath)
	if err != nil {
		return
	}
	_, _ = luax.Call(thread, hook, character, m)
}

func (m *Map) RemovePlayer(playerID uint32) error {
	if m.objects[constant.ObjectTypeCharacter] == nil {
		return fmt.Errorf("no players on map")
	}

	if _, exists := m.objects[constant.ObjectTypeCharacter][playerID]; !exists {
		return fmt.Errorf("player %d not found on map", playerID)
	}

	character := m.objects[constant.ObjectTypeCharacter][playerID].(*Character)
	m.collectOwnedFieldDrops(character)

	m.callMapLifecycleScript(character, "on_map_leave")

	for _, summon := range character.GetSummons() {
		if summon == nil || summon.Map != m || summon.OID == 0 {
			continue
		}
		m.RemoveSummon(summon.OID, false)
	}

	delete(m.objects[constant.ObjectTypeCharacter], playerID)

	character.SuspendTimers()
	character.SetPartySearchConfig(nil)
	character.Map = nil
	m.controllerTable.LeavePlayer(character)

	m.listener.OnPlayerRemoved(m, character)

	return nil
}

func (m *Map) GetMapID() uint32 { return m.id }

func (m *Map) GetObject(objectType constant.ObjectType, id uint32) Object {
	if m == nil {
		return nil
	}
	bucket := m.objects[objectType]
	if bucket == nil {
		return nil
	}
	return bucket[id]
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

func (m *Map) GetPartyMembers(partyID uint32) []*Character {
	if m == nil {
		return nil
	}
	players := m.GetAllPlayers()
	if len(players) == 0 {
		return nil
	}
	out := make([]*Character, 0, len(players))
	for _, obj := range players {
		ch, ok := obj.(*Character)
		if !ok || ch == nil {
			continue
		}
		pid := ch.GetPartyID()
		if pid != nil && *pid == partyID {
			out = append(out, ch)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
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

func (m *Map) PointBelow(point types.Point[int16]) *types.Point[int16] {
	if m == nil || m.Wz == nil {
		return nil
	}
	return m.Wz.PointBelow(point)
}

func (m *Map) initNpcs() {
	if m.objects[constant.ObjectTypeNpc] == nil {
		m.objects[constant.ObjectTypeNpc] = make(map[uint32]Object)
	}

	for _, wz := range m.Wz.NpcSpawns {
		oid := m.allocateOID()
		npc := &Npc{
			ObjectCore: ObjectCore{
				OID:       oid,
				Position:  types.Point[int16]{X: wz.BaseSpawn.Position.X, Y: wz.BaseSpawn.Position.Y},
				GameWorld: m.GameWorld,
				Map:       m,
			},
			Wz: &wz,
		}
		npc.ObjectCore.self = npc
		npc.initTimers()
		m.objects[constant.ObjectTypeNpc][oid] = npc
	}
}

func (m *Map) GetReactor(oid uint32) *Reactor {
	if m.objects[constant.ObjectTypeReactor] == nil {
		return nil
	}
	if reactor, ok := m.objects[constant.ObjectTypeReactor][oid].(*Reactor); ok {
		return reactor
	}
	return nil
}

func (m *Map) GetReactors() map[uint32]Object {
	if m.objects[constant.ObjectTypeReactor] == nil {
		return make(map[uint32]Object)
	}
	return m.objects[constant.ObjectTypeReactor]
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
	s.ObjectCore.GameWorld = m.GameWorld
	if m.objects[constant.ObjectTypeSummon] == nil {
		m.objects[constant.ObjectTypeSummon] = make(map[uint32]Object)
	}
	if s.OID == 0 {
		s.OID = m.allocateOID()
	}
	if s.Map == nil {
		s.Map = m
	}
	s.initTimers()
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
	mist.ObjectCore.GameWorld = m.GameWorld
	if m.objects[constant.ObjectTypeMist] == nil {
		m.objects[constant.ObjectTypeMist] = make(map[uint32]Object)
	}
	if mist.OID == 0 {
		mist.OID = m.allocateOID()
	}
	if mist.Map == nil {
		mist.Map = m
	}
	mist.initTimers()
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
	door.ObjectCore.GameWorld = m.GameWorld
	if m.objects[constant.ObjectTypeDoor] == nil {
		m.objects[constant.ObjectTypeDoor] = make(map[uint32]Object)
	}
	if door.OID == 0 {
		door.OID = m.allocateOID()
	}
	if door.Map == nil {
		door.Map = m
	}
	door.initTimers()
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

	if notifyMysticCounterpart && skillID == constant.SkillMysticDoor && counterpartMapWZID != 0 && m.GameWorld != nil {
		if gw := m.GameWorld; gw != nil {
			gw.GetMapSystem().RemoveReturnDoor(ownerID, uint32(skillID), counterpartMapWZID)
		}
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

func (m *Map) initMobs() {
	for spawnId, spawnWz := range m.Wz.MobSpawns {
		m.MobSpawns[spawnId] = &MobSpawn{
			Wz:            &spawnWz,
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

	pointBelow := m.Wz.PointBelow(position)
	spawnPosition := position
	if pointBelow != nil {
		spawnPosition = *pointBelow
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
			OID:       oid,
			Position:  spawnPosition,
			GameWorld: m.GameWorld,
			Map:       m,
		},
		Wz: &npcSpawn,
	}
	npc.ObjectCore.self = npc
	npc.initTimers()

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

func (m *Map) SpawnMob(mobId uint32, position types.Point[int16], mobSpawn *MobSpawn, spawnType constant.MobSpawnType, link uint32) (*Mob, error) {
	oid := m.allocateOID()

	mobWz, ok := m.GameWorld.GetResources().Monsters[mobId]
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

	position.Y = position.Y - 1
	spawnPoint := position
	if pointBelow := m.Wz.PointBelow(position); pointBelow != nil {
		spawnPoint = *pointBelow
	}

	mob := &Mob{
		Listener: m.mobListener,
		LifeCore: LifeCore{
			ObjectCore: ObjectCore{
				OID:       oid,
				Position:  spawnPoint,
				GameWorld: m.GameWorld,
				Map:       m,
			},
			hp:     uint32(max(0, mobWz.MaxHP)),
			mp:     uint32(max(0, mobWz.MaxMP)),
			BaseHp: uint32(max(0, mobWz.MaxHP)),
			BaseMp: uint32(max(0, mobWz.MaxMP)),
			Stance: 5,
		},
		Foothold:  footholdID,
		Wz:        mobWz,
		Spawn:     mobSpawn,
		ExpRate:   100,
		DropRate:  100,
		Homing:    make(map[uint32]*Homing),
		accDamage: make(map[int64]map[uint32]uint64),
	}
	if mob.Listener == nil {
		panic("SpawnMob: mob listener must not be nil")
	}
	mob.Skills = NewMobSkillContainer(mob)
	mob.Buffs = NewMobBuffContainer(mob)
	mob.LifeCore.ObjectCore.self = mob
	mob.initTimers()
	if spawnType == constant.MobSpawnTypeFake {
		mob.Fake = true
	}

	if m.objects[constant.ObjectTypeMob] == nil {
		m.objects[constant.ObjectTypeMob] = make(map[uint32]Object)
	}

	m.objects[constant.ObjectTypeMob][oid] = mob
	m.listener.OnMobSpawned(m, mob, spawnType, link)
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
	mob.ClearTimers()
	delete(m.objects[constant.ObjectTypeMob], mobID)

	mob.Buffs.Clear()
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

func (m *Map) GetObjectsIn(filter constant.ObjectType, bounds types.Rect[int32]) []Object {
	out := make([]Object, 0)
	for _, obj := range m.GetObjects(filter) {
		if obj == nil {
			continue
		}
		pos := obj.GetPosition()
		point := types.Point[int32]{X: int32(pos.X), Y: int32(pos.Y)}
		if bounds.ContainsPoint(point) {
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

	fp := item.GetFieldPlacement()
	if fp == nil {
		return fmt.Errorf("item has no field placement")
	}

	dropPoint, ok := m.Wz.DropPoint(fp.Position)
	if !ok {
		dropPoint = fp.SpawnedPoint
	}

	fp.Position = dropPoint
	fp.OID = oid
	fp.Owner = ownerID
	fp.DropType = dropType
	fp.Map = m
	m.registerFieldDropTimers(fp, dropType)

	if m.objects[constant.ObjectTypeItem] == nil {
		m.objects[constant.ObjectTypeItem] = make(map[uint32]Object)
	}
	mapObj, ok := item.(Object)
	if !ok {
		return fmt.Errorf("item must implement Object")
	}
	if fp.ObjectCore != nil {
		fp.ObjectCore.self = mapObj
		fp.initTimers()
	}
	m.objects[constant.ObjectTypeItem][oid] = mapObj
	m.listener.OnItemSpawned(m, item, fp)

	owner := m.GetPlayer(ownerID)
	m.activateItemReactors(item, owner)

	return nil
}

func (m *Map) SpawnMeso(count int32, position types.Point[int16], ownerID uint32, dropType constant.DropType, playerDrop bool) (*Meso, error) {
	oid := m.allocateOID()

	dropPoint, ok := m.Wz.DropPoint(position)
	if !ok {
		dropPoint = position
	}

	meso := NewMeso(count, dropPoint, ownerID, dropType, oid, m.GameWorld, m)

	fp := meso.GetFieldPlacement()
	if fp != nil {
		fp.PlayerDrop = playerDrop
		m.registerFieldDropTimers(fp, dropType)
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

func (m *Map) LootItem(obj Object, character *Character, position types.Point[int16]) constant.LootResult {
	if m.objects[constant.ObjectTypeItem] == nil {
		return constant.LootFailedItemNotFound
	}

	switch item := obj.(type) {
	case Item:
		fp := item.GetFieldPlacement()
		if fp == nil {
			return constant.LootFailedInvalidItem
		}

		if !m.canLootFieldDrop(fp, character) {
			return constant.LootFailedNoOwnership
		}

		invenType := item.GetInventoryType()
		inven := character.Inventory[invenType]
		model := item.GetModel()

		if !inven.IsFree(model, item.GetCount()) {
			return constant.LootFailedInventoryFull
		}

		if _, err := character.AddItem(item, false); err != nil {
			log.Printf("Failed to add item: %v", err)
		}
		return constant.LootSuccess

	case *Meso:
		fp := item.GetFieldPlacement()
		if fp == nil {
			return constant.LootFailedInvalidItem
		}

		if !m.canLootFieldDrop(fp, character) {
			return constant.LootFailedNoOwnership
		}

		mesoCount := item.GetCount32()
		cap := math.MaxInt32 - character.Meso
		if int32(mesoCount) > cap {
			return constant.LootFailedMesoFull
		}

		character.GainMeso(int32(mesoCount))
		return constant.LootSuccess

	default:
		return constant.LootFailedInvalidItem
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
