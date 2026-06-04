package server

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
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

func (l *MapListenerImpl) OnPlayerAdded(ctx actor.Context, mapInstance *entity.Map, character *entity.Character, init bool) {
	if mapInstance == nil {
		return
	}

	if init {
		characterDTO := character.ToFullDTO()
		loginPacket := &response.Login{
			Channel:   l.gs.config.ChannelId,
			Character: characterDTO,
		}
		character.Send(loginPacket, types.SEND_POLICY_ENCRYPT)
		if character.Listener != nil {
			character.Listener.OnShowGuildInfo(character)
			character.Listener.OnShowAllianceInfo(character)
		}
		character.Buffs.EmitAllBuffAddedEvents()
	} else {
		characterDTO := character.ToDTO()
		warpPacket := &response.Warp{
			Character: characterDTO,
			Channel:   0,
		}
		character.Send(warpPacket, types.SEND_POLICY_ENCRYPT)
	}
	if character.IsHidden() {
		character.Send(&response.SuperHide{Hidden: true}, types.SEND_POLICY_ENCRYPT)
	}

	if kl := character.KeyLayout(); kl != nil {
		character.Send(&response.KeyMap{Slots: kl.Bindings()}, types.SEND_POLICY_ENCRYPT)
	}

	for _, obj := range mapInstance.GetObjects(constant.ObjectTypeCharacter) {
		if viewer, ok := obj.(*entity.Character); ok {
			character.SendSpawnSyncToViewer(viewer)
		}
	}

	for _, obj := range mapInstance.GetObjects(constant.ObjectTypeObject) {
		obj.SendSpawnSyncToViewer(character)
	}

	SyncPartyMemberHPOnMapEnter(mapInstance, character, nil)

	if l.gs != nil && character != nil && l.gs.characterRuntime != nil {
		_ = l.gs.characterRuntime.SetMapPID(character.GetID(), mapInstance.GetActorPID())
	}

	if l.gs != nil {
		l.gs.party.SendPartySilentAsync(ctx, character).Run()
	}
}

func SyncPartyMemberHPOnMapEnter(mapInstance *entity.Map, character *entity.Character, effectivePartyID *uint32) {
	if mapInstance == nil || character == nil {
		return
	}
	var partyID uint32
	inParty := false
	if effectivePartyID != nil {
		partyID = *effectivePartyID
		inParty = true
	} else if pid := character.GetPartyID(); pid != nil {
		partyID = *pid
		inParty = true
	}
	if !inParty {
		return
	}
	character.Listener.OnPartyMemberHPChanged(character, nil)
	for _, obj := range mapInstance.GetObjects(constant.ObjectTypeCharacter) {
		peer, ok := obj.(*entity.Character)
		if !ok || peer == nil || peer.GetID() == character.GetID() {
			continue
		}
		pPeer := peer.GetPartyID()
		if pPeer == nil || *pPeer != partyID {
			continue
		}
		peer.Listener.OnPartyMemberHPChanged(peer, character)
	}
}

func (l *MapListenerImpl) OnPlayerRemoved(mapInstance *entity.Map, character *entity.Character) {
	if mapInstance == nil {
		return
	}

	leavePacket := &response.LeavePlayer{
		ID: character.GetID(),
	}
	mapInstance.Broadcast(leavePacket, nil)

	if l.gs != nil && character != nil && l.gs.characterRuntime != nil {
		_ = l.gs.characterRuntime.SetMapPID(character.GetID(), nil)
	}
}

func (l *MapListenerImpl) OnPlayerMoved(mapInstance *entity.Map, character *entity.Character) {
	if mapInstance == nil {
		return
	}

}

func (l *MapListenerImpl) OnPlayerMove(mapInstance *entity.Map, character *entity.Character, startPoint types.Vector2[int16], fragments []dto.MoveFragment) {
	if mapInstance == nil {
		return
	}

	characterDTO := character.ToDTO()

	movePacket := &response.Move{
		Character:  characterDTO,
		Fragments:  fragments,
		StartPoint: startPoint,
	}

	character.Broadcast(movePacket, nil)
}

func (l *MapListenerImpl) OnPlayerChat(mapInstance *entity.Map, character *entity.Character, message string) {
	if mapInstance == nil {
		return
	}

}

func (l *MapListenerImpl) OnItemSpawned(mapInstance *entity.Map, item entity.Item, placement *entity.FieldPlacement) {
	if mapInstance == nil {
		return
	}

	placementObj := item.GetFieldPlacement()

	spawnPacket := &response.SpawnItem{
		ID:           placement.OID,
		Animation:    constant.DropItemAnimationTypeLooting,
		DropType:     placement.DropType,
		ItemModel:    item.GetModel(),
		Expiration:   item.GetExpiration(),
		Position:     placementObj.Position,
		OwnerID:      placement.Owner,
		SpawnedPoint: placement.SpawnedPoint,
		IsPlayerDrop: true,
	}

	mapInstance.Broadcast(spawnPacket, nil)
}

func (l *MapListenerImpl) OnMesoSpawned(mapInstance *entity.Map, meso *entity.Meso) {
	if mapInstance == nil {
		return
	}

	fp := meso.GetFieldPlacement()

	spawnPacket := &response.SpawnMeso{
		ID:           fp.OID,
		Animation:    constant.DropItemAnimationTypeLooting,
		DropType:     fp.DropType,
		Count:        meso.Count,
		OwnerID:      fp.Owner,
		Position:     fp.Position,
		SpawnedPoint: fp.SpawnedPoint,
		IsPlayerDrop: true,
	}

	mapInstance.Broadcast(spawnPacket, nil)
}

func (l *MapListenerImpl) OnItemRemoved(mapInstance *entity.Map, itemID uint32, looterID uint32, mode constant.RemoveItemType) {
	if mapInstance == nil {
		return
	}

	removePacket := &response.RemoveItem{
		Mode:        mode,
		OID:         itemID,
		CharacterId: looterID,
	}

	mapInstance.Broadcast(removePacket, nil)
}

func (l *MapListenerImpl) OnMobSpawned(mapInstance *entity.Map, mob *entity.Mob, spawnType constant.MobSpawnType, link uint32) {
	if mapInstance == nil {
		return
	}

	mobDTO := mob.ToDTO()

	spawnPacket := &response.SpawnMob{
		Mob:       mobDTO,
		SpawnType: spawnType,
		Link:      link,
	}

	mapInstance.Broadcast(spawnPacket, nil)
}

func (l *MapListenerImpl) OnMobRemoved(mapInstance *entity.Map, mob *entity.Mob, animationType constant.MobDieAnimationType) {
	if mapInstance == nil {
		return
	}

	removePacket := &response.DieMob{
		OID:           mob.OID,
		AnimationType: animationType,
	}

	mapInstance.Broadcast(removePacket, nil)
}

func (l *MapListenerImpl) OnMobControllerChange(mob *entity.Mob, before *entity.Character, after *entity.Character, aggro bool) {
	switch {
	case before == nil && after != nil:
		after.Send(&response.StartControlMob{Mob: mob.ToDTO(), Aggro: aggro}, types.SEND_POLICY_ENCRYPT)
	case before != nil && after == nil:
		before.Send(&response.StopControlMob{OID: mob.OID}, types.SEND_POLICY_ENCRYPT)
	case before != nil && after != nil:
		before.Send(&response.StopControlMob{OID: mob.OID}, types.SEND_POLICY_ENCRYPT)
		after.Send(&response.StartControlMob{Mob: mob.ToDTO(), Aggro: aggro}, types.SEND_POLICY_ENCRYPT)
	}
}

func (l *MapListenerImpl) OnMobMoved(mapInstance *entity.Map, mob *entity.Mob, isAggroed bool, centerSplit int8, skill1 uint8, skill2 uint8, skill3 uint8, skill4 uint8, startPoint types.Vector2[int16], movements []dto.MoveFragment) {
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
		OID:         mob.OID,
		StartPoint:  startPoint,
		Movements:   movements,
	}

	mob.Broadcast(movePacket, nil)
}

func (l *MapListenerImpl) broadcastAttack(character *entity.Character, packet types.Packet) {
	if character == nil || character.GetMap() == nil {
		return
	}

	character.Broadcast(packet, nil)
}

func (l *MapListenerImpl) OnAttack(mapInstance *entity.Map, character *entity.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	l.broadcastAttack(character, &response.Attack{
		CharacterId: character.GetID(),
		AttackInfo:  attackPayload.ToAttackInfo(),
		SkillLevel:  skillLevel,
	})
}

func (l *MapListenerImpl) OnRangedAttack(mapInstance *entity.Map, character *entity.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	l.broadcastAttack(character, &response.RangedAttack{
		CharacterId: character.GetID(),
		AttackInfo:  attackPayload.ToAttackInfo(),
		SkillLevel:  skillLevel,
		CashBullet:  0,
	})
}

func (l *MapListenerImpl) OnMagicAttack(mapInstance *entity.Map, character *entity.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	l.broadcastAttack(character, &response.MagicAttack{
		CharacterId: character.GetID(),
		AttackInfo:  attackPayload.ToAttackInfo(),
		SkillLevel:  skillLevel,
	})
}

func (l *MapListenerImpl) OnMistSpawned(mapInstance *entity.Map, mist *entity.Mist) {
	if mapInstance == nil || mist == nil {
		return
	}
	skillID := uint32(0)
	if mist.SkillWz != nil {
		skillID = mist.SkillWz.ID
	}
	pkt := &response.SpawnMist{
		OID:        mist.OID,
		Type:       mist.MistType,
		MobMist:    mist.MobMist,
		CauserID:   mist.Causer,
		SkillID:    skillID,
		SkillLevel: mist.SkillLevel,
		SkillDelay: mist.SkillDelay,
		Bounds:     mist.Bounds,
		MobSkill:   mist.MobSkill,
	}
	mist.Broadcast(pkt, nil)
}

func (l *MapListenerImpl) OnMistRemoved(mapInstance *entity.Map, mist *entity.Mist) {
	if mapInstance == nil || mist == nil {
		return
	}
	pkt := &response.RemoveMist{
		OID:      mist.OID,
		Eruption: false,
	}
	mist.Broadcast(pkt, nil)
}

func (l *MapListenerImpl) OnDoorRemoved(mapInstance *entity.Map, door *entity.Door, animated bool) {
	if mapInstance == nil || door == nil {
		return
	}
	door.Broadcast(&response.RemoveDoor{
		OwnerID:  door.OwnerID,
		Animated: animated,
	}, nil)
	door.Broadcast(&response.SpawnPortal{
		DestMapID:   response.DisabledPortalMapID,
		SourceMapID: response.DisabledPortalMapID,
		SkillID:     0,
		Position:    nil,
	}, nil)
}

func (l *MapListenerImpl) OnMobHomingRemoved(mapInstance *entity.Map, mob *entity.Mob, removed *entity.Homing, causer *entity.Character) {
	if causer == nil {
		return
	}
	causer.Send(&response.CancelHomingBeacon{}, types.SEND_POLICY_ENCRYPT)
}

func (l *MapListenerImpl) OnMobHomingSet(mapInstance *entity.Map, mob *entity.Mob, homing *entity.Homing, causer *entity.Character) {
	if causer == nil {
		return
	}
	x := int32(1)
	if ld := homing.SkillWz.GetLevelData(int(homing.SkillLevel)); ld != nil && ld.X > 0 {
		x = int32(ld.X)
	}
	causer.Send(&response.CancelHomingBeacon{}, types.SEND_POLICY_ENCRYPT)
	causer.Send(&response.GiveHomingBeacon{
		SkillID: homing.SkillWz.ID,
		MobOID:  mob.GetOID(),
		X:       x,
	}, types.SEND_POLICY_ENCRYPT)
}
