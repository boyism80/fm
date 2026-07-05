package entity

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/boyism80/fm/services/game/wz"
)

type QuestStatusType uint8

const (
	QuestStatusNotStarted QuestStatusType = 0
	QuestStatusStarted    QuestStatusType = 1
	QuestStatusCompleted  QuestStatusType = 2
)

type Quest struct {
	container      *QuestContainer
	Wz             *wz.Quest
	QuestID        uint32
	Status         QuestStatusType
	MobKills       map[uint32]int
	StatusRecord   string
	Unknown2       map[string]string
	CompletionTime time.Time
	Forfeited      int
}

func (qp *Quest) WiresToClient() bool {
	if qp == nil {
		return false
	}
	return qp.Wz != nil
}

func (qp *Quest) IsStarted() bool {
	if qp == nil {
		return false
	}
	return qp.Status == QuestStatusStarted
}

func (qp *Quest) MatchesState(state int) bool {
	if qp == nil {
		return state == int(QuestStatusNotStarted)
	}
	return int(qp.Status) == state
}

func (qp *Quest) MeetsMobCounts(mobs []wz.QuestMobCount) bool {
	for _, mob := range mobs {
		kills := 0
		if qp != nil && qp.MobKills != nil {
			kills = qp.MobKills[mob.MobID]
		}
		if kills < mob.Count {
			return false
		}
	}
	return true
}

func (qp *Quest) IsCompletable(ch *Character) bool {
	if qp == nil || ch == nil || ch.Quests == nil || !qp.IsStarted() || qp.Wz == nil {
		return false
	}
	if qp.Wz.Meta.Blocked {
		return false
	}
	return qp.meetsPhaseRequirements(ch, qp.Wz.Complete, QuestPrepareOpts{NpcID: nil})
}

func (qp *Quest) StartedWirePayload() string {
	if qp == nil || qp.Wz == nil {
		return ""
	}
	if len(qp.Wz.OrderedMobIDs()) > 0 {
		return qp.mobKillEncodedString()
	}
	return qp.StatusRecord
}

func (qp *Quest) RecordMobKill(mobID uint32) bool {
	if qp == nil || !qp.IsStarted() || qp.Wz == nil || mobID == 0 {
		return false
	}
	required := qp.Wz.RelevantMobs()
	need, ok := required[mobID]
	if !ok || need <= 0 {
		return false
	}
	if qp.MobKills == nil {
		qp.MobKills = make(map[uint32]int)
	}
	if qp.MobKills[mobID] >= need {
		return false
	}
	qp.MobKills[mobID]++
	return true
}

func (qp *Quest) mobKillEncodedString() string {
	if qp == nil || qp.Wz == nil {
		return ""
	}
	ids := qp.Wz.OrderedMobIDs()
	var b strings.Builder
	for _, id := range ids {
		kills := 0
		if qp.MobKills != nil {
			kills = qp.MobKills[id]
		}
		fmt.Fprintf(&b, "%03d", kills)
	}
	return b.String()
}

func (qp *Quest) MobKillCountsOrdered() []uint16 {
	if qp == nil || qp.Wz == nil {
		return nil
	}
	ids := qp.Wz.OrderedMobIDs()
	if len(ids) == 0 {
		return nil
	}
	out := make([]uint16, len(ids))
	for i, id := range ids {
		kills := 0
		if qp.MobKills != nil {
			kills = qp.MobKills[id]
		}
		out[i] = uint16(kills)
	}
	return out
}

func (ch *Character) OnQuestMobKilled(mobID uint32) {
	if ch == nil || mobID == 0 || ch.Quests == nil {
		return
	}
	ch.Quests.ForEach(func(_ uint32, qp *Quest) {
		if qp == nil || !qp.RecordMobKill(mobID) {
			return
		}
		ch.Listener.OnQuestProgress(ch, qp)
	})
}

func (qp *Quest) InitMobKillCounters() {
	if qp == nil || qp.Wz == nil {
		return
	}
	mobs := qp.Wz.RelevantMobs()
	if len(mobs) == 0 {
		return
	}
	if qp.MobKills == nil {
		qp.MobKills = make(map[uint32]int)
	}
	for mobID := range mobs {
		if _, ok := qp.MobKills[mobID]; !ok {
			qp.MobKills[mobID] = 0
		}
	}
}

func (qp *Quest) ApplyPhaseMeta(meta questPhaseMeta, qc *QuestContainer) {
	if qp == nil {
		return
	}
	if meta.info != "" {
		qp.StatusRecord = meta.info
	}
	if len(meta.chainQuests) > 0 && qc != nil {
		qc.applyQuestChainActions(meta.chainQuests)
	}
	if len(meta.infoNumberQuests) > 0 && qc != nil {
		for _, refID := range meta.infoNumberQuests {
			qc.applyInfoNumberAction(refID)
		}
	}
}

func (qp *Quest) CanForfeit() bool {
	if qp == nil || !qp.IsStarted() || qp.Wz == nil {
		return false
	}
	switch qp.QuestID {
	case 20000, 20010, 20015, 20020:
		return false
	default:
		return true
	}
}

func (qp *Quest) CanRestoreLostItem(ch *Character, itemID uint32) bool {
	if qp == nil || qp.Wz == nil || itemID == 0 || ch == nil {
		return false
	}
	if !qp.IsStarted() {
		return false
	}
	if ch.HasItem(itemID) {
		return false
	}
	for _, act := range qp.Wz.Start.Actions {
		if act.Kind != wz.QuestActItem {
			continue
		}
		for _, item := range act.Items {
			if item.ItemID == itemID && item.Count > 0 {
				return true
			}
		}
	}
	return false
}

func (qp *Quest) RestoreLostItem(ch *Character, itemID uint32) error {
	if !qp.CanRestoreLostItem(ch, itemID) {
		return ErrQuestRestoreItem
	}
	if ch == nil || qp.Wz == nil {
		return ErrQuestRestoreItem
	}
	var count uint16
	found := false
	for _, act := range qp.Wz.Start.Actions {
		if act.Kind != wz.QuestActItem {
			continue
		}
		for _, item := range act.Items {
			if item.ItemID == itemID && item.Count > 0 {
				count = uint16(item.Count)
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		return ErrQuestRestoreItem
	}
	spec := FlowSpec{
		Reward: FlowSide{
			Items: map[uint32]uint16{itemID: count},
		},
	}
	if ch.Exchange(spec) != FlowOK {
		return ErrQuestFlowFailed
	}
	return nil
}

func (qp *Quest) Complete(ch *Character, opts QuestPrepareOpts) error {
	if qp == nil || ch == nil || ch.Quests == nil || qp.Wz == nil {
		return ErrQuestNotCompletable
	}
	wireNPC := uint32(0)
	if opts.NpcID != nil {
		wireNPC = *opts.NpcID
	}
	var meta questPhaseMeta
	if !opts.Force {
		spec, m, err := ch.Quests.prepareComplete(qp, ch, opts)
		if err != nil {
			return err
		}
		meta = m
		if ch.Exchange(spec) != FlowOK {
			return ErrQuestFlowFailed
		}
		qp.ApplyPhaseMeta(meta, ch.Quests)
	} else if !qp.IsStarted() {
		return ErrQuestNotCompletable
	}
	qp.Status = QuestStatusCompleted
	qp.CompletionTime = time.Now()
	nextQuestID := uint32(0)
	if qp.Wz != nil {
		nextQuestID = qp.Wz.NextQuestID()
	}
	ch.Listener.OnQuestCompleted(ch, qp, wireNPC, nextQuestID)
	return nil
}

func (qp *Quest) Forfeit(ch *Character) error {
	if qp == nil || ch == nil || ch.Quests == nil {
		return ErrQuestNotForfeitable
	}
	if !qp.CanForfeit() {
		return ErrQuestNotForfeitable
	}
	snapshot := &Quest{
		QuestID:   qp.QuestID,
		Status:    QuestStatusNotStarted,
		Forfeited: qp.Forfeited + 1,
	}
	ch.Quests.Remove(qp.QuestID)
	ch.Listener.OnQuestForfeited(ch, snapshot)
	return nil
}

type QuestContainer struct {
	owner    *Character
	progress map[uint32]*Quest
}

func NewQuestContainer(owner *Character) *QuestContainer {
	return &QuestContainer{
		owner:    owner,
		progress: make(map[uint32]*Quest),
	}
}

func (qc *QuestContainer) Get(questID uint32) *Quest {
	if qc == nil {
		return nil
	}
	return qc.progress[questID]
}

func (qc *QuestContainer) questDef(questID uint32) *wz.Quest {
	if qc == nil || qc.owner == nil || qc.owner.GameWorld == nil {
		return nil
	}
	resources := qc.owner.GameWorld.GetResources()
	if resources == nil {
		return nil
	}
	return resources.GetQuest(questID)
}

func (qc *QuestContainer) Create(questID uint32, status QuestStatusType) *Quest {
	if qc == nil || questID == 0 {
		return nil
	}
	if qc.progress[questID] != nil {
		return nil
	}
	qp := &Quest{
		container: qc,
		Wz:        qc.questDef(questID),
		QuestID:   questID,
		Status:    status,
		MobKills:  make(map[uint32]int),
		Unknown2:  make(map[string]string),
	}
	qc.progress[questID] = qp
	return qp
}

func (qc *QuestContainer) Remove(questID uint32) {
	if qc == nil {
		return
	}
	delete(qc.progress, questID)
}

func (qc *QuestContainer) Clear(questID uint32) bool {
	if qc == nil {
		return false
	}
	existing := qc.progress[questID]
	if existing == nil {
		return false
	}
	delete(qc.progress, questID)
	if qc.owner != nil && qc.owner.Listener != nil && existing.WiresToClient() {
		qc.owner.Listener.OnQuestForfeited(qc.owner, &Quest{
			QuestID: questID,
			Status:  QuestStatusNotStarted,
		})
	}
	return true
}

func (qc *QuestContainer) ClearAll() int {
	if qc == nil {
		return 0
	}
	snapshots := make([]*Quest, 0, len(qc.progress))
	for id, qp := range qc.progress {
		if qp == nil {
			continue
		}
		snapshots = append(snapshots, &Quest{
			QuestID: id,
			Status:  QuestStatusNotStarted,
			Wz:      qp.Wz,
		})
	}
	count := len(snapshots)
	qc.progress = make(map[uint32]*Quest)
	if qc.owner == nil || qc.owner.Listener == nil {
		return count
	}
	for _, snap := range snapshots {
		if snap.WiresToClient() {
			qc.owner.Listener.OnQuestForfeited(qc.owner, &Quest{
				QuestID: snap.QuestID,
				Status:  QuestStatusNotStarted,
			})
		}
	}
	return count
}

func (qc *QuestContainer) ForEach(fn func(questID uint32, qp *Quest)) {
	if qc == nil {
		return
	}
	for id, qp := range qc.progress {
		if qp == nil {
			continue
		}
		fn(id, qp)
	}
}

func (qc *QuestContainer) Unknown2QuestInfo() map[uint16]string {
	result := make(map[uint16]string)
	qc.ForEach(func(questID uint32, qp *Quest) {
		if len(qp.Unknown2) == 0 {
			return
		}
		if questID > 0xFFFF {
			return
		}
		result[uint16(questID)] = serializeUnknown2(qp.Unknown2)
	})
	return result
}

func serializeUnknown2(fields map[string]string) string {
	if len(fields) == 0 {
		return ""
	}
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var b strings.Builder
	for i, key := range keys {
		if i > 0 {
			b.WriteByte(';')
		}
		b.WriteString(key)
		b.WriteByte('=')
		b.WriteString(fields[key])
	}
	return b.String()
}

func parseUnknown2(data string) map[string]string {
	if data == "" {
		return nil
	}
	fields := make(map[string]string)
	for _, part := range strings.Split(data, ";") {
		if part == "" {
			continue
		}
		key, value, ok := strings.Cut(part, "=")
		if !ok || key == "" {
			continue
		}
		fields[key] = value
	}
	if len(fields) == 0 {
		return nil
	}
	return fields
}
