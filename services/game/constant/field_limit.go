package constant

type FieldLimit uint32

const (
	FieldLimitJump              FieldLimit = 0x1
	FieldLimitMovementSkill     FieldLimit = 0x2
	FieldLimitSummoningBag      FieldLimit = 0x4
	FieldLimitMysticDoor        FieldLimit = 0x8
	FieldLimitMigrate           FieldLimit = 0x10
	FieldLimitReturnScroll      FieldLimit = 0x20
	FieldLimitTeleportStone     FieldLimit = 0x40
	FieldLimitMinigame          FieldLimit = 0x80
	FieldLimitMapReturnScroll   FieldLimit = 0x100
	FieldLimitMount             FieldLimit = 0x200
	FieldLimitPotion            FieldLimit = 0x400
	FieldLimitPartyLeaderChange FieldLimit = 0x800
	FieldLimitWeddingInvitation FieldLimit = 0x2000
	FieldLimitWeather           FieldLimit = 0x4000
	FieldLimitLieDetector       FieldLimit = 0x10000
	FieldLimitFallDown          FieldLimit = 0x20000
	FieldLimitSummonNpc         FieldLimit = 0x40000
	FieldLimitNoFallDamage      FieldLimit = 0x100000
	FieldLimitDrop              FieldLimit = 0x400000
)
