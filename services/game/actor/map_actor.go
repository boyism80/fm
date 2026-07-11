package actor

import (
	"github.com/boyism80/fm/core/clock"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/ensure"
	"github.com/boyism80/fm/core/luax"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/actor/timers"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	lua "github.com/yuin/gopher-lua"
)

type MapActor struct {
	Map       *entity.Map
	GameWorld entity.GameWorld
	scheduler *scheduler.TimerScheduler
	timerReg  *TimerRegistry
}

func (a *MapActor) Receive(ctx actor.Context) {
	msg := ctx.Message()
	if env, ok := msg.(*ensure.EnsureDeliver); ok {
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
		a.onRemoveCharacter(ctx, m)
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
	case *ResetMap:
		a.onResetMap(ctx, m)
	case *ResetMapAck:
		a.onResetMapAck(m)
	case *RunOnMap:
		a.onRunOnMap(ctx, m)
	case *RunOnMapAck:
		a.onRunOnMapAck(m)
	case *c_actor.RunObjectTimer:
		a.onRunObjectTimer(ctx, m)
	case *c_actor.RunReactorRespawn:
		a.onRunReactorRespawn(m)
	case *TimerTick:
		a.onTimerTick(ctx, m)
	case *SyncParty:
		a.onSyncParty(m)
	case *ClearPartyByPartyID:
		a.onClearPartyByPartyID(m)
	case *SyncCharacterPartyState:
		a.onSyncCharacterPartyState(m)
	case *PartyMemberLeft:
		a.onPartyMemberLeft(m)
	case *PartyDisband:
		a.onPartyDisband(m)
	case *DeliverPartyInvite:
		a.onDeliverPartyInvite(m)
	case *DeliverMultiChat:
		a.onDeliverMultiChat(m)
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
	case *DeliverBuddyChannelUpdate:
		a.onDeliverBuddyChannelUpdate(m)
	case *DeliverBuddyListUpdate:
		a.onDeliverBuddyListUpdate(m)
	case *DeliverBuddyAddRequest:
		a.onDeliverBuddyAddRequest(m)
	case *DeliverGuildInvite:
		a.onDeliverGuildInvite(m)
	case *DeliverGuildNewMember:
		a.onDeliverGuildNewMember(m)
	case *DeliverGuildLeaveSelf:
		a.onDeliverGuildLeaveSelf(m)
	case *DeliverGuildExpelSelf:
		a.onDeliverGuildExpelSelf(m)
	case *DeliverGuildMemberLeft:
		a.onDeliverGuildMemberLeft(m)
	case *DeliverGuildRankTitleChange:
		a.onDeliverGuildRankTitleChange(m)
	case *DeliverGuildMemberRankChange:
		a.onDeliverGuildMemberRankChange(m)
	case *DeliverGuildEmblemChange:
		a.onDeliverGuildEmblemChange(m)
	case *DeliverGuildNoticeChange:
		a.onDeliverGuildNoticeChange(m)
	case *DeliverGuildCapacityChange:
		a.onDeliverGuildCapacityChange(m)
	case *DeliverGuildMemberOnlineChange:
		a.onDeliverGuildMemberOnlineChange(m)
	case *DeliverGuildMemberFieldsChange:
		a.onDeliverGuildMemberFieldsChange(m)
	case *DeliverGuildDisbandSelf:
		a.onDeliverGuildDisbandSelf(m)
	case *DeliverGuildMessage:
		a.onDeliverGuildMessage(m)
	case *DeliverAllianceCreate:
		a.onDeliverAllianceCreate(m)
	case *DeliverAllianceDisband:
		a.onDeliverAllianceDisband(m)
	case *DeliverAllianceGuildLeft:
		a.onDeliverAllianceGuildLeft(m)
	case *DeliverAllianceInvite:
		a.onDeliverAllianceInvite(m)
	case *DeliverAllianceGuildAdded:
		a.onDeliverAllianceGuildAdded(m)
	case *DeliverAllianceInfoBroadcast:
		a.onDeliverAllianceInfoBroadcast(m)
	case *DeliverAllianceNoticeChanged:
		a.onDeliverAllianceNoticeChanged(m)
	case *DeliverAllianceLeaderChanged:
		a.onDeliverAllianceLeaderChanged(m)
	case *DeliverAllianceMemberRankChanged:
		a.onDeliverAllianceMemberRankChanged(m)
	case *DeliverAllianceMemberOnlineChange:
		a.onDeliverAllianceMemberOnlineChange(m)
	case *DeliverAllianceMemberFieldsChange:
		a.onDeliverAllianceMemberFieldsChange(m)
	case *DeliverMessage:
		a.onDeliverMessage(m)
	case *SaveMapCharacters:
		a.onSaveMapCharacters(ctx)
	default:
		return
	}
}

func (a *MapActor) handleEnsureDeliver(ctx actor.Context, env *ensure.EnsureDeliver) {
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
	if a.Map == nil || characterID == 0 {
		return false
	}
	return a.Map.GetPlayer(characterID) != nil
}

func (a *MapActor) onResumeLua(msg *ResumeLua) {
	if msg.Root == nil || msg.Thread == nil {
		return
	}
	resumeArgs := make([]interface{}, len(msg.Args))
	for i, a := range msg.Args {
		resumeArgs[i] = a
	}
	state, _, _ := luax.Resume(msg.Root, msg.Thread, "", resumeArgs...)
	if state == lua.ResumeOK {
		cfg, ok := luax.GetConfiguration(msg.Thread)
		if ok && cfg.KeepAlive {
			luax.Close(msg.Thread)
		}
	}
}

func (a *MapActor) onResetMap(ctx actor.Context, msg *ResetMap) {
	ok := false
	if a.Map != nil {
		a.Map.Reset()
		ok = true
	}
	if msg == nil || msg.ReplyTo == nil {
		return
	}
	ctx.Send(msg.ReplyTo, &ResetMapAck{
		Ok:     ok,
		Root:   msg.Root,
		Thread: msg.Thread,
	})
}

func (a *MapActor) onResetMapAck(msg *ResetMapAck) {
	if msg == nil || msg.Root == nil || msg.Thread == nil {
		return
	}
	args := []lua.LValue{lua.LBool(msg.Ok)}
	a.onResumeLua(&ResumeLua{Root: msg.Root, Thread: msg.Thread, Args: args})
}

func runOnMapResumeValues(ok bool, result lua.LValue, errMsg string) []lua.LValue {
	if !ok {
		return []lua.LValue{lua.LBool(false), lua.LNil, lua.LString(errMsg)}
	}
	if result == nil {
		result = lua.LNil
	}
	return []lua.LValue{lua.LBool(true), result, lua.LNil}
}

func (a *MapActor) onRunOnMap(ctx actor.Context, msg *RunOnMap) {
	if msg == nil || msg.ReplyTo == nil {
		return
	}
	system := ctx.ActorSystem()
	if a.Map == nil {
		ctx.Send(msg.ReplyTo, &RunOnMapAck{
			Root:   msg.CallerRoot,
			Thread: msg.CallerThread,
			Values: runOnMapResumeValues(false, lua.LNil, "map not found"),
		})
		return
	}
	a.Map.RunScript(ctx, msg.ScriptPath, msg.FuncName, msg.Args).Then(func(v interface{}) (interface{}, error) {
		if system == nil || system.Root == nil {
			return nil, nil
		}
		vals := luax.ResultValues(v)
		result := lua.LNil
		if len(vals) > 0 && vals[0] != nil {
			result = vals[0]
		}
		values := runOnMapResumeValues(true, result, "")
		system.Root.Send(msg.ReplyTo, &RunOnMapAck{
			Root:   msg.CallerRoot,
			Thread: msg.CallerThread,
			Values: values,
		})
		return nil, nil
	}).OnError(func(err error) {
		if system == nil || system.Root == nil {
			return
		}
		values := runOnMapResumeValues(false, lua.LNil, err.Error())
		system.Root.Send(msg.ReplyTo, &RunOnMapAck{
			Root:   msg.CallerRoot,
			Thread: msg.CallerThread,
			Values: values,
		})
	})
}

func (a *MapActor) onRunOnMapAck(msg *RunOnMapAck) {
	if msg == nil || msg.Root == nil || msg.Thread == nil {
		return
	}
	args := msg.Values
	if args == nil {
		args = runOnMapResumeValues(false, lua.LNil, "run_on_map: empty response")
	}
	a.onResumeLua(&ResumeLua{Root: msg.Root, Thread: msg.Thread, Args: args})
}

func (a *MapActor) onHandlerPacket(ctx actor.Context, msg *c_actor.HandlePacket) {
	err := core.ExecutePacketHandler(ctx, a.GameWorld, msg.Client, msg.Opcode, msg.Data, msg.LogicActorPID)
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
	if a.Map == nil {
		return
	}
	a.Map.AddPlayer(ctx, msg.Character.GetID(), msg.Character, msg.SpawnPoint, msg.Init)
	msg.Character.ResumeTimers(ctx.Self())
	msg.Character.Listener.OnPartyMemberFieldsChanged(msg.Character)
	if msg.Init {
		msg.Character.SendBuddyLoginSync()
	}
}

func (a *MapActor) onRemoveCharacter(ctx actor.Context, msg *RemoveCharacter) {
	if a.Map == nil {
		if ctx.Sender() != nil {
			ctx.Respond(struct{}{})
		}
		return
	}
	_ = a.Map.RemovePlayer(msg.CharacterID)
	if ctx.Sender() != nil {
		ctx.Respond(struct{}{})
	}
}

func (a *MapActor) onWarpCharacter(ctx actor.Context, msg *WarpCharacter) {
	if a.Map == nil {
		return
	}
	a.Map.AddPlayer(ctx, msg.Character.GetID(), msg.Character, msg.Portal, false)
	msg.Character.ResumeTimers(ctx.Self())
	msg.Character.Listener.OnPartyMemberFieldsChanged(msg.Character)
}

func (a *MapActor) onRemoveDoor(msg *RemoveDoor) {
	if a.Map == nil || msg == nil {
		return
	}
	a.Map.RemoveDoorByOwnerSkill(msg.OwnerID, constant.SkillID(msg.SkillID), true)
}

func (a *MapActor) onRequestSpawnDoor(ctx actor.Context, msg *RequestSpawnDoor) {
	if msg == nil || msg.ReplyTo == nil || a.Map == nil || a.Map.Wz == nil {
		return
	}
	portalID, townPos, ok := a.Map.TryAcquireMysticReturnPortal(msg.PartyOwnerSlot)
	if !ok {
		ctx.Send(msg.ReplyTo, &ResponseSpawnDoor{
			Ok:          false,
			CharacterID: msg.CharacterID,
			OwnerID:     msg.OwnerID,
			SkillID:     msg.SkillID,
			Field:       msg.Field,
			PartyID:     msg.PartyID,
		})
		return
	}
	committed := false
	defer func() {
		if !committed {
			a.Map.ReleaseMysticReturnPortal(portalID)
		}
	}()
	wz := a.Map.Wz
	returnEp := entity.DoorEndpoint{
		MapID:    uint32(wz.ID),
		PortalID: portalID,
		Position: townPos,
	}
	door := entity.NewDoor(
		msg.OwnerID,
		msg.SkillID,
		msg.Field,
		returnEp,
		msg.PartyID,
	)
	a.Map.AddDoor(door)
	committed = true
	ctx.Send(msg.ReplyTo, &ResponseSpawnDoor{
		Ok:          true,
		CharacterID: msg.CharacterID,
		OwnerID:     msg.OwnerID,
		SkillID:     msg.SkillID,
		Return:      returnEp,
		Field:       msg.Field,
		PartyID:     msg.PartyID,
	})
}

func (a *MapActor) onResponseSpawnDoor(msg *ResponseSpawnDoor) {
	if msg == nil || a.Map == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		if msg.Ok && a.Map.Wz != nil && a.Map.GameWorld != nil {
			a.Map.GameWorld.GetMapSystem().RemoveReturnDoor(msg.OwnerID, uint32(msg.SkillID), uint32(a.Map.Wz.ReturnMapId))
		}
		return
	}
	gw := ch.GameWorld
	if !msg.Ok {
		if gw != nil {
			ch.Listener.OnMessage(ch, constant.MsgPinkText, constant.DoorNoTownPortalMessage)
		}
		return
	}
	door := ch.SpawnFieldMapDoor(msg.SkillID, msg.Return, msg.Field)
	if door == nil && gw != nil {
		gw.GetMapSystem().RemoveReturnDoor(msg.OwnerID, uint32(msg.SkillID), uint32(ch.GetMap().Wz.ReturnMapId))
	}
}

func (a *MapActor) onRunReactorRespawn(msg *c_actor.RunReactorRespawn) {
	if a.Map == nil || msg == nil {
		return
	}

	reactorSpawn := a.Map.GetReactorSpawn(msg.SpawnID)
	if reactorSpawn == nil || reactorSpawn.Spawned {
		return
	}

	_, _ = a.Map.SpawnReactor(reactorSpawn)
}

func (a *MapActor) onRunObjectTimer(ctx actor.Context, msg *c_actor.RunObjectTimer) {
	if a.Map == nil {
		return
	}
	obj := a.Map.GetObject(constant.ObjectType(msg.ObjectType), msg.ID)
	if obj == nil {
		return
	}
	entry := obj.GetTimerEntry(msg.Key)
	if entry == nil {
		return
	}
	if entry.Callback != nil {
		entry.Callback()
	}
	if entry.Repeat {
		obj.RescheduleTimer(msg.Key)
	} else {
		obj.RemoveTimer(msg.Key)
	}
}

func (a *MapActor) onStarted(ctx actor.Context) {
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
	if a.Map != nil {
		a.Map.ClearLuaRoot()
	}
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
	RegisterTimer[*timers.CharacterSaveTimer](a.timerReg)
	RegisterTimer[*timers.PartySearchTimer](a.timerReg)
}

func (a *MapActor) onTimerTick(ctx actor.Context, msg *TimerTick) {
	if a.Map == nil {
		return
	}
	handler := a.timerReg.GetHandler(msg.HandlerName)
	if handler == nil {
		return
	}
	if err := handler.Handle(ctx, a.Map); err != nil {
		log.Printf("Timer handler %s error: %v", handler.GetName(), err)
	}
}

func (a *MapActor) onSyncParty(msg *SyncParty) {
	if a.Map == nil || msg == nil || msg.Party == nil {
		return
	}
	partyID := msg.Party.GetPartyId()
	memberSet := make(map[uint32]struct{}, len(msg.Party.GetMembers()))
	for _, member := range msg.Party.GetMembers() {
		if member == nil {
			continue
		}
		memberSet[member.GetCharacterId()] = struct{}{}
	}
	for _, obj := range a.Map.GetAllPlayers() {
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
	if a.Map == nil || msg == nil {
		return
	}
	for _, obj := range a.Map.GetAllPlayers() {
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
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
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

func (a *MapActor) onPartyMemberLeft(msg *PartyMemberLeft) {
	if a.Map == nil || msg == nil {
		return
	}
	a.Map.ApplyPartyLeaveDoorSync(msg.LeaverID)
}

func (a *MapActor) onPartyDisband(msg *PartyDisband) {
	if a.Map == nil || msg == nil {
		return
	}
	a.Map.ApplyPartyDisbandDoorSync(msg.FormerMemberIDs)
}

func (a *MapActor) onSaveMapCharacters(ctx actor.Context) {
	sender := ctx.Sender()
	if sender == nil {
		return
	}
	ack := &SaveMapCharactersAck{}
	if a.Map == nil || a.GameWorld == nil {
		ctx.Respond(ack)
		return
	}
	if a.Map.Wz != nil {
		ack.MapID = uint32(a.Map.Wz.ID)
	}
	allPlayers := a.Map.GetAllPlayers()
	if len(allPlayers) == 0 {
		ctx.Respond(ack)
		return
	}
	chars := make([]*entity.Character, 0, len(allPlayers))
	for _, obj := range allPlayers {
		if ch, ok := obj.(*entity.Character); ok && ch != nil {
			chars = append(chars, ch)
		}
	}
	ack.Saved = len(chars)
	if len(chars) == 0 {
		ctx.Respond(ack)
		return
	}
	p := a.GameWorld.SaveAsync(ctx, chars)
	if p == nil {
		ack.Err = "nil save promise"
		ctx.Respond(ack)
		return
	}
	var saveErr error
	p.OnError(func(err error) {
		saveErr = err
	}).Finally(func() {
		if saveErr != nil {
			ack.Err = saveErr.Error()
		}
		// Finally runs outside the actor receive turn, so ctx.Respond can lose the original sender/future context.
		// Send ack explicitly to the captured sender PID.
		system := ctx.ActorSystem()
		if system == nil || system.Root == nil {
			return
		}
		system.Root.Send(sender, ack)
	})
}

func (a *MapActor) ensureNotOnMap(msg *ensure.EnsureDeliver) {
	if a.GameWorld != nil {
		a.GameWorld.EnsureRedispatch(msg)
	}
}

func (a *MapActor) ensureFinish(ctx actor.Context, msg *ensure.EnsureDeliver, ok bool, reason string) {
	if msg == nil {
		return
	}
	if a.GameWorld != nil {
		a.GameWorld.EnsureComplete(msg.CorrelationID)
	}
	if msg.Caller == nil {
		return
	}
	ctx.Send(msg.Caller, &ensure.EnsureResult{
		CorrelationID: msg.CorrelationID,
		OK:            ok,
		Reason:        reason,
	})
}

func (a *MapActor) onDeliverPartyInvite(msg *DeliverPartyInvite) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnPartyInvite(ch, msg.PartyID, msg.InviterName, msg.PartySearch)
}

func (a *MapActor) onDeliverMultiChat(msg *DeliverMultiChat) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnMultiChat(ch, msg.Mode, msg.SenderName, msg.Message)
}

func (a *MapActor) onDeliverPartyStatusMessage(msg *DeliverPartyStatusMessage) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnPartyStatusMessage(ch, msg.Code, msg.Name)
}

func (a *MapActor) onDeliverPartyUpdateJoin(msg *DeliverPartyUpdateJoin) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnPartyUpdateJoin(ch, msg.ForChannel, msg.PartyID, msg.JoinName, msg.LeaderID, msg.Members)
	if ch.GetID() != msg.JoinCharacterID {
		return
	}
	ch.Listener.OnPartyMemberHPChanged(ch, nil)
	for _, obj := range a.Map.GetObjects(constant.ObjectTypeCharacter) {
		peer, ok := obj.(*entity.Character)
		if !ok || peer == nil || peer.GetID() == ch.GetID() {
			continue
		}
		pPeer := peer.GetPartyID()
		if pPeer == nil || *pPeer != msg.PartyID {
			continue
		}
		peer.Listener.OnPartyMemberHPChanged(peer, ch)
	}
}

func (a *MapActor) onDeliverPartyUpdateLeave(msg *DeliverPartyUpdateLeave) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	if msg.Expelled {
		ch.Listener.OnPartyUpdateExpel(ch, msg.ForChannel, msg.PartyID, msg.TargetID, msg.TargetName, msg.LeaderID, msg.Members)
		return
	}
	ch.Listener.OnPartyUpdateLeave(ch, msg.ForChannel, msg.PartyID, msg.TargetID, msg.TargetName, msg.LeaderID, msg.Members)
}

func (a *MapActor) onDeliverPartyUpdateDisband(msg *DeliverPartyUpdateDisband) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnPartyUpdateDisband(ch, msg.PartyID, msg.LeaderID)
}

func (a *MapActor) onDeliverPartyUpdateLeaderChange(msg *DeliverPartyUpdateLeaderChange) {
	if a.Map == nil || msg == nil || msg.NewLeaderCharacterID == 0 {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnPartyUpdateLeaderChange(ch, msg.NewLeaderCharacterID, msg.ByDisconnect)
}

func (a *MapActor) onDeliverPartyUpdateLogOnOff(msg *DeliverPartyUpdateLogOnOff) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnPartyUpdateLogOnOff(ch, msg.ForChannel, msg.PartyID, msg.LeaderID, msg.Members)
}

func (a *MapActor) onDeliverPartyUpdateSilent(msg *DeliverPartyUpdateSilent) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnPartyUpdateSilent(ch, msg.ForChannel, msg.PartyID, msg.LeaderID, msg.Members)
}

func (a *MapActor) onDeliverBuddyChannelUpdate(msg *DeliverBuddyChannelUpdate) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.RecipientCharacterID)
	if ch == nil {
		return
	}
	ch.BuddyList().SetChannel(msg.BuddyCharacterID, msg.Channel)
	ch.Listener.OnBuddyChannelUpdate(ch, msg.BuddyCharacterID, msg.Channel)
}

func (a *MapActor) onDeliverBuddyListUpdate(msg *DeliverBuddyListUpdate) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.RecipientCharacterID)
	if ch == nil {
		return
	}
	action := pconst.BuddyListSyncAction(msg.SyncAction)
	bl := ch.BuddyList()
	if action == pconst.BuddyListSyncDelete {
		bl.Clear()
	}
	for _, entry := range msg.Entries {
		if entry.CharacterID == 0 {
			continue
		}
		channel := entry.Channel
		if channel < 0 {
			channel = -1
		}
		bl.Upsert(entity.BuddyListEntry{
			CharacterID: entry.CharacterID,
			Name:        entry.Name,
			Group:       entry.Group,
			Pending:     entry.Pending,
			Channel:     channel,
		})
	}
	ch.Listener.OnBuddyListUpdate(ch, action, msg.Entries)
}

func (a *MapActor) onDeliverBuddyAddRequest(msg *DeliverBuddyAddRequest) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.RecipientCharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnBuddyAddRequest(ch, msg.FromCharacterID, msg.FromName)
}

func (a *MapActor) onDeliverGuildInvite(msg *DeliverGuildInvite) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	if _, inGuild := ch.GetGuildID(); inGuild {
		if a.GameWorld != nil {
			a.GameWorld.GetDispatchSystem().SendTo(msg.InviterCharacterID, &DeliverGuildMessage{
				CharacterID: msg.InviterCharacterID,
				Code:        pconst.GuildResponseAlreadyInGuild,
			})
		}
		return
	}
	now := clock.Now()
	for id, expiresAt := range ch.GuildInvites {
		if !now.Before(expiresAt) {
			delete(ch.GuildInvites, id)
		}
	}
	if len(ch.GuildInvites) > 0 {
		if a.GameWorld != nil {
			a.GameWorld.GetDispatchSystem().SendTo(msg.InviterCharacterID, &DeliverMessage{
				CharacterID: msg.InviterCharacterID,
				MessageType: constant.MsgPinkText,
				Message:     constant.GuildInviteTargetBusyMessage,
			})
		}
		return
	}
	ch.GuildInvites[msg.GuildID] = now.Add(constant.GuildInviteDuration)
	ch.Listener.OnGuildInvite(ch, msg.GuildID, msg.InviterName)
}

func (a *MapActor) onDeliverGuildNewMember(msg *DeliverGuildNewMember) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildNewMember(ch, msg.GuildID, msg.Member)
}

func (a *MapActor) onDeliverGuildLeaveSelf(msg *DeliverGuildLeaveSelf) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildLeaveSelf(ch)
}

func (a *MapActor) onDeliverGuildExpelSelf(msg *DeliverGuildExpelSelf) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildExpelledSelf(ch, msg.GuildID)
}

func (a *MapActor) onDeliverGuildMemberLeft(msg *DeliverGuildMemberLeft) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildMemberLeft(ch, msg.GuildID, msg.TargetID, msg.TargetName, msg.WasExpelled)
}

func (a *MapActor) onDeliverGuildRankTitleChange(msg *DeliverGuildRankTitleChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildRankTitleChange(ch, msg.GuildID, msg.RankTitles)
}

func (a *MapActor) onDeliverGuildMemberRankChange(msg *DeliverGuildMemberRankChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildMemberRankChange(ch, msg.GuildID, msg.TargetID, msg.GuildRank)
}

func (a *MapActor) onDeliverGuildEmblemChange(msg *DeliverGuildEmblemChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildEmblemChange(ch, msg.GuildID, msg.LogoBG, msg.LogoBGColor, msg.Logo, msg.LogoColor)
}

func (a *MapActor) onDeliverGuildNoticeChange(msg *DeliverGuildNoticeChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildNoticeChange(ch, msg.GuildID, msg.Notice)
}

func (a *MapActor) onDeliverGuildCapacityChange(msg *DeliverGuildCapacityChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildCapacityChange(ch, msg.GuildID, msg.Capacity)
}

func (a *MapActor) onDeliverGuildMemberOnlineChange(msg *DeliverGuildMemberOnlineChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildMemberOnlineChange(ch, msg.GuildID, msg.SubjectID, msg.Online)
}

func (a *MapActor) onDeliverGuildMemberFieldsChange(msg *DeliverGuildMemberFieldsChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildMemberFieldsChange(ch, msg.GuildID, msg.SubjectID, msg.Level, msg.ClassID)
}

func (a *MapActor) onDeliverGuildDisbandSelf(msg *DeliverGuildDisbandSelf) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildDisbandSelf(ch, msg.GuildID)
}

func (a *MapActor) onDeliverAllianceCreate(msg *DeliverAllianceCreate) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceCreate(ch, msg.Info, msg.Guilds, msg.MembershipGuilds)
}

func (a *MapActor) onDeliverAllianceDisband(msg *DeliverAllianceDisband) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceDisband(ch, msg.AllianceID)
}

func (a *MapActor) onDeliverAllianceInfoBroadcast(msg *DeliverAllianceInfoBroadcast) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceInfoBroadcast(ch, msg.Info)
}

func (a *MapActor) onDeliverAllianceNoticeChanged(msg *DeliverAllianceNoticeChanged) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceNoticeChanged(ch, msg.Info)
}

func (a *MapActor) onDeliverAllianceLeaderChanged(msg *DeliverAllianceLeaderChanged) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceLeaderChanged(ch, msg.AllianceID, msg.OldLeaderID, msg.NewLeaderID, msg.Info, msg.Guilds)
}

func (a *MapActor) onDeliverAllianceMemberRankChanged(msg *DeliverAllianceMemberRankChanged) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceMemberRankChanged(ch, msg.Info, msg.Guilds)
}

func (a *MapActor) onDeliverAllianceGuildLeft(msg *DeliverAllianceGuildLeft) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceGuildLeft(
		ch,
		msg.Info,
		msg.RemovedGuildID,
		msg.RemovedGuild,
		msg.RemovedMembers,
		msg.Expelled,
		msg.Leaving,
	)
}

func (a *MapActor) onDeliverAllianceInvite(msg *DeliverAllianceInvite) {
	if a.Map == nil || msg == nil || a.GameWorld == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	guildID, inGuild := ch.GetGuildID()
	if !inGuild || guildID != msg.TargetGuildID {
		return
	}
	if g := a.GameWorld.GetGuildSystem().Get(guildID); g != nil {
		if _, joined := g.GetAllianceID(); joined {
			return
		}
	}
	if !a.GameWorld.GetGuildSystem().TrySetAllianceInvite(msg.TargetGuildID, msg.AllianceID, msg.ExpiresAt) {
		a.GameWorld.GetDispatchSystem().SendTo(msg.InviterCharacterID, &DeliverMessage{
			CharacterID: msg.InviterCharacterID,
			MessageType: constant.MsgPinkText,
			Message:     constant.GuildInviteTargetBusyMessage,
		})
		return
	}
	ch.Listener.OnAllianceInvite(ch, msg.InviterGuildID, msg.InviterName, msg.AllianceName)
}

func (a *MapActor) onDeliverAllianceGuildAdded(msg *DeliverAllianceGuildAdded) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	var membershipGuild *dto.AllianceMembershipChangeGuild
	if msg.HasMembership {
		membershipGuild = &msg.MembershipGuild
	}
	ch.Listener.OnAllianceGuildAdded(
		ch,
		msg.Info,
		msg.Guilds,
		msg.NewGuildID,
		msg.AddedGuild,
		msg.Members,
		msg.Joining,
		membershipGuild,
	)
}

func (a *MapActor) onDeliverAllianceMemberOnlineChange(msg *DeliverAllianceMemberOnlineChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceMemberOnlineChange(ch, msg.AllianceID, msg.GuildID, msg.SubjectID, msg.Online)
}

func (a *MapActor) onDeliverAllianceMemberFieldsChange(msg *DeliverAllianceMemberFieldsChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceMemberFieldsChange(ch, msg.AllianceID, msg.GuildID, msg.SubjectID, msg.Level, msg.ClassID)
}

func (a *MapActor) onDeliverGuildMessage(msg *DeliverGuildMessage) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildMessage(ch, msg.Code)
}

func (a *MapActor) onDeliverMessage(msg *DeliverMessage) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnMessage(ch, msg.MessageType, msg.Message)
}
