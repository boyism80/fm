package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
)

type ResponseSpawnDoorHandler struct{}

func (ResponseSpawnDoorHandler) New() *ResponseSpawnDoorHandler {
	return &ResponseSpawnDoorHandler{}
}

func (h *ResponseSpawnDoorHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *ResponseSpawnDoor) {
	if msg == nil {
		return
	}
	m := a.GetCharacter(msg.CharacterID)
	if m == nil {
		return
	}
	ch := m.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	gw := ch.GameWorld
	if !msg.Ok {
		if gw != nil {
			ch.Listener.OnMessage(ch, constant.MsgPinkText, constant.DoorNoTownPortalMessage)
		}
		return
	}
	door := ch.SpawnFieldMapDoor(msg.SkillID, msg.Return, msg.Field)
	if door == nil && gw != nil {
		gw.GetMapSystem().RemoveReturnDoor(msg.OwnerID, uint32(msg.SkillID), uint32(ch.GetMap().Wz.ReturnMapId))
	}
}
