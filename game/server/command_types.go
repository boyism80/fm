package server

import (
	"github.com/boyism80/fm/game/client"
)

type Command interface {
	Handle(gameClient *client.GameClient, args ...string) error
	GetCommandName() string
	GetUsage() string
}

type CommandHandlerConstructor[H Command] interface {
	New(*GameServer) H
}
