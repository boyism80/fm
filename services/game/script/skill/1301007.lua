-- Skill name (String.wz/Skill.img.xml): 하이퍼 바디

local combat = require("script/lib/combat")

return {
	on_activated = function(me, skill, params)
		local effect = skill:effect()
		if effect == nil then
			return
		end
		combat.for_each_near_party_member(me, skill, function(ch)
			ch:buff(skill, {
				[BuffFlag.MaxHp] = effect.x,
				[BuffFlag.MaxMp] = effect.x,
			})
			if ch ~= me then
				ch:show_skill_effect(skill, SkillEffectType.Affected)
			end
		end)
	end
}
