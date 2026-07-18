package server

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	c_actor "github.com/boyism80/fm/core/actor"
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

func (s mapSystem) Get(mapID uint32) *entity.Map {
	if s.gs == nil {
		return nil
	}
	s.gs.mapsMutex.RLock()
	defer s.gs.mapsMutex.RUnlock()
	return s.gs.maps[mapID]
}

func (s mapSystem) ResetFromLua(L *lua.LState, mapInstance *entity.Map, actorCtx actor.Context) int {
	if mapInstance == nil {
		if L != nil {
			L.Push(lua.LBool(false))
			return 1
		}
		return 0
	}
	targetPID := mapInstance.GetActorPID()
	if targetPID == nil {
		L.Push(lua.LBool(false))
		return 1
	}
	cfg, _ := luax.GetConfiguration(L)
	callerPID := cfg.MapActorPID
	if actorCtx != nil {
		callerPID = actorCtx.Self()
	}
	if callerPID != nil && callerPID.Equal(targetPID) {
		mapInstance.Reset()
		L.Push(lua.LBool(true))
		return 1
	}
	if callerPID == nil || s.gs == nil {
		L.Push(lua.LBool(false))
		return 1
	}
	root := L.Parent
	if root == nil {
		L.Push(lua.LBool(false))
		return 1
	}
	s.gs.GetRootContext().Send(targetPID, &g_actor.ResetMap{
		ReplyTo: callerPID,
		Root:    root,
		Thread:  L,
	})
	return L.Yield(lua.LNil)
}

func (s mapSystem) RespawnFromLua(L *lua.LState, mapInstance *entity.Map, actorCtx actor.Context, includeNegativeMobTime bool) int {
	if mapInstance == nil {
		if L != nil {
			L.Push(lua.LNumber(0))
			return 1
		}
		return 0
	}
	targetPID := mapInstance.GetActorPID()
	if targetPID == nil {
		L.Push(lua.LNumber(0))
		return 1
	}
	cfg, _ := luax.GetConfiguration(L)
	callerPID := cfg.MapActorPID
	if actorCtx != nil {
		callerPID = actorCtx.Self()
	}
	if callerPID != nil && callerPID.Equal(targetPID) {
		L.Push(lua.LNumber(mapInstance.Respawn(includeNegativeMobTime)))
		return 1
	}
	if callerPID == nil || s.gs == nil {
		L.Push(lua.LNumber(0))
		return 1
	}
	root := L.Parent
	if root == nil {
		L.Push(lua.LNumber(0))
		return 1
	}
	s.gs.GetRootContext().Send(targetPID, &g_actor.RespawnMap{
		ReplyTo:                callerPID,
		Root:                   root,
		Thread:                 L,
		IncludeNegativeMobTime: includeNegativeMobTime,
	})
	return L.Yield(lua.LNil)
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

func (s mapSystem) RunOnMapFromLua(L *lua.LState, actorCtx actor.Context, mapID uint32, scriptPath string, funcName string, args []interface{}) int {
	if L == nil {
		return 0
	}
	targetMap := s.Get(mapID)
	if targetMap == nil {
		return pushRunOnMapResult(L, false, nil, "run_on_map: target map not found")
	}
	targetPID := targetMap.GetActorPID()
	if targetPID == nil {
		return pushRunOnMapResult(L, false, nil, "run_on_map: target map actor not found")
	}
	cfg, _ := luax.GetConfiguration(L)
	callerPID := cfg.MapActorPID
	if actorCtx != nil {
		callerPID = actorCtx.Self()
	}
	if callerPID == nil {
		return pushRunOnMapResult(L, false, nil, "run_on_map: caller actor not found")
	}
	if s.gs == nil {
		return pushRunOnMapResult(L, false, nil, "run_on_map: game server not found")
	}
	root := L.Parent
	if root == nil {
		return pushRunOnMapResult(L, false, nil, "run_on_map: root lua state not found")
	}
	s.gs.GetRootContext().Send(targetPID, &g_actor.RunOnMap{
		ReplyTo:      callerPID,
		CallerRoot:   root,
		CallerThread: L,
		ScriptPath:   scriptPath,
		FuncName:     funcName,
		Args:         args,
	})
	return L.Yield(lua.LNil)
}

func (s mapSystem) Warp(character *entity.Character, targetMap *entity.Map, spawnPoint uint8) error {
	if targetMap == nil {
		return fmt.Errorf("target map is nil")
	}
	if character == nil {
		return fmt.Errorf("character is nil")
	}
	targetPID := targetMap.GetActorPID()
	if targetPID == nil {
		return fmt.Errorf("target map actor not found")
	}
	currentMap := character.GetMap()
	if currentMap != nil {
		if srcPID := currentMap.GetActorPID(); srcPID != nil {
			s.gs.GetRootContext().Send(srcPID, &g_actor.TransferCharacter{
				Character: character,
				TargetPID: targetPID,
				Portal:    spawnPoint,
			})
			return nil
		}
		currentMap.RemovePlayer(character.GetID())
	}
	s.gs.GetRootContext().Send(targetPID, &g_actor.WarpCharacter{
		Character: character,
		Portal:    spawnPoint,
	})
	return nil
}

func (s mapSystem) CreateReturnDoor(ch *entity.Character, skillID gameconst.SkillID) {
	if s.gs == nil || ch == nil {
		return
	}
	m := ch.GetMap()
	if m == nil || m.Wz == nil {
		return
	}
	destMapID := uint32(m.Wz.ReturnMapId)
	if destMapID == 0 || destMapID == uint32(m.Wz.ID) {
		return
	}
	destMap := s.Get(destMapID)
	if destMap == nil {
		ch.Listener.OnMessage(ch, gameconst.MsgPinkText, gameconst.DoorNoTownPortalMessage)
		return
	}
	destPID := destMap.GetActorPID()
	if destPID == nil {
		ch.Listener.OnMessage(ch, gameconst.MsgPinkText, gameconst.DoorNoTownPortalMessage)
		return
	}
	srcPID := m.GetActorPID()
	if srcPID == nil {
		ch.Listener.OnMessage(ch, gameconst.MsgPinkText, gameconst.DoorNoTownPortalMessage)
		return
	}
	root := s.gs.GetRootContext()
	if root == nil {
		return
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
		CharacterID: ch.GetID(),
		OwnerID:     ch.GetID(),
		SkillID:     skillID,
		Field: entity.DoorEndpoint{
			MapID:    uint32(m.Wz.ID),
			PortalID: closestPortalID,
			Position: ch.Position,
		},
		PartyOwnerSlot: slot,
		PartyID:        ch.GetPartyID(),
	})
}

func (s mapSystem) RemoveReturnDoor(ownerID uint32, skillID uint32, counterpartMapWZID uint32) {
	mapInstance := s.Get(counterpartMapWZID)
	if mapInstance == nil {
		return
	}
	pid := mapInstance.GetActorPID()
	if pid == nil {
		return
	}
	if root := s.gs.GetRootContext(); root != nil {
		root.Send(pid, &g_actor.RemoveDoor{
			OwnerID: ownerID,
			SkillID: skillID,
		})
	}
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

func (s schedulerSystem) RunReactorRespawn(pid *actor.PID, spawnID uint32) {
	if s.gs == nil || pid == nil {
		return
	}
	payload := &c_actor.RunReactorRespawn{
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
