package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
)

type Command interface {
	Handle(gameClient *client.GameClient, args ...string) error
	GetCommandName() string
	GetUsage() string
}

type ContextCommand interface {
	HandleWithContext(ctx *core.ClientContext, gameClient *client.GameClient, args ...string) error
}

type CommandHandlerConstructor[H Command] interface {
	New(*GameServer) H
}
