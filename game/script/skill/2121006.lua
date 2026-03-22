-- Skill name (String.wz/Skill.img.xml): 페럴라이즈

function on_attack_2121006(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobStatus.Freeze)
end

function on_activated_2121006(me, skill, params)
end
