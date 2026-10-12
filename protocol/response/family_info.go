package response

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type FamilyInfo struct {
	Reputation      uint32
	TotalReputation uint32
	Juniors         uint16
	LeaderID        uint32
	LeaderName      string
	Precept         string
	UsedFeatures    []internal.FamilyFeature
}

func (p *FamilyInfo) Opcode() uint16 {
	return 0x4A
}

func (p *FamilyInfo) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.Reputation)
	w.WriteU32(p.TotalReputation)
	w.WriteU32(0)
	w.WriteU16(p.Juniors)
	w.WriteU16(2)
	w.WriteU16(p.Juniors)
	w.WriteU32(p.LeaderID)
	w.WriteStr16(p.LeaderName)
	w.WriteStr16(p.Precept)
	w.WriteU32(uint32(len(p.UsedFeatures)))
	for _, f := range p.UsedFeatures {
		w.WriteU32(uint32(f))
		w.WriteU32(1)
	}
	return nil
}

func (p *FamilyInfo) Deserialize(reader *stream.StreamReader) {
	p.Reputation = reader.ReadU32()
	p.TotalReputation = reader.ReadU32()
	reader.Skip(4)
	p.Juniors = reader.ReadU16()
	reader.Skip(4)
	p.LeaderID = reader.ReadU32()
	p.LeaderName = reader.ReadStr16()
	p.Precept = reader.ReadStr16()
	count := int(reader.ReadU32())
	for range count {
		p.UsedFeatures = append(p.UsedFeatures, internal.FamilyFeature(reader.ReadU32()))
		reader.Skip(4)
	}
}
