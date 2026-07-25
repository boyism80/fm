-- Skill name (String.wz/Skill.img.xml): 아이언 월

local combat = require("script/lib/combat")

return {
	on_activated = function(me, skill, params)
		local effect = skill:effect()
		if effect == nil then
			return
		end
		if effect.pdd <= 0 and effect.mdd <= 0 then
			return
		end
		local vals = {
			[BuffFlag.WeaponDef] = effect.pdd,
			[BuffFlag.MagicDef] = effect.mdd,
		}
		combat.for_each_near_party_member(me, skill, function(ch)
			ch:buff(skill, vals)
			if ch ~= me then
				ch:show_skill_effect(skill, SkillEffectType.Affected)
			end
		end)
	end
}
