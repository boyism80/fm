-- Skill name (String.wz/Skill.img.xml): 메이플 용사

function on_activated_1321000(me, skill, params)
	for_each_near_party_member(me, skill, function(ch)
		apply_buff_from_effect(ch, skill, BuffFlag.MapleWarrior, "x")
		if ch ~= me then
			ch:show_skill_effect(skill, SkillEffectType.Affected)
		end
	end)
end
