package msg

import "github.com/boyism80/fm/packet"

type BotConnect struct {
	Port int
}

type BotSendPacket struct {
	Packet packet.Packet
}
