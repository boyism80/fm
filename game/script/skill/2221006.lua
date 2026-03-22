-- Skill name (String.wz/Skill.img.xml): 체인 라이트닝

function on_attack_2221006(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobStatus.Stun)
end

function on_activated_2221006(me, skill, params)
end
