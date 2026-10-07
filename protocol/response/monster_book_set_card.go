package response

import (
	"github.com/boyism80/fm/stream"
)

type MonsterBookSetCard struct {
	Success bool
	CardID  uint32
	Count   uint32
}

func (p *MonsterBookSetCard) Opcode() uint16 {
	return 0x41
}

func (p *MonsterBookSetCard) Serialize(writer *stream.StreamWriter) error {
	writer.WriteBoolean(p.Success)
	if p.Success {
		writer.WriteU32(p.CardID)
		writer.WriteU32(p.Count)
	}
	return nil
}

func (p *MonsterBookSetCard) Deserialize(reader *stream.StreamReader) {
	p.Success = reader.ReadBool()
	if p.Success {
		p.CardID = reader.ReadU32()
		p.Count = reader.ReadU32()
	}
}
