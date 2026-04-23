package request

import (
	"github.com/boyism80/fm/stream"
)

type Secure struct {
	Type    byte
	ParamA  int32
	ParamB  int32
	Payload []byte
}

func (*Secure) Opcode() byte { return 0x0C }

func (p *Secure) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *Secure) Deserialize(reader *stream.StreamReader) {
	n := reader.Remaining()
	if n < 9 {
		if n > 0 {
			p.Payload = reader.Read(n)
		}
		return
	}
	p.Type = reader.ReadU8()
	p.ParamA = reader.Read32()
	p.ParamB = reader.Read32()
	remain := reader.Remaining()
	if remain > 0 {
		p.Payload = reader.Read(remain)
	}
}
