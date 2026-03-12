package entity

import "github.com/boyism80/fm/game/constant"

type BaseStats struct {
	Str uint16
	Dex uint16
	Int uint16
	Luk uint16
}

type BonusStats struct {
	Str            int16
	Dex            int16
	Int            int16
	Luk            int16
	Watk           int16
	Matk           int16
	Wdef           int16
	Mdef           int16
	Acc            int16
	Avoid          int16
	Speed          int16
	Jump           int16
	MaxHpFixed     int16
	MaxMpFixed     int16
	MaxHpPercent   int16
	MaxMpPercent   int16
	MesoMultiplier int16 // 100 = 100%, 0 = use 100
	DropRate       int16 // 100 = 100%, 0 = use 100; applies to both item and meso drop probability
	ExpRate        int16 // 100 = 100%, 0 = no bonus; when > 0, exp is multiplied by ExpRate/100 (e.g. 150 = 1.5x for Holy Symbol)
}

func (ch *Character) GetTotalStr() uint16 {
	total := int32(ch.BaseStats.Str) + int32(ch.BonusStats.Str)
	if total < 0 {
		return 0
	}
	if total > int32(constant.STAT_MAX_STR_DEX_INT_LUK) {
		return constant.STAT_MAX_STR_DEX_INT_LUK
	}
	return uint16(total)
}

func (ch *Character) GetTotalDex() uint16 {
	total := int32(ch.BaseStats.Dex) + int32(ch.BonusStats.Dex)
	if total < 0 {
		return 0
	}
	if total > int32(constant.STAT_MAX_STR_DEX_INT_LUK) {
		return constant.STAT_MAX_STR_DEX_INT_LUK
	}
	return uint16(total)
}

func (ch *Character) GetTotalInt() uint16 {
	total := int32(ch.BaseStats.Int) + int32(ch.BonusStats.Int)
	if total < 0 {
		return 0
	}
	if total > int32(constant.STAT_MAX_STR_DEX_INT_LUK) {
		return constant.STAT_MAX_STR_DEX_INT_LUK
	}
	return uint16(total)
}

func (ch *Character) GetTotalLuk() uint16 {
	total := int32(ch.BaseStats.Luk) + int32(ch.BonusStats.Luk)
	if total < 0 {
		return 0
	}
	if total > int32(constant.STAT_MAX_STR_DEX_INT_LUK) {
		return constant.STAT_MAX_STR_DEX_INT_LUK
	}
	return uint16(total)
}

func (ch *Character) GetMaxHp() uint16 {
	base := int32(ch.Life.BaseHp) + int32(ch.Life.BonusHp) + int32(ch.BonusStats.MaxHpFixed)
	total := base + (base*int32(ch.BonusStats.MaxHpPercent))/100
	if total < 1 {
		return 1
	}
	if total > int32(constant.STAT_MAX_HP_MP) {
		return constant.STAT_MAX_HP_MP
	}
	return uint16(total)
}

func (ch *Character) GetMaxMp() uint16 {
	base := int32(ch.Life.BaseMp) + int32(ch.Life.BonusMp) + int32(ch.BonusStats.MaxMpFixed)
	total := base + (base*int32(ch.BonusStats.MaxMpPercent))/100
	if total < 0 {
		return 0
	}
	if total > int32(constant.STAT_MAX_HP_MP) {
		return constant.STAT_MAX_HP_MP
	}
	return uint16(total)
}

func (ch *Character) notifyStatChange(stat constant.Stat) {
	if ch.Listener == nil {
		return
	}

	stats := make(map[constant.Stat]int32)
	switch stat {
	case constant.STAT_STR:
		stats[constant.STAT_STR] = int32(ch.GetTotalStr())
	case constant.STAT_DEX:
		stats[constant.STAT_DEX] = int32(ch.GetTotalDex())
	case constant.STAT_INT:
		stats[constant.STAT_INT] = int32(ch.GetTotalInt())
	case constant.STAT_LUK:
		stats[constant.STAT_LUK] = int32(ch.GetTotalLuk())
	case constant.STAT_MAX_HP:
		stats[constant.STAT_MAX_HP] = int32(ch.GetMaxHp())
	case constant.STAT_MAX_MP:
		stats[constant.STAT_MAX_MP] = int32(ch.GetMaxMp())
	}

	ch.Listener.OnUpdateStats(stats, false)
}

func (ch *Character) ConsumeMP(amount uint16) bool {
	if ch.Mp < amount {
		return false
	}
	ch.Mp -= amount
	if ch.Listener != nil {
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{
			constant.STAT_MP: int32(ch.Mp),
		}, false)
	}
	return true
}
