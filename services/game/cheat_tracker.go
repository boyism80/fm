package game

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/types"
)

type CheatOffense uint8

const (
	CheatDamageExceeded CheatOffense = iota
	CheatInvalidAttack
	CheatAttackRate
	CheatAttackRange
	CheatAttackWhileDead
	CheatSummonRate
	cheatOffenseCount
)

type cheatRule struct {
	name     string
	points   int
	validity time.Duration
}

var cheatRules = [cheatOffenseCount]cheatRule{
	CheatDamageExceeded:  {"damage_exceeds", 10, 3 * time.Minute},
	CheatInvalidAttack:   {"invalid_attack", 5, time.Minute},
	CheatAttackRate:      {"attack_rate", 3, 30 * time.Second},
	CheatAttackRange:     {"attack_range", 3, 30 * time.Second},
	CheatAttackWhileDead: {"attack_while_dead", 1, 5 * time.Minute},
	CheatSummonRate:      {"summon_rate", 3, 30 * time.Second},
}

// Lag moves when a packet arrives and where its sender stood, so these limits are loose; offenses carry the latency to tell lag from cheating.
const (
	attackRatePerSecond = 8.0
	attackBurstSeconds  = 5
	attackBurst         = attackRatePerSecond * attackBurstSeconds
	summonRatePerSecond = 2.0
	summonBurst         = summonRatePerSecond * attackBurstSeconds
	moveMarginPerSecond = 150
	attackRangeMargin   = moveMarginPerSecond * attackBurstSeconds
)

type cheatEvent struct {
	at     time.Time
	points int
	lagged bool
}

type latencyReporter interface {
	Latency(now time.Time) core.Latency
}

type CheatTracker struct {
	owner      *Character
	events     [cheatOffenseCount][]cheatEvent
	violations [cheatOffenseCount]int
	attacks    paceBucket
	summons    paceBucket
	acted      bool
}

type paceBucket struct {
	tokens     float64
	refilledAt time.Time
}

func (b *paceBucket) over(now time.Time, rate, burst float64) bool {
	if b.refilledAt.IsZero() {
		b.tokens = burst
	} else {
		b.tokens = min(burst, b.tokens+now.Sub(b.refilledAt).Seconds()*rate)
	}
	b.refilledAt = now
	b.tokens--
	if b.tokens < 0 {
		b.tokens = 0
		return true
	}
	return false
}

func (t *CheatTracker) rule(offense CheatOffense) cheatRule {
	rule := cheatRules[offense]
	if override, ok := t.owner.gameWorld.AntiCheat().Offenses[rule.name]; ok {
		if override.Points > 0 {
			rule.points = override.Points
		}
		if override.ValiditySeconds > 0 {
			rule.validity = time.Duration(override.ValiditySeconds) * time.Second
		}
	}
	return rule
}

func (t *CheatTracker) Score() (clean, lagged int) {
	now := clock.Now()
	for offense := range t.events {
		validity := t.rule(CheatOffense(offense)).validity
		kept := t.events[offense][:0]
		for _, event := range t.events[offense] {
			if now.Sub(event.at) >= validity {
				continue
			}
			kept = append(kept, event)
			if event.lagged {
				lagged += event.points
			} else {
				clean += event.points
			}
		}
		t.events[offense] = kept
	}
	return clean, lagged
}

func (t *CheatTracker) Register(actx actor.Context, offense CheatOffense, info dto.AttackInfo, format string, args ...any) {
	cfg := t.owner.gameWorld.AntiCheat()
	rule := t.rule(offense)
	now := clock.Now()

	var latency core.Latency
	if reporter, ok := t.owner.Client().(latencyReporter); ok {
		latency = reporter.Latency(now)
	}
	lagged := latency.Worst() > time.Duration(cfg.LagRTTMilliseconds)*time.Millisecond

	t.events[offense] = append(t.events[offense], cheatEvent{at: now, points: rule.points, lagged: lagged})
	t.violations[offense]++
	clean, laggedPoints := t.Score()
	score := float64(clean) + cfg.LagWeight*float64(laggedPoints)

	log.Printf("cheat offense=%s points=%d count=%d lagged=%t score=%.1f (clean=%d lagged=%d) rtt=%s smoothed=%s pending=%s character=%d name=%q map=%d skill=%d tick=%d: %s",
		rule.name, rule.points, t.violations[offense], lagged, score, clean, laggedPoints,
		latency.Last, latency.Smoothed, latency.Pending,
		t.owner.GetID(), t.owner.GetName(), t.owner.GetMap().TemplateID(), info.Skill, info.LastTick, fmt.Sprintf(format, args...))

	if score < float64(cfg.Threshold) {
		t.acted = false
		return
	}
	if t.acted {
		return
	}
	t.acted = true
	t.act(actx, cfg.Action, score, clean, laggedPoints, latency)
}

func (t *CheatTracker) act(actx actor.Context, action string, score float64, clean, lagged int, latency core.Latency) {
	ch := t.owner
	// The summary goes to the clients as text, which has to encode in EUC-KR: no duration units like µs.
	summary := fmt.Sprintf("%s(%d) score %.1f (clean %d, lagged %d), rtt %dms", ch.GetName(), ch.GetID(), score, clean, lagged, latency.Worst().Milliseconds())
	log.Printf("cheat threshold reached action=%s: %s", action, summary)

	switch action {
	case "notify":
		ch.gameWorld.Systems().Get[*InternalSystem]().NotifyAdminsAsync(actx, "[치트 의심] "+summary).OnError(func(err error) {
			log.Printf("cheat notify character=%d: %v", ch.GetID(), err)
		})
	case "disconnect":
		ch.listener.Disconnect(ch)
	case "suspend":
		ch.gameWorld.Systems().Get[*InternalSystem]().SetAccountBannedAsync(actx, ch, "anti-cheat: "+summary).OnError(func(err error) {
			log.Printf("cheat suspend character=%d: %v", ch.GetID(), err)
		})
		ch.listener.Disconnect(ch)
	}
}

func (t *CheatTracker) PaceSummon(actx actor.Context, summon *Summon) {
	count := float64(max(1, len(t.owner.summons.entries)))
	info := dto.AttackInfo{AttackHeader: dto.AttackHeader{Skill: uint32(summon.skillID)}}
	if t.summons.over(clock.Now(), summonRatePerSecond*count, summonBurst*count) {
		t.Register(actx, CheatSummonRate, info, "more than %.0f attacks per second per summon sustained over %ds", summonRatePerSecond, attackBurstSeconds)
	}
}

func (t *CheatTracker) Rejected(actx actor.Context, info dto.AttackInfo, err error) {
	if errors.Is(err, ErrDamageExceedsMaximum) {
		t.Register(actx, CheatDamageExceeded, info, "%v", err)
		return
	}
	t.Register(actx, CheatInvalidAttack, info, "%v", err)
}

func (t *CheatTracker) Report() string {
	clean, lagged := t.Score()
	report := fmt.Sprintf("clean=%d lagged=%d", clean, lagged)
	for offense, count := range t.violations {
		if count > 0 {
			report += fmt.Sprintf(" %s=%d", cheatRules[offense].name, count)
		}
	}
	return report
}

// Pace runs before validation, so rejected attacks count too.
func (t *CheatTracker) Pace(actx actor.Context, info dto.AttackInfo) {
	if t.attacks.over(clock.Now(), attackRatePerSecond, attackBurst) {
		t.Register(actx, CheatAttackRate, info, "more than %.0f attacks per second sustained over %ds", attackRatePerSecond, attackBurstSeconds)
	}

	skill := t.owner.skills.Get(info.Skill)
	if skill == nil || skill.Level() <= 0 {
		return
	}
	data := skill.wz.GetLevelData(skill.Level())
	if data == nil {
		return
	}
	area := data.Area(int32(t.owner.position.X), int32(t.owner.position.Y), t.owner.stance%2 != 0)
	area.Left -= attackRangeMargin
	area.Right += attackRangeMargin
	area.Top -= attackRangeMargin
	area.Bottom += attackRangeMargin
	var outside []uint32
	for _, target := range info.Damages {
		mob := t.owner.GetMap().GetMob(target.OID)
		if mob == nil || !mob.IsAlive() {
			continue
		}
		if !area.ContainsPoint(types.Point[int32]{X: int32(mob.position.X), Y: int32(mob.position.Y)}) {
			outside = append(outside, target.OID)
		}
	}
	if len(outside) > 0 {
		t.Register(actx, CheatAttackRange, info, "mobs=%v outside, attacker at (%d,%d)", outside, t.owner.position.X, t.owner.position.Y)
	}
}
