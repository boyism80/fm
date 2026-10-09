package server

import (
	"fmt"
	"log"
	"math"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	lua "github.com/yuin/gopher-lua"
)

type DistributeAP struct {
	gs *GameServer
}

func (DistributeAP) New(gs *GameServer) *DistributeAP {
	return &DistributeAP{
		gs: gs,
	}
}

func (h *DistributeAP) Handle(ctx *core.ClientContext, req *request.DistributeAP) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}

	stat := constant.StatType(req.StatType)
	assigned := false
	switch stat {
	case constant.StatTypeHP:
		assigned = character.Points.AssignHPMP(stat, h.apIncrease(character, "get_ap_to_hp", 10))
	case constant.StatTypeMP:
		assigned = character.Points.AssignHPMP(stat, h.apIncrease(character, "get_ap_to_mp", 5))
	default:
		assigned = character.Points.AssignStats([]entity.APEntry{{Stat: stat, Amount: 1}})
	}
	if !assigned {
		character.Listener.OnUpdateStats(character, map[constant.Stat]int32{}, true)
	}
	return nil
}

func (h *DistributeAP) apIncrease(character *entity.Character, funcName string, fallback uint32) uint32 {
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return fallback
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return fallback
	}
	thread, err := luax.NewThread(root, constant.CharacterQueryScriptPath)
	if err != nil {
		return fallback
	}

	value, err := luax.Call(thread, funcName, character)
	if err != nil {
		log.Printf("DistributeAP: %s character=%d: %v", funcName, character.GetID(), err)
		return fallback
	}
	n, ok := value.(lua.LNumber)
	if !ok || n <= 0 {
		return fallback
	}
	if float64(n) >= math.MaxUint32 {
		return math.MaxUint32
	}
	return uint32(n)
}
