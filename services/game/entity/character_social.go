package entity

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
)

type noteRemoveDetail struct {
	Action  internal.NotesAction `json:"action"`
	NoteIDs []uint32             `json:"note_ids,omitempty"`
}

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
	ch.Listener.OnFamilyChanged(ch)

	ch.Notes = notes
	if len(notes) > 0 {
		ch.Listener.OnNotes(ch, notes)
	}
}

// ReceiveWhisper claims the request first so only one channel shows it,
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

			received := !current.LoggedOut()
			if whisper.Mode == internal.WhisperMode_WHISPER_MODE_SEND && received {
				received = current.Listener.OnWhisperReceived(current, whisper.SenderName, whisper.Channel, whisper.Message) == nil
			}
			current.Listener.CompleteWhisperAsync(ctx, current, whisper.RequestID, received).OnError(func(error) {})
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

func (ch *Character) RemoveNotes(actx actor.Context, action internal.NotesAction, noteIDs []uint32) *async.Promise[*internal.NotesReply] {
	detail := noteRemoveDetail{Action: action, NoteIDs: noteIDs}
	id := ch.Unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_NOTE_REMOVE, 0, detail)
	return ch.removeNotes(actx, id, detail)
}

// The reply returns to this actor even if the character has moved maps since, and a
// relogin may have loaded the same request again, so only the holder of the request grants.
func (ch *Character) removeNotes(actx actor.Context, id uint64, detail noteRemoveDetail) *async.Promise[*internal.NotesReply] {
	return ch.Listener.RemoveNotesAsync(actx, ch, id, detail.Action, detail.NoteIDs).Do(func(reply *internal.NotesReply) error {
		ch.GameWorld.GetDispatchSystem().CallCharacter(ch.GetID(), func(_ actor.Context, current *Character) {
			if current == nil || current.Unsettled.Find(id) == nil {
				return
			}
			current.ReceiveNotes(reply.GetNotes(), reply.GetPopulation())
			current.Unsettled.Erase(id)
		})
		return nil
	}).OnError(func(err error) {
		ch.GameWorld.GetDispatchSystem().CallCharacter(ch.GetID(), func(_ actor.Context, current *Character) {
			if current == nil {
				return
			}
			if request := current.Unsettled.Find(id); request != nil {
				request.InFlight = false
			}
		})
		log.Printf("RemoveNotes character=%d id=%d: %v", ch.GetID(), id, err)
	})
}

func (ch *Character) settleNoteRemove(actx actor.Context, request *UnsettledRequest) {
	var detail noteRemoveDetail
	if err := json.Unmarshal([]byte(request.Detail), &detail); err != nil {
		log.Printf("settleNoteRemove character=%d id=%d: %v", ch.GetID(), request.ID, err)
		return
	}
	ch.removeNotes(actx, request.ID, detail)
}
