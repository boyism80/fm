package request

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type PartyOperation struct {
	Operation         constant.PartyOperationCode
	PartyID           uint32
	TargetName        string
	TargetCharacterID uint32
}

func (*PartyOperation) Opcode() byte { return 0x66 }

func (p *PartyOperation) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *PartyOperation) Deserialize(reader *stream.StreamReader) {
	p.Operation = constant.PartyOperationCode(reader.ReadU8())
	switch p.Operation {
	case constant.PartyC2SAcceptInvite:
		p.PartyID = reader.ReadU32()
	case constant.PartyC2SInvite:
		p.TargetName = reader.ReadStr16()
	case constant.PartyC2SExpel:
		p.TargetCharacterID = reader.ReadU32()
	case constant.PartyC2SChangeLeader:
		p.TargetCharacterID = reader.ReadU32()
	}
}
