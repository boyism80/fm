package constant

type StanceKind uint8

const (
	StanceDefault  StanceKind = 0
	StanceWalk     StanceKind = 2
	StanceStanding StanceKind = 4
	StanceJump     StanceKind = 6
	StanceAttack   StanceKind = 8
	StanceProne    StanceKind = 10
	StanceRope     StanceKind = 12
	StanceLadder   StanceKind = 14
	StanceSit      StanceKind = 21
)

const (
	StanceDefaultValue  = 0
	StanceWalkRight     = 2
	StanceWalkLeft      = 3
	StanceStandingRight = 4
	StanceStandingLeft  = 5
	StanceJumpRight     = 6
	StanceJumpLeft      = 7
	StanceAttackRight   = 8
	StanceAttackLeft    = 9
	StanceProneRight    = 10
	StanceProneLeft     = 11
	StanceRopeRight     = 12
	StanceRopeLeft      = 13
	StanceLadderRight   = 14
	StanceLadderLeft    = 15
	StanceSitValue      = 21
)

func AllStanceConstants() map[string]uint8 {
	return map[string]uint8{
		"Default":       StanceDefaultValue,
		"WalkRight":     StanceWalkRight,
		"WalkLeft":      StanceWalkLeft,
		"Walk":          StanceWalkRight,
		"StandingRight": StanceStandingRight,
		"StandingLeft":  StanceStandingLeft,
		"Standing":      StanceStandingRight,
		"JumpRight":     StanceJumpRight,
		"JumpLeft":      StanceJumpLeft,
		"Jump":          StanceJumpRight,
		"AttackRight":   StanceAttackRight,
		"AttackLeft":    StanceAttackLeft,
		"Attack":        StanceAttackRight,
		"ProneRight":    StanceProneRight,
		"ProneLeft":     StanceProneLeft,
		"Prone":         StanceProneRight,
		"RopeRight":     StanceRopeRight,
		"RopeLeft":      StanceRopeLeft,
		"Rope":          StanceRopeRight,
		"LadderRight":   StanceLadderRight,
		"LadderLeft":    StanceLadderLeft,
		"Ladder":        StanceLadderRight,
		"Sit":           StanceSitValue,
	}
}

var StanceGroupByKind = map[StanceKind][]uint8{
	StanceDefault:  {StanceDefaultValue},
	StanceWalk:     {StanceWalkRight, StanceWalkLeft},
	StanceStanding: {StanceStandingRight, StanceStandingLeft},
	StanceJump:     {StanceJumpRight, StanceJumpLeft},
	StanceAttack:   {StanceAttackRight, StanceAttackLeft},
	StanceProne:    {StanceProneRight, StanceProneLeft},
	StanceRope:     {StanceRopeRight, StanceRopeLeft},
	StanceLadder:   {StanceLadderRight, StanceLadderLeft},
	StanceSit:      {StanceSitValue},
}
