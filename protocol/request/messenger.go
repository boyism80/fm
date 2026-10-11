package request

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type Messenger struct {
	Action internal.MessengerAction
	RoomID uint32
	Text   string
}

func (*Messenger) Opcode() byte { return 0x64 }

func (p *Messenger) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(p.Action))
	switch p.Action {
	case internal.MessengerAction_MESSENGER_ACTION_ENTER:
		w.WriteU32(p.RoomID)
	case internal.MessengerAction_MESSENGER_ACTION_INVITE, internal.MessengerAction_MESSENGER_ACTION_DECLINE, internal.MessengerAction_MESSENGER_ACTION_CHAT:
		w.WriteStr16(p.Text)
	}
	return nil
}

func (p *Messenger) Deserialize(r *stream.StreamReader) {
	p.Action = internal.MessengerAction(r.ReadU8())
	switch p.Action {
	case internal.MessengerAction_MESSENGER_ACTION_ENTER:
		p.RoomID = r.ReadU32()
	case internal.MessengerAction_MESSENGER_ACTION_INVITE, internal.MessengerAction_MESSENGER_ACTION_DECLINE, internal.MessengerAction_MESSENGER_ACTION_CHAT:
		p.Text = r.ReadStr16()
	}
}
