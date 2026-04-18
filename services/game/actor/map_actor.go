package actor

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/actor/timers"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

type MapActor struct {
	MapData        *entity.Map
	Context        core.ServerContext
	SaveCharacters func([]*entity.Character) error
	Ensure         EnsureCoordinator
	scheduler      *scheduler.TimerScheduler
	timerReg       *TimerRegistry
}

func (a *MapActor) Receive(ctx actor.Context) {
	msg := ctx.Message()
	if env, ok := msg.(*EnsureDeliver); ok {
		a.handleEnsureDeliver(ctx, env)
		return
	}
	a.dispatch(ctx, msg)
}

func (a *MapActor) dispatch(ctx actor.Context, msg interface{}) {
	switch m := msg.(type) {
	case *actor.Started:
		a.onStarted(ctx)
	case *actor.Stopped:
		a.onStopped(ctx)
	case *c_actor.HandlePacket:
		a.handlePacket(ctx, m)
	case *c_actor.ScheduleTimer:
		a.scheduleTimer(ctx, m)
	case *c_actor.ExecuteTimer:
		a.executeTimer(m)
	case *AddCharacter:
		a.addCharacter(ctx, m)
	case *RemoveCharacter:
		a.removeCharacter(m)
	case *WarpCharacter:
		a.warpCharacter(ctx, m)
	case *SpawnDoor:
		a.spawnTownMysticDoor(m)
	case *RemoveDoor:
		a.removeMysticDoor(m)
	case *ResumeLua:
		a.resumeLua(m)
	case *c_actor.RunCharacterTimer:
		a.runCharacterTimer(ctx, m)
	case *TimerTick:
		a.onTimerTick(ctx, m)
	case *SyncPartySnapshot:
		a.syncPartySnapshot(m)
	case *ClearPartyByPartyID:
		a.clearPartyByPartyID(m)
	case *SyncCharacterPartyState:
		a.handleSyncCharacterPartyState(m)
	case *DeliverPartyInvite:
		a.handleDeliverPartyInvite(m)
	case *DeliverPartyStatusMessage:
		a.handleDeliverPartyStatusMessage(m)
	case *DeliverPartyUpdateJoin:
		a.handleDeliverPartyUpdateJoin(m)
	case *DeliverPartyUpdateLeave:
		a.handleDeliverPartyUpdateLeave(m)
	case *DeliverPartyUpdateDisband:
		a.handleDeliverPartyUpdateDisband(m)
	case *DeliverPartyUpdateLeaderChange:
		a.handleDeliverPartyUpdateLeaderChange(m)
	case *DeliverPartyUpdateLogOnOff:
		a.handleDeliverPartyUpdateLogOnOff(m)
	case *DeliverPartyUpdateSilent:
		a.handleDeliverPartyUpdateSilent(m)
	default:
		return
	}
}

func (a *MapActor) handleEnsureDeliver(ctx actor.Context, env *EnsureDeliver) {
	if env == nil {
		return
	}
	inner := env.Inner
	if inner == nil {
		a.ensureFinish(ctx, env, false, "nil_inner")
		return
	}
	if !a.hasCharacterOnMap(env.CharacterID) {
		a.ensureNotOnMap(env)
		return
	}
	a.dispatch(ctx, inner)
	a.ensureFinish(ctx, env, true, "")
}

func (a *MapActor) hasCharacterOnMap(characterID uint32) bool {
	if a.MapData == nil || characterID == 0 {
		return false
	}
	return a.MapData.GetPlayer(characterID) != nil
}

func (a *MapActor) resumeLua(msg *ResumeLua) {
	if msg.Root == nil || msg.Thread == nil {
		return
	}
	state, _, _ := msg.Root.Resume(msg.Thread, nil)
	if state == lua.ResumeOK {
		luax.ClearThreadPID(msg.Thread)
		msg.Thread.Close()
	}
}

func (a *MapActor) handlePacket(ctx actor.Context, msg *c_actor.HandlePacket) {
	err := core.ExecutePacketHandler(ctx, a.Context, msg.Client, msg.Opcode, msg.Data, msg.LogicActorPID)
	if err != nil {
		log.Printf("Error handling packet 0x%02X: %v", msg.Opcode, err)
	}
}

func (a *MapActor) scheduleTimer(ctx actor.Context, msg *c_actor.ScheduleTimer) {
	if msg.Logic == nil {
		return
	}

	go func() {
		time.Sleep(msg.Interval)
		ctx.Send(ctx.Self(), &c_actor.ExecuteTimer{
			Logic: msg.Logic,
		})
	}()
}

func (a *MapActor) executeTimer(msg *c_actor.ExecuteTimer) {
	if msg.Logic == nil {
		return
	}

	if err := msg.Logic(); err != nil {
		log.Printf("Error executing timer logic: %v", err)
	}
}

func (a *MapActor) addCharacter(ctx actor.Context, msg *AddCharacter) {
	if a.MapData == nil {
		return
	}
	a.MapData.AddPlayer(msg.Character.GetID(), msg.Character, msg.SpawnPoint, msg.Init)
	msg.Character.ResumeTimers(ctx.Self())
	msg.Character.Listener.OnPartyMemberFieldsChanged(msg.Character)
}

func (a *MapActor) removeCharacter(msg *RemoveCharacter) {
	if a.MapData == nil {
		return
	}
	a.MapData.RemovePlayer(msg.CharacterID)
}

func (a *MapActor) warpCharacter(ctx actor.Context, msg *WarpCharacter) {
	if a.MapData == nil {
		return
	}
	a.MapData.AddPlayer(msg.Character.GetID(), msg.Character, msg.Portal, false)
	msg.Character.ResumeTimers(ctx.Self())
	msg.Character.Listener.OnPartyMemberFieldsChanged(msg.Character)
}

func (a *MapActor) spawnTownMysticDoor(msg *SpawnDoor) {
	if a.MapData == nil || msg == nil || a.MapData.Wz == nil {
		return
	}
	wz := a.MapData.Wz
	position, ok := wz.GetSpawnPosition(msg.ReturnPortalID)
	if !ok {
		if p, ok2 := wz.Portals[msg.ReturnPortalID]; ok2 {
			position = p.Position
		} else {
			return
		}
	}
	door := entity.NewDoor(position, msg.OwnerID, msg.SkillID, wz.ID, msg.FieldMapID, msg.ReturnPortalID, msg.FieldPortalID)
	a.MapData.AddDoor(door)
}

func (a *MapActor) removeMysticDoor(msg *RemoveDoor) {
	if a.MapData == nil || msg == nil {
		return
	}
	a.MapData.RemoveMysticDoorByOwnerSkill(msg.OwnerID, constant.SkillID(msg.SkillID), true)
}

func (a *MapActor) runCharacterTimer(ctx actor.Context, msg *c_actor.RunCharacterTimer) {
	if a.MapData == nil {
		return
	}
	ch := a.MapData.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	entry := ch.GetTimerEntry(msg.Key)
	if entry == nil {
		return
	}
	if entry.Callback != nil {
		entry.Callback()
	}
	if entry.Repeat && ch.GetTimerEntry(msg.Key) != nil {
		characterID := msg.CharacterID
		key := msg.Key
		entry.NextFireAt = time.Now().Add(entry.Interval)
		entry.Timer = time.AfterFunc(entry.Interval, func() {
			ctx.Send(ctx.Self(), &c_actor.RunCharacterTimer{CharacterID: characterID, Key: key})
		})
	} else if !entry.Repeat {
		ch.RemoveTimer(msg.Key)
	}
}

func (a *MapActor) onStarted(ctx actor.Context) {
	L := luax.NewState()
	luax.RegisterRootLuaState(ctx.Self().String(), L)

	a.scheduler = scheduler.NewTimerScheduler(ctx)
	a.timerReg = NewTimerRegistry()
	a.registerTimers()

	for _, handler := range a.timerReg.GetAllHandlers() {
		a.scheduler.SendRepeatedly(
			handler.GetInitialDelay(),
			handler.GetInterval(),
			ctx.Self(),
			&TimerTick{
				HandlerName: handler.GetName(),
			},
		)
	}
}

func (a *MapActor) onStopped(ctx actor.Context) {
	luax.UnregisterRootLuaState(ctx.Self().String())
}

func (a *MapActor) registerTimers() {
	RegisterTimer[*timers.MobSpawnTimer](a.timerReg)
	RegisterTimer[*timers.ItemCleanupTimer](a.timerReg)
	RegisterTimer[*timers.CooldownCheckTimer](a.timerReg)
	RegisterTimer[*timers.ClientPingTimer](a.timerReg)
	RegisterTimer[*timers.BuffExpireTimer](a.timerReg)
	RegisterTimer[*timers.MobBuffExpireTimer](a.timerReg)
	RegisterTimer[*timers.MobPoisonTickTimer](a.timerReg)
	RegisterTimer[*timers.MistExpireTimer](a.timerReg)
	RegisterTimer[*timers.MistPoisonTickTimer](a.timerReg)
	if a.SaveCharacters != nil {
		a.timerReg.handlers[timers.CharacterSaveTimerName] = timers.NewCharacterSaveTimer(a.SaveCharacters)
	}
}

func (a *MapActor) onTimerTick(ctx actor.Context, msg *TimerTick) {
	if a.MapData == nil {
		return
	}
	handler := a.timerReg.GetHandler(msg.HandlerName)
	if handler == nil {
		return
	}
	if err := handler.Handle(ctx, a.MapData); err != nil {
		log.Printf("Timer handler %s error: %v", handler.GetName(), err)
	}
}

func (a *MapActor) syncPartySnapshot(msg *SyncPartySnapshot) {
	if a.MapData == nil || msg == nil || msg.Snapshot == nil {
		return
	}
	partyID := msg.Snapshot.GetPartyId()
	if partyID == 0 {
		return
	}
	memberSet := make(map[uint32]struct{}, len(msg.Snapshot.GetMembers()))
	for _, member := range msg.Snapshot.GetMembers() {
		memberSet[member.GetCharacterId()] = struct{}{}
	}
	for _, obj := range a.MapData.GetAllPlayers() {
		ch, ok := obj.(*entity.Character)
		if !ok || ch == nil {
			continue
		}
		if _, exists := memberSet[ch.GetID()]; exists {
			id := partyID
			ch.SetPartyID(&id)
			continue
		}
		if currentPartyID, ok := ch.GetPartyID(); ok && currentPartyID == partyID {
			ch.SetPartyID(nil)
		}
	}
}

func (a *MapActor) clearPartyByPartyID(msg *ClearPartyByPartyID) {
	if a.MapData == nil || msg == nil || msg.PartyID == 0 {
		return
	}
	for _, obj := range a.MapData.GetAllPlayers() {
		ch, ok := obj.(*entity.Character)
		if !ok || ch == nil {
			continue
		}
		if currentPartyID, ok := ch.GetPartyID(); ok && currentPartyID == msg.PartyID {
			ch.SetPartyID(nil)
		}
	}
}

func (a *MapActor) handleSyncCharacterPartyState(msg *SyncCharacterPartyState) {
	if a.MapData == nil || msg == nil {
		return
	}
	ch := a.MapData.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	if msg.PartyID == nil {
		ch.SetPartyID(nil)
	} else {
		id := *msg.PartyID
		ch.SetPartyID(&id)
	}
}

func (a *MapActor) ensureNotOnMap(msg *EnsureDeliver) {
	if a.Ensure != nil {
		a.Ensure.EnsureRedispatch(msg)
	}
}

func (a *MapActor) ensureFinish(ctx actor.Context, msg *EnsureDeliver, ok bool, reason string) {
	if msg == nil {
		return
	}
	if a.Ensure != nil {
		a.Ensure.EnsureComplete(msg.CorrelationID)
	}
	if msg.Caller == nil {
		return
	}
	ctx.Send(msg.Caller, &EnsureResult{
		CorrelationID: msg.CorrelationID,
		OK:            ok,
		Reason:        reason,
	})
}

func (a *MapActor) handleDeliverPartyInvite(msg *DeliverPartyInvite) {
	if a.MapData == nil || msg == nil {
		return
	}
	ch := a.MapData.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	_ = ch.Send(&response.PartyInvite{
		PartyID:     msg.PartyID,
		InviterName: msg.InviterName,
		PartySearch: msg.PartySearch,
	}, types.SEND_POLICY_ENCRYPT)
}

func (a *MapActor) handleDeliverPartyStatusMessage(msg *DeliverPartyStatusMessage) {
	if a.MapData == nil || msg == nil {
		return
	}
	ch := a.MapData.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	_ = ch.Send(&response.PartyStatusMessage{
		Code: msg.Code,
		Name: msg.Name,
	}, types.SEND_POLICY_ENCRYPT)
}

func (a *MapActor) handleDeliverPartyUpdateJoin(msg *DeliverPartyUpdateJoin) {
	if a.MapData == nil || msg == nil {
		return
	}
	ch := a.MapData.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	_ = ch.Send(&response.PartyUpdateJoin{
		ForChannel:           msg.ForChannel,
		PartyID:              msg.PartyID,
		JoiningCharacterName: msg.JoinName,
		LeaderCharacterID:    msg.LeaderID,
		Members:              msg.Members,
	}, types.SEND_POLICY_ENCRYPT)
}

func (a *MapActor) handleDeliverPartyUpdateLeave(msg *DeliverPartyUpdateLeave) {
	if a.MapData == nil || msg == nil {
		return
	}
	ch := a.MapData.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	if msg.Expelled {
		_ = ch.Send(&response.PartyUpdateExpel{
			ForChannel:          msg.ForChannel,
			PartyID:             msg.PartyID,
			TargetCharacterID:   msg.TargetID,
			TargetCharacterName: msg.TargetName,
			LeaderCharacterID:   msg.LeaderID,
			Members:             msg.Members,
		}, types.SEND_POLICY_ENCRYPT)
		return
	}
	_ = ch.Send(&response.PartyUpdateLeave{
		ForChannel:          msg.ForChannel,
		PartyID:             msg.PartyID,
		TargetCharacterID:   msg.TargetID,
		TargetCharacterName: msg.TargetName,
		LeaderCharacterID:   msg.LeaderID,
		Members:             msg.Members,
	}, types.SEND_POLICY_ENCRYPT)
}

func (a *MapActor) handleDeliverPartyUpdateDisband(msg *DeliverPartyUpdateDisband) {
	if a.MapData == nil || msg == nil {
		return
	}
	ch := a.MapData.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	_ = ch.Send(&response.PartyUpdateDisband{
		PartyID:           msg.PartyID,
		LeaderCharacterID: msg.LeaderID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (a *MapActor) handleDeliverPartyUpdateLeaderChange(msg *DeliverPartyUpdateLeaderChange) {
	if a.MapData == nil || msg == nil || msg.NewLeaderCharacterID == 0 {
		return
	}
	ch := a.MapData.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	_ = ch.Send(&response.PartyUpdateLeaderChange{
		NewLeaderCharacterID: msg.NewLeaderCharacterID,
		ByDisconnect:         msg.ByDisconnect,
	}, types.SEND_POLICY_ENCRYPT)
}

func (a *MapActor) handleDeliverPartyUpdateLogOnOff(msg *DeliverPartyUpdateLogOnOff) {
	if a.MapData == nil || msg == nil {
		return
	}
	ch := a.MapData.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	_ = ch.Send(&response.PartyUpdateLogOnOff{
		ForChannel:        msg.ForChannel,
		PartyID:           msg.PartyID,
		LeaderCharacterID: msg.LeaderID,
		Members:           msg.Members,
	}, types.SEND_POLICY_ENCRYPT)
}

func (a *MapActor) handleDeliverPartyUpdateSilent(msg *DeliverPartyUpdateSilent) {
	if a.MapData == nil || msg == nil {
		return
	}
	ch := a.MapData.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	_ = ch.Send(&response.PartyUpdateSilent{
		ForChannel:        msg.ForChannel,
		PartyID:           msg.PartyID,
		LeaderCharacterID: msg.LeaderID,
		Members:           msg.Members,
	}, types.SEND_POLICY_ENCRYPT)
}
