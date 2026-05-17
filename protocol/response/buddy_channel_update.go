package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type BuddyChannelUpdate struct {
	CharacterID uint32
	Channel     int32
}

func (p *BuddyChannelUpdate) Opcode() uint16 {
	return 0x2E
}

func (p *BuddyChannelUpdate) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.BuddyS2CChannelUpdate))
	w.WriteU32(p.CharacterID)
	w.WriteU8(0)
	w.Write32(p.Channel)
	return nil
}

func (p *BuddyChannelUpdate) Deserialize(*stream.StreamReader) {}
