package entity

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/boyism80/fm/core/clock"

	"github.com/boyism80/fm/services/game/wz"
)

type QuestStatusType = wz.QuestStatus

const (
	QuestStatusNotStarted = wz.QuestStatusNotStarted
	QuestStatusStarted    = wz.QuestStatusStarted
	QuestStatusCompleted  = wz.QuestStatusCompleted
)

type Quest struct {
	container      *QuestContainer
	Wz             *wz.Quest
	QuestID        uint32
	Status         QuestStatusType
	MobKills       map[uint32]int
	Deadline       time.Time
	StartTime      time.Time
	StatusRecord   StatusRecord
	RecordEx       map[string]string
	CompletionTime time.Time
	Forfeited      int
}

func (qp *Quest) IsStarted() bool {
	if qp == nil {
		return false
	}
	return qp.Status == QuestStatusStarted
}

func (qp *Quest) Expired() bool {
	if qp == nil || qp.Deadline.IsZero() {
		return false
	}
	return !clock.Now().Before(qp.Deadline)
}

func (qp *Quest) SetDeadline(t time.Time) {
	if qp == nil {
		return
	}
	qp.Deadline = t
}

func (qp *Quest) SetDeadlineAfter(d time.Duration) {
	if qp == nil || d <= 0 {
		return
	}
	qp.Deadline = clock.Now().Add(d)
}

func (qp *Quest) ResetDeadline() {
	if qp == nil {
		return
	}
	qp.Deadline = time.Time{}
}

func (qp *Quest) SetStartTime(t time.Time) {
	if qp == nil {
		return
	}
	qp.StartTime = t
}

func (qp *Quest) ResetStartTime() {
	if qp == nil {
		return
	}
	qp.StartTime = time.Time{}
}

func (qp *Quest) MatchesState(state wz.QuestStatus) bool {
	if qp == nil {
		return state == wz.QuestStatusNotStarted
	}
	return qp.Status == QuestStatusType(state)
}

func (qp *Quest) MeetsMobCounts(mobs map[uint32]int) bool {
	for mobID, count := range mobs {
		kills := 0
		if qp != nil && qp.MobKills != nil {
			kills = qp.MobKills[mobID]
		}
		if kills < count {
			return false
		}
	}
	return true
}

func (qp *Quest) CanComplete(ch *Character, opts QuestPhaseOpts) error {
	if qp == nil || ch == nil || ch.Quests == nil || qp.Wz == nil {
		return ErrQuestNotCompletable
	}
	if !qp.IsStarted() {
		return ErrQuestNotCompletable
	}
	if qp.Expired() {
		return ErrQuestExpired
	}
	if qp.Wz.Meta.Blocked {
		return ErrQuestNotCompletable
	}
	checkOpts := opts
	if qp.Wz.Meta.AutoPreComplete || qp.Wz.Meta.AutoComplete {
		checkOpts.NpcID = nil
	}
	if !phaseRequirementsMet(qp.Wz.Complete, qp.container, qp, checkOpts) {
		return ErrQuestNotCompletable
	}
	actions := ch.buildPhaseActions(qp.Wz.Complete, opts.Selection)
	if actions.Exchange.Valid(ch) != ExchangeOK {
		return ErrQuestNotCompletable
	}
	return nil
}

func (qp *Quest) StartedMobKills() []uint16 {
	if qp == nil || qp.Wz == nil || len(qp.Wz.OrderedMobIDs()) == 0 {
		return nil
	}
	return qp.MobKillCountsOrdered()
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
	spec := ExchangeSpec{
		Reward: ExchangeSide{
			Items: map[uint32]uint16{itemID: count},
		},
	}
	if ch.Exchange(spec) != ExchangeOK {
		return ErrQuestExchangeFailed
	}
	return nil
}

func (qp *Quest) Complete(ch *Character, opts QuestPhaseOpts) error {
	if qp == nil || ch == nil || ch.Quests == nil || qp.Wz == nil {
		return ErrQuestNotCompletable
	}
	if qp.Expired() {
		return ErrQuestExpired
	}
	wireNPC := uint32(0)
	if opts.NpcID != nil {
		wireNPC = *opts.NpcID
	}
	if !opts.Force {
		if err := qp.CanComplete(ch, opts); err != nil {
			return err
		}
		actions := ch.buildPhaseActions(qp.Wz.Complete, opts.Selection)
		if err := actions.Apply(qp, wireNPC, true); err != nil {
			return err
		}
	} else if !qp.IsStarted() {
		return ErrQuestNotCompletable
	}
	qp.Status = QuestStatusCompleted
	qp.CompletionTime = clock.Now()
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

func (qc *QuestContainer) wzDef(questID uint32) *wz.Quest {
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
		Wz:        qc.wzDef(questID),
		QuestID:   questID,
		Status:    status,
		MobKills:  make(map[uint32]int),
		RecordEx:  make(map[string]string),
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
	if qc.owner != nil && qc.owner.Listener != nil && existing.Wz != nil {
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
		if snap.Wz != nil {
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

func (qc *QuestContainer) RecordExWireMap() map[uint16]string {
	result := make(map[uint16]string)
	qc.ForEach(func(questID uint32, qp *Quest) {
		if len(qp.RecordEx) == 0 {
			return
		}
		if questID > 0xFFFF {
			return
		}
		result[uint16(questID)] = formatRecordEx(qp.RecordEx)
	})
	return result
}

func formatRecordEx(fields map[string]string) string {
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

func (qp *Quest) RecordExWire() string {
	if qp == nil {
		return ""
	}
	return formatRecordEx(qp.RecordEx)
}

func (qp *Quest) RecordExField(key string) (string, bool) {
	if qp == nil || key == "" || qp.Wz == nil {
		return "", false
	}
	if len(qp.RecordEx) == 0 {
		return "", false
	}
	value, ok := qp.RecordEx[key]
	return value, ok
}

func (qp *Quest) SetRecordExField(key, value string) bool {
	if qp == nil || key == "" || value == "" || qp.Wz == nil {
		return false
	}
	if qp.container == nil || qp.container.Get(qp.QuestID) == nil {
		return false
	}
	if qp.RecordEx == nil {
		qp.RecordEx = make(map[string]string)
	}
	qp.RecordEx[key] = value
	qp.notifyRecordExChanged()
	return true
}

func (qp *Quest) IncrementRecordExField(key string, delta int) bool {
	if qp == nil || key == "" || delta == 0 {
		return false
	}
	count := 0
	if val, ok := qp.RecordExField(key); ok {
		parsed, err := strconv.Atoi(val)
		if err != nil {
			return false
		}
		count = parsed
	}
	count += delta
	if count < 0 {
		count = 0
	}
	return qp.SetRecordExField(key, strconv.Itoa(count))
}

func (qp *Quest) notifyRecordExChanged() {
	if qp == nil || qp.container == nil {
		return
	}
	ch := qp.container.owner
	if ch == nil || ch.Listener == nil {
		return
	}
	ch.Listener.OnQuestRecordExChanged(ch, qp)
}
