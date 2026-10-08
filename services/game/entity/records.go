package entity

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
	Kind    RecordResetKind
	At      time.Duration
	Every   time.Duration
	Restart bool
}

type Record struct {
	Value     int64
	Text      string
	ExpiresAt time.Time
	UpdatedAt time.Time
}

type Records struct {
	entries map[string]*Record
}

func NewRecords() *Records {
	return &Records{entries: make(map[string]*Record)}
}

func (reset RecordReset) expiresAt(now time.Time) time.Time {
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	switch reset.Kind {
	case RecordResetDaily:
		boundary := midnight.Add(reset.At)
		if boundary.After(now) == false {
			boundary = boundary.AddDate(0, 0, 1)
		}
		return boundary
	case RecordResetWeekly:
		monday := midnight.AddDate(0, 0, -((int(now.Weekday()) + 6) % 7))
		boundary := monday.Add(reset.At)
		if boundary.After(now) == false {
			boundary = boundary.AddDate(0, 0, 7)
		}
		return boundary
	case RecordResetEvery:
		return now.Add(reset.Every)
	}
	return time.Time{}
}

func (r *Record) expired(now time.Time) bool {
	if r.ExpiresAt.IsZero() {
		return false
	}
	return now.Before(r.ExpiresAt) == false
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
		record = &Record{ExpiresAt: reset.expiresAt(now)}
		r.entries[key] = record
	case reset.Restart:
		record.ExpiresAt = reset.expiresAt(now)
	}
	record.UpdatedAt = now
	return record
}

func (r *Records) Get(key string) int64 {
	record := r.find(key)
	if record == nil {
		return 0
	}
	return record.Value
}

func (r *Records) Text(key string) string {
	record := r.find(key)
	if record == nil {
		return ""
	}
	return record.Text
}

func (r *Records) Set(key string, value int64, reset RecordReset) {
	r.touch(key, reset).Value = value
}

func (r *Records) SetText(key string, text string, reset RecordReset) {
	r.touch(key, reset).Text = text
}

func (r *Records) Add(key string, n int64, reset RecordReset) int64 {
	record := r.touch(key, reset)
	record.Value += n
	return record.Value
}

func (r *Records) Remove(key string) {
	delete(r.entries, key)
}

func (r *Records) Clear() {
	clear(r.entries)
}
