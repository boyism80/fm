package response

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
)

type RemoveDebuff struct {
	Diseases []constant.DebuffFlag
}

func (p *RemoveDebuff) Opcode() uint16 { return 0x16 }

func (p *RemoveDebuff) Serialize(writer *stream.StreamWriter) error {
	if p.Diseases == nil {
		p.Diseases = []constant.DebuffFlag{}
	}
	WriteDebuffs(writer, p.Diseases)
	writer.WriteU8(1)
	return nil
}

func (p *RemoveDebuff) Deserialize(_ *stream.StreamReader) {}
