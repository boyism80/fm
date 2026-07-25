-- Skill name (String.wz/Skill.img.xml): 디스펠

local util = require("script/lib/skill")
local combat = require("script/lib/combat")

return {
	on_activated = function(me, skill, params)
		combat.for_each_near_party_member(me, skill, function(ch)
			util.apply_dispel(ch, skill)
			if ch ~= me then
				ch:show_skill_effect(skill, SkillEffectType.Affected)
			end
		end)
	end
}
