package response

import "github.com/boyism80/fm/stream"

// CoconutHit changes a coconut's state after Delay ms; a negative ID changes every coconut.
type CoconutHit struct {
	ID    int16
	Delay int16
	State uint8
}

func (p *CoconutHit) Opcode() uint16 {
	return 0xD7
}

func (p *CoconutHit) Serialize(sw *stream.StreamWriter) error {
	sw.Write16(p.ID)
	sw.Write16(p.Delay)
	sw.WriteU8(p.State)
	return nil
}

func (p *CoconutHit) Deserialize(reader *stream.StreamReader) {
	p.ID = reader.Read16()
	p.Delay = reader.Read16()
	p.State = reader.ReadU8()
}

// CoconutScore shows the coconuts each team has knocked down.
type CoconutScore struct {
	Maple int16
	Story int16
}

func (p *CoconutScore) Opcode() uint16 {
	return 0xD8
}

func (p *CoconutScore) Serialize(sw *stream.StreamWriter) error {
	sw.Write16(p.Maple)
	sw.Write16(p.Story)
	return nil
}

func (p *CoconutScore) Deserialize(reader *stream.StreamReader) {
	p.Maple = reader.Read16()
	p.Story = reader.Read16()
}
