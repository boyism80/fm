package dao

type Character struct {
	Life
	Id            uint32
	Name          string
	Gender        uint8
	SkinColor     uint8
	Face          uint32
	Hair          uint32
	Level         uint8
	Class         uint16
	Rank          uint32
	RankDiff      int32
	ClassRank     uint32
	ClassRankDiff int32
	Admin         bool
	Str           uint16
	Dex           uint16
	Int           uint16
	Luk           uint16
	Hp            uint16
	MaxHp         uint16
	Mp            uint16
	MaxMp         uint16
	AbilityPoint  uint16
	SkillPoint    []uint16
	Exp           uint32
	FamePoint     uint16
	Map           uint32
	SpawnPoint    uint8
	Mega          bool
}

func (ch *Character) IsRanked() bool {
	if ch.Admin {
		return false
	}

	if ch.Level < 30 {
		return false
	}

	return true
}

func (ch *Character) IsEvan() bool {
	return ch.Class == 2001 || (ch.Class >= 2200 && ch.Class <= 2218)
}

func (ch *Character) IsKOC() bool {
	return ch.Class >= 1000 && ch.Class < 2000
}

func (ch *Character) IsMercedes() bool {
	return ch.Class == 2002 || (ch.Class >= 2300 && ch.Class <= 2312)
}

func (ch *Character) IsDemon() bool {
	return ch.Class == 3001 || (ch.Class >= 3100 && ch.Class <= 3112)
}

func (ch *Character) IsAran() bool {
	return (ch.Class >= 2000 && ch.Class <= 2112) && ch.Class != 2001 && ch.Class != 2002
}

func (ch *Character) IsResist() bool {
	return ch.Class >= 3000 && ch.Class <= 3512
}

func (ch *Character) IsAdventurer() bool {
	return ch.Class < 1000
}

func (ch *Character) IsCannon() bool {
	return ch.Class == 1 || ch.Class == 501 || (ch.Class >= 530 && ch.Class <= 532)
}

func (ch *Character) RemainingSkillPoints() uint16 {
	ret := 0
	for _, sp := range ch.SkillPoint {
		if sp > 0 {
			ret++
		}
	}
	return uint16(ret)
}
