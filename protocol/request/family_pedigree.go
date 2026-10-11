package request

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type FamilyPedigree struct {
	Family
}

func (*FamilyPedigree) Opcode() byte { return 0x7a }

func (p *FamilyPedigree) Serialize(w *stream.StreamWriter) error {
	w.WriteStr16(p.Name)
	return nil
}

func (p *FamilyPedigree) Deserialize(r *stream.StreamReader) {
	p.Action = internal.FamilyAction_FAMILY_ACTION_PEDIGREE
	p.Name = r.ReadStr16()
}
