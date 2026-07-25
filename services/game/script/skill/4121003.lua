-- Skill name (String.wz/Skill.img.xml): 쇼다운

local combat = require("script/lib/combat")

return {
	on_mob_buff = function(mob, skill, causer)
		local effect = skill:effect()
		if effect == nil then
			return
		end
		if effect.x <= 0 then
			return
		end
		mob:exp_rate(mob:exp_rate() + effect.x)
		mob:drop_rate(mob:drop_rate() + effect.x)
	end,

	on_mob_unbuff = function(mob, skill, causer)
		local effect = skill:effect()
		if effect == nil then
			return
		end
		if effect.x <= 0 then
			return
		end
		mob:exp_rate(mob:exp_rate() - effect.x)
		mob:drop_rate(mob:drop_rate() - effect.x)
	end,

	on_attack = function(me, skill, damages)
		combat.apply_showdown(me, skill, damages)
	end
}
