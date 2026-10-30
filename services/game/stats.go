package game

import "github.com/boyism80/fm/maple/constant"

type BaseStats struct {
	str uint16
	dex uint16
	int uint16
	luk uint16
}

type BonusStats struct {
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
		return s.base.str+amount <= constant.StatMaxStrDexIntLuk
	case constant.StatTypeDex:
		return s.base.dex+amount <= constant.StatMaxStrDexIntLuk
	case constant.StatTypeInt:
		return s.base.int+amount <= constant.StatMaxStrDexIntLuk
	case constant.StatTypeLuk:
		return s.base.luk+amount <= constant.StatMaxStrDexIntLuk
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
		statUpdate[constant.StatStr] = int32(s.base.str)
	case constant.StatTypeDex:
		newDex := s.base.dex + amount
		if newDex > constant.StatMaxStrDexIntLuk {
			newDex = constant.StatMaxStrDexIntLuk
		}
		s.base.dex = newDex
		statUpdate[constant.StatDex] = int32(s.base.dex)
	case constant.StatTypeInt:
		newInt := s.base.int + amount
		if newInt > constant.StatMaxStrDexIntLuk {
			newInt = constant.StatMaxStrDexIntLuk
		}
		s.base.int = newInt
		statUpdate[constant.StatInt] = int32(s.base.int)
	case constant.StatTypeLuk:
		newLuk := s.base.luk + amount
		if newLuk > constant.StatMaxStrDexIntLuk {
			newLuk = constant.StatMaxStrDexIntLuk
		}
		s.base.luk = newLuk
		statUpdate[constant.StatLuk] = int32(s.base.luk)
	}
}

// removeOne takes one point off a base stat that stays at or above floor.
func (s *Stats) removeOne(stat constant.StatType, floor uint16, statUpdate map[constant.Stat]int32) bool {
	var base *uint16
	var key constant.Stat
	switch stat {
	case constant.StatTypeStr:
		base, key = &s.base.str, constant.StatStr
	case constant.StatTypeDex:
		base, key = &s.base.dex, constant.StatDex
	case constant.StatTypeInt:
		base, key = &s.base.int, constant.StatInt
	case constant.StatTypeLuk:
		base, key = &s.base.luk, constant.StatLuk
	default:
		return false
	}
	if *base <= floor {
		return false
	}
	*base--
	switch stat {
	case constant.StatTypeStr:
		statUpdate[key] = int32(s.base.str)
	case constant.StatTypeDex:
		statUpdate[key] = int32(s.base.dex)
	case constant.StatTypeInt:
		statUpdate[key] = int32(s.base.int)
	case constant.StatTypeLuk:
		statUpdate[key] = int32(s.base.luk)
	}
	return true
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
