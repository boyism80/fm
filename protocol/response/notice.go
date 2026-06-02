package response

import (
	"math"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
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
	sw.WriteU8(uint8(p.Type))
	if p.Type == constant.MsgScrollingTop {
		sw.WriteU8(1)
	}
	sw.WriteStr16(p.Message)

	switch p.Type {
	case constant.MsgSuperMegaphone, constant.MsgHeartMegaphone, constant.MsgSkullSuperMegaphone:
		channel := int(math.Max(1, float64(p.Channel))) - 1
		sw.WriteU8(uint8(channel))
		sw.WriteBoolean(p.MegaEar)

	case constant.MsgLightBlueText, constant.MsgBlueNotice:
		var item uint32
		if p.Channel >= 1_000_000 && p.Channel < 6_000_000 {
			item = uint32(p.Channel)
		}
		sw.WriteU32(item)
	}

	return nil
}

func (p *Notice) Deserialize(sr *stream.StreamReader) {
}
