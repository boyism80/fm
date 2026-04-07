package response

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
)

// CancelRemoteBuff notifies other clients on the map that a character's buff(s) were removed.
type CancelRemoteBuff struct {
	CharacterID int32
	Buffs       []constant.BuffFlag
}

func (p *CancelRemoteBuff) Opcode() uint16 { return 0x91 }

func (p *CancelRemoteBuff) Serialize(writer *stream.StreamWriter) error {
	if p.Buffs == nil {
		p.Buffs = []constant.BuffFlag{}
	}
	writer.Write32(p.CharacterID)
	return WriteMask(writer, SlotsFromBuffFlags(p.Buffs))
}

func (p *CancelRemoteBuff) Deserialize(_ *stream.StreamReader) {}
