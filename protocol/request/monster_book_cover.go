package request

import (
	"github.com/boyism80/fm/stream"
)

type MonsterBookCover struct {
	CardID uint32
}

func (*MonsterBookCover) Opcode() byte { return 0x28 }

func (p *MonsterBookCover) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CardID)
	return nil
}

func (p *MonsterBookCover) Deserialize(reader *stream.StreamReader) {
	p.CardID = reader.ReadU32()
}
