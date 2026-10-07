package request

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type RingAction struct {
	Mode        pconst.RingActionMode
	Name        string
	ItemID      uint32
	Accepted    bool
	CharacterID uint32
	MarriageID  uint32
	Slot        uint32
	Wishes      []string
}

func (*RingAction) Opcode() byte { return 0x73 }

func (p *RingAction) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Mode))
	switch p.Mode {
	case pconst.RingActionPropose:
		writer.WriteStr16(p.Name)
		writer.WriteU32(p.ItemID)
	case pconst.RingActionAnswer:
		writer.WriteBoolean(p.Accepted)
		writer.WriteStr16(p.Name)
		writer.WriteU32(p.CharacterID)
	case pconst.RingActionDropRing:
		writer.WriteU32(p.ItemID)
	case pconst.RingActionInviteGuest:
		writer.WriteStr16(p.Name)
		writer.WriteU32(p.MarriageID)
		writer.WriteU32(p.Slot)
	case pconst.RingActionOpenInvitation:
		writer.WriteU32(p.Slot)
		writer.WriteU32(p.ItemID)
	case pconst.RingActionWeddingWishlist:
		writer.WriteU16(uint16(len(p.Wishes)))
		for _, wish := range p.Wishes {
			writer.WriteStr16(wish)
		}
	}
	return nil
}

func (p *RingAction) Deserialize(reader *stream.StreamReader) {
	p.Mode = pconst.RingActionMode(reader.ReadU8())
	switch p.Mode {
	case pconst.RingActionPropose:
		p.Name = reader.ReadStr16()
		p.ItemID = reader.ReadU32()
	case pconst.RingActionAnswer:
		p.Accepted = reader.ReadBool()
		p.Name = reader.ReadStr16()
		p.CharacterID = reader.ReadU32()
	case pconst.RingActionDropRing:
		p.ItemID = reader.ReadU32()
	case pconst.RingActionInviteGuest:
		p.Name = reader.ReadStr16()
		p.MarriageID = reader.ReadU32()
		p.Slot = reader.ReadU32()
	case pconst.RingActionOpenInvitation:
		p.Slot = reader.ReadU32()
		p.ItemID = reader.ReadU32()
	case pconst.RingActionWeddingWishlist:
		count := reader.ReadU16()
		p.Wishes = make([]string, 0, count)
		for range count {
			p.Wishes = append(p.Wishes, reader.ReadStr16())
		}
	}
}
