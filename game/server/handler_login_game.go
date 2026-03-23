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

// LoginGame handles game login packet requests
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

	character := entity.NewDummyCharacter(ctx.Client, nil, req.PlayerId, name, h.gs)
	character.Listener = NewGameCharacterListener(h.gs, character)

	// Set GM mode (for testing: player ID 1 is GM)
	if req.PlayerId == 1 {
		character.Role = constant.RoleAdmin
		character.Invincible = true
	}

	// Character 생성 시점에는 map이 결정되지 않음 (AddPlayer에서 SetMap 호출됨)
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}
	client.SetCharacter(character)

	// 초기 맵으로 이동 (nil MapActor에서 실제 맵으로)
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

	character.Stance = 0

	// initialMapID에 대응되는 MapActor의 PID를 구해서 WarpCharacter 메시지 전송
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
