package request

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type FamilySummon struct {
	Family
}

func (*FamilySummon) Opcode() byte { return 0x82 }

func (p *FamilySummon) Serialize(w *stream.StreamWriter) error {
	w.WriteStr16(p.Name)
	w.WriteBoolean(p.Accepted)
	return nil
}

func (p *FamilySummon) Deserialize(r *stream.StreamReader) {
	p.Action = internal.FamilyAction_FAMILY_ACTION_SUMMON
	p.Name = r.ReadStr16()
	p.Accepted = r.ReadBool()
}
