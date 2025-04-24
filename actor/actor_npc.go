package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/model"
)

type NpcActor struct {
	Npc     *model.NPC
	Context actor.Context
}
