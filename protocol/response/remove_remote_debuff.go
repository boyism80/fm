package response

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
)

// RemoveRemoteDebuff notifies other clients that a character's debuff(s) were removed (opcode 0x91).
type RemoveRemoteDebuff struct {
	CharacterID int32
	Diseases    []constant.DebuffFlag
}

func (p *RemoveRemoteDebuff) Opcode() uint16 { return 0x91 }

func (p *RemoveRemoteDebuff) Serialize(writer *stream.StreamWriter) error {
	if p.Diseases == nil {
		p.Diseases = []constant.DebuffFlag{}
	}
	writer.Write32(p.CharacterID)
	if err := WriteMask(writer, SlotsFromDebuffFlags(p.Diseases)); err != nil {
		return err
	}
	writer.WriteU8(3)
	writer.WriteU8(1)
	return nil
}

func (p *RemoveRemoteDebuff) Deserialize(_ *stream.StreamReader) {}
