package game

import (
	"time"

	"github.com/boyism80/fm/core/clock"
	pconst "github.com/boyism80/fm/protocol/constant"
)

const (
	populationMinLevel       = 15
	populationDailyCooldown  = 24 * time.Hour
	populationTargetCooldown = 30 * 24 * time.Hour
	populationLimit          = 30000
)

type PopulationHistory struct {
	LastGivenAt time.Time
	Given       map[uint32]time.Time
}

func (ch *Character) GivePopulation(targetID uint32, mode pconst.PopulationMode) {
	m := ch.GetMap()
	if m == nil {
		return
	}

	target := m.GetPlayer(targetID)
	if (mode != pconst.PopulationModeDown && mode != pconst.PopulationModeUp) || target == nil || target == ch || target.AccountID == ch.AccountID {
		ch.Listener.OnPopulationFailed(ch, pconst.PopulationResultWrongName)
		return
	}
	if target.IsHidden() && !ch.HasRoleAtLeast(target.Role) {
		ch.Listener.OnPopulationFailed(ch, pconst.PopulationResultWrongName)
		return
	}
	if ch.GetLevel() < populationMinLevel {
		ch.Listener.OnPopulationFailed(ch, pconst.PopulationResultLevelTooLow)
		return
	}

	now := clock.Now()
	if now.Sub(ch.Stats.PopulationHistory.LastGivenAt) < populationDailyCooldown {
		ch.Listener.OnPopulationFailed(ch, pconst.PopulationResultAlreadyToday)
		return
	}
	if givenAt, ok := ch.Stats.PopulationHistory.Given[targetID]; ok && now.Sub(givenAt) < populationTargetCooldown {
		ch.Listener.OnPopulationFailed(ch, pconst.PopulationResultAlreadyThisMonth)
		return
	}

	delta := int32(-1)
	if mode == pconst.PopulationModeUp {
		delta = 1
	}
	current := int32(target.Stats.Population)
	if current+delta < -populationLimit || current+delta > populationLimit {
		ch.Listener.OnPopulationFailed(ch, pconst.PopulationResultFailed)
		return
	}

	ch.Stats.PopulationHistory.record(targetID, now)
	value := target.Stats.AddPopulation(delta)
	target.Listener.OnPopulationReceived(target, mode, ch.GetName())
	ch.Listener.OnPopulationGiven(ch, mode, target.GetName(), value)
}

func (h *PopulationHistory) record(targetID uint32, now time.Time) {
	for id, givenAt := range h.Given {
		if now.Sub(givenAt) >= populationTargetCooldown {
			delete(h.Given, id)
		}
	}

	h.LastGivenAt = now
	h.Given[targetID] = now
}
