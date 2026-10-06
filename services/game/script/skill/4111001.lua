-- Skill name (String.wz/Skill.img.xml): 메소 업

local combat = require("script/lib/combat")

return {
	on_activated = function(me, skill, params)
		local effect = skill:effect()
		if effect == nil then
			return
		end
		local percent = effect.x or 100
		combat.for_each_near_party_member(me, skill, function(ch)
			ch:buff(skill, BuffFlag.MesoUp, percent)
			if ch ~= me then
				ch:show_skill_effect(skill, SkillEffectType.Affected)
			end
		end)
	end
}
