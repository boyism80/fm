package request

import (
	"github.com/boyism80/fm/stream"
)

type UseDoor struct {
	OwnerCharacterID uint32
	ToTown           bool
}

func (p *UseDoor) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *UseDoor) Deserialize(reader *stream.StreamReader) {
	p.OwnerCharacterID = reader.ReadU32()
	p.ToTown = reader.ReadU8() == 0
}
