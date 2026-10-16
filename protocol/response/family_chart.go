package response

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type FamilyDescendants struct {
	CharacterID uint32
	Count       uint32
}

type FamilyMember struct {
	CharacterID     uint32
	SeniorID        uint32
	Name            string
	Level           uint32
	ClassID         uint32
	Reputation      uint32
	TotalReputation uint32
	ChannelID       int32
}

type FamilyChart struct {
	CharacterID  uint32
	Members      []FamilyMember
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
		w.WriteU32(m.CharacterID)
		w.WriteU32(m.SeniorID)
		w.WriteU16(uint16(m.ClassID))
		w.WriteU8(uint8(m.Level))
		w.WriteBoolean(m.ChannelID >= 0)
		w.WriteU32(m.Reputation)
		w.WriteU32(m.TotalReputation)
		w.WriteU32(0)
		w.WriteU32(0)
		w.Write32(max(0, m.ChannelID))
		w.WriteU32(0)
		w.WriteStr16(m.Name)
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
		m := FamilyMember{
			CharacterID: reader.ReadU32(),
			SeniorID:    reader.ReadU32(),
			ClassID:     uint32(reader.ReadU16()),
			Level:       uint32(reader.ReadU8()),
		}
		online := reader.ReadBool()
		m.Reputation = reader.ReadU32()
		m.TotalReputation = reader.ReadU32()
		reader.Skip(8)
		m.ChannelID = reader.Read32()
		if !online {
			m.ChannelID = -1
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
