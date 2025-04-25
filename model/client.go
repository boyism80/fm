package model

import (
	"net"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
)

type ClientActor struct {
	Conn           net.Conn
	Buffer         []byte
	messageHandler *handler.MessageHandler
	packetHandler  *handler.PacketHandler
}

func NewClientActor(conn net.Conn) actor.Actor {
	act := &ClientActor{
		Conn:           conn,
		Buffer:         []byte{},
		messageHandler: handler.NewMessageHandler(),
		packetHandler:  handler.NewPacketHandler(),
	}

	RegisterClientHandlers(act, act.messageHandler)
	RegisterPacketHandler(act, act.packetHandler)
	return act
}

func (state *ClientActor) Receive(context actor.Context) {
	state.messageHandler.Handle(context)
}

func (state *ClientActor) Invoke(header int, data []byte) error {
	return state.packetHandler.Handle(header, data)
}
