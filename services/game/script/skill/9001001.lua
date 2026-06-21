-- Skill name (String.wz/Skill.img.xml): 헤이스트

local combat = require("script/lib/combat")

function on_activated_9001001(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end

	local vals = {
		[BuffFlag.Speed] = effect.speed,
		[BuffFlag.Jump] = effect.jump,
	}
	combat.for_each_character_in_skill_area(me, skill, function(ch)
		if ch == nil or not ch:is_alive() then
			return
		end
		ch:buff(skill, vals)
		if ch ~= me then
			ch:show_skill_effect(skill, SkillEffectType.Affected)
		end
	end)
end
