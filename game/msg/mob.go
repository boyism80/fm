package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/constant"
)

type MobKill struct {
	AnimationType constant.MobDieAnimationType
}

type MobControllerChange struct {
	Controller *actor.PID
}
