package request

import "github.com/boyism80/fm/stream"

// PartySearchStart payload (C2S 0xB5), matching legacy Java parsing order:
// min level, max level, members needed, class bitmask.
type PartySearchStart struct {
	MinLevel      int32
	MaxLevel      int32
	MembersNeeded int32
	ClassMask     int32
}

func (p *PartySearchStart) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *PartySearchStart) Deserialize(reader *stream.StreamReader) {
	// Packet is fixed 16 bytes (4 x int32) in Java.
	p.MinLevel = reader.Read32()
	p.MaxLevel = reader.Read32()
	p.MembersNeeded = reader.Read32()
	p.ClassMask = reader.Read32()
}
