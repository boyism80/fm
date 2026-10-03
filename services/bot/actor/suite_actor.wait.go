package actor

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/bot/bot"
	"github.com/boyism80/fm/services/bot/conn"
	"github.com/boyism80/fm/stream"
	lua "github.com/yuin/gopher-lua"
)

type outbound interface {
	Opcode() byte
	Serialize(writer *stream.StreamWriter) error
}

type waiter struct {
	id     uint64
	thread *lua.LState
	bot    *bot.Bot
	gen    int
	match  func(any) bool
	done   func(pkt any, ok bool)
	timer  *time.Timer
}

func (a *SuiteActor) Request(thread *lua.LState, b *bot.Bot, pkt outbound, timeout time.Duration, match func(any) bool, done func(pkt any, ok bool)) {
	if a.finished {
		return
	}
	a.nextWait++
	w := &waiter{id: a.nextWait, thread: thread, bot: b, gen: b.Gen, match: match, done: done}
	w.timer = time.AfterFunc(timeout, func() {
		a.actors.Send(a.self, &WaitTimeout{Wait: w.id})
	})
	a.waiters = append(a.waiters, w)

	if pkt == nil {
		return
	}
	if err := b.Send(pkt); err != nil {
		a.Fail(fmt.Sprintf("%s send: %v", b.Name, err))
		w.timer.Stop()
		a.actors.Send(a.self, &WaitTimeout{Wait: w.id})
	}
}

func (a *SuiteActor) Sleep(thread *lua.LState, d time.Duration) {
	if a.finished {
		return
	}
	a.sleeping[thread] = true
	time.AfterFunc(d, func() {
		a.actors.Send(a.self, &SleepDone{Thread: thread})
	})
}

func (a *SuiteActor) Command(b *bot.Bot, text string) error {
	return b.Send(&request.NormalChat{Message: text})
}

func (a *SuiteActor) warpByCommand(thread *lua.LState, b *bot.Bot, text string, mapID uint32, done func(bool)) {
	a.Request(thread, b, &request.NormalChat{Message: text}, a.timeout(), func(pkt any) bool {
		warp, ok := pkt.(*response.Warp)
		return ok && warp.Character != nil && warp.Character.Map == mapID
	}, func(_ any, ok bool) {
		if ok == false {
			done(a.Fail(fmt.Sprintf("%s %s: 이동하지 않음", b.Name, text)))
			return
		}
		done(true)
	})
}

func (a *SuiteActor) timeout() time.Duration {
	return time.Duration(a.cfg.TimeoutMs) * time.Millisecond
}

func (a *SuiteActor) packetReceived(msg *PacketReceived) {
	b := a.bots[msg.Bot]
	if a.finished || msg.Gen != b.Gen {
		return
	}
	pkt, err := conn.Decode(msg.Opcode, msg.Body)
	if err != nil {
		return
	}
	b.Update(pkt)
	for _, fn := range a.hooks {
		fn(b, pkt)
	}

	if notice, ok := pkt.(*response.Notice); ok && strings.Contains(notice.Message, "권한이 부족합니다") {
		a.Fail(fmt.Sprintf("%s: 권한이 부족합니다", b.Name))
		a.endWaiters(b)
		return
	}

	for _, w := range a.waiters {
		if w.bot != b || w.gen != msg.Gen {
			continue
		}
		if w.match != nil && w.match(pkt) == false {
			continue
		}
		a.removeWaiter(w)
		w.done(pkt, true)
		return
	}
}

func (a *SuiteActor) disconnected(msg *Disconnected) {
	b := a.bots[msg.Bot]
	if a.finished || msg.Gen != b.Gen {
		return
	}
	log.Printf("[%s] %s disconnected: %v", a.name, b.Name, msg.Err)
	a.Fail(fmt.Sprintf("%s 연결 끊김: %v", b.Name, msg.Err))
	b.Close()
	a.endWaiters(b)
}

func (a *SuiteActor) waitTimeout(msg *WaitTimeout) {
	for _, w := range a.waiters {
		if w.id != msg.Wait {
			continue
		}
		a.removeWaiter(w)
		w.done(nil, false)
		return
	}
}

func (a *SuiteActor) endWaiters(b *bot.Bot) {
	var ended []*waiter
	for _, w := range a.waiters {
		if w.bot == b {
			ended = append(ended, w)
		}
	}
	for _, w := range ended {
		a.removeWaiter(w)
		w.done(nil, false)
	}
}

func (a *SuiteActor) removeWaiter(w *waiter) {
	w.timer.Stop()
	for i, item := range a.waiters {
		if item == w {
			a.waiters = append(a.waiters[:i], a.waiters[i+1:]...)
			return
		}
	}
}
