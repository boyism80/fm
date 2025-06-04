package listener

import (
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/protocol/resp"
)

type CharacterListener struct {
}

func (l *CharacterListener) OnDialog(ch *entity.Character, msg string, prev bool, next bool) {
	ch.Send(&resp.Dialog{
		NPC:  9001000,
		Text: msg,
		Prev: prev,
		Next: next,
	}, types.SEND_POLICY_ENCRYPT)
}
