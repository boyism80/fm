package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/protocol"
	"github.com/boyism80/fm/game/protocol/req"
)

type MobKill struct {
	AnimationType constant.MobDieAnimationType
}

type MobControllerChange struct {
	Before *actor.PID
	After  *actor.PID
}

type MobMove struct {
	req.MoveMob
	Sender *actor.PID
}

type MobDamaged struct {
	Sender      *actor.PID
	CharacterId uint32
	DamagePairs []protocol.DamagePair
}
