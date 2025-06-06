package actorx

import (
	"github.com/asynkron/protoactor-go/actor"
	common_msg "github.com/boyism80/fm/common/msg"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/msg"
	"github.com/boyism80/fm/game/protocol/resp"
)

type CharacterListener struct {
	Actor *GameClientActor
}

func (l *CharacterListener) GetContext() actor.Context {
	return l.Actor.actorContext
}

func (l *CharacterListener) OnDialog(message string, prev bool, next bool) {
	l.Actor.Send(&resp.Dialog{
		NPC:  9001000,
		Text: message,
		Prev: prev,
		Next: next,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListener) OnDialogYesNo(message string, prev bool, next bool) {

	l.Actor.Send(&resp.DialogYesNo{
		NPC:  9001000,
		Text: message,
		Prev: prev,
		Next: next,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListener) OnDialogAccept(message string, enableEscape bool) {
	l.Actor.Send(&resp.DialogAccept{
		NPC:          9001000,
		Text:         message,
		EnableEscape: enableEscape,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListener) OnDialogList(message string, selections []string) {
	l.Actor.Send(&resp.DialogList{
		NPC:        9001000,
		Text:       message,
		Selections: selections,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListener) OnDialogInput(message string) {
	l.Actor.Send(&resp.DialogInput{
		NPC:  9001000,
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

	ctx.Send(mapActor, &msg.MapBroadcastRange{
		Sender: ctx.Self(),
		Pivot:  position,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.NormalChat{
				CharacterId:       ch.ID,
				Highlight:         false,
				Message:           message,
				DontRecordHistory: dontRecordHistory,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		},
	})
}
