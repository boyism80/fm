package server

import (
	"fmt"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	g_actor "github.com/boyism80/fm/services/game/actor"
	gameconst "github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

type mapSystem struct {
	gs *GameServer
}

type schedulerSystem struct {
	gs *GameServer
}

type partySystem struct {
	gs *GameServer
}

type guildSystem struct {
	gs *GameServer
}

type allianceSystem struct {
	gs *GameServer
}

type dispatchSystem struct {
	gs *GameServer
}

func (gs *GameServer) GetMapSystem() entity.MapSystem {
	if gs == nil {
		return nil
	}
	return gs.mapSystem
}

func (gs *GameServer) GetSchedulerSystem() entity.SchedulerSystem {
	if gs == nil {
		return nil
	}
	return gs.schedulerSystem
}

func (gs *GameServer) GetPartySystem() entity.PartySystem {
	if gs == nil {
		return nil
	}
	return gs.partySystem
}

func (gs *GameServer) GetGuildSystem() entity.GuildSystem {
	if gs == nil {
		return nil
	}
	return gs.guildSystem
}

func (gs *GameServer) GetAllianceSystem() entity.AllianceSystem {
	if gs == nil {
		return nil
	}
	return gs.allianceSystem
}

func (gs *GameServer) GetDispatchSystem() entity.DispatchSystem {
	if gs == nil {
		return nil
	}
	return gs.dispatchSystem
}

func (gs *GameServer) GetStateMachineRegistry() entity.StateMachineRegistry {
	if gs == nil {
		return nil
	}
	return gs.stateMachines
}

func (gs *GameServer) GetCarnivalRegistry() *entity.CarnivalRegistry {
	if gs == nil {
		return nil
	}
	return gs.carnivalRegistry
}

func (s mapSystem) Visit(fn func(*entity.Map)) {
	if s.gs == nil || fn == nil {
		return
	}
	s.gs.mapsMutex.RLock()
	maps := make([]*entity.Map, 0, len(s.gs.maps)+len(s.gs.instanceMaps))
	for _, m := range s.gs.maps {
		if m != nil {
			maps = append(maps, m)
		}
	}
	for _, m := range s.gs.instanceMaps {
		if m != nil {
			maps = append(maps, m)
		}
	}
	s.gs.mapsMutex.RUnlock()
	for _, m := range maps {
		fn(m)
	}
}

func (s mapSystem) Get(mapID uint32) *entity.Map {
	if s.gs == nil {
		return nil
	}
	s.gs.mapsMutex.RLock()
	defer s.gs.mapsMutex.RUnlock()
	return s.gs.maps[mapID]
}

func (s mapSystem) GetInstance(instanceKey uint32) *entity.Map {
	if s.gs == nil {
		return nil
	}
	s.gs.mapsMutex.RLock()
	defer s.gs.mapsMutex.RUnlock()
	return s.gs.instanceMaps[instanceKey]
}

// CreateInstanceMap returns an instance that closes when its last character leaves, or after InstanceLease if nobody ever enters.
func (s mapSystem) CreateInstanceMap(templateID uint32, opts entity.MapInitOpts) (*entity.Map, error) {
	if s.gs == nil {
		return nil, fmt.Errorf("map system not ready")
	}
	if s.gs.resources == nil || s.gs.resources.Maps[templateID] == nil {
		return nil, fmt.Errorf("template map %d not found", templateID)
	}
	if s.gs.mapListener == nil || s.gs.mobListener == nil {
		return nil, fmt.Errorf("map listeners not ready")
	}
	opts.Instance = true
	key := s.gs.nextInstanceID.Add(1)
	name := fmt.Sprintf("map_inst_%d", key)
	pid := s.gs.actorRegistry.PredictPID(name)
	mapInstance := entity.NewMapWithOpts(key, s.gs.mapListener, s.gs.mobListener, templateID, s.gs, pid, opts)

	props := actor.PropsFromProducer(func() actor.Actor {
		return g_actor.NewMapActor(mapInstance, s.gs)
	})
	actual := s.gs.actorRegistry.GetOrCreateActor(name, props)
	if !actual.Equal(pid) {
		s.gs.actorRegistry.StopActor(name, actual)
		return nil, fmt.Errorf("instance map %d actor pid mismatch", key)
	}

	s.gs.mapsMutex.Lock()
	s.gs.instanceMaps[key] = mapInstance
	s.gs.mapsMutex.Unlock()

	if err := mapInstance.StartLease(); err != nil {
		return nil, err
	}
	return mapInstance, nil
}

// RemoveInstanceMap runs once the instance has closed: no character is on it or on the way, and no owner holds it.
func (s mapSystem) RemoveInstanceMap(instanceKey uint32) error {
	if s.gs == nil {
		return fmt.Errorf("map system not ready")
	}
	s.gs.mapsMutex.Lock()
	m := s.gs.instanceMaps[instanceKey]
	if m == nil {
		s.gs.mapsMutex.Unlock()
		return nil
	}
	delete(s.gs.instanceMaps, instanceKey)
	s.gs.mapsMutex.Unlock()

	s.gs.slotMutex.Lock()
	for slot, slotMap := range s.gs.slotInstances {
		if slotMap == m {
			delete(s.gs.slotInstances, slot)
		}
	}
	s.gs.slotMutex.Unlock()

	if m.GetPlayerCount() > 0 {
		return fmt.Errorf("instance map %d closed with %d players", instanceKey, m.GetPlayerCount())
	}
	name := fmt.Sprintf("map_inst_%d", instanceKey)
	s.gs.actorRegistry.PoisonActor(name, m.HomeActorPID())
	return nil
}

// CloseInstance refuses new arrivals, takes down its doors and sends everyone on the map to its return map;
// the map is removed when the last ref goes.
func (s mapSystem) CloseInstance(m *entity.Map) {
	if m.IsInstance() == false {
		return
	}
	m.Close()

	pid := m.LogicActorPID()
	if pid == nil {
		return
	}
	exitMap := s.returnMap(m)
	s.gs.GetRootContext().Send(pid, &g_actor.MapCall{Run: func(ctx actor.Context, _ *g_actor.GameLogicActor) []lua.LValue {
		for _, obj := range m.GetObjects(gameconst.ObjectTypeDoor) {
			if door, ok := obj.(*entity.Door); ok {
				m.RemoveDoor(door.OID, true)
			}
		}

		if exitMap == nil {
			log.Printf("instance map %d (template %d) has no return map; players stay until they leave", m.GetMapID(), m.TemplateID())
			return nil
		}
		for _, obj := range m.GetAllPlayers() {
			ch, ok := obj.(*entity.Character)
			if ok == false || ch.GetMap() != m {
				continue
			}
			if err := s.Warp(ctx, ch, exitMap, 0, nil); err != nil {
				log.Printf("instance map %d: return character %d: %v", m.GetMapID(), ch.GetID(), err)
			}
		}
		return nil
	}})
}

func (s mapSystem) returnMap(m *entity.Map) *entity.Map {
	if m.Wz == nil {
		return nil
	}
	exitID := uint32(m.Wz.ReturnMapId)
	if m.Wz.HasForcedReturn() {
		exitID = uint32(m.Wz.ForcedReturn)
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
func (s mapSystem) SlotInstance(templateID uint32, slot uint32) (*entity.Map, error) {
	key := instanceSlot{templateID: templateID, slot: slot}

	s.gs.slotMutex.Lock()
	defer s.gs.slotMutex.Unlock()
	if found := s.gs.slotInstances[key]; found != nil && found.Closing() == false {
		return found, nil
	}
	created, err := s.CreateInstanceMap(templateID, entity.DefaultMapInitOpts())
	if err != nil {
		return nil, err
	}
	s.gs.slotInstances[key] = created
	return created, nil
}

func (s mapSystem) Reset(L *lua.LState, mapInstance *entity.Map, actorCtx actor.Context) int {
	if mapInstance == nil {
		if L != nil {
			L.Push(lua.LBool(false))
			return 1
		}
		return 0
	}
	if L == nil || s.gs == nil {
		return 0
	}
	targetPID := mapInstance.LogicActorPID()
	if targetPID == nil {
		L.Push(lua.LBool(false))
		return 1
	}
	return (luaMapCall{gs: s.gs}).InvokeAwait(L, actorCtx, targetPID, func(ctx actor.Context, a *g_actor.GameLogicActor) []lua.LValue {
		mapInstance.Reset()
		return []lua.LValue{lua.LBool(true)}
	})
}

func (s mapSystem) Respawn(L *lua.LState, mapInstance *entity.Map, actorCtx actor.Context, includeNegativeMobTime bool) int {
	if mapInstance == nil {
		if L != nil {
			L.Push(lua.LNumber(0))
			return 1
		}
		return 0
	}
	if L == nil || s.gs == nil {
		return 0
	}
	targetPID := mapInstance.LogicActorPID()
	if targetPID == nil {
		L.Push(lua.LNumber(0))
		return 1
	}
	return (luaMapCall{gs: s.gs}).InvokeAwait(L, actorCtx, targetPID, func(ctx actor.Context, a *g_actor.GameLogicActor) []lua.LValue {
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

func (s mapSystem) RunScript(L *lua.LState, actorCtx actor.Context, mapID uint32, scriptPath string, funcName string, args []interface{}) int {
	if L == nil {
		return 0
	}
	targetMap := s.Get(mapID)
	if targetMap == nil {
		return pushRunOnMapResult(L, false, nil, "run_on_map: target map not found")
	}
	targetPID := targetMap.LogicActorPID()
	if targetPID == nil {
		return pushRunOnMapResult(L, false, nil, "run_on_map: target map actor not found")
	}
	if s.gs == nil {
		return pushRunOnMapResult(L, false, nil, "run_on_map: game server not found")
	}
	return (luaMapCall{gs: s.gs}).InvokeAwaitAsync(L, actorCtx, targetPID, func(ctx actor.Context, a *g_actor.GameLogicActor) *async.Promise {
		return targetMap.RunScript(ctx, scriptPath, funcName, args).Then(func(v interface{}) (interface{}, error) {
			vals := luax.ResultValues(v)
			result := lua.LNil
			if len(vals) > 0 && vals[0] != nil {
				result = vals[0]
			}
			return []lua.LValue{lua.LBool(true), result, lua.LNil}, nil
		})
	})
}

func (s mapSystem) Warp(actorCtx actor.Context, character *entity.Character, targetMap *entity.Map, spawnPoint uint8, onEnter func(actor.Context)) error {
	if targetMap == nil {
		return fmt.Errorf("target map is nil")
	}
	if character == nil {
		return fmt.Errorf("character is nil")
	}
	// The ticket is released on every failure below, or carried to the target actor and released after AddPlayer.
	ticket, err := targetMap.Reserve()
	if err != nil {
		return err
	}
	currentMap := character.GetMap()
	targetPID := targetMap.LogicActorPID()
	if targetPID == nil {
		ticket.Release()
		return fmt.Errorf("target map actor not found")
	}
	if currentMap == nil {
		character.Destination = targetMap
		s.gs.GetRootContext().Send(targetPID, &g_actor.WarpCharacter{
			Character: character,
			TargetMap: targetMap,
			Portal:    spawnPoint,
			OnEnter:   onEnter,
			Ticket:    ticket,
		})
		return nil
	}
	sourcePID := currentMap.LogicActorPID()
	if sourcePID == nil {
		ticket.Release()
		return fmt.Errorf("source map actor not found")
	}
	if sourcePID.Equal(targetPID) {
		if actorCtx == nil || actorCtx.Self() == nil || !actorCtx.Self().Equal(sourcePID) {
			ticket.Release()
			return fmt.Errorf("same-owner warp must run on owner actor")
		}
		if err := currentMap.RemovePlayer(character.GetID()); err != nil {
			ticket.Release()
			return err
		}
		err := targetMap.AddPlayer(actorCtx, character.GetID(), character, spawnPoint, false)
		ticket.Release()
		if err != nil {
			return err
		}
		character.ResumeTimers(actorCtx.Self())
		character.Listener.OnPartyMemberFieldsChanged(character)
		if onEnter != nil {
			// Deferred so a Lua caller finishes its own call before onEnter starts on the same root.
			actorCtx.Send(targetPID, &g_actor.MapCall{Run: func(ctx actor.Context, _ *g_actor.GameLogicActor) []lua.LValue {
				onEnter(ctx)
				return nil
			}})
		}
		return nil
	}
	if actorCtx != nil && actorCtx.Self() != nil && actorCtx.Self().Equal(sourcePID) {
		character.Destination = targetMap
		if err := currentMap.RemovePlayer(character.GetID()); err != nil {
			character.Destination = nil
			ticket.Release()
			return err
		}
		actorCtx.Send(targetPID, &g_actor.WarpCharacter{
			Character: character,
			TargetMap: targetMap,
			Portal:    spawnPoint,
			OnEnter:   onEnter,
			Ticket:    ticket,
		})
		return nil
	}
	s.gs.GetRootContext().Send(sourcePID, &g_actor.HandoffCharacter{
		Character: character,
		TargetMap: targetMap,
		Portal:    spawnPoint,
		OnEnter:   onEnter,
		Ticket:    ticket,
	})
	return nil
}

func (s mapSystem) CreateReturnDoor(ch *entity.Character, key entity.DoorKey, skillID gameconst.SkillID) *entity.Map {
	if s.gs == nil || ch == nil {
		return nil
	}
	m := ch.GetMap()
	if m == nil || m.Wz == nil {
		return nil
	}
	destMapID := uint32(m.Wz.ReturnMapId)
	if destMapID == 0 || destMapID == uint32(m.Wz.ID) {
		return nil
	}
	destMap := s.Get(destMapID)
	if destMap == nil {
		ch.Listener.OnMessage(ch, gameconst.MsgPinkText, gameconst.DoorNoTownPortalMessage)
		return nil
	}
	destPID := destMap.LogicActorPID()
	if destPID == nil {
		ch.Listener.OnMessage(ch, gameconst.MsgPinkText, gameconst.DoorNoTownPortalMessage)
		return nil
	}
	srcPID := m.LogicActorPID()
	if srcPID == nil {
		ch.Listener.OnMessage(ch, gameconst.MsgPinkText, gameconst.DoorNoTownPortalMessage)
		return nil
	}
	root := s.gs.GetRootContext()
	if root == nil {
		return nil
	}
	fieldAnchorPt := types.Point[int16]{X: ch.Position.X, Y: ch.Position.Y}
	var closestPortalID uint8
	if id, ok := m.Wz.FindClosestDoorReturnPortalSpawnID(fieldAnchorPt); ok {
		closestPortalID = id
	} else {
		closestPortalID = m.Wz.FindClosestPortalSpawnID(fieldAnchorPt)
	}
	slot := 0
	slot = s.gs.party.PartyMemberIndex(ch.GetID(), ch.GetPartyID())
	root.Send(destPID, &g_actor.RequestSpawnDoor{
		ReplyTo:     srcPID,
		TargetMap:   destMap,
		Key:         key,
		CharacterID: ch.GetID(),
		OwnerID:     ch.GetID(),
		SkillID:     skillID,
		Field: entity.DoorEndpoint{
			Map:      m,
			PortalID: closestPortalID,
			Position: ch.Position,
		},
		PartyOwnerSlot: slot,
		PartyID:        ch.GetPartyID(),
	})
	return destMap
}

// DespawnDoor removes the door on its map's own actor; callers on another actor must not touch that map directly.
func (s mapSystem) DespawnDoor(m *entity.Map, key entity.DoorKey, animated bool, notifyCounterpart bool) {
	if m == nil {
		return
	}
	pid := m.LogicActorPID()
	if pid == nil {
		return
	}
	s.gs.GetRootContext().Send(pid, &g_actor.DespawnDoor{
		Map:               m,
		Key:               key,
		Animated:          animated,
		NotifyCounterpart: notifyCounterpart,
	})
}

func (s schedulerSystem) RunObjectTimer(pid *actor.PID, obj entity.Object, key string) {
	if s.gs == nil || pid == nil || obj == nil || key == "" {
		return
	}
	payload := &c_actor.RunObjectTimer{
		ObjectType: obj.GetObjectType(),
		ID:         obj.GetPK(),
		Key:        key,
	}
	if root := s.gs.GetRootContext(); root != nil {
		root.Send(pid, payload)
	}
}

func (s schedulerSystem) RunReactorRespawn(pid *actor.PID, mapID uint32, spawnID uint32) {
	if s.gs == nil || pid == nil {
		return
	}
	payload := &c_actor.RunReactorRespawn{
		MapID:   mapID,
		SpawnID: spawnID,
	}
	if root := s.gs.GetRootContext(); root != nil {
		root.Send(pid, payload)
	}
}

func (s partySystem) Get(partyID uint32) *entity.Party {
	if s.gs == nil {
		return nil
	}
	return s.gs.party.Get(partyID)
}

func (s guildSystem) Get(guildID uint32) *entity.Guild {
	if s.gs == nil {
		return nil
	}
	return s.gs.guild.Get(guildID)
}

func (s guildSystem) TrySetAllianceInvite(guildID, allianceID uint32, expiresAt time.Time) bool {
	if s.gs == nil {
		return false
	}
	return s.gs.guild.TrySetAllianceInvite(guildID, allianceID, expiresAt)
}

func (s allianceSystem) Get(allianceID uint32) *entity.Alliance {
	if s.gs == nil {
		return nil
	}
	return s.gs.alliance.Get(allianceID)
}

func (s dispatchSystem) SendTo(characterID uint32, msg interface{}) {
	if s.gs == nil || characterID == 0 || msg == nil {
		return
	}
	s.gs.EnsureSend(nil, characterID, msg)
}
