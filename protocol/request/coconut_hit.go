package request

import "github.com/boyism80/fm/stream"

// CoconutHit reports the coconut a melee attack struck; ID is its index among the map's coconuts.
type CoconutHit struct {
	ID    int16
	Delay int16
}

func (*CoconutHit) Opcode() byte {
	return 0xAC
}

func (p *CoconutHit) Serialize(writer *stream.StreamWriter) error {
	writer.Write16(p.ID)
	writer.Write16(p.Delay)
	return nil
}

func (p *CoconutHit) Deserialize(reader *stream.StreamReader) {
	p.ID = reader.Read16()
	p.Delay = reader.Read16()
}
