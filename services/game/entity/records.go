package entity

import (
	"time"

	"github.com/boyism80/fm/core/clock"
)

type RecordPeriod uint8

const (
	RecordPeriodNone RecordPeriod = iota
	RecordPeriodDaily
	RecordPeriodWeekly
)

type Record struct {
	Value     int64
	Text      string
	Period    RecordPeriod
	UpdatedAt time.Time
}

type Records struct {
	entries map[string]*Record
}

func NewRecords() *Records {
	return &Records{entries: make(map[string]*Record)}
}

func (r *Record) expired(now time.Time) bool {
	switch r.Period {
	case RecordPeriodDaily:
		y1, m1, d1 := r.UpdatedAt.Date()
		y2, m2, d2 := now.Date()
		return y1 != y2 || m1 != m2 || d1 != d2
	case RecordPeriodWeekly:
		y1, w1 := r.UpdatedAt.ISOWeek()
		y2, w2 := now.ISOWeek()
		return y1 != y2 || w1 != w2
	}
	return false
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

func (r *Records) touch(key string, period RecordPeriod) *Record {
	record := r.find(key)
	if record == nil {
		record = &Record{}
		r.entries[key] = record
	}
	record.Period = period
	record.UpdatedAt = clock.Now()
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

func (r *Records) Set(key string, value int64, period RecordPeriod) {
	r.touch(key, period).Value = value
}

func (r *Records) SetText(key string, text string, period RecordPeriod) {
	r.touch(key, period).Text = text
}

func (r *Records) Add(key string, n int64, period RecordPeriod) int64 {
	record := r.touch(key, period)
	record.Value += n
	return record.Value
}

func (r *Records) Remove(key string) {
	delete(r.entries, key)
}

func (r *Records) Clear() {
	clear(r.entries)
}
