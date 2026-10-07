package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type UpdateCharacterLook struct {
	Character      *dto.Character
	CrushRing      *dto.Ring
	FriendshipRing *dto.Ring
	MarriageRing   *dto.MarriageRing
}

func writeLookRing(writer *stream.StreamWriter, ring *dto.Ring) {
	if ring == nil {
		writer.WriteU8(0)
		return
	}
	writer.WriteU8(1)
	writer.WriteU64(ring.RingId)
	writer.WriteU64(ring.PartnerId)
	writer.WriteU32(ring.ItemId)
}

func writeMarriageRing(writer *stream.StreamWriter, ring *dto.MarriageRing) {
	if ring == nil {
		writer.WriteU8(0)
		return
	}
	writer.WriteU8(1)
	writer.WriteU32(ring.CharacterID)
	writer.WriteU32(ring.PartnerID)
	writer.WriteU32(ring.ItemID)
}

func (p *UpdateCharacterLook) Serialize(writer *stream.StreamWriter) error {
	if p.Character == nil {
		return nil
	}
	writer.WriteU32(p.Character.ID)
	writer.WriteU8(1)
	p.Character.SerializeLook(writer)
	writeLookRing(writer, p.CrushRing)
	writeLookRing(writer, p.FriendshipRing)
	writeMarriageRing(writer, p.MarriageRing)
	return nil
}

func readLookRing(reader *stream.StreamReader) *dto.Ring {
	if reader.ReadBool() == false {
		return nil
	}
	return &dto.Ring{
		RingId:    reader.ReadU64(),
		PartnerId: reader.ReadU64(),
		ItemId:    reader.ReadU32(),
	}
}

func (p *UpdateCharacterLook) Deserialize(reader *stream.StreamReader) {
	p.Character = &dto.Character{ID: reader.ReadU32()}
	reader.ReadU8()
	p.Character.DeserializeLook(reader)
	p.CrushRing = readLookRing(reader)
	p.FriendshipRing = readLookRing(reader)
	if reader.ReadBool() {
		p.MarriageRing = &dto.MarriageRing{
			CharacterID: reader.ReadU32(),
			PartnerID:   reader.ReadU32(),
			ItemID:      reader.ReadU32(),
		}
	}
}

func (p *UpdateCharacterLook) Opcode() uint16 {
	return 0x8E
}
