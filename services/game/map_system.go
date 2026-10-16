package game

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/async"
	gameconst "github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

type MapSystem struct {
	gw             GameWorld
	server         *core.ServerCore
	mapListener    MapListener
	mobListener    MobListener
	mu             sync.RWMutex
	maps           map[uint32]*Map
	instances      map[uint32]*Map
	nextInstanceID atomic.Uint32
	slotMu         sync.Mutex
	slots          map[instanceSlot]*Map
}

func NewMapSystem(gw GameWorld, server *core.ServerCore, mapListener MapListener, mobListener MobListener) *MapSystem {
	return &MapSystem{
		gw:          gw,
		server:      server,
		mapListener: mapListener,
		mobListener: mobListener,
		maps:        make(map[uint32]*Map),
		instances:   make(map[uint32]*Map),
		slots:       make(map[instanceSlot]*Map),
	}
}

func (s *MapSystem) Start() error { return nil }
func (s *MapSystem) Stop()        {}

// PreCreate creates every fixed map with its actor, and the nil actor that handles packets between maps.
func (s *MapSystem) PreCreate() {
	actors := s.gw.Systems().Get[*c_actor.ActorRegistry]()
	s.subscribeDeadLetters()

	nilMapPID := actors.GetOrCreateActor("map_nil", actor.PropsFromProducer(func() actor.Actor {
		return s.gw.Systems().Get[ActorFactory]().NewMapActor(nil)
	}))
	s.server.SetNilActorPID(nilMapPID)

	mapListener := s.mapListener
	mobListener := s.mobListener
	for mapID := range s.gw.GetResources().Maps {
		name := fmt.Sprintf("map_%d", mapID)
		pid := actors.PredictPID(name)
		mapInstance := NewMap(mapID, mapListener, mobListener, mapID, s.gw, pid)
		s.maps[mapID] = mapInstance

		props := actor.PropsFromProducer(func() actor.Actor {
			return s.gw.Systems().Get[ActorFactory]().NewMapActor(mapInstance)
		})
		actual := actors.GetOrCreateActor(name, props)
		if !actual.Equal(pid) {
			panic(fmt.Sprintf("map %d actor pid mismatch", mapID))
		}
	}
	log.Printf("Pre-created %d map instances", len(s.maps))
}

func (s *MapSystem) Each(fn func(*Map)) {
	if s.gw == nil || fn == nil {
		return
	}
	s.mu.RLock()
	maps := make([]*Map, 0, len(s.maps)+len(s.instances))
	for _, m := range s.maps {
		if m != nil {
			maps = append(maps, m)
		}
	}
	for _, m := range s.instances {
		if m != nil {
			maps = append(maps, m)
		}
	}
	s.mu.RUnlock()
	for _, m := range maps {
		fn(m)
	}
}

func (s *MapSystem) Get(mapID uint32) *Map {
	if s.gw == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.maps[mapID]
}

func (s *MapSystem) Find(sm *StateMachine, mapID uint32) *Map {
	if m := sm.Map(mapID); m != nil {
		return m
	}
	return s.Get(mapID)
}

func (s *MapSystem) GetInstance(instanceKey uint32) *Map {
	if s.gw == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.instances[instanceKey]
}

// CreateInstanceMap returns an instance that closes when its last character leaves, or after InstanceLease if nobody ever enters.
func (s *MapSystem) CreateInstanceMap(templateID uint32, opts MapInitOpts) (*Map, error) {
	if s.gw == nil {
		return nil, fmt.Errorf("map system not ready")
	}
	if s.gw.GetResources() == nil || s.gw.GetResources().Maps[templateID] == nil {
		return nil, fmt.Errorf("template map %d not found", templateID)
	}
	opts.Instance = true
	key := s.nextInstanceID.Add(1)
	name := fmt.Sprintf("map_inst_%d", key)
	pid := s.gw.Systems().Get[*c_actor.ActorRegistry]().PredictPID(name)
	mapInstance := NewMapWithOpts(key, s.mapListener, s.mobListener, templateID, s.gw, pid, opts)

	props := actor.PropsFromProducer(func() actor.Actor {
		return s.gw.Systems().Get[ActorFactory]().NewMapActor(mapInstance)
	})
	actual := s.gw.Systems().Get[*c_actor.ActorRegistry]().GetOrCreateActor(name, props)
	if !actual.Equal(pid) {
		s.gw.Systems().Get[*c_actor.ActorRegistry]().StopActor(name, actual)
		return nil, fmt.Errorf("instance map %d actor pid mismatch", key)
	}

	s.mu.Lock()
	s.instances[key] = mapInstance
	s.mu.Unlock()

	if err := mapInstance.StartLease(); err != nil {
		return nil, err
	}
	return mapInstance, nil
}

func (s *MapSystem) CreateStateMachineMap(templateID uint32, sm *StateMachine, opts MapInitOpts) (*Map, error) {
	if s.gw.GetResources().Maps[templateID] == nil {
		return nil, fmt.Errorf("template map %d not found", templateID)
	}
	opts.Instance = true
	opts.StateMachine = sm
	key := s.nextInstanceID.Add(1)
	mapInstance := NewMapWithOpts(key, s.mapListener, s.mobListener, templateID, s.gw, sm.actorPID, opts)

	s.mu.Lock()
	s.instances[key] = mapInstance
	s.mu.Unlock()
	return mapInstance, nil
}

// RemoveInstanceMap runs once the instance has closed: no character is on it or on the way, and no owner holds it.
func (s *MapSystem) RemoveInstanceMap(instanceKey uint32) error {
	if s.gw == nil {
		return fmt.Errorf("map system not ready")
	}
	s.mu.Lock()
	m := s.instances[instanceKey]
	if m == nil {
		s.mu.Unlock()
		return nil
	}
	delete(s.instances, instanceKey)
	s.mu.Unlock()

	s.slotMu.Lock()
	for slot, slotMap := range s.slots {
		if slotMap == m {
			delete(s.slots, slot)
		}
	}
	s.slotMu.Unlock()

	sm := m.StateMachine()
	if sm != nil {
		s.gw.Systems().Get[*StateMachineSystem]().Send(sm.actorPID, &StateMachineMapRemoved{Map: m})
	}
	if m.GetPlayerCount() > 0 {
		return fmt.Errorf("instance map %d closed with %d players", instanceKey, m.GetPlayerCount())
	}
	if sm != nil {
		return nil
	}
	name := fmt.Sprintf("map_inst_%d", instanceKey)
	s.gw.Systems().Get[*c_actor.ActorRegistry]().PoisonActor(name, m.LogicActorPID())
	return nil
}

// CloseInstance refuses new arrivals, takes down its doors and sends everyone on the map to its return map;
// the map is removed when the last ref goes.
func (s *MapSystem) CloseInstance(m *Map) {
	if !m.IsInstance() {
		return
	}
	m.Close()

	exitMap := s.returnMap(m)
	s.Visit(m, func(ctx actor.Context) {
		for _, obj := range m.GetObjects(gameconst.ObjectTypeDoor) {
			if door, ok := obj.(*Door); ok {
				m.RemoveDoor(door.oid, true)
			}
		}

		if exitMap == nil {
			log.Printf("instance map %d (template %d) has no return map; players stay until they leave", m.GetMapID(), m.TemplateID())
			return
		}
		for _, obj := range m.GetAllPlayers() {
			ch, ok := obj.(*Character)
			if !ok || ch.GetMap() != m {
				continue
			}
			if err := s.Warp(ctx, ch, exitMap, 0, nil); err != nil {
				log.Printf("instance map %d: return character %d: %v", m.GetMapID(), ch.GetID(), err)
			}
		}
	})
}

func (s *MapSystem) Visit(m *Map, run func(ctx actor.Context)) {
	pid := m.LogicActorPID()
	if pid == nil {
		return
	}
	s.gw.Systems().Get[*c_actor.ActorRegistry]().Root().Send(pid, &MapCall{Run: func(ctx actor.Context, _ []*Map) []lua.LValue {
		run(ctx)
		return nil
	}})
}

func (s *MapSystem) Broadcast(m *Map, message types.Packet) {
	s.Visit(m, func(actor.Context) {
		m.Broadcast(message, nil)
	})
}

func (s *MapSystem) returnMap(m *Map) *Map {
	if m.wz == nil {
		return nil
	}
	exitID := uint32(m.wz.ReturnMapId)
	if m.wz.ForcedReturn != 0 {
		exitID = m.wz.ForcedReturn
	}
	if exitID == 0 || exitID == m.TemplateID() {
		return nil
	}
	return s.Get(exitID)
}

type instanceSlot struct {
	templateID uint32
	slot       uint32
}

// SlotInstance finds or creates the instance for (template, slot) under slotMutex, so callers asking for the same slot share one instance.
func (s *MapSystem) SlotInstance(templateID uint32, slot uint32) (*Map, error) {
	key := instanceSlot{templateID: templateID, slot: slot}

	s.slotMu.Lock()
	defer s.slotMu.Unlock()
	if found := s.slots[key]; found != nil && !found.Closing() {
		return found, nil
	}
	created, err := s.CreateInstanceMap(templateID, DefaultMapInitOpts())
	if err != nil {
		return nil, err
	}
	s.slots[key] = created
	return created, nil
}

func (s *MapSystem) Reset(L *lua.LState, mapInstance *Map, actorCtx actor.Context) int {
	if mapInstance == nil {
		if L != nil {
			L.Push(lua.LBool(false))
			return 1
		}
		return 0
	}
	if L == nil || s.gw == nil {
		return 0
	}
	targetPID := mapInstance.LogicActorPID()
	if targetPID == nil {
		L.Push(lua.LBool(false))
		return 1
	}
	return (luaMapCall{gw: s.gw}).InvokeAwait(L, actorCtx, targetPID, func(ctx actor.Context, _ []*Map) []lua.LValue {
		mapInstance.Reset()
		return []lua.LValue{lua.LBool(true)}
	})
}

func (s *MapSystem) Respawn(L *lua.LState, mapInstance *Map, actorCtx actor.Context, includeNegativeMobTime bool) int {
	if mapInstance == nil {
		if L != nil {
			L.Push(lua.LNumber(0))
			return 1
		}
		return 0
	}
	if L == nil || s.gw == nil {
		return 0
	}
	targetPID := mapInstance.LogicActorPID()
	if targetPID == nil {
		L.Push(lua.LNumber(0))
		return 1
	}
	return (luaMapCall{gw: s.gw}).InvokeAwait(L, actorCtx, targetPID, func(ctx actor.Context, _ []*Map) []lua.LValue {
		return []lua.LValue{lua.LNumber(mapInstance.Respawn(includeNegativeMobTime))}
	})
}

func pushRunOnMapResult(L *lua.LState, ok bool, result lua.LValue, errMsg string) int {
	if !ok {
		L.Push(lua.LBool(false))
		L.Push(lua.LNil)
		if errMsg == "" {
			L.Push(lua.LNil)
		} else {
			L.Push(lua.LString(errMsg))
		}
		return 3
	}
	L.Push(lua.LBool(true))
	if result == nil {
		L.Push(lua.LNil)
	} else {
		L.Push(result)
	}
	L.Push(lua.LNil)
	return 3
}

func (s *MapSystem) RunScript(L *lua.LState, actorCtx actor.Context, targetMap *Map, scriptPath string, funcName string, args []interface{}) int {
	if L == nil {
		return 0
	}
	if targetMap == nil {
		return pushRunOnMapResult(L, false, nil, "run_on_map: target map not found")
	}
	targetPID := targetMap.LogicActorPID()
	if targetPID == nil {
		return pushRunOnMapResult(L, false, nil, "run_on_map: target map actor not found")
	}
	if s.gw == nil {
		return pushRunOnMapResult(L, false, nil, "run_on_map: game server not found")
	}
	return (luaMapCall{gw: s.gw}).InvokeAwaitAsync(L, actorCtx, targetPID, func(ctx actor.Context, _ []*Map) *async.Promise[[]lua.LValue] {
		return targetMap.RunScript(ctx, scriptPath, funcName, args).Then(func(vals []lua.LValue) ([]lua.LValue, error) {
			result := lua.LNil
			if len(vals) > 0 && vals[0] != nil {
				result = vals[0]
			}
			return []lua.LValue{lua.LBool(true), result, lua.LNil}, nil
		})
	})
}

func (s *MapSystem) Warp(actorCtx actor.Context, character *Character, targetMap *Map, spawnPoint uint8, onEnter func(actor.Context)) error {
	if targetMap == nil {
		return fmt.Errorf("target map is nil")
	}
	if character == nil {
		return fmt.Errorf("character is nil")
	}
	if sm := targetMap.StateMachine(); sm != nil && !sm.Admits(character) {
		return ErrStateMachineEntryDenied
	}
	// The ticket is released on every failure below, or carried to the target actor and released after AddPlayer.
	ticket, err := targetMap.Reserve()
	if err != nil {
		return err
	}
	currentMap := character.GetMap()
	owner := currentMap
	if owner == nil {
		owner = character.Destination()
	}
	if owner == nil {
		ticket.Release()
		return fmt.Errorf("character %d is on no map", character.GetID())
	}
	ownerPID := owner.LogicActorPID()
	if currentMap == nil || actorCtx == nil || !ownerPID.Equal(actorCtx.Self()) {
		s.gw.Systems().Get[*c_actor.ActorRegistry]().Root().Send(ownerPID, &HandoffCharacter{
			Character: character,
			TargetMap: targetMap,
			Portal:    spawnPoint,
			OnEnter:   onEnter,
			Ticket:    ticket,
		})
		return nil
	}

	targetPID := targetMap.LogicActorPID()
	if ownerPID.Equal(targetPID) {
		if err := currentMap.RemovePlayer(character.GetID()); err != nil {
			ticket.Release()
			return err
		}
		character.Account().Storage().Close()
		err := targetMap.AddPlayer(actorCtx, character.GetID(), character, spawnPoint, false)
		ticket.Release()
		if err != nil {
			return err
		}
		character.ResumeTimers(actorCtx.Self())
		character.listener.OnPartyMemberFieldsChanged(character)
		character.Wedding().NotifySpouseMap(actorCtx, false)
		if onEnter != nil {
			// Deferred so a Lua caller finishes its own call before onEnter starts on the same root.
			actorCtx.Send(targetPID, &MapCall{Run: func(ctx actor.Context, _ []*Map) []lua.LValue {
				onEnter(ctx)
				return nil
			}})
		}
		return nil
	}
	if !character.BeginMove(targetMap) {
		ticket.Release()
		return fmt.Errorf("character %d is already moving", character.GetID())
	}
	if err := currentMap.RemovePlayer(character.GetID()); err != nil {
		character.FinishMove(targetMap)
		ticket.Release()
		return err
	}
	character.Account().Storage().Close()
	actorCtx.Send(targetPID, &WarpCharacter{
		Character: character,
		TargetMap: targetMap,
		Portal:    spawnPoint,
		OnEnter:   onEnter,
		Ticket:    ticket,
	})
	return nil
}

func (s *MapSystem) CreateReturnDoor(ch *Character, key DoorKey, skillID gameconst.SkillID) *Map {
	if s.gw == nil || ch == nil {
		return nil
	}
	m := ch.GetMap()
	if m == nil || m.wz == nil {
		return nil
	}
	destMapID := uint32(m.wz.ReturnMapId)
	if destMapID == 0 || destMapID == uint32(m.wz.ID) {
		return nil
	}
	destMap := s.Get(destMapID)
	if destMap == nil {
		ch.listener.OnMessage(ch, gameconst.MsgPinkText, gameconst.DoorNoTownPortalMessage)
		return nil
	}
	destPID := destMap.LogicActorPID()
	if destPID == nil {
		ch.listener.OnMessage(ch, gameconst.MsgPinkText, gameconst.DoorNoTownPortalMessage)
		return nil
	}
	srcPID := m.LogicActorPID()
	if srcPID == nil {
		ch.listener.OnMessage(ch, gameconst.MsgPinkText, gameconst.DoorNoTownPortalMessage)
		return nil
	}
	root := s.gw.Systems().Get[*c_actor.ActorRegistry]().Root()
	if root == nil {
		return nil
	}
	fieldAnchorPt := types.Point[int16]{X: ch.position.X, Y: ch.position.Y}
	var closestPortalID uint8
	if id, ok := m.wz.FindClosestDoorReturnPortalSpawnID(fieldAnchorPt); ok {
		closestPortalID = id
	} else {
		closestPortalID = m.wz.FindClosestPortalSpawnID(fieldAnchorPt)
	}
	slot := 0
	slot = s.gw.Systems().Get[*PartySystem]().PartyMemberIndex(ch.GetID(), ch.party.ID())
	root.Send(destPID, &RequestSpawnDoor{
		ReplyTo:     srcPID,
		TargetMap:   destMap,
		Key:         key,
		CharacterID: ch.GetID(),
		OwnerID:     ch.GetID(),
		SkillID:     skillID,
		Field: DoorEndpoint{
			Map:      m,
			PortalID: closestPortalID,
			Position: ch.position,
		},
		PartyOwnerSlot: slot,
		PartyID:        ch.party.ID(),
	})
	return destMap
}

// DespawnDoor removes the door on its map's own actor; callers on another actor must not touch that map directly.
func (s *MapSystem) DespawnDoor(m *Map, key DoorKey, animated bool, notifyCounterpart bool) {
	if m == nil {
		return
	}
	pid := m.LogicActorPID()
	if pid == nil {
		return
	}
	s.gw.Systems().Get[*c_actor.ActorRegistry]().Root().Send(pid, &DespawnDoor{
		Map:               m,
		Key:               key,
		Animated:          animated,
		NotifyCounterpart: notifyCounterpart,
	})
}

func (s *MapSystem) SetDoorPartyID(m *Map, key DoorKey, partyID *uint32) {
	if m == nil {
		return
	}
	pid := m.LogicActorPID()
	if pid == nil {
		return
	}
	s.gw.Systems().Get[*c_actor.ActorRegistry]().Root().Send(pid, &SetDoorPartyID{
		Map:     m,
		Key:     key,
		PartyID: partyID,
	})
}
