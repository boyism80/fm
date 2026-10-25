package game

import (
	"time"

	"github.com/boyism80/fm/core/clock"
)

type RecordResetKind uint8

const (
	RecordResetNone RecordResetKind = iota
	RecordResetDaily
	RecordResetWeekly
	RecordResetEvery
)

type RecordReset struct {
	kind    RecordResetKind
	at      time.Duration
	every   time.Duration
	restart bool
}

type Record struct {
	value      int64
	text       string
	expiresAt  time.Time
	recordedAt time.Time
}

type Records struct {
	entries map[string]*Record
}

func NewRecords() *Records {
	return &Records{entries: make(map[string]*Record)}
}

func (reset RecordReset) expiresAt(now time.Time) time.Time {
	now = now.In(time.Local)
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	switch reset.kind {
	case RecordResetDaily:
		boundary := midnight.Add(reset.at)
		if !boundary.After(now) {
			boundary = boundary.AddDate(0, 0, 1)
		}
		return boundary
	case RecordResetWeekly:
		monday := midnight.AddDate(0, 0, -((int(now.Weekday()) + 6) % 7))
		boundary := monday.Add(reset.at)
		if !boundary.After(now) {
			boundary = boundary.AddDate(0, 0, 7)
		}
		return boundary
	case RecordResetEvery:
		return now.Add(reset.every)
	}
	return time.Time{}
}

func (r *Record) expired(now time.Time) bool {
	if r.expiresAt.IsZero() {
		return false
	}
	return !now.Before(r.expiresAt)
}

func (r *Records) find(key string) *Record {
	record := r.entries[key]
	if record == nil {
		return nil
	}
	if record.expired(clock.Now()) {
		delete(r.entries, key)
		return nil
	}
	return record
}

func (r *Records) touch(key string, reset RecordReset) *Record {
	now := clock.Now()
	record := r.find(key)
	switch {
	case record == nil:
		record = &Record{expiresAt: reset.expiresAt(now)}
		r.entries[key] = record
	case reset.restart:
		record.expiresAt = reset.expiresAt(now)
	}
	record.recordedAt = now
	return record
}

func (r *Records) Get(key string) int64 {
	record := r.find(key)
	if record == nil {
		return 0
	}
	return record.value
}

func (r *Records) Text(key string) string {
	record := r.find(key)
	if record == nil {
		return ""
	}
	return record.text
}

func (r *Records) Set(key string, value int64, reset RecordReset) {
	r.touch(key, reset).value = value
}

func (r *Records) SetText(key string, text string, reset RecordReset) {
	r.touch(key, reset).text = text
}

func (r *Records) Add(key string, n int64, reset RecordReset) int64 {
	record := r.touch(key, reset)
	record.value += n
	return record.value
}

func (r *Records) Remove(key string) {
	delete(r.entries, key)
}

func (r *Records) Clear() {
	clear(r.entries)
}
