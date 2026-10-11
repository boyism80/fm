package request

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type FamilyInvite struct {
	Family
}

func (*FamilyInvite) Opcode() byte { return 0x7c }

func (p *FamilyInvite) Serialize(w *stream.StreamWriter) error {
	w.WriteStr16(p.Name)
	return nil
}

func (p *FamilyInvite) Deserialize(r *stream.StreamReader) {
	p.Action = internal.FamilyAction_FAMILY_ACTION_INVITE
	p.Name = r.ReadStr16()
}
