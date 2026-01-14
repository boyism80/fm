package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/request"
	lua "github.com/yuin/gopher-lua"
)

// Dialog handles dialog packet requests
type Dialog struct {
	gs     *GameServer
	opcode byte
}

func (Dialog) New(gs *GameServer) *Dialog {
	return &Dialog{
		gs:     gs,
		opcode: 0x2B,
	}
}

func (h *Dialog) GetOpcode() byte {
	return h.opcode
}

func (h *Dialog) Handle(ctx *core.ClientContext, req *request.Dialog) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character not found for client")
		return fmt.Errorf("character not found")
	}

	dialog := character.GetCurrentDialog()
	if dialog == nil {
		log.Printf("No active dialog for character %d", character.GetID())
		return fmt.Errorf("no active dialog")
	}

	// Get thread-local LuaState
	// This will be shared among all MapActors running on the same thread
	luaState := luax.GetThreadLocalState()
	if luaState == nil {
		return fmt.Errorf("lua state not available")
	}

	var args []lua.LValue
	switch req.DialogType {
	case constant.DIALOG_TYPE_DEFAULT:
		args = append(args, lua.LBool(req.Next))
	case constant.DIALOG_TYPE_YES_NO:
		args = append(args, lua.LBool(req.Next))
	case constant.DIALOG_TYPE_LIST:
		if req.Next {
			args = append(args, lua.LNumber(req.Selected))
		} else {
			args = append(args, lua.LNil)
		}
	case constant.DIALOG_TYPE_INPUT:
		if req.Next {
			args = append(args, lua.LString(req.Text))
		} else {
			args = append(args, lua.LNil)
		}
	case constant.DIALOG_TYPE_ACCEPT_ESCAPE:
	case constant.DIALOG_TYPE_ACCEPT:
		args = append(args, lua.LBool(req.Next))
	}

	resumeState, err, _ := luaState.Resume(dialog, nil, args...)
	if err != nil {
		log.Printf("Failed to resume dialog: %v", err)
		character.ClearCurrentDialog()
		return fmt.Errorf("failed to resume dialog: %w", err)
	}

	switch resumeState {
	case lua.ResumeOK:
		character.ClearCurrentDialog()
	case lua.ResumeYield:
	case lua.ResumeError:
		log.Printf("Dialog error for character %d: %v", character.GetID(), err)
		character.ClearCurrentDialog()
		return fmt.Errorf("dialog error: %w", err)
	}

	return nil
}
