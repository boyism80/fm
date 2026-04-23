package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	lua "github.com/yuin/gopher-lua"
)

type Dialog struct {
	gs *GameServer
}

func (Dialog) New(gs *GameServer) *Dialog {
	return &Dialog{
		gs: gs,
	}
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

	mapInstance := character.GetMap()
	if mapInstance == nil {
		return fmt.Errorf("lua state not available")
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return fmt.Errorf("lua state not available")
	}
	thread := character.GetDialog()
	if thread == nil {
		log.Printf("No active dialog for character %d", character.GetID())
		return fmt.Errorf("no active dialog")
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

	luax.SetConfiguration(thread, luax.Configuration{
		ActorContext: ctx.ActorContext,
	})
	resumeState, err := luax.Resume(root, thread, args...)
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
		log.Printf("Dialog error for character %d", character.GetID())
		character.ClearCurrentDialog()
		return fmt.Errorf("dialog error")
	}

	return nil
}
