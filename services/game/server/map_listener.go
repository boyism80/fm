package server

import (
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/types"
)

type MapListenerImpl struct {
	gs *GameServer
}

func NewGameMapListener(gs *GameServer) *MapListenerImpl {
	return &MapListenerImpl{
		gs: gs,
	}
}

func (l *MapListenerImpl) OnPlayerAdded(mapInstance *game.Map, character *game.Character, spawnPoint uint8, init bool) {
	if mapInstance == nil {
		return
	}

	if init {
		characterDTO := character.BuildFullS2C()
		characterDTO.SpawnPoint = spawnPoint
		loginPacket := &response.Login{
			Channel:       l.gs.config.ChannelId,
			FieldKey:      character.FieldKey(),
			CharacterInfo: response.CharacterInfo{Character: characterDTO},
		}
		character.Send(loginPacket, types.SEND_POLICY_ENCRYPT)
		character.Listener().OnShowGuildInfo(character)
		character.Listener().OnShowAllianceInfo(character)
		character.Buffs().ShowAll()
		character.Listener().OnSkillMacros(character)
	} else {
		characterDTO := character.BuildS2C()
		characterDTO.SpawnPoint = spawnPoint
		warpPacket := &response.Warp{
			Character: characterDTO,
			Channel:   0,
			FieldKey:  character.FieldKey(),
		}
		character.Send(warpPacket, types.SEND_POLICY_ENCRYPT)
	}
	if character.IsHidden() {
		character.Send(&response.SuperHide{Hidden: true}, types.SEND_POLICY_ENCRYPT)
	}
	if character.Pets().Active() != nil {
		character.Send(&response.SpawnPet{CharacterID: character.GetID(), Pet: character.Pets().Active().BuildS2C()}, types.SEND_POLICY_ENCRYPT)
		if init {
			character.Pets().ShowExceptions()
		}
	}

	character.Send(&response.KeyMap{Slots: character.KeyLayout().Bindings()}, types.SEND_POLICY_ENCRYPT)
	if init && character.Pets().HPItem() != 0 {
		character.Send(&response.PetAutoHP{ItemID: character.Pets().HPItem()}, types.SEND_POLICY_ENCRYPT)
	}
	if init && character.Pets().MPItem() != 0 {
		character.Send(&response.PetAutoMP{ItemID: character.Pets().MPItem()}, types.SEND_POLICY_ENCRYPT)
	}

	if mapInstance.Wz() != nil && mapInstance.Wz().HasClock {
		_, _, _, hour, minute, second := clock.DateTimeTable(clock.Now())
		character.Send(&response.ClockTime{
			Hour:   uint8(hour),
			Minute: uint8(minute),
			Second: uint8(second),
		}, types.SEND_POLICY_ENCRYPT)
	}

	if itemID, name, ok := mapInstance.Song(); ok {
		character.Send(&response.CashSong{ItemID: itemID, Name: name}, types.SEND_POLICY_ENCRYPT)
	}
	if itemID, message, ok := mapInstance.CurrentWeather(); ok {
		character.Send(&response.Weather{ItemID: itemID, Message: message}, types.SEND_POLICY_ENCRYPT)
	}
	if y, ok := mapInstance.HealerY(); ok {
		character.Send(&response.HealerMove{Y: y}, types.SEND_POLICY_ENCRYPT)
	}

	for _, obj := range mapInstance.GetObjectsNear(character.GetPosition(), constant.ObjectTypeCharacter, &game.SearchOption{IncludeHidden: true}) {
		if viewer, ok := obj.(*game.Character); ok {
			character.SendSpawnSyncToViewer(viewer)
		}
	}

	for _, obj := range mapInstance.GetObjectsNear(character.GetPosition(), constant.ObjectTypeObject, nil) {
		obj.SendSpawnSyncToViewer(character)
	}

	mapInstance.SyncPartyHP(character)

	if match := character.Carnival().Match(); match != nil {
		mapID := mapInstance.GetMapID()
		if mapID == match.FieldMapID() || mapID == match.ReviveMapID() {
			match.NotifyCarnivalStart(character)
		}
	}
}

func (l *MapListenerImpl) OnPlayerRemoved(mapInstance *game.Map, character *game.Character) {
	if mapInstance == nil {
		return
	}

	leavePacket := &response.LeavePlayer{
		ID: character.GetID(),
	}
	mapInstance.BroadcastNear(character.GetPosition(), leavePacket, nil)
}

func (l *MapListenerImpl) OnPlayerMoved(mapInstance *game.Map, character *game.Character) {
	if mapInstance == nil {
		return
	}

}

func (l *MapListenerImpl) OnPlayerMove(mapInstance *game.Map, character *game.Character, startPoint types.Vector2[int16], fragments []dto.MoveFragment) {
	if mapInstance == nil {
		return
	}

	characterDTO := character.BuildS2C()

	movePacket := &response.Move{
		Character:  characterDTO,
		Fragments:  fragments,
		StartPoint: startPoint,
	}

	character.Broadcast(movePacket, nil)
}

func (l *MapListenerImpl) OnPlayerChat(mapInstance *game.Map, character *game.Character, message string) {
	if mapInstance == nil {
		return
	}

}

func (l *MapListenerImpl) OnItemSpawned(mapInstance *game.Map, item game.Item, placement *game.FieldPlacement) {
	if mapInstance == nil {
		return
	}

	placementObj := item.GetFieldPlacement()

	spawnPacket := &response.SpawnItem{
		ID:           placement.GetOID(),
		Animation:    constant.DropItemAnimationTypeLooting,
		DropType:     placement.DropType(),
		ItemModel:    item.GetModel(),
		Expiration:   item.GetExpiration(),
		Position:     placementObj.GetPosition(),
		OwnerID:      placement.Owner(),
		SpawnedPoint: placement.SpawnedPoint(),
		IsPlayerDrop: placement.PlayerDrop(),
	}

	if placement.Quest() == 0 {
		mapInstance.BroadcastNear(placementObj.GetPosition(), spawnPacket, nil)
		return
	}

	for _, obj := range mapInstance.GetObjectsNear(placementObj.GetPosition(), constant.ObjectTypeCharacter, &game.SearchOption{IncludeHidden: true}) {
		viewer, ok := obj.(*game.Character)
		if !ok {
			continue
		}
		if !viewer.Quests().NeedsItem(placement.Quest(), item.GetModel().GetID()) {
			continue
		}
		viewer.Send(spawnPacket, types.SEND_POLICY_ENCRYPT)
	}
}

func (l *MapListenerImpl) OnMesoSpawned(mapInstance *game.Map, meso *game.Meso) {
	if mapInstance == nil {
		return
	}

	fp := meso.GetFieldPlacement()

	spawnPacket := &response.SpawnMeso{
		ID:           fp.GetOID(),
		Animation:    constant.DropItemAnimationTypeLooting,
		DropType:     fp.DropType(),
		Count:        meso.Count(),
		OwnerID:      fp.Owner(),
		Position:     fp.GetPosition(),
		SpawnedPoint: fp.SpawnedPoint(),
		IsPlayerDrop: fp.PlayerDrop(),
	}

	mapInstance.BroadcastNear(fp.GetPosition(), spawnPacket, nil)
}

func (l *MapListenerImpl) OnItemRemoved(mapInstance *game.Map, itemID uint32, looterID uint32, mode constant.RemoveItemType, position types.Vector2[int16]) {
	if mapInstance == nil {
		return
	}

	removePacket := &response.RemoveItem{
		Mode:        mode,
		OID:         itemID,
		CharacterId: looterID,
	}

	mapInstance.BroadcastNear(position, removePacket, nil)
}

func (l *MapListenerImpl) OnMobSpawned(mapInstance *game.Map, mob *game.Mob, spawnType constant.MobSpawnType, link uint32) {
	if mapInstance == nil {
		return
	}

	mobDTO := mob.BuildS2C()

	spawnPacket := &response.SpawnMob{
		Mob:       mobDTO,
		SpawnType: spawnType,
		Link:      link,
	}

	mob.Broadcast(spawnPacket, nil)
}

func (l *MapListenerImpl) OnMobRemoved(mapInstance *game.Map, mob *game.Mob, animationType constant.MobDieAnimationType) {
	if mapInstance == nil {
		return
	}

	removePacket := &response.DieMob{
		OID:           mob.GetOID(),
		AnimationType: animationType,
	}

	mapInstance.BroadcastNear(mob.GetPosition(), removePacket, nil)
}

func (l *MapListenerImpl) OnReactorSpawned(mapInstance *game.Map, reactor *game.Reactor) {
	if mapInstance == nil || reactor == nil {
		return
	}

	reactor.Broadcast(&response.SpawnReactor{
		Reactor: reactor.BuildS2C(),
	}, nil)
}

func (l *MapListenerImpl) OnReactorRemoved(mapInstance *game.Map, reactor *game.Reactor) {
	if mapInstance == nil || reactor == nil {
		return
	}

	mapInstance.BroadcastNear(reactor.GetPosition(), &response.DestroyReactor{
		Reactor: reactor.BuildS2C(),
	}, nil)
}

func (l *MapListenerImpl) OnReactorTriggered(mapInstance *game.Map, reactor *game.Reactor, stance int32) {
	if mapInstance == nil || reactor == nil {
		return
	}

	reactor.Broadcast(&response.TriggerReactor{
		Reactor: reactor.BuildS2C(),
		Stance:  stance,
	}, nil)
}

func (l *MapListenerImpl) OnMusicChanged(mapInstance *game.Map, song string) {
	if mapInstance == nil || song == "" {
		return
	}

	mapInstance.Broadcast(&response.EnvironmentChange{
		Mode: response.EnvironmentChangeModeMusic,
		Env:  song,
	}, nil)
}

func (l *MapListenerImpl) OnMapMessage(mapInstance *game.Map, messageType constant.ServerMessageType, message string) {
	if mapInstance == nil || message == "" {
		return
	}

	mapInstance.Broadcast(&response.Notice{
		Type:    messageType,
		Message: message,
	}, nil)
}

func (l *MapListenerImpl) OnWeather(mapInstance *game.Map, itemID uint32, message string) {
	mapInstance.Broadcast(&response.Weather{
		ItemID:  itemID,
		Message: message,
	}, nil)
}

func (l *MapListenerImpl) OnOXQuiz(mapInstance *game.Map, ask bool, set uint8, id uint16) {
	mapInstance.Broadcast(&response.OXQuiz{
		Ask: ask,
		Set: set,
		ID:  id,
	}, nil)
}

func (l *MapListenerImpl) OnHealerMove(mapInstance *game.Map, y int16) {
	mapInstance.Broadcast(&response.HealerMove{Y: y}, nil)
}

func (l *MapListenerImpl) OnCoconutHit(mapInstance *game.Map, id int16, delay int16, state uint8) {
	mapInstance.Broadcast(&response.CoconutHit{ID: id, Delay: delay, State: state}, nil)
}

func (l *MapListenerImpl) OnCoconutScore(mapInstance *game.Map, maple int16, story int16) {
	mapInstance.Broadcast(&response.CoconutScore{Maple: maple, Story: story}, nil)
}

func (l *MapListenerImpl) OnSnowballState(mapInstance *game.Map, state uint8, snowmanHP [2]int32, position [2]int16, stage [2]uint8, damage [3]int16) {
	mapInstance.Broadcast(&response.SnowballState{State: state, SnowmanHP: snowmanHP, Position: position, Stage: stage, Damage: damage}, nil)
}

func (l *MapListenerImpl) OnSnowballHit(mapInstance *game.Map, target uint8, damage int16, delay int16) {
	mapInstance.Broadcast(&response.SnowballHit{Target: target, Damage: damage, Delay: delay}, nil)
}

func (l *MapListenerImpl) OnSnowballMessage(mapInstance *game.Map, team uint8, message uint8) {
	mapInstance.Broadcast(&response.SnowballMessage{Team: team, Message: message}, nil)
}

func (l *MapListenerImpl) OnSongPlayed(mapInstance *game.Map, itemID uint32, name string) {
	mapInstance.Broadcast(&response.CashSong{
		ItemID: itemID,
		Name:   name,
	}, nil)
}

func (l *MapListenerImpl) OnYellowChat(mapInstance *game.Map, message string) {
	mapInstance.Broadcast(&response.YellowChat{
		Message: message,
	}, nil)
}

func (l *MapListenerImpl) OnWeddingEffect(mapInstance *game.Map, groomID uint32, brideID uint32) {
	mapInstance.Broadcast(&response.WeddingCouple{
		GroomID: groomID,
		BrideID: brideID,
	}, nil)
	mapInstance.Broadcast(&response.WeddingEffect{}, nil)
}

func (l *MapListenerImpl) OnClearEffect(mapInstance *game.Map) {
	if mapInstance == nil {
		return
	}

	mapInstance.ShowEffect("quest/party/clear")
	mapInstance.PlaySound("Party1/Clear")
	mapInstance.Broadcast(&response.EnvironmentChange{
		Mode: response.EnvironmentChangeModeObjectState,
		Env:  "gate",
	}, nil)
}

func (l *MapListenerImpl) OnShowEffect(mapInstance *game.Map, path string) {
	if mapInstance == nil || path == "" {
		return
	}
	mapInstance.Broadcast(&response.EnvironmentChange{
		Mode: response.EnvironmentChangeModeMapEffect,
		Env:  path,
	}, nil)
}

func (l *MapListenerImpl) OnPlaySound(mapInstance *game.Map, path string) {
	if mapInstance == nil || path == "" {
		return
	}
	mapInstance.Broadcast(&response.EnvironmentChange{
		Mode: response.EnvironmentChangeModeSound,
		Env:  path,
	}, nil)
}

func (l *MapListenerImpl) OnTremble(mapInstance *game.Map, typ uint8, delay int32) {
	mapInstance.Broadcast(&response.Tremble{
		Type:  typ,
		Delay: delay,
	}, nil)
}

func (l *MapListenerImpl) OnMobControllerChange(mob *game.Mob, before *game.Character, after *game.Character, aggro bool) {
	switch {
	case before == nil && after != nil:
		after.Send(&response.StartControlMob{Mob: mob.BuildS2C(), Aggro: aggro}, types.SEND_POLICY_ENCRYPT)
	case before != nil && after == nil:
		before.Send(&response.StopControlMob{OID: mob.GetOID()}, types.SEND_POLICY_ENCRYPT)
	case before != nil && after != nil:
		before.Send(&response.StopControlMob{OID: mob.GetOID()}, types.SEND_POLICY_ENCRYPT)
		after.Send(&response.StartControlMob{Mob: mob.BuildS2C(), Aggro: aggro}, types.SEND_POLICY_ENCRYPT)
	}
}

func (l *MapListenerImpl) OnMobMoved(mapInstance *game.Map, mob *game.Mob, isAggroed bool, centerSplit int8, skill1 uint8, skill2 uint8, skill3 uint8, skill4 uint8, startPoint types.Vector2[int16], movements []dto.MoveFragment) {
	if mapInstance == nil {
		return
	}

	movePacket := &response.MoveMob{
		IsAggroed:   isAggroed,
		CenterSplit: centerSplit,
		Skill1:      skill1,
		Skill2:      skill2,
		Skill3:      skill3,
		Skill4:      skill4,
		OID:         mob.GetOID(),
		StartPoint:  startPoint,
		Movements:   movements,
	}

	mob.Broadcast(movePacket, nil)
}

func (l *MapListenerImpl) OnAttack(mapInstance *game.Map, character *game.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	if character == nil || character.GetMap() == nil {
		return
	}
	character.Broadcast(&response.Attack{
		CharacterId: character.GetID(),
		AttackInfo:  attackPayload.ToAttackInfo(),
		SkillLevel:  skillLevel,
	}, nil)
}

func (l *MapListenerImpl) OnRangedAttack(mapInstance *game.Map, character *game.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	if character == nil || character.GetMap() == nil {
		return
	}
	character.Broadcast(&response.RangedAttack{
		CharacterId: character.GetID(),
		AttackInfo:  attackPayload.ToAttackInfo(),
		SkillLevel:  skillLevel,
		CashBullet:  0,
	}, nil)
}

func (l *MapListenerImpl) OnMagicAttack(mapInstance *game.Map, character *game.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	if character == nil || character.GetMap() == nil {
		return
	}
	character.Broadcast(&response.MagicAttack{
		CharacterId: character.GetID(),
		AttackInfo:  attackPayload.ToAttackInfo(),
		SkillLevel:  skillLevel,
	}, nil)
}

func (l *MapListenerImpl) OnMistSpawned(mapInstance *game.Map, mist *game.Mist) {
	if mapInstance == nil || mist == nil {
		return
	}
	skillID := uint32(0)
	if mist.SkillWz() != nil {
		skillID = mist.SkillWz().ID
	}
	pkt := &response.SpawnMist{
		OID:        mist.GetOID(),
		Type:       mist.MistType(),
		MobMist:    mist.MobMist(),
		CauserID:   mist.Causer(),
		SkillID:    skillID,
		SkillLevel: mist.SkillLevel(),
		SkillDelay: mist.SkillDelay(),
		Bounds:     mist.Bounds(),
		MobSkill:   mist.MobMist(),
	}
	mist.Broadcast(pkt, nil)
}

func (l *MapListenerImpl) OnMistRemoved(mapInstance *game.Map, mist *game.Mist) {
	if mapInstance == nil || mist == nil {
		return
	}
	pkt := &response.RemoveMist{
		OID:      mist.GetOID(),
		Eruption: false,
	}
	mist.Broadcast(pkt, nil)
}

func (l *MapListenerImpl) OnMessageBoxSpawned(mapInstance *game.Map, box *game.MessageBox) {
	box.Broadcast(&response.SpawnMessageBox{
		OID:    box.GetOID(),
		ItemID: box.ItemID(),
		Text:   box.Text(),
		Owner:  box.Owner(),
		X:      box.GetPosition().X,
		Y:      box.GetPosition().Y,
	}, nil)
}

func (l *MapListenerImpl) OnMessageBoxRemoved(mapInstance *game.Map, box *game.MessageBox) {
	box.Broadcast(&response.RemoveMessageBox{OID: box.GetOID(), Animated: true}, nil)
}

func (l *MapListenerImpl) OnDoorRemoved(mapInstance *game.Map, door *game.Door, animated bool) {
	if mapInstance == nil || door == nil {
		return
	}
	door.Broadcast(&response.RemoveDoor{
		OwnerID:  door.OwnerID(),
		Animated: animated,
	}, nil)
	door.Broadcast(&response.SpawnPortal{
		DestMapID:   response.DisabledPortalMapID,
		SourceMapID: response.DisabledPortalMapID,
		Position:    nil,
	}, nil)
}

func (l *MapListenerImpl) OnMobHomingRemoved(mapInstance *game.Map, mob *game.Mob, removed *game.Homing, causer *game.Character) {
	if causer == nil {
		return
	}
	causer.Send(&response.CancelHomingBeacon{}, types.SEND_POLICY_ENCRYPT)
}

func (l *MapListenerImpl) OnMobHomingSet(mapInstance *game.Map, mob *game.Mob, homing *game.Homing, causer *game.Character) {
	if causer == nil {
		return
	}
	x := int32(1)
	if ld := homing.SkillWz().GetLevelData(int(homing.SkillLevel())); ld != nil && ld.X > 0 {
		x = int32(ld.X)
	}
	causer.Send(&response.CancelHomingBeacon{}, types.SEND_POLICY_ENCRYPT)
	causer.Send(&response.GiveHomingBeacon{
		SkillID: homing.SkillWz().ID,
		MobOID:  mob.GetOID(),
		X:       x,
	}, types.SEND_POLICY_ENCRYPT)
}
