package resp

import (
	"math"

	"github.com/boyism80/fm/core/stream"
)

type Hint struct {
	Text   string
	Width  uint16
	Height uint16
}

func (p *Hint) Opcode() uint16 {
	return 0x9F
}

func (p *Hint) Serialize(sw *stream.StreamWriter) error {
	w := p.Width
	if w == 0 {
		w = uint16(len(p.Text) * 10)
	}
	w = uint16(math.Max(float64(w), 40))
	h := uint16(math.Max(float64(p.Height), 5))

	sw.WriteStr16(p.Text)
	sw.WriteU16(w)
	sw.WriteU16(h)
	sw.WriteU8(1)
	return nil
}

func (p *Hint) Deserialize(sr *stream.StreamReader) error {
	return nil
}
