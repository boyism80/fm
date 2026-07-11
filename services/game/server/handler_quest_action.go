package server

import (
	"errors"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type QuestAction struct {
	gs *GameServer
}

func (QuestAction) New(gs *GameServer) *QuestAction {
	return &QuestAction{
		gs: gs,
	}
}

func (h *QuestAction) Handle(ctx *core.ClientContext, req *request.QuestAction) error {
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

	if character.Quests == nil {
		log.Printf("Quest container is nil for character %d", character.GetID())
		return fmt.Errorf("quest container is nil")
	}

	resources := character.GameWorld.GetResources()
	if resources == nil {
		return fmt.Errorf("resources unavailable")
	}

	questDef := resources.GetQuest(uint32(req.QuestID))
	if questDef == nil {
		log.Printf("Quest definition not found: %d", req.QuestID)
		character.Listener.OnUnlockAction(character)
		return nil
	}

	quest := character.Quests.Get(questDef.ID)

	switch req.Mode {
	case request.QuestActionRestoreLostItem:
		if !quest.CanRestoreLostItem(character, req.ItemID) {
			character.Listener.OnUnlockAction(character)
			return nil
		}
		if err := quest.RestoreLostItem(character, req.ItemID); err != nil {
			character.Listener.OnUnlockAction(character)
			return nil
		}
		character.Listener.OnUnlockAction(character)

	case request.QuestActionStart:
		if questDef.HasStartScript() {
			character.Listener.OnUnlockAction(character)
			return nil
		}
		npc := req.NPCID
		if _, err := character.Quests.Start(uint32(req.QuestID), entity.QuestPhaseOpts{NpcID: &npc}); err != nil {
			character.Listener.OnUnlockAction(character)
			return nil
		}

	case request.QuestActionComplete:
		if questDef.HasEndScript() {
			character.Listener.OnUnlockAction(character)
			return nil
		}
		npc := req.NPCID
		if err := quest.Complete(character, entity.QuestPhaseOpts{NpcID: &npc, Selection: req.Selection}); err != nil {
			character.Listener.OnUnlockAction(character)
			return nil
		}

	case request.QuestActionForfeit:
		if err := quest.Forfeit(character); err != nil {
			if errors.Is(err, entity.ErrQuestNotForfeitable) {
				character.Listener.OnMessage(character, constant.MsgPopup, "You may not forfeit this quest.")
			}
			character.Listener.OnUnlockAction(character)
			return nil
		}

	case request.QuestActionScriptedStart:
		if character.GetDialog() != nil {
			character.Listener.OnUnlockAction(character)
			return nil
		}
		if err := character.Quests.CanStart(uint32(req.QuestID), entity.QuestPhaseOpts{
			NpcID: nil,
		}); err != nil {
			character.Listener.OnUnlockAction(character)
			return nil
		}
		if err := character.RunQuestScript(ctx.ActorContext, questDef.ID, req.NPCID, "on_start"); err != nil {
			log.Printf("scripted quest start %d: %v", questDef.ID, err)
			character.Listener.OnUnlockAction(character)
		}

	case request.QuestActionScriptedEnd:
		if character.GetDialog() != nil {
			character.Listener.OnUnlockAction(character)
			return nil
		}
		if quest == nil || quest.CanComplete(character, entity.QuestPhaseOpts{}) != nil {
			character.Listener.OnUnlockAction(character)
			return nil
		}
		if err := character.RunQuestScript(ctx.ActorContext, questDef.ID, req.NPCID, "on_end"); err != nil {
			log.Printf("scripted quest end %d: %v", questDef.ID, err)
			character.Listener.OnUnlockAction(character)
		}

	default:
		character.Listener.OnUnlockAction(character)
	}

	return nil
}
