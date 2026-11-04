package server

import (
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) OnDialog(ch *game.Character, npc uint32, message string, prev bool, next bool) {
	dialogPacket := &response.Dialog{
		NPC:  npc,
		Type: constant.DialogTypeDefault,
		Text: message,
		Prev: prev,
		Next: next,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogYesNo(ch *game.Character, npc uint32, message string, prev bool, next bool) {
	dialogPacket := &response.DialogYesNo{
		NPC:  npc,
		Text: message,
		Prev: prev,
		Next: next,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogAccept(ch *game.Character, npc uint32, message string, enableEscape bool) {
	dialogPacket := &response.DialogAccept{
		NPC:          npc,
		Text:         message,
		EnableEscape: enableEscape,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogList(ch *game.Character, npc uint32, message string, selections []string) {
	dialogPacket := &response.DialogList{
		NPC:        npc,
		Text:       message,
		Selections: selections,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogStyle(ch *game.Character, npc uint32, message string, styles []uint32) {
	ch.Send(&response.DialogStyle{
		NPC:    npc,
		Text:   message,
		Styles: styles,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogInput(ch *game.Character, npc uint32, message string) {
	dialogPacket := &response.DialogInput{
		NPC:  npc,
		Text: message,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogNumber(ch *game.Character, npc uint32, message string, def, min, max int32) {
	ch.Send(&response.DialogNumber{
		NPC:     npc,
		Text:    message,
		Default: def,
		Min:     min,
		Max:     max,
	}, types.SEND_POLICY_ENCRYPT)
}
