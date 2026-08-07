package entity

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	"github.com/robfig/cron/v3"
	_ "time/tzdata"
)

var (
	stateMachineCronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	stateMachineCronLoc    *time.Location
)

func init() {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		loc = time.Local
	}
	stateMachineCronLoc = loc
}

func ParseStateMachineCron(expr string) (cron.Schedule, error) {
	if expr == "" {
		return nil, fmt.Errorf("cron expression is empty")
	}
	return stateMachineCronParser.Parse(expr)
}

func NextStateMachineCronDelay(sched cron.Schedule, from time.Time) time.Duration {
	if sched == nil {
		return 0
	}
	now := from.In(stateMachineCronLoc)
	next := sched.Next(now)
	if next.IsZero() {
		return 0
	}
	d := next.Sub(now)
	if d < time.Millisecond {
		next = sched.Next(next)
		d = next.Sub(now)
	}
	if d < time.Millisecond {
		return 0
	}
	return d
}

func (sm *StateMachine) After(id string, ms int64, hook string) {
	if sm == nil || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil {
		return
	}
	if id == "" {
		id = hook
	}
	if hook == "" || ms <= 0 {
		return
	}
	sm.Group.GameWorld.SendStateMachineMessage(sm.ActorPID, &ScheduleStateMachineAfter{
		ID:           id,
		Milliseconds: ms,
		Hook:         hook,
	})
}

func (sm *StateMachine) AfterAsync(ctx actor.Context, id string, ms int64, hook string) *async.Promise {
	if sm == nil || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil || ctx == nil {
		return nil
	}
	if id == "" {
		id = hook
	}
	if hook == "" || ms <= 0 {
		return nil
	}
	gw := sm.Group.GameWorld
	pid := sm.ActorPID
	return async.Ask(ctx, pid, 10*time.Second, func(replyTo *actor.PID) {
		gw.SendStateMachineMessage(pid, &ScheduleStateMachineAfter{
			ID:           id,
			Milliseconds: ms,
			Hook:         hook,
			ReplyTo:      replyTo,
		})
	})
}

func (sm *StateMachine) Cron(id, expr, hook string) {
	if sm == nil || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil {
		return
	}
	if id == "" {
		id = hook
	}
	if hook == "" || expr == "" {
		return
	}
	sm.Group.GameWorld.SendStateMachineMessage(sm.ActorPID, &ScheduleStateMachineCron{
		ID:   id,
		Expr: expr,
		Hook: hook,
	})
}

func (sm *StateMachine) CronAsync(ctx actor.Context, id, expr, hook string) *async.Promise {
	if sm == nil || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil || ctx == nil {
		return nil
	}
	if id == "" {
		id = hook
	}
	if hook == "" || expr == "" {
		return nil
	}
	gw := sm.Group.GameWorld
	pid := sm.ActorPID
	return async.Ask(ctx, pid, 10*time.Second, func(replyTo *actor.PID) {
		gw.SendStateMachineMessage(pid, &ScheduleStateMachineCron{
			ID:      id,
			Expr:    expr,
			Hook:    hook,
			ReplyTo: replyTo,
		})
	})
}

func (sm *StateMachine) CancelSchedule(id string) {
	if sm == nil || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil {
		return
	}
	sm.Group.GameWorld.SendStateMachineMessage(sm.ActorPID, &CancelStateMachineNamedSchedule{ID: id})
}

func (sm *StateMachine) CancelScheduleAsync(ctx actor.Context, id string) *async.Promise {
	if sm == nil || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil || ctx == nil {
		return nil
	}
	gw := sm.Group.GameWorld
	pid := sm.ActorPID
	return async.Ask(ctx, pid, 10*time.Second, func(replyTo *actor.PID) {
		gw.SendStateMachineMessage(pid, &CancelStateMachineNamedSchedule{
			ID:      id,
			ReplyTo: replyTo,
		})
	})
}

func (sm *StateMachine) CancelAllSchedules() {
	sm.CancelSchedule("")
}

type ScheduleStateMachineAfter struct {
	ID           string
	Milliseconds int64
	Hook         string
	ReplyTo      *actor.PID
}

type ScheduleStateMachineCron struct {
	ID      string
	Expr    string
	Hook    string
	ReplyTo *actor.PID
}

type CancelStateMachineNamedSchedule struct {
	ID      string
	ReplyTo *actor.PID
}

type StateMachineNamedTimeout struct {
	ID      string
	Version uint64
}
