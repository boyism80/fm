package request

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type NoteRemoval struct {
	ID     uint32
	Kind   pconst.NoteKind
	ItemID uint32
}

type Notes struct {
	Action   internal.NotesAction
	Target   string
	Message  string
	Kind     pconst.NoteKind
	Removals []NoteRemoval
}

func (*Notes) Opcode() byte { return 0x6D }

func (p *Notes) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(p.Action))
	switch p.Action {
	case internal.NotesAction_NOTES_ACTION_SEND:
		w.WriteStr16(p.Target)
		w.WriteStr16(p.Message)
		w.WriteU8(uint8(p.Kind))
	case internal.NotesAction_NOTES_ACTION_REMOVE:
		w.WriteU8(uint8(len(p.Removals)))
		w.WriteU16(0)
		for _, n := range p.Removals {
			w.WriteU32(n.ID)
			w.WriteU8(uint8(n.Kind))
			if n.Kind == pconst.NoteKindItem {
				w.WriteU32(n.ItemID)
			}
		}
	}
	return nil
}

func (p *Notes) Deserialize(r *stream.StreamReader) {
	p.Action = internal.NotesAction(r.ReadU8())
	switch p.Action {
	case internal.NotesAction_NOTES_ACTION_SEND:
		p.Target = r.ReadStr16()
		p.Message = r.ReadStr16()
		p.Kind = pconst.NoteKind(r.ReadU8())
	case internal.NotesAction_NOTES_ACTION_REMOVE:
		count := int(r.ReadU8())
		r.Skip(2)
		for i := 0; i < count; i++ {
			n := NoteRemoval{ID: r.ReadU32(), Kind: pconst.NoteKind(r.ReadU8())}
			if n.Kind == pconst.NoteKindItem {
				n.ItemID = r.ReadU32()
			}
			p.Removals = append(p.Removals, n)
		}
	}
}
