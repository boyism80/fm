package request

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type FamilyPrecept struct {
	Family
}

func (*FamilyPrecept) Opcode() byte { return 0x81 }

func (p *FamilyPrecept) Serialize(w *stream.StreamWriter) error {
	w.WriteStr16(p.Name)
	return nil
}

func (p *FamilyPrecept) Deserialize(r *stream.StreamReader) {
	p.Action = internal.FamilyAction_FAMILY_ACTION_PRECEPT
	p.Name = r.ReadStr16()
}
