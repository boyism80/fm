package entity

import "github.com/boyism80/fm/services/game/constant"

func (ch *Character) GetMaxHp() uint32 {
	base := int32(ch.BaseHp) + ch.BonusHp
	percent := int32(ch.Stats.Bonus.MaxHpPercent)
	if _, v, ok := ch.Buffs.GetBuffValue(constant.BuffFlagMaxHp); ok {
		percent += v
	}
	total := base + (base*percent)/100
	if total < 1 {
		return 1
	}
	if total > int32(constant.StatMaxHPMP) {
		return constant.StatMaxHPMP
	}
	return uint32(total)
}

func (ch *Character) GetMaxMp() uint32 {
	base := int32(ch.BaseMp) + ch.BonusMp
	percent := int32(ch.Stats.Bonus.MaxMpPercent)
	if _, v, ok := ch.Buffs.GetBuffValue(constant.BuffFlagMaxMp); ok {
		percent += v
	}
	total := base + (base*percent)/100
	if total < 0 {
		return 0
	}
	if total > int32(constant.StatMaxHPMP) {
		return constant.StatMaxHPMP
	}
	return uint32(total)
}

func (ch *Character) GetStatValue(stat constant.Stat) (int32, bool) {
	switch stat {
	case constant.StatLevel:
		return int32(ch.level), true
	case constant.StatEXP:
		return int32(ch.exp), true
	case constant.StatClass:
		return int32(ch.Class), true
	case constant.StatStr:
		return int32(ch.Stats.TotalStr()), true
	case constant.StatDex:
		return int32(ch.Stats.TotalDex()), true
	case constant.StatInt:
		return int32(ch.Stats.TotalInt()), true
	case constant.StatLuk:
		return int32(ch.Stats.TotalLuk()), true
	case constant.StatHP:
		return int32(ch.GetHp()), true
	case constant.StatMaxHP:
		return int32(ch.GetMaxHp()), true
	case constant.StatMP:
		return int32(ch.GetMp()), true
	case constant.StatMaxMP:
		return int32(ch.GetMaxMp()), true
	case constant.StatAvailableAP:
		return int32(ch.Points.AP), true
	case constant.StatAvailableSP:
		return int32(ch.Points.SP), true
	case constant.StatPopulation:
		return int32(ch.Stats.Population), true
	case constant.StatMeso:
		return ch.Inventory.Meso, true
	default:
		return 0, false
	}
}

func (ch *Character) notifyStatChange(stat constant.Stat) {
	stats := make(map[constant.Stat]int32)
	switch stat {
	case constant.StatStr:
		stats[constant.StatStr] = int32(ch.Stats.TotalStr())
	case constant.StatDex:
		stats[constant.StatDex] = int32(ch.Stats.TotalDex())
	case constant.StatInt:
		stats[constant.StatInt] = int32(ch.Stats.TotalInt())
	case constant.StatLuk:
		stats[constant.StatLuk] = int32(ch.Stats.TotalLuk())
	case constant.StatMaxHP:
		stats[constant.StatMaxHP] = int32(ch.GetMaxHp())
	case constant.StatMaxMP:
		stats[constant.StatMaxMP] = int32(ch.GetMaxMp())
	}

	ch.Listener.OnUpdateStats(ch, stats, false)
}

func (ch *Character) ConsumeMP(amount uint32) bool {
	if ch.GetMp() < amount {
		return false
	}
	ch.SetMp(ch.GetMp()-amount, false)
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatMP: int32(ch.GetMp()),
	}, false)
	return true
}

func (ch *Character) ConsumeHP(amount uint32) bool {
	if ch.GetHp() <= amount {
		return false
	}
	ch.SetHp(ch.GetHp()-amount, false)
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatHP: int32(ch.GetHp()),
	}, false)
	return true
}
