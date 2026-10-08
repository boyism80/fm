package entity

import "github.com/boyism80/fm/services/game/constant"

type BaseStats struct {
	Str uint16
	Dex uint16
	Int uint16
	Luk uint16
}

type BonusStats struct {
	Str                int16
	Dex                int16
	Int                int16
	Luk                int16
	Watk               int16
	Matk               int16
	Wdef               int16
	Mdef               int16
	Acc                int16
	Avoid              int16
	Speed              int16
	Jump               int16
	MaxHpPercent       int16
	MaxMpPercent       int16
	MesoMultiplier     int16
	DropRate           int16
	ExpRate            int16
	PotionHealRate     int16
	PotionDurationRate int16
}

type Stats struct {
	owner      *Character
	Base       BaseStats
	Bonus      BonusStats
	Population uint16
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
		newStr := s.Base.Str + amount
		if newStr > constant.StatMaxStrDexIntLuk {
			newStr = constant.StatMaxStrDexIntLuk
		}
		s.Base.Str = newStr
		statUpdate[constant.StatStr] = int32(s.TotalStr())
	case constant.StatTypeDex:
		newDex := s.Base.Dex + amount
		if newDex > constant.StatMaxStrDexIntLuk {
			newDex = constant.StatMaxStrDexIntLuk
		}
		s.Base.Dex = newDex
		statUpdate[constant.StatDex] = int32(s.TotalDex())
	case constant.StatTypeInt:
		newInt := s.Base.Int + amount
		if newInt > constant.StatMaxStrDexIntLuk {
			newInt = constant.StatMaxStrDexIntLuk
		}
		s.Base.Int = newInt
		statUpdate[constant.StatInt] = int32(s.TotalInt())
	case constant.StatTypeLuk:
		newLuk := s.Base.Luk + amount
		if newLuk > constant.StatMaxStrDexIntLuk {
			newLuk = constant.StatMaxStrDexIntLuk
		}
		s.Base.Luk = newLuk
		statUpdate[constant.StatLuk] = int32(s.TotalLuk())
	}
}

func (s *Stats) TotalStr() uint16 {
	total := int32(s.Base.Str) + int32(s.Bonus.Str)
	if total < 0 {
		return 0
	}
	if total > int32(constant.StatMaxStrDexIntLuk) {
		return constant.StatMaxStrDexIntLuk
	}
	return uint16(total)
}

func (s *Stats) TotalDex() uint16 {
	total := int32(s.Base.Dex) + int32(s.Bonus.Dex)
	if total < 0 {
		return 0
	}
	if total > int32(constant.StatMaxStrDexIntLuk) {
		return constant.StatMaxStrDexIntLuk
	}
	return uint16(total)
}

func (s *Stats) TotalInt() uint16 {
	total := int32(s.Base.Int) + int32(s.Bonus.Int)
	if total < 0 {
		return 0
	}
	if total > int32(constant.StatMaxStrDexIntLuk) {
		return constant.StatMaxStrDexIntLuk
	}
	return uint16(total)
}

func (s *Stats) TotalLuk() uint16 {
	total := int32(s.Base.Luk) + int32(s.Bonus.Luk)
	if total < 0 {
		return 0
	}
	if total > int32(constant.StatMaxStrDexIntLuk) {
		return constant.StatMaxStrDexIntLuk
	}
	return uint16(total)
}

func (s *Stats) PotionHealMultiplierPercent() int {
	r := s.Bonus.PotionHealRate
	if r <= 0 {
		return 100
	}
	return int(r)
}

func (s *Stats) PotionDurationMultiplierPercent() int {
	r := s.Bonus.PotionDurationRate
	if r <= 0 {
		return 100
	}
	return int(r)
}

func (s *Stats) setPopulation(value int32) {
	if value < 0 {
		value = 0
	}
	if value > 65535 {
		value = 65535
	}
	s.Population = uint16(value)
	s.owner.Listener.OnUpdateStats(s.owner, map[constant.Stat]int32{
		constant.StatPopulation: int32(s.Population),
	}, false)
}

func (s *Stats) gainPopulation(amount int32) {
	if amount <= 0 {
		return
	}
	sum := uint32(s.Population) + uint32(amount)
	if sum > 65535 {
		sum = 65535
	}
	s.Population = uint16(sum)
	s.owner.Listener.OnUpdateStats(s.owner, map[constant.Stat]int32{
		constant.StatPopulation: int32(s.Population),
	}, false)
}

func (s *Stats) losePopulation(amount int32) {
	if amount <= 0 {
		return
	}
	if uint32(amount) >= uint32(s.Population) {
		s.Population = 0
	} else {
		s.Population -= uint16(amount)
	}
	s.owner.Listener.OnUpdateStats(s.owner, map[constant.Stat]int32{
		constant.StatPopulation: int32(s.Population),
	}, false)
}
