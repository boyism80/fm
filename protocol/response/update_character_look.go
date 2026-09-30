package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type UpdateCharacterLook struct {
	Character       *dto.Character
	CrushRings      []*dto.Ring
	FriendshipRings []*dto.Ring
	MarriageRings   []*dto.Ring
}

func writeLookRings(writer *stream.StreamWriter, rings []*dto.Ring) {
	count := 0
	for _, ring := range rings {
		if ring != nil {
			count++
		}
	}
	writer.WriteU8(uint8(count))
	for _, ring := range rings {
		if ring == nil {
			continue
		}
		writer.WriteU64(ring.RingId)
		writer.WriteU64(ring.PartnerId)
		writer.WriteU32(ring.ItemId)
	}
}

func writeMarriageRings(writer *stream.StreamWriter, characterID uint32, rings []*dto.Ring) {
	count := 0
	for _, ring := range rings {
		if ring != nil {
			count++
		}
	}
	writer.WriteU8(uint8(count))
	for _, ring := range rings {
		if ring == nil {
			continue
		}
		writer.WriteU32(characterID)
		writer.WriteU32(ring.PartnerChrId)
		writer.WriteU32(ring.ItemId)
	}
}

func (p *UpdateCharacterLook) Serialize(writer *stream.StreamWriter) error {
	if p.Character == nil {
		return nil
	}
	writer.WriteU32(p.Character.ID)
	writer.WriteU8(1)
	p.Character.SerializeLook(writer)
	writeLookRings(writer, p.CrushRings)
	writeLookRings(writer, p.FriendshipRings)
	writeMarriageRings(writer, p.Character.ID, p.MarriageRings)
	return nil
}

func (p *UpdateCharacterLook) Deserialize(reader *stream.StreamReader) {
}

func (p *UpdateCharacterLook) Opcode() uint16 {
	return 0x8E
}
