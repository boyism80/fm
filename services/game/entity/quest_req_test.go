package entity

import (
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/core/ensure"
	"github.com/boyism80/fm/services/game/wz"
)

type stubQuestGameWorld struct {
	resources *wz.Resources
}

func (s *stubQuestGameWorld) EnsureRedispatch(d *ensure.EnsureDeliver) {
}

func (s *stubQuestGameWorld) EnsureComplete(correlationID uint64) {
}

func (s *stubQuestGameWorld) GetPacketHandler() *core.PacketHandler {
	return nil
}

func (s *stubQuestGameWorld) GetResources() *wz.Resources {
	return s.resources
}

func (s *stubQuestGameWorld) GetExpRate() int {
	return 1
}

func (s *stubQuestGameWorld) GetDropRate() int {
	return 1
}

func (s *stubQuestGameWorld) GetMesoRate() int {
	return 1
}

func (s *stubQuestGameWorld) SaveAsync(ctx actor.Context, chars []*Character) *async.Promise {
	return nil
}

func (s *stubQuestGameWorld) GetMapSystem() MapSystem {
	return nil
}

func (s *stubQuestGameWorld) GetSchedulerSystem() SchedulerSystem {
	return nil
}

func (s *stubQuestGameWorld) GetGuildSystem() GuildSystem {
	return nil
}

func (s *stubQuestGameWorld) GetAllianceSystem() AllianceSystem {
	return nil
}

func (s *stubQuestGameWorld) GetPartySystem() PartySystem {
	return nil
}

func (s *stubQuestGameWorld) GetDispatchSystem() DispatchSystem {
	return nil
}

func newQuestContainerForTest(defs map[uint32]*wz.Quest) *QuestContainer {
	world := &stubQuestGameWorld{
		resources: &wz.Resources{Quests: defs},
	}
	return NewQuestContainer(&Character{
		LifeCore: LifeCore{
			ObjectCore: ObjectCore{
				GameWorld: world,
			},
		},
	})
}

func TestParseQuestEventTime(t *testing.T) {

	cases := []struct {
		raw   string
		ok    bool
		year  int
		month time.Month
		day   int
		hour  int
	}{
		{raw: "2007010100", ok: true, year: 2007, month: time.January, day: 1, hour: 0},
		{raw: "2007081600", ok: true, year: 2007, month: time.August, day: 16, hour: 0},
		{raw: "2007092000", ok: true, year: 2007, month: time.September, day: 20, hour: 0},
		{raw: "200701010", ok: false},
		{raw: "20070101000", ok: false},
		{raw: "abcd010100", ok: false},
	}

	for _, tc := range cases {
		got, ok := parseQuestEventTime(tc.raw)
		if ok != tc.ok {
			t.Fatalf("parseQuestEventTime(%q) ok=%v want %v", tc.raw, ok, tc.ok)
		}
		if !tc.ok {
			continue
		}
		if got.Year() != tc.year || got.Month() != tc.month || got.Day() != tc.day || got.Hour() != tc.hour {
			t.Fatalf("parseQuestEventTime(%q)=%v", tc.raw, got)
		}
	}
}

func TestMeetsQuestEventTimeRequirement(t *testing.T) {
	clock.Reset()
	t.Cleanup(clock.Reset)
	if err := clock.SetAbsolute(time.Date(2007, 1, 1, 12, 0, 0, 0, time.Local)); err != nil {
		t.Fatalf("SetAbsolute: %v", err)
	}

	qp := &Quest{}
	if !qp.meetsQuestEventTimeRequirement(wz.QuestReqTimeStart, "2007010100") {
		t.Fatal("expected start window open")
	}
	if qp.meetsQuestEventTimeRequirement(wz.QuestReqTimeStart, "2007010200") {
		t.Fatal("expected start window closed")
	}
	if !qp.meetsQuestEventTimeRequirement(wz.QuestReqTimeEnd, "2007010200") {
		t.Fatal("expected end window open")
	}
	if qp.meetsQuestEventTimeRequirement(wz.QuestReqTimeEnd, "2007010100") {
		t.Fatal("expected end window closed")
	}
}

func TestMeetsDayByDayRequirement(t *testing.T) {
	clock.Reset()
	t.Cleanup(clock.Reset)
	if err := clock.SetAbsolute(time.Date(2007, 1, 2, 10, 0, 0, 0, time.Local)); err != nil {
		t.Fatalf("SetAbsolute: %v", err)
	}

	started := &Quest{Status: QuestStatusStarted}
	if !started.meetsDayByDayRequirement() {
		t.Fatal("started quest should pass dayByDay")
	}

	completedToday := &Quest{
		Status:         QuestStatusCompleted,
		CompletionTime: time.Date(2007, 1, 2, 8, 0, 0, 0, time.Local),
	}
	if completedToday.meetsDayByDayRequirement() {
		t.Fatal("completed today should block dayByDay")
	}

	completedYesterday := &Quest{
		Status:         QuestStatusCompleted,
		CompletionTime: time.Date(2007, 1, 1, 23, 0, 0, 0, time.Local),
	}
	if !completedYesterday.meetsDayByDayRequirement() {
		t.Fatal("completed yesterday should allow dayByDay")
	}
}

func TestApplyInfoNumberActionCompletesKnownQuest(t *testing.T) {
	t.Parallel()

	def := &wz.Quest{ID: 2166}
	qc := newQuestContainerForTest(map[uint32]*wz.Quest{2166: def})
	ref := qc.Create(2166, QuestStatusStarted)
	qc.applyInfoNumberAction(2166)
	if ref.Status != QuestStatusCompleted {
		t.Fatalf("status=%d want completed", ref.Status)
	}
}
