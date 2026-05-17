package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type BuddyListUpdate struct {
	Action  constant.BuddyListSyncAction
	Entries []BuddyEntry
}

func (p *BuddyListUpdate) Opcode() uint16 {
	return 0x2E
}

func (p *BuddyListUpdate) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(p.Action))
	w.WriteU8(uint8(len(p.Entries)))
	for _, e := range p.Entries {
		w.WriteU32(e.CharacterID)
		w.WriteStaticStr(e.Name, buddyNameFieldLen)
		w.WriteBoolean(e.Pending)
		w.Write32(e.Channel)
		w.WriteStaticStr(e.Group, buddyGroupFieldLen)
	}
	for range p.Entries {
		w.Write32(0)
	}
	return nil
}

func (p *BuddyListUpdate) Deserialize(*stream.StreamReader) {}
