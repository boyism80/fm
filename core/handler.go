package core

import (
	"github.com/boyism80/fm/stream"
)

type Request interface {
	Deserialize(reader *stream.StreamReader)
	Opcode() byte
}

type RequestPtr[T any] interface {
	Request
	*T
}

type Handler[T Request] interface {
	Handle(ctx *ClientContext, req T) error
}

type HandlerConstructor[S any, H Handler[T], T RequestPtr[U], U any] interface {
	New(S) H
}
