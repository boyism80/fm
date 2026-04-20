package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
)

type NpcClick struct {
	gs     *GameServer
	opcode byte
}

func (NpcClick) New(gs *GameServer) *NpcClick {
	return &NpcClick{
		gs:     gs,
		opcode: 0x29,
	}
}

func (h *NpcClick) GetOpcode() byte {
	return h.opcode
}

func (h *NpcClick) Handle(ctx *core.ClientContext, req *request.NpcClick) error {
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
		log.Printf("Map not found")
		return fmt.Errorf("map not found")
	}

	npcs := mapInstance.GetNpcs()
	npcInterface, exists := npcs[req.OID]
	if !exists {
		log.Printf("NPC %d not found on map", req.OID)
		return fmt.Errorf("npc %d not found", req.OID)
	}

	npc, ok := npcInterface.(*entity.Npc)
	if !ok {
		log.Printf("Invalid NPC type for OID %d", req.OID)
		return fmt.Errorf("invalid NPC type")
	}

	resources := h.gs.GetResources()
	if resources == nil {
		log.Printf("Resources not available")
		return fmt.Errorf("resources not available")
	}

	npcID := npc.Wz.ID
	shop := resources.GetShop(npcID)

	if shop != nil {

		character.CurrentShopID = npcID
		packet := &response.OpenNpcShop{
			ShopID: int32(npcID),
			Shop:   shop,
			Items:  resources.Items,
		}
		if err := character.Send(packet, types.SEND_POLICY_ENCRYPT); err != nil {
			log.Printf("Failed to send OPEN_NPC_SHOP packet: %v", err)
			return err
		}
		return nil
	}

	root := mapInstance.GetLuaRoot()
	if root == nil {
		log.Printf("No lua root state for map %d", mapInstance.GetMapID())
		return fmt.Errorf("lua state not available")
	}
	scriptPath := fmt.Sprintf("script/npc/%d.lua", npcID)
	luaThread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		log.Printf("Failed to create NPC script thread: %v", err)
		return err
	}
	luax.SetConfiguration(luaThread, luax.Configuration{
		ActorContext: ctx.ActorContext,
	})
	_, err = luax.CallThread(luaThread, "on_start", character)
	if err != nil {
		log.Printf("Failed to execute NPC script: %v", err)
		return err
	}
	return nil
}
