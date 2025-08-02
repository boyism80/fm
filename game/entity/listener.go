package entity

import (
	"github.com/boyism80/fm/game/protocol"
)

type CharacterListener interface {
	OnDialog(npc uint32, message string, prev bool, next bool)
	OnDialogYesNo(npc uint32, message string, prev bool, next bool)
	OnDialogAccept(npc uint32, message string, enableEscape bool)
	OnDialogList(npc uint32, message string, selections []string)
	OnDialogInput(npc uint32, message string)
	OnChat(message string, highlight bool, dontRecordHistory bool)
	OnAttack(attackInfo protocol.AttackInfo)
	OnMesoChanged(meso int32)
	OnMessage(message string)
}
