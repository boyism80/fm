package resp

import (
	"math"

	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/constant"
)

type Notice struct {
	Type    constant.ServerMessageType
	Channel int
	Message string
	MegaEar bool
}

func (p *Notice) Opcode() uint16 {
	return 0x33
}

func (p *Notice) Serialize(sw *stream.StreamWriter) error {
	if err := sw.WriteU8(uint8(p.Type)); err != nil {
		return err
	}
	if p.Type == constant.MSG_SCROLLING_TOP {
		if err := sw.WriteU8(1); err != nil {
			return err
		}
	}
	if err := sw.WriteStr16(p.Message); err != nil {
		return err
	}

	switch p.Type {
	case constant.MSG_SUPER_MEGAPHONE, constant.MSG_HEART_MEGAPHONE, constant.MSG_SKULL_SUPER_MEGAPHONE:
		channel := int(math.Max(1, float64(p.Channel))) - 1
		if err := sw.WriteU8(uint8(channel)); err != nil {
			return err
		}
		if err := sw.WriteBoolean(p.MegaEar); err != nil {
			return err
		}

	case constant.MSG_LIGHT_BLUE_TEXT, constant.MSG_BLUE_NOTICE:
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
