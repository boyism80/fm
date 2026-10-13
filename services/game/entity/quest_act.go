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
	Class     *wz.Class
	Selection *uint32
	Forfeited bool
	NpcID     uint32
}

func (qc *Quests) grant(qp *Quest, actions wz.QuestActions, opts questActionOpts) error {
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
	if actions.BuffItemID != 0 && !opts.Forfeited {
		if !ch.RunQuestHook(qp.QuestID, "on_buff_item", actions.BuffItemID) {
			ch.UseItemEffect(actions.BuffItemID)
		}
	}
	if actions.PetTameness > 0 && ch.Pets.Active != nil && !opts.Forfeited {
		ch.Pets.Active.AddCloseness(actions.PetTameness)
	}
	if actions.PetSpeed > 0 && ch.Pets.Active != nil && !opts.Forfeited {
		ch.Pets.Active.Item.Speed = uint16(actions.PetSpeed)
		ch.Listener.OnPetUpdated(ch, ch.Pets.Active.Item)
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

func (qc *Quests) buildPhaseExchange(actions wz.QuestActions, opts questActionOpts) ExchangeSpec {
	spec := ExchangeSpec{}
	if qc == nil {
		return spec
	}
	if !opts.Forfeited {
		if actions.Money > 0 {
			spec.Reward.Meso += int32(actions.Money)
		} else if actions.Money < 0 {
			spec.Cost.Meso += int32(-actions.Money)
		}
		if actions.Exp > 0 {
			spec.Reward.Exp += uint32(actions.Exp)
		}
		if actions.Pop > 0 {
			spec.Reward.Population += int32(actions.Pop)
		} else if actions.Pop < 0 {
			spec.Cost.Population += int32(-actions.Pop)
		}
		if actions.SP > 0 {
			spec.Reward.SkillPoint += uint16(actions.SP)
		}
	}
	if len(actions.Item) > 0 {
		qc.appendPhaseActItems(&spec, opts.Class, actions.Item, opts.Selection)
	}
	if len(actions.Skills) == 0 {
		return spec
	}
	if !opts.Class.InQuestClasses(actions.SkillClasses) {
		return spec
	}
	for _, skill := range actions.Skills {
		if skill.SkillID == 0 {
			continue
		}
		if !opts.Class.InQuestClasses(skill.Classes) {
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

func (qc *Quests) Start(questID uint32, opts QuestPhaseOpts) (*Quest, error) {
	if qc == nil || questID == 0 || qc.owner == nil {
		return nil, ErrQuestNotStartable
	}
	if !opts.Force {
		if err := qc.CanStart(questID, opts); err != nil {
			return nil, err
		}
	}

	qp := qc.Get(questID)
	created := qp == nil
	if created {
		qp = qc.Create(questID, QuestStatusStarted)
	}
	if qp == nil {
		return nil, ErrQuestNotStartable
	}

	wireNPC := uint32(0)
	if opts.NpcID != nil {
		wireNPC = *opts.NpcID
	}
	if qp.IsStarted() && !created {
		qc.notifyQuestStart(qp, wireNPC, opts)
		return qp, nil
	}

	before := *qp
	qp.Status = QuestStatusStarted
	qp.StatusRecord.WriteString("")
	if qp.Wz.Meta.TimeLimit2 > 0 {
		qp.SetDeadline(clock.Now().Add(time.Duration(qp.Wz.Meta.TimeLimit2) * time.Second))
	} else {
		qp.ResetDeadline()
	}
	qp.MobKills = make(map[uint32]int)
	qp.InitMobKillCounters()
	if !opts.Force {
		err := qc.grant(qp, qp.Wz.Start.Actions, questActionOpts{
			Class:     qc.owner.Class,
			Forfeited: qp.Forfeited > 0,
			NpcID:     wireNPC,
		})
		if err != nil {
			if created {
				qc.Remove(questID)
			} else {
				*qp = before
			}
			return nil, err
		}
	}
	qc.notifyQuestStart(qp, wireNPC, opts)
	qc.RunAutoTriggers(nil, AutoQuestTriggerInfoStart, questID)
	return qp, nil
}

func (qc *Quests) notifyQuestStart(qp *Quest, wireNPC uint32, opts QuestPhaseOpts) {
	if qp == nil || qc.owner == nil {
		return
	}
	ch := qc.owner
	ch.Listener.OnQuestStarted(ch, qp, wireNPC)
	if opts.Record == nil || *opts.Record == "" {
		return
	}
	qp.StatusRecord.WriteString(*opts.Record)
	ch.Listener.OnQuestProgress(ch, qp)
}

func (qc *Quests) CanStart(questID uint32, opts QuestPhaseOpts) error {
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
		case QuestStatusNotStarted:
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
	if !qc.requirementsMet(def.Start.Requirements, existing, checkOpts) {
		return ErrQuestNotStartable
	}
	exchange := qc.buildPhaseExchange(def.Start.Actions, questActionOpts{
		Class:     qc.owner.Class,
		Forfeited: existing != nil && existing.Forfeited > 0,
	})
	if exchange.Cost.ValidCost(qc.owner) != ExchangeOK {
		return ErrQuestNotStartable
	}
	return nil
}

func (qc *Quests) broadcastNpcAct(npcID uint32, act string) {
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

func (qc *Quests) updateLinkedQuests(quests map[uint32]wz.QuestStatus) {
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
				qc.Create(questID, QuestStatusStarted)
			} else {
				existing.Status = QuestStatusStarted
			}
		case wz.QuestStatusCompleted:
			existing := qc.Get(questID)
			if existing == nil {
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

func (qc *Quests) completeLinkedQuest(refID uint32) {
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

func (qc *Quests) appendPhaseActItems(spec *ExchangeSpec, class *wz.Class, items []wz.QuestActionItem, selection *uint32) {
	if qc == nil || spec == nil {
		return
	}

	var randomPool []uint32
	for _, item := range items {
		if !qc.matchesRewardItem(item, class) || item.Count <= 0 || item.Prop <= 0 {
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
		if !qc.matchesRewardItem(item, class) || item.Count == 0 {
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
		case item.Prop == wz.QuestRewardPropAlways:
		case item.Prop == wz.QuestRewardPropSelect:
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
		if item.Period > 0 {
			if spec.Reward.Period == nil {
				spec.Reward.Period = make(map[uint32]time.Duration)
			}
			spec.Reward.Period[item.ItemID] = time.Duration(item.Period) * time.Minute
		}
	}
}

func (qc *Quests) completedQuestCount() int {
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

func (qc *Quests) matchesRewardItem(item wz.QuestRewardItem, class *wz.Class) bool {
	if qc == nil || qc.owner == nil {
		return false
	}
	if item.Gender <= 1 && item.Gender != int(qc.owner.look.Gender) {
		return false
	}
	if item.Classes == nil {
		return true
	}
	for _, code := range item.Classes {
		if code/100 == int(class.ID)/100 {
			return true
		}
	}
	return false
}
