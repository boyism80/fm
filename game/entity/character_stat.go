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
	base := int32(ch.BaseHp) + int32(ch.BonusHp)
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
	base := int32(ch.BaseMp) + int32(ch.BonusMp)
	total := base + (base*int32(ch.BonusStats.MaxMpPercent))/100
	if total < 0 {
		return 0
	}
	if total > int32(constant.STAT_MAX_HP_MP) {
		return constant.STAT_MAX_HP_MP
	}
	return uint16(total)
}

func (ch *Character) GetStatValue(stat constant.Stat) (int32, bool) {
	switch stat {
	case constant.STAT_LEVEL:
		return int32(ch.level), true
	case constant.STAT_EXP:
		return int32(ch.exp), true
	case constant.STAT_CLASS:
		return int32(ch.Class), true
	case constant.STAT_STR:
		return int32(ch.GetTotalStr()), true
	case constant.STAT_DEX:
		return int32(ch.GetTotalDex()), true
	case constant.STAT_INT:
		return int32(ch.GetTotalInt()), true
	case constant.STAT_LUK:
		return int32(ch.GetTotalLuk()), true
	case constant.STAT_HP:
		return int32(ch.Hp), true
	case constant.STAT_MAX_HP:
		return int32(ch.GetMaxHp()), true
	case constant.STAT_MP:
		return int32(ch.Mp), true
	case constant.STAT_MAX_MP:
		return int32(ch.GetMaxMp()), true
	case constant.STAT_AVAILABLE_AP:
		return int32(ch.AbilityPoint), true
	case constant.STAT_AVAILABLE_SP:
		return int32(ch.SkillPoint), true
	case constant.STAT_FAME:
		return int32(ch.famePoint), true
	case constant.STAT_MESO:
		return ch.Meso, true
	default:
		return 0, false
	}
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

func (ch *Character) ConsumeHP(amount uint16) bool {
	if ch.Hp <= amount {
		return false
	}
	ch.Hp -= amount
	if ch.Listener != nil {
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{
			constant.STAT_HP: int32(ch.Hp),
		}, false)
	}
	return true
}
