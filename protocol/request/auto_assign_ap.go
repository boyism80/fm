package request

import (
	"github.com/boyism80/fm/stream"
)

type AutoAssignAPEntry struct {
	Stat   uint32
	Amount uint32
}

type AutoAssignAP struct {
	Tick    uint32
	Entries []AutoAssignAPEntry
}

func (*AutoAssignAP) Opcode() byte { return 0x47 }

func (p *AutoAssignAP) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Tick)
	writer.WriteU32(uint32(len(p.Entries)))
	for _, e := range p.Entries {
		writer.WriteU32(e.Stat)
		writer.WriteU32(e.Amount)
	}
	return nil
}

func (p *AutoAssignAP) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	count := reader.ReadU32()
	for range count {
		p.Entries = append(p.Entries, AutoAssignAPEntry{
			Stat:   reader.ReadU32(),
			Amount: reader.ReadU32(),
		})
	}
}
