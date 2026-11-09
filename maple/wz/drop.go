package wz

type Drop struct {
	Item  uint32
	Money uint32
	Prob  float32
	Min   uint16
	Max   uint16
	Quest uint32 // only those with this quest started see the drop
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
}

// codedGlobalDrops are the global drops the original server added in code, not in its drop table, so they are there
// whether or not GlobalReward.img.xml is.
var codedGlobalDrops = []GlobalDrop{
	// The monster card, the only source of the match card set.
	{Drop: Drop{Item: 4030012, Prob: 1000.0 / 999999, Min: 1, Max: 1}, Continent: -1},
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
