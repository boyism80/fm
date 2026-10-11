package response

import (
	"time"

	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/util"
)

type Note struct {
	ID      uint32
	Sender  string
	Message string
	SentAt  time.Time
	Kind    pconst.NoteKind
}

type Notes struct {
	Action pconst.NotesPacket
	Status internal.NoteResult
	Notes  []Note
}

func (p *Notes) Opcode() uint16 {
	return 0x1D
}

func (p *Notes) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(p.Action))
	switch p.Action {
	case pconst.NotesPacketList:
		w.WriteU8(uint8(len(p.Notes)))
		for _, n := range p.Notes {
			w.WriteU32(n.ID)
			w.WriteStr16(n.Sender)
			w.WriteStr16(n.Message)
			w.WriteDateTime(n.SentAt)
			w.WriteU8(uint8(n.Kind))
		}
	case pconst.NotesPacketFailed:
		w.WriteU8(uint8(p.Status))
	}
	return nil
}

func (p *Notes) Deserialize(reader *stream.StreamReader) {
	p.Action = pconst.NotesPacket(reader.ReadU8())
	switch p.Action {
	case pconst.NotesPacketList:
		count := int(reader.ReadU8())
		for range count {
			p.Notes = append(p.Notes, Note{
				ID:      reader.ReadU32(),
				Sender:  reader.ReadStr16(),
				Message: reader.ReadStr16(),
				SentAt:  util.FromFileTime(reader.ReadU64()),
				Kind:    pconst.NoteKind(reader.ReadU8()),
			})
		}
	case pconst.NotesPacketFailed:
		p.Status = internal.NoteResult(reader.ReadU8())
	}
}
