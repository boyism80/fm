package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type KeyMap struct {
	UseEmptyDefaultBranch bool
	Slots                 map[int]dto.KeyBinding
}

func (p *KeyMap) Opcode() uint16 {
	return 0xFD
}

func (p *KeyMap) Serialize(writer *stream.StreamWriter) error {
	if p.UseEmptyDefaultBranch || len(p.Slots) == 0 {
		writer.WriteU8(1)
		return nil
	}
	writer.WriteU8(0)
	for i := 0; i < 89; i++ {
		if b, ok := p.Slots[i]; ok {
			writer.WriteU8(b.Type)
			writer.Write32(b.Action)
		} else {
			writer.WriteU8(0)
			writer.Write32(0)
		}
	}
	return nil
}

func (p *KeyMap) Deserialize(reader *stream.StreamReader) {}
