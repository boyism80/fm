package game

import "github.com/boyism80/fm/maple/constant"

type BaseStats struct {
	str uint16
	dex uint16
	int uint16
	luk uint16
}

type BonusStats struct {
	str                int16
	dex                int16
	int                int16
	luk                int16
	watk               int16
	matk               int16
	wdef               int16
	mdef               int16
	acc                int16
	avoid              int16
	speed              int16
	jump               int16
	maxHpPercent       int16
	maxMpPercent       int16
	mesoMultiplier     int16
	dropRate           int16
	expRate            int16
	potionHealRate     int16
	potionDurationRate int16
}

type Stats struct {
	owner             *Character
	base              BaseStats
	bonus             BonusStats
	population        int16
	populationHistory PopulationHistory
}

func (s *Stats) canAdd(stat constant.StatType, amount uint16) bool {
	switch stat {
	case constant.StatTypeStr:
		return s.TotalStr()+amount <= constant.StatMaxStrDexIntLuk
	case constant.StatTypeDex:
		return s.TotalDex()+amount <= constant.StatMaxStrDexIntLuk
	case constant.StatTypeInt:
		return s.TotalInt()+amount <= constant.StatMaxStrDexIntLuk
	case constant.StatTypeLuk:
		return s.TotalLuk()+amount <= constant.StatMaxStrDexIntLuk
	default:
		return false
	}
}

func (s *Stats) add(stat constant.StatType, amount uint16, statUpdate map[constant.Stat]int32) {
	switch stat {
	case constant.StatTypeStr:
		newStr := s.base.str + amount
		if newStr > constant.StatMaxStrDexIntLuk {
			newStr = constant.StatMaxStrDexIntLuk
		}
		s.base.str = newStr
		statUpdate[constant.StatStr] = int32(s.TotalStr())
	case constant.StatTypeDex:
		newDex := s.base.dex + amount
		if newDex > constant.StatMaxStrDexIntLuk {
			newDex = constant.StatMaxStrDexIntLuk
		}
		s.base.dex = newDex
		statUpdate[constant.StatDex] = int32(s.TotalDex())
	case constant.StatTypeInt:
		newInt := s.base.int + amount
		if newInt > constant.StatMaxStrDexIntLuk {
			newInt = constant.StatMaxStrDexIntLuk
		}
		s.base.int = newInt
		statUpdate[constant.StatInt] = int32(s.TotalInt())
	case constant.StatTypeLuk:
		newLuk := s.base.luk + amount
		if newLuk > constant.StatMaxStrDexIntLuk {
			newLuk = constant.StatMaxStrDexIntLuk
		}
		s.base.luk = newLuk
		statUpdate[constant.StatLuk] = int32(s.TotalLuk())
	}
}

func (s *Stats) TotalStr() uint16 {
	total := int32(s.base.str) + int32(s.bonus.str)
	if total < 0 {
		return 0
	}
	if total > int32(constant.StatMaxStrDexIntLuk) {
		return constant.StatMaxStrDexIntLuk
	}
	return uint16(total)
}

func (s *Stats) TotalDex() uint16 {
	total := int32(s.base.dex) + int32(s.bonus.dex)
	if total < 0 {
		return 0
	}
	if total > int32(constant.StatMaxStrDexIntLuk) {
		return constant.StatMaxStrDexIntLuk
	}
	return uint16(total)
}

func (s *Stats) TotalInt() uint16 {
	total := int32(s.base.int) + int32(s.bonus.int)
	if total < 0 {
		return 0
	}
	if total > int32(constant.StatMaxStrDexIntLuk) {
		return constant.StatMaxStrDexIntLuk
	}
	return uint16(total)
}

func (s *Stats) TotalLuk() uint16 {
	total := int32(s.base.luk) + int32(s.bonus.luk)
	if total < 0 {
		return 0
	}
	if total > int32(constant.StatMaxStrDexIntLuk) {
		return constant.StatMaxStrDexIntLuk
	}
	return uint16(total)
}

func (s *Stats) PotionHealMultiplierPercent() int {
	r := s.bonus.potionHealRate
	if r <= 0 {
		return 100
	}
	return int(r)
}

func (s *Stats) PotionDurationMultiplierPercent() int {
	r := s.bonus.potionDurationRate
	if r <= 0 {
		return 100
	}
	return int(r)
}

func (s *Stats) setPopulation(value int32) {
	value = max(-populationLimit, min(populationLimit, value))
	if int16(value) == s.population {
		return
	}

	s.population = int16(value)
	s.owner.listener.OnUpdateStats(s.owner, map[constant.Stat]int32{
		constant.StatPopulation: value,
	}, false)
}

func (s *Stats) AddPopulation(delta int32) int32 {
	s.setPopulation(int32(s.population) + delta)
	return int32(s.population)
}

func (s *Stats) validatePopulation(costPopulation, rewardPopulation int32) ExchangeResult {
	if costPopulation < 0 || rewardPopulation < 0 {
		return ExchangeLackCost
	}
	if costPopulation == 0 && rewardPopulation == 0 {
		return ExchangeOK
	}
	current := int32(s.population)
	if costPopulation > current {
		return ExchangeLackCost
	}
	if rewardPopulation > populationLimit-(current-costPopulation) {
		return ExchangeLackCapacity
	}
	return ExchangeOK
}
