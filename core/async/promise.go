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

type chain struct {
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

type Task struct {
	c *chain
}

type Promise[T any] struct {
	c *chain
}

func NewTask(ctx actor.Context, perStepTimeout time.Duration) *Task {
	return &Task{c: &chain{ctx: ctx, perStepTimeout: perStepTimeout, settled: true}}
}

func NewDeferredTask(ctx actor.Context) *Task {
	return &Task{c: &chain{ctx: ctx}}
}

func NewDeferred[T any](ctx actor.Context) *Promise[T] {
	return &Promise[T]{c: &chain{ctx: ctx}}
}

func (t *Task) Do(fn func() error) *Task {
	t.c.appendStep(promiseStep{fn: func(v interface{}) (interface{}, error) {
		return v, fn()
	}})
	return t
}

func (t *Task) DoAsync(fn func() error) *Task {
	t.c.appendStep(promiseStep{async: true, fn: func(v interface{}) (interface{}, error) {
		return v, fn()
	}})
	return t
}

func (t *Task) Then[U any](fn func() (U, error)) *Promise[U] {
	t.c.appendStep(promiseStep{fn: func(interface{}) (interface{}, error) {
		return fn()
	}})
	return &Promise[U]{c: t.c}
}

func (t *Task) ThenAsync[U any](fn func() (U, error)) *Promise[U] {
	t.c.appendStep(promiseStep{async: true, fn: func(interface{}) (interface{}, error) {
		return fn()
	}})
	return &Promise[U]{c: t.c}
}

func (t *Task) ThenRPC[U any](call func(context.Context) (U, error), use func(U) error) *Promise[U] {
	timeout := t.c.perStepTimeout
	next := t.ThenAsync(func() (U, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return call(ctx)
	})
	if use == nil {
		return next
	}
	return next.Do(use)
}

func (t *Task) OnError(fn func(error)) *Task {
	if fn != nil {
		t.c.onErrorSet(fn)
	}
	return t
}

func (t *Task) Finally(fn func()) *Task {
	if fn != nil {
		t.c.finallySet(fn)
	}
	return t
}

func (t *Task) Complete() {
	t.c.setResult(nil)
}

func (t *Task) SetError(err error) {
	t.c.setError(err)
}

func (p *Promise[T]) Task() *Task {
	return &Task{c: p.c}
}

func (p *Promise[T]) Then[U any](fn func(T) (U, error)) *Promise[U] {
	p.c.appendStep(promiseStep{fn: func(v interface{}) (interface{}, error) {
		value, _ := v.(T)
		return fn(value)
	}})
	return &Promise[U]{c: p.c}
}

func (p *Promise[T]) Do(fn func(T) error) *Promise[T] {
	p.c.appendStep(promiseStep{fn: func(v interface{}) (interface{}, error) {
		value, _ := v.(T)
		return v, fn(value)
	}})
	return p
}

func (p *Promise[T]) DoAsync(fn func(T) error) *Promise[T] {
	p.c.appendStep(promiseStep{async: true, fn: func(v interface{}) (interface{}, error) {
		value, _ := v.(T)
		return v, fn(value)
	}})
	return p
}

func (p *Promise[T]) ThenAsync[U any](fn func(T) (U, error)) *Promise[U] {
	p.c.appendStep(promiseStep{async: true, fn: func(v interface{}) (interface{}, error) {
		value, _ := v.(T)
		return fn(value)
	}})
	return &Promise[U]{c: p.c}
}

func (p *Promise[T]) ThenRPC[U any](call func(context.Context) (U, error), use func(U) error) *Promise[U] {
	return p.Task().ThenRPC(call, use)
}

func (p *Promise[T]) OnError(fn func(error)) *Promise[T] {
	if fn != nil {
		p.c.onErrorSet(fn)
	}
	return p
}

func (p *Promise[T]) Finally(fn func()) *Promise[T] {
	if fn != nil {
		p.c.finallySet(fn)
	}
	return p
}

func (p *Promise[T]) SetResult(value T) {
	p.c.setResult(value)
}

func (p *Promise[T]) SetError(err error) {
	p.c.setError(err)
}

func (p *Promise[T]) Completed() bool {
	p.c.mu.Lock()
	defer p.c.mu.Unlock()
	return p.c.settled || p.c.rejected
}

func (p *Promise[T]) Result() (T, error) {
	p.c.mu.Lock()
	defer p.c.mu.Unlock()
	value, _ := p.c.value.(T)
	return value, p.c.err
}

func (p *chain) appendStep(step promiseStep) {
	p.mu.Lock()
	if p.rejected {
		p.mu.Unlock()
		return
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
}

func (p *chain) scheduleKick(gen uint64) {
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

func (p *chain) kick(gen uint64) {
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

func (p *chain) onErrorSet(fn func(error)) {
	p.mu.Lock()
	if p.rejected {
		err := p.err
		p.mu.Unlock()
		fn(err)
		return
	}
	p.onError = fn
	p.mu.Unlock()
}

func (p *chain) finallySet(fn func()) {
	p.mu.Lock()
	if p.done {
		p.mu.Unlock()
		fn()
		return
	}
	p.finally = fn
	p.mu.Unlock()
}

func (p *chain) setResult(value interface{}) {
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

func (p *chain) setError(err error) {
	if err == nil {
		p.setResult(nil)
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

func (p *chain) isDone() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.done || p.rejected
}

func (p *chain) driveFrom(i int, value interface{}) {
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

func (p *chain) completeSuccess(value interface{}) {
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

func (p *chain) completeError(err error) {
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
