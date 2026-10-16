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
	lastGivenAt time.Time
	given       map[uint32]time.Time
}

func (ch *Character) GivePopulation(targetID uint32, mode pconst.PopulationMode) {
	m := ch.GetMap()
	if m == nil {
		return
	}

	target := m.GetPlayer(targetID)
	if (mode != pconst.PopulationModeDown && mode != pconst.PopulationModeUp) || target == nil || target == ch || target.accountID == ch.accountID {
		ch.listener.OnPopulationFailed(ch, pconst.PopulationResultWrongName)
		return
	}
	if target.IsHidden() && !ch.HasRoleAtLeast(target.role) {
		ch.listener.OnPopulationFailed(ch, pconst.PopulationResultWrongName)
		return
	}
	if ch.GetLevel() < populationMinLevel {
		ch.listener.OnPopulationFailed(ch, pconst.PopulationResultLevelTooLow)
		return
	}

	now := clock.Now()
	if now.Sub(ch.stats.populationHistory.lastGivenAt) < populationDailyCooldown {
		ch.listener.OnPopulationFailed(ch, pconst.PopulationResultAlreadyToday)
		return
	}
	if givenAt, ok := ch.stats.populationHistory.given[targetID]; ok && now.Sub(givenAt) < populationTargetCooldown {
		ch.listener.OnPopulationFailed(ch, pconst.PopulationResultAlreadyThisMonth)
		return
	}

	delta := int32(-1)
	if mode == pconst.PopulationModeUp {
		delta = 1
	}
	current := int32(target.stats.population)
	if current+delta < -populationLimit || current+delta > populationLimit {
		ch.listener.OnPopulationFailed(ch, pconst.PopulationResultFailed)
		return
	}

	ch.stats.populationHistory.record(targetID, now)
	value := target.stats.AddPopulation(delta)
	target.listener.OnPopulationReceived(target, mode, ch.GetName())
	ch.listener.OnPopulationGiven(ch, mode, target.GetName(), value)
}

func (h *PopulationHistory) record(targetID uint32, now time.Time) {
	for id, givenAt := range h.given {
		if now.Sub(givenAt) >= populationTargetCooldown {
			delete(h.given, id)
		}
	}

	h.lastGivenAt = now
	h.given[targetID] = now
}
