-- Skill name (String.wz/Skill.img.xml): 헤이스트

function on_activated_4101004(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end

	local vals = {
		[BuffFlag.Speed] = effect.speed,
		[BuffFlag.Jump] = effect.jump,
	}
	for_each_near_party_member(me, skill, function(ch)
		ch:buff(skill, vals)
		if ch ~= me then
			ch:show_skill_effect(skill, SkillEffectType.Affected)
		end
	end)
end
