package response

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type FamilyDescendants struct {
	CharacterID uint32
	Count       uint32
}

type FamilyChart struct {
	CharacterID  uint32
	Members      []*internal.FamilyMember
	GrandJuniors []FamilyDescendants
	MemberCount  uint32
	UsedFeatures []internal.FamilyFeature
}

func (p *FamilyChart) Opcode() uint16 {
	return 0x49
}

func (p *FamilyChart) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.CharacterID)

	w.WriteU32(uint32(len(p.Members)))
	for _, m := range p.Members {
		w.WriteU32(m.GetCharacterId())
		w.WriteU32(m.GetSeniorId())
		w.WriteU16(uint16(m.GetClassId()))
		w.WriteU8(uint8(m.GetLevel()))
		w.WriteBoolean(m.GetChannelId() >= 0)
		w.WriteU32(m.GetReputation())
		w.WriteU32(m.GetTotalReputation())
		w.WriteU32(0)
		w.WriteU32(0)
		w.Write32(max(0, m.GetChannelId()))
		w.WriteU32(0)
		w.WriteStr16(m.GetName())
	}

	w.WriteU32(uint32(2 + len(p.GrandJuniors)))
	w.WriteU32(0)
	w.WriteU32(uint32(len(p.GrandJuniors)))
	w.Write32(-1)
	w.WriteU32(p.MemberCount)
	for _, d := range p.GrandJuniors {
		w.WriteU32(d.CharacterID)
		w.WriteU32(d.Count)
	}

	w.WriteU32(uint32(len(p.UsedFeatures)))
	for _, f := range p.UsedFeatures {
		w.WriteU32(uint32(f))
		w.WriteU32(1)
	}
	w.WriteU16(2)
	return nil
}

func (p *FamilyChart) Deserialize(reader *stream.StreamReader) {
	p.CharacterID = reader.ReadU32()
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
		p.Members = append(p.Members, m)
	}

	pairs := int(reader.ReadU32())
	for range pairs {
		key := reader.Read32()
		value := reader.ReadU32()
		switch key {
		case 0:
		case -1:
			p.MemberCount = value
		default:
			p.GrandJuniors = append(p.GrandJuniors, FamilyDescendants{CharacterID: uint32(key), Count: value})
		}
	}
	features := int(reader.ReadU32())
	for range features {
		p.UsedFeatures = append(p.UsedFeatures, internal.FamilyFeature(reader.ReadU32()))
		reader.Skip(4)
	}
	reader.Skip(2)
}
