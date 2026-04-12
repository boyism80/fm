package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	g_actor "github.com/boyism80/fm/game/actor"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/protocol/request"
)

type LoginGame struct {
	gs     *GameServer
	opcode byte
}

func (LoginGame) New(gs *GameServer) *LoginGame {
	return &LoginGame{
		gs:     gs,
		opcode: 0x06,
	}
}

func (h *LoginGame) GetOpcode() byte {
	return h.opcode
}

func (h *LoginGame) Handle(ctx *core.ClientContext, req *request.LoginGame) error {
	name := "채승현"
	if req.PlayerId != 1 {
		name = "채진영"
	}

	character := entity.NewDummyCharacter(ctx.Client, h.gs.characterListener, req.PlayerId, name, h.gs)

	if req.PlayerId == 1 {
		character.Role = constant.RoleAdmin
		character.Invincible = true
	}

	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}
	client.SetCharacter(character)

	initialMapID, ok := h.gs.resources.NameToMap("헤네시스")
	if !ok {
		return fmt.Errorf("initial map name not found")
	}
	initialSpawnPoint := uint8(1)

	mapInstance := h.gs.GetMap(initialMapID)
	if mapInstance == nil {
		log.Printf("Initial map %d not found", initialMapID)
		return fmt.Errorf("initial map %d not found", initialMapID)
	}

	wz := mapInstance.Wz
	if wz == nil {
		log.Printf("Wz not found for map %d", initialMapID)
		return fmt.Errorf("wz not found for map %d", initialMapID)
	}

	if _, ok := wz.Portals[initialSpawnPoint]; !ok {
		log.Printf("Portal %d not found in map %d", initialSpawnPoint, initialMapID)
		return fmt.Errorf("portal %d not found in map %d", initialSpawnPoint, initialMapID)
	}

	character.Stance = constant.StanceDefaultValue

	targetMapPID := mapInstance.GetActorPID()
	if targetMapPID == nil {
		return fmt.Errorf("MapActor PID not found for map %d", initialMapID)
	}

	rootContext := h.gs.GetServer().GetRootContext()
	if rootContext == nil {
		return fmt.Errorf("rootContext not set")
	}

	rootContext.Send(targetMapPID, &g_actor.AddCharacter{
		Character:  character,
		SpawnPoint: initialSpawnPoint,
		Init:       true,
	})

	return nil
}
