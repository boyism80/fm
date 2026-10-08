package entity

import "github.com/boyism80/fm/services/game/constant"

type Points struct {
	owner    *Character
	AP       uint16
	SP       uint16
	HPMPUsed uint16
}

type APEntry struct {
	Stat   constant.StatType
	Amount uint16
}

func (p *Points) SetAP(v uint16, notify bool) {
	p.AP = v
	if notify {
		p.owner.Listener.OnUpdateStats(p.owner, map[constant.Stat]int32{
			constant.StatAvailableAP: int32(p.AP),
		}, false)
	}
}

func (p *Points) AssignStats(entries []APEntry) bool {
	total := uint32(0)
	for _, e := range entries {
		if p.owner.Stats.canAdd(e.Stat, e.Amount) == false {
			return false
		}
		total += uint32(e.Amount)
	}
	if total > uint32(p.AP) {
		return false
	}

	statUpdate := map[constant.Stat]int32{}
	for _, e := range entries {
		p.owner.Stats.add(e.Stat, e.Amount, statUpdate)
	}
	p.AP -= uint16(total)
	statUpdate[constant.StatAvailableAP] = int32(p.AP)
	p.owner.Listener.OnUpdateStats(p.owner, statUpdate, true)
	return true
}

func (p *Points) AssignHPMP(stat constant.StatType, increase uint32) bool {
	if p.AP == 0 || p.HPMPUsed >= constant.HpAPUsedMax {
		return false
	}

	statUpdate := map[constant.Stat]int32{}
	switch stat {
	case constant.StatTypeHP:
		if p.owner.GetMaxHp() >= constant.StatMaxHPMP {
			return false
		}
		p.owner.AddBaseHp(increase, false)
		statUpdate[constant.StatHP] = int32(p.owner.GetHp())
		statUpdate[constant.StatMaxHP] = int32(p.owner.GetMaxHp())
	case constant.StatTypeMP:
		if p.owner.GetMaxMp() >= constant.StatMaxHPMP {
			return false
		}
		p.owner.AddBaseMp(increase, false)
		statUpdate[constant.StatMP] = int32(p.owner.GetMp())
		statUpdate[constant.StatMaxMP] = int32(p.owner.GetMaxMp())
	default:
		return false
	}

	p.HPMPUsed++
	p.AP--
	statUpdate[constant.StatAvailableAP] = int32(p.AP)
	p.owner.Listener.OnUpdateStats(p.owner, statUpdate, true)
	return true
}

func (p *Points) SetSP(v uint16, notify bool) {
	p.SP = v
	if notify {
		p.owner.Listener.OnUpdateStats(p.owner, map[constant.Stat]int32{
			constant.StatAvailableSP: int32(p.SP),
		}, false)
	}
}

func (p *Points) DistributeSkill(skillID uint32) bool {
	if p.SP == 0 || p.owner.canLearn(skillID) == false {
		return false
	}

	skillEntry := p.owner.Skills.Get(skillID)
	if skillEntry == nil {
		if p.owner.GameWorld == nil {
			return false
		}
		resources := p.owner.GameWorld.GetResources()
		if resources == nil {
			return false
		}
		wzSkill := resources.GetSkill(skillID)
		if wzSkill == nil {
			return false
		}
		skillEntry = NewSkillEntry(p.owner, wzSkill, 0, wzSkill.DefaultMasterLevel())
	}

	if skillEntry.Level() >= min(skillEntry.MasterLevel, skillEntry.Wz.MaxLevel) {
		return false
	}

	p.SP--
	if skillEntry.Level() == 0 {
		skillEntry.SetLevel(1)
		p.owner.Skills.Bind(skillID, skillEntry)
	} else {
		skillEntry.SetLevel(skillEntry.Level() + 1)
	}
	p.owner.Listener.OnUpdateStats(p.owner, map[constant.Stat]int32{
		constant.StatAvailableSP: int32(p.SP),
	}, false)
	return true
}
