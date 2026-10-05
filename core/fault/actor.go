package fault

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

func (in *Injector) ReceiverMiddleware(name string) actor.ReceiverMiddleware {
	return func(next actor.ReceiverFunc) actor.ReceiverFunc {
		return func(c actor.ReceiverContext, envelope *actor.MessageEnvelope) {
			switch envelope.Message.(type) {
			case actor.SystemMessage, actor.AutoReceiveMessage:
				next(c, envelope)
				return
			}

			delay, mode := in.find(name)
			time.Sleep(delay)
			if mode != None {
				return
			}
			next(c, envelope)
		}
	}
}
