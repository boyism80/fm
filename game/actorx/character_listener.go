package actorx

import (
	common_msg "github.com/boyism80/fm/common/msg"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/msg"
	"github.com/boyism80/fm/game/protocol/resp"
)

type CharacterListener struct {
	Actor *GameClientActor
}

func (l *CharacterListener) OnDialog(npc uint32, message string, prev bool, next bool) {
	l.Actor.Send(&resp.Dialog{
		NPC:  npc,
		Text: message,
		Prev: prev,
		Next: next,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListener) OnDialogYesNo(npc uint32, message string, prev bool, next bool) {

	l.Actor.Send(&resp.DialogYesNo{
		NPC:  npc,
		Text: message,
		Prev: prev,
		Next: next,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListener) OnDialogAccept(npc uint32, message string, enableEscape bool) {
	l.Actor.Send(&resp.DialogAccept{
		NPC:          npc,
		Text:         message,
		EnableEscape: enableEscape,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListener) OnDialogList(npc uint32, message string, selections []string) {
	l.Actor.Send(&resp.DialogList{
		NPC:        npc,
		Text:       message,
		Selections: selections,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListener) OnDialogInput(npc uint32, message string) {
	l.Actor.Send(&resp.DialogInput{
		NPC:  npc,
		Text: message,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListener) OnChat(message string, highlight bool, dontRecordHistory bool) {
	client := l.Actor
	ch := client.ch
	ctx := l.Actor.actorContext
	position := ch.Position
	mapActor := client.serverContext.MapActors[ch.Map]
	if mapActor == nil {
		return
	}

	ctx.Send(mapActor, &msg.MapBroadcast{
		Sender: ctx.Self(),
		Pivot:  position,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.NormalChat{
				CharacterId:       ch.ID,
				Highlight:         highlight,
				Message:           message,
				DontRecordHistory: dontRecordHistory,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		},
	})
}

func (l *CharacterListener) OnMesoChanged(meso int32) {
	l.Actor.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_MESO: meso,
		},
	}, types.SEND_POLICY_ENCRYPT)
}
