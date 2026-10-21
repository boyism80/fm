package game

import (
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/types"
)

// These only record. Lag shifts when a packet arrives and where its sender stood, so a limit
// that blocks has to be chosen from what real players produce; until then each one is as
// loose as the damageHistoryWindow allows.
const (
	// The fastest attack animations repeat about every 120 ms.
	attackRatePerSecond = 8.0
	// A client that lagged for the whole window sends it all at once, at its normal pace.
	attackBurst = attackRatePerSecond * float64(damageHistoryWindow/time.Second)
	// A fast character covers about this much ground per second.
	moveMarginPerSecond = 150
	attackRangeMargin   = moveMarginPerSecond * int32(damageHistoryWindow/time.Second)
)

// attackAudit counts how fast and how far a character attacks, and keeps the log of every
// attack the server did not accept as sent.
type attackAudit struct {
	owner      *Character
	tokens     float64
	refilledAt time.Time
	violations map[string]int
}

// observe looks at an attack before it is validated: attacks that are rejected still count toward the pace.
func (a *attackAudit) observe(info dto.AttackInfo) {
	now := clock.Now()
	if a.refilledAt.IsZero() {
		a.tokens = attackBurst
	} else {
		a.tokens = min(attackBurst, a.tokens+now.Sub(a.refilledAt).Seconds()*attackRatePerSecond)
	}
	a.refilledAt = now
	a.tokens--
	if a.tokens < 0 {
		a.tokens = 0
		a.violation("rate", info, "more than %.0f attacks per second sustained over %s", attackRatePerSecond, damageHistoryWindow)
	}

	skill := a.owner.skills.Get(info.Skill)
	if skill == nil || skill.Level() <= 0 {
		return
	}
	data := skill.wz.GetLevelData(skill.Level())
	if data == nil {
		return
	}
	area := data.Area(int32(a.owner.position.X), int32(a.owner.position.Y), a.owner.stance%2 != 0)
	area.Left -= attackRangeMargin
	area.Right += attackRangeMargin
	area.Top -= attackRangeMargin
	area.Bottom += attackRangeMargin
	for _, target := range info.Damages {
		mob := a.owner.GetMap().GetMob(target.OID)
		if mob == nil || !mob.IsAlive() {
			continue
		}
		if !area.ContainsPoint(types.Point[int32]{X: int32(mob.position.X), Y: int32(mob.position.Y)}) {
			a.violation("range", info, "mob=%d at (%d,%d), attacker at (%d,%d)", target.OID, mob.position.X, mob.position.Y, a.owner.position.X, a.owner.position.Y)
		}
	}
}

// violation counts the kind and logs who sent which attack. LastTick tells one request from another.
func (a *attackAudit) violation(kind string, info dto.AttackInfo, format string, args ...any) {
	if a.violations == nil {
		a.violations = map[string]int{}
	}
	a.violations[kind]++
	log.Printf("attack-audit kind=%s count=%d character=%d name=%q map=%d skill=%d tick=%d: %s",
		kind, a.violations[kind], a.owner.GetID(), a.owner.GetName(), a.owner.GetMap().TemplateID(), info.Skill, info.LastTick, fmt.Sprintf(format, args...))
}
