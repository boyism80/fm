package handler

import (
	"fmt"
	"reflect"

	"github.com/boyism80/fm/packet"
	"github.com/boyism80/fm/stream"
)

type PacketHandler struct {
	handlers map[int]func(p interface{})
	types    map[int]reflect.Type
}

func NewPacketHandler() *PacketHandler {
	return &PacketHandler{
		handlers: make(map[int]func(p interface{})),
	}
}

func (state *PacketHandler) Register(packetType int, handler func(p interface{})) {
	state.handlers[packetType] = handler
}

func RegisterPacketHandler[T any, P any](header int, target *T, h *PacketHandler, fn func(*T, *P)) {
	h.Register(header, func(p interface{}) {
		fn(target, p.(*P))
	})

	var p *P
	h.types[header] = reflect.TypeOf(p)
}

func (state *PacketHandler) Handle(header int, data []byte) {
	packetType, exists := state.types[header]
	if !exists {
		fmt.Printf("No handler found for header %d\n", header)
		return
	}

	ptr := reflect.New(packetType).Interface()
	packet, ok := ptr.(packet.Packet)
	if !ok {
		fmt.Println("Object does not implement packet.Packet interface")
		return
	}

	reader := stream.NewStreamReader(data, stream.LittleEndian)
	err := packet.Deserialize(reader)
	if err != nil {
		fmt.Printf("Failed to deserialize data: %v\n", err)
		return
	}

	handler, exists := state.handlers[header]
	if !exists {
		fmt.Printf("No handler registered for header %d\n", header)
		return
	}

	handler(packet)
}
