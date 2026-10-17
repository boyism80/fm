package request

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type Buddy struct {
	Mode        constant.BuddyMode
	Name        string
	Group       string
	CharacterID uint32
}

func (*Buddy) Opcode() byte { return 0x6C }

func (p *Buddy) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Mode))
	switch p.Mode {
	case constant.BuddyAdd:
		writer.WriteStr16(p.Name)
		writer.WriteStr16(p.Group)
	case constant.BuddyAccept, constant.BuddyDelete:
		writer.WriteU32(p.CharacterID)
	}
	return nil
}

func (p *Buddy) Deserialize(reader *stream.StreamReader) {
	p.Mode = constant.BuddyMode(reader.ReadU8())
	switch p.Mode {
	case constant.BuddyAdd:
		p.Name = reader.ReadStr16()
		p.Group = reader.ReadStr16()
	case constant.BuddyAccept, constant.BuddyDelete:
		p.CharacterID = reader.ReadU32()
	}
}
