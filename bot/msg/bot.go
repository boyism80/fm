package msg

import "github.com/boyism80/fm/core/types"

type BotConnect struct {
	Port int
}

type BotSendPacket struct {
	Packet types.Packet
}

type BotClose struct{}
