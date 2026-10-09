package actor

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/bot/bot"
)

func (a *SuiteActor) cleanup() {
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
		if _, moving := a.moving[b]; moving || a.entering[b] {
			continue
		}
		if b.Gen == 0 {
			continue
		}
		commands := []string{"/인벤토리초기화", "/메소초기화"}
		if !reset[b.Map] {
			reset[b.Map] = true
			commands = append(commands, "/맵리셋")
		}
		for _, text := range commands {
			if err := a.Command(b, text); err != nil {
				log.Printf("[%s] cleanup %s %s: %v", a.name, b.Name, text, err)
			}
		}
	}
	time.AfterFunc(500*time.Millisecond, func() {
		a.actors.Send(a.self, &CleanupDone{})
	})
}

func (a *SuiteActor) close(ctx actor.Context) {
	for _, b := range a.bots {
		if _, moving := a.moving[b]; moving || a.entering[b] {
			continue
		}
		b.Close()
	}
	elapsed := time.Since(a.started)
	log.Printf("[%s] seat %d finished in %s (%d failures)", a.name, a.seat, elapsed.Round(time.Millisecond), len(a.failures))
	ctx.Send(ctx.Parent(), &SuiteFinished{
		Name:     a.name,
		Seat:     a.seat,
		Serial:   a.suite.Serial,
		Passed:   len(a.failures) == 0,
		Skipped:  a.skipped,
		Infra:    a.infra,
		Failures: a.failures,
		Elapsed:  elapsed,
	})

	a.closed = true
	if len(a.entering) > 0 || len(a.moving) > 0 {
		log.Printf("[%s] seat %d waits for %d bots still connecting", a.name, a.seat, len(a.entering)+len(a.moving))
		return
	}
	ctx.Stop(a.self)
}

func (a *SuiteActor) closeLate(b *bot.Bot) {
	b.Close()
	if !a.closed || len(a.entering) > 0 || len(a.moving) > 0 {
		return
	}
	a.actors.Stop(a.self)
}
