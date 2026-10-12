package request

import "github.com/boyism80/fm/stream"

type FamilyRemoveSenior struct{}

func (*FamilyRemoveSenior) Opcode() byte { return 0x7e }

func (p *FamilyRemoveSenior) Serialize(_ *stream.StreamWriter) error {
	return nil
}

func (p *FamilyRemoveSenior) Deserialize(_ *stream.StreamReader) {
}
