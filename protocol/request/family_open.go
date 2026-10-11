package request

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type FamilyOpen struct {
	Family
}

func (*FamilyOpen) Opcode() byte { return 0x7b }

func (p *FamilyOpen) Serialize(_ *stream.StreamWriter) error {
	return nil
}

func (p *FamilyOpen) Deserialize(_ *stream.StreamReader) {
	p.Action = internal.FamilyAction_FAMILY_ACTION_OPEN
}
