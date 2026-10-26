package response

import "github.com/boyism80/fm/stream"

// HealerMove moves the field's healer to Y; the client animates it from where it last was.
type HealerMove struct {
	Y int16
}

func (p *HealerMove) Opcode() uint16 {
	return 0xD9
}

func (p *HealerMove) Serialize(sw *stream.StreamWriter) error {
	sw.Write16(p.Y)
	return nil
}

func (p *HealerMove) Deserialize(reader *stream.StreamReader) {
	p.Y = reader.Read16()
}
