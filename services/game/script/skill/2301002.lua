-- Skill name (String.wz/Skill.img.xml): 힐

function on_activated_2301002(me, skill, params)
	local amount = get_heal_recovery_amount(me, skill)
	if amount == 0 then
		return
	end
	for_each_near_party_member(me, skill, function(ch)
		ch:add_hp(amount)
		if ch ~= me then
			ch:show_skill_effect(skill, SkillEffectType.Affected)
		end
	end)
end
