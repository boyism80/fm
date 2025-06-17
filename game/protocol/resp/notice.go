package resp

import (
	"math"

	"github.com/boyism80/fm/common/stream"
)

// ServerMessageType defines the different SERVERMESSAGE types.
type ServerMessageType uint8

const (
	MSG_NOTICE                ServerMessageType = 0  // [Notice]
	MSG_POPUP                 ServerMessageType = 1  // Popup
	MSG_MEGAPHONE             ServerMessageType = 2  // Megaphone
	MSG_SUPER_MEGAPHONE       ServerMessageType = 3  // Super Megaphone
	MSG_SCROLLING_TOP         ServerMessageType = 4  // Scrolling message at top
	MSG_PINK_TEXT             ServerMessageType = 5  // Pink Text
	MSG_LIGHT_BLUE_TEXT       ServerMessageType = 6  // Lightblue Text
	MSG_ITEM_MEGAPHONE        ServerMessageType = 8  // Item megaphone
	MSG_HEART_MEGAPHONE       ServerMessageType = 9  // Heart megaphone
	MSG_SKULL_SUPER_MEGAPHONE ServerMessageType = 10 // Skull Super megaphone
	MSG_GREEN_MEGAPHONE       ServerMessageType = 11 // Green megaphone message
	MSG_THREE_MEGAPHONE_LINES ServerMessageType = 12 // Three-line megaphone text
	MSG_EOF                   ServerMessageType = 13 // End of file
	MSG_ANI                   ServerMessageType = 14 // Ani msg
	MSG_RED_GACHAPON_BOX      ServerMessageType = 15 // Red Gachapon box
	MSG_BLUE_NOTICE           ServerMessageType = 18 // Blue Notice
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
	if p.Type == MSG_SCROLLING_TOP {
		if err := sw.WriteU8(1); err != nil {
			return err
		}
	}
	if err := sw.WriteStr16(p.Message); err != nil {
		return err
	}

	switch p.Type {
	case MSG_SUPER_MEGAPHONE, MSG_HEART_MEGAPHONE, MSG_SKULL_SUPER_MEGAPHONE:
		channel := int(math.Max(1, float64(p.Channel))) - 1
		if err := sw.WriteU8(uint8(channel)); err != nil {
			return err
		}
		if err := sw.WriteBoolean(p.MegaEar); err != nil {
			return err
		}

	case MSG_LIGHT_BLUE_TEXT, MSG_BLUE_NOTICE:
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
