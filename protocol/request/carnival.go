package request

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type Carnival struct {
	Tab constant.CarnivalTab
	Num int32
}

func (*Carnival) Opcode() byte { return 0xB1 }

func (p *Carnival) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *Carnival) Deserialize(reader *stream.StreamReader) {
	p.Tab = constant.CarnivalTab(reader.ReadU8())
	p.Num = reader.Read32()
}
