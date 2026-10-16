package game

import (
	"encoding/json"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

// Note is one social note. Internal stays the source of truth and is copied in.
type Note struct {
	id      uint32
	sender  string
	message string
	sentAt  time.Time
	gift    uint32
}

func (n Note) ID() uint32        { return n.id }
func (n Note) Sender() string    { return n.sender }
func (n Note) Message() string   { return n.message }
func (n Note) SentAt() time.Time { return n.sentAt }
func (n Note) Gift() uint32      { return n.gift }

func NotesFromProto(notes []*internal.Note) []Note {
	out := make([]Note, 0, len(notes))
	for _, pb := range notes {
		if pb == nil {
			continue
		}
		out = append(out, Note{
			id:      pb.GetId(),
			sender:  pb.GetSenderName(),
			message: pb.GetMessage(),
			sentAt:  time.UnixMilli(pb.GetSentAtUnixMs()),
			gift:    pb.GetGift(),
		})
	}
	return out
}

type Notes struct {
	owner   *Character
	entries []Note
}

type noteRemoveDetail struct {
	Action  internal.NotesAction `json:"action"`
	NoteIDs []uint32             `json:"note_ids,omitempty"`
}

func (n *Notes) Receive(notes []*internal.Note, population uint32) {
	n.entries = NotesFromProto(notes)
	if population > 0 {
		n.owner.stats.AddPopulation(int32(population))
	}
}

func (n *Notes) Remove(actx actor.Context, action internal.NotesAction, noteIDs []uint32) *async.Promise[*internal.NotesReply] {
	detail := noteRemoveDetail{Action: action, NoteIDs: noteIDs}
	id := n.owner.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_NOTE_REMOVE, 0, detail)
	return n.remove(actx, id, detail)
}

// The reply returns to this actor even if the character has moved maps since, and a
// relogin may have loaded the same request again, so only the holder of the request grants.
func (n *Notes) remove(actx actor.Context, id uint64, detail noteRemoveDetail) *async.Promise[*internal.NotesReply] {
	return n.owner.gameWorld.Systems().Get[*InternalSystem]().RemoveNotesAsync(actx, n.owner, id, detail.Action, detail.NoteIDs).Do(func(reply *internal.NotesReply) error {
		n.owner.gameWorld.Systems().Get[*DispatchSystem]().Visit(n.owner.GetID(), func(_ actor.Context, current *Character) {
			if current == nil || current.unsettled.Find(id) == nil {
				return
			}
			current.notes.Receive(reply.GetNotes(), reply.GetPopulation())
			current.unsettled.Erase(id)
		})
		return nil
	}).OnError(func(err error) {
		n.owner.gameWorld.Systems().Get[*DispatchSystem]().Visit(n.owner.GetID(), func(_ actor.Context, current *Character) {
			if current == nil {
				return
			}
			if request := current.unsettled.Find(id); request != nil {
				request.inFlight = false
			}
		})
		log.Printf("RemoveNotes character=%d id=%d: %v", n.owner.GetID(), id, err)
	})
}

func (n *Notes) settleRemove(actx actor.Context, request *UnsettledRequest) {
	var detail noteRemoveDetail
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		log.Printf("settleNoteRemove character=%d id=%d: %v", n.owner.GetID(), request.id, err)
		return
	}
	n.remove(actx, request.id, detail)
}
