package core

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/stream"
)

type Request interface {
	Deserialize(reader *stream.StreamReader)
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
	GetServerContext() ServerContext
}

func Bind[S ServerRegistry, C HandlerConstructor[S, H, T, U], H Handler[T], T RequestPtr[U], U any](registry S) {
	var constructor C
	handler := constructor.New(registry)

	opcode := handler.GetOpcode()

	handlerFunc := func(ctx *ClientContext, data []byte) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("deserialize panic for packet 0x%X from %s: %v", opcode, ctx.Client.GetConnection().RemoteAddr(), r)
				err = fmt.Errorf("deserialize: %v", r)
			}
		}()

		reader := stream.NewStreamReader(&data, stream.LittleEndian)

		var req T = new(U)
		req.Deserialize(reader)

		if GetPacketLogEnabled() {
			log.Printf("recv 0x%04X: %+v", opcode, req)
		}

		return handler.Handle(ctx, req)
	}

	registry.GetServerContext().GetPacketHandler().RegisterHandler(int(opcode), handlerFunc)
	log.Printf("Registered packet handler 0x%X for type %T", opcode, handler)
}
