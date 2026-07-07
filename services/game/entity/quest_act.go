package entity

import (
	"errors"
	"github.com/boyism80/fm/core/clock"
	"math/rand"
	"time"

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
	ErrQuestFlowFailed     = errors.New("quest flow failed")
)

type questPhaseMeta struct {
	info             string
	chainQuests      []wz.QuestStateRef
	infoNumberQuests []uint32
}

type QuestPrepareOpts struct {
	NpcID     *uint32
	Selection *uint32
	Force     bool
}

func (qc *QuestContainer) Start(questID uint32, opts QuestPrepareOpts) (*Quest, error) {
	if qc == nil || questID == 0 || qc.owner == nil {
		return nil, ErrQuestNotStartable
	}
	def := qc.questDef(questID)
	if def == nil {
		return nil, ErrQuestNotStartable
	}
	wireNPC := uint32(0)
	if opts.NpcID != nil {
		wireNPC = *opts.NpcID
	}

	var meta questPhaseMeta
	if !opts.Force {
		spec, m, err := qc.prepareStart(def, qc.owner, opts)
		if err != nil {
			return nil, err
		}
		meta = m
		if qc.owner.Exchange(spec) != FlowOK {
			return nil, ErrQuestFlowFailed
		}
	}

	qp := qc.Get(def.ID)
	if qp != nil && qp.IsStarted() {
		if !opts.Force {
			qp.ApplyPhaseMeta(meta, qc)
		}
		if qp.WiresToClient() {
			qc.owner.Listener.OnQuestStarted(qc.owner, qp, wireNPC)
		}
		return qp, nil
	}

	forfeited := 0
	if qp != nil {
		if qp.Status == QuestStatusCompleted {
			if !def.Meta.Repeatable && !opts.Force {
				return nil, ErrQuestInvalidState
			}
			forfeited = qp.Forfeited
		} else if !opts.Force {
			return nil, ErrQuestInvalidState
		} else {
			forfeited = qp.Forfeited
		}
	} else {
		qp = qc.Create(def.ID, QuestStatusStarted)
		if qp == nil {
			existing := qc.Get(def.ID)
			if existing != nil && existing.IsStarted() && !opts.Force {
				existing.ApplyPhaseMeta(meta, qc)
				if existing.WiresToClient() {
					qc.owner.Listener.OnQuestStarted(qc.owner, existing, wireNPC)
				}
				return existing, nil
			}
			return nil, ErrQuestInvalidState
		}
	}

	qp.Status = QuestStatusStarted
	qp.Forfeited = forfeited
	qp.StatusRecord.WriteString("")
	qp.ClearDeadline()
	qp.MobKills = make(map[uint32]int)
	qp.InitMobKillCounters()
	if !opts.Force {
		qp.ApplyPhaseMeta(meta, qc)
	}
	if qp.WiresToClient() {
		qc.owner.Listener.OnQuestStarted(qc.owner, qp, wireNPC)
	}
	return qp, nil
}

func (qc *QuestContainer) IsStartable(questID uint32, opts QuestPrepareOpts) bool {
	def := qc.questDef(questID)
	if def == nil {
		return false
	}
	_, _, err := qc.prepareStart(def, qc.owner, opts)
	return err == nil
}

func (qc *QuestContainer) prepareStart(
	def *wz.Quest,
	ch *Character,
	opts QuestPrepareOpts,
) (FlowSpec, questPhaseMeta, error) {
	if qc == nil || def == nil || ch == nil {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotStartable
	}
	if def.Meta.Blocked {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotStartable
	}
	existing := qc.Get(def.ID)
	if existing != nil {
		switch existing.Status {
		case QuestStatusStarted:
			return FlowSpec{}, questPhaseMeta{}, ErrQuestNotStartable
		case QuestStatusCompleted:
			if !def.Meta.Repeatable {
				return FlowSpec{}, questPhaseMeta{}, ErrQuestNotStartable
			}
		default:
			return FlowSpec{}, questPhaseMeta{}, ErrQuestNotStartable
		}
	}
	checkOpts := opts
	if def.Meta.AutoStart || def.Meta.AutoAccept {
		checkOpts.NpcID = nil
	}
	view := existing
	if view == nil {
		view = &Quest{
			container: qc,
			Wz:        def,
			QuestID:   def.ID,
			Status:    QuestStatusNotStarted,
		}
	}
	if !view.meetsPhaseRequirements(ch, def.Start, checkOpts) {
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
	opts QuestPrepareOpts,
) (FlowSpec, questPhaseMeta, error) {
	if qc == nil || ch == nil {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotCompletable
	}
	if qp == nil || qp.Wz == nil || !qp.IsStarted() {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotCompletable
	}
	if qp.IsDeadlineExpired() {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestExpired
	}
	if qp.Wz.Meta.Blocked {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotCompletable
	}
	checkOpts := opts
	if qp.Wz.Meta.AutoPreComplete || qp.Wz.Meta.AutoComplete {
		checkOpts.NpcID = nil
	}
	if !qp.meetsPhaseRequirements(ch, qp.Wz.Complete, checkOpts) {
		return FlowSpec{}, questPhaseMeta{}, ErrQuestNotCompletable
	}
	spec, meta := ch.buildPhaseFlow(qp.Wz.Complete, opts.Selection)
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
				if qc.questDef(ref.QuestID) == nil {
					continue
				}
				qc.Create(ref.QuestID, QuestStatusStarted)
			} else {
				existing.Status = QuestStatusStarted
			}
		case 2:
			existing := qc.Get(ref.QuestID)
			if existing == nil {
				if qc.questDef(ref.QuestID) == nil {
					continue
				}
				created := qc.Create(ref.QuestID, QuestStatusCompleted)
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

func (qp *Quest) meetsPhaseRequirements(
	ch *Character,
	phase wz.QuestPhase,
	opts QuestPrepareOpts,
) bool {
	if qp == nil || ch == nil {
		return false
	}
	if opts.Force {
		return true
	}
	for _, req := range phase.Requirements {
		if !qp.meetsRequirement(ch, phase, req, opts) {
			return false
		}
	}
	return true
}

func (qp *Quest) meetsRequirement(
	ch *Character,
	phase wz.QuestPhase,
	req wz.QuestRequirement,
	opts QuestPrepareOpts,
) bool {
	if qp == nil || ch == nil {
		return false
	}
	switch req.Kind {
	case wz.QuestReqNPC:
		if opts.NpcID == nil {
			return true
		}
		required := uint32(req.IntValue)
		if required == 0 {
			return true
		}
		wireNPC := *opts.NpcID
		if wireNPC != 0 && wireNPC != required {
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
			if npc.Wz.ID == required {
				return true
			}
		}
		return false
	case wz.QuestReqLvMin:
		return int(ch.GetLevel()) >= req.IntValue
	case wz.QuestReqLvMax:
		return int(ch.GetLevel()) <= req.IntValue
	case wz.QuestReqClass:
		return ch.matchesQuestClass(req.Classes)
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
		if qp.container == nil {
			return false
		}
		for _, ref := range req.Quests {
			if !qp.container.Get(ref.QuestID).MatchesState(ref.State) {
				return false
			}
		}
		return true
	case wz.QuestReqStartScript, wz.QuestReqEndScript:
		return true
	case wz.QuestReqPop:
		return int(ch.population) >= req.IntValue
	case wz.QuestReqFieldEnter:
		mapID := req.IntValue
		if mapID <= 0 {
			return true
		}
		mapInst := ch.GetMap()
		if mapInst == nil {
			return false
		}
		return int(mapInst.GetMapID()) == mapID
	case wz.QuestReqNormalAutoStart:
		return true
	case wz.QuestReqInterval:
		if qp.Status != QuestStatusCompleted {
			return true
		}
		if qp.CompletionTime.IsZero() {
			return true
		}
		minutes := req.IntValue
		if minutes <= 0 {
			return true
		}
		return !questRequirementNow().Before(qp.CompletionTime.Add(time.Duration(minutes) * time.Minute))
	case wz.QuestReqDayByDay:
		return qp.meetsDayByDayRequirement()
	case wz.QuestReqTimeStart, wz.QuestReqTimeEnd:
		return qp.meetsQuestEventTimeRequirement(req.Kind, req.StrValue)
	case wz.QuestReqInfoNumber:
		return qp.meetsInfoNumberRequirement(req.IntValue, phase)
	case wz.QuestReqInfo:
		return true
	case wz.QuestReqSkill, wz.QuestReqPet,
		wz.QuestReqPetTamenessMin, wz.QuestReqMBMin, wz.QuestReqMBCard,
		wz.QuestReqSubClassFlags,
		wz.QuestReqPartyQuestS, wz.QuestReqQuestComplete:
		return false
	default:
		if req.IntValue != 0 || req.StrValue != "" || len(req.Items) > 0 {
			return false
		}
		return true
	}
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
		if len(act.ApplicableClasses) > 0 && !ch.matchesQuestClass(act.ApplicableClasses) {
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
		case wz.QuestActPop:
			if act.IntValue > 0 {
				spec.Reward.Population += int32(act.IntValue)
			} else if act.IntValue < 0 {
				spec.Cost.Population += int32(-act.IntValue)
			}
		case wz.QuestActNextQuest, wz.QuestActSkill,
			wz.QuestActBuffItemID, wz.QuestActSP, wz.QuestActNPCAct:
		case wz.QuestActInfoNumber:
			if act.IntValue > 0 {
				meta.infoNumberQuests = append(meta.infoNumberQuests, uint32(act.IntValue))
			}
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

	var randomPool []uint32
	for _, item := range items {
		if !ch.matchesQuestRewardItem(item) || item.Count <= 0 || !item.Prop.IsWeightedRandom() {
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
		if !ch.matchesQuestRewardItem(item) || item.Count == 0 {
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

func (ch *Character) matchesQuestClass(classes []int) bool {
	if len(classes) == 0 {
		return true
	}
	classCode := int(ch.Class)
	for _, req := range classes {
		if classCode == req {
			return true
		}
		if req == 0 && classCode == 0 {
			return true
		}
		if req%100 == 0 && classCode/100 == req/100 {
			return true
		}
	}
	return false
}

func (ch *Character) matchesQuestRewardItem(item wz.QuestRewardItem) bool {
	if ch == nil {
		return false
	}
	if item.Gender <= 1 && item.Gender != int(ch.gender) {
		return false
	}
	if item.Class <= 0 {
		return true
	}
	classCode := int(ch.Class)
	for _, codec := range questClassesBy5ByteEncoding(item.Class) {
		if codec/100 == classCode/100 {
			return true
		}
	}
	if item.ClassEx > 0 {
		for _, codec := range questClassesBySimpleEncoding(item.ClassEx) {
			if (codec/100)%10 == (classCode/100)%10 {
				return true
			}
		}
	}
	return false
}

func questClassesBy5ByteEncoding(encoded int) []int {
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

func questClassesBySimpleEncoding(encoded int) []int {
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
