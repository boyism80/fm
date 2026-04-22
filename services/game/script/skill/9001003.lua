-- Skill name (String.wz/Skill.img.xml): 블레스

function on_activated_9001003(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end

	local vals = {
		[BuffFlag.WeaponDef] = effect.pdd,
		[BuffFlag.MagicDef] = effect.mdd,
		[BuffFlag.Acc] = effect.acc,
		[BuffFlag.Avoid] = effect.eva,
	}
	for_each_character_in_skill_area(me, skill, function(ch)
		if ch == nil or not ch:is_alive() then
			return
		end
		ch:buff(skill, vals)
		if ch ~= me then
			ch:show_skill_effect(skill, SkillEffectType.Affected)
		end
	end)
end
