package clock

import (
	"errors"
	"fmt"
	"time"
)

const DateTimeLayout = "2006-01-02 15:04:05"

var ErrInvalidDateTime = errors.New("invalid datetime")

func ParseDateTime(raw string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, ErrInvalidDateTime
	}
	parsed, err := time.ParseInLocation(DateTimeLayout, raw, time.Local)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %v", ErrInvalidDateTime, err)
	}
	return parsed, nil
}

func FormatDateTime(t time.Time) string {
	return t.In(time.Local).Format(DateTimeLayout)
}

func SetAbsoluteString(raw string) error {
	parsed, err := ParseDateTime(raw)
	if err != nil {
		return err
	}
	return SetAbsolute(parsed)
}

func ApplyDateTime(reset bool, raw string) error {
	if reset {
		Reset()
		return nil
	}
	return SetAbsoluteString(raw)
}

func DateTimeTable(t time.Time) (year, month, day, hour, minute, second int) {
	local := t.In(time.Local)
	year, mon, day := local.Date()
	hour, minute, second = local.Clock()
	return year, int(mon), day, hour, minute, second
}
