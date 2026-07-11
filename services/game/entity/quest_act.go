package entity

import (
	"errors"
	"math/rand"
	"time"

	"github.com/boyism80/fm/core/clock"

	"github.com/boyism80/fm/services/game/wz"
)

var (
	ErrQuestNotFound       = errors.New("quest not found")
	ErrQuestInvalidState   = errors.New("quest invalid state")
	ErrQuestNotStartable   = errors.New("quest not startable")
	ErrQuestNotCompletable = errors.New("quest not completable")
	ErrQuestNotForfeitable = errors.New("quest not forfeitable")
	ErrQuestExpired        = errors.New("quest expired")
	ErrQuestRestoreItem    = errors.New("quest restore item unavailable")
	ErrQuestExchangeFailed = errors.New("quest exchange failed")
)

type questPhaseActions struct {
	Exchange               ExchangeSpec
	info                   string
	linkedQuests           map[uint32]wz.QuestStatus
	linkedQuestsToComplete []uint32
	skillGrants            []wz.QuestRewardSkill
	npcActs                []string
}

func (a questPhaseActions) Apply(qp *Quest, npcID uint32, grantSkills bool) error {
	if qp == nil {
		return nil
	}
	qc := qp.container
	if qc == nil {
		return nil
	}
	ch := qc.owner
	if ch != nil && ch.Exchange(a.Exchange) != ExchangeOK {
		return ErrQuestExchangeFailed
	}
	if a.info != "" {
		qp.StatusRecord.WriteString(a.info)
	}
	if len(a.linkedQuests) > 0 {
		qc.updateLinkedQuests(a.linkedQuests)
	}
	for _, refID := range a.linkedQuestsToComplete {
		qc.completeLinkedQuest(refID)
	}
	if grantSkills {
		qc.grantSkillRewards(a.skillGrants)
	}
	qc.broadcastNpcActs(npcID, a.npcActs)
	return nil
}

type QuestPhaseOpts struct {
	NpcID     *uint32
	Selection *uint32
	Force     bool
	Record    *string
}

func (qc *QuestContainer) Start(questID uint32, opts QuestPhaseOpts) (*Quest, error) {
	if qc == nil || questID == 0 || qc.owner == nil {
		return nil, ErrQuestNotStartable
	}
	def := qc.wzDef(questID)
	if def == nil {
		if !opts.Force {
			return nil, ErrQuestNotStartable
		}
		qp := qc.Get(questID)
		if qp == nil {
			qp = qc.Create(questID, QuestStatusStarted)
			if qp == nil {
				qp = qc.Get(questID)
			}
		}
		if qp == nil {
			return nil, ErrQuestInvalidState
		}
		qp.Status = QuestStatusStarted
		if opts.Record != nil {
			qp.StatusRecord.WriteString(*opts.Record)
		}
		qp.ResetDeadline()
		if qp.MobKills == nil {
			qp.MobKills = make(map[uint32]int)
		}
		qc.RunAutoTriggers(nil, AutoQuestTriggerLogin, 0)
		qc.RunAutoTriggers(nil, AutoQuestTriggerLevelUp, 0)
		return qp, nil
	} else {
		wireNPC := uint32(0)
		if opts.NpcID != nil {
			wireNPC = *opts.NpcID
		}
		existing := qc.Get(def.ID)
		if existing != nil && existing.IsStarted() {
			if !opts.Force {
				return nil, ErrQuestNotStartable
			}
			qc.notifyQuestStart(existing, wireNPC, opts)
			return existing, nil
		}
		return qc.startWZQuest(def, existing, wireNPC, opts)
	}
}

func (qc *QuestContainer) startWZQuest(def *wz.Quest, existing *Quest, wireNPC uint32, opts QuestPhaseOpts) (*Quest, error) {
	ch := qc.owner
	var actions questPhaseActions
	if !opts.Force {
		if err := qc.CanStart(def.ID, opts); err != nil {
			return nil, err
		}
		actions = ch.buildPhaseActions(def.Start, nil)
	}

	forfeited := 0
	if existing != nil {
		if existing.Status == QuestStatusCompleted {
			if !def.Meta.Repeatable && !opts.Force {
				return nil, ErrQuestInvalidState
			}
			forfeited = existing.Forfeited
		} else if !opts.Force {
			return nil, ErrQuestInvalidState
		} else {
			forfeited = existing.Forfeited
		}
	}

	qp := existing
	if qp == nil {
		qp = qc.Create(def.ID, QuestStatusStarted)
		if qp == nil {
			raced := qc.Get(def.ID)
			if raced == nil || !raced.IsStarted() {
				return nil, ErrQuestInvalidState
			}
			if !opts.Force {
				if err := qc.CanStart(def.ID, opts); err != nil {
					return nil, err
				}
				actions = ch.buildPhaseActions(def.Start, nil)
				if err := actions.Apply(raced, wireNPC, raced.Forfeited == 0); err != nil {
					return nil, err
				}
			}
			qc.notifyQuestStart(raced, wireNPC, opts)
			return raced, nil
		}
	}

	qp.Status = QuestStatusStarted
	qp.Forfeited = forfeited
	qp.StatusRecord.WriteString("")
	if def.Meta.TimeLimit2 > 0 {
		qp.SetDeadlineAfter(time.Duration(def.Meta.TimeLimit2) * time.Second)
	} else {
		qp.ResetDeadline()
	}
	qp.MobKills = make(map[uint32]int)
	qp.InitMobKillCounters()
	if !opts.Force {
		if err := actions.Apply(qp, wireNPC, forfeited == 0); err != nil {
			return nil, err
		}
	}
	qc.notifyQuestStart(qp, wireNPC, opts)
	return qp, nil
}

func (qc *QuestContainer) notifyQuestStart(qp *Quest, wireNPC uint32, opts QuestPhaseOpts) {
	if qp == nil || qc.owner == nil {
		return
	}
	ch := qc.owner
	if qp.Wz != nil && ch.Listener != nil {
		ch.Listener.OnQuestStarted(ch, qp, wireNPC)
	}
	if opts.Record == nil || *opts.Record == "" {
		return
	}
	qp.StatusRecord.WriteString(*opts.Record)
	if qp.Wz != nil && ch.Listener != nil {
		ch.Listener.OnQuestProgress(ch, qp)
	}
}

func (qc *QuestContainer) CanStart(questID uint32, opts QuestPhaseOpts) error {
	if qc == nil || qc.owner == nil {
		return ErrQuestNotStartable
	}
	def := qc.wzDef(questID)
	if def == nil {
		return ErrQuestNotStartable
	}
	existing := qc.Get(def.ID)
	if def.Meta.Blocked {
		return ErrQuestNotStartable
	}
	if existing != nil {
		switch existing.Status {
		case QuestStatusStarted:
			return ErrQuestNotStartable
		case QuestStatusCompleted:
			if !def.Meta.Repeatable {
				return ErrQuestNotStartable
			}
		default:
			return ErrQuestNotStartable
		}
	}
	checkOpts := opts
	if def.Meta.AutoStart || def.Meta.AutoAccept {
		checkOpts.NpcID = nil
	}
	if !phaseRequirementsMet(def.Start, qc, existing, checkOpts) {
		return ErrQuestNotStartable
	}
	actions := qc.owner.buildPhaseActions(def.Start, nil)
	if actions.Exchange.Valid(qc.owner) != ExchangeOK {
		return ErrQuestNotStartable
	}
	return nil
}

func (qc *QuestContainer) broadcastNpcActs(npcID uint32, acts []string) {
	if qc == nil || qc.owner == nil || npcID == 0 || len(acts) == 0 {
		return
	}
	mapInst := qc.owner.GetMap()
	if mapInst == nil {
		return
	}
	for _, obj := range mapInst.GetNpcs() {
		npc, ok := obj.(*Npc)
		if !ok || npc.Wz == nil || npc.Wz.BaseSpawn == nil {
			continue
		}
		if npc.Wz.ID != npcID {
			continue
		}
		for _, act := range acts {
			npc.ShowEffect(act)
		}
		return
	}
}

func (qc *QuestContainer) updateLinkedQuests(quests map[uint32]wz.QuestStatus) {
	if qc == nil || len(quests) == 0 {
		return
	}
	for questID, state := range quests {
		switch state {
		case wz.QuestStatusNotStarted:
			qc.Remove(questID)
		case wz.QuestStatusStarted:
			existing := qc.Get(questID)
			if existing == nil {
				if qc.wzDef(questID) == nil {
					continue
				}
				qc.Create(questID, QuestStatusStarted)
			} else {
				existing.Status = QuestStatusStarted
			}
		case wz.QuestStatusCompleted:
			existing := qc.Get(questID)
			if existing == nil {
				if qc.wzDef(questID) == nil {
					continue
				}
				created := qc.Create(questID, QuestStatusCompleted)
				if created != nil {
					created.CompletionTime = clock.Now()
				}
			} else {
				existing.Status = QuestStatusCompleted
				existing.CompletionTime = clock.Now()
			}
		}
	}
}

func (qc *QuestContainer) completeLinkedQuest(refID uint32) {
	if qc == nil || refID == 0 {
		return
	}
	refQP := qc.Get(refID)
	if refQP == nil || !refQP.IsStarted() {
		return
	}
	refQP.Status = QuestStatusCompleted
	refQP.CompletionTime = clock.Now()
}

func (ch *Character) buildPhaseActions(phase wz.QuestPhase, selection *uint32) questPhaseActions {
	actions := questPhaseActions{}
	if ch == nil {
		return actions
	}

	for _, act := range phase.Actions {
		if len(act.ApplicableClasses) > 0 && !matchesQuestClass(ch.Class, act.ApplicableClasses) {
			continue
		}
		switch act.Kind {
		case wz.QuestActMoney:
			if act.IntValue > 0 {
				actions.Exchange.Reward.Meso += int32(act.IntValue)
			} else if act.IntValue < 0 {
				actions.Exchange.Cost.Meso += int32(-act.IntValue)
			}
		case wz.QuestActEXP:
			if act.IntValue > 0 {
				actions.Exchange.Reward.Exp += uint32(act.IntValue)
			}
		case wz.QuestActItem:
			appendPhaseActItems(ch, &actions.Exchange, act.Items, selection)
		case wz.QuestActInfo:
			if act.StrValue != "" {
				actions.info = act.StrValue
			}
		case wz.QuestActQuest:
			if len(act.Quests) == 0 {
				break
			}
			if actions.linkedQuests == nil {
				actions.linkedQuests = make(map[uint32]wz.QuestStatus, len(act.Quests))
			}
			for questID, state := range act.Quests {
				actions.linkedQuests[questID] = state
			}
		case wz.QuestActPop:
			if act.IntValue > 0 {
				actions.Exchange.Reward.Population += int32(act.IntValue)
			} else if act.IntValue < 0 {
				actions.Exchange.Cost.Population += int32(-act.IntValue)
			}
		case wz.QuestActSkill:
			actions.skillGrants = append(actions.skillGrants, act.Skills...)
		case wz.QuestActNPCAct:
			if act.StrValue != "" {
				actions.npcActs = append(actions.npcActs, act.StrValue)
			}
		case wz.QuestActNextQuest,
			wz.QuestActBuffItemID, wz.QuestActSP:
		case wz.QuestActInfoNumber:
			if act.IntValue > 0 {
				actions.linkedQuestsToComplete = append(actions.linkedQuestsToComplete, uint32(act.IntValue))
			}
		}
	}

	return actions
}

func appendPhaseActItems(ch *Character, spec *ExchangeSpec, items []wz.QuestRewardItem, selection *uint32) {
	if spec == nil || ch == nil {
		return
	}

	var randomPool []uint32
	for _, item := range items {
		if ch.Quests == nil || !ch.Quests.matchesRewardItem(item) || item.Count <= 0 || !item.Prop.IsWeightedRandom() {
			continue
		}
		for i := 0; i < item.Prop.RandomWeight(); i++ {
			randomPool = append(randomPool, item.ItemID)
		}
	}

	randomPick := uint32(0)
	if len(randomPool) > 0 {
		randomPick = randomPool[rand.Intn(len(randomPool))]
	}

	selectionIdx := -1
	if selection != nil {
		selectionIdx = int(*selection)
	}

	extNum := 0
	for _, item := range items {
		if ch.Quests == nil || !ch.Quests.matchesRewardItem(item) || item.Count == 0 {
			continue
		}
		if item.Count < 0 {
			if spec.Cost.Items == nil {
				spec.Cost.Items = make(map[uint32]uint16)
			}
			spec.Cost.Items[item.ItemID] += uint16(-item.Count)
			continue
		}

		switch {
		case item.Prop.IsAlways():
		case item.Prop.IsSelection():
			if selectionIdx < 0 {
				continue
			}
			if selectionIdx != extNum {
				extNum++
				continue
			}
			extNum++
		default:
			if item.ItemID != randomPick {
				continue
			}
		}

		if spec.Reward.Items == nil {
			spec.Reward.Items = make(map[uint32]uint16)
		}
		spec.Reward.Items[item.ItemID] += uint16(item.Count)
	}
}

func matchesQuestClass(classID uint16, codes []int) bool {
	if len(codes) == 0 {
		return true
	}
	class := int(classID)
	for _, code := range codes {
		if code == 0 && class == 0 {
			return true
		}
		if code == class {
			return true
		}
		if code%100 == 0 && class/100 == code/100 {
			return true
		}
	}
	return false
}

func (ch *Character) meetsQuestSkillRequirement(skillID uint32, acquire int) bool {
	if ch == nil || skillID == 0 {
		return true
	}
	mustAcquire := acquire > 0
	var wzSkill *wz.Skill
	if ch.GameWorld != nil {
		resources := ch.GameWorld.GetResources()
		if resources != nil {
			wzSkill = resources.GetSkill(skillID)
		}
	}
	entry := (*SkillEntry)(nil)
	if ch.Skills != nil {
		entry = ch.Skills.Get(skillID)
	}
	skillLevel := 0
	masterLevel := 0
	if entry != nil {
		skillLevel = entry.Level()
		masterLevel = entry.MasterLevel
	}
	if mustAcquire {
		if wzSkill != nil && wzSkill.IsFourthJob() {
			return masterLevel > 0
		}
		return skillLevel > 0
	}
	return skillLevel == 0 && masterLevel == 0
}

func (qc *QuestContainer) completedQuestCount() int {
	if qc == nil {
		return 0
	}
	count := 0
	qc.ForEach(func(questID uint32, qp *Quest) {
		if qp == nil || qp.Status != QuestStatusCompleted {
			return
		}
		if questID > 99999 {
			return
		}
		count++
	})
	return count
}

func (qc *QuestContainer) matchesRewardItem(item wz.QuestRewardItem) bool {
	if qc == nil || qc.owner == nil {
		return false
	}
	ch := qc.owner
	if item.Gender <= 1 && item.Gender != int(ch.gender) {
		return false
	}
	if item.Class <= 0 {
		return true
	}
	classCode := int(ch.Class)
	for _, codec := range classesBy5ByteEncoding(item.Class) {
		if codec/100 == classCode/100 {
			return true
		}
	}
	if item.ClassEx > 0 {
		for _, codec := range classesBySimpleEncoding(item.ClassEx) {
			if (codec/100)%10 == (classCode/100)%10 {
				return true
			}
		}
	}
	return false
}

func classesBy5ByteEncoding(encoded int) []int {
	ret := make([]int, 0, 8)
	if encoded&0x1 != 0 {
		ret = append(ret, 0)
	}
	if encoded&0x2 != 0 {
		ret = append(ret, 100)
	}
	if encoded&0x4 != 0 {
		ret = append(ret, 200)
	}
	if encoded&0x8 != 0 {
		ret = append(ret, 300)
	}
	if encoded&0x10 != 0 {
		ret = append(ret, 400)
	}
	if encoded&0x20 != 0 {
		ret = append(ret, 500)
	}
	if encoded&0x400 != 0 {
		ret = append(ret, 1000)
	}
	if encoded&0x800 != 0 {
		ret = append(ret, 1100)
	}
	if encoded&0x1000 != 0 {
		ret = append(ret, 1200)
	}
	if encoded&0x2000 != 0 {
		ret = append(ret, 1300)
	}
	if encoded&0x4000 != 0 {
		ret = append(ret, 1400)
	}
	if encoded&0x8000 != 0 {
		ret = append(ret, 1500)
	}
	if encoded&0x20000 != 0 {
		ret = append(ret, 2001, 2200)
	}
	if encoded&0x100000 != 0 {
		ret = append(ret, 2000, 2001)
	}
	if encoded&0x200000 != 0 {
		ret = append(ret, 2100)
	}
	if encoded&0x400000 != 0 {
		ret = append(ret, 2001, 2200)
	}
	if encoded&0x40000000 != 0 {
		ret = append(ret, 3000, 3200, 3300, 3500)
	}
	return ret
}

func classesBySimpleEncoding(encoded int) []int {
	ret := make([]int, 0, 4)
	if encoded&0x1 != 0 {
		ret = append(ret, 200)
	}
	if encoded&0x2 != 0 {
		ret = append(ret, 300)
	}
	if encoded&0x4 != 0 {
		ret = append(ret, 400)
	}
	if encoded&0x8 != 0 {
		ret = append(ret, 500)
	}
	return ret
}

func (qc *QuestContainer) grantSkillRewards(rewards []wz.QuestRewardSkill) {
	if qc == nil || len(rewards) == 0 {
		return
	}
	for _, reward := range rewards {
		qc.grantSkill(reward)
	}
}

func (qc *QuestContainer) grantSkill(reward wz.QuestRewardSkill) {
	if qc == nil || qc.owner == nil || reward.SkillID == 0 {
		return
	}
	if !qc.canGrantSkill(reward) {
		return
	}
	ch := qc.owner
	if ch.GameWorld == nil {
		return
	}
	resources := ch.GameWorld.GetResources()
	if resources == nil {
		return
	}
	wzSkill := resources.GetSkill(reward.SkillID)
	if wzSkill == nil {
		return
	}
	targetLevel := skillTargetLevel(0, reward.SkillLevel)
	targetMaster := skillTargetMaster(0, reward.MasterLevel, wzSkill)
	entry := ch.Skills.Get(reward.SkillID)
	if entry == nil {
		entry = NewSkillEntry(ch, wzSkill, targetLevel, targetMaster)
		entry.Expiration = time.Time{}
		ch.Skills.Register(reward.SkillID, entry)
	} else {
		targetLevel = skillTargetLevel(entry.Level(), reward.SkillLevel)
		targetMaster = skillTargetMaster(entry.MasterLevel, reward.MasterLevel, wzSkill)
		if targetLevel != entry.Level() || targetMaster != entry.MasterLevel {
			entry.SetLevelAndMaster(targetLevel, targetMaster)
		}
	}
}

func (qc *QuestContainer) canGrantSkill(reward wz.QuestRewardSkill) bool {
	if qc == nil || qc.owner == nil {
		return false
	}
	ch := qc.owner
	if reward.SkillID/10000 == 0 && !ch.IsBeginner() {
		return false
	}
	return matchesQuestClass(ch.Class, reward.Classes)
}

func skillTargetLevel(current int, rewardLevel int) int {
	if rewardLevel <= 0 {
		if current > 0 {
			return current
		}
		return 1
	}
	if current > rewardLevel {
		return current
	}
	return rewardLevel
}

func skillTargetMaster(current int, rewardMaster int, wzSkill *wz.Skill) int {
	resolved := rewardMaster
	if resolved <= 0 && wzSkill != nil {
		if wzSkill.MasterLevel > 0 {
			resolved = wzSkill.MasterLevel
		} else if wzSkill.MaxLevel > 0 {
			resolved = wzSkill.MaxLevel
		}
	}
	if current > resolved {
		return current
	}
	return resolved
}
