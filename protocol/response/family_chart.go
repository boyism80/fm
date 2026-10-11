package response

import (
	"slices"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type FamilyChart struct {
	State *internal.FamilyState
}

func (p *FamilyChart) Opcode() uint16 {
	return 0x49
}

func (p *FamilyChart) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.State.CharacterId)

	w.WriteU32(uint32(len(p.State.Members)))
	for _, m := range p.State.Members {
		w.WriteU32(m.CharacterId)
		w.WriteU32(m.SeniorId)
		w.WriteU16(uint16(m.ClassId))
		w.WriteU8(uint8(m.Level))
		w.WriteBoolean(m.ChannelId >= 0)
		w.WriteU32(m.Reputation)
		w.WriteU32(m.TotalReputation)
		w.WriteU32(0)
		w.WriteU32(0)
		w.Write32(max(0, m.ChannelId))
		w.WriteU32(0)
		w.WriteStr16(m.Name)
	}

	var grandJuniors []*internal.FamilyMember
	for _, m := range p.State.Members {
		if slices.Contains(p.State.Juniors, m.SeniorId) {
			grandJuniors = append(grandJuniors, m)
		}
	}
	w.WriteU32(uint32(2 + len(grandJuniors)))
	w.WriteU32(0)
	w.WriteU32(uint32(len(grandJuniors)))
	w.Write32(-1)
	w.WriteU32(p.State.MemberCount)
	for _, m := range grandJuniors {
		w.WriteU32(m.CharacterId)
		w.WriteU32(m.Descendants)
	}

	w.WriteU32(uint32(len(p.State.UsedFeatures)))
	for _, f := range p.State.UsedFeatures {
		w.WriteU32(uint32(f))
		w.WriteU32(1)
	}
	w.WriteU16(2)
	return nil
}

func (p *FamilyChart) Deserialize(reader *stream.StreamReader) {
	p.State = &internal.FamilyState{CharacterId: reader.ReadU32()}
	count := int(reader.ReadU32())
	for range count {
		m := &internal.FamilyMember{
			CharacterId: reader.ReadU32(),
			SeniorId:    reader.ReadU32(),
			ClassId:     uint32(reader.ReadU16()),
			Level:       uint32(reader.ReadU8()),
		}
		online := reader.ReadBool()
		m.Reputation = reader.ReadU32()
		m.TotalReputation = reader.ReadU32()
		reader.Skip(8)
		m.ChannelId = reader.Read32()
		if !online {
			m.ChannelId = -1
		}
		reader.Skip(4)
		m.Name = reader.ReadStr16()
		p.State.Members = append(p.State.Members, m)
	}

	pairs := int(reader.ReadU32())
	for range pairs {
		key := reader.Read32()
		value := reader.ReadU32()
		if key == -1 {
			p.State.MemberCount = value
		}
	}
	features := int(reader.ReadU32())
	for range features {
		p.State.UsedFeatures = append(p.State.UsedFeatures, internal.FamilyFeature(reader.ReadU32()))
		reader.Skip(4)
	}
	reader.Skip(2)
}
