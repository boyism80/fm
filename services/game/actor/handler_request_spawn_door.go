package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/entity"
)

type RequestSpawnDoorHandler struct{}

func (RequestSpawnDoorHandler) New() *RequestSpawnDoorHandler {
	return &RequestSpawnDoorHandler{}
}

func (h *RequestSpawnDoorHandler) Handle(ctx actor.Context, a *MapActor, msg *RequestSpawnDoor) {
	if msg == nil || msg.ReplyTo == nil {
		return
	}
	var targetMap *entity.Map
	for _, m := range a.Maps() {
		if m != nil && m.GetMapID() == msg.TargetMapID {
			targetMap = m
			break
		}
	}
	if targetMap == nil || targetMap.Wz == nil {
		return
	}
	portalID, townPos, ok := targetMap.TryAcquireMysticReturnPortal(msg.PartyOwnerSlot)
	if !ok {
		ctx.Send(msg.ReplyTo, &ResponseSpawnDoor{
			Ok:          false,
			CharacterID: msg.CharacterID,
			OwnerID:     msg.OwnerID,
			SkillID:     msg.SkillID,
			Field:       msg.Field,
			PartyID:     msg.PartyID,
		})
		return
	}
	committed := false
	defer func() {
		if !committed {
			targetMap.ReleaseMysticReturnPortal(portalID)
		}
	}()
	wz := targetMap.Wz
	returnEp := entity.DoorEndpoint{
		MapID:    uint32(wz.ID),
		PortalID: portalID,
		Position: townPos,
	}
	door := entity.NewDoor(
		msg.OwnerID,
		msg.SkillID,
		msg.Field,
		returnEp,
		msg.PartyID,
	)
	targetMap.AddDoor(door)
	committed = true
	ctx.Send(msg.ReplyTo, &ResponseSpawnDoor{
		Ok:          true,
		CharacterID: msg.CharacterID,
		OwnerID:     msg.OwnerID,
		SkillID:     msg.SkillID,
		Return:      returnEp,
		Field:       msg.Field,
		PartyID:     msg.PartyID,
	})
}
