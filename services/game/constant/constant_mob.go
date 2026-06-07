package constant

type MobSpawnType int8

const (
	MobSpawnTypeNone    MobSpawnType = -1
	MobSpawnTypeAnimate MobSpawnType = -2
	MobSpawnTypeRevive  MobSpawnType = -3
	MobSpawnTypeFake    MobSpawnType = -4
)

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
