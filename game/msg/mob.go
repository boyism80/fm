package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/protocol/req"
)

type MobKill struct {
	AnimationType constant.MobDieAnimationType
}

type MobControllerChange struct {
	Controller *actor.PID
}

type MobMove struct {
	req.MoveMob
	Sender *actor.PID
}
