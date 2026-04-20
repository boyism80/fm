package constant

type MobSpawnType int8

const (
	MOB_SPAWN_TYPE_NONE    MobSpawnType = -1
	MOB_SPAWN_TYPE_ANIMATE MobSpawnType = -2
)

type MobDieAnimationType uint8

const (
	MOB_DIE_ANIMATION_TYPE_DISAPPEAR MobDieAnimationType = 0
	MOB_DIE_ANIMATION_TYPE_FADE_OUT  MobDieAnimationType = 1
)

func AllMobDieAnimationTypes() map[string]MobDieAnimationType {
	return map[string]MobDieAnimationType{
		"Disappear": MOB_DIE_ANIMATION_TYPE_DISAPPEAR,
		"FadeOut":   MOB_DIE_ANIMATION_TYPE_FADE_OUT,
	}
}
