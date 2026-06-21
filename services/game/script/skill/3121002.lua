-- Skill name (String.wz/Skill.img.xml): 샤프 아이즈

local combat = require("script/lib/combat")

function on_activated_3121002(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.x <= 0 or effect.y <= 0 then
		return
	end
	local packed = effect.x * 256 + (effect.y % 256)
	combat.for_each_near_party_member(me, skill, function(ch)
		ch:buff(skill, BuffFlag.SharpEyes, packed)
		if ch ~= me then
			ch:show_skill_effect(skill, SkillEffectType.Affected)
		end
	end)
end
