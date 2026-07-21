package async

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

func Ask(ctx actor.Context, target *actor.PID, timeout time.Duration, send func(replyTo *actor.PID)) *Promise {
	p := NewDeferred()
	if ctx == nil || target == nil || send == nil {
		p.SetError(fmt.Errorf("async.Ask: invalid args"))
		return p
	}
	if ctx.ActorSystem() == nil {
		p.SetError(fmt.Errorf("async.Ask: actor system is nil"))
		return p
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	p.ctx = ctx
	p.perStepTimeout = timeout

	f := actor.NewFuture(ctx.ActorSystem(), timeout)
	send(f.PID())
	ctx.ReenterAfter(f, func(res interface{}, err error) {
		if err != nil {
			p.SetError(err)
			return
		}
		p.SetResult(res)
	})
	return p
}
