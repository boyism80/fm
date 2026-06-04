package request

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type AllianceOperation struct {
	Operation             constant.AllianceC2SOperation
	AllianceName          string
	TargetGuildLeaderName string
	TargetGuildID         uint32
	AllianceID            uint32
	NewLeaderID           uint32
	RankTitles            [5]string
	TargetCharacterID     uint32
	RankChangePromote     bool
	Notice                string
}

func (*AllianceOperation) Opcode() byte {
	return 0x78
}

func (p *AllianceOperation) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *AllianceOperation) Deserialize(reader *stream.StreamReader) {
	p.Operation = constant.AllianceC2SOperation(reader.ReadU8())
	switch p.Operation {
	case constant.AllianceC2SCreate:
		p.AllianceName = reader.ReadStr16()

	case constant.AllianceC2SInvite:
		p.TargetGuildLeaderName = reader.ReadStr16()

	case constant.AllianceC2SExpel:
		p.TargetGuildID = reader.ReadU32()
		if reader.Remaining() >= 4 {
			p.AllianceID = reader.ReadU32()
		}

	case constant.AllianceC2SChangeLeader:
		p.NewLeaderID = reader.ReadU32()

	case constant.AllianceC2SChangeRankTitles:
		for i := range 5 {
			p.RankTitles[i] = reader.ReadStr16()
		}

	case constant.AllianceC2SChangeMemberRank:
		p.TargetCharacterID = reader.ReadU32()
		p.RankChangePromote = reader.ReadBool()

	case constant.AllianceC2SChangeNotice:
		p.Notice = reader.ReadStr16()

	case constant.AllianceC2SLoadInfo, constant.AllianceC2SLeave, constant.AllianceC2SAcceptInvite, constant.AllianceC2SDenyInvite:
	default:
	}
}
