package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type PartyUpdateLeaderChange struct {
	NewLeaderCharacterID uint32
	ByDisconnect         bool
}

func (p *PartyUpdateLeaderChange) Opcode() uint16 {
	return 0x2D
}

func (p *PartyUpdateLeaderChange) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.PartyS2CLeaderChange))
	w.WriteU32(p.NewLeaderCharacterID)
	if p.ByDisconnect {
		w.WriteU8(1)
	} else {
		w.WriteU8(0)
	}
	return nil
}

func (p *PartyUpdateLeaderChange) Deserialize(*stream.StreamReader) {}
