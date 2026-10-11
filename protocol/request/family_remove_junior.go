package request

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type FamilyRemoveJunior struct {
	Family
}

func (*FamilyRemoveJunior) Opcode() byte { return 0x7d }

func (p *FamilyRemoveJunior) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.TargetID)
	return nil
}

func (p *FamilyRemoveJunior) Deserialize(r *stream.StreamReader) {
	p.Action = internal.FamilyAction_FAMILY_ACTION_REMOVE_JUNIOR
	p.TargetID = r.ReadU32()
}
