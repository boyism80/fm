-- Skill name (String.wz/Skill.img.xml): 타임 리프

local combat = require("script/lib/combat")

function on_activated_5121010(me, skill, params)
	combat.for_each_near_party_member(me, skill, function(ch)
		local skills = ch:skills()
		if skills == nil then
			return
		end
		for id, sk in pairs(skills) do
			if id ~= Skill.TimeLeap and sk ~= nil then
				sk:cooldown(0)
			end
		end
		if ch ~= me then
			ch:show_skill_effect(skill, SkillEffectType.Affected)
		end
	end)
end
