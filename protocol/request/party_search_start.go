package request

import "github.com/boyism80/fm/stream"

type PartySearchStart struct {
	MinLevel      int32
	MaxLevel      int32
	MembersNeeded int32
	ClassMask     int32
}

func (*PartySearchStart) Opcode() byte { return 0xB5 }

func (p *PartySearchStart) Serialize(writer *stream.StreamWriter) error {
	writer.Write32(p.MinLevel)
	writer.Write32(p.MaxLevel)
	writer.Write32(p.MembersNeeded)
	writer.Write32(p.ClassMask)
	return nil
}

func (p *PartySearchStart) Deserialize(reader *stream.StreamReader) {
	p.MinLevel = reader.Read32()
	p.MaxLevel = reader.Read32()
	p.MembersNeeded = reader.Read32()
	p.ClassMask = reader.Read32()
}
