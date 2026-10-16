package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type RequestSpawnDoorHandler struct{}

func (RequestSpawnDoorHandler) New() *RequestSpawnDoorHandler {
	return &RequestSpawnDoorHandler{}
}

func (h *RequestSpawnDoorHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *game.RequestSpawnDoor) {
	if msg == nil || msg.ReplyTo == nil {
		return
	}
	var targetMap *game.Map
	for _, m := range a.Maps() {
		if m == msg.TargetMap {
			targetMap = m
			break
		}
	}
	if targetMap == nil || targetMap.Wz() == nil {
		return
	}
	portalID, townPos, ok := targetMap.FindMysticReturnPortal(msg.PartyOwnerSlot)
	if !ok {
		ctx.Send(msg.ReplyTo, &game.ResponseSpawnDoor{
			Ok:          false,
			Key:         msg.Key,
			CharacterID: msg.CharacterID,
			OwnerID:     msg.OwnerID,
			SkillID:     msg.SkillID,
			Field:       msg.Field,
			PartyID:     msg.PartyID,
		})
		return
	}
	returnEp := game.DoorEndpoint{
		Map:      targetMap,
		PortalID: portalID,
		Position: townPos,
	}
	door := game.NewDoor(
		msg.Key,
		msg.OwnerID,
		msg.SkillID,
		msg.Field,
		returnEp,
		msg.PartyID,
	)
	targetMap.AddDoor(door)
	ctx.Send(msg.ReplyTo, &game.ResponseSpawnDoor{
		Ok:          true,
		Key:         msg.Key,
		CharacterID: msg.CharacterID,
		OwnerID:     msg.OwnerID,
		SkillID:     msg.SkillID,
		Return:      returnEp,
		Field:       msg.Field,
		PartyID:     msg.PartyID,
	})
}
