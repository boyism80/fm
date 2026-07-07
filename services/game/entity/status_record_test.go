package entity

import (
	"testing"
	"time"
)

func TestStatusRecordWriteAsInt(t *testing.T) {
	t.Parallel()

	var r StatusRecord
	r.WriteInt(100)
	if r.AsString() != "100" {
		t.Fatalf("string=%q", r.AsString())
	}
	v, err := r.AsInt()
	if err != nil || v != 100 {
		t.Fatalf("int=%d err=%v", v, err)
	}
}

func TestStatusRecordWriteAsBool(t *testing.T) {
	t.Parallel()

	var r StatusRecord
	r.WriteBool(true)
	if r.AsString() != "1" {
		t.Fatalf("string=%q", r.AsString())
	}
	ok, err := r.AsBool()
	if err != nil || !ok {
		t.Fatalf("bool=%v err=%v", ok, err)
	}
	r.WriteBool(false)
	ok, err = r.AsBool()
	if err != nil || ok {
		t.Fatalf("bool=%v err=%v", ok, err)
	}
}

func TestStatusRecordWriteAsDateTime(t *testing.T) {
	t.Parallel()

	want := time.Date(2026, 7, 6, 12, 30, 0, 0, time.Local)
	var r StatusRecord
	r.WriteDateTime(want)
	got, err := r.AsDateTime()
	if err != nil {
		t.Fatalf("AsDateTime: %v", err)
	}
	if !got.Equal(want) {
		t.Fatalf("got=%v want=%v", got, want)
	}
}
