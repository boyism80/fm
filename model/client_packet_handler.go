package model

import (
	"log"

	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/packet"
)

func RegisterPacketHandler(client *ClientActor, h *handler.PacketHandler) {
	handler.RegisterPacketHandler(1, client, h, onMove)
	handler.RegisterPacketHandler(2, client, h, onAttack)
}

func onMove(client *ClientActor, request *packet.MovePacket) {
	log.Printf("onMove")
}

func onAttack(client *ClientActor, request *packet.MovePacket) {
	log.Printf("onAttack")
}
