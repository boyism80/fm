package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/services/game"
)

type ResponseSpawnDoorHandler struct{}

func (ResponseSpawnDoorHandler) New() *ResponseSpawnDoorHandler {
	return &ResponseSpawnDoorHandler{}
}

func (h *ResponseSpawnDoorHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *game.ResponseSpawnDoor) {
	if msg == nil {
		return
	}
	var ch *game.Character
	if m := a.GetCharacter(msg.CharacterID); m != nil {
		ch = m.GetPlayer(msg.CharacterID)
	}
	if !msg.Ok {
		if ch != nil {
			ch.Listener().OnMessage(ch, constant.MsgPinkText, constant.DoorNoTownPortalMessage)
		}
		return
	}

	if ch == nil || ch.Doors().SpawnField(msg.Key, msg.SkillID, msg.Return, msg.Field) == nil {
		msg.Return.Map.GameWorld().Systems().Get[*game.MapSystem]().DespawnDoor(msg.Return.Map, msg.Key, true, false)
	}
}
