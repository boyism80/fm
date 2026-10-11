package response

import (
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type Messenger struct {
	Action   internal.MessengerAction
	Position uint8
	Name     string
	Text     string
	RoomID   uint32
	Success  bool
	Channel  uint16
	Look     []byte
}

func (p *Messenger) Opcode() uint16 {
	return 0xEE
}

func (p *Messenger) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(p.Action))
	switch p.Action {
	case internal.MessengerAction_MESSENGER_ACTION_ENTER:
		w.WriteU8(p.Position)
		w.Write(p.Look)
		w.WriteStr16(p.Name)
		w.WriteU8(uint8(p.Channel))
		w.WriteBoolean(true)
	case internal.MessengerAction_MESSENGER_ACTION_AVATAR:
		w.WriteU8(p.Position)
		w.Write(p.Look)
	case internal.MessengerAction_MESSENGER_ACTION_SELF_ENTER, internal.MessengerAction_MESSENGER_ACTION_LEAVE:
		w.WriteU8(p.Position)
	case internal.MessengerAction_MESSENGER_ACTION_INVITE:
		w.WriteStr16(p.Name)
		w.WriteU8(0)
		w.WriteU32(p.RoomID)
		w.WriteU8(0)
	case internal.MessengerAction_MESSENGER_ACTION_INVITE_RESULT, internal.MessengerAction_MESSENGER_ACTION_DECLINE:
		w.WriteStr16(p.Name)
		w.WriteBoolean(p.Success)
	case internal.MessengerAction_MESSENGER_ACTION_CHAT:
		w.WriteStr16(p.Text)
	}
	return nil
}

func (p *Messenger) Deserialize(reader *stream.StreamReader) {
	p.Action = internal.MessengerAction(reader.ReadU8())
	switch p.Action {
	case internal.MessengerAction_MESSENGER_ACTION_ENTER, internal.MessengerAction_MESSENGER_ACTION_AVATAR:
		p.Position = reader.ReadU8()
		rest := reader.Read(reader.Remaining())
		look := stream.NewStreamReader(&rest, stream.LittleEndian)
		(&dto.Character{}).DeserializeLook(look)
		p.Look = rest[:len(rest)-look.Remaining()]
		if p.Action == internal.MessengerAction_MESSENGER_ACTION_ENTER {
			p.Name = look.ReadStr16()
			p.Channel = uint16(look.ReadU8())
			look.ReadBool()
		}
	case internal.MessengerAction_MESSENGER_ACTION_SELF_ENTER, internal.MessengerAction_MESSENGER_ACTION_LEAVE:
		p.Position = reader.ReadU8()
	case internal.MessengerAction_MESSENGER_ACTION_INVITE:
		p.Name = reader.ReadStr16()
		reader.ReadU8()
		p.RoomID = reader.ReadU32()
		reader.ReadU8()
	case internal.MessengerAction_MESSENGER_ACTION_INVITE_RESULT, internal.MessengerAction_MESSENGER_ACTION_DECLINE:
		p.Name = reader.ReadStr16()
		p.Success = reader.ReadBool()
	case internal.MessengerAction_MESSENGER_ACTION_CHAT:
		p.Text = reader.ReadStr16()
	}
}
