package constant

type SummonMovementType uint8

const (
	SummonMoveStationary       SummonMovementType = 0
	SummonMoveFollow           SummonMovementType = 1
	SummonMoveWalkStationary   SummonMovementType = 2
	SummonMoveCircleFollow     SummonMovementType = 3
	SummonMoveCircleStationary SummonMovementType = 4
)

type SummonType uint8

const (
	SummonTypePuppet SummonType = iota
	SummonTypeNormal
	SummonTypeBuff
)

func AllSummonMovementTypes() map[string]SummonMovementType {
	return map[string]SummonMovementType{
		"Stationary":       SummonMoveStationary,
		"Follow":           SummonMoveFollow,
		"WalkStationary":   SummonMoveWalkStationary,
		"CircleFollow":     SummonMoveCircleFollow,
		"CircleStationary": SummonMoveCircleStationary,
	}
}

func AllSummonTypes() map[string]SummonType {
	return map[string]SummonType{
		"Puppet": SummonTypePuppet,
		"Normal": SummonTypeNormal,
		"Buff":   SummonTypeBuff,
	}
}
