package game

import (
	"sort"
	"strings"
	"time"

	"github.com/boyism80/fm/core/clock"

	"github.com/boyism80/fm/maple/wz"
)

type QuestStatusType = wz.QuestStatus

const (
	QuestStatusNotStarted = wz.QuestStatusNotStarted
	QuestStatusStarted    = wz.QuestStatusStarted
	QuestStatusCompleted  = wz.QuestStatusCompleted
)

type Quest struct {
	quests         *Quests
	wz             *wz.Quest
	questID        uint32
	status         QuestStatusType
	mobKills       map[uint32]int
	deadline       time.Time
	startTime      time.Time
	statusRecord   StatusRecord
	recordEx       map[string]string
	completionTime time.Time
	forfeited      int
}

func (qp *Quest) IsStarted() bool {
	if qp == nil {
		return false
	}
	return qp.status == QuestStatusStarted
}

func (qp *Quest) Expired() bool {
	if qp == nil || qp.deadline.IsZero() {
		return false
	}
	return !clock.Now().Before(qp.deadline)
}

func (qp *Quest) SetDeadline(t time.Time) {
	if qp == nil {
		return
	}
	qp.deadline = t
}

func (qp *Quest) ResetDeadline() {
	if qp == nil {
		return
	}
	qp.deadline = time.Time{}
}

func (qp *Quest) SetStartTime(t time.Time) {
	if qp == nil {
		return
	}
	qp.startTime = t
}

func (qp *Quest) ResetStartTime() {
	if qp == nil {
		return
	}
	qp.startTime = time.Time{}
}

func (qp *Quest) MatchesState(state wz.QuestStatus) bool {
	if qp == nil {
		return state == wz.QuestStatusNotStarted
	}
	return qp.status == QuestStatusType(state)
}

func (qp *Quest) CanComplete(ch *Character, opts QuestPhaseOpts) error {
	if qp == nil || ch == nil {
		return ErrQuestNotCompletable
	}
	if !qp.IsStarted() {
		return ErrQuestNotCompletable
	}
	if qp.Expired() {
		return ErrQuestExpired
	}
	if qp.wz.Meta.Blocked {
		return ErrQuestNotCompletable
	}
	checkOpts := opts
	if qp.wz.Meta.AutoPreComplete || qp.wz.Meta.AutoComplete {
		checkOpts.NpcID = nil
	}
	if !qp.quests.passesRequirements(qp.wz.Complete.Requirements, qp, checkOpts) {
		return ErrQuestNotCompletable
	}
	exchange := ch.quests.buildPhaseExchange(qp.wz.Complete.Actions, questActionOpts{
		classID:   ch.class,
		selection: opts.Selection,
	})
	if exchange.Cost.ValidCost(ch) != ExchangeOK {
		return ErrQuestNotCompletable
	}
	return nil
}

func (qp *Quest) StartedMobKills() []uint16 {
	if qp == nil || len(qp.wz.MobIDs) == 0 {
		return nil
	}
	return qp.MobKillCountsOrdered()
}

func (qp *Quest) RecordMobKill(mobID uint32) bool {
	if qp == nil || !qp.IsStarted() || mobID == 0 {
		return false
	}
	required := qp.wz.RelevantMobs()
	need, ok := required[mobID]
	if !ok || need <= 0 {
		return false
	}
	if qp.mobKills == nil {
		qp.mobKills = make(map[uint32]int)
	}
	if qp.mobKills[mobID] >= need {
		return false
	}
	qp.mobKills[mobID]++
	return true
}

func (qp *Quest) MobKillCountsOrdered() []uint16 {
	if qp == nil {
		return nil
	}
	if len(qp.wz.MobIDs) == 0 {
		return nil
	}
	out := make([]uint16, len(qp.wz.MobIDs))
	for i, id := range qp.wz.MobIDs {
		kills := 0
		if qp.mobKills != nil {
			kills = qp.mobKills[id]
		}
		out[i] = uint16(kills)
	}
	return out
}

func (qc *Quests) RecordMobKill(mobID uint32) {
	if mobID == 0 {
		return
	}
	qc.ForEach(func(_ uint32, qp *Quest) {
		if qp == nil || !qp.RecordMobKill(mobID) {
			return
		}
		qc.owner.listener.OnQuestProgress(qc.owner, qp)
	})
}

func (qp *Quest) InitMobKillCounters() {
	if qp == nil {
		return
	}
	mobs := qp.wz.RelevantMobs()
	if len(mobs) == 0 {
		return
	}
	if qp.mobKills == nil {
		qp.mobKills = make(map[uint32]int)
	}
	for mobID := range mobs {
		if _, ok := qp.mobKills[mobID]; !ok {
			qp.mobKills[mobID] = 0
		}
	}
}

func (qp *Quest) CanForfeit() bool {
	if qp == nil || !qp.IsStarted() {
		return false
	}
	switch qp.questID {
	case 20000, 20010, 20015, 20020:
		return false
	default:
		return true
	}
}

func (qp *Quest) CanRestoreLostItem(ch *Character, itemID uint32) bool {
	if qp == nil || itemID == 0 || ch == nil {
		return false
	}
	if !qp.IsStarted() {
		return false
	}
	if ch.inventory.HasItem(itemID) {
		return false
	}
	for _, item := range qp.wz.Start.Actions.Item {
		if item.ItemID == itemID && item.Count > 0 {
			return true
		}
	}
	return false
}

func (qp *Quest) RestoreLostItem(ch *Character, itemID uint32) error {
	if !qp.CanRestoreLostItem(ch, itemID) {
		return ErrQuestRestoreItem
	}
	if ch == nil {
		return ErrQuestRestoreItem
	}
	var count uint16
	found := false
	for _, item := range qp.wz.Start.Actions.Item {
		if item.ItemID == itemID && item.Count > 0 {
			count = uint16(item.Count)
			found = true
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
	if qp == nil || ch == nil {
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
		if err := ch.quests.grant(qp, qp.wz.Complete.Actions, questActionOpts{
			classID:   ch.class,
			selection: opts.Selection,
			npcID:     wireNPC,
		}); err != nil {
			return err
		}
	} else if !qp.IsStarted() {
		return ErrQuestNotCompletable
	}
	qp.status = QuestStatusCompleted
	qp.completionTime = clock.Now()
	ch.listener.OnQuestCompleted(ch, qp, wireNPC, qp.wz.NextQuestID())
	return nil
}

func (qp *Quest) Forfeit(ch *Character) error {
	if qp == nil || ch == nil {
		return ErrQuestNotForfeitable
	}
	if !qp.CanForfeit() {
		return ErrQuestNotForfeitable
	}
	qp.status = QuestStatusNotStarted
	qp.forfeited++
	qp.statusRecord.WriteString("")
	qp.mobKills = make(map[uint32]int)
	qp.ResetDeadline()
	ch.listener.OnQuestForfeited(ch, qp)
	return nil
}

type Quests struct {
	owner    *Character
	progress map[uint32]*Quest
}

func NewQuests(owner *Character) *Quests {
	return &Quests{
		owner:    owner,
		progress: make(map[uint32]*Quest),
	}
}

func (qc *Quests) Get(questID uint32) *Quest {
	if qc == nil {
		return nil
	}
	return qc.progress[questID]
}

func (qc *Quests) wzDef(questID uint32) *wz.Quest {
	if qc == nil || qc.owner == nil || qc.owner.gameWorld == nil {
		return nil
	}
	resources := qc.owner.gameWorld.GetResources()
	if resources == nil {
		return nil
	}
	return resources.GetQuest(questID)
}

func (qc *Quests) Create(questID uint32, status QuestStatusType) *Quest {
	if qc == nil || questID == 0 {
		return nil
	}
	if qc.progress[questID] != nil {
		return nil
	}
	def := qc.wzDef(questID)
	if def == nil {
		return nil
	}
	qp := &Quest{
		quests:   qc,
		wz:       def,
		questID:  questID,
		status:   status,
		mobKills: make(map[uint32]int),
		recordEx: make(map[string]string),
	}
	qc.progress[questID] = qp
	return qp
}

func (qc *Quests) Remove(questID uint32) {
	if qc == nil {
		return
	}
	delete(qc.progress, questID)
}

func (qc *Quests) Clear(questID uint32) bool {
	if qc == nil {
		return false
	}
	existing := qc.progress[questID]
	if existing == nil {
		return false
	}
	delete(qc.progress, questID)
	if qc.owner != nil {
		qc.owner.listener.OnQuestForfeited(qc.owner, &Quest{
			questID: questID,
			status:  QuestStatusNotStarted,
		})
	}
	return true
}

func (qc *Quests) ClearAll() int {
	if qc == nil {
		return 0
	}
	snapshots := make([]*Quest, 0, len(qc.progress))
	for id, qp := range qc.progress {
		if qp == nil {
			continue
		}
		snapshots = append(snapshots, &Quest{
			questID: id,
			status:  QuestStatusNotStarted,
		})
	}
	count := len(snapshots)
	qc.progress = make(map[uint32]*Quest)
	if qc.owner == nil {
		return count
	}
	for _, snap := range snapshots {
		qc.owner.listener.OnQuestForfeited(qc.owner, snap)
	}
	return count
}

func (qc *Quests) ForEach(fn func(questID uint32, qp *Quest)) {
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

func (qc *Quests) CompletedCount() int {
	if qc == nil {
		return 0
	}
	n := 0
	qc.ForEach(func(_ uint32, qp *Quest) {
		if qp != nil && qp.status == QuestStatusCompleted {
			n++
		}
	})
	return n
}

func (qc *Quests) RecordExWireMap() map[uint16]string {
	result := make(map[uint16]string)
	qc.ForEach(func(questID uint32, qp *Quest) {
		if len(qp.recordEx) == 0 {
			return
		}
		if questID > 0xFFFF {
			return
		}
		result[uint16(questID)] = formatRecordEx(qp.recordEx)
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
	return formatRecordEx(qp.recordEx)
}

func (qp *Quest) RecordExField(key string) (string, bool) {
	if qp == nil || key == "" {
		return "", false
	}
	if len(qp.recordEx) == 0 {
		return "", false
	}
	value, ok := qp.recordEx[key]
	return value, ok
}

func (qp *Quest) SetRecordExField(key, value string) bool {
	if qp == nil || key == "" || value == "" {
		return false
	}
	if qp.quests == nil || qp.quests.Get(qp.questID) == nil {
		return false
	}
	if qp.recordEx == nil {
		qp.recordEx = make(map[string]string)
	}
	qp.recordEx[key] = value
	if qp.quests != nil {
		ch := qp.quests.owner
		if ch != nil && ch.listener != nil {
			ch.listener.OnQuestRecordExChanged(ch, qp)
		}
	}
	return true
}

func (qc *Quests) NeedsItem(questID uint32, itemID uint32) bool {
	qp := qc.Get(questID)
	if !qp.IsStarted() {
		return false
	}
	for _, act := range qp.wz.Complete.Actions.Item {
		if act.ItemID == itemID && act.Count < 0 {
			return !qc.owner.inventory.HasItemCount(itemID, uint16(-act.Count))
		}
	}
	return true
}
