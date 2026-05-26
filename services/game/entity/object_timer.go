package entity

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type ObjectTimer struct {
	Timer      *time.Timer
	Interval   time.Duration
	Repeat     bool
	Callback   func()
	NextFireAt time.Time
	Remaining  time.Duration
}

func (obj *ObjectCore) initTimers() {
	obj.timers = make(map[string]*ObjectTimer)
}

func (obj *ObjectCore) AddTimer(key string, interval time.Duration, repeat bool, callback func()) bool {
	if _, exists := obj.timers[key]; exists {
		obj.RemoveTimer(key)
	}
	if obj.GameWorld == nil || obj.self == nil {
		return false
	}
	mapInstance := obj.Map
	if mapInstance == nil {
		return false
	}
	pid := mapInstance.GetActorPID()
	if pid == nil {
		return false
	}
	entry := &ObjectTimer{
		Interval:   interval,
		Repeat:     repeat,
		Callback:   callback,
		NextFireAt: time.Now().Add(interval),
	}
	k := key
	entry.Timer = time.AfterFunc(interval, func() {
		obj.GameWorld.DispatchRunObjectTimer(pid, obj.self, k)
	})
	obj.timers[key] = entry
	return true
}

func (obj *ObjectCore) RemoveTimer(key string) bool {
	entry := obj.timers[key]
	if entry == nil {
		return false
	}
	if entry.Timer != nil {
		entry.Timer.Stop()
	}
	delete(obj.timers, key)
	return true
}

func (obj *ObjectCore) GetTimerEntry(key string) *ObjectTimer {
	return obj.timers[key]
}

func (obj *ObjectCore) ClearTimers() {
	for key, entry := range obj.timers {
		if entry != nil && entry.Timer != nil {
			entry.Timer.Stop()
		}
		delete(obj.timers, key)
	}
}

func (obj *ObjectCore) SuspendTimers() {
	now := time.Now()
	for _, entry := range obj.timers {
		if entry == nil || entry.Timer == nil {
			continue
		}
		entry.Remaining = entry.NextFireAt.Sub(now)
		if entry.Remaining < 0 {
			entry.Remaining = 0
		}
		entry.Timer.Stop()
		entry.Timer = nil
	}
}

func (obj *ObjectCore) ResumeTimers(pid *actor.PID) {
	if pid == nil || obj.GameWorld == nil || obj.self == nil {
		return
	}
	for key, entry := range obj.timers {
		if entry == nil || entry.Timer != nil {
			continue
		}
		duration := entry.Remaining
		if duration <= 0 {
			duration = entry.Interval
		}
		entry.Remaining = 0
		entry.NextFireAt = time.Now().Add(duration)
		k := key
		entry.Timer = time.AfterFunc(duration, func() {
			obj.GameWorld.DispatchRunObjectTimer(pid, obj.self, k)
		})
	}
}
