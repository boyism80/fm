package game

import "github.com/boyism80/fm/maple/wz"

type CarnivalMembership struct {
	owner *Character
	team  *CarnivalTeam
}

func (c *CarnivalMembership) Match() *CarnivalMatch {
	team := c.Team()
	if team == nil {
		return nil
	}
	return team.Match
}

// The team outlives the match slot so reward NPCs can read the result after Conclude.
func (c *CarnivalMembership) Team() *CarnivalTeam {
	if c.team == nil || !c.team.HasMember(c.owner.GetID()) {
		return nil
	}
	return c.team
}

func (c *CarnivalMembership) useItem(wzConsume *wz.Consume) {
	team := c.Team()
	if team == nil || team.Match == nil {
		return
	}
	if wzConsume.CP > 0 {
		team.AddCP(c.owner, wzConsume.CP)
	}
	if wzConsume.NuffSkillID == 0 {
		return
	}
	enemy := team.Match.enemyTeam(team.TeamID)
	if enemy == nil {
		return
	}
	m := c.owner.GetMap()
	if m == nil {
		return
	}
	enemy.Debuff(m, wzConsume.NuffSkillID)
}
