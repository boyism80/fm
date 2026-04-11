package constant

// StanceKind is the canonical stance group value (movement state).
// Used for stance_of checks; each kind may map to one or two raw stance bytes (left/right).
type StanceKind uint8

const (
	StanceDefault  StanceKind = 0
	StanceWalk     StanceKind = 2  // 2=right, 3=left
	StanceStanding StanceKind = 4  // 4=right, 5=left
	StanceJump     StanceKind = 6  // 6=right, 7=left
	StanceAttack   StanceKind = 8  // 8=right, 9=left
	StanceProne    StanceKind = 10 // 10=right, 11=left
	StanceRope     StanceKind = 12 // 12=right, 13=left
	StanceLadder   StanceKind = 14 // 14=right, 15=left
	StanceSit      StanceKind = 21
)

// Individual stance values (including Left/Right); same numeric value as protocol.
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

// AllStanceConstants returns name -> value for every stance (for Lua injection).
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

// StanceGroupByKind maps a stance kind to the raw stance bytes that belong to it.
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
