package entity

import (
	"testing"

	"github.com/boyism80/fm/services/game/wz"
)

func TestForceStartInternalQuest(t *testing.T) {
	t.Parallel()

	qc := newQuestContainerForTest(map[uint32]*wz.Quest{})
	qp, err := qc.ForceStart(7631, "100")
	if err != nil {
		t.Fatalf("ForceStart: %v", err)
	}
	if qp.Wz != nil {
		t.Fatal("expected nil Wz")
	}
	if qp.StatusRecord.AsString() != "100" {
		t.Fatalf("record=%q", qp.StatusRecord.AsString())
	}
	if qp.WiresToClient() {
		t.Fatal("internal quest should not wire to client")
	}
}

func TestMeetsInfoNumberRequirementInternalQuest(t *testing.T) {
	t.Parallel()

	qc := newQuestContainerForTest(map[uint32]*wz.Quest{})
	_, err := qc.ForceStart(7631, "100")
	if err != nil {
		t.Fatalf("ForceStart: %v", err)
	}

	def1049 := &wz.Quest{
		ID: 1049,
		Start: wz.QuestPhase{
			Requirements: []wz.QuestRequirement{
				{Kind: wz.QuestReqInfo, InfoStrings: []string{"100", ""}},
				{Kind: wz.QuestReqInfoNumber, IntValue: 7631},
			},
		},
	}
	qp := &Quest{
		container: qc,
		Wz:        def1049,
		QuestID:   1049,
		Status:    QuestStatusNotStarted,
	}
	phase := def1049.Start
	if !qp.meetsInfoNumberRequirement(7631, phase) {
		t.Fatal("expected warrior branch to pass")
	}
	phase.Requirements[0].InfoStrings = []string{"200", ""}
	if qp.meetsInfoNumberRequirement(7631, phase) {
		t.Fatal("expected mage branch to fail")
	}
}

func TestMeetsInfoNumberRequirementChecksCatalogQuest(t *testing.T) {
	t.Parallel()

	def2166 := &wz.Quest{ID: 2166}
	qc := newQuestContainerForTest(map[uint32]*wz.Quest{2166: def2166})
	qp := &Quest{
		container: qc,
		Wz:        def2166,
		QuestID:   2166,
		Status:    QuestStatusStarted,
	}
	phase := wz.QuestPhase{}
	if qp.meetsInfoNumberRequirement(2166, phase) {
		t.Fatal("expected missing status record to fail")
	}

	ref := qc.Create(2166, QuestStatusStarted)
	ref.StatusRecord.WriteString("1")
	if !qp.meetsInfoNumberRequirement(2166, phase) {
		t.Fatal("expected started ref with record to pass")
	}
}

func TestApplyInfoNumberActionInternalQuest(t *testing.T) {
	t.Parallel()

	qc := newQuestContainerForTest(map[uint32]*wz.Quest{})
	ref, err := qc.ForceStart(7631, "100")
	if err != nil {
		t.Fatalf("ForceStart: %v", err)
	}
	qc.applyInfoNumberAction(7631)
	if ref.Status != QuestStatusCompleted {
		t.Fatalf("status=%d want completed", ref.Status)
	}
}
