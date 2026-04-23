package request

import (
	"github.com/boyism80/fm/stream"
)

type Aggro struct{}

func (*Aggro) Opcode() byte { return 0x96 }

func (p *Aggro) Serialize(writer *stream.StreamWriter) error {

	return nil
}

func (p *Aggro) Deserialize(reader *stream.StreamReader) {

}
