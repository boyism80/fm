package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type AddCharacterHandler struct{}

func (AddCharacterHandler) New() *AddCharacterHandler {
	return &AddCharacterHandler{}
}

func (h *AddCharacterHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *game.AddCharacter) {
	if msg == nil || msg.Character == nil || msg.TargetMap == nil {
		return
	}
	msg.TargetMap.AddPlayer(ctx, msg.Character.GetID(), msg.Character, msg.SpawnPoint, msg.Init)
	msg.Character.FinishMove(msg.TargetMap)
	if msg.Character.LoggedOut() {
		_ = msg.TargetMap.LogoutPlayer(ctx, msg.Character.GetID())
		return
	}
	msg.Character.ResumeTimers(ctx.Self())
	msg.Character.Wedding().NotifySpouseMap(ctx, msg.Init)
	if !msg.Init {
		a.GameWorld.Systems().Get[*game.PartySystem]().UpdateMemberAsync(nil, msg.Character)
		return
	}
	msg.Character.NotifyLevelClassChanged()
	msg.Character.SocialLogin(msg.Notes)
	msg.Character.Buddies().SyncLogin()
	msg.Character.Duey().CheckArrivals(ctx)
	msg.Character.Unsettled().Settle(ctx)
	msg.Character.Unsettled().ScheduleSettle()
}
