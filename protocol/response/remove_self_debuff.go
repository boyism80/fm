package response

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type RemoveSelfDebuff struct {
	Diseases []constant.DebuffFlag
}

func (p *RemoveSelfDebuff) Opcode() uint16 { return 0x16 }

func (p *RemoveSelfDebuff) Serialize(writer *stream.StreamWriter) error {
	if p.Diseases == nil {
		p.Diseases = []constant.DebuffFlag{}
	}
	WriteDebuffs(writer, p.Diseases)
	writer.WriteU8(1)
	return nil
}

func (p *RemoveSelfDebuff) Deserialize(_ *stream.StreamReader) {}
