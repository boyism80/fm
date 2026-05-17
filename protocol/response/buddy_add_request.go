package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type BuddyAddRequest struct {
	FromCharacterID uint32
	FromName        string
}

func (p *BuddyAddRequest) Opcode() uint16 {
	return 0x2E
}

func (p *BuddyAddRequest) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.BuddyS2CAddRequest))
	w.WriteU32(p.FromCharacterID)
	w.WriteStr16(p.FromName)
	w.WriteU32(p.FromCharacterID)
	w.WriteStaticStr(p.FromName, buddyNameFieldLen)
	w.WriteU8(1)
	w.Write32(0)
	w.WriteStaticStr(constant.BuddyDefaultGroup, buddyGroupFieldLen)
	w.WriteU8(0)
	return nil
}

func (p *BuddyAddRequest) Deserialize(*stream.StreamReader) {}
