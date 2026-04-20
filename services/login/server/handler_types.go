package server

import (
	"github.com/boyism80/fm/core"
)

type HandlerConstructor[H core.Handler[T], T core.RequestPtr[U], U any] interface {
	New(*LoginServer) H
}
