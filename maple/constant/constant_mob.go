package constant

type MobSpawnType int8

const (
	MobSpawnTypeNone    MobSpawnType = -1
	MobSpawnTypeAnimate MobSpawnType = -2
	MobSpawnTypeRevive  MobSpawnType = -3
	MobSpawnTypeFake    MobSpawnType = -4
)

func AllMobSpawnTypes() map[string]MobSpawnType {
	return map[string]MobSpawnType{
		"None":    MobSpawnTypeNone,
		"Animate": MobSpawnTypeAnimate,
		"Revive":  MobSpawnTypeRevive,
		"Fake":    MobSpawnTypeFake,
	}
}

type MobDieAnimationType uint8

const (
	MobDieAnimationTypeDisappear MobDieAnimationType = 0
	MobDieAnimationTypeFadeOut   MobDieAnimationType = 1
)

func AllMobDieAnimationTypes() map[string]MobDieAnimationType {
	return map[string]MobDieAnimationType{
		"Disappear": MobDieAnimationTypeDisappear,
		"FadeOut":   MobDieAnimationTypeFadeOut,
	}
}

type MobDamageDisplayType uint8

const (
	MobDamageDisplayNormal     MobDamageDisplayType = 0
	MobDamageDisplayAllyShowHp MobDamageDisplayType = 1
	MobDamageDisplayAllySilent MobDamageDisplayType = 2
)

func (t MobDamageDisplayType) IncludesHpMaxHp() bool {
	return t == MobDamageDisplayAllyShowHp || t == MobDamageDisplayAllySilent
}
