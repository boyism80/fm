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

func (accepter *Character) ShouldSkipInvitePendingForPartySearch(partyID uint32) bool {
	if accepter == nil || partyID == 0 || accepter.GameWorld == nil {
		return false
	}
	var p *Party
	if gw := accepter.GameWorld; gw != nil {
		p = gw.GetPartySystem().Get(partyID)
	}
	if p == nil {
		return false
	}
	mapInst := accepter.GetMap()
	if mapInst == nil {
		return false
	}
	leaderID := p.GetLeaderCharacterId()
	if leaderID == 0 {
		return false
	}
	leader := mapInst.GetPlayer(leaderID)
	if leader == nil {
		return false
	}
	lpid := leader.GetPartyID()
	if lpid == nil || *lpid != partyID {
		return false
	}
	cfg := leader.GetPartySearchConfig()
	if cfg == nil {
		return false
	}
	if accepter.GetPartyID() != nil {
		return false
	}
	lvl := int32(accepter.GetLevel())
	if lvl < cfg.MinLevel || lvl > cfg.MaxLevel {
		return false
	}
	if !leader.HasRoleAtLeast(constant.RoleAdmin) && accepter.HasRoleAtLeast(constant.RoleAdmin) {
		return false
	}
	if !MatchesPartySearchClassMask(accepter, cfg.ClassMask) {
		return false
	}
	return true
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
