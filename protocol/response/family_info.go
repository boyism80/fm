package response

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type FamilyInfo struct {
	State *internal.FamilyState
}

func (p *FamilyInfo) Opcode() uint16 {
	return 0x4A
}

func (p *FamilyInfo) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.State.Reputation)
	w.WriteU32(p.State.TotalReputation)
	w.WriteU32(0)
	w.WriteU16(uint16(len(p.State.Juniors)))
	w.WriteU16(2)
	w.WriteU16(uint16(len(p.State.Juniors)))
	w.WriteU32(p.State.LeaderId)
	w.WriteStr16(p.State.LeaderName)
	w.WriteStr16(p.State.Notice)
	w.WriteU32(uint32(len(p.State.UsedFeatures)))
	for _, f := range p.State.UsedFeatures {
		w.WriteU32(uint32(f))
		w.WriteU32(1)
	}
	return nil
}

func (p *FamilyInfo) Deserialize(reader *stream.StreamReader) {
	p.State = &internal.FamilyState{
		Reputation:      reader.ReadU32(),
		TotalReputation: reader.ReadU32(),
	}
	reader.Skip(4)
	p.State.Juniors = make([]uint32, reader.ReadU16())
	reader.Skip(4)
	p.State.LeaderId = reader.ReadU32()
	p.State.LeaderName = reader.ReadStr16()
	p.State.Notice = reader.ReadStr16()
	count := int(reader.ReadU32())
	for range count {
		p.State.UsedFeatures = append(p.State.UsedFeatures, internal.FamilyFeature(reader.ReadU32()))
		reader.Skip(4)
	}
}
