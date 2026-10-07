package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type SpawnPet struct {
	CharacterID uint32
	Pet         *dto.ActivePet
	Reason      constant.PetRemoveReason
}

func (p *SpawnPet) Opcode() uint16 {
	return 0x75
}

func (p *SpawnPet) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteBoolean(p.Pet != nil)
	if p.Pet != nil {
		p.Pet.Serialize(writer)
	} else {
		writer.WriteU8(uint8(p.Reason))
	}
	return nil
}

func (p *SpawnPet) Deserialize(reader *stream.StreamReader) {
	p.CharacterID = reader.ReadU32()
	if reader.ReadBool() {
		p.Pet = &dto.ActivePet{}
		p.Pet.Deserialize(reader)
	} else {
		p.Reason = constant.PetRemoveReason(reader.ReadU8())
	}
}
