package request

import "github.com/boyism80/fm/stream"

type MonsterCarnival struct {
	Tab uint8
	Num int32
}

func (*MonsterCarnival) Opcode() byte { return 0xB1 }

func (p *MonsterCarnival) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *MonsterCarnival) Deserialize(reader *stream.StreamReader) {
	p.Tab = reader.ReadU8()
	p.Num = reader.Read32()
}
