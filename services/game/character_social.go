package game

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type Whisper struct {
	RequestID  string
	Mode       internal.WhisperMode
	SenderName string
	SenderRole constant.CharacterRole
	Channel    int32
	Message    string
}

func (ch *Character) SocialLogin(notes []Note) {
	ch.listener.OnFamilyPrivileges(ch)
	ch.listener.OnFamilyChanged(ch)

	ch.notes.entries = notes
	if len(notes) > 0 {
		ch.listener.OnNotes(ch, notes)
	}
}

func (ch *Character) NotifyMemberFieldsChanged() {
	ch.gameWorld.Systems().Get[*PartySystem]().UpdateMemberAsync(nil, ch)
	ch.gameWorld.Systems().Get[*AllianceSystem]().NotifyMemberFieldsChanged(ch)
}

// ReceiveWhisper claims the request first so only one channel shows it,
// then reports whether the message actually reached the client. The claim
// reply returns to this actor even if the character has moved maps since.
func (ch *Character) ReceiveWhisper(ctx actor.Context, whisper Whisper) {
	if ch.IsHidden() && whisper.SenderRole < ch.role {
		return
	}

	ch.gameWorld.Systems().Get[*InternalSystem]().ClaimWhisperAsync(ctx, ch, whisper.RequestID).Do(func(reply *internal.CompleteWhisperReply) error {
		if !reply.Ok {
			return nil
		}

		ch.gameWorld.Systems().Get[*DispatchSystem]().Visit(ch.GetID(), func(ctx actor.Context, current *Character) {
			if current == nil {
				return
			}

			received := !current.LoggedOut()
			if whisper.Mode == internal.WhisperMode_WHISPER_MODE_SEND && received {
				received = current.listener.OnWhisperReceived(current, whisper.SenderName, whisper.Channel, whisper.Message) == nil
			}
			current.gameWorld.Systems().Get[*InternalSystem]().CompleteWhisperAsync(ctx, current, whisper.RequestID, received).OnError(func(error) {})
		})
		return nil
	}).OnError(func(error) {})
}
