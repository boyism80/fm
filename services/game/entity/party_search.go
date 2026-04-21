package entity

import "github.com/boyism80/fm/services/game/constant"

type PartySearchConfig struct {
	MinLevel      int32
	MaxLevel      int32
	MembersNeeded int32
	ClassMask     int32
}

func (ch *Character) SetPartySearchConfig(cfg *PartySearchConfig) {
	if ch == nil {
		return
	}
	ch.partySearchConfig = cfg
}

func (ch *Character) GetPartySearchConfig() *PartySearchConfig {
	if ch == nil {
		return nil
	}
	return ch.partySearchConfig
}

func MatchesPartySearchClassMask(ch *Character, classMask int32) bool {
	if ch == nil {
		return false
	}
	if (classMask & 0x1) != 0 {
		return true
	}
	classCodesByMask := map[int32]constant.ClassType{
		0x4:     constant.ClassWarrior,
		0x8:     constant.ClassFighter,
		0x10:    constant.ClassPage,
		0x20:    constant.ClassSpearman,
		0x40:    constant.ClassMagician,
		0x80:    constant.ClassFpWizard,
		0x100:   constant.ClassIlWizard,
		0x200:   constant.ClassCleric,
		0x400:   constant.ClassPirate,
		0x800:   constant.ClassBrawler,
		0x1000:  constant.ClassGunslinger,
		0x2000:  constant.ClassThief,
		0x4000:  constant.ClassAssassin,
		0x8000:  constant.ClassBandit,
		0x10000: constant.ClassBowman,
		0x20000: constant.ClassHunter,
		0x40000: constant.ClassCrossbowman,
	}
	if (classMask & 0x2) != 0 {
		if ch.ClassOf(constant.ClassBeginner) {
			return true
		}
	}
	for bit, classCode := range classCodesByMask {
		if (classMask & bit) == 0 {
			continue
		}
		if ch.ClassOf(classCode) {
			return true
		}
	}
	return false
}
