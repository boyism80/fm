package server

import (
	"context"
	"fmt"
	"strconv"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/rpc"
	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type CharacterListenerImpl struct {
	gs *GameServer
}

func (l *CharacterListenerImpl) OnChat(ch *game.Character, message string, highlight bool, dontRecordHistory bool) {
	chatPacket := &response.NormalChat{
		CharacterId:       ch.GetID(),
		Message:           message,
		Highlight:         highlight,
		DontRecordHistory: dontRecordHistory,
	}

	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(chatPacket, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnMesoChanged(ch *game.Character, meso int32) {
	ch.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.StatMeso: meso,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnKeyMap(ch *game.Character) {
	ch.Send(&response.KeyMap{Slots: ch.KeyLayout().Bindings()}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnSkillMacros(ch *game.Character) {
	macros := ch.Skills().Macros().Entries()
	entries := make([]response.SkillMacroEntry, 0, len(macros))
	for _, macro := range macros {
		entries = append(entries, response.SkillMacroEntry{Name: macro.Name, Shout: macro.Shout, Skills: macro.Skills})
	}
	ch.Send(&response.SkillMacro{Entries: entries}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnChalkboard(ch *game.Character) {
	ch.Broadcast(&response.Chalkboard{CharacterID: ch.GetID(), Text: ch.Chalkboard()}, nil)
}

func (l *CharacterListenerImpl) OnFaceExpression(ch *game.Character, expression uint32) {
	ch.Broadcast(&response.FaceExpression{CharacterID: ch.GetID(), Expression: expression}, nil)
}

func (l *CharacterListenerImpl) OnMessage(ch *game.Character, messageType constant.ServerMessageType, message string) {
	l.OnNotice(ch, messageType, message, 0, false)
}

func (l *CharacterListenerImpl) OnNotice(ch *game.Character, messageType constant.ServerMessageType, message string, channel int, ear bool) {
	ch.Send(&response.Notice{
		Message: message,
		Type:    messageType,
		Channel: channel,
		MegaEar: ear,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMegaphone(ch *game.Character, megaphone *response.Megaphone) {
	ch.Send(megaphone, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMessageBoxFailed(ch *game.Character) {
	ch.Send(&response.MessageBoxFailed{}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) Disconnect(ch *game.Character) {
	c, ok := ch.Client().(*client.GameClient)
	if !ok {
		return
	}
	_ = c.GetConnection().Close()
}

func (l *CharacterListenerImpl) OnClock(ch *game.Character, seconds int32) {
	ch.Send(&response.Clock{
		Seconds: seconds,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStopClock(ch *game.Character) {
	ch.Send(&response.StopClock{}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnEventInstructions(ch *game.Character) {
	ch.Send(&response.EventInstructions{}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnExpGain(ch *game.Character, exp uint32) {
	expPacket := &response.GainExp{
		Gain:  exp,
		White: false,
	}
	ch.Send(expPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDojoEnergy(ch *game.Character, energy int) {
	ch.Send(&response.SessionValue{
		Key:   "energy",
		Value: strconv.Itoa(energy),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShipState(ch *game.Character, state uint16) {
	ch.Send(&response.ShipState{State: state}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShipBalrog(ch *game.Character) {
	ch.Send(&response.ShipSpecialEffect{Effect: response.ShipSpecialBalrog}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnSitOnChair(ch *game.Character, itemID uint32) {
	ch.Broadcast(&response.ShowChair{
		CharacterID: ch.GetID(),
		ItemID:      itemID,
	}, nil)
}

func (l *CharacterListenerImpl) OnSitOnMapSeat(ch *game.Character, seatID int16) {
	ch.Send(&response.CancelChair{
		ChairID: seatID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStandUp(ch *game.Character) {
	ch.Send(&response.CancelChair{
		ChairID: -1,
	}, types.SEND_POLICY_ENCRYPT)
	ch.Broadcast(&response.ShowChair{
		CharacterID: ch.GetID(),
		ItemID:      0,
	}, nil)
}

func (l *CharacterListenerImpl) OnSnowballKnockback(ch *game.Character) {
	_ = ch.Send(&response.SnowballKnockback{}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnIncubatorResult(ch *game.Character, itemID uint32, count int16) {
	_ = ch.Send(&response.IncubatorResult{
		ItemID: itemID,
		Count:  count,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnUnlockAction(ch *game.Character) {
	ch.Send(&response.UpdateStats{
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnInspect(ch *game.Character, profile *response.CharacterProfile) {
	ch.Send(profile, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMonsterBookCardRegistered(ch *game.Character, cardID uint32, count uint32) {
	ch.Send(&response.MonsterBookSetCard{Success: true, CardID: cardID, Count: count}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMonsterBookCardFull(ch *game.Character) {
	ch.Send(&response.MonsterBookSetCard{}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMonsterBookCover(ch *game.Character, cardID uint32) {
	ch.Send(&response.MonsterBookSetCover{CardID: cardID}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnTeleportStones(ch *game.Character, vip bool) {
	ch.Send(&response.TeleportStoneResult{Result: pconst.TeleportStoneResultList, VIP: vip, Maps: ch.TeleportStones().Slots(vip)}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnTeleportStoneFailed(ch *game.Character, vip bool, result pconst.TeleportStoneResult) {
	ch.Send(&response.TeleportStoneResult{Result: result, VIP: vip}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnScriptError(ch *game.Character, script string, err error) {
	if ch.GetRole() != constant.RoleAdmin {
		return
	}
	l.OnMessage(ch, constant.MsgPinkText, fmt.Sprintf("스크립트 오류 %s: %v", script, err))
}

func (l *CharacterListenerImpl) OnUpdateStats(ch *game.Character, stats map[constant.Stat]int32, unlock bool) {
	ch.Send(&response.UpdateStats{
		Stats:        stats,
		UnlockAction: unlock,
	}, types.SEND_POLICY_ENCRYPT)
	if stats != nil {
		_, hasHP := stats[constant.StatHP]
		_, hasMaxHP := stats[constant.StatMaxHP]
		if m := ch.GetMap(); m != nil && (hasHP || hasMaxHP) {
			m.BroadcastPartyHP(ch)
		}
	}
}

func (l *CharacterListenerImpl) OnPlayerMove(ch *game.Character, startPoint types.Vector2[int16], fragments []dto.MoveFragment) {
	if ch.GetMap() == nil {
		return
	}

	characterDTO := ch.BuildS2C()

	movePacket := &response.Move{
		Character:  characterDTO,
		Fragments:  fragments,
		StartPoint: startPoint,
	}

	ch.Broadcast(movePacket, nil)
}

func (l *CharacterListenerImpl) OnFieldRelocate(ch *game.Character, spawnPoint uint8) {
	ch.Send(&response.FieldRelocate{Portal: spawnPoint}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnUpdateCharacterLook(ch *game.Character) {
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}

	lookPacket := &response.UpdateCharacterLook{
		Character:      ch.BuildS2C(),
		CrushRing:      ch.Inventory().WornRing(ch.Inventory().Rings().Left()),
		FriendshipRing: ch.Inventory().WornRing(ch.Inventory().Rings().Mid()),
		MarriageRing:   ch.Wedding().BuildRingS2C(),
	}

	ch.Broadcast(lookPacket, nil)
	if ch.Messenger().RoomID() != 0 {
		l.gs.systems.Get[*game.DispatchSystem]().Visit(ch.GetID(), func(ctx actor.Context, current *game.Character) {
			if current.Messenger().RoomID() == 0 {
				return
			}
			w := stream.NewStreamWriter(stream.LittleEndian)
			current.BuildS2C().SerializeLook(w)
			req := &internal.MessengerRequest{
				WorldId:     l.gs.config.WorldId,
				CharacterId: current.GetID(),
				ChannelId:   l.gs.config.ChannelId,
				Action:      internal.MessengerAction_MESSENGER_ACTION_AVATAR,
				Look:        w.Bytes(),
			}
			async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
				func(c context.Context) (*internal.MessengerReply, error) {
					return l.gs.systems.Get[*rpc.InternalConnection]().Client.Messenger(c, req)
				},
				nil,
			).OnError(func(error) {})
		})
	}
}

func (l *CharacterListenerImpl) OnNpcAction(ch *game.Character, bytes []byte) {
	ch.Send(&response.NpcAction{
		Bytes: bytes,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnClassChange(ch *game.Character, oldClass constant.ClassType, newClass constant.ClassType, skillPoint uint16, abilityPoint uint16) {
	stats := map[constant.Stat]int32{
		constant.StatClass:       int32(newClass),
		constant.StatAvailableSP: int32(skillPoint),
		constant.StatAvailableAP: int32(abilityPoint),
	}

	ch.Send(&response.UpdateStats{
		Stats:        stats,
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnHiddenChanged(ch *game.Character, hidden bool) {
	ch.Send(&response.SuperHide{Hidden: hidden}, types.SEND_POLICY_ENCRYPT)

	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}

	if hidden {
		ch.Broadcast(&response.LeavePlayer{ID: ch.GetID()}, &game.ObjectBroadcastOption{
			RecipientsRoleBelowPivot: true,
		})
	} else {
		for _, obj := range mapInstance.GetObjectsNear(ch.GetPosition(), constant.ObjectTypeCharacter, &game.SearchOption{IncludeHidden: true}) {
			if viewer, ok := obj.(*game.Character); ok {
				ch.SendSpawnSyncToViewer(viewer)
			}
		}
	}
	if controllerTable := mapInstance.GetControllerTable(); controllerTable != nil {
		controllerTable.Update(ch)
	}
}

func (l *CharacterListenerImpl) OnLieDetectorStarted(ch *game.Character, target string) {
	ch.Send(&response.LieDetectorStarted{Target: target}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnLieDetectorQuestion(ch *game.Character, image []byte, skill bool) {
	ch.Send(&response.LieDetectorQuestion{Image: image, Skill: skill}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnLieDetectorResult(ch *game.Character, result pconst.LieDetectorResult, skill bool) {
	ch.Send(&response.LieDetectorResult{Result: result, Skill: skill}, types.SEND_POLICY_ENCRYPT)
}
