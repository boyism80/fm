package server

import (
	"github.com/boyism80/fm/core/types"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/protocol"
	"github.com/boyism80/fm/game/protocol/resp"
)

// CharacterListenerImpl implements CharacterListener for game server
type CharacterListenerImpl struct {
	gameServer *GameServer
	character  *entity.Character
}

// NewGameCharacterListener creates a new GameCharacterListener instance
func NewGameCharacterListener(gameServer *GameServer, character *entity.Character) *CharacterListenerImpl {
	return &CharacterListenerImpl{
		gameServer: gameServer,
		character:  character,
	}
}

// OnDialog handles dialog events
func (l *CharacterListenerImpl) OnDialog(npc uint32, message string, prev bool, next bool) {
	// Send dialog packet to client
	dialogPacket := &resp.Dialog{
		NPC:  npc,
		Type: constant.DIALOG_TYPE_DEFAULT,
		Text: message,
		Prev: prev,
		Next: next,
	}
	l.character.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

// OnDialogYesNo handles yes/no dialog events
func (l *CharacterListenerImpl) OnDialogYesNo(npc uint32, message string, prev bool, next bool) {
	// Send yes/no dialog packet to client
	dialogPacket := &resp.DialogYesNo{
		NPC:  npc,
		Text: message,
		Prev: prev,
		Next: next,
	}
	l.character.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

// OnDialogAccept handles dialog accept events
func (l *CharacterListenerImpl) OnDialogAccept(npc uint32, message string, enableEscape bool) {
	// Send dialog accept packet to client
	dialogPacket := &resp.DialogAccept{
		NPC:          npc,
		Text:         message,
		EnableEscape: enableEscape,
	}
	l.character.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

// OnDialogList handles dialog list events
func (l *CharacterListenerImpl) OnDialogList(npc uint32, message string, selections []string) {
	// Send dialog list packet to client
	dialogPacket := &resp.DialogList{
		NPC:        npc,
		Text:       message,
		Selections: selections,
	}
	l.character.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

// OnDialogInput handles dialog input events
func (l *CharacterListenerImpl) OnDialogInput(npc uint32, message string) {
	// Send dialog input packet to client
	dialogPacket := &resp.DialogInput{
		NPC:  npc,
		Text: message,
	}
	l.character.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

// OnChat handles chat events by broadcasting chat packet to all players on the map
func (l *CharacterListenerImpl) OnChat(message string, highlight bool, dontRecordHistory bool) {
	// Create chat packet
	chatPacket := &resp.NormalChat{
		CharacterId:       l.character.ID,
		Message:           message,
		Highlight:         highlight,
		DontRecordHistory: dontRecordHistory,
	}

	// Get map instance
	mapInstance := l.gameServer.GetMap(l.character.Map)
	if mapInstance == nil {
		return
	}

	mapInstance.BroadcastToAllPlayers(chatPacket, types.SEND_POLICY_ENCRYPT)
}

// OnAttack handles attack events
func (l *CharacterListenerImpl) OnAttack(attackInfo protocol.AttackInfo) {
	// TODO: Implement attack handling
	// This could involve sending attack packets to other players
}

// OnMesoChanged handles meso change events
func (l *CharacterListenerImpl) OnMesoChanged(meso int32) {
	// TODO: Implement meso change handling
	// This could involve sending a packet to update the client's meso display
	l.character.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_MESO: meso,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

// OnMessage handles general message events
func (l *CharacterListenerImpl) OnMessage(messageType constant.ServerMessageType, message string) {
	// Send notice packet to client
	noticePacket := &resp.Notice{
		Message: message,
		Type:    messageType,
	}
	l.character.Send(noticePacket, types.SEND_POLICY_ENCRYPT)
}

// OnExpGain handles experience gain events
func (l *CharacterListenerImpl) OnExpGain(exp uint32) {
	// Send experience gain packet to client
	expPacket := &resp.GainExp{
		Gain:  exp,
		White: false,
	}
	l.character.Send(expPacket, types.SEND_POLICY_ENCRYPT)
}
