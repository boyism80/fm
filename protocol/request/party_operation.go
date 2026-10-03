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
	writer.WriteU8(uint8(p.Operation))
	switch p.Operation {
	case constant.PartyC2SAcceptInvite:
		writer.WriteU32(p.PartyID)
	case constant.PartyC2SInvite:
		writer.WriteStr16(p.TargetName)
	case constant.PartyC2SExpel:
		writer.WriteU32(p.TargetCharacterID)
	case constant.PartyC2SChangeLeader:
		writer.WriteU32(p.TargetCharacterID)
	}
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
