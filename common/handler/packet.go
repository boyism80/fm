package handler

import (
	"errors"
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
)

type PacketHandler struct {
	handlers  map[int]func(ctx actor.Context, p interface{})
	generator map[int]func() types.Packet
}

func NewPacketHandler() *PacketHandler {
	return &PacketHandler{
		handlers:  map[int]func(ctx actor.Context, p interface{}){},
		generator: map[int]func() types.Packet{},
	}
}

func (state *PacketHandler) Register(packetType int, handler func(ctx actor.Context, p interface{})) {
	state.handlers[packetType] = handler
}

func RegisterPacketHandler[T any, P any](header int, ctx actor.Context, target *T, h *PacketHandler, fn func(actor.Context, *T, *P)) {
	h.Register(header, func(ctx actor.Context, p interface{}) {
		fn(ctx, target, p.(*P))
	})

	h.generator[header] = func() types.Packet {
		var p any = new(P)
		if pkt, ok := p.(types.Packet); ok {
			return pkt
		}

		return nil

	}
}

func (state *PacketHandler) Handle(ctx actor.Context, header int, data []byte) error {
	generator, exists := state.generator[header]
	if !exists {
		return fmt.Errorf("No handler found for header %d\n", header)
	}

	ptr := generator()
	if ptr == nil {
		return errors.New("Object does not implement types.Packet interface")
	}

	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	err := ptr.Deserialize(reader)
	if err != nil {
		return fmt.Errorf("Failed to deserialize data: %v\n", err)
	}

	handler, exists := state.handlers[header]
	if !exists {
		return fmt.Errorf("No handler registered for header %d\n", header)
	}

	handler(ctx, ptr)
	return nil
}
