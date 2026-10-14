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

func (l *CharacterListenerImpl) CompleteWhisperAsync(ctx actor.Context, ch *entity.Character, requestID string, received bool) *async.Promise[*internal.CompleteWhisperReply] {
	req := &internal.CompleteWhisperRequest{
		WorldId:     l.gs.config.WorldId,
		RequestId:   requestID,
		CharacterId: ch.GetID(),
		MapId:       ch.GetMap().TemplateID(),
		Complete:    true,
		Received:    received,
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

func (l *CharacterListenerImpl) SendNoteAsync(ctx actor.Context, ch *entity.Character, recipient string, message string) *async.Promise[*internal.NotesReply] {
	worldID := l.gs.config.WorldId
	characterID := ch.GetID()
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.NotesReply, error) {
			return l.gs.internalClient.Notes(c, &internal.NotesRequest{
				WorldId:     worldID,
				CharacterId: characterID,
				Action:      internal.NotesAction_NOTES_ACTION_SEND,
				TargetName:  recipient,
				Message:     message,
				Gift:        uint32(pconst.NoteKindNormal),
			})
		},
		nil,
	)
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
	id := ch.GetID()
	self := ch.Family.Self()
	info := &response.FamilyInfo{
		Reputation:      self.GetReputation(),
		TotalReputation: self.GetTotalReputation(),
		Juniors:         uint16(len(ch.Family.Tree.Juniors(id))),
		Precept:         ch.Family.Tree.Precept(),
		UsedFeatures:    self.GetUsedFeatures(),
	}
	if leader := ch.Family.Tree.Leader(); leader != nil {
		info.LeaderID = leader.GetCharacterId()
		info.LeaderName = leader.GetName()
	}
	ch.Send(info, types.SEND_POLICY_ENCRYPT)

	remaining := ch.Family.BuffRemaining()
	if remaining > 0 {
		ch.Send(&response.FamilyBuff{
			Feature:   self.GetBuffFeature(),
			ExpRate:   ch.Family.ExpRate(),
			DropRate:  ch.Family.DropRate(),
			Remaining: remaining,
		}, types.SEND_POLICY_ENCRYPT)
	}
}

func (l *CharacterListenerImpl) GetFamilyAsync(ctx actor.Context, ch *entity.Character, name string) *async.Promise[*internal.GetFamilyReply] {
	req := &internal.GetFamilyRequest{
		WorldId: l.gs.config.WorldId,
		Name:    name,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.GetFamilyReply, error) {
			return l.gs.internalClient.GetFamily(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) JoinFamilyAsync(ctx actor.Context, senior *entity.Character, junior *entity.Character) *async.Promise[*internal.FamilyReply] {
	req := &internal.JoinFamilyRequest{
		WorldId: l.gs.config.WorldId,
		Senior: &internal.FamilyMember{
			CharacterId: senior.GetID(),
			Name:        senior.GetName(),
			Level:       uint32(senior.GetLevel()),
			ClassId:     uint32(senior.Class),
		},
		Junior: &internal.FamilyMember{
			CharacterId: junior.GetID(),
			Name:        junior.GetName(),
			Level:       uint32(junior.GetLevel()),
			ClassId:     uint32(junior.Class),
		},
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.FamilyReply, error) {
			return l.gs.internalClient.JoinFamily(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) SeparateFamilyAsync(ctx actor.Context, ch *entity.Character, juniorID uint32) *async.Promise[*internal.FamilyReply] {
	req := &internal.SeparateFamilyRequest{
		WorldId:     l.gs.config.WorldId,
		CharacterId: ch.GetID(),
		JuniorId:    juniorID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.FamilyReply, error) {
			return l.gs.internalClient.SeparateFamily(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) SetFamilyPreceptAsync(ctx actor.Context, ch *entity.Character, precept string) *async.Promise[*internal.FamilyReply] {
	req := &internal.SetFamilyPreceptRequest{
		WorldId:     l.gs.config.WorldId,
		CharacterId: ch.GetID(),
		Precept:     precept,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.FamilyReply, error) {
			return l.gs.internalClient.SetFamilyPrecept(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) UseFamilyFeatureAsync(ctx actor.Context, ch *entity.Character, feature internal.FamilyFeature, targetID uint32, requestID uint64) *async.Promise[*internal.FamilyReply] {
	req := &internal.UseFamilyFeatureRequest{
		WorldId:     l.gs.config.WorldId,
		CharacterId: ch.GetID(),
		Feature:     feature,
		TargetId:    targetID,
		RequestId:   requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.FamilyReply, error) {
			return l.gs.internalClient.UseFamilyFeature(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) FinishFamilyMoveAsync(ctx actor.Context, ch *entity.Character, requestID uint64, complete bool) *async.Promise[*internal.FamilyReply] {
	req := &internal.FinishFamilyMoveRequest{
		WorldId:     l.gs.config.WorldId,
		CharacterId: ch.GetID(),
		RequestId:   requestID,
		Complete:    complete,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.FamilyReply, error) {
			return l.gs.internalClient.FinishFamilyMove(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) OnFamilyChart(ch *entity.Character, tree entity.FamilyTree, characterID uint32) {
	chart := &response.FamilyChart{
		CharacterID:  characterID,
		Members:      tree.Pedigree(characterID),
		MemberCount:  tree.Size(),
		UsedFeatures: tree.Member(characterID).GetUsedFeatures(),
	}
	for _, junior := range tree.Juniors(characterID) {
		for _, grand := range tree.Juniors(junior.GetCharacterId()) {
			chart.GrandJuniors = append(chart.GrandJuniors, response.FamilyDescendants{
				CharacterID: grand.GetCharacterId(),
				Count:       tree.Descendants(grand.GetCharacterId()),
			})
		}
	}
	ch.Send(chart, types.SEND_POLICY_ENCRYPT)
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
