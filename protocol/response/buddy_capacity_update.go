package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type BuddyCapacityUpdate struct {
	Capacity uint8
}

func (p *BuddyCapacityUpdate) Opcode() uint16 {
	return 0x2E
}

func (p *BuddyCapacityUpdate) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.BuddyS2CCapacityUpdate))
	w.WriteU8(p.Capacity)
	return nil
}

func (p *BuddyCapacityUpdate) Deserialize(*stream.StreamReader) {}
