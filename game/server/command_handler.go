package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/game/client"
)

type CommandHandler struct {
	gs       *GameServer
	handlers map[string]Command
}

func NewCommandHandler(gs *GameServer) *CommandHandler {
	return &CommandHandler{
		gs:       gs,
		handlers: make(map[string]Command),
	}
}

func Bind[C CommandHandlerConstructor[H], H Command](ch *CommandHandler) {
	var constructor C
	handler := constructor.New(ch.gs)

	commandName := handler.GetCommandName()
	ch.handlers[commandName] = handler
	log.Printf("Registered command handler: %s", commandName)
}

func (ch *CommandHandler) Handle(gameClient *client.GameClient, params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("no command specified")
	}

	command := params[0]
	args := params[1:]

	handler, ok := ch.handlers[command]
	if !ok {
		return fmt.Errorf("unknown command: %s", command)
	}

	return handler.Handle(gameClient, args...)
}

func (ch *CommandHandler) GetHandlers() map[string]Command {
	return ch.handlers
}

func (gs *GameServer) registerCommandHandlers() {
	Bind[*ListCommands](gs.commandHandler)
	Bind[*CreateItem](gs.commandHandler)
	Bind[*ClearMeso](gs.commandHandler)
	Bind[*GainMeso](gs.commandHandler)
	Bind[*FullMeso](gs.commandHandler)
	Bind[*ChangeMap](gs.commandHandler)
	Bind[*GetPosition](gs.commandHandler)
	Bind[*ChangeHp](gs.commandHandler)
	Bind[*ChangeMp](gs.commandHandler)
	Bind[*ChangeStr](gs.commandHandler)
	Bind[*ChangeDex](gs.commandHandler)
	Bind[*ChangeInt](gs.commandHandler)
	Bind[*ChangeLuk](gs.commandHandler)
	Bind[*ChangeAllStats](gs.commandHandler)
	Bind[*ChangeLevel](gs.commandHandler)
	Bind[*Invincible](gs.commandHandler)
	Bind[*ChangeClass](gs.commandHandler)
	Bind[*SpawnMob](gs.commandHandler)
	Bind[*KillAllMobs](gs.commandHandler)
	Bind[*MasterSkills](gs.commandHandler)
	Bind[*SpawnNpc](gs.commandHandler)
}
