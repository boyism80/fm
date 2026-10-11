package request

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type FamilyAccept struct {
	Family
}

func (*FamilyAccept) Opcode() byte { return 0x7f }

func (p *FamilyAccept) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.TargetID)
	w.WriteStr16(p.Name)
	w.WriteBoolean(p.Accepted)
	return nil
}

func (p *FamilyAccept) Deserialize(r *stream.StreamReader) {
	p.Action = internal.FamilyAction_FAMILY_ACTION_ACCEPT
	p.TargetID = r.ReadU32()
	p.Name = r.ReadStr16()
	p.Accepted = r.ReadBool()
}
