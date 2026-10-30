package game

import "github.com/boyism80/fm/maple/constant"

func (ch *Character) GetMaxHp() uint32 {
	base := int32(ch.baseHp) + ch.bonusHp
	percent := int32(ch.stats.bonus.maxHpPercent)
	if _, v, ok := ch.buffs.GetBuffValue(constant.BuffFlagMaxHp); ok {
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
	base := int32(ch.baseMp) + ch.bonusMp
	percent := int32(ch.stats.bonus.maxMpPercent)
	if _, v, ok := ch.buffs.GetBuffValue(constant.BuffFlagMaxMp); ok {
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

type statTotals struct {
	STR, DEX, INT, LUK, MAD int
}

// totalStats adds every equipped item to the base stats, as the legacy server's
// total stats did; MAD is that INT plus equipped and buffed magic attack.
func (ch *Character) totalStats() statTotals {
	t := statTotals{
		STR: int(ch.stats.base.str),
		DEX: int(ch.stats.base.dex),
		INT: int(ch.stats.base.int),
		LUK: int(ch.stats.base.luk),
	}
	for slot, item := range ch.inventory.equipped {
		// Cash appearance items have no stats.
		if slot <= -100 || item == nil {
			continue
		}
		e := item.GetEquipmentCore().BuildEquipmentS2C()
		if e == nil {
			continue
		}
		t.STR += int(e.Str)
		t.DEX += int(e.Dex)
		t.INT += int(e.Int)
		t.LUK += int(e.Luk)
		t.MAD += int(e.MAD)
	}
	t.MAD += t.INT + int(ch.buffs.value(constant.BuffFlagMagicAtk))
	return t
}

func (ch *Character) GetStatValue(stat constant.Stat) (int32, bool) {
	switch stat {
	case constant.StatLevel:
		return int32(ch.level), true
	case constant.StatEXP:
		return int32(ch.exp), true
	case constant.StatClass:
		return int32(ch.class), true
	case constant.StatStr:
		return int32(ch.stats.base.str), true
	case constant.StatDex:
		return int32(ch.stats.base.dex), true
	case constant.StatInt:
		return int32(ch.stats.base.int), true
	case constant.StatLuk:
		return int32(ch.stats.base.luk), true
	case constant.StatHP:
		return int32(ch.GetHp()), true
	case constant.StatMaxHP:
		return int32(ch.GetMaxHp()), true
	case constant.StatMP:
		return int32(ch.GetMp()), true
	case constant.StatMaxMP:
		return int32(ch.GetMaxMp()), true
	case constant.StatAvailableAP:
		return int32(ch.points.ap), true
	case constant.StatAvailableSP:
		return int32(ch.points.sp), true
	case constant.StatPopulation:
		return int32(ch.stats.population), true
	case constant.StatMeso:
		return ch.inventory.meso, true
	default:
		return 0, false
	}
}

func (ch *Character) notifyStatChange(stat constant.Stat) {
	stats := make(map[constant.Stat]int32)
	switch stat {
	case constant.StatStr:
		stats[constant.StatStr] = int32(ch.stats.base.str)
	case constant.StatDex:
		stats[constant.StatDex] = int32(ch.stats.base.dex)
	case constant.StatInt:
		stats[constant.StatInt] = int32(ch.stats.base.int)
	case constant.StatLuk:
		stats[constant.StatLuk] = int32(ch.stats.base.luk)
	case constant.StatMaxHP:
		stats[constant.StatMaxHP] = int32(ch.GetMaxHp())
	case constant.StatMaxMP:
		stats[constant.StatMaxMP] = int32(ch.GetMaxMp())
	}

	ch.listener.OnUpdateStats(ch, stats, false)
}

func (ch *Character) ConsumeMP(amount uint32) bool {
	if ch.GetMp() < amount {
		return false
	}
	ch.SetMp(ch.GetMp()-amount, false)
	ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatMP: int32(ch.GetMp()),
	}, false)
	return true
}

func (ch *Character) ConsumeHP(amount uint32) bool {
	if ch.GetHp() <= amount {
		return false
	}
	ch.SetHp(ch.GetHp()-amount, false)
	ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatHP: int32(ch.GetHp()),
	}, false)
	return true
}
