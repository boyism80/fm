package async

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type promiseResult struct {
	Value interface{}
	Err   error
}

type promiseStep struct {
	Do  func() (interface{}, error)
	Use func(interface{}) error
}

type Promise struct {
	ctx            actor.Context
	perStepTimeout time.Duration
	steps          []promiseStep
	onError        func(error)
	finally        func()

	mu       sync.Mutex
	terminal bool // SetResult, SetError, or natural completion ran terminal handling
}

// NewPromise builds a Promise. If ctx is non-nil, each step's Use callback runs on that actor via ReenterAfter
// after Do completes. If ctx is nil, the entire chain runs in one background goroutine: each Use runs immediately
// after its Do in that same goroutine (no actor hop).
func NewPromise(ctx actor.Context, perStepTimeout time.Duration) *Promise {
	return &Promise{
		ctx:            ctx,
		perStepTimeout: perStepTimeout,
	}
}

func (p *Promise) Then(do func() (interface{}, error), use func(interface{}) error) *Promise {
	p.steps = append(p.steps, promiseStep{Do: do, Use: use})
	return p
}

func ThenRPC[T any](p *Promise, call func(context.Context) (T, error), use func(T) error) *Promise {
	return p.Then(
		func() (interface{}, error) {
			ctx, cancel := context.WithTimeout(context.Background(), p.perStepTimeout)
			defer cancel()
			return call(ctx)
		},
		func(v interface{}) error {
			return use(v.(T))
		},
	)
}

func (p *Promise) OnError(fn func(error)) *Promise {
	p.onError = fn
	return p
}

func (p *Promise) Finally(fn func()) *Promise {
	p.finally = fn
	return p
}

func (p *Promise) Run() {
	if len(p.steps) == 0 {
		p.completeSuccess()
		return
	}
	p.runStep(0)
}

// SetResult ends the promise successfully from outside the step chain: runs Finally once, skips remaining steps.
// Safe to call from any goroutine; later natural completion or SetError become no-ops.
func (p *Promise) SetResult() {
	p.completeSuccess()
}

// SetError ends the promise with an error from outside the step chain: runs OnError (if set) and Finally once.
// Passing nil is treated like SetResult (success terminal). Safe to call from any goroutine.
func (p *Promise) SetError(err error) {
	if err == nil {
		p.completeSuccess()
		return
	}
	p.completeError(err)
}

func (p *Promise) isTerminal() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.terminal
}

func (p *Promise) completeSuccess() {
	p.mu.Lock()
	if p.terminal {
		p.mu.Unlock()
		return
	}
	p.terminal = true
	p.mu.Unlock()
	if p.finally != nil {
		p.finally()
	}
}

func (p *Promise) completeError(err error) {
	p.mu.Lock()
	if p.terminal {
		p.mu.Unlock()
		return
	}
	p.terminal = true
	p.mu.Unlock()
	if p.onError != nil {
		p.onError(err)
	}
	if p.finally != nil {
		p.finally()
	}
}

func (p *Promise) runStep(i int) {
	if p.isTerminal() {
		return
	}
	if i >= len(p.steps) {
		p.completeSuccess()
		return
	}
	if p.ctx == nil {
		go p.runChainFrom(i)
		return
	}
	step := p.steps[i]
	f := actor.NewFuture(p.ctx.ActorSystem(), p.perStepTimeout)
	go func() {
		v, err := step.Do()
		p.ctx.Send(f.PID(), &promiseResult{Value: v, Err: err})
	}()
	p.ctx.ReenterAfter(f, func(res interface{}, ferr error) {
		if p.isTerminal() {
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
		if err := step.Use(pr.Value); err != nil {
			p.completeError(err)
			return
		}
		p.runStep(i + 1)
	})
}

// runChainFrom executes steps [i..] in the current goroutine: Do then Use for each step.
func (p *Promise) runChainFrom(i int) {
	for i < len(p.steps) {
		if p.isTerminal() {
			return
		}
		step := p.steps[i]
		v, err := step.Do()
		if err != nil {
			p.completeError(err)
			return
		}
		if p.isTerminal() {
			return
		}
		if err := step.Use(v); err != nil {
			p.completeError(err)
			return
		}
		i++
	}
	if !p.isTerminal() {
		p.completeSuccess()
	}
}
