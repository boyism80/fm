package request

import "github.com/boyism80/fm/stream"

type DenyGuildRequest struct {
	Code        uint8
	InviterName string
}

func (*DenyGuildRequest) Opcode() byte {
	return 0x69
}

func (p *DenyGuildRequest) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *DenyGuildRequest) Deserialize(reader *stream.StreamReader) {
	p.Code = reader.ReadU8()
	p.InviterName = reader.ReadStr16()
}
