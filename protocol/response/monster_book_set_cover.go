package response

import (
	"github.com/boyism80/fm/stream"
)

type MonsterBookSetCover struct {
	CardID uint32
}

func (p *MonsterBookSetCover) Opcode() uint16 {
	return 0x42
}

func (p *MonsterBookSetCover) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CardID)
	return nil
}

func (p *MonsterBookSetCover) Deserialize(reader *stream.StreamReader) {
	p.CardID = reader.ReadU32()
}
