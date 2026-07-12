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

type QuestPhaseOpts struct {
	NpcID     *uint32
	Selection *uint32
	Force     bool
	Record    *string
}

type questActionOpts struct {
	ClassID       uint16
	Selection     *uint32
	IncludeSkills bool
	NpcID         uint32
}

func (qc *QuestContainer) applyActions(qp *Quest, actions wz.QuestActions, opts questActionOpts) error {
	if qp == nil {
		return nil
	}
	if qc == nil {
		return nil
	}
	ch := qc.owner
	exchange := qc.buildPhaseExchange(actions, opts)
	if ch != nil && ch.Exchange(exchange) != ExchangeOK {
		return ErrQuestExchangeFailed
	}
	if actions.Info != "" {
		qp.StatusRecord.WriteString(actions.Info)
	}
	if len(actions.Quests) > 0 {
		qc.updateLinkedQuests(actions.Quests)
	}
	if actions.InfoNumber > 0 {
		qc.completeLinkedQuest(actions.InfoNumber)
	}
	if actions.NPCAct != "" {
		qc.broadcastNpcAct(opts.NpcID, actions.NPCAct)
	}
	return nil
}

func (qc *QuestContainer) buildPhaseExchange(actions wz.QuestActions, opts questActionOpts) ExchangeSpec {
	spec := ExchangeSpec{}
	if qc == nil {
		return spec
	}
	if actions.Money > 0 {
		spec.Reward.Meso += int32(actions.Money)
	} else if actions.Money < 0 {
		spec.Cost.Meso += int32(-actions.Money)
	}
	if actions.Exp > 0 {
		spec.Reward.Exp += uint32(actions.Exp)
	}
	if len(actions.Item) > 0 {
		qc.appendPhaseActItems(&spec, opts.ClassID, actions.Item, opts.Selection)
	}
	if actions.Pop > 0 {
		spec.Reward.Population += int32(actions.Pop)
	} else if actions.Pop < 0 {
		spec.Cost.Population += int32(-actions.Pop)
	}
	if !opts.IncludeSkills || len(actions.Skills) == 0 {
		return spec
	}
	if len(actions.SkillJobs) > 0 && !matchesQuestClass(opts.ClassID, actions.SkillJobs) {
		return spec
	}
	for _, skill := range actions.Skills {
		if skill.SkillID == 0 {
			continue
		}
		if len(skill.Classes) > 0 && !matchesQuestClass(opts.ClassID, skill.Classes) {
			continue
		}
		spec.Reward.Skills = append(spec.Reward.Skills, ExchangeSkill{
			SkillID:     skill.SkillID,
			Level:       skill.SkillLevel,
			MasterLevel: skill.MasterLevel,
		})
	}
	return spec
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
	if !opts.Force {
		if err := qc.CanStart(def.ID, opts); err != nil {
			return nil, err
		}
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
				if err := qc.applyActions(raced, def.Start.Actions, questActionOpts{
					ClassID:       ch.Class,
					IncludeSkills: raced.Forfeited == 0,
					NpcID:         wireNPC,
				}); err != nil {
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
		if err := qc.applyActions(qp, def.Start.Actions, questActionOpts{
			ClassID:       ch.Class,
			IncludeSkills: forfeited == 0,
			NpcID:         wireNPC,
		}); err != nil {
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
	exchange := qc.buildPhaseExchange(def.Start.Actions, questActionOpts{ClassID: qc.owner.Class})
	if exchange.Cost.ValidCost(qc.owner) != ExchangeOK {
		return ErrQuestNotStartable
	}
	return nil
}

func (qc *QuestContainer) broadcastNpcAct(npcID uint32, act string) {
	if qc == nil || qc.owner == nil || npcID == 0 || act == "" {
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
		npc.ShowEffect(act)
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

func (qc *QuestContainer) appendPhaseActItems(spec *ExchangeSpec, classID uint16, items []wz.QuestActionItem, selection *uint32) {
	if qc == nil || spec == nil {
		return
	}

	var randomPool []uint32
	for _, item := range items {
		if !qc.matchesRewardItem(item, classID) || item.Count <= 0 || !item.Prop.IsWeightedRandom() {
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
		if !qc.matchesRewardItem(item, classID) || item.Count == 0 {
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

func (qc *QuestContainer) matchesRewardItem(item wz.QuestRewardItem, classID uint16) bool {
	if qc == nil || qc.owner == nil {
		return false
	}
	if item.Gender <= 1 && item.Gender != int(qc.owner.gender) {
		return false
	}
	if item.Class <= 0 {
		return true
	}
	classCode := int(classID)
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
