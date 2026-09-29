package response

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type CarnivalStart struct {
	Team                constant.CarnivalTeam
	PersonalAvailableCP uint16
	PersonalTotalCP     uint16
	TeamAvailableCP     uint16
	TeamTotalCP         uint16
	EnemyAvailableCP    uint16
	EnemyTotalCP        uint16
}

func (p *CarnivalStart) Opcode() uint16 { return 0xDB }

func (p *CarnivalStart) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Team))
	writer.WriteU16(p.PersonalAvailableCP)
	writer.WriteU16(p.PersonalTotalCP)
	writer.WriteU16(p.TeamAvailableCP)
	writer.WriteU16(p.TeamTotalCP)
	writer.WriteU16(p.EnemyAvailableCP)
	writer.WriteU16(p.EnemyTotalCP)
	writer.WriteU64(0)
	writer.WriteU16(0)
	return nil
}

func (p *CarnivalStart) Deserialize(reader *stream.StreamReader) {}

type CarnivalObtainedCP struct {
	Avail uint16
	Total uint16
}

func (p *CarnivalObtainedCP) Opcode() uint16 { return 0xDC }

func (p *CarnivalObtainedCP) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(p.Avail)
	writer.WriteU16(p.Total)
	return nil
}

func (p *CarnivalObtainedCP) Deserialize(reader *stream.StreamReader) {}

type CarnivalPartyCP struct {
	Team  constant.CarnivalTeam
	Avail uint16
	Total uint16
}

func (p *CarnivalPartyCP) Opcode() uint16 { return 0xDD }

func (p *CarnivalPartyCP) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Team))
	writer.WriteU16(p.Avail)
	writer.WriteU16(p.Total)
	return nil
}

func (p *CarnivalPartyCP) Deserialize(reader *stream.StreamReader) {}

type CarnivalSummon struct {
	Tab  uint8
	Num  uint8
	Name string
}

func (p *CarnivalSummon) Opcode() uint16 { return 0xDE }

func (p *CarnivalSummon) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(p.Tab)
	writer.WriteU8(p.Num)
	writer.WriteStr16(p.Name)
	return nil
}

func (p *CarnivalSummon) Deserialize(reader *stream.StreamReader) {}

type CarnivalDied struct {
	Team   constant.CarnivalTeam
	Name   string
	LostCP uint8
}

func (p *CarnivalDied) Opcode() uint16 { return 0xDF }

func (p *CarnivalDied) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Team))
	writer.WriteStr16(p.Name)
	writer.WriteU8(p.LostCP)
	return nil
}

func (p *CarnivalDied) Deserialize(reader *stream.StreamReader) {}
