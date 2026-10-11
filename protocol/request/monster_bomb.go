package request

import "github.com/boyism80/fm/stream"

type MonsterBomb struct {
	OID uint32
}

func (*MonsterBomb) Opcode() byte { return 0x9A }

func (p *MonsterBomb) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.OID)
	return nil
}

func (p *MonsterBomb) Deserialize(reader *stream.StreamReader) {
	p.OID = reader.ReadU32()
}
