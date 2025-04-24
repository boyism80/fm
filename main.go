package main

import (
	console "github.com/asynkron/goconsole"
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/actor"
	"github.com/boyism80/fm/model"
	"github.com/boyism80/fm/msg"
)

func main() {
	system := protoactor.NewActorSystem()
	props := protoactor.PropsFromProducer(func() protoactor.Actor { return actor.NewCharacterActor(&model.Character{}) })

	pid := system.Root.Spawn(props)
	system.Root.Send(pid, &msg.ObjectMove{X: 1, Y: 1})
	system.Root.Send(pid, &msg.LifeAddHp{Hp: 100})
	_, _ = console.ReadLine()
}
