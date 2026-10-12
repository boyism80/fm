package request

import "github.com/boyism80/fm/stream"

type FamilyOpen struct{}

func (*FamilyOpen) Opcode() byte { return 0x7b }

func (p *FamilyOpen) Serialize(_ *stream.StreamWriter) error {
	return nil
}

func (p *FamilyOpen) Deserialize(_ *stream.StreamReader) {
}
