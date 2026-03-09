package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/client"
)

type Script struct {
	gs *GameServer
}

func (*Script) New(gs *GameServer) *Script {
	return &Script{
		gs: gs,
	}
}

func (h *Script) GetCommandName() string {
	return "스크립트"
}

func (h *Script) GetUsage() string {
	return "- script/script.lua의 on_script(me) 실행"
}

func (h *Script) Handle(gameClient *client.GameClient, args ...string) error {
	return fmt.Errorf("client context is required")
}

func (h *Script) HandleWithContext(ctx *core.ClientContext, gameClient *client.GameClient, args ...string) error {
	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	logicActorPID := ctx.LogicActorPID
	if logicActorPID == nil {
		return fmt.Errorf("logic actor pid not found")
	}

	root := luax.GetRootLuaState(logicActorPID.String())
	if root == nil {
		return fmt.Errorf("lua root state not found")
	}

	_, thread, err := luax.Call(root, "script/script.lua", "on_script", character)
	if thread != nil {
		thread.Close()
	}
	if err != nil {
		return fmt.Errorf("failed to execute on_script: %w", err)
	}

	return nil
}
