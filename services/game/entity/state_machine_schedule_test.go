package entity

import (
	"testing"
	"time"
)

func TestParseStateMachineCron(t *testing.T) {
	sched, err := ParseStateMachineCron("5,15,25,35,45,55 * * * *")
	if err != nil {
		t.Fatal(err)
	}
	from := time.Date(2026, 8, 7, 12, 0, 0, 0, stateMachineCronLoc)
	delay := NextStateMachineCronDelay(sched, from)
	if delay <= 0 {
		t.Fatalf("expected positive delay, got %v", delay)
	}
	next := from.Add(delay)
	if next.Minute() != 5 {
		t.Fatalf("expected next minute 5, got %v", next)
	}
}
