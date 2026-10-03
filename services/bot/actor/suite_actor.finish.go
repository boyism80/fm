package actor

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

func (a *SuiteActor) finish() {
	if a.finished {
		return
	}
	a.finished = true
	a.deadline.Stop()
	for _, w := range a.waiters {
		w.timer.Stop()
	}
	a.waiters = nil

	reset := make(map[uint32]bool)
	for _, b := range a.bots {
		if b.Gen == 0 {
			continue
		}
		commands := []string{"/인벤토리초기화", "/메소초기화"}
		if reset[b.Map] == false {
			reset[b.Map] = true
			commands = append(commands, "/맵리셋")
		}
		for _, text := range commands {
			if err := a.Command(b, text); err != nil {
				log.Printf("[%s] cleanup %s %s: %v", a.suite.Name, b.Name, text, err)
			}
		}
	}
	time.AfterFunc(500*time.Millisecond, func() {
		a.root.Send(a.self, &CleanupDone{})
	})
}

func (a *SuiteActor) close(ctx actor.Context) {
	for _, b := range a.bots {
		b.Close()
	}
	elapsed := time.Since(a.started)
	log.Printf("[%s] seat %d finished in %s (%d failures)", a.suite.Name, a.seat, elapsed.Round(time.Millisecond), len(a.failures))
	ctx.Send(ctx.Parent(), &SuiteFinished{
		Name:     a.suite.Name,
		Seat:     a.seat,
		Serial:   a.suite.Serial,
		Passed:   len(a.failures) == 0,
		Infra:    a.infra,
		Failures: a.failures,
		Elapsed:  elapsed,
	})
	ctx.Stop(a.self)
}
