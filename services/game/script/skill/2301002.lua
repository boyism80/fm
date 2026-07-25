-- Skill name (String.wz/Skill.img.xml): 힐

local util = require("script/lib/skill")
local combat = require("script/lib/combat")

return {
	on_activated = function(me, skill, params)
		local amount = util.get_heal_recovery_amount(me, skill)
		if amount == 0 then
			return
		end
		combat.for_each_near_party_member(me, skill, function(ch)
			ch:add_hp(amount)
			if ch ~= me then
				ch:show_skill_effect(skill, SkillEffectType.Affected)
			end
		end)
	end
}
