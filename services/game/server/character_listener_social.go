package server

import (
	"context"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) OnWhisperFailed(ch *entity.Character, name string) {
	ch.Send(&response.Whisper{
		Mode: internal.WhisperMode_WHISPER_MODE_SEND_RESULT,
		Name: name,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnWhisperSent(ch *entity.Character, name string) {
	ch.Send(&response.Whisper{
		Mode:    internal.WhisperMode_WHISPER_MODE_SEND_RESULT,
		Name:    name,
		Success: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnWhisperFound(ch *entity.Character, reply *internal.WhisperReply, buddy bool) {
	packet := &response.Whisper{
		Mode:     internal.WhisperMode_WHISPER_MODE_FIND_RESULT,
		Name:     reply.Name,
		Location: pconst.WhisperLocationChannel,
		Channel:  reply.ChannelId,
	}
	if buddy {
		packet.Mode = internal.WhisperMode_WHISPER_MODE_FIND_BUDDY_RESULT
	}
	if reply.CashShop {
		packet.Location = pconst.WhisperLocationCashShop
		packet.Channel = -1
	} else if reply.ChannelId == int32(l.gs.config.ChannelId) {
		packet.Location = pconst.WhisperLocationMap
		packet.MapID = reply.MapId
	}
	ch.Send(packet, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnWhisperReceived(ch *entity.Character, sender string, channel int32, message string) error {
	return ch.Send(&response.Whisper{
		Mode:    internal.WhisperMode_WHISPER_MODE_RECEIVE,
		Name:    sender,
		Channel: channel,
		Message: message,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) ClaimWhisperAsync(ctx actor.Context, ch *entity.Character, requestID string) *async.Promise[*internal.CompleteWhisperReply] {
	req := &internal.CompleteWhisperRequest{
		WorldId:     l.gs.config.WorldId,
		RequestId:   requestID,
		CharacterId: ch.GetID(),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.CompleteWhisperReply, error) {
			return l.gs.internalClient.CompleteWhisper(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) CompleteWhisperAsync(ctx actor.Context, ch *entity.Character, requestID string, delivered bool) *async.Promise[*internal.CompleteWhisperReply] {
	req := &internal.CompleteWhisperRequest{
		WorldId:     l.gs.config.WorldId,
		RequestId:   requestID,
		CharacterId: ch.GetID(),
		MapId:       ch.GetMap().TemplateID(),
		Complete:    true,
		Delivered:   delivered,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.CompleteWhisperReply, error) {
			return l.gs.internalClient.CompleteWhisper(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) OnMessengerJoined(ch *entity.Character, position uint8, members []*internal.MessengerMember, queued []entity.MessengerEvent) {
	ch.Send(&response.Messenger{
		Action:   internal.MessengerAction_MESSENGER_ACTION_SELF_ENTER,
		Position: position,
	}, types.SEND_POLICY_ENCRYPT)
	for _, m := range members {
		ch.Send(&response.Messenger{
			Action:   internal.MessengerAction_MESSENGER_ACTION_ENTER,
			Position: uint8(m.Position),
			Name:     m.Name,
			Channel:  uint16(m.ChannelId),
			Look:     m.Look,
		}, types.SEND_POLICY_ENCRYPT)
	}
	for _, event := range queued {
		l.OnMessengerEvent(ch, event)
	}
}

func (l *CharacterListenerImpl) OnMessengerInviteResult(ch *entity.Character, name string, success bool) {
	ch.Send(&response.Messenger{
		Action:  internal.MessengerAction_MESSENGER_ACTION_INVITE_RESULT,
		Name:    name,
		Success: success,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMessengerEvent(ch *entity.Character, event entity.MessengerEvent) {
	ch.Send(&response.Messenger{
		Action:   event.Action,
		Position: event.Position,
		Name:     event.Name,
		Text:     event.Text,
		RoomID:   event.RoomID,
		Channel:  event.Channel,
		Look:     event.Look,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPopulationFailed(ch *entity.Character, result pconst.PopulationResult) {
	ch.Send(&response.Population{
		Status: result,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPopulationGiven(ch *entity.Character, mode pconst.PopulationMode, targetName string, value int32) {
	ch.Send(&response.Population{
		Status: pconst.PopulationResultGiven,
		Mode:   mode,
		Name:   targetName,
		Value:  value,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPopulationReceived(ch *entity.Character, mode pconst.PopulationMode, giverName string) {
	ch.Send(&response.Population{
		Status: pconst.PopulationResultReceived,
		Mode:   mode,
		Name:   giverName,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnNotes(ch *entity.Character, notes []*internal.Note) {
	packet := &response.Notes{Action: pconst.NotesPacketList}
	for _, n := range notes {
		packet.Notes = append(packet.Notes, response.Note{
			ID:      n.Id,
			Sender:  n.SenderName,
			Message: n.Message,
			SentAt:  time.UnixMilli(n.SentAtUnixMs),
			Kind:    pconst.NoteKind(n.Gift),
		})
	}
	ch.Send(packet, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnNoteSent(ch *entity.Character) {
	ch.Send(&response.Notes{
		Action: pconst.NotesPacketSent,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnNoteFailed(ch *entity.Character, result internal.NoteResult) {
	ch.Send(&response.Notes{
		Action: pconst.NotesPacketFailed,
		Status: result,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnNoteArrived(ch *entity.Character) {
	ch.Send(&response.Notes{
		Action: pconst.NotesPacketRefresh,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFamilyPrivileges(ch *entity.Character) {
	ch.Send(&response.FamilyPrivileges{}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFamilyChanged(ch *entity.Character) {
	ch.Send(&response.FamilyInfo{
		State: ch.Family.State,
	}, types.SEND_POLICY_ENCRYPT)

	remaining := ch.Family.BuffRemaining()
	if remaining > 0 {
		ch.Send(&response.FamilyBuff{
			Feature:   ch.Family.State.BuffFeature,
			ExpRate:   ch.Family.ExpRate(),
			DropRate:  ch.Family.DropRate(),
			Remaining: remaining,
		}, types.SEND_POLICY_ENCRYPT)
	}
}

func (l *CharacterListenerImpl) LoadFamilyAsync(ctx actor.Context, ch *entity.Character) *async.Promise[*internal.FamilyReply] {
	req := &internal.FamilyRequest{
		WorldId:     l.gs.config.WorldId,
		CharacterId: ch.GetID(),
		ChannelId:   l.gs.config.ChannelId,
		Action:      internal.FamilyAction_FAMILY_ACTION_OPEN,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.FamilyReply, error) {
			return l.gs.internalClient.Family(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) OnFamilyChart(ch *entity.Character, state *internal.FamilyState) {
	ch.Send(&response.FamilyChart{
		State: state,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFamilyResult(ch *entity.Character, result internal.FamilyResult) {
	ch.Send(&response.FamilyResult{
		Code: result,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFamilyInvite(ch *entity.Character, inviterID uint32, name string) {
	ch.Send(&response.FamilyInvite{
		InviterID: inviterID,
		Name:      name,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFamilyInviteResult(ch *entity.Character, accepted bool, name string) {
	ch.Send(&response.FamilyInviteResult{
		Accepted: accepted,
		Name:     name,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFamilyJoined(ch *entity.Character, name string) {
	ch.Send(&response.FamilyJoined{
		Name: name,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFamilySummon(ch *entity.Character, name string, mapName string) {
	ch.Send(&response.FamilySummon{
		Name:    name,
		MapName: mapName,
	}, types.SEND_POLICY_ENCRYPT)
}
