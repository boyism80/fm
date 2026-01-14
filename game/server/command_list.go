package server

import (
	"fmt"
	"sort"

	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
)

type ListCommands struct {
	gs *GameServer
}

func (*ListCommands) New(gs *GameServer) *ListCommands {
	return &ListCommands{
		gs: gs,
	}
}

func (h *ListCommands) GetCommandName() string {
	return "명령어"
}

func (h *ListCommands) GetUsage() string {
	return "- 사용 가능한 명령어 목록 표시"
}

func (h *ListCommands) Handle(gameClient *client.GameClient, args ...string) error {
	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	type commandInfo struct {
		name  string
		usage string
	}

	var commands []commandInfo
	handlers := h.gs.commandHandler.GetHandlers()
	for cmdName, cmd := range handlers {
		if cmdName != "명령어" {
			commands = append(commands, commandInfo{
				name:  cmdName,
				usage: cmd.GetUsage(),
			})
		}
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].name < commands[j].name
	})

	character.Listener.OnMessage(constant.MSG_LIGHT_BLUE_TEXT, "=== 사용 가능한 명령어 목록 ===")
	for i, cmd := range commands {
		line := fmt.Sprintf("%d. /%s %s", i+1, cmd.name, cmd.usage)
		character.Listener.OnMessage(constant.MSG_LIGHT_BLUE_TEXT, line)
	}
	character.Listener.OnMessage(constant.MSG_LIGHT_BLUE_TEXT, fmt.Sprintf("총 %d개의 명령어가 있습니다.", len(commands)))

	return nil
}
