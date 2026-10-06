package entity

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/wz"
)

type AutoQuestTrigger int

const (
	AutoQuestTriggerMapEnter AutoQuestTrigger = iota
	AutoQuestTriggerLogin
	AutoQuestTriggerLevelUp
	AutoQuestTriggerInfoStart
)

func (qc *QuestContainer) RunAutoTriggers(actx actor.Context, trigger AutoQuestTrigger, key uint32) {
	if qc == nil {
		return
	}
	if qc.owner == nil || qc.owner.GameWorld == nil {
		return
	}
	resources := qc.owner.GameWorld.GetResources()
	if resources == nil {
		return
	}
	quests := qc.autoCandidates(resources, trigger, key)
	if len(quests) == 0 {
		return
	}
	for _, def := range quests {
		if def == nil {
			continue
		}
		qc.tryRunAutoQuest(actx, def)
	}
}

func (qc *QuestContainer) autoCandidates(resources *wz.Resources, trigger AutoQuestTrigger, key uint32) []*wz.Quest {
	if resources == nil {
		return nil
	}
	switch trigger {
	case AutoQuestTriggerMapEnter:
		return resources.GetQuestsByStartFieldEnter(key)
	case AutoQuestTriggerLogin, AutoQuestTriggerLevelUp:
		return resources.GetQuestsByAutoStart()
	case AutoQuestTriggerInfoStart:
		return resources.GetQuestsByAutoStartInfo(key)
	default:
		return nil
	}
}

func (qc *QuestContainer) tryRunAutoQuest(actx actor.Context, def *wz.Quest) {
	if qc == nil || def == nil || qc.owner == nil || def.Meta.Blocked {
		return
	}
	autoOpts := QuestPhaseOpts{NpcID: nil}

	if qp := qc.Get(def.ID); qp != nil && qp.IsStarted() {
		if (def.Meta.AutoPreComplete || def.Meta.AutoComplete) && !def.HasEndScript() {
			_ = qp.Complete(qc.owner, autoOpts)
		}
		return
	}

	scriptNpcID := uint32(0)
	if def.Start.Requirements.NPC > 0 {
		scriptNpcID = def.Start.Requirements.NPC
	}

	if def.HasStartScript() {
		if qc.owner.Dialog.Thread() != nil {
			return
		}
		if err := qc.CanStart(def.ID, autoOpts); err != nil {
			return
		}
		_ = qc.owner.RunQuestScript(actx, def.ID, scriptNpcID, "on_start")
		return
	}

	if _, err := qc.Start(def.ID, autoOpts); err != nil {
		return
	}
	if !def.Meta.AutoPreComplete && !def.Meta.AutoComplete {
		return
	}
	if def.HasEndScript() {
		return
	}
	qp := qc.Get(def.ID)
	if qp == nil {
		return
	}
	_ = qp.Complete(qc.owner, autoOpts)
}
