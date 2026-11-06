package constant

type SkillEffectType uint8

const (
	SkillEffectTypeCast     SkillEffectType = 1
	SkillEffectTypeAffected SkillEffectType = 2
	SkillEffectTypePortal   SkillEffectType = 7
	// SkillEffectTypeNineSpirit plays the effect of the item in the skill ID field, as for the Horntail kill buff.
	SkillEffectTypeNineSpirit SkillEffectType = 0x0B
)
