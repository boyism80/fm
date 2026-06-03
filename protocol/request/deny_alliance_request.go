package request

import "github.com/boyism80/fm/stream"

type DenyAllianceRequest struct{}

func (*DenyAllianceRequest) Opcode() byte {
	return 0x79
}

func (p *DenyAllianceRequest) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *DenyAllianceRequest) Deserialize(reader *stream.StreamReader) {}
