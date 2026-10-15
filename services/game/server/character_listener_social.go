package server

import (
	"time"

	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) OnWhisperFailed(ch *game.Character, name string) {
	ch.Send(&response.Whisper{
		Mode: internal.WhisperMode_WHISPER_MODE_SEND_RESULT,
		Name: name,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnWhisperSent(ch *game.Character, name string) {
	ch.Send(&response.Whisper{
		Mode:    internal.WhisperMode_WHISPER_MODE_SEND_RESULT,
		Name:    name,
		Success: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnWhisperFound(ch *game.Character, reply *internal.WhisperReply, buddy bool) {
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

func (l *CharacterListenerImpl) OnWhisperReceived(ch *game.Character, sender string, channel int32, message string) error {
	return ch.Send(&response.Whisper{
		Mode:    internal.WhisperMode_WHISPER_MODE_RECEIVE,
		Name:    sender,
		Channel: channel,
		Message: message,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMessengerJoined(ch *game.Character, position uint8, members []*internal.MessengerMember, queued []game.MessengerEvent) {
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

func (l *CharacterListenerImpl) OnMessengerInviteResult(ch *game.Character, name string, success bool) {
	ch.Send(&response.Messenger{
		Action:  internal.MessengerAction_MESSENGER_ACTION_INVITE_RESULT,
		Name:    name,
		Success: success,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMessengerEvent(ch *game.Character, event game.MessengerEvent) {
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

func (l *CharacterListenerImpl) OnPopulationFailed(ch *game.Character, result pconst.PopulationResult) {
	ch.Send(&response.Population{
		Status: result,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPopulationGiven(ch *game.Character, mode pconst.PopulationMode, targetName string, value int32) {
	ch.Send(&response.Population{
		Status: pconst.PopulationResultGiven,
		Mode:   mode,
		Name:   targetName,
		Value:  value,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPopulationReceived(ch *game.Character, mode pconst.PopulationMode, giverName string) {
	ch.Send(&response.Population{
		Status: pconst.PopulationResultReceived,
		Mode:   mode,
		Name:   giverName,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnNotes(ch *game.Character, notes []*internal.Note) {
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

func (l *CharacterListenerImpl) OnNoteSent(ch *game.Character) {
	ch.Send(&response.Notes{
		Action: pconst.NotesPacketSent,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnNoteFailed(ch *game.Character, result internal.NoteResult) {
	ch.Send(&response.Notes{
		Action: pconst.NotesPacketFailed,
		Status: result,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnNoteArrived(ch *game.Character) {
	ch.Send(&response.Notes{
		Action: pconst.NotesPacketRefresh,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFamilyPrivileges(ch *game.Character) {
	ch.Send(&response.FamilyPrivileges{}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFamilyChanged(ch *game.Character) {
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

func (l *CharacterListenerImpl) OnFamilyChart(ch *game.Character, tree game.FamilyTree, characterID uint32) {
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

func (l *CharacterListenerImpl) OnFamilyResult(ch *game.Character, result internal.FamilyResult) {
	ch.Send(&response.FamilyResult{
		Code: result,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFamilyInvite(ch *game.Character, inviterID uint32, name string) {
	ch.Send(&response.FamilyInvite{
		InviterID: inviterID,
		Name:      name,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFamilyInviteResult(ch *game.Character, accepted bool, name string) {
	ch.Send(&response.FamilyInviteResult{
		Accepted: accepted,
		Name:     name,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFamilyJoined(ch *game.Character, name string) {
	ch.Send(&response.FamilyJoined{
		Name: name,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFamilySummon(ch *game.Character, name string, mapName string) {
	ch.Send(&response.FamilySummon{
		Name:    name,
		MapName: mapName,
	}, types.SEND_POLICY_ENCRYPT)
}
