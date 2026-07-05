package entity

import (
	"errors"
	"time"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

var (
	ErrQuestNotFound       = errors.New("quest not found")
	ErrQuestInvalidState   = errors.New("quest invalid state")
	ErrQuestNotStartable   = errors.New("quest not startable")
	ErrQuestNotCompletable = errors.New("quest not completable")
	ErrQuestNotForfeitable = errors.New("quest not forfeitable")
	ErrQuestRestoreItem    = errors.New("quest restore item unavailable")
	ErrQuestFlowFailed     = errors.New("quest flow failed")
)

type questPhaseMeta struct {
	info        string
	chainQuests []wz.QuestStateRef
}

type QuestPrepareOpts struct {
	IgnoreScriptRequirement bool
	SkipNPCRequirement      bool
	Force                   bool
}

func isQuestForfeitAllowed(questID uint32) bool {
	switch questID {
	case 20000, 20010, 20015, 20020:
		return false
	default:
		return true
	}
}

func (qc *QuestContainer) Start(def *wz.Quest, npcID uint32, opts QuestPrepareOpts) (*Quest, error) {
	ch := qc.owner
	if qc == nil || def == nil || ch == nil {
		return nil, ErrQuestNotStartable
	}

	var meta questPhaseMeta
	if !opts.Force {
		spec, m, err := qc.prepareStart(def, ch, npcID, opts)
		if err != nil {
			return nil, err
		}
		meta = m
		if ch.Exchange(spec) != FlowOK {
			return nil, ErrQuestFlowFailed
		}
	}

	qp := qc.Get(def.ID)
	if qp != nil && qp.IsStarted() {
		if !opts.Force {
			qp.ApplyPhaseMeta(meta, qc)
		}
		ch.Listener.OnQuestStarted(ch, qp, npcID)
		return qp, nil
	}

	forfeited := 0
	if qp != nil {
		if !opts.Force {
			return nil, ErrQuestInvalidState
		}
		forfeited = qp.Forfeited
	} else {
		qp = qc.Create(def, QuestStatusStarted)
		if qp == nil {
			existing := qc.Get(def.ID)
			if existing != nil && existing.IsStarted() && !opts.Force {
				existing.ApplyPhaseMeta(meta, qc)
				ch.Listener.OnQuestStarted(ch, existing, npcID)
				return existing, nil
			}
			return nil, ErrQuestInvalidState
		}
	}

	qp.Status = QuestStatusStarted
	qp.Forfeited = forfeited
	qp.InitMobKillCounters()
	if !opts.Force {
		qp.ApplyPhaseMeta(meta, qc)
	}
	ch.Listener.OnQuestStarted(ch, qp, npcID)
	return qp, nil
}

func (qc *QuestContainer) IsStartable(def *wz.Quest, npcID uint32, opts QuestPrepareOpts) bool {
	_, _, err := qc.prepareStart(def, qc.owner, npcID, opts)
	return err == nil
}

func (qc *QuestContainer) prepareStart(
	def *wz.Quest,
	ch *Character,
	npcID uint32,
	opts QuestPrepareOpts,
) (FlowSpec, questPhaseMeta, error) {
	if qc == nil || def == nil || ch == nil {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotStartable
	}
	if !opts.IgnoreScriptRequirement && def.HasStartScript() {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotStartable
	}
	if def.Meta.Blocked {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotStartable
	}
	existing := qc.Get(def.ID)
	if existing != nil && existing.Status != QuestStatusNotStarted {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotStartable
	}
	skipNPC := def.Meta.AutoStart || opts.SkipNPCRequirement
	if !meetsPhaseRequirements(def.Start, ch, existing, qc, npcID, skipNPC) {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotStartable
	}
	spec, meta := ch.buildPhaseFlow(def.Start, nil)
	if ch.ValidateFlow(spec) != FlowOK {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotStartable
	}
	return spec, meta, nil
}

func (qc *QuestContainer) prepareComplete(
	qp *Quest,
	ch *Character,
	npcID uint32,
	selection *uint32,
	opts QuestPrepareOpts,
) (FlowSpec, questPhaseMeta, error) {
	if qc == nil || ch == nil {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotCompletable
	}
	if qp == nil || qp.Wz == nil || !qp.IsStarted() {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotCompletable
	}
	def := qp.Wz
	if !opts.IgnoreScriptRequirement && def.HasEndScript() {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotCompletable
	}
	if def.Meta.Blocked {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotCompletable
	}
	skipNPC := def.Meta.AutoPreComplete || def.Meta.AutoComplete || opts.SkipNPCRequirement
	if !meetsPhaseRequirements(def.Complete, ch, qp, qc, npcID, skipNPC) {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotCompletable
	}
	spec, meta := ch.buildPhaseFlow(def.Complete, selection)
	if ch.ValidateFlow(spec) != FlowOK {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotCompletable
	}
	return spec, meta, nil
}

func (qc *QuestContainer) applyQuestChainActions(refs []wz.QuestStateRef) {
	if qc == nil {
		return
	}
	for _, ref := range refs {
		switch ref.State {
		case 0:
			qc.Remove(ref.QuestID)
		case 1:
			existing := qc.Get(ref.QuestID)
			if existing == nil {
				def := qc.questDef(ref.QuestID)
				if def != nil {
					qc.Create(def, QuestStatusStarted)
				}
			} else {
				existing.Status = QuestStatusStarted
			}
		case 2:
			existing := qc.Get(ref.QuestID)
			if existing == nil {
				def := qc.questDef(ref.QuestID)
				if def == nil {
					continue
				}
				created := qc.Create(def, QuestStatusCompleted)
				if created != nil {
					created.CompletionTime = time.Now()
				}
			} else {
				existing.Status = QuestStatusCompleted
				existing.CompletionTime = time.Now()
			}
		}
	}
}

func meetsPhaseRequirements(
	phase wz.QuestPhase,
	ch *Character,
	qp *Quest,
	qc *QuestContainer,
	npcID uint32,
	skipNPC bool,
) bool {
	if ch == nil {
		return false
	}
	for _, req := range phase.Requirements {
		if !meetsRequirement(req, ch, qp, qc, npcID, skipNPC) {
			return false
		}
	}
	return true
}

func meetsRequirement(
	req wz.QuestRequirement,
	ch *Character,
	qp *Quest,
	qc *QuestContainer,
	npcID uint32,
	skipNPC bool,
) bool {
	switch req.Kind {
	case wz.QuestReqNPC:
		if skipNPC {
			return true
		}
		required := uint32(req.IntValue)
		if required == 0 {
			return true
		}
		if npcID != 0 && npcID != required {
			return false
		}
		return ch.mapHasNpcTemplate(required)
	case wz.QuestReqLvMin:
		return int(ch.GetLevel()) >= req.IntValue
	case wz.QuestReqLvMax:
		return int(ch.GetLevel()) <= req.IntValue
	case wz.QuestReqJob:
		return ch.matchesQuestJob(req.Jobs)
	case wz.QuestReqItem:
		for _, item := range req.Items {
			if !ch.HasItemCount(item.ItemID, uint16(item.Count)) {
				return false
			}
		}
		return true
	case wz.QuestReqMob:
		return qp.MeetsMobCounts(req.Mobs)
	case wz.QuestReqQuest:
		if qc == nil {
			return false
		}
		for _, ref := range req.Quests {
			if !qc.Get(ref.QuestID).MatchesState(ref.State) {
				return false
			}
		}
		return true
	case wz.QuestReqStartScript, wz.QuestReqEndScript:
		return true
	case wz.QuestReqPop:
		return int(ch.famePoint) >= req.IntValue
	case wz.QuestReqInterval, wz.QuestReqSkill, wz.QuestReqPet,
		wz.QuestReqPetTamenessMin, wz.QuestReqMBMin, wz.QuestReqMBCard,
		wz.QuestReqSubJobFlags, wz.QuestReqDayByDay, wz.QuestReqNormalAutoStart,
		wz.QuestReqPartyQuestS, wz.QuestReqFieldEnter, wz.QuestReqQuestComplete,
		wz.QuestReqTimeStart, wz.QuestReqTimeEnd:
		return false
	default:
		if req.IntValue != 0 || req.StrValue != "" || len(req.Items) > 0 {
			return false
		}
		return true
	}
}

func startItemGrant(def *wz.Quest, itemID uint32) (uint16, bool) {
	if def == nil {
		return 0, false
	}
	for _, act := range def.Start.Actions {
		if act.Kind != wz.QuestActItem {
			continue
		}
		for _, item := range act.Items {
			if item.ItemID == itemID && item.Count > 0 {
				return uint16(item.Count), true
			}
		}
	}
	return 0, false
}

func (ch *Character) buildPhaseFlow(
	phase wz.QuestPhase,
	selection *uint32,
) (FlowSpec, questPhaseMeta) {
	spec := FlowSpec{}
	meta := questPhaseMeta{}
	if ch == nil {
		return spec, meta
	}

	for _, act := range phase.Actions {
		if len(act.ApplicableJobs) > 0 && !ch.matchesQuestJob(act.ApplicableJobs) {
			continue
		}
		switch act.Kind {
		case wz.QuestActMoney:
			if act.IntValue > 0 {
				spec.Reward.Meso += int32(act.IntValue)
			} else if act.IntValue < 0 {
				spec.Cost.Meso += int32(-act.IntValue)
			}
		case wz.QuestActEXP:
			if act.IntValue > 0 {
				spec.Reward.Exp += uint32(act.IntValue)
			}
		case wz.QuestActItem:
			appendPhaseActItems(ch, &spec, act.Items, selection)
		case wz.QuestActInfo:
			if act.StrValue != "" {
				meta.info = act.StrValue
			}
		case wz.QuestActQuest:
			meta.chainQuests = append(meta.chainQuests, act.Quests...)
		case wz.QuestActNextQuest, wz.QuestActSkill, wz.QuestActPop,
			wz.QuestActBuffItemID, wz.QuestActInfoNumber, wz.QuestActSP, wz.QuestActNPCAct:
		}
	}

	return spec, meta
}

func appendPhaseActItems(
	ch *Character,
	spec *FlowSpec,
	items []wz.QuestRewardItem,
	selection *uint32,
) {
	if spec == nil || ch == nil {
		return
	}

	var propChoices []wz.QuestRewardItem
	var equipRewards []wz.QuestRewardItem
	for _, item := range items {
		if !questItemMatchesCharacter(ch, item) {
			continue
		}
		if item.Count == 0 {
			continue
		}
		if item.Count < 0 {
			if spec.Cost.Items == nil {
				spec.Cost.Items = make(map[uint32]uint16)
			}
			spec.Cost.Items[item.ItemID] += uint16(-item.Count)
			continue
		}
		if item.Prop >= 0 {
			propChoices = append(propChoices, item)
			continue
		}
		if constant.GetInventoryTypeByItemID(item.ItemID) == constant.InventoryTypeEquipment {
			equipRewards = append(equipRewards, item)
			continue
		}
		if spec.Reward.Items == nil {
			spec.Reward.Items = make(map[uint32]uint16)
		}
		spec.Reward.Items[item.ItemID] += uint16(item.Count)
	}

	if len(propChoices) > 0 {
		idx := 0
		if selection != nil {
			idx = int(*selection)
		}
		if idx >= 0 && idx < len(propChoices) {
			item := propChoices[idx]
			if spec.Reward.Items == nil {
				spec.Reward.Items = make(map[uint32]uint16)
			}
			spec.Reward.Items[item.ItemID] += uint16(item.Count)
		}
	}

	if len(equipRewards) == 0 {
		return
	}
	if len(equipRewards) == 1 {
		item := equipRewards[0]
		if spec.Reward.Items == nil {
			spec.Reward.Items = make(map[uint32]uint16)
		}
		spec.Reward.Items[item.ItemID] += uint16(item.Count)
		return
	}
	if selection == nil {
		return
	}
	idx := int(*selection)
	if idx < 0 || idx >= len(equipRewards) {
		return
	}
	item := equipRewards[idx]
	if spec.Reward.Items == nil {
		spec.Reward.Items = make(map[uint32]uint16)
	}
	spec.Reward.Items[item.ItemID] += uint16(item.Count)
}

func (ch *Character) mapHasNpcTemplate(npcTemplateID uint32) bool {
	if ch == nil || npcTemplateID == 0 {
		return false
	}
	mapInst := ch.GetMap()
	if mapInst == nil {
		return false
	}
	for _, obj := range mapInst.GetNpcs() {
		npc, ok := obj.(*Npc)
		if !ok || npc.Wz == nil || npc.Wz.BaseSpawn == nil {
			continue
		}
		if npc.Wz.ID == npcTemplateID {
			return true
		}
	}
	return false
}

func (ch *Character) matchesQuestJob(jobs []int) bool {
	if len(jobs) == 0 {
		return true
	}
	job := int(ch.Class)
	for _, req := range jobs {
		if job == req {
			return true
		}
		if req == 0 && job == 0 {
			return true
		}
		if req%100 == 0 && job/100 == req/100 {
			return true
		}
	}
	return false
}

func questItemMatchesCharacter(ch *Character, item wz.QuestRewardItem) bool {
	if ch == nil {
		return false
	}
	if item.Job >= 0 && !ch.matchesQuestJob([]int{item.Job}) {
		return false
	}
	if item.Gender <= 1 && item.Gender != int(ch.gender) {
		return false
	}
	return true
}
