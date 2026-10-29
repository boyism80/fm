package actor

import (
	"sync/atomic"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/maple/wz"
)

// AgentActor runs the suites a cluster master assigns, one suite per seat. report is called on the actor's
// goroutine with the run of the suite and a SuiteStarted or SuiteFinished.
type AgentActor struct {
	cfg    *config.Bot
	wz     *wz.Resources
	atlas  *Atlas
	report func(run string, msg any)
	bots   string
	botSeq atomic.Int32
	runs   map[int]string
}

func NewAgentActor(cfg *config.Bot, resources *wz.Resources, report func(run string, msg any)) *AgentActor {
	return &AgentActor{
		cfg:    cfg,
		wz:     resources,
		atlas:  NewAtlas(resources),
		report: report,
		runs:   make(map[int]string),
	}
}

func (a *AgentActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *AssignSuite:
		// Bot names come from Bots, so the sequence restarts with each run's names.
		if msg.Bots != a.bots {
			a.bots = msg.Bots
			a.botSeq.Store(0)
		}
		a.runs[msg.Seat] = msg.Run
		a.report(msg.Run, &SuiteStarted{Suite: msg.Suite.Name, Seat: msg.Seat})
		ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return NewSuiteActor(a.cfg, a.wz, a.atlas, msg.Suite, msg.Seat, msg.Bots, &a.botSeq)
		}))
	case *SuiteFinished:
		run := a.runs[msg.Seat]
		delete(a.runs, msg.Seat)
		a.report(run, msg)
	case *RunStop:
		for _, child := range ctx.Children() {
			ctx.Send(child, &SuiteStop{})
		}
	}
}
