package async

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

const promiseKickTimeout = 5 * time.Second

type promiseResult struct {
	Value interface{}
	Err   error
}

type promiseStep struct {
	async bool
	fn    func(interface{}) (interface{}, error)
}

type Promise struct {
	ctx            actor.Context
	perStepTimeout time.Duration
	steps          []promiseStep
	onError        func(error)
	finally        func()

	mu       sync.Mutex
	settled  bool
	rejected bool
	done     bool
	started  bool
	value    interface{}
	err      error
	kickGen  uint64
}

func NewPromise(ctx actor.Context, perStepTimeout time.Duration) *Promise {
	return &Promise{
		ctx:            ctx,
		perStepTimeout: perStepTimeout,
		settled:        true,
		value:          nil,
	}
}

func NewDeferred(ctx actor.Context) *Promise {
	return &Promise{ctx: ctx}
}

func (p *Promise) Then(fn func(interface{}) (interface{}, error)) *Promise {
	return p.appendStep(promiseStep{async: false, fn: fn})
}

func (p *Promise) ThenAsync(fn func(interface{}) (interface{}, error)) *Promise {
	return p.appendStep(promiseStep{async: true, fn: fn})
}

func (p *Promise) appendStep(step promiseStep) *Promise {
	if step.fn == nil {
		return p
	}
	p.mu.Lock()
	if p.rejected {
		p.mu.Unlock()
		return p
	}
	p.steps = append(p.steps, step)
	var gen uint64
	if p.settled {
		p.kickGen++
		gen = p.kickGen
	}
	p.mu.Unlock()
	if gen != 0 {
		p.scheduleKick(gen)
	}
	return p
}

func (p *Promise) scheduleKick(gen uint64) {
	if p.ctx == nil {
		go func() {
			runtime.Gosched()
			p.kick(gen)
		}()
		return
	}
	f := actor.NewFuture(p.ctx.ActorSystem(), promiseKickTimeout)
	p.ctx.ReenterAfter(f, func(interface{}, error) {
		p.kick(gen)
	})
	p.ctx.Send(f.PID(), &promiseResult{})
}

func (p *Promise) kick(gen uint64) {
	p.mu.Lock()
	// A newer kick owns the chain; running this one too would execute every step twice.
	if p.kickGen != gen || p.started || !p.settled || p.rejected || p.done {
		p.mu.Unlock()
		return
	}
	p.started = true
	value := p.value
	p.mu.Unlock()
	p.driveFrom(0, value)
}

func ThenRPC[T any](p *Promise, call func(context.Context) (T, error), use func(T) error) *Promise {
	p.ThenAsync(func(_ interface{}) (interface{}, error) {
		ctx, cancel := context.WithTimeout(context.Background(), p.perStepTimeout)
		defer cancel()
		return call(ctx)
	})
	return p.Then(func(v interface{}) (interface{}, error) {
		return v, use(v.(T))
	})
}

func (p *Promise) OnError(fn func(error)) *Promise {
	if fn == nil {
		return p
	}
	p.mu.Lock()
	if p.rejected {
		err := p.err
		p.mu.Unlock()
		fn(err)
		return p
	}
	p.onError = fn
	p.mu.Unlock()
	return p
}

func (p *Promise) Finally(fn func()) *Promise {
	if fn == nil {
		return p
	}
	p.mu.Lock()
	if p.done {
		p.mu.Unlock()
		fn()
		return p
	}
	p.finally = fn
	p.mu.Unlock()
	return p
}

func (p *Promise) SetResult(value interface{}) {
	p.mu.Lock()
	if p.settled || p.rejected || p.done {
		p.mu.Unlock()
		return
	}
	p.settled = true
	p.value = value
	hasSteps := len(p.steps) > 0
	var gen uint64
	if hasSteps {
		p.kickGen++
		gen = p.kickGen
	}
	p.mu.Unlock()
	if hasSteps {
		p.scheduleKick(gen)
	}
}

func (p *Promise) SetError(err error) {
	if err == nil {
		p.SetResult(nil)
		return
	}
	p.mu.Lock()
	if p.settled || p.rejected || p.done {
		p.mu.Unlock()
		return
	}
	p.rejected = true
	p.err = err
	onError := p.onError
	finally := p.finally
	p.done = true
	p.mu.Unlock()
	if onError != nil {
		onError(err)
	}
	if finally != nil {
		finally()
	}
}

func (p *Promise) Completed() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.settled || p.rejected
}

func (p *Promise) Result() (interface{}, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.value, p.err
}

func (p *Promise) isDone() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.done || p.rejected
}

func (p *Promise) driveFrom(i int, value interface{}) {
	if p.isDone() {
		return
	}
	p.mu.Lock()
	if i >= len(p.steps) {
		p.mu.Unlock()
		p.completeSuccess(value)
		return
	}
	step := p.steps[i]
	ctx := p.ctx
	p.mu.Unlock()

	if step.async {
		if ctx != nil {
			// Give the Future slightly more time than the gRPC context so the async
			// step can complete and deliver its result before the Future times out.
			const futureGrace = 5 * time.Second
			f := actor.NewFuture(ctx.ActorSystem(), p.perStepTimeout+futureGrace)
			go func() {
				v, err := step.fn(value)
				ctx.Send(f.PID(), &promiseResult{Value: v, Err: err})
			}()
			ctx.ReenterAfter(f, func(res interface{}, ferr error) {
				if p.isDone() {
					return
				}
				if ferr != nil {
					p.completeError(ferr)
					return
				}
				pr, ok := res.(*promiseResult)
				if !ok {
					p.completeError(fmt.Errorf("async.Promise: unexpected result type %T", res))
					return
				}
				if pr.Err != nil {
					p.completeError(pr.Err)
					return
				}
				p.driveFrom(i+1, pr.Value)
			})
			return
		}
		go func() {
			v, err := step.fn(value)
			if err != nil {
				p.completeError(err)
				return
			}
			p.driveFrom(i+1, v)
		}()
		return
	}

	v, err := step.fn(value)
	if err != nil {
		p.completeError(err)
		return
	}
	p.driveFrom(i+1, v)
}

func (p *Promise) completeSuccess(value interface{}) {
	p.mu.Lock()
	if p.done || p.rejected {
		p.mu.Unlock()
		return
	}
	p.done = true
	p.settled = true
	p.value = value
	finally := p.finally
	p.mu.Unlock()
	if finally != nil {
		finally()
	}
}

func (p *Promise) completeError(err error) {
	p.mu.Lock()
	if p.done || p.rejected {
		p.mu.Unlock()
		return
	}
	p.rejected = true
	p.done = true
	p.err = err
	onError := p.onError
	finally := p.finally
	p.mu.Unlock()
	if onError != nil {
		onError(err)
	}
	if finally != nil {
		finally()
	}
}
