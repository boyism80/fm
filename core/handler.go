package core

import (
	"log"

	"github.com/boyism80/fm/stream"
)

type Request interface {
	Deserialize(reader *stream.StreamReader) error
}

type RequestPtr[T any] interface {
	Request
	*T
}

type Handler[T Request] interface {
	Handle(ctx *ClientContext, req T) error
	GetOpcode() byte
}

type HandlerConstructor[S any, H Handler[T], T RequestPtr[U], U any] interface {
	New(S) H
}

type ServerRegistry interface {
	GetServer() *Server
}

func Bind[S ServerRegistry, C HandlerConstructor[S, H, T, U], H Handler[T], T RequestPtr[U], U any](registry S) {
	var constructor C
	handler := constructor.New(registry)

	opcode := handler.GetOpcode()

	handlerFunc := func(ctx *ClientContext, data []byte) error {
		reader := stream.NewStreamReader(&data, stream.LittleEndian)

		var req T = new(U)

		if err := req.Deserialize(reader); err != nil {
			log.Printf("Failed to deserialize packet 0x%X from %s: %v", opcode, ctx.Client.GetConnection().RemoteAddr(), err)
			return err
		}

		return handler.Handle(ctx, req)
	}

	registry.GetServer().RegisterPacketHandler(int(opcode), handlerFunc)
	log.Printf("Registered packet handler 0x%X for type %T", opcode, handler)
}
