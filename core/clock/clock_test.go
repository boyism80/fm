package clock

import (
	"testing"
	"time"
)

func TestSetAbsoluteAndReset(t *testing.T) {
	Reset()
	defer Reset()

	target := Now().Add(3 * time.Hour)
	if err := SetAbsolute(target); err != nil {
		t.Fatalf("SetAbsolute: %v", err)
	}
	diff := Now().Sub(target)
	if diff < -time.Second || diff > time.Second {
		t.Fatalf("Now() drift from target: %v", diff)
	}

	Reset()
	if Offset() != 0 {
		t.Fatalf("expected zero offset after reset, got %v", Offset())
	}
}

func TestApplyDateTime(t *testing.T) {
	Reset()
	defer Reset()

	if err := ApplyDateTime(true, ""); err != nil {
		t.Fatalf("ApplyDateTime reset: %v", err)
	}
	if Offset() != 0 {
		t.Fatalf("expected zero offset")
	}

	target := time.Date(2007, 1, 1, 12, 0, 0, 0, time.Local)
	if err := ApplyDateTime(false, FormatDateTime(target)); err != nil {
		t.Fatalf("ApplyDateTime: %v", err)
	}
	got := Now()
	if got.Year() != 2007 || got.Month() != time.January || got.Day() != 1 || got.Hour() != 12 {
		t.Fatalf("unexpected Now(): %v", got)
	}
}

func TestParseTimespan(t *testing.T) {
	d, err := ParseTimespan("3:00:00")
	if err != nil {
		t.Fatalf("ParseTimespan: %v", err)
	}
	if d != 3*time.Hour {
		t.Fatalf("expected 3h, got %v", d)
	}

	d, err = ParseTimespan("10.12:30:00")
	if err != nil {
		t.Fatalf("ParseTimespan days: %v", err)
	}
	expected := 10*24*time.Hour + 12*time.Hour + 30*time.Minute
	if d != expected {
		t.Fatalf("expected %v, got %v", expected, d)
	}
}
