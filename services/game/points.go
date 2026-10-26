package game

import (
	"log"
	"math"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	lua "github.com/yuin/gopher-lua"
)

type Points struct {
	owner    *Character
	ap       uint16
	sp       uint16
	hpmpUsed uint16
}

type APEntry struct {
	Stat   constant.StatType
	Amount uint16
}

func (p *Points) SetAP(v uint16, notify bool) {
	p.ap = v
	if notify {
		p.owner.listener.OnUpdateStats(p.owner, map[constant.Stat]int32{
			constant.StatAvailableAP: int32(p.ap),
		}, false)
	}
}

func (p *Points) AssignAll(entries []APEntry) bool {
	total := uint32(0)
	for _, e := range entries {
		total += uint32(e.Amount)
	}
	if total != uint32(p.ap) {
		return false
	}
	return p.AssignStats(entries)
}

func (p *Points) AssignStats(entries []APEntry) bool {
	combined := map[constant.StatType]uint32{}
	total := uint32(0)
	for _, e := range entries {
		combined[e.Stat] += uint32(e.Amount)
		total += uint32(e.Amount)
	}
	if total > uint32(p.ap) {
		return false
	}
	for stat, amount := range combined {
		if amount > math.MaxUint16 || !p.owner.stats.canAdd(stat, uint16(amount)) {
			return false
		}
	}

	statUpdate := map[constant.Stat]int32{}
	for _, e := range entries {
		p.owner.stats.add(e.Stat, e.Amount, statUpdate)
	}
	p.ap -= uint16(total)
	statUpdate[constant.StatAvailableAP] = int32(p.ap)
	p.owner.listener.OnUpdateStats(p.owner, statUpdate, true)
	return true
}

func (p *Points) AssignHPMP(stat constant.StatType) bool {
	if p.ap == 0 || p.hpmpUsed >= constant.HpAPUsedMax {
		return false
	}

	statUpdate := map[constant.Stat]int32{}
	switch stat {
	case constant.StatTypeHP:
		if int64(p.owner.baseHp)+int64(p.owner.bonusHp) >= int64(constant.StatMaxHPMP) {
			return false
		}
		p.owner.AddBaseHp(p.hpmpIncrease("get_ap_to_hp", 10), false)
		statUpdate[constant.StatHP] = int32(p.owner.GetHp())
		statUpdate[constant.StatMaxHP] = int32(p.owner.GetMaxHp())
	case constant.StatTypeMP:
		if int64(p.owner.baseMp)+int64(p.owner.bonusMp) >= int64(constant.StatMaxHPMP) {
			return false
		}
		p.owner.AddBaseMp(p.hpmpIncrease("get_ap_to_mp", 5), false)
		statUpdate[constant.StatMP] = int32(p.owner.GetMp())
		statUpdate[constant.StatMaxMP] = int32(p.owner.GetMaxMp())
	default:
		return false
	}

	p.hpmpUsed++
	p.ap--
	statUpdate[constant.StatAvailableAP] = int32(p.ap)
	p.owner.listener.OnUpdateStats(p.owner, statUpdate, true)
	return true
}

// ResetAP moves one point from one stat to another, as the AP reset scroll does.
func (p *Points) ResetAP(to, from constant.StatType) bool {
	if to == from {
		return false
	}
	switch to {
	case constant.StatTypeHP:
		if p.hpmpUsed >= constant.HpAPUsedMax || int64(p.owner.baseHp)+int64(p.owner.bonusHp) >= int64(constant.StatMaxHPMP) {
			return false
		}
	case constant.StatTypeMP:
		if p.hpmpUsed >= constant.HpAPUsedMax || int64(p.owner.baseMp)+int64(p.owner.bonusMp) >= int64(constant.StatMaxHPMP) {
			return false
		}
	default:
		if !p.owner.stats.canAdd(to, 1) {
			return false
		}
	}

	statUpdate := map[constant.Stat]int32{}
	switch from {
	case constant.StatTypeHP:
		loss := p.hpmpIncrease("get_ap_reset_hp", 10)
		if p.hpmpUsed == 0 || p.owner.baseHp <= loss {
			return false
		}
		p.owner.SetBaseHp(p.owner.baseHp-loss, false)
		p.hpmpUsed--
	case constant.StatTypeMP:
		loss := p.hpmpIncrease("get_ap_reset_mp", 5)
		if p.hpmpUsed == 0 || p.owner.baseMp <= loss {
			return false
		}
		p.owner.SetBaseMp(p.owner.baseMp-loss, false)
		p.hpmpUsed--
	default:
		if !p.owner.stats.removeOne(from, p.resetFloor(from), statUpdate) {
			return false
		}
	}

	switch to {
	case constant.StatTypeHP:
		p.owner.AddBaseHp(p.hpmpIncrease("get_ap_to_hp", 10), false)
		p.hpmpUsed++
	case constant.StatTypeMP:
		p.owner.AddBaseMp(p.hpmpIncrease("get_ap_to_mp", 5), false)
		p.hpmpUsed++
	default:
		p.owner.stats.add(to, 1, statUpdate)
	}
	if to == constant.StatTypeHP || to == constant.StatTypeMP || from == constant.StatTypeHP || from == constant.StatTypeMP {
		statUpdate[constant.StatHP] = int32(p.owner.GetHp())
		statUpdate[constant.StatMaxHP] = int32(p.owner.GetMaxHp())
		statUpdate[constant.StatMP] = int32(p.owner.GetMp())
		statUpdate[constant.StatMaxMP] = int32(p.owner.GetMaxMp())
	}
	p.owner.listener.OnUpdateStats(p.owner, statUpdate, true)
	return true
}

// resetFloor keeps 4 of each stat and the minimum a class needs of its main stat, as the original server did.
func (p *Points) resetFloor(stat constant.StatType) uint16 {
	class := p.owner.class
	switch {
	case stat == constant.StatTypeStr && (class.Of(constant.ClassWarrior) || class.Of(constant.ClassDawnWarrior1)):
		return 35
	case stat == constant.StatTypeDex && (class.Of(constant.ClassBowman) || class.Of(constant.ClassThief) || class.Of(constant.ClassWindArcher1) || class.Of(constant.ClassNightWalker1)):
		return 25
	case stat == constant.StatTypeInt && (class.Of(constant.ClassMagician) || class.Of(constant.ClassBlazeWizard1)):
		return 20
	}
	return 4
}

func (p *Points) hpmpIncrease(query string, fallback uint32) uint32 {
	m := p.owner.GetMap()
	if m == nil {
		return fallback
	}
	root := m.GetLuaRoot()
	if root == nil {
		return fallback
	}
	thread, err := luax.NewThread(root, constant.CharacterQueryScriptPath)
	if err != nil {
		return fallback
	}
	value, err := luax.Call(thread, query, p.owner)
	if err != nil {
		log.Printf("AssignHPMP: %s character=%d: %v", query, p.owner.GetID(), err)
		return fallback
	}
	n, ok := value.(lua.LNumber)
	if !ok || n <= 0 {
		return fallback
	}
	if float64(n) >= math.MaxUint32 {
		return math.MaxUint32
	}
	return uint32(n)
}

func (p *Points) SetSP(v uint16, notify bool) {
	p.sp = v
	if notify {
		p.owner.listener.OnUpdateStats(p.owner, map[constant.Stat]int32{
			constant.StatAvailableSP: int32(p.sp),
		}, false)
	}
}

func (p *Points) DistributeSkill(skillID uint32) bool {
	wzSkill := p.owner.gameWorld.GetResources().GetSkill(skillID)
	if wzSkill == nil || !wzSkill.LearnableBy(p.owner.class) {
		return false
	}

	// The distributable beginner skills share their own budget, level - 1 up to 6, instead of SP.
	beginner := wzSkill.Beginner && wzSkill.Distributable
	if beginner {
		used := 0
		p.owner.skills.ForEach(func(_ uint32, entry *SkillEntry) {
			if entry.wz.Beginner && entry.wz.Distributable {
				used += entry.Level()
			}
		})
		if used >= min(int(p.owner.level)-1, 6) {
			return false
		}
	} else if p.sp == 0 {
		return false
	}

	skillEntry := p.owner.skills.Learned(skillID)
	if skillEntry == nil {
		skillEntry = NewSkillEntry(p.owner, wzSkill, 0, wzSkill.DefaultMasterLevel())
	}

	if !skillEntry.wz.Distributable {
		return false
	}
	if skillEntry.Level() >= min(skillEntry.masterLevel, skillEntry.wz.MaxLevel) {
		return false
	}

	if !beginner {
		p.sp--
	}
	if skillEntry.Level() == 0 {
		skillEntry.SetLevel(1)
		p.owner.skills.Bind(skillID, skillEntry)
	} else {
		skillEntry.SetLevel(skillEntry.Level() + 1)
	}
	p.owner.listener.OnUpdateStats(p.owner, map[constant.Stat]int32{
		constant.StatAvailableSP: int32(p.sp),
	}, false)
	return true
}
