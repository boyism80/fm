package constant

const (
	MobSkillCastActionMin = 21
	MobSkillCastActionMax = 25
)

func MobSkillCastActionInRange(action int) bool {
	return action >= MobSkillCastActionMin && action <= MobSkillCastActionMax
}
