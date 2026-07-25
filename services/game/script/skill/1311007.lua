-- Skill name (String.wz/Skill.img.xml): 파워 크래쉬

local combat = require("script/lib/combat")

return {
	on_activated = function(me, skill, params)
		combat.for_each_mob_in_skill_area(me, skill, function(mob)
			mob:clear_buffs(MobBuff.WeaponAttackUp)
		end)
	end
}
