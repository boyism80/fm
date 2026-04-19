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
		a.onHandlerPacket(ctx, m)
	case *c_actor.ScheduleTimer:
		a.onScheduleTimer(ctx, m)
	case *c_actor.ExecuteTimer:
		a.onExecuteTimer(m)
	case *AddCharacter:
		a.onAddCharacter(ctx, m)
	case *RemoveCharacter:
		a.onRemoveCharacter(m)
	case *WarpCharacter:
		a.onWarpCharacter(ctx, m)
	case *RequestSpawnDoor:
		a.onRequestSpawnDoor(ctx, m)
	case *ResponseSpawnDoor:
		a.onResponseSpawnDoor(m)
	case *RemoveDoor:
		a.onRemoveDoor(m)
	case *ResumeLua:
		a.onResumeLua(m)
	case *c_actor.RunCharacterTimer:
		a.onRunCharacterTimer(ctx, m)
	case *TimerTick:
		a.onTimerTick(ctx, m)
	case *SyncPartySnapshot:
		a.onSyncPartySnapshot(m)
	case *ClearPartyByPartyID:
		a.onClearPartyByPartyID(m)
	case *SyncCharacterPartyState:
		a.onSyncCharacterPartyState(m)
	case *DeliverPartyInvite:
		a.onDeliverPartyInvite(m)
	case *DeliverPartyStatusMessage:
		a.onDeliverPartyStatusMessage(m)
	case *DeliverPartyUpdateJoin:
		a.onDeliverPartyUpdateJoin(m)
	case *DeliverPartyUpdateLeave:
		a.onDeliverPartyUpdateLeave(m)
	case *DeliverPartyUpdateDisband:
		a.onDeliverPartyUpdateDisband(m)
	case *DeliverPartyUpdateLeaderChange:
		a.onDeliverPartyUpdateLeaderChange(m)
	case *DeliverPartyUpdateLogOnOff:
		a.onDeliverPartyUpdateLogOnOff(m)
	case *DeliverPartyUpdateSilent:
		a.onDeliverPartyUpdateSilent(m)
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

func (a *MapActor) onResumeLua(msg *ResumeLua) {
	if msg.Root == nil || msg.Thread == nil {
		return
	}
	state, _, _ := msg.Root.Resume(msg.Thread, nil)
	if state == lua.ResumeOK {
		luax.ClearThreadPID(msg.Thread)
		msg.Thread.Close()
	}
}

func (a *MapActor) onHandlerPacket(ctx actor.Context, msg *c_actor.HandlePacket) {
	err := core.ExecutePacketHandler(ctx, a.Context, msg.Client, msg.Opcode, msg.Data, msg.LogicActorPID)
	if err != nil {
		log.Printf("Error handling packet 0x%02X: %v", msg.Opcode, err)
	}
}

func (a *MapActor) onScheduleTimer(ctx actor.Context, msg *c_actor.ScheduleTimer) {
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

func (a *MapActor) onExecuteTimer(msg *c_actor.ExecuteTimer) {
	if msg.Logic == nil {
		return
	}

	if err := msg.Logic(); err != nil {
		log.Printf("Error executing timer logic: %v", err)
	}
}

func (a *MapActor) onAddCharacter(ctx actor.Context, msg *AddCharacter) {
	if a.MapData == nil {
		return
	}
	a.MapData.AddPlayer(msg.Character.GetID(), msg.Character, msg.SpawnPoint, msg.Init)
	msg.Character.ResumeTimers(ctx.Self())
	msg.Character.Listener.OnPartyMemberFieldsChanged(msg.Character)
}

func (a *MapActor) onRemoveCharacter(msg *RemoveCharacter) {
	if a.MapData == nil {
		return
	}
	a.MapData.RemovePlayer(msg.CharacterID)
}

func (a *MapActor) onWarpCharacter(ctx actor.Context, msg *WarpCharacter) {
	if a.MapData == nil {
		return
	}
	a.MapData.AddPlayer(msg.Character.GetID(), msg.Character, msg.Portal, false)
	msg.Character.ResumeTimers(ctx.Self())
	msg.Character.Listener.OnPartyMemberFieldsChanged(msg.Character)
}

func (a *MapActor) onRemoveDoor(msg *RemoveDoor) {
	if a.MapData == nil || msg == nil {
		return
	}
	a.MapData.RemoveDoorByOwnerSkill(msg.OwnerID, constant.SkillID(msg.SkillID), true)
}

func (a *MapActor) onRequestSpawnDoor(ctx actor.Context, msg *RequestSpawnDoor) {
	if msg == nil || msg.ReplyTo == nil || a.MapData == nil || a.MapData.Wz == nil {
		return
	}
	portalID, townPos, ok := a.MapData.TryAcquireMysticReturnPortal(msg.PartyOwnerSlot)
	if !ok {
		ctx.Send(msg.ReplyTo, &ResponseSpawnDoor{
			Ok:                 false,
			CharacterID:        msg.CharacterID,
			OwnerID:            msg.OwnerID,
			SkillID:            msg.SkillID,
			FieldMapID:         msg.FieldMapID,
			FieldPortalID:      msg.FieldPortalID,
			PartyID:            msg.PartyID,
			FieldAnchor:        msg.FieldAnchor,
			ReturnPortalID:     0,
			TownPortalPosition: types.Vector2[int16]{},
		})
		return
	}
	committed := false
	defer func() {
		if !committed {
			a.MapData.ReleaseMysticReturnPortal(portalID)
		}
	}()
	wz := a.MapData.Wz
	door := entity.NewDoor(
		townPos,
		msg.OwnerID,
		msg.SkillID,
		uint32(wz.ID),
		msg.FieldMapID,
		portalID,
		msg.FieldPortalID,
		msg.PartyID,
		msg.FieldAnchor,
		townPos,
	)
	a.MapData.AddDoor(door)
	committed = true
	ctx.Send(msg.ReplyTo, &ResponseSpawnDoor{
		Ok:                 true,
		CharacterID:        msg.CharacterID,
		OwnerID:            msg.OwnerID,
		SkillID:            msg.SkillID,
		ReturnPortalID:     portalID,
		TownPortalPosition: townPos,
		FieldMapID:         msg.FieldMapID,
		FieldPortalID:      msg.FieldPortalID,
		PartyID:            msg.PartyID,
		FieldAnchor:        msg.FieldAnchor,
	})
}

func (a *MapActor) onResponseSpawnDoor(msg *ResponseSpawnDoor) {
	if msg == nil || a.MapData == nil {
		return
	}
	ch := a.MapData.GetPlayer(msg.CharacterID)
	if ch == nil {
		if msg.Ok && a.MapData.Wz != nil && a.MapData.Context != nil {
			a.MapData.Context.NotifyDoorRemove(msg.OwnerID, uint32(msg.SkillID), uint32(a.MapData.Wz.ReturnMapId))
		}
		return
	}
	gc := ch.Context
	if !msg.Ok {
		if gc != nil {
			ch.Listener.OnMessage(ch, constant.MSG_PINK_TEXT, constant.DoorNoTownPortalMessage)
		}
		return
	}
	door := ch.SpawnFieldMapDoor(msg.SkillID, msg.ReturnPortalID, msg.TownPortalPosition, msg.FieldPortalID)
	if door == nil && gc != nil {
		gc.NotifyDoorRemove(msg.OwnerID, uint32(msg.SkillID), uint32(ch.GetMap().Wz.ReturnMapId))
	}
}

func (a *MapActor) onRunCharacterTimer(ctx actor.Context, msg *c_actor.RunCharacterTimer) {
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

func (a *MapActor) onSyncPartySnapshot(msg *SyncPartySnapshot) {
	if a.MapData == nil || msg == nil || msg.Snapshot == nil {
		return
	}
	partyID := msg.Snapshot.GetPartyId()
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
		if cur := ch.GetPartyID(); cur != nil && *cur == partyID {
			ch.SetPartyID(nil)
		}
	}
}

func (a *MapActor) onClearPartyByPartyID(msg *ClearPartyByPartyID) {
	if a.MapData == nil || msg == nil {
		return
	}
	for _, obj := range a.MapData.GetAllPlayers() {
		ch, ok := obj.(*entity.Character)
		if !ok || ch == nil {
			continue
		}
		if cur := ch.GetPartyID(); cur != nil && *cur == msg.PartyID {
			ch.SetPartyID(nil)
		}
	}
}

func (a *MapActor) onSyncCharacterPartyState(msg *SyncCharacterPartyState) {
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

func (a *MapActor) onDeliverPartyInvite(msg *DeliverPartyInvite) {
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

func (a *MapActor) onDeliverPartyStatusMessage(msg *DeliverPartyStatusMessage) {
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

func (a *MapActor) onDeliverPartyUpdateJoin(msg *DeliverPartyUpdateJoin) {
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

func (a *MapActor) onDeliverPartyUpdateLeave(msg *DeliverPartyUpdateLeave) {
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

func (a *MapActor) onDeliverPartyUpdateDisband(msg *DeliverPartyUpdateDisband) {
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

func (a *MapActor) onDeliverPartyUpdateLeaderChange(msg *DeliverPartyUpdateLeaderChange) {
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

func (a *MapActor) onDeliverPartyUpdateLogOnOff(msg *DeliverPartyUpdateLogOnOff) {
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

func (a *MapActor) onDeliverPartyUpdateSilent(msg *DeliverPartyUpdateSilent) {
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
