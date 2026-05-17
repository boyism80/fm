package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type BuddyStatus struct {
	Code constant.BuddyStatusCode
}

func (p *BuddyStatus) Opcode() uint16 {
	return 0x2E
}

func (p *BuddyStatus) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(p.Code))
	return nil
}

func (p *BuddyStatus) Deserialize(*stream.StreamReader) {}
