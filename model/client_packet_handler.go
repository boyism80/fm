package model

import (
	"log"

	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/packet/req"
)

func RegisterPacketHandler(client *ClientActor, h *handler.PacketHandler) {
	handler.RegisterPacketHandler(0x0A, client, h, onPong)
}

func onPong(client *ClientActor, request *req.Pong) {
	log.Printf("pong")
}
