package constant

import "time"

type PetSkill uint16

const (
	PetSkillPickupItem   PetSkill = 0x01
	PetSkillLongRange    PetSkill = 0x02
	PetSkillDropSweep    PetSkill = 0x04
	PetSkillIgnorePickup PetSkill = 0x08
	PetSkillPickupAll    PetSkill = 0x10
	PetSkillConsumeHP    PetSkill = 0x20
	PetSkillConsumeMP    PetSkill = 0x40
)

var PetSkillRequires = map[PetSkill]PetSkill{
	PetSkillLongRange: PetSkillDropSweep,
	PetSkillDropSweep: PetSkillPickupItem,
	PetSkillPickupAll: PetSkillPickupItem,
}

type PetRemoveReason uint8

const (
	PetRemoveReasonNone    PetRemoveReason = 0
	PetRemoveReasonHungry  PetRemoveReason = 1
	PetRemoveReasonExpired PetRemoveReason = 2
)

const (
	PetMaxLevel          = 30
	PetMaxCloseness      = 30000
	PetMaxFullness       = 100
	PetStarveFullness    = 5
	PetStarvedFullness   = 15
	PetFoodClosenessRate = 51
	PetCalledByNameBonus = 10
	PetHungerInterval    = 164 * time.Second
	PetLootClientRange   = 100
	PetLootRange         = 800
	PetNameMinBytes      = 3
	PetNameMaxBytes      = 12
)

var PetClosenessByLevel = [PetMaxLevel]uint16{0, 1, 3, 6, 14, 31, 60, 108, 181, 287, 434, 632, 891, 1224, 1642, 2161, 2793, 3557, 4467, 5542, 6801, 8263, 9950, 11882, 14084, 16578, 19391, 22547, 26074, 30000}
