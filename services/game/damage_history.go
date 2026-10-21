package game

import (
	"time"

	"github.com/boyism80/fm/core/clock"
)

// damageHistoryWindow is how far the client's view of the character may lag behind the server's.
// An attack is bounded by the loosest state the character had during this window.
const damageHistoryWindow = 5 * time.Second

const damageHistoryLimit = 64

// damageSnapshot is a state the character had until the given time.
type damageSnapshot struct {
	stats     damageStats
	modifiers damageModifiers
	until     time.Time
}

// damageHistory keeps the states a character left recently, so a packet that was
// sent from an older state is still judged against it. The current state is not
// stored: it is computed when the attack is validated.
type damageHistory struct {
	owner *Character
	past  []damageSnapshot
}

// changing keeps the state that is about to change. Call it before every mutation
// of what damageStats reads.
func (h *damageHistory) changing() {
	if h.owner.GetMap() == nil {
		return
	}

	now := clock.Now()
	stats, modifiers := h.owner.damageStats()
	if n := len(h.past); n > 0 && h.past[n-1].stats == stats && h.past[n-1].modifiers == modifiers {
		h.past[n-1].until = now
	} else {
		h.past = append(h.past, damageSnapshot{stats: stats, modifiers: modifiers, until: now})
	}
	h.prune(now)
}

// recent returns the states the character left within the window.
func (h *damageHistory) recent() []damageSnapshot {
	h.prune(clock.Now())
	return h.past
}

func (h *damageHistory) prune(now time.Time) {
	expired := 0
	for expired < len(h.past) && now.Sub(h.past[expired].until) > damageHistoryWindow {
		expired++
	}
	if over := len(h.past) - expired - damageHistoryLimit; over > 0 {
		expired += over
	}
	h.past = h.past[expired:]
}

// base replaces what damageStats produces and keeps the fields an attack layers on top.
func (m damageModifiers) base(snapshot damageModifiers) damageModifiers {
	layered := m
	m = snapshot
	m.ProjectilePAD = layered.ProjectilePAD
	m.HomingPercent = layered.HomingPercent
	m.ShadowPercent = layered.ShadowPercent
	m.ExtraPercent = layered.ExtraPercent
	m.AssassinateSeconds = layered.AssassinateSeconds
	m.SkillDamageOverride = layered.SkillDamageOverride
	m.ComboPercent = max(m.ComboPercent, layered.ComboPercent)
	return m
}
