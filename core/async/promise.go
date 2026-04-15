package async

import (
	"context"
	"fmt"
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
}

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
		p.runFinally()
		return
	}
	p.runStep(0)
}

func (p *Promise) runFinally() {
	if p.finally != nil {
		p.finally()
	}
}

func (p *Promise) fail(err error) {
	if p.onError != nil {
		p.onError(err)
	}
	p.runFinally()
}

func (p *Promise) runStep(i int) {
	if i >= len(p.steps) {
		p.runFinally()
		return
	}
	step := p.steps[i]
	f := actor.NewFuture(p.ctx.ActorSystem(), p.perStepTimeout)
	go func() {
		v, err := step.Do()
		p.ctx.Send(f.PID(), &promiseResult{Value: v, Err: err})
	}()
	p.ctx.ReenterAfter(f, func(res interface{}, ferr error) {
		if ferr != nil {
			p.fail(ferr)
			return
		}
		pr, ok := res.(*promiseResult)
		if !ok {
			p.fail(fmt.Errorf("async.Promise: unexpected result type %T", res))
			return
		}
		if pr.Err != nil {
			p.fail(pr.Err)
			return
		}
		if err := step.Use(pr.Value); err != nil {
			p.fail(err)
			return
		}
		p.runStep(i + 1)
	})
}
