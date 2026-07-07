package entity

import (
	"testing"
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/services/game/wz"
)

func TestPrepareCompleteRejectsExpiredDeadline(t *testing.T) {
	t.Parallel()

	clock.Reset()
	t.Cleanup(clock.Reset)

	def := &wz.Quest{ID: 29002}
	qc := newQuestContainerForTest(map[uint32]*wz.Quest{29002: def})
	qp := qc.Create(29002, QuestStatusStarted)
	if qp == nil {
		t.Fatal("create failed")
	}
	qp.SetDeadline(clock.Now().Add(-time.Minute))

	_, _, err := qc.prepareComplete(qp, qc.owner, QuestPrepareOpts{})
	if err != ErrQuestExpired {
		t.Fatalf("err=%v want %v", err, ErrQuestExpired)
	}
}

func TestCompleteRejectsExpiredDeadline(t *testing.T) {
	t.Parallel()

	clock.Reset()
	t.Cleanup(clock.Reset)

	def := &wz.Quest{ID: 29002}
	qc := newQuestContainerForTest(map[uint32]*wz.Quest{29002: def})
	qp := qc.Create(29002, QuestStatusStarted)
	if qp == nil {
		t.Fatal("create failed")
	}
	qp.SetDeadline(clock.Now().Add(-time.Minute))

	err := qp.Complete(qc.owner, QuestPrepareOpts{})
	if err != ErrQuestExpired {
		t.Fatalf("err=%v want %v", err, ErrQuestExpired)
	}
}

func TestPrepareCompleteAllowsBeforeDeadline(t *testing.T) {
	t.Parallel()

	clock.Reset()
	t.Cleanup(clock.Reset)

	def := &wz.Quest{ID: 29002}
	qc := newQuestContainerForTest(map[uint32]*wz.Quest{29002: def})
	qp := qc.Create(29002, QuestStatusStarted)
	if qp == nil {
		t.Fatal("create failed")
	}
	qp.SetDeadline(clock.Now().Add(time.Hour))

	_, _, err := qc.prepareComplete(qp, qc.owner, QuestPrepareOpts{})
	if err != nil {
		t.Fatalf("prepareComplete: %v", err)
	}
}

func TestCanForfeitWhenExpired(t *testing.T) {
	t.Parallel()

	clock.Reset()
	t.Cleanup(clock.Reset)

	def := &wz.Quest{ID: 29002}
	qc := newQuestContainerForTest(map[uint32]*wz.Quest{29002: def})
	qp := qc.Create(29002, QuestStatusStarted)
	if qp == nil {
		t.Fatal("create failed")
	}
	qp.SetDeadline(clock.Now().Add(-time.Minute))

	if !qp.CanForfeit() {
		t.Fatal("expected forfeitable when expired")
	}
}

func TestStartClearsDeadline(t *testing.T) {
	t.Parallel()

	clock.Reset()
	t.Cleanup(clock.Reset)

	qc := newQuestContainerForTest(map[uint32]*wz.Quest{})
	qp, err := qc.ForceStart(7631, "100")
	if err != nil {
		t.Fatalf("ForceStart: %v", err)
	}
	qp.SetDeadline(clock.Now().Add(time.Hour))
	qc.Remove(7631)

	qp, err = qc.ForceStart(7631, "100")
	if err != nil {
		t.Fatalf("ForceStart: %v", err)
	}
	if !qp.Deadline.IsZero() {
		t.Fatalf("deadline=%v want zero", qp.Deadline)
	}
}

func TestIsCompletableFalseWhenExpired(t *testing.T) {
	t.Parallel()

	clock.Reset()
	t.Cleanup(clock.Reset)

	def := &wz.Quest{ID: 29002}
	qc := newQuestContainerForTest(map[uint32]*wz.Quest{29002: def})
	qp := qc.Create(29002, QuestStatusStarted)
	if qp == nil {
		t.Fatal("create failed")
	}
	qp.SetDeadline(clock.Now().Add(-time.Minute))

	if qp.IsCompletable(qc.owner) {
		t.Fatal("expected not completable")
	}
}
