package request

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildOperation struct {
	Operation     constant.GuildOperationCode
	GuildName     string
	TargetName    string
	GuildID       uint32
	CharacterID   uint32
	CharacterName string
	RankTitles    [5]string
	NewMemberRank uint8
	LogoBG        uint16
	LogoBGColor   uint8
	Logo          uint16
	LogoColor     uint8
	Notice        string
	SkillID       uint32
	NewLeaderID   uint32
}

func (*GuildOperation) Opcode() byte {
	return 0x68
}

func (p *GuildOperation) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *GuildOperation) Deserialize(reader *stream.StreamReader) {
	p.Operation = constant.GuildOperationCode(reader.ReadU8())
	switch p.Operation {
	case constant.GuildC2SCreate:
		p.GuildName = reader.ReadStr16()
	case constant.GuildC2SInvite:
		p.TargetName = reader.ReadStr16()
	case constant.GuildC2SAcceptInvite:
		p.GuildID = reader.ReadU32()
		p.CharacterID = reader.ReadU32()
	case constant.GuildC2SLeave, constant.GuildC2SExpel:
		p.CharacterID = reader.ReadU32()
		p.CharacterName = reader.ReadStr16()
	case constant.GuildC2SChangeRankTitles:
		for i := range 5 {
			p.RankTitles[i] = reader.ReadStr16()
		}
	case constant.GuildC2SChangeMemberRank:
		p.CharacterID = reader.ReadU32()
		p.NewMemberRank = reader.ReadU8()
	case constant.GuildC2SChangeEmblem:
		p.LogoBG = reader.ReadU16()
		p.LogoBGColor = reader.ReadU8()
		p.Logo = reader.ReadU16()
		p.LogoColor = reader.ReadU8()
	case constant.GuildC2SChangeNotice:
		p.Notice = reader.ReadStr16()
	case constant.GuildC2SPurchaseSkill, constant.GuildC2SActivateSkill:
		p.SkillID = reader.ReadU32()
	case constant.GuildC2SChangeLeader:
		p.NewLeaderID = reader.ReadU32()
	}
}
