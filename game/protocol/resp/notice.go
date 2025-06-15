package resp

import (
	"math"

	"github.com/boyism80/fm/common/stream"
)

// ServerMessageType defines the different SERVERMESSAGE types.
type ServerMessageType uint8

const (
	MsgNotice              ServerMessageType = 0  // [Notice]
	MsgPopup               ServerMessageType = 1  // Popup
	MsgMegaphone           ServerMessageType = 2  // Megaphone
	MsgSuperMegaphone      ServerMessageType = 3  // Super Megaphone
	MsgScrollingTop        ServerMessageType = 4  // Scrolling message at top
	MsgPinkText            ServerMessageType = 5  // Pink Text
	MsgLightBlueText       ServerMessageType = 6  // Lightblue Text
	MsgItemMegaphone       ServerMessageType = 8  // Item megaphone
	MsgHeartMegaphone      ServerMessageType = 9  // Heart megaphone
	MsgSkullSuperMegaphone ServerMessageType = 10 // Skull Super megaphone
	MsgGreenMegaphone      ServerMessageType = 11 // Green megaphone message
	MsgThreeMegaphoneLines ServerMessageType = 12 // Three-line megaphone text
	MsgEOF                 ServerMessageType = 13 // End of file
	MsgAni                 ServerMessageType = 14 // Ani msg
	MsgRedGachaponBox      ServerMessageType = 15 // Red Gachapon box
	MsgBlueNotice          ServerMessageType = 18 // Blue Notice
)

type Notice struct {
	Type    ServerMessageType
	Channel int
	Message string
	MegaEar bool
}

func (p *Notice) Serialize(sw *stream.StreamWriter) error {
	if err := sw.WriteU16(0x33); err != nil {
		return err
	}
	if err := sw.WriteU8(uint8(p.Type)); err != nil {
		return err
	}
	if p.Type == MsgScrollingTop {
		if err := sw.WriteU8(1); err != nil {
			return err
		}
	}
	if err := sw.WriteStr16(p.Message); err != nil {
		return err
	}

	switch p.Type {
	case MsgSuperMegaphone, MsgHeartMegaphone, MsgSkullSuperMegaphone:
		channel := int(math.Max(1, float64(p.Channel))) - 1
		if err := sw.WriteU8(uint8(channel)); err != nil {
			return err
		}
		if err := sw.WriteBoolean(p.MegaEar); err != nil {
			return err
		}

	case MsgLightBlueText, MsgBlueNotice:
		var item uint32
		if p.Channel >= 1_000_000 && p.Channel < 6_000_000 {
			item = uint32(p.Channel)
		}
		if err := sw.WriteU32(item); err != nil {
			return err
		}
	}

	return nil
}

func (p *Notice) Deserialize(sr *stream.StreamReader) error {
	return nil
}
