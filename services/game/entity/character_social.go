package entity

import (
	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
)

type Whisper struct {
	RequestID  string
	Mode       internal.WhisperMode
	SenderName string
	SenderRole constant.CharacterRole
	Channel    int32
	Message    string
}

func (ch *Character) SocialLogin(notes []*internal.Note) {
	ch.Listener.OnFamilyPrivileges(ch)
	if ch.Family.State != nil {
		ch.Listener.OnFamilyChanged(ch)
	}

	ch.Notes = notes
	if len(notes) > 0 {
		ch.Listener.OnNotes(ch, notes)
	}
}

// ReceiveWhisper claims the request first so only one channel delivers it,
// then reports whether the message actually reached the client. The claim
// reply returns to this actor even if the character has moved maps since.
func (ch *Character) ReceiveWhisper(ctx actor.Context, whisper Whisper) {
	if ch.IsHidden() && whisper.SenderRole < ch.Role {
		return
	}

	ch.Listener.ClaimWhisperAsync(ctx, ch, whisper.RequestID).Do(func(reply *internal.CompleteWhisperReply) error {
		if !reply.Ok {
			return nil
		}

		ch.GameWorld.GetDispatchSystem().CallCharacter(ch.GetID(), func(ctx actor.Context, current *Character) {
			if current == nil {
				return
			}

			delivered := !current.LoggedOut()
			if whisper.Mode == internal.WhisperMode_WHISPER_MODE_SEND && delivered {
				delivered = current.Listener.OnWhisperReceived(current, whisper.SenderName, whisper.Channel, whisper.Message) == nil
			}
			current.Listener.CompleteWhisperAsync(ctx, current, whisper.RequestID, delivered).OnError(func(error) {})
		})
		return nil
	}).OnError(func(error) {})
}

// Notes mirror internal's social_notes, which stays the source of truth.
func (ch *Character) ReceiveNotes(notes []*internal.Note, population uint32) {
	ch.Notes = notes
	if population > 0 {
		ch.Stats.AddPopulation(int32(population))
	}
}
