package request

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type FamilyRemoveSenior struct {
	Family
}

func (*FamilyRemoveSenior) Opcode() byte { return 0x7e }

func (p *FamilyRemoveSenior) Serialize(_ *stream.StreamWriter) error {
	return nil
}

func (p *FamilyRemoveSenior) Deserialize(_ *stream.StreamReader) {
	p.Action = internal.FamilyAction_FAMILY_ACTION_REMOVE_SENIOR
}
