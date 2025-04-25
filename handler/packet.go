package handler

import (
	"errors"
	"fmt"

	"github.com/boyism80/fm/packet"
	"github.com/boyism80/fm/stream"
)

type PacketHandler struct {
	handlers  map[int]func(p interface{})
	generator map[int]func() packet.Packet
}

func NewPacketHandler() *PacketHandler {
	return &PacketHandler{
		handlers:  map[int]func(p interface{}){},
		generator: map[int]func() packet.Packet{},
	}
}

func (state *PacketHandler) Register(packetType int, handler func(p interface{})) {
	state.handlers[packetType] = handler
}

func RegisterPacketHandler[T any, P any](header int, target *T, h *PacketHandler, fn func(*T, *P)) {
	h.Register(header, func(p interface{}) {
		fn(target, p.(*P))
	})

	h.generator[header] = func() packet.Packet {
		var p any = new(P)
		if pkt, ok := p.(packet.Packet); ok {
			return pkt
		}

		return nil

	}
}

func (state *PacketHandler) Handle(header int, data []byte) error {
	generator, exists := state.generator[header]
	if !exists {
		return fmt.Errorf("No handler found for header %d\n", header)
	}

	ptr := generator()
	if ptr == nil {
		return errors.New("Object does not implement packet.Packet interface")
	}

	reader := stream.NewStreamReader(data, stream.LittleEndian)
	err := ptr.Deserialize(reader)
	if err != nil {
		return fmt.Errorf("Failed to deserialize data: %v\n", err)
	}

	handler, exists := state.handlers[header]
	if !exists {
		return fmt.Errorf("No handler registered for header %d\n", header)
	}

	handler(ptr)
	return nil
}
