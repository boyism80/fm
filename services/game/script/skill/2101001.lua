-- Skill name (String.wz/Skill.img.xml): 메디테이션

local util = require("script/lib/skill")
local combat = require("script/lib/combat")

return {
	on_activated = function(me, skill, params)
		combat.for_each_near_party_member(me, skill, function(ch)
			util.apply_buff_from_effect(ch, skill, BuffFlag.MagicAtk, "mad")
			if ch ~= me then
				ch:show_skill_effect(skill, SkillEffectType.Affected)
			end
		end)
	end
}
