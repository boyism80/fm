package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/domain"
)

type ResponseSpawnDoorHandler struct{}

func (ResponseSpawnDoorHandler) New() *ResponseSpawnDoorHandler {
	return &ResponseSpawnDoorHandler{}
}

func (h *ResponseSpawnDoorHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *domain.ResponseSpawnDoor) {
	if msg == nil {
		return
	}
	var ch *domain.Character
	if m := a.GetCharacter(msg.CharacterID); m != nil {
		ch = m.GetPlayer(msg.CharacterID)
	}
	if !msg.Ok {
		if ch != nil {
			ch.Listener.OnMessage(ch, constant.MsgPinkText, constant.DoorNoTownPortalMessage)
		}
		return
	}

	if ch == nil || ch.Doors.SpawnField(msg.Key, msg.SkillID, msg.Return, msg.Field) == nil {
		msg.Return.Map.GameWorld.Systems().Get[*domain.MapSystem]().DespawnDoor(msg.Return.Map, msg.Key, true, false)
	}
}
