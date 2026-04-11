package constant

// MobSpawnType represents how a monster spawns.
type MobSpawnType int8

const (
	MOB_SPAWN_TYPE_NONE    MobSpawnType = -1 // No spawn animation
	MOB_SPAWN_TYPE_ANIMATE MobSpawnType = -2 // With spawn animation
)

// MobDieAnimationType represents how a monster dies.
type MobDieAnimationType uint8

const (
	MOB_DIE_ANIMATION_TYPE_DISAPPEAR MobDieAnimationType = 0 // Instant disappear
	MOB_DIE_ANIMATION_TYPE_FADE_OUT  MobDieAnimationType = 1 // Fade out animation
)

func AllMobDieAnimationTypes() map[string]MobDieAnimationType {
	return map[string]MobDieAnimationType{
		"Disappear": MOB_DIE_ANIMATION_TYPE_DISAPPEAR,
		"FadeOut":   MOB_DIE_ANIMATION_TYPE_FADE_OUT,
	}
}
