package wz

type Drop struct {
	Item  uint32
	Money uint32
	Prob  float32
	Min   uint16
	Max   uint16
}

type GlobalDropType uint8

const (
	GlobalDropTypeAll GlobalDropType = iota
	GlobalDropTypeBoss
	GlobalDropTypeExplosive
	GlobalDropTypeRegular
)

type GlobalDrop struct {
	Drop
	Continent int32
	Type      GlobalDropType
	Quest     uint32
}

func (d GlobalDrop) IsAvailableAt(mapID uint32) bool {
	switch {
	case d.Continent < 0:
		return true
	case d.Continent < 10:
		return mapID/100000000 == uint32(d.Continent)
	case d.Continent < 100:
		return mapID/10000000 == uint32(d.Continent)
	default:
		return mapID/1000000 == uint32(d.Continent)
	}
}

func (d GlobalDrop) IsAvailableFor(mob *Mob) bool {
	switch d.Type {
	case GlobalDropTypeAll:
		return true
	case GlobalDropTypeBoss:
		return mob.Boss
	case GlobalDropTypeExplosive:
		return mob.ExplosiveReward
	case GlobalDropTypeRegular:
		return !mob.Boss && !mob.ExplosiveReward
	default:
		return false
	}
}
