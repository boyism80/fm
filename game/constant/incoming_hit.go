package constant

// IncomingHitType is the first byte of the TakeDamage client packet (after update tick).
//
// Negative values are reserved hit sources that do not select a mob attack slot:
//   -4 mist, -3 environmental / map hazard, -2 map debuff, -1 touch / collision with mob.
//
// Non-negative values (0, 1, 2, …) are not extra enum members: they are the 0-based index
// into the mob’s WZ attack list (attack1 → 0, attack2 → 1, attack3 → 2, …). The client
// sends the slot index of the attack that hit the player.
//
// Physical vs magic is not encoded by this byte. Use Mob WZ attackN/info magic, or the
// following HitElement byte, together with that attack index — not the index value alone.
type IncomingHitType int8

const (
	IncomingHitMist      IncomingHitType = -4
	IncomingHitEnv       IncomingHitType = -3
	IncomingHitMapDebuff IncomingHitType = -2
	IncomingHitCollide   IncomingHitType = -1
)

func (t IncomingHitType) MobAttackIndex() (index int, ok bool) {
	if t >= 0 {
		return int(t), true
	}
	return 0, false
}

// AllIncomingHitConstants maps PascalCase keys to TakeDamage packet hit-type values for Lua (IncomingHit global).
func AllIncomingHitConstants() map[string]int32 {
	return map[string]int32{
		"Mist":      int32(IncomingHitMist),
		"Env":       int32(IncomingHitEnv),
		"MapDebuff": int32(IncomingHitMapDebuff),
		"Collide":   int32(IncomingHitCollide),
	}
}
