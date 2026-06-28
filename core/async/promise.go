package async

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

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

func NewDeferred() *Promise {
	return &Promise{}
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
	go func() {
		for {
			runtime.Gosched()
			p.mu.Lock()
			current := p.kickGen
			if current != gen {
				gen = current
				p.mu.Unlock()
				continue
			}
			settled := p.settled
			value := p.value
			rejected := p.rejected
			done := p.done
			p.mu.Unlock()
			if !settled || rejected || done {
				return
			}
			p.driveFrom(0, value)
			return
		}
	}()
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
			f := actor.NewFuture(ctx.ActorSystem(), p.perStepTimeout)
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
