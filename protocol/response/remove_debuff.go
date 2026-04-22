package response

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type RemoveDebuff struct {
	CharacterID int32
	Diseases    []constant.DebuffFlag
}

func (p *RemoveDebuff) Opcode() uint16 { return 0x91 }

func (p *RemoveDebuff) Serialize(writer *stream.StreamWriter) error {
	if p.Diseases == nil {
		p.Diseases = []constant.DebuffFlag{}
	}
	writer.Write32(p.CharacterID)
	WriteDebuffs(writer, p.Diseases)
	writer.WriteU8(3)
	writer.WriteU8(1)
	return nil
}

func (p *RemoveDebuff) Deserialize(_ *stream.StreamReader) {}
