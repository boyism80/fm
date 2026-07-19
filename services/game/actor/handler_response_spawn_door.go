package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
)

type ResponseSpawnDoorHandler struct{}

func (ResponseSpawnDoorHandler) New() *ResponseSpawnDoorHandler {
	return &ResponseSpawnDoorHandler{}
}

func (h *ResponseSpawnDoorHandler) Handle(ctx actor.Context, a *MapActor, msg *ResponseSpawnDoor) {
	if msg == nil || a.Map == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		if msg.Ok && a.Map.Wz != nil && a.Map.GameWorld != nil {
			a.Map.GameWorld.GetMapSystem().RemoveReturnDoor(msg.OwnerID, uint32(msg.SkillID), uint32(a.Map.Wz.ReturnMapId))
		}
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
