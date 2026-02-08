package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

// NpcClick handles NPC click packet requests
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

	mapInstance := h.gs.GetMap(character.GetMap())
	if mapInstance == nil {
		log.Printf("Map %d not found", character.GetMap())
		return fmt.Errorf("map %d not found", character.GetMap())
	}

	npcs := mapInstance.GetNpcs()
	npcInterface, exists := npcs[req.OID]
	if !exists {
		log.Printf("NPC %d not found on map %d", req.OID, character.GetMap())
		return fmt.Errorf("npc %d not found", req.OID)
	}

	// Type assert to entity.Npc
	npc, ok := npcInterface.(*entity.Npc)
	if !ok {
		log.Printf("Invalid NPC type for OID %d", req.OID)
		return fmt.Errorf("invalid NPC type")
	}

	// Check if NPC has a shop
	resources := h.gs.GetResources()
	if resources == nil {
		log.Printf("Resources not available")
		return fmt.Errorf("resources not available")
	}

	npcID := npc.Wz.ID
	shop := resources.GetShop(npcID)

	if shop != nil {
		// NPC has a shop, send OPEN_NPC_SHOP packet
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

	// NPC doesn't have a shop, execute script
	pid := ctx.LogicActorID
	rootState := luax.GetRootState(pid)
	if rootState == nil {
		log.Printf("No lua root state for actor %s", pid)
		return fmt.Errorf("lua state not available")
	}
	scriptPath := fmt.Sprintf("script/npc/%d.lua", npcID)
	luaThread, err := luax.NewThread(rootState, scriptPath)
	if err != nil {
		log.Printf("Failed to create NPC script thread: %v", err)
		return err
	}
	if err := h.gs.ExecuteScript(rootState, luaThread, "on_start", character); err != nil {
		log.Printf("Failed to execute NPC script: %v", err)
		return err
	}

	return nil
}
