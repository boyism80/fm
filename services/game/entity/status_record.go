package entity

import (
	"strconv"
	"time"

	"github.com/boyism80/fm/core/clock"
)

type StatusRecord struct {
	s string
}

func (r StatusRecord) AsString() string {
	return r.s
}

func (r StatusRecord) IsEmpty() bool {
	return r.s == ""
}

func (r *StatusRecord) WriteString(v string) {
	r.s = v
}

func (r *StatusRecord) WriteInt(v int) {
	r.s = strconv.Itoa(v)
}

func (r *StatusRecord) WriteBool(v bool) {
	if v {
		r.s = "1"
	} else {
		r.s = "0"
	}
}

func (r *StatusRecord) WriteDateTime(v time.Time) {
	if v.IsZero() {
		r.s = ""
	} else {
		r.s = strconv.FormatInt(v.UnixMilli(), 10)
	}
}

func (r StatusRecord) AsInt() (int, error) {
	return strconv.Atoi(r.s)
}

func (r StatusRecord) AsBool() (bool, error) {
	switch r.s {
	case "1":
		return true, nil
	case "0":
		return false, nil
	default:
		return false, strconv.ErrSyntax
	}
}

func (r StatusRecord) AsDateTime() (time.Time, error) {
	if r.s == "" {
		return time.Time{}, strconv.ErrSyntax
	}
	ms, err := strconv.ParseInt(r.s, 10, 64)
	if err != nil {
		parsed, parseErr := clock.ParseDateTime(r.s)
		if parseErr != nil {
			return time.Time{}, err
		}
		return parsed, nil
	}
	return time.UnixMilli(ms), nil
}
