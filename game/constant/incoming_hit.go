package constant

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

func AllIncomingHitConstants() map[string]int32 {
	return map[string]int32{
		"Mist":      int32(IncomingHitMist),
		"Env":       int32(IncomingHitEnv),
		"MapDebuff": int32(IncomingHitMapDebuff),
		"Collide":   int32(IncomingHitCollide),
	}
}
