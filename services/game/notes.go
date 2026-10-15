package game

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

// Entries mirror internal's social_notes, which stays the source of truth.
type Notes struct {
	owner   *Character
	Entries []*internal.Note
}

type noteRemoveDetail struct {
	Action  internal.NotesAction `json:"action"`
	NoteIDs []uint32             `json:"note_ids,omitempty"`
}

func (n *Notes) Receive(notes []*internal.Note, population uint32) {
	n.Entries = notes
	if population > 0 {
		n.owner.Stats.AddPopulation(int32(population))
	}
}

func (n *Notes) Remove(actx actor.Context, action internal.NotesAction, noteIDs []uint32) *async.Promise[*internal.NotesReply] {
	detail := noteRemoveDetail{Action: action, NoteIDs: noteIDs}
	id := n.owner.Unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_NOTE_REMOVE, 0, detail)
	return n.remove(actx, id, detail)
}

// The reply returns to this actor even if the character has moved maps since, and a
// relogin may have loaded the same request again, so only the holder of the request grants.
func (n *Notes) remove(actx actor.Context, id uint64, detail noteRemoveDetail) *async.Promise[*internal.NotesReply] {
	return n.owner.GameWorld.Systems().Get[*InternalSystem]().RemoveNotesAsync(actx, n.owner, id, detail.Action, detail.NoteIDs).Do(func(reply *internal.NotesReply) error {
		n.owner.GameWorld.Systems().Get[*DispatchSystem]().Visit(n.owner.GetID(), func(_ actor.Context, current *Character) {
			if current == nil || current.Unsettled.Find(id) == nil {
				return
			}
			current.Notes.Receive(reply.GetNotes(), reply.GetPopulation())
			current.Unsettled.Erase(id)
		})
		return nil
	}).OnError(func(err error) {
		n.owner.GameWorld.Systems().Get[*DispatchSystem]().Visit(n.owner.GetID(), func(_ actor.Context, current *Character) {
			if current == nil {
				return
			}
			if request := current.Unsettled.Find(id); request != nil {
				request.InFlight = false
			}
		})
		log.Printf("RemoveNotes character=%d id=%d: %v", n.owner.GetID(), id, err)
	})
}

func (n *Notes) settleRemove(actx actor.Context, request *UnsettledRequest) {
	var detail noteRemoveDetail
	if err := json.Unmarshal([]byte(request.Detail), &detail); err != nil {
		log.Printf("settleNoteRemove character=%d id=%d: %v", n.owner.GetID(), request.ID, err)
		return
	}
	n.remove(actx, request.ID, detail)
}
